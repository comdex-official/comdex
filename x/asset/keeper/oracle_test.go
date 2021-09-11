package keeper

import (
	"github.com/bandprotocol/bandchain-packet/obi"
	bandpacket "github.com/bandprotocol/bandchain-packet/packet"
	"github.com/comdex-official/comdex/x/asset/types"
	"github.com/stretchr/testify/require"
	"strconv"
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
	require.False(suite.T(), hasMarket)
	require.False(suite.T(), found)
	require.Equal(suite.T(), types.Market{}, actualMarket)
	require.Equal(suite.T(), cap(actualMarkets), 0)

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
	require.Equal(suite.T(), expectedMarkets, actualMarkets)
}

func (suite *KeeperTestSuite) TestSetAndGetPriceForMarkets() {
	suite.SetupTest()

	symbol := "ABCD"
	expectedPrice := uint64(0)

	actualPrice, found := suite.assetKeeper.GetPriceForMarket(suite.ctx, symbol)
	require.False(suite.T(), found)
	require.Equal(suite.T(), expectedPrice, actualPrice)

	expectedPrice = 42
	suite.assetKeeper.SetPriceForMarket(suite.ctx, symbol, expectedPrice)
	actualPrice, found = suite.assetKeeper.GetPriceForMarket(suite.ctx, symbol)
	require.True(suite.T(), found)
	require.Equal(suite.T(), expectedPrice, actualPrice)
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
		require.False(suite.T(), hasMarketForAsset)
		require.False(suite.T(), found)
		require.Equalf(suite.T(), types.Market{}, actualMarketForAsset, "market found when no market should exist.")
	}

	validateMarketForAssetExists := func(expectedAssetId uint64) {
		hasMarketForAsset := suite.assetKeeper.HasMarketForAsset(suite.ctx, expectedAssetId)
		actualMarketForAsset, found := suite.assetKeeper.GetMarketForAsset(suite.ctx, expectedAssetId)
		require.True(suite.T(), hasMarketForAsset)
		require.True(suite.T(), found)
		require.Equalf(suite.T(), expectedMarket, actualMarketForAsset, "Markets do not match for asset id %d", assetId)
	}

	// test get market for asset when no market for asset exists yet
	validateMarketForAssetDoesNotExist(assetId)

	//set market for asset
	suite.assetKeeper.SetMarketForAsset(suite.ctx, assetId, expectedMarket.Symbol)

	// validate asset for market was set but market does not exist yet
	hasMarketForAsset := suite.assetKeeper.HasMarketForAsset(suite.ctx, assetId)
	actualMarketForAsset, found := suite.assetKeeper.GetMarketForAsset(suite.ctx, assetId)
	require.True(suite.T(), hasMarketForAsset)
	require.False(suite.T(), found)
	require.Equal(suite.T(), types.Market{}, actualMarketForAsset)

	// set market
	suite.assetKeeper.SetMarket(suite.ctx, expectedMarket)
	validateMarketForAssetExists(assetId)

	// delete market for asset
	suite.assetKeeper.DeleteMarketForAsset(suite.ctx, assetId)
	validateMarketForAssetDoesNotExist(assetId)
}

func (suite *KeeperTestSuite) TestGetPriceForAsset() {
	suite.SetupTest()

	assetId := uint64(10)
	market := types.Market{
		Symbol:   "ABCD",
		ScriptID: 0,
	}
	expectedPrice := uint64(0)

	//test with no asset and market in the store
	actualPrice, found := suite.assetKeeper.GetPriceForAsset(suite.ctx, assetId)
	require.False(suite.T(), found)
	require.Equal(suite.T(), expectedPrice, actualPrice)
	
	//add market, set market for the asset, set price for market and validate
	suite.assetKeeper.SetMarket(suite.ctx, market)
	suite.assetKeeper.SetMarketForAsset(suite.ctx, assetId, market.Symbol)
	suite.assetKeeper.SetPriceForMarket(suite.ctx, market.Symbol, expectedPrice)
	actualPrice, found = suite.assetKeeper.GetPriceForAsset(suite.ctx, assetId)
	require.True(suite.T(), found)
	require.Equal(suite.T(), expectedPrice, actualPrice)
}

func (suite * KeeperTestSuite) TestSetAndGetCalldataId()  {
	suite.SetupTest()

	expectedCalldataId := uint64(0)
	actualCalldataId := suite.assetKeeper.GetCalldataID(suite.ctx)
	require.Equal(suite.T(), expectedCalldataId, actualCalldataId)

	expectedCalldataId = 42
	suite.assetKeeper.SetCalldataID(suite.ctx, expectedCalldataId)
	actualCalldataId = suite.assetKeeper.GetCalldataID(suite.ctx)
	require.Equal(suite.T(), expectedCalldataId, actualCalldataId)
}

func (suite *KeeperTestSuite) TestSetAndGetCalldata() {
	suite.SetupTest()

	callDataKey := uint64(69)
	expectedCalldata := types.Calldata{
		Symbols:    []string{"ABCD", "EFGH"},
		Multiplier: 3,
	}

	actualCalldata, found := suite.assetKeeper.GetCalldata(suite.ctx, callDataKey)
	require.False(suite.T(), found)
	require.Equal(suite.T(), types.Calldata{}, actualCalldata)

	suite.assetKeeper.SetCalldata(suite.ctx, callDataKey, expectedCalldata)
	actualCalldata, found = suite.assetKeeper.GetCalldata(suite.ctx, callDataKey)
	require.True(suite.T(), found)
	require.Equal(suite.T(), expectedCalldata, actualCalldata)

	suite.assetKeeper.DeleteCalldata(suite.ctx, callDataKey)
	actualCalldata, found = suite.assetKeeper.GetCalldata(suite.ctx, callDataKey)
	require.False(suite.T(), found)
	require.Equal(suite.T(), types.Calldata{}, actualCalldata)
}

func (suite *KeeperTestSuite) TestOnRecvPacket()  {
	suite.SetupTest()

	symbolsAndPrices := struct{
		Symbols []string
		marketPrices []uint64
	} {
		Symbols: []string{"ABCD", "EFGH"},
		marketPrices: []uint64{42, 13},
	}
	calldataId := uint64(12)
	calldata := types.Calldata{
		Symbols:    []string{symbolsAndPrices.Symbols[0], symbolsAndPrices.Symbols[1]},
		Multiplier: 2,
	}

	result := types.Result{
		Rates: []uint64{42, 13},
	}
	encodedResult, err := obi.Encode(result)
	require.NoError(suite.T(), err)

	validateSymbolsAndMarketPrices := func(symbols []string, prices []uint64) {
		for i, symbol := range symbols {
			price, found := suite.assetKeeper.GetPriceForMarket(suite.ctx, symbol)
			require.True(suite.T(), found)
			require.Equalf(suite.T(), prices[i], price, "found invalid market price for symbol %s", symbol)
		}
	}

	//initialize market prices for symbols
	for i, symbol := range symbolsAndPrices.Symbols {
		suite.assetKeeper.SetPriceForMarket(suite.ctx, symbol, symbolsAndPrices.marketPrices[i])
	}

	//validate OnRecvPacket throws on invalid clientId
	oracleResponsePacketData := bandpacket.NewOracleResponsePacketData("abc", 123, 3, 1257894000, 1257895000, bandpacket.RESOLVE_STATUS_SUCCESS, encodedResult)
	err = suite.assetKeeper.OnRecvPacket(suite.ctx, oracleResponsePacketData)
	require.Error(suite.T(), err)

	//set a valid id in oracle response
	oracleResponsePacketData.ClientID = strconv.FormatUint(calldataId, 10)

	//set the oracle result status to failure
	oracleResponsePacketData.ResolveStatus = bandpacket.RESOLVE_STATUS_FAILURE
	oracleResponsePacketData.Result = nil

	//set calldata
	suite.assetKeeper.SetCalldata(suite.ctx, calldataId, calldata)

	//validate calldata is deleted on orcale response failure status and symbol prices do not change
	err = suite.assetKeeper.OnRecvPacket(suite.ctx, oracleResponsePacketData)
	require.NoError(suite.T(), err)
	validateSymbolsAndMarketPrices(symbolsAndPrices.Symbols, symbolsAndPrices.marketPrices)
	_, found := suite.assetKeeper.GetCalldata(suite.ctx, calldataId)
	require.False(suite.T(), found)

	//set the oracle result status to success
	oracleResponsePacketData.ResolveStatus = bandpacket.RESOLVE_STATUS_SUCCESS
	oracleResponsePacketData.Result = encodedResult

	//validate OnRecvPacket throws when no calldata found
	err = suite.assetKeeper.OnRecvPacket(suite.ctx, oracleResponsePacketData)
	require.Error(suite.T(), err)
	validateSymbolsAndMarketPrices(symbolsAndPrices.Symbols, symbolsAndPrices.marketPrices)

	//set calldata
	suite.assetKeeper.SetCalldata(suite.ctx, calldataId, calldata)

	//validate symbol prices are updated and calldata is deleted
	err = suite.assetKeeper.OnRecvPacket(suite.ctx, oracleResponsePacketData)
	require.NoError(suite.T(), err)
	validateSymbolsAndMarketPrices(symbolsAndPrices.Symbols, result.Rates)
	_, found = suite.assetKeeper.GetCalldata(suite.ctx, calldataId)
	require.False(suite.T(), found)
}