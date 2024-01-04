package keeper

import (
	"github.com/comdex-official/comdex/x/rwa/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetRwaUser(ctx sdk.Context, rwaUser types.RwaUser) error {
	var (
		store = k.Store(ctx)
		key   = types.RwaUserKey(rwaUser.AccountAddress)
		value = k.cdc.MustMarshal(&rwaUser)
	)

	store.Set(key, value)
	return nil
}

func (k Keeper) GetRwaUSer(ctx sdk.Context, address string) (rwaUser types.RwaUser, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.RwaUserKey(address)
		value = store.Get(key)
	)

	if value == nil {
		return rwaUser, false
	}

	k.cdc.MustUnmarshal(value, &rwaUser)
	return rwaUser, true
}

func (k Keeper) GetCounterparty(ctx sdk.Context, ID uint64) (counterParty types.Counterparty, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.CounterPartyKey(ID)
		value = store.Get(key)
	)

	if value == nil {
		return counterParty, false
	}

	k.cdc.MustUnmarshal(value, &counterParty)
	return counterParty, true
}

func (k Keeper) SetCounterParty(ctx sdk.Context, counterParty types.Counterparty) error {
	var (
		store = k.Store(ctx)
		key   = types.CounterPartyKey(counterParty.Id)
		value = k.cdc.MustMarshal(&counterParty)
	)

	store.Set(key, value)
	return nil
}

func (k Keeper) UpdateCounterParty(ctx sdk.Context, counterParty types.Counterparty) error {
	var (
		store = k.Store(ctx)
		key   = types.CounterPartyKey(counterParty.Id)
		value = k.cdc.MustMarshal(&counterParty)
	)

	store.Set(key, value)
	return nil
}

func (k Keeper) UpdateRwaUser(ctx sdk.Context, rwaUser types.RwaUser) error {
	var (
		store = k.Store(ctx)
		key   = types.RwaUserKey(rwaUser.AccountAddress)
		value = k.cdc.MustMarshal(&rwaUser)
	)

	store.Set(key, value)
	return nil
}

func (k Keeper) UpdateRwaUserCounterPartyList(ctx sdk.Context, rwaUser types.RwaUser) error {
	var (
		store = k.Store(ctx)
		key   = types.RwaUserKey(rwaUser.AccountAddress)
		value = k.cdc.MustMarshal(&rwaUser)
	)

	store.Set(key, value)
	return nil
}

func (k Keeper) SetInvoice(ctx sdk.Context, invoice types.Invoice) error {
	var (
		store = k.Store(ctx)
		key   = types.InvoiceKey(invoice.Id)
		value = k.cdc.MustMarshal(&invoice)
	)

	store.Set(key, value)
	return nil
}

func (k Keeper) GetInvoice(ctx sdk.Context, ID uint64) (invoice types.Invoice, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.InvoiceKey(ID)
		value = store.Get(key)
	)

	if value == nil {
		return invoice, false
	}

	k.cdc.MustUnmarshal(value, &invoice)
	return invoice, true
}
