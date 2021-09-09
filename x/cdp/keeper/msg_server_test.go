package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"testing"

	"github.com/comdex-official/comdex/x/cdp/types"
)

func setupMsgServer(t testing.TB) (types.MsgServiceServer) {
	keeper := setupKeeper(t)
	return NewMsgServerImpl(*keeper)
}

func TestMsgServer_MsgCreateCDP(t *testing.T) {
	ctxc := setupctx(t)
	_,err:= msgServer.MsgCreateCDP(msgServer{},sdk.WrapSDKContext(ctxc), types.NewMsgCreateCDPRequest(sdk.AccAddress{},sdk.Coin{}, sdk.Coin{}, "cmdx-a"))
	if err == nil{
		t.Error()
	}
}

func TestMsgServer_MsgDepositCollateral(t *testing.T) {
	ctxc := setupctx(t)
	_,err:= msgServer.MsgDepositCollateral(msgServer{},sdk.WrapSDKContext(ctxc), types.NewMsgDepositCollateralRequest(sdk.AccAddress{},sdk.Coin{}, "cmdx-a"))
	if err == nil{
		t.Error()
	}
}

func TestMsgServer_MsgWithdrawCollateral(t *testing.T) {
	ctxc := setupctx(t)
	_,err:= msgServer.MsgWithdrawCollateral(msgServer{},sdk.WrapSDKContext(ctxc), types.NewMsgWithdrawCollateralRequest(sdk.AccAddress{},sdk.Coin{}, "cmdx-a"))
	if err == nil{
		t.Error()
	}
}

func TestMsgServer_MsgDrawDebt(t *testing.T) {
	ctxc := setupctx(t)
	_,err:= msgServer.MsgDrawDebt(msgServer{},sdk.WrapSDKContext(ctxc), types.NewMsgDrawDebtRequest(sdk.AccAddress{}, "cmdx-a", sdk.Coin{}))
	if err == nil{
		t.Error()
	}
}

func TestMsgServer_MsgRepayDebt(t *testing.T) {
	ctxc := setupctx(t)
	_,err:= msgServer.MsgRepayDebt(msgServer{},sdk.WrapSDKContext(ctxc), types.NewMsgRepayDebtRequest(sdk.AccAddress{}, "cmdx-a", sdk.Coin{}))
	if err == nil{
		t.Error()
	}
}

func TestMsgServer_MsgLiquidateCDP(t *testing.T) {
	ctxc := setupctx(t)
	_,err:= msgServer.MsgLiquidateCDP(msgServer{},sdk.WrapSDKContext(ctxc), types.NewMsgLiquidateCDPRequest(sdk.AccAddress{},"cmdx-a"))
	if err == nil{
		t.Error()
	}
}
