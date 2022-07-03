package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	assettypes "github.com/comdex-official/comdex/x/asset/types"

	_ "github.com/stretchr/testify/suite"
)

func NewSubApp(s *KeeperTestSuite, appName string) error {
	return s.app.AssetKeeper.AddAppRecords(s.ctx, assettypes.AppData{
		Name:             appName,
		ShortName:        appName,
		MinGovDeposit:    sdk.NewInt(0),
		GovTimeInSeconds: 0,
		GenesisToken:     []assettypes.MintGenesisToken{},
	})
}

func GetAppIDByAppName(s *KeeperTestSuite, appName string) uint64 {
	apps, found := s.app.AssetKeeper.GetApps(s.ctx)
	s.Require().True(found)

	var appID uint64
	for _, app := range apps {
		if app.Name == appName {
			appID = app.Id
			break
		}
	}
	return appID
}

func NewAddAsset(s *KeeperTestSuite, name, denom string) error {
	return s.app.AssetKeeper.AddAssetRecords(s.ctx, assettypes.Asset{
		Name:                  name,
		Denom:                 denom,
		Decimals:              1000000,
		IsOnChain:             true,
		IsOraclePriceRequired: true,
	})
}

func GetAssetIDByDenom(s *KeeperTestSuite, denom string) uint64 {
	assets := s.app.AssetKeeper.GetAssets(s.ctx)

	var assetID uint64
	for _, asset := range assets {
		if asset.Denom == denom {
			assetID = asset.Id
			break
		}
	}
	return assetID
}
