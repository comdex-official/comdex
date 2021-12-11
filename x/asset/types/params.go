package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var (
	AccountAddressPrefix       = "comdex"
	AccountPubKeyPrefix        = AccountAddressPrefix + sdk.PrefixPublic
	ValidatorAddressPrefix     = AccountAddressPrefix + sdk.PrefixValidator + sdk.PrefixOperator
	ValidatorPubKeyPrefix      = ValidatorAddressPrefix + sdk.PrefixPublic
	ConsensusNodeAddressPrefix = AccountAddressPrefix + sdk.PrefixValidator + sdk.PrefixConsensus
	ConsensusNodePubKeyPrefix  = ConsensusNodeAddressPrefix + sdk.PrefixPublic
)

// SetAccountAddressPrefixes sets the global prefix to be used when serializing addresses to bech32 strings.
func SetAccountAddressPrefixes() {

	config := sdk.GetConfig()

	config.SetBech32PrefixForAccount(AccountAddressPrefix, AccountPubKeyPrefix)
	config.SetBech32PrefixForValidator(ValidatorAddressPrefix, ValidatorPubKeyPrefix)
	config.SetBech32PrefixForConsensusNode(ConsensusNodeAddressPrefix, ConsensusNodePubKeyPrefix)

	config.Seal()
}

const (
	DefaultAdmin = "comdex1le3hcr9mqhpeutr83tqu52taztz4len8ydvlqu"
)

var (
	KeyAdmin = []byte("Admin")
)

var (
	_ paramstypes.ParamSet = (*Params)(nil)
)

func NewParams(admin string) Params {
	return Params{
		Admin: admin,
	}
}

func DefaultParams() Params {
	return NewParams(
		DefaultAdmin,
	)
}

func ParamKeyTable() paramstypes.KeyTable {
	return paramstypes.NewKeyTable().RegisterParamSet(&Params{})
}

func (m *Params) ParamSetPairs() paramstypes.ParamSetPairs {
	return paramstypes.ParamSetPairs{
		paramstypes.NewParamSetPair(
			KeyAdmin,
			m.Admin,
			validateAdmin,
		),
	}
}

func validateAdmin(v interface{}) error {
	value, ok := v.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type %T", v)
	}

	if value == "" {
		return fmt.Errorf("admin cannot be empty")
	}
	if _, err := sdk.AccAddressFromBech32(value); err != nil {
		return errors.Wrapf(err, "invalid admin %s", value)
	}

	return nil
}

func (m *Params) Validate() error {
	if m.Admin == "" {
		return fmt.Errorf("admin cannot be empty")
	}
	if _, err := sdk.AccAddressFromBech32(m.Admin); err != nil {
		return errors.Wrapf(err, "invalid admin %s", m.Admin)
	}

	return nil
}
