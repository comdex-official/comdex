package types

import (
	"reflect"
	"testing"
)

func TestNewGenesisState(t *testing.T) {
	a := GenesisState{}.ValidateGenesis
	if a != nil {
		t.Error("required nil")
	}
}

func TestDefaultGenesisState(t *testing.T) {
	a := DefaultGenesisState()
	if reflect.TypeOf(a) != reflect.TypeOf(&GenesisState{}) {
		t.Error()
	}
}

func TestValidateGenesis(t *testing.T) {
	a := GenesisState{}.ValidateGenesis
	if a != nil {
		t.Errorf("required nil")
	}
}


