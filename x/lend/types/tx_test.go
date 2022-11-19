package types_test

import (
	"testing"

	"github.com/petrichormoney/petri/x/lend/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestNewMsgLend(t *testing.T) {
	testCases := []struct {
		name     string
		msg      *types.MsgLend
		isErrExp bool
	}{
		{
			name:     "empty from",
			msg:      types.NewMsgLend("", 1, sdk.NewCoin("upetri", sdk.NewInt(1000000)), 1, 1),
			isErrExp: true,
		},
		{
			name: "amountIn zero",
			msg: types.NewMsgLend(
				"cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t",
				1,
				sdk.NewCoin("upetri", sdk.NewInt(0)),
				1,
				1,
			),
			isErrExp: true,
		},
		{
			name: "assetID zero",
			msg: types.NewMsgLend(
				"cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t",
				0,
				sdk.NewCoin("upetri", sdk.NewInt(1000000)),
				1,
				1,
			),
			isErrExp: true,
		},
		{
			name: "poolID zero",
			msg: types.NewMsgLend(
				"cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t",
				1,
				sdk.NewCoin("upetri", sdk.NewInt(1000000)),
				0,
				1,
			),
			isErrExp: true,
		},
		{
			name: "appID zero",
			msg: types.NewMsgLend(
				"cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t",
				1,
				sdk.NewCoin("upetri", sdk.NewInt(1000000)),
				1,
				0,
			),
			isErrExp: true,
		},
		{
			name: "valid case",
			msg: types.NewMsgLend(
				"cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t",
				1,
				sdk.NewCoin("upetri", sdk.NewInt(1000000)),
				1,
				1,
			),
			isErrExp: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.msg.Route(), types.RouterKey)
			require.Equal(t, tc.msg.Type(), types.TypeLendAssetRequest)

			err := tc.msg.ValidateBasic()

			if tc.isErrExp {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestNewMsgWithdraw(t *testing.T) {
	testCases := []struct {
		name     string
		msg      *types.MsgWithdraw
		isErrExp bool
	}{
		{
			name: "empty from",
			msg: types.NewMsgWithdraw(
				"",
				1,
				sdk.NewCoin("upetri", sdk.NewInt(1000000)),
			),
			isErrExp: true,
		},
		{
			name: "lendID zero",
			msg: types.NewMsgWithdraw(
				"cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t",
				0,
				sdk.NewCoin("upetri", sdk.NewInt(1000000)),
			),
			isErrExp: true,
		},
		{
			name: "amount zero",
			msg: types.NewMsgWithdraw(
				"cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t",
				0,
				sdk.NewCoin("upetri", sdk.NewInt(0)),
			),
			isErrExp: true,
		},
		{
			name: "valid case",
			msg: types.NewMsgWithdraw(
				"cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t",
				1,
				sdk.NewCoin("upetri", sdk.NewInt(1000000)),
			),
			isErrExp: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.msg.Route(), types.RouterKey)
			require.Equal(t, tc.msg.Type(), types.TypeWithdrawAssetRequest)

			err := tc.msg.ValidateBasic()

			if tc.isErrExp {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestNewMsgDeposit(t *testing.T) {
	testCases := []struct {
		name     string
		msg      *types.MsgDeposit
		isErrExp bool
	}{
		{
			name: "empty from",
			msg: types.NewMsgDeposit(
				"",
				1,
				sdk.NewCoin("upetri", sdk.NewInt(1000000)),
			),
			isErrExp: true,
		},
		{
			name: "lendID zero",
			msg: types.NewMsgDeposit(
				"cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t",
				0,
				sdk.NewCoin("upetri", sdk.NewInt(1000000)),
			),
			isErrExp: true,
		},
		{
			name: "amount zero",
			msg: types.NewMsgDeposit(
				"cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t",
				0,
				sdk.NewCoin("upetri", sdk.NewInt(0)),
			),
			isErrExp: true,
		},
		{
			name: "valid case",
			msg: types.NewMsgDeposit(
				"cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t",
				1,
				sdk.NewCoin("upetri", sdk.NewInt(1000000)),
			),
			isErrExp: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.msg.Route(), types.RouterKey)
			require.Equal(t, tc.msg.Type(), types.TypeDepositAssetRequest)

			err := tc.msg.ValidateBasic()

			if tc.isErrExp {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestNewMsgCloseLend(t *testing.T) {
	testCases := []struct {
		name     string
		msg      *types.MsgCloseLend
		isErrExp bool
	}{
		{
			name: "empty from",
			msg: types.NewMsgCloseLend(
				"",
				1,
			),
			isErrExp: true,
		},
		{
			name: "lendID zero",
			msg: types.NewMsgCloseLend(
				"cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t",
				0,
			),
			isErrExp: true,
		},
		{
			name: "valid case",
			msg: types.NewMsgCloseLend(
				"cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t",
				1,
			),
			isErrExp: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.msg.Route(), types.RouterKey)
			require.Equal(t, tc.msg.Type(), types.TypeCloseLendAssetRequest)

			err := tc.msg.ValidateBasic()

			if tc.isErrExp {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestNewMsgBorrow(t *testing.T) {
	testCases := []struct {
		name     string
		msg      *types.MsgBorrow
		isErrExp bool
	}{
		{
			name: "empty from",
			msg: types.NewMsgBorrow(
				"",
				1,
				1,
				false,
				sdk.NewCoin("upetri", sdk.NewInt(1000000)),
				sdk.NewCoin("upetri", sdk.NewInt(1000000)),
			),
			isErrExp: true,
		},
		{
			name: "amountIn zero",
			msg: types.NewMsgBorrow(
				"cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t",
				1,
				1,
				false,
				sdk.NewCoin("upetri", sdk.NewInt(0)),
				sdk.NewCoin("upetri", sdk.NewInt(1000000)),
			),
			isErrExp: true,
		},
		{
			name: "amountOut zero",
			msg: types.NewMsgBorrow(
				"cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t",
				1,
				1,
				false,
				sdk.NewCoin("upetri", sdk.NewInt(1000000)),
				sdk.NewCoin("upetri", sdk.NewInt(0)),
			),
			isErrExp: true,
		},
		{
			name: "lendID zero",
			msg: types.NewMsgBorrow(
				"cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t",
				0,
				1,
				false,
				sdk.NewCoin("upetri", sdk.NewInt(1000000)),
				sdk.NewCoin("upetri", sdk.NewInt(1000000)),
			),
			isErrExp: true,
		},
		{
			name: "pairID zero",
			msg: types.NewMsgBorrow(
				"cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t",
				1,
				0,
				false,
				sdk.NewCoin("upetri", sdk.NewInt(1000000)),
				sdk.NewCoin("upetri", sdk.NewInt(1000000)),
			),
			isErrExp: true,
		},
		{
			name: "valid case",
			msg: types.NewMsgBorrow(
				"cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t",
				1,
				1,
				false,
				sdk.NewCoin("upetri", sdk.NewInt(1000000)),
				sdk.NewCoin("upetri", sdk.NewInt(1000000)),
			),
			isErrExp: false,
		},
		{
			name: "valid case",
			msg: types.NewMsgBorrow(
				"cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t",
				1,
				1,
				true,
				sdk.NewCoin("upetri", sdk.NewInt(1000000)),
				sdk.NewCoin("upetri", sdk.NewInt(1000000)),
			),
			isErrExp: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.msg.Route(), types.RouterKey)
			require.Equal(t, tc.msg.Type(), types.TypeBorrowAssetRequest)

			err := tc.msg.ValidateBasic()

			if tc.isErrExp {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestNewMsgRepay(t *testing.T) {
	testCases := []struct {
		name     string
		msg      *types.MsgRepay
		isErrExp bool
	}{
		{
			name: "empty from",
			msg: types.NewMsgRepay(
				"",
				1,
				sdk.NewCoin("upetri", sdk.NewInt(1000000)),
			),
			isErrExp: true,
		},
		{
			name: "borrowID zero",
			msg: types.NewMsgRepay(
				"cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t",
				0,
				sdk.NewCoin("upetri", sdk.NewInt(1000000)),
			),
			isErrExp: true,
		},
		{
			name: "amount zero",
			msg: types.NewMsgRepay(
				"cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t",
				1,
				sdk.NewCoin("upetri", sdk.NewInt(0)),
			),
			isErrExp: true,
		},
		{
			name: "valid case",
			msg: types.NewMsgRepay(
				"cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t",
				1,
				sdk.NewCoin("upetri", sdk.NewInt(1000000)),
			),
			isErrExp: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.msg.Route(), types.RouterKey)
			require.Equal(t, tc.msg.Type(), types.TypeRepayAssetRequest)

			err := tc.msg.ValidateBasic()

			if tc.isErrExp {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestNewMsgDraw(t *testing.T) {
	testCases := []struct {
		name     string
		msg      *types.MsgDraw
		isErrExp bool
	}{
		{
			name: "empty from",
			msg: &types.MsgDraw{
				Borrower: "",
				BorrowId: 1,
				Amount:   sdk.NewCoin("upetri", sdk.NewInt(1000000)),
			},
			isErrExp: true,
		},
		{
			name: "borrowID zero",
			msg: &types.MsgDraw{
				Borrower: "cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t",
				BorrowId: 0,
				Amount:   sdk.NewCoin("upetri", sdk.NewInt(1000000)),
			},
			isErrExp: true,
		},
		{
			name: "amount zero",
			msg: &types.MsgDraw{
				Borrower: "cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t",
				BorrowId: 1,
				Amount:   sdk.NewCoin("upetri", sdk.NewInt(0)),
			},
			isErrExp: true,
		},
		{
			name: "valid case",
			msg: &types.MsgDraw{
				Borrower: "cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t",
				BorrowId: 1,
				Amount:   sdk.NewCoin("upetri", sdk.NewInt(1000000)),
			},
			isErrExp: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.msg.Route(), types.RouterKey)
			require.Equal(t, tc.msg.Type(), types.TypeDrawAssetRequest)

			err := tc.msg.ValidateBasic()

			if tc.isErrExp {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestNewMsgDepositBorrow(t *testing.T) {
	testCases := []struct {
		name     string
		msg      *types.MsgDepositBorrow
		isErrExp bool
	}{
		{
			name: "empty from",
			msg: &types.MsgDepositBorrow{
				Borrower: "",
				BorrowId: 1,
				Amount:   sdk.NewCoin("upetri", sdk.NewInt(1000000)),
			},
			isErrExp: true,
		},
		{
			name: "borrowID zero",
			msg: &types.MsgDepositBorrow{
				Borrower: "cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t",
				BorrowId: 0,
				Amount:   sdk.NewCoin("upetri", sdk.NewInt(1000000)),
			},
			isErrExp: true,
		},
		{
			name: "amount zero",
			msg: &types.MsgDepositBorrow{
				Borrower: "cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t",
				BorrowId: 1,
				Amount:   sdk.NewCoin("upetri", sdk.NewInt(0)),
			},
			isErrExp: true,
		},
		{
			name: "valid case",
			msg: &types.MsgDepositBorrow{
				Borrower: "cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t",
				BorrowId: 1,
				Amount:   sdk.NewCoin("upetri", sdk.NewInt(1000000)),
			},
			isErrExp: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.msg.Route(), types.RouterKey)
			require.Equal(t, tc.msg.Type(), types.TypeDepositBorrowAssetRequest)

			err := tc.msg.ValidateBasic()

			if tc.isErrExp {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestNewMsgCloseBorrow(t *testing.T) {
	testCases := []struct {
		name     string
		msg      *types.MsgCloseBorrow
		isErrExp bool
	}{
		{
			name: "empty from",
			msg: types.NewMsgCloseBorrow(
				"",
				1,
			),
			isErrExp: true,
		},
		{
			name: "lendID zero",
			msg: types.NewMsgCloseBorrow(
				"cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t",
				0,
			),
			isErrExp: true,
		},
		{
			name: "valid case",
			msg: types.NewMsgCloseBorrow(
				"cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t",
				1,
			),
			isErrExp: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.msg.Route(), types.RouterKey)
			require.Equal(t, tc.msg.Type(), types.TypeCloseBorrowAssetRequest)

			err := tc.msg.ValidateBasic()

			if tc.isErrExp {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestNewMsgBorrowAlternate(t *testing.T) {
	testCases := []struct {
		name     string
		msg      *types.MsgBorrowAlternate
		isErrExp bool
	}{
		{
			name: "empty from",
			msg: types.NewMsgBorrowAlternate(
				"",
				1,
				1,
				sdk.NewCoin("upetri", sdk.NewInt(1000000)),
				1,
				false,
				sdk.NewCoin("uatom", sdk.NewInt(1000000)),
				1,
			),
			isErrExp: true,
		},
		{
			name: "assetID zero",
			msg: types.NewMsgBorrowAlternate(
				"cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t",
				0,
				1,
				sdk.NewCoin("upetri", sdk.NewInt(1000000)),
				1,
				false,
				sdk.NewCoin("uatom", sdk.NewInt(1000000)),
				1,
			),
			isErrExp: true,
		},
		{
			name: "poolID zero",
			msg: types.NewMsgBorrowAlternate(
				"cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t",
				1,
				0,
				sdk.NewCoin("upetri", sdk.NewInt(1000000)),
				1,
				false,
				sdk.NewCoin("uatom", sdk.NewInt(1000000)),
				1,
			),
			isErrExp: true,
		},
		{
			name: "pairID zero",
			msg: types.NewMsgBorrowAlternate(
				"cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t",
				1,
				1,
				sdk.NewCoin("upetri", sdk.NewInt(1000000)),
				0,
				false,
				sdk.NewCoin("uatom", sdk.NewInt(1000000)),
				1,
			),
			isErrExp: true,
		},
		{
			name: "appID zero",
			msg: types.NewMsgBorrowAlternate(
				"cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t",
				1,
				1,
				sdk.NewCoin("upetri", sdk.NewInt(1000000)),
				1,
				false,
				sdk.NewCoin("uatom", sdk.NewInt(1000000)),
				0,
			),
			isErrExp: true,
		},
		{
			name: "amountIn zero",
			msg: types.NewMsgBorrowAlternate(
				"cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t",
				1,
				1,
				sdk.NewCoin("upetri", sdk.NewInt(0)),
				1,
				false,
				sdk.NewCoin("uatom", sdk.NewInt(1000000)),
				1,
			),
			isErrExp: true,
		},
		{
			name: "amountOut zero",
			msg: types.NewMsgBorrowAlternate(
				"cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t",
				1,
				1,
				sdk.NewCoin("upetri", sdk.NewInt(1000000)),
				1,
				false,
				sdk.NewCoin("uatom", sdk.NewInt(0)),
				1,
			),
			isErrExp: true,
		},
		{
			name: "valid case",
			msg: types.NewMsgBorrowAlternate(
				"cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t",
				1,
				1,
				sdk.NewCoin("upetri", sdk.NewInt(1000000)),
				1,
				false,
				sdk.NewCoin("uatom", sdk.NewInt(1000000)),
				1,
			),
			isErrExp: false,
		},
		{
			name: "valid case",
			msg: types.NewMsgBorrowAlternate(
				"cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t",
				1,
				1,
				sdk.NewCoin("upetri", sdk.NewInt(1000000)),
				1,
				true,
				sdk.NewCoin("uatom", sdk.NewInt(1000000)),
				1,
			),
			isErrExp: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.msg.Route(), types.RouterKey)
			require.Equal(t, tc.msg.Type(), types.TypeBorrowAlternateAssetRequest)

			err := tc.msg.ValidateBasic()

			if tc.isErrExp {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestNewMsgFundModuleAccounts(t *testing.T) {
	testCases := []struct {
		name     string
		msg      *types.MsgFundModuleAccounts
		isErrExp bool
	}{
		{
			name: "empty from",
			msg: types.NewMsgFundModuleAccounts(
				1,
				1,
				"",
				sdk.NewCoin("upetri", sdk.NewInt(1000000)),
			),
			isErrExp: true,
		},
		/*{
			name: "empty moduleName",
			msg: types.NewMsgFundModuleAccounts(
				sdk.NewInt(),
				1,
				"cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t",
				sdk.NewCoin("upetri", sdk.NewInt(1000000)),
			),
			isErrExp: true,
		},*/
		{
			name: "amount zero",
			msg: types.NewMsgFundModuleAccounts(
				3,
				1,
				"cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t",
				sdk.NewCoin("upetri", sdk.NewInt(0)),
			),
			isErrExp: true,
		},
		{
			name: "valid case",
			msg: types.NewMsgFundModuleAccounts(
				1,
				1,
				"cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t",
				sdk.NewCoin("upetri", sdk.NewInt(1000000)),
			),
			isErrExp: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.msg.Route(), types.RouterKey)
			require.Equal(t, tc.msg.Type(), types.TypeFundModuleAccountRequest)

			err := tc.msg.ValidateBasic()

			if tc.isErrExp {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
