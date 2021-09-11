package keeper

import (
	"github.com/comdex-official/comdex/x/asset/types"
	"github.com/stretchr/testify/require"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *KeeperTestSuite) TestSetAndGetPairId()  {
	suite.SetupTest()

	expectedPairId := uint64(0)

	//test get pair id when no pair id saved yet
	actualPairId := suite.assetKeeper.GetPairID(suite.ctx)
	require.Equal(suite.T(), expectedPairId, actualPairId)

	expectedPairId = 2

	suite.assetKeeper.SetPairID(suite.ctx, expectedPairId)
	actualPairId = suite.assetKeeper.GetPairID(suite.ctx)
	require.Equal(suite.T(), expectedPairId, actualPairId)
}

func (suite *KeeperTestSuite) TestSetAndGetPair() {
	suite.SetupTest()

	expectedPairs := []types.Pair{
		{
			Id:       	 10,
			AssetIn:     13,
			AssetOut:    42,
			LiquidationRatio: sdk.NewDec(12),
		},
		{
			Id:       	 42,
			AssetIn:     99,
			AssetOut:    24,
			LiquidationRatio: sdk.NewDec(20),
		},
	}

	//test without any pair saved in the store
	actualPair, found := suite.assetKeeper.GetPair(suite.ctx, expectedPairs[0].Id)
	actualPairs := suite.assetKeeper.GetPairs(suite.ctx)
	require.False(suite.T(), found)
	require.Equal(suite.T(), types.Pair{}, actualPair)
	require.Equal(suite.T(), cap(actualPairs), 0)

	//set pairs
	for _, pair := range expectedPairs {
		suite.assetKeeper.SetPair(suite.ctx, pair)
	}

	//test that the pairs exist
	for _, pair := range expectedPairs {
		actualPair, found = suite.assetKeeper.GetPair(suite.ctx, pair.Id)
		require.Truef(suite.T(), found, "GetPair found returns false when pair %d exists.", pair.Id)
		require.Equalf(suite.T(), pair, actualPair, "GetPair returns an invalid pair for id %d.", pair.Id)
	}

	//test get pairs return all pairs
	actualPairs = suite.assetKeeper.GetPairs(suite.ctx)
	require.Equal(suite.T(), expectedPairs, actualPairs)
}
