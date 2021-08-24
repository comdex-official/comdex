package types

import (
	"github.com/cosmos/cosmos-sdk/x/params/types"
	"reflect"
	"testing"
)

func TestParamKeyTable(t *testing.T) {
	a := reflect.TypeOf(ParamKeyTable())
	var b types.KeyTable
	if a != reflect.TypeOf(b){
		t.Errorf("expected %v got %v", a, b)
	}
}
