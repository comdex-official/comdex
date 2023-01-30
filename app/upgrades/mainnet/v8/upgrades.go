package v8

import (
	"github.com/comdex-official/comdex/app/wasm/bindings"
	assetkeeper "github.com/comdex-official/comdex/x/asset/keeper"
	assettypes "github.com/comdex-official/comdex/x/asset/types"
	auctionkeeper "github.com/comdex-official/comdex/x/auction/keeper"
	lendkeeper "github.com/comdex-official/comdex/x/lend/keeper"
	"github.com/comdex-official/comdex/x/lend/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

func Dec(s string) sdk.Dec {
	dec, err := sdk.NewDecFromStr(s)
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
		Decimals:              sdk.NewInt(1000000),
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
		Decimals:              sdk.NewInt(1000000),
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
		CAssetID:             21,
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
		CAssetID:             22,
	}
	lendKeeper.SetAssetRatesParams(ctx, stATOMRatesParams)

	auctionParamsLend := types.AuctionParams{
		AppId:                  3,
		AuctionDurationSeconds: 18000,
		Buffer:                 Dec("1.15"),
		Cusp:                   Dec("0.7"),
		Step:                   sdk.NewInt(360),
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
		Cusp:                   Dec("0.75"),
		Step:                   360,
		PriceFunctionType:      1,
		SurplusID:              1,
		DebtID:                 2,
		DutchID:                3,
		BidDurationSeconds:     3600,
	}
	err = auctionKeeper.AddAuctionParams(ctx, &auctionParams)
}

func CreateUpgradeHandler800(
	mm *module.Manager,
	configurator module.Configurator,
	assetKeeper assetkeeper.Keeper,
	lendKeeper lendkeeper.Keeper,
	auctionKeeper auctionkeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		newVM, err := mm.RunMigrations(ctx, configurator, fromVM)
		UpdateAuctionParams(ctx, assetKeeper, lendKeeper, auctionKeeper)
		return newVM, err
	}
}
