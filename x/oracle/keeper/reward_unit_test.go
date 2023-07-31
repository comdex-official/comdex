package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPrependOjoIfUnique(t *testing.T) {
	require := require.New(t)
	tcs := []struct {
		in  []string
		out []string
	}{
		// Should prepend "uojo" to a slice of denoms, unless it is already present.
		{[]string{}, []string{"uojo"}},
		{[]string{"a"}, []string{"uojo", "a"}},
		{[]string{"x", "a", "heeeyyy"}, []string{"uojo", "x", "a", "heeeyyy"}},
		{[]string{"x", "a", "uojo"}, []string{"x", "a", "uojo"}},
	}
	for i, tc := range tcs {
		require.Equal(tc.out, prependOjoIfUnique(tc.in), i)
	}
}
