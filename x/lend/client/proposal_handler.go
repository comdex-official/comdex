package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"

	"github.com/comdex-official/comdex/x/lend/client/cli"
	"github.com/comdex-official/comdex/x/lend/client/rest"
)

var (
	AddLendPairsHandler        = govclient.NewProposalHandler(cli.CmdAddNewLendPairsProposal, rest.AddNewPairsProposalRESTHandler)
	UpdateLendPairsHandler     = govclient.NewProposalHandler(cli.CmdUpdateLendPairProposal, rest.UpdatePairProposalRESTHandler)
	AddPoolHandler             = govclient.NewProposalHandler(cli.CmdAddPoolProposal, rest.AddPoolProposalRESTHandler)
	AddAssetToPairHandler      = govclient.NewProposalHandler(cli.CmdAddAssetToPairProposal, rest.AddAssetToPairProposalRESTHandler)
	AddAssetRatesParamsHandler = govclient.NewProposalHandler(cli.CmdAddNewAssetRatesParamsProposal, rest.AddNewAssetRatesParamsProposalRESTHandler)
	AddAuctionParamsHandler    = govclient.NewProposalHandler(cli.CmdAddNewAuctionParamsProposal, rest.AddNewAuctionParamsProposalRESTHandler)
)
