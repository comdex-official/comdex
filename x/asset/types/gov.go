package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	ProposalTypeUpdateLiquidationRatio = "UpdateLiquidationRatioProposal"
)

func init() {
	govtypes.RegisterProposalType(ProposalTypeUpdateLiquidationRatio)
	govtypes.RegisterProposalTypeCodec(&UpdateLiquidationRatioProposal{}, "/comdex/UpdateLiquidationRatioProposal")
}

var (
	_ govtypes.Content = (*UpdateLiquidationRatioProposal)(nil)
)

func (m *UpdateLiquidationRatioProposal) GetTitle() string       { return m.Title }
func (m *UpdateLiquidationRatioProposal) GetDescription() string { return m.Description }
func (m *UpdateLiquidationRatioProposal) ProposalRoute() string  { return RouterKey }
func (m *UpdateLiquidationRatioProposal) ProposalType() string {
	return ProposalTypeUpdateLiquidationRatio
}

func (m *UpdateLiquidationRatioProposal) ValidateBasic() error {

	return nil
}

func (m *UpdateLiquidationRatioProposal) GetSigners() []sdk.AccAddress {

	from, err := sdk.AccAddressFromBech32(m.From)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}

func NewUpdateLiquidationRatio(title, description, from, liquidationRatio string) *govtypes.Content {
	return &UpdateLiquidationRatioProposal{
		Title:            title,
		Description:      description,
		LiquidationRatio: liquidationRatio,
		From:             from,
	}
}
