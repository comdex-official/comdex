package expected

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	assettypes "github.com/comdex-official/comdex/x/asset/types"
)

type BankKeeper interface {
	SpendableCoins(ctx sdk.Context, address sdk.AccAddress) sdk.Coins
}

type AssetKeeper interface {
	GetPair(ctx sdk.Context, id uint64) (pair assettypes.Pair, found bool)
}
