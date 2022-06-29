package asset

import (
	"github.com/comdex-official/comdex/x/asset/keeper"
	"github.com/comdex-official/comdex/x/asset/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

func NewHandler(k keeper.Keeper) sdk.Handler {
	server := keeper.NewMsgServiceServer(k)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgAddAssetRequest:
			res, err := server.MsgAddAsset(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgUpdateAssetRequest:
			res, err := server.MsgUpdateAsset(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgAddPairRequest:
			res, err := server.MsgAddPair(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		default:
			return nil, errors.Wrapf(types.ErrorUnknownMsgType, "%T", msg)
		}
	}
}

func NewUpdateAssetProposalHandler(k keeper.Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case *types.AddAssetsProposal:
			return handleAddAssetProposal(ctx, k, c)
		case *types.UpdateAssetProposal:
			return handleUpdateAssetProposal(ctx, k, c)
		case *types.AddPairsProposal:
			return handleAddPairsProposal(ctx, k, c)
		case *types.UpdateGovTimeInAppMappingProposal:
			return handleUpdateGovTimeInAppMappingProposal(ctx, k, c)
		case *types.AddAppMappingProposal:
			return handleAddAppMappingProposal(ctx, k, c)
		case *types.AddAssetMappingProposal:
			return handleAddAssetMappingProposal(ctx, k, c)
		case *types.AddExtendedPairsVaultProposal:
			return handleExtendedPairsVaultProposal(ctx, k, c)

		default:
			return errors.Wrapf(types.ErrorUnknownProposalType, "%T", c)
		}
	}
}

func handleAddAssetProposal(ctx sdk.Context, k keeper.Keeper, p *types.AddAssetsProposal) error {
	return k.HandleProposalAddAsset(ctx, p)
}

func handleUpdateAssetProposal(ctx sdk.Context, k keeper.Keeper, p *types.UpdateAssetProposal) error {
	return k.HandleProposalUpdateAsset(ctx, p)
}

func handleAddPairsProposal(ctx sdk.Context, k keeper.Keeper, p *types.AddPairsProposal) error {
	return k.HandleProposalAddPairs(ctx, p)
}

func handleUpdateGovTimeInAppMappingProposal(ctx sdk.Context, k keeper.Keeper, p *types.UpdateGovTimeInAppMappingProposal) error {
	return k.HandleUpdateGovTimeInAppMapping(ctx, p)
}

func handleAddAppMappingProposal(ctx sdk.Context, k keeper.Keeper, p *types.AddAppMappingProposal) error {
	return k.HandleAddAppMappingRecords(ctx, p)
}
func handleAddAssetMappingProposal(ctx sdk.Context, k keeper.Keeper, p *types.AddAssetMappingProposal) error {
	return k.HandleAddAssetMappingRecords(ctx, p)
}

func handleExtendedPairsVaultProposal(ctx sdk.Context, k keeper.Keeper, p *types.AddExtendedPairsVaultProposal) error {
	return k.HandleAddExtendedPairsVaultRecords(ctx, p)
}
