package keeper

import (
	"context"
	"fmt"

	"github.com/comdex-official/comdex/x/liquidation/types"
	rewardstypes "github.com/comdex-official/comdex/x/rewards/types"
	vaultypes "github.com/comdex-official/comdex/x/vault/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ types.MsgServer = msgServer{}

type msgServer struct {
	Keeper
}

func NewMsgServer(keeper Keeper) types.MsgServer {
	return &msgServer{
		Keeper: keeper,
	}
}

func (k msgServer) MsgLiquidateVault(c context.Context, msg *types.MsgLiquidateVaultRequest) (*types.MsgLiquidateVaultResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	appID, found := k.GetAppIDByAppForLiquidation(ctx, msg.AppId)
	if !found {
		return nil, types.ErrAppIDDoesNotExists
	}
	esmStatus, found := k.esm.GetESMStatus(ctx, appID)
	status := false
	if found {
		status = esmStatus.Status
	}
	klwsParams, _ := k.esm.GetKillSwitchData(ctx, appID)
	if klwsParams.BreakerEnable || status {
		return nil, fmt.Errorf("kill Switch Or ESM is enabled For Liquidation, AppID %d", appID)
	}
	// Call Vault Data
	vault, found := k.vault.GetVault(ctx, msg.VaultId)
	if !found {
		return nil, vaultypes.ErrorVaultDoesNotExist
	}
	if vault.AppId != appID {
		return nil, fmt.Errorf("vault and app id mismatch for vault ID %d", vault.Id)
	}
	extPair, _ := k.asset.GetPairsVault(ctx, vault.ExtendedPairVaultID)
	pair, _ := k.asset.GetPair(ctx, extPair.PairId)
	assetIn, _ := k.asset.GetAsset(ctx, pair.AssetIn)

	totalRate, err := k.market.CalcAssetPrice(ctx, assetIn.Id, vault.AmountIn)
	if err != nil {
		return nil, err
	}
	totalIn := totalRate

	liqRatio := extPair.MinCr
	totalOut := vault.AmountOut.Add(vault.InterestAccumulated).Add(vault.ClosingFeeAccumulated)
	collateralizationRatio, err := k.vault.CalculateCollateralizationRatio(ctx, vault.ExtendedPairVaultID, vault.AmountIn, totalOut)
	if err != nil {
		return nil, err
	}
	if collateralizationRatio.LT(liqRatio) {
		// calculate interest and update vault
		totalDebt := vault.AmountOut.Add(vault.InterestAccumulated)
		err1 := k.rewards.CalculateVaultInterest(ctx, vault.AppId, vault.ExtendedPairVaultID, vault.Id, totalDebt, vault.BlockHeight, vault.BlockTime.Unix())
		if err1 != nil {
			return nil, err
		}
		vault, _ := k.vault.GetVault(ctx, vault.Id)
		totalFees := vault.InterestAccumulated.Add(vault.ClosingFeeAccumulated)
		totalOut := vault.AmountOut.Add(vault.InterestAccumulated).Add(vault.ClosingFeeAccumulated)
		collateralizationRatio, err := k.vault.CalculateCollateralizationRatio(ctx, vault.ExtendedPairVaultID, vault.AmountIn, totalOut)
		if err != nil {
			return nil, err
		}
		err = k.CreateLockedVault(ctx, vault, totalIn, collateralizationRatio, msg.AppId, totalFees)
		if err != nil {
			return nil, err
		}
		k.vault.DeleteVault(ctx, vault.Id)
		var rewards rewardstypes.VaultInterestTracker
		rewards.AppMappingId = msg.AppId
		rewards.VaultId = vault.Id
		k.rewards.DeleteVaultInterestTracker(ctx, rewards)
		k.vault.DeleteAddressFromAppExtendedPairVaultMapping(ctx, vault.ExtendedPairVaultID, vault.Id, msg.AppId)
	}
	return &types.MsgLiquidateVaultResponse{}, nil
}

func (k msgServer) MsgLiquidateBorrow(c context.Context, msg *types.MsgLiquidateBorrowRequest) (*types.MsgLiquidateBorrowResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	borrowPos, found := k.lend.GetBorrow(ctx, msg.BorrowId)
	if !found {
		return nil, types.BorrowDoesNotExist
	}
	if borrowPos.IsLiquidated {
		return nil, types.BorrowPosAlreadyLiquidated
	}
	lendPair, _ := k.lend.GetLendPair(ctx, borrowPos.PairID)
	lendPos, found := k.lend.GetLend(ctx, borrowPos.LendingID)
	if !found {
		return nil, fmt.Errorf("lend Pos Not Found for ID %d", borrowPos.LendingID)
	}

	// calculating and updating the interest accumulated before checking for liquidations
	err := k.lend.MsgCalculateBorrowInterest(ctx, lendPos.Owner, borrowPos.ID)
	if err != nil {
		return nil, err
	}
	borrowPos, _ = k.lend.GetBorrow(ctx, msg.BorrowId)
	if !borrowPos.StableBorrowRate.Equal(sdk.ZeroDec()) {
		borrowPos, err = k.lend.ReBalanceStableRates(ctx, borrowPos)
		if err != nil {
			return nil, err
		}
	}

	killSwitchParams, _ := k.esm.GetKillSwitchData(ctx, lendPos.AppID)
	if killSwitchParams.BreakerEnable {
		return nil, fmt.Errorf("kill Switch is enabled in Liquidation for ID %d", lendPos.AppID)
	}
	// calculating and updating the interest accumulated before checking for liquidations
	err1 := k.lend.MsgCalculateBorrowInterest(ctx, lendPos.Owner, borrowPos.ID)
	if err1 != nil {
		return nil, err1
	}
	pool, _ := k.lend.GetPool(ctx, lendPos.PoolID)
	assetIn, _ := k.asset.GetAsset(ctx, lendPair.AssetIn)
	assetOut, _ := k.asset.GetAsset(ctx, lendPair.AssetOut)

	var currentCollateralizationRatio sdk.Dec
	var firstTransitAssetID, secondTransitAssetID uint64
	// for getting transit assets details
	for _, data := range pool.AssetData {
		if data.AssetTransitType == 2 {
			firstTransitAssetID = data.AssetID
		}
		if data.AssetTransitType == 3 {
			secondTransitAssetID = data.AssetID
		}
	}

	liqThreshold, _ := k.lend.GetAssetRatesParams(ctx, lendPair.AssetIn)
	liqThresholdBridgedAssetOne, _ := k.lend.GetAssetRatesParams(ctx, firstTransitAssetID)
	liqThresholdBridgedAssetTwo, _ := k.lend.GetAssetRatesParams(ctx, secondTransitAssetID)
	firstBridgedAsset, _ := k.asset.GetAsset(ctx, firstTransitAssetID)
	// there are three possible cases
	// 	a. if borrow is from same pool
	//  b. if borrow is from first transit asset
	//  c. if borrow is from second transit asset
	if borrowPos.BridgedAssetAmount.Amount.Equal(sdk.ZeroInt()) { // first condition
		currentCollateralizationRatio, err = k.lend.CalculateCollateralizationRatio(ctx, borrowPos.AmountIn.Amount, assetIn, borrowPos.AmountOut.Amount.Add(borrowPos.InterestAccumulated.TruncateInt()), assetOut)
		if err != nil {
			return nil, err
		}
		if sdk.Dec.GT(currentCollateralizationRatio, liqThreshold.LiquidationThreshold) {
			// after checking the currentCollateralizationRatio with LiquidationThreshold if borrow is to be liquidated then
			// CreateLockedBorrow function is called
			lockedVault, err := k.CreateLockedBorrow(ctx, borrowPos, currentCollateralizationRatio, lendPos.AppID)
			if err != nil {
				return nil, err
			}
			borrowPos.IsLiquidated = true // isLiquidated flag is set to true
			k.lend.SetBorrow(ctx, borrowPos)
			err = k.UpdateLockedBorrows(ctx, lockedVault)
			if err != nil {
				return nil, err
			}
		}
	} else {
		if borrowPos.BridgedAssetAmount.Denom == firstBridgedAsset.Denom {
			currentCollateralizationRatio, _ = k.lend.CalculateCollateralizationRatio(ctx, borrowPos.AmountIn.Amount, assetIn, borrowPos.AmountOut.Amount.Add(borrowPos.InterestAccumulated.TruncateInt()), assetOut)
			if sdk.Dec.GT(currentCollateralizationRatio, liqThreshold.LiquidationThreshold.Mul(liqThresholdBridgedAssetOne.LiquidationThreshold)) {
				lockedVault, err := k.CreateLockedBorrow(ctx, borrowPos, currentCollateralizationRatio, lendPos.AppID)
				if err != nil {
					return nil, err
				}
				borrowPos.IsLiquidated = true
				k.lend.SetBorrow(ctx, borrowPos)
				err = k.UpdateLockedBorrows(ctx, lockedVault)
				if err != nil {
					return nil, err
				}
			}
		} else {
			currentCollateralizationRatio, _ = k.lend.CalculateCollateralizationRatio(ctx, borrowPos.AmountIn.Amount, assetIn, borrowPos.AmountOut.Amount.Add(borrowPos.InterestAccumulated.TruncateInt()), assetOut)

			if sdk.Dec.GT(currentCollateralizationRatio, liqThreshold.LiquidationThreshold.Mul(liqThresholdBridgedAssetTwo.LiquidationThreshold)) {
				lockedVault, err := k.CreateLockedBorrow(ctx, borrowPos, currentCollateralizationRatio, lendPos.AppID)
				if err != nil {
					return nil, err
				}
				borrowPos.IsLiquidated = true
				k.lend.SetBorrow(ctx, borrowPos)
				err = k.UpdateLockedBorrows(ctx, lockedVault)
				if err != nil {
					return nil, err
				}
			}
		}
	}
	return &types.MsgLiquidateBorrowResponse{}, nil
}
