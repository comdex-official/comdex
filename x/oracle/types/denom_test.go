package types_test

import (
	"testing"

	"github.com/comdex-official/comdex/x/oracle/types"
	"github.com/stretchr/testify/require"
)

func TestDenomString(t *testing.T) {
	testCases := []struct {
		denom       types.Denom
		expectedStr string
	}{
		{
			denom:       types.DenomOjo,
			expectedStr: "base_denom: uojo\nsymbol_denom: ojo\nexponent: 6\n",
		},
		{
			denom:       types.DenomLuna,
			expectedStr: "base_denom: ibc/0EF15DF2F02480ADE0BB6E85D9EBB5DAEA2836D3860E9F97F9AADE4F57A31AA0\nsymbol_denom: LUNA\nexponent: 6\n",
		},
		{
			denom:       types.DenomAtom,
			expectedStr: "base_denom: ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2\nsymbol_denom: ATOM\nexponent: 6\n",
		},
	}

	for _, testCase := range testCases {
		require.Equal(t, testCase.expectedStr, testCase.denom.String())
	}
}

func TestDenomEqual(t *testing.T) {
	testCases := []struct {
		denom         types.Denom
		denomCompared types.Denom
		equal         bool
	}{
		{
			denom:         types.DenomOjo,
			denomCompared: types.DenomOjo,
			equal:         true,
		},
		{
			denom:         types.DenomOjo,
			denomCompared: types.DenomLuna,
			equal:         false,
		},
		{
			denom:         types.DenomLuna,
			denomCompared: types.DenomLuna,
			equal:         true,
		},
		{
			denom:         types.DenomAtom,
			denomCompared: types.DenomAtom,
			equal:         true,
		},
		{
			denom:         types.DenomAtom,
			denomCompared: types.DenomLuna,
			equal:         false,
		},
	}

	for _, testCase := range testCases {
		require.Equal(t, testCase.equal, testCase.denom.Equal(&testCase.denomCompared))
	}
}

func TestDenomListString(t *testing.T) {
	testCases := []struct {
		denomList   types.DenomList
		expectedStr string
	}{
		{
			denomList:   types.DenomList{types.DenomOjo},
			expectedStr: "base_denom: uojo\nsymbol_denom: ojo\nexponent: 6",
		},
		{
			denomList:   types.DenomList{types.DenomAtom, types.DenomLuna},
			expectedStr: "base_denom: ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2\nsymbol_denom: ATOM\nexponent: 6\n\nbase_denom: ibc/0EF15DF2F02480ADE0BB6E85D9EBB5DAEA2836D3860E9F97F9AADE4F57A31AA0\nsymbol_denom: LUNA\nexponent: 6",
		},
	}

	for _, testCase := range testCases {
		require.Equal(t, testCase.expectedStr, testCase.denomList.String())
	}
}

func TestDenomListContains(t *testing.T) {
	testCases := []struct {
		denomList    types.DenomList
		denomSymbol  string
		symbolInList bool
	}{
		{
			denomList:    types.DenomList{types.DenomOjo},
			denomSymbol:  types.DenomOjo.SymbolDenom,
			symbolInList: true,
		},
		{
			denomList:    types.DenomList{types.DenomOjo},
			denomSymbol:  types.DenomLuna.SymbolDenom,
			symbolInList: false,
		},
		{
			denomList:    types.DenomList{types.DenomOjo, types.DenomAtom},
			denomSymbol:  types.DenomLuna.SymbolDenom,
			symbolInList: false,
		},
		{
			denomList:    types.DenomList{types.DenomOjo, types.DenomAtom},
			denomSymbol:  types.DenomAtom.SymbolDenom,
			symbolInList: true,
		},
		{
			denomList:    types.DenomList{types.DenomOjo, types.DenomAtom, types.DenomLuna},
			denomSymbol:  types.DenomLuna.SymbolDenom,
			symbolInList: true,
		},
	}

	for _, testCase := range testCases {
		require.Equal(t, testCase.symbolInList, testCase.denomList.Contains(testCase.denomSymbol))
	}
}

// TestNormalize makes sure that the Normalize function
// updates the SymbolDenom of a DenomList to all caps.
func TestNormalize(t *testing.T) {
	testCases := []struct {
		denomList       types.DenomList
		expectedSymbols map[string]struct{}
	}{
		{
			denomList: types.DenomList{types.DenomOjo},
			expectedSymbols: map[string]struct{}{
				"OJO": {},
			},
		},
		{
			denomList: types.DenomList{types.DenomAtom, types.DenomLuna},
			expectedSymbols: map[string]struct{}{
				"ATOM": {},
				"LUNA": {},
			},
		},
	}

	for _, testCase := range testCases {
		normalized := testCase.denomList.Normalize()
		for _, v := range normalized {
			_, ok := testCase.expectedSymbols[v.SymbolDenom]
			require.True(t, ok)
		}
	}
}
