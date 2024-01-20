package v5_0_0 //nolint:revive,stylecheck

import (
	sdkmath "cosmossdk.io/math"
	lendkeeper "github.com/comdex-official/comdex/x/lend/keeper"
	"github.com/comdex-official/comdex/x/lend/types"
	liquidationkeeper "github.com/comdex-official/comdex/x/liquidation/keeper"
	vaultkeeper "github.com/comdex-official/comdex/x/vault/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetVaultLengthCounter - Set vault length for liquidation check
func SetVaultLengthCounter(
	ctx sdk.Context,
	vaultkeeper vaultkeeper.Keeper,
) {
	var count uint64
	appExtendedPairVaultData, found := vaultkeeper.GetAppMappingData(ctx, 2)
	if found {
		for _, data := range appExtendedPairVaultData {
			count += uint64(len(data.VaultIds))
		}
	}
	vaultkeeper.SetLengthOfVault(ctx, count)
}

// FuncMigrateLiquidatedBorrow -  Migrate all liquidated borrow to new borrow struct and make is_liquidated field to true
func FuncMigrateLiquidatedBorrow(ctx sdk.Context, k lendkeeper.Keeper, liqK liquidationkeeper.Keeper) error {
	liqBorrow := liqK.GetLockedVaultByApp(ctx, 3)
	for _, v := range liqBorrow {
		if v.AmountIn.GT(sdkmath.ZeroInt()) && v.AmountOut.GT(sdkmath.ZeroInt()) {
			borrowMetaData := v.GetBorrowMetaData()
			pair, _ := k.GetLendPair(ctx, v.ExtendedPairId)
			assetIn, _ := k.Asset.GetAsset(ctx, pair.AssetIn)
			assetOut, _ := k.Asset.GetAsset(ctx, pair.AssetOut)
			amountIn := sdk.NewCoin(assetIn.Denom, v.AmountIn)
			amountOut := sdk.NewCoin(assetOut.Denom, v.AmountOut)
			var cpoolName string
			if pair.AssetOutPoolID == 1 {
				cpoolName = "CMDX-ATOM-CMST"
			} else {
				cpoolName = "OSMO-ATOM-CMST"
			}

			globalIndex, _ := sdkmath.LegacyNewDecFromStr("0.002")

			newBorrow := types.BorrowAsset{
				ID:                  v.OriginalVaultId,
				LendingID:           borrowMetaData.LendingId,
				IsStableBorrow:      borrowMetaData.IsStableBorrow,
				PairID:              v.ExtendedPairId,
				AmountIn:            amountIn,
				AmountOut:           amountOut,
				BridgedAssetAmount:  borrowMetaData.BridgedAssetAmount,
				BorrowingTime:       ctx.BlockTime(),
				StableBorrowRate:    borrowMetaData.StableBorrowRate,
				InterestAccumulated: sdkmath.LegacyNewDecFromInt(v.InterestAccumulated),
				GlobalIndex:         globalIndex,
				ReserveGlobalIndex:  sdkmath.LegacyOneDec(),
				LastInteractionTime: ctx.BlockTime(),
				CPoolName:           cpoolName,
				IsLiquidated:        false,
			}
			lend, _ := k.GetLend(ctx, newBorrow.LendingID)
			k.UpdateBorrowStats(ctx, pair, newBorrow.IsStableBorrow, v.AmountOut, true)

			poolAssetLBMappingData, _ := k.GetAssetStatsByPoolIDAndAssetID(ctx, pair.AssetOutPoolID, pair.AssetOut)
			poolAssetLBMappingData.BorrowIds = append(poolAssetLBMappingData.BorrowIds, newBorrow.ID)
			k.SetAssetStatsByPoolIDAndAssetID(ctx, poolAssetLBMappingData)

			k.SetUserBorrowIDCounter(ctx, newBorrow.ID)
			k.SetBorrow(ctx, newBorrow)

			mappingData, _ := k.GetUserLendBorrowMapping(ctx, lend.Owner, newBorrow.LendingID)
			mappingData.BorrowId = append(mappingData.BorrowId, newBorrow.ID)
			k.SetUserLendBorrowMapping(ctx, mappingData)
		} else {
			// delete faulty lockedVault
			liqK.DeleteLockedVault(ctx, 3, v.LockedVaultId)
		}
	}
	return nil
}
