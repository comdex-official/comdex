package lend

import (
	"fmt"

	bam "github.com/cosmos/cosmos-sdk/baseapp"

	errorsmod "cosmossdk.io/errors"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/pkg/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/comdex-official/comdex/x/lend/keeper"
	"github.com/comdex-official/comdex/x/lend/types"
)

// NewHandler ...
func NewHandler(k keeper.Keeper) bam.MsgServiceHandler {
	server := keeper.NewMsgServerImpl(k)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgLend:
			res, err := server.Lend(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgWithdraw:
			res, err := server.Withdraw(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgDeposit:
			res, err := server.Deposit(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgCloseLend:
			res, err := server.CloseLend(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgBorrow:
			res, err := server.Borrow(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgDraw:
			res, err := server.Draw(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgDepositBorrow:
			res, err := server.DepositBorrow(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgCloseBorrow:
			res, err := server.CloseBorrow(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgRepay:
			res, err := server.Repay(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgBorrowAlternate:
			res, err := server.BorrowAlternate(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgFundModuleAccounts:
			res, err := server.FundModuleAccounts(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgCalculateInterestAndRewards:
			res, err := server.CalculateInterestAndRewards(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgFundReserveAccounts:
			res, err := server.FundReserveAccounts(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgRepayWithdraw:
			res, err := server.RepayWithdraw(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, errorsmod.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func NewLendHandler(k keeper.Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case *types.LendPairsProposal:
			return handleAddWhitelistedPairsProposal(ctx, k, c)
		case *types.MultipleLendPairsProposal:
			return handleAddMultipleWhitelistedPairsProposal(ctx, k, c)
		case *types.AddPoolsProposal:
			return handleAddPoolProposal(ctx, k, c)
		case *types.AddAssetToPairProposal:
			return handleAddAssetToPairProposal(ctx, k, c)
		case *types.AddMultipleAssetToPairProposal:
			return handleAddMultipleAssetToPairProposal(ctx, k, c)
		case *types.AddAssetRatesParams:
			return handleAddAssetRatesParamsProposal(ctx, k, c)
		case *types.AddAuctionParamsProposal:
			return HandleAddAuctionParamsProposal(ctx, k, c)
		case *types.AddPoolPairsProposal:
			return handleAddPoolPairsProposal(ctx, k, c)
		case *types.AddAssetRatesPoolPairsProposal:
			return handleAddAssetRatesPoolPairsProposal(ctx, k, c)
		case *types.AddPoolDepreciateProposal:
			return handlePoolDepreciateProposal(ctx, k, c)
		case *types.AddEModePairsProposal:
			return handleEModePairsProposal(ctx, k, c)

		default:
			return errors.Wrapf(types.ErrorUnknownProposalType, "%T", c)
		}
	}
}

func handleAddWhitelistedPairsProposal(ctx sdk.Context, k keeper.Keeper, p *types.LendPairsProposal) error {
	return k.HandleAddWhitelistedPairsRecords(ctx, p)
}

func handleAddMultipleWhitelistedPairsProposal(ctx sdk.Context, k keeper.Keeper, p *types.MultipleLendPairsProposal) error {
	return k.HandleMultipleAddWhitelistedPairsRecords(ctx, p)
}

func handleAddPoolProposal(ctx sdk.Context, k keeper.Keeper, p *types.AddPoolsProposal) error {
	return k.HandleAddPoolRecords(ctx, p)
}

func handleAddAssetToPairProposal(ctx sdk.Context, k keeper.Keeper, p *types.AddAssetToPairProposal) error {
	return k.HandleAddAssetToPairRecords(ctx, p)
}

func handleAddMultipleAssetToPairProposal(ctx sdk.Context, k keeper.Keeper, p *types.AddMultipleAssetToPairProposal) error {
	return k.HandleAddMultipleAssetToPairRecords(ctx, p)
}

func handleAddAssetRatesParamsProposal(ctx sdk.Context, k keeper.Keeper, p *types.AddAssetRatesParams) error {
	return k.HandleAddAssetRatesParamsRecords(ctx, p)
}

func HandleAddAuctionParamsProposal(ctx sdk.Context, k keeper.Keeper, p *types.AddAuctionParamsProposal) error {
	return k.HandleAddAuctionParamsRecords(ctx, p)
}

func handleAddPoolPairsProposal(ctx sdk.Context, k keeper.Keeper, p *types.AddPoolPairsProposal) error {
	return k.HandleAddPoolPairsRecords(ctx, p)
}

func handleAddAssetRatesPoolPairsProposal(ctx sdk.Context, k keeper.Keeper, p *types.AddAssetRatesPoolPairsProposal) error {
	return k.HandleAddAssetRatesPoolPairsRecords(ctx, p)
}

func handlePoolDepreciateProposal(ctx sdk.Context, k keeper.Keeper, p *types.AddPoolDepreciateProposal) error {
	return k.HandlePoolDepreciateProposal(ctx, p)
}

func handleEModePairsProposal(ctx sdk.Context, k keeper.Keeper, p *types.AddEModePairsProposal) error {
	return k.HandleEModePairsProposal(ctx, p)
}
