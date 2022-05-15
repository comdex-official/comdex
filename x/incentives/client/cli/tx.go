package cli

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"

	// "github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/comdex-official/comdex/x/incentives/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	DefaultRelativePacketTimeoutTimestamp = uint64((time.Duration(10) * time.Minute).Nanoseconds())
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		NewCreateGaugeCmd(),
	)
	return cmd
}

// NewLockTokensCmd lock tokens into bonding pool from user's account.
func NewCreateGaugeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-gauge [gauge-type-id] [trigger-duration] [deposit-amount] [total-triggers]",
		Short: "create new gauge",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txf := tx.NewFactoryCLI(clientCtx, cmd.Flags()).WithTxConfig(clientCtx.TxConfig).WithAccountRetriever(clientCtx.AccountRetriever)

			gaugeTypeId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("parse gauge-type-id: %w", err)
			}

			triggerDuration, err := time.ParseDuration(args[1])
			if err != nil {
				return fmt.Errorf("parse trigger-duration: %w", err)
			}

			depositAmount, err := sdk.ParseCoinNormalized(args[2])
			if err != nil {
				return err
			}

			totalTriggers, err := strconv.ParseUint(args[3], 10, 64)
			if err != nil {
				return fmt.Errorf("parse gauge-type-id: %w", err)
			}

			msg := types.NewMsgCreateGauge(
				clientCtx.GetFromAddress(),
				gaugeTypeId,
				triggerDuration,
				depositAmount,
				totalTriggers,
			)

			switch msg.GaugeTypeId {
			case types.LiquidityGaugeTypeId:
				gaugeExtraData, err := NewBuildLiquidityGaugeExtraData(cmd)
				if err != nil {
					return err
				}
				msg.Kind = &gaugeExtraData
			}

			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msg)
		},
	}

	cmd.Flags().AddFlagSet(FlagSetCreateGauge())
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func NewBuildLiquidityGaugeExtraData(cmd *cobra.Command) (types.MsgCreateGauge_LiquidityMetaData, error) {
	poolId, err := cmd.Flags().GetUint64(FlagPoolId)
	if err != nil {
		return types.MsgCreateGauge_LiquidityMetaData{}, err
	}
	if poolId == 0 {
		return types.MsgCreateGauge_LiquidityMetaData{}, fmt.Errorf("%s required but not specified / pool-id cannot be 0", FlagPoolId)
	}

	isMasterPool, err := cmd.Flags().GetBool(FlagIsMasterPool)
	if err != nil {
		return types.MsgCreateGauge_LiquidityMetaData{}, err
	}

	childPoolIdsCombined, err := cmd.Flags().GetString(FlagChildPoolIds)
	if err != nil {
		return types.MsgCreateGauge_LiquidityMetaData{}, err
	}
	childPoolIds := []uint64{}
	if childPoolIdsCombined != "" {
		childPoolIdsStrs := strings.Split(childPoolIdsCombined, ",")
		for _, poolIdStr := range childPoolIdsStrs {
			poolId, err := strconv.ParseUint(poolIdStr, 10, 64)
			if err != nil {
				return types.MsgCreateGauge_LiquidityMetaData{}, err
			}
			childPoolIds = append(childPoolIds, poolId)
		}
	}

	startTime := time.Time{}
	timeStr, err := cmd.Flags().GetString(FlagStartTime)
	if err != nil {
		return types.MsgCreateGauge_LiquidityMetaData{}, err
	}
	if timeStr == "" { // empty start time
		return types.MsgCreateGauge_LiquidityMetaData{}, fmt.Errorf("%s required but not specified", FlagStartTime)
	} else if timeUnix, err := strconv.ParseInt(timeStr, 10, 64); err == nil { // unix time
		startTime = time.Unix(timeUnix, 0)
	} else if timeRFC, err := time.Parse(time.RFC3339, timeStr); err == nil { // RFC time
		startTime = timeRFC
	} else { // invalid input
		return types.MsgCreateGauge_LiquidityMetaData{}, errors.New("invalid start time format")
	}

	lockDuration, err := cmd.Flags().GetDuration(FlagLockDuration)
	if err != nil {
		return types.MsgCreateGauge_LiquidityMetaData{}, err
	}
	liquidityGaugeExtraData := types.MsgCreateGauge_LiquidityMetaData{
		LiquidityMetaData: &types.LiquidtyGaugeMetaData{
			PoolId:       poolId,
			IsMasterPool: isMasterPool,
			ChildPoolIds: childPoolIds,
			StartTime:    startTime,
			LockDuration: lockDuration,
		},
	}
	return liquidityGaugeExtraData, nil
}
