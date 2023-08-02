package client

import (
	"github.com/comdex-official/comdex/x/auctionsV2/client/cli"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
)

var AuctionsV2Handler = []govclient.ProposalHandler{
	govclient.NewProposalHandler(cli.NewAddAuctionParamsProposal),
}
