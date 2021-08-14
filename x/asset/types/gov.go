package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	ProposalTypeAddPair = "AddPair"
)

func init() {
	govtypes.RegisterProposalType(ProposalTypeAddPair)
	govtypes.RegisterProposalTypeCodec(&AddPairProposal{}, fmt.Sprintf("comdex/%s", ProposalTypeAddPair))
}

var (
	_ govtypes.Content = (*AddPairProposal)(nil)
)

func (m *AddPairProposal) GetTitle() string       { return m.Title }
func (m *AddPairProposal) GetDescription() string { return m.Description }
func (m *AddPairProposal) ProposalRoute() string  { return RouterKey }
func (m *AddPairProposal) ProposalType() string   { return ProposalTypeAddPair }

func (m *AddPairProposal) ValidateBasic() error {
	if err := govtypes.ValidateAbstract(m); err != nil {
		return err
	}
	if m.DenomIn == "" {
		return fmt.Errorf("denom_in cannot be empty")
	}
	if err := sdk.ValidateDenom(m.DenomIn); err != nil {
		return errors.Wrapf(err, "invalid denom_in %s", m.DenomIn)
	}
	if m.DenomOut == "" {
		return fmt.Errorf("denom_out cannot be empty")
	}
	if err := sdk.ValidateDenom(m.DenomOut); err != nil {
		return errors.Wrapf(err, "invalid denom_out %s", m.DenomOut)
	}
	if m.LiquidationRatio.IsNil() {
		return fmt.Errorf("liquidation_ratio cannot be nil")
	}
	if m.LiquidationRatio.IsNegative() {
		return fmt.Errorf("liquidation_ratio cannot be negative")
	}

	return nil
}
