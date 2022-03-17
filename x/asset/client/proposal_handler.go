package client

import (
	"github.com/comdex-official/comdex/x/asset/client/cli"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
)

var (
	AddAssetsHandler = []govclient.ProposalHandler{
		govclient.NewProposalHandler(cli.NewCmdSubmitAddAssetsProposal, nil),
		govclient.NewProposalHandler(cli.NewCmdSubmitUpdateAssetProposal, nil),
		govclient.NewProposalHandler(cli.NewCmdSubmitAddPairsProposal, nil),
		govclient.NewProposalHandler(cli.NewCmdSubmitUpdatePairProposal, nil),
	}
)
