package client

import (
	"github.com/comdex-official/comdex/x/liquidationsV2/client/cli"
	"github.com/comdex-official/comdex/x/liquidationsV2/client/rest"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
)

var LiquidationsV2Handler = []govclient.ProposalHandler{
	govclient.NewProposalHandler(cli.NewCmdSubmitWhitelistingLiquidationProposal, rest.AddNewWhitelistingLiquidationRESTHandler),
}
