package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgFetchPriceData = "fetch_price_data"

var (
	_ sdk.Msg = &MsgFetchPriceData{}

	// FetchPriceResultStoreKeyPrefix is a prefix for storing result
	FetchPriceResultStoreKeyPrefix = "fetch_price_result"

	// LastFetchPriceIDKey is the key for the last request id
	LastFetchPriceIDKey = "fetch_price_last_id"

	TempFetchPriceIDKey = "fetch_price_temp_id"

	// FetchPriceClientIDKey is query request identifier
	FetchPriceClientIDKey = "fetch_price_id"

	LastBlockheightKey = "last_blockheight"

	OracleValidationResultKey = "Oracle_Validation_Result"
)

// NewMsgFetchPriceData creates a new FetchPrice message
func NewMsgFetchPriceData(
	creator string,
	oracleScriptID OracleScriptID,
	sourceChannel string,
	calldata *FetchPriceCallData,
	askCount uint64,
	minCount uint64,
	feeLimit sdk.Coins,
	prepareGas uint64,
	executeGas uint64,
) *MsgFetchPriceData {
	return &MsgFetchPriceData{
		ClientID:       FetchPriceClientIDKey,
		Creator:        creator,
		OracleScriptID: uint64(oracleScriptID),
		SourceChannel:  sourceChannel,
		Calldata:       calldata,
		AskCount:       askCount,
		MinCount:       minCount,
		FeeLimit:       feeLimit,
		PrepareGas:     prepareGas,
		ExecuteGas:     executeGas,
	}
}

// Route returns the message route
func (m *MsgFetchPriceData) Route() string {
	return RouterKey
}

// Type returns the message type
func (m *MsgFetchPriceData) Type() string {
	return TypeMsgFetchPriceData
}

// GetSigners returns the message signers
func (m *MsgFetchPriceData) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(m.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// GetSignBytes returns the signed bytes from the message
func (m *MsgFetchPriceData) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic check the basic message validation
func (m *MsgFetchPriceData) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if m.SourceChannel == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid source channel")
	}
	return nil
}

// FetchPriceResultStoreKey is a function to generate key for each result in store
func FetchPriceResultStoreKey(requestID OracleRequestID) []byte {
	return append(KeyPrefix(FetchPriceResultStoreKeyPrefix), int64ToBytes(int64(requestID))...)
}
