package keeper_test

import (
	assetTypes "github.com/comdex-official/comdex/x/asset/types"
	"github.com/comdex-official/comdex/x/locker/keeper"
	lockerTypes "github.com/comdex-official/comdex/x/locker/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (s *KeeperTestSuite) AddAppAsset() {
	assetKeeper, ctx := &s.assetKeeper, &s.ctx
	msg1 := []assetTypes.AppMapping{{
		Name:             "cswap",
		ShortName:        "cswap",
		MinGovDeposit:    sdk.NewIntFromUint64(10000000),
		GovTimeInSeconds: 900},
		{
			Name:             "commodo",
			ShortName:        "commodo",
			MinGovDeposit:    sdk.NewIntFromUint64(10000000),
			GovTimeInSeconds: 900},
	}
	err := assetKeeper.AddAppMappingRecords(*ctx, msg1...)
	s.Require().NoError(err)

	msg2 := []assetTypes.Asset{
		{Name: "CMDX",
			Denom:     "ucmdx",
			Decimals:  1000000,
			IsOnchain: true}, {Name: "CMST",
			Denom:     "ucmst",
			Decimals:  1000000,
			IsOnchain: true}, {Name: "HARBOR",
			Denom:     "uharbor",
			Decimals:  1000000,
			IsOnchain: true},
	}
	err = assetKeeper.AddAssetRecords(*ctx, msg2...)
	s.Require().NoError(err)

}

func (s *KeeperTestSuite) TestCreateLocker() {
	userAddress := "cosmos1q7q90qsl9g0gl2zz0njxwv2a649yqrtyxtnv3v"
	lockerKeeper, ctx := &s.lockerKeeper, &s.ctx
	s.AddAppAsset()
	server := keeper.NewMsgServiceServer(*lockerKeeper)

	//Add whitelisted App Asset combinations
	for _, tc := range []struct {
		name string
		msg  lockerTypes.MsgAddWhiteListedAssetRequest
	}{
		{"Whitelist : App1 Asset 1",
			lockerTypes.MsgAddWhiteListedAssetRequest{
				From:         userAddress,
				AppMappingId: 1,
				AssetId:      1,
			},
		},
		{"Whitelist : App1 Asset 2",
			lockerTypes.MsgAddWhiteListedAssetRequest{
				From:         userAddress,
				AppMappingId: 1,
				AssetId:      2,
			},
		},
		{"Whitelist : App2 Asset 1",
			lockerTypes.MsgAddWhiteListedAssetRequest{
				From:         userAddress,
				AppMappingId: 2,
				AssetId:      1,
			},
		},
	} {
		s.Run(tc.name, func() {
			_, err := server.MsgAddWhiteListedAsset(sdk.WrapSDKContext(*ctx), &tc.msg)
			s.Require().NoError(err)
		})
	}

	// create lockers for App Asset combination , query locker and validate
	for _, tc := range []struct {
		name          string
		msg           lockerTypes.MsgCreateLockerRequest
		fundAmount    uint64
		expectedError bool
		query         lockerTypes.QueryLockerInfoRequest
	}{
		{"CreateLocker : Insufficient balance App1 Asset 1",
			lockerTypes.MsgCreateLockerRequest{
				Depositor:    userAddress,
				Amount:       sdk.NewIntFromUint64(1000000),
				AssetId:      1,
				AppMappingId: 1,
			},
			100000,
			true,
			lockerTypes.QueryLockerInfoRequest{
				Id: "cswap1",
			},
		},
		{"CreateLocker : App1 Asset 1",
			lockerTypes.MsgCreateLockerRequest{
				Depositor:    userAddress,
				Amount:       sdk.NewIntFromUint64(1000000),
				AssetId:      1,
				AppMappingId: 1,
			},
			1000000,
			false,
			lockerTypes.QueryLockerInfoRequest{
				Id: "cswap1",
			},
		},
		{"CreateLocker : Duplicate locker App1 Asset1 should fail",
			lockerTypes.MsgCreateLockerRequest{
				Depositor:    userAddress,
				Amount:       sdk.NewIntFromUint64(1000000),
				AssetId:      1,
				AppMappingId: 1,
			},
			1000000,
			true,
			lockerTypes.QueryLockerInfoRequest{
				Id: "cswap1",
			},
		},
		{"CreateLocker : App1 Asset 2",
			lockerTypes.MsgCreateLockerRequest{
				Depositor:    userAddress,
				Amount:       sdk.NewIntFromUint64(2000000),
				AssetId:      2,
				AppMappingId: 1,
			},
			2000000,
			false,
			lockerTypes.QueryLockerInfoRequest{
				Id: "cswap2",
			},
		},
		{"CreateLocker : App2 Asset 1",
			lockerTypes.MsgCreateLockerRequest{
				Depositor:    userAddress,
				Amount:       sdk.NewIntFromUint64(9900000),
				AssetId:      1,
				AppMappingId: 2,
			},
			9900000,
			false,
			lockerTypes.QueryLockerInfoRequest{
				Id: "commodo1",
			},
		},
	} {
		s.Run(tc.name, func() {
			if tc.msg.AssetId == 1 {
				s.fundAddr(userAddress, sdk.NewCoin("ucmdx", sdk.NewIntFromUint64(tc.fundAmount)))
			} else {
				s.fundAddr(userAddress, sdk.NewCoin("ucmst", sdk.NewIntFromUint64(tc.fundAmount)))
			}
			_, err := server.MsgCreateLocker(sdk.WrapSDKContext(*ctx), &tc.msg)
			if tc.expectedError {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				lockerInfo, err := s.querier.QueryLockerInfo(sdk.WrapSDKContext(*ctx), &tc.query)
				s.Require().NoError(err)
				s.Require().Equal(lockerInfo.LockerInfo.Depositor, tc.msg.Depositor)
				s.Require().Equal(lockerInfo.LockerInfo.AppMappingId, tc.msg.AppMappingId)
				s.Require().Equal(lockerInfo.LockerInfo.AssetDepositId, tc.msg.AssetId)
				s.Require().Equal(lockerInfo.LockerInfo.NetBalance, tc.msg.Amount)
			}
		})
	}

}

func (s *KeeperTestSuite) TestDepositLocker() {
	userAddress := "cosmos1q7q90qsl9g0gl2zz0njxwv2a649yqrtyxtnv3v"
	lockerKeeper, ctx := &s.lockerKeeper, &s.ctx
	//s.AddAppAsset()
	s.TestCreateLocker()
	server := keeper.NewMsgServiceServer(*lockerKeeper)
	for _, tc := range []struct {
		name          string
		msg           lockerTypes.MsgDepositAssetRequest
		query         lockerTypes.QueryLockerInfoRequest
		fundAmount    uint64
		expectedError bool
	}{
		{"DepositLocker : Insufficient Balance App1 Asset 1",
			lockerTypes.MsgDepositAssetRequest{
				Depositor:    userAddress,
				LockerId:     "cswap1",
				Amount:       sdk.NewIntFromUint64(4000000000),
				AssetId:      1,
				AppMappingId: 1,
			},
			lockerTypes.QueryLockerInfoRequest{
				Id: "cswap1",
			},
			9900000,
			true,
		},
		{"DepositLocker : App1 Asset 1",
			lockerTypes.MsgDepositAssetRequest{
				Depositor:    userAddress,
				LockerId:     "cswap1",
				Amount:       sdk.NewIntFromUint64(4000000),
				AssetId:      1,
				AppMappingId: 1,
			},
			lockerTypes.QueryLockerInfoRequest{
				Id: "cswap1",
			},
			9900000,
			false,
		},
		{"DepositLocker : App2 Asset 1",
			lockerTypes.MsgDepositAssetRequest{
				Depositor:    userAddress,
				LockerId:     "commodo1",
				Amount:       sdk.NewIntFromUint64(4000000),
				AssetId:      1,
				AppMappingId: 2,
			},
			lockerTypes.QueryLockerInfoRequest{
				Id: "commodo1",
			},
			9900000,
			false,
		},
	} {
		s.Run(tc.name, func() {
			if tc.msg.AssetId == 1 {
				s.fundAddr(userAddress, sdk.NewCoin("ucmdx", sdk.NewIntFromUint64(tc.fundAmount)))
			} else {
				s.fundAddr(userAddress, sdk.NewCoin("ucmst", sdk.NewIntFromUint64(tc.fundAmount)))
			}
			lockerInfo, err := s.querier.QueryLockerInfo(sdk.WrapSDKContext(*ctx), &tc.query)
			s.Require().NoError(err)
			previousNetAmount := lockerInfo.LockerInfo.NetBalance
			_, err = server.MsgDepositAsset(sdk.WrapSDKContext(*ctx), &tc.msg)
			if tc.expectedError {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				lockerInfo, err := s.querier.QueryLockerInfo(sdk.WrapSDKContext(*ctx), &tc.query)
				s.Require().NoError(err)
				s.Require().Equal(lockerInfo.LockerInfo.Depositor, tc.msg.Depositor)
				s.Require().Equal(lockerInfo.LockerInfo.AppMappingId, tc.msg.AppMappingId)
				s.Require().Equal(lockerInfo.LockerInfo.AssetDepositId, tc.msg.AssetId)
				s.Require().Equal(lockerInfo.LockerInfo.NetBalance, tc.msg.Amount.Add(previousNetAmount))
			}
		})
	}
}

func (s *KeeperTestSuite) TestWithdrawLocker() {
	userAddress := "cosmos1q7q90qsl9g0gl2zz0njxwv2a649yqrtyxtnv3v"
	lockerKeeper, ctx := &s.lockerKeeper, &s.ctx
	//s.AddAppAsset()
	s.TestCreateLocker()
	server := keeper.NewMsgServiceServer(*lockerKeeper)
	for _, tc := range []struct {
		name          string
		msg           lockerTypes.MsgWithdrawAssetRequest
		query         lockerTypes.QueryLockerInfoRequest
		expectedError bool
		partial       bool
	}{
		{"WithdrawLocker : Insufficient Balance App1 Asset 1",
			lockerTypes.MsgWithdrawAssetRequest{
				Depositor:    userAddress,
				LockerId:     "cswap1",
				Amount:       sdk.NewIntFromUint64(10000000),
				AssetId:      1,
				AppMappingId: 1,
			},
			lockerTypes.QueryLockerInfoRequest{
				Id: "cswap1",
			},
			true,
			true,
		},
		{"WithdrawLocker : Partial withdraw App1 Asset 1",
			lockerTypes.MsgWithdrawAssetRequest{
				Depositor:    userAddress,
				LockerId:     "cswap1",
				Amount:       sdk.NewIntFromUint64(100000),
				AssetId:      1,
				AppMappingId: 1,
			},
			lockerTypes.QueryLockerInfoRequest{
				Id: "cswap1",
			},
			false,
			true,
		},
		{"WithdrawLocker : Full Withdraw App1 Asset 1",
			lockerTypes.MsgWithdrawAssetRequest{
				Depositor:    userAddress,
				LockerId:     "cswap1",
				Amount:       sdk.NewIntFromUint64(900000),
				AssetId:      1,
				AppMappingId: 1,
			},
			lockerTypes.QueryLockerInfoRequest{
				Id: "cswap1",
			},
			false,
			false,
		},
		{"WithdrawLocker : Full Withdraw App2 Asset 1",
			lockerTypes.MsgWithdrawAssetRequest{
				Depositor:    userAddress,
				LockerId:     "commodo1",
				Amount:       sdk.NewIntFromUint64(9900000),
				AssetId:      1,
				AppMappingId: 2,
			},
			lockerTypes.QueryLockerInfoRequest{
				Id: "commodo1",
			},
			false,
			false,
		},
	} {
		s.Run(tc.name, func() {
			lockerInfo, err := s.querier.QueryLockerInfo(sdk.WrapSDKContext(*ctx), &tc.query)
			s.Require().NoError(err)
			previousNetAmount := lockerInfo.LockerInfo.NetBalance
			_, err = server.MsgWithdrawAsset(sdk.WrapSDKContext(*ctx), &tc.msg)
			if tc.expectedError {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				lockerInfo, err := s.querier.QueryLockerInfo(sdk.WrapSDKContext(*ctx), &tc.query)
				s.Require().NoError(err)
				s.Require().Equal(lockerInfo.LockerInfo.Depositor, tc.msg.Depositor)
				s.Require().Equal(lockerInfo.LockerInfo.AppMappingId, tc.msg.AppMappingId)
				s.Require().Equal(lockerInfo.LockerInfo.AssetDepositId, tc.msg.AssetId)
				if tc.partial {
					s.Require().Equal(lockerInfo.LockerInfo.NetBalance, previousNetAmount.Sub(tc.msg.Amount))
				} else {
					s.Require().Equal(lockerInfo.LockerInfo.NetBalance, sdk.NewIntFromUint64(0))
				}
			}
		})
	}

}
