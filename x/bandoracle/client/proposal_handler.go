package client

import (
	"github.com/comdex-official/comdex/x/bandoracle/client/cli"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
)

var (
	AddFetchPriceHandler = govclient.NewProposalHandler(cli.NewCmdSubmitFetchPriceProposal, nil)
)
