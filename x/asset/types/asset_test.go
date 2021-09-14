package types

import (
	"github.com/stretchr/testify/require"
	"testing"
)


func TestValidateAsset(t *testing.T) {
	invalidAsset := []Asset{
		{0, "abc", "def", 1},
		{1, "", "def", 1},
		{1, "qwertyuioplkjhgfd", "def", 1},
		{1, "abc", "", 1},
		{1,"abcd","\u5586",1},
		{1, "abc", "def", -1},
	}

	validAsset := Asset{
		1, "sankalp", "singh", 1,
	}

	for _, asset := range invalidAsset {
		err := asset.Validate()
		require.Error(t, err)
	}

	err := validAsset.Validate()
	require.NoError(t, err)
}
