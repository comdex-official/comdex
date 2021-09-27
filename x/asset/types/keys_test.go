package types

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
)

func TestAssetKey(t *testing.T) {
	keyMatch := AssetKey(uint64(rand.Int()))
	var match []uint8
	fmt.Printf("assetTest %T", keyMatch)
	if reflect.TypeOf(keyMatch) != reflect.TypeOf(match){
		t.Errorf("type misMatched")
	}
}

func TestAssetForDenomKey(t *testing.T) {
	keyMatch := AssetForDenomKey("AssetKey")
	var match []uint8
	fmt.Printf("denom key %T", keyMatch)
	if reflect.TypeOf(keyMatch) != reflect.TypeOf(match){
		t.Errorf("type misMatched")
	}
}

func TestPairKey(t *testing.T) {
	keyMatch := PairKey(uint64(rand.Int()))
	var match []uint8
	fmt.Printf("pair key %T", keyMatch)
	if reflect.TypeOf(keyMatch) != reflect.TypeOf(match){
		t.Errorf("type misMatched")
	}
}


