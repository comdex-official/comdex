package expected

import (
	"context"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	rewardstypes "github.com/comdex-official/comdex/x/rewards/types"
	vaulttypes "github.com/comdex-official/comdex/x/vault/types"
)

type AccountKeeper interface {
	GetAccount(ctx context.Context, addr sdk.AccAddress) sdk.AccountI 
	GetModuleAddress(moduleName string) sdk.AccAddress
}

// BankKeeper is the expected bank keeper.
type BankKeeper interface {
	GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin 
	GetAllBalances(ctx context.Context, addr sdk.AccAddress) sdk.Coins
	SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
}

type RewardsKeeper interface {
	GetAppIDByApp(ctx sdk.Context, appID uint64) (uint64, bool)
	CalculationOfRewards(ctx sdk.Context, amount sdkmath.Int, lsr sdkmath.LegacyDec, bTime int64) (sdkmath.LegacyDec, error)
	GetVaultInterestTracker(ctx sdk.Context, id, appID uint64) (vault rewardstypes.VaultInterestTracker, found bool)
	SetVaultInterestTracker(ctx sdk.Context, vault rewardstypes.VaultInterestTracker)
}

type VaultKeeper interface {
	GetAppExtendedPairVaultMappingData(ctx sdk.Context, appMappingID uint64, pairVaultsID uint64) (appExtendedPairVaultData vaulttypes.AppExtendedPairVaultMappingData, found bool)
	GetVault(ctx sdk.Context, id uint64) (vault vaulttypes.Vault, found bool)
	SetVault(ctx sdk.Context, vault vaulttypes.Vault)
}

type Bandoraclekeeper interface {
	SetCheckFlag(ctx sdk.Context, flag bool)
}
