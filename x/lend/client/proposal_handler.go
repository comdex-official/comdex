package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"

	"github.com/comdex-official/comdex/x/lend/client/cli"
)

var (
	AddLendPairsHandler           = govclient.NewProposalHandler(cli.CmdAddNewLendPairsProposal)
	AddPoolHandler                = govclient.NewProposalHandler(cli.CmdAddPoolProposal)
	AddAssetToPairHandler         = govclient.NewProposalHandler(cli.CmdAddAssetToPairProposal)
	AddMultipleAssetToPairHandler = govclient.NewProposalHandler(cli.CmdAddMultipleAssetToPairProposal)
	AddAssetRatesParamsHandler    = govclient.NewProposalHandler(cli.CmdAddNewAssetRatesParamsProposal)
	AddAuctionParamsHandler       = govclient.NewProposalHandler(cli.CmdAddNewAuctionParamsProposal)
	AddMultipleLendPairsHandler   = govclient.NewProposalHandler(cli.CmdAddNewMultipleLendPairsProposal)
	AddPoolPairsHandler           = govclient.NewProposalHandler(cli.CmdAddPoolPairsProposal)
	AddAssetRatesPoolPairsHandler = govclient.NewProposalHandler(cli.CmdAddAssetRatesPoolPairsProposal)
)
