package types

import (
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type PriceStamps []PriceStamp

func NewPriceStamp(exchangeRate sdk.Dec, denom string, blockNum uint64) *PriceStamp {
	return &PriceStamp{
		ExchangeRate: &sdk.DecCoin{
			Amount: exchangeRate,
			Denom:  denom,
		},
		BlockNum: blockNum,
	}
}

// Decs returns the exchange rates in sdk.Dec format
func (p *PriceStamps) Decs() []sdk.Dec {
	decs := []sdk.Dec{}
	for _, priceStamp := range *p {
		decs = append(decs, priceStamp.ExchangeRate.Amount)
	}
	return decs
}

// FilterByBlock filters the PriceStamps by block number
func (p *PriceStamps) FilterByBlock(blockNum uint64) *PriceStamps {
	priceStamps := PriceStamps{}
	for _, priceStamp := range *p {
		if priceStamp.BlockNum == blockNum {
			priceStamps = append(priceStamps, priceStamp)
		}
	}
	return &priceStamps
}

// FilterByDenom filters the PriceStamps by denom
func (p *PriceStamps) FilterByDenom(denom string) *PriceStamps {
	priceStamps := PriceStamps{}
	for _, priceStamp := range *p {
		if priceStamp.ExchangeRate.Denom == denom {
			priceStamps = append(priceStamps, priceStamp)
		}
	}
	return &priceStamps
}

// Sort sorts the PriceStamps by block number and denom
func (p *PriceStamps) Sort() *PriceStamps {
	priceStamps := *p
	sort.Slice(
		priceStamps,
		func(i, j int) bool {
			if priceStamps[i].BlockNum == priceStamps[j].BlockNum {
				return priceStamps[i].ExchangeRate.Denom < priceStamps[j].ExchangeRate.Denom
			}
			return priceStamps[i].BlockNum < priceStamps[j].BlockNum
		},
	)
	return &priceStamps
}

// NewestPrices returns PriceStamps at the newest block number
func (p *PriceStamps) NewestPrices() *PriceStamps {
	blockNum := p.NewestBlockNum()
	return p.FilterByBlock(blockNum)
}

// NewestBlockNum returns the newest block number in the PriceStamps
func (p *PriceStamps) NewestBlockNum() uint64 {
	blockNum := uint64(0)
	for _, price := range *p {
		if price.BlockNum > blockNum {
			blockNum = price.BlockNum
		}
	}
	return blockNum
}
