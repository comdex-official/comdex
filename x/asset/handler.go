package asset

//goland:noinspection GoLinter
import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/comdex-official/comdex/x/asset/keeper"
	"github.com/comdex-official/comdex/x/asset/types"
)

// NewHandler ...
func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		_ = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
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
		case *types.UpdatePairProposal:
			return handleUpdatePairProposal(ctx, k, c)
		case *types.UpdateGovTimeInAppProposal:
			return handleUpdateGovTimeInAppProposal(ctx, k, c)
		case *types.AddAppProposal:
			return handleAddAppProposal(ctx, k, c)
		case *types.AddAssetInAppProposal:
			return handleAddAssetInAppProposal(ctx, k, c)

		default:
			return sdkerrors.Wrapf(types.ErrorUnknownProposalType, "%T", c)
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

func handleUpdatePairProposal(ctx sdk.Context, k keeper.Keeper, p *types.UpdatePairProposal) error {
	return k.HandleProposalUpdatePair(ctx, p)
}

func handleUpdateGovTimeInAppProposal(ctx sdk.Context, k keeper.Keeper, p *types.UpdateGovTimeInAppProposal) error {
	return k.HandleUpdateGovTimeInApp(ctx, p)
}

func handleAddAppProposal(ctx sdk.Context, k keeper.Keeper, p *types.AddAppProposal) error {
	return k.HandleAddAppRecords(ctx, p)
}

func handleAddAssetInAppProposal(ctx sdk.Context, k keeper.Keeper, p *types.AddAssetInAppProposal) error {
	return k.HandleAddAssetInAppRecords(ctx, p)
}
