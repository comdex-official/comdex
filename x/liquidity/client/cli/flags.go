package cli

// DONTCOVER

import (
	"strings"

	flag "github.com/spf13/pflag"
)

const (
	FlagPairID         = "pair-id"
	FlagDisabled       = "disabled"
	FlagPoolCoinDenom  = "pool-coin-denom"
	FlagReserveAddress = "reserve-address"
	FlagDenoms         = "denoms"
	FlagOrderLifespan  = "order-lifespan"
	FlagNumTicks       = "num-ticks"
)

func flagSetPools() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.String(FlagPairID, "", "The pair id")
	fs.String(FlagDisabled, "", "Whether the pool is disabled or not")

	return fs
}

func flagSetPool() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.String(FlagPoolCoinDenom, "", "The denomination of the pool coin")
	fs.String(FlagReserveAddress, "", "The bech-32 encoded address of the reserve account")

	return fs
}

func flagSetPairs() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.StringSlice(FlagDenoms, []string{}, "Coin denominations to query")

	return fs
}

func flagSetOrders() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.String(FlagPairID, "", "The pair id")

	return fs
}

func flagSetOrder() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.Duration(FlagOrderLifespan, 0, "Duration the order lives until it is expired; an order will be executed for at least one batch, even if the lifespan is 0; valid time units are ns|us|ms|s|m|h")

	return fs
}

func ParseStringSliceFromString(s string, separator string) ([]string, error) {
	stringSlice := strings.Split(s, separator)

	parsedStrings := make([]string, 0, len(stringSlice))
	for _, s := range stringSlice {
		s = strings.TrimSpace(s)

		parsedStrings = append(parsedStrings, s)
	}
	return parsedStrings, nil
}
