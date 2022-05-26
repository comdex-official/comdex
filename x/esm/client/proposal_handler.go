package client

import (
	"github.com/comdex-official/comdex/x/esm/client/cli"
	"github.com/comdex-official/comdex/x/esm/client/rest"

	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
)

var (
	ToggleEsmHandler   = govclient.NewProposalHandler(cli.NewCmdSubmitToggleEsmProposal, rest.NewCmdSubmitToggleEsmProposalRESTHandler)
)
