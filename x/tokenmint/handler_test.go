package tokenmint_test

import (
	"strings"
	"testing"

	errorsmod "cosmossdk.io/errors"
	"github.com/comdex-official/comdex/app"
	"github.com/comdex-official/comdex/x/tokenmint"
	"github.com/comdex-official/comdex/x/tokenmint/keeper"
	"github.com/comdex-official/comdex/x/tokenmint/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestInvalidMsg(t *testing.T) {
	app1 := app.Setup(t, false)

	app1.TokenmintKeeper = keeper.NewKeeper(
		app1.AppCodec(), app1.GetKey(types.StoreKey), app1.BankKeeper, &app1.AssetKeeper)
	h := tokenmint.NewHandler(app1.TokenmintKeeper)

	res, err := h(sdk.NewContext(nil, tmproto.Header{}, false, nil), testdata.NewTestMsg())
	require.Error(t, err)
	require.Nil(t, res)

	_, _, log := errorsmod.ABCIInfo(err, false)
	require.True(t, strings.Contains(log, "unknown request"))
}

func (s *ModuleTestSuite) TestMsgMintNewTokensRequest() {
	handler := tokenmint.NewHandler(s.keeper)
	addr1 := s.addr(1)
	appID1 := s.CreateNewApp("appone")
	s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	s.CreateNewAsset("ASSETTWO", "uasset2", 2000000)
	msg := types.NewMsgMintNewTokensRequest(
		addr1.String(), appID1, 1,
	)

	_, err := handler(s.ctx, msg)
	s.Require().Error(err)
}
