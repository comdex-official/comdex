package v14

import (
	"fmt"
	commonkeeper "github.com/comdex-official/comdex/x/common/keeper"
	commontypes "github.com/comdex-official/comdex/x/common/types"
	lendkeeper "github.com/comdex-official/comdex/x/lend/keeper"
	tokenfactorykeeper "github.com/comdex-official/comdex/x/tokenfactory/keeper"
	tokenfactorytypes "github.com/comdex-official/comdex/x/tokenfactory/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	auctionkeeperskip "github.com/skip-mev/block-sdk/x/auction/keeper"
	auctionmoduleskiptypes "github.com/skip-mev/block-sdk/x/auction/types"
	"strings"
)

// We now charge 2 million gas * gas price to create a denom.
const NewDenomCreationGasConsume uint64 = 2_000_000

func CreateUpgradeHandlerV14(
	mm *module.Manager,
	configurator module.Configurator,
	commonkeeper commonkeeper.Keeper,
	auctionkeeperskip auctionkeeperskip.Keeper,
	lendKeeper lendkeeper.Keeper,
	tokenfactorykeeper tokenfactorykeeper.Keeper,
	accountKeeper authkeeper.AccountKeeper,

) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {

		ctx.Logger().Info("Applying test net upgrade - v14.0.0")

		vm, err := mm.RunMigrations(ctx, configurator, fromVM)
		if err != nil {
			return vm, err
		}
		moduleAccI := accountKeeper.GetModuleAccount(ctx, authtypes.FeeCollectorName)
		moduleAcc := moduleAccI.(*authtypes.ModuleAccount)
		moduleAcc.Permissions = []string{authtypes.Burner}
		accountKeeper.SetModuleAccount(ctx, moduleAcc)

		ctx.Logger().Info("set common module params")
		commonkeeper.SetParams(ctx, commontypes.DefaultParams())

		ctx.Logger().Info("setting default params for MEV module (x/auction)")
		if err = setDefaultMEVParams(ctx, auctionkeeperskip); err != nil {
			return nil, err
		}

		// x/TokenFactory
		// Use denom creation gas consumption instead of fee for contract developers
		ctx.Logger().Info("setting params for Tokenfactory module (x/tokenfactory)")
		updatedTf := tokenfactorytypes.Params{
			DenomCreationFee:        nil,
			DenomCreationGasConsume: NewDenomCreationGasConsume,
		}

		if err := tokenfactorykeeper.SetParams(ctx, updatedTf); err != nil {
			return vm, err
		}
		ctx.Logger().Info(fmt.Sprintf("updated tokenfactory params to %v", updatedTf))

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
