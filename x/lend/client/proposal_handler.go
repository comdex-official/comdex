package client

import (
	"github.com/comdex-official/comdex/x/lend/client/cli"
	"github.com/comdex-official/comdex/x/lend/client/rest"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
)

var (
	AddLendPairsHandler       = govclient.NewProposalHandler(cli.CmdAddNewLendPairsProposal, rest.AddNewPairsProposalRESTHandler)
	UpdateLendPairsHandler    = govclient.NewProposalHandler(cli.CmdUpdateLendPairProposal, rest.UpdatePairProposalRESTHandler)
	AddPoolHandler            = govclient.NewProposalHandler(cli.CmdAddPoolProposal, rest.AddPoolProposalRESTHandler)
	AddAssetToPairHandler     = govclient.NewProposalHandler(cli.CmdAddAssetToPairProposal, rest.AddAssetToPairProposalRESTHandler)
	AddAssetRatesStatsHandler = govclient.NewProposalHandler(cli.CmdAddNewAssetRatesStatsProposal, rest.AddNewAssetRatesStatsProposalRESTHandler)
	AddAuctionParamsHandler   = govclient.NewProposalHandler(cli.CmdAddNewAuctionParamsProposal, rest.AddNewAuctionParamsProposalRESTHandler)
)
