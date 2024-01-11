package keeper_test

import (
	"github.com/comdex-official/comdex/app/wasm/bindings"
	assetTypes "github.com/comdex-official/comdex/x/asset/types"
	liquidationTypes "github.com/comdex-official/comdex/x/liquidation/types"
	markettypes "github.com/comdex-official/comdex/x/market/types"
	vaultKeeper1 "github.com/comdex-official/comdex/x/vault/keeper"
	vaultTypes "github.com/comdex-official/comdex/x/vault/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (s *KeeperTestSuite) AddPairAndExtendedPairVault1() {
	assetKeeper, liquidationKeeper, ctx := &s.assetKeeper, &s.liquidationKeeper, &s.ctx

	for _, tc := range []struct {
		name              string
		pair              assetTypes.Pair
		extendedPairVault bindings.MsgAddExtendedPairsVault
		asset1            uint64
		asset2            uint64
	}{
		{
			"Add Pair , Extended Pair Vault : cmdx cmst",
			assetTypes.Pair{
				AssetIn:  1,
				AssetOut: 2,
			},
			bindings.MsgAddExtendedPairsVault{
				AppID:               1,
				PairID:              1,
				StabilityFee:        sdk.MustNewDecFromStr("0.01"),
				ClosingFee:          sdk.MustNewDecFromStr("0"),
				LiquidationPenalty:  sdk.MustNewDecFromStr("0.12"),
				DrawDownFee:         sdk.MustNewDecFromStr("0.01"),
				IsVaultActive:       true,
				DebtCeiling:         sdk.NewInt(1000000000000),
				DebtFloor:           sdk.NewInt(1000000),
				IsStableMintVault:   false,
				MinCr:               sdk.MustNewDecFromStr("1.5"),
				PairName:            "CMDX-B",
				AssetOutOraclePrice: true,
				AssetOutPrice:       1000000,
				MinUsdValueLeft:     1000000,
			},
			1,
			2,
		},
	} {
		s.Run(tc.name, func() {
			err := assetKeeper.AddPairsRecords(*ctx, tc.pair)
			s.Require().NoError(err)

			err = assetKeeper.WasmAddExtendedPairsVaultRecords(*ctx, &tc.extendedPairVault)
			s.Require().NoError(err)

			err = liquidationKeeper.WasmWhitelistAppIDLiquidation(*ctx, 1)
			s.Require().NoError(err)

			s.SetInitialOraclePriceForID(tc.asset1, tc.asset2)
		})
	}
}

func (s *KeeperTestSuite) SetOraclePrice(assetID uint64, price uint64) {
	market := markettypes.TimeWeightedAverage{
		AssetID:       assetID,
		ScriptID:      12,
		Twa:           price,
		CurrentIndex:  0,
		IsPriceActive: true,
		PriceValue:    []uint64{price},
	}
	s.app.MarketKeeper.SetTwa(s.ctx, market)
}

func (s *KeeperTestSuite) SetInitialOraclePriceForID(asset1 uint64, asset2 uint64) {
	s.SetOraclePrice(asset1, 2000000)
	s.SetOraclePrice(asset2, 1000000)
}

func (s *KeeperTestSuite) ChangeOraclePrice(asset uint64) {
	s.SetOraclePrice(asset, 1000000)
}

func (s *KeeperTestSuite) CreateVault() {
	userAddress1 := "cosmos1q7q90qsl9g0gl2zz0njxwv2a649yqrtyxtnv3v"
	userAddress2 := "cosmos1hm7w7dnvdnra78pz9qxysy7u4tuhc3fnpjmyj7"
	vaultKeeper, ctx := &s.vaultKeeper, &s.ctx
	s.AddAppAsset()
	s.AddPairAndExtendedPairVault1()
	server := vaultKeeper1.NewMsgServer(*vaultKeeper)

	for index, tc := range []struct {
		name string
		msg  vaultTypes.MsgCreateRequest
	}{
		{
			"Create Vault : AppID 1 extended pair 1 user address 1",
			vaultTypes.MsgCreateRequest{
				From:                userAddress1,
				AppId:               1,
				ExtendedPairVaultId: 1,
				AmountIn:            sdk.NewIntFromUint64(1000000),
				AmountOut:           sdk.NewIntFromUint64(1000000),
			},
		},
		{
			"Create Vault : AppID 1 extended pair 1 user address 2",
			vaultTypes.MsgCreateRequest{
				From:                userAddress2,
				AppId:               1,
				ExtendedPairVaultId: 1,
				AmountIn:            sdk.NewIntFromUint64(1000000),
				AmountOut:           sdk.NewIntFromUint64(1000000),
			},
		},
	} {
		s.Run(tc.name, func() {
			_, err := server.MsgCreate(sdk.WrapSDKContext(*ctx), &tc.msg)
			s.Require().NoError(err)
			res, err := s.vaultQuerier.QueryAllVaults(sdk.WrapSDKContext(*ctx), &vaultTypes.QueryAllVaultsRequest{})
			s.Require().NoError(err)
			_, err = s.vaultQuerier.QueryVaultInfoByVaultID(sdk.WrapSDKContext(*ctx), &vaultTypes.QueryVaultInfoByVaultIDRequest{Id: res.Vault[index].Id})
			s.Require().NoError(err)
		})
	}
}

func (s *KeeperTestSuite) GetVaultCount() int {
	ctx := &s.ctx
	res, err := s.vaultQuerier.QueryAllVaults(sdk.WrapSDKContext(*ctx), &vaultTypes.QueryAllVaultsRequest{})
	s.Require().NoError(err)
	return len(res.Vault)
}

func (s *KeeperTestSuite) GetVaultCountForExtendedPairIDbyAppID(appID, extID uint64) int {
	vaultKeeper1, ctx := &s.vaultKeeper, &s.ctx
	res, found := vaultKeeper1.GetAppExtendedPairVaultMappingData(*ctx, appID, extID)
	s.Require().True(found)
	return len(res.VaultIds)
}

func (s *KeeperTestSuite) GetAssetPrice(id uint64) sdk.Dec {
	price, err := s.app.MarketKeeper.GetLatestPrice(s.ctx, id)
	s.Suite.NoError(err)
	price1 := sdk.NewDecFromInt(sdk.NewIntFromUint64(price))
	return price1
}

func (s *KeeperTestSuite) AddAppAsset() {
	userAddress1 := "cosmos1q7q90qsl9g0gl2zz0njxwv2a649yqrtyxtnv3v"
	userAddress2 := "cosmos1hm7w7dnvdnra78pz9qxysy7u4tuhc3fnpjmyj7"
	addr, err := sdk.AccAddressFromBech32(userAddress1)
	s.Require().NoError(err)
	addr2, err := sdk.AccAddressFromBech32(userAddress2)
	s.Require().NoError(err)
	genesisSupply := sdk.NewIntFromUint64(1000000)
	assetKeeper, ctx := &s.assetKeeper, &s.ctx
	msg1 := assetTypes.AppData{
		Name:             "cswap",
		ShortName:        "cswap",
		MinGovDeposit:    sdk.NewIntFromUint64(10000000),
		GovTimeInSeconds: 900,
		GenesisToken: []assetTypes.MintGenesisToken{
			{
				AssetId:       3,
				GenesisSupply: genesisSupply,
				IsGovToken:    true,
				Recipient:     userAddress1,
			},
			{
				AssetId:       2,
				GenesisSupply: genesisSupply,
				IsGovToken:    true,
				Recipient:     userAddress1,
			},
		},
	}
	err = assetKeeper.AddAppRecords(*ctx, msg1)
	s.Require().NoError(err)

	for _, tc := range []struct {
		name string
		msg  assetTypes.Asset
	}{
		{
			"Add Asset 1",
			assetTypes.Asset{
				Name:          "CMDX",
				Denom:         "ucmdx",
				Decimals:      sdk.NewInt(1000000),
				IsOnChain:     true,
				IsCdpMintable: true,
			},
		},
		{
			"Add Asset 2",
			assetTypes.Asset{
				Name:          "CMST",
				Denom:         "ucmst",
				Decimals:      sdk.NewInt(1000000),
				IsOnChain:     true,
				IsCdpMintable: true,
			},
		},
		{
			"Add Asset 3",
			assetTypes.Asset{
				Name:          "HARBOR",
				Denom:         "uharbor",
				Decimals:      sdk.NewInt(1000000),
				IsOnChain:     true,
				IsCdpMintable: true,
			},
		},
	} {
		s.Run(tc.name, func() {
			err := assetKeeper.AddAssetRecords(*ctx, tc.msg)
			s.Require().NoError(err)
			s.fundAddr(addr, sdk.NewCoin(tc.msg.Denom, sdk.NewInt(1000000)))
			s.fundAddr(addr2, sdk.NewCoin(tc.msg.Denom, sdk.NewInt(1000000)))
		})
	}
}

func (s *KeeperTestSuite) TestLiquidateVaults1() {
	liquidationKeeper, ctx := &s.liquidationKeeper, &s.ctx
	s.CreateVault()
	currentVaultsCount := 2
	s.Require().Equal(s.GetVaultCount(), currentVaultsCount)
	s.Require().Equal(s.GetVaultCountForExtendedPairIDbyAppID(uint64(1), 1), currentVaultsCount)
	beforeVault, found := s.vaultKeeper.GetVault(*ctx, 1)
	s.Require().True(found)

	s.AddAuctionParams()

	// Liquidation shouldn't happen as price not changed
	err := liquidationKeeper.LiquidateVaults(*ctx)
	s.Require().NoError(err)
	id := liquidationKeeper.GetLockedVaultID(*ctx)
	s.Require().Equal(id, uint64(0))

	// Liquidation should happen as price changed
	s.ChangeOraclePrice(1)
	err = liquidationKeeper.LiquidateVaults(*ctx)
	s.Require().NoError(err)
	id = liquidationKeeper.GetLockedVaultID(*ctx)
	s.Require().Equal(id, uint64(2))
	s.Require().Equal(s.GetVaultCount(), currentVaultsCount-2)
	s.Require().Equal(s.GetVaultCountForExtendedPairIDbyAppID(uint64(1), 1), currentVaultsCount-2)

	lockedVault := liquidationKeeper.GetLockedVaults(*ctx)
	s.Require().Equal(lockedVault[0].OriginalVaultId, beforeVault.Id)
	s.Require().Equal(lockedVault[0].ExtendedPairId, beforeVault.ExtendedPairVaultID)
	s.Require().Equal(lockedVault[0].Owner, beforeVault.Owner)
	s.Require().Equal(lockedVault[0].AmountIn, beforeVault.AmountIn)
	s.Require().Equal(lockedVault[0].AmountOut, beforeVault.AmountOut)
	s.Require().Equal(lockedVault[0].UpdatedAmountOut, sdk.ZeroInt())
	s.Require().Equal(lockedVault[0].Initiator, liquidationTypes.ModuleName)
	s.Require().Equal(lockedVault[0].IsAuctionInProgress, true)
	s.Require().Equal(lockedVault[0].IsAuctionComplete, false)
	s.Require().Equal(lockedVault[0].SellOffHistory, []string(nil))
	price, err := s.app.MarketKeeper.CalcAssetPrice(*ctx, uint64(1), beforeVault.AmountIn)
	s.Require().NoError(err)
	s.Require().Equal(lockedVault[0].CollateralToBeAuctioned, price)
	s.Require().Equal(lockedVault[0].CrAtLiquidation, sdk.NewDecFromInt(lockedVault[0].AmountIn).Mul(s.GetAssetPrice(1)).Quo(sdk.NewDecFromInt(lockedVault[0].AmountOut).Mul(s.GetAssetPrice(2))))
}

func (s *KeeperTestSuite) TestUpdateLockedVaults() {
	s.TestLiquidateVaults1()
	liquidationKeeper, ctx := &s.liquidationKeeper, &s.ctx

	lockedVault1 := liquidationKeeper.GetLockedVaults(*ctx)
	s.Require().Equal(lockedVault1[0].CrAtLiquidation, sdk.NewDecFromInt(lockedVault1[0].AmountIn).Mul(s.GetAssetPrice(1)).Quo(sdk.NewDecFromInt(lockedVault1[0].AmountOut).Mul(s.GetAssetPrice(2))))
}

func (s *KeeperTestSuite) TestSetFlags() {
	liquidationKeeper, ctx := &s.liquidationKeeper, &s.ctx
	s.TestUpdateLockedVaults()
	err := liquidationKeeper.SetFlagIsAuctionInProgress(*ctx, 1, 1, true)
	s.Require().NoError(err)
	lockedVault, found := liquidationKeeper.GetLockedVault(*ctx, 1, 1)
	s.Require().True(found)
	s.Require().True(lockedVault.IsAuctionInProgress)
	err = liquidationKeeper.SetFlagIsAuctionInProgress(*ctx, 1, 1, false)
	s.Require().NoError(err)
	lockedVault, found = liquidationKeeper.GetLockedVault(*ctx, 1, 1)
	s.Require().True(found)
	s.Require().False(lockedVault.IsAuctionInProgress)
	err = liquidationKeeper.SetFlagIsAuctionComplete(*ctx, 1, 1, true)
	s.Require().NoError(err)
	lockedVault, found = liquidationKeeper.GetLockedVault(*ctx, 1, 1)
	s.Require().True(found)
	s.Require().True(lockedVault.IsAuctionComplete)
	err = liquidationKeeper.SetFlagIsAuctionComplete(*ctx, 1, 1, false)
	s.Require().NoError(err)
	lockedVault, found = liquidationKeeper.GetLockedVault(*ctx, 1, 1)
	s.Require().True(found)
	s.Require().False(lockedVault.IsAuctionComplete)
}

func (s *KeeperTestSuite) TestDeleteLockedVault() {
	s.TestUpdateLockedVaults()
	liquidationKeeper, ctx := &s.liquidationKeeper, &s.ctx
	vault, found := liquidationKeeper.GetLockedVault(*ctx, 1, 1)
	s.Require().True(found)
	err := liquidationKeeper.CreateLockedVaultHistory(*ctx, vault)
	s.Require().NoError(err)
	id := liquidationKeeper.GetLockedVaultIDHistory(*ctx)
	s.Require().Equal(id, uint64(1))
	liquidationKeeper.DeleteLockedVault(*ctx, 1, 1)
	_, found = liquidationKeeper.GetLockedVault(*ctx, 1, 1)
	s.Require().False(found)
}
