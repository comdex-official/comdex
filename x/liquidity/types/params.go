package types

import (
	"time"

	storetypes "cosmossdk.io/store/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

const (
	DefaultFarmingQueueDuration = time.Hour * 24

	FeeCollectorAddressPrefix = "FeeCollectorAddress"

	PoolReserveAddressPrefix          = "PoolReserveAddress"
	PairSwapFeeCollectorAddressPrefix = "PairSwapFeeCollectorAddress"
	PairEscrowAddressPrefix           = "PairEscrowAddress"
	ModuleAddressNameSplitter         = "|"
)

const (
	CreatePoolGas      = storetypes.Gas(67500)
	CancelOrderGas     = storetypes.Gas(65000)
	CancelAllOrdersGas = storetypes.Gas(74000)
	FarmGas            = storetypes.Gas(62300)
	UnfarmGas          = storetypes.Gas(69000)
)

// GlobalEscrowAddress is an escrow for deposit/withdraw requests.
var GlobalEscrowAddress = DeriveAddress(AddressType32Bytes, ModuleName, "GlobalEscrow")

var _ paramstypes.ParamSet = (*Params)(nil)

func ParamKeyTable() paramstypes.KeyTable {
	return paramstypes.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams returns a default params for the liquidity module.
func DefaultParams() Params {
	return Params{}
}

// ParamSetPairs implements ParamSet.
func (params *Params) ParamSetPairs() paramstypes.ParamSetPairs {
	return paramstypes.ParamSetPairs{}
}

// Validate validates Params.
func (params Params) Validate() error {
	for _, field := range []struct {
		val          interface{}
		validateFunc func(i interface{}) error
	}{} {
		if err := field.validateFunc(field.val); err != nil {
			return err
		}
	}
	return nil
}
