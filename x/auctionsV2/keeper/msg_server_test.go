package keeper_test

import (
	"fmt"
	"github.com/comdex-official/comdex/app/wasm/bindings"
	utils "github.com/comdex-official/comdex/types"
	assetTypes "github.com/comdex-official/comdex/x/asset/types"
	"github.com/comdex-official/comdex/x/auctionsV2"
	auctionsV2types "github.com/comdex-official/comdex/x/auctionsV2/types"
	collectortypes "github.com/comdex-official/comdex/x/collector/types"
	lendKeeper "github.com/comdex-official/comdex/x/lend/keeper"
	lendtypes "github.com/comdex-official/comdex/x/lend/types"
	"github.com/comdex-official/comdex/x/liquidationsV2/types"
	markettypes "github.com/comdex-official/comdex/x/market/types"
	vaultKeeper1 "github.com/comdex-official/comdex/x/vault/keeper"
	vaultTypes "github.com/comdex-official/comdex/x/vault/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (s *KeeperTestSuite) AddAppAssets() {

	assetOneID := s.CreateNewAsset("ASSETONE", "uasset1", 2000000)
	assetTwoID := s.CreateNewAsset("ASSETTWO", "uasset2", 2000000)
	assetThreeID := s.CreateNewAsset("ASSETTHREE", "uasset3", 1000000)
	assetFourID := s.CreateNewAsset("ASSETFOUR", "uasset4", 2000000)
	cAssetOneID := s.CreateNewAsset("CASSETONE", "ucasset1", 1000000)
	cAssetTwoID := s.CreateNewAsset("CASSETTWO", "ucasset2", 2000000)
	cAssetThreeID := s.CreateNewAsset("CASSETTHRE", "ucasset3", 2000000)
	cAssetFourID := s.CreateNewAsset("CASSETFOUR", "ucasset4", 2000000)

	var (
		assetDataPoolOne []*lendtypes.AssetDataPoolMapping
		assetDataPoolTwo []*lendtypes.AssetDataPoolMapping
	)
	assetDataPoolOneAssetOne := &lendtypes.AssetDataPoolMapping{
		AssetID:          assetOneID,
		AssetTransitType: 3,
		SupplyCap:        sdk.NewDec(5000000000000000000),
	}
	assetDataPoolOneAssetTwo := &lendtypes.AssetDataPoolMapping{
		AssetID:          assetTwoID,
		AssetTransitType: 1,
		SupplyCap:        sdk.NewDec(1000000000000000000),
	}
	assetDataPoolOneAssetThree := &lendtypes.AssetDataPoolMapping{
		AssetID:          assetThreeID,
		AssetTransitType: 2,
		SupplyCap:        sdk.NewDec(5000000000000000000),
	}
	assetDataPoolTwoAssetFour := &lendtypes.AssetDataPoolMapping{
		AssetID:          assetFourID,
		AssetTransitType: 1,
		SupplyCap:        sdk.NewDec(3000000000000000000),
	}

	assetDataPoolOne = append(assetDataPoolOne, assetDataPoolOneAssetOne, assetDataPoolOneAssetTwo, assetDataPoolOneAssetThree)
	assetDataPoolTwo = append(assetDataPoolTwo, assetDataPoolTwoAssetFour, assetDataPoolOneAssetOne, assetDataPoolOneAssetThree)

	s.AddAssetRatesStats(assetThreeID, newDec("0.8"), newDec("0.002"), newDec("0.06"), newDec("0.6"), true, newDec("0.04"), newDec("0.04"), newDec("0.06"), newDec("0.8"), newDec("0.85"), newDec("0.025"), newDec("0.025"), newDec("0.1"), cAssetThreeID)
	s.AddAssetRatesStats(assetOneID, newDec("0.75"), newDec("0.002"), newDec("0.07"), newDec("1.25"), false, newDec("0.0"), newDec("0.0"), newDec("0.0"), newDec("0.7"), newDec("0.75"), newDec("0.05"), newDec("0.05"), newDec("0.2"), cAssetOneID)
	s.AddAssetRatesPoolPairs(assetTwoID, newDec("0.5"), newDec("0.002"), newDec("0.08"), newDec("2.0"), false, newDec("0.0"), newDec("0.0"), newDec("0.0"), newDec("0.5"), newDec("0.55"), newDec("0.05"), newDec("0.05"), newDec("0.2"), cAssetTwoID, "cmdx", "CMDX-ATOM-CMST", assetDataPoolOne, 1000000, false)
	s.AddAssetRatesPoolPairs(assetFourID, newDec("0.65"), newDec("0.002"), newDec("0.08"), newDec("1.5"), false, newDec("0.0"), newDec("0.0"), newDec("0.0"), newDec("0.6"), newDec("0.65"), newDec("0.05"), newDec("0.05"), newDec("0.2"), cAssetFourID, "osmo", "OSMO-ATOM-CMST", assetDataPoolTwo, 1000000, false)

	_ = s.CreateNewApp("cswap", "cswap")
	_ = s.CreateNewApp("harbor", "hbr")
	appThreeID := s.CreateNewApp("commodo", "cmdo")
	s.fundAddr(sdk.MustAccAddressFromBech32("cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t"), sdk.NewCoins(sdk.NewCoin("uasset1", newInt(1000000000000000))))
	s.fundAddr(sdk.MustAccAddressFromBech32("cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t"), sdk.NewCoins(sdk.NewCoin("uasset2", newInt(1000000000000000))))
	s.fundAddr(sdk.MustAccAddressFromBech32("cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t"), sdk.NewCoins(sdk.NewCoin("uasset3", newInt(1000000000000000))))
	s.fundAddr(sdk.MustAccAddressFromBech32("cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t"), sdk.NewCoins(sdk.NewCoin("uasset4", newInt(1000000000000000))))

	s.fundAddr(sdk.MustAccAddressFromBech32("cosmos1hm7w7dnvdnra78pz9qxysy7u4tuhc3fnpjmyj7"), sdk.NewCoins(sdk.NewCoin("uasset1", newInt(1000000000000000))))
	s.fundAddr(sdk.MustAccAddressFromBech32("cosmos1hm7w7dnvdnra78pz9qxysy7u4tuhc3fnpjmyj7"), sdk.NewCoins(sdk.NewCoin("uasset2", newInt(1000000000000000))))
	s.fundAddr(sdk.MustAccAddressFromBech32("cosmos1hm7w7dnvdnra78pz9qxysy7u4tuhc3fnpjmyj7"), sdk.NewCoins(sdk.NewCoin("uasset3", newInt(13000000))))

	msgLend1 := lendtypes.NewMsgLend("cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t", assetOneID, sdk.NewCoin("uasset1", newInt(3000000000)), 1, appThreeID)
	msgLend2 := lendtypes.NewMsgLend("cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t", assetTwoID, sdk.NewCoin("uasset2", newInt(10000000000)), 1, appThreeID)
	msgLend3 := lendtypes.NewMsgLend("cosmos1hm7w7dnvdnra78pz9qxysy7u4tuhc3fnpjmyj7", assetOneID, sdk.NewCoin("uasset1", newInt(10000000000)), 1, appThreeID)

	msg3 := lendtypes.NewMsgFundModuleAccounts(1, assetOneID, "cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t", sdk.NewCoin("uasset1", newInt(10000000000)))
	msg4 := lendtypes.NewMsgFundModuleAccounts(1, assetTwoID, "cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t", sdk.NewCoin("uasset2", newInt(10000000000)))
	msg5 := lendtypes.NewMsgFundModuleAccounts(1, assetThreeID, "cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t", sdk.NewCoin("uasset3", newInt(120000000)))
	msg7 := lendtypes.NewMsgFundModuleAccounts(2, assetOneID, "cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t", sdk.NewCoin("uasset1", newInt(10000000000)))
	msg8 := lendtypes.NewMsgFundModuleAccounts(2, assetFourID, "cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t", sdk.NewCoin("uasset4", newInt(10000000000)))

	lendkeeper := &s.lendKeeper
	server := lendKeeper.NewMsgServerImpl(*lendkeeper)

	_, _ = server.Lend(sdk.WrapSDKContext(s.ctx), msgLend1)
	_, _ = server.Lend(sdk.WrapSDKContext(s.ctx), msgLend2)
	_, _ = server.Lend(sdk.WrapSDKContext(s.ctx), msgLend3)
	_, _ = server.FundModuleAccounts(sdk.WrapSDKContext(s.ctx), msg3)
	_, _ = server.FundModuleAccounts(sdk.WrapSDKContext(s.ctx), msg4)
	_, _ = server.FundModuleAccounts(sdk.WrapSDKContext(s.ctx), msg5)
	_, _ = server.FundModuleAccounts(sdk.WrapSDKContext(s.ctx), msg7)
	_, _ = server.FundModuleAccounts(sdk.WrapSDKContext(s.ctx), msg8)

	msg2 := lendtypes.NewMsgBorrow("cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t", 1, 1, false, sdk.NewCoin("ucasset1", newInt(100000000)), sdk.NewCoin("uasset2", newInt(70000000)))
	_, err := server.Borrow(sdk.WrapSDKContext(s.ctx), msg2)
	s.Require().NoError(err)

	msg22 := lendtypes.NewMsgBorrow("cosmos1hm7w7dnvdnra78pz9qxysy7u4tuhc3fnpjmyj7", 3, 1, false, sdk.NewCoin("ucasset1", newInt(1000000000)), sdk.NewCoin("uasset2", newInt(700000000)))
	_, err = server.Borrow(sdk.WrapSDKContext(s.ctx), msg22)
	s.Require().NoError(err)

	pair := assetTypes.Pair{AssetIn: 2, AssetOut: 3}
	extendedPairVault := bindings.MsgAddExtendedPairsVault{
		AppID:               2,
		PairID:              1,
		StabilityFee:        sdk.MustNewDecFromStr("0.01"),
		ClosingFee:          sdk.MustNewDecFromStr("0"),
		LiquidationPenalty:  sdk.MustNewDecFromStr("0.12"),
		DrawDownFee:         sdk.MustNewDecFromStr("0.01"),
		IsVaultActive:       true,
		DebtCeiling:         sdk.NewInt(1000000000000),
		DebtFloor:           sdk.NewInt(1000000),
		IsStableMintVault:   false,
		MinCr:               sdk.MustNewDecFromStr("1.5"),
		PairName:            "CMDX-B",
		AssetOutOraclePrice: true,
		AssetOutPrice:       1000000,
		MinUsdValueLeft:     1000000,
	}

	assetkeeper, ctx := &s.assetKeeper, &s.ctx
	err = assetkeeper.AddPairsRecords(*ctx, pair)
	s.Require().NoError(err)

	err = assetkeeper.WasmAddExtendedPairsVaultRecords(*ctx, &extendedPairVault)
	s.Require().NoError(err)

	// set liquidation whitelisting
	dutchAuctionParams := types.DutchAuctionParam{
		Premium:         newDec("1.2"),
		Discount:        newDec("0.7"),
		DecrementFactor: sdk.NewInt(1),
	}
	englishAuctionParams := types.EnglishAuctionParam{DecrementFactor: sdk.NewInt(1)}

	liqWhitelistingHbr := types.LiquidationWhiteListing{
		AppId:               2,
		Initiator:           true,
		IsDutchActivated:    true,
		DutchAuctionParam:   &dutchAuctionParams,
		IsEnglishActivated:  true,
		EnglishAuctionParam: &englishAuctionParams,
		KeeeperIncentive:    newDec("0.1"),
	}
	s.liquidationKeeper.SetLiquidationWhiteListing(s.ctx, liqWhitelistingHbr)

	liqWhitelistingCmdo := types.LiquidationWhiteListing{
		AppId:               3,
		Initiator:           true,
		IsDutchActivated:    true,
		DutchAuctionParam:   &dutchAuctionParams,
		IsEnglishActivated:  false,
		EnglishAuctionParam: nil,
		KeeeperIncentive:    newDec("0.1"),
	}
	s.liquidationKeeper.SetLiquidationWhiteListing(s.ctx, liqWhitelistingCmdo)

	auctionParams := auctionsV2types.AuctionParams{
		AuctionDurationSeconds: 3600,
		Step:                   newDec("0.1"),
		WithdrawalFee:          newDec("0.0"),
		ClosingFee:             newDec("0.0"),
		MinUsdValueLeft:        100000,
		BidFactor:              newDec("0.1"),
		LiquidationPenalty:     newDec("0.1"),
		AuctionBonus:           newDec("0.0"),
	}

	s.addAuctionParams(auctionParams)

}

func (s *KeeperTestSuite) CreateVault() {
	userAddress1 := "cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t"
	userAddress2 := "cosmos1hm7w7dnvdnra78pz9qxysy7u4tuhc3fnpjmyj7"
	vaultKeeper, ctx := &s.vaultKeeper, &s.ctx
	s.AddAppAssets()
	server := vaultKeeper1.NewMsgServer(*vaultKeeper)

	for index, tc := range []struct {
		name string
		msg  vaultTypes.MsgCreateRequest
	}{
		{
			"Create Vault : AppID 1 extended pair 1 user address 1",
			vaultTypes.MsgCreateRequest{
				From:                userAddress1,
				AppId:               2,
				ExtendedPairVaultId: 1,
				AmountIn:            sdk.NewIntFromUint64(1000000),
				AmountOut:           sdk.NewIntFromUint64(1000000),
			},
		},
		{
			"Create Vault : AppID 1 extended pair 1 user address 2",
			vaultTypes.MsgCreateRequest{
				From:                userAddress2,
				AppId:               2,
				ExtendedPairVaultId: 1,
				AmountIn:            sdk.NewIntFromUint64(1000000),
				AmountOut:           sdk.NewIntFromUint64(1000000),
			},
		},
	} {
		s.Run(tc.name, func() {
			_, err := server.MsgCreate(sdk.WrapSDKContext(*ctx), &tc.msg)
			s.Require().NoError(err)
			res, err := s.vaultQuerier.QueryAllVaults(sdk.WrapSDKContext(*ctx), &vaultTypes.QueryAllVaultsRequest{})
			s.Require().NoError(err)
			_, err = s.vaultQuerier.QueryVaultInfoByVaultID(sdk.WrapSDKContext(*ctx), &vaultTypes.QueryVaultInfoByVaultIDRequest{Id: res.Vault[index].Id})
			s.Require().NoError(err)
		})
	}
}

func (s *KeeperTestSuite) GetVaultCount() int {
	ctx := &s.ctx
	res, err := s.vaultQuerier.QueryAllVaults(sdk.WrapSDKContext(*ctx), &vaultTypes.QueryAllVaultsRequest{})
	s.Require().NoError(err)
	return len(res.Vault)
}

func (s *KeeperTestSuite) GetBorrowsCount() int {
	ctx := &s.ctx
	res, err := s.lendQuerier.QueryBorrows(sdk.WrapSDKContext(*ctx), &lendtypes.QueryBorrowsRequest{})
	s.Require().NoError(err)
	return len(res.Borrows)
}

func (s *KeeperTestSuite) GetVaultCountForExtendedPairIDbyAppID(appID, extID uint64) int {
	vaultKeeper, ctx := &s.vaultKeeper, &s.ctx
	res, found := vaultKeeper.GetAppExtendedPairVaultMappingData(*ctx, appID, extID)
	s.Require().True(found)
	return len(res.VaultIds)
}

func (s *KeeperTestSuite) ChangeOraclePrice(asset uint64) {
	s.SetOraclePrice(asset, 1000000)
}

func (s *KeeperTestSuite) ChangeOraclePrice1(asset uint64) {
	s.SetOraclePrice(asset, 1800000)
}

func (s *KeeperTestSuite) SetOraclePrice(assetID uint64, price uint64) {
	market := markettypes.TimeWeightedAverage{
		AssetID:       assetID,
		ScriptID:      12,
		Twa:           price,
		CurrentIndex:  0,
		IsPriceActive: true,
		PriceValue:    []uint64{price},
	}
	s.app.MarketKeeper.SetTwa(s.ctx, market)
}

func (s *KeeperTestSuite) TestLiquidateVaults() {
	liquidationKeeper, ctx := &s.liquidationKeeper, &s.ctx
	s.CreateVault()
	currentVaultsCount := 2
	s.Require().Equal(s.GetVaultCount(), currentVaultsCount)
	s.Require().Equal(s.GetVaultCountForExtendedPairIDbyAppID(2, 1), currentVaultsCount)
	beforeVault, found := s.vaultKeeper.GetVault(*ctx, 1)
	s.Require().True(found)

	// Liquidation shouldn't happen as price not changed
	err := liquidationKeeper.Liquidate(*ctx)
	s.Require().NoError(err)
	id := liquidationKeeper.GetLockedVaultID(*ctx)
	s.Require().Equal(id, uint64(0))

	// Liquidation should happen as price changed
	s.ChangeOraclePrice(2)
	err = liquidationKeeper.Liquidate(*ctx)
	s.Require().NoError(err)
	id = liquidationKeeper.GetLockedVaultID(*ctx)
	s.Require().Equal(id, uint64(2))
	s.Require().Equal(s.GetVaultCount(), currentVaultsCount-2)
	s.Require().Equal(s.GetVaultCountForExtendedPairIDbyAppID(2, 1), currentVaultsCount-2)

	lockedVault := liquidationKeeper.GetLockedVaults(*ctx)
	s.Require().Equal(lockedVault[0].OriginalVaultId, beforeVault.Id)
	s.Require().Equal(lockedVault[0].ExtendedPairId, beforeVault.ExtendedPairVaultID)
	s.Require().Equal(lockedVault[0].Owner, beforeVault.Owner)
	s.Require().Equal(lockedVault[0].CollateralToken.Amount, beforeVault.AmountIn)
	s.Require().Equal(lockedVault[0].DebtToken.Amount, beforeVault.AmountOut)
	s.Require().Equal(lockedVault[0].TargetDebt.Amount, lockedVault[0].DebtToken.Amount.Add(sdk.NewDecFromInt(beforeVault.AmountOut).Mul(newDec("0.12")).TruncateInt()))
	s.Require().Equal(lockedVault[0].FeeToBeCollected, sdk.NewDecFromInt(beforeVault.AmountOut).Mul(newDec("0.12")).TruncateInt())
	s.Require().Equal(lockedVault[0].IsDebtCmst, false)
	s.Require().Equal(lockedVault[0].CollateralAssetId, uint64(2))
	s.Require().Equal(lockedVault[0].DebtAssetId, uint64(3))
	price, err := s.app.MarketKeeper.CalcAssetPrice(*ctx, 2, beforeVault.AmountIn)
	s.Require().NoError(err)
	s.Require().Equal(lockedVault[0].CollateralToBeAuctioned.Amount, price.TruncateInt())

	// get auctions and tally them
	auctions := s.app.NewaucKeeper.GetAuctions(s.ctx)
	s.Require().Equal(len(auctions), 2)
	s.Require().Equal(auctions[0].AppId, lockedVault[0].AppId)
	s.Require().Equal(auctions[0].AuctionType, lockedVault[0].AuctionType)
	s.Require().Equal(auctions[0].CollateralAssetId, lockedVault[0].CollateralAssetId)
	s.Require().Equal(auctions[0].DebtAssetId, lockedVault[0].DebtAssetId)
	s.Require().Equal(auctions[0].BonusAmount, lockedVault[0].BonusToBeGiven)
	s.Require().Equal(auctions[0].LockedVaultId, lockedVault[0].LockedVaultId)
	s.Require().Equal(auctions[0].CollateralToken, lockedVault[0].CollateralToken)
	s.Require().Equal(auctions[0].DebtToken, lockedVault[0].TargetDebt)

	twaDataCollateral, _ := s.app.MarketKeeper.GetTwa(s.ctx, lockedVault[0].CollateralAssetId)
	liquidationWhitelistingAppData, _ := s.app.NewliqKeeper.GetLiquidationWhiteListing(s.ctx, lockedVault[0].AppId)
	CollateralTokenInitialPrice := s.app.NewaucKeeper.GetCollalteralTokenInitialPrice(sdk.NewIntFromUint64(twaDataCollateral.Twa), liquidationWhitelistingAppData.DutchAuctionParam.Premium)
	s.Require().Equal(auctions[0].CollateralTokenAuctionPrice, CollateralTokenInitialPrice)

}

func (s *KeeperTestSuite) TestLiquidateBorrows() {
	liquidationKeeper, ctx := &s.liquidationKeeper, &s.ctx
	s.AddAppAssets()
	currentBorrowsCount := 2
	s.Require().Equal(s.GetBorrowsCount(), currentBorrowsCount)

	beforeBorrow, found := s.lendKeeper.GetBorrow(*ctx, 1)
	s.Require().True(found)

	beforeLend, found := s.lendKeeper.GetLend(*ctx, beforeBorrow.LendingID)
	s.Require().True(found)

	// Liquidation shouldn't happen as price not changed
	err := liquidationKeeper.Liquidate(*ctx)
	s.Require().NoError(err)
	id := liquidationKeeper.GetLockedVaultID(*ctx)
	s.Require().Equal(id, uint64(0))

	assetStatsLend, _ := s.lendKeeper.GetAssetStatsByPoolIDAndAssetID(*ctx, 1, 1)
	s.Require().Equal(len(assetStatsLend.LendIds), 2)
	s.Require().Equal(len(assetStatsLend.BorrowIds), 0)
	s.Require().Equal(assetStatsLend.TotalBorrowed, sdk.NewInt(0))
	s.Require().Equal(assetStatsLend.TotalLend, sdk.NewInt(13000000000))

	assetStatsBorrow, _ := s.lendKeeper.GetAssetStatsByPoolIDAndAssetID(*ctx, 1, 2)
	s.Require().Equal(len(assetStatsBorrow.LendIds), 1)
	s.Require().Equal(len(assetStatsBorrow.BorrowIds), 2)
	s.Require().Equal(assetStatsBorrow.TotalBorrowed, sdk.NewInt(770000000))
	s.Require().Equal(assetStatsBorrow.TotalLend, sdk.NewInt(10000000000))

	modBalInitial, _ := s.lendKeeper.GetModuleBalanceByPoolID(*ctx, 1)

	// Liquidation should happen as price changed
	s.ChangeOraclePrice1(1)
	err = liquidationKeeper.Liquidate(*ctx)
	s.Require().NoError(err)
	id = liquidationKeeper.GetLockedVaultID(*ctx)
	s.Require().Equal(id, uint64(2))
	s.Require().Equal(s.GetBorrowsCount(), currentBorrowsCount)

	lockedVault := liquidationKeeper.GetLockedVaults(*ctx)
	s.Require().Equal(lockedVault[0].OriginalVaultId, beforeBorrow.ID)
	s.Require().Equal(lockedVault[0].ExtendedPairId, beforeBorrow.PairID)
	s.Require().Equal(lockedVault[0].Owner, beforeLend.Owner)
	s.Require().Equal(lockedVault[0].DebtToken.Amount, beforeBorrow.AmountOut.Amount)
	s.Require().Equal(lockedVault[0].TargetDebt.Amount, lockedVault[0].DebtToken.Amount.Add(sdk.NewDecFromInt(beforeBorrow.AmountOut.Amount).Mul(newDec("0.05")).TruncateInt()))
	s.Require().Equal(lockedVault[0].FeeToBeCollected, sdk.NewDecFromInt(beforeBorrow.AmountOut.Amount).Mul(newDec("0.05")).TruncateInt())
	s.Require().Equal(lockedVault[0].IsDebtCmst, false)
	s.Require().Equal(lockedVault[0].CollateralAssetId, uint64(1))
	s.Require().Equal(lockedVault[0].DebtAssetId, uint64(2))

	// get data of total borrow and lend and tally
	assetStatsLend, _ = s.lendKeeper.GetAssetStatsByPoolIDAndAssetID(*ctx, 1, 1)
	s.Require().Equal(len(assetStatsLend.LendIds), 2)
	s.Require().Equal(len(assetStatsLend.BorrowIds), 0)
	s.Require().Equal(assetStatsLend.TotalBorrowed, sdk.NewInt(0))
	s.Require().Equal(assetStatsLend.TotalLend, sdk.NewInt(11900000000))

	assetStatsBorrow, _ = s.lendKeeper.GetAssetStatsByPoolIDAndAssetID(*ctx, 1, 2)
	s.Require().Equal(len(assetStatsBorrow.LendIds), 1)
	s.Require().Equal(len(assetStatsBorrow.BorrowIds), 2)
	s.Require().Equal(assetStatsBorrow.TotalBorrowed, sdk.NewInt(0))
	s.Require().Equal(assetStatsBorrow.TotalLend, sdk.NewInt(10000000000))

	afterBorrow, found := s.lendKeeper.GetBorrow(*ctx, 1)
	s.Require().True(found)
	s.Require().Equal(afterBorrow.IsLiquidated, true)

	modBalFinal, _ := s.lendKeeper.GetModuleBalanceByPoolID(*ctx, 1)
	s.Require().Equal(modBalInitial.ModuleBalanceStats[0].Balance.Amount.Sub(modBalFinal.ModuleBalanceStats[0].Balance.Amount), sdk.NewInt(1100000000))

	// get auctions and tally them
	auctions := s.app.NewaucKeeper.GetAuctions(s.ctx)
	s.Require().Equal(len(auctions), 2)
	s.Require().Equal(auctions[0].AppId, lockedVault[0].AppId)
	s.Require().Equal(auctions[0].AuctionType, lockedVault[0].AuctionType)
	s.Require().Equal(auctions[0].CollateralAssetId, lockedVault[0].CollateralAssetId)
	s.Require().Equal(auctions[0].DebtAssetId, lockedVault[0].DebtAssetId)
	s.Require().Equal(auctions[0].BonusAmount, lockedVault[0].BonusToBeGiven)
	s.Require().Equal(auctions[0].LockedVaultId, lockedVault[0].LockedVaultId)
	s.Require().Equal(auctions[0].CollateralToken, lockedVault[0].CollateralToken)
	s.Require().Equal(auctions[0].DebtToken, lockedVault[0].TargetDebt)

	twaDataCollateral, _ := s.app.MarketKeeper.GetTwa(s.ctx, lockedVault[0].CollateralAssetId)
	liquidationWhitelistingAppData, _ := s.app.NewliqKeeper.GetLiquidationWhiteListing(s.ctx, lockedVault[0].AppId)
	CollateralTokenInitialPrice := s.app.NewaucKeeper.GetCollalteralTokenInitialPrice(sdk.NewIntFromUint64(twaDataCollateral.Twa), liquidationWhitelistingAppData.DutchAuctionParam.Premium)
	s.Require().Equal(auctions[0].CollateralTokenAuctionPrice, CollateralTokenInitialPrice)
}

func (s *KeeperTestSuite) TestPlaceMarketBidForVaults() {
	s.TestLiquidateVaults()
	//auctionKeeper := &s.keeper
	bidder := "cosmos1hm7w7dnvdnra78pz9qxysy7u4tuhc3fnpjmyj7"
	testCases := []struct {
		Name    string
		Msg     auctionsV2types.MsgPlaceMarketBidRequest
		ExpErr  error
		ExpResp *auctionsV2types.MsgPlaceMarketBidResponse
	}{
		{
			Name:    "auction does not exist",
			Msg:     *auctionsV2types.NewMsgPlaceMarketBid(bidder, 10, sdk.NewCoin("uasset2", sdk.NewInt(100000))),
			ExpErr:  sdkerrors.ErrNotFound,
			ExpResp: nil,
		},
		{
			Name:    "dust amount",
			Msg:     *auctionsV2types.NewMsgPlaceMarketBid(bidder, 1, sdk.NewCoin("uasset3", sdk.NewInt(1110000))),
			ExpErr:  auctionsV2types.ErrCannotLeaveDebtLessThanDust,
			ExpResp: nil,
		},
		{
			Name:    "success valid case partial",
			Msg:     *auctionsV2types.NewMsgPlaceMarketBid(bidder, 1, sdk.NewCoin("uasset3", sdk.NewInt(100000))),
			ExpErr:  nil,
			ExpResp: &auctionsV2types.MsgPlaceMarketBidResponse{},
		},
		{
			Name:    "success valid case full",
			Msg:     *auctionsV2types.NewMsgPlaceMarketBid(bidder, 1, sdk.NewCoin("uasset3", sdk.NewInt(1020000))),
			ExpErr:  nil,
			ExpResp: &auctionsV2types.MsgPlaceMarketBidResponse{},
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.Name, func() {

			ctx := sdk.WrapSDKContext(s.ctx)
			resp, err := s.auctionMsgServer.MsgPlaceMarketBid(ctx, &tc.Msg)
			if tc.ExpErr != nil {
				s.Require().Error(err)
				s.Require().EqualError(err, tc.ExpErr.Error())
				s.Require().Equal(tc.ExpResp, resp)
			} else {
				s.Require().NoError(err)
				s.Require().NotNil(resp)
				s.Require().Equal(tc.ExpResp, resp)

			}
		})
	}
}

func (s *KeeperTestSuite) TestPlaceMarketBidForBorrows() {
	s.TestLiquidateBorrows()
	bidder := "cosmos1hm7w7dnvdnra78pz9qxysy7u4tuhc3fnpjmyj7"
	testCases := []struct {
		Name    string
		Msg     auctionsV2types.MsgPlaceMarketBidRequest
		ExpErr  error
		ExpResp *auctionsV2types.MsgPlaceMarketBidResponse
	}{
		{
			Name:    "auction does not exist",
			Msg:     *auctionsV2types.NewMsgPlaceMarketBid(bidder, 10, sdk.NewCoin("uasset2", sdk.NewInt(100000))),
			ExpErr:  sdkerrors.ErrNotFound,
			ExpResp: nil,
		},
		{
			Name:    "dust amount",
			Msg:     *auctionsV2types.NewMsgPlaceMarketBid(bidder, 1, sdk.NewCoin("uasset2", sdk.NewInt(73450000))),
			ExpErr:  auctionsV2types.ErrCannotLeaveDebtLessThanDust,
			ExpResp: nil,
		},
		{
			Name:    "success valid case partial",
			Msg:     *auctionsV2types.NewMsgPlaceMarketBid(bidder, 1, sdk.NewCoin("uasset2", sdk.NewInt(53000000))),
			ExpErr:  nil,
			ExpResp: &auctionsV2types.MsgPlaceMarketBidResponse{},
		},
		{
			Name:    "success valid case full",
			Msg:     *auctionsV2types.NewMsgPlaceMarketBid(bidder, 1, sdk.NewCoin("uasset2", sdk.NewInt(20500000))),
			ExpErr:  nil,
			ExpResp: &auctionsV2types.MsgPlaceMarketBidResponse{},
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.Name, func() {

			ctx := sdk.WrapSDKContext(s.ctx)
			resp, err := s.auctionMsgServer.MsgPlaceMarketBid(ctx, &tc.Msg)
			if tc.ExpErr != nil {
				s.Require().Error(err)
				s.Require().EqualError(err, tc.ExpErr.Error())
				s.Require().Equal(tc.ExpResp, resp)
			} else {
				s.Require().NoError(err)
				s.Require().NotNil(resp)
				s.Require().Equal(tc.ExpResp, resp)

			}
		})
	}
}

func (s *KeeperTestSuite) TestDepositLimitBid() {
	s.TestLiquidateVaults()
	auctionKeeper := &s.keeper
	bidder := "cosmos1hm7w7dnvdnra78pz9qxysy7u4tuhc3fnpjmyj7"

	testCases := []struct {
		Name    string
		Msg     auctionsV2types.MsgDepositLimitBidRequest
		ExpErr  error
		ExpResp *auctionsV2types.MsgDepositLimitBidResponse
	}{
		{
			Name:    "asset does not exist",
			Msg:     *auctionsV2types.NewMsgDepositLimitBid(bidder, 10, 2, sdk.NewInt(2), sdk.NewCoin("uasset2", sdk.NewInt(1000000))),
			ExpErr:  assetTypes.ErrorAssetDoesNotExist,
			ExpResp: nil,
		},
		{
			Name:    "asset does not exist",
			Msg:     *auctionsV2types.NewMsgDepositLimitBid(bidder, 1, 20, sdk.NewInt(2), sdk.NewCoin("uasset2", sdk.NewInt(1000000))),
			ExpErr:  assetTypes.ErrorAssetDoesNotExist,
			ExpResp: nil,
		},
		{
			Name:    "asset denom does not exist",
			Msg:     *auctionsV2types.NewMsgDepositLimitBid(bidder, 1, 2, sdk.NewInt(2), sdk.NewCoin("uasset1", sdk.NewInt(1000000))),
			ExpErr:  auctionsV2types.ErrorUnknownDebtToken,
			ExpResp: nil,
		},
		{
			Name:    "success valid case",
			Msg:     *auctionsV2types.NewMsgDepositLimitBid(bidder, 1, 2, sdk.NewInt(2), sdk.NewCoin("uasset2", sdk.NewInt(1000000))),
			ExpErr:  nil,
			ExpResp: &auctionsV2types.MsgDepositLimitBidResponse{},
		},
	}
	for _, tc := range testCases {
		tc := tc
		s.Run(tc.Name, func() {

			ctx := sdk.WrapSDKContext(s.ctx)
			resp, err := s.auctionMsgServer.MsgDepositLimitBid(ctx, &tc.Msg)
			if tc.ExpErr != nil {
				s.Require().Error(err)
				s.Require().EqualError(err, tc.ExpErr.Error())
				s.Require().Equal(tc.ExpResp, resp)
			} else {
				s.Require().NoError(err)
				s.Require().NotNil(resp)
				s.Require().Equal(tc.ExpResp, resp)

				userLimitAuctionBid, found := auctionKeeper.GetUserLimitBidData(s.ctx, 2, 1, sdk.NewInt(2), bidder)
				s.Require().Equal(found, true)
				userLimitAuctionBidByPremium, found := auctionKeeper.GetUserLimitBidDataByPremium(s.ctx, 2, 1, sdk.NewInt(2))
				s.Require().Equal(found, true)

				s.Require().Equal(userLimitAuctionBid, userLimitAuctionBidByPremium[0])

			}
		})
	}

}

func (s *KeeperTestSuite) TestCancelLimitBid() {
	s.TestDepositLimitBid()
	bidder := "cosmos1hm7w7dnvdnra78pz9qxysy7u4tuhc3fnpjmyj7"

	auctionKeeper := &s.keeper
	id := auctionKeeper.GetLimitAuctionBidID(s.ctx)
	s.Require().Equal(id, uint64(1))
	_, found := auctionKeeper.GetUserLimitBidData(s.ctx, 2, 1, sdk.NewInt(2), bidder)
	s.Require().Equal(found, true)

	testCases := []struct {
		Name    string
		Msg     auctionsV2types.MsgCancelLimitBidRequest
		ExpErr  error
		ExpResp *auctionsV2types.MsgCancelLimitBidResponse
	}{
		{
			Name:    "asset does not exist",
			Msg:     *auctionsV2types.NewMsgCancelLimitBid(bidder, 10, 2, sdk.NewInt(2)),
			ExpErr:  auctionsV2types.ErrBidNotFound,
			ExpResp: nil,
		},
		{
			Name:    "success valid case",
			Msg:     *auctionsV2types.NewMsgCancelLimitBid(bidder, 1, 2, sdk.NewInt(2)),
			ExpErr:  nil,
			ExpResp: &auctionsV2types.MsgCancelLimitBidResponse{},
		},
	}
	for _, tc := range testCases {
		tc := tc
		s.Run(tc.Name, func() {

			ctx := sdk.WrapSDKContext(s.ctx)
			resp, err := s.auctionMsgServer.MsgCancelLimitBid(ctx, &tc.Msg)
			if tc.ExpErr != nil {
				s.Require().Error(err)
				s.Require().EqualError(err, tc.ExpErr.Error())
				s.Require().Equal(tc.ExpResp, resp)
			} else {
				s.Require().NoError(err)
				s.Require().NotNil(resp)
				s.Require().Equal(tc.ExpResp, resp)

				_, found = auctionKeeper.GetUserLimitBidData(s.ctx, 2, 1, sdk.NewInt(2), bidder)
				s.Require().Equal(found, false)
				_, found = auctionKeeper.GetUserLimitBidDataByPremium(s.ctx, 2, 1, sdk.NewInt(2))
				s.Require().Equal(found, false)
			}
		})
	}

}

func (s *KeeperTestSuite) TestWithdrawLimitBid() {
	s.TestDepositLimitBid()
	bidder := "cosmos1hm7w7dnvdnra78pz9qxysy7u4tuhc3fnpjmyj7"

	auctionKeeper := &s.keeper
	id := auctionKeeper.GetLimitAuctionBidID(s.ctx)
	s.Require().Equal(id, uint64(1))
	_, found := auctionKeeper.GetUserLimitBidData(s.ctx, 2, 1, sdk.NewInt(2), bidder)
	s.Require().Equal(found, true)

	testCases := []struct {
		Name    string
		Msg     auctionsV2types.MsgWithdrawLimitBidRequest
		ExpErr  error
		ExpResp *auctionsV2types.MsgWithdrawLimitBidResponse
	}{
		{
			Name:    "asset does not exist",
			Msg:     *auctionsV2types.NewMsgWithdrawLimitBid(bidder, 10, 2, sdk.NewInt(2), sdk.NewCoin("uasset2", sdk.NewInt(500000))),
			ExpErr:  auctionsV2types.ErrBidNotFound,
			ExpResp: nil,
		},
		{
			Name:    "success valid case",
			Msg:     *auctionsV2types.NewMsgWithdrawLimitBid(bidder, 1, 2, sdk.NewInt(2), sdk.NewCoin("uasset2", sdk.NewInt(500000))),
			ExpErr:  nil,
			ExpResp: &auctionsV2types.MsgWithdrawLimitBidResponse{},
		},
	}
	for _, tc := range testCases {
		tc := tc
		s.Run(tc.Name, func() {

			ctx := sdk.WrapSDKContext(s.ctx)
			resp, err := s.auctionMsgServer.MsgWithdrawLimitBid(ctx, &tc.Msg)
			if tc.ExpErr != nil {
				s.Require().Error(err)
				s.Require().EqualError(err, tc.ExpErr.Error())
				s.Require().Equal(tc.ExpResp, resp)
			} else {
				s.Require().NoError(err)
				s.Require().NotNil(resp)
				s.Require().Equal(tc.ExpResp, resp)

				GetUserLimitBidData, found := auctionKeeper.GetUserLimitBidData(s.ctx, 2, 1, sdk.NewInt(2), bidder)
				s.Require().Equal(found, true)
				s.Require().Equal(GetUserLimitBidData.DebtToken.Amount, sdk.NewInt(500000))

				_, found = auctionKeeper.GetUserLimitBidDataByPremium(s.ctx, 2, 1, sdk.NewInt(2))
				s.Require().Equal(found, true)
			}
		})
	}
}

func (s *KeeperTestSuite) TestDebtActivator() {
	collectorKeeper := &s.collectorKeeper
	liquidationKeeper := &s.liquidationKeeper
	s.AddAppAssets()
	s.WasmSetCollectorLookupTableAndAuctionControlForDebt()

	err := collectorKeeper.SetNetFeeCollectedData(s.ctx, uint64(2), 2, sdk.NewIntFromUint64(4700000))
	s.Require().NoError(err)
	k, ctx := &s.liquidationKeeper, &s.ctx
	err = k.Liquidate(*ctx)
	s.Require().NoError(err)
	lockedVault := liquidationKeeper.GetLockedVaults(s.ctx)
	s.Require().Equal(lockedVault[0].OriginalVaultId, uint64(0))
	s.Require().Equal(lockedVault[0].ExtendedPairId, uint64(0))
	s.Require().Equal(lockedVault[0].Owner, "")
	s.Require().Equal(lockedVault[0].CollateralAssetId, uint64(2))
	s.Require().Equal(lockedVault[0].DebtAssetId, uint64(3))
	s.Require().Equal(lockedVault[0].InitiatorType, "debt")

}

func (s *KeeperTestSuite) TestSurplusActivator() {
	collectorKeeper := &s.collectorKeeper
	liquidationKeeper := &s.liquidationKeeper
	s.AddAppAssets()
	s.WasmSetCollectorLookupTableAndAuctionControlForSurplus()
	err := s.app.BankKeeper.MintCoins(s.ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin("uasset2", sdk.NewInt(10000000))))
	s.Require().NoError(err)
	err = s.app.BankKeeper.SendCoinsFromModuleToModule(s.ctx, types.ModuleName, collectortypes.ModuleName, sdk.NewCoins(sdk.NewCoin("uasset2", sdk.NewInt(10000000))))
	s.Require().NoError(err)

	err = collectorKeeper.SetNetFeeCollectedData(s.ctx, uint64(2), 2, sdk.NewIntFromUint64(100000000))
	s.Require().NoError(err)
	k, ctx := &s.liquidationKeeper, &s.ctx
	err = k.Liquidate(*ctx)
	s.Require().NoError(err)
	lockedVault := liquidationKeeper.GetLockedVaults(s.ctx)
	s.Require().Equal(lockedVault[0].OriginalVaultId, uint64(0))
	s.Require().Equal(lockedVault[0].ExtendedPairId, uint64(0))
	s.Require().Equal(lockedVault[0].Owner, "")
	s.Require().Equal(lockedVault[0].CollateralAssetId, uint64(2))
	s.Require().Equal(lockedVault[0].DebtAssetId, uint64(3))
	s.Require().Equal(lockedVault[0].InitiatorType, "surplus")
}

func (s *KeeperTestSuite) TestDebtAuctionBid() {
	s.TestDebtActivator()
	bidder := "cosmos1hm7w7dnvdnra78pz9qxysy7u4tuhc3fnpjmyj7"

	testCases := []struct {
		Name    string
		Msg     auctionsV2types.MsgPlaceMarketBidRequest
		ExpErr  error
		ExpResp *auctionsV2types.MsgPlaceMarketBidResponse
	}{
		{
			Name:    "auction does not exist",
			Msg:     *auctionsV2types.NewMsgPlaceMarketBid(bidder, 10, sdk.NewCoin("uasset2", sdk.NewInt(100000))),
			ExpErr:  sdkerrors.ErrNotFound,
			ExpResp: nil,
		},
		{
			Name:    "bidding amount is greater than maximum bidding amount",
			Msg:     *auctionsV2types.NewMsgPlaceMarketBid(bidder, 1, sdk.NewCoin("uasset3", sdk.NewInt(53000000))),
			ExpErr:  auctionsV2types.ErrorMaxBidAmount,
			ExpResp: nil,
		},
		{
			Name:    "success valid case",
			Msg:     *auctionsV2types.NewMsgPlaceMarketBid(bidder, 1, sdk.NewCoin("uasset3", sdk.NewInt(200000))),
			ExpErr:  nil,
			ExpResp: &auctionsV2types.MsgPlaceMarketBidResponse{},
		},
		{
			Name:    "bid should be less than or equal to 180000",
			Msg:     *auctionsV2types.NewMsgPlaceMarketBid(bidder, 1, sdk.NewCoin("uasset3", sdk.NewInt(200000))),
			ExpErr:  fmt.Errorf("bid should be less than or equal to 180000 : not found"),
			ExpResp: nil,
		},
		{
			Name:    "success valid case 2",
			Msg:     *auctionsV2types.NewMsgPlaceMarketBid(bidder, 1, sdk.NewCoin("uasset3", sdk.NewInt(180000))),
			ExpErr:  nil,
			ExpResp: &auctionsV2types.MsgPlaceMarketBidResponse{},
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.Name, func() {

			ctx := sdk.WrapSDKContext(s.ctx)
			resp, err := s.auctionMsgServer.MsgPlaceMarketBid(ctx, &tc.Msg)
			if tc.ExpErr != nil {
				s.Require().Error(err)
				s.Require().EqualError(err, tc.ExpErr.Error())
				s.Require().Equal(tc.ExpResp, resp)
			} else {
				s.Require().NoError(err)
				s.Require().NotNil(resp)
				s.Require().Equal(tc.ExpResp, resp)

			}
		})
	}

	//auctionKeeper := &s.keeper
	//debtAuction := auctionKeeper.GetAuctions(s.ctx)
	//fmt.Println("debtAuction", debtAuction)
	//fmt.Println("debtAuction[0].AuctionId", debtAuction[0].AuctionId)
	//fmt.Println("debtAuction[0].CollateralToken", debtAuction[0].CollateralToken)
	//fmt.Println("debtAuction[0].DebtToken", debtAuction[0].DebtToken)
	//fmt.Println("debtAuction[0].ActiveBiddingId", debtAuction[0].ActiveBiddingId)
	//fmt.Println("debtAuction[0].BiddingIds", debtAuction[0].BiddingIds)
	//fmt.Println("debtAuction[0].CollateralTokenAuctionPrice", debtAuction[0].CollateralTokenAuctionPrice)
	//fmt.Println("debtAuction[0].CollateralTokenOraclePrice", debtAuction[0].CollateralTokenOraclePrice)
	//fmt.Println("debtAuction[0].DebtTokenOraclePrice", debtAuction[0].DebtTokenOraclePrice)
	//fmt.Println("debtAuction[0].LockedVaultId", debtAuction[0].LockedVaultId)
	//fmt.Println("debtAuction[0].StartTime", debtAuction[0].StartTime)
	//fmt.Println("debtAuction[0].EndTime", debtAuction[0].EndTime)
	//fmt.Println("debtAuction[0].AppId", debtAuction[0].AppId)
	//fmt.Println("debtAuction[0].AuctionType", debtAuction[0].AuctionType)
	//fmt.Println("debtAuction[0].CollateralAssetId", debtAuction[0].CollateralAssetId)
	//fmt.Println("debtAuction[0].DebtAssetId", debtAuction[0].DebtAssetId)
	//fmt.Println("debtAuction[0].BonusAmount", debtAuction[0].BonusAmount)

}

func (s *KeeperTestSuite) TestSurplusAuctionBid() {
	s.TestSurplusActivator()
	bidder := "cosmos1hm7w7dnvdnra78pz9qxysy7u4tuhc3fnpjmyj7"

	testCases := []struct {
		Name    string
		Msg     auctionsV2types.MsgPlaceMarketBidRequest
		ExpErr  error
		ExpResp *auctionsV2types.MsgPlaceMarketBidResponse
	}{
		{
			Name:    "auction does not exist",
			Msg:     *auctionsV2types.NewMsgPlaceMarketBid(bidder, 10, sdk.NewCoin("uasset3", sdk.NewInt(100000))),
			ExpErr:  sdkerrors.ErrNotFound,
			ExpResp: nil,
		},
		{
			Name:    "success valid case",
			Msg:     *auctionsV2types.NewMsgPlaceMarketBid(bidder, 1, sdk.NewCoin("uasset3", sdk.NewInt(200000))),
			ExpErr:  nil,
			ExpResp: &auctionsV2types.MsgPlaceMarketBidResponse{},
		},
		{
			Name:    "bid should be greater than or equal to 220000",
			Msg:     *auctionsV2types.NewMsgPlaceMarketBid(bidder, 1, sdk.NewCoin("uasset3", sdk.NewInt(200000))),
			ExpErr:  fmt.Errorf("bid should be greater than or equal to 220000 : not found"),
			ExpResp: nil,
		},
		{
			Name:    "success valid case 2",
			Msg:     *auctionsV2types.NewMsgPlaceMarketBid(bidder, 1, sdk.NewCoin("uasset3", sdk.NewInt(220000))),
			ExpErr:  nil,
			ExpResp: &auctionsV2types.MsgPlaceMarketBidResponse{},
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.Name, func() {

			ctx := sdk.WrapSDKContext(s.ctx)
			resp, err := s.auctionMsgServer.MsgPlaceMarketBid(ctx, &tc.Msg)
			if tc.ExpErr != nil {
				s.Require().Error(err)
				s.Require().EqualError(err, tc.ExpErr.Error())
				s.Require().Equal(tc.ExpResp, resp)
			} else {
				s.Require().NoError(err)
				s.Require().NotNil(resp)
				s.Require().Equal(tc.ExpResp, resp)

			}
		})
	}

}

func (s *KeeperTestSuite) TestAuctionIterator() {
	s.ctx = s.ctx.WithBlockTime(utils.ParseTime("2023-06-01T12:00:00Z"))
	s.TestLiquidateVaults()

	// debt auction
	collectorKeeper := &s.collectorKeeper
	s.WasmSetCollectorLookupTableAndAuctionControlForDebt()
	err := collectorKeeper.SetNetFeeCollectedData(s.ctx, uint64(2), 2, sdk.NewIntFromUint64(4700000))
	s.Require().NoError(err)
	k := &s.liquidationKeeper
	err = k.Liquidate(s.ctx)
	s.Require().NoError(err)

	auction, _ := s.keeper.GetAuction(s.ctx, 1)
	s.Require().Equal(auction.CollateralTokenAuctionPrice, sdk.NewDecFromInt(sdk.NewInt(1200000)))

	// update auction price
	s.ctx = s.ctx.WithBlockTime(utils.ParseTime("2023-06-01T12:30:00Z"))
	auctionsV2.BeginBlocker(s.ctx, s.keeper)
	auction, _ = s.keeper.GetAuction(s.ctx, 1)
	s.Require().Equal(auction.CollateralTokenAuctionPrice, sdk.NewDecFromInt(sdk.NewInt(1020000)))

	s.ctx = s.ctx.WithBlockTime(utils.ParseTime("2023-06-01T12:59:59Z"))
	auctionsV2.BeginBlocker(s.ctx, s.keeper)
	auction, _ = s.keeper.GetAuction(s.ctx, 1)
	s.Require().Equal(auction.CollateralTokenAuctionPrice, sdk.NewDecFromInt(sdk.NewInt(840100)))

	// restart auction
	s.ctx = s.ctx.WithBlockTime(utils.ParseTime("2023-06-01T13:30:00Z"))
	auctionsV2.BeginBlocker(s.ctx, s.keeper)
	auction, _ = s.keeper.GetAuction(s.ctx, 1)
	s.Require().Equal(auction.CollateralTokenAuctionPrice, sdk.NewDecFromInt(sdk.NewInt(1200000)))
}

func (s *KeeperTestSuite) TestAuctionIteratorSurplus() {
	s.ctx = s.ctx.WithBlockTime(utils.ParseTime("2023-06-01T12:00:00Z"))
	s.TestLiquidateVaults()
	//winnerAddress := "cosmos1hm7w7dnvdnra78pz9qxysy7u4tuhc3fnpjmyj7"

	collectorKeeper := &s.collectorKeeper
	k := &s.liquidationKeeper
	// surplus auction

	//tokenmintKeeper := &s.tokenmintKeeper
	//server := tokenmintKeeper1.NewMsgServer(*tokenmintKeeper)
	//msg1 := tokenminttypes.MsgMintNewTokensRequest{
	//	From:    winnerAddress,
	//	AppId:   1,
	//	AssetId: 3,
	//}
	//_, err := server.MsgMintNewTokens(sdk.WrapSDKContext(s.ctx), &msg1)
	//s.Require().NoError(err)

	s.WasmSetCollectorLookupTableAndAuctionControlForSurplus()
	err := s.app.BankKeeper.MintCoins(s.ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin("uasset2", sdk.NewInt(10000000))))
	s.Require().NoError(err)
	err = s.app.BankKeeper.SendCoinsFromModuleToModule(s.ctx, types.ModuleName, collectortypes.ModuleName, sdk.NewCoins(sdk.NewCoin("uasset2", sdk.NewInt(10000000))))
	s.Require().NoError(err)
	err = collectorKeeper.SetNetFeeCollectedData(s.ctx, uint64(2), 2, sdk.NewIntFromUint64(100000000))
	s.Require().NoError(err)
	k, ctx := &s.liquidationKeeper, &s.ctx
	err = k.Liquidate(*ctx)
	s.Require().NoError(err)

	auction, _ := s.keeper.GetAuction(s.ctx, 1)
	s.Require().Equal(auction.CollateralTokenAuctionPrice, sdk.NewDecFromInt(sdk.NewInt(1200000)))

	// update auction price and bid for english auction

	s.ctx = s.ctx.WithBlockTime(utils.ParseTime("2023-06-01T12:30:00Z"))
	auctionsV2.BeginBlocker(s.ctx, s.keeper)

	msg := auctionsV2types.MsgPlaceMarketBidRequest{
		AuctionId: 3,
		Bidder:    "cosmos1hm7w7dnvdnra78pz9qxysy7u4tuhc3fnpjmyj7",
		Amount:    sdk.NewCoin("uasset3", sdk.NewInt(220000)),
	}
	_, err = s.auctionMsgServer.MsgPlaceMarketBid(sdk.WrapSDKContext(s.ctx), &msg)
	s.Require().NoError(err)

	// restart auction
	s.ctx = s.ctx.WithBlockTime(utils.ParseTime("2023-06-01T13:30:00Z"))
	auctionsV2.BeginBlocker(s.ctx, s.keeper)
	auction, _ = s.keeper.GetAuction(s.ctx, 1)
	s.Require().Equal(auction.CollateralTokenAuctionPrice, sdk.NewDecFromInt(sdk.NewInt(1200000)))
}

func (s *KeeperTestSuite) TestLimitBid() {
	s.ctx = s.ctx.WithBlockTime(utils.ParseTime("2023-06-01T12:00:00Z"))
	s.TestLiquidateVaults()
	auctions := s.app.NewaucKeeper.GetAuctions(s.ctx)
	s.Require().Equal(len(auctions), 2)

	bidder := "cosmos1hm7w7dnvdnra78pz9qxysy7u4tuhc3fnpjmyj7"

	auctionKeeper := &s.keeper
	liquidationKeeper := &s.liquidationKeeper
	err := liquidationKeeper.MsgAppReserveFundsFn(s.ctx, bidder, 2, 3, sdk.NewCoin("uasset3", sdk.NewInt(5990000)))
	s.Require().NoError(err)

	err = auctionKeeper.DepositLimitAuctionBid(s.ctx, bidder, 2, 3, sdk.NewInt(9), sdk.NewCoin("uasset3", sdk.NewInt(7000000)))
	s.Require().NoError(err)

	a, _ := auctionKeeper.GetUserLimitBidDataByPremium(s.ctx, 3, 2, sdk.NewInt(9))
	s.Require().Equal(len(a), 1)

	s.ctx = s.ctx.WithBlockTime(utils.ParseTime("2023-06-01T12:49:00Z"))
	auctionsV2.BeginBlocker(s.ctx, s.keeper)

	a, _ = auctionKeeper.GetUserLimitBidDataByPremium(s.ctx, 3, 2, sdk.NewInt(9))
	s.Require().Equal(len(a), 1)

	auctions = s.app.NewaucKeeper.GetAuctions(s.ctx)
	s.Require().Equal(len(auctions), 0)
}
