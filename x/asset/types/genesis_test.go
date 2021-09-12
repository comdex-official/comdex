package types

import (
	"github.com/stretchr/testify/require"
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
	err := ValidateGenesis(&GenesisState{
		Assets:          nil,
		Markets:         nil,
		Pairs:           nil,
		Params:          Params{},
		ValidateGenesis: nil,
	})
	require.NoError(t, err)
}