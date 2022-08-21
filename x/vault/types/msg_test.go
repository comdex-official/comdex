package types_test

import (
	"testing"

	"github.com/comdex-official/comdex/x/vault/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestNewMsgCreateRequest(t *testing.T) {
	testCases := []struct {
		name     string
		msg      *types.MsgCreateRequest
		isErrExp bool
	}{
		{
			name: "empty from",
			msg: types.NewMsgCreateRequest(
				sdk.AccAddress([]byte("")), 1, 1, sdk.NewInt(100), sdk.NewInt(50),
			),
			isErrExp: true,
		},
		{
			name: "amountIn nil",
			msg: types.NewMsgCreateRequest(
				sdk.MustAccAddressFromBech32("cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t"), 1, 1, sdk.Int{}, sdk.NewInt(50),
			),
			isErrExp: true,
		},
		{
			name: "amountIn negative",
			msg: types.NewMsgCreateRequest(
				sdk.MustAccAddressFromBech32("cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t"), 1, 1, sdk.NewInt(-123), sdk.NewInt(50),
			),
			isErrExp: true,
		},
		{
			name: "amountIn zero",
			msg: types.NewMsgCreateRequest(
				sdk.MustAccAddressFromBech32("cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t"), 1, 1, sdk.NewInt(0), sdk.NewInt(50),
			),
			isErrExp: true,
		},
		{
			name: "amountOut nil",
			msg: types.NewMsgCreateRequest(
				sdk.MustAccAddressFromBech32("cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t"), 1, 1, sdk.NewInt(50), sdk.Int{},
			),
			isErrExp: true,
		},
		{
			name: "amountOut negative",
			msg: types.NewMsgCreateRequest(
				sdk.MustAccAddressFromBech32("cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t"), 1, 1, sdk.NewInt(50), sdk.NewInt(-123),
			),
			isErrExp: true,
		},
		{
			name: "amountOut zero",
			msg: types.NewMsgCreateRequest(
				sdk.MustAccAddressFromBech32("cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t"), 1, 1, sdk.NewInt(50), sdk.NewInt(0),
			),
			isErrExp: true,
		},
		{
			name: "valid case",
			msg: types.NewMsgCreateRequest(
				sdk.MustAccAddressFromBech32("cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t"), 1, 1, sdk.NewInt(50), sdk.NewInt(25),
			),
			isErrExp: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {

			require.Equal(t, tc.msg.Route(), types.RouterKey)
			require.Equal(t, tc.msg.Type(), types.TypeMsgCreateRequest)

			err := tc.msg.ValidateBasic()

			if tc.isErrExp {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestNewMsgDepositRequest(t *testing.T) {
	testCases := []struct {
		name     string
		msg      *types.MsgDepositRequest
		isErrExp bool
	}{
		{
			name: "empty from",
			msg: types.NewMsgDepositRequest(
				sdk.AccAddress([]byte("")), 1, 1, 1, sdk.NewInt(50),
			),
			isErrExp: true,
		},
		{
			name: "vaultID zero",
			msg: types.NewMsgDepositRequest(
				sdk.MustAccAddressFromBech32("cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t"), 1, 1, 0, sdk.NewInt(50),
			),
			isErrExp: true,
		},
		{
			name: "amount nil",
			msg: types.NewMsgDepositRequest(
				sdk.MustAccAddressFromBech32("cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t"), 1, 1, 1, sdk.Int{},
			),
			isErrExp: true,
		},
		{
			name: "amount negative",
			msg: types.NewMsgDepositRequest(
				sdk.MustAccAddressFromBech32("cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t"), 1, 1, 1, sdk.NewInt(-50),
			),
			isErrExp: true,
		},
		{
			name: "amount Zero",
			msg: types.NewMsgDepositRequest(
				sdk.MustAccAddressFromBech32("cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t"), 1, 1, 1, sdk.NewInt(0),
			),
			isErrExp: true,
		},
		{
			name: "valid case",
			msg: types.NewMsgDepositRequest(
				sdk.MustAccAddressFromBech32("cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t"), 1, 1, 1, sdk.NewInt(25),
			),
			isErrExp: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {

			require.Equal(t, tc.msg.Route(), types.RouterKey)
			require.Equal(t, tc.msg.Type(), types.TypeMsgDepositRequest)

			err := tc.msg.ValidateBasic()

			if tc.isErrExp {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestNewMsgWithdrawRequest(t *testing.T) {
	testCases := []struct {
		name     string
		msg      *types.MsgWithdrawRequest
		isErrExp bool
	}{
		{
			name: "empty from",
			msg: types.NewMsgWithdrawRequest(
				sdk.AccAddress([]byte("")), 1, 1, 1, sdk.NewInt(50),
			),
			isErrExp: true,
		},
		{
			name: "vaultID zero",
			msg: types.NewMsgWithdrawRequest(
				sdk.MustAccAddressFromBech32("cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t"), 1, 1, 0, sdk.NewInt(50),
			),
			isErrExp: true,
		},
		{
			name: "amount nil",
			msg: types.NewMsgWithdrawRequest(
				sdk.MustAccAddressFromBech32("cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t"), 1, 1, 1, sdk.Int{},
			),
			isErrExp: true,
		},
		{
			name: "amount negative",
			msg: types.NewMsgWithdrawRequest(
				sdk.MustAccAddressFromBech32("cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t"), 1, 1, 1, sdk.NewInt(-50),
			),
			isErrExp: true,
		},
		{
			name: "amount Zero",
			msg: types.NewMsgWithdrawRequest(
				sdk.MustAccAddressFromBech32("cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t"), 1, 1, 1, sdk.NewInt(0),
			),
			isErrExp: true,
		},
		{
			name: "valid case",
			msg: types.NewMsgWithdrawRequest(
				sdk.MustAccAddressFromBech32("cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t"), 1, 1, 1, sdk.NewInt(25),
			),
			isErrExp: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {

			require.Equal(t, tc.msg.Route(), types.RouterKey)
			require.Equal(t, tc.msg.Type(), types.TypeMsgWithdrawRequest)

			err := tc.msg.ValidateBasic()

			if tc.isErrExp {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestNewMsgDrawRequest(t *testing.T) {
	testCases := []struct {
		name     string
		msg      *types.MsgDrawRequest
		isErrExp bool
	}{
		{
			name: "empty from",
			msg: types.NewMsgDrawRequest(
				sdk.AccAddress([]byte("")), 1, 1, 1, sdk.NewInt(50),
			),
			isErrExp: true,
		},
		{
			name: "vaultID zero",
			msg: types.NewMsgDrawRequest(
				sdk.MustAccAddressFromBech32("cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t"), 1, 1, 0, sdk.NewInt(50),
			),
			isErrExp: true,
		},
		{
			name: "amount nil",
			msg: types.NewMsgDrawRequest(
				sdk.MustAccAddressFromBech32("cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t"), 1, 1, 1, sdk.Int{},
			),
			isErrExp: true,
		},
		{
			name: "amount negative",
			msg: types.NewMsgDrawRequest(
				sdk.MustAccAddressFromBech32("cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t"), 1, 1, 1, sdk.NewInt(-50),
			),
			isErrExp: true,
		},
		{
			name: "amount Zero",
			msg: types.NewMsgDrawRequest(
				sdk.MustAccAddressFromBech32("cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t"), 1, 1, 1, sdk.NewInt(0),
			),
			isErrExp: true,
		},
		{
			name: "valid case",
			msg: types.NewMsgDrawRequest(
				sdk.MustAccAddressFromBech32("cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t"), 1, 1, 1, sdk.NewInt(25),
			),
			isErrExp: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {

			require.Equal(t, tc.msg.Route(), types.RouterKey)
			require.Equal(t, tc.msg.Type(), types.TypeMsgDrawRequest)

			err := tc.msg.ValidateBasic()

			if tc.isErrExp {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestNewMsgRepayRequest(t *testing.T) {
	testCases := []struct {
		name     string
		msg      *types.MsgRepayRequest
		isErrExp bool
	}{
		{
			name: "empty from",
			msg: types.NewMsgRepayRequest(
				sdk.AccAddress([]byte("")), 1, 1, 1, sdk.NewInt(50),
			),
			isErrExp: true,
		},
		{
			name: "vaultID zero",
			msg: types.NewMsgRepayRequest(
				sdk.MustAccAddressFromBech32("cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t"), 1, 1, 0, sdk.NewInt(50),
			),
			isErrExp: true,
		},
		{
			name: "amount nil",
			msg: types.NewMsgRepayRequest(
				sdk.MustAccAddressFromBech32("cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t"), 1, 1, 1, sdk.Int{},
			),
			isErrExp: true,
		},
		{
			name: "amount negative",
			msg: types.NewMsgRepayRequest(
				sdk.MustAccAddressFromBech32("cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t"), 1, 1, 1, sdk.NewInt(-50),
			),
			isErrExp: true,
		},
		{
			name: "amount Zero",
			msg: types.NewMsgRepayRequest(
				sdk.MustAccAddressFromBech32("cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t"), 1, 1, 1, sdk.NewInt(0),
			),
			isErrExp: true,
		},
		{
			name: "valid case",
			msg: types.NewMsgRepayRequest(
				sdk.MustAccAddressFromBech32("cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t"), 1, 1, 1, sdk.NewInt(25),
			),
			isErrExp: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {

			require.Equal(t, tc.msg.Route(), types.RouterKey)
			require.Equal(t, tc.msg.Type(), types.TypeMsgRepayRequest)

			err := tc.msg.ValidateBasic()

			if tc.isErrExp {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestNewMsgLiquidateRequest(t *testing.T) {
	testCases := []struct {
		name     string
		msg      *types.MsgCloseRequest
		isErrExp bool
	}{
		{
			name: "empty from",
			msg: types.NewMsgLiquidateRequest(
				sdk.AccAddress([]byte("")), 1, 1, 1,
			),
			isErrExp: true,
		},
		{
			name: "vaultID zero",
			msg: types.NewMsgLiquidateRequest(
				sdk.MustAccAddressFromBech32("cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t"), 1, 1, 0,
			),
			isErrExp: true,
		},
		{
			name: "valid case",
			msg: types.NewMsgLiquidateRequest(
				sdk.MustAccAddressFromBech32("cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t"), 1, 1, 1,
			),
			isErrExp: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {

			require.Equal(t, tc.msg.Route(), types.RouterKey)
			require.Equal(t, tc.msg.Type(), types.TypeMsgLiquidateRequest)

			err := tc.msg.ValidateBasic()

			if tc.isErrExp {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestNewMsgCreateStableMintRequest(t *testing.T) {
	testCases := []struct {
		name     string
		msg      *types.MsgCreateStableMintRequest
		isErrExp bool
	}{
		{
			name: "empty from",
			msg: types.NewMsgCreateStableMintRequest(
				sdk.AccAddress([]byte("")), 1, 1, sdk.NewInt(50),
			),
			isErrExp: true,
		},
		{
			name: "amount nil",
			msg: types.NewMsgCreateStableMintRequest(
				sdk.MustAccAddressFromBech32("cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t"), 1, 1, sdk.Int{},
			),
			isErrExp: true,
		},
		{
			name: "amount negative",
			msg: types.NewMsgCreateStableMintRequest(
				sdk.MustAccAddressFromBech32("cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t"), 1, 1, sdk.NewInt(-50),
			),
			isErrExp: true,
		},
		{
			name: "amount Zero",
			msg: types.NewMsgCreateStableMintRequest(
				sdk.MustAccAddressFromBech32("cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t"), 1, 1, sdk.NewInt(0),
			),
			isErrExp: true,
		},
		{
			name: "valid case",
			msg: types.NewMsgCreateStableMintRequest(
				sdk.MustAccAddressFromBech32("cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t"), 1, 1, sdk.NewInt(25),
			),
			isErrExp: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {

			require.Equal(t, tc.msg.Route(), types.RouterKey)
			require.Equal(t, tc.msg.Type(), types.TypeMsgCreateStableMintRequest)

			err := tc.msg.ValidateBasic()

			if tc.isErrExp {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestNewMsgDepositStableMintRequest(t *testing.T) {
	testCases := []struct {
		name     string
		msg      *types.MsgDepositStableMintRequest
		isErrExp bool
	}{
		{
			name: "empty from",
			msg: types.NewMsgDepositStableMintRequest(
				sdk.AccAddress([]byte("")), 1, 1, sdk.NewInt(50), 1,
			),
			isErrExp: true,
		},
		{
			name: "amount nil",
			msg: types.NewMsgDepositStableMintRequest(
				sdk.MustAccAddressFromBech32("cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t"), 1, 1, sdk.Int{}, 1,
			),
			isErrExp: true,
		},
		{
			name: "amount negative",
			msg: types.NewMsgDepositStableMintRequest(
				sdk.MustAccAddressFromBech32("cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t"), 1, 1, sdk.NewInt(-50), 1,
			),
			isErrExp: true,
		},
		{
			name: "amount Zero",
			msg: types.NewMsgDepositStableMintRequest(
				sdk.MustAccAddressFromBech32("cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t"), 1, 1, sdk.NewInt(0), 1,
			),
			isErrExp: true,
		},
		{
			name: "valid case",
			msg: types.NewMsgDepositStableMintRequest(
				sdk.MustAccAddressFromBech32("cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t"), 1, 1, sdk.NewInt(25), 1,
			),
			isErrExp: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {

			require.Equal(t, tc.msg.Route(), types.RouterKey)
			require.Equal(t, tc.msg.Type(), types.TypeMsgDepositStableMintRequest)

			err := tc.msg.ValidateBasic()

			if tc.isErrExp {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestNewMsgWithdrawStableMintRequest(t *testing.T) {
	testCases := []struct {
		name     string
		msg      *types.MsgWithdrawStableMintRequest
		isErrExp bool
	}{
		{
			name: "empty from",
			msg: types.NewMsgWithdrawStableMintRequest(
				sdk.AccAddress([]byte("")), 1, 1, sdk.NewInt(50), 1,
			),
			isErrExp: true,
		},
		{
			name: "amount nil",
			msg: types.NewMsgWithdrawStableMintRequest(
				sdk.MustAccAddressFromBech32("cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t"), 1, 1, sdk.Int{}, 1,
			),
			isErrExp: true,
		},
		{
			name: "amount negative",
			msg: types.NewMsgWithdrawStableMintRequest(
				sdk.MustAccAddressFromBech32("cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t"), 1, 1, sdk.NewInt(-50), 1,
			),
			isErrExp: true,
		},
		{
			name: "amount Zero",
			msg: types.NewMsgWithdrawStableMintRequest(
				sdk.MustAccAddressFromBech32("cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t"), 1, 1, sdk.NewInt(0), 1,
			),
			isErrExp: true,
		},
		{
			name: "valid case",
			msg: types.NewMsgWithdrawStableMintRequest(
				sdk.MustAccAddressFromBech32("cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t"), 1, 1, sdk.NewInt(25), 1,
			),
			isErrExp: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {

			require.Equal(t, tc.msg.Route(), types.RouterKey)
			require.Equal(t, tc.msg.Type(), types.TypeMsgWithdrawStableMintRequest)

			err := tc.msg.ValidateBasic()

			if tc.isErrExp {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
