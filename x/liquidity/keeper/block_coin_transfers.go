package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func (k Keeper) IsCoinTransferBlocked(ctx sdk.Context, coin sdk.Coin) bool {
	return !k.bankKeeper.IsSendEnabledCoin(ctx, coin)
}

func (k Keeper) BlockCoinTransferIfNotAlready(ctx sdk.Context, coin sdk.Coin) {
	if !k.IsCoinTransferBlocked(ctx, coin) {
		bankParams := k.bankKeeper.GetParams(ctx)
		bankParams.SendEnabled = append(bankParams.SendEnabled,
			banktypes.NewSendEnabled(coin.Denom, false),
		)
		k.bankKeeper.SetParams(ctx, bankParams)
	}
}
