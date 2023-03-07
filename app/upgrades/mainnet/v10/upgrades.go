package v10

import (
	auctiontypes "github.com/comdex-official/comdex/x/auction/types"
	esmtypes "github.com/comdex-official/comdex/x/esm/types"
	lendtypes "github.com/comdex-official/comdex/x/lend/types"
	liquidationtypes "github.com/comdex-official/comdex/x/liquidation/types"
	liquiditytypes "github.com/comdex-official/comdex/x/liquidity/types"
	lockertypes "github.com/comdex-official/comdex/x/locker/types"
	rewardstypes "github.com/comdex-official/comdex/x/rewards/types"
	tokenminttypes "github.com/comdex-official/comdex/x/tokenmint/types"
	vaulttypes "github.com/comdex-official/comdex/x/vault/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	ica "github.com/cosmos/ibc-go/v4/modules/apps/27-interchain-accounts"
	icacontrollertypes "github.com/cosmos/ibc-go/v4/modules/apps/27-interchain-accounts/controller/types"
	icahosttypes "github.com/cosmos/ibc-go/v4/modules/apps/27-interchain-accounts/host/types"
	icatypes "github.com/cosmos/ibc-go/v4/modules/apps/27-interchain-accounts/types"
)

func CreateUpgradeHandlerV10(
	mm *module.Manager,
	configurator module.Configurator,
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
		return vm, err
	}
}
