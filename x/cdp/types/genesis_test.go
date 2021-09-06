package types

import (
	"reflect"
	"testing"
)

func TestDefaultGenesis(t *testing.T) {
	a := DefaultGenesis()
	if reflect.TypeOf(a) != reflect.TypeOf(&GenesisState{}) {
		t.Error()
	}
}

func TestValidate(t *testing.T) {
	a := GenesisState{}.Validate()
	if a != nil {
		t.Errorf("required nil")
	}
}