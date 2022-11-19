package types_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"

	utils "github.com/petrichormoney/petri/types"
	"github.com/petrichormoney/petri/x/liquidity/types"
)

var testAddr = sdk.AccAddress(crypto.AddressHash([]byte("test")))

func TestMsgCreatePair(t *testing.T) {
	for _, tc := range []struct {
		name        string
		malleate    func(msg *types.MsgCreatePair)
		expectedErr string
	}{
		{
			"happy case",
			func(msg *types.MsgCreatePair) {},
			"",
		},
		{
			"invalid creator",
			func(msg *types.MsgCreatePair) {
				msg.Creator = "invalidaddr"
			},
			"invalid creator address: decoding bech32 failed: invalid separator index -1: invalid address",
		},
		{
			"invalid base coin denom",
			func(msg *types.MsgCreatePair) {
				msg.BaseCoinDenom = "invaliddenom!"
			},
			"invalid denom: invaliddenom!: invalid request",
		},
		{
			"invalid quote coin denom",
			func(msg *types.MsgCreatePair) {
				msg.QuoteCoinDenom = "invaliddenom!"
			},
			"invalid denom: invaliddenom!: invalid request",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			msg := types.NewMsgCreatePair(1, testAddr, "denom1", "denom2")
			tc.malleate(msg)
			require.Equal(t, types.TypeMsgCreatePair, msg.Type())
			require.Equal(t, types.RouterKey, msg.Route())
			err := msg.ValidateBasic()
			if tc.expectedErr == "" {
				require.NoError(t, err)
				signers := msg.GetSigners()
				require.Len(t, signers, 1)
				require.Equal(t, msg.GetCreator(), signers[0])
			} else {
				require.EqualError(t, err, tc.expectedErr)
			}
		})
	}
}

func TestMsgCreatePool(t *testing.T) {
	for _, tc := range []struct {
		name        string
		malleate    func(msg *types.MsgCreatePool)
		expectedErr string
	}{
		{
			"happy case",
			func(msg *types.MsgCreatePool) {},
			"", // empty means no error expected
		},
		{
			"invalid pair id",
			func(msg *types.MsgCreatePool) {
				msg.PairId = 0
			},
			"pair id must not be 0: invalid request",
		},
		{
			"invalid creator",
			func(msg *types.MsgCreatePool) {
				msg.Creator = "invalidaddr"
			},
			"invalid creator address: decoding bech32 failed: invalid separator index -1: invalid address",
		},
		{
			"invalid deposit coins",
			func(msg *types.MsgCreatePool) {
				msg.DepositCoins = sdk.Coins{utils.ParseCoin("0denom1"), utils.ParseCoin("1000000denom2")}
			},
			"coin 0denom1 amount is not positive",
		},
		{
			"invalid deposit coins",
			func(msg *types.MsgCreatePool) {
				msg.DepositCoins = sdk.Coins{utils.ParseCoin("1000000denom1"), utils.ParseCoin("0denom2")}
			},
			"coin denom2 amount is not positive",
		},
		{
			"invalid deposit coins",
			func(msg *types.MsgCreatePool) {
				msg.DepositCoins = utils.ParseCoins("1000000denom1,1000000denom2,1000000denom3")
			},
			"wrong number of deposit coins: 3: invalid request",
		},
		{
			"too large deposit coins",
			func(msg *types.MsgCreatePool) {
				msg.DepositCoins = utils.ParseCoins("100000000000000000000000000000000000000000denom1,100000000000000000000000000000000000000000denom2")
			},
			"deposit coin 100000000000000000000000000000000000000000denom1 is bigger than the max amount 10000000000000000000000000000000000000000: invalid request",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			msg := types.NewMsgCreatePool(1, testAddr, 1, utils.ParseCoins("1000000denom1,1000000denom2"))
			tc.malleate(msg)
			require.Equal(t, types.TypeMsgCreatePool, msg.Type())
			require.Equal(t, types.RouterKey, msg.Route())
			err := msg.ValidateBasic()
			if tc.expectedErr == "" {
				require.NoError(t, err)
				signers := msg.GetSigners()
				require.Len(t, signers, 1)
				require.Equal(t, msg.GetCreator(), signers[0])
			} else {
				require.EqualError(t, err, tc.expectedErr)
			}
		})
	}
}

func TestMsgDeposit(t *testing.T) {
	testCases := []struct {
		name        string
		malleate    func(msg *types.MsgDeposit)
		expectedErr string
	}{
		{
			"happy case",
			func(msg *types.MsgDeposit) {},
			"", // empty means no error expected
		},
		{
			"invalid depositor",
			func(msg *types.MsgDeposit) {
				msg.Depositor = "invalidaddr"
			},
			"invalid depositor address: decoding bech32 failed: invalid separator index -1: invalid address",
		},
		{
			"invalid pool id",
			func(msg *types.MsgDeposit) {
				msg.PoolId = 0
			},
			"pool id must not be 0: invalid request",
		},
		{
			"invalid deposit coins",
			func(msg *types.MsgDeposit) {
				msg.DepositCoins = sdk.Coins{utils.ParseCoin("0denom1"), utils.ParseCoin("1000000denom2")}
			},
			"coin 0denom1 amount is not positive",
		},
		{
			"invalid deposit coins",
			func(msg *types.MsgDeposit) {
				msg.DepositCoins = sdk.Coins{utils.ParseCoin("1000000denom1"), utils.ParseCoin("0denom2")}
			},
			"coin denom2 amount is not positive",
		},
		{
			"invalid deposit coins",
			func(msg *types.MsgDeposit) {
				msg.DepositCoins = utils.ParseCoins("1000000denom1,1000000denom2,1000000denom3")
			},
			"wrong number of deposit coins: 3: invalid request",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			msg := types.NewMsgDeposit(1, testAddr, 1, utils.ParseCoins("1000000denom1,1000000denom2"))
			tc.malleate(msg)
			require.Equal(t, types.TypeMsgDeposit, msg.Type())
			require.Equal(t, types.RouterKey, msg.Route())
			err := msg.ValidateBasic()
			if tc.expectedErr == "" {
				require.NoError(t, err)
				signers := msg.GetSigners()
				require.Len(t, signers, 1)
				require.Equal(t, msg.GetDepositor(), signers[0])
			} else {
				require.EqualError(t, err, tc.expectedErr)
			}
		})
	}
}

func TestMsgWithdraw(t *testing.T) {
	for _, tc := range []struct {
		name        string
		malleate    func(msg *types.MsgWithdraw)
		expectedErr string
	}{
		{
			"happy case",
			func(msg *types.MsgWithdraw) {},
			"", // empty means no error expected
		},
		{
			"invalid withdrawer",
			func(msg *types.MsgWithdraw) {
				msg.Withdrawer = "invalidaddr"
			},
			"invalid withdrawer address: decoding bech32 failed: invalid separator index -1: invalid address",
		},
		{
			"invalid pool id",
			func(msg *types.MsgWithdraw) {
				msg.PoolId = 0
			},
			"pool id must not be 0: invalid request",
		},
		{
			"invalid pool coin",
			func(msg *types.MsgWithdraw) {
				msg.PoolCoin = utils.ParseCoin("0pool1")
			},
			"pool coin must be positive: invalid request",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			msg := types.NewMsgWithdraw(1, testAddr, 1, utils.ParseCoin("1000000pool1"))
			tc.malleate(msg)
			require.Equal(t, types.TypeMsgWithdraw, msg.Type())
			require.Equal(t, types.RouterKey, msg.Route())
			err := msg.ValidateBasic()
			if tc.expectedErr == "" {
				require.NoError(t, err)
				signers := msg.GetSigners()
				require.Len(t, signers, 1)
				require.Equal(t, msg.GetWithdrawer(), signers[0])
			} else {
				require.EqualError(t, err, tc.expectedErr)
			}
		})
	}
}

func TestMsgLimitOrder(t *testing.T) {
	orderLifespan := 20 * time.Second
	for _, tc := range []struct {
		name        string
		malleate    func(msg *types.MsgLimitOrder)
		expectedErr string
	}{
		{
			"happy case",
			func(msg *types.MsgLimitOrder) {},
			"", // empty means no error expected
		},
		{
			"invalid orderer",
			func(msg *types.MsgLimitOrder) {
				msg.Orderer = "invalidaddr"
			},
			"invalid orderer address: decoding bech32 failed: invalid separator index -1: invalid address",
		},
		{
			"invalid pair id",
			func(msg *types.MsgLimitOrder) {
				msg.PairId = 0
			},
			"pair id must not be 0: invalid request",
		},
		{
			"invalid direction",
			func(msg *types.MsgLimitOrder) {
				msg.Direction = 0
			},
			"invalid order direction: ORDER_DIRECTION_UNSPECIFIED: invalid request",
		},
		{
			"invalid offer coin",
			func(msg *types.MsgLimitOrder) {
				msg.OfferCoin = utils.ParseCoin("0denom1")
			},
			"offer coin 0denom1 is smaller than the min amount 100: invalid request",
		},
		{
			"small offer coin amount",
			func(msg *types.MsgLimitOrder) {
				msg.OfferCoin = utils.ParseCoin("10denom1")
			},
			"offer coin 10denom1 is smaller than the min amount 100: invalid request",
		},
		{
			"insufficient offer coin amount",
			func(msg *types.MsgLimitOrder) {
				msg.OfferCoin = utils.ParseCoin("1000000denom2")
				msg.Price = utils.ParseDec("10")
				msg.Amount = sdk.NewInt(1000000)
			},
			"1000000denom2 is less than 10000000denom2: insufficient offer coin",
		},
		{
			"invalid demand coin denom",
			func(msg *types.MsgLimitOrder) {
				msg.DemandCoinDenom = "invaliddenom!"
			},
			"invalid demand coin denom: invalid denom: invaliddenom!",
		},
		{
			"same offer coin denom and demand coin denom",
			func(msg *types.MsgLimitOrder) {
				msg.OfferCoin = utils.ParseCoin("1000000denom1")
				msg.DemandCoinDenom = "denom1"
			},
			"offer coin denom and demand coin denom must not be same: invalid request",
		},
		{
			"invalid price",
			func(msg *types.MsgLimitOrder) {
				msg.Price = utils.ParseDec("0")
			},
			"price must be positive: invalid request",
		},
		{
			"zero order amount",
			func(msg *types.MsgLimitOrder) {
				msg.Amount = sdk.ZeroInt()
			},
			"order amount 0 is smaller than the min amount 100: invalid request",
		},
		{
			"small order amount",
			func(msg *types.MsgLimitOrder) {
				msg.Amount = sdk.NewInt(10)
			},
			"order amount 10 is smaller than the min amount 100: invalid request",
		},
		{
			"invalid order lifespan",
			func(msg *types.MsgLimitOrder) {
				msg.OrderLifespan = -1
			},
			"order lifespan must not be negative: -1ns: invalid request",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			msg := types.NewMsgLimitOrder(
				1, testAddr, 1, types.OrderDirectionBuy, utils.ParseCoin("1000000denom2"),
				"denom1", utils.ParseDec("1.0"), sdk.NewInt(1000000), orderLifespan)
			tc.malleate(msg)
			require.Equal(t, types.TypeMsgLimitOrder, msg.Type())
			require.Equal(t, types.RouterKey, msg.Route())
			err := msg.ValidateBasic()
			if tc.expectedErr == "" {
				require.NoError(t, err)
				signers := msg.GetSigners()
				require.Len(t, signers, 1)
				require.Equal(t, msg.GetOrderer(), signers[0])
			} else {
				require.EqualError(t, err, tc.expectedErr)
			}
		})
	}
}

func TestMsgMarketOrder(t *testing.T) {
	orderLifespan := 20 * time.Second
	for _, tc := range []struct {
		name        string
		malleate    func(msg *types.MsgMarketOrder)
		expectedErr string
	}{
		{
			"happy case",
			func(msg *types.MsgMarketOrder) {},
			"", // empty means no error expected
		},
		{
			"invalid orderer",
			func(msg *types.MsgMarketOrder) {
				msg.Orderer = "invalidaddr"
			},
			"invalid orderer address: decoding bech32 failed: invalid separator index -1: invalid address",
		},
		{
			"invalid pair id",
			func(msg *types.MsgMarketOrder) {
				msg.PairId = 0
			},
			"pair id must not be 0: invalid request",
		},
		{
			"invalid direction",
			func(msg *types.MsgMarketOrder) {
				msg.Direction = 0
			},
			"invalid order direction: ORDER_DIRECTION_UNSPECIFIED: invalid request",
		},
		{
			"zero offer coin",
			func(msg *types.MsgMarketOrder) {
				msg.OfferCoin = utils.ParseCoin("0denom1")
			},
			"offer coin 0denom1 is smaller than the min amount 100: invalid request",
		},
		{
			"small offer coin amount",
			func(msg *types.MsgMarketOrder) {
				msg.OfferCoin = utils.ParseCoin("10denom1")
			},
			"offer coin 10denom1 is smaller than the min amount 100: invalid request",
		},
		{
			"invalid demand coin denom",
			func(msg *types.MsgMarketOrder) {
				msg.DemandCoinDenom = "invaliddenom!"
			},
			"invalid demand coin denom: invalid denom: invaliddenom!",
		},
		{
			"same offer coin denom and demand coin denom",
			func(msg *types.MsgMarketOrder) {
				msg.OfferCoin = utils.ParseCoin("1000000denom1")
				msg.DemandCoinDenom = "denom1"
			},
			"offer coin denom and demand coin denom must not be same: invalid request",
		},
		{
			"zero order amount",
			func(msg *types.MsgMarketOrder) {
				msg.Amount = sdk.ZeroInt()
			},
			"order amount 0 is smaller than the min amount 100: invalid request",
		},
		{
			"small order amount",
			func(msg *types.MsgMarketOrder) {
				msg.Amount = sdk.NewInt(10)
			},
			"order amount 10 is smaller than the min amount 100: invalid request",
		},
		{
			"invalid order lifespan",
			func(msg *types.MsgMarketOrder) {
				msg.OrderLifespan = -1
			},
			"order lifespan must not be negative: -1ns: invalid request",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			msg := types.NewMsgMarketOrder(
				1, testAddr, 1, types.OrderDirectionBuy, utils.ParseCoin("1000000denom1"),
				"denom2", sdk.NewInt(1000000), orderLifespan)
			tc.malleate(msg)
			require.Equal(t, types.TypeMsgMarketOrder, msg.Type())
			require.Equal(t, types.RouterKey, msg.Route())
			err := msg.ValidateBasic()
			if tc.expectedErr == "" {
				require.NoError(t, err)
				signers := msg.GetSigners()
				require.Len(t, signers, 1)
				require.Equal(t, msg.GetOrderer(), signers[0])
			} else {
				require.EqualError(t, err, tc.expectedErr)
			}
		})
	}
}

func TestMsgCancelOrder(t *testing.T) {
	for _, tc := range []struct {
		name        string
		malleate    func(msg *types.MsgCancelOrder)
		expectedErr string
	}{
		{
			"happy case",
			func(msg *types.MsgCancelOrder) {},
			"", // empty means no error expected
		},
		{
			"invalid orderer",
			func(msg *types.MsgCancelOrder) {
				msg.Orderer = "invalidaddr"
			},
			"invalid orderer address: decoding bech32 failed: invalid separator index -1: invalid address",
		},
		{
			"invalid pair id",
			func(msg *types.MsgCancelOrder) {
				msg.PairId = 0
			},
			"pair id must not be 0: invalid request",
		},
		{
			"invalid order id",
			func(msg *types.MsgCancelOrder) {
				msg.OrderId = 0
			},
			"order id must not be 0: invalid request",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			msg := types.NewMsgCancelOrder(1, testAddr, 1, 1)
			tc.malleate(msg)
			require.Equal(t, types.TypeMsgCancelOrder, msg.Type())
			require.Equal(t, types.RouterKey, msg.Route())
			err := msg.ValidateBasic()
			if tc.expectedErr == "" {
				require.NoError(t, err)
				signers := msg.GetSigners()
				require.Len(t, signers, 1)
				require.Equal(t, msg.GetOrderer(), signers[0])
			} else {
				require.EqualError(t, err, tc.expectedErr)
			}
		})
	}
}

func TestMsgCancelAllOrders(t *testing.T) {
	for _, tc := range []struct {
		name        string
		malleate    func(msg *types.MsgCancelAllOrders)
		expectedErr string
	}{
		{
			"happy case",
			func(msg *types.MsgCancelAllOrders) {},
			"",
		},
		{
			"invalid orderer",
			func(msg *types.MsgCancelAllOrders) {
				msg.Orderer = "invalidaddr"
			},
			"invalid orderer address: decoding bech32 failed: invalid separator index -1: invalid address",
		},
		{
			"invalid pair ids",
			func(msg *types.MsgCancelAllOrders) {
				msg.PairIds = []uint64{0}
			},
			"pair id must not be 0: invalid request",
		},
		{
			"invalid pair ids",
			func(msg *types.MsgCancelAllOrders) {
				msg.PairIds = []uint64{1, 1}
			},
			"duplicate pair id presents in the pair id list",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			msg := types.NewMsgCancelAllOrders(1, sdk.AccAddress(crypto.AddressHash([]byte("orderer"))), []uint64{1, 2, 3})
			tc.malleate(msg)
			require.Equal(t, types.TypeMsgCancelAllOrders, msg.Type())
			require.Equal(t, types.RouterKey, msg.Route())
			err := msg.ValidateBasic()
			if tc.expectedErr == "" {
				require.NoError(t, err)
				signers := msg.GetSigners()
				require.Len(t, signers, 1)
				require.Equal(t, msg.GetOrderer(), signers[0])
			} else {
				require.EqualError(t, err, tc.expectedErr)
			}
		})
	}
}

func TestMsgFarm(t *testing.T) {
	for _, tc := range []struct {
		name        string
		malleate    func(msg *types.MsgFarm)
		expectedErr string
	}{
		{
			"happy case",
			func(msg *types.MsgFarm) {},
			"",
		},
		{
			"invalid farmer",
			func(msg *types.MsgFarm) {
				msg.Farmer = "invalidaddr"
			},
			"invalid withdrawer address: decoding bech32 failed: invalid separator index -1: invalid address",
		},
		{
			"invalid pool id",
			func(msg *types.MsgFarm) {
				msg.PoolId = 0
			},
			"pool id must not be 0: invalid request",
		},
		{
			"invalid pool coin",
			func(msg *types.MsgFarm) {
				msg.FarmingPoolCoin = utils.ParseCoin("0pool1-1")
			},
			"coin must be positive: invalid request",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			msg := types.NewMsgFarm(1, 1, sdk.AccAddress(crypto.AddressHash([]byte("orderer"))), utils.ParseCoin("123pool1-1"))
			tc.malleate(msg)
			require.Equal(t, types.TypeMsgFarm, msg.Type())
			require.Equal(t, types.RouterKey, msg.Route())
			err := msg.ValidateBasic()
			if tc.expectedErr == "" {
				require.NoError(t, err)
				signers := msg.GetSigners()
				require.Len(t, signers, 1)
				require.Equal(t, msg.GetFarmer(), signers[0])
			} else {
				require.EqualError(t, err, tc.expectedErr)
			}
		})
	}
}

func TestMsgUnfarm(t *testing.T) {
	for _, tc := range []struct {
		name        string
		malleate    func(msg *types.MsgUnfarm)
		expectedErr string
	}{
		{
			"happy case",
			func(msg *types.MsgUnfarm) {},
			"",
		},
		{
			"invalid farmer",
			func(msg *types.MsgUnfarm) {
				msg.Farmer = "invalidaddr"
			},
			"invalid withdrawer address: decoding bech32 failed: invalid separator index -1: invalid address",
		},
		{
			"invalid pool id",
			func(msg *types.MsgUnfarm) {
				msg.PoolId = 0
			},
			"pool id must not be 0: invalid request",
		},
		{
			"invalid pool coin",
			func(msg *types.MsgUnfarm) {
				msg.UnfarmingPoolCoin = utils.ParseCoin("0pool1-1")
			},
			"coin must be positive: invalid request",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			msg := types.NewMsgUnfarm(1, 1, sdk.AccAddress(crypto.AddressHash([]byte("orderer"))), utils.ParseCoin("123pool1-1"))
			tc.malleate(msg)
			require.Equal(t, types.TypeMsgUnfarm, msg.Type())
			require.Equal(t, types.RouterKey, msg.Route())
			err := msg.ValidateBasic()
			if tc.expectedErr == "" {
				require.NoError(t, err)
				signers := msg.GetSigners()
				require.Len(t, signers, 1)
				require.Equal(t, msg.GetFarmer(), signers[0])
			} else {
				require.EqualError(t, err, tc.expectedErr)
			}
		})
	}
}
