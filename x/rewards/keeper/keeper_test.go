package keeper_test

import (
	"encoding/binary"
	"github.com/petrichormoney/petri/app/wasm/bindings"
	assettypes "github.com/petrichormoney/petri/x/asset/types"
	collectorKeeper "github.com/petrichormoney/petri/x/collector/keeper"
	rewardsKeeper "github.com/petrichormoney/petri/x/rewards/keeper"
	rewardstypes "github.com/petrichormoney/petri/x/rewards/types"
	vaultKeeper "github.com/petrichormoney/petri/x/vault/keeper"
	"testing"

	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	chain "github.com/petrichormoney/petri/app"
	assetKeeper "github.com/petrichormoney/petri/x/asset/keeper"
	lockerKeeper "github.com/petrichormoney/petri/x/locker/keeper"
	lockerTypes "github.com/petrichormoney/petri/x/locker/types"
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
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (s *KeeperTestSuite) SetupTest() {
	s.app = chain.Setup(false)
	s.ctx = s.app.BaseApp.NewContext(false, tmproto.Header{})
	s.assetKeeper = s.app.AssetKeeper
	s.lockerKeeper = s.app.LockerKeeper
	s.querier = rewardsKeeper.QueryServer{Keeper: s.rewardsKeeper}
	s.msgServer = rewardsKeeper.NewMsgServerImpl(s.rewardsKeeper)
	s.collector = s.app.CollectorKeeper
	s.rewardsKeeper = s.app.Rewardskeeper
	s.vaultKeeper = s.app.VaultKeeper
}

func (s *KeeperTestSuite) fundAddr(addr string, amt sdk.Coin) {
	s.T().Helper()
	err := s.app.BankKeeper.MintCoins(s.ctx, lockerTypes.ModuleName, sdk.NewCoins(amt))
	s.Require().NoError(err)
	addr1, err := sdk.AccAddressFromBech32(addr)
	err = s.app.BankKeeper.SendCoinsFromModuleToAccount(s.ctx, lockerTypes.ModuleName, addr1, sdk.NewCoins(amt))
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
