package keeper_test

import (
	"github.com/comdex-official/comdex/app"
	assetkeeper "github.com/comdex-official/comdex/x/asset/keeper"
	"github.com/comdex-official/comdex/x/vault/keeper"
	"github.com/comdex-official/comdex/x/vault/types"
	types2 "github.com/comdex-official/comdex/x/vault/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtime "github.com/tendermint/tendermint/types/time"
	"testing"
)
type MsgTestSuite struct {
	suite.Suite
	keeper keeper.Keeper
	assetKeeper assetkeeper.Keeper
	app    app.TestApp
	ctx    sdk.Context
	prq    query.PageRequest
}
func (suite *MsgTestSuite) SetupTest() {
	testApp := app.NewTestApp()
	k := testApp.GetVaultKeeper()
	ctx := testApp.NewContext(true, tmproto.Header{Height: 1, Time: tmtime.Now()})
	suite.app = testApp
	suite.keeper = k
	suite.ctx = ctx

	return
}
func TestMsgTestSuite(t *testing.T) {
	suite.Run(t, new(MsgTestSuite))
}

func (suite *MsgTestSuite)TestMsgCreate() {
	app.SetAccountAddressPrefixes()
	msgreq := types.MsgCreateRequest{
		From:      "comdex1yples84d8avjlmegn90663mmjs4tardw45af6s",
		PairID:    1,
		AmountIn:  sdk.Int{},
		AmountOut: sdk.Int{},
	}

	msgServer := keeper.NewMsgServiceServer(suite.keeper)
	_,err := msgServer.MsgCreate(sdk.WrapSDKContext(suite.ctx), &msgreq)
	suite.Error(err)

	msgrequest := types.MsgCreateRequest{
		From:      "comdex1yples84d8avjlmegn90663mmjs4tardw45af6v",
		PairID:    1,
		AmountIn:  sdk.Int(sdk.NewInt(100)),
		AmountOut: sdk.Int(sdk.NewInt(100)),
	}
	_,err = msgServer.MsgCreate(sdk.WrapSDKContext(suite.ctx), &msgrequest)
	suite.Error(err)
}

func (suite *MsgTestSuite) TestMsgDeposit(){
	msgServer := keeper.NewMsgServiceServer(suite.keeper)

	msgDepositReq := types.MsgDepositRequest{
		From:   "comdex1yples84d8avjlmegn90663mmjs4tardw45af6a",
		ID:     1,
		Amount: sdk.Int(sdk.NewInt(100)),
	}
	_,err := msgServer.MsgDeposit(sdk.WrapSDKContext(suite.ctx), &msgDepositReq)
	suite.Error(err)

	msgDepositReq = types.MsgDepositRequest{
		From:   "comdex1yples84d8avjlmegn90663mmjs4tardw45af6v",
		ID:     1,
		Amount: sdk.Int(sdk.NewInt(100)),
	}
	_,err = msgServer.MsgDeposit(sdk.WrapSDKContext(suite.ctx), &msgDepositReq)
	suite.Error(err)

	vault := types2.Vault{ID: 1, PairID: 1, Owner: "comdex1yples84d8avjlmegn90663mmjs4tardw45af6a", AmountIn: sdk.NewInt(100), AmountOut: sdk.NewInt(100)}
	suite.keeper.SetVault(suite.ctx, vault)

	_,err = msgServer.MsgDeposit(sdk.WrapSDKContext(suite.ctx), &msgDepositReq)
	suite.Error(err)

	vault = types2.Vault{ID: 1, PairID: 1, Owner: "comdex1yples84d8avjlmegn90663mmjs4tardw45af6v", AmountIn: sdk.NewInt(100), AmountOut: sdk.NewInt(100)}
	suite.keeper.SetVault(suite.ctx, vault)

	_,err = msgServer.MsgDeposit(sdk.WrapSDKContext(suite.ctx), &msgDepositReq)
	suite.Error(err)
}

func (suite *MsgTestSuite) TestMsgWithdraw(){
	msgServer := keeper.NewMsgServiceServer(suite.keeper)

	msgWithdrawReq := types.MsgWithdrawRequest{
		From:   "comdex1yples84d8avjlmegn90663mmjs4tardw45af6a",
		ID:     1,
		Amount: sdk.Int(sdk.NewInt(100)),
	}
	_,err := msgServer.MsgWithdraw(sdk.WrapSDKContext(suite.ctx), &msgWithdrawReq )
	suite.Error(err)

	msgWithdrawReq = types.MsgWithdrawRequest{
		From:   "comdex1yples84d8avjlmegn90663mmjs4tardw45af6v",
		ID:     1,
		Amount: sdk.Int(sdk.NewInt(100)),
	}
	_,err = msgServer.MsgWithdraw(sdk.WrapSDKContext(suite.ctx), &msgWithdrawReq )
	suite.Error(err)

	vault := types2.Vault{ID: 1, PairID: 1, Owner: "comdex1yples84d8avjlmegn90663mmjs4tardw45af6a", AmountIn: sdk.NewInt(100), AmountOut: sdk.NewInt(100)}
	suite.keeper.SetVault(suite.ctx, vault)

	_,err = msgServer.MsgWithdraw(sdk.WrapSDKContext(suite.ctx), &msgWithdrawReq )
	suite.Error(err)

	vault = types2.Vault{ID: 1, PairID: 1, Owner: "comdex1yples84d8avjlmegn90663mmjs4tardw45af6v", AmountIn: sdk.NewInt(100), AmountOut: sdk.NewInt(100)}
	suite.keeper.SetVault(suite.ctx, vault)

	_,err = msgServer.MsgWithdraw(sdk.WrapSDKContext(suite.ctx), &msgWithdrawReq )
	suite.Error(err)
}

func (suite *MsgTestSuite) TestMsgDraw(){
	msgServer := keeper.NewMsgServiceServer(suite.keeper)

	msgDrawReq := types.MsgDrawRequest{
		From:   "comdex1yples84d8avjlmegn90663mmjs4tardw45af6a",
		ID:     1,
		Amount: sdk.Int(sdk.NewInt(100)),
	}

	_,err := msgServer.MsgDraw(sdk.WrapSDKContext(suite.ctx), &msgDrawReq)
	suite.Error(err)

	msgDrawReq = types.MsgDrawRequest{
		From:   "comdex1yples84d8avjlmegn90663mmjs4tardw45af6v",
		ID:     1,
		Amount: sdk.Int(sdk.NewInt(100)),
	}

	_,err = msgServer.MsgDraw(sdk.WrapSDKContext(suite.ctx), &msgDrawReq)
	suite.Error(err)

	vault := types2.Vault{ID: 1, PairID: 1, Owner: "comdex1yples84d8avjlmegn90663mmjs4tardw45af6a", AmountIn: sdk.NewInt(100), AmountOut: sdk.NewInt(100)}
	suite.keeper.SetVault(suite.ctx, vault)

	_,err = msgServer.MsgDraw(sdk.WrapSDKContext(suite.ctx), &msgDrawReq)
	suite.Error(err)

	vault = types2.Vault{ID: 1, PairID: 1, Owner: "comdex1yples84d8avjlmegn90663mmjs4tardw45af6v", AmountIn: sdk.NewInt(100), AmountOut: sdk.NewInt(100)}
	suite.keeper.SetVault(suite.ctx, vault)

	_,err = msgServer.MsgDraw(sdk.WrapSDKContext(suite.ctx), &msgDrawReq)
	suite.Error(err)
}

func (suite *MsgTestSuite) TestMsgRepay(){
	msgServer := keeper.NewMsgServiceServer(suite.keeper)

	msgRepayReq := types.MsgRepayRequest{
		From:   "comdex1yples84d8avjlmegn90663mmjs4tardw45af6a",
		ID:     1,
		Amount: sdk.Int(sdk.NewInt(100)),
	}

	_,err := msgServer.MsgRepay(sdk.WrapSDKContext(suite.ctx), &msgRepayReq)
	suite.Error(err)

	msgRepayReq = types.MsgRepayRequest{
		From:   "comdex1yples84d8avjlmegn90663mmjs4tardw45af6v",
		ID:     1,
		Amount: sdk.Int(sdk.NewInt(100)),
	}

	_,err = msgServer.MsgRepay(sdk.WrapSDKContext(suite.ctx), &msgRepayReq)
	suite.Error(err)

	vault := types2.Vault{ID: 1, PairID: 1, Owner: "comdex1yples84d8avjlmegn90663mmjs4tardw45af6a", AmountIn: sdk.NewInt(100), AmountOut: sdk.NewInt(100)}
	suite.keeper.SetVault(suite.ctx, vault)

	_,err = msgServer.MsgRepay(sdk.WrapSDKContext(suite.ctx), &msgRepayReq)
	suite.Error(err)

	vault = types2.Vault{ID: 1, PairID: 1, Owner: "comdex1yples84d8avjlmegn90663mmjs4tardw45af6v", AmountIn: sdk.NewInt(100), AmountOut: sdk.NewInt(10)}
	suite.keeper.SetVault(suite.ctx, vault)

	_,err = msgServer.MsgRepay(sdk.WrapSDKContext(suite.ctx), &msgRepayReq)
	suite.Error(err)

	vault = types2.Vault{ID: 1, PairID: 1, Owner: "comdex1yples84d8avjlmegn90663mmjs4tardw45af6v", AmountIn: sdk.NewInt(100), AmountOut: sdk.NewInt(100)}
	suite.keeper.SetVault(suite.ctx, vault)

	_,err = msgServer.MsgRepay(sdk.WrapSDKContext(suite.ctx), &msgRepayReq)
	suite.Error(err)
}