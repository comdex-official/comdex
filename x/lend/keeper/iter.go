package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"fmt"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/x/lend/types"
)

// IterateLends To calculate pending rewards from last interaction
func (k Keeper) IterateLends(ctx sdk.Context, ID uint64) (sdkmath.LegacyDec, error) {
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
			RewardsAccumulated: sdkmath.LegacyZeroDec(),
		}
	}

	// Adding interest to existing rewards accumulated
	lendRewardsTracker.RewardsAccumulated = lendRewardsTracker.RewardsAccumulated.Add(interestPerBlock)

	// initializing new variable newInterestPerBlock
	newInterestPerInteraction := sdkmath.ZeroInt()

	// checking if the rewards accumulated is greater than equal to 1
	if lendRewardsTracker.RewardsAccumulated.GTE(sdkmath.LegacyOneDec()) {
		newInterestPerInteraction = lendRewardsTracker.RewardsAccumulated.TruncateInt()
		newRewardDec := sdkmath.LegacyNewDecFromInt(newInterestPerInteraction)
		lendRewardsTracker.RewardsAccumulated = lendRewardsTracker.RewardsAccumulated.Sub(newRewardDec) // not losing decimal precision
	}
	k.SetLendRewardTracker(ctx, lendRewardsTracker) // setting the remaining decimal part
	pool, _ := k.GetPool(ctx, lend.PoolID)
	asset, _ := k.Asset.GetAsset(ctx, lend.AssetID)
	poolAssetLBMappingData, _ := k.GetAssetStatsByPoolIDAndAssetID(ctx, lend.PoolID, lend.AssetID)
	if newInterestPerInteraction.GT(sdkmath.ZeroInt()) {
		allReserveStats, found := k.GetAllReserveStatsByAssetID(ctx, lend.AssetID)
		if !found {
			allReserveStats = types.AllReserveStats{
				AssetID:                        lend.AssetID,
				AmountOutFromReserveToLenders:  sdkmath.ZeroInt(),
				AmountOutFromReserveForAuction: sdkmath.ZeroInt(),
				AmountInFromLiqPenalty:         sdkmath.ZeroInt(),
				AmountInFromRepayments:         sdkmath.ZeroInt(),
				TotalAmountOutToLenders:        sdkmath.ZeroInt(),
			}
		}
		if newInterestPerInteraction.GT(poolAssetLBMappingData.TotalInterestAccumulated) {
			modBal := k.ModuleBalance(ctx, types.ModuleName, asset.Denom)
			if modBal.LT(newInterestPerInteraction) {
				return sdkmath.LegacyDec{}, types.ErrorInsufficientCTokensForRewards
			}

			// return sdkmath.LegacyDec{}, types.ErrorInsufficientCTokensForRewards
			// check reserve moduleBalance
			// from reserve to pool and mint cToken
			// give cToken to user
			// update log of token used from reserve

			lend.AvailableToBorrow = lend.AvailableToBorrow.Add(newInterestPerInteraction)

			Amount := sdk.NewCoin(asset.Denom, newInterestPerInteraction)
			assetRatesStat, _ := k.GetAssetRatesParams(ctx, lend.AssetID)

			cAsset, _ := k.Asset.GetAsset(ctx, assetRatesStat.CAssetID)
			cToken := sdk.NewCoin(cAsset.Denom, Amount.Amount)
			addr, _ := sdk.AccAddressFromBech32(lend.Owner)
			// taking amount from reserve and minting cTokens
			err := k.bank.SendCoinsFromModuleToModule(ctx, types.ModuleName, pool.ModuleName, sdk.NewCoins(Amount))
			if err != nil {
				return sdkmath.LegacyDec{}, err
			}
			err = k.bank.MintCoins(ctx, pool.ModuleName, sdk.NewCoins(cToken))
			if err != nil {
				return sdkmath.LegacyDec{}, err
			}

			err = k.bank.SendCoinsFromModuleToAccount(ctx, pool.ModuleName, addr, sdk.NewCoins(cToken))
			if err != nil {
				return sdkmath.LegacyDec{}, err
			}
			lend.TotalRewards = lend.TotalRewards.Add(cToken.Amount)
			poolAssetLBMappingData.TotalLend = poolAssetLBMappingData.TotalLend.Add(cToken.Amount)
			k.SetAssetStatsByPoolIDAndAssetID(ctx, poolAssetLBMappingData)
			k.SetLend(ctx, lend)

			// updating reserve
			newAmount := newInterestPerInteraction.Quo(sdkmath.NewIntFromUint64(types.Uint64Two))
			reserve, found := k.GetReserveBuybackAssetData(ctx, lend.AssetID)
			if !found {
				reserve.AssetID = lend.AssetID
				reserve.BuybackAmount = sdkmath.ZeroInt()
				reserve.ReserveAmount = sdkmath.ZeroInt()
			}

			reserve.BuybackAmount = reserve.BuybackAmount.Sub(newAmount)
			reserve.ReserveAmount = reserve.ReserveAmount.Sub(newAmount)

			k.SetReserveBuybackAssetData(ctx, reserve)

			allReserveStats.AmountOutFromReserveToLenders = allReserveStats.AmountOutFromReserveToLenders.Add(newInterestPerInteraction)
			allReserveStats.TotalAmountOutToLenders = allReserveStats.TotalAmountOutToLenders.Add(newInterestPerInteraction)
			k.SetAllReserveStatsByAssetID(ctx, allReserveStats)
		} else {
			// updating user's balance
			lend.AvailableToBorrow = lend.AvailableToBorrow.Add(newInterestPerInteraction)
			Amount := sdk.NewCoin(asset.Denom, newInterestPerInteraction)
			assetRatesStat, _ := k.GetAssetRatesParams(ctx, lend.AssetID)

			cAsset, _ := k.Asset.GetAsset(ctx, assetRatesStat.CAssetID)
			cToken := sdk.NewCoin(cAsset.Denom, Amount.Amount)

			addr, _ := sdk.AccAddressFromBech32(lend.Owner)
			err := k.bank.SendCoinsFromModuleToAccount(ctx, pool.ModuleName, addr, sdk.NewCoins(cToken))
			if err != nil {
				return sdkmath.LegacyDec{}, err
			}
			lend.TotalRewards = lend.TotalRewards.Add(cToken.Amount)
			// subtracting newInterestPerInteraction from global lend and interest accumulated
			poolAssetLBMappingData.TotalInterestAccumulated = poolAssetLBMappingData.TotalInterestAccumulated.Sub(newInterestPerInteraction)
			poolAssetLBMappingData.TotalLend = poolAssetLBMappingData.TotalLend.Add(cToken.Amount)
			k.SetAssetStatsByPoolIDAndAssetID(ctx, poolAssetLBMappingData)
			k.SetLend(ctx, lend)
			allReserveStats.TotalAmountOutToLenders = allReserveStats.TotalAmountOutToLenders.Add(newInterestPerInteraction)
			k.SetAllReserveStatsByAssetID(ctx, allReserveStats)
		}
	}

	return indexGlobalCurrent, nil
}

func (k Keeper) IterateBorrow(ctx sdk.Context, ID uint64) (sdkmath.LegacyDec, sdkmath.LegacyDec, error) {
	// to calculate borrow interest on existing borrow positions
	// also calculate the amount going to reserve pool for that borrow position

	borrow, _ := k.GetBorrow(ctx, ID)
	pair, _ := k.GetLendPair(ctx, borrow.PairID)
	reserveRates, err := k.GetReserveRate(ctx, pair.AssetOutPoolID, pair.AssetOut)
	if err != nil {
		return sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), err
	}
	currBorrowAPR, _ := k.GetBorrowAPRByAssetID(ctx, pair.AssetOutPoolID, pair.AssetOut, borrow.IsStableBorrow)
	interestPerInteraction, indexGlobalCurrent, reservePoolAmountPerInteraction, reserveGlobalIndex, err := k.CalculateBorrowInterest(ctx, borrow.AmountOut.Amount.String(), currBorrowAPR, reserveRates, borrow)
	if err != nil {
		return sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), err
	}

	if !borrow.IsStableBorrow {
		borrow.InterestAccumulated = borrow.InterestAccumulated.Add(interestPerInteraction)
	} else {
		stableInterestPerBlock, err := k.CalculateStableInterest(ctx, borrow.AmountOut.Amount.String(), borrow)
		if err != nil {
			return sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), err
		}
		borrow.InterestAccumulated = borrow.InterestAccumulated.Add(stableInterestPerBlock)
	}

	// if reserve pool records are not found for the borrowId then a new reservePoolRecords is generated for that borrow ID (on first interaction)
	reservePoolRecords, found := k.GetBorrowInterestTracker(ctx, borrow.ID)
	if !found {
		reservePoolRecords = types.BorrowInterestTracker{
			BorrowingId:         borrow.ID,
			ReservePoolInterest: sdkmath.LegacyZeroDec(),
		}
	}
	if reservePoolAmountPerInteraction.GT(sdkmath.LegacyZeroDec()) {
		reservePoolRecords.ReservePoolInterest = reservePoolRecords.ReservePoolInterest.Add(reservePoolAmountPerInteraction)
	}
	k.SetBorrowInterestTracker(ctx, reservePoolRecords)
	k.SetBorrow(ctx, borrow)
	return indexGlobalCurrent, reserveGlobalIndex, nil
}

func (k Keeper) CalculateStableInterest(ctx sdk.Context, amount string, borrow types.BorrowAsset) (sdkmath.LegacyDec, error) {
	currentTime := ctx.BlockTime().Unix()

	prevInterestTime := borrow.LastInteractionTime.Unix()
	if prevInterestTime == int64(types.Uint64Zero) {
		prevInterestTime = currentTime
	}
	secondsElapsed := currentTime - prevInterestTime
	if secondsElapsed < int64(types.Uint64Zero) {
		return sdkmath.LegacyZeroDec(), errorsmod.Wrap(types.ErrNegativeTimeElapsed, fmt.Sprintf("%d seconds", secondsElapsed))
	}
	yearsElapsed := sdkmath.LegacyNewDec(secondsElapsed).QuoInt64(types.SecondsPerYear)
	amt, _ := sdkmath.LegacyNewDecFromStr(amount)
	perc := borrow.StableBorrowRate
	newAmount := amt.Mul(perc).Mul(yearsElapsed)
	return newAmount, nil
}

func (k Keeper) CalculateLendReward(ctx sdk.Context, amount string, rate sdkmath.LegacyDec, lend types.LendAsset) (sdkmath.LegacyDec, sdkmath.LegacyDec, error) {
	currentTime := ctx.BlockTime().Unix()
	lastInteraction := lend.LastInteractionTime
	globalIndex := lend.GlobalIndex
	prevInterestTime := lastInteraction.Unix()
	if prevInterestTime == int64(types.Uint64Zero) {
		prevInterestTime = currentTime
	}
	secondsElapsed := currentTime - prevInterestTime
	if secondsElapsed < int64(types.Uint64Zero) {
		return sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), errorsmod.Wrap(types.ErrNegativeTimeElapsed, fmt.Sprintf("%d seconds", secondsElapsed))
	}
	yearsElapsed := sdkmath.LegacyNewDec(secondsElapsed).QuoInt64(types.SecondsPerYear)
	amt, _ := sdkmath.LegacyNewDecFromStr(amount)

	effectiveRate := rate.Mul(yearsElapsed)
	factor1 := sdkmath.LegacyOneDec().Add(effectiveRate)
	indexGlobalCurrent := globalIndex.Mul(factor1)
	factor2 := indexGlobalCurrent.Quo(globalIndex)
	liabilityCurrent := amt.Mul(factor2)

	newAmount := liabilityCurrent.Sub(amt)
	return newAmount, indexGlobalCurrent, nil
}

func (k Keeper) CalculateBorrowInterest(ctx sdk.Context, amount string, rate, reserveRate sdkmath.LegacyDec, borrow types.BorrowAsset) (sdkmath.LegacyDec, sdkmath.LegacyDec, sdkmath.LegacyDec, sdkmath.LegacyDec, error) {
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
		return sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), errorsmod.Wrap(types.ErrNegativeTimeElapsed, fmt.Sprintf("%d seconds", secondsElapsed))
	}
	yearsElapsed := sdkmath.LegacyNewDec(secondsElapsed).QuoInt64(types.SecondsPerYear)
	amt, _ := sdkmath.LegacyNewDecFromStr(amount)

	// for calculating interest accrued per interaction
	effectiveRate := rate.Mul(yearsElapsed)
	factor1 := sdkmath.LegacyOneDec().Add(effectiveRate)
	indexGlobalCurrent := globalIndex.Mul(factor1)
	factor2 := indexGlobalCurrent.Quo(globalIndex)
	liabilityCurrent := amt.Mul(factor2)

	newAmount := liabilityCurrent.Sub(amt)

	// for calculating amount to reserve pool accrued per interaction
	reserveEffectiveRate := reserveRate.Mul(yearsElapsed)
	reserveFactor1 := sdkmath.LegacyOneDec().Add(reserveEffectiveRate)
	reserveIndexGlobalCurrent := reserveGlobalIndex.Mul(reserveFactor1)
	reserveFactor2 := reserveIndexGlobalCurrent.Quo(reserveGlobalIndex)
	reserveLiabilityCurrent := amt.Mul(reserveFactor2)

	newAmountReservePool := reserveLiabilityCurrent.Sub(amt)

	return newAmount, indexGlobalCurrent, newAmountReservePool, reserveIndexGlobalCurrent, nil
}

func (k Keeper) ReBalanceStableRates(ctx sdk.Context, borrowPos types.BorrowAsset) (types.BorrowAsset, error) {
	pair, found := k.GetLendPair(ctx, borrowPos.PairID)
	if !found {
		return borrowPos, types.ErrorPairNotFound
	}
	assetStats, found := k.UpdateAPR(ctx, pair.AssetOutPoolID, pair.AssetOut)
	if !found {
		return borrowPos, types.ErrorAssetRatesParamsNotFound
	}
	utilizationRatio, err := k.GetUtilisationRatioByPoolIDAndAssetID(ctx, pair.AssetOutPoolID, pair.AssetOut)
	if err != nil {
		return borrowPos, err
	}
	perc1, _ := sdkmath.LegacyNewDecFromStr(types.Perc1)                       // 20%
	perc2, _ := sdkmath.LegacyNewDecFromStr(types.Perc2)                       // 90%
	if borrowPos.StableBorrowRate.GTE(assetStats.StableBorrowApr.Add(perc1)) { // condition 1, 𝑆 ≥ 𝑆𝑡 + 20%, S is the rate at which you borrowed, and St is the current stable rate in the system.
		borrowPos.StableBorrowRate = assetStats.StableBorrowApr
	} else if (borrowPos.StableBorrowRate.Add(perc1)).LTE(assetStats.StableBorrowApr) || utilizationRatio.GTE(perc2) { // condition 2, 𝑆 + 20% ≤ 𝑆𝑡 ∨ 𝑢𝑡𝑖𝑙𝑖𝑧𝑎𝑡𝑖𝑜𝑛 ≥ 90%
		borrowPos.StableBorrowRate = assetStats.StableBorrowApr
	}

	return borrowPos, nil
}

func (k Keeper) IterateLendsForQuery(ctx sdk.Context) ([]types.PoolInterest, error) {
	var (
		poolInterest     types.PoolInterest
		poolInterestData types.PoolInterestData
	)

	pools := k.GetPools(ctx)
	var allPoolInterest []types.PoolInterest
	for _, pool := range pools {
		var v []types.PoolInterestData
		poolInterest.PoolID = pool.PoolID
		for _, data := range pool.AssetData {
			lbMap, _ := k.GetAssetStatsByPoolIDAndAssetID(ctx, pool.PoolID, data.AssetID)
			totalInt := sdkmath.LegacyZeroDec()
			for _, ID := range lbMap.LendIds {
				lend, _ := k.GetLend(ctx, ID)
				lendAPR, _ := k.GetLendAPRByAssetIDAndPoolID(ctx, lend.PoolID, lend.AssetID)
				interestPerBlock, _, _ := k.CalculateLendReward(ctx, lend.AmountIn.Amount.String(), lendAPR, lend)
				totalInt = totalInt.Add(interestPerBlock)
			}
			poolInterestData.LendInterest = totalInt.TruncateInt()
			poolInterestData.AssetID = data.AssetID
			v = append(v, poolInterestData)
		}
		poolInterest.PoolInterestData = v
		allPoolInterest = append(allPoolInterest, poolInterest)
	}

	return allPoolInterest, nil
}

func (k Keeper) IterateBorrowsForQuery(ctx sdk.Context) ([]types.PoolInterestB, error) {
	var (
		poolInterest     types.PoolInterestB
		poolInterestData types.PoolInterestDataB
	)

	pools := k.GetPools(ctx)
	var allPoolInterest []types.PoolInterestB
	for _, pool := range pools {
		var v []types.PoolInterestDataB
		poolInterest.PoolID = pool.PoolID
		for _, data := range pool.AssetData {
			lbMap, _ := k.GetAssetStatsByPoolIDAndAssetID(ctx, pool.PoolID, data.AssetID)
			totalInt := sdkmath.LegacyZeroDec()
			for _, ID := range lbMap.BorrowIds {
				borrow, _ := k.GetBorrow(ctx, ID)
				pair, _ := k.GetLendPair(ctx, borrow.PairID)
				reserveRates, err := k.GetReserveRate(ctx, pair.AssetOutPoolID, pair.AssetOut)
				if err != nil {
					return []types.PoolInterestB{}, err
				}
				currBorrowAPR, _ := k.GetBorrowAPRByAssetID(ctx, pair.AssetOutPoolID, pair.AssetOut, borrow.IsStableBorrow)
				interestPerInteraction, _, _, _, err := k.CalculateBorrowInterest(ctx, borrow.AmountOut.Amount.String(), currBorrowAPR, reserveRates, borrow)
				if err != nil {
					return []types.PoolInterestB{}, err
				}
				totalInt = totalInt.Add(interestPerInteraction)
			}
			poolInterestData.BorrowInterest = totalInt.TruncateInt()
			poolInterestData.AssetID = data.AssetID
			v = append(v, poolInterestData)
		}
		poolInterest.PoolInterestData = v
		allPoolInterest = append(allPoolInterest, poolInterest)
	}

	return allPoolInterest, nil
}
