package client

import (
	"github.com/comdex-official/comdex/x/asset/client/cli"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
)

var AddAsssetHandler = govclient.NewProposalHandler(cli.NewCmdSubmitAddAssetsProposal, nil)
