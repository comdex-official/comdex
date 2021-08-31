package main

import (
	"os"

	"github.com/cosmos/cosmos-sdk/server"

	comdex "github.com/comdex-official/comdex/app"
	nodecmd "github.com/comdex-official/comdex/node/cmd"
	servercmd "github.com/cosmos/cosmos-sdk/server/cmd"
	sdk "github.com/cosmos/cosmos-sdk/types"

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
