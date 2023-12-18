package v13_test

// import (
// 	"testing"

// 	"github.com/stretchr/testify/suite"

// 	"github.com/comdex-official/comdex/app"
// 	v13 "github.com/comdex-official/comdex/app/upgrades/mainnet/v13"
// )

// type UpgradeTestSuite struct {
// 	app.KeeperTestHelper
// }

// func (s *UpgradeTestSuite) SetupTest() {
// 	s.Setup()
// }

// func TestKeeperTestSuite(t *testing.T) {
// 	suite.Run(t, new(UpgradeTestSuite))
// }

// // Ensures the test does not error out.
// func (s *UpgradeTestSuite) TestUpgrade() {
// 	s.Setup()

// 	preUpgradeChecks(s)

// 	upgradeHeight := int64(5)
// 	s.ConfirmUpgradeSucceeded(v13.UpgradeName, upgradeHeight)

// 	postUpgradeChecks(s)
// }

// func preUpgradeChecks(s *UpgradeTestSuite) {

// }

// func postUpgradeChecks(s *UpgradeTestSuite) {

// 	// Ensure the gov params have MinInitialDepositRatio added
// 	gp := s.App.GovKeeper.GetParams(s.Ctx)
// 	s.Require().Equal(gp.MinInitialDepositRatio, "0.200000000000000000")
// }
