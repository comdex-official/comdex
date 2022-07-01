package expected

import (
	assettypes "github.com/comdex-official/comdex/x/asset/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type BankKeeper interface {
	BurnCoins(ctx sdk.Context, name string, coins sdk.Coins) error
	MintCoins(ctx sdk.Context, name string, coins sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, address sdk.AccAddress, name string, coins sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, name string, address sdk.AccAddress, coins sdk.Coins) error

	SendCoinsFromModuleToModule(
		ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins,
	) error

	SpendableCoins(ctx sdk.Context, address sdk.AccAddress) sdk.Coins
}

type AssetKeeper interface {
	GetApp(ctx sdk.Context, ID uint64) (assettypes.AppData, bool)
	GetApps(ctx sdk.Context) ([]assettypes.AppData, bool)
	GetAsset(ctx sdk.Context, ID uint64) (assettypes.Asset, bool)
	GetAssetForDenom(ctx sdk.Context, denom string) (assettypes.Asset, bool)
	GetMintGenesisTokenData(ctx sdk.Context, appID, assetID uint64) (assettypes.MintGenesisToken, bool)
}
