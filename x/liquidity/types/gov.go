package types

import (
	"fmt"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	ProposalUpdateGenericParams = "UpdateGenericParams"
)

func init() {
	govtypes.RegisterProposalType(ProposalUpdateGenericParams)
	govtypes.RegisterProposalTypeCodec(&UpdateGenericParamsProposal{}, "comdex/UpdateGenericParams")
}

var (
	_ govtypes.Content = &UpdateGenericParamsProposal{}
)

func NewUpdateGenericParamsProposal(
	title, description string,
	appID uint64,
	keys, value []string,
) govtypes.Content {
	return &UpdateGenericParamsProposal{
		Title:       title,
		Description: description,
		AppId:       appID,
		Keys:        keys,
		Values:      value,
	}
}
func (p *UpdateGenericParamsProposal) GetTitle() string {
	return p.Title
}

func (p *UpdateGenericParamsProposal) GetDescription() string {
	return p.Description
}
func (p *UpdateGenericParamsProposal) ProposalRoute() string { return RouterKey }

func (p *UpdateGenericParamsProposal) ProposalType() string { return ProposalUpdateGenericParams }

func (p *UpdateGenericParamsProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(p)
	if err != nil {
		return err
	}

	if p.AppId <= 0 {
		return ErrInvalidAppID
	}

	if len(p.Keys) == 0 || len(p.Values) == 0 {
		return ErrorEmptyKeyValueForGenericParams
	}

	if len(p.Keys) != len(p.Values) {
		return ErrorLengthMismatch
	}

	for _, key := range p.Keys {
		keyFound := false
		for _, uKey := range UpdatableKeys {
			if uKey == key {
				keyFound = true
			}
		}
		if !keyFound {
			return fmt.Errorf("invalid key for update: %s", key)
		}
	}

	return nil
}
