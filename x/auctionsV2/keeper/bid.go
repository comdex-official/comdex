package keeper

import (
	assettypes "github.com/comdex-official/comdex/x/asset/types"
	"github.com/comdex-official/comdex/x/auctionsV2/types"
	auctionsV2types "github.com/comdex-official/comdex/x/auctionsV2/types"
	collectortypes "github.com/comdex-official/comdex/x/collector/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	protobuftypes "github.com/gogo/protobuf/types"
)

func (k Keeper) PlaceDutchAuctionBid(ctx sdk.Context, auctionID uint64, bidder sdk.AccAddress, bid sdk.Coin, auctionData types.Auction) error {
	//The bid is in debt token - This is different from the earliar auction model at comdex
	if bid.Amount.Equal(sdk.ZeroInt()) {
		return sdkerrors.Wrapf(sdkerrors.ErrNotFound, "Bid amount can't be Zero")
	}
	liquidationWhitelistingAppData, _ := k.LiquidationsV2.GetLiquidationWhiteListing(ctx, auctionData.AppId)

	if bid.Denom != auctionData.DebtToken.Denom {
		return sdkerrors.Wrapf(sdkerrors.ErrNotFound, "Bid token is not the debt token ", bid.Denom)
	}
	liquidationData, _ := k.LiquidationsV2.GetLockedVault(ctx, auctionData.AppId, auctionData.LockedVaultId)
	//Price data of the token from market module
	debtToken, _ := k.market.GetTwa(ctx, auctionData.DebtAssetId)
	debtPrice := sdk.NewDecFromInt(sdk.NewInt(int64(debtToken.Twa)))
	//Price data of the token from market module
	collateralToken, _ := k.market.GetTwa(ctx, auctionData.CollateralAssetId)
	collateralPrice := sdk.NewDecFromInt(sdk.NewInt(int64(collateralToken.Twa)))

	//only if debt token is CMST , we consider it as $1
	if liquidationData.IsDebtCmst {
		debtPrice = sdk.NewDecFromInt(sdk.NewInt(int64(1000000)))

	}
	isBidFinalBid := false
	//If user has sent a bigger bid than the target amount ,
	if bid.Amount.GTE(auctionData.DebtToken.Amount) {
		bid.Amount = auctionData.DebtToken.Amount
		isBidFinalBid = true
		// bidPercent := 0
		debtuDollar, collateralTokenQuanitity, _ := k.vault.GetAmountOfOtherToken(ctx, auctionData.DebtAssetId, debtPrice, bid.Amount, auctionData.CollateralAssetId, auctionData.CollateralTokenAuctionPrice)
		//From auction bonus quantity , use the available quantity to calculate the collateral value
		_, collateralTokenQuanitityForBonus, _ := k.vault.GetAmountOfOtherToken(ctx, auctionData.DebtAssetId, debtPrice, liquidationData.BonusToBeGiven, auctionData.CollateralAssetId, auctionData.CollateralTokenAuctionPrice)

		//Checking if the auction bonus and the collateral to be given to user isnt more than available colalteral
		totalCollateralTokenQuanitity := collateralTokenQuanitity.Add(collateralTokenQuanitityForBonus)
		if totalCollateralTokenQuanitity.LTE(auctionData.CollateralToken.Amount) {
			//If everything is correct

			//Take Debt Token from user ,
			if bid.Amount.GT(sdk.ZeroInt()) {
				err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, bidder, auctionsV2types.ModuleName, sdk.NewCoins(sdk.NewCoin(auctionData.DebtToken.Denom, bid.Amount)))
				if err != nil {
					return err
				}
			}

			//Send Collateral To bidder
			if bid.Amount.GT(sdk.ZeroInt()) {
				err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, auctionsV2types.ModuleName, bidder, sdk.NewCoins(sdk.NewCoin(auctionData.CollateralToken.Denom, totalCollateralTokenQuanitity)))
				if err != nil {
					return err
				}
			}

			//Burn Debt Token,
			liquidationPenalty := sdk.NewCoin(auctionData.DebtToken.Denom, liquidationData.FeeToBeCollected)
			tokensToBurn := auctionData.DebtToken.Sub(liquidationPenalty)

			if tokensToBurn.Amount.GT(sdk.ZeroInt()) {
				err := k.bankKeeper.BurnCoins(ctx, auctionsV2types.ModuleName, sdk.NewCoins(tokensToBurn))
				if err != nil {
					return err
				}
			}

			//Send rest tokens to the user
			OwnerLeftOverCapital := auctionData.CollateralToken.Amount.Sub(totalCollateralTokenQuanitity)
			if bid.Amount.GT(sdk.ZeroInt()) {
				err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, auctionsV2types.ModuleName, bidder, sdk.NewCoins(sdk.NewCoin(auctionData.CollateralToken.Denom, OwnerLeftOverCapital)))
				if err != nil {
					return err
				}
			}
			//Add bid data to struct
			//Creating user bid struct
			bidding_id, err := k.CreateUserBid(ctx, auctionData.AppId, string(bidder), auctionID, sdk.NewCoin(auctionData.CollateralToken.Denom, totalCollateralTokenQuanitity), sdk.NewCoin(auctionData.DebtToken.Denom, bid.Amount), "dutch")
			if err != nil {
				return err
			}
			//Based on app type call perform specific function - external , internal and /or keeper incentive
			//See if this was keeper initiated transaction- then incentivisation will be in place based on the percentage
			//For apps that are external to comdex chain
			if liquidationData.InitiatorType == "external" {

				//but if an app is external - will have to check the auction bonus , liquidation penalty , module account mechanism

			} else if liquidationData.InitiatorType == "vault" {
				//Check if they are initiated through a keeper, if so they will be incentivised
				if liquidationData.IsInternalKeeper {

					keeperIncentive := (liquidationWhitelistingAppData.KeeeperIncentive.Mul(sdk.NewDecFromInt(liquidationPenalty.Amount))).TruncateInt()
					if keeperIncentive.GT(sdk.ZeroInt()) {
						liquidationPenalty = liquidationPenalty.Sub(sdk.NewCoin(auctionData.DebtToken.Denom, keeperIncentive))
						err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, auctionsV2types.ModuleName, sdk.AccAddress(liquidationData.InternalKeeperAddress), sdk.NewCoins(sdk.NewCoin(auctionData.DebtToken.Denom, keeperIncentive)))
						if err != nil {
							return err
						}

					}
				}
				//Send Liquidation Penalty to the Collector Module
				if liquidationPenalty.Amount.GT(sdk.ZeroInt()) {
					err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, auctionsV2types.ModuleName, collectortypes.ModuleName, sdk.NewCoins(liquidationPenalty))
					if err != nil {
						return err
					}
				}
				//Update Collector Data for CMST
				// Updating fees data in collector
				err = k.collector.SetNetFeeCollectedData(ctx, auctionData.AppId, auctionData.CollateralAssetId, liquidationPenalty.Amount)
				if err != nil {
					return err
				}
				//Updating mapping data of vault
				k.vault.UpdateTokenMintedAmountLockerMapping(ctx, auctionData.AppId, liquidationData.ExtendedPairId, tokensToBurn.Amount, false)
				k.vault.UpdateCollateralLockedAmountLockerMapping(ctx, auctionData.AppId, liquidationData.ExtendedPairId, auctionData.CollateralToken.Amount, false)

			} else if liquidationData.InitiatorType == "borrow" {
				//Check if they are initiated through a keeper, if so they will be incentivised

			}

			//Add bidder data in auction
			bidOwnerMapppingData := auctionsV2types.BidOwnerMapping{bidding_id, string(bidder)}
			auctionData.BiddingIds = append(auctionData.BiddingIds, &bidOwnerMapppingData)
			//Savinga auction data to auction historical
			auctionHistoricalData := auctionsV2types.AuctionHistorical{auctionID, &auctionData, &liquidationData}
			k.SetAuctionHistorical(ctx, auctionHistoricalData)
			//Close Auction
			k.DeleteAuction(ctx, auctionData)
			//Delete liquidation Data
			k.LiquidationsV2.DeleteLockedVault(ctx, liquidationData.LockedVaultId)

		} else {
			//This means that there is less collateral available .
			//So we first try to compensate the difference through the liquidation penalty

			//check the difference in collateral -
			//check if nullifing liquidation penalty helps
			//if yes - go for it

			//else call the module account to give funds to compensate the user.

			leftOverCollateral := auctionData.CollateralToken.Amount

			_, debtTokenAgainstLeftOverCollateral, _ := k.vault.GetAmountOfOtherToken(ctx, auctionData.CollateralAssetId, auctionData.CollateralTokenAuctionPrice, leftOverCollateral, auctionData.DebtAssetId, debtPrice)

			//Amount to call from reserve account for adjusting the auction target debt
			debtGettingLeft := auctionData.DebtToken.Sub(sdk.NewCoin(auctionData.DebtToken.Denom, debtTokenAgainstLeftOverCollateral))

			//Calling reserve account for debt adjustment : debtGettingLeft

			//Taking debtTokenAgainstLeftOverCollateral from user
			//Sending leftOverCollateral to the user
			//Burn Debt Token,
			//Creating user bid struct

			//Based on app type call perform specific function - external , internal and /or keeper incentive
			//See if this was keeper initiated transaction- then incentivisation will be in place based on the percentage
			//For apps that are external to comdex chain

			//Add bidder data in auction

		}

	} else {
		//if bid amount is less than the target bid

		//Checking if bid isnt leaving dust amount less than allowed -for collateral & debt

		//Calculating collateral token value from bid(debt) token value
		debtuDollar, collateralTokenQuanitity, _ := k.vault.GetAmountOfOtherToken(ctx, auctionData.DebtAssetId, debtPrice, bid.Amount, auctionData.CollateralAssetId, collateralPrice)
		//From auction bonus quantity , use the available quantity to calculate the collateral value

		//Checking bid.Amount -> to targetbid ratio
		//using that ratio data to calculate  auction bonus to be given for the bid
		//first taking the debt percentage data
		//then calculating the collateral token data
		_, collateralTokenQuanitityForBonus, _ := k.vault.GetAmountOfOtherToken(ctx, auctionData.DebtAssetId, debtPrice, liquidationData.BonusToBeGiven, auctionData.CollateralAssetId, auctionData.CollateralTokenAuctionPrice)

		if collateralTokenQuanitity.Add(collateralTokenQuanitityForBonus).LTE(auctionData.CollateralToken.Amount) {
			//If there is sufficient collalteral

		} else {

			//Not sure if this condition will arise in which partial bids also arent able to be fulfilled due to shortage of collateral token
			// Technically this case will also close the auction
		}

		//Deducting auction bonus value from liquidation data also for next bid.
	}
	//Deducting the auction bonus

	//Now checking if the bid is not the final bid, we will check the dust amount left by the bidder
	//if the dust check passes, it is good to go.
	//Dust check for debt token

	return nil
}

func (k Keeper) CreateUserBid(ctx sdk.Context, appID uint64, BidderAddress string, auctionID uint64, collateralToken sdk.Coin, debtToken sdk.Coin, bidType string) (bidding_id uint64, err error) {

	userBidId := k.GetUserBidID(ctx)
	bidding := auctionsV2types.Bid{
		BiddingId:             userBidId + 1,
		AuctionId:             auctionID,
		CollateralTokenAmount: collateralToken,
		DebtTokenAmount:       debtToken,
		BidderAddress:         BidderAddress,
		BiddingTimestamp:      ctx.BlockHeader().Time,
		AppId:                 appID,
		BidType:               bidType,
	}
	k.SetUserBidID(ctx, bidding.BiddingId)
	err = k.SetUserBid(ctx, bidding)
	if err != nil {
		return bidding_id, err
	}
	return bidding.BiddingId, nil
}

func (k Keeper) PlaceEnglishAuctionBid(ctx sdk.Context, auctionID uint64, bidder sdk.AccAddress, bid sdk.Coin, auctionData types.Auction) error {
	return nil
}

func (k Keeper) SetLimitAuctionBidID(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.LimitAuctionBidIDKey
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)

	store.Set(key, value)
}

func (k Keeper) GetLimitAuctionBidID(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = types.LimitAuctionBidIDKey
		value = store.Get(key)
	)

	if value == nil {
		return 0
	}

	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue()
}

func (k Keeper) SetUserLimitBidData(ctx sdk.Context, userLimitBidData types.LimitOrderBid, debtTokenID, collateralTokenID uint64, premium string) {
	var (
		store = k.Store(ctx)
		key   = types.UserLimitBidKey(debtTokenID, collateralTokenID, premium, userLimitBidData.BidderAddress)
		value = k.cdc.MustMarshal(&userLimitBidData)
	)

	store.Set(key, value)
}

func (k Keeper) GetUserLimitBidData(ctx sdk.Context, debtTokenID, collateralTokenID uint64, premium, address string) (mappingData types.LimitOrderBid, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.UserLimitBidKey(debtTokenID, collateralTokenID, premium, address)
		value = store.Get(key)
	)

	if value == nil {
		return mappingData, false
	}

	k.cdc.MustUnmarshal(value, &mappingData)
	return mappingData, true
}

func (k Keeper) DeleteUserLimitBidData(ctx sdk.Context, debtTokenID, collateralTokenID uint64, premium, address string) {
	var (
		store = k.Store(ctx)
		key   = types.UserLimitBidKey(debtTokenID, collateralTokenID, premium, address)
	)

	store.Delete(key)
}

func (k Keeper) GetUserLimitBidDataByPremium(ctx sdk.Context, debtTokenID, collateralTokenID uint64, premium string) (biddingData []types.LimitOrderBid, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.UserLimitBidKeyForPremium(debtTokenID, collateralTokenID, premium)
		iter  = sdk.KVStorePrefixIterator(store, key)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var mapData types.LimitOrderBid
		k.cdc.MustUnmarshal(iter.Value(), &mapData)
		biddingData = append(biddingData, mapData)
	}
	if biddingData == nil {
		return nil, false
	}

	return biddingData, true
}

func (k Keeper) DepositLimitAuctionBid(ctx sdk.Context, bidder string, CollateralTokenId, DebtTokenId uint64, PremiumDiscount string, amount sdk.Coin) error {
	id := k.GetLimitAuctionBidID(ctx)
	bidderAddr, err := sdk.AccAddressFromBech32(bidder)
	if err != nil {
		return nil
	}
	premiumDiscount, err := sdk.NewDecFromStr(PremiumDiscount)
	if err != nil {
		return err
	}
	collateralAsset, found := k.asset.GetAsset(ctx, CollateralTokenId)
	if !found {
		return assettypes.ErrorAssetDoesNotExist
	}
	collateralAssetToken := sdk.NewCoin(collateralAsset.Denom, sdk.NewInt(0))
	userLimitBid, found := k.GetUserLimitBidData(ctx, DebtTokenId, CollateralTokenId, PremiumDiscount, bidder)
	if !found {
		userLimitBid = types.LimitOrderBid{
			LimitOrderBiddingId: id + 1,
			BidderAddress:       bidder,
			CollateralToken:     collateralAssetToken, // zero
			DebtToken:           amount,               // user's balance
			BiddingId:           nil,
			PremiumDiscount:     premiumDiscount,
		}
	} else {
		userLimitBid.DebtToken = userLimitBid.DebtToken.Add(amount)
	}

	// send tokens from user to the auction module
	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, bidderAddr, types.ModuleName, sdk.NewCoins(amount))
	if err != nil {
		return err
	}

	// Set ID and LimitBid Data
	k.SetLimitAuctionBidID(ctx, userLimitBid.LimitOrderBiddingId)
	k.SetUserLimitBidData(ctx, userLimitBid, DebtTokenId, CollateralTokenId, PremiumDiscount)
	return nil
}

func (k Keeper) CancelLimitAuctionBid(ctx sdk.Context, bidder string, DebtTokenId, CollateralTokenId uint64, PremiumDiscount string) error {
	userLimitBid, found := k.GetUserLimitBidData(ctx, DebtTokenId, CollateralTokenId, PremiumDiscount, bidder)
	if !found {
		// return err not found
	}

	bidderAddr, err := sdk.AccAddressFromBech32(bidder)
	if err != nil {
		return err
	}
	// return all the tokens back to the user
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, bidderAddr, sdk.NewCoins(userLimitBid.DebtToken))
	if err != nil {
		return err
	}

	// delete userLimitBid from KV store
	k.DeleteUserLimitBidData(ctx, DebtTokenId, CollateralTokenId, PremiumDiscount, bidder)

	return nil
}

func (k Keeper) WithdrawLimitAuctionBid(ctx sdk.Context, bidder string, CollateralTokenId, DebtTokenId uint64, PremiumDiscount string, amount sdk.Coin) error {
	userLimitBid, found := k.GetUserLimitBidData(ctx, DebtTokenId, CollateralTokenId, PremiumDiscount, bidder)
	if !found {
		// return err not found
	}

	bidderAddr, err := sdk.AccAddressFromBech32(bidder)
	if err != nil {
		return err
	}

	if amount.Amount.Equal(userLimitBid.DebtToken.Amount) {
		err := k.CancelLimitAuctionBid(ctx, bidder, DebtTokenId, CollateralTokenId, PremiumDiscount)
		if err != nil {
			return err
		}
		return nil
	}
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, bidderAddr, sdk.NewCoins(amount))
	if err != nil {
		return err
	}
	userLimitBid.DebtToken.Amount = userLimitBid.DebtToken.Amount.Sub(amount.Amount)
	k.SetUserLimitBidData(ctx, userLimitBid, DebtTokenId, CollateralTokenId, PremiumDiscount)
	return nil
}
