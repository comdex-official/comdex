package vault_test

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	chain "github.com/comdex-official/comdex/app"
	"github.com/comdex-official/comdex/app/wasm/bindings"
	assettypes "github.com/comdex-official/comdex/x/asset/types"
	markettypes "github.com/comdex-official/comdex/x/market/types"
	"github.com/comdex-official/comdex/x/vault/keeper"
	"github.com/comdex-official/comdex/x/vault/types"
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
	app := chain.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	suite.app = app
	suite.ctx = ctx
	suite.keeper = suite.app.VaultKeeper
	suite.querier = keeper.QueryServer{Keeper: suite.keeper}
}

// Below are useful helpers to write test code easily.
func (suite *ModuleTestSuite) addr(addrNum int) sdk.AccAddress {
	addr := make(sdk.AccAddress, 20)
	binary.PutVarint(addr, int64(addrNum))
	return addr
}

func (suite *ModuleTestSuite) fundAddr(addr sdk.AccAddress, amt sdk.Coins) {
	suite.T().Helper()
	err := suite.app.BankKeeper.MintCoins(suite.ctx, types.ModuleName, amt)
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, types.ModuleName, addr, amt)
	suite.Require().NoError(err)
}

func (s *ModuleTestSuite) getBalances(addr sdk.AccAddress) sdk.Coins {
	return s.app.BankKeeper.GetAllBalances(s.ctx, addr)
}

func (s *ModuleTestSuite) CreateNewApp(appName string) uint64 {
	err := s.app.AssetKeeper.AddAppRecords(s.ctx, assettypes.AppData{
		Name:             appName,
		ShortName:        appName,
		MinGovDeposit:    sdk.NewInt(0),
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
		Decimals:              1000000,
		IsOnChain:             true,
		IsOraclePriceRequired: true,
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

func (s *ModuleTestSuite) CreateNewPair(addr sdk.Address, assetIn, assetOut uint64) uint64 {
	err := s.app.AssetKeeper.AddPairsRecords(s.ctx, assettypes.Pair{
		AssetIn:  assetIn,
		AssetOut: assetOut,
	})
	s.Suite.NoError(err)
	pairs := s.app.AssetKeeper.GetPairs(s.ctx)
	var pairID uint64
	for _, pair := range pairs {
		if pair.AssetIn == assetIn && pair.AssetOut == assetOut {
			pairID = pair.Id
			break
		}
	}
	s.Require().NotZero(pairID)
	return pairID
}

func (s *ModuleTestSuite) CreateNewExtendedVaultPair(
	pairName string,
	appMappingID, pairID uint64,
	isStableMintVault, isVaultActive bool,
) uint64 {
	err := s.app.AssetKeeper.WasmAddExtendedPairsVaultRecords(s.ctx, &bindings.MsgAddExtendedPairsVault{
		AppID:               appMappingID,
		PairID:              pairID,
		StabilityFee:        sdk.NewDecWithPrec(2, 2), // 0.02
		ClosingFee:          sdk.NewDec(0),
		LiquidationPenalty:  sdk.NewDecWithPrec(15, 2), // 0.15
		DrawDownFee:         sdk.NewDecWithPrec(1, 2),  // 0.01
		IsVaultActive:       isVaultActive,
		DebtCeiling:         1000000000000000000,
		DebtFloor:           100000000,
		IsStableMintVault:   isStableMintVault,
		MinCr:               sdk.NewDecWithPrec(23, 1), // 2.3
		PairName:            pairName,
		AssetOutOraclePrice: true,
		AssetOutPrice:       1000000,
		MinUsdValueLeft:     1000000,
	})
	s.Suite.Require().NoError(err)

	extendedVaultPairs, found := s.app.AssetKeeper.GetPairsVaults(s.ctx)
	s.Suite.Require().True(found)

	var extendedVaultPairID uint64
	for _, extendedVaultPair := range extendedVaultPairs {
		if extendedVaultPair.PairName == pairName && extendedVaultPair.AppId == appMappingID {
			extendedVaultPairID = extendedVaultPair.Id
			break
		}
	}
	s.Require().NotZero(extendedVaultPairID)
	return extendedVaultPairID
}
