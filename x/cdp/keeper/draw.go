package keeper

import (
	"fmt"
	"github.com/comdex-official/comdex/x/cdp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k Keeper) AddPrincipal(ctx sdk.Context, owner string, collateralType string, principal sdk.Coin) error {
	ownerAddrs,_ := sdk.AccAddressFromBech32(owner)
	cdp, found:= k.GetCdpByOwnerAndCollateralType(ctx, ownerAddrs, collateralType)
	if !found {
		return sdkerrors.Wrapf(types.ErrCdpNotFound, "owner %s, denom %s", owner, collateralType)
	}

	err:= k.ValidatePrincipalDraw(ctx, principal, cdp.Principal.Denom)
	if err != nil {
		return err
	}

	err= k.ValidateCollateralizationRatio(ctx, cdp.Collateral, cdp.Type, cdp.Principal.Add(principal))
	if err != nil {
		return err
	}

	err= k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(principal))
	if err != nil {
		panic( err)
	}
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, ownerAddrs, sdk.NewCoins(principal))
	if err != nil {
		panic( err)
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCDPDraw,
			sdk.NewAttribute(sdk.AttributeKeyAmount, principal.String()),
			sdk.NewAttribute(types.AttributeKeyCdpID, fmt.Sprintf("%d", cdp.Id)),
		),
	)

	cdp.Principal= cdp.Principal.Add(principal)

	return k.SetCdp(ctx, cdp)

}

func (k Keeper) RepayPrincipal(ctx sdk.Context, owner sdk.AccAddress, collateralType string, payment sdk.Coin) error {
	cdp, found:= k.GetCdpByOwnerAndCollateralType(ctx, owner, collateralType)
	if !found{
		return sdkerrors.Wrapf(types.ErrCdpNotFound, "owner %s, denom %s", owner, collateralType)
	}

	err:= k.ValidatePaymentCoins(ctx, cdp, payment)
	if err != nil {
		return err
	}

	err= k.ValidateBalance(ctx, payment, owner)
	if err != nil {
		return err
	}

	err= k.bankKeeper.SendCoinsFromAccountToModule(ctx, owner, types.ModuleName, sdk.NewCoins(payment))
	if err != nil {
		return err
	}

	err= k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(payment))
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCDPRepay,
			sdk.NewAttribute(sdk.AttributeKeyAmount, payment.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyCdpID, fmt.Sprintf("%d", cdp.Id)),
			),
		)

	cdp.Principal= cdp.Principal.Sub(payment)


	if cdp.Principal.IsZero(){
		k.ReturnCollateral(ctx, cdp)
		k.RemoveCdpOwnerIndex(ctx, cdp)
		err= k.DeleteCdpAndCollateralRatioIndex(ctx, cdp)
		if err != nil {
			return err
		}
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeCdpClose,
				sdk.NewAttribute(types.AttributeKeyCdpID, fmt.Sprintf("%d", cdp.Id)),
			),
		)
		return nil
	}


	return k.SetCdp(ctx, cdp)
}