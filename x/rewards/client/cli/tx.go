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

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	// "github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/comdex-official/comdex/x/rewards/types"
)

var (
	DefaultRelativePacketTimeoutTimestamp = uint64((time.Duration(10) * time.Minute).Nanoseconds())
)

const (
	flagPacketTimeoutTimestamp = "packet-timeout-timestamp"
	listSeparator              = ","
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

	// this line is used by starport scaffolding # 1
	cmd.AddCommand(
		NewCreateGaugeCmd(),
		txWhitelistAsset(),
		txRemoveWhitelistAsset(),
		txWhitelistAppIdVault(),
		txRemoveWhitelistAppIdVault(),
		txActivateExternalRewardsLockers(),
		txActivateExternalVaultsLockers(),
	)

	return cmd
}

// NewCreateGaugeCmd implemets create-gauge cli transaction command.
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
				msg.AppId = appID
				msg.Kind = &gaugeExtraData
			}

			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msg)
		},
	}

	cmd.Flags().AddFlagSet(FlagSetCreateGauge())
	flags.AddTxFlagsToCmd(cmd)
	_ = cmd.MarkFlagRequired(FlagStartTime)
	return cmd
}

func txWhitelistAsset() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "whitelist-asset [app_mapping_Id] [asset_Id]",
		Short: "Add Whitelisted assetId for Locker savings rewards",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			app_mapping_Id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			asset_Id, err := ParseUint64SliceFromString(args[1], ",")
			if err != nil {
				return err
			}

			var newAssetId []uint64
			for i := range asset_Id {
				newAssetId = append(newAssetId, asset_Id[i])
			}

			msg := types.NewMsgWhitelistAsset(app_mapping_Id, ctx.GetFromAddress(), newAssetId)

			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd

}

func txRemoveWhitelistAsset() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-whitelist-asset [app_mapping_Id] [asset_Id]",
		Short: "Remove Whitelisted assetId for Locker savings rewards",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			app_mapping_Id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			asset_Id, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgRemoveWhitelistAsset(app_mapping_Id, ctx.GetFromAddress(), asset_Id)

			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd

}

func txWhitelistAppIdVault() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "whitelist-app-id-vault-interest [app_mapping_Id]",
		Short: "whitelist app id vault interest",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			appMappingId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgWhitelistAppIdVault(
				appMappingId,
				ctx.GetFromAddress(),
			)

			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd

}

func txRemoveWhitelistAppIdVault() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-whitelist-app-id-vault-interest [app_mapping_Id] ",
		Short: "remove whitelist app id vault interest",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			appMappingId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgRemoveWhitelistAppIdVault(
				appMappingId,
				ctx.GetFromAddress(),
			)

			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd

}

func txActivateExternalRewardsLockers() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "activate-external-rewards-locker [app_mapping_Id] [asset_id] [total_rewards] [duration_days] [min_lockup_time_seconds]",
		Short: "activate external rewards for locker",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			appMappingId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			asset_Id, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			total_rewards, err := sdk.ParseCoinNormalized(args[2])
			if err != nil {
				return err
			}

			duration_days, err := strconv.ParseInt(args[3], 10, 64)
			if err != nil {
				return err
			}

			min_lockup_time_seconds, err := strconv.ParseInt(args[4], 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgActivateExternalRewardsLockers(
				appMappingId,
				asset_Id,
				total_rewards,
				duration_days,
				min_lockup_time_seconds,
				ctx.GetFromAddress(),
			)

			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func txActivateExternalVaultsLockers() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "activate-external-rewards-vault [app_mapping_Id] [extended_pair_id] [total_rewards] [duration_days] [min_lockup_time_seconds]",
		Short: "activate external reward for vault extended_pair_id",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			appMappingId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			extended_pair_id, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			total_rewards, err := sdk.ParseCoinNormalized(args[2])
			if err != nil {
				return err
			}

			duration_days, err := strconv.ParseInt(args[3], 10, 64)
			if err != nil {
				return err
			}

			min_lockup_time_seconds, err := strconv.ParseInt(args[4], 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgActivateExternalVaultLockers(
				appMappingId,
				extended_pair_id,
				total_rewards,
				duration_days,
				min_lockup_time_seconds,
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
	childPoolIds := []uint64{}
	if childPoolIdsCombined != "" {
		childPoolIdsStrs := strings.Split(childPoolIdsCombined, ",")
		for _, poolIDStr := range childPoolIdsStrs {
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
