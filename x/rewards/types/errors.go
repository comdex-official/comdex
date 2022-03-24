package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrorUnknownProposalType           = sdkerrors.Register(ModuleName, 101, "Proposal type is invalid")
	ErrorInvalidAssetDenoms            = sdkerrors.Register(ModuleName, 102, "Invalid asset denoms")
	ErrorMintingRewardPairAlreadyExist = sdkerrors.Register(ModuleName, 103, "rewards for given pair  already exists")
	ErrorInvalidStartTime              = sdkerrors.Register(ModuleName, 104, "starttime should be atleast 10 minutes after current time")
	ErrorMintingRewardNotFound         = sdkerrors.Register(ModuleName, 105, "minting rewards id invalid/not found")
	ErrorMintingRewardAlreadyActive    = sdkerrors.Register(ModuleName, 106, "minting rewards already active")
	ErrorMintingRewardAlreadyDisabled  = sdkerrors.Register(ModuleName, 107, "minting rewards already disabled")
	ErrorMintingRewardExpired          = sdkerrors.Register(ModuleName, 108, "minting rewards expired")
	ErrorDepositAlreadyMade            = sdkerrors.Register(ModuleName, 109, "deposit already made")
	ErrorUnauthorized                  = sdkerrors.Register(ModuleName, 110, "unauthorized to update start time")
	ErrorPriceNotFound                 = sdkerrors.Register(ModuleName, 111, "asset price not found")
)

var (
	ErrorUnknownMsgType = sdkerrors.Register(ModuleName, 301, "unknown message type")
)
