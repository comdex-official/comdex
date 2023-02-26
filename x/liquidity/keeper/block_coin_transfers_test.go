package keeper_test

import (
	utils "github.com/comdex-official/comdex/types"
	_ "github.com/stretchr/testify/suite"
)

func (s *KeeperTestSuite) TestBlockCoinTransferIfNotAlready() {
	coin1 := utils.ParseCoin("1000coin1")
	isBlocked := s.keeper.IsCoinTransferBlocked(s.ctx, coin1)
	s.Require().False(isBlocked)
	s.keeper.BlockCoinTransferIfNotAlready(s.ctx, coin1)
	isBlocked = s.keeper.IsCoinTransferBlocked(s.ctx, coin1)
	s.Require().True(isBlocked)
}
