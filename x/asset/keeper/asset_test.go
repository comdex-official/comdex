package keeper

import (
	"github.com/comdex-official/comdex/x/asset/types"
	"github.com/stretchr/testify/require"
)

func (suite *KeeperTestSuite) TestSetAndGetAssetID() {
	suite.SetupTest()

	expectedAssetId := uint64(0)

	//test get asset when no asset id saved yet
	actualAssetId := suite.assetKeeper.GetAssetID(suite.ctx)
	require.Equal(suite.T(), expectedAssetId, actualAssetId)

	expectedAssetId = 2

	suite.assetKeeper.SetAssetID(suite.ctx, expectedAssetId)
	actualAssetId = suite.assetKeeper.GetAssetID(suite.ctx)
	require.Equal(suite.T(), expectedAssetId, actualAssetId)
}

func (suite *KeeperTestSuite) TestSetAndGetAsset() {
	suite.SetupTest()

	expectedAssets := []types.Asset{
		{
			Id:       10,
			Name:     "assetName1",
			Denom:    "assetDenom1",
			Decimals: 2,
		},
		{
			Id:       42,
			Name:     "assetName2",
			Denom:    "assetDenom2",
			Decimals: 5,
		},
	}

	//test without any saved in the store
	hasAsset := suite.assetKeeper.HasAsset(suite.ctx, expectedAssets[0].Id)
	actualAsset, found := suite.assetKeeper.GetAsset(suite.ctx, expectedAssets[0].Id)
	actualAssets := suite.assetKeeper.GetAssets(suite.ctx)
	require.False(suite.T(), hasAsset)
	require.False(suite.T(), found)
	require.Equal(suite.T(), types.Asset{}, actualAsset)
	require.Equal(suite.T(), cap(actualAssets), 0)

	//set assets
	for _, asset := range expectedAssets {
		suite.assetKeeper.SetAsset(suite.ctx, asset)
	}

	//test that the set assets exist
	for _, asset := range expectedAssets {
		hasAsset = suite.assetKeeper.HasAsset(suite.ctx, asset.Id)
		actualAsset, found = suite.assetKeeper.GetAsset(suite.ctx, asset.Id)
		require.Truef(suite.T(), hasAsset, "HasAsset returns false when asset %d exists.", asset.Id)
		require.Truef(suite.T(), found, "GetAsset found returns false when asset %d exists.", asset.Id)
		require.Equalf(suite.T(), asset, actualAsset, "GetAsset returns an invalid asset for id %d.", asset.Id)
	}

	//test get markets
	actualAssets = suite.assetKeeper.GetAssets(suite.ctx)
	require.Equal(suite.T(), expectedAssets, actualAssets)
}

func (suite *KeeperTestSuite) TestSetAndGetAssetForDenom() {
	suite.SetupTest()

	expectedAsset := types.Asset{
		Id:       5,
		Name:     "Test Asset",
		Denom:    "cmdx",
		Decimals: 5,
	}

	validateAssetForDenomDoesNotExist := func(denom string) {
		actualAssetForDenom, found := suite.assetKeeper.GetAssetForDenom(suite.ctx, denom)
		hasAssetForDenom := suite.assetKeeper.HasAssetForDenom(suite.ctx, denom)
		require.Falsef(suite.T(), found, "GetAssetForDenom found returns true when no asset for denom %s.", denom)
		require.Falsef(suite.T(), hasAssetForDenom, "HasAssetForDenom for denom returns true when no asset for denom %s.", denom)
		require.Equalf(suite.T(), types.Asset{}, actualAssetForDenom, "found an asset for denom %s when no asset was saved in the store.", denom)
	}

	validateAssetForDenomExists := func(denom string) {
		hasAssetForDenom := suite.assetKeeper.HasAssetForDenom(suite.ctx, denom)
		actualAssetForDenom, found := suite.assetKeeper.GetAssetForDenom(suite.ctx, denom)
		require.Truef(suite.T(), found, "Asset for denom %s not found.", denom)
		require.Truef(suite.T(), hasAssetForDenom, "Asset for denom %s not found.", denom)
		require.Equalf(suite.T(), expectedAsset, actualAssetForDenom, "Assets do not match for denom %s.", denom)
	}

	//test when no asset for denom exist in the store yet
	validateAssetForDenomDoesNotExist(expectedAsset.Denom)

	//save asset for denom
	suite.assetKeeper.SetAssetForDenom(suite.ctx, expectedAsset.Denom, expectedAsset.Id)

	//test get asset when no corresponding asset for denom exists yet
	hasAssetForDenom := suite.assetKeeper.HasAssetForDenom(suite.ctx, expectedAsset.Denom)
	actualAssetForDenom, found := suite.assetKeeper.GetAssetForDenom(suite.ctx, expectedAsset.Denom)
	require.False(suite.T(), found)
	require.True(suite.T(), hasAssetForDenom)
	require.Equal(suite.T(), types.Asset{}, actualAssetForDenom)

	//save the corresponding asset
	suite.assetKeeper.SetAsset(suite.ctx, expectedAsset)
	validateAssetForDenomExists(expectedAsset.Denom)

	// delete asset for denom
	suite.assetKeeper.DeleteAssetForDenom(suite.ctx, expectedAsset.Denom)
	validateAssetForDenomDoesNotExist(expectedAsset.Denom)
}