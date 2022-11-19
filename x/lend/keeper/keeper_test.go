package keeper_test

import (
	"encoding/binary"
	"testing"

	assettypes "github.com/petrichormoney/petri/x/asset/types"
	markettypes "github.com/petrichormoney/petri/x/market/types"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	chain "github.com/petrichormoney/petri/app"
	"github.com/petrichormoney/petri/x/lend/keeper"
	"github.com/petrichormoney/petri/x/lend/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
	s.keeper = s.app.LendKeeper
	s.querier = keeper.QueryServer{Keeper: s.keeper}
	s.msgServer = keeper.NewMsgServerImpl(s.keeper)
}

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

func newDec(i string) sdk.Dec {
	dec, _ := sdk.NewDecFromStr(i)
	return dec
}

//func (s *KeeperTestSuite) SetOraclePrice(symbol string, price uint64) {
//	var (
//		store = s.app.MarketKeeper.Store(s.ctx)
//		key   = markettypes.PriceForMarketKey(symbol)
//	)
//	value := s.app.AppCodec().MustMarshal(
//		&protobuftypes.UInt64Value{
//			Value: price,
//		},
//	)
//	store.Set(key, value)
//}

func (s *KeeperTestSuite) CreateNewAsset(name, denom string, price uint64) uint64 {
	err := s.app.AssetKeeper.AddAssetRecords(s.ctx, assettypes.Asset{
		Name:                  name,
		Denom:                 denom,
		Decimals:              sdk.NewInt(1000000),
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

	//market := markettypes.Market{
	//	Symbol:   name,
	//	ScriptID: 12,
	//	Rates:    price,
	//}
	//s.app.MarketKeeper.SetMarket(s.ctx, market)

	// exists := s.app.MarketKeeper.HasMarketForAsset(s.ctx, assetID)

	twa1 := markettypes.TimeWeightedAverage{
		AssetID:       1,
		ScriptID:      10,
		Twa:           1000000,
		CurrentIndex:  1,
		IsPriceActive: true,
		PriceValue:    nil,
	}
	twa2 := markettypes.TimeWeightedAverage{
		AssetID:       2,
		ScriptID:      10,
		Twa:           1000000,
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
		Twa:           1000000,
		CurrentIndex:  1,
		IsPriceActive: true,
		PriceValue:    nil,
	}

	s.app.MarketKeeper.SetTwa(s.ctx, twa1)
	s.app.MarketKeeper.SetTwa(s.ctx, twa2)
	s.app.MarketKeeper.SetTwa(s.ctx, twa3)
	s.app.MarketKeeper.SetTwa(s.ctx, twa4)
	// s.Suite.Require().False(exists)
	// s.app.MarketKeeper.SetMarketForAsset(s.ctx, assetID, name)

	return assetID
}

func (s *KeeperTestSuite) CreateNewPool(moduleName, cPoolName string, assetData []*types.AssetDataPoolMapping) uint64 {
	err := s.app.LendKeeper.AddPoolRecords(s.ctx, types.Pool{
		ModuleName: moduleName,
		CPoolName:  cPoolName,
		AssetData:  assetData,
	})
	s.Require().NoError(err)

	pools := s.app.LendKeeper.GetPools(s.ctx)
	var poolID uint64
	for _, pool := range pools {
		poolID = pool.PoolID
	}
	s.Require().NotZero(poolID)

	return poolID
}

func (s *KeeperTestSuite) AddAssetRatesStats(AssetID uint64, UOptimal, Base, Slope1, Slope2 sdk.Dec, EnableStableBorrow bool, StableBase, StableSlope1, StableSlope2, LTV, LiquidationThreshold, LiquidationPenalty, LiquidationBonus, ReserveFactor sdk.Dec, CAssetID uint64) uint64 {
	err := s.app.LendKeeper.AddAssetRatesParams(s.ctx, types.AssetRatesParams{
		AssetID:              AssetID,
		UOptimal:             UOptimal,
		Base:                 Base,
		Slope1:               Slope1,
		Slope2:               Slope2,
		EnableStableBorrow:   EnableStableBorrow,
		StableBase:           StableBase,
		StableSlope1:         StableSlope1,
		StableSlope2:         StableSlope2,
		Ltv:                  LTV,
		LiquidationThreshold: LiquidationThreshold,
		LiquidationPenalty:   LiquidationPenalty,
		LiquidationBonus:     LiquidationBonus,
		ReserveFactor:        ReserveFactor,
		CAssetID:             CAssetID,
	})
	s.Require().NoError(err)
	return AssetID
}

func (s *KeeperTestSuite) AddExtendedLendPair(AssetIn, AssetOut uint64, IsInterPool bool, AssetOutPoolID, MinUsdValueLeft uint64) uint64 {
	err := s.app.LendKeeper.AddLendPairsRecords(s.ctx, types.Extended_Pair{
		AssetIn:         AssetIn,
		AssetOut:        AssetOut,
		IsInterPool:     IsInterPool,
		AssetOutPoolID:  AssetOutPoolID,
		MinUsdValueLeft: MinUsdValueLeft,
	})
	s.Require().NoError(err)
	pairs := s.app.LendKeeper.GetLendPairs(s.ctx)
	var pairID uint64
	for _, pair := range pairs {
		if pair.AssetIn == AssetIn && pair.AssetOut == AssetOut && pair.IsInterPool == IsInterPool {
			pairID = pair.Id
			break
		}
	}
	s.Require().NotZero(pairID)
	return pairID
}

func (s *KeeperTestSuite) AddAssetToPair(AssetID, PoolID uint64, PairID []uint64) {
	err := s.app.LendKeeper.AddAssetToPair(s.ctx, types.AssetToPairMapping{
		AssetID: AssetID,
		PoolID:  PoolID,
		PairID:  PairID,
	})
	s.Require().NoError(err)
}

func (s *KeeperTestSuite) CreateNewApp(appName, shortName string) uint64 {
	err := s.app.AssetKeeper.AddAppRecords(s.ctx, assettypes.AppData{
		Name:             appName,
		ShortName:        shortName,
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
