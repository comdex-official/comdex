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
	require.Equalf(suite.T(), expectedAssetId, actualAssetId, "asset ids do not match.")

	expectedAssetId = 2

	suite.assetKeeper.SetAssetID(suite.ctx, expectedAssetId)
	actualAssetId = suite.assetKeeper.GetAssetID(suite.ctx)
	require.Equalf(suite.T(), expectedAssetId, actualAssetId, "asset ids do not match.")
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
	require.False(suite.T(), hasAsset, "HasAsset returns true when no asset.")
	require.False(suite.T(), found, "GetAsset found returns true when no asset.")
	require.Equal(suite.T(), types.Asset{}, actualAsset, "Get Asset returns an asset when no asset.")
	require.Equalf(suite.T(), cap(actualAssets), 0, "More than 0 assets found when none were saved")

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
	require.Equalf(suite.T(), expectedAssets, actualAssets, "Assets returned do not match the saved assets.")
}

func (suite *KeeperTestSuite) TestSetAndGetAssetForDenom() {
	suite.SetupTest()

	expectedAsset := types.Asset{
		Id:       5,
		Name:     "Test Asset",
		Denom:    "cmdx",
		Decimals: 5,
	}

	validateAssetForDenomDoesNotExist := func(expectedDenom string) {
		actualAssetForDenom, found := suite.assetKeeper.GetAssetForDenom(suite.ctx, expectedDenom)
		hasAssetForDenom := suite.assetKeeper.HasAssetForDenom(suite.ctx, expectedDenom)
		require.False(suite.T(), found, "GetAssetForDenom found returns true when no asset.")
		require.False(suite.T(), hasAssetForDenom, "HasAssetForDenom for denom returns true when no asset.")
		require.Equal(suite.T(), types.Asset{}, actualAssetForDenom, "found an asset for denom when no asset was saved in the store.")
	}

	validateAssetForDenomExists := func(expectedDenom string) {
		hasAssetForDenom := suite.assetKeeper.HasAssetForDenom(suite.ctx, expectedDenom)
		actualAssetForDenom, found := suite.assetKeeper.GetAssetForDenom(suite.ctx, expectedDenom)
		require.True(suite.T(), found, "GetAssetForDenom found returns false when corresponding asset exists.")
		require.True(suite.T(), hasAssetForDenom, "HasAssetForDenom for denom returns false when corresponding asset exists.")
		require.Equal(suite.T(), expectedAsset, actualAssetForDenom, "Assets do not match.")
	}

	//test when no asset for denom exist in the store yet
	validateAssetForDenomDoesNotExist(expectedAsset.Denom)

	//save asset for denom
	suite.assetKeeper.SetAssetForDenom(suite.ctx, expectedAsset.Denom, expectedAsset.Id)

	//test get asset when no corresponding asset for denom exists yet
	hasAssetForDenom := suite.assetKeeper.HasAssetForDenom(suite.ctx, expectedAsset.Denom)
	actualAssetForDenom, found := suite.assetKeeper.GetAssetForDenom(suite.ctx, expectedAsset.Denom)
	require.False(suite.T(), found, "GetAssetForDenom found returns true when no asset.")
	require.True(suite.T(), hasAssetForDenom, "HasAssetForDenom for denom returns false when asset exists.")
	require.Equal(suite.T(), types.Asset{}, actualAssetForDenom, "found an invalid asset when no asset was saved in the store.")

	//save the corresponding asset
	suite.assetKeeper.SetAsset(suite.ctx, expectedAsset)
	validateAssetForDenomExists(expectedAsset.Denom)

	// delete asset for denom
	suite.assetKeeper.DeleteAssetForDenom(suite.ctx, expectedAsset.Denom)
	validateAssetForDenomDoesNotExist(expectedAsset.Denom)
}