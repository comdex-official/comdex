package testutil

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestLiquidityIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(LiquidityIntegrationTestSuite))
}
