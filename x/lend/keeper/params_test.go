package keeper_test

import (
	"github.com/comdex-official/comdex/app"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"testing"

	"github.com/comdex-official/comdex/x/lend/types"
	"github.com/stretchr/testify/require"
)

func TestGetParams(t *testing.T) {

	comdexApp := app.Setup(false)
	ctx := comdexApp.BaseApp.NewContext(false, tmproto.Header{})

	params := types.DefaultParams()

	comdexApp.LendKeeper.SetParams(ctx, params)

	require.EqualValues(t, params, comdexApp.LendKeeper.GetParams(ctx))
}
