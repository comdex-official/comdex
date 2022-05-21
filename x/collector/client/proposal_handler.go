package client

import (
	"github.com/comdex-official/comdex/x/collector/client/cli"
	"github.com/comdex-official/comdex/x/collector/client/rest"

	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
)

var (
	AddLookupTableParamsHandlers   = govclient.NewProposalHandler(cli.NewCmdLookupTableParams, rest.NewCmdLookupTableParamsRESTHandler)
	AddAuctionControlParamsHandler = govclient.NewProposalHandler(cli.NewCmdAuctionControlProposal, rest.NewCmdAuctionTableAppRESTHandler)
)
