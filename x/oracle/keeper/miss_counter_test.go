package keeper_test

import (
	"math/rand"

	"github.com/comdex-official/comdex/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func (s *IntegrationTestSuite) TestMissCounter() {
	app, ctx := s.app, s.ctx
	missCounter := uint64(rand.Intn(100))

	s.Require().Equal(app.OracleKeeper.GetMissCounter(ctx, valAddr), uint64(0))
	app.OracleKeeper.SetMissCounter(ctx, valAddr, missCounter)
	s.Require().Equal(app.OracleKeeper.GetMissCounter(ctx, valAddr), missCounter)

	app.OracleKeeper.DeleteMissCounter(ctx, valAddr)
	s.Require().Equal(app.OracleKeeper.GetMissCounter(ctx, valAddr), uint64(0))
}

func (s *IntegrationTestSuite) TestIterateMissCounters() {
	keeper, ctx := s.app.OracleKeeper, s.ctx
	missCounters := []types.MissCounter{
		{ValidatorAddress: valAddr.String(), MissCounter: uint64(4)},
		{ValidatorAddress: valAddr2.String(), MissCounter: uint64(87)},
	}

	for _, mc := range missCounters {
		operator, _ := sdk.ValAddressFromBech32(mc.ValidatorAddress)
		keeper.SetMissCounter(ctx, operator, mc.MissCounter)
	}

	newCounters := []types.MissCounter{}
	keeper.IterateMissCounters(ctx, func(operator sdk.ValAddress, missCounter uint64) (stop bool) {
		newCounters = append(newCounters, types.MissCounter{
			ValidatorAddress: operator.String(),
			MissCounter:      missCounter,
		})

		return false
	})
	require.Equal(s.T(), len(missCounters), len(newCounters))

FOUND:
	for _, oldCounter := range missCounters {
		for _, newCounter := range newCounters {
			if oldCounter.ValidatorAddress == newCounter.ValidatorAddress {
				s.Require().Equal(oldCounter.MissCounter, newCounter.MissCounter)
				continue FOUND
			}
		}
		s.T().Errorf("did not find match for miss counter: %+v", oldCounter)
	}
}
