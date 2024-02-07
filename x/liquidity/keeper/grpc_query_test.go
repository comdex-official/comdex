package keeper_test

import (
	"time"

	errorsmod "cosmossdk.io/errors"

	utils "github.com/comdex-official/comdex/types"
	"github.com/comdex-official/comdex/x/liquidity/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *KeeperTestSuite) TestParams() {
	resp, err := s.querier.Params(sdk.WrapSDKContext(s.ctx), &types.QueryParamsRequest{})
	s.Require().NoError(err)
	s.Require().Equal(&types.QueryParamsResponse{}, resp)
}

func (s *KeeperTestSuite) TestGenericParams() {
	appID1 := s.CreateNewApp("appone")

	testCases := []struct {
		Name   string
		Req    *types.QueryGenericParamsRequest
		ExpErr error
	}{
		{
			Name:   "error empty request",
			Req:    nil,
			ExpErr: status.Error(codes.InvalidArgument, "empty request"),
		},
		{
			Name:   "error app id 0",
			Req:    &types.QueryGenericParamsRequest{},
			ExpErr: status.Error(codes.InvalidArgument, "app id cannot be 0"),
		},
		{
			Name:   "error app id invalid",
			Req:    &types.QueryGenericParamsRequest{AppId: 6969},
			ExpErr: errorsmod.Wrapf(types.ErrInvalidAppID, "app id 6969 not found"),
		},
		{
			Name:   "success",
			Req:    &types.QueryGenericParamsRequest{AppId: appID1},
			ExpErr: nil,
		},
	}

	for _, tc := range testCases {
		ctx := sdk.WrapSDKContext(s.ctx)
		s.Run(tc.Name, func() {
			resp, err := s.querier.GenericParams(ctx, tc.Req)
			if tc.ExpErr != nil {
				s.Require().Error(err)
				s.Require().EqualError(err, tc.ExpErr.Error())
			} else {
				s.Require().NoError(err)
				s.Require().Equal(tc.Req.AppId, resp.Params.AppId)
			}
		})
	}
}

func (s *KeeperTestSuite) TestPools() {
	addr1 := s.addr(1)

	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "uasset2", 1000000)

	pair := s.CreateNewLiquidityPair(appID1, addr1, asset1.Denom, asset2.Denom)
	pool := s.CreateNewLiquidityPool(appID1, pair.Id, addr1, "1000000000000uasset1,1000000000000uasset2")

	ctx := sdk.WrapSDKContext(s.ctx)
	_, err := s.querier.Pools(ctx, nil)
	s.Require().Error(err)
	s.Require().EqualError(err, status.Error(codes.InvalidArgument, "empty request").Error())

	_, err = s.querier.Pools(ctx, &types.QueryPoolsRequest{})
	s.Require().Error(err)
	s.Require().EqualError(err, status.Error(codes.InvalidArgument, "app id cannot be 0").Error())

	resp, err := s.querier.Pools(ctx, &types.QueryPoolsRequest{AppId: appID1, PairId: pair.Id})
	s.Require().NoError(err)
	s.Require().Len(resp.Pools, 1)

	s.Require().Equal(resp.Pools[0].Id, pool.Id)
	s.Require().Equal(resp.Pools[0].PairId, pool.PairId)
	s.Require().Equal(resp.Pools[0].ReserveAddress, pool.ReserveAddress)
	s.Require().Equal(resp.Pools[0].PoolCoinDenom, pool.PoolCoinDenom)
	s.Require().Equal(resp.Pools[0].Balances.BaseCoin, s.getBalance(pool.GetReserveAddress(), resp.Pools[0].Balances.BaseCoin.Denom))
	s.Require().Equal(resp.Pools[0].Balances.QuoteCoin, s.getBalance(pool.GetReserveAddress(), resp.Pools[0].Balances.QuoteCoin.Denom))
	s.Require().Equal(resp.Pools[0].LastDepositRequestId, pool.LastDepositRequestId)
	s.Require().Equal(resp.Pools[0].LastWithdrawRequestId, pool.LastWithdrawRequestId)
	s.Require().Equal(resp.Pools[0].AppId, pool.AppId)
	s.Require().Equal(resp.Pools[0].PoolCoinSupply, s.getBalance(addr1, pool.PoolCoinDenom).Amount)
}

func (s *KeeperTestSuite) TestPool() {
	addr1 := s.addr(1)

	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "uasset2", 1000000)

	pair := s.CreateNewLiquidityPair(appID1, addr1, asset1.Denom, asset2.Denom)
	pool := s.CreateNewLiquidityPool(appID1, pair.Id, addr1, "1000000000000uasset1,1000000000000uasset2")

	testCases := []struct {
		Name   string
		Req    *types.QueryPoolRequest
		ExpErr error
	}{
		{
			Name:   "error empty request",
			Req:    nil,
			ExpErr: status.Error(codes.InvalidArgument, "empty request"),
		},
		{
			Name:   "error app id 0",
			Req:    &types.QueryPoolRequest{},
			ExpErr: status.Error(codes.InvalidArgument, "app id cannot be 0"),
		},
		{
			Name:   "error pool id 0",
			Req:    &types.QueryPoolRequest{AppId: appID1},
			ExpErr: status.Error(codes.InvalidArgument, "pool id cannot be 0"),
		},
		{
			Name:   "error pool id invalid",
			Req:    &types.QueryPoolRequest{AppId: appID1, PoolId: 123},
			ExpErr: status.Errorf(codes.NotFound, "pool 123 doesn't exist"),
		},
		{
			Name:   "success",
			Req:    &types.QueryPoolRequest{AppId: appID1, PoolId: pool.Id},
			ExpErr: nil,
		},
	}

	ctx := sdk.WrapSDKContext(s.ctx)
	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			resp, err := s.querier.Pool(ctx, tc.Req)
			if tc.ExpErr != nil {
				s.Require().Error(err)
				s.Require().EqualError(err, tc.ExpErr.Error())
			} else {
				s.Require().NoError(err)
				s.Require().Equal(resp.Pool.Id, pool.Id)
				s.Require().Equal(resp.Pool.PairId, pool.PairId)
				s.Require().Equal(resp.Pool.ReserveAddress, pool.ReserveAddress)
				s.Require().Equal(resp.Pool.PoolCoinDenom, pool.PoolCoinDenom)
				s.Require().Equal(resp.Pool.Balances.BaseCoin, s.getBalance(pool.GetReserveAddress(), resp.Pool.Balances.BaseCoin.Denom))
				s.Require().Equal(resp.Pool.Balances.QuoteCoin, s.getBalance(pool.GetReserveAddress(), resp.Pool.Balances.QuoteCoin.Denom))
				s.Require().Equal(resp.Pool.LastDepositRequestId, pool.LastDepositRequestId)
				s.Require().Equal(resp.Pool.LastWithdrawRequestId, pool.LastWithdrawRequestId)
				s.Require().Equal(resp.Pool.AppId, pool.AppId)
				s.Require().Equal(resp.Pool.PoolCoinSupply, s.getBalance(addr1, pool.PoolCoinDenom).Amount)
			}
		})
	}
}

func (s *KeeperTestSuite) TestPoolByReserveAddress() {
	addr1 := s.addr(1)

	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "uasset2", 1000000)

	pair := s.CreateNewLiquidityPair(appID1, addr1, asset1.Denom, asset2.Denom)
	pool := s.CreateNewLiquidityPool(appID1, pair.Id, addr1, "1000000000000uasset1,1000000000000uasset2")

	testCases := []struct {
		Name   string
		Req    *types.QueryPoolByReserveAddressRequest
		ExpErr error
	}{
		{
			Name:   "error empty request",
			Req:    nil,
			ExpErr: status.Error(codes.InvalidArgument, "empty request"),
		},
		{
			Name:   "error app id 0",
			Req:    &types.QueryPoolByReserveAddressRequest{},
			ExpErr: status.Error(codes.InvalidArgument, "app id cannot be 0"),
		},
		{
			Name:   "error empty reserve address",
			Req:    &types.QueryPoolByReserveAddressRequest{AppId: appID1, ReserveAddress: ""},
			ExpErr: status.Error(codes.InvalidArgument, "empty reserve account address"),
		},
		{
			Name:   "error reserve address invalid",
			Req:    &types.QueryPoolByReserveAddressRequest{AppId: appID1, ReserveAddress: "1234"},
			ExpErr: status.Errorf(codes.InvalidArgument, "reserve account address 1234 is not valid"),
		},
		{
			Name:   "error no pool for address",
			Req:    &types.QueryPoolByReserveAddressRequest{AppId: appID1, ReserveAddress: addr1.String()},
			ExpErr: status.Errorf(codes.NotFound, "pool by %s doesn't exist", addr1.String()),
		},
		{
			Name:   "success",
			Req:    &types.QueryPoolByReserveAddressRequest{AppId: appID1, ReserveAddress: pool.GetReserveAddress().String()},
			ExpErr: nil,
		},
	}

	ctx := sdk.WrapSDKContext(s.ctx)
	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			resp, err := s.querier.PoolByReserveAddress(ctx, tc.Req)
			if tc.ExpErr != nil {
				s.Require().Error(err)
				s.Require().EqualError(err, tc.ExpErr.Error())
			} else {
				s.Require().NoError(err)
				s.Require().Equal(resp.Pool.Id, pool.Id)
				s.Require().Equal(resp.Pool.PairId, pool.PairId)
				s.Require().Equal(resp.Pool.ReserveAddress, pool.ReserveAddress)
				s.Require().Equal(resp.Pool.PoolCoinDenom, pool.PoolCoinDenom)
				s.Require().Equal(resp.Pool.Balances.BaseCoin, s.getBalance(pool.GetReserveAddress(), resp.Pool.Balances.BaseCoin.Denom))
				s.Require().Equal(resp.Pool.Balances.QuoteCoin, s.getBalance(pool.GetReserveAddress(), resp.Pool.Balances.QuoteCoin.Denom))
				s.Require().Equal(resp.Pool.LastDepositRequestId, pool.LastDepositRequestId)
				s.Require().Equal(resp.Pool.LastWithdrawRequestId, pool.LastWithdrawRequestId)
				s.Require().Equal(resp.Pool.AppId, pool.AppId)
				s.Require().Equal(resp.Pool.PoolCoinSupply, s.getBalance(addr1, pool.PoolCoinDenom).Amount)
			}
		})
	}
}

func (s *KeeperTestSuite) TestPoolByPoolCoinDenom() {
	addr1 := s.addr(1)

	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "uasset2", 1000000)

	pair := s.CreateNewLiquidityPair(appID1, addr1, asset1.Denom, asset2.Denom)
	pool := s.CreateNewLiquidityPool(appID1, pair.Id, addr1, "1000000000000uasset1,1000000000000uasset2")

	testCases := []struct {
		Name   string
		Req    *types.QueryPoolByPoolCoinDenomRequest
		ExpErr error
	}{
		{
			Name:   "error empty request",
			Req:    nil,
			ExpErr: status.Error(codes.InvalidArgument, "empty request"),
		},
		{
			Name:   "error app id 0",
			Req:    &types.QueryPoolByPoolCoinDenomRequest{},
			ExpErr: status.Error(codes.InvalidArgument, "app id cannot be 0"),
		},
		{
			Name:   "error empty pool coin denom",
			Req:    &types.QueryPoolByPoolCoinDenomRequest{AppId: appID1, PoolCoinDenom: ""},
			ExpErr: status.Error(codes.InvalidArgument, "empty pool coin denom"),
		},
		{
			Name:   "error invalid pool coin denom",
			Req:    &types.QueryPoolByPoolCoinDenomRequest{AppId: appID1, PoolCoinDenom: "pool1"},
			ExpErr: status.Errorf(codes.InvalidArgument, "failed to parse pool coin denom: pool1 is not a pool coin denom"),
		},
		{
			Name:   "error no pool for denom",
			Req:    &types.QueryPoolByPoolCoinDenomRequest{AppId: appID1, PoolCoinDenom: "pool1-2"},
			ExpErr: status.Errorf(codes.NotFound, "pool 2 doesn't exist"),
		},
		{
			Name:   "success",
			Req:    &types.QueryPoolByPoolCoinDenomRequest{AppId: appID1, PoolCoinDenom: "pool1-1"},
			ExpErr: nil,
		},
	}

	ctx := sdk.WrapSDKContext(s.ctx)
	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			resp, err := s.querier.PoolByPoolCoinDenom(ctx, tc.Req)
			if tc.ExpErr != nil {
				s.Require().Error(err)
				s.Require().EqualError(err, tc.ExpErr.Error())
			} else {
				s.Require().NoError(err)
				s.Require().Equal(resp.Pool.Id, pool.Id)
				s.Require().Equal(resp.Pool.PairId, pool.PairId)
				s.Require().Equal(resp.Pool.ReserveAddress, pool.ReserveAddress)
				s.Require().Equal(resp.Pool.PoolCoinDenom, pool.PoolCoinDenom)
				s.Require().Equal(resp.Pool.Balances.BaseCoin, s.getBalance(pool.GetReserveAddress(), resp.Pool.Balances.BaseCoin.Denom))
				s.Require().Equal(resp.Pool.Balances.QuoteCoin, s.getBalance(pool.GetReserveAddress(), resp.Pool.Balances.QuoteCoin.Denom))
				s.Require().Equal(resp.Pool.LastDepositRequestId, pool.LastDepositRequestId)
				s.Require().Equal(resp.Pool.LastWithdrawRequestId, pool.LastWithdrawRequestId)
				s.Require().Equal(resp.Pool.AppId, pool.AppId)
				s.Require().Equal(resp.Pool.PoolCoinSupply, s.getBalance(addr1, pool.PoolCoinDenom).Amount)
			}
		})
	}
}

func (s *KeeperTestSuite) TestPairs() {
	addr1 := s.addr(1)

	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "uasset2", 1000000)

	pair := s.CreateNewLiquidityPair(appID1, addr1, asset1.Denom, asset2.Denom)

	testCases := []struct {
		Name   string
		Req    *types.QueryPairsRequest
		ExpErr error
	}{
		{
			Name:   "error empty request",
			Req:    nil,
			ExpErr: status.Error(codes.InvalidArgument, "empty request"),
		},
		{
			Name:   "error app id 0",
			Req:    &types.QueryPairsRequest{},
			ExpErr: status.Error(codes.InvalidArgument, "app id cannot be 0"),
		},
		{
			Name:   "error too many denoms",
			Req:    &types.QueryPairsRequest{AppId: appID1, Denoms: []string{"uasset1", "uasset2", "uasset3"}},
			ExpErr: status.Errorf(codes.InvalidArgument, "too many denoms to query: 3"),
		},
		{
			Name:   "success query by denom 1",
			Req:    &types.QueryPairsRequest{AppId: appID1, Denoms: []string{"uasset1"}},
			ExpErr: nil,
		},
		{
			Name:   "success query by denom 2",
			Req:    &types.QueryPairsRequest{AppId: appID1, Denoms: []string{"uasset2"}},
			ExpErr: nil,
		},
		{
			Name:   "success query both denoms",
			Req:    &types.QueryPairsRequest{AppId: appID1, Denoms: []string{"uasset2", "uasset1"}},
			ExpErr: nil,
		},
	}

	ctx := sdk.WrapSDKContext(s.ctx)
	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			resp, err := s.querier.Pairs(ctx, tc.Req)
			if tc.ExpErr != nil {
				s.Require().Error(err)
				s.Require().EqualError(err, tc.ExpErr.Error())
			} else {
				s.Require().NoError(err)
				s.Require().Equal(resp.Pairs[0].Id, pair.Id)
				s.Require().Equal(resp.Pairs[0].BaseCoinDenom, pair.BaseCoinDenom)
				s.Require().Equal(resp.Pairs[0].QuoteCoinDenom, pair.QuoteCoinDenom)
				s.Require().Equal(resp.Pairs[0].EscrowAddress, pair.EscrowAddress)
				s.Require().Equal(resp.Pairs[0].LastOrderId, pair.LastOrderId)
				s.Require().Equal(resp.Pairs[0].LastPrice, pair.LastPrice)
				s.Require().Equal(resp.Pairs[0].CurrentBatchId, pair.CurrentBatchId)
				s.Require().Equal(resp.Pairs[0].SwapFeeCollectorAddress, pair.SwapFeeCollectorAddress)
				s.Require().Equal(resp.Pairs[0].AppId, pair.AppId)
			}
		})
	}
}

func (s *KeeperTestSuite) TestPair() {
	addr1 := s.addr(1)

	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "uasset2", 1000000)

	pair := s.CreateNewLiquidityPair(appID1, addr1, asset1.Denom, asset2.Denom)

	testCases := []struct {
		Name   string
		Req    *types.QueryPairRequest
		ExpErr error
	}{
		{
			Name:   "error empty request",
			Req:    nil,
			ExpErr: status.Error(codes.InvalidArgument, "empty request"),
		},
		{
			Name:   "error app id 0",
			Req:    &types.QueryPairRequest{},
			ExpErr: status.Error(codes.InvalidArgument, "app id cannot be 0"),
		},
		{
			Name:   "error pair id 0",
			Req:    &types.QueryPairRequest{AppId: appID1},
			ExpErr: status.Error(codes.InvalidArgument, "pair id cannot be 0"),
		},
		{
			Name:   "error pair id invalid",
			Req:    &types.QueryPairRequest{AppId: appID1, PairId: 6969},
			ExpErr: status.Errorf(codes.NotFound, "pair 6969 doesn't exist"),
		},
		{
			Name:   "success",
			Req:    &types.QueryPairRequest{AppId: appID1, PairId: pair.Id},
			ExpErr: nil,
		},
	}

	ctx := sdk.WrapSDKContext(s.ctx)
	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			resp, err := s.querier.Pair(ctx, tc.Req)
			if tc.ExpErr != nil {
				s.Require().Error(err)
				s.Require().EqualError(err, tc.ExpErr.Error())
			} else {
				s.Require().NoError(err)
				s.Require().Equal(resp.Pair.Id, pair.Id)
				s.Require().Equal(resp.Pair.BaseCoinDenom, pair.BaseCoinDenom)
				s.Require().Equal(resp.Pair.QuoteCoinDenom, pair.QuoteCoinDenom)
				s.Require().Equal(resp.Pair.EscrowAddress, pair.EscrowAddress)
				s.Require().Equal(resp.Pair.LastOrderId, pair.LastOrderId)
				s.Require().Equal(resp.Pair.LastPrice, pair.LastPrice)
				s.Require().Equal(resp.Pair.CurrentBatchId, pair.CurrentBatchId)
				s.Require().Equal(resp.Pair.SwapFeeCollectorAddress, pair.SwapFeeCollectorAddress)
				s.Require().Equal(resp.Pair.AppId, pair.AppId)
			}
		})
	}
}

func (s *KeeperTestSuite) TestDepositRequests() {
	addr1 := s.addr(1)

	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "uasset2", 1000000)

	pair := s.CreateNewLiquidityPair(appID1, addr1, asset1.Denom, asset2.Denom)
	pool := s.CreateNewLiquidityPool(appID1, pair.Id, addr1, "1000000000000uasset1,1000000000000uasset2")

	req := s.Deposit(appID1, pool.Id, addr1, "1000000uasset1,1000000uasset2")

	testCases := []struct {
		Name   string
		Req    *types.QueryDepositRequestsRequest
		ExpErr error
	}{
		{
			Name:   "error empty request",
			Req:    nil,
			ExpErr: status.Error(codes.InvalidArgument, "empty request"),
		},
		{
			Name:   "error app id 0",
			Req:    &types.QueryDepositRequestsRequest{},
			ExpErr: status.Error(codes.InvalidArgument, "app id cannot be 0"),
		},
		{
			Name:   "error pool id 0",
			Req:    &types.QueryDepositRequestsRequest{AppId: appID1},
			ExpErr: status.Error(codes.InvalidArgument, "pool id cannot be 0"),
		},
		{
			Name:   "success",
			Req:    &types.QueryDepositRequestsRequest{AppId: appID1, PoolId: pool.Id},
			ExpErr: nil,
		},
	}

	ctx := sdk.WrapSDKContext(s.ctx)
	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			resp, err := s.querier.DepositRequests(ctx, tc.Req)
			if tc.ExpErr != nil {
				s.Require().Error(err)
				s.Require().EqualError(err, tc.ExpErr.Error())
			} else {
				s.Require().NoError(err)
				s.Require().Equal(resp.DepositRequests[0].Id, req.Id)
				s.Require().Equal(resp.DepositRequests[0].PoolId, req.PoolId)
				s.Require().Equal(resp.DepositRequests[0].MsgHeight, req.MsgHeight)
				s.Require().Equal(resp.DepositRequests[0].Depositor, req.Depositor)
				s.Require().Equal(resp.DepositRequests[0].DepositCoins, req.DepositCoins)
				s.Require().Equal(resp.DepositRequests[0].AcceptedCoins, req.AcceptedCoins)
				s.Require().Equal(resp.DepositRequests[0].MintedPoolCoin, req.MintedPoolCoin)
				s.Require().Equal(resp.DepositRequests[0].Status, req.Status)
				s.Require().Equal(resp.DepositRequests[0].AppId, req.AppId)
			}
		})
	}
}

func (s *KeeperTestSuite) TestDepositRequest() {
	addr1 := s.addr(1)

	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "uasset2", 1000000)

	pair := s.CreateNewLiquidityPair(appID1, addr1, asset1.Denom, asset2.Denom)
	pool := s.CreateNewLiquidityPool(appID1, pair.Id, addr1, "1000000000000uasset1,1000000000000uasset2")

	req := s.Deposit(appID1, pool.Id, addr1, "1000000uasset1,1000000uasset2")

	testCases := []struct {
		Name   string
		Req    *types.QueryDepositRequestRequest
		ExpErr error
	}{
		{
			Name:   "error empty request",
			Req:    nil,
			ExpErr: status.Error(codes.InvalidArgument, "empty request"),
		},
		{
			Name:   "error app id 0",
			Req:    &types.QueryDepositRequestRequest{},
			ExpErr: status.Error(codes.InvalidArgument, "app id cannot be 0"),
		},
		{
			Name:   "error pool id 0",
			Req:    &types.QueryDepositRequestRequest{AppId: appID1},
			ExpErr: status.Error(codes.InvalidArgument, "pool id cannot be 0"),
		},
		{
			Name:   "error id 0",
			Req:    &types.QueryDepositRequestRequest{AppId: appID1, PoolId: pool.Id},
			ExpErr: status.Error(codes.InvalidArgument, "id cannot be 0"),
		},
		{
			Name:   "error id invalid",
			Req:    &types.QueryDepositRequestRequest{AppId: appID1, PoolId: pool.Id, Id: 6969},
			ExpErr: status.Errorf(codes.NotFound, "deposit request of pool id 1 and request id 6969 doesn't exist or deleted"),
		},
		{
			Name:   "success",
			Req:    &types.QueryDepositRequestRequest{AppId: appID1, PoolId: pool.Id, Id: req.Id},
			ExpErr: nil,
		},
	}

	ctx := sdk.WrapSDKContext(s.ctx)
	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			resp, err := s.querier.DepositRequest(ctx, tc.Req)
			if tc.ExpErr != nil {
				s.Require().Error(err)
				s.Require().EqualError(err, tc.ExpErr.Error())
			} else {
				s.Require().NoError(err)
				s.Require().Equal(resp.DepositRequest.Id, req.Id)
				s.Require().Equal(resp.DepositRequest.PoolId, req.PoolId)
				s.Require().Equal(resp.DepositRequest.MsgHeight, req.MsgHeight)
				s.Require().Equal(resp.DepositRequest.Depositor, req.Depositor)
				s.Require().Equal(resp.DepositRequest.DepositCoins, req.DepositCoins)
				s.Require().Equal(resp.DepositRequest.AcceptedCoins, req.AcceptedCoins)
				s.Require().Equal(resp.DepositRequest.MintedPoolCoin, req.MintedPoolCoin)
				s.Require().Equal(resp.DepositRequest.Status, req.Status)
				s.Require().Equal(resp.DepositRequest.AppId, req.AppId)
			}
		})
	}
}

func (s *KeeperTestSuite) TestWithdrawRequests() {
	addr1 := s.addr(1)

	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "uasset2", 1000000)

	pair := s.CreateNewLiquidityPair(appID1, addr1, asset1.Denom, asset2.Denom)
	pool := s.CreateNewLiquidityPool(appID1, pair.Id, addr1, "1000000000000uasset1,1000000000000uasset2")

	_ = s.Deposit(appID1, pool.Id, addr1, "1000000uasset1,1000000uasset2")
	s.nextBlock()
	reqWdrw := s.Withdraw(appID1, pool.Id, addr1, utils.ParseCoin("100000pool1-1"))

	testCases := []struct {
		Name   string
		Req    *types.QueryWithdrawRequestsRequest
		ExpErr error
	}{
		{
			Name:   "error empty request",
			Req:    nil,
			ExpErr: status.Error(codes.InvalidArgument, "empty request"),
		},
		{
			Name:   "error app id 0",
			Req:    &types.QueryWithdrawRequestsRequest{},
			ExpErr: status.Error(codes.InvalidArgument, "app id cannot be 0"),
		},
		{
			Name:   "error pool id 0",
			Req:    &types.QueryWithdrawRequestsRequest{AppId: appID1},
			ExpErr: status.Error(codes.InvalidArgument, "pool id cannot be 0"),
		},
		{
			Name:   "success",
			Req:    &types.QueryWithdrawRequestsRequest{AppId: appID1, PoolId: pool.Id},
			ExpErr: nil,
		},
	}

	ctx := sdk.WrapSDKContext(s.ctx)
	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			resp, err := s.querier.WithdrawRequests(ctx, tc.Req)
			if tc.ExpErr != nil {
				s.Require().Error(err)
				s.Require().EqualError(err, tc.ExpErr.Error())
			} else {
				s.Require().NoError(err)
				s.Require().Equal(resp.WithdrawRequests[0].Id, reqWdrw.Id)
				s.Require().Equal(resp.WithdrawRequests[0].PoolId, reqWdrw.PoolId)
				s.Require().Equal(resp.WithdrawRequests[0].MsgHeight, reqWdrw.MsgHeight)
				s.Require().Equal(resp.WithdrawRequests[0].Withdrawer, reqWdrw.Withdrawer)
				s.Require().Equal(resp.WithdrawRequests[0].PoolCoin, reqWdrw.PoolCoin)
				s.Require().Equal(resp.WithdrawRequests[0].WithdrawnCoins, reqWdrw.WithdrawnCoins)
				s.Require().Equal(resp.WithdrawRequests[0].Status, reqWdrw.Status)
				s.Require().Equal(resp.WithdrawRequests[0].AppId, reqWdrw.AppId)
			}
		})
	}
}

func (s *KeeperTestSuite) TestWithdrawRequest() {
	addr1 := s.addr(1)

	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "uasset2", 1000000)

	pair := s.CreateNewLiquidityPair(appID1, addr1, asset1.Denom, asset2.Denom)
	pool := s.CreateNewLiquidityPool(appID1, pair.Id, addr1, "1000000000000uasset1,1000000000000uasset2")

	_ = s.Deposit(appID1, pool.Id, addr1, "1000000uasset1,1000000uasset2")
	s.nextBlock()
	reqWdrw := s.Withdraw(appID1, pool.Id, addr1, utils.ParseCoin("100000pool1-1"))

	testCases := []struct {
		Name   string
		Req    *types.QueryWithdrawRequestRequest
		ExpErr error
	}{
		{
			Name:   "error empty request",
			Req:    nil,
			ExpErr: status.Error(codes.InvalidArgument, "empty request"),
		},
		{
			Name:   "error app id 0",
			Req:    &types.QueryWithdrawRequestRequest{},
			ExpErr: status.Error(codes.InvalidArgument, "app id cannot be 0"),
		},
		{
			Name:   "error pool id 0",
			Req:    &types.QueryWithdrawRequestRequest{AppId: appID1},
			ExpErr: status.Error(codes.InvalidArgument, "pool id cannot be 0"),
		},
		{
			Name:   "error  id 0",
			Req:    &types.QueryWithdrawRequestRequest{AppId: appID1, PoolId: pool.Id},
			ExpErr: status.Error(codes.InvalidArgument, "id cannot be 0"),
		},
		{
			Name:   "error  id invalid",
			Req:    &types.QueryWithdrawRequestRequest{AppId: appID1, PoolId: pool.Id, Id: 123},
			ExpErr: status.Errorf(codes.NotFound, "withdraw request of pool id 1 and request id 123 doesn't exist or deleted"),
		},
		{
			Name:   "success",
			Req:    &types.QueryWithdrawRequestRequest{AppId: appID1, PoolId: pool.Id, Id: reqWdrw.Id},
			ExpErr: nil,
		},
	}

	ctx := sdk.WrapSDKContext(s.ctx)
	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			resp, err := s.querier.WithdrawRequest(ctx, tc.Req)
			if tc.ExpErr != nil {
				s.Require().Error(err)
				s.Require().EqualError(err, tc.ExpErr.Error())
			} else {
				s.Require().NoError(err)
				s.Require().Equal(resp.WithdrawRequest.Id, reqWdrw.Id)
				s.Require().Equal(resp.WithdrawRequest.PoolId, reqWdrw.PoolId)
				s.Require().Equal(resp.WithdrawRequest.MsgHeight, reqWdrw.MsgHeight)
				s.Require().Equal(resp.WithdrawRequest.Withdrawer, reqWdrw.Withdrawer)
				s.Require().Equal(resp.WithdrawRequest.PoolCoin, reqWdrw.PoolCoin)
				s.Require().Equal(resp.WithdrawRequest.WithdrawnCoins, reqWdrw.WithdrawnCoins)
				s.Require().Equal(resp.WithdrawRequest.Status, reqWdrw.Status)
				s.Require().Equal(resp.WithdrawRequest.AppId, reqWdrw.AppId)
			}
		})
	}
}

func (s *KeeperTestSuite) TestOrders() {
	addr1 := s.addr(1)

	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "uasset2", 1000000)

	pair := s.CreateNewLiquidityPair(appID1, addr1, asset1.Denom, asset2.Denom)
	_ = s.CreateNewLiquidityPool(appID1, pair.Id, addr1, "1000000000000uasset1,1000000000000uasset2")

	order := s.LimitOrder(appID1, addr1, pair.Id, types.OrderDirectionSell, newDec(1), newInt(100000), time.Second*10)

	testCases := []struct {
		Name   string
		Req    *types.QueryOrdersRequest
		ExpErr error
	}{
		{
			Name:   "error empty request",
			Req:    nil,
			ExpErr: status.Error(codes.InvalidArgument, "empty request"),
		},
		{
			Name:   "error app id 0",
			Req:    &types.QueryOrdersRequest{},
			ExpErr: status.Error(codes.InvalidArgument, "app id cannot be 0"),
		},
		{
			Name:   "error pair id 0",
			Req:    &types.QueryOrdersRequest{AppId: appID1},
			ExpErr: status.Error(codes.InvalidArgument, "pair id cannot be 0"),
		},
		{
			Name:   "success",
			Req:    &types.QueryOrdersRequest{AppId: appID1, PairId: pair.Id},
			ExpErr: nil,
		},
	}

	ctx := sdk.WrapSDKContext(s.ctx)
	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			resp, err := s.querier.Orders(ctx, tc.Req)
			if tc.ExpErr != nil {
				s.Require().Error(err)
				s.Require().EqualError(err, tc.ExpErr.Error())
			} else {
				s.Require().NoError(err)
				s.Require().Equal(resp.Orders[0].Id, order.Id)
				s.Require().Equal(resp.Orders[0].PairId, order.PairId)
				s.Require().Equal(resp.Orders[0].MsgHeight, order.MsgHeight)
				s.Require().Equal(resp.Orders[0].Orderer, order.Orderer)
				s.Require().Equal(resp.Orders[0].Direction, order.Direction)
				s.Require().Equal(resp.Orders[0].OfferCoin, order.OfferCoin)
				s.Require().Equal(resp.Orders[0].RemainingOfferCoin, order.RemainingOfferCoin)
				s.Require().Equal(resp.Orders[0].ReceivedCoin, order.ReceivedCoin)
				s.Require().Equal(resp.Orders[0].Price, order.Price)
				s.Require().Equal(resp.Orders[0].Amount, order.Amount)
				s.Require().Equal(resp.Orders[0].OpenAmount, order.OpenAmount)
				s.Require().Equal(resp.Orders[0].BatchId, order.BatchId)
				s.Require().Equal(resp.Orders[0].ExpireAt, order.ExpireAt)
				s.Require().Equal(resp.Orders[0].Status, order.Status)
				s.Require().Equal(resp.Orders[0].AppId, order.AppId)
			}
		})
	}
}

func (s *KeeperTestSuite) TestOrder() {
	addr1 := s.addr(1)

	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "uasset2", 1000000)

	pair := s.CreateNewLiquidityPair(appID1, addr1, asset1.Denom, asset2.Denom)
	_ = s.CreateNewLiquidityPool(appID1, pair.Id, addr1, "1000000000000uasset1,1000000000000uasset2")

	order := s.LimitOrder(appID1, addr1, pair.Id, types.OrderDirectionSell, newDec(1), newInt(100000), time.Second*10)

	testCases := []struct {
		Name   string
		Req    *types.QueryOrderRequest
		ExpErr error
	}{
		{
			Name:   "error empty request",
			Req:    nil,
			ExpErr: status.Error(codes.InvalidArgument, "empty request"),
		},
		{
			Name:   "error app id 0",
			Req:    &types.QueryOrderRequest{},
			ExpErr: status.Error(codes.InvalidArgument, "app id cannot be 0"),
		},
		{
			Name:   "error pair id 0",
			Req:    &types.QueryOrderRequest{AppId: appID1},
			ExpErr: status.Error(codes.InvalidArgument, "pair id cannot be 0"),
		},
		{
			Name:   "error  id 0",
			Req:    &types.QueryOrderRequest{AppId: appID1, PairId: pair.Id},
			ExpErr: status.Error(codes.InvalidArgument, "id cannot be 0"),
		},
		{
			Name:   "error  id invalid",
			Req:    &types.QueryOrderRequest{AppId: appID1, PairId: pair.Id, Id: 123},
			ExpErr: status.Errorf(codes.NotFound, "order 1 in pair 123 not found"),
		},
		{
			Name:   "success",
			Req:    &types.QueryOrderRequest{AppId: appID1, PairId: pair.Id, Id: order.Id},
			ExpErr: nil,
		},
	}

	ctx := sdk.WrapSDKContext(s.ctx)
	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			resp, err := s.querier.Order(ctx, tc.Req)
			if tc.ExpErr != nil {
				s.Require().Error(err)
				s.Require().EqualError(err, tc.ExpErr.Error())
			} else {
				s.Require().NoError(err)
				s.Require().Equal(resp.Order.Id, order.Id)
				s.Require().Equal(resp.Order.PairId, order.PairId)
				s.Require().Equal(resp.Order.MsgHeight, order.MsgHeight)
				s.Require().Equal(resp.Order.Orderer, order.Orderer)
				s.Require().Equal(resp.Order.Direction, order.Direction)
				s.Require().Equal(resp.Order.OfferCoin, order.OfferCoin)
				s.Require().Equal(resp.Order.RemainingOfferCoin, order.RemainingOfferCoin)
				s.Require().Equal(resp.Order.ReceivedCoin, order.ReceivedCoin)
				s.Require().Equal(resp.Order.Price, order.Price)
				s.Require().Equal(resp.Order.Amount, order.Amount)
				s.Require().Equal(resp.Order.OpenAmount, order.OpenAmount)
				s.Require().Equal(resp.Order.BatchId, order.BatchId)
				s.Require().Equal(resp.Order.ExpireAt, order.ExpireAt)
				s.Require().Equal(resp.Order.Status, order.Status)
				s.Require().Equal(resp.Order.AppId, order.AppId)
			}
		})
	}
}

func (s *KeeperTestSuite) TestOrdersByOrderer() {
	addr1 := s.addr(1)

	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "uasset2", 1000000)

	pair := s.CreateNewLiquidityPair(appID1, addr1, asset1.Denom, asset2.Denom)
	_ = s.CreateNewLiquidityPool(appID1, pair.Id, addr1, "1000000000000uasset1,1000000000000uasset2")

	order := s.LimitOrder(appID1, addr1, pair.Id, types.OrderDirectionSell, newDec(1), newInt(100000), time.Second*10)

	testCases := []struct {
		Name   string
		Req    *types.QueryOrdersByOrdererRequest
		ExpErr error
	}{
		{
			Name:   "error empty request",
			Req:    nil,
			ExpErr: status.Error(codes.InvalidArgument, "empty request"),
		},
		{
			Name:   "error app id 0",
			Req:    &types.QueryOrdersByOrdererRequest{},
			ExpErr: status.Error(codes.InvalidArgument, "app id cannot be 0"),
		},
		{
			Name:   "error invalid orderer",
			Req:    &types.QueryOrdersByOrdererRequest{AppId: appID1, PairId: pair.Id, Orderer: "123"},
			ExpErr: status.Errorf(codes.InvalidArgument, "orderer address 123 is invalid"),
		},
		{
			Name:   "success only by orderer",
			Req:    &types.QueryOrdersByOrdererRequest{AppId: appID1, Orderer: addr1.String()},
			ExpErr: nil,
		},
		{
			Name:   "success orderer and pairID",
			Req:    &types.QueryOrdersByOrdererRequest{AppId: appID1, PairId: pair.Id, Orderer: addr1.String()},
			ExpErr: nil,
		},
	}

	ctx := sdk.WrapSDKContext(s.ctx)
	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			resp, err := s.querier.OrdersByOrderer(ctx, tc.Req)
			if tc.ExpErr != nil {
				s.Require().Error(err)
				s.Require().EqualError(err, tc.ExpErr.Error())
			} else {
				s.Require().NoError(err)
				s.Require().Equal(resp.Orders[0].Id, order.Id)
				s.Require().Equal(resp.Orders[0].PairId, order.PairId)
				s.Require().Equal(resp.Orders[0].MsgHeight, order.MsgHeight)
				s.Require().Equal(resp.Orders[0].Orderer, order.Orderer)
				s.Require().Equal(resp.Orders[0].Direction, order.Direction)
				s.Require().Equal(resp.Orders[0].OfferCoin, order.OfferCoin)
				s.Require().Equal(resp.Orders[0].RemainingOfferCoin, order.RemainingOfferCoin)
				s.Require().Equal(resp.Orders[0].ReceivedCoin, order.ReceivedCoin)
				s.Require().Equal(resp.Orders[0].Price, order.Price)
				s.Require().Equal(resp.Orders[0].Amount, order.Amount)
				s.Require().Equal(resp.Orders[0].OpenAmount, order.OpenAmount)
				s.Require().Equal(resp.Orders[0].BatchId, order.BatchId)
				s.Require().Equal(resp.Orders[0].ExpireAt, order.ExpireAt)
				s.Require().Equal(resp.Orders[0].Status, order.Status)
				s.Require().Equal(resp.Orders[0].AppId, order.AppId)
			}
		})
	}
}

func (s *KeeperTestSuite) TestFarmer() {
	addr1 := s.addr(1)

	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "uasset2", 1000000)

	pair := s.CreateNewLiquidityPair(appID1, addr1, asset1.Denom, asset2.Denom)
	pool := s.CreateNewLiquidityPool(appID1, pair.Id, addr1, "1000000000000uasset1,1000000000000uasset2")

	liquidityProvider1 := s.addr(2)
	s.Deposit(appID1, pool.Id, liquidityProvider1, "1000000000uasset1,1000000000uasset2")
	s.nextBlock()
	s.Require().True(utils.ParseCoins("10000000000pool1-1").Equal(s.getBalances(liquidityProvider1)))

	s.ctx = s.ctx.WithBlockTime(s.ctx.BlockTime())
	msg := types.NewMsgFarm(appID1, pool.Id, liquidityProvider1, utils.ParseCoin("5000000000pool1-1"))
	err := s.keeper.Farm(s.ctx, msg)
	s.Require().NoError(err)

	s.ctx = s.ctx.WithBlockTime(s.ctx.BlockTime().Add(types.DefaultFarmingQueueDuration).Add(time.Hour * 1))
	s.nextBlock()

	msg = types.NewMsgFarm(appID1, pool.Id, liquidityProvider1, utils.ParseCoin("5000000000pool1-1"))
	err = s.keeper.Farm(s.ctx, msg)
	s.Require().NoError(err)

	testCases := []struct {
		Name   string
		Req    *types.QueryFarmerRequest
		ExpErr error
	}{
		{
			Name:   "error empty request",
			Req:    nil,
			ExpErr: status.Error(codes.InvalidArgument, "empty request"),
		},
		{
			Name:   "error app id 0",
			Req:    &types.QueryFarmerRequest{},
			ExpErr: status.Error(codes.InvalidArgument, "app id cannot be 0"),
		},
		{
			Name:   "error invalid depositor",
			Req:    &types.QueryFarmerRequest{AppId: appID1, PoolId: pool.Id, Farmer: "123"},
			ExpErr: status.Errorf(codes.InvalidArgument, "farmer address 123 is invalid"),
		},
		{
			Name:   "error pool id invalid",
			Req:    &types.QueryFarmerRequest{AppId: appID1, PoolId: 123, Farmer: liquidityProvider1.String()},
			ExpErr: types.ErrInvalidPoolID,
		},
		{
			Name:   "success",
			Req:    &types.QueryFarmerRequest{AppId: appID1, PoolId: pool.Id, Farmer: liquidityProvider1.String()},
			ExpErr: nil,
		},
	}

	ctx := sdk.WrapSDKContext(s.ctx)
	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			resp, err := s.querier.Farmer(ctx, tc.Req)
			if tc.ExpErr != nil {
				s.Require().Error(err)
				s.Require().EqualError(err, tc.ExpErr.Error())
			} else {
				s.Require().NoError(err)
				s.Require().Equal(resp.ActivePoolCoin, utils.ParseCoin("5000000000pool1-1"))
				s.Require().Equal(resp.QueuedPoolCoin[0].PoolCoin, utils.ParseCoin("5000000000pool1-1"))
			}
		})
	}
}

func (s *KeeperTestSuite) TestDeserializePoolCoin() {
	addr1 := s.addr(1)

	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "uasset2", 1000000)

	pair := s.CreateNewLiquidityPair(appID1, addr1, asset1.Denom, asset2.Denom)
	pool := s.CreateNewLiquidityPool(appID1, pair.Id, addr1, "1000000000000uasset1,1000000000000uasset2")

	testCases := []struct {
		Name   string
		Req    *types.QueryDeserializePoolCoinRequest
		ExpErr error
	}{
		{
			Name:   "error empty request",
			Req:    nil,
			ExpErr: status.Error(codes.InvalidArgument, "empty request"),
		},
		{
			Name:   "error app id 0",
			Req:    &types.QueryDeserializePoolCoinRequest{},
			ExpErr: status.Error(codes.InvalidArgument, "app id cannot be 0"),
		},
		{
			Name:   "success",
			Req:    &types.QueryDeserializePoolCoinRequest{AppId: appID1, PoolId: pool.Id, PoolCoinAmount: 10000000000},
			ExpErr: nil,
		},
	}

	ctx := sdk.WrapSDKContext(s.ctx)
	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			resp, err := s.querier.DeserializePoolCoin(ctx, tc.Req)
			if tc.ExpErr != nil {
				s.Require().Error(err)
				s.Require().EqualError(err, tc.ExpErr.Error())
			} else {
				s.Require().NoError(err)
				s.Require().Equal(resp.Coins[0], utils.ParseCoin("1000000000uasset2"))
				s.Require().Equal(resp.Coins[1], utils.ParseCoin("1000000000uasset1"))
			}
		})
	}
}

func (s *KeeperTestSuite) TestPoolIncentives() {
	addr1 := s.addr(1)

	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "uasset2", 1000000)

	pair := s.CreateNewLiquidityPair(appID1, addr1, asset1.Denom, asset2.Denom)
	pool := s.CreateNewLiquidityPool(appID1, pair.Id, addr1, "1000000000000uasset1,1000000000000uasset2")

	testCases := []struct {
		Name   string
		Req    *types.QueryPoolsIncentivesRequest
		ExpErr error
	}{
		{
			Name:   "error empty request",
			Req:    nil,
			ExpErr: status.Error(codes.InvalidArgument, "empty request"),
		},
		{
			Name:   "error app id 0",
			Req:    &types.QueryPoolsIncentivesRequest{},
			ExpErr: status.Error(codes.InvalidArgument, "app id cannot be 0"),
		},
		{
			Name:   "success only by orderer",
			Req:    &types.QueryPoolsIncentivesRequest{AppId: appID1},
			ExpErr: nil,
		},
	}

	ctx := sdk.WrapSDKContext(s.ctx)
	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			resp, err := s.querier.PoolIncentives(ctx, tc.Req)
			if tc.ExpErr != nil {
				s.Require().Error(err)
				s.Require().EqualError(err, tc.ExpErr.Error())
			} else {
				s.Require().NoError(err)
				s.Require().Equal(resp.PoolIncentives[0].PoolId, pool.Id)
				s.Require().Equal(resp.PoolIncentives[0].MasterPool, false)
				s.Require().Equal(resp.PoolIncentives[0].TotalRewards, utils.ParseCoin("0ucmdx"))
				s.Require().Equal(resp.PoolIncentives[0].DistributedRewards, utils.ParseCoin("0ucmdx"))
				s.Require().Equal(resp.PoolIncentives[0].TotalEpochs, uint64(1))
				s.Require().Equal(resp.PoolIncentives[0].EpochDuration, time.Hour*24)
				s.Require().Equal(resp.PoolIncentives[0].IsSwapFee, true)
				s.Require().Equal(resp.PoolIncentives[0].AppId, appID1)
			}
		})
	}
}

func (s *KeeperTestSuite) TestFarmedPoolCoin() {
	addr1 := s.addr(1)

	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "uasset2", 1000000)

	pair := s.CreateNewLiquidityPair(appID1, addr1, asset1.Denom, asset2.Denom)
	pool := s.CreateNewLiquidityPool(appID1, pair.Id, addr1, "1000000000000uasset1,1000000000000uasset2")

	liquidityProvider1 := s.addr(2)
	s.Deposit(appID1, pool.Id, liquidityProvider1, "1000000000uasset1,1000000000uasset2")
	s.nextBlock()
	s.Require().True(utils.ParseCoins("10000000000pool1-1").Equal(s.getBalances(liquidityProvider1)))

	s.ctx = s.ctx.WithBlockTime(s.ctx.BlockTime())
	msg := types.NewMsgFarm(appID1, pool.Id, liquidityProvider1, utils.ParseCoin("5000000000pool1-1"))
	err := s.keeper.Farm(s.ctx, msg)
	s.Require().NoError(err)

	s.ctx = s.ctx.WithBlockTime(s.ctx.BlockTime().Add(types.DefaultFarmingQueueDuration).Add(time.Hour * 1))
	s.nextBlock()

	msg = types.NewMsgFarm(appID1, pool.Id, liquidityProvider1, utils.ParseCoin("5000000000pool1-1"))
	err = s.keeper.Farm(s.ctx, msg)
	s.Require().NoError(err)

	testCases := []struct {
		Name   string
		Req    *types.QueryFarmedPoolCoinRequest
		ExpErr error
	}{
		{
			Name:   "error empty request",
			Req:    nil,
			ExpErr: status.Error(codes.InvalidArgument, "empty request"),
		},
		{
			Name:   "error app id 0",
			Req:    &types.QueryFarmedPoolCoinRequest{},
			ExpErr: status.Error(codes.InvalidArgument, "app id cannot be 0"),
		},
		{
			Name:   "error pool id invalid",
			Req:    &types.QueryFarmedPoolCoinRequest{AppId: appID1, PoolId: 123},
			ExpErr: errorsmod.Wrapf(types.ErrInvalidPoolID, "pool id 123 is invalid"),
		},
		{
			Name:   "success only by orderer",
			Req:    &types.QueryFarmedPoolCoinRequest{AppId: appID1, PoolId: pool.Id},
			ExpErr: nil,
		},
	}

	ctx := sdk.WrapSDKContext(s.ctx)
	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			resp, err := s.querier.FarmedPoolCoin(ctx, tc.Req)
			if tc.ExpErr != nil {
				s.Require().Error(err)
				s.Require().EqualError(err, tc.ExpErr.Error())
			} else {
				s.Require().NoError(err)
				s.Require().Equal(resp.Coin, utils.ParseCoin("10000000000pool1-1"))
			}
		})
	}
}
