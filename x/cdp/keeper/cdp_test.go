package keeper_test

import (
	"errors"
	"github.com/comdex-official/comdex/app"
	"github.com/comdex-official/comdex/x/cdp/keeper"
	"github.com/comdex-official/comdex/x/cdp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtime "github.com/tendermint/tendermint/types/time"
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
	_, addrs := app.GeneratePrivKeyAddressPairs(1)
	cdp := types.NewCDP(types.DefaultIndex, addrs[0],sdk.NewCoin("xrp", sdk.NewInt(1)), "xrp-a", sdk.NewCoin("usdx", sdk.NewInt(1)))
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

	addrs, err := sdk.AccAddressFromBech32("abc")

	cdp := types.NewCDP(types.DefaultIndex, addrs,sdk.NewCoin("xrp", sdk.NewInt(200000000)), "xrp-a", sdk.NewCoin("usdx", sdk.NewInt(500000000)))
	suite.keeper.SetCDP(suite.ctx, cdp)

	params := types.CollateralParam{
		CollateralDenom:  "xrp",
		DebtDenom:        "usdx",
		Type:             "xrp-a",
		LiquidationRatio: sdk.Dec{},
	}
	cparams := [] types.CollateralParam{params}
	suite.keeper.SetParams(suite.ctx,types.Params{cparams})

	ak := suite.app.GetAccountKeeper()
	acc := ak.NewAccountWithAddress(suite.ctx, addrs)

	bk := suite.app.GetBankKeeper()
	bk.AddCoins(suite.ctx,acc.GetAddress(),sdk.Coins{sdk.NewCoin("cmdx", sdk.NewInt(200000000)), sdk.NewCoin("btc", sdk.NewInt(500000000))})
	ak.SetAccount(suite.ctx, acc)

	err = suite.keeper.AddCdp(suite.ctx, addrs, sdk.NewCoin("cmdx", sdk.NewInt(200000000)), sdk.NewCoin("usdx", sdk.NewInt(10000000)), "btc-a")
	suite.Require().True(errors.Is(err, types.ErrorCdpNotFound))

	err = suite.keeper.AddCdp(suite.ctx, addrs, sdk.NewCoin("xrp", sdk.NewInt(200000000)), sdk.NewCoin("usdx", sdk.NewInt(10000000)), "xrp-a")
	suite.Error(err)

}

func (suite *CdpTestSuite) TestVerifyCollateralAndDebt(){
	params := types.CollateralParam{
		CollateralDenom:  "xrp",
		DebtDenom:        "usdx",
		Type:             "cmdx-a",
		LiquidationRatio: sdk.Dec{},
	}
	cparams := [] types.CollateralParam{params}
	suite.keeper.SetParams(suite.ctx,types.Params{cparams})
	suite.keeper.GetCollateralParam(suite.ctx, "cmdx-a")

	err := suite.keeper.VerifyCollateralAndDebt(suite.ctx, sdk.NewCoin("cmdx", sdk.NewInt(500000)), sdk.NewCoin("usdx", sdk.NewInt(10000000)), "cmdx-a" )
	suite.Require().True(errors.Is(err, types.ErrorInvalidCollateral))

	err = suite.keeper.VerifyCollateralAndDebt(suite.ctx, sdk.NewCoin("xrp", sdk.NewInt(500000)), sdk.NewCoin("cmdx", sdk.NewInt(10000000)), "cmdx-a" )
	suite.Require().True(errors.Is(err, types.ErrorInvalidDebt))
}

func (suite *CdpTestSuite) TestDepositCollateral() {

	_, addrs := app.GeneratePrivKeyAddressPairs(2)

	err := suite.keeper.DepositCollateral(suite.ctx, addrs[0], sdk.NewCoin("cmdx", sdk.NewInt(200000000) ),"cmdx-a")
	suite.Require().True(errors.Is(err, types.ErrorCdpNotFound))

	cdp := types.NewCDP(types.DefaultIndex, addrs[0],sdk.NewCoin("cmdx", sdk.NewInt(100)), "cmdx-a", sdk.NewCoin("usdx", sdk.NewInt(10)))
	suite.keeper.SetCDP(suite.ctx, cdp)
	suite.keeper.IndexCDPByOwner(suite.ctx, cdp)

	err = suite.keeper.DepositCollateral(suite.ctx, addrs[0], sdk.NewCoin("eth", sdk.NewInt(200000000) ),"cmdx-a")
	suite.Require().True(errors.Is(err, types.ErrorInvalidCollateral))

}

func (suite *CdpTestSuite) TestWithdrawCollateral() {
	_, addrs := app.GeneratePrivKeyAddressPairs(2)

	params := types.CollateralParam{
		CollateralDenom:  "cmdx",
		DebtDenom:        "usdx",
		Type:             "cmdx-a",
		LiquidationRatio: sdk.Dec{},
	}
	cparams := [] types.CollateralParam{params}
	suite.keeper.SetParams(suite.ctx,types.Params{cparams})
	suite.keeper.GetCollateralParam(suite.ctx, "cmdx-a")

	cdp := types.NewCDP(types.DefaultIndex, addrs[0],sdk.NewCoin("cmdx", sdk.NewInt(1)), "cmdx-a", sdk.NewCoin("usdx", sdk.NewInt(1)))
	suite.keeper.SetCDP(suite.ctx, cdp)
	suite.keeper.IndexCDPByOwner(suite.ctx, cdp)

	err := suite.keeper.WithdrawCollateral(suite.ctx, addrs[0], sdk.NewCoin("abc", sdk.NewInt(200000000)), "abc-a" )
	suite.Require().True(errors.Is(err, types.ErrorCdpNotFound))

	err = suite.keeper.WithdrawCollateral(suite.ctx, addrs[0], sdk.NewCoin("usdx", sdk.NewInt(200000000)) , "cmdx-a")
	suite.Require().True(errors.Is(err, types.ErrorInvalidCollateral))

	err = suite.keeper.WithdrawCollateral(suite.ctx, addrs[0], sdk.NewCoin("cmdx", sdk.NewInt(200000000)) , "cmdx-a")
	suite.Require().True(errors.Is(err, types.ErrorInvalidWithdrawAmount))


}

func (suite *CdpTestSuite) TestDrawDebt(){
	_, addrs := app.GeneratePrivKeyAddressPairs(1)
	err := suite.keeper.DrawDebt(suite.ctx, addrs[0],  "cmdx", sdk.NewCoin("cmdx", sdk.NewInt(200000)))
	suite.Require().True(errors.Is(err, types.ErrorCdpNotFound))

	cdp := types.NewCDP(types.DefaultIndex, addrs[0],sdk.NewCoin("cmdx", sdk.NewInt(10)), "cmdx-a", sdk.NewCoin("usdx", sdk.NewInt(10)))
	suite.keeper.SetCDP(suite.ctx, cdp)
	suite.keeper.IndexCDPByOwner(suite.ctx, cdp)

	err = suite.keeper.DrawDebt(suite.ctx, addrs[0],  "cmdx-a", sdk.NewCoin("xrp", sdk.NewInt(200000)))
	suite.Require().True(errors.Is(err, types.ErrorInvalidDebt))

}

func (suite *CdpTestSuite) TestRepayDebt(){
	_, addrs := app.GeneratePrivKeyAddressPairs(1)
	err := suite.keeper.RepayDebt(suite.ctx, addrs[0],  "cmdx", sdk.NewCoin("cmdx", sdk.NewInt(200000)))
	suite.Require().True(errors.Is(err, types.ErrorCdpNotFound))

	cdp := types.NewCDP(types.DefaultIndex, addrs[0],sdk.NewCoin("cmdx", sdk.NewInt(100)), "cmdx-a", sdk.NewCoin("usdx", sdk.NewInt(10)))
	suite.keeper.SetCDP(suite.ctx, cdp)
	suite.keeper.IndexCDPByOwner(suite.ctx, cdp)

	err = suite.keeper.RepayDebt(suite.ctx, addrs[0],  "cmdx-a", sdk.NewCoin("cmdx", sdk.NewInt(20)))
	suite.Require().True(errors.Is(err, types.ErrorInvalidPayment))

	err = suite.keeper.RepayDebt(suite.ctx, addrs[0],  "cmdx-a", sdk.NewCoin("usdx", sdk.NewInt(20)))
	suite.Require().True(errors.Is(err, types.ErrorInvalidAmount))

	err = suite.keeper.RepayDebt(suite.ctx, addrs[0],  "cmdx-a", sdk.NewCoin("usdx", sdk.NewInt(10)))
	suite.Error(err)

	err = suite.keeper.RepayDebt(suite.ctx, addrs[0],  "cmdx-a", sdk.NewCoin("usdx", sdk.NewInt(5)))
	suite.Error(err)

}

func (suite *CdpTestSuite) TestAttemptLiquidation(){
	_, addrs := app.GeneratePrivKeyAddressPairs(1)
	err := suite.keeper.AttemptLiquidation(suite.ctx,addrs[0], "cmdx-a" )
	suite.Require().True(errors.Is(err, types.ErrorCdpNotFound))
	params := types.CollateralParam{
		CollateralDenom:  "cmdx",
		DebtDenom:        "usdx",
		Type:             "cmdx-a",
		LiquidationRatio: sdk.Dec{},
	}
	cparams := [] types.CollateralParam{params}
	suite.keeper.SetParams(suite.ctx,types.Params{cparams})
	cdp := types.NewCDP(types.DefaultIndex, addrs[0],sdk.NewCoin("cmdx", sdk.NewInt(1)), "cmdx-a", sdk.NewCoin("usdx", sdk.NewInt(10)))
	suite.keeper.SetCDP(suite.ctx, cdp)
	suite.keeper.IndexCDPByOwner(suite.ctx, cdp)

	err = suite.keeper.AttemptLiquidation(suite.ctx,addrs[0], "cmdx-a" )
	suite.Error(err)

}

func (suite *CdpTestSuite) TestVerifyBalance(){
	_, addrs := app.GeneratePrivKeyAddressPairs(1)
	err := suite.keeper.VerifyBalance(suite.ctx,sdk.NewCoin("cmdx", sdk.NewInt(200000)) , addrs[0] )
	suite.Require().True(errors.Is(err, types.ErrorAccountNotFound))

	cdp := types.NewCDP(types.DefaultIndex, addrs[0],sdk.NewCoin("cmdx", sdk.NewInt(1)), "cmdx-a", sdk.NewCoin("usdx", sdk.NewInt(10)))
	suite.keeper.SetCDP(suite.ctx, cdp)
	suite.keeper.IndexCDPByOwner(suite.ctx, cdp)

	ak := suite.app.GetAccountKeeper()
	acc := ak.NewAccountWithAddress(suite.ctx, addrs[0])

	bk := suite.app.GetBankKeeper()
	bk.SpendableCoins(suite.ctx, addrs[0]).AmountOf("cmdx")
	ak.SetAccount(suite.ctx, acc)

	err = suite.keeper.VerifyBalance(suite.ctx,sdk.NewCoin("cmdx", sdk.NewInt(20000000)) , addrs[0] )
	suite.Require().True(errors.Is(err, types.ErrorInsufficientBalance))

}

func (suite *CdpTestSuite) TestKeeper_VerifyCollateralizationRatio() {

	suite.keeper.SetParams(suite.ctx, types.Params{})
	err := suite.keeper.VerifyCollateralizationRatio(suite.ctx, sdk.NewCoin("cmdx", sdk.NewInt(200000)), sdk.NewCoin("cmdx", sdk.NewInt(200000)), "cmdx")
	suite.Require().True(errors.Is(err, types.ErrorCollateralNotFound))
	params := types.CollateralParam{
		CollateralDenom:  "cmdx",
		DebtDenom:        "usdx",
		Type:             "cmdx-a",
		LiquidationRatio: sdk.Dec{},
	}
	cparams := []types.CollateralParam{params}
	suite.keeper.SetParams(suite.ctx, types.Params{cparams})
	err = suite.keeper.VerifyCollateralizationRatio(suite.ctx, sdk.NewCoin("cmdx", sdk.NewInt(200000)), sdk.NewCoin("cmdx", sdk.NewInt(200000)), "cmdx")
	suite.Require().True(errors.Is(err, types.ErrorCollateralNotFound))

}

func (suite *CdpTestSuite) TestKeeper_VerifyLiquidation(){
	suite.keeper.SetParams(suite.ctx,types.Params{})
	err := suite.keeper.VerifyLiquidation(suite.ctx,sdk.NewCoin("cmdx", sdk.NewInt(200000)), sdk.NewCoin("cmdx", sdk.NewInt(200000)),"cmdx")
	suite.Error(err)
}

func (suite *CdpTestSuite) TestIndexCDPByOwner(){
	_, addrs := app.GeneratePrivKeyAddressPairs(1)
	cdp := types.NewCDP(types.DefaultIndex, addrs[0],sdk.NewCoin("cmdx", sdk.NewInt(1)), "cmdx-a", sdk.NewCoin("usdx", sdk.NewInt(1)))
	suite.keeper.SetCDP(suite.ctx, cdp)
	suite.keeper.IndexCDPByOwner(suite.ctx, cdp)
	suite.keeper.GetOwnerCDPList(suite.ctx, addrs[0])
	suite.keeper.IndexCDPByOwner(suite.ctx, cdp)
}

func (suite *CdpTestSuite) TestKeeper_GetCDPByOwnerAndCollateralType() {
	_, addrs := app.GeneratePrivKeyAddressPairs(1)
	cdp := types.NewCDP(types.DefaultIndex, addrs[0],sdk.NewCoin("cmdx", sdk.NewInt(1)), "cmdx-a", sdk.NewCoin("usdx", sdk.NewInt(1)))
	suite.keeper.IndexCDPByOwner(suite.ctx, cdp)
	suite.keeper.GetCDP(suite.ctx,1 )
	suite.keeper.GetCDPByOwnerAndCollateralType(suite.ctx, addrs[0], "cmdx-a")
}

