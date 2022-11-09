package types

import (
	"fmt"
)

// DefaultGenesis returns the default Capability genesis state.
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:          DefaultParams(),
		AppGenesisState: []AppGenesisState{},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (genState GenesisState) Validate() error {
	if err := genState.Params.Validate(); err != nil {
		return fmt.Errorf("invalid params: %w", err)
	}

	for _, appState := range genState.AppGenesisState {
		if err := appState.GenericParams.Validate(); err != nil {
			return fmt.Errorf("invalid generic params: %w", err)
		}

		pairMap := map[uint64]Pair{}
		for i, pair := range appState.Pairs {
			if err := pair.Validate(); err != nil {
				return fmt.Errorf("invalid pair at index %d: %w", i, err)
			}
			if pair.Id > appState.LastPairId {
				return fmt.Errorf("pair at index %d has an id greater than last pair id: %d", i, pair.Id)
			}
			if _, ok := pairMap[pair.Id]; ok {
				return fmt.Errorf("pair at index %d has a duplicate id: %d", i, pair.Id)
			}
			pairMap[pair.Id] = pair
		}
		poolMap := map[uint64]Pool{}
		for i, pool := range appState.Pools {
			if err := pool.Validate(); err != nil {
				return fmt.Errorf("invalid pool at index %d: %w", i, err)
			}
			if pool.Id > appState.LastPoolId {
				return fmt.Errorf("pool at index %d has an id greater than last pool id: %d", i, pool.Id)
			}
			if _, ok := pairMap[pool.PairId]; !ok {
				return fmt.Errorf("pool at index %d has unknown pair id: %d", i, pool.PairId)
			}
			if _, ok := poolMap[pool.Id]; ok {
				return fmt.Errorf("pool at index %d has a duplicate pool id: %d", i, pool.Id)
			}
			poolMap[pool.Id] = pool
		}
		depositReqSet := map[uint64]map[uint64]struct{}{}
		for i, req := range appState.DepositRequests {
			if err := req.Validate(); err != nil {
				return fmt.Errorf("invalid deposit request at index %d: %w", i, err)
			}
			pool, ok := poolMap[req.PoolId]
			if !ok {
				return fmt.Errorf("deposit request at index %d has unknown pool id: %d", i, req.PoolId)
			}
			if req.MintedPoolCoin.Denom != pool.PoolCoinDenom {
				return fmt.Errorf("deposit request at index %d has wrong minted pool coin: %s", i, req.MintedPoolCoin)
			}
			pair := pairMap[pool.PairId]
			if req.DepositCoins.AmountOf(pair.BaseCoinDenom).IsZero() ||
				req.DepositCoins.AmountOf(pair.QuoteCoinDenom).IsZero() {
				return fmt.Errorf("deposit request at index %d has wrong deposit coins: %s", i, req.DepositCoins)
			}
			if set, ok := depositReqSet[req.PoolId]; ok {
				if _, ok := set[req.Id]; ok {
					return fmt.Errorf("deposit request at index %d has a duplicate id: %d", i, req.Id)
				}
			} else {
				depositReqSet[req.PoolId] = map[uint64]struct{}{}
			}
			depositReqSet[req.PoolId][req.Id] = struct{}{}
		}
		withdrawReqSet := map[uint64]map[uint64]struct{}{}
		for i, req := range appState.WithdrawRequests {
			if err := req.Validate(); err != nil {
				return fmt.Errorf("invalid withdraw request at index %d: %w", i, err)
			}
			pool, ok := poolMap[req.PoolId]
			if !ok {
				return fmt.Errorf("withdraw request at index %d has unknown pool id: %d", i, req.PoolId)
			}
			if req.PoolCoin.Denom != pool.PoolCoinDenom {
				return fmt.Errorf("withdraw request at index %d has wrong pool coin: %s", i, req.PoolCoin)
			}
			if set, ok := withdrawReqSet[req.PoolId]; ok {
				if _, ok := set[req.Id]; ok {
					return fmt.Errorf("withdraw request at index %d has a duplicate id: %d", i, req.Id)
				}
			} else {
				withdrawReqSet[req.PoolId] = map[uint64]struct{}{}
			}
			withdrawReqSet[req.PoolId][req.Id] = struct{}{}
		}
		orderSet := map[uint64]map[uint64]struct{}{}
		for i, order := range appState.Orders {
			if err := order.Validate(); err != nil {
				return fmt.Errorf("invalid order at index %d: %w", i, err)
			}
			pair, ok := pairMap[order.PairId]
			if !ok {
				return fmt.Errorf("order at index %d has unknown pair id: %d", i, order.PairId)
			}
			if order.BatchId > pair.CurrentBatchId {
				return fmt.Errorf("order at index %d has a batch id greater than its pair's current batch id: %d", i, order.BatchId)
			}
			var offerCoinDenom, demandCoinDenom string
			switch order.Direction {
			case OrderDirectionBuy:
				offerCoinDenom, demandCoinDenom = pair.QuoteCoinDenom, pair.BaseCoinDenom
			case OrderDirectionSell:
				offerCoinDenom, demandCoinDenom = pair.BaseCoinDenom, pair.QuoteCoinDenom
			}
			if order.OfferCoin.Denom != offerCoinDenom {
				return fmt.Errorf("order at index %d has wrong offer coin denom: %s != %s", i, order.OfferCoin.Denom, offerCoinDenom)
			}
			if order.ReceivedCoin.Denom != demandCoinDenom {
				return fmt.Errorf("order at index %d has wrong demand coin denom: %s != %s", i, order.OfferCoin.Denom, demandCoinDenom)
			}
			if set, ok := orderSet[order.PairId]; ok {
				if _, ok := set[order.Id]; ok {
					return fmt.Errorf("order at index %d has a duplicate id: %d", i, order.Id)
				}
			} else {
				orderSet[order.PairId] = map[uint64]struct{}{}
			}
			orderSet[order.PairId][order.Id] = struct{}{}
		}
		activeFarmerMap := map[string]ActiveFarmer{}
		for i := len(appState.ActiveFarmers) - 1; i >= 0; i-- {
			activeFarmer := appState.ActiveFarmers[i]
			if activeFarmer.FarmedPoolCoin.IsPositive() {
				if err := activeFarmer.Validate(); err != nil {
					return fmt.Errorf("invalid active farmer at index %d: %w", i, err)
				}
				if _, ok := poolMap[activeFarmer.PoolId]; !ok {
					return fmt.Errorf("active farmer at index %d has unknown pool id: %d", i, activeFarmer.PoolId)
				}
				if _, ok := activeFarmerMap[activeFarmer.Farmer]; ok {
					continue
					// return fmt.Errorf("active farmer at index %d has a duplicate farmer : %s", i, activeFarmer.Farmer)
				}
				activeFarmerMap[activeFarmer.Farmer] = activeFarmer
			}
		}

		for i, queuedFarmer := range appState.QueuedFarmers {
			if err := queuedFarmer.Validate(); err != nil {
				return fmt.Errorf("invalid queued farmer at index %d: %w", i, err)
			}
			if _, ok := poolMap[queuedFarmer.PoolId]; !ok {
				return fmt.Errorf("active farmer at index %d has unknown pool id: %d", i, queuedFarmer.PoolId)
			}
		}
	}

	return nil
}
