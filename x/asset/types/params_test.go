package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

var (
	addr = []sdk.AccAddress{
		sdk.AccAddress("test1"),
		sdk.AccAddress(""),
		sdk.AccAddress("zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz"),
	}
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

func TestValidateParams(t *testing.T) {
	var _ = sdk.AccAddress("")
	invalidParam := []Params{
		{"", IBCParams{
			Port:    "abc",
			Version: "def",
		}, OracleParams{1, 1, 1}},

		{"",IBCParams{},OracleParams{}},

		{"str", IBCParams{
			Port:    "",
			Version: "",
		}, OracleParams{1, 1, 1}},

		{"str", IBCParams{"abc", "dcf"}, OracleParams{},
		},
	}

	validParam := Params{
		"str", IBCParams{"abc", "def"}, OracleParams{1, 1, 1},
	}

	for _, param := range invalidParam {
		err := param.Validate()
		require.Error(t, err)
	}

	err := validParam.Validate()
	require.NoError(t, err)
}
