package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	chain "github.com/petrichormoney/petri/app"
	assetKeeper "github.com/petrichormoney/petri/x/asset/keeper"
	auctionTypes "github.com/petrichormoney/petri/x/auction/types"
	"github.com/petrichormoney/petri/x/liquidation/keeper"
	liquidationKeeper "github.com/petrichormoney/petri/x/liquidation/keeper"
	"github.com/petrichormoney/petri/x/liquidation/types"
	marketKeeper "github.com/petrichormoney/petri/x/market/keeper"
	vaultKeeper "github.com/petrichormoney/petri/x/vault/keeper"
	vaultTypes "github.com/petrichormoney/petri/x/vault/types"
)

type KeeperTestSuite struct {
	suite.Suite

	app               *chain.App
	ctx               sdk.Context
	vaultKeeper       vaultKeeper.Keeper
	assetKeeper       assetKeeper.Keeper
	liquidationKeeper liquidationKeeper.Keeper
	marketKeeper      marketKeeper.Keeper
	querier           liquidationKeeper.QueryServer
	vaultQuerier      vaultKeeper.QueryServer
	msgServer         types.MsgServer
	vaultMsgServer    vaultTypes.MsgServer
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
	s.querier = keeper.QueryServer{Keeper: s.liquidationKeeper}
	s.msgServer = keeper.NewMsgServer(s.liquidationKeeper)
	s.vaultMsgServer = vaultKeeper.NewMsgServer(s.vaultKeeper)
	s.vaultQuerier = vaultKeeper.QueryServer{Keeper: s.vaultKeeper}
	s.marketKeeper = s.app.MarketKeeper
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
	err := s.app.BankKeeper.MintCoins(s.ctx, types.ModuleName, amt1)
	s.Require().NoError(err)
	err = s.app.BankKeeper.SendCoinsFromModuleToAccount(s.ctx, types.ModuleName, addr, amt1)
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

func (s *KeeperTestSuite) AddAuctionParams() {
	ctx := &s.ctx
	auctionParams := auctionTypes.AuctionParams{
		AppId:                  1,
		AuctionDurationSeconds: 300,
		Buffer:                 sdk.MustNewDecFromStr("1.2"),
		Cusp:                   sdk.MustNewDecFromStr("0.6"),
		Step:                   sdk.NewIntFromUint64(1),
		PriceFunctionType:      1,
		SurplusId:              1,
		DebtId:                 2,
		DutchId:                3,
		BidDurationSeconds:     300,
	}
	s.app.AuctionKeeper.SetAuctionParams(*ctx, auctionParams)
}
