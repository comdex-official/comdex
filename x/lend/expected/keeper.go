package expected

import (
	"github.com/petrichormoney/petri/x/liquidation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	assettypes "github.com/petrichormoney/petri/x/asset/types"
	esmtypes "github.com/petrichormoney/petri/x/esm/types"
	markettypes "github.com/petrichormoney/petri/x/market/types"
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
	GetTwa(ctx sdk.Context, id uint64) (twa markettypes.TimeWeightedAverage, found bool)
	CalcAssetPrice(ctx sdk.Context, id uint64, amt sdk.Int) (price sdk.Dec, err error)
}

type BandOracleKeeper interface {
	GetOracleValidationResult(ctx sdk.Context) bool
}

type AssetKeeper interface {
	GetAsset(ctx sdk.Context, id uint64) (assettypes.Asset, bool)
	GetApp(ctx sdk.Context, id uint64) (assettypes.AppData, bool)
	SetApp(ctx sdk.Context, app assettypes.AppData)
	SetAppID(ctx sdk.Context, id uint64)
}

type EsmKeeper interface {
	GetKillSwitchData(ctx sdk.Context, appID uint64) (esmtypes.KillSwitchParams, bool)
}

type LiquidationKeeper interface {
	GetLockedVaultByApp(ctx sdk.Context, appID uint64) (lockedVault []types.LockedVault)
}

type AuctionKeeper interface {
	LendDutchActivator(ctx sdk.Context, lockedVault types.LockedVault) error
	StartLendDutchAuction(
		ctx sdk.Context,
		outFlowToken sdk.Coin,
		inFlowToken sdk.Coin,
		appID uint64,
		assetInID, assetOutID uint64,
		lockedVaultID uint64,
		lockedVaultOwner string,
		liquidationPenalty sdk.Dec,
	) error
}
