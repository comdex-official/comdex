package keeper

import (
	"fmt"
	"github.com/comdex-official/comdex/x/lend/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"strconv"
)

func (k Keeper) IterateLends(ctx sdk.Context) error {
	lends, _ := k.GetLends(ctx)
	for _, v := range lends.LendIDs {
		lend, _ := k.GetLend(ctx, v)
		lendAPY, err := k.GetLendAPRByAssetIDAndPoolID(ctx, lend.PoolID, lend.AssetID)
		if err != nil {
			return err
		}
		interestPerBlock, err := k.CalculateRewards(ctx, lend.AmountIn.Amount.String(), lendAPY)
		if err != nil {
			return err
		}
		if interestPerBlock.GT(sdk.ZeroInt()) {

			updatedLend := types.LendAsset{
				ID:                 lend.ID,
				AssetID:            lend.AssetID,
				PoolID:             lend.PoolID,
				Owner:              lend.Owner,
				AmountIn:           lend.AmountIn,
				LendingTime:        lend.LendingTime,
				UpdatedAmountIn:    lend.UpdatedAmountIn.Add(interestPerBlock),
				AvailableToBorrow:  lend.AvailableToBorrow.Add(interestPerBlock),
				Reward_Accumulated: lend.Reward_Accumulated.Add(interestPerBlock),
				AppID:              lend.AppID,
			}

			pool, _ := k.GetPool(ctx, lend.PoolID)
			asset, _ := k.GetAsset(ctx, lend.AssetID)
			Amount := sdk.NewCoin(asset.Denom, interestPerBlock)
			assetRatesStat, found := k.GetAssetRatesStats(ctx, lend.AssetID)
			if !found {
				return sdkerrors.Wrap(types.ErrorAssetRatesStatsNotFound, strconv.FormatUint(lend.AssetID, 10))
			}
			cAsset, _ := k.GetAsset(ctx, assetRatesStat.CAssetID)
			cToken := sdk.NewCoin(cAsset.Denom, Amount.Amount)
			if err != nil {
				return err
			}
			addr, _ := sdk.AccAddressFromBech32(lend.Owner)
			err = k.SendCoinFromModuleToAccount(ctx, pool.ModuleName, addr, cToken)
			if err != nil {
				return err
			}
			k.SetLend(ctx, updatedLend)
		}
	}
	return nil
}

func (k Keeper) IterateBorrows(ctx sdk.Context) error {
	borrows, _ := k.GetBorrows(ctx)
	for _, v := range borrows.BorrowIDs {
		borrow, _ := k.GetBorrow(ctx, v)
		pair, _ := k.GetLendPair(ctx, borrow.PairID)

		borrowAPY, _ := k.GetBorrowAPRByAssetID(ctx, pair.AssetOutPoolID, pair.AssetOut, borrow.IsStableBorrow)
		interestPerBlock, err := k.CalculateRewards(ctx, borrow.AmountOut.Amount.String(), borrowAPY)
		if err != nil {
			return err
		}
		if interestPerBlock.GT(sdk.ZeroInt()) {

			updatedBorrow := types.BorrowAsset{
				ID:                   borrow.ID,
				LendingID:            borrow.LendingID,
				IsStableBorrow:       borrow.IsStableBorrow,
				PairID:               borrow.PairID,
				AmountIn:             borrow.AmountIn,
				AmountOut:            borrow.AmountOut,
				BridgedAssetAmount:   borrow.BridgedAssetAmount,
				BorrowingTime:        borrow.BorrowingTime,
				StableBorrowRate:     borrow.StableBorrowRate,
				UpdatedAmountOut:     borrow.UpdatedAmountOut.Add(interestPerBlock),
				Interest_Accumulated: borrow.Interest_Accumulated.Add(interestPerBlock),
			}
			k.SetBorrow(ctx, updatedBorrow)
		}
	}
	return nil
}

func (k Keeper) CalculateRewards(ctx sdk.Context, amount string, rate sdk.Dec) (sdk.Int, error) {

	currentTime := ctx.BlockTime().Unix()

	prevInterestTime := k.GetLastInterestTime(ctx)
	if prevInterestTime == 0 {
		prevInterestTime = currentTime
	}
	secondsElapsed := currentTime - prevInterestTime
	if secondsElapsed < 0 {
		return sdk.ZeroInt(), sdkerrors.Wrap(types.ErrNegativeTimeElapsed, fmt.Sprintf("%d seconds", secondsElapsed))
	}

	yearsElapsed := sdk.NewDec(secondsElapsed).QuoInt64(types.SecondsPerYear).MustFloat64()
	amtFloat, _ := strconv.ParseFloat(amount, 64)
	perc := rate.String()
	b, _ := sdk.NewDecFromStr(perc)

	newAmount := amtFloat * b.MustFloat64() * (yearsElapsed)
	return sdk.NewInt(int64(newAmount)), nil
}

func (k Keeper) RebalanceStableRates(ctx sdk.Context) error {
	borrows, _ := k.GetBorrows(ctx)
	for _, v := range borrows.BorrowIDs {

		borrowPos, _ := k.GetBorrow(ctx, v)
		if borrowPos.IsStableBorrow {
			pair, _ := k.GetLendPair(ctx, borrowPos.PairID)
			assetStats, _ := k.UpdateAPR(ctx, pair.AssetOutPoolID, pair.AssetOut)
			utilizationRatio, _ := k.GetUtilisationRatioByPoolIDAndAssetID(ctx, pair.AssetOutPoolID, pair.AssetOut)
			perc1, _ := sdk.NewDecFromStr("0.2")
			perc2, _ := sdk.NewDecFromStr("0.9")
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
