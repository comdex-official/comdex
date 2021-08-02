package keeper

import (
	"fmt"
	"github.com/comdex-official/comdex/x/cdp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k Keeper) DepositCollateral(ctx sdk.Context, owner string, collateral sdk.Coin, collateralType string) error {
	err := k.ValidateCollateral(ctx, collateral, collateralType)
	if err != nil {
		return err
	}

	ownerAddrs, _ := sdk.AccAddressFromBech32(owner)
	cdp, found := k.GetCdpByOwnerAndCollateralType(ctx, ownerAddrs, collateralType)
	if !found {
		return sdkerrors.Wrapf(types.ErrCdpNotFound, "owner %s, collateral %s", owner, collateral.Denom)
	}

	err = k.ValidateBalance(ctx, collateral, ownerAddrs)
	if err != nil {
		return err
	}

	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, ownerAddrs, types.ModuleName, sdk.NewCoins(collateral))
	if err != nil {
		return err
	}

	cdp.Collateral = cdp.Collateral.Add(collateral)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCDPDeposit,
			sdk.NewAttribute(sdk.AttributeKeyAmount, collateral.String()),
			sdk.NewAttribute(types.AttributeKeyCdpID, fmt.Sprintf("%d", cdp.Id)),
		),
	)

	return k.SetCdp(ctx, cdp)
}

func (k Keeper) WithdrawCollateral(ctx sdk.Context, owner string, collateral sdk.Coin, collateralType string) error {
	err := k.ValidateCollateral(ctx, collateral, collateralType)
	if err != nil {
		return err
	}

	ownerAddrs, err := sdk.AccAddressFromBech32(owner)
	if err != nil {
		return err
	}

	cdp, found := k.GetCdpByOwnerAndCollateralType(ctx, ownerAddrs, collateralType)
	if !found {
		return sdkerrors.Wrapf(types.ErrCdpNotFound, "owner %s, collateral %s", owner, collateral.Denom)
	}

	if collateral.Amount.GT(cdp.Collateral.Amount) {
		return sdkerrors.Wrapf(types.ErrInvalidWithdrawAmount, "collateral %s, deposit %s", collateral, cdp.Collateral.Amount)
	}

	liquidationRatio := k.getLiquidationRatio(ctx, collateralType)
	collateralizationRatio, err := k.CalculateCollateralizationRatio(ctx, cdp.Collateral.Sub(collateral), cdp.Type, cdp.Principal, spot)
	if err != nil {
		return err
	}

	if collateralizationRatio.LT(liquidationRatio) {
		return sdkerrors.Wrapf(types.ErrInvalidCollateralRatio, "collateral %s, collateral ratio %s, liquidation ration %s", collateral.Denom, collateralizationRatio, liquidationRatio)
	}

	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, ownerAddrs, sdk.NewCoins(collateral))
	if err != nil {
		return err
	}

	cdp.Collateral = cdp.Collateral.Sub(collateral)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCDPWithdrawal,
			sdk.NewAttribute(sdk.AttributeKeyAmount, collateral.String()),
			sdk.NewAttribute(types.AttributeKeyCdpID, fmt.Sprintf("%d", cdp.Id)),
		),
	)

	return k.SetCdp(ctx, cdp)

}
