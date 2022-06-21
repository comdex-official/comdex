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

func (k Keeper) IncreaseLockedVaultAmountIn(ctx sdk.Context, lockedVaultID uint64, amount sdk.Int) error {
	lockedVault, found := k.GetLockedVault(ctx, lockedVaultID)
	if !found {
		return auctiontypes.ErrorVaultNotFound
	}
	lockedVault.AmountIn = lockedVault.AmountIn.Add(amount)
	k.SetLockedVault(ctx, lockedVault)
	return nil
}

func (k Keeper) DecreaseLockedVaultAmountIn(ctx sdk.Context, lockedVaultID uint64, amount sdk.Int) (isZero bool, err error) {
	lockedVault, found := k.GetLockedVault(ctx, lockedVaultID)
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

func (k Keeper) DecreaseLockedVaultAmountOut(ctx sdk.Context, lockedVaultID uint64, amount sdk.Int) error {
	lockedVault, found := k.GetLockedVault(ctx, lockedVaultID)
	if !found {
		return auctiontypes.ErrorVaultNotFound
	}
	lockedVault.AmountIn = lockedVault.AmountOut.Sub(amount)
	k.SetLockedVault(ctx, lockedVault)
	return nil
}

func (k Keeper) AddAuctionParams(ctx sdk.Context, appID, auctionDurationSeconds uint64, buffer, cusp sdk.Dec, step, priceFunctionType, surplusID, debtID, dutchID uint64, bidDurationSeconds uint64) error {
	newStep := sdk.NewIntFromUint64(step)
	auctionParams := auctiontypes.AuctionParams{
		AppId:                  appID,
		AuctionDurationSeconds: auctionDurationSeconds,
		Buffer:                 buffer,
		Cusp:                   cusp,
		Step:                   newStep,
		PriceFunctionType:      priceFunctionType,
		SurplusId:              surplusID,
		DebtId:                 debtID,
		DutchId:                dutchID,
		BidDurationSeconds:     bidDurationSeconds,
	}

	k.SetAuctionParams(ctx, auctionParams)

	return nil
}

func (k Keeper) makeFalseForFlags(ctx sdk.Context, appID, assetID uint64) error {
	auctionLookupTable, found := k.GetAuctionMappingForApp(ctx, appID)
	if !found {
		return auctiontypes.ErrorInvalidAddress
	}
	for i, assetToAuction := range auctionLookupTable.AssetIdToAuctionLookup {
		if assetToAuction.AssetId == assetID {
			auctionLookupTable.AssetIdToAuctionLookup[i].IsAuctionActive = false
			err := k.SetAuctionMappingForApp(ctx, auctionLookupTable)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
