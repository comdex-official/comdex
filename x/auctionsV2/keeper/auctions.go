package keeper

import (
	utils "github.com/comdex-official/comdex/types"
	"time"

	auctiontypes "github.com/comdex-official/comdex/x/auction/types"
	"github.com/comdex-official/comdex/x/auctionsV2/types"
	liquidationtypes "github.com/comdex-official/comdex/x/liquidationsV2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) AuctionActivator(ctx sdk.Context, liquidationData liquidationtypes.LockedVault) error {

	//Dutch Auction Model Followed for auction type true
	if liquidationData.AuctionType {

		//Trigger Dutch Auction
		err := k.DutchAuctionActivator(ctx, liquidationData)
		if err != nil {

		}

	} else {
		//Trigger English Auction
		err := k.EnglishAuctionActivator(ctx, liquidationData)
		if err != nil {

		}
	}

	return nil
}

func (k Keeper) DutchAuctionActivator(ctx sdk.Context, liquidationData liquidationtypes.LockedVault) error {

	//Getting previous auction ID
	auctionID := k.GetAuctionID(ctx)

	//Price Calculation Function to determine auction different stage price
	liquidationWhitelistingAppData, _ := k.LiquidationsV2.GetLiquidationWhiteListing(ctx, liquidationData.AppId)

	if !liquidationWhitelistingAppData.IsDutchActivated {
		return types.ErrDutchAuctionDisabled
	}
	dutchAuctionParams := liquidationWhitelistingAppData.DutchAuctionParam

	pair, _ := k.asset.GetPair(ctx, liquidationData.PairId)

	twaDataCollateral, found := k.market.GetTwa(ctx, pair.AssetIn)
	if !found || !twaDataCollateral.IsPriceActive {
		return auctiontypes.ErrorPrices
	}
	twaDataDebt, found := k.market.GetTwa(ctx, pair.AssetOut)
	if !found || !twaDataDebt.IsPriceActive {
		return auctiontypes.ErrorPrices
	}
	//Checking if DEBT  token is CMST  then setting its price to $1 , else all tokens price will come from oracle.
	if liquidationData.IsDebtCmst {
		twaDataDebt.Twa = 1000000
	}

	//Some params will come from the specific app and they could be configured by them ,
	//rest of the params like auction duration and fees and other params will be based on comdex to edit based on governance
	//Understanding different Params:
	//Premium : Initial Price i.e price of the collateral at which the auction will start
	//Discount: Final Price , i.e less than the oracle price of the collateral asset and at this , auction would end
	//Decrement Factor:     Linear decrease in the price of the collateral every block is governed by this.
	CollateralTokenInitialPrice := k.GetCollalteralTokenInitialPrice(sdk.NewIntFromUint64(twaDataCollateral.Twa), dutchAuctionParams.Premium)
	// CollateralTokenEndPrice := k.getOutflowTokenEndPrice(CollateralTokenInitialPrice, dutchAuctionParams.Cusp)
	auctionParams, _ := k.GetAuctionParams(ctx)

	//Saving liquidation data to the auction struct
	auctionData := types.Auctions{
		AuctionId:                   auctionID + 1,
		CollateralToken:             liquidationData.CollateralToken,
		DebtToken:                   liquidationData.TargetDebt,
		CollateralTokenAuctionPrice: CollateralTokenInitialPrice,
		CollateralTokenOraclePrice:  sdk.NewDecFromInt(sdk.NewInt(int64(twaDataCollateral.Twa))),
		DebtTokenOraclePrice:        sdk.NewDecFromInt(sdk.NewInt(int64(twaDataDebt.Twa))),
		LockedVaultId:               liquidationData.LockedVaultId,
		StartTime:                   ctx.BlockTime(),
		EndTime:                     ctx.BlockTime().Add(time.Second * time.Duration(auctionParams.AuctionDurationSeconds)),
		AppId:                       liquidationData.AppId,
		AuctionType:                 liquidationData.AuctionType,
	}

	k.SetAuctionID(ctx, auctionData.AuctionId)
	err := k.SetAuction(ctx, auctionData)
	if err != nil {
		return err
	}

	return nil
}

// AUCTIONITERATOR
// -> DUCTHAUCTIONITERATOR
// -> ENGLISHAUCTIONITERATOR
func (k Keeper) EnglishAuctionActivator(ctx sdk.Context, liquidationData liquidationtypes.LockedVault) error {

	//Getting previous auction ID
	auctionID := k.GetAuctionID(ctx)

	//Price Calculation Function to determine auction different stage price
	liquidationWhitelistingAppData, _ := k.LiquidationsV2.GetLiquidationWhiteListing(ctx, liquidationData.AppId)

	if !liquidationWhitelistingAppData.IsEnglishActivated {
		return types.ErrEnglishAuctionDisabled
	}
	// englishAuctionParams := liquidationWhitelistingAppData.EnglishAuctionParam
	auctionParams, _ := k.GetAuctionParams(ctx)

	auctionData := types.Auctions{
		AuctionId:       auctionID + 1,
		CollateralToken: liquidationData.CollateralToken,
		DebtToken:       liquidationData.TargetDebt,
		// CollateralTokenAuctionPrice: CollateralTokenInitialPrice,
		// CollateralTokenOraclePrice:  sdk.NewDecFromInt(sdk.NewInt(int64(twaDataCollateral.Twa))),
		// DebtTokenOraclePrice:        sdk.NewDecFromInt(sdk.NewInt(int64(twaDataDebt.Twa))),
		LockedVaultId: liquidationData.LockedVaultId,
		StartTime:     ctx.BlockTime(),
		EndTime:       ctx.BlockTime().Add(time.Second * time.Duration(auctionParams.AuctionDurationSeconds)),
		AppId:         liquidationData.AppId,
		AuctionType:   liquidationData.AuctionType,
	}
	k.SetAuctionID(ctx, auctionData.AuctionId)
	err := k.SetAuction(ctx, auctionData)
	if err != nil {
		return err
	}

	return nil

}

// DutchAuctionsIterator iterates over existing active dutch auctions and does 2 main job
// First: if auction time is complete and target not reached with collateral available then Restart
// Second: if not restarting update the price
func (k Keeper) DutchAuctionsIterator(ctx sdk.Context) error {
	dutchAuctions := k.GetAuctions(ctx)
	// SET current price of inflow token and outflow token

	for _, dutchAuction := range dutchAuctions {
		lockedVault, found := k.LiquidationsV2.GetLockedVault(ctx, dutchAuction.AppId, dutchAuction.LockedVaultId)
		if !found {
			return auctiontypes.ErrorInvalidLockedVault
		}
		_ = utils.ApplyFuncIfNoError(ctx, func(ctx sdk.Context) error {
			// First case to check if we have to restart the auction
			if ctx.BlockTime().After(dutchAuction.EndTime) {
				// restart
				err := k.RestartDutchAuctions(ctx, dutchAuction, lockedVault)
				if err != nil {
					return err
				}

			} else {
				// Second case to only reduce the price
				err := k.UpdateDutchAuctionsPrice(ctx, dutchAuction, lockedVault)
				if err != nil {
					return err
				}
			}
			return nil
		})
	}
	return nil
}

func (k Keeper) RestartDutchAuctions(ctx sdk.Context, dutchAuction types.Auctions, lockedVault liquidationtypes.LockedVault) error {
	esmStatus, found := k.esm.GetESMStatus(ctx, dutchAuction.AppId)
	auctionParams, _ := k.GetAuctionParams(ctx)
	status := false
	if found {
		status = esmStatus.Status
	}
	if lockedVault.IsInternalKeeper {
		if status {
			// code goes here for collateral redemption in case of esm
			return nil
		}
	}
	liquidationWhitelistingAppData, _ := k.LiquidationsV2.GetLiquidationWhiteListing(ctx, dutchAuction.AppId)

	if !liquidationWhitelistingAppData.IsDutchActivated {
		return types.ErrDutchAuctionDisabled
	}
	dutchAuctionParams := liquidationWhitelistingAppData.DutchAuctionParam

	pair, _ := k.asset.GetPair(ctx, lockedVault.PairId)
	twaDataCollateral, found := k.market.GetTwa(ctx, pair.AssetIn)
	if !found || !twaDataCollateral.IsPriceActive {
		return auctiontypes.ErrorPrices
	}

	collateralTokenCurrentPrice := twaDataCollateral.Twa
	timeNow := ctx.BlockTime()
	dutchAuction.StartTime = timeNow
	dutchAuction.EndTime = timeNow.Add(time.Second * time.Duration(auctionParams.AuctionDurationSeconds))
	collateralTokenInitialPrice := k.GetCollalteralTokenInitialPrice(sdk.NewIntFromUint64(collateralTokenCurrentPrice), dutchAuctionParams.Premium)
	dutchAuction.CollateralTokenAuctionPrice = collateralTokenInitialPrice
	err := k.SetAuction(ctx, dutchAuction)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) UpdateDutchAuctionsPrice(ctx sdk.Context, dutchAuction types.Auctions, lockedVault liquidationtypes.LockedVault) error {

	auctionParams, _ := k.GetAuctionParams(ctx)
	liquidationWhitelistingAppData, _ := k.LiquidationsV2.GetLiquidationWhiteListing(ctx, dutchAuction.AppId)

	if !liquidationWhitelistingAppData.IsDutchActivated {
		return types.ErrDutchAuctionDisabled
	}
	dutchAuctionParams := liquidationWhitelistingAppData.DutchAuctionParam

	pair, _ := k.asset.GetPair(ctx, lockedVault.PairId)
	twaDataDebt, found := k.market.GetTwa(ctx, pair.AssetOut)
	if !found || !twaDataDebt.IsPriceActive {
		return auctiontypes.ErrorPrices
	}
	//Checking if DEBT  token is CMST  then setting its price to $1 , else all tokens price will come from oracle.
	if lockedVault.IsDebtCmst {
		twaDataDebt.Twa = 1000000
	}

	numerator := dutchAuction.CollateralTokenAuctionPrice.Mul(sdk.NewDecFromInt(sdk.NewIntFromUint64(auctionParams.AuctionDurationSeconds))) //cmdx
	CollateralTokenAuctionEndPrice := k.GetCollateralTokenEndPrice(dutchAuction.CollateralTokenAuctionPrice, dutchAuctionParams.Discount)

	denominator := dutchAuction.CollateralTokenAuctionPrice.Sub(CollateralTokenAuctionEndPrice)
	resultant := numerator.Quo(denominator)
	tau := sdk.NewInt(resultant.TruncateInt64())
	dur := ctx.BlockTime().Sub(dutchAuction.StartTime)
	seconds := sdk.NewInt(int64(dur.Seconds()))
	collateralTokenAuctionPrice := k.GetPriceFromLinearDecreaseFunction(dutchAuction.CollateralTokenAuctionPrice, tau, seconds)
	dutchAuction.DebtTokenOraclePrice = sdk.NewDec(int64(twaDataDebt.Twa))
	dutchAuction.CollateralTokenAuctionPrice = collateralTokenAuctionPrice
	err := k.SetAuction(ctx, dutchAuction)
	if err != nil {
		return err
	}
	return nil
}
