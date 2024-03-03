package v14

import (
	commonkeeper "github.com/comdex-official/comdex/x/common/keeper"
	commontypes "github.com/comdex-official/comdex/x/common/types"
	lendkeeper "github.com/comdex-official/comdex/x/lend/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	auctionkeeperskip "github.com/skip-mev/block-sdk/x/auction/keeper"
	auctionmoduleskiptypes "github.com/skip-mev/block-sdk/x/auction/types"
	"strings"
)

func CreateUpgradeHandlerV14(
	mm *module.Manager,
	configurator module.Configurator,
	commonkeeper commonkeeper.Keeper,
	auctionkeeperskip auctionkeeperskip.Keeper,
	lendKeeper lendkeeper.Keeper,

) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {

		ctx.Logger().Info("Applying test net upgrade - v14.0.0")

		vm, err := mm.RunMigrations(ctx, configurator, fromVM)
		if err != nil {
			return vm, err
		}
		ctx.Logger().Info("set common module params")
		commonkeeper.SetParams(ctx, commontypes.DefaultParams())

		ctx.Logger().Info("setting default params for MEV module (x/auction)")
		if err = setDefaultMEVParams(ctx, auctionkeeperskip); err != nil {
			return nil, err
		}

		//TODO: uncomment this before mainnet upgrade
		//UpdateLendParams(ctx, lendKeeper)
		return vm, err
	}
}

func setDefaultMEVParams(ctx sdk.Context, auctionkeeperskip auctionkeeperskip.Keeper) error {
	nativeDenom := getChainBondDenom(ctx.ChainID())

	// Skip MEV (x/auction)
	return auctionkeeperskip.SetParams(ctx, auctionmoduleskiptypes.Params{
		MaxBundleSize:          auctionmoduleskiptypes.DefaultMaxBundleSize,
		EscrowAccountAddress:   authtypes.NewModuleAddress(auctionmoduleskiptypes.ModuleName), // TODO: revisit
		ReserveFee:             sdk.NewCoin(nativeDenom, sdk.NewInt(10)),
		MinBidIncrement:        sdk.NewCoin(nativeDenom, sdk.NewInt(5)),
		FrontRunningProtection: auctionmoduleskiptypes.DefaultFrontRunningProtection,
		ProposerFee:            auctionmoduleskiptypes.DefaultProposerFee,
	})
}

// getChainBondDenom returns expected bond denom based on chainID.
func getChainBondDenom(chainID string) string {
	if strings.HasPrefix(chainID, "comdex-") {
		return "ucmdx"
	}
	return "stake"
}

func UpdateLendParams(
	ctx sdk.Context,
	lendKeeper lendkeeper.Keeper,
) {
	assetRatesParamsStAtom, _ := lendKeeper.GetAssetRatesParams(ctx, 14)
	assetRatesParamsStAtom.CAssetID = 23
	lendKeeper.SetAssetRatesParams(ctx, assetRatesParamsStAtom)
}
