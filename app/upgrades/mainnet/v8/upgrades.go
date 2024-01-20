package v8

import (
	"context"
	sdkmath "cosmossdk.io/math"
	"fmt"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/comdex-official/comdex/app/wasm/bindings"
	assetkeeper "github.com/comdex-official/comdex/x/asset/keeper"
	assettypes "github.com/comdex-official/comdex/x/asset/types"
	auctionkeeper "github.com/comdex-official/comdex/x/auction/keeper"
	lendkeeper "github.com/comdex-official/comdex/x/lend/keeper"
	"github.com/comdex-official/comdex/x/lend/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

func UpdateExtendedPairVaultsAndAsset(ctx sdk.Context, assetKeeper assetkeeper.Keeper) {
	extPairs := []*bindings.MsgUpdatePairsVault{
		{
			AppID: 2, ExtPairID: 2, StabilityFee: sdkmath.LegacyMustNewDecFromStr("1"), ClosingFee: sdkmath.LegacyZeroDec(), LiquidationPenalty: sdkmath.LegacyMustNewDecFromStr("0.15"),
			DrawDownFee: sdkmath.LegacyMustNewDecFromStr("0.005"), IsVaultActive: true, DebtCeiling: sdkmath.NewInt(250000000000), DebtFloor: sdkmath.NewInt(50000000), MinCr: sdkmath.LegacyMustNewDecFromStr("1.4"),
			MinUsdValueLeft: 100000,
		},
		{
			AppID: 2, ExtPairID: 3, StabilityFee: sdkmath.LegacyMustNewDecFromStr("0.5"), ClosingFee: sdkmath.LegacyZeroDec(), LiquidationPenalty: sdkmath.LegacyMustNewDecFromStr("0.15"),
			DrawDownFee: sdkmath.LegacyMustNewDecFromStr("0.005"), IsVaultActive: true, DebtCeiling: sdkmath.NewInt(350000000000), DebtFloor: sdkmath.NewInt(50000000), MinCr: sdkmath.LegacyMustNewDecFromStr("1.7"),
			MinUsdValueLeft: 100000,
		},
		{
			AppID: 2, ExtPairID: 4, StabilityFee: sdkmath.LegacyMustNewDecFromStr("0.25"), ClosingFee: sdkmath.LegacyZeroDec(), LiquidationPenalty: sdkmath.LegacyMustNewDecFromStr("0.15"),
			DrawDownFee: sdkmath.LegacyMustNewDecFromStr("0.005"), IsVaultActive: true, DebtCeiling: sdkmath.NewInt(400000000000), DebtFloor: sdkmath.NewInt(50000000), MinCr: sdkmath.LegacyMustNewDecFromStr("2"),
			MinUsdValueLeft: 100000,
		},
		{
			AppID: 2, ExtPairID: 5, StabilityFee: sdkmath.LegacyMustNewDecFromStr("1"), ClosingFee: sdkmath.LegacyZeroDec(), LiquidationPenalty: sdkmath.LegacyMustNewDecFromStr("0.15"),
			DrawDownFee: sdkmath.LegacyMustNewDecFromStr("0.005"), IsVaultActive: true, DebtCeiling: sdkmath.NewInt(250000000000), DebtFloor: sdkmath.NewInt(50000000), MinCr: sdkmath.LegacyMustNewDecFromStr("1.5"),
			MinUsdValueLeft: 100000,
		},
		{
			AppID: 2, ExtPairID: 6, StabilityFee: sdkmath.LegacyMustNewDecFromStr("0.5"), ClosingFee: sdkmath.LegacyZeroDec(), LiquidationPenalty: sdkmath.LegacyMustNewDecFromStr("0.15"),
			DrawDownFee: sdkmath.LegacyMustNewDecFromStr("0.005"), IsVaultActive: true, DebtCeiling: sdkmath.NewInt(350000000000), DebtFloor: sdkmath.NewInt(50000000), MinCr: sdkmath.LegacyMustNewDecFromStr("1.8"),
			MinUsdValueLeft: 100000,
		},
		{
			AppID: 2, ExtPairID: 7, StabilityFee: sdkmath.LegacyMustNewDecFromStr("0.25"), ClosingFee: sdkmath.LegacyZeroDec(), LiquidationPenalty: sdkmath.LegacyMustNewDecFromStr("0.15"),
			DrawDownFee: sdkmath.LegacyMustNewDecFromStr("0.005"), IsVaultActive: true, DebtCeiling: sdkmath.NewInt(400000000000), DebtFloor: sdkmath.NewInt(50000000), MinCr: sdkmath.LegacyMustNewDecFromStr("2.1"),
			MinUsdValueLeft: 100000,
		},
	}
	for _, extPair := range extPairs {
		err := assetKeeper.WasmUpdatePairsVault(ctx, extPair)
		if err != nil {
			fmt.Println("err in updating extended pair ", extPair.ExtPairID)
		}
	}
	asset, found := assetKeeper.GetAsset(ctx, 17)
	if found {
		asset.Denom = "ibc/2ABB3F0A1DA07D7F83D5004A4A16A4D4A264067AA85E15A4885D0AB8C0E4587B"
		assetKeeper.SetAsset(ctx, asset)
	}
}

func Dec(s string) sdkmath.LegacyDec {
	dec, err := sdkmath.LegacyNewDecFromStr(s)
	if err != nil {
		panic(err)
	}
	return dec
}

func UpdateAuctionParams(
	ctx sdk.Context,
	assetKeeper assetkeeper.Keeper,
	lendKeeper lendkeeper.Keeper,
	auctionKeeper auctionkeeper.Keeper,
) {
	// Add cAssets for USDC and stATOM
	// Add Asset Rates for OSMO, USDC, stATOM
	// Update auction params for lend module and Harbor app

	// Adding cAssets
	cUSDC := assettypes.Asset{
		Name:                  "CAXLUSDC",
		Denom:                 "ucaxlusdc",
		Decimals:              sdkmath.NewInt(1000000),
		IsOnChain:             true,
		IsOraclePriceRequired: false,
		IsCdpMintable:         true,
	}
	err := assetKeeper.AddAssetRecords(ctx, cUSDC)
	if err != nil {
		return
	}

	cstATOM := assettypes.Asset{
		Name:                  "CSTATOM",
		Denom:                 "ucstatom",
		Decimals:              sdkmath.NewInt(1000000),
		IsOnChain:             true,
		IsOraclePriceRequired: false,
		IsCdpMintable:         true,
	}
	err = assetKeeper.AddAssetRecords(ctx, cstATOM)
	if err != nil {
		return
	}
	// Adding Asset Rates
	OSMORatesParams := types.AssetRatesParams{
		AssetID:              4,
		UOptimal:             Dec("0.65"),
		Base:                 Dec("0.002"),
		Slope1:               Dec("0.08"),
		Slope2:               Dec("1.5"),
		EnableStableBorrow:   false,
		StableBase:           Dec("0.0"),
		StableSlope1:         Dec("0.0"),
		StableSlope2:         Dec("0.0"),
		Ltv:                  Dec("0.65"),
		LiquidationThreshold: Dec("0.70"),
		LiquidationPenalty:   Dec("0.075"),
		LiquidationBonus:     Dec("0.075"),
		ReserveFactor:        Dec("0.2"),
		CAssetID:             8,
	}
	lendKeeper.SetAssetRatesParams(ctx, OSMORatesParams)
	axlUSDCRatesParams := types.AssetRatesParams{
		AssetID:              10,
		UOptimal:             Dec("0.80"),
		Base:                 Dec("0.002"),
		Slope1:               Dec("0.06"),
		Slope2:               Dec("0.6"),
		EnableStableBorrow:   false,
		StableBase:           Dec("0.0"),
		StableSlope1:         Dec("0.0"),
		StableSlope2:         Dec("0.0"),
		Ltv:                  Dec("0.80"),
		LiquidationThreshold: Dec("0.85"),
		LiquidationPenalty:   Dec("0.05"),
		LiquidationBonus:     Dec("0.05"),
		ReserveFactor:        Dec("0.2"),
		CAssetID:             22,
	}
	lendKeeper.SetAssetRatesParams(ctx, axlUSDCRatesParams)

	stATOMRatesParams := types.AssetRatesParams{
		AssetID:              14,
		UOptimal:             Dec("0.60"),
		Base:                 Dec("0.002"),
		Slope1:               Dec("0.08"),
		Slope2:               Dec("1.60"),
		EnableStableBorrow:   false,
		StableBase:           Dec("0.0"),
		StableSlope1:         Dec("0.0"),
		StableSlope2:         Dec("0.0"),
		Ltv:                  Dec("0.6"),
		LiquidationThreshold: Dec("0.65"),
		LiquidationPenalty:   Dec("0.075"),
		LiquidationBonus:     Dec("0.075"),
		ReserveFactor:        Dec("0.2"),
		CAssetID:             23,
	}
	lendKeeper.SetAssetRatesParams(ctx, stATOMRatesParams)

	auctionParamsLend := types.AuctionParams{
		AppId:                  3,
		AuctionDurationSeconds: 18000,
		Buffer:                 Dec("1.15"),
		Cusp:                   Dec("0.7"),
		Step:                   sdkmath.NewInt(360),
		PriceFunctionType:      1,
		DutchId:                3,
		BidDurationSeconds:     3600,
	}
	err = lendKeeper.AddAuctionParamsData(ctx, auctionParamsLend)
	if err != nil {
		return
	}

	auctionParams := bindings.MsgAddAuctionParams{
		AppID:                  2,
		AuctionDurationSeconds: 18000,
		Buffer:                 Dec("1.15"),
		Cusp:                   Dec("0.70"),
		Step:                   360,
		PriceFunctionType:      1,
		SurplusID:              1,
		DebtID:                 2,
		DutchID:                3,
		BidDurationSeconds:     3600,
	}
	err = auctionKeeper.AddAuctionParams(ctx, &auctionParams)
	if err != nil {
		return
	}
}

func CreateUpgradeHandler810(
	mm *module.Manager,
	configurator module.Configurator,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		newVM, err := mm.RunMigrations(ctx, configurator, fromVM)
		return newVM, err
	}
}

func CreateUpgradeHandler811(
	mm *module.Manager,
	configurator module.Configurator,
	assetKeeper assetkeeper.Keeper,
	lendKeeper lendkeeper.Keeper,
	auctionKeeper auctionkeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		newVM, err := mm.RunMigrations(ctx, configurator, fromVM)
		UpdateAuctionParams(ctx, assetKeeper, lendKeeper, auctionKeeper)
		UpdateExtendedPairVaultsAndAsset(ctx, assetKeeper)
		return newVM, err
	}
}
