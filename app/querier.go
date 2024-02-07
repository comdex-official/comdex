// this file used from osmosis, ref: https://github.com/CosmosContracts/juno/v16/blob/2ce971f4c6aa85d3ef7ba33d60e0ae74b923ab83/app/keepers/querier.go
// Original Author: https://github.com/nicolaslara

package app

import (
	storetypes "cosmossdk.io/store/types"
)

// QuerierWrapper is a local wrapper around BaseApp that exports only the Queryable interface.
// This is used to pass the baseApp to Async ICQ without exposing all methods
type QuerierWrapper struct {
	querier storetypes.Queryable
}

// var _ storetypes.Queryable = QuerierWrapper{}

func NewQuerierWrapper(querier storetypes.Queryable) QuerierWrapper {
	return QuerierWrapper{querier: querier}
}

func (q QuerierWrapper) Query(req storetypes.RequestQuery) storetypes.ResponseQuery {
	ResponseQuery, _ := q.querier.Query(&req)
	return *ResponseQuery
}
