package types

import (
	"github.com/cosmos/cosmos-sdk/x/params/types"
	"reflect"
	"testing"
)

func TestParamKeyTable(t *testing.T) {
	paramKeyTable := reflect.TypeOf(ParamKeyTable())
	var newParamKeyTable types.KeyTable
	if paramKeyTable != reflect.TypeOf(newParamKeyTable){
		t.Errorf("expected %v got %v", paramKeyTable, newParamKeyTable)
	}
}
