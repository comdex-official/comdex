package main

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"os"

	"github.com/cosmos/cosmos-sdk/server"
	servercmd "github.com/cosmos/cosmos-sdk/server/cmd"

	comdex "github.com/comdex-official/comdex/app"
	nodecmd "github.com/comdex-official/comdex/node/cmd"

	"github.com/comdex-official/comdex/app"
)

func main() {

	config := sdk.GetConfig()
	app.SetAccountAddressPrefixes(config)
	config.Seal()

	root, _ := nodecmd.NewRootCmd()
	if err := servercmd.Execute(root, comdex.DefaultNodeHome); err != nil {
		switch e := err.(type) {
		case server.ErrorCode:
			os.Exit(e.Code)
		default:
			os.Exit(1)
		}
	}
}
