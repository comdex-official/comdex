package keeper

import (
	"fmt"
	"github.com/comdex-official/comdex/x/lend/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k Keeper) IterateLends(ctx sdk.Context) error {
	lends, _ := k.GetLends(ctx)
	for _, v := range lends.LendIDs {
		lend, _ := k.GetLend(ctx, v)
		lendAPR, err := k.GetLendAPRByAssetIDAndPoolID(ctx, lend.PoolID, lend.AssetID)
		if err != nil {
			continue
		}
		interestPerBlock, err := k.CalculateRewards(ctx, lend.AmountIn.Amount.String(), lendAPR)
		if err != nil {
			continue
		}
		lendRewardsTracker, found := k.GetLendRewardTracker(ctx, lend.ID)
		if !found {
			lendRewardsTracker = types.LendRewardsTracker{
				LendingId:          lend.ID,
				RewardsAccumulated: sdk.ZeroDec(),
			}
		}
		lendRewardsTracker.RewardsAccumulated = lendRewardsTracker.RewardsAccumulated.Add(interestPerBlock)
		newInterestPerBlock := sdk.ZeroInt()
		if lendRewardsTracker.RewardsAccumulated.GTE(sdk.OneDec()) {
			newInterestPerBlock = lendRewardsTracker.RewardsAccumulated.TruncateInt()
			newRewardDec := sdk.NewDec(newInterestPerBlock.Int64())
			lendRewardsTracker.RewardsAccumulated = lendRewardsTracker.RewardsAccumulated.Sub(newRewardDec)
		}
		k.SetLendRewardTracker(ctx, lendRewardsTracker)
		if newInterestPerBlock.GT(sdk.ZeroInt()) {
			lend.UpdatedAmountIn = lend.UpdatedAmountIn.Add(newInterestPerBlock)
			lend.AvailableToBorrow = lend.AvailableToBorrow.Add(newInterestPerBlock)
			lend.Reward_Accumulated = lend.Reward_Accumulated.Add(newInterestPerBlock)

			pool, _ := k.GetPool(ctx, lend.PoolID)
			asset, _ := k.GetAsset(ctx, lend.AssetID)
			Amount := sdk.NewCoin(asset.Denom, newInterestPerBlock)
			assetRatesStat, found := k.GetAssetRatesStats(ctx, lend.AssetID)
			if !found {
				continue
			}
			cAsset, _ := k.GetAsset(ctx, assetRatesStat.CAssetID)
			cToken := sdk.NewCoin(cAsset.Denom, Amount.Amount)
			if err != nil {
				continue
			}
			addr, _ := sdk.AccAddressFromBech32(lend.Owner)
			err = k.SendCoinFromModuleToAccount(ctx, pool.ModuleName, addr, cToken)
			if err != nil {
				continue
			}
			k.SetLend(ctx, lend)
		}
	}
	return nil
}

func (k Keeper) IterateBorrow(ctx sdk.Context, ID uint64) (sdk.Dec, sdk.Dec, error) {
	borrow, _ := k.GetBorrow(ctx, ID)
	pair, _ := k.GetLendPair(ctx, borrow.PairID)
	reserveRates, err := k.GetReserveRate(ctx, pair.AssetOutPoolID, pair.AssetOut)
	if err != nil {
		return sdk.ZeroDec(), sdk.ZeroDec(), err
	}
	currBorrowAPR, _ := k.GetBorrowAPRByAssetID(ctx, pair.AssetOutPoolID, pair.AssetOut, borrow.IsStableBorrow)
	interestPerBlock, indexGlobalCurrent, reservePoolAmountPerBlock, reserveGlobalIndex, err := k.CalculateBorrowInterest(ctx, borrow.AmountOut.Amount.String(), currBorrowAPR, reserveRates, borrow)
	if err != nil {
		return sdk.ZeroDec(), sdk.ZeroDec(), err
	}

	borrowInterestTracker, found := k.GetBorrowInterestTracker(ctx, borrow.ID)
	if !found {
		borrowInterestTracker = types.BorrowInterestTracker{
			BorrowingId:         borrow.ID,
			InterestAccumulated: sdk.ZeroDec(),
		}
	}
	borrowInterestTracker.InterestAccumulated = borrowInterestTracker.InterestAccumulated.Add(interestPerBlock)
	newInterestPerBlock := sdk.ZeroInt()
	if borrowInterestTracker.InterestAccumulated.GTE(sdk.OneDec()) {
		newInterestPerBlock = borrowInterestTracker.InterestAccumulated.TruncateInt()
		newRewardDec := sdk.NewDec(newInterestPerBlock.Int64())
		borrowInterestTracker.InterestAccumulated = borrowInterestTracker.InterestAccumulated.Sub(newRewardDec)
	}
	k.SetBorrowInterestTracker(ctx, borrowInterestTracker)

	reservePoolRecords, found := k.GetReservePoolRecordsForBorrow(ctx, borrow.ID)
	if !found {
		reservePoolRecords = types.ReservePoolRecordsForBorrow{
			ID:                  borrow.ID,
			InterestAccumulated: sdk.ZeroDec(),
		}
	}
	if reservePoolAmountPerBlock.GT(sdk.ZeroDec()) {
		reservePoolRecords.InterestAccumulated = reservePoolRecords.InterestAccumulated.Add(reservePoolAmountPerBlock)
	}
	k.SetReservePoolRecordsForBorrow(ctx, reservePoolRecords)
	if newInterestPerBlock.GT(sdk.ZeroInt()) {
		borrow.UpdatedAmountOut = borrow.UpdatedAmountOut.Add(newInterestPerBlock)
		borrow.Interest_Accumulated = borrow.Interest_Accumulated.Add(newInterestPerBlock)
		k.SetBorrow(ctx, borrow)
	}
	return indexGlobalCurrent, reserveGlobalIndex, nil
}

func (k Keeper) CalculateRewards(ctx sdk.Context, amount string, rate sdk.Dec) (sdk.Dec, error) {

	currentTime := ctx.BlockTime().Unix()

	prevInterestTime := k.GetLastInterestTime(ctx)
	if prevInterestTime == int64(types.Uint64Zero) {
		prevInterestTime = currentTime
	}
	secondsElapsed := currentTime - prevInterestTime
	if secondsElapsed < int64(types.Uint64Zero) {
		return sdk.ZeroDec(), sdkerrors.Wrap(types.ErrNegativeTimeElapsed, fmt.Sprintf("%d seconds", secondsElapsed))
	}
	yearsElapsed := sdk.NewDec(secondsElapsed).QuoInt64(types.SecondsPerYear)
	amt, _ := sdk.NewDecFromStr(amount)
	perc := rate
	newAmount := amt.Mul(perc).Mul(yearsElapsed)
	return newAmount, nil
}

func (k Keeper) CalculateBorrowInterest(ctx sdk.Context, amount string, rate, reserveRate sdk.Dec, borrow types.BorrowAsset) (sdk.Dec, sdk.Dec, sdk.Dec, sdk.Dec, error) {
	currentTime := ctx.BlockTime().Unix()
	lastInteraction := borrow.LastInteractionTime
	globalIndex := borrow.GlobalIndex
	reserveGlobalIndex := borrow.ReserveGlobalIndex
	prevInterestTime := lastInteraction.Unix()
	if prevInterestTime == int64(types.Uint64Zero) {
		prevInterestTime = currentTime
	}
	secondsElapsed := currentTime - prevInterestTime
	if secondsElapsed < int64(types.Uint64Zero) {
		return sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec(), sdkerrors.Wrap(types.ErrNegativeTimeElapsed, fmt.Sprintf("%d seconds", secondsElapsed))
	}
	yearsElapsed := sdk.NewDec(secondsElapsed).QuoInt64(types.SecondsPerYear)
	amt, _ := sdk.NewDecFromStr(amount)

	// for calculating interest accrued per interaction
	effectiveRate := rate.Mul(yearsElapsed)
	factor1 := sdk.OneDec().Add(effectiveRate)
	indexGlobalCurrent := globalIndex.Mul(factor1)
	factor2 := indexGlobalCurrent.Quo(globalIndex)
	liabilityCurrent := amt.Mul(factor2)

	newAmount := liabilityCurrent.Sub(amt)

	// for calculating amount to reserve pool accrued per interaction
	reserveEffectiveRate := reserveRate.Mul(yearsElapsed)
	reserveFactor1 := sdk.OneDec().Add(reserveEffectiveRate)
	reserveIndexGlobalCurrent := reserveGlobalIndex.Mul(reserveFactor1)
	reserveFactor2 := reserveIndexGlobalCurrent.Quo(reserveGlobalIndex)
	reserveLiabilityCurrent := amt.Mul(reserveFactor2)

	newAmountReservePool := reserveLiabilityCurrent.Sub(amt)

	return newAmount, indexGlobalCurrent, newAmountReservePool, reserveIndexGlobalCurrent, nil
}

func (k Keeper) ReBalanceStableRates(ctx sdk.Context) error {
	borrows, _ := k.GetBorrows(ctx)

	for _, v := range borrows.BorrowIDs {
		borrowPos, found := k.GetBorrow(ctx, v)
		if !found {
			continue
		}
		if borrowPos.IsStableBorrow {
			pair, found := k.GetLendPair(ctx, borrowPos.PairID)
			if !found {
				continue
			}
			assetStats, found := k.UpdateAPR(ctx, pair.AssetOutPoolID, pair.AssetOut)
			if !found {
				continue
			}
			utilizationRatio, err := k.GetUtilisationRatioByPoolIDAndAssetID(ctx, pair.AssetOutPoolID, pair.AssetOut)
			if err != nil {
				continue
			}
			perc1, _ := sdk.NewDecFromStr(types.Perc1)
			perc2, _ := sdk.NewDecFromStr(types.Perc2)
			if borrowPos.StableBorrowRate.GTE(assetStats.StableBorrowApr.Add(perc1)) {
				borrowPos.StableBorrowRate = assetStats.StableBorrowApr
				k.SetBorrow(ctx, borrowPos)
			} else if utilizationRatio.GT(perc2) && (borrowPos.StableBorrowRate.Add(perc1)).LTE(assetStats.StableBorrowApr) {
				borrowPos.StableBorrowRate = assetStats.StableBorrowApr
				k.SetBorrow(ctx, borrowPos)
			}
		}
	}
	return nil
}
