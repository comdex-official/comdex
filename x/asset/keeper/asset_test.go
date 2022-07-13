package keeper_test

import (
	"github.com/comdex-official/comdex/app/wasm/bindings"
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
			false,
		},
		{
			"Add Duplicate App name cswap werd",
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
				},
			},
			true,
		},
		{
			"Add Duplicate short name werd cswap",
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
				},
			},
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
			false,
		},
	} {
		s.Run(tc.name, func() {
			err := assetKeeper.AddAppRecords(*ctx, tc.msg)
			if tc.isErrorExpected {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
			}
		})
	}
}

func (s *KeeperTestSuite) TestAddAssetRecords() {

	assetKeeper, ctx := &s.assetKeeper, &s.ctx
	for _, tc := range []struct {
		name            string
		msg             assetTypes.Asset
		isErrorExpected bool
	}{
		{"Add Asset cmdx ucmdx",
			assetTypes.Asset{Name: "CMDX",
				Denom:     "ucmdx",
				Decimals:  1000000,
				IsOnChain: true},
			false,
		},
		{"Add Asset:  Duplicate Asset Name  2 cmdx uosmo",
			assetTypes.Asset{Name: "CMDX",
				Denom:     "uosmo",
				Decimals:  1000000,
				IsOnChain: true},
			true,
		},
		{"Add Asset : Duplicate Denom 2 osmo ucmdx",
			assetTypes.Asset{Name: "OSMO",
				Denom:     "ucmdx",
				Decimals:  1000000,
				IsOnChain: true},
			true,
		},
		{"Add Asset 2",
			assetTypes.Asset{Name: "CMST",
				Denom:     "ucmst",
				Decimals:  1000000,
				IsOnChain: true},
			false,
		},
		{"Add Asset 3",
			assetTypes.Asset{Name: "HARBOR",
				Denom:     "uharbor",
				Decimals:  1000000,
				IsOnChain: true},
			false,
		},
	} {
		s.Run(tc.name, func() {
			err := assetKeeper.AddAssetRecords(*ctx, tc.msg)
			if tc.isErrorExpected {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
			}

		})
	}
}

func (s *KeeperTestSuite) TestAddPairAndExtendedPairVault1() {

	assetKeeper, ctx := &s.assetKeeper, &s.ctx
	s.TestAddApp()
	s.TestAddAssetRecords()
	for _, tc := range []struct {
		name                    string
		pair                    assetTypes.Pair
		extendedPairVault       bindings.MsgAddExtendedPairsVault
		symbol1                 string
		symbol2                 string
		isErrorExpectedForPair  bool
		isErrorExpectedForVault bool
	}{
		{"Add Pair , Extended Pair Vault : cmdx cmst",
			assetTypes.Pair{
				AssetIn:  1,
				AssetOut: 2,
			},
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
			false,
		},
	} {
		s.Run(tc.name, func() {
			err := assetKeeper.AddPairsRecords(*ctx, tc.pair)
			if tc.isErrorExpectedForPair {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
			}
			err = assetKeeper.WasmAddExtendedPairsVaultRecords(*ctx, &tc.extendedPairVault)
			if tc.isErrorExpectedForVault {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
			}

		})
	}
}
