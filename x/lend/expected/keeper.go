package expected

import (
	"github.com/comdex-official/comdex/x/asset/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

type BankKeeper interface {
	BurnCoins(ctx sdk.Context, name string, coins sdk.Coins) error
	MintCoins(ctx sdk.Context, name string, coins sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, address sdk.AccAddress, name string, coins sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, name string, address sdk.AccAddress, coins sdk.Coins) error
	SpendableCoins(ctx sdk.Context, address sdk.AccAddress) sdk.Coins
	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	SendCoinsFromModuleToModule(
		ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins,
	) error
	SendCoins(ctx sdk.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) error
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
}

type AccountKeeper interface {
	GetModuleAddress(name string) sdk.AccAddress
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
}

type MarketKeeper interface {
	GetPriceForAsset(ctx sdk.Context, id uint64) (uint64, bool)
}

type BandoracleKeeper interface {
	GetOracleValidationResult(ctx sdk.Context) bool
}

type AssetKeeper interface {
	GetWhitelistAsset(ctx sdk.Context, id uint64) (asset types.ExtendedAsset, found bool)
	GetWhitelistPair(ctx sdk.Context, id uint64) (pair types.ExtendedPairLend, found bool)
	GetPair(ctx sdk.Context, id uint64) (pair types.Pair, found bool)
	GetAsset(ctx sdk.Context, id uint64) (asset types.Asset, found bool)
}

type Marketkeeper interface {
	GetPriceForAsset(ctx sdk.Context, id uint64) (uint64, bool)
}
