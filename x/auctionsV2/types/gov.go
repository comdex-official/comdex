package types

import (
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

const (
	ProposalDutchAutoBidParams = "DutchAutoBidParams"
)

func init() {
	govtypes.RegisterProposalType(ProposalDutchAutoBidParams)
}

var (
	_ govtypes.Content = &DutchAutoBidParamsProposal{}
)

func NewDutchAutoBidParamsProposal(title, description string, dutchAutoBidParams AuctionParams) govtypes.Content {
	return &DutchAutoBidParamsProposal{
		Title:         title,
		Description:   description,
		AuctionParams: dutchAutoBidParams,
	}
}

func (m *DutchAutoBidParamsProposal) ProposalRoute() string {
	return RouterKey
}

func (m *DutchAutoBidParamsProposal) ProposalType() string {
	return ProposalDutchAutoBidParams
}

func (m *DutchAutoBidParamsProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(m)
	if err != nil {
		return err
	}

	if err := m.AuctionParams.Validate(); err != nil {
		return err
	}

	return nil
}

func (m *AuctionParams) Validate() error {
	//TODO: add conditions
	return nil
}
