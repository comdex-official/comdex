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

}

func (suite *CdpTestSuite) TestDepositCollateral() {
	_, addrs := app.GeneratePrivKeyAddressPairs(1)
	err := suite.keeper.DepositCollateral(suite.ctx, addrs[0], sdk.NewCoin("cmdx", sdk.NewInt(200000000) ),"cmdx-a")
	suite.Require().True(errors.Is(err, types.ErrorCdpNotFound))

}

func (suite *CdpTestSuite) TestWithdrawCollateral() {

	_, addrs := app.GeneratePrivKeyAddressPairs(1)
	err := suite.keeper.WithdrawCollateral(suite.ctx,addrs[0],sdk.NewCoin("cmdx", sdk.NewInt(200000000) ),"cmdx-a")
	suite.Require().True(errors.Is(err, types.ErrorCdpNotFound))



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




func (suite *CdpTestSuite) Test(){
	_, addrs := app.GeneratePrivKeyAddressPairs(1)
	err := suite.keeper.VerifyBalance(suite.ctx,sdk.NewCoin("cmdx", sdk.NewInt(200000)) , addrs[0] )
	suite.Require().True(errors.Is(err, types.ErrorAccountNotFound))
}