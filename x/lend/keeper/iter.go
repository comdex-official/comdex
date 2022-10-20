package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/comdex-official/comdex/x/lend/types"
)

// To calculate pending rewards from last interaction
func (k Keeper) IterateLends(ctx sdk.Context, ID uint64) (sdk.Dec, error) {
	// to calculate lend rewards on the amount lent
	// check if the interest accumulated is sufficient for that assetID and poolID
	// send the cTokens to the lenders address if less than interest accumulated
	// if the user is claiming for the first time then a new lendRewardsTracker is created for that user

	lend, _ := k.GetLend(ctx, ID)
	lendAPR, _ := k.GetLendAPRByAssetIDAndPoolID(ctx, lend.PoolID, lend.AssetID)

	interestPerBlock, indexGlobalCurrent, _ := k.CalculateLendReward(ctx, lend.AmountIn.Amount.String(), lendAPR, lend)

	lendRewardsTracker, found := k.GetLendRewardTracker(ctx, lend.ID)
	if !found {
		lendRewardsTracker = types.LendRewardsTracker{
			LendingId:          lend.ID,
			RewardsAccumulated: sdk.ZeroDec(),
		}
	}

	// Adding interest to existing rewards accumulated
	lendRewardsTracker.RewardsAccumulated = lendRewardsTracker.RewardsAccumulated.Add(interestPerBlock)

	// initializing new variable newInterestPerBlock
	newInterestPerInteraction := sdk.ZeroInt()

	// checking if the rewards accumulated is greater than equal to 1
	if lendRewardsTracker.RewardsAccumulated.GTE(sdk.OneDec()) {
		newInterestPerInteraction = lendRewardsTracker.RewardsAccumulated.TruncateInt()
		newRewardDec := sdk.NewDec(newInterestPerInteraction.Int64())
		lendRewardsTracker.RewardsAccumulated = lendRewardsTracker.RewardsAccumulated.Sub(newRewardDec) // not losing decimal precision
	}
	k.SetLendRewardTracker(ctx, lendRewardsTracker) // setting the remaining decimal part

	// checking if sufficient cTokens are there to give out as rewards to user
	poolAssetLBMappingData, _ := k.GetAssetStatsByPoolIDAndAssetID(ctx, lend.PoolID, lend.AssetID)
	if newInterestPerInteraction.GT(poolAssetLBMappingData.TotalInterestAccumulated) {
		return sdk.Dec{}, types.ErrorInsufficientCTokensForRewards
	}
	if newInterestPerInteraction.GT(sdk.ZeroInt()) {
		// updating user's balance
		lend.AvailableToBorrow = lend.AvailableToBorrow.Add(newInterestPerInteraction)

		pool, _ := k.GetPool(ctx, lend.PoolID)
		asset, _ := k.GetAsset(ctx, lend.AssetID)
		Amount := sdk.NewCoin(asset.Denom, newInterestPerInteraction)
		assetRatesStat, _ := k.GetAssetRatesParams(ctx, lend.AssetID)

		cAsset, _ := k.GetAsset(ctx, assetRatesStat.CAssetID)
		cToken := sdk.NewCoin(cAsset.Denom, Amount.Amount)

		addr, _ := sdk.AccAddressFromBech32(lend.Owner)
		err := k.SendCoinFromModuleToAccount(ctx, pool.ModuleName, addr, cToken)
		if err != nil {
			return sdk.Dec{}, err
		}
		// subtracting newInterestPerInteraction from global lend and interest accumulated
		poolAssetLBMappingData.TotalInterestAccumulated = poolAssetLBMappingData.TotalInterestAccumulated.Sub(newInterestPerInteraction)
		//poolAssetLBMappingData.TotalLend = poolAssetLBMappingData.TotalLend.Sub(newInterestPerInteraction)
		k.SetAssetStatsByPoolIDAndAssetID(ctx, poolAssetLBMappingData)
		k.SetLend(ctx, lend)
	}

	return indexGlobalCurrent, nil
}

func (k Keeper) IterateBorrow(ctx sdk.Context, ID uint64) (sdk.Dec, sdk.Dec, error) {
	// to calculate borrow interest on existing borrow positions
	// also calculate the amount going to reserve pool for that borrow position

	borrow, _ := k.GetBorrow(ctx, ID)
	pair, _ := k.GetLendPair(ctx, borrow.PairID)
	reserveRates, err := k.GetReserveRate(ctx, pair.AssetOutPoolID, pair.AssetOut)
	if err != nil {
		return sdk.ZeroDec(), sdk.ZeroDec(), err
	}
	currBorrowAPR, _ := k.GetBorrowAPRByAssetID(ctx, pair.AssetOutPoolID, pair.AssetOut, borrow.IsStableBorrow)
	interestPerInteraction, indexGlobalCurrent, reservePoolAmountPerInteraction, reserveGlobalIndex, err := k.CalculateBorrowInterest(ctx, borrow.AmountOut.Amount.String(), currBorrowAPR, reserveRates, borrow)
	if err != nil {
		return sdk.ZeroDec(), sdk.ZeroDec(), err
	}

	if !borrow.IsStableBorrow {
		borrow.InterestAccumulated = borrow.InterestAccumulated.Add(interestPerInteraction)
	} else {
		stableInterestPerBlock, err := k.CalculateStableInterest(ctx, borrow.AmountOut.Amount.String(), borrow)
		if err != nil {
			return sdk.ZeroDec(), sdk.ZeroDec(), err
		}
		borrow.InterestAccumulated = borrow.InterestAccumulated.Add(stableInterestPerBlock)
	}

	// if reserve pool records are not found for the borrowId then a new reservePoolRecords is generated for that borrow ID (on first interaction)
	reservePoolRecords, found := k.GetBorrowInterestTracker(ctx, borrow.ID)
	if !found {
		reservePoolRecords = types.BorrowInterestTracker{
			BorrowingId:         borrow.ID,
			ReservePoolInterest: sdk.ZeroDec(),
		}
	}
	if reservePoolAmountPerInteraction.GT(sdk.ZeroDec()) {
		reservePoolRecords.ReservePoolInterest = reservePoolRecords.ReservePoolInterest.Add(reservePoolAmountPerInteraction)
	}
	k.SetBorrowInterestTracker(ctx, reservePoolRecords)
	k.SetBorrow(ctx, borrow)
	return indexGlobalCurrent, reserveGlobalIndex, nil
}

func (k Keeper) CalculateStableInterest(ctx sdk.Context, amount string, borrow types.BorrowAsset) (sdk.Dec, error) {
	currentTime := ctx.BlockTime().Unix()

	prevInterestTime := borrow.LastInteractionTime.Unix()
	if prevInterestTime == int64(types.Uint64Zero) {
		prevInterestTime = currentTime
	}
	secondsElapsed := currentTime - prevInterestTime
	if secondsElapsed < int64(types.Uint64Zero) {
		return sdk.ZeroDec(), sdkerrors.Wrap(types.ErrNegativeTimeElapsed, fmt.Sprintf("%d seconds", secondsElapsed))
	}
	yearsElapsed := sdk.NewDec(secondsElapsed).QuoInt64(types.SecondsPerYear)
	amt, _ := sdk.NewDecFromStr(amount)
	perc := borrow.StableBorrowRate
	newAmount := amt.Mul(perc).Mul(yearsElapsed)
	return newAmount, nil
}

func (k Keeper) CalculateLendReward(ctx sdk.Context, amount string, rate sdk.Dec, lend types.LendAsset) (sdk.Dec, sdk.Dec, error) {
	currentTime := ctx.BlockTime().Unix()
	lastInteraction := lend.LastInteractionTime
	globalIndex := lend.GlobalIndex
	prevInterestTime := lastInteraction.Unix()
	if prevInterestTime == int64(types.Uint64Zero) {
		prevInterestTime = currentTime
	}
	secondsElapsed := currentTime - prevInterestTime
	if secondsElapsed < int64(types.Uint64Zero) {
		return sdk.ZeroDec(), sdk.ZeroDec(), sdkerrors.Wrap(types.ErrNegativeTimeElapsed, fmt.Sprintf("%d seconds", secondsElapsed))
	}
	yearsElapsed := sdk.NewDec(secondsElapsed).QuoInt64(types.SecondsPerYear)
	amt, _ := sdk.NewDecFromStr(amount)

	effectiveRate := rate.Mul(yearsElapsed)
	factor1 := sdk.OneDec().Add(effectiveRate)
	indexGlobalCurrent := globalIndex.Mul(factor1)
	factor2 := indexGlobalCurrent.Quo(globalIndex)
	liabilityCurrent := amt.Mul(factor2)

	newAmount := liabilityCurrent.Sub(amt)
	return newAmount, indexGlobalCurrent, nil
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

	for _, v := range borrows {
		borrowPos, found := k.GetBorrow(ctx, v)
		if !found {
			continue
		}
		if !borrowPos.IsLiquidated {
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
	}
	return nil
}
