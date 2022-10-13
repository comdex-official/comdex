package keeper_test

import (
	"encoding/binary"
	"testing"

	"github.com/comdex-official/comdex/app/wasm/bindings"
	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	chain "github.com/comdex-official/comdex/app"
	"github.com/comdex-official/comdex/x/vault/keeper"
	"github.com/comdex-official/comdex/x/vault/types"

	assettypes "github.com/comdex-official/comdex/x/asset/types"
	markettypes "github.com/comdex-official/comdex/x/market/types"
)

type KeeperTestSuite struct {
	suite.Suite

	app       *chain.App
	ctx       sdk.Context
	keeper    keeper.Keeper
	querier   keeper.QueryServer
	msgServer types.MsgServer
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (s *KeeperTestSuite) SetupTest() {
	s.app = chain.Setup(false)
	s.ctx = s.app.BaseApp.NewContext(false, tmproto.Header{})
	s.keeper = s.app.VaultKeeper
	s.querier = keeper.QueryServer{Keeper: s.keeper}
	s.msgServer = keeper.NewMsgServer(s.keeper)
}

// Below are just shortcuts to frequently-used functions.
func (s *KeeperTestSuite) getBalances(addr sdk.AccAddress) sdk.Coins {
	return s.app.BankKeeper.GetAllBalances(s.ctx, addr)
}

func (s *KeeperTestSuite) getBalance(addr sdk.AccAddress, denom string) sdk.Coin {
	return s.app.BankKeeper.GetBalance(s.ctx, addr, denom)
}

func (s *KeeperTestSuite) sendCoins(fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) {
	s.T().Helper()
	err := s.app.BankKeeper.SendCoins(s.ctx, fromAddr, toAddr, amt)
	s.Require().NoError(err)
}

// Below are useful helpers to write test code easily.
func (s *KeeperTestSuite) addr(addrNum int) sdk.AccAddress {
	addr := make(sdk.AccAddress, 20)
	binary.PutVarint(addr, int64(addrNum))
	return addr
}

func (s *KeeperTestSuite) fundAddr(addr sdk.AccAddress, amt sdk.Coins) {
	s.T().Helper()
	err := s.app.BankKeeper.MintCoins(s.ctx, types.ModuleName, amt)
	s.Require().NoError(err)
	err = s.app.BankKeeper.SendCoinsFromModuleToAccount(s.ctx, types.ModuleName, addr, amt)
	s.Require().NoError(err)
}

func newInt(i int64) sdk.Int {
	return sdk.NewInt(i)
}

func (s *KeeperTestSuite) CreateNewApp(appName string) uint64 {
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

func (s *KeeperTestSuite) CreateNewAsset(name, denom string, price uint64) uint64 {
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

func (s *KeeperTestSuite) CreateNewPair(addr sdk.Address, assetIn, assetOut uint64) uint64 {
	_, err := s.app.AssetKeeper.NewAddPair(s.ctx, &assettypes.MsgAddPairRequest{
		From:     addr.String(),
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

func (s *KeeperTestSuite) CreateNewExtendedVaultPair(
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
