package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	ProposalTypeUpdateAdmin = "UpdateOracleAdmin"
)

func init() {
	govtypes.RegisterProposalType(ProposalTypeUpdateAdmin)
	govtypes.RegisterProposalTypeCodec(&UpdateAdminProposal{}, fmt.Sprintf("comdex/market/%s", ProposalTypeUpdateAdmin))
}

var (
	_ govtypes.Content = (*UpdateAdminProposal)(nil)
)

func (m *UpdateAdminProposal) GetTitle() string       { return m.Title }
func (m *UpdateAdminProposal) GetDescription() string { return m.Description }
func (m *UpdateAdminProposal) ProposalRoute() string  { return RouterKey }
func (m *UpdateAdminProposal) ProposalType() string   { return ProposalTypeUpdateAdmin }

func (m *UpdateAdminProposal) ValidateBasic() error {
	if err := govtypes.ValidateAbstract(m); err != nil {
		return err
	}
	if m.Address == "" {
		return fmt.Errorf("address cannot be empty")
	}
	if _, err := sdk.AccAddressFromBech32(m.Address); err != nil {
		return errors.Wrapf(err, "invalid address %s", m.Address)
	}

	return nil
}
