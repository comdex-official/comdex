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

// Liquidate Vaults function can liquidate all vaults created using the vault module.
//All vauts are looped and check if their underlying app has enabled liquidations.

func (k Keeper) LiquidateVaults(ctx sdk.Context) error {
	params := k.GetParams(ctx)

	//This allows us to loop over a slice of vaults per block , which doesnt stresses the abci.
	//Eg: if there exists 1,000,000 vaults  and the batch size is 100,000. then at every block 100,000 vaults will be looped and it will take
	//a total of 10 blocks to loop over all vaults.
	liquidationOffsetHolder, found := k.GetLiquidationOffsetHolder(ctx, types.VaultLiquidationsOffsetPrefix)
	if !found {
		liquidationOffsetHolder = types.NewLiquidationOffsetHolder(0)
	}
	// Fetching all  vaults
	totalVaults := k.vault.GetVaults(ctx)
	// Getting length of all vaults
	lengthOfVaults := int(k.vault.GetLengthOfVault(ctx))
	// Creating start and end slice
	start, end := types.GetSliceStartEndForLiquidations(lengthOfVaults, int(liquidationOffsetHolder.CurrentOffset), int(params.LiquidationBatchSize))
	if start == end {
		liquidationOffsetHolder.CurrentOffset = 0
		start, end = types.GetSliceStartEndForLiquidations(lengthOfVaults, int(liquidationOffsetHolder.CurrentOffset), int(params.LiquidationBatchSize))
	}
	newVaults := totalVaults[start:end]
	for _, vault := range newVaults {
		_ = utils.ApplyFuncIfNoError(ctx, func(ctx sdk.Context) error {

			//Checking ESM status and / or kill switch status
			esmStatus, found := k.esm.GetESMStatus(ctx, vault.AppId)
			klwsParams, _ := k.esm.GetKillSwitchData(ctx, vault.AppId)
			if (found && esmStatus.Status) || klwsParams.BreakerEnable {
				ctx.Logger().Error("Kill Switch Or ESM is enabled For Liquidation, liquidate_vaults.go for AppID %d", vault.AppId)
				continue
			}

			//Checking if app has enabled liquidations or not
			_, found = k.GetAppIDByAppForLiquidation(ctx, vault.AppId)

			if !found {
				return fmt.Errorf("Liquidation not enabled for App ID  %d", vault.AppId)
			}

			// Checking extended pair vault data for Minimum collateralisation ratio
			extPair, _ := k.asset.GetPairsVault(ctx, vault.ExtendedPairVaultID)
			liqRatio := extPair.MinCr
			totalOut := vault.AmountOut.Add(vault.InterestAccumulated).Add(vault.ClosingFeeAccumulated)
			collateralizationRatio, err := k.vault.CalculateCollateralizationRatio(ctx, vault.ExtendedPairVaultID, vault.AmountIn, totalOut)
			if err != nil {
				return fmt.Errorf("error Calculating CR in Liquidation, liquidate_vaults.go for vault ID %d", vault.Id)
			}
			if collateralizationRatio.LT(liqRatio) {
				totalDebt := vault.AmountOut.Add(vault.InterestAccumulated)
				err1 := k.rewards.CalculateVaultInterest(ctx, vault.AppId, vault.ExtendedPairVaultID, vault.Id, totalDebt, vault.BlockHeight, vault.BlockTime.Unix())
				if err1 != nil {
					return fmt.Errorf("error Calculating vault interest in Liquidation, liquidate_vaults.go for vaultID %d", vault.Id)
				}
				//Callling vault to use the updated values of the vault
				vault, _ := k.vault.GetVault(ctx, vault.Id)

				totalOut := vault.AmountOut.Add(vault.InterestAccumulated).Add(vault.ClosingFeeAccumulated)
				collateralizationRatio, err := k.vault.CalculateCollateralizationRatio(ctx, vault.ExtendedPairVaultID, vault.AmountIn, totalOut)
				if err != nil {
					return fmt.Errorf("error Calculating CR in Liquidation, liquidate_vaults.go for vaultID %d", vault.Id)
				}

				//Creating locked vault struct , which will trigger auction
				err = k.CreateLockedVault(ctx, vault.Id, vault.ExtendedPairVaultID, vault.Owner, vault.AmountIn, totalOut, collateralizationRatio, vault.AppId, false, false, "", "")
				if err != nil {
					return fmt.Errorf("error Creating Locked Vaults in Liquidation, liquidate_vaults.go for Vault %d", vault.Id)
				}
				length := k.vault.GetLengthOfVault(ctx)
				k.vault.SetLengthOfVault(ctx, length-1)

				//Removing data from existing structs
				k.vault.DeleteVault(ctx, vault.Id)
				var rewards rewardstypes.VaultInterestTracker
				rewards.AppMappingId = vault.AppId
				rewards.VaultId = vault.Id
				k.rewards.DeleteVaultInterestTracker(ctx, rewards)
				k.vault.DeleteAddressFromAppExtendedPairVaultMapping(ctx, vault.ExtendedPairVaultID, vault.Id, vault.AppId)
			}
			return nil
		})
	}

	liquidationOffsetHolder.CurrentOffset = uint64(end)
	k.SetLiquidationOffsetHolder(ctx, types.VaultLiquidationsOffsetPrefix, liquidationOffsetHolder)

	return nil

}

func (k Keeper) CreateLockedVault(ctx sdk.Context, OriginalVaultId, ExtendedPairId uint64, Owner string, AmountIn sdk.Int, AmountOut sdk.Int, collateralizationRatio sdk.Dec, appID uint64, isInternalKeeper bool, isExternalKeeper bool, internalKeeperAddress string, externalKeeperAddress string) error {
	lockedVaultID := k.GetLockedVaultID(ctx)

	value := types.LockedVault{
		LockedVaultId:                lockedVaultID + 1,
		AppId:                        appID,
		OriginalVaultId:              OriginalVaultId,
		ExtendedPairId:               ExtendedPairId,
		Owner:                        Owner,
		AmountIn:                     AmountIn,
		AmountOut:                    AmountOut,
		CurrentCollaterlisationRatio: collateralizationRatio,
		CollateralToBeAuctioned:      AmountIn,
		TargetDebt:                   AmountOut,
		LiquidationTimestamp:         ctx.BlockTime(),
		IsInternalKeeper:             false,
		InternalKeeperAddress:        "",
		IsExternalKeeper:             "",
		ExternalKeeperAddress:        "",
	
	}

	k.SetLockedVault(ctx, value)
	k.SetLockedVaultID(ctx, value.LockedVaultId)

	return nil
}
