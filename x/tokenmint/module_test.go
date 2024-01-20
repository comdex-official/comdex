package tokenmint_test

import (
	sdkmath "cosmossdk.io/math"
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"

	chain "github.com/comdex-official/comdex/app"
	assettypes "github.com/comdex-official/comdex/x/asset/types"
	markettypes "github.com/comdex-official/comdex/x/market/types"
	"github.com/comdex-official/comdex/x/tokenmint/keeper"
)

type ModuleTestSuite struct {
	suite.Suite

	app     *chain.App
	ctx     sdk.Context
	keeper  keeper.Keeper
	querier keeper.QueryServer
	addrs   []sdk.AccAddress
}

func TestModuleTestSuite(t *testing.T) {
	suite.Run(t, new(ModuleTestSuite))
}

func (suite *ModuleTestSuite) SetupTest() {
	app := chain.Setup(suite.T(), false)
	ctx := app.BaseApp.NewContext(false)

	suite.app = app
	suite.ctx = ctx
	suite.keeper = suite.app.TokenmintKeeper
	suite.querier = keeper.QueryServer{Keeper: suite.keeper}
}

// Below are useful helpers to write test code easily.
func (suite *ModuleTestSuite) addr(addrNum int) sdk.AccAddress {
	addr := make(sdk.AccAddress, 20)
	binary.PutVarint(addr, int64(addrNum))
	return addr
}

func (s *ModuleTestSuite) CreateNewApp(appName string) uint64 {
	err := s.app.AssetKeeper.AddAppRecords(s.ctx, assettypes.AppData{
		Name:             appName,
		ShortName:        appName,
		MinGovDeposit:    sdkmath.NewInt(0),
		GovTimeInSeconds: 0,
		GenesisToken:     []assettypes.MintGenesisToken{},
	})
	s.Require().NoError(err)
	found := s.app.AssetKeeper.HasAppForName(s.ctx, appName)
	s.Require().True(found)

	apps, found := s.app.AssetKeeper.GetApps(s.ctx)
	s.Require().True(found)
	var appID uint64
	for _, app := range apps {
		if app.Name == appName {
			appID = app.Id
			break
		}
	}
	s.Require().NotZero(appID)
	return appID
}

func (s *ModuleTestSuite) CreateNewAsset(name, denom string, price uint64) uint64 {
	err := s.app.AssetKeeper.AddAssetRecords(s.ctx, assettypes.Asset{
		Name:                  name,
		Denom:                 denom,
		Decimals:              sdkmath.NewInt(1000000),
		IsOnChain:             true,
		IsOraclePriceRequired: true,
		IsCdpMintable:         true,
	})
	s.Require().NoError(err)
	assets := s.app.AssetKeeper.GetAssets(s.ctx)
	var assetID uint64
	for _, asset := range assets {
		if asset.Denom == denom {
			assetID = asset.Id
			break
		}
	}
	s.Require().NotZero(assetID)

	twa1 := markettypes.TimeWeightedAverage{
		AssetID:       1,
		ScriptID:      10,
		Twa:           12000000,
		CurrentIndex:  1,
		IsPriceActive: true,
		PriceValue:    nil,
	}
	twa2 := markettypes.TimeWeightedAverage{
		AssetID:       2,
		ScriptID:      10,
		Twa:           100000,
		CurrentIndex:  1,
		IsPriceActive: true,
		PriceValue:    nil,
	}
	twa3 := markettypes.TimeWeightedAverage{
		AssetID:       3,
		ScriptID:      10,
		Twa:           1000000,
		CurrentIndex:  1,
		IsPriceActive: true,
		PriceValue:    nil,
	}
	twa4 := markettypes.TimeWeightedAverage{
		AssetID:       4,
		ScriptID:      10,
		Twa:           2500000,
		CurrentIndex:  1,
		IsPriceActive: true,
		PriceValue:    nil,
	}
	s.app.MarketKeeper.SetTwa(s.ctx, twa1)
	s.app.MarketKeeper.SetTwa(s.ctx, twa2)
	s.app.MarketKeeper.SetTwa(s.ctx, twa3)
	s.app.MarketKeeper.SetTwa(s.ctx, twa4)

	return assetID
}
