package types

import (
	"fmt"

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
	if m.AssetIn == 0 {
		return fmt.Errorf("asset_in cannot be zero")
	}
	if m.AssetOut == 0 {
		return fmt.Errorf("asset_out cannot be zero")
	}
	if m.LiquidationRatio.IsNil() {
		return fmt.Errorf("liquidation_ratio cannot be nil")
	}
	if m.LiquidationRatio.IsNegative() {
		return fmt.Errorf("liquidation_ratio cannot be negative")
	}

	return nil
}
