package keeper_test

import (
	assetTypes "github.com/comdex-official/comdex/x/asset/types"
	markettypes "github.com/comdex-official/comdex/x/market/types"
	vaultKeeper1 "github.com/comdex-official/comdex/x/vault/keeper"
	vaultTypes "github.com/comdex-official/comdex/x/vault/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	protobuftypes "github.com/gogo/protobuf/types"
)

/*
func (k *Keeper) AddAppMappingRecords(ctx sdk.Context, records ...types.AppMapping) error
func (k *Keeper) AddAssetRecords(ctx sdk.Context, records ...types.Asset) error
func (k *Keeper) AddPairsRecords(ctx sdk.Context, records ...types.Pair) error
func (k *Keeper) AddExtendedPairsVaultRecords(ctx sdk.Context, records ...types.ExtendedPairVault) error
func (k Keeper) WhitelistAppId(ctx sdk.Context, appMappingId uint64) error
*/

func (s *KeeperTestSuite) AddAppIDAssetID() {
	assetKeeper, ctx := &s.assetKeeper, &s.ctx
	msg1 := assetTypes.AppMapping{
		Name:             "cswap",
		ShortName:        "cswap",
		MinGovDeposit:    sdk.NewIntFromUint64(10000000),
		GovTimeInSeconds: 900,
	}
	err := assetKeeper.AddAppMappingRecords(*ctx, msg1)
	s.Require().NoError(err)

	msg2 := []assetTypes.Asset{
		{Name: "CMDX",
			Denom:     "ucmdx",
			Decimals:  1000000,
			IsOnChain: true}, {Name: "CMST",
			Denom:     "ucmst",
			Decimals:  1000000,
			IsOnChain: true}, {Name: "HARBOR",
			Denom:     "uharbor",
			Decimals:  1000000,
			IsOnChain: true},
	}
	err = assetKeeper.AddAssetRecords(*ctx, msg2...)
	s.Require().NoError(err)
}

func (s *KeeperTestSuite) AddPairAndExtendedPairVault() {

	assetKeeper, liquidationKeeper, ctx := &s.assetKeeper, &s.liquidationKeeper, &s.ctx

	msg3 := assetTypes.Pair{
		AssetIn:  1,
		AssetOut: 2,
	}
	err := assetKeeper.AddPairsRecords(*ctx, msg3)
	s.Require().NoError(err)

	msg4 := assetTypes.ExtendedPairVault{
		AppMappingId:        1,
		PairId:              1,
		StabilityFee:        sdk.MustNewDecFromStr("0.01"),
		ClosingFee:          sdk.MustNewDecFromStr("0"),
		LiquidationPenalty:  sdk.MustNewDecFromStr("0.12"),
		DrawDownFee:         sdk.MustNewDecFromStr("0.01"),
		IsVaultActive:       true,
		DebtCeiling:         sdk.NewIntFromUint64(1000000000000),
		DebtFloor:           sdk.NewIntFromUint64(1000000),
		IsStableMintVault:   false,
		MinCr:               sdk.MustNewDecFromStr("1.5"),
		PairName:            "CMDX-B",
		AssetOutOraclePrice: true,
		AssetOutPrice:       1000000,
		MinUsdValueLeft:     1000000,
	}
	err = assetKeeper.AddExtendedPairsVaultRecords(*ctx, msg4)
	s.Require().NoError(err)

	err = liquidationKeeper.WhitelistAppID(*ctx, 1)
	s.Require().NoError(err)

}

func (s *KeeperTestSuite) SetOraclePrice(symbol string, price uint64) {
	var (
		store = s.app.MarketKeeper.Store(s.ctx)
		key   = markettypes.PriceForMarketKey(symbol)
	)
	value := s.app.AppCodec().MustMarshal(
		&protobuftypes.UInt64Value{
			Value: price,
		},
	)
	store.Set(key, value)
}
func (s *KeeperTestSuite) SetOraclePriceForSymbols() {
	s.SetOraclePrice("ucmdx", 2000000)
	s.SetOraclePrice("ucmst", 1000000)
}
func (s *KeeperTestSuite) CreateVault() {
	userAddress := "cosmos1q7q90qsl9g0gl2zz0njxwv2a649yqrtyxtnv3v"
	vaultKeeper, ctx := &s.vaultKeeper, &s.ctx
	addr, err := sdk.AccAddressFromBech32(userAddress)
	s.Require().NoError(err)

	s.AddAppIDAssetID()
	s.AddPairAndExtendedPairVault()
	s.SetOraclePriceForSymbols()
	msg5 := vaultTypes.MsgCreateRequest{
		From:                userAddress,
		AppId:               1,
		ExtendedPairVaultId: 1,
		AmountIn:            sdk.NewIntFromUint64(1000000),
		AmountOut:           sdk.NewIntFromUint64(1000000),
	}
	server := vaultKeeper1.NewMsgServer(*vaultKeeper)
	s.fundAddr(addr, sdk.NewCoin("ucmdx", sdk.NewInt(1000000)))
	_, err = server.MsgCreate(sdk.WrapSDKContext(*ctx), &msg5)
	s.Require().NoError(err)
}

func (s *KeeperTestSuite) TestLiquidateVaults() {
	liquidationKeeper, ctx := &s.liquidationKeeper, &s.ctx
	s.CreateVault()
	err := liquidationKeeper.LiquidateVaults(*ctx)
	s.Require().NoError(err)
}
