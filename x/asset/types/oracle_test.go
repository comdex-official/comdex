package types

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidateMarket(t *testing.T) {
	invalidMarket := []Market{
		{"", 1},
		{"stringMarket", 1},
		{"strMkt", 0},
	}

	validMarket := Market{
		"strMkt", 1,
	}

	for _, market := range invalidMarket {
		err := market.Validate()
		require.Error(t, err)
	}

	err := validMarket.Validate()
	require.NoError(t, err)
}
