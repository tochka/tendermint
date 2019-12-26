package counter

import (
	"encoding/binary"
	"fmt"

	"github.com/tendermint/tendermint/abci/example/code"
	"github.com/tendermint/tendermint/abci/types"
)

type CounterApplication struct {
	types.BaseApplication

	hashCount     map[string]int
	txCount       map[string]int
	serial        map[string]bool
	defaultSerial bool
}

func NewCounterApplication(serial bool) *CounterApplication {
	return &CounterApplication{
		defaultSerial: serial,
		serial:        make(map[string]bool),
		txCount:       make(map[string]int),
		hashCount:     make(map[string]int),
	}
}

func (app *CounterApplication) InitChain(req types.RequestInitChain) types.ResponseInitChain {
	if req.ChainId == "" {
		panic("chain id is empty")
	}
	app.serial[req.ChainId] = app.defaultSerial
	return types.ResponseInitChain{}
}

func (app *CounterApplication) Info(req types.RequestInfo) types.ResponseInfo {
	if req.ChainId == "" {
		panic("chain id is empty")
	}
	return types.ResponseInfo{Data: fmt.Sprintf("{\"hashes\":%v,\"txs\":%v}", app.hashCount[req.ChainId], app.txCount[req.ChainId])}
}

func (app *CounterApplication) SetOption(req types.RequestSetOption) types.ResponseSetOption {
	if req.ChainId == "" {
		panic("chain id is empty")
	}
	key, value := req.Key, req.Value
	if key == "serial" && value == "on" {
		app.serial[req.ChainId] = true
	} else {
		/*
			TODO Panic and have the ABCI server pass an exception.
			The client can call SetOptionSync() and get an `error`.
			return types.ResponseSetOption{
				Error: fmt.Sprintf("Unknown key (%s) or value (%s)", key, value),
			}
		*/
		return types.ResponseSetOption{}
	}

	return types.ResponseSetOption{}
}

func (app *CounterApplication) DeliverTx(req types.RequestDeliverTx) types.ResponseDeliverTx {
	if req.ChainId == "" {
		panic("chain id is empty")
	}

	if app.Serial(req.ChainId) {
		if len(req.Tx) > 8 {
			return types.ResponseDeliverTx{
				Code: code.CodeTypeEncodingError,
				Log:  fmt.Sprintf("Max tx size is 8 bytes, got %d", len(req.Tx))}
		}
		tx8 := make([]byte, 8)
		copy(tx8[len(tx8)-len(req.Tx):], req.Tx)
		txValue := binary.BigEndian.Uint64(tx8)
		if txValue != uint64(app.txCount[req.ChainId]) {
			return types.ResponseDeliverTx{
				Code: code.CodeTypeBadNonce,
				Log:  fmt.Sprintf("Invalid nonce. Expected %v, got %v", app.txCount[req.ChainId], txValue)}
		}
	}
	app.txCount[req.ChainId]++
	return types.ResponseDeliverTx{Code: code.CodeTypeOK}
}

func (app *CounterApplication) CheckTx(req types.RequestCheckTx) types.ResponseCheckTx {
	if req.ChainId == "" {
		panic("chain id is empty")
	}
	if app.Serial(req.ChainId) {
		if len(req.Tx) > 8 {
			return types.ResponseCheckTx{
				Code: code.CodeTypeEncodingError,
				Log:  fmt.Sprintf("Max tx size is 8 bytes, got %d", len(req.Tx))}
		}
		tx8 := make([]byte, 8)
		copy(tx8[len(tx8)-len(req.Tx):], req.Tx)
		txValue := binary.BigEndian.Uint64(tx8)
		if txValue < uint64(app.txCount[req.ChainId]) {
			return types.ResponseCheckTx{
				Code: code.CodeTypeBadNonce,
				Log:  fmt.Sprintf("Invalid nonce. Expected >= %v, got %v", app.txCount[req.ChainId], txValue)}
		}
	}
	return types.ResponseCheckTx{Code: code.CodeTypeOK}
}

func (app *CounterApplication) Commit(req types.RequestCommit) (resp types.ResponseCommit) {
	if req.ChainId == "" {
		panic("chain id is empty")
	}
	app.hashCount[req.ChainId]++
	if app.txCount[req.ChainId] == 0 {
		return types.ResponseCommit{}
	}
	hash := make([]byte, 8)
	binary.BigEndian.PutUint64(hash, uint64(app.txCount[req.ChainId]))
	return types.ResponseCommit{Data: hash}
}

func (app *CounterApplication) Query(reqQuery types.RequestQuery) types.ResponseQuery {
	if reqQuery.ChainId == "" {
		panic("chain id is empty")
	}
	switch reqQuery.Path {
	case "hash":
		return types.ResponseQuery{Value: []byte(fmt.Sprintf("%v", app.hashCount[reqQuery.ChainId]))}
	case "tx":
		return types.ResponseQuery{Value: []byte(fmt.Sprintf("%v", app.txCount[reqQuery.ChainId]))}
	default:
		return types.ResponseQuery{Log: fmt.Sprintf("Invalid query path. Expected hash or tx, got %v", reqQuery.Path)}
	}
}

func (app *CounterApplication) Serial(chainID string) bool {
	if chainID == "" {
		panic("chain id is empty")
	}
	v, ok := app.serial[chainID]
	if !ok {
		return app.defaultSerial
	}
	return v
}
