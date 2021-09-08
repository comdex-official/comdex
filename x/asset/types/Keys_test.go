package types

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestAssetKey(t *testing.T) {
	get := AssetKey(1)
	want := []byte{17,0, 0, 0, 0, 0, 0, 0, 1}
	if cmp.Equal(get, want) != true {
		t.Errorf("wanted %v got %v", get, want)
	}
}
func TestCalldataKey(t *testing.T) {
	get := CalldataKey(1)
	want := []byte{18, 0, 0, 0, 0, 0, 0, 0, 1}
	if cmp.Equal(get, want) != true {
		t.Errorf("wanted %v got %v", get, want)
	}
}
func TestMarketForAssetKey(t *testing.T) {
	get := MarketForAssetKey(1)
	want := []byte{34, 0, 0, 0, 0, 0, 0, 0, 1}
	if cmp.Equal(get, want) != true {
		t.Errorf("wanted %v got %v", get, want)
	}
}
func TestPairKey(t *testing.T) {
	get := PairKey(1)
	want := []byte{20, 0, 0, 0, 0, 0, 0, 0, 1}
	if cmp.Equal(get, want) != true {
		t.Errorf("wanted %v got %v", get, want)
	}
}
func TestPriceForMarketKey(t *testing.T) {
	get := PriceForMarketKey("a")
	want := []byte{35, 97}
	if cmp.Equal(get, want) != true {
		t.Errorf("wanted %v got %v", get, want)
	}
}
func TestAssetForDenomKey(t *testing.T) {
	get := AssetForDenomKey("a")
	want := []byte{33, 97}
	if cmp.Equal(get, want) != true {
		t.Errorf("wanted %v got %v", get, want)
	}
}
func TestMarketKey(t *testing.T) {
	get := MarketKey("a")
	want := []byte{19, 97}
	if cmp.Equal(get, want) != true {
		t.Errorf("wanted %v got %v", get, want)
	}
}
