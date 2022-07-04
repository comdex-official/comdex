package keeper_test

import (
	"encoding/binary"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	chain "github.com/comdex-official/comdex/app"
	"github.com/comdex-official/comdex/x/liquidity"
	"github.com/comdex-official/comdex/x/liquidity/amm"
	"github.com/comdex-official/comdex/x/liquidity/keeper"
	"github.com/comdex-official/comdex/x/liquidity/types"
)

type KeeperTestSuite struct {
	suite.Suite

	app       *chain.App
	ctx       sdk.Context
	keeper    keeper.Keeper
	querier   keeper.Querier
	msgServer types.MsgServer
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (s *KeeperTestSuite) SetupTest() {
	s.app = chain.Setup(false)
	s.ctx = s.app.BaseApp.NewContext(false, tmproto.Header{})
	s.keeper = s.app.LiquidityKeeper
	s.querier = keeper.Querier{Keeper: s.keeper}
	s.msgServer = keeper.NewMsgServerImpl(s.keeper)
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

func (s *KeeperTestSuite) nextBlock() {
	liquidity.EndBlocker(s.ctx, s.keeper)
	liquidity.BeginBlocker(s.ctx, s.keeper)
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

func (s *KeeperTestSuite) createPair(appID uint64, creator sdk.AccAddress, baseCoinDenom, quoteCoinDenom string, fund bool) types.Pair {
	s.T().Helper()
	params, err := s.keeper.GetGenericParams(s.ctx, appID)
	s.Require().NoError(err)
	if fund {
		s.fundAddr(creator, params.PairCreationFee)
	}
	msg := types.NewMsgCreatePair(appID, creator, baseCoinDenom, quoteCoinDenom)
	s.Require().NoError(msg.ValidateBasic())
	pair, err := s.keeper.CreatePair(s.ctx, msg, false)
	s.Require().NoError(err)
	return pair
}

func (s *KeeperTestSuite) createPool(appID uint64, creator sdk.AccAddress, pairId uint64, depositCoins sdk.Coins, fund bool) types.Pool {
	s.T().Helper()
	params, err := s.keeper.GetGenericParams(s.ctx, appID)
	s.Require().NoError(err)
	if fund {
		s.fundAddr(creator, depositCoins.Add(params.PoolCreationFee...))
	}
	msg := types.NewMsgCreatePool(appID, creator, pairId, depositCoins)
	s.Require().NoError(msg.ValidateBasic())
	pool, err := s.keeper.CreatePool(s.ctx, msg)
	s.Require().NoError(err)
	return pool
}

func (s *KeeperTestSuite) deposit(appID uint64, depositor sdk.AccAddress, poolId uint64, depositCoins sdk.Coins, fund bool) types.DepositRequest {
	s.T().Helper()
	if fund {
		s.fundAddr(depositor, depositCoins)
	}
	req, err := s.keeper.Deposit(s.ctx, types.NewMsgDeposit(appID, depositor, poolId, depositCoins))
	s.Require().NoError(err)
	return req
}

func (s *KeeperTestSuite) withdraw(appID uint64, withdrawer sdk.AccAddress, poolId uint64, poolCoin sdk.Coin) types.WithdrawRequest {
	s.T().Helper()
	req, err := s.keeper.Withdraw(s.ctx, types.NewMsgWithdraw(appID, withdrawer, poolId, poolCoin))
	s.Require().NoError(err)
	return req
}

func (s *KeeperTestSuite) limitOrder(
	appID uint64, orderer sdk.AccAddress, pairId uint64, dir types.OrderDirection,
	price sdk.Dec, amt sdk.Int, orderLifespan time.Duration, fund bool) types.Order {
	s.T().Helper()

	params, err := s.keeper.GetGenericParams(s.ctx, appID)
	s.Require().NoError(err)
	pair, found := s.keeper.GetPair(s.ctx, appID, pairId)
	s.Require().True(found)
	var ammDir amm.OrderDirection
	var offerCoinDenom, demandCoinDenom string
	switch dir {
	case types.OrderDirectionBuy:
		ammDir = amm.Buy
		offerCoinDenom, demandCoinDenom = pair.QuoteCoinDenom, pair.BaseCoinDenom
	case types.OrderDirectionSell:
		ammDir = amm.Sell
		offerCoinDenom, demandCoinDenom = pair.BaseCoinDenom, pair.QuoteCoinDenom
	}
	offerCoin := sdk.NewCoin(offerCoinDenom, amm.OfferCoinAmount(ammDir, price, amt))
	offerCoin.Amount = offerCoin.Amount.Add(offerCoin.Amount.ToDec().Mul(params.SwapFeeRate).TruncateInt())
	if fund {
		s.fundAddr(orderer, sdk.NewCoins(offerCoin))
	}
	msg := types.NewMsgLimitOrder(
		appID, orderer, pairId, dir, offerCoin, demandCoinDenom,
		price, amt, orderLifespan)
	s.Require().NoError(msg.ValidateBasic())
	req, err := s.keeper.LimitOrder(s.ctx, msg)
	s.Require().NoError(err)
	return req
}

func (s *KeeperTestSuite) buyLimitOrder(
	appID uint64, orderer sdk.AccAddress, pairId uint64, price sdk.Dec,
	amt sdk.Int, orderLifespan time.Duration, fund bool) types.Order {
	s.T().Helper()
	return s.limitOrder(
		appID, orderer, pairId, types.OrderDirectionBuy, price, amt, orderLifespan, fund)
}

func (s *KeeperTestSuite) sellLimitOrder(
	appID uint64, orderer sdk.AccAddress, pairId uint64, price sdk.Dec,
	amt sdk.Int, orderLifespan time.Duration, fund bool) types.Order {
	s.T().Helper()
	return s.limitOrder(
		appID, orderer, pairId, types.OrderDirectionSell, price, amt, orderLifespan, fund)
}

func (s *KeeperTestSuite) marketOrder(
	appID uint64, orderer sdk.AccAddress, pairId uint64, dir types.OrderDirection,
	amt sdk.Int, orderLifespan time.Duration, fund bool) types.Order {
	s.T().Helper()
	pair, found := s.keeper.GetPair(s.ctx, appID, pairId)
	s.Require().True(found)
	s.Require().NotNil(pair.LastPrice)
	lastPrice := *pair.LastPrice
	params, err := s.keeper.GetGenericParams(s.ctx, appID)
	s.Require().NoError(err)
	var offerCoin sdk.Coin
	var demandCoinDenom string
	switch dir {
	case types.OrderDirectionBuy:
		maxPrice := lastPrice.Mul(sdk.OneDec().Add(params.MaxPriceLimitRatio))
		offerCoin = sdk.NewCoin(pair.QuoteCoinDenom, amm.OfferCoinAmount(amm.Buy, maxPrice, amt))
		demandCoinDenom = pair.BaseCoinDenom
	case types.OrderDirectionSell:
		offerCoin = sdk.NewCoin(pair.BaseCoinDenom, amt)
		demandCoinDenom = pair.QuoteCoinDenom
	}
	offerCoin.Amount = offerCoin.Amount.Add(offerCoin.Amount.ToDec().Mul(params.SwapFeeRate).TruncateInt())
	if fund {
		s.fundAddr(orderer, sdk.NewCoins(offerCoin))
	}
	msg := types.NewMsgMarketOrder(
		appID, orderer, pairId, dir, offerCoin, demandCoinDenom,
		amt, orderLifespan)
	s.Require().NoError(msg.ValidateBasic())
	req, err := s.keeper.MarketOrder(s.ctx, msg)
	s.Require().NoError(err)
	return req
}

func (s *KeeperTestSuite) buyMarketOrder(
	appID uint64, orderer sdk.AccAddress, pairId uint64,
	amt sdk.Int, orderLifespan time.Duration, fund bool) types.Order {
	s.T().Helper()
	return s.marketOrder(
		appID, orderer, pairId, types.OrderDirectionBuy, amt, orderLifespan, fund)
}

//nolint
func (s *KeeperTestSuite) sellMarketOrder(
	appID uint64, orderer sdk.AccAddress, pairId uint64,
	amt sdk.Int, orderLifespan time.Duration, fund bool) types.Order {
	s.T().Helper()
	return s.marketOrder(
		appID, orderer, pairId, types.OrderDirectionSell, amt, orderLifespan, fund)
}

//nolint
func (s *KeeperTestSuite) cancelOrder(appID uint64, orderer sdk.AccAddress, pairId, orderId uint64) {
	s.T().Helper()
	err := s.keeper.CancelOrder(s.ctx, types.NewMsgCancelOrder(appID, orderer, pairId, orderId))
	s.Require().NoError(err)
}

func (s *KeeperTestSuite) cancelAllOrders(appID uint64, orderer sdk.AccAddress, pairIds []uint64) {
	s.T().Helper()
	err := s.keeper.CancelAllOrders(s.ctx, types.NewMsgCancelAllOrders(appID, orderer, pairIds))
	s.Require().NoError(err)
}

func coinEq(exp, got sdk.Coin) (bool, string, string, string) {
	return exp.IsEqual(got), "expected:\t%v\ngot:\t\t%v", exp.String(), got.String()
}

func coinsEq(exp, got sdk.Coins) (bool, string, string, string) {
	return exp.IsEqual(got), "expected:\t%v\ngot:\t\t%v", exp.String(), got.String()
}

func intEq(exp, got sdk.Int) (bool, string, string, string) {
	return exp.Equal(got), "expected:\t%v\ngot:\t\t%v", exp.String(), got.String()
}

func decEq(exp, got sdk.Dec) (bool, string, string, string) {
	return exp.Equal(got), "expected:\t%v\ngot:\t\t%v", exp.String(), got.String()
}

func newInt(i int64) sdk.Int {
	return sdk.NewInt(i)
}
