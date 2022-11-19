package keeper_test

import (
	"fmt"
	"github.com/petrichormoney/petri/app/wasm/bindings"
	utils "github.com/petrichormoney/petri/types"
	assetTypes "github.com/petrichormoney/petri/x/asset/types"
	lockerkeeper "github.com/petrichormoney/petri/x/locker/keeper"
	lockertypes "github.com/petrichormoney/petri/x/locker/types"
	markettypes "github.com/petrichormoney/petri/x/market/types"
	"github.com/petrichormoney/petri/x/rewards"
	keeper "github.com/petrichormoney/petri/x/rewards/keeper"
	"github.com/petrichormoney/petri/x/rewards/types"
	vaultkeeper "github.com/petrichormoney/petri/x/vault/keeper"
	vaulttypes "github.com/petrichormoney/petri/x/vault/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func (s *KeeperTestSuite) AddAppAsset() {
	assetKeeper, ctx := &s.assetKeeper, &s.ctx
	msg1 := assetTypes.AppData{
		Name:             "cswap",
		ShortName:        "cswap",
		MinGovDeposit:    sdk.NewIntFromUint64(10000000),
		GovTimeInSeconds: 900,
	}
	err := assetKeeper.AddAppRecords(*ctx, msg1)
	s.Require().NoError(err)

	msg2 := assetTypes.AppData{
		Name:             "commodo",
		ShortName:        "comdo",
		MinGovDeposit:    sdk.NewIntFromUint64(10000000),
		GovTimeInSeconds: 900,
	}
	err = assetKeeper.AddAppRecords(*ctx, msg2)
	s.Require().NoError(err)

	msg3 := assetTypes.Asset{
		Name:          "CMDX",
		Denom:         "upetri",
		Decimals:      sdk.NewInt(1000000),
		IsOnChain:     true,
		IsCdpMintable: true,
	}

	err = assetKeeper.AddAssetRecords(*ctx, msg3)
	s.Require().NoError(err)
	market1 := markettypes.TimeWeightedAverage{
		AssetID:       1,
		ScriptID:      12,
		Twa:           1000000,
		CurrentIndex:  0,
		IsPriceActive: true,
		PriceValue:    []uint64{1000000},
	}
	s.app.MarketKeeper.SetTwa(s.ctx, market1)
	_, err = s.app.MarketKeeper.GetLatestPrice(s.ctx, 1)
	s.Suite.NoError(err)

	msg4 := assetTypes.Asset{
		Name:          "FUST",
		Denom:         "ucmst",
		Decimals:      sdk.NewInt(1000000),
		IsOnChain:     true,
		IsCdpMintable: true,
	}
	err = assetKeeper.AddAssetRecords(*ctx, msg4)
	s.Require().NoError(err)

	market2 := markettypes.TimeWeightedAverage{
		AssetID:       2,
		ScriptID:      12,
		Twa:           1000000,
		CurrentIndex:  0,
		IsPriceActive: true,
		PriceValue:    []uint64{1000000},
	}
	s.app.MarketKeeper.SetTwa(s.ctx, market2)
	_, err = s.app.MarketKeeper.GetLatestPrice(s.ctx, 2)
	s.Suite.NoError(err)

	msg5 := assetTypes.Asset{
		Name:      "HARBOR",
		Denom:     "uharbor",
		Decimals:  sdk.NewInt(1000000),
		IsOnChain: true,
	}
	err = assetKeeper.AddAssetRecords(*ctx, msg5)
	s.Require().NoError(err)

	market3 := markettypes.TimeWeightedAverage{
		AssetID:       3,
		ScriptID:      12,
		Twa:           1000000,
		CurrentIndex:  0,
		IsPriceActive: true,
		PriceValue:    []uint64{1000000},
	}
	s.app.MarketKeeper.SetTwa(s.ctx, market3)
	_, err = s.app.MarketKeeper.GetLatestPrice(s.ctx, 3)
	s.Suite.NoError(err)
}

func (s *KeeperTestSuite) AddCollectorLookupTable() {
	collectorKeeper, ctx := &s.collector, &s.ctx
	msg1 := bindings.MsgSetCollectorLookupTable{
		AppID:            1,
		CollectorAssetID: 1,
		SecondaryAssetID: 3,
		SurplusThreshold: sdk.NewInt(10000000),
		DebtThreshold:    sdk.NewInt(5000000),
		LockerSavingRate: sdk.MustNewDecFromStr("0.1"),
		LotSize:          sdk.NewInt(2000000),
		BidFactor:        sdk.MustNewDecFromStr("0.01"),
		DebtLotSize:      sdk.NewInt(2000000),
	}
	err := collectorKeeper.WasmSetCollectorLookupTable(*ctx, &msg1)
	s.Require().NoError(err)

	msg2 := bindings.MsgSetCollectorLookupTable{
		AppID:            1,
		CollectorAssetID: 2,
		SecondaryAssetID: 3,
		SurplusThreshold: sdk.NewInt(10000000),
		DebtThreshold:    sdk.NewInt(5000000),
		LockerSavingRate: sdk.MustNewDecFromStr("0.1"),
		LotSize:          sdk.NewInt(2000000),
		BidFactor:        sdk.MustNewDecFromStr("0.01"),
		DebtLotSize:      sdk.NewInt(2000000),
	}
	err1 := collectorKeeper.WasmSetCollectorLookupTable(*ctx, &msg2)
	s.Require().NoError(err1)

	msg3 := bindings.MsgSetCollectorLookupTable{
		AppID:            2,
		CollectorAssetID: 1,
		SecondaryAssetID: 3,
		SurplusThreshold: sdk.NewInt(10000000),
		DebtThreshold:    sdk.NewInt(5000000),
		LockerSavingRate: sdk.MustNewDecFromStr("0.1"),
		LotSize:          sdk.NewInt(2000000),
		BidFactor:        sdk.MustNewDecFromStr("0.01"),
		DebtLotSize:      sdk.NewInt(2000000),
	}
	err2 := collectorKeeper.WasmSetCollectorLookupTable(*ctx, &msg3)
	s.Require().NoError(err2)
}

func (s *KeeperTestSuite) TestCreateLocker() {
	s.ctx = s.ctx.WithBlockTime(utils.ParseTime("2022-03-01T12:00:00Z"))
	s.ctx = s.ctx.WithBlockHeight(10)
	userAddress := "cosmos1q7q90qsl9g0gl2zz0njxwv2a649yqrtyxtnv3v"
	s.AddAppAsset()
	s.AddCollectorLookupTable()
	lockerKeeper, ctx := &s.lockerKeeper, &s.ctx
	server := lockerkeeper.NewMsgServer(*lockerKeeper)
	for _, tc := range []struct {
		name string
		msg  lockertypes.MsgAddWhiteListedAssetRequest
	}{
		{
			"Whitelist : App1 Asset 1",
			lockertypes.MsgAddWhiteListedAssetRequest{
				From:    userAddress,
				AppId:   1,
				AssetId: 1,
			},
		},
		{
			"Whitelist : App1 Asset 2",
			lockertypes.MsgAddWhiteListedAssetRequest{
				From:    userAddress,
				AppId:   1,
				AssetId: 2,
			},
		},
		{
			"Whitelist : App2 Asset 1",
			lockertypes.MsgAddWhiteListedAssetRequest{
				From:    userAddress,
				AppId:   2,
				AssetId: 1,
			},
		},
	} {
		s.Run(tc.name, func() {
			_, err := lockerKeeper.AddWhiteListedAsset(*ctx, &tc.msg)
			s.Require().NoError(err)
		})
	}
	msg2 := lockertypes.MsgCreateLockerRequest{
		Depositor: userAddress,
		Amount:    sdk.NewInt(1000000000),
		AssetId:   1,
		AppId:     1,
	}

	s.fundAddr(userAddress, sdk.NewCoin("upetri", sdk.NewIntFromUint64(1000000000)))
	_, err := server.MsgCreateLocker(sdk.WrapSDKContext(*ctx), &msg2)
	s.Require().NoError(err)
}

func (s *KeeperTestSuite) TestCreateExtRewardsLocker() {
	s.ctx = s.ctx.WithBlockTime(utils.ParseTime("2022-03-01T12:00:00Z"))
	s.ctx = s.ctx.WithBlockHeight(10)
	userAddress := "cosmos1q7q90qsl9g0gl2zz0njxwv2a649yqrtyxtnv3v"
	amt, _ := sdk.NewIntFromString("1000000000000000000000")
	s.fundAddr(userAddress, sdk.NewCoin("weth", amt))

	s.TestCreateLocker()
	rewardsKeeper, ctx := &s.rewardsKeeper, &s.ctx
	server := keeper.NewMsgServerImpl(*rewardsKeeper)
	for _, tc := range []struct {
		name          string
		msg           types.ActivateExternalRewardsLockers
		expectedError bool
		ExpErr        error
	}{
		{
			"ActivateExternalRewardsLockers : success",
			types.ActivateExternalRewardsLockers{
				AppMappingId:         1,
				AssetId:              1,
				TotalRewards:         sdk.NewCoin("weth", amt),
				DurationDays:         5,
				Depositor:            userAddress,
				MinLockupTimeSeconds: 0,
			},
			false,
			nil,
		},
	} {
		s.Run(tc.name, func() {

			_, err := server.ExternalRewardsLockers(sdk.WrapSDKContext(*ctx), &tc.msg)
			if tc.ExpErr != nil {
				s.Require().Error(err)
				s.Require().EqualError(err, tc.ExpErr.Error())
			} else {
				s.Require().NoError(err)
				availableBalances := s.getBalances(sdk.MustAccAddressFromBech32(userAddress))
				fmt.Println("bal when created ext rewards", availableBalances)
			}
		})
	}
	s.ctx = s.ctx.WithBlockTime(utils.ParseTime("2022-03-02T12:10:00Z"))
	s.ctx = s.ctx.WithBlockHeight(11)
	req := abci.RequestBeginBlock{}
	rewards.BeginBlocker(*ctx, req, *rewardsKeeper)
	availableBalances := s.getBalances(sdk.MustAccAddressFromBech32(userAddress))
	fmt.Println("bal at first day", availableBalances)
	s.ctx = s.ctx.WithBlockTime(utils.ParseTime("2022-03-03T12:11:00Z"))
	s.ctx = s.ctx.WithBlockHeight(12)
	rewards.BeginBlocker(*ctx, req, *rewardsKeeper)
	availableBalances = s.getBalances(sdk.MustAccAddressFromBech32(userAddress))
	fmt.Println("bal at second day", availableBalances)
}

func (s *KeeperTestSuite) TestCreateVault() {
	s.ctx = s.ctx.WithBlockTime(utils.ParseTime("2022-03-01T12:00:00Z"))
	s.ctx = s.ctx.WithBlockHeight(10)
	userAddress := "cosmos1q7q90qsl9g0gl2zz0njxwv2a649yqrtyxtnv3v"
	userAddress1 := "cosmos1kwtdrjkwu6y87vlylaeatzmc5p4jhvn7qwqnkp"

	addr1 := s.addr(1)
	s.AddAppAsset()

	pairID := s.CreateNewPair(addr1, 1, 2)
	extendedVaultPairID1 := s.CreateNewExtendedVaultPair("CMDX-C", 1, pairID, false, true)

	vaultKeeper, ctx := &s.vaultKeeper, &s.ctx
	server := vaultkeeper.NewMsgServer(*vaultKeeper)

	msg2 := vaulttypes.MsgCreateRequest{
		From:                userAddress,
		AppId:               1,
		ExtendedPairVaultId: extendedVaultPairID1,
		AmountIn:            sdk.NewInt(1000000000),
		AmountOut:           sdk.NewInt(200000000),
	}

	s.fundAddr(userAddress, sdk.NewCoin("upetri", sdk.NewIntFromUint64(1000000000)))
	_, err := server.MsgCreate(sdk.WrapSDKContext(*ctx), &msg2)
	s.Require().NoError(err)

	msg3 := vaulttypes.MsgCreateRequest{
		From:                userAddress1,
		AppId:               1,
		ExtendedPairVaultId: extendedVaultPairID1,
		AmountIn:            sdk.NewInt(1000000000),
		AmountOut:           sdk.NewInt(100000000),
	}

	s.fundAddr(userAddress1, sdk.NewCoin("upetri", sdk.NewIntFromUint64(1000000000)))
	_, err = server.MsgCreate(sdk.WrapSDKContext(*ctx), &msg3)
	s.Require().NoError(err)
}

func (s *KeeperTestSuite) TestCreateExtRewardsVault() {
	s.ctx = s.ctx.WithBlockTime(utils.ParseTime("2022-03-01T12:00:00Z"))
	s.ctx = s.ctx.WithBlockHeight(10)
	userAddress := "cosmos1q7q90qsl9g0gl2zz0njxwv2a649yqrtyxtnv3v"
	amt, _ := sdk.NewIntFromString("1000000000000000000000")
	s.fundAddr(userAddress, sdk.NewCoin("btc", amt))

	s.TestCreateVault()
	rewardsKeeper, ctx := &s.rewardsKeeper, &s.ctx
	server := keeper.NewMsgServerImpl(*rewardsKeeper)
	for _, tc := range []struct {
		name          string
		msg           types.ActivateExternalRewardsVault
		expectedError bool
		ExpErr        error
	}{
		{
			"ActivateExternalRewardsLockers : success",
			types.ActivateExternalRewardsVault{
				AppMappingId:         1,
				ExtendedPairId:       1,
				TotalRewards:         sdk.NewCoin("btc", amt),
				DurationDays:         3,
				Depositor:            userAddress,
				MinLockupTimeSeconds: 0,
			},
			false,
			nil,
		},
	} {
		s.Run(tc.name, func() {

			_, err := server.ExternalRewardsVault(sdk.WrapSDKContext(*ctx), &tc.msg)
			if tc.ExpErr != nil {
				s.Require().Error(err)
				s.Require().EqualError(err, tc.ExpErr.Error())
			} else {
				s.Require().NoError(err)
				availableBalances := s.getBalances(sdk.MustAccAddressFromBech32(userAddress))
				fmt.Println("bal when created ext rewards", availableBalances)
			}
		})
	}
	userAddress1 := "cosmos1kwtdrjkwu6y87vlylaeatzmc5p4jhvn7qwqnkp"
	s.ctx = s.ctx.WithBlockTime(utils.ParseTime("2022-03-02T12:10:00Z"))
	s.ctx = s.ctx.WithBlockHeight(11)
	req := abci.RequestBeginBlock{}
	rewards.BeginBlocker(*ctx, req, *rewardsKeeper)
	availableBalances := s.getBalances(sdk.MustAccAddressFromBech32(userAddress))
	fmt.Println("bal at first day", availableBalances)
	availableBalances1 := s.getBalances(sdk.MustAccAddressFromBech32(userAddress1))
	fmt.Println("bal at first day second user", availableBalances1)
	s.ctx = s.ctx.WithBlockTime(utils.ParseTime("2022-03-03T12:11:00Z"))
	s.ctx = s.ctx.WithBlockHeight(12)
	rewards.BeginBlocker(*ctx, req, *rewardsKeeper)
	availableBalances = s.getBalances(sdk.MustAccAddressFromBech32(userAddress))
	fmt.Println("bal at second day", availableBalances)
	rewards.BeginBlocker(*ctx, req, *rewardsKeeper)
	availableBalances1 = s.getBalances(sdk.MustAccAddressFromBech32(userAddress1))
	fmt.Println("bal at second day second user", availableBalances1)
}
