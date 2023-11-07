package v13

import (
	assetkeeper "github.com/comdex-official/comdex/x/asset/keeper"
	auctionkeeperold "github.com/comdex-official/comdex/x/auction/keeper"
	auctionkeeper "github.com/comdex-official/comdex/x/auctionsV2/keeper"
	auctionsV2types "github.com/comdex-official/comdex/x/auctionsV2/types"
	liquidationkeeperold "github.com/comdex-official/comdex/x/liquidation/keeper"
	liquidationtypesold "github.com/comdex-official/comdex/x/liquidation/types"
	liquidationkeeper "github.com/comdex-official/comdex/x/liquidationsV2/keeper"
	liquidationtypes "github.com/comdex-official/comdex/x/liquidationsV2/types"
	marketkeeper "github.com/comdex-official/comdex/x/market/keeper"
	vaultkeeper "github.com/comdex-official/comdex/x/vault/keeper"
	vaultTypes "github.com/comdex-official/comdex/x/vault/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"time"
)

func ReturnCoin(ctx sdk.Context, assetKeeper assetkeeper.Keeper, assetID uint64, amount sdk.Int) sdk.Coin {
	asset, _ := assetKeeper.GetAsset(ctx, assetID)
	return sdk.NewCoin(asset.Denom, amount)
}

func MigrateAuctionsHarbor(
	ctx sdk.Context,
	assetKeeper assetkeeper.Keeper,
	auctionKeeperOld auctionkeeperold.Keeper,
	auctionKeeper auctionkeeper.Keeper,
	liquidationKeeperOld liquidationkeeperold.Keeper,
	liquidationKeeper liquidationkeeper.Keeper,
	marketKeeper marketkeeper.Keeper,
	vaultKeeper vaultkeeper.Keeper,
) {
	// first get all the auctions and their locked vaults and set them in new LiquidationsV2,
	// if locked vault is not available create a new locked vault
	// get all the auction and set them accordingly in new auctionV2

	auctionsOld := auctionKeeperOld.GetDutchAuctions(ctx, 2)
	if len(auctionsOld) != 0 {
		for _, auction := range auctionsOld {
			liquidationData, found := liquidationKeeperOld.GetLockedVault(ctx, 2, auction.LockedVaultId)
			if !found {
				//todo
				// create logic for generating locked vault from auction data
				var vault vaultTypes.Vault
				userVaults, _ := vaultKeeper.GetUserAppMappingData(ctx, auction.VaultOwner.String(), 2)
				// loop into vaults and check if asset in and asset out are matching
				for _, userVaultMap := range userVaults {
					extPair, _ := assetKeeper.GetPairsVault(ctx, userVaultMap.ExtendedPairId)
					pair, _ := assetKeeper.GetPair(ctx, extPair.PairId)
					if pair.AssetIn == auction.AssetInId && pair.AssetOut == auction.AssetOutId {
						vault, _ = vaultKeeper.GetVault(ctx, userVaultMap.VaultId)
					}
				}
				extPair, _ := assetKeeper.GetPairsVault(ctx, vault.ExtendedPairVaultID)
				pair, _ := assetKeeper.GetPair(ctx, extPair.PairId)
				assetIn, _ := assetKeeper.GetAsset(ctx, pair.AssetIn)
				totalIn, _ := marketKeeper.CalcAssetPrice(ctx, assetIn.Id, vault.AmountIn)

				totalFees := vault.InterestAccumulated.Add(vault.ClosingFeeAccumulated)
				liquidationData = liquidationtypesold.LockedVault{
					AppId:                        2,
					OriginalVaultId:              vault.Id,
					ExtendedPairId:               vault.ExtendedPairVaultID,
					Owner:                        vault.Owner,
					AmountIn:                     vault.AmountIn,
					AmountOut:                    vault.AmountOut,
					UpdatedAmountOut:             vault.AmountOut.Add(vault.InterestAccumulated),
					Initiator:                    "liquidationV1",
					IsAuctionComplete:            false,
					IsAuctionInProgress:          true,
					CrAtLiquidation:              sdk.ZeroDec(),
					CurrentCollaterlisationRatio: sdk.ZeroDec(),
					CollateralToBeAuctioned:      totalIn,
					LiquidationTimestamp:         ctx.BlockTime(),
					SellOffHistory:               nil,
					InterestAccumulated:          totalFees,
					Kind:                         nil,
				}
			}

			lockedVaultID := liquidationKeeper.GetLockedVaultID(ctx)
			extPair, _ := assetKeeper.GetPairsVault(ctx, liquidationData.ExtendedPairId)
			pair, _ := assetKeeper.GetPair(ctx, extPair.PairId)
			feesToBeCollected := sdk.NewDecFromInt(liquidationData.AmountOut).Mul(extPair.LiquidationPenalty).TruncateInt()
			newLockedVault := liquidationtypes.LockedVault{
				LockedVaultId:                lockedVaultID + 1,
				AppId:                        2,
				OriginalVaultId:              liquidationData.OriginalVaultId,
				ExtendedPairId:               liquidationData.ExtendedPairId,
				Owner:                        liquidationData.Owner,
				CollateralToken:              ReturnCoin(ctx, assetKeeper, pair.AssetIn, liquidationData.AmountIn),
				DebtToken:                    ReturnCoin(ctx, assetKeeper, pair.AssetOut, liquidationData.AmountOut),
				CurrentCollaterlisationRatio: liquidationData.CrAtLiquidation,
				CollateralToBeAuctioned:      ReturnCoin(ctx, assetKeeper, pair.AssetIn, liquidationData.AmountIn),
				TargetDebt:                   ReturnCoin(ctx, assetKeeper, pair.AssetOut, liquidationData.AmountOut.Add(feesToBeCollected)),
				LiquidationTimestamp:         liquidationData.LiquidationTimestamp,
				IsInternalKeeper:             false,
				InternalKeeperAddress:        "",
				ExternalKeeperAddress:        "",
				FeeToBeCollected:             feesToBeCollected,
				BonusToBeGiven:               sdk.ZeroInt(),
				InitiatorType:                "vault",
				AuctionType:                  true,
				IsDebtCmst:                   true,
				CollateralAssetId:            pair.AssetIn,
				DebtAssetId:                  pair.AssetOut,
			}
			// set locked vault id
			liquidationKeeper.SetLockedVaultID(ctx, newLockedVault.LockedVaultId)
			// set new locked vault
			liquidationKeeper.SetLockedVault(ctx, newLockedVault)

			// now migrate old auctions data to new module
			auctionID := auctionKeeper.GetAuctionID(ctx)

			twaDataCollateral, found := marketKeeper.GetTwa(ctx, newLockedVault.CollateralAssetId)
			if !found || !twaDataCollateral.IsPriceActive {
				return
			}
			twaDataDebt, found := marketKeeper.GetTwa(ctx, newLockedVault.DebtAssetId)
			if !found || !twaDataDebt.IsPriceActive {
				return
			}
			liquidationWhitelistingAppData, _ := liquidationKeeper.GetLiquidationWhiteListing(ctx, liquidationData.AppId)
			dutchAuctionParams := liquidationWhitelistingAppData.DutchAuctionParam
			auctionParams, _ := auctionKeeper.GetAuctionParams(ctx)

			CollateralTokenInitialPrice := auctionKeeper.GetCollalteralTokenInitialPrice(sdk.NewIntFromUint64(twaDataCollateral.Twa), dutchAuctionParams.Premium)

			newAuction := auctionsV2types.Auction{
				AuctionId:                   auctionID + 1,
				CollateralToken:             auction.OutflowTokenCurrentAmount, // outflow
				DebtToken:                   auction.InflowTokenCurrentAmount,  // inflow
				ActiveBiddingId:             0,
				BiddingIds:                  nil,
				CollateralTokenAuctionPrice: CollateralTokenInitialPrice,
				CollateralTokenOraclePrice:  sdk.NewDecFromInt(sdk.NewInt(int64(twaDataCollateral.Twa))),
				DebtTokenOraclePrice:        sdk.NewDecFromInt(sdk.NewInt(int64(twaDataDebt.Twa))),
				LockedVaultId:               newLockedVault.LockedVaultId,
				StartTime:                   ctx.BlockTime(),
				EndTime:                     ctx.BlockTime().Add(time.Second * time.Duration(auctionParams.AuctionDurationSeconds)),
				AppId:                       newLockedVault.AppId,
				AuctionType:                 newLockedVault.AuctionType,
				CollateralAssetId:           newLockedVault.CollateralAssetId,
				DebtAssetId:                 newLockedVault.DebtAssetId,
				BonusAmount:                 newLockedVault.BonusToBeGiven,
				CollateralTokenInitialPrice: CollateralTokenInitialPrice,
			}
			// update auction ID
			auctionKeeper.SetAuctionID(ctx, newAuction.AuctionId)

			// migrate old auctions to new module
			err := auctionKeeper.SetAuction(ctx, newAuction)
			if err != nil {
				return
			}

			// delete old auctions and locked vaults
			// first create a history of both of them then delete
			liquidationKeeperOld.SetLockedVaultHistory(ctx, liquidationData, liquidationData.LockedVaultId)
			err = auctionKeeperOld.SetHistoryDutchAuction(ctx, auction)
			if err != nil {
				return
			}
		}
	}
}
