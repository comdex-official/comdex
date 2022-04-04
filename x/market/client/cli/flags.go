package cli

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibcclienttypes "github.com/cosmos/ibc-go/v2/modules/core/02-client/types"
	"github.com/spf13/cobra"
)

const (
	flagLiquidationRatio       = "liquidation-ratio"
	flagScriptID               = "script-id"
	flagPacketTimeoutHeight    = "packet-timeout-height"
	flagFeeLimit               = "fee-limit"

)

func GetLiquidationRatio(cmd *cobra.Command) (sdk.Dec, error) {
	s, err := cmd.Flags().GetString(flagLiquidationRatio)
	if err != nil {
		return sdk.Dec{}, err
	}

	return sdk.NewDecFromStr(s)
}

func GetPacketTimeoutHeight(cmd *cobra.Command) (ibcclienttypes.Height, error) {
	s, err := cmd.Flags().GetString(flagPacketTimeoutHeight)
	if err != nil {
		return ibcclienttypes.Height{}, err
	}

	return ibcclienttypes.ParseHeight(s)
}

func GetFeeLimit(cmd *cobra.Command) (sdk.Coins, error) {
	s, err := cmd.Flags().GetString(flagFeeLimit)
	if err != nil {
		return sdk.Coins{}, err
	}

	return sdk.ParseCoinsNormalized(s)
}
