package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestParamKeyTable(t *testing.T) {
	require.NotNil(t, ParamKeyTable())
}

func TestValidateVotePeriod(t *testing.T) {
	err := validateVotePeriod("invalidUint64")
	require.ErrorContains(t, err, "invalid parameter type: string")

	err = validateVotePeriod(uint64(0))
	require.ErrorContains(t, err, "oracle parameter VotePeriod must be > 0")

	err = validateVotePeriod(uint64(10))
	require.Nil(t, err)
}

func TestValidateVoteThreshold(t *testing.T) {
	tcs := []struct {
		name   string
		t      sdk.Dec
		errMsg string
	}{
		{"fail: negative", sdk.MustNewDecFromStr("-1"), "threshold must be"},
		{"fail: zero", sdk.ZeroDec(), "threshold must be"},
		{"fail: less than 0.33", sdk.MustNewDecFromStr("0.3"), "threshold must be"},
		{"fail: equal 0.33", sdk.MustNewDecFromStr("0.33"), "threshold must be"},
		{"fail: more than 1", sdk.MustNewDecFromStr("1.1"), "threshold must be"},
		{"fail: more than 1", sdk.MustNewDecFromStr("10"), "threshold must be"},
		{"fail: max precision 2", sdk.MustNewDecFromStr("0.333"), "maximum 2 decimals"},
		{"fail: max precision 2", sdk.MustNewDecFromStr("0.401"), "maximum 2 decimals"},
		{"fail: max precision 2", sdk.MustNewDecFromStr("0.409"), "maximum 2 decimals"},
		{"fail: max precision 2", sdk.MustNewDecFromStr("0.4009"), "maximum 2 decimals"},
		{"fail: max precision 2", sdk.MustNewDecFromStr("0.999"), "maximum 2 decimals"},

		{"ok: 1", sdk.MustNewDecFromStr("1"), ""},
		{"ok: 0.34", sdk.MustNewDecFromStr("0.34"), ""},
		{"ok: 0.99", sdk.MustNewDecFromStr("0.99"), ""},
	}

	for _, tc := range tcs {
		err := validateVoteThreshold(tc.t)
		if tc.errMsg == "" {
			require.NoError(t, err, "test_case", tc.name)
		} else {
			require.ErrorContains(t, err, tc.errMsg, tc.name)
		}
	}
}

func TestValidateRewardBand(t *testing.T) {
	err := validateRewardBand("invalidSdkType")
	require.ErrorContains(t, err, "invalid parameter type: string")

	err = validateRewardBand(sdk.MustNewDecFromStr("-0.31"))
	require.ErrorContains(t, err, "oracle parameter RewardBand must be between [0, 1]")

	err = validateRewardBand(sdk.MustNewDecFromStr("40.0"))
	require.ErrorContains(t, err, "oracle parameter RewardBand must be between [0, 1]")

	err = validateRewardBand(sdk.OneDec())
	require.Nil(t, err)
}

func TestValidateRewardDistributionWindow(t *testing.T) {
	err := validateRewardDistributionWindow("invalidUint64")
	require.ErrorContains(t, err, "invalid parameter type: string")

	err = validateRewardDistributionWindow(uint64(0))
	require.ErrorContains(t, err, "oracle parameter RewardDistributionWindow must be > 0")

	err = validateRewardDistributionWindow(uint64(10))
	require.Nil(t, err)
}

func TestValidateDenomList(t *testing.T) {
	err := validateDenomList("invalidUint64")
	require.ErrorContains(t, err, "invalid parameter type: string")

	err = validateDenomList(DenomList{
		{BaseDenom: ""},
	})
	require.ErrorContains(t, err, "oracle parameter AcceptList Denom must have BaseDenom")

	err = validateDenomList(DenomList{
		{BaseDenom: DenomOjo.BaseDenom, SymbolDenom: ""},
	})
	require.ErrorContains(t, err, "oracle parameter AcceptList Denom must have SymbolDenom")

	err = validateDenomList(DenomList{
		{BaseDenom: DenomOjo.BaseDenom, SymbolDenom: DenomOjo.SymbolDenom},
	})
	require.Nil(t, err)
}

func TestValidateSlashFraction(t *testing.T) {
	err := validateSlashFraction("invalidSdkType")
	require.ErrorContains(t, err, "invalid parameter type: string")

	err = validateSlashFraction(sdk.MustNewDecFromStr("-0.31"))
	require.ErrorContains(t, err, "oracle parameter SlashFraction must be between [0, 1]")

	err = validateSlashFraction(sdk.MustNewDecFromStr("40.0"))
	require.ErrorContains(t, err, "oracle parameter SlashFraction must be between [0, 1]")

	err = validateSlashFraction(sdk.OneDec())
	require.Nil(t, err)
}

func TestValidateSlashWindow(t *testing.T) {
	err := validateSlashWindow("invalidUint64")
	require.ErrorContains(t, err, "invalid parameter type: string")

	err = validateSlashWindow(uint64(0))
	require.ErrorContains(t, err, "oracle parameter SlashWindow must be > 0")

	err = validateSlashWindow(uint64(10))
	require.Nil(t, err)
}

func TestValidateMinValidPerWindow(t *testing.T) {
	err := validateMinValidPerWindow("invalidSdkType")
	require.ErrorContains(t, err, "invalid parameter type: string")

	err = validateMinValidPerWindow(sdk.MustNewDecFromStr("-0.31"))
	require.ErrorContains(t, err, "oracle parameter MinValidPerWindow must be between [0, 1]")

	err = validateMinValidPerWindow(sdk.MustNewDecFromStr("40.0"))
	require.ErrorContains(t, err, "oracle parameter MinValidPerWindow must be between [0, 1]")

	err = validateMinValidPerWindow(sdk.OneDec())
	require.Nil(t, err)
}

func TestParamsEqual(t *testing.T) {
	p1 := DefaultParams()
	err := p1.Validate()
	require.NoError(t, err)

	// minus vote period
	p1.VotePeriod = 0
	err = p1.Validate()
	require.Error(t, err)

	// small vote threshold
	p2 := DefaultParams()
	p2.VoteThreshold = sdk.ZeroDec()
	err = p2.Validate()
	require.Error(t, err)

	// negative reward band
	p3 := DefaultParams()
	p3.RewardBands[0].RewardBand = sdk.NewDecWithPrec(-1, 2)
	err = p3.Validate()
	require.Error(t, err)

	// negative slash fraction
	p4 := DefaultParams()
	p4.SlashFraction = sdk.NewDec(-1)
	err = p4.Validate()
	require.Error(t, err)

	// negative min valid per window
	p5 := DefaultParams()
	p5.MinValidPerWindow = sdk.NewDec(-1)
	err = p5.Validate()
	require.Error(t, err)

	// small slash window
	p6 := DefaultParams()
	p6.SlashWindow = 0
	err = p6.Validate()
	require.Error(t, err)

	// slash window not a multiple of vote period
	p7 := DefaultParams()
	p7.SlashWindow = 7
	err = p7.Validate()
	require.Error(t, err)

	// small distribution window
	p8 := DefaultParams()
	p8.RewardDistributionWindow = 0
	err = p8.Validate()
	require.Error(t, err)

	// empty name
	p9 := DefaultParams()
	p9.AcceptList[0].BaseDenom = ""
	p9.AcceptList[0].SymbolDenom = "ATOM"
	err = p9.Validate()
	require.Error(t, err)

	// empty
	p10 := DefaultParams()
	p10.AcceptList[0].BaseDenom = "uatom"
	p10.AcceptList[0].SymbolDenom = ""
	err = p10.Validate()
	require.Error(t, err)

	p11 := DefaultParams()
	require.NotNil(t, p11.ParamSetPairs())
	require.NotNil(t, p11.String())
}
