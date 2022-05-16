package types

import (
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	ProposalLookupTableParams = "LookupTableParams"
)

func init() {
	govtypes.RegisterProposalType(ProposalLookupTableParams)
	govtypes.RegisterProposalTypeCodec(&LookupTableParams{}, "comdex/LookupTableParams")
}

var (
	_ govtypes.Content = &LookupTableParams{}
)

func NewLookupTableParamsProposal(title, description string, lookupTableData []AccumulatorLookupTable) govtypes.Content {
	return &LookupTableParams{
		Title:           title,
		Description:     description,
		LookupTableData: lookupTableData,
	}
}

func (p *LookupTableParams) ProposalRoute() string { return RouterKey }

func (p *LookupTableParams) ProposalType() string { return ProposalLookupTableParams }

func (p *LookupTableParams) ValidateBasic() error {

	err := govtypes.ValidateAbstract(p)
	if err != nil {
		return err
	}
	return nil
}
