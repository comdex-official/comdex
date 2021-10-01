package types_test

import (
	"github.com/comdex-official/comdex/app"
	"github.com/comdex-official/comdex/x/asset/types"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)


func TestNewParams(t *testing.T) {
	params := types.Params{
		Admin: "",
	}
	if reflect.TypeOf(params) != reflect.TypeOf(types.Params{}) {
		t.Error()
	}
}

func TestDefaultParams(t *testing.T) {
	defParam := types.NewParams(types.DefaultAdmin)
	if reflect.TypeOf(defParam) != reflect.TypeOf(types.Params{
		Admin: "",
	}) {
		t.Error()
	}
}

func TestParams(t *testing.T) {
	app.SetAccountAddressPrefixes()
	tests := []struct {
		description string
		admin       string
		expectPass  bool
	}{
		{"validCondition","comdex1yples84d8avjlmegn90663mmjs4tardw45af6v",true},
		{"emptyAdmin", "", false},
		{"invalidAddress","sentinelles84d8avjlmegn90663mmjs4tardw45af6v",false},
	}

	for _, tc := range tests{
		msg := types.Params {
			tc.admin,
		}
		if tc.expectPass {
			require.NoError(t, msg.Validate() , "test: %v", tc.description)
		} else {
			require.Error(t, msg.Validate(), "test: %v", tc.description)
		}
	}
}
