package keeper

import (
	"github.com/comdex-official/comdex/x/cdp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"testing"
)

func TestKeeper_QueryCDP(t *testing.T) {

	ctx := sdk.WrapSDKContext(setupctx(t))
	keeper := setupKeeper(t)
	queryCDP := types.QueryCDPRequest{
		CollateralType: "cmdx-a",
		Owner:          "",
	}
	Keeper.QueryCDP(*keeper,ctx,&queryCDP )
}
