package keeper_test

import (
	"fmt"
	"github.com/comdex-official/comdex/app/wasm/bindings"
	utils "github.com/comdex-official/comdex/types"
	assetTypes "github.com/comdex-official/comdex/x/asset/types"
	lendkeeper "github.com/comdex-official/comdex/x/lend/keeper"
	lendtypes "github.com/comdex-official/comdex/x/lend/types"
	lockerkeeper "github.com/comdex-official/comdex/x/locker/keeper"
	lockertypes "github.com/comdex-official/comdex/x/locker/types"
	markettypes "github.com/comdex-official/comdex/x/market/types"
	"github.com/comdex-official/comdex/x/rewards"
	keeper "github.com/comdex-official/comdex/x/rewards/keeper"
	"github.com/comdex-official/comdex/x/rewards/types"
	vaultkeeper "github.com/comdex-official/comdex/x/vault/keeper"
	vaulttypes "github.com/comdex-official/comdex/x/vault/types"
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
		Denom:         "ucmdx",
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
		Name:          "CMST",
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

func Dec(s string) sdk.Dec {
	dec, err := sdk.NewDecFromStr(s)
	if err != nil {
		panic(err)
	}
	return dec
}

func (s *KeeperTestSuite) AddAppAssetLend() {
	assetKeeper, ctx := &s.assetKeeper, &s.ctx
	lendKeeper := &s.lendKeeper
	msg1 := assetTypes.AppData{
		Name:             "cswap",
		ShortName:        "cswap",
		MinGovDeposit:    sdk.NewIntFromUint64(10000000),
		GovTimeInSeconds: 900,
	}
	err := assetKeeper.AddAppRecords(*ctx, msg1)
	s.Require().NoError(err)

	msg2 := assetTypes.AppData{
		Name:             "harbor",
		ShortName:        "harbor",
		MinGovDeposit:    sdk.NewIntFromUint64(10000000),
		GovTimeInSeconds: 900,
	}
	err = assetKeeper.AddAppRecords(*ctx, msg2)
	s.Require().NoError(err)
	msg3 := assetTypes.AppData{
		Name:             "commodo",
		ShortName:        "comdo",
		MinGovDeposit:    sdk.NewIntFromUint64(10000000),
		GovTimeInSeconds: 900,
	}
	err = assetKeeper.AddAppRecords(*ctx, msg3)
	s.Require().NoError(err)

	msg4 := assetTypes.Asset{
		Name:          "ATOM",
		Denom:         "uatom",
		Decimals:      sdk.NewInt(1000000),
		IsOnChain:     true,
		IsCdpMintable: true,
	}

	err = assetKeeper.AddAssetRecords(*ctx, msg4)
	s.Require().NoError(err)
	market1 := markettypes.TimeWeightedAverage{
		AssetID:       1,
		ScriptID:      12,
		Twa:           12000000,
		CurrentIndex:  0,
		IsPriceActive: true,
		PriceValue:    []uint64{12000000},
	}
	s.app.MarketKeeper.SetTwa(s.ctx, market1)
	_, err = s.app.MarketKeeper.GetLatestPrice(s.ctx, 1)
	s.Suite.NoError(err)

	msg5 := assetTypes.Asset{
		Name:          "CMDX",
		Denom:         "ucmdx",
		Decimals:      sdk.NewInt(1000000),
		IsOnChain:     true,
		IsCdpMintable: true,
	}

	err = assetKeeper.AddAssetRecords(*ctx, msg5)
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

	msg6 := assetTypes.Asset{
		Name:          "CMST",
		Denom:         "ucmst",
		Decimals:      sdk.NewInt(1000000),
		IsOnChain:     true,
		IsCdpMintable: true,
	}
	err = assetKeeper.AddAssetRecords(*ctx, msg6)
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

	msg7 := assetTypes.Asset{
		Name:      "HARBOR",
		Denom:     "uharbor",
		Decimals:  sdk.NewInt(1000000),
		IsOnChain: true,
	}
	err = assetKeeper.AddAssetRecords(*ctx, msg7)
	s.Require().NoError(err)

	market4 := markettypes.TimeWeightedAverage{
		AssetID:       4,
		ScriptID:      12,
		Twa:           1000000,
		CurrentIndex:  0,
		IsPriceActive: true,
		PriceValue:    []uint64{1000000},
	}
	s.app.MarketKeeper.SetTwa(s.ctx, market4)
	_, err = s.app.MarketKeeper.GetLatestPrice(s.ctx, 4)
	s.Suite.NoError(err)

	msg11 := assetTypes.Asset{
		Name:      "CATOM",
		Denom:     "ucatom",
		Decimals:  sdk.NewInt(1000000),
		IsOnChain: true,
	}
	err = assetKeeper.AddAssetRecords(*ctx, msg11)
	s.Require().NoError(err)

	msg12 := assetTypes.Asset{
		Name:      "CCMDX",
		Denom:     "uccmdx",
		Decimals:  sdk.NewInt(1000000),
		IsOnChain: true,
	}
	err = assetKeeper.AddAssetRecords(*ctx, msg12)
	s.Require().NoError(err)

	msg13 := assetTypes.Asset{
		Name:      "CCMST",
		Denom:     "uccmst",
		Decimals:  sdk.NewInt(1000000),
		IsOnChain: true,
	}
	err = assetKeeper.AddAssetRecords(*ctx, msg13)
	s.Require().NoError(err)

	cmstRatesParams := lendtypes.AssetRatesParams{
		AssetID:              3,
		UOptimal:             Dec("0.8"),
		Base:                 Dec("0.002"),
		Slope1:               Dec("0.06"),
		Slope2:               Dec("0.6"),
		EnableStableBorrow:   false,
		StableBase:           Dec("0.0"),
		StableSlope1:         Dec("0.0"),
		StableSlope2:         Dec("0.0"),
		Ltv:                  Dec("0.8"),
		LiquidationThreshold: Dec("0.85"),
		LiquidationPenalty:   Dec("0.025"),
		LiquidationBonus:     Dec("0.025"),
		ReserveFactor:        Dec("0.1"),
		CAssetID:             7,
	}
	lendKeeper.SetAssetRatesParams(s.ctx, cmstRatesParams)
	atomRatesParams := lendtypes.AssetRatesParams{
		AssetID:              1,
		UOptimal:             Dec("0.75"),
		Base:                 Dec("0.002"),
		Slope1:               Dec("0.07"),
		Slope2:               Dec("1.25"),
		EnableStableBorrow:   false,
		StableBase:           Dec("0.0"),
		StableSlope1:         Dec("0.0"),
		StableSlope2:         Dec("0.0"),
		Ltv:                  Dec("0.7"),
		LiquidationThreshold: Dec("0.75"),
		LiquidationPenalty:   Dec("0.05"),
		LiquidationBonus:     Dec("0.05"),
		ReserveFactor:        Dec("0.2"),
		CAssetID:             5,
	}
	lendKeeper.SetAssetRatesParams(s.ctx, atomRatesParams)

	cmdxRatesParams := lendtypes.AssetRatesParams{
		AssetID:              2,
		UOptimal:             Dec("0.5"),
		Base:                 Dec("0.002"),
		Slope1:               Dec("0.08"),
		Slope2:               Dec("2.0"),
		EnableStableBorrow:   false,
		StableBase:           Dec("0.0"),
		StableSlope1:         Dec("0.0"),
		StableSlope2:         Dec("0.0"),
		Ltv:                  Dec("0.5"),
		LiquidationThreshold: Dec("0.55"),
		LiquidationPenalty:   Dec("0.05"),
		LiquidationBonus:     Dec("0.05"),
		ReserveFactor:        Dec("0.2"),
		CAssetID:             6,
	}
	lendKeeper.SetAssetRatesParams(s.ctx, cmdxRatesParams)

	var (
		assetDataCMDXPool []*lendtypes.AssetDataPoolMapping
	)
	assetDataPoolOneAssetOne := &lendtypes.AssetDataPoolMapping{
		AssetID:          1,
		AssetTransitType: 3,
		SupplyCap:        sdk.NewDec(5000000000000),
	}
	assetDataPoolOneAssetTwo := &lendtypes.AssetDataPoolMapping{
		AssetID:          2,
		AssetTransitType: 1,
		SupplyCap:        sdk.NewDec(1000000000000),
	}
	assetDataPoolOneAssetThree := &lendtypes.AssetDataPoolMapping{
		AssetID:          3,
		AssetTransitType: 2,
		SupplyCap:        sdk.NewDec(5000000000000),
	}

	assetDataCMDXPool = append(assetDataCMDXPool, assetDataPoolOneAssetOne, assetDataPoolOneAssetTwo, assetDataPoolOneAssetThree)
	cmdxPool := lendtypes.Pool{
		ModuleName: "cmdx",
		CPoolName:  "CMDX-ATOM-CMST",
		AssetData:  assetDataCMDXPool,
	}
	err = lendKeeper.AddPoolRecords(s.ctx, cmdxPool)
	if err != nil {
		panic(err)
	}

	cmdxcmstPair := lendtypes.Extended_Pair{ // 1
		AssetIn:         2,
		AssetOut:        3,
		IsInterPool:     false,
		AssetOutPoolID:  1,
		MinUsdValueLeft: 100000,
	}
	err = lendKeeper.AddLendPairsRecords(s.ctx, cmdxcmstPair)
	if err != nil {
		panic(err)
	}
	cmdxatomPair := lendtypes.Extended_Pair{ // 2
		AssetIn:         2,
		AssetOut:        1,
		IsInterPool:     false,
		AssetOutPoolID:  1,
		MinUsdValueLeft: 100000,
	}
	err = lendKeeper.AddLendPairsRecords(s.ctx, cmdxatomPair)
	if err != nil {
		panic(err)
	}
	atomcmdxPair := lendtypes.Extended_Pair{ // 3
		AssetIn:         1,
		AssetOut:        2,
		IsInterPool:     false,
		AssetOutPoolID:  1,
		MinUsdValueLeft: 100000,
	}
	err = lendKeeper.AddLendPairsRecords(s.ctx, atomcmdxPair)
	if err != nil {
		panic(err)
	}
	atomcmstPair := lendtypes.Extended_Pair{ // 4
		AssetIn:         1,
		AssetOut:        3,
		IsInterPool:     false,
		AssetOutPoolID:  1,
		MinUsdValueLeft: 100000,
	}
	err = lendKeeper.AddLendPairsRecords(s.ctx, atomcmstPair)
	if err != nil {
		panic(err)
	}
	cmstcmdxPair := lendtypes.Extended_Pair{ // 5
		AssetIn:         3,
		AssetOut:        2,
		IsInterPool:     false,
		AssetOutPoolID:  1,
		MinUsdValueLeft: 100000,
	}
	err = lendKeeper.AddLendPairsRecords(s.ctx, cmstcmdxPair)
	if err != nil {
		panic(err)
	}
	cmstatomPair := lendtypes.Extended_Pair{ // 6
		AssetIn:         3,
		AssetOut:        1,
		IsInterPool:     false,
		AssetOutPoolID:  1,
		MinUsdValueLeft: 100000,
	}
	err = lendKeeper.AddLendPairsRecords(s.ctx, cmstatomPair)
	if err != nil {
		panic(err)
	}

	// Adding Lend Pair Mapping
	map1 := lendtypes.AssetToPairMapping{
		PoolID:  1,
		AssetID: 1,
		PairID:  []uint64{3, 4},
	}
	lendKeeper.SetAssetToPair(s.ctx, map1)
	map2 := lendtypes.AssetToPairMapping{
		PoolID:  1,
		AssetID: 2,
		PairID:  []uint64{1, 2},
	}
	lendKeeper.SetAssetToPair(s.ctx, map2)
	map3 := lendtypes.AssetToPairMapping{
		PoolID:  1,
		AssetID: 3,
		PairID:  []uint64{5, 6},
	}
	lendKeeper.SetAssetToPair(s.ctx, map3)

	auctionParams := lendtypes.AuctionParams{
		AppId:                  3,
		AuctionDurationSeconds: 21600,
		Buffer:                 Dec("1.2"),
		Cusp:                   Dec("0.7"),
		Step:                   sdk.NewInt(360),
		PriceFunctionType:      1,
		DutchId:                3,
		BidDurationSeconds:     3600,
	}
	err = lendKeeper.AddAuctionParamsData(s.ctx, auctionParams)
	if err != nil {
		return
	}

	userAddress := "cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t"
	s.fundAddr2(sdk.AccAddress(userAddress), sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewIntFromUint64(100000000000000))))
	err = lendKeeper.FundModAcc(s.ctx, 1, 1, sdk.AccAddress(userAddress), sdk.NewCoin("uatom", sdk.NewInt(100000000000)))
	s.Require().NoError(err)
	s.fundAddr2(sdk.AccAddress(userAddress), sdk.NewCoins(sdk.NewCoin("ucmdx", sdk.NewIntFromUint64(100000000000000))))

	err = lendKeeper.FundModAcc(s.ctx, 1, 2, sdk.AccAddress(userAddress), sdk.NewCoin("ucmdx", sdk.NewInt(1000000000000)))
	s.Require().NoError(err)
	s.fundAddr2(sdk.AccAddress(userAddress), sdk.NewCoins(sdk.NewCoin("ucmst", sdk.NewIntFromUint64(100000000000000))))

	err = lendKeeper.FundModAcc(s.ctx, 1, 3, sdk.AccAddress(userAddress), sdk.NewCoin("ucmst", sdk.NewInt(100000000000)))
	s.Require().NoError(err)
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

	s.fundAddr(userAddress, sdk.NewCoin("ucmdx", sdk.NewIntFromUint64(1000000000)))
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

	s.fundAddr(userAddress, sdk.NewCoin("ucmdx", sdk.NewIntFromUint64(1000000000)))
	_, err := server.MsgCreate(sdk.WrapSDKContext(*ctx), &msg2)
	s.Require().NoError(err)

	msg3 := vaulttypes.MsgCreateRequest{
		From:                userAddress1,
		AppId:               1,
		ExtendedPairVaultId: extendedVaultPairID1,
		AmountIn:            sdk.NewInt(1000000000),
		AmountOut:           sdk.NewInt(100000000),
	}

	s.fundAddr(userAddress1, sdk.NewCoin("ucmdx", sdk.NewIntFromUint64(1000000000)))
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

func (s *KeeperTestSuite) TestCreateLend() {
	s.ctx = s.ctx.WithBlockTime(utils.ParseTime("2022-03-01T12:00:00Z"))
	s.ctx = s.ctx.WithBlockHeight(10)
	userAddress := "cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t"

	s.AddAppAssetLend()
	lendKeeper, ctx := &s.lendKeeper, &s.ctx
	server := lendkeeper.NewMsgServerImpl(*lendKeeper)

	msg2 := lendtypes.MsgBorrowAlternate{
		Lender:         "cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t",
		AssetId:        1,
		PoolId:         1,
		AmountIn:       sdk.NewCoin("uatom", sdk.NewInt(100000000000)),
		PairId:         3,
		IsStableBorrow: false,
		AmountOut:      sdk.NewCoin("ucmdx", sdk.NewInt(1000000000)),
		AppId:          3,
	}

	s.fundAddr(userAddress, sdk.NewCoin("uatom", sdk.NewIntFromUint64(1000000000000)))
	_, err := server.BorrowAlternate(sdk.WrapSDKContext(*ctx), &msg2)
	s.Require().NoError(err)
}

//func (s *KeeperTestSuite) TestCreateExtRewardsLend() {
//	s.ctx = s.ctx.WithBlockTime(utils.ParseTime("2022-03-01T12:00:00Z"))
//	s.ctx = s.ctx.WithBlockHeight(10)
//	userAddress := "cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t"
//	amt, _ := sdk.NewIntFromString("1000000000000000000000")
//	s.fundAddr2(sdk.AccAddress(userAddress), sdk.NewCoins(sdk.NewCoin("uatom", amt)))
//
//	s.TestCreateLend()
//	rewardsKeeper, ctx := &s.rewardsKeeper, &s.ctx
//	server := keeper.NewMsgServerImpl(*rewardsKeeper)
//	s.fundAddr(userAddress, sdk.NewCoin("ucmst", sdk.NewInt(1234567890)))
//	availableBalances := s.getBalances(sdk.MustAccAddressFromBech32(userAddress))
//	fmt.Println("bal at t0 ", availableBalances)
//
//	for _, tc := range []struct {
//		name          string
//		msg           types.ActivateExternalRewardsLend
//		expectedError bool
//		ExpErr        error
//	}{
//		{
//			"ActivateExternalRewardsLockers : success",
//			types.ActivateExternalRewardsLend{
//				AppMappingId:         3,
//				CPoolId:              1,
//				AssetId:              []uint64{2},
//				CSwapAppId:           1,
//				CSwapMinLockAmount:   0,
//				TotalRewards:         sdk.NewCoin("ucmst", sdk.NewInt(1234567890)),
//				MasterPoolId:         1,
//				DurationDays:         4,
//				MinLockupTimeSeconds: 1,
//				Depositor:            "cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t",
//			},
//			false,
//			nil,
//		},
//	} {
//		s.Run(tc.name, func() {
//
//			_, err := server.ExternalRewardsLend(sdk.WrapSDKContext(*ctx), &tc.msg)
//			if tc.ExpErr != nil {
//				s.Require().Error(err)
//				s.Require().EqualError(err, tc.ExpErr.Error())
//			} else {
//				s.Require().NoError(err)
//				availableBalances := s.getBalances(sdk.MustAccAddressFromBech32(userAddress))
//				fmt.Println("bal when created ext rewards", availableBalances)
//			}
//		})
//	}
//
//	s.ctx = s.ctx.WithBlockTime(utils.ParseTime("2022-03-02T12:10:00Z"))
//	s.ctx = s.ctx.WithBlockHeight(11)
//	req := abci.RequestBeginBlock{}
//	rewards.BeginBlocker(*ctx, req, *rewardsKeeper)
//	availableBalances = s.getBalances(sdk.MustAccAddressFromBech32(userAddress))
//	fmt.Println("bal at first day", availableBalances)
//
//	s.ctx = s.ctx.WithBlockTime(utils.ParseTime("2022-03-03T12:11:00Z"))
//	s.ctx = s.ctx.WithBlockHeight(12)
//	rewards.BeginBlocker(*ctx, req, *rewardsKeeper)
//	availableBalances = s.getBalances(sdk.MustAccAddressFromBech32(userAddress))
//	fmt.Println("bal at second day", availableBalances)
//	rewards.BeginBlocker(*ctx, req, *rewardsKeeper)
//
//	s.ctx = s.ctx.WithBlockTime(utils.ParseTime("2022-03-04T12:12:00Z"))
//	s.ctx = s.ctx.WithBlockHeight(15)
//	rewards.BeginBlocker(*ctx, req, *rewardsKeeper)
//	availableBalances = s.getBalances(sdk.MustAccAddressFromBech32(userAddress))
//	fmt.Println("bal at third day", availableBalances)
//	rewards.BeginBlocker(*ctx, req, *rewardsKeeper)
//
//	s.ctx = s.ctx.WithBlockTime(utils.ParseTime("2022-03-05T12:15:00Z"))
//	s.ctx = s.ctx.WithBlockHeight(15)
//	rewards.BeginBlocker(*ctx, req, *rewardsKeeper)
//	availableBalances = s.getBalances(sdk.MustAccAddressFromBech32(userAddress))
//	fmt.Println("bal at fourth day", availableBalances)
//	rewards.BeginBlocker(*ctx, req, *rewardsKeeper)
//
//}
