package keeper_test

import (
	"github.com/comdex-official/comdex/app/wasm/bindings"
	utils "github.com/comdex-official/comdex/types"
	assetTypes "github.com/comdex-official/comdex/x/asset/types"
	lockertypes "github.com/comdex-official/comdex/x/locker/types"
	keeper "github.com/comdex-official/comdex/x/rewards/keeper"
	"github.com/comdex-official/comdex/x/rewards/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
		Name:      "CMDX",
		Denom:     "ucmdx",
		Decimals:  sdk.NewInt(1000000),
		IsOnChain: true,
	}

	err = assetKeeper.AddAssetRecords(*ctx, msg3)
	s.Require().NoError(err)

	msg4 := assetTypes.Asset{
		Name:      "CMST",
		Denom:     "ucmst",
		Decimals:  sdk.NewInt(1000000),
		IsOnChain: true,
	}
	err = assetKeeper.AddAssetRecords(*ctx, msg4)
	s.Require().NoError(err)

	msg5 := assetTypes.Asset{
		Name:      "HARBOR",
		Denom:     "uharbor",
		Decimals:  sdk.NewInt(1000000),
		IsOnChain: true,
	}
	err = assetKeeper.AddAssetRecords(*ctx, msg5)
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
	userAddress := "cosmos1q7q90qsl9g0gl2zz0njxwv2a649yqrtyxtnv3v"
	lockerKeeper, ctx := &s.lockerKeeper, &s.ctx
	s.AddAppAsset()
	s.AddCollectorLookupTable()
	locker := lockertypes.Locker{
		LockerId:           1,
		Depositor:          userAddress,
		ReturnsAccumulated: sdk.ZeroInt(),
		NetBalance:         sdk.NewIntFromUint64(1000000),
		CreatedAt:          utils.ParseTime("2022-03-01T12:00:00Z"),
		AssetDepositId:     1,
		IsLocked:           false,
		AppId:              1,
		BlockHeight:        10,
		BlockTime:          utils.ParseTime("2022-03-01T12:00:00Z"),
	}
	lockerKeeper.SetLocker(*ctx, locker)
	_, found := lockerKeeper.GetLocker(*ctx, 1)
	s.Require().True(found)
}

func (s *KeeperTestSuite) TestCreateExtRewardsLocker() {
	userAddress := "cosmos1q7q90qsl9g0gl2zz0njxwv2a649yqrtyxtnv3v"

	s.TestCreateLocker()
	rewardsKeeper, ctx := &s.rewardsKeeper, &s.ctx
	amt, _ := sdk.NewIntFromString("1000000000000000000000")
	server := keeper.NewMsgServerImpl(*rewardsKeeper)
	for _, tc := range []struct {
		name          string
		msg           types.ActivateExternalRewardsLockers
		expectedError bool
		query         types.QueryExternalRewardsLockersRequest
		ExpErr        error
	}{
		{
			"WithdrawLocker : success",
			types.ActivateExternalRewardsLockers{
				AppMappingId:         1,
				AssetId:              1,
				TotalRewards:         sdk.NewCoin("weth", amt),
				DurationDays:         5,
				Depositor:            userAddress,
				MinLockupTimeSeconds: 0,
			},
			false,
			types.QueryExternalRewardsLockersRequest{},
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
				_, err := s.querier.QueryExternalRewardsLockers(sdk.WrapSDKContext(*ctx), &tc.query)
				s.Require().NoError(err)
			}
		})
	}
}
