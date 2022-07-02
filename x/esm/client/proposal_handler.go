package client

import (
	"github.com/comdex-official/comdex/x/esm/client/cli"
	"github.com/comdex-official/comdex/x/esm/client/rest"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
)

var (
	AddESMTriggerParams = govclient.NewProposalHandler(cli.CmdAddESMTriggerParamsProposal, rest.AddESMTriggerParamsProposal)
)
