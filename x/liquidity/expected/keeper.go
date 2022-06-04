package expected

import (
	"time"

	assettypes "github.com/comdex-official/comdex/x/asset/types"
	rewardstypes "github.com/comdex-official/comdex/x/rewards/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// AccountKeeper is the expected account keeper.
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
	GetModuleAddress(moduleName string) sdk.AccAddress
}

// BankKeeper is the expected bank keeper.
type BankKeeper interface {
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	GetSupply(ctx sdk.Context, denom string) sdk.Coin
	IterateTotalSupply(ctx sdk.Context, cb func(sdk.Coin) bool)
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	SendCoins(ctx sdk.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	InputOutputCoins(ctx sdk.Context, inputs []banktypes.Input, outputs []banktypes.Output) error
}

type AssetKeeper interface {
	HasAssetForDenom(ctx sdk.Context, denom string) bool
	GetAssetForDenom(ctx sdk.Context, denom string) (asset assettypes.Asset, found bool)
}

type MarketKeeper interface {
	GetPriceForAsset(ctx sdk.Context, id uint64) (uint64, bool)
}

type RewardsKeeper interface {
	GetAllGaugesByGaugeTypeID(ctx sdk.Context, gaugeTypeID uint64) (gauges []rewardstypes.Gauge)
	GetEpochInfoByDuration(ctx sdk.Context, duration time.Duration) (epochInfo rewardstypes.EpochInfo, found bool)
	CreateNewGauge(ctx sdk.Context, msg *rewardstypes.MsgCreateGauge, forSwapFee bool) error
}
