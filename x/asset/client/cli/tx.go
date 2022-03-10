package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/spf13/cobra"
	"strconv"

	"github.com/comdex-official/comdex/x/asset/types"
)

func txAddAsset() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-asset [name] [denom] [decimals]",
		Short: "Add an asset",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			decimals, err := strconv.ParseInt(args[2], 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgAddAssetRequest(
				ctx.FromAddress,
				args[0],
				args[1],
				decimals,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func txUpdateAsset() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-asset [id]",
		Short: "Update an asset",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			name, err := cmd.Flags().GetString(flagName)
			if err != nil {
				return err
			}

			denom, err := cmd.Flags().GetString(flagDenom)
			if err != nil {
				return err
			}

			decimals, err := cmd.Flags().GetInt64(flagDecimals)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateAssetRequest(
				ctx.FromAddress,
				id,
				name,
				denom,
				decimals,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	cmd.Flags().String(flagName, "", "name")
	cmd.Flags().String(flagDenom, "", "denomination")
	cmd.Flags().Int64(flagDecimals, -1, "decimals")

	return cmd
}

func txAddPair() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-pair [asset-in] [asset-out] [liquidation-ratio]",
		Short: "Add a pair",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			assetIn, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			assetOut, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			liquidationRatio, err := sdk.NewDecFromStr(args[2])
			if err != nil {
				return err
			}

			msg := types.NewMsgAddPairRequest(
				ctx.FromAddress,
				assetIn,
				assetOut,
				liquidationRatio,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func txUpdatePair() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-pair [id]",
		Short: "Update a pair",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			liquidationRatio, err := GetLiquidationRatio(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdatePairRequest(
				ctx.FromAddress,
				id,
				liquidationRatio,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	cmd.Flags().String(flagLiquidationRatio, "", "liquidation ratio")

	return cmd
}

func NewCmdSubmitAddAssetsProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-assets [name] [Denom] [Decimals]",
		Args:  cobra.ExactArgs(3),
		Short: "Submit assets",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			names, err := ParseStringFromString(args[0], ",")
			if err != nil {
				return err
			}
			denoms, err :=  ParseStringFromString(args[1], ",")
			if err != nil {
				return err
			}

			decimals ,err := ParseUint64SliceFromString(args[2], ",")
			if err != nil {
				return err
			}

			title, err := cmd.Flags().GetString(cli.FlagTitle)
			if err != nil {
				return err
			}

			description, err := cmd.Flags().GetString(cli.FlagDescription)
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()

			var assets []types.Asset
			for i, _ := range names {
				assets = append(assets, types.Asset{
					Name: names[i],
					Denom: denoms[i],
					Decimals: decimals[i],
				})
			}

			depositStr, err := cmd.Flags().GetString(cli.FlagDeposit)
			if err != nil {
				return err
			}
			deposit, err := sdk.ParseCoinsNormalized(depositStr)
			if err != nil {
				return err
			}

			content := types.NewUpdateLiquidationRatioProposal(title, description, assets)

			msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
			if err != nil {
				return err
			}

			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(cli.FlagTitle, "", "title of proposal")
	cmd.Flags().String(cli.FlagDescription, "", "description of proposal")
	cmd.Flags().String(cli.FlagDeposit, "", "deposit of proposal")
	_ = cmd.MarkFlagRequired(cli.FlagTitle)
	_ = cmd.MarkFlagRequired(cli.FlagDescription)

	return cmd
}
