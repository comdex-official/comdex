package keeper

import (
	"github.com/comdex-official/comdex/x/cdp/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"sort"
)

func (k Keeper) AddCdp(ctx sdk.Context, owner sdk.AccAddress, collateral sdk.Coin, debt sdk.Coin, collateralType string) error {
	err := k.VerifyCollateralAndDebt(ctx, collateral, debt, collateralType)
	if err != nil {
		return err
	}

	err = k.VerifyBalance(ctx, collateral, owner)
	if err != nil {
		return err
	}

	err = k.VerifyCollateralizationRatio(ctx, collateral, debt, collateralType)
	if err != nil {
		return err
	}

	id := k.GetNextCdpID(ctx)

	cdp := types.NewCDP(id, owner, collateral, collateralType, debt)

	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, owner, types.ModuleName, sdk.NewCoins(collateral))
	if err != nil {
		return err
	}

	err = k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(debt))
	if err != nil {
		return err
	}

	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, owner, sdk.NewCoins(debt))
	if err != nil {
		return err
	}

	k.SetCDP(ctx, cdp)
	k.IndexCDPByOwner(ctx, cdp)
	k.SetNextCdpId(ctx, id+1)

	return nil
}

func (k Keeper) DepositCollateral(ctx sdk.Context, owner sdk.AccAddress, collateral sdk.Coin, collateralType string) error {
	cdp, found := k.GetCDPByOwnerAndCollateralType(ctx, owner, collateralType)
	if !found {
		return sdkerrors.Wrapf(types.ErrorCdpNotFound, "owner %s, collateral %s", owner, collateral.Denom)
	}

	if collateral.Denom != cdp.Collateral.Denom {
		return sdkerrors.Wrapf(types.ErrorInvalidCollateral, "collateral given %s , collateral required %s", collateral.Denom, cdp.Collateral.Denom)
	}

	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, owner, types.ModuleName, sdk.NewCoins(collateral))
	if err != nil {
		return err
	}

	cdp.Collateral = cdp.Collateral.Add(collateral)

	k.SetCDP(ctx, cdp)
	return nil
}

func (k Keeper) WithdrawCollateral(ctx sdk.Context, owner sdk.AccAddress, collateral sdk.Coin, collateralType string) error {
	cdp, found := k.GetCDPByOwnerAndCollateralType(ctx, owner, collateralType)
	if !found {
		return sdkerrors.Wrapf(types.ErrorCdpNotFound, "owner %s, collateral %s", owner, collateral.Denom)
	}

	if collateral.Denom != cdp.Collateral.Denom {
		return sdkerrors.Wrapf(types.ErrorInvalidCollateral, "collateral given %s , collateral required %s", collateral.Denom, cdp.Collateral.Denom)
	}

	if collateral.Amount.GT(cdp.Collateral.Amount) {
		return sdkerrors.Wrapf(types.ErrorInvalidWithdrawAmount, "collateral %s, deposited %s", collateral, cdp.Collateral.Amount)
	}

	liquidationRatio, err := k.GetLiquidationRatio(ctx, collateralType)
	if err != nil {
		return err
	}
	collateralizationRatio, err := k.CalculateCollateralizationRatio(ctx, cdp.Collateral.Sub(collateral), cdp.Type, cdp.Debt, types.Spot)
	if err != nil {
		return err
	}

	if collateralizationRatio.LT(liquidationRatio) {
		return sdkerrors.Wrapf(types.ErrorInvalidCollateralRatio, "collateral %s, collateral ratio %s, liquidation ration %s", collateral.Denom, collateralizationRatio, liquidationRatio)
	}

	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, owner, sdk.NewCoins(collateral))
	if err != nil {
		return err
	}

	cdp.Collateral = cdp.Collateral.Sub(collateral)

	k.SetCDP(ctx, cdp)
	return nil

}

func (k Keeper) DrawDebt(ctx sdk.Context, owner sdk.AccAddress, collateralType string, debt sdk.Coin) error {
	cdp, found := k.GetCDPByOwnerAndCollateralType(ctx, owner, collateralType)
	if !found {
		return sdkerrors.Wrapf(types.ErrorCdpNotFound, "owner %s, denom %s", owner, collateralType)
	}

	if debt.Denom != cdp.Debt.Denom {
		return sdkerrors.Wrapf(types.ErrorInvalidDebt, "requested %s, expected %s", debt.Denom, cdp.Debt.Denom)
	}

	err := k.VerifyCollateralizationRatio(ctx, cdp.Collateral, cdp.Debt.Add(debt), cdp.Type)
	if err != nil {
		return err
	}

	err = k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(debt))
	if err != nil {
		return err
	}
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, owner, sdk.NewCoins(debt))
	if err != nil {
		return err
	}

	cdp.Debt = cdp.Debt.Add(debt)
	k.SetCDP(ctx, cdp)

	return nil
}

func (k Keeper) RepayDebt(ctx sdk.Context, owner sdk.AccAddress, collateralType string, debt sdk.Coin) error {
	cdp, found := k.GetCDPByOwnerAndCollateralType(ctx, owner, collateralType)
	if !found {
		return sdkerrors.Wrapf(types.ErrorCdpNotFound, "owner %s, denom %s", owner, collateralType)
	}

	if cdp.Debt.Denom != debt.Denom {
		return sdkerrors.Wrapf(types.ErrorInvalidPayment, "cdp %d: expected %s, got %s", cdp.Id, cdp.Debt.Denom, debt.Denom)
	}

	if debt.Amount.GT(cdp.Debt.Amount) {
		return sdkerrors.Wrapf(types.ErrorInvalidAmount, "debt amount %s greater than present in cdp %s", debt.Amount, cdp.Debt.Amount)
	}

	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, owner, types.ModuleName, sdk.NewCoins(debt))
	if err != nil {
		return err
	}

	err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(debt))
	if err != nil {
		return err
	}

	cdp.Debt = cdp.Debt.Sub(debt)

	if cdp.Debt.IsZero() {
		err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, owner, sdk.NewCoins(cdp.Collateral))
		if err != nil {
			return err
		}
		err = k.DeleteCDP(ctx, cdp)
		if err != nil {
			return err
		}
		return nil
	}

	k.SetCDP(ctx, cdp)

	return nil
}

func (k Keeper) AttemptLiquidation(ctx sdk.Context, owner sdk.AccAddress, collateralType string) error {
	cdp, found := k.GetCDPByOwnerAndCollateralType(ctx, owner, collateralType)
	if !found {
		return sdkerrors.Wrapf(types.ErrorCdpNotFound, "owner %s, denom %s", owner, collateralType)
	}
	err := k.VerifyLiquidation(ctx, cdp.Collateral, cdp.Debt, cdp.Type)
	if err != nil {
		return err
	}

	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, owner, types.ModuleName, sdk.NewCoins(cdp.Debt))
	if err != nil {
		return err
	}
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, owner, sdk.NewCoins(cdp.Collateral))
	if err != nil {
		return err
	}
	err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(cdp.Debt))
	if err != nil {
		return err
	}
	return k.DeleteCDP(ctx, cdp)
}

func (k Keeper) VerifyBalance(ctx sdk.Context, amount sdk.Coin, sender sdk.AccAddress) error {
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

func (k Keeper) VerifyCollateralAndDebt(ctx sdk.Context, collateral sdk.Coin, debt sdk.Coin, collateralType string) error {
	collateralParam, found := k.GetCollateralParam(ctx, collateralType)
	if !found {
		return sdkerrors.Wrapf(types.ErrorCdpNotFound, "%s cdp does not exist", collateralType)
	}

	if collateralParam.CollateralDenom != collateral.Denom {
		return sdkerrors.Wrapf(types.ErrorInvalidCollateral, "collateral given %s , collateral required %s", collateral.Denom, collateralParam.CollateralDenom)
	}

	if collateralParam.DebtDenom != debt.Denom {
		return sdkerrors.Wrapf(types.ErrorInvalidDebt, "collateral given %s , collateral required %s", collateral.Denom, collateralParam.CollateralDenom)
	}
	return nil
}

func (k Keeper) VerifyCollateralizationRatio(ctx sdk.Context, collateral sdk.Coin, debt sdk.Coin, collateralType string) error {
	collateralizationRatio, err := k.CalculateCollateralizationRatio(ctx, collateral, collateralType, debt, types.Spot)
	if err != nil {
		return err
	}
	liquidationRatio, err := k.GetLiquidationRatio(ctx, collateralType)
	if err != nil {
		return err
	}
	if collateralizationRatio.LT(liquidationRatio) {
		return sdkerrors.Wrapf(types.ErrorInvalidCollateralRatio, "collateral %s, collateral ratio %s, liquidation ratio %s", collateral.Denom, collateralizationRatio, liquidationRatio)
	}

	return nil
}

func (k Keeper) GetCDPByOwnerAndCollateralType(ctx sdk.Context, owner sdk.AccAddress, collateralType string) (types.CDP, bool) {
	ownerCDPList, found := k.GetOwnerCDPList(ctx, owner)
	if !found {
		return types.CDP{}, false
	}

	for _, ownedCDP := range ownerCDPList.OwnedCDPs {
		if collateralType == ownedCDP.CollateralType {
			cdp, found := k.GetCDP(ctx, ownedCDP.Id)
			if found {
				return cdp, true
			}
		}
	}

	return types.CDP{}, false
}

func (k Keeper) GetCDP(ctx sdk.Context, cdpID uint64) (types.CDP, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.CdpKey(cdpID))
	if bz == nil {
		return types.CDP{}, false
	}

	var cdp types.CDP
	k.cdc.MustUnmarshalBinaryBare(bz, &cdp)
	return cdp, true
}

func (k Keeper) SetCDP(ctx sdk.Context, cdp types.CDP) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryBare(&cdp)
	store.Set(types.CdpKey(cdp.Id), bz)
	return
}

func (k Keeper) DeleteCDP(ctx sdk.Context, cdp types.CDP) error {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.CdpKey(cdp.Id))
	return nil
}

func (k Keeper) GetNextCdpID(ctx sdk.Context) (id uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.CdpIdKey)

	if bz == nil {
		return 1
	}
	id = types.GetCdpIDFromBytes(bz)
	return id
}

func (k Keeper) IndexCDPByOwner(ctx sdk.Context, cdp types.CDP) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.CdpIdIndexKeyPrefix)
	ownerAddrs, _ := sdk.AccAddressFromBech32(cdp.Owner)
	ownerCDPList, found := k.GetOwnerCDPList(ctx, ownerAddrs)
	ownedCDP := types.OwnedCDP{Id: cdp.Id, CollateralType: cdp.Type}
	if !found {
		idBytes := k.cdc.MustMarshalBinaryBare(&types.OwnerCDPList{[]types.OwnedCDP{ownedCDP}})
		store.Set(ownerAddrs, idBytes)
		return
	}

	ownedCDPs := append(ownerCDPList.OwnedCDPs, ownedCDP)
	sort.Slice(ownedCDPs, func(i, j int) bool { return ownedCDPs[i].Id < ownedCDPs[j].Id })
	ownerCDPList.OwnedCDPs = ownedCDPs
	store.Set(ownerAddrs, k.cdc.MustMarshalBinaryBare(&ownerCDPList))
}

func (k Keeper) SetNextCdpId(ctx sdk.Context, id uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.CdpIdKey, types.GetCdpIDBytes(id))
}

func (k Keeper) GetOwnerCDPList(ctx sdk.Context, owner sdk.AccAddress) (types.OwnerCDPList, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.CdpIdIndexKeyPrefix)
	bz := store.Get(owner)
	if bz == nil {
		return types.OwnerCDPList{[]types.OwnedCDP{}}, false
	}

	var ownerCDPList types.OwnerCDPList
	k.cdc.MustUnmarshalBinaryBare(bz, &ownerCDPList)
	return ownerCDPList, true
}

func (k Keeper) VerifyLiquidation(ctx sdk.Context, collateral sdk.Coin, debt sdk.Coin, collateralType string) error {
	collateralizationRatio, err := k.CalculateCollateralizationRatio(ctx, collateral, collateralType, debt, types.Spot)
	if err != nil {
		return err
	}
	liquidationRatio, err := k.GetLiquidationRatio(ctx, collateralType)
	if err != nil {
		return err
	}
	if collateralizationRatio.GT(liquidationRatio) {
		return sdkerrors.Wrapf(types.ErrorLowCollateralizationRatio, "collateral %s, collateral ratio %s, liquidation ratio %s", collateral.Denom, collateralizationRatio, liquidationRatio)
	}
	return nil
}

func (k Keeper) CalculateCollateralizationRatio(ctx sdk.Context, collateral sdk.Coin, collateralType string, debt sdk.Coin, pfType types.PricefeedType) (sdk.Dec, error) {
	//TODO update when the price of token is available from oracle
	collateralTokenPrice := sdk.NewInt(5)
	debtTokenPrice := sdk.NewInt(1)
	collateralWorth := sdk.NewDecFromInt(collateral.Amount.Mul(collateralTokenPrice))
	debtPrice := sdk.NewDecFromInt(debt.Amount.Mul(debtTokenPrice))

	return collateralWorth.Quo(debtPrice), nil
}
