package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"

	"github.com/comdex-official/comdex/x/lend/client/cli"
	"github.com/comdex-official/comdex/x/lend/client/rest"
)

var (
	AddLendPairsHandler       = govclient.NewProposalHandler(cli.CmdAddNewLendPairsProposal, rest.AddNewPairsProposalRESTHandler)
	AddPoolHandler            = govclient.NewProposalHandler(cli.CmdAddPoolProposal, rest.AddPoolProposalRESTHandler)
	AddAssetToPairHandler     = govclient.NewProposalHandler(cli.CmdAddAssetToPairProposal, rest.AddAssetToPairProposalRESTHandler)
	AddAssetRatesStatsHandler = govclient.NewProposalHandler(cli.CmdAddNewAssetRatesStatsProposal, rest.AddWNewAssetRatesStatsProposalRESTHandler)
	AddAuctionParamsHandler   = govclient.NewProposalHandler(cli.CmdAddNewAuctionParamsProposal, rest.AddNewAuctionParamsProposalRESTHandler)
)
