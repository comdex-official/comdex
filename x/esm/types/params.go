package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var (
	_ paramstypes.ParamSet = (*Params)(nil)

	KeyAdmin     = []byte("admin")
	DefaultAdmin = []string{"petri1tadhnvwa0sqzwr3m60f7dsjw4ua77qsz3ptcyw"}
)

const (
	DepositESMGas              = sdk.Gas(36329)
	ExecuteESMGas              = sdk.Gas(23554)
	MsgKillSwitchGas           = sdk.Gas(76473)
	MsgCollateralRedemptionGas = sdk.Gas(37559)
)

func NewParams(admin []string) Params {
	return Params{
		Admin: admin,
	}
}

func DefaultParams() Params {
	return NewParams(DefaultAdmin)
}

func ParamKeyTable() paramstypes.KeyTable {
	return paramstypes.NewKeyTable().RegisterParamSet(&Params{})
}

func (k Params) ParamSetPairs() paramstypes.ParamSetPairs {
	return paramstypes.ParamSetPairs{
		paramstypes.NewParamSetPair(
			KeyAdmin,
			k.Admin,
			validateAdmin,
		),
	}
}

func (k Params) Validate() error {
	if len(k.Admin) == 0 {
		return fmt.Errorf("admin cannot be empty")
	}
	for _, addr := range k.Admin {
		if _, err := sdk.AccAddressFromBech32(addr); err != nil {
			return errors.Wrapf(err, "invalid admin %s", addr)
		}
	}

	return nil
}

func validateAdmin(v interface{}) error {
	if v == "" {
		return fmt.Errorf("admin cannot be empty")
	}

	return nil
}
