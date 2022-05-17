package client

import (
	"github.com/comdex-official/comdex/x/lend/client/cli"
	"github.com/comdex-official/comdex/x/lend/client/rest"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
)

var (
	AddLendPairsHandler    = govclient.NewProposalHandler(cli.CmdAddWNewLendPairsProposal, rest.AddNewPairsProposalRESTHandler)
	UpdateLendPairsHandler = govclient.NewProposalHandler(cli.CmdUpdateLendPairProposal, rest.UpdateNewNewPairsProposalRESTHandler)
)
