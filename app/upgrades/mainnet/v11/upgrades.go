package v11

import (
	assetkeeper "github.com/comdex-official/comdex/x/asset/keeper"
	auctiontypes "github.com/comdex-official/comdex/x/auction/types"
	lendtypes "github.com/comdex-official/comdex/x/lend/types"
	liquiditykeeper "github.com/comdex-official/comdex/x/liquidity/keeper"
	liquiditytypes "github.com/comdex-official/comdex/x/liquidity/types"
	lockertypes "github.com/comdex-official/comdex/x/locker/types"
	rewardskeeper "github.com/comdex-official/comdex/x/rewards/keeper"
	rewardstypes "github.com/comdex-official/comdex/x/rewards/types"
	tokenminttypes "github.com/comdex-official/comdex/x/tokenmint/types"
	vaulttypes "github.com/comdex-official/comdex/x/vault/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	icahostkeeper "github.com/cosmos/ibc-go/v4/modules/apps/27-interchain-accounts/host/keeper"
	icahosttypes "github.com/cosmos/ibc-go/v4/modules/apps/27-interchain-accounts/host/types"
	icatypes "github.com/cosmos/ibc-go/v4/modules/apps/27-interchain-accounts/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v4/modules/apps/transfer/types"
)

// An error occurred during the creation of the CMST/STJUNO pair, as it was mistakenly created in the Harbor app (ID-2) instead of the cSwap app (ID-1).
// As a result, the transaction fee was charged to the creator of the pair, who is entitled to a refund.
// The provided code is designed to initiate the refund process.
// The transaction hash for the pair creation is EF408AD53B8BB0469C2A593E4792CB45552BD6495753CC2C810A1E4D82F3982F.
// MintSan - https://www.mintscan.io/comdex/txs/EF408AD53B8BB0469C2A593E4792CB45552BD6495753CC2C810A1E4D82F3982F

func RefundFeeForAccidentallyCreatedPirToOwner(
	ctx sdk.Context,
	liquidityKeeper liquiditykeeper.Keeper,
	assetKeeper assetkeeper.Keeper,
	bankKeeper bankkeeper.Keeper,
) {
	allApps, found := assetKeeper.GetApps(ctx)
	if !found {
		panic("apps not found")
	}
	harborAppID := 0
	for _, app := range allApps {
		if app.Name == "harbor" {
			harborAppID = int(app.Id)
			break
		}
	}
	if harborAppID == 0 {
		panic("harbor app not found")
	}
	liquidityParams, err := liquidityKeeper.GetGenericParams(ctx, uint64(harborAppID))
	if err != nil {
		panic(err.Error())
	}

	refundAmount := sdk.NewCoin("ucmdx", sdk.NewInt(2000000000))

	feeCollectorAddress := sdk.MustAccAddressFromBech32(liquidityParams.FeeCollectorAddress)
	feeCollectorCmdxBalance := bankKeeper.GetBalance(ctx, feeCollectorAddress, "ucmdx")
	if feeCollectorCmdxBalance.Amount.LT(refundAmount.Amount) {
		panic("isufficient balance in fee collector of harbor app")
	}

	pairCreatorAddress := sdk.MustAccAddressFromBech32("comdex19wmd9xzhjnvmr90z8r3cnjlns5kgem9qglt2ll")
	pairCreatorCmdxBalance := bankKeeper.GetBalance(ctx, pairCreatorAddress, "ucmdx")
	err = bankKeeper.SendCoins(ctx, feeCollectorAddress, pairCreatorAddress, sdk.NewCoins(refundAmount))
	if err != nil {
		panic(err)
	}
	pairCreatorCmdxBalanceNew := bankKeeper.GetBalance(ctx, pairCreatorAddress, "ucmdx")

	if !pairCreatorCmdxBalance.Add(refundAmount).IsEqual(pairCreatorCmdxBalanceNew) {
		panic("account balance invariant after pair creation fee refund")
	}
}

func CreateUpgradeHandlerV11(
	mm *module.Manager,
	configurator module.Configurator,
	liquidityKeeper liquiditykeeper.Keeper,
	assetKeeper assetkeeper.Keeper,
	bankKeeper bankkeeper.Keeper,
	accountKeeper authkeeper.AccountKeeper,
	rewardsKeeper rewardskeeper.Keeper,
	icahostkeeper icahostkeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {

		ctx.Logger().Info("Applying main net upgrade - v.11.0.1")

		fromVM[icatypes.ModuleName] = mm.Modules[icatypes.ModuleName].ConsensusVersion()

		hostParams := icahosttypes.Params{
			HostEnabled: true,
			AllowMessages: []string{
				sdk.MsgTypeURL(&ibctransfertypes.MsgTransfer{}),
				sdk.MsgTypeURL(&banktypes.MsgSend{}),
				sdk.MsgTypeURL(&stakingtypes.MsgDelegate{}),
				sdk.MsgTypeURL(&stakingtypes.MsgBeginRedelegate{}),
				sdk.MsgTypeURL(&stakingtypes.MsgCreateValidator{}),
				sdk.MsgTypeURL(&stakingtypes.MsgEditValidator{}),
				sdk.MsgTypeURL(&stakingtypes.MsgUndelegate{}),
				sdk.MsgTypeURL(&distrtypes.MsgWithdrawDelegatorReward{}),
				sdk.MsgTypeURL(&distrtypes.MsgSetWithdrawAddress{}),
				sdk.MsgTypeURL(&distrtypes.MsgWithdrawValidatorCommission{}),
				sdk.MsgTypeURL(&distrtypes.MsgFundCommunityPool{}),
				sdk.MsgTypeURL(&govtypes.MsgVote{}),
				sdk.MsgTypeURL(&auctiontypes.MsgPlaceSurplusBidRequest{}),
				sdk.MsgTypeURL(&auctiontypes.MsgPlaceDebtBidRequest{}),
				sdk.MsgTypeURL(&auctiontypes.MsgPlaceDutchBidRequest{}),
				sdk.MsgTypeURL(&auctiontypes.MsgPlaceDutchLendBidRequest{}),
				sdk.MsgTypeURL(&lendtypes.MsgLend{}),
				sdk.MsgTypeURL(&lendtypes.MsgWithdraw{}),
				sdk.MsgTypeURL(&lendtypes.MsgDeposit{}),
				sdk.MsgTypeURL(&lendtypes.MsgCloseLend{}),
				sdk.MsgTypeURL(&lendtypes.MsgBorrow{}),
				sdk.MsgTypeURL(&lendtypes.MsgDraw{}),
				sdk.MsgTypeURL(&lendtypes.MsgRepay{}),
				sdk.MsgTypeURL(&lendtypes.MsgDepositBorrow{}),
				sdk.MsgTypeURL(&lendtypes.MsgCloseBorrow{}),
				sdk.MsgTypeURL(&lendtypes.MsgBorrowAlternate{}),
				sdk.MsgTypeURL(&lendtypes.MsgCalculateInterestAndRewards{}),
				sdk.MsgTypeURL(&liquiditytypes.MsgCreatePair{}),
				sdk.MsgTypeURL(&liquiditytypes.MsgCreatePool{}),
				sdk.MsgTypeURL(&liquiditytypes.MsgCreateRangedPool{}),
				sdk.MsgTypeURL(&liquiditytypes.MsgDeposit{}),
				sdk.MsgTypeURL(&liquiditytypes.MsgWithdraw{}),
				sdk.MsgTypeURL(&liquiditytypes.MsgLimitOrder{}),
				sdk.MsgTypeURL(&liquiditytypes.MsgMarketOrder{}),
				sdk.MsgTypeURL(&liquiditytypes.MsgMMOrder{}),
				sdk.MsgTypeURL(&liquiditytypes.MsgCancelOrder{}),
				sdk.MsgTypeURL(&liquiditytypes.MsgCancelAllOrders{}),
				sdk.MsgTypeURL(&liquiditytypes.MsgCancelMMOrder{}),
				sdk.MsgTypeURL(&liquiditytypes.MsgFarm{}),
				sdk.MsgTypeURL(&liquiditytypes.MsgUnfarm{}),
				sdk.MsgTypeURL(&lockertypes.MsgCreateLockerRequest{}),
				sdk.MsgTypeURL(&lockertypes.MsgDepositAssetRequest{}),
				sdk.MsgTypeURL(&lockertypes.MsgWithdrawAssetRequest{}),
				sdk.MsgTypeURL(&lockertypes.MsgCloseLockerRequest{}),
				sdk.MsgTypeURL(&lockertypes.MsgLockerRewardCalcRequest{}),
				sdk.MsgTypeURL(&rewardstypes.MsgCreateGauge{}),
				sdk.MsgTypeURL(&tokenminttypes.MsgMintNewTokensRequest{}),
				sdk.MsgTypeURL(&vaulttypes.MsgCreateRequest{}),
				sdk.MsgTypeURL(&vaulttypes.MsgDepositRequest{}),
				sdk.MsgTypeURL(&vaulttypes.MsgWithdrawRequest{}),
				sdk.MsgTypeURL(&vaulttypes.MsgDrawRequest{}),
				sdk.MsgTypeURL(&vaulttypes.MsgRepayRequest{}),
				sdk.MsgTypeURL(&vaulttypes.MsgCloseRequest{}),
				sdk.MsgTypeURL(&vaulttypes.MsgDepositAndDrawRequest{}),
				sdk.MsgTypeURL(&vaulttypes.MsgCreateStableMintRequest{}),
				sdk.MsgTypeURL(&vaulttypes.MsgDepositStableMintRequest{}),
				sdk.MsgTypeURL(&vaulttypes.MsgWithdrawStableMintRequest{}),
				sdk.MsgTypeURL(&vaulttypes.MsgVaultInterestCalcRequest{}),
			},
		}
		icahostkeeper.SetParams(ctx, hostParams)

		vm, err := mm.RunMigrations(ctx, configurator, fromVM)
		if err != nil {
			return nil, err
		}
		RefundFeeForAccidentallyCreatedPirToOwner(ctx, liquidityKeeper, assetKeeper, bankKeeper)
		DistributeRewards(ctx, accountKeeper, bankKeeper, rewardsKeeper)
		return vm, err
	}
}
