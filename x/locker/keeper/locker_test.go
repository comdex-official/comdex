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
	msg3 := lockerTypes.MsgAddWhiteListedAssetRequest{
		From:         userAddress,
		AppMappingId: 1,
		AssetId:      1,
	}
	server := keeper.NewMsgServiceServer(*lockerKeeper)
	_, err := server.MsgAddWhiteListedAsset(sdk.WrapSDKContext(*ctx), &msg3)
	s.Require().NoError(err)

	msg4 := lockerTypes.MsgCreateLockerRequest{
		Depositor:    userAddress,
		Amount:       sdk.NewIntFromUint64(1000000),
		AssetId:      1,
		AppMappingId: 1,
	}

	//Insufficient balance check
	_, err = server.MsgCreateLocker(sdk.WrapSDKContext(*ctx), &msg4)
	s.Require().Error(err)

	//create locker for first asset
	s.fundAddr(userAddress, sdk.NewCoin("ucmdx", sdk.NewIntFromUint64(1000000)))
	_, err = server.MsgCreateLocker(sdk.WrapSDKContext(*ctx), &msg4)
	s.Require().NoError(err)

	qmsg1 := lockerTypes.QueryLockerInfoRequest{
		Id: "cswap1",
	}
	lockerInfo, err := s.querier.QueryLockerInfo(sdk.WrapSDKContext(*ctx), &qmsg1)
	s.Require().NoError(err)
	s.Require().Equal(lockerInfo.LockerInfo.Depositor, msg4.Depositor)
	s.Require().Equal(lockerInfo.LockerInfo.AppMappingId, msg4.AppMappingId)
	s.Require().Equal(lockerInfo.LockerInfo.AssetDepositId, msg4.AssetId)
	s.Require().Equal(lockerInfo.LockerInfo.NetBalance, msg4.Amount)

	//create locker two times on same asset should fail
	_, err = server.MsgCreateLocker(sdk.WrapSDKContext(*ctx), &msg4)
	s.Require().Error(err)

	//create locker for second asset
	msg5 := lockerTypes.MsgAddWhiteListedAssetRequest{
		From:         userAddress,
		AppMappingId: 1,
		AssetId:      2,
	}
	_, err = server.MsgAddWhiteListedAsset(sdk.WrapSDKContext(*ctx), &msg5)
	s.Require().NoError(err)
	msg6 := lockerTypes.MsgCreateLockerRequest{
		Depositor:    userAddress,
		Amount:       sdk.NewIntFromUint64(2000000),
		AssetId:      2,
		AppMappingId: 1,
	}
	s.fundAddr(userAddress, sdk.NewCoin("ucmst", sdk.NewIntFromUint64(2000000)))
	_, err = server.MsgCreateLocker(sdk.WrapSDKContext(*ctx), &msg6)
	s.Require().NoError(err)
	qmsg2 := lockerTypes.QueryLockerInfoRequest{
		Id: "cswap2",
	}
	lockerInfo, err = s.querier.QueryLockerInfo(sdk.WrapSDKContext(*ctx), &qmsg2)
	s.Require().NoError(err)
	s.Require().Equal(lockerInfo.LockerInfo.Depositor, msg6.Depositor)
	s.Require().Equal(lockerInfo.LockerInfo.AppMappingId, msg6.AppMappingId)
	s.Require().Equal(lockerInfo.LockerInfo.AssetDepositId, msg6.AssetId)
	s.Require().Equal(lockerInfo.LockerInfo.NetBalance, msg6.Amount)

	//create locker for different app and same asset

	msg8 := lockerTypes.MsgAddWhiteListedAssetRequest{
		From:         userAddress,
		AppMappingId: 2,
		AssetId:      1,
	}
	_, err = server.MsgAddWhiteListedAsset(sdk.WrapSDKContext(*ctx), &msg8)
	s.Require().NoError(err)

	msg7 := lockerTypes.MsgCreateLockerRequest{
		Depositor:    userAddress,
		Amount:       sdk.NewIntFromUint64(9900000),
		AssetId:      1,
		AppMappingId: 2,
	}

	//Insufficient balance check
	_, err = server.MsgCreateLocker(sdk.WrapSDKContext(*ctx), &msg7)
	s.Require().Error(err)

	//create locker for first asset
	s.fundAddr(userAddress, sdk.NewCoin("ucmdx", sdk.NewIntFromUint64(9900000)))
	_, err = server.MsgCreateLocker(sdk.WrapSDKContext(*ctx), &msg7)
	s.Require().NoError(err)

	qmsg3 := lockerTypes.QueryLockerInfoRequest{
		Id: "commodo1",
	}
	lockerInfo, err = s.querier.QueryLockerInfo(sdk.WrapSDKContext(*ctx), &qmsg3)
	s.Require().NoError(err)
	s.Require().Equal(lockerInfo.LockerInfo.Depositor, msg7.Depositor)
	s.Require().Equal(lockerInfo.LockerInfo.AppMappingId, msg7.AppMappingId)
	s.Require().Equal(lockerInfo.LockerInfo.AssetDepositId, msg7.AssetId)
	s.Require().Equal(lockerInfo.LockerInfo.NetBalance, msg7.Amount)

}

func (s *KeeperTestSuite) TestDepositLocker() {
	userAddress := "cosmos1q7q90qsl9g0gl2zz0njxwv2a649yqrtyxtnv3v"
	lockerKeeper, ctx := &s.lockerKeeper, &s.ctx
	//s.AddAppAsset()
	s.TestCreateLocker()
	msg1 := lockerTypes.MsgDepositAssetRequest{
		Depositor:    userAddress,
		LockerId:     "cswap1",
		Amount:       sdk.NewIntFromUint64(4000000),
		AssetId:      1,
		AppMappingId: 1,
	}
	server := keeper.NewMsgServiceServer(*lockerKeeper)

	//Insufficient balance
	_, err := server.MsgDepositAsset(sdk.WrapSDKContext(*ctx), &msg1)
	s.Require().Error(err)

	//Deposit Asset
	s.fundAddr(userAddress, sdk.NewCoin("ucmdx", sdk.NewIntFromUint64(4000000)))
	_, err = server.MsgDepositAsset(sdk.WrapSDKContext(*ctx), &msg1)
	s.Require().NoError(err)
	qmsg1 := lockerTypes.QueryLockerInfoRequest{
		Id: "cswap1",
	}
	lockerInfo, err := s.querier.QueryLockerInfo(sdk.WrapSDKContext(*ctx), &qmsg1)
	s.Require().NoError(err)
	s.Require().Equal(lockerInfo.LockerInfo.Depositor, msg1.Depositor)
	s.Require().Equal(lockerInfo.LockerInfo.AppMappingId, msg1.AppMappingId)
	s.Require().Equal(lockerInfo.LockerInfo.AssetDepositId, msg1.AssetId)
	s.Require().Equal(lockerInfo.LockerInfo.NetBalance, msg1.Amount.Add(sdk.NewIntFromUint64(1000000)))

}

func (s *KeeperTestSuite) TestWithdrawLocker() {
	userAddress := "cosmos1q7q90qsl9g0gl2zz0njxwv2a649yqrtyxtnv3v"
	lockerKeeper, ctx := &s.lockerKeeper, &s.ctx
	//s.AddAppAsset()
	s.TestCreateLocker()
	msg1 := lockerTypes.MsgWithdrawAssetRequest{
		Depositor:    userAddress,
		LockerId:     "cswap1",
		Amount:       sdk.NewIntFromUint64(10000000),
		AssetId:      1,
		AppMappingId: 1,
	}
	server := keeper.NewMsgServiceServer(*lockerKeeper)

	//Insufficient balance
	_, err := server.MsgWithdrawAsset(sdk.WrapSDKContext(*ctx), &msg1)
	s.Require().Error(err)

	//Partial Withdraw Asset
	msg2 := lockerTypes.MsgWithdrawAssetRequest{
		Depositor:    userAddress,
		LockerId:     "cswap1",
		Amount:       sdk.NewIntFromUint64(100000),
		AssetId:      1,
		AppMappingId: 1,
	}
	_, err = server.MsgWithdrawAsset(sdk.WrapSDKContext(*ctx), &msg2)
	s.Require().NoError(err)
	qmsg1 := lockerTypes.QueryLockerInfoRequest{
		Id: "cswap1",
	}
	lockerInfo, err := s.querier.QueryLockerInfo(sdk.WrapSDKContext(*ctx), &qmsg1)
	s.Require().NoError(err)
	s.Require().Equal(lockerInfo.LockerInfo.NetBalance, sdk.NewIntFromUint64(900000))

	//Full Withdraw Asset
	msg3 := lockerTypes.MsgWithdrawAssetRequest{
		Depositor:    userAddress,
		LockerId:     "cswap1",
		Amount:       sdk.NewIntFromUint64(900000),
		AssetId:      1,
		AppMappingId: 1,
	}
	_, err = server.MsgWithdrawAsset(sdk.WrapSDKContext(*ctx), &msg3)
	s.Require().NoError(err)
	lockerInfo, err = s.querier.QueryLockerInfo(sdk.WrapSDKContext(*ctx), &qmsg1)
	s.Require().NoError(err)
	s.Require().Equal(lockerInfo.LockerInfo.NetBalance, sdk.NewIntFromUint64(0))

}
