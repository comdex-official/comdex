package types

import (
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
)

func TestAssetKey(t *testing.T) {
	key := AssetKey(uint64(rand.Int()))
	id := len(key)
	require.Equal(t, id,9)
}

func TestPairKey(t *testing.T) {
	key := PairKey(uint64(rand.Int()))
	id := len(key)
	require.Equal(t, id,9)
}
