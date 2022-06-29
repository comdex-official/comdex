package types


import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func NewGenesisState(collections []Collection) *GenesisState {
	return &GenesisState{
		Collections: collections,
	}
}

func ValidateGenesis(data GenesisState) error {
	for _, c := range data.Collections {
		creator, err := sdk.AccAddressFromBech32(c.Denom.Creator)
		if err != nil {
			return err
		}
		if creator.Empty() {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing denom creator")
		}
		if err := ValidateDenomID(c.Denom.Id); err != nil {
			return err
		}
		if err := ValidateDenomSymbol(c.Denom.Symbol); err != nil {
			return err
		}
		if err := ValidateName(c.Denom.Name); err != nil {
			return err
		}
		if err := ValidateDescription(c.Denom.Description); err != nil {
			return err
		}
		if err := ValidateURI(c.Denom.PreviewURI); err != nil {
			return err
		}

		for _, nft := range c.NFTs {
			if nft.GetOwner().Empty() {
				return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing nft owner")
			}

			if err := ValidateNFTID(nft.GetID()); err != nil {
				return err
			}
			if err := ValidateName(nft.GetName()); err != nil {
				return err
			}
			if err := ValidateDescription(nft.GetDescription()); err != nil {
				return err
			}

			if err := ValidateURI(nft.GetMediaURI()); err != nil {
				return err
			}
			if err := ValidateURI(nft.GetPreviewURI()); err != nil {
				return err
			}
		}
	}
	return nil
}
