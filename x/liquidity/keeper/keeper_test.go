package keeper_test

import (
	"encoding/binary"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	chain "github.com/comdex-official/comdex/app"
	assettypes "github.com/comdex-official/comdex/x/asset/types"
	"github.com/comdex-official/comdex/x/liquidity"
	"github.com/comdex-official/comdex/x/liquidity/amm"
	"github.com/comdex-official/comdex/x/liquidity/keeper"
	"github.com/comdex-official/comdex/x/liquidity/types"
	markettypes "github.com/comdex-official/comdex/x/market/types"

	utils "github.com/comdex-official/comdex/types"
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
	s.app = chain.Setup(s.T(), false)
	s.ctx = s.app.BaseApp.NewContext(false)
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
	liquidity.EndBlocker(s.ctx, s.keeper, s.app.AssetKeeper)
	liquidity.BeginBlocker(s.ctx, s.keeper, s.app.AssetKeeper)
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

func newInt(i int64) sdkmath.Int {
	return sdkmath.NewInt(i)
}

func newDec(i int64) sdkmath.LegacyDec {
	return sdkmath.LegacyNewDec(i)
}

func coinEq(exp, got sdk.Coin) (bool, string, string, string) {
	return exp.IsEqual(got), "expected:\t%v\ngot:\t\t%v", exp.String(), got.String()
}

func coinsEq(exp, got sdk.Coins) (bool, string, string, string) {
	return exp.Equal(got), "expected:\t%v\ngot:\t\t%v", exp.String(), got.String()
}

func intEq(exp, got sdkmath.Int) (bool, string, string, string) {
	return exp.Equal(got), "expected:\t%v\ngot:\t\t%v", exp.String(), got.String()
}

func decEq(exp, got sdkmath.LegacyDec) (bool, string, string, string) {
	return exp.Equal(got), "expected:\t%v\ngot:\t\t%v", exp.String(), got.String()
}

func (s *KeeperTestSuite) CreateNewApp(appName string) uint64 {
	err := s.app.AssetKeeper.AddAppRecords(s.ctx, assettypes.AppData{
		Name:             strings.ToLower(appName),
		ShortName:        strings.ToLower(appName),
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

func (s *KeeperTestSuite) CreateNewAsset(name, denom string, price uint64) assettypes.Asset {
	err := s.app.AssetKeeper.AddAssetRecords(s.ctx, assettypes.Asset{
		Name:                  name,
		Denom:                 denom,
		Decimals:              sdkmath.NewInt(1000000),
		IsOnChain:             true,
		IsOraclePriceRequired: true,
	})
	s.Require().NoError(err)
	assets := s.app.AssetKeeper.GetAssets(s.ctx)
	var assetObj assettypes.Asset
	for _, asset := range assets {
		if asset.Denom == denom {
			assetObj = asset
			break
		}
	}
	s.Require().NotZero(assetObj.Id)

	market := markettypes.TimeWeightedAverage{
		AssetID:       assetObj.Id,
		ScriptID:      12,
		Twa:           price,
		CurrentIndex:  0,
		IsPriceActive: true,
		PriceValue:    []uint64{price},
	}
	s.app.MarketKeeper.SetTwa(s.ctx, market)
	_, err = s.app.MarketKeeper.GetLatestPrice(s.ctx, assetObj.Id)
	s.Suite.NoError(err)

	return assetObj
}

func (s *KeeperTestSuite) CreateNewLiquidityPair(appID uint64, creator sdk.AccAddress, baseCoinDenom, quoteCoinDenom string) types.Pair {
	params, err := s.keeper.GetGenericParams(s.ctx, appID)
	s.Require().NoError(err)

	s.fundAddr(creator, params.PairCreationFee)

	msg := types.NewMsgCreatePair(appID, creator, baseCoinDenom, quoteCoinDenom)
	pair, err := s.keeper.CreatePair(s.ctx, msg, false)

	s.Require().NoError(err)
	s.Require().IsType(types.Pair{}, pair)

	return pair
}

func (s *KeeperTestSuite) CreateNewLiquidityPool(appID, pairID uint64, creator sdk.AccAddress, depositCoins string) types.Pool {
	params, err := s.keeper.GetGenericParams(s.ctx, appID)
	s.Require().NoError(err)

	parsedDepositCoins := utils.ParseCoins(depositCoins)

	s.fundAddr(creator, params.PoolCreationFee)
	s.fundAddr(creator, parsedDepositCoins)
	msg := types.NewMsgCreatePool(appID, creator, pairID, parsedDepositCoins)
	pool, err := s.keeper.CreatePool(s.ctx, msg)
	s.Require().NoError(err)
	s.Require().IsType(types.Pool{}, pool)

	return pool
}

func (s *KeeperTestSuite) CreateNewLiquidityRangedPool(appID, pairID uint64, creator sdk.AccAddress, depositCoins string, minPrice, maxPrice, initialPrice sdkmath.LegacyDec) types.Pool {
	params, err := s.keeper.GetGenericParams(s.ctx, appID)
	s.Require().NoError(err)

	parsedDepositCoins := utils.ParseCoins(depositCoins)

	s.fundAddr(creator, params.PoolCreationFee)
	s.fundAddr(creator, parsedDepositCoins)
	msg := types.NewMsgCreateRangedPool(appID, creator, pairID, parsedDepositCoins, minPrice, maxPrice, initialPrice)
	pool, err := s.keeper.CreateRangedPool(s.ctx, msg)
	s.Require().NoError(err)
	s.Require().IsType(types.Pool{}, pool)

	return pool
}

func (s *KeeperTestSuite) Deposit(appID, poolID uint64, depositor sdk.AccAddress, depositCoins string) types.DepositRequest {
	msg := types.NewMsgDeposit(
		appID, depositor, poolID, utils.ParseCoins(depositCoins),
	)
	s.fundAddr(depositor, msg.DepositCoins)
	req, err := s.keeper.Deposit(s.ctx, msg)
	s.Require().NoError(err)
	s.Require().IsType(types.DepositRequest{}, req)
	return req
}

func (s *KeeperTestSuite) Withdraw(appID, poolID uint64, withdrawer sdk.AccAddress, poolCoin sdk.Coin) types.WithdrawRequest {
	msg := types.NewMsgWithdraw(
		appID, withdrawer, poolID, poolCoin,
	)
	req, err := s.keeper.Withdraw(s.ctx, msg)
	s.Require().NoError(err)
	s.Require().IsType(types.WithdrawRequest{}, req)
	return req
}

func (s *KeeperTestSuite) LimitOrder(
	appID uint64,
	orderer sdk.AccAddress,
	pairId uint64,
	dir types.OrderDirection,
	price sdkmath.LegacyDec,
	amt sdkmath.Int,
	orderLifespan time.Duration,
) types.Order {
	s.T().Helper()

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

	params, err := s.keeper.GetGenericParams(s.ctx, appID)
	s.Require().NoError(err)

	offerCoin = offerCoin.Add(sdk.NewCoin(offerCoin.Denom, sdkmath.LegacyNewDec(offerCoin.Amount.Int64()).Mul(params.SwapFeeRate).RoundInt()))
	s.fundAddr(orderer, sdk.NewCoins(offerCoin))

	msg := types.NewMsgLimitOrder(
		appID, orderer, pairId, dir, offerCoin, demandCoinDenom, price, amt, orderLifespan,
	)
	s.Require().NoError(msg.ValidateBasic())
	req, err := s.keeper.LimitOrder(s.ctx, msg)
	s.Require().NoError(err)
	return req
}

func (s *KeeperTestSuite) MarketOrder(
	appID uint64,
	orderer sdk.AccAddress,
	pairId uint64,
	dir types.OrderDirection,
	amt sdkmath.Int,
	orderLifespan time.Duration,
) types.Order {
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
		maxPrice := lastPrice.Mul(sdkmath.LegacyOneDec().Add(params.MaxPriceLimitRatio))
		offerCoin = sdk.NewCoin(pair.QuoteCoinDenom, amm.OfferCoinAmount(amm.Buy, maxPrice, amt))
		demandCoinDenom = pair.BaseCoinDenom
	case types.OrderDirectionSell:
		offerCoin = sdk.NewCoin(pair.BaseCoinDenom, amt)
		demandCoinDenom = pair.QuoteCoinDenom
	}

	offerCoin = offerCoin.Add(sdk.NewCoin(offerCoin.Denom, sdkmath.LegacyNewDec(offerCoin.Amount.Int64()).Mul(params.SwapFeeRate).RoundInt()))
	s.fundAddr(orderer, sdk.NewCoins(offerCoin))

	msg := types.NewMsgMarketOrder(
		appID, orderer, pairId, dir, offerCoin, demandCoinDenom, amt, orderLifespan,
	)
	s.Require().NoError(msg.ValidateBasic())
	req, err := s.keeper.MarketOrder(s.ctx, msg)
	s.Require().NoError(err)
	return req
}

func (s *KeeperTestSuite) MarketMakingOrder(
	orderer sdk.AccAddress, appID, pairId uint64,
	maxSellPrice, minSellPrice sdkmath.LegacyDec, sellAmt sdkmath.Int,
	maxBuyPrice, minBuyPrice sdkmath.LegacyDec, buyAmt sdkmath.Int,
	orderLifespan time.Duration, fund bool) []types.Order {
	s.T().Helper()
	params, err := s.keeper.GetGenericParams(s.ctx, appID)
	s.Require().NoError(err)

	pair, found := s.keeper.GetPair(s.ctx, appID, pairId)
	s.Require().True(found)

	maxNumTicks := int(params.MaxNumMarketMakingOrderTicks)
	tickPrec := int(params.TickPrecision)

	var buyTicks, sellTicks []types.MMOrderTick
	offerBaseCoin := sdk.NewInt64Coin(pair.BaseCoinDenom, 0)
	offerQuoteCoin := sdk.NewInt64Coin(pair.QuoteCoinDenom, 0)
	if buyAmt.IsPositive() {
		buyTicks = types.MMOrderTicks(
			types.OrderDirectionBuy, minBuyPrice, maxBuyPrice, buyAmt, maxNumTicks, tickPrec)
		for _, tick := range buyTicks {
			offerQuoteCoin = offerQuoteCoin.AddAmount(tick.OfferCoinAmount)
		}
	}
	if sellAmt.IsPositive() {
		sellTicks = types.MMOrderTicks(
			types.OrderDirectionSell, minSellPrice, maxSellPrice, sellAmt, maxNumTicks, tickPrec)
		for _, tick := range sellTicks {
			offerBaseCoin = offerBaseCoin.AddAmount(tick.OfferCoinAmount)
		}
	}
	s.fundAddr(orderer, sdk.NewCoins(offerBaseCoin, offerQuoteCoin))
	msg := types.NewMsgMMOrder(
		appID,
		orderer, pairId,
		maxSellPrice, minSellPrice, sellAmt,
		maxBuyPrice, minBuyPrice, buyAmt,
		orderLifespan)
	s.Require().NoError(msg.ValidateBasic())
	orders, err := s.keeper.MMOrder(s.ctx, msg)
	s.Require().NoError(err)

	index, found := s.keeper.GetMMOrderIndex(s.ctx, orderer, appID, pairId)
	s.Require().True(found)
	s.Require().Equal(orderer.String(), index.Orderer)
	s.Require().Equal(pairId, index.PairId)
	s.Require().True(len(index.OrderIds) <= maxNumTicks*2)
	s.Require().True(len(index.OrderIds) == len(orders))
	for i, order := range orders {
		s.Require().Equal(order.Id, index.OrderIds[i])
	}
	return orders
}

func (s *KeeperTestSuite) Farm(appID, poolID uint64, farmer sdk.AccAddress, farmingCoin string) {
	msg := types.NewMsgFarm(
		appID, poolID, farmer, utils.ParseCoin(farmingCoin),
	)
	s.fundAddr(farmer, sdk.NewCoins(msg.FarmingPoolCoin))
	err := s.keeper.Farm(s.ctx, msg)
	s.Require().NoError(err)
}
