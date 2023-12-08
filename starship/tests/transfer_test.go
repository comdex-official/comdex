package main

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func (s *TestSuite) TestChainTokenTransfer() {
	chain1, err := s.chainClients.GetChainClient("comdex-1")
	s.Require().NoError(err)

	keyName := "test-transfer"
	address, err := chain1.CreateRandWallet(keyName)
	s.Require().NoError(err)

	denom, err := chain1.GetChainDenom()
	s.Require().NoError(err)

	s.TransferTokens(chain1, address, 2345000, denom)

	// Verify the address recived the token
	balance, err := chain1.Client.QueryBalanceWithDenomTraces(context.Background(), sdk.MustAccAddressFromBech32(address), nil)
	s.Require().NoError(err)

	// Assert correct transfers
	s.Assert().Len(balance, 1)
	s.Assert().Equal(balance.Denoms(), []string{denom})
	s.Assert().Equal(balance[0].Amount, sdk.NewInt(2345000))
}

func (s *TestSuite) TestChainIBCTransfer() {
	chain2, err := s.chainClients.GetChainClient("gaia-1")
	s.Require().NoError(err)
	chain1, err := s.chainClients.GetChainClient("comdex-1")
	s.Require().NoError(err)

	keyName := "test-ibc-transfer"
	address, err := chain1.CreateRandWallet(keyName)
	s.Require().NoError(err)

	// Tranfer atom to comdex chain
	s.IBCTransferTokens(chain2, chain1, address, 12345000)

	// Verify the address recived the token
	balances, err := banktypes.NewQueryClient(chain1.Client).AllBalances(context.Background(), &banktypes.QueryAllBalancesRequest{
		Address: address,
	})
	s.Require().NoError(err)

	// Assert correct transfers
	s.Assert().Len(balances.Balances.Denoms(), 1)
	s.Assert().Equal(balances.Balances[0].Amount.Uint64(), uint64(12345000))
	s.Assert().Contains(balances.Balances.Denoms()[0], "ibc/")
}
