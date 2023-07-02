package cli

import (
	"github.com/comdex-official/comdex/x/nft/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                types.ModuleName,
		Short:              "Querying commands for the NFT module",
		DisableFlagParsing: true,
	}

	queryCmd.AddCommand(
		QueryNFT(),
		QueryOwnerNFTs(),
	)

	return queryCmd
}

func QueryNFT() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nft [denom-id] [nft-id]",
		Short: "Query nft by denom-id and nft-id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			denomID := args[0]
			nftID := args[1]

			queryClient := types.NewQueryClient(ctx)

			res, err := queryClient.NFT(cmd.Context(), &types.QueryNFTRequest{
				DenomId: denomID,
				Id:      nftID,
			})
			if err != nil {
				return err
			}
			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func QueryOwnerNFTs() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "owner-nft [denom-id] [owner]",
		Short: "Query nft by denom-id and owner",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			denomID := args[0]
			owner := args[1]

			queryClient := types.NewQueryClient(ctx)

			res, err := queryClient.OwnerNFTs(cmd.Context(), &types.QueryOwnerNFTsRequest{
				DenomId:    denomID,
				Owner:      owner,
				Pagination: nil,
			})
			if err != nil {
				return err
			}
			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
