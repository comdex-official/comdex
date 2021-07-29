package keeper

import (
	"github.com/comdex-official/comdex/x/cdp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k Keeper) AddCdp(ctx sdk.Context, owner sdk.AccAddress, collateral sdk.Coin, principal sdk.Coin, collateralType string) error {
	// validation
	/*err := k.ValidateCollateral(ctx, collateral, collateralType)
	if err != nil {
		return err
	}*/
	err := k.ValidateBalance(ctx, collateral, owner)
	if err != nil {
		return err
	}

}

func (k Keeper) ValidateBalance(ctx sdk.Context, amount sdk.Coin, sender sdk.AccAddress) error {
	acc := k.accountKeeper.GetAccount(ctx, sender)
	if acc == nil {
		return sdkerrors.Wrapf(types.ErrAccountNotFound, "address: %s", sender)
	}

	spendableBalance := k.bankKeeper.SpendableCoins(ctx, sender).AmountOf(amount.Denom)
	if spendableBalance.LT(amount.Amount) {
		return sdkerrors.Wrapf(types.ErrInsufficientBalance, "%s < %s", sdk.NewCoin(amount.Denom, spendableBalance), amount)
	}

	return nil
}
