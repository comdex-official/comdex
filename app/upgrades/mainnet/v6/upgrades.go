package v6

import (
	"fmt"

	assetkeeper "github.com/comdex-official/comdex/x/asset/keeper"
	assettypes "github.com/comdex-official/comdex/x/asset/types"
	lendkeeper "github.com/comdex-official/comdex/x/lend/keeper"
	"github.com/comdex-official/comdex/x/lend/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

func Dec(s string) sdk.Dec {
	dec, err := sdk.NewDecFromStr(s)
	if err != nil {
		panic(err)
	}
	return dec
}

func InitializeLendStates(
	ctx sdk.Context,
	assetKeeper assetkeeper.Keeper,
	lendKeeper lendkeeper.Keeper,
) {
	// Add Commodo App
	// Add Asset Rates for CMST, ATOM, CMDX
	// Add Lend Pool
	// Add Lend Pair
	// Add Lend Asset Pair

	// Adding Commodo App
	app := assettypes.AppData{Name: "commodo", ShortName: "cmdo", MinGovDeposit: sdk.ZeroInt(), GovTimeInSeconds: 0, GenesisToken: []assettypes.MintGenesisToken{}}
	err := assetKeeper.AddAppRecords(ctx, app)
	if err != nil {
		panic(err)
	}

	// Adding Asset Rates
	cmstRatesParams := types.AssetRatesParams{
		AssetID:              3,
		UOptimal:             Dec("0.8"),
		Base:                 Dec("0.002"),
		Slope1:               Dec("0.06"),
		Slope2:               Dec("0.6"),
		EnableStableBorrow:   false,
		StableBase:           Dec("0.0"),
		StableSlope1:         Dec("0.0"),
		StableSlope2:         Dec("0.0"),
		Ltv:                  Dec("0.8"),
		LiquidationThreshold: Dec("0.85"),
		LiquidationPenalty:   Dec("0.025"),
		LiquidationBonus:     Dec("0.025"),
		ReserveFactor:        Dec("0.1"),
		CAssetID:             7,
	}
	lendKeeper.SetAssetRatesParams(ctx, cmstRatesParams)
	atomRatesParams := types.AssetRatesParams{
		AssetID:              1,
		UOptimal:             Dec("0.75"),
		Base:                 Dec("0.002"),
		Slope1:               Dec("0.07"),
		Slope2:               Dec("1.25"),
		EnableStableBorrow:   false,
		StableBase:           Dec("0.0"),
		StableSlope1:         Dec("0.0"),
		StableSlope2:         Dec("0.0"),
		Ltv:                  Dec("0.7"),
		LiquidationThreshold: Dec("0.75"),
		LiquidationPenalty:   Dec("0.05"),
		LiquidationBonus:     Dec("0.05"),
		ReserveFactor:        Dec("0.2"),
		CAssetID:             5,
	}
	lendKeeper.SetAssetRatesParams(ctx, atomRatesParams)

	cmdxRatesParams := types.AssetRatesParams{
		AssetID:              2,
		UOptimal:             Dec("0.5"),
		Base:                 Dec("0.002"),
		Slope1:               Dec("0.08"),
		Slope2:               Dec("2.0"),
		EnableStableBorrow:   false,
		StableBase:           Dec("0.0"),
		StableSlope1:         Dec("0.0"),
		StableSlope2:         Dec("0.0"),
		Ltv:                  Dec("0.5"),
		LiquidationThreshold: Dec("0.55"),
		LiquidationPenalty:   Dec("0.05"),
		LiquidationBonus:     Dec("0.05"),
		ReserveFactor:        Dec("0.2"),
		CAssetID:             6,
	}
	lendKeeper.SetAssetRatesParams(ctx, cmdxRatesParams)

	// Adding Lend Pool
	var (
		assetDataCMDXPool []*types.AssetDataPoolMapping
	)
	assetDataPoolOneAssetOne := &types.AssetDataPoolMapping{
		AssetID:          1,
		AssetTransitType: 3,
		SupplyCap:        sdk.NewDec(5000000000000),
	}
	assetDataPoolOneAssetTwo := &types.AssetDataPoolMapping{
		AssetID:          2,
		AssetTransitType: 1,
		SupplyCap:        sdk.NewDec(1000000000000),
	}
	assetDataPoolOneAssetThree := &types.AssetDataPoolMapping{
		AssetID:          3,
		AssetTransitType: 2,
		SupplyCap:        sdk.NewDec(5000000000000),
	}

	assetDataCMDXPool = append(assetDataCMDXPool, assetDataPoolOneAssetOne, assetDataPoolOneAssetTwo, assetDataPoolOneAssetThree)
	cmdxPool := types.Pool{
		ModuleName: "cmdx",
		CPoolName:  "CMDX-ATOM-CMST",
		AssetData:  assetDataCMDXPool,
	}
	err = lendKeeper.AddPoolRecords(ctx, cmdxPool)
	if err != nil {
		panic(err)
	}

	// Adding Lend Pair
	cmdxcmstPair := types.Extended_Pair{ // 1
		AssetIn:         2,
		AssetOut:        3,
		IsInterPool:     false,
		AssetOutPoolID:  1,
		MinUsdValueLeft: 100000,
	}
	err = lendKeeper.AddLendPairsRecords(ctx, cmdxcmstPair)
	if err != nil {
		panic(err)
	}
	cmdxatomPair := types.Extended_Pair{ // 2
		AssetIn:         2,
		AssetOut:        1,
		IsInterPool:     false,
		AssetOutPoolID:  1,
		MinUsdValueLeft: 100000,
	}
	err = lendKeeper.AddLendPairsRecords(ctx, cmdxatomPair)
	if err != nil {
		panic(err)
	}
	atomcmdxPair := types.Extended_Pair{ // 3
		AssetIn:         1,
		AssetOut:        2,
		IsInterPool:     false,
		AssetOutPoolID:  1,
		MinUsdValueLeft: 100000,
	}
	err = lendKeeper.AddLendPairsRecords(ctx, atomcmdxPair)
	if err != nil {
		panic(err)
	}
	atomcmstPair := types.Extended_Pair{ // 4
		AssetIn:         1,
		AssetOut:        3,
		IsInterPool:     false,
		AssetOutPoolID:  1,
		MinUsdValueLeft: 100000,
	}
	err = lendKeeper.AddLendPairsRecords(ctx, atomcmstPair)
	if err != nil {
		panic(err)
	}
	cmstcmdxPair := types.Extended_Pair{ // 5
		AssetIn:         3,
		AssetOut:        2,
		IsInterPool:     false,
		AssetOutPoolID:  1,
		MinUsdValueLeft: 100000,
	}
	err = lendKeeper.AddLendPairsRecords(ctx, cmstcmdxPair)
	if err != nil {
		panic(err)
	}
	cmstatomPair := types.Extended_Pair{ // 6
		AssetIn:         3,
		AssetOut:        1,
		IsInterPool:     false,
		AssetOutPoolID:  1,
		MinUsdValueLeft: 100000,
	}
	err = lendKeeper.AddLendPairsRecords(ctx, cmstatomPair)
	if err != nil {
		panic(err)
	}

	// Adding Lend Pair Mapping
	map1 := types.AssetToPairMapping{
		PoolID:  1,
		AssetID: 1,
		PairID:  []uint64{3, 4},
	}
	lendKeeper.SetAssetToPair(ctx, map1)
	map2 := types.AssetToPairMapping{
		PoolID:  1,
		AssetID: 2,
		PairID:  []uint64{1, 2},
	}
	lendKeeper.SetAssetToPair(ctx, map2)
	map3 := types.AssetToPairMapping{
		PoolID:  1,
		AssetID: 3,
		PairID:  []uint64{5, 6},
	}
	lendKeeper.SetAssetToPair(ctx, map3)

	auctionParams := types.AuctionParams{
		AppId:                  3,
		AuctionDurationSeconds: 21600,
		Buffer:                 Dec("1.2"),
		Cusp:                   Dec("0.7"),
		Step:                   sdk.NewInt(360),
		PriceFunctionType:      1,
		DutchId:                3,
		BidDurationSeconds:     3600,
	}
	err = lendKeeper.AddAuctionParamsData(ctx, auctionParams)
	if err != nil {
		return
	}
}

// CreateUpgradeHandler creates an SDK upgrade handler for v5
func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	slashingkeeper slashingkeeper.Keeper,
	mintkeeper mintkeeper.Keeper,
	bankkeeper bankkeeper.Keeper,
	stakingkeeper stakingkeeper.Keeper,
	assetKeeper assetkeeper.Keeper,
	lendKeeper lendkeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("Running revert of tombstoning")

		err := RevertCosTombstoning(
			ctx,
			slashingkeeper,
			mintkeeper,
			bankkeeper,
			stakingkeeper,
		)
		if err != nil {
			panic(fmt.Sprintf("failed to revert tombstoning: %s", err))
		}

		ctx.Logger().Info("Running module migrations for v6.0.0...")
		newVM, err := mm.RunMigrations(ctx, configurator, fromVM)
		InitializeLendStates(ctx, assetKeeper, lendKeeper)
		return newVM, err
	}
}

func InitializeLendReservesStates(
	ctx sdk.Context,
	lendKeeper lendkeeper.Keeper,
) {
	dataAsset1, _ := lendKeeper.GetReserveBuybackAssetData(ctx, 1)
	dataAsset2, _ := lendKeeper.GetReserveBuybackAssetData(ctx, 2)
	dataAsset3, _ := lendKeeper.GetReserveBuybackAssetData(ctx, 3)
	reserveStat1 := types.AllReserveStats{
		AssetID:                        1,
		AmountOutFromReserveToLenders:  sdk.ZeroInt(),
		AmountOutFromReserveForAuction: sdk.ZeroInt(),
		AmountInFromLiqPenalty:         sdk.ZeroInt(),
		AmountInFromRepayments:         dataAsset1.BuybackAmount.Add(dataAsset1.ReserveAmount),
		TotalAmountOutToLenders:        sdk.ZeroInt(),
	}

	reserveStat2 := types.AllReserveStats{
		AssetID:                        2,
		AmountOutFromReserveToLenders:  sdk.ZeroInt(),
		AmountOutFromReserveForAuction: sdk.ZeroInt(),
		AmountInFromLiqPenalty:         sdk.ZeroInt(),
		AmountInFromRepayments:         dataAsset2.BuybackAmount.Add(dataAsset2.ReserveAmount),
		TotalAmountOutToLenders:        sdk.ZeroInt(),
	}

	reserveStat3 := types.AllReserveStats{
		AssetID:                        3,
		AmountOutFromReserveToLenders:  sdk.ZeroInt(),
		AmountOutFromReserveForAuction: sdk.ZeroInt(),
		AmountInFromLiqPenalty:         sdk.ZeroInt(),
		AmountInFromRepayments:         dataAsset3.BuybackAmount.Add(dataAsset3.ReserveAmount),
		TotalAmountOutToLenders:        sdk.ZeroInt(),
	}
	lendKeeper.SetAllReserveStatsByAssetID(ctx, reserveStat1)
	lendKeeper.SetAllReserveStatsByAssetID(ctx, reserveStat2)
	lendKeeper.SetAllReserveStatsByAssetID(ctx, reserveStat3)
}

func CreateUpgradeHandler610(
	mm *module.Manager,
	configurator module.Configurator,
	lendKeeper lendkeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {

		newVM, err := mm.RunMigrations(ctx, configurator, fromVM)
		InitializeLendReservesStates(ctx, lendKeeper)
		return newVM, err
	}
}
