package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/petrichormoney/petri/x/vault/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestValidateVault(t *testing.T) {
	testCases := []struct {
		name   string
		vault  types.Vault
		expErr bool
	}{
		{
			name: "empty pair id",
			vault: types.Vault{
				ExtendedPairVaultID: 0,
				Owner:               "cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t",
				AmountIn:            sdk.NewInt(10000),
				AmountOut:           sdk.NewInt(5000),
			},
			expErr: true,
		},
		{
			name: "empty owner",
			vault: types.Vault{
				ExtendedPairVaultID: 1,
				Owner:               "",
				AmountIn:            sdk.NewInt(10000),
				AmountOut:           sdk.NewInt(5000),
			},
			expErr: true,
		},
		{
			name: "invalid owner address",
			vault: types.Vault{
				ExtendedPairVaultID: 1,
				Owner:               "cosmos....",
				AmountIn:            sdk.NewInt(10000),
				AmountOut:           sdk.NewInt(5000),
			},
			expErr: true,
		},
		{
			name: "amount_in nil",
			vault: types.Vault{
				ExtendedPairVaultID: 1,
				Owner:               "cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t",
				AmountIn:            sdk.Int{},
				AmountOut:           sdk.NewInt(5000),
			},
			expErr: true,
		},
		{
			name: "amount_in negative",
			vault: types.Vault{
				ExtendedPairVaultID: 1,
				Owner:               "cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t",
				AmountIn:            sdk.NewInt(-123),
				AmountOut:           sdk.NewInt(5000),
			},
			expErr: true,
		},
		{
			name: "amount_out nil",
			vault: types.Vault{
				ExtendedPairVaultID: 1,
				Owner:               "cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t",
				AmountIn:            sdk.NewInt(10000),
				AmountOut:           sdk.Int{},
			},
			expErr: true,
		},
		{
			name: "amount in negative",
			vault: types.Vault{
				ExtendedPairVaultID: 1,
				Owner:               "cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t",
				AmountIn:            sdk.NewInt(10000),
				AmountOut:           sdk.NewInt(-5000),
			},
			expErr: true,
		},
		{
			name: "valid case",
			vault: types.Vault{
				ExtendedPairVaultID: 1,
				Owner:               "cosmos1yq8lgssgxlx9smjhes6ryjasmqmd3ts2559g0t",
				AmountIn:            sdk.NewInt(10000),
				AmountOut:           sdk.NewInt(5000),
			},
			expErr: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := tc.vault.Validate()

			if tc.expErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
