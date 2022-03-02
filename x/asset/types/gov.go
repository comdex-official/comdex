package types

import (
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	ProposalTypeUpdateLiquidationRatio = "UpdateLiquidationRatio"
)

func init() {
	govtypes.RegisterProposalType(ProposalTypeUpdateLiquidationRatio)
	govtypes.RegisterProposalTypeCodec(&UpdateLiquidationRatio{}, "comdex/UpdateLiquidationRatioProposal")
}

var (
	_ govtypes.Content = (*UpdateLiquidationRatio)(nil)
)

func (m *UpdateLiquidationRatio) GetTitle() string       { return m.Title }
func (m *UpdateLiquidationRatio) GetDescription() string { return m.Description }
func (m *UpdateLiquidationRatio) ProposalRoute() string  { return RouterKey }
func (m *UpdateLiquidationRatio) ProposalType() string {
	return ProposalTypeUpdateLiquidationRatio
}

func (m *UpdateLiquidationRatio) ValidateBasic() error {

	return nil
}

/*func (m *UpdateLiquidationRatioProposal) GetSigners() []sdk.AccAddress {

	from, err := sdk.AccAddressFromBech32(m.From)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}*/

func NewUpdateLiquidationRatio(title, description, liquidationRatio string) *govtypes.Content {
	return &UpdateLiquidationRatio{
		Title:            title,
		Description:      description,
		LiquidationRatio: liquidationRatio,
	}
}
