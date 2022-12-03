package types_test

import (
	"github.com/comdex-official/comdex/x/tokenmint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewMsgCreateRequest(t *testing.T) {
	testCases := []struct {
		name     string
		msg      *types.MsgMintNewTokensRequest
		isErrExp bool
	}{
		{
			name: "empty from",
			msg: types.NewMsgMintNewTokensRequest(
				sdk.AccAddress([]byte("")).String(), 1, 1,
			),
			isErrExp: true,
		},
		{
			name: "appID 0",
			msg: types.NewMsgMintNewTokensRequest(
				sdk.AccAddress([]byte("abc")).String(), 0, 1,
			),
			isErrExp: true,
		},
		{
			name: "assetID 0",
			msg: types.NewMsgMintNewTokensRequest(
				sdk.AccAddress([]byte("abc")).String(), 1, 0,
			),
			isErrExp: true,
		},
		{
			name: "valid case",
			msg: types.NewMsgMintNewTokensRequest(
				sdk.AccAddress([]byte("abc")).String(), 1, 1,
			),
			isErrExp: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.msg.Route(), types.RouterKey)
			require.Equal(t, tc.msg.Type(), types.TypeMsgMintNewTokensRequest)

			err := tc.msg.ValidateBasic()

			if tc.isErrExp {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
