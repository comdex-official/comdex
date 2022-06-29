package keeper

import (
	"fmt"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/comdex-official/comdex/x/nft/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Keeper struct {
	storeKey sdk.StoreKey
	cdc      codec.BinaryCodec
}

func NewKeeper(cdc codec.BinaryCodec, storeKey sdk.StoreKey) Keeper {
	return Keeper{
		storeKey: storeKey,
		cdc:      cdc,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) CreateDenom(
	ctx sdk.Context, id, symbol, name, schema string,
	creator sdk.AccAddress, description, previewUri string) error {

	if k.HasDenomID(ctx, id) {
		return sdkerrors.Wrapf(types.ErrInvalidDenom, "denomID %s has already exists", id)
	}

	if k.HasDenomSymbol(ctx, symbol) {
		return sdkerrors.Wrapf(types.ErrInvalidDenom, "denomSymbol %s has already exists", symbol)
	}
	err := k.SetDenom(ctx, types.NewDenom(id, symbol, name, schema, creator, description, previewUri))
	if err != nil {
		return err
	}
	k.setDenomOwner(ctx, id, creator)
	return nil
}

func (k Keeper) UpdateDenom(ctx sdk.Context, id, name, description, previewURI string, sender sdk.AccAddress) error {
	if !k.HasDenomID(ctx, id) {
		return sdkerrors.Wrapf(types.ErrInvalidDenom, "denom id %s not exists", id)
	}
	denom, err := k.AuthorizeDenomCreator(ctx, id, sender)
	if err != nil {
		return err
	}
	if len(name) > 0 && name != types.DoNotModify {
		denom.Name = name
	}
	if len(description) > 0 && description != types.DoNotModify {
		denom.Description = description
	}
	if len(previewURI) > 0 && previewURI != types.DoNotModify {
		denom.PreviewURI = previewURI
	}
	return k.SetDenom(ctx, denom)
}

func (k Keeper) TransferDenomOwner(ctx sdk.Context, id string, curOwner, newOwner sdk.AccAddress) error {
	denom, err := k.GetDenom(ctx, id)
	if err != nil {
		return sdkerrors.Wrapf(types.ErrInvalidDenom, "denom ID %s not exists", id)
	}

	if curOwner.String() != denom.Creator {
		return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "unauthorized address %s", curOwner.String())
	}

	denom.Creator = newOwner.String()

	err = k.SetDenom(ctx, denom)
	if err != nil {
		return err
	}
	k.swapDenomOwner(ctx, id, curOwner, newOwner)
	return nil
}

func (k Keeper) MintNFT(
	ctx sdk.Context, denomID, nftID string,
	metadata types.Metadata, data string, transferable, extensible, nsfw bool,
	royaltyShare sdk.Dec, sender, recipient sdk.AccAddress) error {
	if !k.HasDenomID(ctx, denomID) {
		return sdkerrors.Wrapf(types.ErrInvalidDenom, "denomID %s not exists", denomID)
	}

	if !k.HasPermissionToMint(ctx, denomID, sender) {
		return sdkerrors.Wrapf(types.ErrUnauthorized, "only creator of denom has permission to mint")
	}

	// not working need to make change by replacing the nftID with nft name
	if k.HasNFT(ctx, denomID, nftID) {
		return sdkerrors.Wrapf(types.ErrNFTAlreadyExists, "NFT %s already exists in collection %s", nftID, denomID)
	}

	k.setNFT(ctx, denomID, types.NewNFT(
		nftID,
		metadata,
		data,
		transferable,
		extensible,
		recipient,
		ctx.BlockHeader().Time,
		nsfw,
		royaltyShare,
	))
	k.setOwner(ctx, denomID, nftID, recipient)
	k.increaseSupply(ctx, denomID)
	return nil
}

func (k Keeper) EditNFT(ctx sdk.Context, denomID, nftID string, owner sdk.AccAddress) error {
	if !k.HasDenomID(ctx, denomID) {
		return sdkerrors.Wrapf(types.ErrInvalidDenom, "denomID %s not exists", denomID)
	}
	_, err := k.GetDenom(ctx, denomID)
	if err != nil {
		return err
	}

	nft, err := k.Authorize(ctx, denomID, nftID, owner)
	if err != nil {
		return err
	}

	k.setNFT(ctx, denomID, nft)
	return nil
}

func (k Keeper) TransferOwnership(ctx sdk.Context, denomID, nftID string, srcOwner, dstOwner sdk.AccAddress) error {
	if !k.HasDenomID(ctx, denomID) {
		return sdkerrors.Wrapf(types.ErrInvalidDenom, "denomID %s not exists", denomID)
	}

	nft, err := k.Authorize(ctx, denomID, nftID, srcOwner)
	if err != nil {
		return err
	}
	if !nft.IsTransferable() {
		return sdkerrors.Wrap(types.ErrNotTransferable, nft.GetID())
	}

	nft.Owner = dstOwner.String()

	k.setNFT(ctx, denomID, nft)
	k.swapOwner(ctx, denomID, nftID, srcOwner, dstOwner)
	return nil
}

func (k Keeper) BurnNFT(ctx sdk.Context,
	denomID, nftID string,
	owner sdk.AccAddress) error {
	if !k.HasDenomID(ctx, denomID) {
		return sdkerrors.Wrapf(types.ErrInvalidDenom, "denomID %s not exists", denomID)
	}

	nft, err := k.Authorize(ctx, denomID, nftID, owner)
	if err != nil {
		return err
	}

	k.deleteNFT(ctx, denomID, nft)
	k.deleteOwner(ctx, denomID, nftID, owner)
	k.decreaseSupply(ctx, denomID)
	return nil
}
