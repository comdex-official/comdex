package keeper

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	assettypes "github.com/comdex-official/comdex/x/asset/types"
	"github.com/comdex-official/comdex/x/lend/types"
)

func (k Keeper) VerifyCollateralizationRatio(
	ctx sdk.Context,
	amountIn sdkmath.Int,
	assetIn assettypes.Asset,
	amountOut sdkmath.Int,
	assetOut assettypes.Asset,
	liquidationThreshold sdkmath.LegacyDec,
) error {
	collateralizationRatio, err := k.CalculateCollateralizationRatio(ctx, amountIn, assetIn, amountOut, assetOut)
	if err != nil {
		return err
	}

	if collateralizationRatio.GT(liquidationThreshold) {
		return types.ErrorInvalidCollateralizationRatio
	}

	return nil
}

func (k Keeper) CalculateCollateralizationRatio(
	ctx sdk.Context,
	amountIn sdkmath.Int,
	assetIn assettypes.Asset,
	amountOut sdkmath.Int,
	assetOut assettypes.Asset,
) (sdkmath.LegacyDec, error) {
	totalIn, err := k.Market.CalcAssetPrice(ctx, assetIn.Id, amountIn)
	if err != nil {
		return sdkmath.LegacyZeroDec(), err
	}
	totalOut, err := k.Market.CalcAssetPrice(ctx, assetOut.Id, amountOut)
	if err != nil {
		return sdkmath.LegacyZeroDec(), err
	}

	return totalOut.Quo(totalIn), nil
}
