package types

import (
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	ProposalToggleEsm             = "ToggleEsm"
)

func init() {
	govtypes.RegisterProposalType(ProposalToggleEsm)
	govtypes.RegisterProposalTypeCodec(&ToggleEsmProposal{}, "comdex/ToggleEsm")

}

var (
	_ govtypes.Content = &ToggleEsmProposal{}
)

func NewToggleEsmProposal(title, description string, esmActive EsmActive ) govtypes.Content {
	return &ToggleEsmProposal{
		Title:       title,
		Description: description,
		EsmActive:      esmActive,
	}
}
func (p *ToggleEsmProposal) GetTitle() string {
	return p.Title
}

func (p *ToggleEsmProposal) GetDescription() string {
	return p.Description
}
func (p *ToggleEsmProposal) ProposalRoute() string { return RouterKey }

func (p *ToggleEsmProposal) ProposalType() string { return ProposalToggleEsm }

func (p *ToggleEsmProposal) ValidateBasic() error {

	err := govtypes.ValidateAbstract(p)
	if err != nil {
		return err
	}

	return nil
}