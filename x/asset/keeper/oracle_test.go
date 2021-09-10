package keeper

import (
	"github.com/comdex-official/comdex/x/asset/types"
	"github.com/stretchr/testify/require"
)

func (suite *KeeperTestSuite) TestSetAndGetMarkets() {
	suite.SetupTest()

	expectedMarkets := []types.Market{
		{
			Symbol: "ABCD",
			ScriptID: 5,
		}, {
			Symbol: "EFGH",
			ScriptID: 10,
		},
	}

	// test that no market set in the store yet
	hasMarket := suite.assetKeeper.HasMarket(suite.ctx, expectedMarkets[0].Symbol)
	actualMarket, found := suite.assetKeeper.GetMarket(suite.ctx, expectedMarkets[0].Symbol)
	actualMarkets := suite.assetKeeper.GetMarkets(suite.ctx)
	require.False(suite.T(), hasMarket, "HasMarket returns true when no market exists.")
	require.False(suite.T(), found, "GetMarket found returns true when no market exists.")
	require.Equal(suite.T(), types.Market{}, actualMarket, "GetMarket returns a market when no market exists.")
	require.Equalf(suite.T(), cap(actualMarkets), 0, "More than 0 markets found when none were saved")

	//set markets
	for _, market := range expectedMarkets {
		suite.assetKeeper.SetMarket(suite.ctx, market)
	}

	//test that the set markets exist
	for _, market := range expectedMarkets {
		hasMarket = suite.assetKeeper.HasMarket(suite.ctx, market.Symbol)
		actualMarket, found = suite.assetKeeper.GetMarket(suite.ctx, market.Symbol)
		require.Truef(suite.T(), hasMarket, "HasMarket returns false when market %s exists.", market.Symbol)
		require.Truef(suite.T(), found, "GetMarket found returns false when market %s exists.", market.Symbol)
		require.Equalf(suite.T(), market, actualMarket, "GetMarket returns an invalid market for symbol %s.", market.Symbol)
	}

	//test get markets
	actualMarkets = suite.assetKeeper.GetMarkets(suite.ctx)
	require.Equalf(suite.T(), expectedMarkets, actualMarkets, "Markets returned do not match the saved markets.")
}

func (suite *KeeperTestSuite) TestSetAndGetPriceForMarkets() {
	suite.SetupTest()

	symbol := "ABCD"
	expectedPrice := uint64(0)

	actualPrice, found := suite.assetKeeper.GetPriceForMarket(suite.ctx, symbol)
	require.False(suite.T(), found, "price found for market when none was saved.")
	require.Equal(suite.T(), expectedPrice, actualPrice, "market prices do not match.")

	expectedPrice = 42
	suite.assetKeeper.SetPriceForMarket(suite.ctx, symbol, expectedPrice)
	actualPrice, found = suite.assetKeeper.GetPriceForMarket(suite.ctx, symbol)
	require.True(suite.T(), found, "price not found for the saved market.")
	require.Equal(suite.T(), expectedPrice, actualPrice, "market prices do not match.")
}

func (suite *KeeperTestSuite) TestSetAndGetMarketForAsset() {
	suite.SetupTest()

	assetId := uint64(42)
	expectedMarket := types.Market{
		Symbol:   "ABCD",
		ScriptID: 1,
	}

	validateMarketForAssetDoesNotExist := func(expectedAssetId uint64) {
		hasMarketForAsset := suite.assetKeeper.HasMarketForAsset(suite.ctx, expectedAssetId)
		actualMarketForAsset, found := suite.assetKeeper.GetMarketForAsset(suite.ctx, expectedAssetId)
		require.False(suite.T(), hasMarketForAsset, "HasMarketForAsset returns true when market for asset does not exist.")
		require.False(suite.T(), found, "GetMarketForAsset found returns true when market for asset does not exist.")
		require.Equalf(suite.T(), types.Market{}, actualMarketForAsset, "GetMarketForAsset returns a market when no market should exist.")
	}

	validateMarketForAssetExists := func(expectedAssetId uint64) {
		hasMarketForAsset := suite.assetKeeper.HasMarketForAsset(suite.ctx, expectedAssetId)
		actualMarketForAsset, found := suite.assetKeeper.GetMarketForAsset(suite.ctx, expectedAssetId)
		require.True(suite.T(), hasMarketForAsset, "HasMarketForAsset returns false when market for asset exists.")
		require.True(suite.T(), found, "GetMarketForAsset found returns false when market for asset exists.")
		require.Equalf(suite.T(), expectedMarket, actualMarketForAsset, "Markets do not match for asset id %d", assetId)
	}

	// test get market for asset when no market for asset exists yet
	validateMarketForAssetDoesNotExist(assetId)

	//set market for asset
	suite.assetKeeper.SetMarketForAsset(suite.ctx, assetId, expectedMarket.Symbol)

	// validate asset for market was set but market does not exist yet
	hasMarketForAsset := suite.assetKeeper.HasMarketForAsset(suite.ctx, assetId)
	actualMarketForAsset, found := suite.assetKeeper.GetMarketForAsset(suite.ctx, assetId)
	require.True(suite.T(), hasMarketForAsset, "HasMarketForAsset returns false when market for asset exists.")
	require.False(suite.T(), found, "GetMarketForAsset found returns true when no market exists.")
	require.Equalf(suite.T(), types.Market{}, actualMarketForAsset, "GetMarketForAsset returns a market when no market should exist.")

	// set market
	suite.assetKeeper.SetMarket(suite.ctx, expectedMarket)
	validateMarketForAssetExists(assetId)

	// delete market for asset
	suite.assetKeeper.DeleteMarketForAsset(suite.ctx, assetId)
	validateMarketForAssetDoesNotExist(assetId)
}

func (suite *KeeperTestSuite) TestGetPriceForAsset() {}