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

	totalIn, err := k.CalcAssetPrice(ctx, assetIn.Id, amountIn)
	if err != nil {
		return sdk.ZeroDec(), err
	}
	totalOut, err := k.CalcAssetPrice(ctx, assetOut.Id, amountOut)
	if err != nil {
		return sdk.ZeroDec(), err
	}
	// assetInPrice, found := k.GetPriceForAsset(ctx, assetIn.Id)
	// if !found {
	// 	return sdk.ZeroDec(), types.ErrorPriceInDoesNotExist
	// }

	// assetOutPrice, found := k.GetPriceForAsset(ctx, assetOut.Id)
	// if !found {
	// 	return sdk.ZeroDec(), types.ErrorPriceOutDoesNotExist
	// }

	// totalIn := amountIn.Mul(sdk.NewIntFromUint64(assetInPrice)).ToDec()
	// if totalIn.LTE(sdk.ZeroDec()) {
	// 	return sdk.ZeroDec(), types.ErrorInvalidAmountIn
	// }

	// totalOut := amountOut.Mul(sdk.NewIntFromUint64(assetOutPrice)).ToDec()
	// if totalOut.LTE(sdk.ZeroDec()) {
	// 	return sdk.ZeroDec(), types.ErrorInvalidAmountOut
	// }

	return sdk.NewDecFromInt(totalOut).Quo(sdk.NewDecFromInt(totalIn)), nil
}
