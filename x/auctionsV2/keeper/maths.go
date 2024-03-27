package keeper

import sdk "github.com/cosmos/cosmos-sdk/types"

func Multiply(a, b sdk.Dec) sdk.Dec {
	return a.Mul(b)
}

func (k Keeper) GetCollalteralTokenInitialPrice(price sdk.Int, premium sdk.Dec) sdk.Dec {
	result := premium.Mul(sdk.NewDecFromInt(price))
	return result
}

func (k Keeper) GetPriceFromLinearDecreaseFunction(CollateralTokenAuctionPrice sdk.Dec, timeToReachZeroPrice, timeElapsed sdk.Int) sdk.Dec {
	timeDifference := timeToReachZeroPrice.Sub(timeElapsed)
	resultantPrice := CollateralTokenAuctionPrice.Mul(sdk.NewDecFromInt(timeDifference))
	currentPrice := resultantPrice.Quo(sdk.NewDecFromInt(timeToReachZeroPrice))
	return currentPrice
}

func (k Keeper) GetCollateralTokenEndPrice(price, cusp sdk.Dec) sdk.Dec {
	result := Multiply(price, cusp)
	return result
}
