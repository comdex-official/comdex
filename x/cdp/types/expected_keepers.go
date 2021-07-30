package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	paramTypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// AccountKeeper expected interface for the account keeper (noalias)
type AccountKeeper interface {
	NewAccount(sdk.Context, types.AccountI) types.AccountI
	NewAccountWithAddress(ctx sdk.Context, addr sdk.AccAddress) types.AccountI

	GetAccount(ctx sdk.Context, addr sdk.AccAddress) types.AccountI
	GetAllAccounts(ctx sdk.Context) []types.AccountI
	SetAccount(ctx sdk.Context, acc types.AccountI)

	IterateAccounts(ctx sdk.Context, process func(types.AccountI) bool)

	ValidatePermissions(macc types.ModuleAccountI) error

	GetModuleAddress(moduleName string) sdk.AccAddress
	GetModuleAddressAndPermissions(moduleName string) (addr sdk.AccAddress, permissions []string)
	GetModuleAccountAndPermissions(ctx sdk.Context, moduleName string) (types.ModuleAccountI, []string)
	GetModuleAccount(ctx sdk.Context, moduleName string) types.ModuleAccountI
	SetModuleAccount(ctx sdk.Context, macc types.ModuleAccountI)
}

type BankKeeper interface {
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromModuleToModule(ctx sdk.Context, senderPool, recipientPool string, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error

	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error

	MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	BurnCoins(ctx sdk.Context, name string, amt sdk.Coins) error

	SetDenomMetaData(ctx sdk.Context, denomMetaData bankTypes.Metadata)

	// Only needed for simulation interface matching
	// TODO: Look into golang syntax to make this "Everything in stakingtypes.bankkeeper + extra funcs"
	// I think it has to do with listing another interface as the first line here?
	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
	SetBalances(ctx sdk.Context, addr sdk.AccAddress, balances sdk.Coins) error
	LockedCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	GetSupply(ctx sdk.Context) bankexported.SupplyI
	UndelegateCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	DelegateCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
}

// ParamSubspace defines the expected Subspace interfacace
type ParamSubspace interface {
	HasKeyTable() bool
	WithKeyTable(table paramTypes.KeyTable) paramTypes.Subspace
	Get(ctx sdk.Context, key []byte, ptr interface{})
	GetParamSet(ctx sdk.Context, ps paramTypes.ParamSet)
	SetParamSet(ctx sdk.Context, ps paramTypes.ParamSet)
}
