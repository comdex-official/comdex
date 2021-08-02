package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/x/cdp/types"
)

type (
	Keeper struct {
		cdc           codec.Marshaler
		storeKey      sdk.StoreKey
		memKey        sdk.StoreKey
		accountKeeper types.AccountKeeper
		bankKeeper    types.BankKeeper
		paramSpace    types.ParamSubspace
	}
)

func NewKeeper(
	cdc codec.Marshaler,
	storeKey,
	memKey sdk.StoreKey,

) *Keeper {
	return &Keeper{
		cdc:      cdc,
		storeKey: storeKey,
		memKey:   memKey,
	}
}

func (k Keeper) IncrementTotalPrincipal(ctx sdk.Context, collateralType string, principal sdk.Coin) {
	total := k.GetTotalPrincipal(ctx, collateralType, principal.Denom)
	total = total.Add(principal.Amount)
	k.SetTotalPrincipal(ctx, collateralType, principal.Denom, total)
}

func (k Keeper) GetTotalPrincipal(ctx sdk.Context, collateralType, principalDenom string) (total sdk.Int) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PrincipalKeyPrefix)
	bz := store.Get([]byte(collateralType + principalDenom))

	if bz == nil {
		k.SetTotalPrincipal(ctx, collateralType, principalDenom, sdk.ZeroInt())
		return sdk.ZeroInt()
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &total)
	return
}

func (k Keeper) SetTotalPrincipal(ctx sdk.Context, collateralType, principalDenom string, total sdk.Int) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PrincipalKeyPrefix)
	_, found := k.GetCollateralTypePrefix(ctx, collateralType)
	if !found {
		fmt.Sprintf("collateral not found")
	}
	store.Set([]byte(collateralType+principalDenom), k.cdc.MustMarshalBinaryLengthPrefixed(total))

}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
