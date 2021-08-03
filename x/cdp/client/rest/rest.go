package rest

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/gorilla/mux"
)

func registerQueryRoutes(ctx client.Context, router *mux.Router) {
	router.HandleFunc("/cdps/id/{id}", queryCdp(ctx)).
		Methods("GET")
	router.HandleFunc("/params", queryParams(ctx)).
		Methods("GET")
	router.HandleFunc("/deposits/collateral_type/owner", queryDeposits(ctx)).
		Methods("GET")

}

func registerTxRoutes(ctx client.Context, router *mux.Router) {
	router.HandleFunc("/cdps/id", postCdp(ctx)).
		Methods("POST")
	router.HandleFunc("/deposit", postDeposit(ctx)).
		Methods("POST")
	router.HandleFunc("/withdraw", postWithdraw(ctx)).
		Methods("POST")
	router.HandleFunc("/draw", postDrawDebt(ctx)).
		Methods("POST")
	router.HandleFunc("/repay", postRepay(ctx)).
		Methods("POST")
	router.HandleFunc("/liquidate", postLiquidate(ctx)).
		Methods("POST")

}
func RegisterRoutes(ctx client.Context, router *mux.Router) {
	registerQueryRoutes(ctx, router)
	registerTxRoutes(ctx, router)
}
