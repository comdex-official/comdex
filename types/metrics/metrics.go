package metrics

import (
	"github.com/armon/go-metrics"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	missCounterLabel          = "miss_counter"
	exchangeRateLabel         = "exchange_rate"
	medianPriceLabel          = "median_price"
	medianDeviationPriceLabel = "median_deviation_price"
)

// RecordMissCounter records the miss counter gauge for a validator
func RecordMissCounter(operator sdk.ValAddress, missCounter uint64) {
	metrics.SetGaugeWithLabels(
		[]string{missCounterLabel},
		float32(missCounter),
		[]metrics.Label{{Name: "address", Value: operator.String()}},
	)
}

// RecordExchangeRate records the exchange rate gauge for a denom
func RecordExchangeRate(denom string, exchangeRate sdk.Dec) {
	metrics.SetGaugeWithLabels(
		[]string{exchangeRateLabel},
		float32(exchangeRate.MustFloat64()),
		[]metrics.Label{{Name: "denom", Value: denom}},
	)
}

// RecordAggregateExchangeRate records the median price gauge for a denom
func RecordMedianPrice(denom string, price sdk.Dec) {
	metrics.SetGaugeWithLabels(
		[]string{medianPriceLabel},
		float32(price.MustFloat64()),
		[]metrics.Label{{Name: "denom", Value: denom}},
	)
}

// RecordAggregateExchangeRate records the median deviation price gauge for a denom
func RecordMedianDeviationPrice(denom string, price sdk.Dec) {
	metrics.SetGaugeWithLabels(
		[]string{medianDeviationPriceLabel},
		float32(price.MustFloat64()),
		[]metrics.Label{{Name: "denom", Value: denom}},
	)
}
