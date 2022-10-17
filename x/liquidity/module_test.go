package liquidity_test

import (
	"encoding/binary"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	chain "github.com/comdex-official/comdex/app"
	utils "github.com/comdex-official/comdex/types"
	assettypes "github.com/comdex-official/comdex/x/asset/types"
	"github.com/comdex-official/comdex/x/liquidity"
	"github.com/comdex-official/comdex/x/liquidity/amm"
	"github.com/comdex-official/comdex/x/liquidity/keeper"
	"github.com/comdex-official/comdex/x/liquidity/types"
	markettypes "github.com/comdex-official/comdex/x/market/types"
)

type ModuleTestSuite struct {
	suite.Suite

	app     *chain.App
	ctx     sdk.Context
	keeper  keeper.Keeper
	querier keeper.Querier
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
	suite.keeper = suite.app.LiquidityKeeper
	suite.querier = keeper.Querier{Keeper: suite.keeper}
}

// Below are useful helpers to write test code easily.
func (suite *ModuleTestSuite) addr(addrNum int) sdk.AccAddress {
	addr := make(sdk.AccAddress, 20)
	binary.PutVarint(addr, int64(addrNum))
	return addr
}

func (s *ModuleTestSuite) nextBlock() {
	liquidity.EndBlocker(s.ctx, s.keeper)
	liquidity.BeginBlocker(s.ctx, s.keeper)
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

func (s *ModuleTestSuite) CreateNewAsset(name, denom string, price uint64) assettypes.Asset {
	err := s.app.AssetKeeper.AddAssetRecords(s.ctx, assettypes.Asset{
		Name:                  name,
		Denom:                 denom,
		Decimals:              1000000,
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

func (s *ModuleTestSuite) CreateNewLiquidityPair(appID uint64, creator sdk.AccAddress, baseCoinDenom, quoteCoinDenom string) types.Pair {
	params, err := s.keeper.GetGenericParams(s.ctx, appID)
	s.Require().NoError(err)

	s.fundAddr(creator, params.PairCreationFee)

	msg := types.NewMsgCreatePair(appID, creator, baseCoinDenom, quoteCoinDenom)
	pair, err := s.keeper.CreatePair(s.ctx, msg, false)

	s.Require().NoError(err)
	s.Require().IsType(types.Pair{}, pair)

	return pair
}

func (s *ModuleTestSuite) CreateNewLiquidityPool(appID, pairID uint64, creator sdk.AccAddress, depositCoins string) types.Pool {
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

func (s *ModuleTestSuite) Deposit(appID, poolID uint64, depositor sdk.AccAddress, depositCoins string) types.DepositRequest {
	msg := types.NewMsgDeposit(
		appID, depositor, poolID, utils.ParseCoins(depositCoins),
	)
	s.fundAddr(depositor, msg.DepositCoins)
	req, err := s.keeper.Deposit(s.ctx, msg)
	s.Require().NoError(err)
	s.Require().IsType(types.DepositRequest{}, req)
	return req
}

func (s *ModuleTestSuite) Withdraw(appID, poolID uint64, withdrawer sdk.AccAddress, poolCoin sdk.Coin) types.WithdrawRequest {
	msg := types.NewMsgWithdraw(
		appID, withdrawer, poolID, poolCoin,
	)
	req, err := s.keeper.Withdraw(s.ctx, msg)
	s.Require().NoError(err)
	s.Require().IsType(types.WithdrawRequest{}, req)
	return req
}

func (s *ModuleTestSuite) LimitOrder(
	appID uint64,
	orderer sdk.AccAddress,
	pairId uint64,
	dir types.OrderDirection,
	price sdk.Dec,
	amt sdk.Int,
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

	offerCoin = offerCoin.Add(sdk.NewCoin(offerCoin.Denom, offerCoin.Amount.ToDec().Mul(params.SwapFeeRate).RoundInt()))
	s.fundAddr(orderer, sdk.NewCoins(offerCoin))

	msg := types.NewMsgLimitOrder(
		appID, orderer, pairId, dir, offerCoin, demandCoinDenom, price, amt, orderLifespan,
	)
	s.Require().NoError(msg.ValidateBasic())
	req, err := s.keeper.LimitOrder(s.ctx, msg)
	s.Require().NoError(err)
	return req
}

func (s *ModuleTestSuite) MarketOrder(
	appID uint64,
	orderer sdk.AccAddress,
	pairId uint64,
	dir types.OrderDirection,
	amt sdk.Int,
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
		maxPrice := lastPrice.Mul(sdk.OneDec().Add(params.MaxPriceLimitRatio))
		offerCoin = sdk.NewCoin(pair.QuoteCoinDenom, amm.OfferCoinAmount(amm.Buy, maxPrice, amt))
		demandCoinDenom = pair.BaseCoinDenom
	case types.OrderDirectionSell:
		offerCoin = sdk.NewCoin(pair.BaseCoinDenom, amt)
		demandCoinDenom = pair.QuoteCoinDenom
	}

	offerCoin = offerCoin.Add(sdk.NewCoin(offerCoin.Denom, offerCoin.Amount.ToDec().Mul(params.SwapFeeRate).RoundInt()))
	s.fundAddr(orderer, sdk.NewCoins(offerCoin))

	msg := types.NewMsgMarketOrder(
		appID, orderer, pairId, dir, offerCoin, demandCoinDenom, amt, orderLifespan,
	)
	s.Require().NoError(msg.ValidateBasic())
	req, err := s.keeper.MarketOrder(s.ctx, msg)
	s.Require().NoError(err)
	return req
}
