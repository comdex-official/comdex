package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/comdex-official/comdex/x/oracle/types"
)

func (s *IntegrationTestSuite) TestIterateAllHistoricPrices() {
	keeper, ctx := s.app.OracleKeeper, s.ctx

	historicPrices := []types.PriceStamp{
		{BlockNum: 10, ExchangeRate: &sdk.DecCoin{
			Denom: "ucmdx", Amount: sdk.MustNewDecFromStr("20.45"),
		}},
		{BlockNum: 11, ExchangeRate: &sdk.DecCoin{
			Denom: "ucmdx", Amount: sdk.MustNewDecFromStr("20.44"),
		}},
		{BlockNum: 10, ExchangeRate: &sdk.DecCoin{
			Denom: "btc", Amount: sdk.MustNewDecFromStr("1200.56"),
		}},
		{BlockNum: 11, ExchangeRate: &sdk.DecCoin{
			Denom: "btc", Amount: sdk.MustNewDecFromStr("1200.19"),
		}},
	}

	for _, hp := range historicPrices {
		keeper.SetHistoricPrice(ctx, hp.ExchangeRate.Denom, hp.BlockNum, hp.ExchangeRate.Amount)
	}

	newPrices := []types.PriceStamp{}
	keeper.IterateAllHistoricPrices(
		ctx,
		func(historicPrice types.PriceStamp) bool {
			newPrices = append(newPrices, historicPrice)
			return false
		},
	)

	s.Require().Equal(len(historicPrices), len(newPrices))

	// Verify that the historic prices from IterateAllHistoricPrices equal
	// the ones set by SetHistoricPrice
FOUND:
	for _, oldPrice := range historicPrices {
		for _, newPrice := range newPrices {
			if oldPrice.BlockNum == newPrice.BlockNum && oldPrice.ExchangeRate.Denom == newPrice.ExchangeRate.Denom {
				s.Require().Equal(oldPrice.ExchangeRate.Amount, newPrice.ExchangeRate.Amount)
				continue FOUND
			}
		}
		s.T().Errorf("did not find match for historic price: %+v", oldPrice)
	}
}

func (s *IntegrationTestSuite) TestIterateAllMedianPrices() {
	keeper, ctx := s.app.OracleKeeper, s.ctx
	medians := sdk.DecCoins{
		{Denom: "ucmdx", Amount: sdk.MustNewDecFromStr("20.44")},
		{Denom: "atom", Amount: sdk.MustNewDecFromStr("2.66")},
		{Denom: "osmo", Amount: sdk.MustNewDecFromStr("13.64")},
	}

	for _, m := range medians {
		keeper.SetHistoricMedian(ctx, m.Denom, uint64(s.ctx.BlockHeight()), m.Amount)
	}

	newMedians := []types.PriceStamp{}
	keeper.IterateAllMedianPrices(
		ctx,
		func(median types.PriceStamp) bool {
			newMedians = append(newMedians, median)
			return false
		},
	)
	require.Equal(s.T(), len(medians), len(newMedians))

FOUND:
	for _, oldMedian := range medians {
		for _, newMedian := range newMedians {
			if oldMedian.Denom == newMedian.ExchangeRate.Denom {
				s.Require().Equal(oldMedian.Amount, newMedian.ExchangeRate.Amount)
				continue FOUND
			}
		}
		s.T().Errorf("did not find match for median price: %+v", oldMedian)
	}
}

func (s *IntegrationTestSuite) TestIterateAllMedianDeviationPrices() {
	keeper, ctx := s.app.OracleKeeper, s.ctx
	medians := sdk.DecCoins{
		{Denom: "ucmdx", Amount: sdk.MustNewDecFromStr("21.44")},
		{Denom: "atom", Amount: sdk.MustNewDecFromStr("3.66")},
		{Denom: "osmo", Amount: sdk.MustNewDecFromStr("14.64")},
	}

	for _, m := range medians {
		keeper.SetHistoricMedianDeviation(ctx, m.Denom, uint64(s.ctx.BlockHeight()), m.Amount)
	}

	newMedians := []types.PriceStamp{}
	keeper.IterateAllMedianDeviationPrices(
		ctx,
		func(median types.PriceStamp) bool {
			newMedians = append(newMedians, median)
			return false
		},
	)
	require.Equal(s.T(), len(medians), len(newMedians))

FOUND:
	for _, oldMedian := range medians {
		for _, newMedian := range newMedians {
			if oldMedian.Denom == newMedian.ExchangeRate.Denom {
				s.Require().Equal(oldMedian.Amount, newMedian.ExchangeRate.Amount)
				continue FOUND
			}
		}
		s.T().Errorf("did not find match for median price: %+v", oldMedian)
	}
}
