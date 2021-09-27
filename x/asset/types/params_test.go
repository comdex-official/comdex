package types

import (
	"reflect"
	"testing"
)

func TestNewParams(t *testing.T) {
	params := Params{
		Admin:  "",
	}
	if reflect.TypeOf(params) != reflect.TypeOf(Params{}) {
		t.Error()
	}
}

func TestDefaultParams(t *testing.T) {
	defParam := NewParams(DefaultAdmin)
	if reflect.TypeOf(defParam) != reflect.TypeOf(Params{
		Admin:  "",
	}) {
		t.Error()
	}
}
