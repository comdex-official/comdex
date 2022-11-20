package v5

import (
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/petrichormoney/petri/app/wasm/bindings"
	assetkeeper "github.com/petrichormoney/petri/x/asset/keeper"
	assettypes "github.com/petrichormoney/petri/x/asset/types"
	auctionkeeper "github.com/petrichormoney/petri/x/auction/keeper"
	collectorkeeper "github.com/petrichormoney/petri/x/collector/keeper"
	liquidationkeeper "github.com/petrichormoney/petri/x/liquidation/keeper"
	liquiditykeeper "github.com/petrichormoney/petri/x/liquidity/keeper"
	liquiditytypes "github.com/petrichormoney/petri/x/liquidity/types"
	lockerkeeper "github.com/petrichormoney/petri/x/locker/keeper"
	lockertypes "github.com/petrichormoney/petri/x/locker/types"
	rewardskeeper "github.com/petrichormoney/petri/x/rewards/keeper"
	rewardstypes "github.com/petrichormoney/petri/x/rewards/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	ica "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts"
	icacontrollertypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/controller/types"
	icahosttypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/host/types"
	icatypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
)

func InitializeStates(
	ctx sdk.Context,
	assetKeeper assetkeeper.Keeper,
	liquidityKeeper liquiditykeeper.Keeper,
	collectorKeeper collectorkeeper.Keeper,
	auctionKeeper auctionkeeper.Keeper,
	lockerKeeper lockerkeeper.Keeper,
	rewardsKeeper rewardskeeper.Keeper,
	liquidationKeeper liquidationkeeper.Keeper,
) {
	genesisToken := assettypes.MintGenesisToken{
		AssetId:       9,
		GenesisSupply: sdk.NewIntFromUint64(500000000000000),
		IsGovToken:    true,
		Recipient:     "petri1tadhnvwa0sqzwr3m60f7dsjw4ua77qsz3ptcyw",
	}
	var gToken []assettypes.MintGenesisToken
	gToken = append(gToken, genesisToken)

	apps := []assettypes.AppData{
		{Name: "atom", ShortName: "uatom", MinGovDeposit: sdk.ZeroInt(), GovTimeInSeconds: 0, GenesisToken: []assettypes.MintGenesisToken{}},
		{Name: "harbor", ShortName: "hbr", MinGovDeposit: sdk.NewInt(10000000000), GovTimeInSeconds: 259200, GenesisToken: gToken},
	}
	for _, app := range apps {
		err := assetKeeper.AddAppRecords(ctx, app)
		if err != nil {
			panic(err)
		}
	}
	assetKeeper.SetGenesisTokenForApp(ctx, 2, 9)

	assets := []assettypes.Asset{
		{Name: "ATOM", Denom: "ibc/961FA3E54F5DCCA639F37A7C45F7BBE41815579EF1513B5AFBEFCFEB8F256352", Decimals: sdk.NewInt(1000000), IsOnChain: false, IsOraclePriceRequired: true, IsCdpMintable: false},
		{Name: "FURY", Denom: "ufury", Decimals: sdk.NewInt(1000000), IsOnChain: false, IsOraclePriceRequired: true, IsCdpMintable: false},
		{Name: "PETRI", Denom: "upetri", Decimals: sdk.NewInt(1000000000000), IsOnChain: false, IsOraclePriceRequired: true, IsCdpMintable: false},
		{Name: "FUST", Denom: "ufust", Decimals: sdk.NewInt(1000000), IsOnChain: true, IsOraclePriceRequired: true, IsCdpMintable: true},
		{Name: "OSMO", Denom: "ibc/0471F1C4E7AFD3F07702BEF6DC365268D64570F7C1FDC98EA6098DD6DE59817B", Decimals: sdk.NewInt(1000000), IsOnChain: false, IsOraclePriceRequired: true, IsCdpMintable: false},
		{Name: "sATOM", Denom: "usatom", Decimals: sdk.NewInt(1000000), IsOnChain: true, IsOraclePriceRequired: false, IsCdpMintable: true},
		{Name: "sFURY", Denom: "usfury", Decimals: sdk.NewInt(1000000), IsOnChain: true, IsOraclePriceRequired: false, IsCdpMintable: true},
		{Name: "sFUST", Denom: "uscmst", Decimals: sdk.NewInt(1000000), IsOnChain: true, IsOraclePriceRequired: false, IsCdpMintable: true},
		{Name: "sOSMO", Denom: "usosmo", Decimals: sdk.NewInt(1000000), IsOnChain: true, IsOraclePriceRequired: false, IsCdpMintable: true},
		{Name: "sPETRI", Denom: "uspetri", Decimals: sdk.NewInt(1000000), IsOnChain: true, IsOraclePriceRequired: false, IsCdpMintable: true},
		{Name: "HARBOR", Denom: "uharbor", Decimals: sdk.NewInt(1000000), IsOnChain: true, IsOraclePriceRequired: false, IsCdpMintable: false},
		{Name: "AXLUSDC", Denom: "ibc/E1616E7C19EA474C565737709A628D6F8A23FF9D3E9A7A6871306CF5E0A5341E", Decimals: sdk.NewInt(1000000), IsOnChain: false, IsOraclePriceRequired: true, IsCdpMintable: false},
	}

	for _, asset := range assets {
		err := assetKeeper.AddAssetRecords(ctx, asset)
		if err != nil {
			panic(err)
		}
	}
	// add pairs
	pairs := []assettypes.Pair{
		{AssetIn: 1, AssetOut: 3},
		{AssetIn: 2, AssetOut: 3},
		{AssetIn: 4, AssetOut: 3},
		{AssetIn: 10, AssetOut: 3},
	}

	for _, pair := range pairs {
		err := assetKeeper.AddPairsRecords(ctx, pair)
		if err != nil {
			panic(err)
		}
	}

	// add extended pairs
	extPairs := []*bindings.MsgAddExtendedPairsVault{
		{
			AppID: 2, PairID: 4, StabilityFee: sdk.ZeroDec(), ClosingFee: sdk.ZeroDec(), LiquidationPenalty: sdk.ZeroDec(),
			DrawDownFee: sdk.MustNewDecFromStr("0.001"), IsVaultActive: true, DebtCeiling: sdk.NewInt(40000000000000), DebtFloor: sdk.NewInt(1000000), IsStableMintVault: true, MinCr: sdk.MustNewDecFromStr("1"),
			PairName: "AXL-USDC-FUST", AssetOutOraclePrice: false, AssetOutPrice: 1000000, MinUsdValueLeft: 10000000000,
		},
	}

	for _, extPair := range extPairs {
		err := assetKeeper.WasmAddExtendedPairsVaultRecords(ctx, extPair)
		if err != nil {
			panic(err)
		}
	}
	// add collector params
	collector := bindings.MsgSetCollectorLookupTable{
		AppID:            2,
		CollectorAssetID: 3,
		SecondaryAssetID: 9,
		SurplusThreshold: sdk.NewInt(50000000000),
		DebtThreshold:    sdk.NewInt(0),
		LockerSavingRate: sdk.MustNewDecFromStr("0.00"),
		LotSize:          sdk.NewInt(10000000000),
		BidFactor:        sdk.MustNewDecFromStr("0.01"),
		DebtLotSize:      sdk.NewInt(1000000000000),
	}

	err := collectorKeeper.WasmSetCollectorLookupTable(ctx, &collector)
	if err != nil {
		panic(err)
	}

	// add auction params

	auctionParam := bindings.MsgAddAuctionParams{
		AppID:                  2,
		AuctionDurationSeconds: 21600,
		BidDurationSeconds:     3600,
		Buffer:                 sdk.MustNewDecFromStr("1.2"),
		Cusp:                   sdk.MustNewDecFromStr("0.75"),
		DebtID:                 2,
		DutchID:                3,
		PriceFunctionType:      1,
		Step:                   360,
		SurplusID:              1,
	}
	err = auctionKeeper.AddAuctionParams(ctx, &auctionParam)
	if err != nil {
		panic(err)
	}
	// add auction mapping
	auction := bindings.MsgSetAuctionMappingForApp{
		AppID:                2,
		AssetIDs:             3,
		IsSurplusAuctions:    false,
		IsDebtAuctions:       false,
		IsDistributor:        false,
		AssetOutOraclePrices: false,
		AssetOutPrices:       1000000,
	}

	err = collectorKeeper.WasmSetAuctionMappingForApp(ctx, &auction)
	if err != nil {
		panic(err)
	}

	// whitlist cmst for locker
	locker := lockertypes.MsgAddWhiteListedAssetRequest{
		AppId:   2,
		AssetId: 3,
	}
	_, err = lockerKeeper.AddWhiteListedAsset(ctx, &locker)
	if err != nil {
		panic(err)
	}
	// whielist for locker rewards
	reward := rewardstypes.WhitelistAsset{
		AppMappingId: 2,
		AssetId:      3,
	}
	_, err = rewardsKeeper.Whitelist(ctx, &reward)
	if err != nil {
		panic(err)
	}
	// whitlist for vaultInterest
	vInterest := rewardstypes.WhitelistAppIdVault{
		AppMappingId: 2,
	}
	_, err = rewardsKeeper.WhitelistAppVault(ctx, &vInterest)
	if err != nil {
		panic(err)
	}

	// whitlist for liquidation
	err = liquidationKeeper.WasmWhitelistAppIDLiquidation(ctx, 2)
	if err != nil {
		panic(err)
	}

	type LiquidityPair struct {
		AppID          uint64
		From           string
		BaseCoinDenom  string
		QuoteCoinDenom string
	}

	liquidityPairs := []LiquidityPair{
		{AppID: 1, From: "petri12gfx7e3p08ljrwhq4lxz0360czcv9jpzajlytv", BaseCoinDenom: "upetri", QuoteCoinDenom: "ibc/961FA3E54F5DCCA639F37A7C45F7BBE41815579EF1513B5AFBEFCFEB8F256352"},
		{AppID: 1, From: "petri12gfx7e3p08ljrwhq4lxz0360czcv9jpzajlytv", BaseCoinDenom: "upetri", QuoteCoinDenom: "ibc/0471F1C4E7AFD3F07702BEF6DC365268D64570F7C1FDC98EA6098DD6DE59817B"},
		{AppID: 1, From: "petri12gfx7e3p08ljrwhq4lxz0360czcv9jpzajlytv", BaseCoinDenom: "upetri", QuoteCoinDenom: "ufust"},
		{AppID: 1, From: "petri12gfx7e3p08ljrwhq4lxz0360czcv9jpzajlytv", BaseCoinDenom: "upetri", QuoteCoinDenom: "uharbor"},
	}

	for _, lpair := range liquidityPairs {
		msg := liquiditytypes.NewMsgCreatePair(
			lpair.AppID, sdk.MustAccAddressFromBech32(lpair.From), lpair.BaseCoinDenom, lpair.QuoteCoinDenom,
		)
		_, err := liquidityKeeper.CreatePair(ctx, msg, true)
		if err != nil {
			panic(err)
		}
	}
}

// CreateUpgradeHandler creates an SDK upgrade handler for v5
func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	wasmKeeper wasmkeeper.Keeper,
	assetKeeper assetkeeper.Keeper,
	liquidityKeeper liquiditykeeper.Keeper,
	collectorKeeper collectorkeeper.Keeper,
	auctionKeeper auctionkeeper.Keeper,
	lockerKeeper lockerkeeper.Keeper,
	rewardsKeeper rewardskeeper.Keeper,
	liquidationKeeper liquidationkeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		// Refs:
		// - https://docs.cosmos.network/master/building-modules/upgrade.html#registering-migrations
		// - https://docs.cosmos.network/master/migrations/chain-upgrade-guide-044.html#chain-upgrade

		// Deleting these modules from the upgrades current state
		// Add Interchain Accounts host module
		// set the ICS27 consensus version so InitGenesis is not run
		fromVM[icatypes.ModuleName] = mm.Modules[icatypes.ModuleName].ConsensusVersion()

		// create ICS27 Controller submodule params, controller module not enabled.
		controllerParams := icacontrollertypes.Params{}

		// create ICS27 Host submodule params

		// create ICS27 Host submodule params
		hostParams := icahosttypes.Params{
			HostEnabled: true,
			AllowMessages: []string{
				sdk.MsgTypeURL(&ibctransfertypes.MsgTransfer{}),
				sdk.MsgTypeURL(&banktypes.MsgSend{}),
				sdk.MsgTypeURL(&stakingtypes.MsgDelegate{}),
				sdk.MsgTypeURL(&stakingtypes.MsgBeginRedelegate{}),
				sdk.MsgTypeURL(&stakingtypes.MsgCreateValidator{}),
				sdk.MsgTypeURL(&stakingtypes.MsgEditValidator{}),
				sdk.MsgTypeURL(&stakingtypes.MsgUndelegate{}),
				sdk.MsgTypeURL(&distrtypes.MsgWithdrawDelegatorReward{}),
				sdk.MsgTypeURL(&distrtypes.MsgSetWithdrawAddress{}),
				sdk.MsgTypeURL(&distrtypes.MsgWithdrawValidatorCommission{}),
				sdk.MsgTypeURL(&distrtypes.MsgFundCommunityPool{}),
				sdk.MsgTypeURL(&govtypes.MsgVote{}),
			},
		}
		// No changes in existing module and their states,
		// This upgrades adds new modules and new states in the existing store

		icamodule, correctTypecast := mm.Modules[icatypes.ModuleName].(ica.AppModule)
		if !correctTypecast {
			panic("mm.Modules[icatypes.ModuleName] is not of type ica.AppModule")
		}
		icamodule.InitModule(ctx, controllerParams, hostParams)
		newVM, err := mm.RunMigrations(ctx, configurator, fromVM)
		if err != nil {
			return newVM, err
		}

		InitializeStates(ctx, assetKeeper, liquidityKeeper, collectorKeeper, auctionKeeper, lockerKeeper, rewardsKeeper, liquidationKeeper)

		// update wasm to permission
		wasmParams := wasmKeeper.GetParams(ctx)
		wasmParams.CodeUploadAccess = wasmtypes.AllowNobody
		wasmKeeper.SetParams(ctx, wasmParams)

		return newVM, err
	}
}
