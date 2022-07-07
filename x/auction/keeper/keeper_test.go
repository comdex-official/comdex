package keeper_test

import (
	"time"

	"testing"

	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	chain "github.com/comdex-official/comdex/app"
	assetKeeper "github.com/comdex-official/comdex/x/asset/keeper"
	"github.com/comdex-official/comdex/x/auction/keeper"
	"github.com/comdex-official/comdex/x/auction/types"
	collectKeeper "github.com/comdex-official/comdex/x/collector/keeper"
	liquidationKeeper "github.com/comdex-official/comdex/x/liquidation/keeper"
	liquidationTypes "github.com/comdex-official/comdex/x/liquidation/types"
	marketKeeper "github.com/comdex-official/comdex/x/market/keeper"
	tokenmintKeeper "github.com/comdex-official/comdex/x/tokenmint/keeper"
	vaultKeeper "github.com/comdex-official/comdex/x/vault/keeper"
	vaultTypes "github.com/comdex-official/comdex/x/vault/types"
)

type KeeperTestSuite struct {
	suite.Suite

	app                  *chain.App
	ctx                  sdk.Context
	vaultKeeper          vaultKeeper.Keeper
	assetKeeper          assetKeeper.Keeper
	liquidationKeeper    liquidationKeeper.Keeper
	tokenmintKeeper      tokenmintKeeper.Keeper
	marketKeeper         marketKeeper.Keeper
	collectorKeeper      collectKeeper.Keeper
	liquidationQuerier   liquidationKeeper.QueryServer
	vaultQuerier         vaultKeeper.QueryServer
	liquidationMsgServer liquidationTypes.MsgServer
	vaultMsgServer       vaultTypes.MsgServer
	keeper               keeper.Keeper
	auctionMsgServer     types.MsgServer
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (s *KeeperTestSuite) SetupTest() {
	s.app = chain.Setup(false)
	s.ctx = s.app.BaseApp.NewContext(false, tmproto.Header{})
	s.vaultKeeper = s.app.VaultKeeper
	s.liquidationKeeper = s.app.LiquidationKeeper
	s.assetKeeper = s.app.AssetKeeper
	s.collectorKeeper = s.app.CollectorKeeper
	s.liquidationQuerier = liquidationKeeper.QueryServer{Keeper: s.liquidationKeeper}
	s.liquidationMsgServer = liquidationKeeper.NewMsgServer(s.liquidationKeeper)
	s.vaultMsgServer = vaultKeeper.NewMsgServer(s.vaultKeeper)
	s.vaultQuerier = vaultKeeper.QueryServer{Keeper: s.vaultKeeper}
	s.marketKeeper = s.app.MarketKeeper
	s.keeper = s.app.AuctionKeeper
	s.auctionMsgServer = keeper.NewMsgServiceServer(s.keeper)
	s.tokenmintKeeper = s.app.TokenmintKeeper
}
func (s *KeeperTestSuite) getBalance(addr string, denom string) (coin sdk.Coin, err error) {
	addr1, err := sdk.AccAddressFromBech32(addr)
	if err != nil {
		return coin, err
	}
	return s.app.BankKeeper.GetBalance(s.ctx, addr1, denom), nil
}

//
//// Below are just shortcuts to frequently-used functions.
//func (s *KeeperTestSuite) getBalances(addr sdk.AccAddress) sdk.Coins {
//	return s.app.bankKeeper.GetAllBalances(s.ctx, addr)
//}
//
//func (s *KeeperTestSuite) getBalance(addr sdk.AccAddress, denom string) sdk.Coin {
//	return s.app.bankKeeper.GetBalance(s.ctx, addr, denom)
//}
//
//func (s *KeeperTestSuite) sendCoins(fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) {
//	s.T().Helper()
//	err := s.app.bankKeeper.SendCoins(s.ctx, fromAddr, toAddr, amt)
//	s.Require().NoError(err)
//}
//
//func (s *KeeperTestSuite) nextBlock() {
//	liquidity.EndBlocker(s.ctx, s.keeper)
//	liquidity.BeginBlocker(s.ctx, s.keeper)
//}
//
//// Below are useful helpers to write test code easily.
//func (s *KeeperTestSuite) addr(addrNum int) sdk.AccAddress {
//	addr := make(sdk.AccAddress, 20)
//	binary.PutVarint(addr, int64(addrNum))
//	return addr
//}

func (s *KeeperTestSuite) fundAddr(addr sdk.AccAddress, amt sdk.Coin) {
	amt1 := sdk.NewCoins(amt)
	s.T().Helper()
	err := s.app.BankKeeper.MintCoins(s.ctx, liquidationTypes.ModuleName, amt1)
	s.Require().NoError(err)
	err = s.app.BankKeeper.SendCoinsFromModuleToAccount(s.ctx, liquidationTypes.ModuleName, addr, amt1)
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
