package wasm

import (
	sdkmath "cosmossdk.io/math"
	"testing"

	"github.com/comdex-official/comdex/app/wasm"
	"github.com/comdex-official/comdex/app/wasm/bindings"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestWhitelistAssetLocker(t *testing.T) {
	actor := RandomAccountAddress()
	comdex, ctx := SetupCustomApp(t)
	AddAppAsset(comdex, *ctx)
	querier := wasm.NewQueryPlugin(&comdex.AssetKeeper,
		&comdex.LockerKeeper,
		&comdex.TokenmintKeeper,
		&comdex.Rewardskeeper,
		&comdex.CollectorKeeper,
		&comdex.LiquidationKeeper,
		&comdex.EsmKeeper,
		&comdex.VaultKeeper,
		&comdex.LendKeeper,
		&comdex.LiquidityKeeper,
		&comdex.MarketKeeper)
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

func TestAddMsgAddExtendedPairsVault(t *testing.T) {
	actor := RandomAccountAddress()
	comdex, ctx := SetupCustomApp(t)
	AddPair(comdex, *ctx)
	querier := wasm.NewQueryPlugin(&comdex.AssetKeeper,
		&comdex.LockerKeeper,
		&comdex.TokenmintKeeper,
		&comdex.Rewardskeeper,
		&comdex.CollectorKeeper,
		&comdex.LiquidationKeeper,
		&comdex.EsmKeeper,
		&comdex.VaultKeeper,
		&comdex.LendKeeper,
		&comdex.LiquidityKeeper,
		&comdex.MarketKeeper)
	for _, tc := range []struct {
		name            string
		msg             *bindings.MsgAddExtendedPairsVault
		isErrorExpected bool
	}{
		{
			name: "Add Extended Pair Vaultr",
			msg: &bindings.MsgAddExtendedPairsVault{
				AppID:               1,
				PairID:              1,
				StabilityFee:        sdkmath.LegacyMustNewDecFromStr("0.01"),
				ClosingFee:          sdkmath.LegacyMustNewDecFromStr("0"),
				LiquidationPenalty:  sdkmath.LegacyMustNewDecFromStr("0.12"),
				DrawDownFee:         sdkmath.LegacyMustNewDecFromStr("0.01"),
				IsVaultActive:       true,
				DebtCeiling:         sdkmath.NewInt(1000000000000),
				DebtFloor:           sdkmath.NewInt(1000000),
				IsStableMintVault:   false,
				MinCr:               sdkmath.LegacyMustNewDecFromStr("1.5"),
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

func TestMsgSetCollectorLookupTable(t *testing.T) {
	actor := RandomAccountAddress()
	comdex, ctx := SetupCustomApp(t)
	AddPair(comdex, *ctx)
	querier := wasm.NewQueryPlugin(&comdex.AssetKeeper,
		&comdex.LockerKeeper,
		&comdex.TokenmintKeeper,
		&comdex.Rewardskeeper,
		&comdex.CollectorKeeper,
		&comdex.LiquidationKeeper,
		&comdex.EsmKeeper,
		&comdex.VaultKeeper,
		&comdex.LendKeeper,
		&comdex.LiquidityKeeper,
		&comdex.MarketKeeper)
	for _, tc := range []struct {
		name            string
		msg             *bindings.MsgSetCollectorLookupTable
		isErrorExpected bool
	}{
		{
			name: "Add  Collector Lookup Table",
			msg: &bindings.MsgSetCollectorLookupTable{
				AppID:            1,
				CollectorAssetID: 2,
				SecondaryAssetID: 3,
				SurplusThreshold: sdkmath.NewInt(10000000),
				DebtThreshold:    sdkmath.NewInt(5000000),
				LockerSavingRate: sdkmath.LegacyMustNewDecFromStr("0.1"),
				LotSize:          sdkmath.NewInt(2000000),
				BidFactor:        sdkmath.LegacyMustNewDecFromStr("0.01"),
				DebtLotSize:      sdkmath.NewInt(2000000),
			},
			isErrorExpected: false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			found, _ := querier.CollectorLookupTableQueryCheck(*ctx, tc.msg.AppID, tc.msg.CollectorAssetID, tc.msg.SecondaryAssetID)

			require.True(t, found)
			err := wasm.MsgSetCollectorLookupTable(comdex.CollectorKeeper, *ctx, actor, tc.msg)
			if tc.isErrorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				found, _ := querier.CollectorLookupTableQueryCheck(*ctx, tc.msg.AppID, tc.msg.CollectorAssetID, tc.msg.SecondaryAssetID)
				require.False(t, found)
			}
		})
	}
}

func TestMsgSetAuctionMappingForApp(t *testing.T) {
	actor := RandomAccountAddress()
	comdex, ctx := SetupCustomApp(t)
	AddPair(comdex, *ctx)
	querier := wasm.NewQueryPlugin(&comdex.AssetKeeper,
		&comdex.LockerKeeper,
		&comdex.TokenmintKeeper,
		&comdex.Rewardskeeper,
		&comdex.CollectorKeeper,
		&comdex.LiquidationKeeper,
		&comdex.EsmKeeper,
		&comdex.VaultKeeper,
		&comdex.LendKeeper,
		&comdex.LiquidityKeeper,
		&comdex.MarketKeeper)
	for _, tc := range []struct {
		name            string
		msg             *bindings.MsgSetAuctionMappingForApp
		isErrorExpected bool
	}{
		{
			name: "Add  Collector Lookup Table",
			msg: &bindings.MsgSetAuctionMappingForApp{
				AppID:                1,
				AssetIDs:             uint64(2),
				IsSurplusAuctions:    true,
				IsDebtAuctions:       false,
				IsDistributor:        false,
				AssetOutOraclePrices: false,
				AssetOutPrices:       1000000,
			},
			isErrorExpected: false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			found, _ := querier.AuctionMappingForAppQueryCheck(*ctx, tc.msg.AppID)

			require.True(t, found)
			err := wasm.MsgSetAuctionMappingForApp(comdex.CollectorKeeper, *ctx, actor, tc.msg)
			if tc.isErrorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				found, _ := querier.AuctionMappingForAppQueryCheck(*ctx, tc.msg.AppID)
				require.True(t, found)
			}
		})
	}
}

func TestMsgUpdateCollectorLookupTable(t *testing.T) {
	actor := RandomAccountAddress()
	comdex, ctx := SetupCustomApp(t)
	AddPair(comdex, *ctx)
	AddCollectorLookuptable(comdex, *ctx)
	querier := wasm.NewQueryPlugin(&comdex.AssetKeeper,
		&comdex.LockerKeeper,
		&comdex.TokenmintKeeper,
		&comdex.Rewardskeeper,
		&comdex.CollectorKeeper,
		&comdex.LiquidationKeeper,
		&comdex.EsmKeeper,
		&comdex.VaultKeeper,
		&comdex.LendKeeper,
		&comdex.LiquidityKeeper,
		&comdex.MarketKeeper)
	for _, tc := range []struct {
		name            string
		msg             *bindings.MsgUpdateCollectorLookupTable
		isErrorExpected bool
	}{
		{
			name: "Add  Collector Lookup Table",
			msg: &bindings.MsgUpdateCollectorLookupTable{
				AppID:            1,
				AssetID:          2,
				SurplusThreshold: sdkmath.NewInt(9999),
				DebtThreshold:    sdkmath.NewInt(99),
				LSR:              sdkmath.LegacyMustNewDecFromStr("0.001"),
				LotSize:          sdkmath.NewInt(100),
				BidFactor:        sdkmath.LegacyMustNewDecFromStr("0.00001"),
				DebtLotSize:      sdkmath.NewInt(300),
			},
			isErrorExpected: false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			found, _ := querier.UpdateCollectorLookupTableQueryCheck(*ctx, tc.msg.AppID, tc.msg.AssetID)

			require.True(t, found)
			err := wasm.MsgUpdateCollectorLookupTable(comdex.CollectorKeeper, *ctx, actor, tc.msg)
			if tc.isErrorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				found, _ := querier.UpdateCollectorLookupTableQueryCheck(*ctx, tc.msg.AppID, tc.msg.AssetID)
				require.True(t, found)
			}
		})
	}
}

func TestMsgUpdatePairsVault(t *testing.T) {
	actor := RandomAccountAddress()
	comdex, ctx := SetupCustomApp(t)
	AddPair(comdex, *ctx)
	AddExtendedPairVault(comdex, *ctx)
	querier := wasm.NewQueryPlugin(&comdex.AssetKeeper,
		&comdex.LockerKeeper,
		&comdex.TokenmintKeeper,
		&comdex.Rewardskeeper,
		&comdex.CollectorKeeper,
		&comdex.LiquidationKeeper,
		&comdex.EsmKeeper,
		&comdex.VaultKeeper,
		&comdex.LendKeeper,
		&comdex.LiquidityKeeper,
		&comdex.MarketKeeper)
	for _, tc := range []struct {
		name            string
		msg             *bindings.MsgUpdatePairsVault
		isErrorExpected bool
	}{
		{
			name: "Add  Collector Lookup Table",
			msg: &bindings.MsgUpdatePairsVault{
				AppID:              1,
				ExtPairID:          1,
				StabilityFee:       sdkmath.LegacyMustNewDecFromStr("0.4"),
				ClosingFee:         sdkmath.LegacyMustNewDecFromStr("233.23"),
				LiquidationPenalty: sdkmath.LegacyMustNewDecFromStr("0.56"),
				DrawDownFee:        sdkmath.LegacyMustNewDecFromStr("0.29"),
				DebtCeiling:        sdkmath.NewInt(1000000000),
				DebtFloor:          sdkmath.NewInt(1000),
				MinCr:              sdkmath.LegacyMustNewDecFromStr("1.8"),
				MinUsdValueLeft:    100000000,
			},
			isErrorExpected: false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			found, _ := querier.UpdatePairsVaultQueryCheck(*ctx, tc.msg.AppID, tc.msg.ExtPairID)

			require.True(t, found)
			err := wasm.MsgUpdatePairsVault(comdex.AssetKeeper, *ctx, actor, tc.msg)
			if tc.isErrorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				found, _ := querier.UpdatePairsVaultQueryCheck(*ctx, tc.msg.AppID, tc.msg.ExtPairID)
				require.True(t, found)
			}
		})
	}
}

// func MsgWhitelistAppIDLiquidation(liquidationKeeper liquidationkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress,
//	a *bindings.MsgWhitelistAppIDLiquidation)

func TestMsgWhitelistAppIDLiquidation(t *testing.T) {
	actor := RandomAccountAddress()
	comdex, ctx := SetupCustomApp(t)
	AddPair(comdex, *ctx)
	querier := wasm.NewQueryPlugin(&comdex.AssetKeeper,
		&comdex.LockerKeeper,
		&comdex.TokenmintKeeper,
		&comdex.Rewardskeeper,
		&comdex.CollectorKeeper,
		&comdex.LiquidationKeeper,
		&comdex.EsmKeeper,
		&comdex.VaultKeeper,
		&comdex.LendKeeper,
		&comdex.LiquidityKeeper,
		&comdex.MarketKeeper)
	for _, tc := range []struct {
		name            string
		msg             *bindings.MsgWhitelistAppIDLiquidation
		isErrorExpected bool
	}{
		{
			name: "Add  Collector Lookup Table",
			msg: &bindings.MsgWhitelistAppIDLiquidation{
				AppID: 1,
			},
			isErrorExpected: false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			found, _ := querier.WasmWhitelistAppIDLiquidationQueryCheck(*ctx, tc.msg.AppID)

			require.True(t, found)
			err := wasm.MsgWhitelistAppIDLiquidation(comdex.LiquidationKeeper, *ctx, actor, tc.msg)
			if tc.isErrorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				found, _ := querier.WasmWhitelistAppIDLiquidationQueryCheck(*ctx, tc.msg.AppID)
				require.False(t, found)
			}
		})
	}
}

// func MsgRemoveWhitelistAppIDLiquidation(liquidationKeeper liquidationkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress,
//	a *bindings.MsgRemoveWhitelistAppIDLiquidation)

func TestMsgRemoveWhitelistAppIDLiquidation(t *testing.T) {
	actor := RandomAccountAddress()
	comdex, ctx := SetupCustomApp(t)
	AddPair(comdex, *ctx)
	WhitelistAppIDLiquidation(comdex, *ctx)
	querier := wasm.NewQueryPlugin(&comdex.AssetKeeper,
		&comdex.LockerKeeper,
		&comdex.TokenmintKeeper,
		&comdex.Rewardskeeper,
		&comdex.CollectorKeeper,
		&comdex.LiquidationKeeper,
		&comdex.EsmKeeper,
		&comdex.VaultKeeper,
		&comdex.LendKeeper,
		&comdex.LiquidityKeeper,
		&comdex.MarketKeeper)
	for _, tc := range []struct {
		name            string
		msg             *bindings.MsgRemoveWhitelistAppIDLiquidation
		isErrorExpected bool
	}{
		{
			name: "Add  Collector Lookup Table",
			msg: &bindings.MsgRemoveWhitelistAppIDLiquidation{
				AppID: 1,
			},
			isErrorExpected: false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			found, _ := querier.WasmRemoveWhitelistAppIDLiquidationQueryCheck(*ctx, tc.msg.AppID)

			require.True(t, found)
			err := wasm.MsgRemoveWhitelistAppIDLiquidation(comdex.LiquidationKeeper, *ctx, actor, tc.msg)
			if tc.isErrorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				found, _ := querier.WasmRemoveWhitelistAppIDLiquidationQueryCheck(*ctx, tc.msg.AppID)
				require.False(t, found)
			}
		})
	}
}

func TestMsgAddAuctionParams(t *testing.T) {
	actor := RandomAccountAddress()
	comdex, ctx := SetupCustomApp(t)
	AddPair(comdex, *ctx)
	querier := wasm.NewQueryPlugin(&comdex.AssetKeeper,
		&comdex.LockerKeeper,
		&comdex.TokenmintKeeper,
		&comdex.Rewardskeeper,
		&comdex.CollectorKeeper,
		&comdex.LiquidationKeeper,
		&comdex.EsmKeeper,
		&comdex.VaultKeeper,
		&comdex.LendKeeper,
		&comdex.LiquidityKeeper,
		&comdex.MarketKeeper)
	for _, tc := range []struct {
		name            string
		msg             *bindings.MsgAddAuctionParams
		isErrorExpected bool
	}{
		{
			name: "Add Auction Params",
			msg: &bindings.MsgAddAuctionParams{
				AppID:                  1,
				AuctionDurationSeconds: 300,
				Buffer:                 sdkmath.LegacyMustNewDecFromStr("1.2"),
				Cusp:                   sdkmath.LegacyMustNewDecFromStr("0.6"),
				Step:                   1,
				PriceFunctionType:      1,
				SurplusID:              1,
				DebtID:                 2,
				DutchID:                3,
				BidDurationSeconds:     300,
			},
			isErrorExpected: false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			found, _ := querier.AuctionMappingForAppQueryCheck(*ctx, tc.msg.AppID)

			require.True(t, found)
			err := wasm.MsgAddAuctionParams(comdex.AuctionKeeper, *ctx, actor, tc.msg)
			if tc.isErrorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				found, _ := querier.AuctionMappingForAppQueryCheck(*ctx, tc.msg.AppID)
				require.True(t, found)
			}
		})
	}
}

// func MsgBurnGovTokensForApp(tokenMintKeeper tokenmintkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress,
//	a *bindings.MsgBurnGovTokensForApp)

func TestMsgBurnGovTokensForApp(t *testing.T) {
	actor := RandomAccountAddress()
	userAddress := "cosmos1q7q90qsl9g0gl2zz0njxwv2a649yqrtyxtnv3v"
	addr, _ := sdk.AccAddressFromBech32(userAddress)
	comdex, ctx := SetupCustomApp(t)
	MsgMintNewTokens(comdex, *ctx)

	for _, tc := range []struct {
		name            string
		msg             *bindings.MsgBurnGovTokensForApp
		isErrorExpected bool
	}{
		{
			name: "Add Auction Params",
			msg: &bindings.MsgBurnGovTokensForApp{
				AppID:  1,
				From:   addr,
				Amount: sdk.NewCoin("uharbor", sdkmath.NewInt(100)),
			},
			isErrorExpected: false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := wasm.MsgBurnGovTokensForApp(comdex.TokenmintKeeper, *ctx, actor, tc.msg)
			if tc.isErrorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
