package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	paramTypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// AccountKeeper expected interface for the account keeper (noalias)
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) types.AccountI
	SetAccount(ctx sdk.Context, acc types.AccountI)
}

type BankKeeper interface {
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error

	MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	BurnCoins(ctx sdk.Context, name string, amt sdk.Coins) error
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
}

// ParamSubspace defines the expected Subspace interfacace
type ParamSubspace interface {
	HasKeyTable() bool
	WithKeyTable(table paramTypes.KeyTable) paramTypes.Subspace
	Get(ctx sdk.Context, key []byte, ptr interface{})
	GetParamSet(ctx sdk.Context, ps paramTypes.ParamSet)
	SetParamSet(ctx sdk.Context, ps paramTypes.ParamSet)
}
