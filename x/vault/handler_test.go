package vault_test

import (
	"strings"
	"testing"

	_ "github.com/stretchr/testify/suite"

	"github.com/comdex-official/comdex/app"
	"github.com/comdex-official/comdex/x/vault"
	"github.com/comdex-official/comdex/x/vault/keeper"
	"github.com/comdex-official/comdex/x/vault/types"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

func TestInvalidMsg(t *testing.T) {
	app := app.Setup(false)

	app.VaultKeeper = keeper.NewKeeper(
		app.AppCodec(), app.GetKey(types.StoreKey), app.BankKeeper, &app.AssetKeeper, &app.MarketKeeper, &app.CollectorKeeper, &app.EsmKeeper,
		app.TokenmintKeeper, app.Rewardskeeper)
	h := vault.NewHandler(app.VaultKeeper)

	res, err := h(sdk.NewContext(nil, tmproto.Header{}, false, nil), testdata.NewTestMsg())
	require.Error(t, err)
	require.Nil(t, res)

	_, _, log := sdkerrors.ABCIInfo(err, false)
	require.True(t, strings.Contains(log, "unknown message type"))
}

func (s *ModuleTestSuite) TestMsgCreate() {
	handler := vault.NewHandler(s.keeper)
	addr1 := s.addr(1)
	appID1 := s.CreateNewApp("appone")
	asseOneID := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asseTwoID := s.CreateNewAsset("ASSETTWO", "uasset2", 2000000)
	pairID := s.CreateNewPair(addr1, asseOneID, asseTwoID)
	extendedVaultPairID1 := s.CreateNewExtendedVaultPair("CMDX-C", appID1, pairID, false, true)

	msg := types.NewMsgCreateRequest(
		addr1, appID1, extendedVaultPairID1, sdk.NewInt(1000000000), sdk.NewInt(200000000),
	)

	s.fundAddr(addr1, sdk.NewCoins(sdk.NewCoin("uasset1", msg.AmountIn)))

	_, err := handler(s.ctx, msg)
	s.Require().NoError(err)

	availableBalances := s.getBalances(sdk.MustAccAddressFromBech32(addr1.String()))
	s.Require().True(sdk.NewCoins(sdk.NewCoin("uasset2", sdk.NewInt(198000000))).IsEqual(availableBalances))

	vaults := s.keeper.GetVaults(s.ctx)
	s.Require().Len(vaults, 1)
	s.Require().Equal(uint64(1), vaults[0].Id)
	s.Require().Equal(addr1.String(), vaults[0].Owner)
	s.Require().Equal(msg.AmountIn, vaults[0].AmountIn)
	s.Require().Equal(msg.AmountOut, vaults[0].AmountOut)
	s.Require().Equal(extendedVaultPairID1, vaults[0].ExtendedPairVaultID)
	s.Require().Equal(appID1, vaults[0].AppId)
}

func (s *ModuleTestSuite) TestMsgDeposit() {
	handler := vault.NewHandler(s.keeper)
	addr1 := s.addr(1)

	appID1 := s.CreateNewApp("appone")
	asseOneID := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asseTwoID := s.CreateNewAsset("ASSETTWO", "uasset2", 2000000)
	pairID := s.CreateNewPair(addr1, asseOneID, asseTwoID)
	extendedVaultPairID1 := s.CreateNewExtendedVaultPair("CMDX-C", appID1, pairID, false, true)

	msgCreate := types.NewMsgCreateRequest(addr1, appID1, extendedVaultPairID1, sdk.NewInt(1000000000), sdk.NewInt(200000000))
	s.fundAddr(sdk.MustAccAddressFromBech32(addr1.String()), sdk.NewCoins(sdk.NewCoin("uasset1", msgCreate.AmountIn)))

	_, err := handler(s.ctx, msgCreate)
	s.Require().NoError(err)

	msgDeposit := types.NewMsgDepositRequest(
		addr1, appID1, extendedVaultPairID1, 1, sdk.NewInt(69000000),
	)
	s.fundAddr(sdk.MustAccAddressFromBech32(addr1.String()), sdk.NewCoins(sdk.NewCoin("uasset1", msgDeposit.Amount)))

	_, err = handler(s.ctx, msgDeposit)
	s.Require().NoError(err)

	vaults := s.keeper.GetVaults(s.ctx)
	s.Require().Len(vaults, 1)
	s.Require().Equal(uint64(1), vaults[0].Id)
	s.Require().Equal(addr1.String(), vaults[0].Owner)
	s.Require().Equal(msgCreate.AmountIn.Add(msgDeposit.Amount), vaults[0].AmountIn)
	s.Require().Equal(msgCreate.AmountOut, vaults[0].AmountOut)
	s.Require().Equal(extendedVaultPairID1, vaults[0].ExtendedPairVaultID)
	s.Require().Equal(appID1, vaults[0].AppId)
}

func (s *ModuleTestSuite) TestMsgWithdraw() {
	handler := vault.NewHandler(s.keeper)
	addr1 := s.addr(1)

	appID1 := s.CreateNewApp("appone")
	asseOneID := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asseTwoID := s.CreateNewAsset("ASSETTWO", "uasset2", 2000000)
	pairID := s.CreateNewPair(addr1, asseOneID, asseTwoID)
	extendedVaultPairID1 := s.CreateNewExtendedVaultPair("CMDX-C", appID1, pairID, false, true)

	msgCreate := types.NewMsgCreateRequest(addr1, appID1, extendedVaultPairID1, sdk.NewInt(1000000000), sdk.NewInt(200000000))
	s.fundAddr(sdk.MustAccAddressFromBech32(addr1.String()), sdk.NewCoins(sdk.NewCoin("uasset1", msgCreate.AmountIn)))

	_, err := handler(s.ctx, msgCreate)
	s.Require().NoError(err)

	msgWithdraw := types.NewMsgWithdrawRequest(
		addr1, appID1, extendedVaultPairID1, 1, sdk.NewInt(50000000),
	)
	_, err = handler(s.ctx, msgWithdraw)
	s.Require().NoError(err)

	vaults := s.keeper.GetVaults(s.ctx)
	s.Require().Len(vaults, 1)
	s.Require().Equal(uint64(1), vaults[0].Id)
	s.Require().Equal(addr1.String(), vaults[0].Owner)
	s.Require().Equal(msgCreate.AmountIn.Sub(msgWithdraw.Amount), vaults[0].AmountIn)
	s.Require().Equal(msgCreate.AmountOut, vaults[0].AmountOut)
	s.Require().Equal(extendedVaultPairID1, vaults[0].ExtendedPairVaultID)
	s.Require().Equal(appID1, vaults[0].AppId)
}

func (s *ModuleTestSuite) TestMsgDraw() {
	handler := vault.NewHandler(s.keeper)
	addr1 := s.addr(1)

	appID1 := s.CreateNewApp("appone")
	asseOneID := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asseTwoID := s.CreateNewAsset("ASSETTWO", "uasset2", 2000000)
	pairID := s.CreateNewPair(addr1, asseOneID, asseTwoID)
	extendedVaultPairID1 := s.CreateNewExtendedVaultPair("CMDX-C", appID1, pairID, false, true)

	msgCreate := types.NewMsgCreateRequest(addr1, appID1, extendedVaultPairID1, sdk.NewInt(1000000000), sdk.NewInt(200000000))
	s.fundAddr(sdk.MustAccAddressFromBech32(addr1.String()), sdk.NewCoins(sdk.NewCoin("uasset1", msgCreate.AmountIn)))

	_, err := handler(s.ctx, msgCreate)
	s.Require().NoError(err)

	msgDraw := types.NewMsgDrawRequest(
		addr1, appID1, extendedVaultPairID1, 1, sdk.NewInt(10000000),
	)
	_, err = handler(s.ctx, msgDraw)
	s.Require().NoError(err)

	vaults := s.keeper.GetVaults(s.ctx)
	s.Require().Len(vaults, 1)
	s.Require().Equal(uint64(1), vaults[0].Id)
	s.Require().Equal(addr1.String(), vaults[0].Owner)
	s.Require().Equal(msgCreate.AmountIn, vaults[0].AmountIn)
	s.Require().Equal(msgCreate.AmountOut.Add(msgDraw.Amount), vaults[0].AmountOut)
	s.Require().Equal(extendedVaultPairID1, vaults[0].ExtendedPairVaultID)
	s.Require().Equal(appID1, vaults[0].AppId)
}

func (s *ModuleTestSuite) TestMsgRepay() {
	handler := vault.NewHandler(s.keeper)
	addr1 := s.addr(1)

	appID1 := s.CreateNewApp("appone")
	asseOneID := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asseTwoID := s.CreateNewAsset("ASSETTWO", "uasset2", 2000000)
	pairID := s.CreateNewPair(addr1, asseOneID, asseTwoID)
	extendedVaultPairID1 := s.CreateNewExtendedVaultPair("CMDX-C", appID1, pairID, false, true)

	msgCreate := types.NewMsgCreateRequest(addr1, appID1, extendedVaultPairID1, sdk.NewInt(1000000000), sdk.NewInt(200000000))
	s.fundAddr(sdk.MustAccAddressFromBech32(addr1.String()), sdk.NewCoins(sdk.NewCoin("uasset1", msgCreate.AmountIn)))

	_, err := handler(s.ctx, msgCreate)
	s.Require().NoError(err)

	msgRepay := types.NewMsgRepayRequest(
		addr1, appID1, extendedVaultPairID1, 1, sdk.NewInt(100000000),
	)
	_, err = handler(s.ctx, msgRepay)
	s.Require().NoError(err)

	vaults := s.keeper.GetVaults(s.ctx)
	s.Require().Len(vaults, 1)
	s.Require().Equal(uint64(1), vaults[0].Id)
	s.Require().Equal(addr1.String(), vaults[0].Owner)
	s.Require().Equal(msgCreate.AmountIn, vaults[0].AmountIn)
	s.Require().Equal(msgCreate.AmountOut.Sub(msgRepay.Amount), vaults[0].AmountOut)
	s.Require().Equal(extendedVaultPairID1, vaults[0].ExtendedPairVaultID)
	s.Require().Equal(appID1, vaults[0].AppId)
}

func (s *ModuleTestSuite) TestMsgClose() {
	handler := vault.NewHandler(s.keeper)
	addr1 := s.addr(1)

	appID1 := s.CreateNewApp("appone")
	asseOneID := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asseTwoID := s.CreateNewAsset("ASSETTWO", "uasset2", 2000000)
	pairID := s.CreateNewPair(addr1, asseOneID, asseTwoID)
	extendedVaultPairID1 := s.CreateNewExtendedVaultPair("CMDX-C", appID1, pairID, false, true)

	msgCreate := types.NewMsgCreateRequest(addr1, appID1, extendedVaultPairID1, sdk.NewInt(1000000000), sdk.NewInt(200000000))
	s.fundAddr(sdk.MustAccAddressFromBech32(addr1.String()), sdk.NewCoins(sdk.NewCoin("uasset1", msgCreate.AmountIn)))

	_, err := handler(s.ctx, msgCreate)
	s.Require().NoError(err)

	s.fundAddr(sdk.MustAccAddressFromBech32(addr1.String()), sdk.NewCoins(sdk.NewCoin("uasset2", sdk.NewInt(2000000))))

	msgRepay := types.NewMsgLiquidateRequest(
		addr1, appID1, extendedVaultPairID1, 1,
	)
	_, err = handler(s.ctx, msgRepay)
	s.Require().NoError(err)

	vaults := s.keeper.GetVaults(s.ctx)
	s.Require().Len(vaults, 0)
}

func (s *ModuleTestSuite) TestMsgCreateStableMint() {
	handler := vault.NewHandler(s.keeper)
	addr1 := s.addr(1)

	appID1 := s.CreateNewApp("appone")
	asseOneID := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asseTwoID := s.CreateNewAsset("ASSETTWO", "uasset2", 2000000)
	pairID := s.CreateNewPair(addr1, asseOneID, asseTwoID)
	extendedVaultPairID1 := s.CreateNewExtendedVaultPair("CMDX-C", appID1, pairID, true, true)

	msgStableMintCreate := types.NewMsgCreateStableMintRequest(
		addr1, appID1, extendedVaultPairID1, sdk.NewInt(100000000),
	)
	s.fundAddr(sdk.MustAccAddressFromBech32(addr1.String()), sdk.NewCoins(sdk.NewCoin("uasset1", msgStableMintCreate.Amount)))

	_, err := handler(s.ctx, msgStableMintCreate)
	s.Require().NoError(err)

	stableMintVaults := s.querier.GetStableMintVaults(s.ctx)
	s.Require().Equal(uint64(1), stableMintVaults[0].Id)
	s.Require().Equal(msgStableMintCreate.Amount, stableMintVaults[0].AmountIn)
	s.Require().Equal(msgStableMintCreate.Amount, stableMintVaults[0].AmountOut)
	s.Require().Equal(appID1, stableMintVaults[0].AppId)
	s.Require().Equal(extendedVaultPairID1, stableMintVaults[0].ExtendedPairVaultID)
}

func (s *ModuleTestSuite) TestMsgDepositStableMint() {
	handler := vault.NewHandler(s.keeper)
	addr1 := s.addr(1)

	appID1 := s.CreateNewApp("appone")
	asseOneID := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asseTwoID := s.CreateNewAsset("ASSETTWO", "uasset2", 2000000)
	pairID := s.CreateNewPair(addr1, asseOneID, asseTwoID)
	extendedVaultPairID1 := s.CreateNewExtendedVaultPair("CMDX-C", appID1, pairID, true, true)

	msgStableMintCreate := types.NewMsgCreateStableMintRequest(
		addr1, appID1, extendedVaultPairID1, sdk.NewInt(100000000),
	)
	s.fundAddr(sdk.MustAccAddressFromBech32(addr1.String()), sdk.NewCoins(sdk.NewCoin("uasset1", msgStableMintCreate.Amount)))

	_, err := handler(s.ctx, msgStableMintCreate)
	s.Require().NoError(err)

	msgStableMintDeposit := types.NewMsgDepositStableMintRequest(
		addr1, appID1, extendedVaultPairID1, sdk.NewInt(2000000000), 1,
	)
	s.fundAddr(sdk.MustAccAddressFromBech32(addr1.String()), sdk.NewCoins(sdk.NewCoin("uasset1", msgStableMintDeposit.Amount)))

	_, err = handler(s.ctx, msgStableMintDeposit)
	s.Require().NoError(err)

	stableMintVaults := s.querier.GetStableMintVaults(s.ctx)
	s.Require().Equal(uint64(1), stableMintVaults[0].Id)
	s.Require().Equal(msgStableMintCreate.Amount.Add(msgStableMintDeposit.Amount), stableMintVaults[0].AmountIn)
	s.Require().Equal(msgStableMintCreate.Amount.Add(msgStableMintDeposit.Amount), stableMintVaults[0].AmountOut)
	s.Require().Equal(appID1, stableMintVaults[0].AppId)
	s.Require().Equal(extendedVaultPairID1, stableMintVaults[0].ExtendedPairVaultID)
}

func (s *ModuleTestSuite) TestMsgWithdrawStableMint() {
	handler := vault.NewHandler(s.keeper)
	addr1 := s.addr(1)

	appID1 := s.CreateNewApp("appone")
	asseOneID := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asseTwoID := s.CreateNewAsset("ASSETTWO", "uasset2", 2000000)
	pairID := s.CreateNewPair(addr1, asseOneID, asseTwoID)
	extendedVaultPairID1 := s.CreateNewExtendedVaultPair("CMDX-C", appID1, pairID, true, true)

	msgStableMintCreate := types.NewMsgCreateStableMintRequest(
		addr1, appID1, extendedVaultPairID1, sdk.NewInt(1000000000),
	)
	s.fundAddr(sdk.MustAccAddressFromBech32(addr1.String()), sdk.NewCoins(sdk.NewCoin("uasset1", msgStableMintCreate.Amount)))

	_, err := handler(s.ctx, msgStableMintCreate)
	s.Require().NoError(err)

	msgStableMintWithdraw := types.NewMsgWithdrawStableMintRequest(
		addr1, appID1, extendedVaultPairID1, sdk.NewInt(500000000), 1,
	)
	s.fundAddr(sdk.MustAccAddressFromBech32(addr1.String()), sdk.NewCoins(sdk.NewCoin("uasset1", msgStableMintWithdraw.Amount)))

	_, err = handler(s.ctx, msgStableMintWithdraw)
	s.Require().NoError(err)

	stableMintVaults := s.querier.GetStableMintVaults(s.ctx)
	s.Require().Equal(uint64(1), stableMintVaults[0].Id)
	s.Require().Equal(sdk.NewInt(505000000), stableMintVaults[0].AmountIn)
	s.Require().Equal(sdk.NewInt(505000000), stableMintVaults[0].AmountOut)
	s.Require().Equal(appID1, stableMintVaults[0].AppId)
	s.Require().Equal(extendedVaultPairID1, stableMintVaults[0].ExtendedPairVaultID)
}
