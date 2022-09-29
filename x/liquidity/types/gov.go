package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	ProposalUpdateGenericParams    = "UpdateGenericParams"
	ProposalCreateNewLiquidityPair = "CreateNewLiquidityPair"
)

func init() {
	govtypes.RegisterProposalType(ProposalUpdateGenericParams)
	govtypes.RegisterProposalType(ProposalCreateNewLiquidityPair)
	govtypes.RegisterProposalTypeCodec(&UpdateGenericParamsProposal{}, "comdex/UpdateGenericParams")
	govtypes.RegisterProposalTypeCodec(&CreateNewLiquidityPairProposal{}, "comdex/CreateNewLiquidityPair")
}

var (
	_ govtypes.Content = &UpdateGenericParamsProposal{}
	_ govtypes.Content = &CreateNewLiquidityPairProposal{}
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

func NewCreateLiquidityPairProposal(
	title, description string,
	from sdk.AccAddress,
	appID uint64,
	baseCoinDenom, quoteCoinDenom string,
) govtypes.Content {
	return &CreateNewLiquidityPairProposal{
		Title:          title,
		Description:    description,
		AppId:          appID,
		BaseCoinDenom:  baseCoinDenom,
		QuoteCoinDenom: quoteCoinDenom,
		From:           from.String(),
	}
}

func (p *CreateNewLiquidityPairProposal) GetTitle() string {
	return p.Title
}

func (p *CreateNewLiquidityPairProposal) GetDescription() string {
	return p.Description
}
func (p *CreateNewLiquidityPairProposal) ProposalRoute() string { return RouterKey }

func (p *CreateNewLiquidityPairProposal) ProposalType() string { return ProposalCreateNewLiquidityPair }

func (p *CreateNewLiquidityPairProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(p)
	if err != nil {
		return err
	}

	if p.AppId <= 0 {
		return ErrInvalidAppID
	}

	return nil
}
