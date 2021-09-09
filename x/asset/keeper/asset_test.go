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
	require.Equalf(suite.T(), actualAssetId, expectedAssetId, "asset ids do not match.")

	expectedAssetId = uint64(2)

	suite.assetKeeper.SetAssetID(suite.ctx, expectedAssetId)
	actualAssetId = suite.assetKeeper.GetAssetID(suite.ctx)
	require.Equalf(suite.T(), actualAssetId, expectedAssetId, "asset ids do not match.")
}

func (suite *KeeperTestSuite) TestSetAndGetAsset() {
	suite.SetupTest()

	expectedAsset := types.Asset{
		Id:       uint64(42),
		Name:     "Test Name",
		Denom:    "cmdx",
		Decimals: 4,
	}

	//test without any saved in the store
	hasAsset := suite.assetKeeper.HasAsset(suite.ctx, expectedAsset.Id)
	actualAsset, found := suite.assetKeeper.GetAsset(suite.ctx, expectedAsset.Id)
	require.False(suite.T(), hasAsset, "HasAsset returns true when no asset.")
	require.False(suite.T(), found, "GetAsset found returns true when no asset.")
	require.Equal(suite.T(), types.Asset{}, actualAsset, "Get Asset returns an asset when no asset.")

	suite.assetKeeper.SetAsset(suite.ctx, expectedAsset)

	hasAsset = suite.assetKeeper.HasAsset(suite.ctx, expectedAsset.Id)
	actualAsset, found = suite.assetKeeper.GetAsset(suite.ctx, expectedAsset.Id)
	require.True(suite.T(), hasAsset, "HasAsset returns false after setting an asset.")
	require.True(suite.T(), found, "GetAsset found returns false when no asset.")
	require.Equal(suite.T(), expectedAsset, actualAsset, "Get Asset returns an invalid asset.")
}

func (suite *KeeperTestSuite) TestSetAndGetAssets() {
	suite.SetupTest()

	expectedAssets := []types.Asset{
		{
			Id:       uint64(10),
			Name:     "assetName1",
			Denom:    "assetDenom1",
			Decimals: 2,
		},
		{
			Id:       uint64(42),
			Name:     "assetName2",
			Denom:    "assetDenom2",
			Decimals: 5,
		},
	}

	//test without any saved asset
	actualAssets := suite.assetKeeper.GetAssets(suite.ctx)
	require.Equalf(suite.T(), cap(actualAssets), 0, "More than 0 assets found when none were saved")

	for _, asset := range expectedAssets {
		suite.assetKeeper.SetAsset(suite.ctx, asset)
	}

	actualAssets = suite.assetKeeper.GetAssets(suite.ctx)
	require.Equalf(suite.T(), actualAssets, expectedAssets, "Assets returned do not match the saved assets")
}

func (suite *KeeperTestSuite) TestSetAndGetAssetForDenom() {
	suite.SetupTest()

	expectedAsset := types.Asset{
		Id:       uint64(5),
		Name:     "Test Asset",
		Denom:    "cmdx",
		Decimals: 5,
	}

	//test when no asset for denom exist in the store yet
	actualAssetForDenom, found := suite.assetKeeper.GetAssetForDenom(suite.ctx, expectedAsset.Denom)
	hasAssetForDenom := suite.assetKeeper.HasAssetForDenom(suite.ctx, expectedAsset.Denom)
	require.False(suite.T(), found, "GetAssetForDenom found returns true when no asset.")
	require.False(suite.T(), hasAssetForDenom, "HasAssetForDenom for denom returns true when no asset.")
	require.Equal(suite.T(), types.Asset{}, actualAssetForDenom, "found an asset for denom when no asset was saved in the store.")

	//save asset for denom
	suite.assetKeeper.SetAssetForDenom(suite.ctx, expectedAsset.Denom, expectedAsset.Id)

	//test get asset when no corresponding asset for denom exists yet
	hasAssetForDenom = suite.assetKeeper.HasAssetForDenom(suite.ctx, expectedAsset.Denom)
	actualAssetForDenom, found = suite.assetKeeper.GetAssetForDenom(suite.ctx, expectedAsset.Denom)
	require.False(suite.T(), found, "GetAssetForDenom found returns true when no asset.")
	require.True(suite.T(), hasAssetForDenom, "HasAssetForDenom for denom returns false when no corresponding asset.")
	require.Equal(suite.T(), types.Asset{}, actualAssetForDenom, "found an invalid asset when no asset was saved in the store.")

	//save the corresponding asset
	suite.assetKeeper.SetAsset(suite.ctx, expectedAsset)

	//test when both asset and corresponding denom keys ares stored
	hasAssetForDenom = suite.assetKeeper.HasAssetForDenom(suite.ctx, expectedAsset.Denom)
	actualAssetForDenom, found = suite.assetKeeper.GetAssetForDenom(suite.ctx, expectedAsset.Denom)
	require.True(suite.T(), found, "GetAssetForDenom found returns false when corresponding asset exists.")
	require.True(suite.T(), hasAssetForDenom, "HasAssetForDenom for denom returns false when no asset.")
	require.Equal(suite.T(), expectedAsset, actualAssetForDenom, "Assets do not match.")

	// delete asset for denom
	suite.assetKeeper.DeleteAssetForDenom(suite.ctx, expectedAsset.Denom)

	//test denom key is gone from the store but the corresponding asset still remains
	hasAssetForDenom = suite.assetKeeper.HasAssetForDenom(suite.ctx, expectedAsset.Denom)
	actualAssetForDenom, found = suite.assetKeeper.GetAssetForDenom(suite.ctx, expectedAsset.Denom)
	actualAssetFound := suite.assetKeeper.HasAsset(suite.ctx, expectedAsset.Id)
	require.False(suite.T(), found, "GetAssetForDenom found returns true when no asset.")
	require.False(suite.T(), hasAssetForDenom, "HasAssetForDenom for denom returns true when no asset.")
	require.True(suite.T(), actualAssetFound, "HasAsset returns false when asset exists.")
	require.Equal(suite.T(), types.Asset{}, actualAssetForDenom, "found an invalid asset.")
}