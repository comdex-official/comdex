package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	appparams "github.com/comdex-official/comdex/app/params"
	"github.com/comdex-official/comdex/x/oracle/types"
)

// init, to set sdk bech32 address prefix for tests
func init() {
	appparams.SetAddressPrefixes()
}

func TestAggregateExchangeRatePrevoteString(t *testing.T) {
	addr := sdk.ValAddress(sdk.AccAddress([]byte("addr1_______________")))
	aggregateVoteHash := types.GetAggregateVoteHash("salt", "OJO:100,ATOM:100", addr)
	aggregateExchangeRatePreVote := types.NewAggregateExchangeRatePrevote(
		aggregateVoteHash,
		addr,
		100,
	)

	require.Equal(t, "hash: ccd44c2be8cec771f4bc8a0b33895bd44e3459b9\nvoter: ojovaloper1v9jxgu33ta047h6lta047h6lta047h6ludnc0y\nsubmit_block: 100\n", aggregateExchangeRatePreVote.String())
}

func TestAggregateExchangeRateVoteString(t *testing.T) {
	aggregateExchangeRatePreVote := types.NewAggregateExchangeRateVote(
		sdk.DecCoins{
			sdk.NewDecCoinFromDec(types.OjoDenom, sdk.OneDec()),
		},
		sdk.ValAddress(sdk.AccAddress([]byte("addr1_______________"))),
	)

	require.Equal(t, "exchangerates:\n    - denom: uojo\n      amount: \"1.000000000000000000\"\nvoter: ojovaloper1v9jxgu33ta047h6lta047h6lta047h6ludnc0y\n", aggregateExchangeRatePreVote.String())
}

func TestParseExchangeRateDecCoins(t *testing.T) {
	valid := "uojo:123.0,uatom:123.123"
	_, err := types.ParseExchangeRateDecCoins(valid)
	require.NoError(t, err)

	duplicatedDenom := "uojo:100.0,uatom:123.123,uatom:121233.123"
	_, err = types.ParseExchangeRateDecCoins(duplicatedDenom)
	require.Error(t, err)

	invalidCoins := "123.123"
	_, err = types.ParseExchangeRateDecCoins(invalidCoins)
	require.Error(t, err)

	invalidCoinsWithValid := "uojo:123.0,123.1"
	_, err = types.ParseExchangeRateDecCoins(invalidCoinsWithValid)
	require.Error(t, err)

	zeroCoinsWithValid := "uojo:0.0,uatom:123.1"
	_, err = types.ParseExchangeRateDecCoins(zeroCoinsWithValid)
	require.Error(t, err)

	negativeCoinsWithValid := "uojo:-1234.5,uatom:123.1"
	_, err = types.ParseExchangeRateDecCoins(negativeCoinsWithValid)
	require.Error(t, err)

	multiplePricesPerRate := "uojo:123: uojo:456,uusdc:789"
	_, err = types.ParseExchangeRateDecCoins(multiplePricesPerRate)
	require.Error(t, err)

	res, err := types.ParseExchangeRateDecCoins("")
	require.Nil(t, err)
	require.Nil(t, res)
}
