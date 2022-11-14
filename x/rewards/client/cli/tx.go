package cli

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"

	"github.com/comdex-official/comdex/x/rewards/types"
)

// GetTxCmd returns the transaction commands for this module .
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "rewards",
		Short:                      fmt.Sprintf("%s transactions subcommands", "rewards"),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		NewCreateGaugeCmd(),
		txActivateExternalRewardsLockers(),
		txActivateExternalRewardsVaults(),
		txActivateExternalRewardsLend(),
	)

	return cmd
}

// NewCreateGaugeCmd implements create-gauge cli transaction command.
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

			gaugeTypeID, err := strconv.ParseUint(args[0], 10, 64)
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

			var startTime time.Time
			timeStr, err := cmd.Flags().GetString(FlagStartTime)
			if err != nil {
				return err
			}
			if timeStr == "" { // empty start time
				return fmt.Errorf("%s cannot be empty string", FlagStartTime)
			} else if timeUnix, err := strconv.ParseInt(timeStr, 10, 64); err == nil { // unix time
				startTime = time.Unix(timeUnix, 0)
			} else if timeRFC, err := time.Parse(time.RFC3339, timeStr); err == nil { // RFC time
				startTime = timeRFC
			} else { // invalid input
				return errors.New("invalid start time format")
			}

			msg := types.NewMsgCreateGauge(
				0,
				clientCtx.GetFromAddress(),
				startTime,
				gaugeTypeID,
				triggerDuration,
				depositAmount,
				totalTriggers,
			)

			switch msg.GaugeTypeId {
			case types.LiquidityGaugeTypeID:
				gaugeExtraData, err := NewBuildLiquidityGaugeExtraData(cmd)
				if err != nil {
					return err
				}
				appID, err := cmd.Flags().GetUint64(FlagAppID)
				if err != nil {
					return err
				}
				msg.AppId = appID
				msg.Kind = &gaugeExtraData
			}

			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msg)
		},
	}

	cmd.Flags().AddFlagSet(FlagSetCreateGauge())
	flags.AddTxFlagsToCmd(cmd)
	_ = cmd.MarkFlagRequired(FlagStartTime)
	return cmd
}

func txActivateExternalRewardsLockers() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "activate-external-rewards-locker [appMappingID] [asset_id] [totalRewards] [durationDays] [minLockupTimeSeconds]",
		Short: "activate external rewards for locker",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			appMappingID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			assetID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			totalRewards, err := sdk.ParseCoinNormalized(args[2])
			if err != nil {
				return err
			}

			durationDays, err := strconv.ParseInt(args[3], 10, 64)
			if err != nil {
				return err
			}

			minLockupTimeSeconds, err := strconv.ParseInt(args[4], 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgActivateExternalRewardsLockers(
				appMappingID,
				assetID,
				totalRewards,
				durationDays,
				minLockupTimeSeconds,
				ctx.GetFromAddress(),
			)

			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func txActivateExternalRewardsVaults() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "activate-external-rewards-vault [appMappingID] [extendedPairID] [totalRewards] [durationDays] [minLockupTimeSeconds]",
		Short: "activate external reward for vault extendedPairID",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			appMappingID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			extendedPairID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			totalRewards, err := sdk.ParseCoinNormalized(args[2])
			if err != nil {
				return err
			}

			durationDays, err := strconv.ParseInt(args[3], 10, 64)
			if err != nil {
				return err
			}

			minLockupTimeSeconds, err := strconv.ParseInt(args[4], 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgActivateExternalRewardsVault(
				appMappingID,
				extendedPairID,
				totalRewards,
				durationDays,
				minLockupTimeSeconds,
				ctx.GetFromAddress(),
			)

			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// NewBuildLiquidityGaugeExtraData sanitizes cli input data for MsgCreateGauge_LiquidityMetaData.
func NewBuildLiquidityGaugeExtraData(cmd *cobra.Command) (types.MsgCreateGauge_LiquidityMetaData, error) {
	poolID, err := cmd.Flags().GetUint64(FlagPoolID)
	if err != nil {
		return types.MsgCreateGauge_LiquidityMetaData{}, err
	}
	if poolID == 0 {
		return types.MsgCreateGauge_LiquidityMetaData{}, fmt.Errorf("%s required but not specified / pool-id cannot be 0", FlagPoolID)
	}

	appID, err := cmd.Flags().GetUint64(FlagAppID)
	if err != nil {
		return types.MsgCreateGauge_LiquidityMetaData{}, err
	}
	if appID == 0 {
		return types.MsgCreateGauge_LiquidityMetaData{}, fmt.Errorf("%s required but not specified / app-id cannot be 0", FlagAppID)
	}

	isMasterPool, err := cmd.Flags().GetBool(FlagIsMasterPool)
	if err != nil {
		return types.MsgCreateGauge_LiquidityMetaData{}, err
	}

	childPoolIdsCombined, err := cmd.Flags().GetString(FlagChildPoolIds)
	if err != nil {
		return types.MsgCreateGauge_LiquidityMetaData{}, err
	}
	var childPoolIds []uint64
	if childPoolIdsCombined != "" {
		childPoolIdsStr := strings.Split(childPoolIdsCombined, ",")
		for _, poolIDStr := range childPoolIdsStr {
			poolID, err := strconv.ParseUint(poolIDStr, 10, 64)
			if err != nil {
				return types.MsgCreateGauge_LiquidityMetaData{}, err
			}
			childPoolIds = append(childPoolIds, poolID)
		}
	}

	liquidityGaugeExtraData := types.MsgCreateGauge_LiquidityMetaData{
		LiquidityMetaData: &types.LiquidtyGaugeMetaData{
			PoolId:       poolID,
			IsMasterPool: isMasterPool,
			ChildPoolIds: childPoolIds,
		},
	}
	return liquidityGaugeExtraData, nil
}

func txActivateExternalRewardsLend() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "activate-external-rewards-lend [appMappingID] [cPoolID] [assetIDs] [cSwapAppID] [cSwapMinLockAmount] [totalRewards] [masterPoolID] [duration] [minLockupTimeSeconds]",
		Short: "activate external reward for borrowers",
		Args:  cobra.ExactArgs(9),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			appMappingID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			cPoolID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			assetIDs, err := ParseUint64SliceFromString(args[2], ",")
			if err != nil {
				return err
			}

			cSwapAppID, err := strconv.ParseUint(args[3], 10, 64)
			if err != nil {
				return err
			}

			cSwapMinLockAmount, err := strconv.ParseUint(args[4], 10, 64)
			if err != nil {
				return err
			}

			totalRewards, err := sdk.ParseCoinNormalized(args[5])
			if err != nil {
				return err
			}

			masterPoolID, err := strconv.ParseInt(args[6], 10, 64)
			if err != nil {
				return err
			}

			durationDays, err := strconv.ParseInt(args[7], 10, 64)
			if err != nil {
				return err
			}

			minLockupTimeSeconds, err := strconv.ParseInt(args[8], 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgActivateExternalRewardsLend(
				appMappingID,
				cPoolID,
				assetIDs,
				cSwapAppID,
				cSwapMinLockAmount,
				totalRewards,
				masterPoolID,
				durationDays,
				minLockupTimeSeconds,
				ctx.GetFromAddress(),
			)

			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
