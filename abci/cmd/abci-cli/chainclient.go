package main

import (
	abcicli "github.com/tendermint/tendermint/abci/client"
	"github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
)

type chainClient struct {
	chainID string
	client  abcicli.Client
}

func (cc chainClient) SetResponseCallback(c abcicli.Callback) {
	cc.client.SetResponseCallback(c)
}

func (cc chainClient) Error() error {
	return cc.client.Error()
}

func (cc chainClient) FlushAsync() *abcicli.ReqRes {
	return cc.client.FlushAsync()
}

func (cc chainClient) EchoAsync(msg string) *abcicli.ReqRes {
	return cc.client.EchoAsync(msg)
}

func (cc chainClient) InfoAsync(req types.RequestInfo) *abcicli.ReqRes {
	if req.ChainId == "" {
		req.ChainId = cc.chainID
	}
	return cc.client.InfoAsync(req)
}

func (cc chainClient) SetOptionAsync(req types.RequestSetOption) *abcicli.ReqRes {
	if req.ChainId == "" {
		req.ChainId = cc.chainID
	}
	return cc.client.SetOptionAsync(req)
}

func (cc chainClient) DeliverTxAsync(req types.RequestDeliverTx) *abcicli.ReqRes {
	if req.ChainId == "" {
		req.ChainId = cc.chainID
	}
	return cc.client.DeliverTxAsync(req)
}

func (cc chainClient) CheckTxAsync(req types.RequestCheckTx) *abcicli.ReqRes {
	if req.ChainId == "" {
		req.ChainId = cc.chainID
	}
	return cc.client.CheckTxAsync(req)
}

func (cc chainClient) QueryAsync(req types.RequestQuery) *abcicli.ReqRes {
	if req.ChainId == "" {
		req.ChainId = cc.chainID
	}
	return cc.client.QueryAsync(req)
}

func (cc chainClient) CommitAsync(req types.RequestCommit) *abcicli.ReqRes {
	if req.ChainId == "" {
		req.ChainId = cc.chainID
	}
	return cc.client.CommitAsync(req)
}

func (cc chainClient) InitChainAsync(req types.RequestInitChain) *abcicli.ReqRes {
	if req.ChainId == "" {
		req.ChainId = cc.chainID
	}
	return cc.client.InitChainAsync(req)
}

func (cc chainClient) BeginBlockAsync(req types.RequestBeginBlock) *abcicli.ReqRes {
	if req.ChainId == "" {
		req.ChainId = cc.chainID
	}
	return cc.client.BeginBlockAsync(req)
}

func (cc chainClient) EndBlockAsync(req types.RequestEndBlock) *abcicli.ReqRes {
	if req.ChainId == "" {
		req.ChainId = cc.chainID
	}
	return cc.client.EndBlockAsync(req)
}

func (cc chainClient) FlushSync() error {
	return cc.client.FlushSync()
}

func (cc chainClient) EchoSync(msg string) (*types.ResponseEcho, error) {
	return cc.client.EchoSync(msg)
}

func (cc chainClient) InfoSync(req types.RequestInfo) (*types.ResponseInfo, error) {
	if req.ChainId == "" {
		req.ChainId = cc.chainID
	}
	return cc.client.InfoSync(req)
}

func (cc chainClient) SetOptionSync(req types.RequestSetOption) (*types.ResponseSetOption, error) {
	if req.ChainId == "" {
		req.ChainId = cc.chainID
	}
	return cc.client.SetOptionSync(req)
}

func (cc chainClient) DeliverTxSync(req types.RequestDeliverTx) (*types.ResponseDeliverTx, error) {
	if req.ChainId == "" {
		req.ChainId = cc.chainID
	}
	return cc.client.DeliverTxSync(req)
}

func (cc chainClient) CheckTxSync(req types.RequestCheckTx) (*types.ResponseCheckTx, error) {
	if req.ChainId == "" {
		req.ChainId = cc.chainID
	}
	return cc.client.CheckTxSync(req)
}

func (cc chainClient) QuerySync(req types.RequestQuery) (*types.ResponseQuery, error) {
	if req.ChainId == "" {
		req.ChainId = cc.chainID
	}
	return cc.client.QuerySync(req)
}

func (cc chainClient) CommitSync(req types.RequestCommit) (*types.ResponseCommit, error) {
	if req.ChainId == "" {
		req.ChainId = cc.chainID
	}
	return cc.client.CommitSync(req)
}

func (cc chainClient) InitChainSync(req types.RequestInitChain) (*types.ResponseInitChain, error) {
	if req.ChainId == "" {
		req.ChainId = cc.chainID
	}
	return cc.client.InitChainSync(req)
}

func (cc chainClient) BeginBlockSync(req types.RequestBeginBlock) (*types.ResponseBeginBlock, error) {
	if req.ChainId == "" {
		req.ChainId = cc.chainID
	}
	return cc.client.BeginBlockSync(req)
}

func (cc chainClient) EndBlockSync(req types.RequestEndBlock) (*types.ResponseEndBlock, error) {
	if req.ChainId == "" {
		req.ChainId = cc.chainID
	}
	return cc.client.EndBlockSync(req)
}

func (cc chainClient) Start() error {
	return cc.client.Start()
}

func (cc chainClient) OnStart() error {
	return cc.client.OnStart()
}

func (cc chainClient) Stop() error {
	return cc.client.Stop()
}

func (cc chainClient) OnStop() {
	cc.client.OnStop()
}

func (cc chainClient) Reset() error {
	return cc.client.Reset()
}

func (cc chainClient) OnReset() error {
	return cc.client.OnReset()
}

func (cc chainClient) IsRunning() bool {
	return cc.client.IsRunning()
}

func (cc chainClient) Quit() <-chan struct{} {
	return cc.client.Quit()
}

func (cc chainClient) String() string {
	return cc.client.String()
}

func (cc chainClient) SetLogger(l log.Logger) {
	cc.client.SetLogger(l)
}
