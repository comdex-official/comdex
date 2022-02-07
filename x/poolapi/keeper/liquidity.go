package keeper

import (
	"context"

	"github.com/comdex-official/comdex/x/vault/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k *Keeper) GetLiquidity(ctx sdk.Context, pool_id uint64) (uint64, bool) {
	pool, found := k.liquiditykeeper.GetPool(ctx, pool_id)
	if !found {
		return 0, false
	}

	pool_metadata := k.liquiditykeeper.GetPoolMetaData(ctx, pool)
	reserve_coins := pool_metadata.ReserveCoins

	var pool_liquidity uint64 = 0
	for _, coin := range reserve_coins {
		amount := reserve_coins.AmountOf(coin.Denom)
		price_of_coin, _ := k.oraclekeeper.GetPriceForMarket(ctx, coin.Denom)

		price_of_all_coins := price_of_coin * amount.Uint64()
		pool_liquidity = pool_liquidity + price_of_all_coins
	}

	return pool_liquidity, true
}

func (k *Keeper) PoolLiquidity(ctx sdk.Context, pool_id uint64) (uint64, bool) {
	return k.GetLiquidity(ctx, pool_id)
}

func (k *Keeper) TotalLiquidity(ctx sdk.Context) (uint64, bool) {
	pools := k.liquiditykeeper.GetAllPools(ctx)

	var (
		total_liquidity uint64 = 0
	)

	for i := 0; i < len(pools); i++ {
		pool_id := pools[i].Id
		liquidity_of_ith_pool, found := k.GetLiquidity(ctx, pool_id)
		if !found {
			return 0, false
		}
		total_liquidity = total_liquidity + liquidity_of_ith_pool
	}

	return total_liquidity, true
}

func (k *Keeper) GetTotalCollateral(c context.Context) (uint64, bool) {
	req := &types.QueryAllVaultsRequest{}
	vaults, err := k.vaultkeeper.QueryAllVaults(c, req)
	if err != nil {
		return 0, false
	}

	var (
		total_liquidity uint64 = 0
		ctx                    = sdk.UnwrapSDKContext(c)
	)

	for _, vault := range vaults.VaultsInfo {
		collateral := vault.Collateral
		amount := collateral.Amount
		denom := collateral.Denom
		price_of_coin, found := k.oraclekeeper.GetPriceForMarket(ctx, denom)
		if !found {
			return 0, false
		}

		price_of_collateral := price_of_coin * amount.Uint64()
		total_liquidity = total_liquidity + price_of_collateral
	}

	return total_liquidity, true
}
