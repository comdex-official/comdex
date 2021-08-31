package rest

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/gorilla/mux"
)

func registerQueryRoutes(ctx client.Context, router *mux.Router) {
	router.HandleFunc("/cdps/id", queryCdp(ctx)).
		Methods("GET")
	router.HandleFunc("/params", queryParams(ctx)).
		Methods("GET")

}

func registerTxRoutes(ctx client.Context, router *mux.Router) {
	router.HandleFunc("/cdps/id", createCdp(ctx)).
		Methods("POST")
	router.HandleFunc("/deposit", createDeposit(ctx)).
		Methods("POST")
	router.HandleFunc("/withdraw", createWithdraw(ctx)).
		Methods("POST")
	router.HandleFunc("/draw", createDrawDebt(ctx)).
		Methods("POST")
	router.HandleFunc("/repay", createRepay(ctx)).
		Methods("POST")
	router.HandleFunc("/cdp/close", closeCDP(ctx)).
		Methods("POST")

}
func RegisterRoutes(ctx client.Context, router *mux.Router) {
	registerQueryRoutes(ctx, router)
	registerTxRoutes(ctx, router)
}
