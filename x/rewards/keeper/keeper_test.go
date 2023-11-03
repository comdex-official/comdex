package keeper_test

import (
	"encoding/binary"
	"strings"
	"testing"

	"github.com/comdex-official/comdex/app/wasm/bindings"
	utils "github.com/comdex-official/comdex/types"
	assettypes "github.com/comdex-official/comdex/x/asset/types"
	collectorKeeper "github.com/comdex-official/comdex/x/collector/keeper"
	lendkeeper "github.com/comdex-official/comdex/x/lend/keeper"
	"github.com/comdex-official/comdex/x/liquidity"
	"github.com/comdex-official/comdex/x/liquidity/types"
	markettypes "github.com/comdex-official/comdex/x/market/types"
	rewardsKeeper "github.com/comdex-official/comdex/x/rewards/keeper"
	rewardstypes "github.com/comdex-official/comdex/x/rewards/types"
	vaultKeeper "github.com/comdex-official/comdex/x/vault/keeper"

	"github.com/stretchr/testify/suite"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	chain "github.com/comdex-official/comdex/app"
	assetKeeper "github.com/comdex-official/comdex/x/asset/keeper"
	lockerKeeper "github.com/comdex-official/comdex/x/locker/keeper"
	lockerTypes "github.com/comdex-official/comdex/x/locker/types"
)

type KeeperTestSuite struct {
	suite.Suite

	app           *chain.App
	ctx           sdk.Context
	assetKeeper   assetKeeper.Keeper
	lockerKeeper  lockerKeeper.Keeper
	querier       rewardsKeeper.QueryServer
	msgServer     rewardstypes.MsgServer
	collector     collectorKeeper.Keeper
	rewardsKeeper rewardsKeeper.Keeper
	vaultKeeper   vaultKeeper.Keeper
	lendKeeper    lendkeeper.Keeper
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (s *KeeperTestSuite) SetupTest() {
	s.app = chain.Setup(s.T(), false)
	s.ctx = s.app.BaseApp.NewContext(false, tmproto.Header{})
	s.assetKeeper = s.app.AssetKeeper
	s.lockerKeeper = s.app.LockerKeeper
	s.querier = rewardsKeeper.QueryServer{Keeper: s.rewardsKeeper}
	s.msgServer = rewardsKeeper.NewMsgServerImpl(s.rewardsKeeper)
	s.collector = s.app.CollectorKeeper
	s.rewardsKeeper = s.app.Rewardskeeper
	s.vaultKeeper = s.app.VaultKeeper
	s.lendKeeper = s.app.LendKeeper
}

func (s *KeeperTestSuite) fundAddr(addr string, amt sdk.Coin) {
	s.T().Helper()
	err := s.app.BankKeeper.MintCoins(s.ctx, lockerTypes.ModuleName, sdk.NewCoins(amt))
	s.Require().NoError(err)
	addr1, err := sdk.AccAddressFromBech32(addr)
	err = s.app.BankKeeper.SendCoinsFromModuleToAccount(s.ctx, lockerTypes.ModuleName, addr1, sdk.NewCoins(amt))
	s.Require().NoError(err)
}

func (s *KeeperTestSuite) fundAddr2(addr sdk.AccAddress, amt sdk.Coins) {
	s.T().Helper()
	err := s.app.BankKeeper.MintCoins(s.ctx, types.ModuleName, amt)
	s.Require().NoError(err)
	err = s.app.BankKeeper.SendCoinsFromModuleToAccount(s.ctx, types.ModuleName, addr, amt)
	s.Require().NoError(err)
}

func (s *KeeperTestSuite) getBalances(addr sdk.AccAddress) sdk.Coins {
	return s.app.BankKeeper.GetAllBalances(s.ctx, addr)
}

func (s *KeeperTestSuite) CreateNewPair(addr sdk.Address, assetIn, assetOut uint64) uint64 {
	err := s.app.AssetKeeper.AddPairsRecords(s.ctx, assettypes.Pair{
		AssetIn:  assetIn,
		AssetOut: assetOut,
	})
	s.Suite.NoError(err)
	pairs := s.app.AssetKeeper.GetPairs(s.ctx)
	var pairID uint64
	for _, pair := range pairs {
		if pair.AssetIn == assetIn && pair.AssetOut == assetOut {
			pairID = pair.Id
			break
		}
	}
	s.Require().NotZero(pairID)
	return pairID
}

func (s *KeeperTestSuite) CreateNewApp(appName string) uint64 {
	err := s.app.AssetKeeper.AddAppRecords(s.ctx, assettypes.AppData{
		Name:             strings.ToLower(appName),
		ShortName:        strings.ToLower(appName),
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

func (s *KeeperTestSuite) CreateNewAsset(name, denom string, price uint64) assettypes.Asset {
	err := s.app.AssetKeeper.AddAssetRecords(s.ctx, assettypes.Asset{
		Name:                  name,
		Denom:                 denom,
		Decimals:              sdk.NewInt(1000000),
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

func (s *KeeperTestSuite) addr(addrNum int) sdk.AccAddress {
	addr := make(sdk.AccAddress, 20)
	binary.PutVarint(addr, int64(addrNum))
	return addr
}

func (s *KeeperTestSuite) CreateNewExtendedVaultPair(
	pairName string,
	appMappingID, pairID uint64,
	isStableMintVault, isVaultActive bool,
) uint64 {
	err := s.app.AssetKeeper.WasmAddExtendedPairsVaultRecords(s.ctx, &bindings.MsgAddExtendedPairsVault{
		AppID:               appMappingID,
		PairID:              pairID,
		StabilityFee:        sdk.NewDecWithPrec(2, 2), // 0.02
		ClosingFee:          sdk.NewDec(0),
		LiquidationPenalty:  sdk.NewDecWithPrec(15, 2), // 0.15
		DrawDownFee:         sdk.NewDecWithPrec(1, 2),  // 0.01
		IsVaultActive:       isVaultActive,
		DebtCeiling:         sdk.NewInt(1000000000000000000),
		DebtFloor:           sdk.NewInt(100000000),
		IsStableMintVault:   isStableMintVault,
		MinCr:               sdk.NewDecWithPrec(23, 1), // 2.3
		PairName:            pairName,
		AssetOutOraclePrice: true,
		AssetOutPrice:       1000000,
		MinUsdValueLeft:     1000000,
	})
	s.Suite.Require().NoError(err)

	extendedVaultPairs, found := s.app.AssetKeeper.GetPairsVaults(s.ctx)
	s.Suite.Require().True(found)

	var extendedVaultPairID uint64
	for _, extendedVaultPair := range extendedVaultPairs {
		if extendedVaultPair.PairName == pairName && extendedVaultPair.AppId == appMappingID {
			extendedVaultPairID = extendedVaultPair.Id
			break
		}
	}
	s.Require().NotZero(extendedVaultPairID)
	return extendedVaultPairID
}

func (s *KeeperTestSuite) CreateNewLiquidityPair(appID uint64, creator sdk.AccAddress, baseCoinDenom, quoteCoinDenom string) types.Pair {
	params, err := s.app.LiquidityKeeper.GetGenericParams(s.ctx, appID)
	s.Require().NoError(err)

	s.fundAddr2(creator, params.PairCreationFee)

	msg := types.NewMsgCreatePair(appID, creator, baseCoinDenom, quoteCoinDenom)
	pair, err := s.app.LiquidityKeeper.CreatePair(s.ctx, msg, false)

	s.Require().NoError(err)
	s.Require().IsType(types.Pair{}, pair)

	return pair
}

func (s *KeeperTestSuite) CreateNewLiquidityPool(appID, pairID uint64, creator sdk.AccAddress, depositCoins string) types.Pool {
	params, err := s.app.LiquidityKeeper.GetGenericParams(s.ctx, appID)
	s.Require().NoError(err)

	parsedDepositCoins := utils.ParseCoins(depositCoins)

	s.fundAddr2(creator, params.PoolCreationFee)
	s.fundAddr2(creator, parsedDepositCoins)
	msg := types.NewMsgCreatePool(appID, creator, pairID, parsedDepositCoins)
	pool, err := s.app.LiquidityKeeper.CreatePool(s.ctx, msg)
	s.Require().NoError(err)
	s.Require().IsType(types.Pool{}, pool)

	return pool
}

func (s *KeeperTestSuite) Deposit(appID, poolID uint64, depositor sdk.AccAddress, depositCoins string) types.DepositRequest {
	msg := types.NewMsgDeposit(
		appID, depositor, poolID, utils.ParseCoins(depositCoins),
	)
	s.fundAddr2(depositor, msg.DepositCoins)
	req, err := s.app.LiquidityKeeper.Deposit(s.ctx, msg)
	s.Require().NoError(err)
	s.Require().IsType(types.DepositRequest{}, req)
	return req
}

func (s *KeeperTestSuite) Farm(appID, poolID uint64, farmer sdk.AccAddress, farmingCoin string) {
	msg := types.NewMsgFarm(
		appID, poolID, farmer, utils.ParseCoin(farmingCoin),
	)
	s.fundAddr2(farmer, sdk.NewCoins(msg.FarmingPoolCoin))
	err := s.app.LiquidityKeeper.Farm(s.ctx, msg)
	s.Require().NoError(err)
}

func (s *KeeperTestSuite) nextBlock() {
	liquidity.EndBlocker(s.ctx, s.app.LiquidityKeeper, s.app.AssetKeeper)
	liquidity.BeginBlocker(s.ctx, s.app.LiquidityKeeper, s.app.AssetKeeper)
}
