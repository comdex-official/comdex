package keeper

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

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
		price_of_coin, _ := k.oracle.GetPriceForMarket(ctx, coin.Denom)

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

	var (
		total_liquidity uint64 = 0
		ctx                    = sdk.UnwrapSDKContext(c)
	)
	vaults := k.vault.GetVaults(ctx)

	for _, vault := range vaults {
		pair, _ := k.asset.GetPair(ctx, vault.PairID)
		assetIn, _ := k.asset.GetAsset(ctx, pair.AssetIn)
		collateral := sdk.NewCoin(assetIn.Denom, vault.AmountIn)
		amount := collateral.Amount
		denom := collateral.Denom
		price_of_coin, _ := k.oracle.GetPriceForMarket(ctx, denom)

		price_of_collateral := price_of_coin * amount.Uint64()
		total_liquidity = total_liquidity + price_of_collateral
	}

	return total_liquidity, true
}

type BankTotal struct {
	Height string
	Result []string
}

type Inflation struct {
	Height string
	Result string
}

type StakingPool struct {
	Height string
	Result []string
}

func (k *Keeper) GetAPR(c context.Context) (float64, bool) {
	var client http.Client
	var (
		apr float64 = 0.0
		ctx         = sdk.UnwrapSDKContext(c)
	)
	banktotal_res, err := client.Get("https://api-comdex.zenchainlabs.io/bank/total/ucmdx")
	if err != nil {
		return 0.0, false
	}

	inflation_res, err := http.Get("https://api-comdex.zenchainlabs.io/minting/inflation")
	if err != nil {
		return 0.0, false
	}

	stakingtokens_res, err := http.Get("https://api-comdex.zenchainlabs.io/staking/pool")
	if err != nil {
		return 0.0, false
	}

	decoder1 := json.NewDecoder(banktotal_res.Body)
	json1 := &BankTotal{}

	err = decoder1.Decode(json1)
	if err != nil {
		return 0.0, false
	}

	denom := json1.Result[0]
	amount, _ := strconv.ParseUint(json1.Result[1], 10, 64)
	denom_price, _ := k.oracle.GetPriceForMarket(ctx, denom)
	banktotal := amount * denom_price

	decoder2 := json.NewDecoder(inflation_res.Body)
	json2 := &Inflation{}

	err = decoder2.Decode(json2)
	if err != nil {
		return 0.0, false
	}

	inflation_string := json2.Result
	inflation, _ := strconv.ParseFloat(inflation_string, 64)

	decoder3 := json.NewDecoder(stakingtokens_res.Body)
	json3 := &StakingPool{}

	err = decoder3.Decode(json3)
	if err != nil {
		return 0.0, false
	}

	bondedtokens_string := json3.Result[1]
	bondedtokens, _ := strconv.ParseUint(bondedtokens_string, 10, 64)

	apr = inflation * float64(banktotal) / float64(bondedtokens)

	return apr, true
}
