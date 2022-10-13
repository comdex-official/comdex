package client

import (
	"github.com/comdex-official/comdex/x/rewards/client/cli"
	"github.com/comdex-official/comdex/x/rewards/client/rest"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
)

var AddRewardsHandler = []govclient.ProposalHandler{
	govclient.NewProposalHandler(cli.SubmitAddLendExternalRewardsProposal, rest.AddLendExternalRewardsProposalRESTHandler),
}
