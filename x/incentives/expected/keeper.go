package expected

import (
	"github.com/comdex-official/comdex/x/incentives/types"
	liquiditytypes "github.com/comdex-official/comdex/x/liquidity/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
	// Methods imported from account should be defined here
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
}

type LiquidityKeeper interface {
	GetPool(ctx sdk.Context, id uint64) (pool liquiditytypes.Pool, found bool)
	GetFarmingRewardsData(ctx sdk.Context, liquidityGaugeData types.LiquidtyGaugeMetaData) []types.RewardDistributionDataCollector
}
