package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"math/rand"
	"reflect"
	"testing"
)
var (
	id = rand.Uint64()
	addr = sdk.AccAddress("comdex1yples84d8avjl")
	test []byte

)
func TestVaultKey(t *testing.T) {
	vtype := VaultKey(id)
	if reflect.TypeOf(vtype) != reflect.TypeOf(test){
		t.Errorf("test failed")
	}
}

func TestVaultForAddressByPair(t *testing.T) {
	atype := VaultForAddressByPair(addr,1)
	if reflect.TypeOf(atype) != reflect.TypeOf(test){
		t.Errorf("test failed")
	}
}