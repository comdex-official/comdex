package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"

	"github.com/comdex-official/comdex/x/liquidity/client/cli"
)

var LiquidityProposalHandler = []govclient.ProposalHandler{
	govclient.NewProposalHandler(cli.NewCmdUpdateGenericParamsProposal),
	govclient.NewProposalHandler(cli.NewCmdCreateNewLiquidityPairProposal),
}
