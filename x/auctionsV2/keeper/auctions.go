package keeper

import (
	"time"

	utils "github.com/comdex-official/comdex/types"

	"github.com/comdex-official/comdex/x/auctionsV2/types"
	auctionsV2types "github.com/comdex-official/comdex/x/auctionsV2/types"
	collectortypes "github.com/comdex-official/comdex/x/collector/types"
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

	// pair, _ := k.asset.GetPair(ctx, liquidationData.PairId)

	twaDataCollateral, found := k.market.GetTwa(ctx, liquidationData.CollateralAssetId)
	if !found || !twaDataCollateral.IsPriceActive {
		return auctionsV2types.ErrorPriceNotFound
	}
	twaDataDebt, found := k.market.GetTwa(ctx, liquidationData.DebtAssetId)
	if !found || !twaDataDebt.IsPriceActive {
		return auctionsV2types.ErrorPriceNotFound
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
	auctionData := types.Auction{
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
		CollateralAssetId:           liquidationData.CollateralAssetId,
		DebtAssetId:                 liquidationData.DebtAssetId,
		BonusAmount:                 liquidationData.BonusToBeGiven,
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
	auctionData := types.Auction{
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

func (k Keeper) AuctionIterator(ctx sdk.Context) error {

	auctions := k.GetAuctions(ctx)
	//Dutch Auction Model Followed for auction type true
	for _, auction := range auctions {
		_ = utils.ApplyFuncIfNoError(ctx, func(ctx sdk.Context) error {

			//Dutch Auction= true
			if auction.AuctionType {

				esmStatus, found := k.esm.GetESMStatus(ctx, auction.AppId)
				//FOr esm , can also check vault as initiator exists or not just to be sure
				if found && esmStatus.Status {
					//Checking if auction price is supposed to be reduced or restared

					//Checking condition

					if ctx.BlockTime().After(auction.EndTime) {
						//If restart - DO ESM specific operation
						//Most Probably Close Auction

						//Check here if initiator is vault , then for vault do esm trigger option accordingly
						liquidationData, _ := k.LiquidationsV2.GetLockedVault(ctx, auction.AppId, auction.LockedVaultId)
						if liquidationData.InitiatorType == "vault" {

							err := k.TriggerEsm(ctx, auction, liquidationData)
							if err != nil {
								return err
							}

						}

					} else {
						//Else reduce - normal operation
						err := k.UpdateDutchAuction(ctx, auction)
						if err != nil {
							return err
						}

					}

				} else if !found || !esmStatus.Status {
					//This app is not eligible for ESM
					//Continue normal operation

					//DO update
					//then check if to be restarred , then restart

					if ctx.BlockTime().After(auction.EndTime) {
						//Restart
						err := k.RestartDutchAuction(ctx, auction)
						if err != nil {
							return err
						}

					} else {
						//Else udate price params
						err := k.UpdateDutchAuction(ctx, auction)
						if err != nil {
							return err
						}

					}

				}

			} else {
				//English Auction=false
				//Check if auction time has ended, then close auction - if there is at least a single bid
				//English auction does not require price so no important operation
				if ctx.BlockTime().After(auction.EndTime) {

					if auction.ActiveBiddingId != uint64(0) {
						//If atleast there is one bidding on the auction
						err := k.CloseEnglishAuction(ctx, auction)
						if err != nil {
							return err
						}

					} else {
						//Restart the auction by updating the end time param
						err := k.RestartEnglishAuction(ctx, auction)
						if err != nil {
							return err
						}
					}

				}

			}
			return nil
		})
	}

	return nil
}

// DutchAuctionsIterator iterates over existing active dutch auctions and does 2 main job
// First: if auction time is complete and target not reached with collateral available then Restart
// Second: if not restarting update the price
// func (k Keeper) DutchAuctionsIterator(ctx sdk.Context) error {
// 	dutchAuctions := k.GetAuctions(ctx)
// 	// SET current price of inflow token and outflow token

// 	for _, dutchAuction := range dutchAuctions {
// 		lockedVault, found := k.LiquidationsV2.GetLockedVault(ctx, dutchAuction.AppId, dutchAuction.LockedVaultId)
// 		if !found {
// 			return auctiontypes.ErrorInvalidLockedVault
// 		}
// 		_ = utils.ApplyFuncIfNoError(ctx, func(ctx sdk.Context) error {
// 			// First case to check if we have to restart the auction
// 			if ctx.BlockTime().After(dutchAuction.EndTime) {
// 				// restart
// 				err := k.RestartDutchAuctions(ctx, dutchAuction, lockedVault)
// 				if err != nil {
// 					return err
// 				}

// 			} else {
// 				// Second case to only reduce the price
// 				err := k.UpdateDutchAuctionPrice(ctx, dutchAuction)
// 				if err != nil {
// 					return err
// 				}
// 			}
// 			return nil
// 		})
// 	}
// 	return nil
// }

func (k Keeper) RestartDutchAuction(ctx sdk.Context, dutchAuction types.Auction) error {
	auctionParams, _ := k.GetAuctionParams(ctx)
	liquidationWhitelistingAppData, _ := k.LiquidationsV2.GetLiquidationWhiteListing(ctx, dutchAuction.AppId)

	dutchAuctionParams := liquidationWhitelistingAppData.DutchAuctionParam

	twaDataCollateral, found := k.market.GetTwa(ctx, dutchAuction.CollateralAssetId)
	if !found || !twaDataCollateral.IsPriceActive {
		return auctionsV2types.ErrorPriceNotFound
	}
	twaDataDebt, found := k.market.GetTwa(ctx, dutchAuction.DebtAssetId)
	if !found || !twaDataDebt.IsPriceActive {
		return auctionsV2types.ErrorPriceNotFound
	}
	liquidationData, _ := k.LiquidationsV2.GetLockedVault(ctx, dutchAuction.AppId, dutchAuction.LockedVaultId)
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

	//Saving liquidation data to the auction struct
	//Only updating necessary params

	dutchAuction.CollateralTokenAuctionPrice = CollateralTokenInitialPrice
	dutchAuction.CollateralTokenOraclePrice = sdk.NewDecFromInt(sdk.NewInt(int64(twaDataCollateral.Twa)))
	dutchAuction.DebtTokenOraclePrice = sdk.NewDecFromInt(sdk.NewInt(int64(twaDataDebt.Twa)))
	// dutchAuction.StartTime = ctx.BlockTime()
	dutchAuction.EndTime = ctx.BlockTime().Add(time.Second * time.Duration(auctionParams.AuctionDurationSeconds))

	err := k.SetAuction(ctx, dutchAuction)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) UpdateDutchAuction(ctx sdk.Context, dutchAuction types.Auction) error {
	auctionParams, _ := k.GetAuctionParams(ctx)
	liquidationWhitelistingAppData, _ := k.LiquidationsV2.GetLiquidationWhiteListing(ctx, dutchAuction.AppId)

	dutchAuctionParams := liquidationWhitelistingAppData.DutchAuctionParam

	twaDataCollateral, found := k.market.GetTwa(ctx, dutchAuction.CollateralAssetId)
	if !found || !twaDataCollateral.IsPriceActive {
		return auctionsV2types.ErrorPriceNotFound
	}
	twaDataDebt, found := k.market.GetTwa(ctx, dutchAuction.DebtAssetId)
	if !found || !twaDataDebt.IsPriceActive {
		return auctionsV2types.ErrorPriceNotFound
	}
	liquidationData, _ := k.LiquidationsV2.GetLockedVault(ctx, dutchAuction.AppId, dutchAuction.LockedVaultId)
	//Checking if DEBT  token is CMST  then setting its price to $1 , else all tokens price will come from oracle.
	if liquidationData.IsDebtCmst {
		twaDataDebt.Twa = 1000000
	}

	//Now calculating the auction price of the Collateral Token
	dutchAuction.CollateralTokenOraclePrice = sdk.NewDecFromInt(sdk.NewInt(int64(twaDataCollateral.Twa)))
	dutchAuction.DebtTokenOraclePrice = sdk.NewDecFromInt(sdk.NewInt(int64(twaDataDebt.Twa)))

	numerator := dutchAuction.CollateralTokenAuctionPrice.Mul(sdk.NewDecFromInt(sdk.NewIntFromUint64(auctionParams.AuctionDurationSeconds))) //cmdx
	CollateralTokenAuctionEndPrice := k.GetCollateralTokenEndPrice(dutchAuction.CollateralTokenAuctionPrice, dutchAuctionParams.Discount)
	denominator := dutchAuction.CollateralTokenAuctionPrice.Sub(CollateralTokenAuctionEndPrice)
	timeToReachZeroPrice := numerator.Quo(denominator)
	timeElapsed := ctx.BlockTime().Sub(dutchAuction.StartTime)
	// Example: CollateralTokenAuctionPrice = 1.2
	// AuctionDurationSeconds = 10 unit
	// numerator = 1.2*10 = 12, CollateralTokenAuctionEndPrice = 1.2*0.7 = 0.84
	// denominator = 1.2- 0.84 = 0.36
	// timeToReachZeroPrice = 12/0.36 = 33.3 unit
	// now assuming auction just started
	// timeElapsed = 0
	// collateralTokenAuctionPrice = GetPriceFromLinearDecreaseFunction(1.2, 33.3, 0)
	// timeDifference = 33.3- 0 = 33.3
	// resultantPrice = 1.2 *33.3
	// currentPrice = 1.2*33.3/33.3 = 1.2 unit
	collateralTokenAuctionPrice := k.GetPriceFromLinearDecreaseFunction(dutchAuction.CollateralTokenAuctionPrice, sdk.NewInt(timeToReachZeroPrice.TruncateInt64()), sdk.NewInt(int64(timeElapsed.Seconds())))
	dutchAuction.CollateralTokenAuctionPrice = collateralTokenAuctionPrice

	err := k.SetAuction(ctx, dutchAuction)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) RestartEnglishAuction(ctx sdk.Context, englishAuction types.Auction) error {

	auctionParams, _ := k.GetAuctionParams(ctx)
	englishAuction.EndTime = ctx.BlockTime().Add(time.Second * time.Duration(auctionParams.AuctionDurationSeconds))
	err := k.SetAuction(ctx, englishAuction)
	if err != nil {
		return err
	}
	return nil

}

func (k Keeper) CloseEnglishAuction(ctx sdk.Context, englishAuction types.Auction) error {

	//Send Collateral To the user
	//Delete Auction Data

	return nil

}

func (k Keeper) TriggerEsm(ctx sdk.Context, auctionData types.Auction, liquidationData liquidationtypes.LockedVault) error {

	//Check if liquidation penalty has been recovered
	debtCollected := liquidationData.TargetDebt.Sub(auctionData.DebtToken)
	collateralAuctioned := liquidationData.CollateralToken.Amount.Sub(auctionData.CollateralToken.Amount)
	tokensToTransfer := debtCollected
	//If more debt collected, send liquidation penalty to collector, and open the vault from the rest amount and update params
	if debtCollected.Amount.GT(liquidationData.FeeToBeCollected) {
		//Send Liquidation Penalty to the Collector Module
		tokensToTransfer = sdk.NewCoin(auctionData.DebtToken.Denom, liquidationData.FeeToBeCollected)
		//burning rest collected tokens
		tokensToBurn := debtCollected.Amount.Sub(liquidationData.FeeToBeCollected)
		if tokensToBurn.GT(sdk.ZeroInt()) {
			err := k.bankKeeper.BurnCoins(ctx, auctionsV2types.ModuleName, sdk.NewCoins(sdk.NewCoin(auctionData.DebtToken.Denom, tokensToBurn)))
			if err != nil {
				return err
			}
		}
		//updating token minted
		//updating collateral locked data
		k.vault.UpdateTokenMintedAmountLockerMapping(ctx, auctionData.AppId, liquidationData.ExtendedPairId, tokensToBurn, false)
	}

	err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, auctionsV2types.ModuleName, collectortypes.ModuleName, sdk.NewCoins(tokensToTransfer))
	if err != nil {
		return err
	}
	//Update Collector Data for CMST
	// Updating fees data in collector
	err = k.collector.SetNetFeeCollectedData(ctx, auctionData.AppId, auctionData.CollateralAssetId, tokensToTransfer.Amount)
	if err != nil {
		return err
	}
	//Opening vault
	k.vault.CreateNewVault(ctx, liquidationData.Owner, auctionData.AppId, liquidationData.ExtendedPairId, auctionData.CollateralToken.Amount, auctionData.DebtToken.Amount)
	k.vault.UpdateCollateralLockedAmountLockerMapping(ctx, auctionData.AppId, liquidationData.ExtendedPairId, collateralAuctioned, false)

	return nil

}

func (k Keeper) LimitOrderBid(ctx sdk.Context) error {
	// Get Auctions One by One and for that particular auction check the current discount
	// if we find any active limit bid for that premium then we will execute it and update both

	auctions := k.GetAuctions(ctx)
	for _, auction := range auctions {
		_ = utils.ApplyFuncIfNoError(ctx, func(ctx sdk.Context) error {
			if auction.CollateralTokenOraclePrice.GT(auction.CollateralTokenAuctionPrice) {
				premium := (auction.CollateralTokenOraclePrice.Sub(auction.CollateralTokenAuctionPrice)).Quo(auction.CollateralTokenOraclePrice)
				premiumPerc := premium.Mul(sdk.NewDecFromInt(sdk.NewInt(100)))
				biddingData, found := k.GetUserLimitBidDataByPremium(ctx, auction.DebtAssetId, auction.CollateralAssetId, premiumPerc.TruncateInt().String())
				if !found {
					return nil
				}
				// Here we will check if the auction amount is greater than individual bids or vice versa,
				// in any of the case update both user bids and individual auctions

				for _, individualBids := range biddingData {
					addr, _ := sdk.AccAddressFromBech32(individualBids.BidderAddress)
					if individualBids.DebtToken.Amount.GTE(auction.DebtToken.Amount) {
						//User has more tokens than target debt, so their bid will close the auction
						///Placing a user bid
						bidding_id, err := k.PlaceDutchAuctionBid(ctx, auction.AuctionId, addr, individualBids.DebtToken, auction, true)
						if err != nil {
							return err
						}
						if individualBids.DebtToken.Amount.Equal(auction.DebtToken.Amount) {
							k.DeleteUserLimitBidData(ctx, auction.DebtAssetId, auction.CollateralAssetId, premiumPerc.TruncateInt().String(), individualBids.BidderAddress)
							return nil
						}
						individualBids.DebtToken.Amount = individualBids.DebtToken.Amount.Sub(auction.DebtToken.Amount)
						individualBids.BiddingId = append(individualBids.BiddingId, bidding_id)
						k.SetUserLimitBidData(ctx, individualBids, auction.DebtAssetId, auction.CollateralAssetId, premiumPerc.TruncateInt().String())
					} else {
						_, err := k.PlaceDutchAuctionBid(ctx, auction.AuctionId, addr, individualBids.DebtToken, auction, true)
						if err != nil {
							return err
						}
						k.DeleteUserLimitBidData(ctx, auction.DebtAssetId, auction.CollateralAssetId, premiumPerc.TruncateInt().String(), individualBids.BidderAddress)
					}

				}

			}
			return nil
		})
	}
	return nil
}
