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
	for _, v := range lends.LendIds {
		lend, _ := k.GetLend(ctx, v)
		lendAPY, _ := k.GetLendAPYByAssetId(ctx, lend.PoolId, lend.AssetId)
		interestPerBlock, err := k.CalculateRewards(ctx, lend.AmountIn.Amount, lendAPY)
		if err != nil {
			return err
		}

		updatedLend := types.LendAsset{
			ID:                 lend.ID,
			AssetId:            lend.AssetId,
			PoolId:             lend.PoolId,
			Owner:              lend.Owner,
			AmountIn:           lend.AmountIn,
			LendingTime:        lend.LendingTime,
			UpdatedAmountIn:    lend.UpdatedAmountIn.Add(interestPerBlock),
			Reward_Accumulated: lend.Reward_Accumulated.Add(interestPerBlock),
		}

		pool, _ := k.GetPool(ctx, lend.PoolId)
		asset, _ := k.GetAsset(ctx, lend.AssetId)
		Amount := sdk.NewCoin(asset.Denom, interestPerBlock)
		cToken, err := k.ExchangeToken(ctx, Amount, asset.Name)
		if err != nil {
			return err
		}
		addr, _ := sdk.AccAddressFromBech32(lend.Owner)
		cTokens := sdk.NewCoins(cToken)
		err = k.bank.SendCoinsFromModuleToAccount(ctx, pool.ModuleName, addr, cTokens)
		if err != nil {
			return err
		}
		k.SetLend(ctx, updatedLend)

	}
	return nil
}

func (k Keeper) IterateBorrows(ctx sdk.Context) error {
	borrows, _ := k.GetBorrows(ctx)
	for _, v := range borrows.BorrowIds {
		borrow, _ := k.GetBorrow(ctx, v)
		pair, _ := k.GetLendPair(ctx, borrow.PairID)

		borrowAPY, _ := k.GetBorrowAPYByAssetId(ctx, pair.AssetOutPoolId, pair.AssetOut, borrow.IsStableBorrow)
		interestPerBlock, err := k.CalculateRewards(ctx, borrow.AmountOut.Amount, borrowAPY)
		if err != nil {
			return err
		}

		updatedBorrow := types.BorrowAsset{
			ID:                   borrow.ID,
			LendingID:            borrow.LendingID,
			IsStableBorrow:       borrow.IsStableBorrow,
			PairID:               borrow.PairID,
			AmountIn:             borrow.AmountIn,
			AmountOut:            borrow.AmountOut,
			BorrowingTime:        borrow.BorrowingTime,
			StableBorrowRate:     borrow.StableBorrowRate,
			UpdatedAmountOut:     borrow.UpdatedAmountOut.Add(interestPerBlock),
			Interest_Accumulated: borrow.Interest_Accumulated.Add(interestPerBlock),
		}
		k.SetBorrow(ctx, updatedBorrow)

	}
	return nil
}

func (k Keeper) CalculateRewards(ctx sdk.Context, amount sdk.Int, rate sdk.Dec) (sdk.Int, error) {

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
	amtFloat, _ := strconv.ParseFloat(amount.String(), 64)
	perc := rate.String()
	b, _ := sdk.NewDecFromStr(perc)

	newAmount := amtFloat * b.MustFloat64() * (yearsElapsed)
	return sdk.NewInt(int64(newAmount)), nil
}
