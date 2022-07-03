package testutil

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestVaultIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(VaultIntegrationTestSuite))
}
