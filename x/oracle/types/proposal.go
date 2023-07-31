package types

import (
	gov "github.com/cosmos/cosmos-sdk/x/gov/types"
)

var (
	proposalTypeMsgGovUpdateParams = MsgGovUpdateParams{}.Type()
)

func init() {
	gov.RegisterProposalType(proposalTypeMsgGovUpdateParams)
}

// Implements Proposal Interface
var _ gov.Content = &MsgGovUpdateParams{}

// GetTitle returns the title of a community pool spend proposal.
func (msg *MsgGovUpdateParams) GetTitle() string { return msg.Title }

// GetDescription returns the description of a community pool spend proposal.
func (msg *MsgGovUpdateParams) GetDescription() string { return msg.Description }

// GetDescription returns the routing key of a community pool spend proposal.
func (msg *MsgGovUpdateParams) ProposalRoute() string { return RouterKey }

// ProposalType returns the type of a community pool spend proposal.
func (msg *MsgGovUpdateParams) ProposalType() string { return proposalTypeMsgGovUpdateParams }
