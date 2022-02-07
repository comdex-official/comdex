package expected

import (
	"context"

	vaulttypes "github.com/comdex-official/comdex/x/vault/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	liquiditytypes "github.com/tendermint/liquidity/x/liquidity/types"
)

type LiquidityKeeper interface {
	GetPool(ctx sdk.Context, poolID uint64) (liquiditytypes.Pool, bool)
	GetPoolMetaData(ctx sdk.Context, pool liquiditytypes.Pool) liquiditytypes.PoolMetadata
	GetAllPools(ctx sdk.Context) (pools []liquiditytypes.Pool)
}

type OracleKeeper interface {
	GetPriceForMarket(ctx sdk.Context, symbol string) (uint64, bool)
}

type VaultKeeper interface {
	QueryAllVaults(c context.Context, req *vaulttypes.QueryAllVaultsRequest) (*vaulttypes.QueryAllVaultsResponse, error)
}
