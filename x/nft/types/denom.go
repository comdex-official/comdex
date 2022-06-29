package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func NewDenom(id, symbol, name, schema string, creator sdk.AccAddress, description, previewURI string) Denom {
	return Denom{
		Id:          id,
		Symbol:      symbol,
		Name:        name,
		Schema:      schema,
		Creator:     creator.String(),
		Description: description,
		PreviewURI:  previewURI,
	}
}

func ValidateDenomID(denomID string) error {
	denomID = strings.TrimSpace(denomID)
	if len(denomID) < MinIDLen || len(denomID) > MaxIDLen {
		return sdkerrors.Wrapf(ErrInvalidDenom, "invalid denom ID %s, length  must be between [%d, %d]", denomID, MinIDLen, MaxIDLen)
	}
	if !IsBeginWithAlpha(denomID) || !IsAlphaNumeric(denomID) {
		return sdkerrors.Wrapf(ErrInvalidDenom, "invalid denom ID %s, only accepts alphanumeric characters,and begin with an english letter", denomID)
	}
	return nil
}
func ValidateDenomSymbol(denomSymbol string) error {
	denomSymbol = strings.TrimSpace(denomSymbol)
	if len(denomSymbol) < MinDenomLen || len(denomSymbol) > MaxDenomLen {
		return sdkerrors.Wrapf(ErrInvalidDenom, "invalid denom symbol %s, only accepts value [%d, %d]", denomSymbol, MinDenomLen, MaxDenomLen)
	}
	if !IsBeginWithAlpha(denomSymbol) || !IsAlpha(denomSymbol) {
		return sdkerrors.Wrapf(ErrInvalidDenom, "invalid denom symbol %s, only accepts alphabetic characters", denomSymbol)
	}
	return nil
}

func ValidateName(name string) error {
	name = strings.TrimSpace(name)
	if len(name) > MaxNameLen {
		return sdkerrors.Wrapf(ErrInvalidName, "invalid name %s, length must be less than %d", name, MaxNameLen)
	}
	return nil
}

func ValidateDescription(description string) error {
	description = strings.TrimSpace(description)
	if len(description) > MaxDescriptionLen {
		return sdkerrors.Wrapf(ErrInvalidDescription, "invalid description %s, length must be less than %d", description, MaxDescriptionLen)
	}
	return nil
}

func ValidateURI(uri string) error {
	uri = strings.TrimSpace(uri)
	if len(uri) > MaxURILen {
		return sdkerrors.Wrapf(ErrInvalidURI, "invalid uri %s, length must be less than %d", uri, MaxURILen)
	}
	return nil
}