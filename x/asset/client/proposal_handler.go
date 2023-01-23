package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"

	"github.com/comdex-official/comdex/x/asset/client/cli"
	"github.com/comdex-official/comdex/x/asset/client/rest"
)

var AddAssetsHandler = []govclient.ProposalHandler{
	govclient.NewProposalHandler(cli.NewCmdSubmitAddAssetsProposal, rest.AddNewAssetsProposalRESTHandler),
	govclient.NewProposalHandler(cli.NewCmdSubmitUpdateAssetProposal, rest.UpdateNewAssetProposalRESTHandler),
	govclient.NewProposalHandler(cli.NewCmdSubmitAddPairsProposal, rest.AddNewPairsProposalRESTHandler),
	govclient.NewProposalHandler(cli.NewCmdSubmitUpdatePairProposal, rest.UpdateNewPairProposalRESTHandler),
	govclient.NewProposalHandler(cli.NewCmdSubmitAddAppProposal, rest.AddNewAppProposalRESTHandler),
	govclient.NewProposalHandler(cli.NewCmdSubmitAddAssetInAppProposal, rest.AddNewAssetInAppProposalRESTHandler),
	govclient.NewProposalHandler(cli.NewCmdSubmitUpdateGovTimeInAppProposal, rest.UpdateNewGovTimeInAppProposalRESTHandler),
	govclient.NewProposalHandler(cli.NewCmdSubmitAddMultipleAssetsProposal, rest.AddNewAssetsProposalRESTHandler),
	govclient.NewProposalHandler(cli.NewCmdSubmitAddMultiplePairsProposal, rest.AddNewPairsProposalRESTHandler),
}
