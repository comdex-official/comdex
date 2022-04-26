package client

import (
	"github.com/comdex-official/comdex/x/lend/client/cli"
	"github.com/comdex-official/comdex/x/lend/client/rest"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
)

var (
	AddWhitelistedAssetsHandler = govclient.NewProposalHandler(cli.NewCmdSubmitAddWhitelistedAssetsProposal, rest.AddNewWhitelistedAssetsProposalRESTHandler)
)
