package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	ProposalTypeUpdateAdmin = "UpdateAdmin"
)

func init() {
	govtypes.RegisterProposalType(ProposalTypeUpdateAdmin)
	govtypes.RegisterProposalTypeCodec(&UpdateLiquidationRatioProposal{}, fmt.Sprintf("comdex/asset/%s", ProposalTypeUpdateAdmin))
}

var (
	_ govtypes.Content = (*UpdateLiquidationRatioProposal)(nil)
)

func (m *UpdateLiquidationRatioProposal) GetTitle() string       { return m.Title }
func (m *UpdateLiquidationRatioProposal) GetDescription() string { return m.Description }
func (m *UpdateLiquidationRatioProposal) ProposalRoute() string  { return RouterKey }
func (m *UpdateLiquidationRatioProposal) ProposalType() string   { return ProposalTypeUpdateAdmin }

func (m *UpdateLiquidationRatioProposal) ValidateBasic() error {

	return nil
}

func (m *UpdateLiquidationRatioProposal) GetSigners() []sdk.AccAddress {
	return nil
}
