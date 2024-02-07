package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	"encoding/binary"
	"testing"
	"time"

	assettypes "github.com/comdex-official/comdex/x/asset/types"
	markettypes "github.com/comdex-official/comdex/x/market/types"

	collectorTypes "github.com/comdex-official/comdex/x/collector/types"

	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"

	chain "github.com/comdex-official/comdex/app"
	assetKeeper "github.com/comdex-official/comdex/x/asset/keeper"
	auctionKeeper "github.com/comdex-official/comdex/x/auction/keeper"
	collectorKeeper "github.com/comdex-official/comdex/x/collector/keeper"
	tokenmintKeeper "github.com/comdex-official/comdex/x/tokenmint/keeper"
	tokenmintTypes "github.com/comdex-official/comdex/x/tokenmint/types"
)

type KeeperTestSuite struct {
	suite.Suite

	app             *chain.App
	ctx             sdk.Context
	assetKeeper     assetKeeper.Keeper
	collectorKeeper collectorKeeper.Keeper
	auctionKeeper   auctionKeeper.Keeper
	tokenmintKeeper tokenmintKeeper.Keeper
	querier         tokenmintKeeper.QueryServer
	msgServer       tokenmintTypes.MsgServer
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (s *KeeperTestSuite) SetupTest() {
	s.app = chain.Setup(s.T(), false)
	s.ctx = s.app.BaseApp.NewContext(false)
	s.collectorKeeper = s.app.CollectorKeeper
	s.assetKeeper = s.app.AssetKeeper
	s.auctionKeeper = s.app.AuctionKeeper
	s.tokenmintKeeper = s.app.TokenmintKeeper
	s.querier = tokenmintKeeper.QueryServer{Keeper: s.tokenmintKeeper}
	s.msgServer = tokenmintKeeper.NewMsgServer(s.tokenmintKeeper)
}

// // Below are just shortcuts to frequently-used functions.
//
//	func (s *KeeperTestSuite) getBalances(addr sdk.AccAddress) sdk.Coins {
//		return s.app.bankKeeper.GetAllBalances(s.ctx, addr)
//	}
func (s *KeeperTestSuite) getBalance(addr string, denom string) (coin sdk.Coin, err error) {
	addr1, err := sdk.AccAddressFromBech32(addr)
	if err != nil {
		return coin, err
	}
	return s.app.BankKeeper.GetBalance(s.ctx, addr1, denom), nil
}

//
//func (s *KeeperTestSuite) sendCoins(fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) {
//	s.T().Helper()
//	err := s.app.bankKeeper.SendCoins(s.ctx, fromAddr, toAddr, amt)
//	s.Require().NoError(err)
//}

func (s *KeeperTestSuite) addr(addrNum int) sdk.AccAddress {
	addr := make(sdk.AccAddress, 20)
	binary.PutVarint(addr, int64(addrNum))
	return addr
}

func (s *KeeperTestSuite) fundAddr(addr string, amt sdk.Coin) {
	s.T().Helper()
	err := s.app.BankKeeper.MintCoins(s.ctx, collectorTypes.ModuleName, sdk.NewCoins(amt))
	s.Require().NoError(err)
	addr1, err := sdk.AccAddressFromBech32(addr)
	err = s.app.BankKeeper.SendCoinsFromModuleToAccount(s.ctx, collectorTypes.ModuleName, addr1, sdk.NewCoins(amt))
	s.Require().NoError(err)
}

func (s *KeeperTestSuite) advanceseconds(dur int64) {
	s.ctx = s.ctx.WithBlockTime(s.ctx.BlockTime().Add(time.Second * time.Duration(dur)))
}

// ParseCoins parses and returns sdk.Coins.
func ParseCoin(s string) sdk.Coin {
	coins, err := sdk.ParseCoinNormalized(s)
	if err != nil {
		panic(err)
	}
	return coins
}

func (s *KeeperTestSuite) CreateNewApp(appName string) uint64 {
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

func (s *KeeperTestSuite) CreateNewAsset(name, denom string, price uint64) uint64 {
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

	market := markettypes.TimeWeightedAverage{
		AssetID:       assetID,
		ScriptID:      12,
		Twa:           price,
		CurrentIndex:  0,
		IsPriceActive: true,
		PriceValue:    []uint64{price},
	}
	s.app.MarketKeeper.SetTwa(s.ctx, market)
	_, err = s.app.MarketKeeper.GetLatestPrice(s.ctx, assetID)
	s.Suite.NoError(err)

	return assetID
}
