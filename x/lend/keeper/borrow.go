package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)


func (k Keeper) GetBorrowerBorrows(ctx sdk.Context, borrowerAddr sdk.AccAddress) sdk.Coins {
	/*prefix := types.CreateAdjustedBorrowKeyNoDenom(borrowerAddr)
	totalBorrowed := sdk.NewCoins()

	iterator := func(key, val []byte) error {
		// get borrow denom from key
		denom := types.DenomFromKeyWithAddress(key, types.KeyPrefixAdjustedBorrow)

		// get borrowed amount
		var adjustedAmount sdk.Dec
		if err := adjustedAmount.Unmarshal(val); err != nil {
			// improperly marshaled borrow amount should never happen
			panic(err)
		}

		// apply interest scalar
		amount := adjustedAmount.Mul(k.getInterestScalar(ctx, denom)).Ceil().TruncateInt()

		// add to totalBorrowed
		totalBorrowed = totalBorrowed.Add(sdk.NewCoin(denom, amount))
		return nil
	}*/

	return sdk.Coins{}
}
