package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func NewParams(
	securityAddress []string, contractGasLimit uint64,
) Params {
	return Params{
		SecurityAddress:  securityAddress,
		ContractGasLimit: contractGasLimit,
	}
}

// default minting module parameters
func DefaultParams() Params {
	return Params{
		SecurityAddress:  []string{"comdex1nh4gxgzq7hw8fvtkxjg4kpfqmsq65szqxxdqye"},
		ContractGasLimit: uint64(1000000000),
	}
}

// validate params
func (p Params) Validate() error {
	minimumGas := uint64(100_000)
	if p.ContractGasLimit < minimumGas {
		return errorsmod.Wrapf(
			sdkerrors.ErrInvalidRequest,
			"invalid contract gas limit: %d. Must be above %d", p.ContractGasLimit, minimumGas,
		)
	}

	for _, addr := range p.SecurityAddress {
		// Valid address check
		if _, err := sdk.AccAddressFromBech32(addr); err != nil {
			return errorsmod.Wrapf(
				sdkerrors.ErrInvalidAddress,
				"invalid security address: %s", err.Error(),
			)
		}

		// duplicate address check
		count := 0
		for _, addr2 := range p.SecurityAddress {
			if addr == addr2 {
				count++
			}

			if count > 1 {
				return errorsmod.Wrapf(
					sdkerrors.ErrInvalidAddress,
					"duplicate contract address: %s", addr,
				)
			}
		}
	}

	return nil
}
