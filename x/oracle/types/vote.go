package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"gopkg.in/yaml.v3"
)

func NewAggregateExchangeRatePrevote(
	hash AggregateVoteHash,
	voter sdk.ValAddress,
	submitBlock uint64,
) AggregateExchangeRatePrevote {
	return AggregateExchangeRatePrevote{
		Hash:        hash.String(),
		Voter:       voter.String(),
		SubmitBlock: submitBlock,
	}
}

// String implement stringify
func (v AggregateExchangeRatePrevote) String() string {
	out, _ := yaml.Marshal(v)
	return string(out)
}

func NewAggregateExchangeRateVote(
	decCoins sdk.DecCoins,
	voter sdk.ValAddress,
) AggregateExchangeRateVote {
	return AggregateExchangeRateVote{
		ExchangeRates: decCoins,
		Voter:         voter.String(),
	}
}

// String implement stringify
func (v AggregateExchangeRateVote) String() string {
	out, _ := yaml.Marshal(v)
	return string(out)
}

// ParseExchangeRateDecCoins DecCoins parser
func ParseExchangeRateDecCoins(tuplesStr string) (sdk.DecCoins, error) {
	if len(tuplesStr) == 0 {
		return nil, nil
	}

	decCoinsStrs := strings.Split(tuplesStr, ",")
	decCoins := make(sdk.DecCoins, len(decCoinsStrs))

	duplicateCheckMap := make(map[string]bool)
	for i, decCoinStr := range decCoinsStrs {
		denomAmountStr := strings.Split(decCoinStr, ":")
		if len(denomAmountStr) != 2 {
			return nil, fmt.Errorf("invalid exchange rate %s", decCoinStr)
		}

		dec, err := sdk.NewDecFromStr(denomAmountStr[1])
		if err != nil {
			return nil, err
		}
		if !dec.IsPositive() {
			return nil, ErrInvalidOraclePrice
		}

		denom := strings.ToUpper(denomAmountStr[0])

		decCoins[i] = sdk.NewDecCoinFromDec(denom, dec)

		if _, ok := duplicateCheckMap[denom]; ok {
			return nil, fmt.Errorf("duplicated denom %s", denom)
		}

		duplicateCheckMap[denom] = true
	}

	return decCoins, nil
}
