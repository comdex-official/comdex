package expected

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	rewardstypes "github.com/comdex-official/comdex/x/rewards/types"
	vaulttypes "github.com/comdex-official/comdex/x/vault/types"
)

type RewardsKeeper interface {
	GetAppIDByApp(ctx sdk.Context, appID uint64) (uint64, bool)
	CalculationOfRewards(ctx sdk.Context, amount sdk.Int, lsr sdk.Dec, bTime int64) (sdk.Dec, error)
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
