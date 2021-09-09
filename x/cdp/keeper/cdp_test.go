package keeper_test

import (
	"errors"
	"github.com/comdex-official/comdex/app"
	"github.com/comdex-official/comdex/x/cdp/keeper"
	"github.com/comdex-official/comdex/x/cdp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtime "github.com/tendermint/tendermint/types/time"
	"math/big"
	"testing"
)

type CdpTestSuite struct {
	suite.Suite
	keeper keeper.Keeper
	app    app.TestApp
	ctx    sdk.Context
}

func (suite *CdpTestSuite) SetupTest() {
	testApp := app.NewTestApp()
	k := testApp.GetCDPKeeper()
	ctx := testApp.NewContext(true, tmproto.Header{Height: 1, Time: tmtime.Now()})
	suite.app = testApp
	suite.keeper = k
	suite.ctx = ctx

	return
}
func TestCdpTestSuite(t *testing.T) {
	suite.Run(t, new(CdpTestSuite))
}



func (suite *CdpTestSuite)TestCdp_GetNextCdpID() {
	id := suite.keeper.GetNextCdpID(suite.ctx)
	suite.Equal(types.DefaultIndex, id)
}

func (suite *CdpTestSuite) TestCdp_SetGet() {
	addrs, _ := sdk.AccAddressFromBech32("abc")
	cdp := types.NewCDP(types.DefaultIndex, addrs,sdk.NewCoin("cmdx", sdk.NewInt(1)), "cmdx-a", sdk.NewCoin("uscx", sdk.NewInt(1)))
	suite.keeper.SetCDP(suite.ctx, cdp)

	t, found := suite.keeper.GetCDP(suite.ctx, types.DefaultIndex)
	suite.True(found)
	suite.Equal(cdp, t)
	_, found = suite.keeper.GetCDP(suite.ctx, 100)
	suite.False(found)
	suite.keeper.DeleteCDP(suite.ctx, cdp)
	_, found = suite.keeper.GetCDP(suite.ctx, 100)
	suite.False(found)

	suite.keeper.SetNextCdpId(suite.ctx,1)

}

func (suite *CdpTestSuite) TestAddCdp() {

	addrs, _ := sdk.AccAddressFromBech32("abc")


	cdp := types.NewCDP(types.DefaultIndex, addrs,sdk.NewCoin("cmdx", sdk.NewInt(200000000)), "cmdx-a", sdk.NewCoin("uscx", sdk.NewInt(500000000)))
	suite.keeper.SetCDP(suite.ctx, cdp)

	params := types.CollateralParam{
		CollateralDenom:  "cmdx",
		DebtDenom:        "uscx",
		Type:             "cmdx-a",
		LiquidationRatio: sdk.Dec{},
	}
	cparams := [] types.CollateralParam{params}
	suite.keeper.SetParams(suite.ctx,types.Params{cparams})

	err := suite.keeper.AddCdp(suite.ctx, addrs, sdk.NewCoin("xprt", sdk.NewInt(200000000)), sdk.NewCoin("uscx", sdk.NewInt(10000000)), "btc-a")
	suite.Require().True(errors.Is(err, types.ErrorCdpNotFound))

}

func (suite *CdpTestSuite) TestVerifyCollateralAndDebt(){
	params := types.CollateralParam{
		CollateralDenom:  "cmdx",
		DebtDenom:        "uscx",
		Type:             "cmdx-a",
		LiquidationRatio: sdk.Dec{},
	}
	cparams := [] types.CollateralParam{params}
	suite.keeper.SetParams(suite.ctx,types.Params{cparams})
	suite.keeper.GetCollateralParam(suite.ctx, "cmdx-a")

	err := suite.keeper.VerifyCollateralAndDebt(suite.ctx, sdk.NewCoin("xprt", sdk.NewInt(500000)), sdk.NewCoin("uscx", sdk.NewInt(10000000)), "cmdx-a" )
	suite.Require().True(errors.Is(err, types.ErrorInvalidCollateral))

	err = suite.keeper.VerifyCollateralAndDebt(suite.ctx, sdk.NewCoin("cmdx", sdk.NewInt(500000)), sdk.NewCoin("xprt", sdk.NewInt(10000000)), "cmdx-a" )
	suite.Require().True(errors.Is(err, types.ErrorInvalidDebt))
}

func (suite *CdpTestSuite) TestDepositCollateral() {

	_, addrs := app.GeneratePrivKeyAddressPairs(2)

	err := suite.keeper.DepositCollateral(suite.ctx, addrs[0], sdk.NewCoin("cmdx", sdk.NewInt(200000000) ),"cmdx-a")
	suite.Require().True(errors.Is(err, types.ErrorCdpNotFound))

	cdp := types.NewCDP(types.DefaultIndex, addrs[0],sdk.NewCoin("cmdx", sdk.NewInt(100)), "cmdx-a", sdk.NewCoin("uscx", sdk.NewInt(10)))
	suite.keeper.SetCDP(suite.ctx, cdp)
	suite.keeper.IndexCDPByOwner(suite.ctx, cdp)

	err = suite.keeper.DepositCollateral(suite.ctx, addrs[0], sdk.NewCoin("eth", sdk.NewInt(200000000) ),"cmdx-a")
	suite.Require().True(errors.Is(err, types.ErrorInvalidCollateral))

}

func (suite *CdpTestSuite) TestWithdrawCollateral() {
	_, addrs := app.GeneratePrivKeyAddressPairs(2)

	params := types.CollateralParam{
		CollateralDenom:  "cmdx",
		DebtDenom:        "uscx",
		Type:             "cmdx-a",
		LiquidationRatio: sdk.Dec{},
	}
	cparams := [] types.CollateralParam{params}
	suite.keeper.SetParams(suite.ctx,types.Params{cparams})
	suite.keeper.GetCollateralParam(suite.ctx, "cmdx-a")

	cdp := types.NewCDP(types.DefaultIndex, addrs[0],sdk.NewCoin("cmdx", sdk.NewInt(1)), "cmdx-a", sdk.NewCoin("uscx", sdk.NewInt(1)))
	suite.keeper.SetCDP(suite.ctx, cdp)
	suite.keeper.IndexCDPByOwner(suite.ctx, cdp)

	err := suite.keeper.WithdrawCollateral(suite.ctx, addrs[0], sdk.NewCoin("xprt", sdk.NewInt(200000000)), "xprt-a" )
	suite.Require().True(errors.Is(err, types.ErrorCdpNotFound))

	err = suite.keeper.WithdrawCollateral(suite.ctx, addrs[0], sdk.NewCoin("uscx", sdk.NewInt(200000000)) , "cmdx-a")
	suite.Require().True(errors.Is(err, types.ErrorInvalidCollateral))

	err = suite.keeper.WithdrawCollateral(suite.ctx, addrs[0], sdk.NewCoin("cmdx", sdk.NewInt(200000000)) , "cmdx-a")
	suite.Require().True(errors.Is(err, types.ErrorInvalidWithdrawAmount))


}

func (suite *CdpTestSuite) TestDrawDebt(){
	_, addrs := app.GeneratePrivKeyAddressPairs(1)
	err := suite.keeper.DrawDebt(suite.ctx, addrs[0],  "cmdx", sdk.NewCoin("cmdx", sdk.NewInt(200000)))
	suite.Require().True(errors.Is(err, types.ErrorCdpNotFound))

	cdp := types.NewCDP(types.DefaultIndex, addrs[0],sdk.NewCoin("cmdx", sdk.NewInt(10)), "cmdx-a", sdk.NewCoin("uscx", sdk.NewInt(10)))
	suite.keeper.SetCDP(suite.ctx, cdp)
	suite.keeper.IndexCDPByOwner(suite.ctx, cdp)

	err = suite.keeper.DrawDebt(suite.ctx, addrs[0],  "cmdx-a", sdk.NewCoin("xprt", sdk.NewInt(200000)))
	suite.Require().True(errors.Is(err, types.ErrorInvalidDebt))

}

func (suite *CdpTestSuite) TestRepayDebt(){
	_, addrs := app.GeneratePrivKeyAddressPairs(1)
	err := suite.keeper.RepayDebt(suite.ctx, addrs[0],  "cmdx", sdk.NewCoin("cmdx", sdk.NewInt(200000)))
	suite.Require().True(errors.Is(err, types.ErrorCdpNotFound))

	cdp := types.NewCDP(types.DefaultIndex, addrs[0],sdk.NewCoin("cmdx", sdk.NewInt(100)), "cmdx-a", sdk.NewCoin("uscx", sdk.NewInt(10)))
	suite.keeper.SetCDP(suite.ctx, cdp)
	suite.keeper.IndexCDPByOwner(suite.ctx, cdp)

	err = suite.keeper.RepayDebt(suite.ctx, addrs[0],  "cmdx-a", sdk.NewCoin("cmdx", sdk.NewInt(20)))
	suite.Require().True(errors.Is(err, types.ErrorInvalidPayment))

	err = suite.keeper.RepayDebt(suite.ctx, addrs[0],  "cmdx-a", sdk.NewCoin("uscx", sdk.NewInt(20)))
	suite.Require().True(errors.Is(err, types.ErrorInvalidAmount))

	err = suite.keeper.RepayDebt(suite.ctx, addrs[0],  "cmdx-a", sdk.NewCoin("uscx", sdk.NewInt(10)))
	suite.Error(err)

}

func (suite *CdpTestSuite) TestAttemptLiquidation(){
	_, addrs := app.GeneratePrivKeyAddressPairs(1)
	err := suite.keeper.AttemptLiquidation(suite.ctx,addrs[0], "cmdx-a" )
	suite.Require().True(errors.Is(err, types.ErrorCdpNotFound))
	params := types.CollateralParam{
		CollateralDenom:  "cmdx",
		DebtDenom:        "uscx",
		Type:             "cmdx-a",
		LiquidationRatio: sdk.Dec{},
	}
	cparams := [] types.CollateralParam{params}
	suite.keeper.SetParams(suite.ctx,types.Params{cparams})
	cdp := types.NewCDP(types.DefaultIndex, addrs[0],sdk.NewCoin("cmdx", sdk.NewInt(1)), "cmdx-a", sdk.NewCoin("uscx", sdk.NewInt(10)))
	suite.keeper.SetCDP(suite.ctx, cdp)
	suite.keeper.IndexCDPByOwner(suite.ctx, cdp)

	err = suite.keeper.AttemptLiquidation(suite.ctx,addrs[0], "cmdx-a" )
	suite.Error(err)

}

func (suite *CdpTestSuite) TestVerifyBalance(){
	addrs, _ := sdk.AccAddressFromBech32("abc")
	_, addr:= app.GeneratePrivKeyAddressPairs(1)
	err := suite.keeper.VerifyBalance(suite.ctx,sdk.NewCoin("cmdx", sdk.NewInt(200000)) , addrs )
	suite.Require().True(errors.Is(err, types.ErrorAccountNotFound))

	_= banktypes.Balance{
		Address: "bytes,1,opt,name=address,proto3",
		Coins:   sdk.Coins{sdk.NewCoin("cmdx", sdk.NewInt(4))},
	}

	accountKeeper := suite.app.GetAccountKeeper()
	acc := accountKeeper.NewAccountWithAddress(suite.ctx, addr[0])
	accountKeeper.SetAccount(suite.ctx, acc)

	err = suite.keeper.VerifyBalance(suite.ctx,sdk.NewCoin("cmdx", sdk.NewInt(3)) , addr[0] )
	suite.Require().True(errors.Is(err, types.ErrorInsufficientBalance))

}

func (suite *CdpTestSuite) TestKeeper_VerifyCollateralizationRatio() {

	suite.keeper.SetParams(suite.ctx, types.Params{})
	params := types.CollateralParam{
		CollateralDenom:  "cmdx",
		DebtDenom:        "usdx",
		Type:             "cmdx-a",
		LiquidationRatio: sdk.NewDecFromBigInt(big.NewInt(150)),
	}
	cparams := []types.CollateralParam{params}
	suite.keeper.SetParams(suite.ctx, types.Params{cparams})
	err := suite.keeper.VerifyCollateralizationRatio(suite.ctx, sdk.NewCoin("cmdx", sdk.NewInt(200000)), sdk.NewCoin("cmdx", sdk.NewInt(200000)), "abc-a")
	suite.Require().True(errors.Is(err, types.ErrorCollateralNotFound))

}

func (suite *CdpTestSuite) TestKeeper_VerifyLiquidation(){
	suite.keeper.SetParams(suite.ctx,types.Params{})
	err := suite.keeper.VerifyLiquidation(suite.ctx,sdk.NewCoin("cmdx", sdk.NewInt(200000)), sdk.NewCoin("cmdx", sdk.NewInt(200000)),"cmdx-a")
	suite.Error(err)
}

func (suite *CdpTestSuite) TestIndexCDPByOwner(){
	_, addrs := app.GeneratePrivKeyAddressPairs(1)
	cdp := types.NewCDP(types.DefaultIndex, addrs[0],sdk.NewCoin("cmdx", sdk.NewInt(1)), "cmdx-a", sdk.NewCoin("uscx", sdk.NewInt(1)))
	suite.keeper.SetCDP(suite.ctx, cdp)
	suite.keeper.IndexCDPByOwner(suite.ctx, cdp)
	suite.keeper.GetOwnerCDPList(suite.ctx, addrs[0])
	suite.keeper.IndexCDPByOwner(suite.ctx, cdp)
}

func (suite *CdpTestSuite) TestKeeper_GetCDPByOwnerAndCollateralType() {
	_, addrs := app.GeneratePrivKeyAddressPairs(1)
	cdp := types.NewCDP(types.DefaultIndex, addrs[0],sdk.NewCoin("cmdx", sdk.NewInt(1)), "cmdx-a", sdk.NewCoin("uscx", sdk.NewInt(1)))
	suite.keeper.IndexCDPByOwner(suite.ctx, cdp)
	suite.keeper.GetCDP(suite.ctx,1 )
	suite.keeper.GetCDPByOwnerAndCollateralType(suite.ctx, addrs[0], "cmdx-a")
}