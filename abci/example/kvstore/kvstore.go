package kvstore

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/tendermint/tendermint/abci/example/code"
	"github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/version"
	dbm "github.com/tendermint/tm-db"
)

var (
	stateKey        = []byte("stateKey:")
	kvPairPrefixKey = []byte("kvPairKey:")

	ProtocolVersion version.Protocol = 0x1
)

type State struct {
	db      dbm.DB `json:"-"`
	Size    int64  `json:"size"`
	Height  int64  `json:"height"`
	AppHash []byte `json:"app_hash"`
}

func loadState(db dbm.DB, chainID string) State {
	chdb := ChainDB{db: db, chainID: chainID}
	stateBytes := chdb.Get(stateKey)
	var state State
	if len(stateBytes) != 0 {
		err := json.Unmarshal(stateBytes, &state)
		if err != nil {
			panic(err)
		}
	}
	state.db = chdb
	return state
}

func saveState(state State) {
	stateBytes, err := json.Marshal(state)
	if err != nil {
		panic(err)
	}
	state.db.Set(stateKey, stateBytes)
}

func prefixKey(key []byte) []byte {
	return append(kvPairPrefixKey, key...)
}

//---------------------------------------------------

var _ types.Application = (*KVStoreApplication)(nil)

type KVStoreApplication struct {
	types.BaseApplication

	db    dbm.DB
	state map[string]*State
}

func NewKVStoreApplication() *KVStoreApplication {
	return &KVStoreApplication{db: dbm.NewMemDB(), state: make(map[string]*State)}
}

func (app *KVStoreApplication) Info(req types.RequestInfo) (resInfo types.ResponseInfo) {
	if req.ChainId == "" {
		panic("chain id is empty")
	}
	return types.ResponseInfo{
		Data:       fmt.Sprintf("{\"size\":%v}", app.State(req.ChainId).Size),
		Version:    version.ABCIVersion,
		AppVersion: ProtocolVersion.Uint64(),
	}
}

// tx is either "key=value" or just arbitrary bytes
func (app *KVStoreApplication) DeliverTx(req types.RequestDeliverTx) types.ResponseDeliverTx {
	if req.ChainId == "" {
		panic("chain id is empty")
	}
	var key, value []byte
	parts := bytes.Split(req.Tx, []byte("="))
	if len(parts) == 2 {
		key, value = parts[0], parts[1]
	} else {
		key, value = req.Tx, req.Tx
	}

	s := app.State(req.ChainId)
	s.db.Set(prefixKey(key), value)
	s.Size += 1

	events := []types.Event{
		{
			Type: "app",
			Attributes: []cmn.KVPair{
				{Key: []byte("creator"), Value: []byte("Cosmoshi Netowoko")},
				{Key: []byte("key"), Value: key},
			},
		},
	}

	return types.ResponseDeliverTx{Code: code.CodeTypeOK, Events: events}
}

func (app *KVStoreApplication) CheckTx(req types.RequestCheckTx) types.ResponseCheckTx {
	if req.ChainId == "" {
		panic("chain id is empty")
	}
	return types.ResponseCheckTx{Code: code.CodeTypeOK, GasWanted: 1}
}

func (app *KVStoreApplication) Commit(req types.RequestCommit) types.ResponseCommit {
	if req.ChainId == "" {
		panic("chain id is empty")
	}

	// Using a memdb - just return the big endian size of the db
	appHash := make([]byte, 8)
	s := app.State(req.ChainId)
	binary.PutVarint(appHash, s.Size)
	s.AppHash = appHash
	s.Height += 1
	saveState(*s)
	return types.ResponseCommit{Data: appHash}
}

// Returns an associated value or nil if missing.
func (app *KVStoreApplication) Query(reqQuery types.RequestQuery) (resQuery types.ResponseQuery) {
	if reqQuery.ChainId == "" {
		panic("chain id is empty")
	}
	s := app.State(reqQuery.ChainId)
	if reqQuery.Prove {
		value := s.db.Get(prefixKey(reqQuery.Data))
		resQuery.Index = -1 // TODO make Proof return index
		resQuery.Key = reqQuery.Data
		resQuery.Value = value
		if value != nil {
			resQuery.Log = "exists"
		} else {
			resQuery.Log = "does not exist"
		}
		return
	} else {
		resQuery.Key = reqQuery.Data
		value := s.db.Get(prefixKey(reqQuery.Data))
		resQuery.Value = value
		if value != nil {
			resQuery.Log = "exists"
		} else {
			resQuery.Log = "does not exist"
		}
		return
	}
}

func (app *KVStoreApplication) State(chainID string) *State {
	if chainID == "" {
		panic("chain id is empty")
	}
	var state *State
	var ok bool
	if state, ok = app.state[chainID]; !ok {
		s := loadState(app.db, chainID)
		state = &s
		app.state[chainID] = state
	}
	return state
}
