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

/*func (suite *CdpTestSuite) TestAddCdp() {
	_, addrs := app.GeneratePrivKeyAddressPairs(2)
	ak := suite.app.GetAccountKeeper()
	acc := ak.NewAccountWithAddress(suite.ctx, addrs[0])
	bk := suite.app.GetBankKeeper()
	bk.AddCoins(suite.ctx,acc.GetAddress(),sdk.Coins{sdk.NewCoin("cmdx", sdk.NewInt(200000000)), sdk.NewCoin("btc", sdk.NewInt(500000000))})
	ak.SetAccount(suite.ctx, acc)
	err := suite.keeper.AddCdp(suite.ctx, addrs[0], sdk.NewCoin("cmdx", sdk.NewInt(200000000) ), sdk.NewCoin("usdx", sdk.NewInt(10000000) ), "btc-a")
	suite.Require().True(errors.Is(err, types.ErrorInvalidCollateral))
	err = suite.keeper.AddCdp(suite.ctx, addrs[0], sdk.NewCoin("cmdx", sdk.NewInt(200000000) ), sdk.NewCoin("usdx", sdk.NewInt(26000000) ), "xrp-a")
	suite.Require().True(errors.Is(err, types.ErrorInvalidCollateralRatio))
	err = suite.keeper.AddCdp(suite.ctx, addrs[0], sdk.NewCoin("cmdx", sdk.NewInt(500000000) ), sdk.NewCoin("usdx", sdk.NewInt(26000000) ), "xrp-a")
	suite.Error(err) // insufficient balance
	err = suite.keeper.AddCdp(suite.ctx, addrs[0], sdk.NewCoin("cmdx", sdk.NewInt(200000000) ), sdk.NewCoin("xusd", sdk.NewInt(10000000) ), "xrp-a")
	suite.Require().True(errors.Is(err, types.ErrorDebtNotSupported))

}*/

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

}

func (suite *CdpTestSuite) TestDepositCollateral() {
	_, addrs := app.GeneratePrivKeyAddressPairs(1)
	err := suite.keeper.DepositCollateral(suite.ctx, addrs[0], sdk.NewCoin("cmdx", sdk.NewInt(200000000) ),"cmdx-a")
	suite.Require().True(errors.Is(err, types.ErrorCdpNotFound))

	//err = suite.keeper.DepositCollateral(suite.ctx, addrs[0], sdk.NewCoin("cmdx", sdk.NewInt(500000000) ),"cmdx-a")
	//suite.Require().True(errors.Is(err, types.ErrorInvalidCollateral))
}

func (suite *CdpTestSuite) TestWithdrawCollateral() {
	_, addrs := app.GeneratePrivKeyAddressPairs(1)
	err := suite.keeper.WithdrawCollateral(suite.ctx,addrs[0],sdk.NewCoin("cmdx", sdk.NewInt(200000000) ),"cmdx-a")
	suite.Require().True(errors.Is(err, types.ErrorCdpNotFound))

	//err = suite.keeper.WithdrawCollateral(suite.ctx,addrs[1],sdk.NewCoin("usdt", sdk.NewInt(0) ),"usdt-a")
	//suite.Require().True(errors.Is(err, types.ErrorInvalidCollateral))


}

func (suite *CdpTestSuite) TestDrawDebt(){
	_, addrs := app.GeneratePrivKeyAddressPairs(1)
	err := suite.keeper.DrawDebt(suite.ctx, addrs[0],  "cmdx", sdk.NewCoin("cmdx", sdk.NewInt(200000)))
	suite.Require().True(errors.Is(err, types.ErrorCdpNotFound))
}

func (suite *CdpTestSuite) TestRepayDebt(){
	_, addrs := app.GeneratePrivKeyAddressPairs(1)
	err := suite.keeper.RepayDebt(suite.ctx, addrs[0],  "cmdx", sdk.NewCoin("cmdx", sdk.NewInt(200000)))
	suite.Require().True(errors.Is(err, types.ErrorCdpNotFound))

	//err = suite.keeper.RepayDebt(suite.ctx, addrs[0],  "cmdx", sdk.NewCoin("cmdx", sdk.NewInt(200000)))
	//suite.Require().True(errors.Is(err, types.ErrorInvalidPayment))

}

func (suite *CdpTestSuite) TestAttemptLiquidation(){
	_, addrs := app.GeneratePrivKeyAddressPairs(1)
	err := suite.keeper.AttemptLiquidation(suite.ctx,addrs[0], "cmdx" )
	suite.Require().True(errors.Is(err, types.ErrorCdpNotFound))

}

func (suite *CdpTestSuite) TestVerifyBalance(){
	_, addrs := app.GeneratePrivKeyAddressPairs(1)
	err := suite.keeper.VerifyBalance(suite.ctx,sdk.NewCoin("cmdx", sdk.NewInt(200000)) , addrs[0] )
	suite.Require().True(errors.Is(err, types.ErrorAccountNotFound))
}

/*func (suite *CdpTestSuite) TestVerifyCollateralAndDebt(){
	err := suite.keeper.VerifyCollateralAndDebt(suite.ctx,sdk.NewCoin("cmdx", sdk.NewInt(200)) , sdk.NewCoin("cmdx", sdk.NewInt(200)), "cmdx"  )
	suite.Require().True(errors.Is(err, types.ErrorCdpNotFound))
}


func (suite *CdpTestSuite) TestVerifyCollateralizationRatio(){
	err := suite.keeper.VerifyCollateralizationRatio(suite.ctx,sdk.NewCoin("cmdx", sdk.NewInt(200000)) , sdk.NewCoin("cmdx", sdk.NewInt(200)),"hi" )
	suite.Require().True(errors.Is(err, types.ErrorInvalidCollateralRatio))
}

func (suite *CdpTestSuite) TestVerifyLiquidation(){
	err := suite.keeper.VerifyLiquidation(suite.ctx,sdk.NewCoin("cmdx", sdk.NewInt(200000)) , sdk.NewCoin("cmdx", sdk.NewInt(200000)),"hi" )
	suite.Require().True(errors.Is(err, types.ErrorLowCollateralizationRatio))
}*/


func (suite *CdpTestSuite) Test(){
	_, addrs := app.GeneratePrivKeyAddressPairs(1)
	err := suite.keeper.VerifyBalance(suite.ctx,sdk.NewCoin("cmdx", sdk.NewInt(200000)) , addrs[0] )
	suite.Require().True(errors.Is(err, types.ErrorAccountNotFound))
}