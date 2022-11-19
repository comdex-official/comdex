package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"

	"github.com/petrichormoney/petri/x/liquidity/client/cli"
	"github.com/petrichormoney/petri/x/liquidity/client/rest"
)

var LiquidityProposalHandler = []govclient.ProposalHandler{
	govclient.NewProposalHandler(cli.NewCmdUpdateGenericParamsProposal, rest.UpdateGenericParamsProposalRESTHandler),
	govclient.NewProposalHandler(cli.NewCmdCreateNewLiquidityPairProposal, rest.CreateNewLiquidityPairProposalRESTHandler),
}
