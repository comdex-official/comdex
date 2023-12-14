package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/nft module sentinel errors
var (
	ErrInvalidCollection  = sdkerrors.Register(ModuleName, 3, "invalid NFT collection")
	ErrUnknownCollection  = sdkerrors.Register(ModuleName, 4, "unknown NFT collection")
	ErrInvalidNFT         = sdkerrors.Register(ModuleName, 5, "invalid NFT")
	ErrNFTAlreadyExists   = sdkerrors.Register(ModuleName, 6, "NFT already exists")
	ErrUnknownNFT         = sdkerrors.Register(ModuleName, 7, "unknown NFT")
	ErrEmptyMetaData      = sdkerrors.Register(ModuleName, 8, "NFT MetaData can't be empty")
	ErrUnauthorized       = sdkerrors.Register(ModuleName, 9, "unauthorized address")
	ErrInvalidDenom       = sdkerrors.Register(ModuleName, 10, "invalid denom")
	ErrInvalidNFTID       = sdkerrors.Register(ModuleName, 11, "invalid ID")
	ErrInvalidNFTMeta     = sdkerrors.Register(ModuleName, 12, "invalid metadata")
	ErrInvalidMediaURI    = sdkerrors.Register(ModuleName, 13, "invalid media URI")
	ErrInvalidPreviewURI  = sdkerrors.Register(ModuleName, 14, "invalid preview URI")
	ErrNotTransferable    = sdkerrors.Register(ModuleName, 15, "nft is not transferable")
	ErrNotEditable        = sdkerrors.Register(ModuleName, 16, "nft is not editable")
	ErrInvalidOption      = sdkerrors.Register(ModuleName, 17, "invalid option")
	ErrInvalidName        = sdkerrors.Register(ModuleName, 18, "invalid name")
	ErrInvalidDescription = sdkerrors.Register(ModuleName, 19, "invalid description")
	ErrInvalidURI         = sdkerrors.Register(ModuleName, 20, "invalid URI")
	ErrInvalidPercentage  = sdkerrors.Register(ModuleName, 21, "invalid percentage")
)
