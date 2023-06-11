package keeper

import (
	assettypes "github.com/comdex-official/comdex/x/asset/types"
	"github.com/comdex-official/comdex/x/auctionsV2/types"
	auctionsV2types "github.com/comdex-official/comdex/x/auctionsV2/types"
	auctiontypes "github.com/comdex-official/comdex/x/auctionsV2/types"
	collectortypes "github.com/comdex-official/comdex/x/collector/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	protobuftypes "github.com/gogo/protobuf/types"
)

func (k Keeper) PlaceDutchAuctionBid(ctx sdk.Context, auctionID uint64, bidder string, bid sdk.Coin, auctionData types.Auction, isAutoBid bool) (bidId uint64, err error) {
	auctionParams, _ := k.GetAuctionParams(ctx)
	if bid.Amount.Equal(sdk.ZeroInt()) {
		return bidId, types.ErrBidCannotBeZero
	}

	bidderAddr, _ := sdk.AccAddressFromBech32(bidder)

	liquidationWhitelistingAppData, _ := k.LiquidationsV2.GetLiquidationWhiteListing(ctx, auctionData.AppId)

	if bid.Denom != auctionData.DebtToken.Denom {
		return bidId, sdkerrors.Wrapf(sdkerrors.ErrNotFound, "Bid token is not the debt token ", bid.Denom)
	}
	liquidationData, _ := k.LiquidationsV2.GetLockedVault(ctx, auctionData.AppId, auctionData.LockedVaultId)
	//Price data of the token from market module
	debtToken, _ := k.market.GetTwa(ctx, auctionData.DebtAssetId)
	debtPrice := sdk.NewDecFromInt(sdk.NewInt(int64(debtToken.Twa)))
	//only if debt token is CMST , we consider it as $1
	if liquidationData.IsDebtCmst {
		debtPrice = sdk.NewDecFromInt(sdk.NewInt(int64(1000000)))
	}
	//Check to update bid.Amount
	fullBid := false

	if bid.Amount.GTE(auctionData.DebtToken.Amount) {
		bid.Amount = auctionData.DebtToken.Amount
		fullBid = true

	}
	_, collateralTokenQuanitity, _ := k.vault.GetAmountOfOtherToken(ctx, auctionData.DebtAssetId, debtPrice, bid.Amount, auctionData.CollateralAssetId, auctionData.CollateralTokenAuctionPrice)
	//From auction bonus quantity , use the available quantity to calculate the collateral value
	_, collateralTokenQuanitityForBonus, _ := k.vault.GetAmountOfOtherToken(ctx, auctionData.DebtAssetId, debtPrice, auctionData.BonusAmount, auctionData.CollateralAssetId, auctionData.CollateralTokenAuctionPrice)
	//Checking if the auction bonus and the collateral to be given to user isnt more than available colalteral
	totalCollateralTokenQuanitity := collateralTokenQuanitity.Add(collateralTokenQuanitityForBonus)
	//If user has sent a bigger bid than the target amount ,
	if fullBid || !totalCollateralTokenQuanitity.LTE(auctionData.CollateralToken.Amount) {

		if !totalCollateralTokenQuanitity.LTE(auctionData.CollateralToken.Amount) {
			//This means that there is less collateral available .
			leftOverCollateral := auctionData.CollateralToken.Amount
			_, debtTokenAgainstLeftOverCollateral, _ := k.vault.GetAmountOfOtherToken(ctx, auctionData.CollateralAssetId, auctionData.CollateralTokenAuctionPrice, leftOverCollateral.Sub(collateralTokenQuanitityForBonus), auctionData.DebtAssetId, debtPrice)
			bid.Amount = debtTokenAgainstLeftOverCollateral
			totalCollateralTokenQuanitity = leftOverCollateral
			//Amount to call from reserve account for adjusting the auction target debt
			//So we call the module account to give funds to compensate the user.
			debtGettingLeft := auctionData.DebtToken.Sub(sdk.NewCoin(auctionData.DebtToken.Denom, debtTokenAgainstLeftOverCollateral))
			//Calling reserve account for debt adjustment : debtGettingLeft
			//Updating the protocol was in loss stuct
			err := k.LiquidationsV2.WithdrawAppReserveFundsFn(ctx, auctionData.AppId, auctionData.DebtAssetId, debtGettingLeft)
			if err != nil {
				return bidId, err
			}
		}
		//Take Debt Token from user ,

		if !isAutoBid {
			if bid.Amount.GT(sdk.ZeroInt()) {
				err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, bidderAddr, auctionsV2types.ModuleName, sdk.NewCoins(sdk.NewCoin(auctionData.DebtToken.Denom, bid.Amount)))
				if err != nil {
					return bidId, err
				}
			}

		}

		//Send Collateral To bidder
		if totalCollateralTokenQuanitity.GT(sdk.ZeroInt()) {
			err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, auctionsV2types.ModuleName, bidderAddr, sdk.NewCoins(sdk.NewCoin(auctionData.CollateralToken.Denom, totalCollateralTokenQuanitity)))
			if err != nil {
				return bidId, err
			}
		}
		//Burn Debt Token,
		liquidationPenalty := sdk.NewCoin(auctionData.DebtToken.Denom, liquidationData.FeeToBeCollected)
		var tokensToBurn sdk.Coin
		if liquidationData.InitiatorType != "external" {
			tokensToBurn = liquidationData.TargetDebt.Sub(liquidationPenalty)
			if tokensToBurn.Amount.GT(sdk.ZeroInt()) {
				err := k.bankKeeper.BurnCoins(ctx, auctionsV2types.ModuleName, sdk.NewCoins(tokensToBurn))
				if err != nil {
					return bidId, err
				}
			}

		}

		//Send rest tokens to the user
		OwnerLeftOverCapital := auctionData.CollateralToken.Amount.Sub(totalCollateralTokenQuanitity)
		if OwnerLeftOverCapital.GT(sdk.ZeroInt()) {
			err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, auctionsV2types.ModuleName, sdk.AccAddress(liquidationData.Owner), sdk.NewCoins(sdk.NewCoin(auctionData.CollateralToken.Denom, OwnerLeftOverCapital)))
			if err != nil {
				return bidId, err
			}
		}
		//Add bid data to struct
		//Creating user bid struct
		biddingId, err := k.CreateUserBid(ctx, auctionData.AppId, bidder, auctionID, sdk.NewCoin(auctionData.CollateralToken.Denom, totalCollateralTokenQuanitity), sdk.NewCoin(auctionData.DebtToken.Denom, bid.Amount), "dutch")
		if err != nil {
			return bidId, err
		}
		//Based on app type call perform specific function - external , internal and /or keeper incentive
		//See if this was keeper initiated transaction- then incentivisation will be in place based on the percentage
		//For apps that are external to comdex chain

		if liquidationData.InitiatorType == "external" {
			//Send Liquidation penalty to the comdex protocol  --  create a kv store like SetAuctionLimitBidFeeData with name SetAuctionExternalFeeData
			//Send debt to the initiator address of the auction
			finalDebtToInitiator := liquidationData.DebtToken.Sub(liquidationPenalty)
			keeperIncentive := (liquidationWhitelistingAppData.KeeeperIncentive.Mul(sdk.NewDecFromInt(liquidationPenalty.Amount))).TruncateInt()
			if keeperIncentive.GT(sdk.ZeroInt()) {
				liquidationPenalty = liquidationPenalty.Sub(sdk.NewCoin(auctionData.DebtToken.Denom, keeperIncentive))
				err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, auctionsV2types.ModuleName, sdk.AccAddress(liquidationData.InternalKeeperAddress), sdk.NewCoins(sdk.NewCoin(auctionData.DebtToken.Denom, keeperIncentive)))
				if err != nil {
					return bidId, err
				}
			}

			// updating feeData for external initiator
			feeData, found := k.GetAuctionLimitBidFeeDataExternal(ctx, auctionData.DebtAssetId)
			if !found {
				feeData.AssetId = auctionData.DebtAssetId
				feeData.Amount = liquidationPenalty.Amount
			} else {
				feeData.Amount = feeData.Amount.Add(liquidationPenalty.Amount)
			}
			err = k.SetAuctionLimitBidFeeDataExternal(ctx, feeData)
			if err != nil {
				return 0, err
			}

			// sending collected debt to the initiator
			externalInitiator, err := sdk.AccAddressFromBech32(liquidationData.ExternalKeeperAddress)
			if err != nil {
				return 0, err
			}

			err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, auctionsV2types.ModuleName, externalInitiator, sdk.NewCoins(finalDebtToInitiator))
			if err != nil {
				return 0, err
			}

			//but if an app is external - will have to check the auction bonus , liquidation penalty , module account mechanism
		} else if liquidationData.InitiatorType == "vault" {
			//Check if they are initiated through a keeper, if so they will be incentivised
			if liquidationData.IsInternalKeeper {

				keeperIncentive := (liquidationWhitelistingAppData.KeeeperIncentive.Mul(sdk.NewDecFromInt(liquidationPenalty.Amount))).TruncateInt()
				if keeperIncentive.GT(sdk.ZeroInt()) {
					liquidationPenalty = liquidationPenalty.Sub(sdk.NewCoin(auctionData.DebtToken.Denom, keeperIncentive))
					err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, auctionsV2types.ModuleName, sdk.AccAddress(liquidationData.InternalKeeperAddress), sdk.NewCoins(sdk.NewCoin(auctionData.DebtToken.Denom, keeperIncentive)))
					if err != nil {
						return bidId, err
					}
				}
			}
			//Send Liquidation Penalty to the Collector Module
			if liquidationPenalty.Amount.GT(sdk.ZeroInt()) {
				err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, auctionsV2types.ModuleName, collectortypes.ModuleName, sdk.NewCoins(liquidationPenalty))
				if err != nil {
					return bidId, err
				}
			}
			//Update Collector Data for CMST
			// Updating fees data in collector
			err = k.collector.SetNetFeeCollectedData(ctx, auctionData.AppId, auctionData.CollateralAssetId, liquidationPenalty.Amount)
			if err != nil {
				return bidId, err
			}
			//Updating mapping data of vault
			k.vault.UpdateTokenMintedAmountLockerMapping(ctx, auctionData.AppId, liquidationData.ExtendedPairId, tokensToBurn.Amount, false)
			k.vault.UpdateCollateralLockedAmountLockerMapping(ctx, auctionData.AppId, liquidationData.ExtendedPairId, liquidationData.CollateralToken.Amount, false)
		} else if liquidationData.InitiatorType == "borrow" {
			//Check if they are initiated through a keeper, if so they will be incentivised
			//TODO:
			// send money back to the debt pool (assetOut pool)
			// liquidation penalty to the reserve and interest to the pool
			// send token to the bidder
			// if cross pool borrow, settle the transit asset to it's native pool
			// close the borrow and update the stats
			err = k.LiquidationsV2.MsgCloseDutchAuctionForBorrow(ctx, liquidationData, auctionData)
			if err != nil {
				return 0, err
			}
		}
		//Add bidder data in auction
		bidOwnerMappingData := auctionsV2types.BidOwnerMapping{BidId: biddingId, BidOwner: bidder}
		auctionData.BiddingIds = append(auctionData.BiddingIds, &bidOwnerMappingData)
		//Saving auction data to auction historical
		auctionHistoricalData := auctionsV2types.AuctionHistorical{AuctionId: auctionID, AuctionHistorical: &auctionData, LockedVault: &liquidationData}
		err = k.SetAuctionHistorical(ctx, auctionHistoricalData)
		if err != nil {
			return 0, err
		}
		//Close Auction
		err = k.DeleteAuction(ctx, auctionData)
		if err != nil {
			return 0, err
		}
		//Delete liquidation Data
		k.LiquidationsV2.DeleteLockedVault(ctx, auctionData.AppId, liquidationData.LockedVaultId)
		bidId = biddingId
	} else {
		//if bid amount is less than the target bid
		//Calculating collateral token value from bid(debt) token value
		_, collateralTokenQuantity, _ := k.vault.GetAmountOfOtherToken(ctx, auctionData.DebtAssetId, debtPrice, bid.Amount, auctionData.CollateralAssetId, auctionData.CollateralTokenAuctionPrice)
		debtLeft := bid.Amount.Sub(bid.Amount)
		debtuDollar, _ := k.CalcDollarValueForToken(ctx, auctionData.DebtAssetId, debtPrice, debtLeft)
		if !(debtuDollar).GT(sdk.NewDecFromInt(sdk.NewIntFromUint64(auctionParams.MinUsdValueLeft))) {
			return bidId, types.ErrCannotLeaveDebtLessThanDust
		}

		//From auction bonus quantity , use the available quantity to calculate the collateral value
		//Checking bid.Amount -> to targetbid ratio
		bidToTargetDebtRatio := (bid.Amount).Quo(auctionData.DebtToken.Amount)
		expectedBonusShareForCurrentBid := liquidationData.BonusToBeGiven.Mul(bidToTargetDebtRatio)
		//If somehow bonus to be given is less than what is there in the protocol
		if expectedBonusShareForCurrentBid.GT(auctionData.BonusAmount) {
			expectedBonusShareForCurrentBid = auctionData.BonusAmount
		}
		//using that ratio data to calculate  auction bonus to be given for the bid
		//first taking the debt percentage data
		//then calculating the collateral token data
		_, collateralTokenQuantityForBonus, _ := k.vault.GetAmountOfOtherToken(ctx, auctionData.DebtAssetId, debtPrice, expectedBonusShareForCurrentBid, auctionData.CollateralAssetId, auctionData.CollateralTokenAuctionPrice)
		totalCollateralTokenQuantity := collateralTokenQuantity.Add(collateralTokenQuantityForBonus)
		if !isAutoBid {
			if bid.Amount.GT(sdk.ZeroInt()) {
				err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, bidderAddr, auctionsV2types.ModuleName, sdk.NewCoins(sdk.NewCoin(auctionData.DebtToken.Denom, bid.Amount)))
				if err != nil {
					return bidId, err
				}
			}
		}
		//Send Collateral To bidder
		if totalCollateralTokenQuantity.GT(sdk.ZeroInt()) {
			err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, auctionsV2types.ModuleName, bidderAddr, sdk.NewCoins(sdk.NewCoin(auctionData.CollateralToken.Denom, totalCollateralTokenQuantity)))
			if err != nil {
				return bidId, err
			}
		}
		biddingId, err := k.CreateUserBid(ctx, auctionData.AppId, bidder, auctionID, sdk.NewCoin(auctionData.CollateralToken.Denom, totalCollateralTokenQuantity), sdk.NewCoin(auctionData.DebtToken.Denom, bid.Amount), "dutch")
		if err != nil {
			return bidId, err
		}
		//Add bidder data in auction
		bidOwnerMappingData := auctionsV2types.BidOwnerMapping{BidId: biddingId, BidOwner: string(bidder)}
		auctionData.BiddingIds = append(auctionData.BiddingIds, &bidOwnerMappingData)

		//Reduce Auction collateral and debt value
		auctionData.CollateralToken.Amount = auctionData.CollateralToken.Amount.Sub(totalCollateralTokenQuantity)
		auctionData.DebtToken.Amount = auctionData.DebtToken.Amount.Sub(bid.Amount)
		auctionData.BonusAmount = auctionData.BonusAmount.Sub(expectedBonusShareForCurrentBid)
		//Set Auction
		err = k.SetAuction(ctx, auctionData)
		if err != nil {
			return 0, err
		}
		bidId = biddingId
	}

	return bidId, nil
}

func (k Keeper) CreateUserBid(ctx sdk.Context, appID uint64, BidderAddress string, auctionID uint64, collateralToken sdk.Coin, debtToken sdk.Coin, bidType string) (bidding_id uint64, err error) {

	userBidId := k.GetUserBidID(ctx)
	bidding := auctionsV2types.Bid{
		BiddingId:             userBidId + 1,
		AuctionId:             auctionID,
		CollateralTokenAmount: collateralToken,
		DebtTokenAmount:       debtToken,
		BidderAddress:         BidderAddress,
		BiddingTimestamp:      ctx.BlockTime(),
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

func (k Keeper) PlaceEnglishAuctionBid(ctx sdk.Context, auctionID uint64, bidder string, bid sdk.Coin, auctionData types.Auction) error {
	auctionParams, _ := k.GetAuctionParams(ctx)
	if bid.Amount.Equal(sdk.ZeroInt()) {
		return types.ErrBidCannotBeZero
	}
	bidderAddr, _ := sdk.AccAddressFromBech32(bidder)

	liquidationData, found := k.LiquidationsV2.GetLockedVault(ctx, auctionData.AppId, auctionData.LockedVaultId)
	if !found {
		return auctiontypes.ErrLiquidationNotFound
	}
	//TokenLastBid is used to get the last bid on the auction from the user
	tokenLastBid := auctionData.DebtToken
	//this is used to save the current collateral data.
	tokenCollateralData := auctionData.CollateralToken
	//bidFrom user is used to know how many token do we need to collect form the user
	bidFromUser := bid
	if liquidationData.InitiatorType == "debt" {
		//In debt bid, the bid comes in form of collateral token , but gets converted interally for easy usecase
		tokenLastBid = auctionData.CollateralToken
		bidFromUser = auctionData.DebtToken
		tokenCollateralData = bid
	}
	if bid.Denom != tokenLastBid.Denom {
		return sdkerrors.Wrapf(sdkerrors.ErrNotFound, "Bid is not in correct denom ", bid.Denom)
	}
	if auctionData.BiddingIds != nil {

		change := auctionParams.BidFactor.MulInt(tokenLastBid.Amount).Ceil().TruncateInt()
		bidAmount := tokenLastBid.Amount.Add(change)
		if bid.Amount.LT(bidAmount) {
			return sdkerrors.Wrapf(sdkerrors.ErrNotFound, "bid should be less than or equal to %d ", bidAmount.Uint64())
		}
		if liquidationData.InitiatorType == "debt" {
			bidAmount = tokenLastBid.Amount.Sub(change)
			if bid.Amount.GT(bidAmount) {
				return sdkerrors.Wrapf(sdkerrors.ErrNotFound, "bid should be less than or equal to %d ", bidAmount.Uint64())
			}
		}
	} else {
		if liquidationData.InitiatorType != "debt" && bid.Amount.LT(tokenLastBid.Amount) {
			return auctiontypes.ErrorLowBidAmount
		}
		if liquidationData.InitiatorType == "debt" && bid.Amount.GT(tokenLastBid.Amount) {
			return auctiontypes.ErrorMaxBidAmount
		}
	}
	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, bidderAddr, auctiontypes.ModuleName, sdk.NewCoins(bidFromUser))
	if err != nil {
		return err
	}
	biddingId, err := k.CreateUserBid(ctx, auctionData.AppId, string(bidder), auctionID, tokenCollateralData, bidFromUser, "english")
	if err != nil {
		return err
	}
	if auctionData.ActiveBiddingId != 0 {
		userBid, err := k.GetUserBid(ctx, auctionData.ActiveBiddingId)
		if err != nil {
			return err
		}
		addr, _ := sdk.AccAddressFromBech32(userBid.BidderAddress)
		err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, auctiontypes.ModuleName, addr, sdk.NewCoins(auctionData.DebtToken))
		if err != nil {
			return err
		}
	}
	if liquidationData.InitiatorType == "debt" {
		auctionData.CollateralToken.Amount = bid.Amount
	} else {
		auctionData.DebtToken.Amount = bid.Amount
	}

	auctionData.ActiveBiddingId = biddingId
	bidIDOwner := &auctiontypes.BidOwnerMapping{BidId: biddingId, BidOwner: bidder}
	auctionData.BiddingIds = append(auctionData.BiddingIds, bidIDOwner)

	err = k.SetAuction(ctx, auctionData)
	if err != nil {
		return err
	}
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

func (k Keeper) SetUserLimitBidData(ctx sdk.Context, userLimitBidData types.LimitOrderBid, debtTokenID, collateralTokenID uint64, premium sdk.Int) {
	var (
		store = k.Store(ctx)
		key   = types.UserLimitBidKey(debtTokenID, collateralTokenID, premium, userLimitBidData.BidderAddress)
		value = k.cdc.MustMarshal(&userLimitBidData)
	)

	store.Set(key, value)
}

func (k Keeper) GetUserLimitBidData(ctx sdk.Context, debtTokenID, collateralTokenID uint64, premium sdk.Int, address string) (mappingData types.LimitOrderBid, found bool) {
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

func (k Keeper) DeleteUserLimitBidData(ctx sdk.Context, debtTokenID, collateralTokenID uint64, premium sdk.Int, address string) {
	var (
		store = k.Store(ctx)
		key   = types.UserLimitBidKey(debtTokenID, collateralTokenID, premium, address)
	)

	store.Delete(key)
}

func (k Keeper) GetUserLimitBidDataByPremium(ctx sdk.Context, debtTokenID, collateralTokenID uint64, premium sdk.Int) (biddingData []types.LimitOrderBid, found bool) {
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

func (k Keeper) DepositLimitAuctionBid(ctx sdk.Context, bidder string, CollateralTokenId, DebtTokenId uint64, PremiumDiscount sdk.Int, amount sdk.Coin) error {
	id := k.GetLimitAuctionBidID(ctx)
	bidderAddr, err := sdk.AccAddressFromBech32(bidder)
	if err != nil {
		return nil
	}

	_, found := k.asset.GetAsset(ctx, CollateralTokenId)
	if !found {
		return assettypes.ErrorAssetDoesNotExist
	}
	userLimitBid, found := k.GetUserLimitBidData(ctx, DebtTokenId, CollateralTokenId, PremiumDiscount, bidder)
	if !found {
		userLimitBid = types.LimitOrderBid{
			LimitOrderBiddingId: id + 1,
			BidderAddress:       bidder,
			DebtToken:           amount, // user's balance
			BiddingId:           nil,
			PremiumDiscount:     PremiumDiscount,
			CollateralTokenId:   CollateralTokenId,
			DebtTokenId:         DebtTokenId,
		}

		k.UpdateUserLimitBidDataForAddress(ctx, userLimitBid, true)
	} else {
		userLimitBid.DebtToken = userLimitBid.DebtToken.Add(amount)
	}

	// send tokens from user to the auction module
	if amount.Amount.GT(sdk.ZeroInt()) {
		err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, bidderAddr, types.ModuleName, sdk.NewCoins(amount))
		if err != nil {
			return err
		}
	}

	// Set ID and LimitBid Data
	k.SetLimitAuctionBidID(ctx, userLimitBid.LimitOrderBiddingId)
	k.SetUserLimitBidData(ctx, userLimitBid, DebtTokenId, CollateralTokenId, PremiumDiscount)

	return nil
}

func (k Keeper) CancelLimitAuctionBid(ctx sdk.Context, bidder string, DebtTokenId, CollateralTokenId uint64, PremiumDiscount sdk.Int) error {
	userLimitBid, found := k.GetUserLimitBidData(ctx, DebtTokenId, CollateralTokenId, PremiumDiscount, bidder)
	if !found {
		// return err not found
		return types.ErrBidNotFound
	}
	auctionParams, _ := k.GetAuctionParams(ctx)

	bidderAddr, err := sdk.AccAddressFromBech32(bidder)
	if err != nil {
		return err
	}
	// return all the tokens back to the user
	if userLimitBid.DebtToken.Amount.GT(sdk.ZeroInt()) {
		feesToBecollected := auctionParams.ClosingFee.Mul(sdk.NewDecFromInt(userLimitBid.DebtToken.Amount)).TruncateInt()
		userLimitBid.DebtToken.Amount = userLimitBid.DebtToken.Amount.Sub(feesToBecollected)

		err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, bidderAddr, sdk.NewCoins(userLimitBid.DebtToken))
		if err != nil {
			return err
		}
		//updating fees in auction data
		feeData, found := k.GetAuctionLimitBidFeeData(ctx, DebtTokenId)
		if !found {
			var feeData types.AuctionFeesCollectionFromLimitBidTx
			feeData.AssetId = DebtTokenId
			feeData.Amount = feesToBecollected
		} else {
			feeData.Amount = feeData.Amount.Add(feesToBecollected)
		}
		err := k.SetAuctionLimitBidFeeData(ctx, feeData)
		if err != nil {
			return err
		}
	}

	// delete userLimitBid from KV store
	k.UpdateUserLimitBidDataForAddress(ctx, userLimitBid, false)
	k.DeleteUserLimitBidData(ctx, DebtTokenId, CollateralTokenId, PremiumDiscount, bidder)

	return nil
}

func (k Keeper) WithdrawLimitAuctionBid(ctx sdk.Context, bidder string, CollateralTokenId, DebtTokenId uint64, PremiumDiscount sdk.Int, amount sdk.Coin) error {
	userLimitBid, found := k.GetUserLimitBidData(ctx, DebtTokenId, CollateralTokenId, PremiumDiscount, bidder)
	if !found {
		// return err not found
	}

	bidderAddr, err := sdk.AccAddressFromBech32(bidder)
	if err != nil {
		return err
	}
	auctionParams, _ := k.GetAuctionParams(ctx)

	if amount.Amount.Equal(userLimitBid.DebtToken.Amount) {
		err := k.CancelLimitAuctionBid(ctx, bidder, DebtTokenId, CollateralTokenId, PremiumDiscount)
		if err != nil {
			return err
		}
		return nil
	}

	// return all the tokens back to the user
	if userLimitBid.DebtToken.Amount.GT(sdk.ZeroInt()) {
		feesToBecollected := auctionParams.WithdrawalFee.Mul(sdk.NewDecFromInt(amount.Amount)).TruncateInt()
		userLimitBid.DebtToken.Amount = userLimitBid.DebtToken.Amount.Sub(feesToBecollected)

		err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, bidderAddr, sdk.NewCoins(amount))
		if err != nil {
			return err
		}
		//updating fees in auction data
		feeData, found := k.GetAuctionLimitBidFeeData(ctx, DebtTokenId)
		if !found {
			var feeData types.AuctionFeesCollectionFromLimitBidTx
			feeData.AssetId = DebtTokenId
			feeData.Amount = feesToBecollected
		} else {
			feeData.Amount = feeData.Amount.Add(feesToBecollected)
		}
		err := k.SetAuctionLimitBidFeeData(ctx, feeData)
		if err != nil {
			return err
		}
	}

	userLimitBid.DebtToken.Amount = userLimitBid.DebtToken.Amount.Sub(amount.Amount)
	k.SetUserLimitBidData(ctx, userLimitBid, DebtTokenId, CollateralTokenId, PremiumDiscount)
	return nil
}

func (k Keeper) CalcDollarValueForToken(ctx sdk.Context, id uint64, rate sdk.Dec, amt sdk.Int) (price sdk.Dec, err error) {
	asset, _ := k.asset.GetAsset(ctx, id)
	numerator := sdk.NewDecFromInt(amt).Mul(rate)
	denominator := sdk.NewDecFromInt(asset.Decimals)
	return numerator.Quo(denominator), nil
}

func (k Keeper) SetUserLimitBidDataForAddress(ctx sdk.Context, userLimitBidData types.LimitOrderBidsForUser) {
	var (
		store = k.Store(ctx)
		key   = types.UserLimitBidKeyForAddress(userLimitBidData.BidderAddress)
		value = k.cdc.MustMarshal(&userLimitBidData)
	)

	store.Set(key, value)
}

func (k Keeper) UpdateUserLimitBidDataForAddress(ctx sdk.Context, userLimitBid types.LimitOrderBid, changeType bool) {
	userLimitBidForAddress, found := k.GetUserLimitBidDataByAddress(ctx, userLimitBid.BidderAddress)

	if changeType {
		if !found {
			userLimitBidForAddress.BidderAddress = userLimitBid.BidderAddress
			var userLimitBidSecondaryKey types.LimitOrderUserKey
			userLimitBidSecondaryKey.CollateralTokenId = userLimitBid.CollateralTokenId
			userLimitBidSecondaryKey.DebtTokenId = userLimitBid.DebtTokenId
			userLimitBidSecondaryKey.PremiumDiscount = userLimitBid.PremiumDiscount
			userLimitBidSecondaryKey.LimitOrderBiddingId = userLimitBid.LimitOrderBiddingId
			userLimitBidForAddress.LimitOrderBidKey = append(userLimitBidForAddress.LimitOrderBidKey, userLimitBidSecondaryKey)
		}
	} else {
		if !found {
			userLimitBidForAddress.BidderAddress = userLimitBid.BidderAddress
		} else {
			for index, individualLimitOrderBid := range userLimitBidForAddress.LimitOrderBidKey {
				if individualLimitOrderBid.LimitOrderBiddingId == userLimitBid.LimitOrderBiddingId {
					userLimitBidForAddress.LimitOrderBidKey = append(userLimitBidForAddress.LimitOrderBidKey[:index], userLimitBidForAddress.LimitOrderBidKey[index+1:]...)
				}
			}
		}
	}

	k.SetUserLimitBidDataForAddress(ctx, userLimitBidForAddress)
}

func (k Keeper) GetUserLimitBidDataByAddress(ctx sdk.Context, address string) (mappingData types.LimitOrderBidsForUser, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.UserLimitBidKeyForAddress(address)
		value = store.Get(key)
	)

	if value == nil {
		return mappingData, false
	}

	k.cdc.MustUnmarshal(value, &mappingData)
	return mappingData, true
}
