package keeper_test

import (
	chain "github.com/comdex-official/comdex/app"
	assetKeeper "github.com/comdex-official/comdex/x/asset/keeper"
	assettypes "github.com/comdex-official/comdex/x/asset/types"
	auctionsV2Keeper "github.com/comdex-official/comdex/x/auctionsV2/keeper"
	auctionsV2types "github.com/comdex-official/comdex/x/auctionsV2/types"
	lendkeeper "github.com/comdex-official/comdex/x/lend/keeper"
	lendtypes "github.com/comdex-official/comdex/x/lend/types"
	liquidationKeeper "github.com/comdex-official/comdex/x/liquidationsV2/keeper"
	"github.com/comdex-official/comdex/x/liquidationsV2/types"
	marketKeeper "github.com/comdex-official/comdex/x/market/keeper"
	markettypes "github.com/comdex-official/comdex/x/market/types"
	vaultKeeper "github.com/comdex-official/comdex/x/vault/keeper"
	vaultTypes "github.com/comdex-official/comdex/x/vault/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"testing"
)

type KeeperTestSuite struct {
	suite.Suite

	app                 *chain.App
	ctx                 sdk.Context
	vaultKeeper         vaultKeeper.Keeper
	assetKeeper         assetKeeper.Keeper
	liquidationKeeper   liquidationKeeper.Keeper
	marketKeeper        marketKeeper.Keeper
	querier             liquidationKeeper.QueryServer
	vaultQuerier        vaultKeeper.QueryServer
	msgServer           types.MsgServer
	vaultMsgServer      vaultTypes.MsgServer
	lendKeeper          lendkeeper.Keeper
	lendQuerier         lendkeeper.QueryServer
	auctionsV2Keeper    auctionsV2Keeper.Keeper
	auctionsV2MsgServer auctionsV2types.MsgServer
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (s *KeeperTestSuite) SetupTest() {
	s.app = chain.Setup(false)
	s.ctx = s.app.BaseApp.NewContext(false, tmproto.Header{})
	s.vaultKeeper = s.app.VaultKeeper
	s.liquidationKeeper = s.app.NewliqKeeper
	s.assetKeeper = s.app.AssetKeeper
	s.querier = liquidationKeeper.QueryServer{Keeper: s.liquidationKeeper}
	s.msgServer = liquidationKeeper.NewMsgServerImpl(s.liquidationKeeper)
	s.vaultMsgServer = vaultKeeper.NewMsgServer(s.vaultKeeper)
	s.vaultQuerier = vaultKeeper.QueryServer{Keeper: s.vaultKeeper}
	s.marketKeeper = s.app.MarketKeeper
	s.lendKeeper = s.app.LendKeeper
	s.lendQuerier = lendkeeper.QueryServer{Keeper: s.lendKeeper}
	s.auctionsV2Keeper = s.app.NewaucKeeper
	s.auctionsV2MsgServer = auctionsV2Keeper.NewMsgServerImpl(s.auctionsV2Keeper)
}

func (s *KeeperTestSuite) CreateNewAsset(name, denom string, twa uint64) uint64 {
	err := s.app.AssetKeeper.AddAssetRecords(s.ctx, assettypes.Asset{
		Name:                  name,
		Denom:                 denom,
		Decimals:              sdk.NewInt(1000000),
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
		AssetID:       assetID,
		ScriptID:      10,
		Twa:           twa,
		CurrentIndex:  1,
		IsPriceActive: true,
		PriceValue:    nil,
	}

	s.app.MarketKeeper.SetTwa(s.ctx, twa1)

	return assetID
}

func newInt(i int64) sdk.Int {
	return sdk.NewInt(i)
}

func newDec(i string) sdk.Dec {
	dec, _ := sdk.NewDecFromStr(i)
	return dec
}

func (s *KeeperTestSuite) AddAssetRatesStats(AssetID uint64, UOptimal, Base, Slope1, Slope2 sdk.Dec, EnableStableBorrow bool, StableBase, StableSlope1, StableSlope2, LTV, LiquidationThreshold, LiquidationPenalty, LiquidationBonus, ReserveFactor sdk.Dec, CAssetID uint64) uint64 {
	err := s.app.LendKeeper.AddAssetRatesParams(s.ctx, lendtypes.AssetRatesParams{
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

func (s *KeeperTestSuite) AddAssetRatesPoolPairs(AssetID uint64, UOptimal, Base, Slope1, Slope2 sdk.Dec, EnableStableBorrow bool, StableBase, StableSlope1, StableSlope2, LTV, LiquidationThreshold, LiquidationPenalty, LiquidationBonus, ReserveFactor sdk.Dec, CAssetID uint64, moduleName, cPoolName string, assetData []*lendtypes.AssetDataPoolMapping, MinUsdValueLeft uint64, IsIsolated bool) uint64 {
	err := s.app.LendKeeper.AddAssetRatesPoolPairs(s.ctx, lendtypes.AssetRatesPoolPairs{
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
		ModuleName:           moduleName,
		CPoolName:            cPoolName,
		AssetData:            assetData,
		MinUsdValueLeft:      MinUsdValueLeft,
		//IsIsolated:           IsIsolated,
	})
	s.Require().NoError(err)
	return AssetID
}

func (s *KeeperTestSuite) fundAddr(addr sdk.AccAddress, amt sdk.Coins) {
	s.T().Helper()
	err := s.app.BankKeeper.MintCoins(s.ctx, types.ModuleName, amt)
	s.Require().NoError(err)
	err = s.app.BankKeeper.SendCoinsFromModuleToAccount(s.ctx, types.ModuleName, addr, amt)
	s.Require().NoError(err)
}

func (s *KeeperTestSuite) addAuctionParams(auctionParams auctionsV2types.AuctionParams) {
	s.T().Helper()
	s.app.NewaucKeeper.SetAuctionParams(s.ctx, auctionParams)
}
