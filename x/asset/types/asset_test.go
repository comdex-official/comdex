package types

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidateAsset(t *testing.T) {
	invalidAsset := []Asset{
		{0, "cmdx", "strDenom", 1},
		{1, "", "strDenom", 1},
		{1, "Name_length_greater", "strDenom", 1},
		{1, "cmdx", "", 1},
		{1, "cmdx", "\u5586", 1},
		{1, "cmdx", "strDenom", -1},
	}

	validAsset := Asset{
		1, "cmdx", "strDenom", 1,
	}

	for _, asset := range invalidAsset {
		err := asset.Validate()
		require.Error(t, err)
	}

	err := validAsset.Validate()
	require.NoError(t, err)
}
