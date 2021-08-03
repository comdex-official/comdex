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

	cdp := types.NewCDP(id, owner, collateral, collateralType, principal, ctx.BlockHeader().Time)

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

func (k Keeper) AttemptLiquidation(ctx sdk.Context, owner string, collateralType string) error {
	//TODO
	return nil
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

func (k Keeper) ValidateCollateral(ctx sdk.Context, collateral sdk.Coin, collateralType string) error {
	collateralParam, found := k.GetCollateral(ctx, collateralType)
	if !found {
		return sdkerrors.Wrapf(types.ErrCdpNotFound, "%s cdp does not exist", collateralType)
	}

	if collateralParam.Denom != collateral.Denom {
		return sdkerrors.Wrapf(types.ErrInvalidCollateral, "collateral given %s , collateral required %s", collateral.Denom, collateralParam.Denom)
	}
	return nil
}

func (k Keeper) ValidateCollateralizationRatio(ctx sdk.Context, collateral sdk.Coin, collateralType string, principal sdk.Coin) error {
	collateralizationRatio, err := k.CalculateCollateralizationRatio(ctx, collateral, collateralType, principal, spot)
	if err != nil {
		return err
	}
	liquidationRatio := k.getLiquidationRatio(ctx, collateralType)

	if collateralizationRatio.LT(liquidationRatio) {
		return sdkerrors.Wrapf(types.ErrInvalidCollateralRatio, "collateral %s, collateral ratio %s, liquidation ratio %s", collateral.Denom, collateralizationRatio, liquidationRatio)
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
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.CdpKeyPrefix)
	cdpPrefix, found := k.GetCollateralTypePrefix(ctx, collateralType)

	if !found {
		return types.CDP{}, false
	}

	bz := store.Get(types.CdpKey(cdpPrefix, cdpID))

	if bz == nil {
		return types.CDP{}, false
	}
	var cdp types.CDP
	k.cdc.MustUnmarshalBinaryBare(bz, &cdp)
	return cdp, true
}

func (k Keeper) SetCdp(ctx sdk.Context, cdp types.CDP) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.CdpKeyPrefix)
	cdpPrefix, found := k.GetCollateralTypePrefix(ctx, cdp.Type)
	if !found {
		sdkerrors.Wrapf(types.ErrDenomPrefixNotFound, "%s", cdp.Collateral.Denom)
	}
	bz := k.cdc.MustMarshalBinaryBare(&cdp)
	store.Set(types.CdpKey(cdpPrefix, cdp.Id), bz)
	return nil
}

func (k Keeper) DeleteCDP(ctx sdk.Context, cdp types.CDP) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.CdpKeyPrefix)
	db, found := k.GetCollateralTypePrefix(ctx, cdp.Type)
	if !found {
		return sdkerrors.Wrapf(types.ErrDenomPrefixNotFound, "%s", cdp.Collateral.Denom)
	}
	store.Delete(types.CdpKey(db, cdp.Id))
	return nil
}

func (k Keeper) GetNextCdpID(ctx sdk.Context) (id uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.CdpIDKey)
	bz := store.Get([]byte{})

	if bz == nil {
		panic("starting cdp id not set in genesis")
	}
	id = types.GetCdpIDFromBytes(bz)
	return
}

func (k Keeper) MintDebtCoins(ctx sdk.Context, moduleAccount string, denom string, principalCoins sdk.Coin) error {
	debtCoins := sdk.NewCoins(sdk.NewCoin(denom, principalCoins.Amount))
	return k.bankKeeper.MintCoins(ctx, moduleAccount, debtCoins)
}

func (k Keeper) IndexCdpByOwner(ctx sdk.Context, cdp types.CDP) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.CdpIDKeyPrefix)
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
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.CdpIDKey)
	store.Set([]byte{}, types.GetCdpIDBytes(id))
}

func (k Keeper) GetCdpIdsByOnwer(ctx sdk.Context, owner sdk.AccAddress) (types.CdpIdList, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.CdpIDKeyPrefix)
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
		return sdkerrors.Wrapf(types.ErrInvalidDebtRequest, "proposed %s, expected %s", principal.Denom, expectedDenom)
	}

	_, found := k.GetDebtParam(ctx, principal.Denom)
	if !found {
		return sdkerrors.Wrapf(types.ErrDebtNotSupported, principal.Denom)
	}
	return nil
}

func (k Keeper) CalculateCollateralizationRatio(ctx sdk.Context, collateral sdk.Coin, collateralType string, principal sdk.Coin, pfType pricefeedType) (sdk.Dec, error) {
	//TODO
	return sdk.NewDec(2), nil
}

type pricefeedType string

const (
	spot        pricefeedType = "spot"
	liquidation pricefeedType = "liquidation"
)

func (pft pricefeedType) IsValid() error {
	switch pft {
	case spot, liquidation:
		return nil
	}
	return fmt.Errorf("invalid pricefeed type: %s", pft)
}
