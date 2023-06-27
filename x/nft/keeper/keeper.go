package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

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

	err := k.SetDenom(ctx, types.NewDenom(id , symbol, name, schema, creator, description, previewUri))
	if err != nil {
		return err
	}
	k.setDenomOwner(ctx, id, creator)
	return nil

}
