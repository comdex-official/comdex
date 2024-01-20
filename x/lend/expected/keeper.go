package expected

import (
	"context"
	auctiontypes "github.com/comdex-official/comdex/x/auction/types"
	"github.com/comdex-official/comdex/x/liquidation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	sdkmath "cosmossdk.io/math"
	assettypes "github.com/comdex-official/comdex/x/asset/types"
	esmtypes "github.com/comdex-official/comdex/x/esm/types"
	markettypes "github.com/comdex-official/comdex/x/market/types"
)

type BankKeeper interface {
	BurnCoins(ctx context.Context, name string, coins sdk.Coins) error
	MintCoins(ctx context.Context, name string, coins sdk.Coins) error
	SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx context.Context, name string, address sdk.AccAddress, coins sdk.Coins) error
	SpendableCoins(ctx context.Context, address sdk.AccAddress) sdk.Coins
	GetAllBalances(ctx context.Context, addr sdk.AccAddress) sdk.Coins
	SendCoinsFromModuleToModule(
		ctx context.Context, senderModule, recipientModule string, amt sdk.Coins,
	) error
	SendCoins(ctx context.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) error
	GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin 
}

type AccountKeeper interface {
	GetModuleAddress(name string) sdk.AccAddress
	GetAccount(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
}

type MarketKeeper interface {
	GetTwa(ctx sdk.Context, id uint64) (twa markettypes.TimeWeightedAverage, found bool)
	CalcAssetPrice(ctx sdk.Context, id uint64, amt sdkmath.Int) (price sdkmath.LegacyDec, err error)
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
	GetLockedVault(ctx sdk.Context, appID, id uint64) (lockedVault types.LockedVault, found bool)
	DeleteLockedVault(ctx sdk.Context, appID, id uint64)
}

type AuctionKeeper interface {
	GetDutchLendAuctions(ctx sdk.Context, appID uint64) (auctions []auctiontypes.DutchAuction)
	SetHistoryDutchLendAuction(ctx sdk.Context, auction auctiontypes.DutchAuction) error
	DeleteDutchLendAuction(ctx sdk.Context, auction auctiontypes.DutchAuction) error
}
