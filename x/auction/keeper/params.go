package keeper

import (
	"github.com/comdex-official/comdex/x/auction/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k *Keeper) LiquidationPenaltyPercent(ctx sdk.Context) (s string) {
	k.paramstore.Get(ctx, types.KeyLiquidationPenaltyPercent, &s)
	return
}

func (k *Keeper) AuctionDiscountPercent(ctx sdk.Context) (s string) {
	k.paramstore.Get(ctx, types.KeyAuctionDiscountPercent, &s)
	return
}

func (k *Keeper) AuctionDurationHours(ctx sdk.Context) (s uint64) {
	k.paramstore.Get(ctx, types.KeyAuctionDurationSeconds, &s)
	return
}

func (k *Keeper) DebtAuctionDecreasePercentage(ctx sdk.Context) (s sdk.Dec) {
	k.paramstore.Get(ctx, types.KeyDebtMintTokenDecreasePercentage, &s)
	return
}

func (k *Keeper) DutchBuffer(ctx sdk.Context) (s sdk.Dec) {
	k.paramstore.Get(ctx, types.KeyBuffer, &s)
	return
}
func (k *Keeper) DutchCusp(ctx sdk.Context) (s sdk.Dec) {
	k.paramstore.Get(ctx, types.KeyCusp, &s)
	return
}
func (k *Keeper) DutchTau(ctx sdk.Context) (s sdk.Int) {
	k.paramstore.Get(ctx, types.KeyTau, &s)
	return
}
func (k *Keeper) DutchDecreasePercentage(ctx sdk.Context) (s sdk.Dec) {
	k.paramstore.Get(ctx, types.KeyDutchDecreasePercentage, &s)
	return
}
func (k *Keeper) DutchChost(ctx sdk.Context) (s sdk.Dec) {
	k.paramstore.Get(ctx, types.KeyChost, &s)
	return
}
func (k *Keeper) DutchStep(ctx sdk.Context) (s sdk.Int) {
	k.paramstore.Get(ctx, types.KeyStep, &s)
	return
}
func (k *Keeper) DutchPriceFunctionType(ctx sdk.Context) (s uint64) {
	k.paramstore.Get(ctx, types.KeyPriceFunctionType, &s)
	return
}

func (k *Keeper) SurplusId(ctx sdk.Context) (s uint64) {
	k.paramstore.Get(ctx, types.KeySurplusId, &s)
	return
}

func (k *Keeper) DebtId(ctx sdk.Context) (s uint64) {
	k.paramstore.Get(ctx, types.KeyDebtId, &s)
	return
}
func (k *Keeper) DutchId(ctx sdk.Context) (s uint64) {
	k.paramstore.Get(ctx, types.KeyDutchId, &s)
	return
}

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.NewParams(
		k.LiquidationPenaltyPercent(ctx),
		k.AuctionDiscountPercent(ctx),
		k.AuctionDurationHours(ctx),
		k.DebtAuctionDecreasePercentage(ctx),
		k.DutchBuffer(ctx),
		k.DutchCusp(ctx),
		k.DutchTau(ctx),
		k.DutchDecreasePercentage(ctx),
		k.DutchChost(ctx),
		k.DutchStep(ctx),
		k.DutchPriceFunctionType(ctx),
		k.SurplusId(ctx),
		k.DebtId(ctx),
		k.DutchId(ctx),
	)
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}
