package keeper_test

import (
	"github.com/comdex-official/comdex/app/wasm/bindings"
	"github.com/comdex-official/comdex/x/asset/keeper"
	assetTypes "github.com/comdex-official/comdex/x/asset/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

/*
func (k *Keeper) AddAppMappingRecords(ctx sdk.Context, records ...types.AppMapping) error
func (k *Keeper) AddAssetRecords(ctx sdk.Context, records ...types.Asset) error
func (k *Keeper) AddPairsRecords(ctx sdk.Context, records ...types.Pair) error
func (k *Keeper) AddExtendedPairsVaultRecords(ctx sdk.Context, records ...types.ExtendedPairVault) error
func (k Keeper) WhitelistAppId(ctx sdk.Context, appMappingId uint64) error
*/

func (s *KeeperTestSuite) TestAddApp() {
	userAddress1 := "cosmos1q7q90qsl9g0gl2zz0njxwv2a649yqrtyxtnv3v"

	genesisSupply := sdk.NewIntFromUint64(1000000)
	assetKeeper, ctx := &s.assetKeeper, &s.ctx

	for _, tc := range []struct {
		name            string
		msg             assetTypes.AppData
		appID           uint64
		isErrorExpected bool
	}{
		{
			"Add App cswap cswap",
			assetTypes.AppData{
				Name:             "cswap",
				ShortName:        "cswap",
				MinGovDeposit:    sdk.NewIntFromUint64(10000000),
				GovTimeInSeconds: 900,
				GenesisToken: []assetTypes.MintGenesisToken{
					{
						3,
						&genesisSupply,
						true,
						userAddress1,
					},
					{
						2,
						&genesisSupply,
						true,
						userAddress1,
					},
				},
			},
			1,
			false,
		},
		{
			"Add Duplicate App name cswap werd",
			assetTypes.AppData{
				Name:             "cswap",
				ShortName:        "werd",
				MinGovDeposit:    sdk.NewIntFromUint64(10000000),
				GovTimeInSeconds: 900,
				GenesisToken: []assetTypes.MintGenesisToken{
					{
						3,
						&genesisSupply,
						true,
						userAddress1,
					},
				},
			},
			2,
			true,
		},
		{
			"Add Duplicate short name werd cswap",
			assetTypes.AppData{
				Name:             "werd",
				ShortName:        "cswap",
				MinGovDeposit:    sdk.NewIntFromUint64(10000000),
				GovTimeInSeconds: 900,
				GenesisToken: []assetTypes.MintGenesisToken{
					{
						3,
						&genesisSupply,
						true,
						userAddress1,
					},
				},
			},
			2,
			true,
		},
		{
			"Add App commodo commodo",
			assetTypes.AppData{
				Name:             "commodo",
				ShortName:        "commodo",
				MinGovDeposit:    sdk.NewIntFromUint64(10000000),
				GovTimeInSeconds: 900,
				GenesisToken: []assetTypes.MintGenesisToken{
					{
						3,
						&genesisSupply,
						true,
						userAddress1,
					},
				},
			},
			2,
			false,
		},
	} {
		s.Run(tc.name, func() {
			err := assetKeeper.AddAppRecords(*ctx, tc.msg)
			if tc.isErrorExpected {
				s.Require().Error(err)
			} else {

				s.Require().NoError(err)
				server := keeper.NewQueryServer(*assetKeeper)
				res, err := server.QueryApp(sdk.WrapSDKContext(*ctx), &assetTypes.QueryAppRequest{Id: tc.appID})
				s.Require().NoError(err)
				s.Require().Equal(res.App.Id, tc.appID)
				s.Require().Equal(res.App.Name, tc.msg.Name)
				s.Require().Equal(res.App.ShortName, tc.msg.ShortName)
				s.Require().Equal(res.App.GovTimeInSeconds, tc.msg.GovTimeInSeconds)
				s.Require().Equal(res.App.MinGovDeposit, tc.msg.MinGovDeposit)
				s.Require().Equal(res.App.GenesisToken[0].AssetId, tc.msg.GenesisToken[0].AssetId)
				s.Require().Equal(res.App.GenesisToken[0].GenesisSupply, tc.msg.GenesisToken[0].GenesisSupply)
				s.Require().Equal(res.App.GenesisToken[0].Recipient, tc.msg.GenesisToken[0].Recipient)
				s.Require().Equal(res.App.GenesisToken[0].IsGovToken, tc.msg.GenesisToken[0].IsGovToken)
			}

		})
	}
}

func (s *KeeperTestSuite) TestQueryApps() {
	s.TestAddApp()
	assetKeeper, ctx := &s.assetKeeper, &s.ctx
	server := keeper.NewQueryServer(*assetKeeper)
	res, err := server.QueryApps(sdk.WrapSDKContext(*ctx), &assetTypes.QueryAppsRequest{})
	s.Require().NoError(err)
	s.Require().Equal(len(res.Apps), 2)
}
func (s *KeeperTestSuite) TestAddAssetRecords() {

	assetKeeper, ctx := &s.assetKeeper, &s.ctx
	for _, tc := range []struct {
		name            string
		msg             assetTypes.Asset
		assetID         uint64
		isErrorExpected bool
	}{
		{"Add Asset cmdx ucmdx",
			assetTypes.Asset{Name: "CMDX",
				Denom:     "ucmdx",
				Decimals:  1000000,
				IsOnChain: true},
			1,
			false,
		},
		//{"Add Asset:  Duplicate Asset Name  2 cmdx uosmo",
		//	assetTypes.Asset{Name: "CMDX",
		//		Denom:     "uosmo",
		//		Decimals:  1000000,
		//		IsOnChain: true},
		//	2,
		//	true,
		//},
		{"Add Asset : Duplicate Denom 2 osmo ucmdx",
			assetTypes.Asset{Name: "OSMO",
				Denom:     "ucmdx",
				Decimals:  1000000,
				IsOnChain: true},
			2,
			true,
		},
		{"Add Asset 2",
			assetTypes.Asset{Name: "CMST",
				Denom:     "ucmst",
				Decimals:  1000000,
				IsOnChain: true},
			2,
			false,
		},
		{"Add Asset 3",
			assetTypes.Asset{Name: "HARBOR",
				Denom:     "uharbor",
				Decimals:  1000000,
				IsOnChain: true},
			3,
			false,
		},
	} {
		s.Run(tc.name, func() {
			err := assetKeeper.AddAssetRecords(*ctx, tc.msg)
			if tc.isErrorExpected {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				server := keeper.NewQueryServer(*assetKeeper)
				res, err := server.QueryAsset(sdk.WrapSDKContext(*ctx), &assetTypes.QueryAssetRequest{Id: tc.assetID})
				s.Require().NoError(err)
				s.Require().Equal(res.Asset.Id, tc.assetID)
				s.Require().Equal(res.Asset.Denom, tc.msg.Denom)
				s.Require().Equal(res.Asset.IsOraclePriceRequired, tc.msg.IsOraclePriceRequired)
				s.Require().Equal(res.Asset.Decimals, tc.msg.Decimals)
				s.Require().Equal(res.Asset.Name, tc.msg.Name)
				s.Require().Equal(res.Asset.IsOnChain, tc.msg.IsOnChain)

			}

		})
	}
}

func (s *KeeperTestSuite) TestQueryAssets() {
	s.TestAddAssetRecords()
	assetKeeper, ctx := &s.assetKeeper, &s.ctx
	server := keeper.NewQueryServer(*assetKeeper)
	res, err := server.QueryAssets(sdk.WrapSDKContext(*ctx), &assetTypes.QueryAssetsRequest{})
	s.Require().NoError(err)
	s.Require().Equal(len(res.Assets), 3)
}

func (s *KeeperTestSuite) TestAddPair() {

	assetKeeper, ctx := &s.assetKeeper, &s.ctx
	s.TestAddApp()
	s.TestAddAssetRecords()
	for _, tc := range []struct {
		name                   string
		pair                   assetTypes.Pair
		symbol1                string
		symbol2                string
		isErrorExpectedForPair bool
		pairID                 uint64
	}{
		{"Add Pair 1: cmdx cmst",
			assetTypes.Pair{
				AssetIn:  1,
				AssetOut: 2,
			},
			"ucmdx",
			"ucmst",
			false,
			1,
		},
		{"Add Duplicate Pair : cmdx cmst",
			assetTypes.Pair{
				AssetIn:  1,
				AssetOut: 2,
			},
			"ucmdx",
			"ucmst",
			true,
			1,
		},
		{"Add Pair 2 : cmst cmdx",
			assetTypes.Pair{
				AssetIn:  2,
				AssetOut: 1,
			},
			"ucmst",
			"ucmdx",
			false,
			2,
		},
		{"Add Pair 3 : cmst harbor",
			assetTypes.Pair{
				AssetIn:  2,
				AssetOut: 3,
			},
			"ucmst",
			"uharbor",
			false,
			3,
		},
	} {
		s.Run(tc.name, func() {
			server := keeper.NewQueryServer(*assetKeeper)
			err := assetKeeper.AddPairsRecords(*ctx, tc.pair)
			if tc.isErrorExpectedForPair {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				res, err := server.QueryPair(sdk.WrapSDKContext(*ctx), &assetTypes.QueryPairRequest{Id: tc.pairID})
				s.Require().NoError(err)
				s.Require().Equal(res.PairInfo.Id, tc.pairID)
				s.Require().Equal(res.PairInfo.AssetIn, tc.pair.AssetIn)
				s.Require().Equal(res.PairInfo.AssetOut, tc.pair.AssetOut)
				s.Require().Equal(res.PairInfo.DenomIn, tc.symbol1)
				s.Require().Equal(res.PairInfo.DenomOut, tc.symbol2)
			}

		})
	}
}
func (s *KeeperTestSuite) TestAddPairAndExtendedPairVault() {

	assetKeeper, ctx := &s.assetKeeper, &s.ctx
	s.TestAddApp()
	s.TestAddAssetRecords()
	s.TestAddPair()
	for _, tc := range []struct {
		name                    string
		extendedPairVault       bindings.MsgAddExtendedPairsVault
		symbol1                 string
		symbol2                 string
		isErrorExpectedForVault bool
		vaultID                 uint64
	}{
		{"Add Pair , Extended Pair Vault : cmdx cmst",

			bindings.MsgAddExtendedPairsVault{
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
				MinUsdValueLeft:     1000000,
			},
			"ucmdx",
			"ucmst",
			false,
			1,
		},
		{"Add Duplicate Pair , Extended Pair Vault : cmdx cmst",

			bindings.MsgAddExtendedPairsVault{
				AppID:               1,
				PairID:              2,
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
				MinUsdValueLeft:     1000000,
			},
			"ucmst",
			"ucmdx",
			false,
			2,
		},
		{"Add Pair , Duplicate Extended Pair Vault : cmdx harbor",

			bindings.MsgAddExtendedPairsVault{
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
				MinUsdValueLeft:     1000000,
			},
			"ucmdx",
			"ucmst",
			true,
			3,
		},
		{"Add Pair , Extended Pair Vault : cmdx harbor",

			bindings.MsgAddExtendedPairsVault{
				AppID:               1,
				PairID:              3,
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
				MinUsdValueLeft:     1000000,
			},
			"ucmst",
			"uharbor",
			false,
			3,
		},
	} {
		s.Run(tc.name, func() {
			server := keeper.NewQueryServer(*assetKeeper)

			err := assetKeeper.WasmAddExtendedPairsVaultRecords(*ctx, &tc.extendedPairVault)

			if tc.isErrorExpectedForVault {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				res2, err := server.QueryExtendedPairVault(sdk.WrapSDKContext(*ctx), &assetTypes.QueryExtendedPairVaultRequest{Id: tc.vaultID})
				s.Require().NoError(err)
				s.Require().Equal(res2.PairVault.AppId, tc.extendedPairVault.AppID)
				s.Require().Equal(res2.PairVault.PairId, tc.extendedPairVault.PairID)
				s.Require().Equal(res2.PairVault.StabilityFee, tc.extendedPairVault.StabilityFee)
				s.Require().Equal(res2.PairVault.ClosingFee, tc.extendedPairVault.ClosingFee)
				s.Require().Equal(res2.PairVault.LiquidationPenalty, tc.extendedPairVault.LiquidationPenalty)
				s.Require().Equal(res2.PairVault.DrawDownFee, tc.extendedPairVault.DrawDownFee)
				s.Require().Equal(res2.PairVault.IsVaultActive, tc.extendedPairVault.IsVaultActive)
				s.Require().Equal(res2.PairVault.DebtCeiling.Uint64(), tc.extendedPairVault.DebtCeiling)
				s.Require().Equal(res2.PairVault.DebtFloor.Uint64(), tc.extendedPairVault.DebtFloor)
				s.Require().Equal(res2.PairVault.IsStableMintVault, tc.extendedPairVault.IsStableMintVault)
				s.Require().Equal(res2.PairVault.MinCr, tc.extendedPairVault.MinCr)
				s.Require().Equal(res2.PairVault.PairName, tc.extendedPairVault.PairName)
				s.Require().Equal(res2.PairVault.AssetOutOraclePrice, tc.extendedPairVault.AssetOutOraclePrice)
				s.Require().Equal(res2.PairVault.AssetOutPrice, tc.extendedPairVault.AssetOutPrice)
				s.Require().Equal(res2.PairVault.MinUsdValueLeft, tc.extendedPairVault.MinUsdValueLeft)
			}

		})
	}
}

func (s *KeeperTestSuite) TestQueryPairsAndExtendedPairVaults() {
	s.TestAddPairAndExtendedPairVault()
	assetKeeper, ctx := &s.assetKeeper, &s.ctx
	server := keeper.NewQueryServer(*assetKeeper)
	res, err := server.QueryPairs(sdk.WrapSDKContext(*ctx), &assetTypes.QueryPairsRequest{})
	s.Require().NoError(err)
	s.Require().Equal(len(res.PairsInfo), 1)
}
