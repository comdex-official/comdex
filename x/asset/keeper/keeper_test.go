package keeper

import (
	comdex "github.com/comdex-official/comdex/app"
	simapp "github.com/comdex-official/comdex/testutil/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"testing"
	"time"
)

type KeeperTestSuite struct {
	suite.Suite

	app *comdex.App
	ctx sdk.Context
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.app = simapp.New("~/.comdex/testapp")
	suite.ctx = suite.app.BaseApp.NewContext(false, tmproto.Header{Height: 1, ChainID: "comdex-test", Time: time.Now().UTC()})
}

func TestKeeperTestSuite(t *testing.T)  {
	suite.Run(t, new(KeeperTestSuite))
}