package expected

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	assettypes "github.com/comdex-official/comdex/x/asset/types"
	esmtypes "github.com/comdex-official/comdex/x/esm/types"
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
	GetApp(ctx sdk.Context, id uint64) (assettypes.AppMapping, bool)
	GetApps(ctx sdk.Context) ([]assettypes.AppMapping, bool)
	GetAsset(ctx sdk.Context, id uint64) (assettypes.Asset, bool)
	GetMintGenesisTokenData(ctx sdk.Context, appId, assetId uint64) (assettypes.MintGenesisToken,bool)
}

type EsmKeeper interface {
	GetTriggerEsm(ctx sdk.Context, appid uint64) (esmtypes.EsmActive, bool)
}