package types

import (
	"reflect"
	"testing"
)


func TestNewIBCParams(t *testing.T) {
	a := IBCParams{}
	if reflect.TypeOf(a) != reflect.TypeOf(IBCParams{}){
		t.Error()
	}
}

func TestDefaultIBCParams(t *testing.T) {
	a := IBCParams{}
	if reflect.TypeOf(a) != reflect.TypeOf(DefaultIBCParams()){
		t.Error()
	}
}

func TestValidate(t *testing.T) {
	a := IBCParams{
		Port:    "",
		Version: "",
	}
	if reflect.TypeOf(a) != reflect.TypeOf(IBCParams{}) {
		t.Errorf("no match")
	}
}

func TestNewOracleParams(t *testing.T) {
	a := OracleParams{}
	if reflect.TypeOf(a) != reflect.TypeOf(NewOracleParams(1,1,1)){
		t.Error()
	}
}

func TestDefaultOracleParams(t *testing.T) {
	a := OracleParams{}
	if reflect.TypeOf(a) != reflect.TypeOf(DefaultOracleParams()){
		t.Error()
	}
}

func TestValidate2(t *testing.T) {
	a := OracleParams{AskCount: DefaultOracleAskCount,
		MinCount: DefaultOracleMinCount}
	if reflect.TypeOf(a) != reflect.TypeOf(OracleParams{AskCount: 1, MinCount: 1}) {
		t.Error()
	}
}

func TestNewParams(t *testing.T) {
	a := Params{
		Admin:  "",
		IBC:    IBCParams{},
		Oracle: OracleParams{},
	}
	if reflect.TypeOf(a) != reflect.TypeOf(Params{}){
			t.Error()
	}
}

func TestDefaultParams(t *testing.T) {
	a := NewParams(DefaultAdmin,IBCParams{},OracleParams{})
	if reflect.TypeOf(a) != reflect.TypeOf(Params{
		Admin:  "",
		IBC:    IBCParams{},
		Oracle: OracleParams{},
	}) {
		t.Error()
	}
}

func TestNewParams2(t *testing.T) {

}
