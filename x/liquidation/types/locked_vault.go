package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m *WhitelistedAppIds) ValidateBasic() error {
	if m.WhitelistedAppMappingIds == nil {
		return fmt.Errorf("WhitelistedAppIds cannot be empty")
	}

	return nil
}

func (m *WhitelistedAppIds) GetSigners() []sdk.AccAddress {
	return nil
}

func (m *WhitelistedAppIds) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}
