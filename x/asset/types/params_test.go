package types

import (
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)




func TestNewIBCParams(t *testing.T) {
	a := IBCParams{}
	if reflect.TypeOf(a) != reflect.TypeOf(IBCParams{}) {
		t.Error()
	}
}

func TestDefaultIBCParams(t *testing.T) {
	a := IBCParams{}
	if reflect.TypeOf(a) != reflect.TypeOf(DefaultIBCParams()) {
		t.Error()
	}
}

func TestValidate(t *testing.T) {
	invalidParams := []IBCParams{
		{"", "abc"},
		{"cmdx", ""},
	}
	validParams := IBCParams{
		"str", "sysStr",
	}
	for _, params := range invalidParams {
		err := params.Validate()
		require.Error(t, err)
	}

	err := validParams.Validate()
	require.NoError(t, err)
}

func TestNewOracleParams(t *testing.T) {
	a := OracleParams{}
	if reflect.TypeOf(a) != reflect.TypeOf(NewOracleParams(1, 1, 1)) {
		t.Error()
	}
}

func TestDefaultOracleParams(t *testing.T) {
	a := OracleParams{}
	if reflect.TypeOf(a) != reflect.TypeOf(DefaultOracleParams()) {
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
	a := Params{
		Admin:  "",
		IBC:    IBCParams{},
		Oracle: OracleParams{},
	}
	if reflect.TypeOf(a) != reflect.TypeOf(Params{}) {
		t.Error()
	}
}

func TestDefaultParams(t *testing.T) {
	a := NewParams(DefaultAdmin, IBCParams{}, OracleParams{})
	if reflect.TypeOf(a) != reflect.TypeOf(Params{
		Admin:  "",
		IBC:    IBCParams{},
		Oracle: OracleParams{},
	}) {
		t.Error()
	}
}

//func TestValidateParams(t *testing.T) {
//	invalidParam := []Params{
//		{},
//	}
//	validParam := Params{
//		{"str",sdk.AccAddress{},5},
//	}
//
//	for _, param := range invalidParam{
//		err := param.Validate()
//		require.Error(t, err)
//	}
//
//	err := validParam.Validate()
//	require.NoError(t, err)
//}
