package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/x/liquidity/types"
)

// RegisterInvariants registers all liquidity module invariants.
func RegisterInvariants(ir sdk.InvariantRegistry, k Keeper) {
	ir.RegisterRoute(types.ModuleName, "deposit-coins-escrow", DepositCoinsEscrowInvariant(k))
	ir.RegisterRoute(types.ModuleName, "pool-coin-escrow", PoolCoinEscrowInvariant(k))
	ir.RegisterRoute(types.ModuleName, "remaining-offer-coin-escrow", RemainingOfferCoinEscrowInvariant(k))
	ir.RegisterRoute(types.ModuleName, "pool-status", PoolStatusInvariant(k))
}

// AllInvariants returns a combined invariant of the liquidity module.
func AllInvariants(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		for _, inv := range []func(Keeper) sdk.Invariant{
			DepositCoinsEscrowInvariant,
			PoolCoinEscrowInvariant,
			RemainingOfferCoinEscrowInvariant,
			PoolStatusInvariant,
		} {
			res, stop := inv(k)(ctx)
			if stop {
				return res, stop
			}
		}
		return "", false
	}
}

// DepositCoinsEscrowInvariant checks that the amount of coins in the global
// escrow address is greater or equal than remaining deposit coins in all
// deposit requests.
func DepositCoinsEscrowInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		allApps, found := k.assetKeeper.GetApps(ctx)
		if !found {
			return sdk.FormatInvariant(
				types.ModuleName, "deposit-coins-escrow",
				fmt.Sprintf("no apps found"),
			), false
		}
		for _, app := range allApps {
			escrowDepositCoins := sdk.Coins{}
			_ = k.IterateAllDepositRequests(ctx, app.Id, func(req types.DepositRequest) (stop bool, err error) {
				if req.Status == types.RequestStatusNotExecuted {
					escrowDepositCoins = escrowDepositCoins.Add(req.DepositCoins...)
				}
				return false, nil
			})
			balances := k.bankKeeper.SpendableCoins(ctx, types.GlobalEscrowAddress)
			broken := !balances.IsAllGTE(escrowDepositCoins)
			if broken {
				return sdk.FormatInvariant(
					types.ModuleName, "deposit-coins-escrow",
					fmt.Sprintf("escrow amount %s is smaller than expected %s", balances, escrowDepositCoins),
				), broken
			}
		}
		return sdk.FormatInvariant(
			types.ModuleName, "deposit-coins-escrow",
			fmt.Sprintf("all good"),
		), false
	}
}

// PoolCoinEscrowInvariant checks that the amount of coins in the global
// escrow address is greater or equal than remaining withdrawing pool
// coins in all withdrawal requests.
func PoolCoinEscrowInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		allApps, found := k.assetKeeper.GetApps(ctx)
		if !found {
			return sdk.FormatInvariant(
				types.ModuleName, "pool-coin-escrow",
				fmt.Sprintf("no apps found"),
			), false
		}
		for _, app := range allApps {
			escrowPoolCoins := sdk.Coins{}
			_ = k.IterateAllWithdrawRequests(ctx, app.Id, func(req types.WithdrawRequest) (stop bool, err error) {
				if req.Status == types.RequestStatusNotExecuted {
					escrowPoolCoins = escrowPoolCoins.Add(req.PoolCoin)
				}
				return false, nil
			})
			balances := k.bankKeeper.SpendableCoins(ctx, types.GlobalEscrowAddress)
			broken := !balances.IsAllGTE(escrowPoolCoins)
			if broken {
				return sdk.FormatInvariant(
					types.ModuleName, "pool-coin-escrow",
					fmt.Sprintf("escrow amount %s is smaller than expected %s", balances, escrowPoolCoins),
				), broken
			}
		}
		return sdk.FormatInvariant(
			types.ModuleName, "pool-coin-escrow",
			fmt.Sprintf("all good"),
		), false
	}
}

// RemainingOfferCoinEscrowInvariant checks that the amount of coins in each pair's
// escrow address is greater or equal than remaining offer coins in the pair's
// orders.
func RemainingOfferCoinEscrowInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		allApps, found := k.assetKeeper.GetApps(ctx)
		if !found {
			return sdk.FormatInvariant(
				types.ModuleName, "remaining-offer-coin-escrow",
				fmt.Sprintf("no apps found"),
			), false
		}

		for _, app := range allApps {
			var (
				count int
				msg   string
			)
			_ = k.IterateAllPairs(ctx, app.Id, func(pair types.Pair) (stop bool, err error) {
				remainingOfferCoins := sdk.Coins{}
				_ = k.IterateOrdersByPair(ctx, app.Id, pair.Id, func(req types.Order) (stop bool, err error) {
					if !req.Status.ShouldBeDeleted() {
						remainingOfferCoins = remainingOfferCoins.Add(req.RemainingOfferCoin)
					}
					return false, nil
				})
				balances := k.bankKeeper.SpendableCoins(ctx, pair.GetEscrowAddress())
				if !balances.IsAllGTE(remainingOfferCoins) {
					count++
					msg += fmt.Sprintf("\tpair %d has %s, which is smaller than %s\n", pair.Id, balances, remainingOfferCoins)
				}
				return false, nil
			})
			broken := count != 0
			if broken {
				return sdk.FormatInvariant(
					types.ModuleName, "remaining-offer-coin-escrow",
					fmt.Sprintf("%d pair(s) with insufficient escrow amount found\n%s", count, msg),
				), broken
			}
		}
		return sdk.FormatInvariant(
			types.ModuleName, "remaining-offer-coin-escrow",
			fmt.Sprintf("all good"),
		), false
	}
}

// PoolStatusInvariant checks that the pools with zero pool coin supply have
// been marked as disabled.
func PoolStatusInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {

		allApps, found := k.assetKeeper.GetApps(ctx)
		if !found {
			return sdk.FormatInvariant(
				types.ModuleName, "pool-status",
				fmt.Sprintf("no apps found"),
			), false
		}

		for _, app := range allApps {
			var (
				count int
				msg   string
			)
			_ = k.IterateAllPools(ctx, app.Id, func(pool types.Pool) (stop bool, err error) {
				if !pool.Disabled {
					ps := k.GetPoolCoinSupply(ctx, pool)
					if ps.IsZero() {
						count++
						msg += fmt.Sprintf("\tpool %d should be disabled, but not\n", pool.Id)
					}
				}
				return false, nil
			})
			broken := count != 0
			if broken {
				return sdk.FormatInvariant(
					types.ModuleName, "pool-status",
					fmt.Sprintf("%d pool(s) with wrong status found\n%s", count, msg),
				), broken
			}
		}
		return sdk.FormatInvariant(
			types.ModuleName, "pool-status",
			fmt.Sprintf("all good"),
		), false
	}
}
