package keepers

import (
	"github.com/CosmWasm/wasmd/x/wasm"
	assettypes "github.com/comdex-official/comdex/x/asset/types"
	auctiontypes "github.com/comdex-official/comdex/x/auction/types"
	auctionsV2types "github.com/comdex-official/comdex/x/auctionsV2/types"
	bandoraclemoduletypes "github.com/comdex-official/comdex/x/bandoracle/types"
	collectortypes "github.com/comdex-official/comdex/x/collector/types"
	esmtypes "github.com/comdex-official/comdex/x/esm/types"
	lendtypes "github.com/comdex-official/comdex/x/lend/types"
	liquidationtypes "github.com/comdex-official/comdex/x/liquidation/types"
	liquidationsV2types "github.com/comdex-official/comdex/x/liquidationsV2/types"
	liquiditytypes "github.com/comdex-official/comdex/x/liquidity/types"
	lockertypes "github.com/comdex-official/comdex/x/locker/types"
	markettypes "github.com/comdex-official/comdex/x/market/types"
	rewardstypes "github.com/comdex-official/comdex/x/rewards/types"
	tokenminttypes "github.com/comdex-official/comdex/x/tokenmint/types"
	vaulttypes "github.com/comdex-official/comdex/x/vault/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	consensusparamtypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	nftkeeper "github.com/cosmos/cosmos-sdk/x/nft/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	packetforwardtypes "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v7/packetforward/types"
	icqtypes "github.com/cosmos/ibc-apps/modules/async-icq/v7/types"
	ibchookstypes "github.com/cosmos/ibc-apps/modules/ibc-hooks/v7/types"
	icahosttypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/host/types"
	ibcfeetypes "github.com/cosmos/ibc-go/v7/modules/apps/29-fee/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	ibchost "github.com/cosmos/ibc-go/v7/modules/core/exported"
)

func (appKeepers *AppKeepers) GenerateKeys() {

	appKeepers.tkeys = sdk.NewTransientStoreKeys(paramstypes.TStoreKey)
	appKeepers.mkeys = sdk.NewMemoryStoreKeys(capabilitytypes.MemStoreKey)

	appKeepers.keys = sdk.NewKVStoreKeys(
		authtypes.StoreKey, banktypes.StoreKey, stakingtypes.StoreKey,
		minttypes.StoreKey, distrtypes.StoreKey, slashingtypes.StoreKey,
		govtypes.StoreKey, paramstypes.StoreKey, ibchost.StoreKey, icahosttypes.StoreKey, upgradetypes.StoreKey,
		evidencetypes.StoreKey, ibctransfertypes.StoreKey, ibcfeetypes.StoreKey, capabilitytypes.StoreKey,
		vaulttypes.StoreKey, assettypes.StoreKey, collectortypes.StoreKey, liquidationtypes.StoreKey,
		markettypes.StoreKey, bandoraclemoduletypes.StoreKey, lockertypes.StoreKey,
		wasm.StoreKey, authzkeeper.StoreKey, auctiontypes.StoreKey, tokenminttypes.StoreKey,
		rewardstypes.StoreKey, feegrant.StoreKey, liquiditytypes.StoreKey, esmtypes.ModuleName, lendtypes.StoreKey,
		liquidationsV2types.StoreKey, auctionsV2types.StoreKey, ibchookstypes.StoreKey, packetforwardtypes.StoreKey,
		icqtypes.StoreKey, consensusparamtypes.StoreKey, crisistypes.StoreKey, nftkeeper.StoreKey,
	)
}

// GetKey returns the KVStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (appkeepers *AppKeepers) GetKey(storeKey string) *storetypes.KVStoreKey {
	return appkeepers.keys[storeKey]
}

// GetTKey returns the TransientStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (appkeepers *AppKeepers) GetTKey(storeKey string) *storetypes.TransientStoreKey {
	return appkeepers.tkeys[storeKey]
}

// GetMemKey returns the MemStoreKey for the provided mem key.
//
// NOTE: This is solely used for testing purposes.
func (appkeepers *AppKeepers) GetMemKey(storeKey string) *storetypes.MemoryStoreKey {
	return appkeepers.mkeys[storeKey]
}

func (appKeepers *AppKeepers) GetKVStoreKey() map[string]*storetypes.KVStoreKey {
	return appKeepers.keys
}

func (appKeepers *AppKeepers) GetTransientStoreKey() map[string]*storetypes.TransientStoreKey {
	return appKeepers.tkeys
}

func (appKeepers *AppKeepers) GetMemoryStoreKey() map[string]*storetypes.MemoryStoreKey {
	return appKeepers.mkeys
}
