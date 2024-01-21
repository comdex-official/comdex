package main

import (
	"os"

	"cosmossdk.io/log"
	servercmd "github.com/cosmos/cosmos-sdk/server/cmd"

	comdex "github.com/comdex-official/comdex/app"
)

func main() {
	comdex.SetAccountAddressPrefixes()

	root := NewRootCmd() //TODO: check wasmd root
	if err := servercmd.Execute(root, "", comdex.DefaultNodeHome); err != nil {
		log.NewLogger(root.OutOrStderr()).Error("failure when running app", "err", err)
		os.Exit(1)
	}
}
