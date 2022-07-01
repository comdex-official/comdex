package expected

import (
	"github.com/comdex-official/comdex/x/asset/types"
	auctiontypes "github.com/comdex-official/comdex/x/auction/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type BankKeeper interface {
	SendCoinsFromAccountToModule(ctx sdk.Context, address sdk.AccAddress, name string, coins sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, name string, address sdk.AccAddress, coins sdk.Coins) error

	SendCoinsFromModuleToModule(
		ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins,
	) error
}

type AssetKeeper interface {
	HasAssetForDenom(ctx sdk.Context, id string) bool
	HasAsset(ctx sdk.Context, id uint64) bool
	GetAssetForDenom(ctx sdk.Context, denom string) (types.Asset, bool)
	GetApp(ctx sdk.Context, id uint64) (types.AppData, bool)
	GetAsset(ctx sdk.Context, id uint64) (types.Asset, bool)
	GetMintGenesisTokenData(ctx sdk.Context, appID, assetID uint64) (mintData types.MintGenesisToken, found bool)
}

type AuctionKeeper interface {
	GetAuctionParams(ctx sdk.Context, appID uint64) (asset auctiontypes.AuctionParams, found bool)
}
