package keeper

import (
	auctiontypes "github.com/comdex-official/comdex/x/auction/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetUUSDFromUSD(ctx sdk.Context, price sdk.Dec) sdk.Dec {
	usdInUUSD := sdk.MustNewDecFromStr("1000000")
	return price.Mul(usdInUUSD)
}
func (k Keeper) GetModuleAccountBalance(ctx sdk.Context, moduleName string, denom string) sdk.Int {
	address := k.account.GetModuleAddress(moduleName)
	return k.bank.GetBalance(ctx, address, denom).Amount
}

func (k Keeper) IncreaseLockedVaultAmountIn(ctx sdk.Context, lockedVaultId uint64, amount sdk.Int) error {
	lockedVault, found := k.GetLockedVault(ctx, lockedVaultId)
	if !found {
		return auctiontypes.ErrorVaultNotFound
	}
	lockedVault.AmountIn = lockedVault.AmountIn.Add(amount)
	k.SetLockedVault(ctx, lockedVault)
	return nil
}

func (k Keeper) DecreaseLockedVaultAmountIn(ctx sdk.Context, lockedVaultId uint64, amount sdk.Int) (isZero bool, err error) {
	lockedVault, found := k.GetLockedVault(ctx, lockedVaultId)
	if !found {
		return false, auctiontypes.ErrorVaultNotFound
	}
	lockedVault.AmountIn = lockedVault.AmountIn.Sub(amount)
	k.SetLockedVault(ctx, lockedVault)
	if lockedVault.AmountIn.IsZero() {
		return true, nil
	}
	return false, nil
}

func (k Keeper) DecreaseLockedVaultAmountOut(ctx sdk.Context, lockedVaultId uint64, amount sdk.Int) error {
	lockedVault, found := k.GetLockedVault(ctx, lockedVaultId)
	if !found {
		return auctiontypes.ErrorVaultNotFound
	}
	lockedVault.AmountIn = lockedVault.AmountOut.Sub(amount)
	k.SetLockedVault(ctx, lockedVault)
	return nil
}

func (k Keeper) AddAuctionParams(ctx sdk.Context, appId, auctionDurationSeconds uint64, buffer, cusp sdk.Dec, step, priceFunctionType, surplusId, debtId, dutchId uint64, bidDurationSeconds uint64) error {
	newStep := sdk.NewIntFromUint64(step)
	auctionParams := auctiontypes.AuctionParams{
		AppId:                  appId,
		AuctionDurationSeconds: auctionDurationSeconds,
		Buffer:                 buffer,
		Cusp:                   cusp,
		Step:                   newStep,
		PriceFunctionType:      priceFunctionType,
		SurplusId:              surplusId,
		DebtId:                 debtId,
		DutchId:                dutchId,
		BidDurationSeconds:     bidDurationSeconds,
	}

	k.SetAuctionParams(ctx, auctionParams)

	return nil
}

func (k Keeper) makeFalseForFlags(ctx sdk.Context, appId, assetId uint64) error {

	auctionLookupTable, found := k.GetAuctionMappingForApp(ctx, appId)
	if !found {
		return auctiontypes.ErrorInvalidAddress
	}
	for i, assetToAuction := range auctionLookupTable.AssetIdToAuctionLookup {
		if assetToAuction.AssetId == assetId {
			auctionLookupTable.AssetIdToAuctionLookup[i].IsAuctionActive = false
			err := k.SetAuctionMappingForApp(ctx, auctionLookupTable)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
