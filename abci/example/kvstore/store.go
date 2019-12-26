package kvstore

import (
	db "github.com/tendermint/tm-db"
	dbm "github.com/tendermint/tm-db"
)

var (
	kvPairChainDBPrefix = []byte{0x00}
)

func chainDBPrefixKey(chainID string, key []byte) []byte {
	r := make([]byte, 0, 2*len(kvPairChainDBPrefix)+len(chainID)+len(key))
	r = append(r, kvPairChainDBPrefix...)
	r = append(r, []byte(chainID)...)
	r = append(r, []byte(kvPairChainDBPrefix)...)
	r = append(r, key...)
	return r
}

// DBs are goroutine safe.
type ChainDB struct {
	chainID string
	db      dbm.DB
}

func (db ChainDB) Get(key []byte) []byte {
	return db.db.Get(chainDBPrefixKey(db.chainID, key))
}

func (db ChainDB) Has(key []byte) bool {
	return db.db.Has(chainDBPrefixKey(db.chainID, key))
}

func (db ChainDB) Set(k []byte, v []byte) {
	db.db.Set(chainDBPrefixKey(db.chainID, k), v)
}
func (db ChainDB) SetSync(k []byte, v []byte) {
	db.db.SetSync(chainDBPrefixKey(db.chainID, k), v)
}

func (db ChainDB) Delete(k []byte) {
	db.db.Delete(chainDBPrefixKey(db.chainID, k))
}
func (db ChainDB) DeleteSync(k []byte) {
	db.db.DeleteSync(chainDBPrefixKey(db.chainID, k))
}

// Iterate over a domain of keys in ascending order. End is exclusive.
// Start must be less than end, or the Iterator is invalid.
// A nil start is interpreted as an empty byteslice.
// If end is nil, iterates up to the last item (inclusive).
// CONTRACT: No writes may happen within a domain while an iterator exists over it.
// CONTRACT: start, end readonly []byte
func (db ChainDB) Iterator(start, end []byte) db.Iterator {
	chs := chainDBPrefixKey(db.chainID, start)
	che := chainDBPrefixKey(db.chainID, end)
	if len(end) == 0 {
		che[len(che)-1]++ // the last byte is split symbole so we iterate over the end
	}
	return chainIterator{itr: db.db.Iterator(chs, che), chainID: db.chainID}
}

// Iterate over a domain of keys in descending order. End is exclusive.
// Start must be less than end, or the Iterator is invalid.
// If start is nil, iterates up to the first/least item (inclusive).
// If end is nil, iterates from the last/greatest item (inclusive).
// CONTRACT: No writes may happen within a domain while an iterator exists over it.
// CONTRACT: start, end readonly []byte
func (db ChainDB) ReverseIterator(start, end []byte) db.Iterator {
	chs := chainDBPrefixKey(db.chainID, start)
	che := chainDBPrefixKey(db.chainID, end)
	if len(end) == 0 {
		che[len(che)-1]++ // the last byte is split symbole so we iterate over the end
	}
	return chainIterator{itr: db.db.ReverseIterator(chs, che), chainID: db.chainID}
}

// Closes the connection.
func (db ChainDB) Close() {
	db.db.Close()
}

// Creates a batch for atomic updates.
func (db ChainDB) NewBatch() db.Batch {
	return &chainDBBatch{batch: db.db.NewBatch(), chainID: db.chainID}
}

// For debugging
func (db ChainDB) Print() {
	db.db.Print()
}

// Stats returns a map of property values for all keys and the size of the cache.
func (db ChainDB) Stats() map[string]string {
	return db.db.Stats()
}

//----------------------------------------
// Batch
type chainDBBatch struct {
	batch   db.Batch
	chainID string
}

// Implements Batch.
func (mBatch *chainDBBatch) Set(key, value []byte) {
	mBatch.batch.Set(chainDBPrefixKey(mBatch.chainID, key), value)
}

// Implements Batch.
func (mBatch *chainDBBatch) Delete(key []byte) {
	mBatch.batch.Delete(chainDBPrefixKey(mBatch.chainID, key))
}

// Implements Batch.
func (mBatch *chainDBBatch) Write() {
	mBatch.batch.Write()
}

// Implements Batch.
func (mBatch *chainDBBatch) WriteSync() {
	mBatch.batch.WriteSync()
}

// Implements Batch.
// Close is no-op for chainDBBatch.
func (mBatch *chainDBBatch) Close() {
	mBatch.batch.Close()
}

//----------------------------------------
// Iterator

type chainIterator struct {
	itr     dbm.Iterator
	chainID string
}

func (itr chainIterator) Domain() (start []byte, end []byte) {
	chs, che := itr.itr.Domain()
	return itr.trimChainPrefix(chs), itr.trimChainPrefix(che)
}

func (itr chainIterator) Valid() bool {
	return itr.itr.Valid()
}

func (itr chainIterator) Next() {
	itr.itr.Next()
}
func (itr chainIterator) Close() {
	itr.itr.Close()
}

func (itr chainIterator) Key() (key []byte) {
	return itr.trimChainPrefix(itr.itr.Key())
}
func (itr chainIterator) Value() (value []byte) {
	return itr.itr.Value()
}

func (itr chainIterator) trimChainPrefix(key []byte) []byte {
	prfx := chainDBPrefixKey(itr.chainID, nil)
	return key[len(prfx):]
}
