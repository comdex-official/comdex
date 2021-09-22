package types

import (
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func TestNewIBCParams(t *testing.T) {
	param := IBCParams{}
	if reflect.TypeOf(param) != reflect.TypeOf(IBCParams{}) {
		t.Error()
	}
}

func TestDefaultIBCParams(t *testing.T) {
	param := IBCParams{}
	if reflect.TypeOf(param) != reflect.TypeOf(DefaultIBCParams()) {
		t.Error()
	}
}

func TestValidate(t *testing.T) {
	invalidParams := []IBCParams{
		{"", "strParam"},
		{"cmdx", ""},
	}
	validParams := IBCParams{
		"cmdx", "strParam",
	}
	for _, params := range invalidParams {
		err := params.Validate()
		require.Error(t, err)
	}

	err := validParams.Validate()
	require.NoError(t, err)
}

func TestNewOracleParams(t *testing.T) {
	orclParam := OracleParams{}
	if reflect.TypeOf(orclParam) != reflect.TypeOf(NewOracleParams(1, 1, 1)) {
		t.Error()
	}
}

func TestDefaultOracleParams(t *testing.T) {
	orclParam:= OracleParams{}
	if reflect.TypeOf(orclParam) != reflect.TypeOf(DefaultOracleParams()) {
		t.Error()
	}
}

func TestValidateOraclePrams(t *testing.T) {
	invalidOracle := []OracleParams{
		{0, 1, 1},
		{1, 0, 1},
	}
	validOracle := OracleParams{
		1, 1, 1,
	}
	for _, params := range invalidOracle {
		err := params.Validate()
		require.Error(t, err)
	}

	err := validOracle.Validate()
	require.NoError(t, err)
}

func TestNewParams(t *testing.T) {
	params := Params{
		Admin:  "",
		IBC:    IBCParams{},
		Oracle: OracleParams{},
	}
	if reflect.TypeOf(params) != reflect.TypeOf(Params{}) {
		t.Error()
	}
}

func TestDefaultParams(t *testing.T) {
	defParam := NewParams(DefaultAdmin, IBCParams{}, OracleParams{})
	if reflect.TypeOf(defParam) != reflect.TypeOf(Params{
		Admin:  "",
		IBC:    IBCParams{},
		Oracle: OracleParams{},
	}) {
		t.Error()
	}
}
