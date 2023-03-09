package v10

import (
	assetkeeper "github.com/comdex-official/comdex/x/asset/keeper"
	auctiontypes "github.com/comdex-official/comdex/x/auction/types"
	esmtypes "github.com/comdex-official/comdex/x/esm/types"
	lendtypes "github.com/comdex-official/comdex/x/lend/types"
	liquidationtypes "github.com/comdex-official/comdex/x/liquidation/types"
	liquiditykeeper "github.com/comdex-official/comdex/x/liquidity/keeper"
	liquiditytypes "github.com/comdex-official/comdex/x/liquidity/types"
	lockertypes "github.com/comdex-official/comdex/x/locker/types"
	rewardstypes "github.com/comdex-official/comdex/x/rewards/types"
	tokenminttypes "github.com/comdex-official/comdex/x/tokenmint/types"
	vaulttypes "github.com/comdex-official/comdex/x/vault/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	ica "github.com/cosmos/ibc-go/v4/modules/apps/27-interchain-accounts"
	icacontrollertypes "github.com/cosmos/ibc-go/v4/modules/apps/27-interchain-accounts/controller/types"
	icahosttypes "github.com/cosmos/ibc-go/v4/modules/apps/27-interchain-accounts/host/types"
	icatypes "github.com/cosmos/ibc-go/v4/modules/apps/27-interchain-accounts/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v4/modules/apps/transfer/types"
)

func DeleteAccidentallyCreatedPairAndRefundPirCreationFeeToOwner(
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

func CreateUpgradeHandlerV10(
	mm *module.Manager,
	configurator module.Configurator,
	liquidityKeeper liquiditykeeper.Keeper,
	assetKeeper assetkeeper.Keeper,
	bankKeeper bankkeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		fromVM[icatypes.ModuleName] = mm.Modules[icatypes.ModuleName].ConsensusVersion()

		// create ICS27 Controller submodule params, controller module not enabled.
		controllerParams := icacontrollertypes.Params{}

		// create ICS27 Host submodule params

		// create ICS27 Host submodule params
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
				sdk.MsgTypeURL(&esmtypes.MsgDepositESM{}),
				sdk.MsgTypeURL(&esmtypes.MsgExecuteESM{}),
				sdk.MsgTypeURL(&esmtypes.MsgKillRequest{}),
				sdk.MsgTypeURL(&esmtypes.MsgCollateralRedemptionRequest{}),
				sdk.MsgTypeURL(&esmtypes.MsgCollateralRedemptionRequest{}),
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
				sdk.MsgTypeURL(&lendtypes.MsgFundModuleAccounts{}),
				sdk.MsgTypeURL(&lendtypes.MsgCalculateInterestAndRewards{}),
				sdk.MsgTypeURL(&lendtypes.MsgFundReserveAccounts{}),
				sdk.MsgTypeURL(&liquidationtypes.MsgLiquidateVaultRequest{}),
				sdk.MsgTypeURL(&liquidationtypes.MsgLiquidateBorrowRequest{}),
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
				sdk.MsgTypeURL(&rewardstypes.ActivateExternalRewardsLockers{}),
				sdk.MsgTypeURL(&rewardstypes.ActivateExternalRewardsVault{}),
				sdk.MsgTypeURL(&rewardstypes.ActivateExternalRewardsLend{}),
				sdk.MsgTypeURL(&rewardstypes.ActivateExternalRewardsStableMint{}),
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
		// No changes in existing module and their states,
		// This upgrades adds new modules and new states in the existing store

		icamodule, correctTypecast := mm.Modules[icatypes.ModuleName].(ica.AppModule)
		if !correctTypecast {
			panic("mm.Modules[icatypes.ModuleName] is not of type ica.AppModule")
		}
		icamodule.InitModule(ctx, controllerParams, hostParams)

		vm, err := mm.RunMigrations(ctx, configurator, fromVM)
		if err != nil {
			return nil, err
		}
		DeleteAccidentallyCreatedPairAndRefundPirCreationFeeToOwner(ctx, liquidityKeeper, assetKeeper, bankKeeper)
		return vm, err
	}
}
