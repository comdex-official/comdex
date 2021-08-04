package keeper

import (
	"fmt"
	"github.com/comdex-official/comdex/x/cdp/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"sort"
)

func (k Keeper) AddCdp(ctx sdk.Context, owner sdk.AccAddress, collateral sdk.Coin, principal sdk.Coin, collateralType string) error {
	err := k.ValidateCollateral(ctx, collateral, collateralType)
	if err != nil {
		return err
	}
	err = k.ValidateBalance(ctx, collateral, owner)
	if err != nil {
		return err
	}

	err = k.ValidateCollateralizationRatio(ctx, collateral, collateralType, principal)
	if err != nil {
		return err
	}

	id := k.GetNextCdpID(ctx)

	cdp := types.NewCDP(id, owner, collateral, collateralType, principal)

	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, owner, types.ModuleName, sdk.NewCoins(collateral))
	if err != nil {
		return err
	}

	err = k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(principal))
	if err != nil {
		return err
	}

	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, owner, sdk.NewCoins(principal))
	if err != nil {
		return err
	}
	k.SetCdp(ctx, cdp)
	k.IndexCdpByOwner(ctx, cdp)
	k.SetNextCdpId(ctx, id+1)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCDPCreated,
			sdk.NewAttribute(types.AttributeKeyCdpID, fmt.Sprintf("%d", cdp.Id)),
		),
	)

	return nil
}

func (k Keeper) DepositCollateral(ctx sdk.Context, owner string, collateral sdk.Coin, collateralType string) error {
	err := k.ValidateCollateral(ctx, collateral, collateralType)
	if err != nil {
		return err
	}

	ownerAddrs, _ := sdk.AccAddressFromBech32(owner)
	cdp, found := k.GetCdpByOwnerAndCollateralType(ctx, ownerAddrs, collateralType)
	if !found {
		return sdkerrors.Wrapf(types.ErrorCdpNotFound, "owner %s, collateral %s", owner, collateral.Denom)
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
		return sdkerrors.Wrapf(types.ErrorCdpNotFound, "owner %s, collateral %s", owner, collateral.Denom)
	}

	if collateral.Amount.GT(cdp.Collateral.Amount) {
		return sdkerrors.Wrapf(types.ErrorInvalidWithdrawAmount, "collateral %s, deposit %s", collateral, cdp.Collateral.Amount)
	}

	liquidationRatio := k.GetLiquidationRatio(ctx, collateralType)
	collateralizationRatio, err := k.CalculateCollateralizationRatio(ctx, cdp.Collateral.Sub(collateral), cdp.Type, cdp.Principal, types.Spot)
	if err != nil {
		return err
	}

	if collateralizationRatio.LT(liquidationRatio) {
		return sdkerrors.Wrapf(types.ErrorInvalidCollateralRatio, "collateral %s, collateral ratio %s, liquidation ration %s", collateral.Denom, collateralizationRatio, liquidationRatio)
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

func (k Keeper) AddPrincipal(ctx sdk.Context, owner string, collateralType string, principal sdk.Coin) error {
	ownerAddrs, err := sdk.AccAddressFromBech32(owner)
	if err != nil {
		return err
	}
	cdp, found := k.GetCdpByOwnerAndCollateralType(ctx, ownerAddrs, collateralType)
	if !found {
		return sdkerrors.Wrapf(types.ErrorCdpNotFound, "owner %s, denom %s", owner, collateralType)
	}

	err = k.ValidatePrincipalDraw(ctx, principal, cdp.Principal.Denom)
	if err != nil {
		return err
	}

	err = k.ValidateCollateralizationRatio(ctx, cdp.Collateral, cdp.Type, cdp.Principal.Add(principal))
	if err != nil {
		return err
	}

	err = k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(principal))
	if err != nil {
		panic(err)
	}
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, ownerAddrs, sdk.NewCoins(principal))
	if err != nil {
		panic(err)
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCDPDraw,
			sdk.NewAttribute(sdk.AttributeKeyAmount, principal.String()),
			sdk.NewAttribute(types.AttributeKeyCdpID, fmt.Sprintf("%d", cdp.Id)),
		),
	)

	cdp.Principal = cdp.Principal.Add(principal)

	return k.SetCdp(ctx, cdp)

}

func (k Keeper) RepayPrincipal(ctx sdk.Context, owner string, collateralType string, payment sdk.Coin) error {
	ownerAddrs, err := sdk.AccAddressFromBech32(owner)
	if err != nil {
		return err
	}
	cdp, found := k.GetCdpByOwnerAndCollateralType(ctx, ownerAddrs, collateralType)
	if !found {
		return sdkerrors.Wrapf(types.ErrorCdpNotFound, "owner %s, denom %s", owner, collateralType)
	}

	err = k.ValidatePaymentCoins(ctx, cdp, payment)
	if err != nil {
		return err
	}

	err = k.ValidateBalance(ctx, payment, ownerAddrs)
	if err != nil {
		return err
	}

	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, ownerAddrs, types.ModuleName, sdk.NewCoins(payment))
	if err != nil {
		return err
	}

	err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(payment))
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

	cdp.Principal = cdp.Principal.Sub(payment)

	if cdp.Principal.IsZero() {
		k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, ownerAddrs, sdk.NewCoins(cdp.Collateral))
		err = k.DeleteCDP(ctx, cdp)
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

func (k Keeper) AttemptLiquidation(ctx sdk.Context, owner string, collateralType string) error {
	//TODO
	return nil
}

func (k Keeper) ValidateBalance(ctx sdk.Context, amount sdk.Coin, sender sdk.AccAddress) error {
	acc := k.accountKeeper.GetAccount(ctx, sender)
	if acc == nil {
		return sdkerrors.Wrapf(types.ErrorAccountNotFound, "address: %s", sender)
	}

	spendableBalance := k.bankKeeper.SpendableCoins(ctx, sender).AmountOf(amount.Denom)
	if spendableBalance.LT(amount.Amount) {
		return sdkerrors.Wrapf(types.ErrorInsufficientBalance, "%s < %s", sdk.NewCoin(amount.Denom, spendableBalance), amount)
	}

	return nil
}

func (k Keeper) ValidateCollateral(ctx sdk.Context, collateral sdk.Coin, collateralType string) error {
	collateralParam, found := k.GetCollateral(ctx, collateralType)
	if !found {
		return sdkerrors.Wrapf(types.ErrorCdpNotFound, "%s cdp does not exist", collateralType)
	}

	if collateralParam.Denom != collateral.Denom {
		return sdkerrors.Wrapf(types.ErrorInvalidCollateral, "collateral given %s , collateral required %s", collateral.Denom, collateralParam.Denom)
	}
	return nil
}

func (k Keeper) ValidateCollateralizationRatio(ctx sdk.Context, collateral sdk.Coin, collateralType string, principal sdk.Coin) error {
	collateralizationRatio, err := k.CalculateCollateralizationRatio(ctx, collateral, collateralType, principal, types.Spot)
	if err != nil {
		return err
	}
	liquidationRatio := k.GetLiquidationRatio(ctx, collateralType)

	if collateralizationRatio.LT(liquidationRatio) {
		return sdkerrors.Wrapf(types.ErrorInvalidCollateralRatio, "collateral %s, collateral ratio %s, liquidation ratio %s", collateral.Denom, collateralizationRatio, liquidationRatio)
	}

	return nil
}

func (k Keeper) GetCdpByOwnerAndCollateralType(ctx sdk.Context, owner sdk.AccAddress, collateralType string) (types.CDP, bool) {
	cdpIdList, found := k.GetCdpIdsByOnwer(ctx, owner)
	if !found {
		return types.CDP{}, false
	}

	for _, id := range cdpIdList.Ids {
		cdp, found := k.GetCDP(ctx, collateralType, id)
		if found {
			return cdp, true
		}
	}

	return types.CDP{}, false
}

func (k Keeper) GetCDP(ctx sdk.Context, collateralType string, cdpID uint64) (types.CDP, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.CdpKey(cdpID))
	if bz == nil {
		return types.CDP{}, false
	}

	var cdp types.CDP
	k.cdc.MustUnmarshalBinaryBare(bz, &cdp)
	return cdp, true
}

func (k Keeper) SetCdp(ctx sdk.Context, cdp types.CDP) error {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryBare(&cdp)
	store.Set(types.CdpKey(cdp.Id), bz)
	return nil
}

func (k Keeper) DeleteCDP(ctx sdk.Context, cdp types.CDP) error {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.CdpKey(cdp.Id))
	return nil
}

func (k Keeper) GetNextCdpID(ctx sdk.Context) (id uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.CdpIDKey)

	if bz == nil {
		panic("starting cdp id not set in genesis")
	}
	id = types.GetCdpIDFromBytes(bz)
	return
}

func (k Keeper) IndexCdpByOwner(ctx sdk.Context, cdp types.CDP) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.CdpIDIndexKeyPrefix)
	ownerAddrs, _ := sdk.AccAddressFromBech32(cdp.Owner)
	cdpIDs, found := k.GetCdpIdsByOnwer(ctx, ownerAddrs)
	if !found {
		idBytes := k.cdc.MustMarshalBinaryBare(&types.CdpIdList{[]uint64{cdp.Id}})
		store.Set(ownerAddrs, idBytes)
		return
	}

	cdpIDList := append(cdpIDs.Ids, cdp.Id)
	sort.Slice(cdpIDs, func(i, j int) bool { return cdpIDList[i] < cdpIDList[j] })
	cdpIDs.Ids = cdpIDList
	store.Set(ownerAddrs, k.cdc.MustMarshalBinaryBare(&cdpIDs))
}

func (k Keeper) SetNextCdpId(ctx sdk.Context, id uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.CdpIDKey, types.GetCdpIDBytes(id))
}

func (k Keeper) GetCdpIdsByOnwer(ctx sdk.Context, owner sdk.AccAddress) (types.CdpIdList, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.CdpIDIndexKeyPrefix)
	bz := store.Get(owner)
	if bz == nil {
		return types.CdpIdList{[]uint64{}}, false
	}

	var cdpIDs types.CdpIdList
	k.cdc.MustUnmarshalBinaryBare(bz, &cdpIDs)
	return cdpIDs, true
}

func (k Keeper) ValidatePrincipalDraw(ctx sdk.Context, principal sdk.Coin, expectedDenom string) error {
	if principal.Denom != expectedDenom {
		return sdkerrors.Wrapf(types.ErrorInvalidDebtRequest, "proposed %s, expected %s", principal.Denom, expectedDenom)
	}

	_, found := k.GetDebtParam(ctx, principal.Denom)
	if !found {
		return sdkerrors.Wrapf(types.ErrorDebtNotSupported, principal.Denom)
	}
	return nil
}

func (k Keeper) ValidatePaymentCoins(ctx sdk.Context, cdp types.CDP, payment sdk.Coin) error {
	if payment.Denom != cdp.Principal.Denom {
		return sdkerrors.Wrapf(types.ErrorInvalidPayment, "cdp %d: expected %s, got %s", cdp.Id, cdp.Principal.Denom, payment.Denom)
	}
	_, found := k.GetDebtParam(ctx, payment.Denom)
	if !found {
		return sdkerrors.Wrapf(types.ErrorInvalidPayment, "payment denom %s not found", payment.Denom)
	}
	return nil
}

func (k Keeper) CalculateCollateralizationRatio(ctx sdk.Context, collateral sdk.Coin, collateralType string, principal sdk.Coin, pfType types.PricefeedType) (sdk.Dec, error) {
	//TODO
	return sdk.NewDec(2), nil
}
