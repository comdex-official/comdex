package types

import (
	"reflect"
	"testing"
)

func TestDefaultGenesis(t *testing.T) {
	defaultGenesis := DefaultGenesis()
	if reflect.TypeOf(defaultGenesis) != reflect.TypeOf(&GenesisState{}) {
		t.Error()
	}
}

func TestValidate(t *testing.T) {
	genesisState := GenesisState{}.Validate()
	if genesisState != nil {
		t.Errorf("required nil")
	}
}
