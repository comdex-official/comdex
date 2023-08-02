package types

import (
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

const (
	ProposalWhitelistLiquidation = "AddWhitelistLiquidation"
)

func init() {
	govtypes.RegisterProposalType(ProposalWhitelistLiquidation)
}

var (
	_ govtypes.Content = &WhitelistLiquidationProposal{}
)

func NewLiquidationWhiteListingProposal(title, description string, liquidationWhiteListing LiquidationWhiteListing) govtypes.Content {
	return &WhitelistLiquidationProposal{
		Title:        title,
		Description:  description,
		Whitelisting: liquidationWhiteListing,
	}
}

func (p *WhitelistLiquidationProposal) GetTitle() string {
	return p.Title
}

func (p *WhitelistLiquidationProposal) GetDescription() string {
	return p.Description
}
func (p *WhitelistLiquidationProposal) ProposalRoute() string { return RouterKey }

func (p *WhitelistLiquidationProposal) ProposalType() string { return ProposalWhitelistLiquidation }

func (p *WhitelistLiquidationProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(p)
	if err != nil {
		return err
	}

	if err := p.Whitelisting.Validate(); err != nil {
		return err
	}

	return nil
}

func (m *LiquidationWhiteListing) Validate() error {
	//TODO: add conditions
	return nil
}
