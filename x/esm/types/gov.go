package types

import govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

var (
	ProposalAddESMTriggerParams = "ProposalAddESMTriggerParams"
)

func init() {
	govtypes.RegisterProposalType(ProposalAddESMTriggerParams)
	govtypes.RegisterProposalTypeCodec(&ESMTriggerParams{}, "comdex/AddESMTriggerParams")
}

var (
	_ govtypes.Content = &ESMTriggerParamsProposal{}
)

func NewAddESMTriggerParamsProposal(title, description string, esmTriggerParams ESMTriggerParams) govtypes.Content {
	return &ESMTriggerParamsProposal{
		Title:            title,
		Description:      description,
		EsmTriggerParams: esmTriggerParams,
	}
}

func (p *ESMTriggerParamsProposal) ProposalRoute() string { return RouterKey }

func (p *ESMTriggerParamsProposal) ProposalType() string { return ProposalAddESMTriggerParams }

func (p *ESMTriggerParamsProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(p)
	if err != nil {
		return err
	}

	return nil
}
