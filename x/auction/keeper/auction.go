package keeper

import (
	"fmt"
	"time"

	vaulttypes "github.com/comdex-official/comdex/x/vault/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	tokenminttypes "github.com/comdex-official/comdex/x/tokenmint/types"

	assettypes "github.com/comdex-official/comdex/x/asset/types"
	auctiontypes "github.com/comdex-official/comdex/x/auction/types"
	collectortypes "github.com/comdex-official/comdex/x/collector/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) IncreaseLockedVaultAmountIn(ctx sdk.Context, lockedVaultId uint64, amount sdk.Int) error {
	lockedVault, found := k.GetLockedVault(ctx, lockedVaultId)
	if !found {
		return auctiontypes.ErrorVaultNotFound
	}
	lockedVault.AmountIn = lockedVault.AmountIn.Add(amount)
	k.SetLockedVault(ctx, lockedVault)
	return nil
}

func (k Keeper) DecreaseLockedVaultAmountIn(ctx sdk.Context, lockedVaultId uint64, amount sdk.Int) error {
	lockedVault, found := k.GetLockedVault(ctx, lockedVaultId)
	if !found {
		return auctiontypes.ErrorVaultNotFound
	}
	lockedVault.AmountIn = lockedVault.AmountIn.Sub(amount)
	k.SetLockedVault(ctx, lockedVault)
	return nil
}

func (k Keeper) DecreaseLockedVaultAmountOut(ctx sdk.Context, lockedVaultId uint64, amount sdk.Int) error {
	lockedVault, found := k.GetLockedVault(ctx, lockedVaultId)
	if !found {
		return auctiontypes.ErrorVaultNotFound
	}
	lockedVault.AmountIn = lockedVault.AmountOut.Sub(amount)
	k.SetLockedVault(ctx, lockedVault)
	return nil
}

func (k Keeper) getInflowTokenAmount(ctx sdk.Context, appId, AssetInId, AssetOutId uint64, lotSize sdk.Int) (status uint64, outflowToken, inflowToken sdk.Coin) {
	emptyCoin := sdk.NewCoin("empty", sdk.NewIntFromUint64(1))
	outflowAsset, found1 := k.GetAsset(ctx, AssetOutId)
	inflowAsset, found2 := k.GetAsset(ctx, AssetInId)
	if !found1 || !found2 {
		return auctiontypes.NoAuction, emptyCoin, emptyCoin
	}

	var inFlowTokenPrice uint64
	collectorAuction, _ := k.GetAuctionMappingForApp(ctx, appId)
	for _, data := range collectorAuction.AssetIdToAuctionLookup {

		if data.AssetOutOraclePrice {
			fmt.Println(data.AssetOutOraclePrice, "value bool price required")
			//If oracle Price required for the assetOut
			inFlowTokenPrice, _ = k.GetPriceForAsset(ctx, AssetInId)

			fmt.Println(inFlowTokenPrice, "should be what is set dollar ")
		} else {
			//If oracle Price is not required for the assetOut
			inFlowTokenPrice = data.AssetOutPrice

		}

	}

	outflowToken = sdk.NewCoin(outflowAsset.Denom, lotSize)
	outflowTokenPrice, found3 := k.GetPriceForAsset(ctx, outflowAsset.Id)
	// inflowTokenPrice, found4 := k.GetPriceForAsset(ctx, inflowAsset.Id)
	if !found3 {
		return auctiontypes.NoAuction, emptyCoin, emptyCoin
	}
	inflowTokenAmount := outflowToken.Amount.Mul(sdk.NewIntFromUint64(outflowTokenPrice)).Quo(sdk.NewIntFromUint64(inFlowTokenPrice))
	inflowToken = sdk.NewCoin(inflowAsset.Denom, inflowTokenAmount)
	return 5, outflowToken, inflowToken
}

func (k Keeper) checkStatusOfNetFeesCollectedAndStartAuction(ctx sdk.Context, appId, assetId uint64, assetToAuction collectortypes.AssetIdToAuctionLookupTable) (status uint64, err error) {
	assetsCollectorDataUnderAppId, found := k.GetCollectorLookupTable(ctx, appId)
	if !found {
		return
	}
	//traverse this to access appId , collector asset id , surplus threshhold , debt threshhold
	for _, collector := range assetsCollectorDataUnderAppId.AssetRateInfo {
		if collector.CollectorAssetId == assetId {
			//collectorLookupTable has surplusThreshhold for all assets

			NetFeeCollectedData, found := k.GetNetFeeCollectedData(ctx, appId)
			if !found {
				return auctiontypes.NoAuction, nil
			}
			//traverse this to access appId , collector asset id , netfees collected
			for _, AssetIdToFeeCollected := range NetFeeCollectedData.AssetIdToFeeCollected {
				if AssetIdToFeeCollected.AssetId == assetId {
					// if netfees <= debt threshhold -lotsize the start debt auction with lot size and debt auction is allowed true
					if AssetIdToFeeCollected.NetFeesCollected.LTE(sdk.NewIntFromUint64(collector.DebtThreshold-collector.LotSize)) && assetToAuction.IsDebtAuction {
						// START DEBT AUCTION .  LOTSIZE AS MINTED FOR SECONDARY ASSET and ACCEPT Collector assetid from user
						//calculate inflow token amount
						assetInId := collector.CollectorAssetId
						assetOutId := collector.SecondaryAssetId
						status, inflowToken, outflowToken := k.getInflowTokenAmount(ctx, appId, assetInId, assetOutId, sdk.NewIntFromUint64(collector.LotSize))
						if status == auctiontypes.NoAuction {
							return auctiontypes.NoAuction, nil
						}
						//Mint the tokens when collector module sends tokens to user
						err := k.StartDebtAuction(ctx, outflowToken, inflowToken, collector.BidFactor, appId, assetId, assetInId, assetOutId)
						if err != nil {
							break
						}
						return auctiontypes.StartedDebtAuction, nil
						// if netfees >= surplus threshhold+lotsize the start surplus auction with lot size and surplus auction is allowed true
					} else if AssetIdToFeeCollected.NetFeesCollected.GTE(sdk.NewIntFromUint64(collector.SurplusThreshold+collector.LotSize)) && assetToAuction.IsDebtAuction {
						// START SURPLUS AUCTION .  WITH COLLECTOR ASSET ID AS token given to user of lot size and secondary asset as received from user and burnt , bid factor
						//calculate inflow token amount
						assetInId := collector.SecondaryAssetId
						assetOutId := collector.CollectorAssetId
						status, inflowToken, outflowToken := k.getInflowTokenAmount(ctx, appId, assetInId, assetOutId, sdk.NewIntFromUint64(collector.LotSize))
						if status == auctiontypes.NoAuction {
							return auctiontypes.NoAuction, nil
						}
						//Transfer balance from collector module to auction module
						_, err := k.GetAmountFromCollector(ctx, appId, assetId, outflowToken.Amount)
						if err != nil {
							return status, err
						}
						err = k.StartSurplusAuction(ctx, outflowToken, inflowToken, collector.BidFactor, appId, assetId, assetInId, assetOutId)
						if err != nil {
							return status, err
						}
						return auctiontypes.StartedSurplusAuction, nil
					} else {
						return auctiontypes.NoAuction, nil
					}
				}
			}
		}
	}
	return auctiontypes.NoAuction, nil
}

func (k Keeper) CreateSurplusAndDebtAuctions(ctx sdk.Context) error {
	appIds, found := k.GetApps(ctx)
	if !found {
		return assettypes.AppIdsDoesntExist
	}
	for _, appId := range appIds {
		//check if auction status for an asset is false
		auctionLookupTable, found := k.GetCollectorAuctionLookupTable(ctx, appId.Id)
		if !found {
			continue
		}
		for _, assetToAuction := range auctionLookupTable.AssetIdToAuctionLookup {
			if assetToAuction.IsSurplusAuction || assetToAuction.IsDebtAuction {
				if !assetToAuction.IsAuctionActive {
					status, err := k.checkStatusOfNetFeesCollectedAndStartAuction(ctx, appId.Id, assetToAuction.AssetId, assetToAuction)
					if err != nil {
						return err
					}
					if status == auctiontypes.StartedDebtAuction {
						assetToAuction.IsAuctionActive = true
					} else if status == auctiontypes.StartedSurplusAuction {
						assetToAuction.IsAuctionActive = true
					} else {
						continue
					}
					err = k.SetCollectorAuctionLookupTable(ctx, auctionLookupTable)
					if err == nil {
						continue
					}
				}
			}
		}
	}
	return nil
}

func (k Keeper) makeFalseForFlags(ctx sdk.Context, appId, assetId uint64) error {

	auctionLookupTable, found := k.GetCollectorAuctionLookupTable(ctx, appId)
	if !found {
		return auctiontypes.ErrorInvalidAddress
	}
	for _, assetToAuction := range auctionLookupTable.AssetIdToAuctionLookup {
		if assetToAuction.AssetId == assetId {
			assetToAuction.IsAuctionActive = false
			err := k.SetCollectorAuctionLookupTable(ctx, auctionLookupTable)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (k Keeper) CloseAndRestartAuctions(ctx sdk.Context) error {
	appIds, found := k.GetApps(ctx)
	if !found {
		return assettypes.AppIdsDoesntExist
	}
	for _, appId := range appIds {
		err := k.CloseSurplusAuctions(ctx, appId.Id)
		if err != nil {
			return err
		}
		err = k.CloseDebtAuctions(ctx, appId.Id)
		if err != nil {
			return err
		}
		err = k.RestartDutchAuctions(ctx, appId.Id)
		if err != nil {
			return err
		}
	}
	err := k.CreateSurplusAndDebtAuctions(ctx)
	if err != nil {
		return err
	}
	err = k.CreateNewDutchAuctions(ctx)
	if err != nil {
		return err
	}
	return nil
}
func (k Keeper) CreateNewDutchAuctions(ctx sdk.Context) error {
	fmt.Println("CreateNewDutchAuctions")
	lockedVaults := k.GetLockedVaults(ctx)
	fmt.Println("CreateNewDutchAuctions", lockedVaults)
	fmt.Println(len(lockedVaults))
	if len(lockedVaults) == 0 {
		return auctiontypes.ErrorInvalidLockedVault
	}
	for _, lockedVault := range lockedVaults {
		fmt.Println("in for loop", lockedVault)
		extendedPair, found := k.GetPairsVault(ctx, lockedVault.ExtendedPairId)
		if !found {
			return auctiontypes.ErrorInvalidPair
		}
		pair, found := k.GetPair(ctx, extendedPair.PairId)
		if !found {
			return auctiontypes.ErrorInvalidPair
		}
		assetIn, found := k.GetAsset(ctx, pair.AssetIn)
		if !found {
			return auctiontypes.ErrorAssetNotFound
		}

		assetOut, found := k.GetAsset(ctx, pair.AssetOut)
		if !found {
			return auctiontypes.ErrorAssetNotFound
		}
		assetInPrice, found := k.GetPriceForAsset(ctx, assetIn.Id)
		if !found {
			return auctiontypes.ErrorPrices
		}
		//assetInPrice is the collateral price
		outflowToken := sdk.NewCoin(assetIn.Denom, lockedVault.CollateralToBeAuctioned.Quo(sdk.NewDecFromInt(sdk.NewIntFromUint64(assetInPrice))).TruncateInt())
		inflowToken := sdk.NewCoin(assetOut.Denom, sdk.ZeroInt())
		fmt.Println("check err")

		extendedPairId := lockedVault.ExtendedPairId
		ExtendedPairVault, found := k.GetPairsVault(ctx, extendedPairId)
		if !found {
			return auctiontypes.ErrorInvalidExtendedPairVault
		}
		liquidationPenalty := ExtendedPairVault.LiquidationPenalty
		fmt.Println("Before CreateNewDutchAuctions !lockedVault.IsAuctionInProgress", lockedVault.IsAuctionInProgress)
		if !lockedVault.IsAuctionInProgress {
			err1 := k.StartDutchAuction(ctx, outflowToken, inflowToken, lockedVault.AppMappingId, assetOut.Id, assetIn.Id, lockedVault.LockedVaultId, lockedVault.Owner, liquidationPenalty)
			if err1 != nil {
				return err1
			}
		}

		// fetch liquidation penalty
		//1.fetch extended pair vault id
		//2.query in asset

	}
	return nil
}

func (k Keeper) CloseSurplusAuctions(ctx sdk.Context, appId uint64) error {
	surplusAuctions := k.GetSurplusAuctions(ctx, appId)
	for _, surplusAuction := range surplusAuctions {
		if ctx.BlockTime().After(surplusAuction.EndTime) {
			if surplusAuction.AuctionStatus == auctiontypes.AuctionStartNoBids {
				err := k.RestartSurplusAuction(ctx, appId, surplusAuction)
				if err != nil {
					return err
				}
			} else {
				err := k.CloseSurplusAuction(ctx, surplusAuction)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// Get all app ids and call RestartDutchAuctions with app id
func (k Keeper) CloseDebtAuctions(ctx sdk.Context, appId uint64) error {
	debtAuctions := k.GetDebtAuctions(ctx, appId)
	for _, debtAuction := range debtAuctions {
		if ctx.BlockTime().After(debtAuction.EndTime) {
			if debtAuction.AuctionStatus == auctiontypes.AuctionStartNoBids {
				err := k.RestartDebtAuction(ctx, appId, debtAuction)
				if err != nil {
					return err
				}
			} else {
				err := k.CloseDebtAuction(ctx, debtAuction)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (k Keeper) RestartDebtAuction(
	ctx sdk.Context,
	appId uint64,
	debtAuction auctiontypes.DebtAuction,
) error {
	status, inflowToken, _ := k.getInflowTokenAmount(ctx, appId, debtAuction.AssetInId, debtAuction.AssetOutId, debtAuction.AuctionedToken.Amount)
	if status == auctiontypes.NoAuction {
		return nil
	}
	auctionParams := k.GetParams(ctx)
	debtAuction.ExpectedUserToken = inflowToken
	debtAuction.EndTime = ctx.BlockTime().Add(time.Second * time.Duration(auctionParams.AuctionDurationSeconds))
	err := k.SetDebtAuction(ctx, debtAuction)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) RestartSurplusAuction(
	ctx sdk.Context,
	appId uint64,
	surplusAuction auctiontypes.SurplusAuction,
) error {
	status, inflowToken, _ := k.getInflowTokenAmount(ctx, appId, surplusAuction.AssetInId, surplusAuction.AssetOutId, surplusAuction.OutflowToken.Amount)
	if status == auctiontypes.NoAuction {
		return nil
	}
	auctionParams := k.GetParams(ctx)
	surplusAuction.InflowToken = inflowToken
	surplusAuction.Bid = inflowToken
	surplusAuction.EndTime = ctx.BlockTime().Add(time.Second * time.Duration(auctionParams.AuctionDurationSeconds))
	err := k.SetSurplusAuction(ctx, surplusAuction)
	if err != nil {
		return err
	}
	return nil
}

//get all app ids and call RestartDutchAuctions with app id
func (k Keeper) RestartDutchAuctions(ctx sdk.Context, appId uint64) error {
	dutchAuctions := k.GetDutchAuctions(ctx, appId)
	auctionParams := k.GetParams(ctx)
	// SET current price of inflow token and outflow token
	for _, dutchAuction := range dutchAuctions {
		lockedVault, _ := k.GetLockedVault(ctx, dutchAuction.LockedVaultId)
		ExtendedPairVault, _ := k.GetPairsVault(ctx, lockedVault.ExtendedPairId)

		var inFlowTokenCurrentPrice uint64
		if ExtendedPairVault.AssetOutOraclePrice {
			fmt.Println(ExtendedPairVault.AssetOutOraclePrice, "value bool price required")
			//If oracle Price required for the assetOut
			inFlowTokenCurrentPrice, _ = k.GetPriceForAsset(ctx, dutchAuction.AssetInId)
	
			fmt.Println(inFlowTokenCurrentPrice, "should be what is set dollar ")
		} else {
			//If oracle Price is not required for the assetOut
			inFlowTokenCurrentPrice = ExtendedPairVault.AssetOutPrice
	
		}
		//inFlowTokenCurrentPrice := sdk.MustNewDecFromStr("1")
		dutchAuction.InflowTokenCurrentPrice = sdk.NewDec(int64(inFlowTokenCurrentPrice))
		tau := sdk.NewInt(int64(auctionParams.AuctionDurationSeconds))
		dur := ctx.BlockTime().Sub(dutchAuction.StartTime)
		seconds := sdk.NewInt(int64(dur.Seconds()))
		outFlowTokenCurrentPrice := k.getPriceFromLinearDecreaseFunction(dutchAuction.OutflowTokenInitialPrice, tau, seconds)
		dutchAuction.OutflowTokenCurrentPrice = outFlowTokenCurrentPrice

		//check if auction need to be restarted
		if ctx.BlockTime().After(dutchAuction.EndTime) || outFlowTokenCurrentPrice.LT(dutchAuction.OutflowTokenEndPrice) {
			//SET initial price fetched from market module and also end price , start time , end time
			//outFlowTokenCurrentPrice := sdk.NewIntFromUint64(10)
			outFlowTokenCurrentPrice, found := k.GetPriceForAsset(ctx, dutchAuction.AssetOutId)
			if !found {
				return auctiontypes.ErrorPrices
			}
			timeNow := ctx.BlockTime()
			dutchAuction.StartTime = timeNow
			dutchAuction.EndTime = timeNow.Add(time.Second * time.Duration(auctionParams.AuctionDurationSeconds))
			outFlowTokenInitialPrice := k.getOutflowTokenInitialPrice(sdk.NewIntFromUint64(outFlowTokenCurrentPrice), auctionParams.Buffer)
			outFlowTokenEndPrice := k.getOutflowTokenEndPrice(outFlowTokenInitialPrice, auctionParams.Cusp)
			dutchAuction.OutflowTokenInitialPrice = outFlowTokenInitialPrice
			dutchAuction.OutflowTokenEndPrice = outFlowTokenEndPrice
			dutchAuction.OutflowTokenCurrentPrice = outFlowTokenInitialPrice
		}
		err := k.SetDutchAuction(ctx, dutchAuction)
		if err != nil {
			return err
		}
	}
	return nil
}

func (k Keeper) StartSurplusAuction(
	ctx sdk.Context,
	outflowToken sdk.Coin,
	inflowToken sdk.Coin,
	bidFactor sdk.Dec,
	appId, assetId uint64,
	assetInId, assetOutId uint64,
) error {

	auctionParams := k.GetParams(ctx)
	auction := auctiontypes.SurplusAuction{
		OutflowToken:     outflowToken,
		InflowToken:      inflowToken,
		ActiveBiddingId:  0,
		Bidder:           nil,
		Bid:              inflowToken,
		EndTime:          ctx.BlockTime().Add(time.Second * time.Duration(auctionParams.AuctionDurationSeconds)),
		BidFactor:        bidFactor,
		BiddingIds:       []*auctiontypes.BidOwnerMapping{},
		AuctionStatus:    auctiontypes.AuctionStartNoBids,
		AppId:            appId,
		AssetId:          assetId,
		AuctionMappingId: auctionParams.SurplusId,
		AssetInId:        assetInId,
		AssetOutId:       assetOutId,
	}
	auction.AuctionId = k.GetAuctionID(ctx) + 1
	k.SetAuctionID(ctx, auction.AuctionId)
	err := k.SetSurplusAuction(ctx, auction)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) StartDebtAuction(
	ctx sdk.Context,
	auctionToken sdk.Coin,
	expectedUserToken sdk.Coin,
	bidFactor sdk.Dec,
	appId, assetId uint64,
	assetInId, assetOutId uint64,
) error {

	auctionParams := k.GetParams(ctx)
	auction := auctiontypes.DebtAuction{
		AuctionedToken:      auctionToken,
		ExpectedMintedToken: auctionToken,
		ExpectedUserToken:   expectedUserToken,
		ActiveBiddingId:     0,
		Bidder:              nil,
		EndTime:             ctx.BlockTime().Add(time.Second * time.Duration(auctionParams.AuctionDurationSeconds)),
		CurrentBidAmount:    sdk.NewCoin(auctionToken.Denom, sdk.NewInt(0)),
		AuctionStatus:       auctiontypes.AuctionStartNoBids,
		AppId:               appId,
		AssetId:             assetId,
		BiddingIds:          []*auctiontypes.BidOwnerMapping{},
		AuctionMappingId:    auctionParams.DebtId,
		BidFactor:           bidFactor,
		AssetInId:           assetInId,
		AssetOutId:          assetOutId,
	}
	auction.AuctionId = k.GetAuctionID(ctx) + 1
	k.SetAuctionID(ctx, auction.AuctionId)
	err := k.SetDebtAuction(ctx, auction)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) StartDutchAuction(
	ctx sdk.Context,
	outFlowToken sdk.Coin,
	inFlowToken sdk.Coin,
	appId uint64,
	assetInId, assetOutId uint64,
	lockedVaultId uint64,
	lockedVaultOwner string,
	liquidationPenalty sdk.Dec,
) error {
	var (
		inFlowTokenPrice  uint64
		outFlowTokenPrice uint64
		// found1            bool
		found2 bool
	)

	lockedVault, _ := k.GetLockedVault(ctx, lockedVaultId)

	var extendedPairVault = lockedVault.ExtendedPairId

	ExtendedPairVault, _ := k.GetPairsVault(ctx, extendedPairVault)

	if ExtendedPairVault.AssetOutOraclePrice {
		fmt.Println(ExtendedPairVault.AssetOutOraclePrice, "value bool price required")
		//If oracle Price required for the assetOut
		inFlowTokenPrice, _ = k.GetPriceForAsset(ctx, assetInId)

		fmt.Println(inFlowTokenPrice, "should be what is set dollar ")
	} else {
		//If oracle Price is not required for the assetOut
		inFlowTokenPrice = ExtendedPairVault.AssetOutPrice

	}

	err := k.SendCoinsFromModuleToModule(ctx, vaulttypes.ModuleName, auctiontypes.ModuleName, sdk.NewCoins(outFlowToken))
	if err != nil {
		fmt.Println("print err", err)
		fmt.Println(outFlowToken)
		return err
	}
	fmt.Println("StartDutchAuction")
	auctionParams := k.GetParams(ctx)
	//need to get real price instead of hard coding
	//calculate target amount of cmst to collect
	if auctiontypes.TestFlag != 1 {
		// inFlowTokenPrice, found1 = k.GetPriceForAsset(ctx, assetInId)
		// if !found1 {
		// 	return auctiontypes.ErrorPrices
		// }
		outFlowTokenPrice, found2 = k.GetPriceForAsset(ctx, assetOutId)
		if !found2 {
			return auctiontypes.ErrorPrices
		}
	} else {
		outFlowTokenPrice = uint64(2)
		inFlowTokenPrice = uint64(10)
	}
	inFlowTokenTargetAmount := k.getInflowTokenTargetAmount(outFlowToken.Amount, sdk.NewIntFromUint64(inFlowTokenPrice), sdk.NewIntFromUint64(outFlowTokenPrice))
	inFlowTokenTarget := sdk.NewCoin(inFlowToken.Denom, inFlowTokenTargetAmount)
	outFlowTokenInitialPrice := k.getOutflowTokenInitialPrice(sdk.NewIntFromUint64(outFlowTokenPrice), auctionParams.Buffer)
	outFlowTokenEndPrice := k.getOutflowTokenEndPrice(outFlowTokenInitialPrice, auctionParams.Cusp)
	vaultOwner, err := sdk.AccAddressFromBech32(lockedVaultOwner)
	if err != nil {
		return err
	}
	timeNow := ctx.BlockTime()
	inFlowTokenCurrentAmount := sdk.NewCoin(inFlowToken.Denom, sdk.NewIntFromUint64(0))
	auction := auctiontypes.DutchAuction{
		OutflowTokenInitAmount:    outFlowToken,
		OutflowTokenCurrentAmount: outFlowToken,
		InflowTokenTargetAmount:   inFlowTokenTarget,
		InflowTokenCurrentAmount:  inFlowTokenCurrentAmount,
		OutflowTokenInitialPrice:  outFlowTokenInitialPrice,
		OutflowTokenCurrentPrice:  outFlowTokenInitialPrice,
		OutflowTokenEndPrice:      outFlowTokenEndPrice,
		InflowTokenCurrentPrice:   sdk.NewDecFromInt(sdk.NewIntFromUint64(inFlowTokenPrice)),
		StartTime:                 timeNow,
		EndTime:                   timeNow.Add(time.Second * time.Duration(auctionParams.AuctionDurationSeconds)),
		AuctionStatus:             auctiontypes.AuctionStartNoBids,
		BiddingIds:                []*auctiontypes.BidOwnerMapping{},
		AuctionMappingId:          auctionParams.DutchId,
		AppId:                     appId,
		AssetInId:                 assetInId,
		AssetOutId:                assetOutId,
		LockedVaultId:             lockedVaultId,
		VaultOwner:                vaultOwner,
		LiquidationPenalty:        liquidationPenalty,
	}
	auction.AuctionId = k.GetAuctionID(ctx) + 1
	fmt.Println("k.GetAuctionID(ctx)", k.GetAuctionID(ctx))
	k.SetAuctionID(ctx, auction.AuctionId)
	fmt.Println("acuction id-----------------", auction.AuctionId)
	err = k.SetDutchAuction(ctx, auction)
	if err != nil {
		return err
	}
	err = k.SetFlagIsAuctionInProgress(ctx, lockedVaultId, true)
	if err != nil {
		return err
	}
	err = k.DecreaseLockedVaultAmountIn(ctx, lockedVaultId, outFlowToken.Amount)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) CloseSurplusAuction(
	ctx sdk.Context,
	surplusAuction auctiontypes.SurplusAuction,
) error {

	if surplusAuction.Bidder != nil {

		highestBidReceived := surplusAuction.Bid

		err := k.SendCoinsFromModuleToAccount(ctx, auctiontypes.ModuleName, surplusAuction.Bidder, sdk.NewCoins(surplusAuction.OutflowToken))
		if err != nil {
			return err
		}
		bidding, err := k.GetSurplusUserBidding(ctx, surplusAuction.Bidder.String(), surplusAuction.AppId, surplusAuction.ActiveBiddingId)
		if err != nil {
			return err
		}
		bidding.BiddingStatus = auctiontypes.SuccessBiddingStatus
		err = k.SetSurplusUserBidding(ctx, bidding)
		if err != nil {
			return err
		}

		if auctiontypes.TestFlag == 1 {
			//following 4 lines used for testing purpose
			err = k.BurnCoins(ctx, auctiontypes.ModuleName, highestBidReceived)
			if err != nil {
				return auctiontypes.ErrorInvalidBurn
			}
		} else {
			//burn tokens by sending bid tokens from auction to tokenmint module and then call burn function
			err = k.SendCoinsFromModuleToModule(ctx, auctiontypes.ModuleName, tokenminttypes.ModuleName, sdk.NewCoins(highestBidReceived))
			if err != nil {
				return err
			}
			err = k.BurnTokensForApp(ctx, surplusAuction.AppId, surplusAuction.AssetId, highestBidReceived.Amount)
			if err != nil {
				fmt.Println(err)
				return err
			}
		}

		for _, biddingId := range surplusAuction.BiddingIds {
			bidding, err := k.GetSurplusUserBidding(ctx, biddingId.BidOwner, surplusAuction.AppId, biddingId.BidId)
			if err != nil {
				continue
			}
			bidding.AuctionStatus = auctiontypes.ClosedAuctionStatus
			err = k.SetSurplusUserBidding(ctx, bidding)
			if err != nil {
				return err
			}
			err = k.DeleteSurplusUserBidding(ctx, bidding)
			if err != nil {
				return err
			}
			err = k.SetHistorySurplusUserBidding(ctx, bidding)
			if err != nil {
				return err
			}
		}
	} else {
		err1 := k.SendCoinsFromModuleToModule(ctx, auctiontypes.ModuleName, collectortypes.ModuleName, sdk.NewCoins(surplusAuction.OutflowToken))
		if err1 != nil {
			return err1
		}
		err2 := k.SetNetFeeCollectedData(ctx, surplusAuction.AppId, surplusAuction.AssetId, surplusAuction.OutflowToken.Amount)
		if err2 != nil {
			return auctiontypes.ErrorUnableToSetNetfees
		}
	}
	err := k.makeFalseForFlags(ctx, surplusAuction.AppId, surplusAuction.AssetId)
	if err != nil {
		return auctiontypes.ErrorUnableToMakeFlagsFalse
	}
	err = k.DeleteSurplusAuction(ctx, surplusAuction)
	if err != nil {
		return err
	}
	err = k.SetHistorySurplusAuction(ctx, surplusAuction)
	if err != nil {
		return err
	}
	//store auctions and user bidding in history after they are deleted
	return nil
}

func (k Keeper) CloseDebtAuction(
	ctx sdk.Context,
	debtAuction auctiontypes.DebtAuction,
) error {

	//If there are bids
	if debtAuction.AuctionStatus != auctiontypes.AuctionStartNoBids {

		if auctiontypes.TestFlag == 1 {
			//following 6 lines used for testing purpose
			err := k.MintCoins(ctx, auctiontypes.ModuleName, debtAuction.CurrentBidAmount)
			err = k.SendCoinsFromModuleToAccount(ctx, auctiontypes.ModuleName, debtAuction.Bidder, sdk.NewCoins(debtAuction.CurrentBidAmount))
			if err != nil {
				return err
			}
		} else {
			//ask token mint to mint new tokens for bidder address
			err := k.MintNewTokensForApp(ctx, debtAuction.AppId, debtAuction.AssetId, debtAuction.Bidder.String(), debtAuction.CurrentBidAmount.Amount)
			if err != nil {
				return err
			}
		}
		bidding, err := k.GetDebtUserBidding(ctx, debtAuction.Bidder.String(), debtAuction.AppId, debtAuction.ActiveBiddingId)
		if err != nil {
			return err
		}
		bidding.BiddingStatus = auctiontypes.SuccessBiddingStatus
		err = k.SetDebtUserBidding(ctx, bidding)
		if err != nil {
			return err
		}
		for _, biddingId := range debtAuction.BiddingIds {
			bidding, err := k.GetDebtUserBidding(ctx, biddingId.BidOwner, debtAuction.AppId, biddingId.BidId)
			if err != nil {
				return err
			}
			bidding.AuctionStatus = auctiontypes.ClosedAuctionStatus
			err = k.SetDebtUserBidding(ctx, bidding)
			if err != nil {
				return err
			}
			err = k.DeleteDebtUserBidding(ctx, bidding)
			if err != nil {
				return err
			}
			err = k.SetHistoryDebtUserBidding(ctx, bidding)
			if err != nil {
				return err
			}
		}
		//send to collector module the amount collected in debt auction
		err = k.SendCoinsFromModuleToModule(ctx, auctiontypes.ModuleName, collectortypes.ModuleName, sdk.NewCoins(debtAuction.ExpectedUserToken))
		if err != nil {
			return err
		}
		err = k.SetNetFeeCollectedData(ctx, debtAuction.AuctionId, debtAuction.AssetId, debtAuction.ExpectedUserToken.Amount)
		if err != nil {
			return auctiontypes.ErrorUnableToSetNetfees
		}
		return auctiontypes.ErrorInvalidBidId
	}
	err := k.makeFalseForFlags(ctx, debtAuction.AppId, debtAuction.AssetId)
	if err != nil {
		return auctiontypes.ErrorUnableToMakeFlagsFalse
	}
	err = k.DeleteDebtAuction(ctx, debtAuction)
	if err != nil {
		return err
	}
	err = k.SetHistoryDebtAuction(ctx, debtAuction)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) CloseDutchAuction(
	ctx sdk.Context,
	dutchAuction auctiontypes.DutchAuction,
) error {

	//delete dutch biddings
	if dutchAuction.AuctionStatus != auctiontypes.AuctionStartNoBids {
		for _, biddingId := range dutchAuction.BiddingIds {
			bidding, err := k.GetDutchUserBidding(ctx, biddingId.BidOwner, dutchAuction.AppId, biddingId.BidId)
			if err != nil {
				return err
			}
			bidding.AuctionStatus = auctiontypes.ClosedAuctionStatus
			err = k.SetDutchUserBidding(ctx, bidding)
			if err != nil {
				return err
			}
			err = k.DeleteDutchUserBidding(ctx, bidding)
			if err != nil {
				return err
			}
			err = k.SetHistoryDutchUserBidding(ctx, bidding)
			if err != nil {
				return err
			}
		}
	}

	fmt.Println("hello4__________________")
	lockedVault, found := k.GetLockedVault(ctx, dutchAuction.LockedVaultId)
	if !found {
		return auctiontypes.ErrorVaultNotFound
	}

	// burn and send target CMST to collector
	burnToken := sdk.NewCoin(dutchAuction.InflowTokenCurrentAmount.Denom, sdk.ZeroInt())
	//doing burn amount  = inflowtokencurrentamount / (1 + liq_penalty)
	burnToken.Amount = burnToken.Amount.Add(k.getBurnAmount(dutchAuction.InflowTokenCurrentAmount.Amount, dutchAuction.LiquidationPenalty))
	fmt.Println("hello7__________________")
	err := k.SendCoinsFromModuleToModule(ctx, auctiontypes.ModuleName, tokenminttypes.ModuleName, sdk.NewCoins(burnToken))
	if err != nil {
		return err
	}
	fmt.Println("hello8__________________")
	fmt.Println(dutchAuction.AppId, dutchAuction.AssetInId, burnToken)

	//TODO UNCOMMENT
	//err = k.tokenmint.BurnTokensForApp(ctx, dutchAuction.AppId, dutchAuction.AssetInId, burnToken.Amount)
	//if err != nil {
	//	return err
	//}

	fmt.Println("hello9__________________")
	//calculate penalty
	penaltyAmount := dutchAuction.InflowTokenCurrentAmount.Amount.Sub(burnToken.Amount)
	err = k.SendCoinsFromModuleToModule(ctx, auctiontypes.ModuleName, collectortypes.ModuleName, sdk.NewCoins(sdk.NewCoin(burnToken.Denom, penaltyAmount)))
	if err != nil {
		return err
	}
	fmt.Println("hello6__________________")
	//call increase function in collector
	err = k.SetNetFeeCollectedData(ctx, dutchAuction.AppId, dutchAuction.AssetInId, penaltyAmount)
	if err != nil {
		return err
	}
	lockedVault.AmountOut = lockedVault.AmountOut.Sub(burnToken.Amount)
	lockedVault.UpdatedAmountOut = lockedVault.UpdatedAmountOut.Sub(dutchAuction.InflowTokenCurrentAmount.Amount)
	fmt.Println("locked vault here amount amount out")
	fmt.Println(burnToken)
	fmt.Println(lockedVault)

	//set sell of history in locked vault
	outFlowToken := dutchAuction.OutflowTokenInitAmount.Sub(dutchAuction.OutflowTokenCurrentAmount)
	sellOfHistory := outFlowToken.String() + dutchAuction.InflowTokenCurrentAmount.String()
	lockedVault.SellOffHistory = append(lockedVault.SellOffHistory, sellOfHistory)
	k.SetLockedVault(ctx, lockedVault)
	fmt.Println("hello5__________________")
	dutchAuction.AuctionStatus = auctiontypes.AuctionEnded

	//update locked vault
	err = k.SetFlagIsAuctionComplete(ctx, dutchAuction.LockedVaultId, true)
	if err != nil {
		return err
	}
	fmt.Println("hello3_________________")
	err = k.SetFlagIsAuctionInProgress(ctx, dutchAuction.LockedVaultId, false)
	if err != nil {
		return err
	}

	err = k.SetDutchAuction(ctx, dutchAuction)
	if err != nil {
		return err
	}
	err = k.DeleteDutchAuction(ctx, dutchAuction)
	if err != nil {
		return err
	}
	err = k.SetHistoryDutchAuction(ctx, dutchAuction)
	if err != nil {
		return err
	}
	var appExtendedPairVaultData vaulttypes.AppExtendedPairVaultMapping
	var extendedPairVaultMapping vaulttypes.ExtendedPairVaultMapping

	appExtpair, _ := k.GetAppExtendedPairVaultMapping(ctx, lockedVault.AppMappingId)
	appExtendedPairVaultData.AppMappingId = lockedVault.AppMappingId
	appExtendedPairVaultData.Counter = appExtpair.Counter

	for _, data := range appExtpair.ExtendedPairVaults {
		if data.ExtendedPairId == lockedVault.ExtendedPairId {
			extendedPairVaultMapping.ExtendedPairId = lockedVault.ExtendedPairId
			extendedPairVaultMapping.VaultIds = data.VaultIds
			extendedPairVaultMapping.CollateralLockedAmount = data.CollateralLockedAmount.Sub(outFlowToken.Amount)
			extendedPairVaultMapping.TokenMintedAmount = data.TokenMintedAmount.Sub(burnToken.Amount)
		}
	}
	appExtendedPairVaultData.ExtendedPairVaults = append(appExtendedPairVaultData.ExtendedPairVaults, &extendedPairVaultMapping)

	err = k.SetAppExtendedPairVaultMapping(ctx, appExtendedPairVaultData)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) CreateNewSurplusBid(ctx sdk.Context, appId, auctionMappingId, auctionId uint64, bidder sdk.AccAddress, bid sdk.Coin) (biddingId uint64, err error) {
	auction, err := k.GetSurplusAuction(ctx, appId, auctionMappingId, auctionId)
	if err != nil {
		return biddingId, err
	}
	bidding := auctiontypes.SurplusBiddings{
		BiddingId:           k.GetUserBiddingID(ctx) + 1,
		AuctionId:           auctionId,
		AuctionStatus:       auctiontypes.ActiveAuctionStatus,
		AuctionedCollateral: auction.OutflowToken,
		Bidder:              bidder.String(),
		Bid:                 bid,
		BiddingTimestamp:    ctx.BlockTime(),
		BiddingStatus:       auctiontypes.PlacedBiddingStatus,
		AppId:               appId,
		AuctionMappingId:    auctionMappingId,
	}
	k.SetUserBiddingID(ctx, bidding.BiddingId)
	err = k.SetSurplusUserBidding(ctx, bidding)
	if err != nil {
		return biddingId, err
	}
	return bidding.BiddingId, nil
}

func (k Keeper) CreateNewDebtBid(ctx sdk.Context, appId, auctionMappingId, auctionId uint64, bidder sdk.AccAddress, bid sdk.Coin, expectedUserToken sdk.Coin) (biddingId uint64, err error) {
	bidding := auctiontypes.DebtBiddings{
		BiddingId:        k.GetUserBiddingID(ctx) + 1,
		AuctionId:        auctionId,
		AuctionStatus:    auctiontypes.ActiveAuctionStatus,
		Bidder:           bidder.String(),
		Bid:              bid,
		BiddingTimestamp: ctx.BlockTime(),
		BiddingStatus:    auctiontypes.PlacedBiddingStatus,
		AppId:            appId,
		AuctionMappingId: auctionMappingId,
		OutflowTokens:    expectedUserToken,
	}
	fmt.Println(auctionMappingId)
	k.SetUserBiddingID(ctx, bidding.BiddingId)

	err = k.SetDebtUserBidding(ctx, bidding)
	if err != nil {
		return biddingId, err
	}

	return bidding.BiddingId, nil
}

func (k Keeper) CreateNewDutchBid(ctx sdk.Context, appId, auctionMappingId, auctionId uint64, bidder sdk.AccAddress, outFlowTokenCoin sdk.Coin, inFlowTokenCoin sdk.Coin) (biddingId uint64, err error) {
	bidding := auctiontypes.DutchBiddings{
		BiddingId:          k.GetUserBiddingID(ctx) + 1,
		AuctionId:          auctionId,
		AuctionStatus:      auctiontypes.ActiveAuctionStatus,
		Bidder:             bidder.String(),
		OutflowTokenAmount: outFlowTokenCoin,
		InflowTokenAmount:  inFlowTokenCoin,
		BiddingTimestamp:   ctx.BlockTime(),
		BiddingStatus:      auctiontypes.SuccessBiddingStatus,
		AppId:              appId,
		AuctionMappingId:   auctionMappingId,
	}
	k.SetUserBiddingID(ctx, bidding.BiddingId)
	err = k.SetDutchUserBidding(ctx, bidding)
	if err != nil {
		return biddingId, err
	}
	return bidding.BiddingId, nil
}

func (k Keeper) PlaceSurplusBid(ctx sdk.Context, appId, auctionMappingId, auctionId uint64, bidder sdk.AccAddress, bid sdk.Coin) error {
	auction, err := k.GetSurplusAuction(ctx, appId, auctionMappingId, auctionId)
	if err != nil {
		return auctiontypes.ErrorInvalidSurplusAuctionId
	}
	if bid.Denom != auction.InflowToken.Denom {
		return auctiontypes.ErrorInvalidBiddingDenom
	}
	//Test this multiplication check if new bid greater than previous bid by bid factor
	if auction.AuctionStatus != auctiontypes.AuctionStartNoBids {
		change := auction.BidFactor.MulInt(auction.Bid.Amount).Ceil().TruncateInt()
		minBidAmount := auction.Bid.Amount.Add(change)
		if bid.Amount.LT(minBidAmount) {
			return auctiontypes.ErrorLowBidAmount
		}
	} else {
		if bid.Amount.LT(auction.Bid.Amount) {
			return auctiontypes.ErrorLowBidAmount
		}
	}
	err = k.SendCoinsFromAccountToModule(ctx, bidder, auctiontypes.ModuleName, sdk.NewCoins(bid))
	if err != nil {
		return err
	}
	fmt.Println(bidder)
	biddingId, err := k.CreateNewSurplusBid(ctx, auctionId, auctionMappingId, auctionId, bidder, bid)
	if err != nil {
		return err
	}
	if auction.AuctionStatus != auctiontypes.AuctionStartNoBids {
		// auction.Bidder as previous bidder . refund previous bidder
		err = k.SendCoinsFromModuleToAccount(ctx, auctiontypes.ModuleName, auction.Bidder, sdk.NewCoins(auction.Bid))
		if err != nil {
			return err
		}
		bidding, _ := k.GetSurplusUserBidding(ctx, auction.Bidder.String(), auction.AppId, auction.ActiveBiddingId)
		bidding.BiddingStatus = auctiontypes.RejectedBiddingStatus
		err = k.SetSurplusUserBidding(ctx, bidding)
		if err != nil {
			return err
		}
	} else {
		auction.AuctionStatus = auctiontypes.AuctionGoingOn
	}
	auction.ActiveBiddingId = biddingId
	var bidIdOwner = &auctiontypes.BidOwnerMapping{BidId: biddingId, BidOwner: bidder.String()}
	auction.BiddingIds = append(auction.BiddingIds, bidIdOwner)
	auction.Bidder = bidder
	auction.Bid = bid
	err = k.SetSurplusAuction(ctx, auction)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) PlaceDebtBid(ctx sdk.Context, appId, auctionMappingId, auctionId uint64, bidder sdk.AccAddress, bid sdk.Coin, expectedUserToken sdk.Coin) error {
	auction, err := k.GetDebtAuction(ctx, appId, auctionMappingId, auctionId)

	if err != nil {
		return auctiontypes.ErrorInvalidDebtAuctionId
	}
	if expectedUserToken.Denom != auction.ExpectedUserToken.Denom {
		return auctiontypes.ErrorInvalidDebtUserExpectedDenom
	}
	if !expectedUserToken.Amount.Equal(auction.ExpectedUserToken.Amount) {
		return auctiontypes.ErrorDebtExpectedUserAmount
	}
	if bid.Denom != auction.ExpectedMintedToken.Denom {
		return auctiontypes.ErrorInvalidDebtMintedDenom
	}

	//Test this multiplication check if new bid greater than previous bid by bid factor
	if auction.AuctionStatus != auctiontypes.AuctionStartNoBids {
		change := auction.BidFactor.MulInt(auction.ExpectedMintedToken.Amount).Ceil().TruncateInt()
		maxBidAmount := auction.ExpectedMintedToken.Amount.Sub(change)
		if bid.Amount.GT(maxBidAmount) {
			return auctiontypes.ErrorMaxBidAmount
		}
	} else {
		if bid.Amount.GT(auction.AuctionedToken.Amount) {
			return auctiontypes.ErrorMaxBidAmount
		}
	}
	err = k.SendCoinsFromAccountToModule(ctx, bidder, auctiontypes.ModuleName, sdk.NewCoins(expectedUserToken))
	if err != nil {
		return err
	}
	fmt.Println(auctionMappingId)
	biddingId, err := k.CreateNewDebtBid(ctx, appId, auctionMappingId, auctionId, bidder, bid, expectedUserToken)
	if err != nil {
		return err
	}
	//If auction gets bid from second time onwards . refund previous bidder
	if auction.AuctionStatus != auctiontypes.AuctionStartNoBids {
		err = k.SendCoinsFromModuleToAccount(ctx, auctiontypes.ModuleName, auction.Bidder, sdk.NewCoins(auction.ExpectedUserToken))
		if err != nil {
			return err
		}
		bidding, _ := k.GetDebtUserBidding(ctx, auction.Bidder.String(), auction.AppId, auction.ActiveBiddingId)
		bidding.BiddingStatus = auctiontypes.RejectedBiddingStatus

		err = k.SetDebtUserBidding(ctx, bidding)
		if err != nil {
			return err
		}
	} else {
		auction.AuctionStatus = auctiontypes.AuctionGoingOn
	}
	auction.ActiveBiddingId = biddingId
	var bidIdOwner = &auctiontypes.BidOwnerMapping{BidId: biddingId, BidOwner: bidder.String()}
	auction.BiddingIds = append(auction.BiddingIds, bidIdOwner)
	auction.Bidder = bidder
	auction.CurrentBidAmount = bid
	auction.ExpectedMintedToken = bid
	err = k.SetDebtAuction(ctx, auction)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) PlaceDutchBid(ctx sdk.Context, appId, auctionMappingId, auctionId uint64, bidder sdk.AccAddress, bid sdk.Coin, max sdk.Dec) error {
	auction, err := k.GetDutchAuction(ctx, appId, auctionMappingId, auctionId)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrNotFound, "auction id %d not found", auctionId)
	}
	if bid.Denom != auction.OutflowTokenCurrentAmount.Denom {
		return sdkerrors.Wrapf(sdkerrors.ErrNotFound, "bid denom %s not found", bid.Denom)
	}
	fmt.Println("auction data----------------------------------")

	if max.LT(auction.OutflowTokenCurrentPrice.Ceil()) {
		return auctiontypes.ErrorInvalidDutchPrice
	}

	fmt.Println("price data----------------------------------")
	// slice tells amount of collateral user should be given
	auctionParams := k.GetParams(ctx)
	//using ceil as we need extract more from users
	outFlowTokenCurrentPrice := auction.OutflowTokenCurrentPrice.Ceil().TruncateInt()
	inFlowTokenCurrentPrice := auction.InflowTokenCurrentPrice.Ceil().TruncateInt()
	slice := sdk.MinInt(bid.Amount, auction.OutflowTokenCurrentAmount.Amount)
	//amount in usd to be given to user
	owe := slice.Mul(outFlowTokenCurrentPrice)
	//required target cmst to raise in usd
	tab := auction.InflowTokenTargetAmount.Amount.Mul(inFlowTokenCurrentPrice).Sub(auction.InflowTokenCurrentAmount.Amount)

	inFlowTokenAmount := slice.ToDec().Mul(outFlowTokenCurrentPrice.ToDec()).Quo(inFlowTokenCurrentPrice.ToDec()).Ceil().TruncateInt()
	inFlowTokenCoin := sdk.NewCoin(auction.InflowTokenTargetAmount.Denom, inFlowTokenAmount)
	//check if bid in usd is greater than required target cmst in usd
	if owe.GT(tab) {
		slice = tab.Quo(auction.OutflowTokenCurrentPrice.Ceil().TruncateInt())
		inFlowTokenCoin.Amount = auction.InflowTokenTargetAmount.Amount.Sub(auction.InflowTokenCurrentAmount.Amount)
	} else if auction.OutflowTokenCurrentAmount.Amount.ToDec().Sub(slice.ToDec()).Mul(outFlowTokenCurrentPrice.ToDec()).TruncateInt().LT(auctionParams.Chost.Ceil().TruncateInt()) {
		//(outflowtokenavailableamount-slice) in usd < chost in usd
		//see if user has balance to buy whole collateral
		userBalanceUsd := k.GetBalance(ctx, bidder, bid.Denom).Amount.Mul(outFlowTokenCurrentPrice)
		collateralAvailableUsd := auction.OutflowTokenCurrentAmount.Amount.Mul(outFlowTokenCurrentPrice)
		fmt.Println("entered  chost")
		fmt.Println("collateral usd")
		fmt.Println(collateralAvailableUsd)
		fmt.Println("user balance usd")
		fmt.Println(userBalanceUsd)
		if userBalanceUsd.LT(collateralAvailableUsd) {
			return auctiontypes.ErrorDutchinsufficientUserBalance
		}
		slice = auction.OutflowTokenCurrentAmount.Amount
		fmt.Println("entered  chost slice")
		fmt.Println(slice)
	}

	outFlowTokenCoin := sdk.NewCoin(auction.OutflowTokenInitAmount.Denom, slice)
	fmt.Println("outflow coin")
	fmt.Println(outFlowTokenCoin)
	fmt.Println("inflow coin")
	fmt.Println(inFlowTokenCoin)
	err = k.SendCoinsFromAccountToModule(ctx, bidder, auctiontypes.ModuleName, sdk.NewCoins(inFlowTokenCoin))
	if err != nil {
		return err
	}
	err = k.SendCoinsFromModuleToAccount(ctx, auctiontypes.ModuleName, bidder, sdk.NewCoins(outFlowTokenCoin))
	if err != nil {
		fmt.Println("hey_________")
		fmt.Println(err)
		return err
	}
	//create user bidding
	biddingId, err := k.CreateNewDutchBid(ctx, appId, auctionMappingId, auctionId, bidder, inFlowTokenCoin, outFlowTokenCoin)
	if err != nil {
		return err
	}
	var bidIdOwner = &auctiontypes.BidOwnerMapping{BidId: biddingId, BidOwner: bidder.String()}
	auction.BiddingIds = append(auction.BiddingIds, bidIdOwner)
	if auction.AuctionStatus == auctiontypes.AuctionStartNoBids {
		auction.AuctionStatus = auctiontypes.AuctionGoingOn
	}

	//calculate inflow amount and outflow amount if  user  transaction successfull
	auction.OutflowTokenCurrentAmount = auction.OutflowTokenCurrentAmount.Sub(outFlowTokenCoin)
	auction.InflowTokenCurrentAmount = auction.InflowTokenCurrentAmount.Add(inFlowTokenCoin)

	//collateral not over but target cmst reached then send remaining collateral to owner
	//if inflow token current amount > InflowTokenTargetAmount
	if auction.InflowTokenCurrentAmount.IsGTE(auction.InflowTokenTargetAmount) {
		//send left overcollateral to vault owner as target cmst reached and also
		total := auction.OutflowTokenCurrentAmount
		err := k.SendCoinsFromModuleToModule(ctx, auctiontypes.ModuleName, vaulttypes.ModuleName, sdk.NewCoins(total))
		if err != nil {
			return err
		}
		err = k.IncreaseLockedVaultAmountIn(ctx, auction.LockedVaultId, total.Amount)
		fmt.Println("sending user amount to vault")
		fmt.Println(total)
		err = k.SetDutchAuction(ctx, auction)
		if err != nil {
			return err
		}
		//remove dutch auction
		fmt.Println("bmw1_______________________")
		err = k.CloseDutchAuction(ctx, auction)
		if err != nil {
			return err
		}
	} else if auction.OutflowTokenCurrentAmount.Amount.IsZero() { //entire collateral sold out
		err = k.SetDutchAuction(ctx, auction)
		if err != nil {
			return err
		}
		//remove dutch auction
		fmt.Println("bmw2_______________________")
		err = k.CloseDutchAuction(ctx, auction)
		if err != nil {
			return err
		}
	} else {
		fmt.Println("bmw3_______________________")
		err = k.SetDutchAuction(ctx, auction)
		if err != nil {
			return err
		}
	}
	return nil
}
