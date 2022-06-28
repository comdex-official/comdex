package keeper_test

import (
	"fmt"

	"github.com/comdex-official/comdex/x/liquidity/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/stretchr/testify/suite"
)

func CreateNewApp(s *KeeperTestSuite, appName string) uint64 {
	err := NewSubApp(s, appName)
	s.Require().NoError(err)

	found := s.app.AssetKeeper.HasAppForName(s.ctx, appName)
	s.Require().True(found)

	appID := GetAppIDByAppName(s, appName)
	s.Require().NotZero(appID)

	return appID
}

func pairWithInvalidAppID(s *KeeperTestSuite, creator sdk.AccAddress, denom1, denom2 string) {
	msg := types.NewMsgCreatePair(1, creator, denom1, denom2)
	s.Require().NoError(msg.ValidateBasic())
	_, err := s.keeper.CreatePair(s.ctx, msg, false)
	s.Require().Error(err)
	s.Require().EqualError(err, "app id 1 not found: app id invalid")
}

func pairWithNonWhilistedBothDenoms(s *KeeperTestSuite, appID uint64, creator sdk.AccAddress, denom1, denom2 string) {
	msg := types.NewMsgCreatePair(appID, creator, denom1, denom2)
	s.Require().NoError(msg.ValidateBasic())
	_, err := s.keeper.CreatePair(s.ctx, msg, false)
	s.Require().Error(err)
	s.Require().EqualError(err, fmt.Sprintf("asset with denom  %s is not white listed: asset not whitelisted", denom1))
}

func pairWithNonWhilistedDenom1(s *KeeperTestSuite, appID uint64, creator sdk.AccAddress, denom1, denom2 string) {
	msg := types.NewMsgCreatePair(appID, creator, denom1, denom2)
	s.Require().NoError(msg.ValidateBasic())
	_, err := s.keeper.CreatePair(s.ctx, msg, false)
	s.Require().Error(err)
	s.Require().EqualError(err, fmt.Sprintf("asset with denom  %s is not white listed: asset not whitelisted", denom2))
}

func pairWithNonWhilistedDenom2(s *KeeperTestSuite, appID uint64, creator sdk.AccAddress, denom1, denom2 string) {
	msg := types.NewMsgCreatePair(appID, creator, denom1, denom2)
	s.Require().NoError(msg.ValidateBasic())
	_, err := s.keeper.CreatePair(s.ctx, msg, false)
	s.Require().Error(err)
	s.Require().EqualError(err, fmt.Sprintf("asset with denom  %s is not white listed: asset not whitelisted", denom1))
}

func pairWithInsufficientPairCreationFee(s *KeeperTestSuite, appID uint64, creator sdk.AccAddress, denom1, denom2 string, params types.GenericParams) {
	msg := types.NewMsgCreatePair(appID, creator, denom1, denom2)
	s.Require().NoError(msg.ValidateBasic())
	_, err := s.keeper.CreatePair(s.ctx, msg, false)
	s.Require().Error(err)
	s.Require().EqualError(err, fmt.Sprintf("insufficient pair creation fee: %d%s is smaller than %s: insufficient funds", 0, params.PairCreationFee[0].Denom, params.PairCreationFee.String()))
}

func pairCreationSuccess(s *KeeperTestSuite, appID, pairID uint64, creator sdk.AccAddress, denom1, denom2 string, params types.GenericParams) {
	msg := types.NewMsgCreatePair(appID, creator, denom1, denom2)
	s.Require().NoError(msg.ValidateBasic())
	pair, err := s.keeper.CreatePair(s.ctx, msg, false)
	s.Require().NoError(err)
	s.Require().IsType(pair, types.Pair{})
	s.Require().Equal(appID, pair.AppId)
	s.Require().Equal(denom1, pair.BaseCoinDenom)
	s.Require().Equal(denom2, pair.QuoteCoinDenom)
	s.Require().Equal(pairID, pair.Id)
}

func pairAlreadyExists(s *KeeperTestSuite, appID uint64, creator sdk.AccAddress, denom1, denom2 string, params types.GenericParams) {
	msg := types.NewMsgCreatePair(appID, creator, denom1, denom2)
	s.Require().NoError(msg.ValidateBasic())
	_, err := s.keeper.CreatePair(s.ctx, msg, false)
	s.Require().Error(err)
	s.Require().EqualError(err, "pair already exists")
}

func (s *KeeperTestSuite) TestPairCreation() {
	apps := []string{"appOne", "appTwo", "appThree"}
	addresses := []sdk.AccAddress{s.addr(0), s.addr(1), s.addr(2)}

	pairWithInvalidAppID(s, addresses[0], "denom1", "denom2")

	appID := CreateNewApp(s, apps[0])

	params, err := s.keeper.GetGenericParams(s.ctx, appID)
	s.Require().NoError(err)

	pairWithNonWhilistedBothDenoms(s, appID, addresses[0], "denom1", "denom2")

	err = NewAddAsset(s, "asset1", "denom1")
	s.Require().NoError(err)
	pairWithNonWhilistedDenom1(s, appID, addresses[0], "denom1", "denom2")

	err = NewAddAsset(s, "asset2", "denom2")
	s.Require().NoError(err)
	pairWithNonWhilistedDenom2(s, appID, addresses[0], "denom3", "denom2")

	err = NewAddAsset(s, "asset3", "denom3")
	s.Require().NoError(err)
	pairWithInsufficientPairCreationFee(s, appID, addresses[0], "denom1", "denom2", params)

	s.fundAddr(addresses[0], params.PairCreationFee)
	pairCreationSuccess(s, appID, 1, addresses[0], "denom1", "denom2", params)

	s.fundAddr(addresses[0], params.PairCreationFee)
	pairAlreadyExists(s, appID, addresses[0], "denom1", "denom2", params)

	app2ID := CreateNewApp(s, apps[1])
	params2, err := s.keeper.GetGenericParams(s.ctx, app2ID)
	s.Require().NoError(err)
	s.fundAddr(addresses[1], params2.PairCreationFee)
	pairCreationSuccess(s, app2ID, 1, addresses[1], "denom1", "denom2", params2)
	pairAlreadyExists(s, app2ID, addresses[1], "denom1", "denom2", params2)

	s.fundAddr(addresses[1], params2.PairCreationFee)
	pairCreationSuccess(s, app2ID, 2, addresses[1], "denom2", "denom1", params2)

	s.fundAddr(addresses[1], params2.PairCreationFee)
	pairCreationSuccess(s, app2ID, 3, addresses[1], "denom2", "denom3", params2)

	s.fundAddr(addresses[0], params.PairCreationFee)
	pairCreationSuccess(s, appID, 2, addresses[0], "denom2", "denom3", params)
}
