//go:build regtest
// +build regtest

package app

import (
	"net/http"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

var (
	begin = make(chan struct{})
	end   = make(chan struct{})
)

func init() {
	// start an http server to unblock a block creation when a request is received
	newBlock := func(w http.ResponseWriter, r *http.Request) {
		begin <- struct{}{}
		<-end
	}
	http.HandleFunc("/newBlock", newBlock)
	go http.ListenAndServe(":8080", nil)
}

func (app *App) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	<-begin
	return app.mm.BeginBlock(ctx, req)
}

// EndBlocker application updates every end block
func (app *App) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	defer func() { end <- struct{}{} }()
	return app.mm.EndBlock(ctx, req)
}
