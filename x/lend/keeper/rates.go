package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	assettypes "github.com/comdex-official/comdex/x/asset/types"
	"github.com/comdex-official/comdex/x/lend/types"
)

func (k Keeper) VerifyCollateralizationRatio(
	ctx sdk.Context,
	amountIn sdk.Int,
	assetIn assettypes.Asset,
	amountOut sdk.Int,
	assetOut assettypes.Asset,
	liquidationThreshold sdk.Dec,
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
	amountIn sdk.Int,
	assetIn assettypes.Asset,
	amountOut sdk.Int,
	assetOut assettypes.Asset,
) (sdk.Dec, error) {
	totalIn, err := k.Market.CalcAssetPrice(ctx, assetIn.Id, amountIn)
	if err != nil {
		return sdk.ZeroDec(), err
	}
	totalOut, err := k.Market.CalcAssetPrice(ctx, assetOut.Id, amountOut)
	if err != nil {
		return sdk.ZeroDec(), err
	}

	return totalOut.Quo(totalIn), nil
}
