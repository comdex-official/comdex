package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"

	"github.com/comdex-official/comdex/x/asset/client/cli"
)

var AddAssetsHandler = []govclient.ProposalHandler{
	govclient.NewProposalHandler(cli.NewCmdSubmitAddAssetsProposal),
	govclient.NewProposalHandler(cli.NewCmdSubmitUpdateAssetProposal),
	govclient.NewProposalHandler(cli.NewCmdSubmitAddPairsProposal),
	govclient.NewProposalHandler(cli.NewCmdSubmitUpdatePairProposal),
	govclient.NewProposalHandler(cli.NewCmdSubmitAddAppProposal),
	govclient.NewProposalHandler(cli.NewCmdSubmitAddAssetInAppProposal),
	govclient.NewProposalHandler(cli.NewCmdSubmitUpdateGovTimeInAppProposal),
	govclient.NewProposalHandler(cli.NewCmdSubmitAddMultipleAssetsProposal),
	govclient.NewProposalHandler(cli.NewCmdSubmitAddMultiplePairsProposal),
	govclient.NewProposalHandler(cli.NewCmdSubmitAddMultipleAssetsPairsProposal),
}
