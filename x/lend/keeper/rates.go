package keeper

import (
	assettypes "github.com/comdex-official/comdex/x/asset/types"
	"github.com/comdex-official/comdex/x/lend/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k *Keeper) VerifyCollaterlizationRatio(
	ctx sdk.Context,
	amountIn sdk.Int,
	assetIn assettypes.Asset,
	amountOut sdk.Int,
	assetOut assettypes.Asset,
	liquidationThreshold sdk.Dec,
) error {
	collaterlizationRatio, err := k.CalculateCollaterlizationRatio(ctx, amountIn, assetIn, amountOut, assetOut)
	if err != nil {
		return err
	}

	if collaterlizationRatio.GT(liquidationThreshold) {
		return types.ErrorInvalidCollateralizationRatio
	}

	return nil
}

func (k *Keeper) CalculateCollaterlizationRatio(
	ctx sdk.Context,
	amountIn sdk.Int,
	assetIn assettypes.Asset,
	amountOut sdk.Int,
	assetOut assettypes.Asset,
) (sdk.Dec, error) {
	assetInPrice, found := k.GetPriceForAsset(ctx, assetIn.Id)
	if !found {
		return sdk.ZeroDec(), types.ErrorPriceInDoesNotExist
	}

	assetOutPrice, found := k.GetPriceForAsset(ctx, assetOut.Id)
	if !found {
		return sdk.ZeroDec(), types.ErrorPriceOutDoesNotExist
	}

	totalIn := amountIn.Mul(sdk.NewIntFromUint64(assetInPrice)).ToDec()
	if totalIn.LTE(sdk.ZeroDec()) {
		return sdk.ZeroDec(), types.ErrorInvalidAmountIn
	}

	totalOut := amountOut.Mul(sdk.NewIntFromUint64(assetOutPrice)).ToDec()
	if totalOut.LTE(sdk.ZeroDec()) {
		return sdk.ZeroDec(), types.ErrorInvalidAmountOut
	}

	return totalOut.Quo(totalIn), nil
}
