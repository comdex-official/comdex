package types

import (
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func TestNewGenesisState(t *testing.T) {
	genesis := GenesisState{}.ValidateGenesis
	if genesis != nil {
		t.Error("required nil")
	}
}

func TestDefaultGenesisState(t *testing.T) {
	genesis := DefaultGenesisState()
	if reflect.TypeOf(genesis) != reflect.TypeOf(&GenesisState{}) {
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
