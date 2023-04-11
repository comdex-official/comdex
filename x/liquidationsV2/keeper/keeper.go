package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/comdex-official/comdex/x/liquidationsV2/expected"
	"github.com/comdex-official/comdex/x/liquidationsV2/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

type Keeper struct {
	cdc        codec.BinaryCodec
	storeKey   sdk.StoreKey
	memKey     sdk.StoreKey
	paramstore paramtypes.Subspace
	account    expected.AccountKeeper
	bank       expected.BankKeeper
	vault      expected.VaultKeeper
	asset      expected.AssetKeeper
	market     expected.MarketKeeper
	esm        expected.EsmKeeper
	rewards    expected.RewardsKeeper
	lend       expected.LendKeeper
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey sdk.StoreKey,
	ps paramtypes.Subspace,
	account expected.AccountKeeper,
	bank expected.BankKeeper,
	asset expected.AssetKeeper,
	vault expected.VaultKeeper,
	market expected.MarketKeeper,
	esm expected.EsmKeeper,
	rewards expected.RewardsKeeper,
	lend expected.LendKeeper,
) Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		cdc:        cdc,
		storeKey:   storeKey,
		memKey:     memKey,
		paramstore: ps,
		account:    account,
		bank:       bank,
		asset:      asset,
		vault:      vault,
		market:     market,
		esm:        esm,
		rewards:    rewards,
		lend:       lend,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) Store(ctx sdk.Context) sdk.KVStore {
	return ctx.KVStore(k.storeKey)
}

// List of functions to be created

//1.  ABCI
//2.		  			Liquidate
//3.					Harbor
//4.					Commodo

//   MsgServer
//5.		HarborLiquidateKeeper
// 6.     	CommodoLiquidateKeeper
//7.		ExternalLiquidateKeeper

//8. Liquidation Whitelisting Proposal

// List of auction functions to be created

// 1. Auction Activator
//		2.Dutch Activator
//		3.English Activator

//4. User Biddings

//5. 					DEPOSIT BID
//6. 					CANCEL BID/ UPDATE BID
//7.					TRIGGER BID  (ABCI)
//8. 					WITHDRAW BID

//9. RESTART AUCTION
//16. END AUCTION
//10. HARBOR AUCTION TRIGGER
//11. HARBOR AUCTION END TRIGGER
//12. VAULT AUCTION END TRIGGER
//13. BORROW AUCTION END TRIGGER
//14. EXTERNAL APPS AUCTION END TRIGGER
//15. INTERNAL LIQUIDATORS INCENTIVISING LOGIC TRIGGER
