package client

import (
	"github.com/comdex-official/comdex/x/collector/client/cli"
	"github.com/comdex-official/comdex/x/collector/client/rest"

	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
)

var (
	AddLookupTableParamsHandlers  = []govclient.ProposalHandler{
		govclient.NewProposalHandler(cli.NewCmdLookupTableParams, rest.NewCmdLookupTableParamsRESTHandler),
		govclient.NewProposalHandler(cli.NewCmdAuctionControlProposal, rest.NewCmdAuctionTableAppRESTHandler),

	}
)
