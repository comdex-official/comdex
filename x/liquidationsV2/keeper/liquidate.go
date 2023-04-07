package keeper

import (
	"fmt"
	utils "github.com/comdex-official/comdex/types"
	"github.com/comdex-official/comdex/x/liquidationsV2/types"
	rewardstypes "github.com/comdex-official/comdex/x/rewards/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) Liquidate(ctx sdk.Context) error {
	err := k.LiquidateVaults(ctx)
	if err != nil {
		return err
	}
	err = k.LiquidateBorrows(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) LiquidateVaults(ctx sdk.Context) error {
	appIds := k.GetAppIdsForLiquidation(ctx)
	params := k.GetParams(ctx)

	for i := range appIds {
		esmStatus, found := k.esm.GetESMStatus(ctx, appIds[i])
		status := false
		if found {
			status = esmStatus.Status
		}
		klwsParams, _ := k.esm.GetKillSwitchData(ctx, appIds[i])
		if klwsParams.BreakerEnable || status {
			ctx.Logger().Error("Kill Switch Or ESM is enabled For Liquidation, liquidate_vaults.go for AppID %d", appIds[i])
			continue
		}

		liquidationOffsetHolder, found := k.GetLiquidationOffsetHolder(ctx, appIds[i], types.VaultLiquidationsOffsetPrefix)
		if !found {
			liquidationOffsetHolder = types.NewLiquidationOffsetHolder(appIds[i], 0)
		}
		totalVaults := k.vault.GetVaults(ctx)
		lengthOfVaults := int(k.vault.GetLengthOfVault(ctx))
		//// get all vaults
		/// range over those vaults
		//// for length of vaults use vault counter
		//// wen inside the vault slice check if the app_id matches with that of app_id[i]

		start, end := types.GetSliceStartEndForLiquidations(lengthOfVaults, int(liquidationOffsetHolder.CurrentOffset), int(params.LiquidationBatchSize))
		if start == end {
			liquidationOffsetHolder.CurrentOffset = 0
			start, end = types.GetSliceStartEndForLiquidations(lengthOfVaults, int(liquidationOffsetHolder.CurrentOffset), int(params.LiquidationBatchSize))
		}

		newVaults := totalVaults[start:end]
		for _, vault := range newVaults {
			_ = utils.ApplyFuncIfNoError(ctx, func(ctx sdk.Context) error {
				if vault.AppId != appIds[i] {
					return fmt.Errorf("vault and app id mismatch in Liquidation, liquidate_vaults.go for vault ID %d", vault.Id)
				}
				extPair, _ := k.asset.GetPairsVault(ctx, vault.ExtendedPairVaultID)
				pair, _ := k.asset.GetPair(ctx, extPair.PairId)
				assetIn, found := k.asset.GetAsset(ctx, pair.AssetIn)
				if !found {
					return fmt.Errorf("asset not found in Liquidation, liquidate_vaults.go for vault ID %d", vault.Id)
				}
				totalRate, err := k.market.CalcAssetPrice(ctx, assetIn.Id, vault.AmountIn)
				if err != nil {
					return fmt.Errorf("error in CalcAssetPrice in Liquidation, liquidate_vaults.go for vault ID %d", vault.Id)
				}
				totalIn := totalRate

				liqRatio := extPair.MinCr
				totalOut := vault.AmountOut.Add(vault.InterestAccumulated).Add(vault.ClosingFeeAccumulated)
				collateralizationRatio, err := k.vault.CalculateCollateralizationRatio(ctx, vault.ExtendedPairVaultID, vault.AmountIn, totalOut)
				if err != nil {
					return fmt.Errorf("error Calculating CR in Liquidation, liquidate_vaults.go for vault ID %d", vault.Id)
				}
				if collateralizationRatio.LT(liqRatio) {
					// calculate interest and update vault
					totalDebt := vault.AmountOut.Add(vault.InterestAccumulated)
					err1 := k.rewards.CalculateVaultInterest(ctx, vault.AppId, vault.ExtendedPairVaultID, vault.Id, totalDebt, vault.BlockHeight, vault.BlockTime.Unix())
					if err1 != nil {
						return fmt.Errorf("error Calculating vault interest in Liquidation, liquidate_vaults.go for vaultID %d", vault.Id)
					}
					vault, _ := k.vault.GetVault(ctx, vault.Id)
					totalFees := vault.InterestAccumulated.Add(vault.ClosingFeeAccumulated)
					totalOut := vault.AmountOut.Add(vault.InterestAccumulated).Add(vault.ClosingFeeAccumulated)
					collateralizationRatio, err := k.vault.CalculateCollateralizationRatio(ctx, vault.ExtendedPairVaultID, vault.AmountIn, totalOut)
					if err != nil {
						return fmt.Errorf("error Calculating CR in Liquidation, liquidate_vaults.go for vaultID %d", vault.Id)
					}
					err = k.CreateLockedVault(ctx, vault, totalIn, collateralizationRatio, appIds[i], totalFees)
					if err != nil {
						return fmt.Errorf("error Creating Locked Vaults in Liquidation, liquidate_vaults.go for Vault %d", vault.Id)
					}
					k.vault.DeleteVault(ctx, vault.Id)
					var rewards rewardstypes.VaultInterestTracker
					rewards.AppMappingId = appIds[i]
					rewards.VaultId = vault.Id
					k.rewards.DeleteVaultInterestTracker(ctx, rewards)
					k.vault.DeleteAddressFromAppExtendedPairVaultMapping(ctx, vault.ExtendedPairVaultID, vault.Id, appIds[i])
				}
				return nil
			})
		}
		liquidationOffsetHolder.CurrentOffset = uint64(end)
		k.SetLiquidationOffsetHolder(ctx, types.VaultLiquidationsOffsetPrefix, liquidationOffsetHolder)
	}
	return nil
}

func (k Keeper) CreateLockedVault(ctx sdk.Context, OriginalVaultId, ExtendedPairId uint64, Owner string, AmountIn, AmountOut sdk.Int, totalIn sdk.Dec, collateralizationRatio sdk.Dec, appID uint64, totalFees sdk.Int) error {
	lockedVaultID := k.GetLockedVaultID(ctx)

	//value := types.LockedVault{
	//	LockedVaultId:           lockedVaultID + 1,
	//	AppId:                   appID,
	//	OriginalVaultId:         vault.Id,
	//	ExtendedPairId:          vault.ExtendedPairVaultID,
	//	Owner:                   vault.Owner,
	//	AmountIn:                vault.AmountIn,
	//	AmountOut:               vault.AmountOut,
	//	UpdatedAmountOut:        sdk.ZeroInt(),
	//	Initiator:               types.ModuleName,
	//	IsAuctionComplete:       false,
	//	IsAuctionInProgress:     false,
	//	CrAtLiquidation:         collateralizationRatio,
	//	CollateralToBeAuctioned: totalIn,
	//	LiquidationTimestamp:    ctx.BlockTime(),
	//	InterestAccumulated:     totalFees,
	//	Kind:                    nil,
	//}

	value := types.LockedVault{
		LockedVaultId:                lockedVaultID + 1,
		AppId:                        appID,
		OriginalVaultId:              OriginalVaultId,
		ExtendedPairId:               ExtendedPairId,
		Owner:                        Owner,
		AmountIn:                     AmountIn,
		AmountOut:                    AmountOut.Add(totalFees),
		CurrentCollaterlisationRatio: collateralizationRatio,
		CollateralToBeAuctioned:      totalIn,
		LiquidationTimestamp:         ctx.BlockTime(),
		IsInternalKeeper:             false,
		InternalKeeperAddress:        "",
		IsExternalKeeper:             "",
		ExternalKeeperAddress:        "",
	}

	k.SetLockedVault(ctx, value)
	k.SetLockedVaultID(ctx, value.LockedVaultId)
	length := k.vault.GetLengthOfVault(ctx)
	k.vault.SetLengthOfVault(ctx, length-1)
	return nil
}
