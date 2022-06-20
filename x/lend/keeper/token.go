package keeper

import (
	"github.com/comdex-official/comdex/x/lend/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k Keeper) ExchangeToken(ctx sdk.Context, token sdk.Coin, tokenName string) (sdk.Coin, error) {
	if !token.IsValid() {
		return sdk.Coin{}, sdkerrors.Wrap(types.ErrInvalidAsset, token.String())
	}

	cTokenDenom := k.FromTokenToCTokenDenom(ctx, tokenName)
	if cTokenDenom == "" {
		return sdk.Coin{}, sdkerrors.Wrap(types.ErrInvalidAsset, tokenName)
	}
	uTokenAmount := token.Amount
	return sdk.NewCoin(cTokenDenom, uTokenAmount), nil
}

func (k Keeper) FromTokenToCTokenDenom(_ sdk.Context, tokenDenom string) string {
	return types.CTokenPrefix + tokenDenom
}

func (k Keeper) GetCTokenSupply(ctx sdk.Context, denom string) sdk.Coin {
	store := ctx.KVStore(k.storeKey)
	key := types.CreateCTokenSupplyKey(denom)
	amount := sdk.ZeroInt()

	if bz := store.Get(key); bz != nil {
		err := amount.Unmarshal(bz)
		if err != nil {
			panic(err)
		}
	}

	return sdk.NewCoin(denom, amount)
}

// setUTokenSupply sets the total supply of a uToken.
func (k Keeper) setCTokenSupply(ctx sdk.Context, coin sdk.Coin) error {
	if !coin.IsValid() {
		return sdkerrors.Wrap(types.ErrInvalidAsset, coin.String())
	}

	key := types.CreateCTokenSupplyKey(coin.Denom)

	// save the new reserve amount
	bz, err := coin.Amount.Marshal()
	if err != nil {
		return err
	}

	ctx.KVStore(k.storeKey).Set(key, bz)
	return nil
}
