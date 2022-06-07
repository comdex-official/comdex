package client

import (
	"github.com/comdex-official/comdex/x/asset/client/cli"
	"github.com/comdex-official/comdex/x/asset/client/rest"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
)

var (
	AddAssetsHandler = []govclient.ProposalHandler{
		govclient.NewProposalHandler(cli.NewCmdSubmitAddAssetsProposal, rest.AddNewAssetsProposalRESTHandler),
		govclient.NewProposalHandler(cli.NewCmdSubmitUpdateAssetProposal, rest.UpdateNewAssetProposalRESTHandler),
		govclient.NewProposalHandler(cli.NewCmdSubmitAddPairsProposal, rest.AddNewPairsProposalRESTHandler),
		govclient.NewProposalHandler(cli.NewCmdSubmitAddWhitelistedAssetsProposal, rest.AddNewWhitelistedAssetsProposalRESTHandler),
		govclient.NewProposalHandler(cli.NewCmdUpdateWhitelistedAssetProposal, rest.UpdateNewWhitelistedAssetsProposalRESTHandler),
		govclient.NewProposalHandler(cli.NewCmdAddWhitelistedPairsProposal, rest.AddNewWhitelistedPairsProposalRESTHandler),
		govclient.NewProposalHandler(cli.NewCmdUpdateWhitelistedPairProposal, rest.UpdateNewWhitelistedPairProposalRESTHandler),
		govclient.NewProposalHandler(cli.NewCmdSubmitAddAppMapingProposal, rest.AddNewAppMappingProposalRESTHandler),
		govclient.NewProposalHandler(cli.NewCmdSubmitAddExtendedPairsVaultProposal, rest.AddExtendedPairsVaultProposalRESTHandler),
		govclient.NewProposalHandler(cli.NewCmdSubmitAddAssetMapingProposal, rest.AddNewAssetMappingProposalRESTHandler),
	}
)
