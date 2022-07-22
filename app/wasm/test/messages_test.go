package wasm

import (
	"github.com/comdex-official/comdex/app/wasm"
	"github.com/comdex-official/comdex/app/wasm/bindings"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestWhitelistAssetLocker(t *testing.T) {
	actor := RandomAccountAddress()
	comdex, ctx := SetupCustomApp()
	AddAppAsset(comdex, *ctx)
	querier := wasm.NewQueryPlugin(&comdex.AssetKeeper,
		&comdex.LockerKeeper,
		&comdex.TokenmintKeeper,
		&comdex.Rewardskeeper,
		&comdex.CollectorKeeper,
		&comdex.LiquidationKeeper,
		&comdex.EsmKeeper)
	for _, tc := range []struct {
		name            string
		msg             *bindings.MsgWhiteListAssetLocker
		isErrorExpected bool
	}{
		{
			name: "Add Whitelist Asset Locker",
			msg: &bindings.MsgWhiteListAssetLocker{
				AppID:   1,
				AssetID: 1,
			},
			isErrorExpected: false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			found, _ := querier.WhiteListedAssetQueryCheck(*ctx, tc.msg.AppID, tc.msg.AssetID)
			require.True(t, found)
			err := wasm.WhiteListAsset(comdex.LockerKeeper, *ctx, actor.String(), tc.msg)
			if tc.isErrorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				found, _ := querier.WhiteListedAssetQueryCheck(*ctx, tc.msg.AppID, tc.msg.AssetID)
				require.False(t, found)
			}
		})
	}
}

//ExtendedPairsVaultRecordsQueryCheck(ctx sdk.Context, appID, pairID uint64, StabilityFee, ClosingFee, DrawDownFee sdk.Dec, DebtCeiling, DebtFloor uint64, PairName string) (found bool, err string)
func TestWhitelistAssetLocker1(t *testing.T) {
	actor := RandomAccountAddress()
	comdex, ctx := SetupCustomApp()
	AddPair(comdex, *ctx)
	querier := wasm.NewQueryPlugin(&comdex.AssetKeeper,
		&comdex.LockerKeeper,
		&comdex.TokenmintKeeper,
		&comdex.Rewardskeeper,
		&comdex.CollectorKeeper,
		&comdex.LiquidationKeeper,
		&comdex.EsmKeeper)
	for _, tc := range []struct {
		name            string
		msg             *bindings.MsgAddExtendedPairsVault
		isErrorExpected bool
	}{
		{
			name: "Add Whitelist Asset Locker",
			msg: &bindings.MsgAddExtendedPairsVault{
				AppID:               1,
				PairID:              1,
				StabilityFee:        sdk.MustNewDecFromStr("0.01"),
				ClosingFee:          sdk.MustNewDecFromStr("0"),
				LiquidationPenalty:  sdk.MustNewDecFromStr("0.12"),
				DrawDownFee:         sdk.MustNewDecFromStr("0.01"),
				IsVaultActive:       true,
				DebtCeiling:         1000000000000,
				DebtFloor:           1000000,
				IsStableMintVault:   false,
				MinCr:               sdk.MustNewDecFromStr("1.5"),
				PairName:            "CMDX-B",
				AssetOutOraclePrice: true,
				AssetOutPrice:       1000000,
				MinUsdValueLeft:     1000000000000,
			},
			isErrorExpected: false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			found, _ := querier.ExtendedPairsVaultRecordsQueryCheck(*ctx, tc.msg.AppID, tc.msg.PairID, tc.msg.StabilityFee, tc.msg.ClosingFee, tc.msg.DrawDownFee, tc.msg.DebtCeiling, tc.msg.DebtFloor, tc.msg.PairName)

			require.True(t, found)
			err := wasm.MsgAddExtendedPairsVault(comdex.AssetKeeper, *ctx, actor, tc.msg)
			if tc.isErrorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				found, _ := querier.ExtendedPairsVaultRecordsQueryCheck(*ctx, tc.msg.AppID, tc.msg.PairID, tc.msg.StabilityFee, tc.msg.ClosingFee, tc.msg.DrawDownFee, tc.msg.DebtCeiling, tc.msg.DebtFloor, tc.msg.PairName)
				require.False(t, found)
			}
		})
	}
}
