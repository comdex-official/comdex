package types

import (
	"strconv"
	"strings"

	sdkmath "cosmossdk.io/math"

	"github.com/cometbft/cometbft/crypto"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"

	"github.com/comdex-official/comdex/x/liquidity/amm"
	"github.com/comdex-official/comdex/x/liquidity/expected"
)

type sendCoinsTxKey struct {
	from, to string
}

type sendCoinsTx struct {
	from, to sdk.AccAddress
	amt      sdk.Coins
}

// BulkSendCoinsOperation holds a list of SendCoins operations for bulk execution.
type BulkSendCoinsOperation struct {
	txSet map[sendCoinsTxKey]*sendCoinsTx
	txs   []*sendCoinsTx
}

// NewBulkSendCoinsOperation returns an empty BulkSendCoinsOperation.
func NewBulkSendCoinsOperation() *BulkSendCoinsOperation {
	return &BulkSendCoinsOperation{
		txSet: map[sendCoinsTxKey]*sendCoinsTx{},
	}
}

// QueueSendCoins queues a BankKeeper.SendCoins operation for later execution.
func (op *BulkSendCoinsOperation) QueueSendCoins(
	fromAddr, toAddr sdk.AccAddress,
	amt sdk.Coins,
) {
	if amt.IsValid() && !amt.IsZero() {
		txKey := sendCoinsTxKey{fromAddr.String(), toAddr.String()}
		tx, ok := op.txSet[txKey]
		if !ok {
			tx = &sendCoinsTx{fromAddr, toAddr, sdk.Coins{}}
			op.txSet[txKey] = tx
			op.txs = append(op.txs, tx)
		}
		tx.amt = tx.amt.Add(amt...)
	}
}

// Run runs BankKeeper.InputOutputCoins once for queued operations.
func (op *BulkSendCoinsOperation) Run(ctx sdk.Context, bankKeeper expected.BankKeeper) error {
	if len(op.txs) > 0 {
		for _, tx := range op.txs {
			if err := bankKeeper.SendCoins(ctx, tx.from, tx.to, tx.amt); err != nil {
				return err
			}
		}
		return nil
	}
	return nil
}

// NewPoolResponse returns a new PoolResponse from given information.
func NewPoolResponse(pool Pool, rx, ry sdk.Coin, poolCoinSupply sdkmath.Int) PoolResponse {
	var price *sdkmath.LegacyDec
	if !pool.Disabled {
		p := pool.AMMPool(rx.Amount, ry.Amount, sdkmath.Int{}).Price()
		price = &p
	}
	return PoolResponse{
		Id:             pool.Id,
		PairId:         pool.PairId,
		ReserveAddress: pool.ReserveAddress,
		PoolCoinDenom:  pool.PoolCoinDenom,
		Balances: PoolBalances{
			BaseCoin:  ry,
			QuoteCoin: rx,
		},
		LastDepositRequestId:  pool.LastDepositRequestId,
		LastWithdrawRequestId: pool.LastWithdrawRequestId,
		AppId:                 pool.AppId,
		Type:                  pool.Type,
		Creator:               pool.Creator,
		PoolCoinSupply:        poolCoinSupply,
		MinPrice:              pool.MinPrice,
		MaxPrice:              pool.MaxPrice,
		Price:                 price,
		Disabled:              pool.Disabled,
	}
}

// IsTooSmallOrderAmount returns whether the order amount is too small for
// matching, based on the order price.
func IsTooSmallOrderAmount(amt sdkmath.Int, price sdkmath.LegacyDec) bool {
	return amt.LT(amm.MinCoinAmount) || price.MulInt(amt).LT(amm.MinCoinAmount.ToLegacyDec())
}

// PriceLimits returns the lowest and the highest price limits with given last price
// and price limit ratio.
func PriceLimits(lastPrice, priceLimitRatio sdkmath.LegacyDec, tickPrec int) (lowestPrice, highestPrice sdkmath.LegacyDec) {
	lowestPrice = amm.PriceToUpTick(lastPrice.Mul(sdkmath.LegacyOneDec().Sub(priceLimitRatio)), tickPrec)
	highestPrice = amm.PriceToDownTick(lastPrice.Mul(sdkmath.LegacyOneDec().Add(priceLimitRatio)), tickPrec)
	return
}

func NewMMOrderIndex(orderer sdk.AccAddress, appID, pairID uint64, orderIds []uint64) MMOrderIndex {
	return MMOrderIndex{
		Orderer:  orderer.String(),
		AppId:    appID,
		PairId:   pairID,
		OrderIds: orderIds,
	}
}

func (index MMOrderIndex) GetOrderer() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(index.Orderer)
	if err != nil {
		panic(err)
	}
	return addr
}

// MMOrderTick holds information about each tick's price and amount of an MMOrder.
type MMOrderTick struct {
	OfferCoinAmount sdkmath.Int
	Price           sdkmath.LegacyDec
	Amount          sdkmath.Int
}

// MMOrderTicks returns fairly distributed tick information with given parameters.
func MMOrderTicks(dir OrderDirection, minPrice, maxPrice sdkmath.LegacyDec, amt sdkmath.Int, maxNumTicks, tickPrec int) (ticks []MMOrderTick) {
	ammDir := amm.OrderDirection(dir)
	if minPrice.Equal(maxPrice) {
		return []MMOrderTick{{OfferCoinAmount: amm.OfferCoinAmount(ammDir, minPrice, amt), Price: minPrice, Amount: amt}}
	}
	gap := maxPrice.Sub(minPrice).QuoInt64(int64(maxNumTicks - 1))
	switch dir {
	case OrderDirectionBuy:
		var prevP sdkmath.LegacyDec
		for i := 0; i < maxNumTicks-1; i++ {
			p := amm.PriceToDownTick(minPrice.Add(gap.MulInt64(int64(i))), tickPrec)
			if prevP.IsNil() || !p.Equal(prevP) {
				ticks = append(ticks, MMOrderTick{
					Price: p,
				})
				prevP = p
			}
		}
		tickAmt := amt.QuoRaw(int64(len(ticks) + 1))
		for i := range ticks {
			ticks[i].Amount = tickAmt
			ticks[i].OfferCoinAmount = amm.OfferCoinAmount(ammDir, ticks[i].Price, ticks[i].Amount)
			amt = amt.Sub(tickAmt)
		}
		ticks = append(ticks, MMOrderTick{
			OfferCoinAmount: amm.OfferCoinAmount(ammDir, maxPrice, amt),
			Price:           maxPrice,
			Amount:          amt,
		})
	case OrderDirectionSell:
		var prevP sdkmath.LegacyDec
		for i := 0; i < maxNumTicks-1; i++ {
			p := amm.PriceToUpTick(maxPrice.Sub(gap.MulInt64(int64(i))), tickPrec)
			if prevP.IsNil() || !p.Equal(prevP) {
				ticks = append(ticks, MMOrderTick{
					Price: p,
				})
				prevP = p
			}
		}
		tickAmt := amt.QuoRaw(int64(len(ticks) + 1))
		for i := range ticks {
			ticks[i].Amount = tickAmt
			ticks[i].OfferCoinAmount = amm.OfferCoinAmount(ammDir, ticks[i].Price, ticks[i].Amount)
			amt = amt.Sub(tickAmt)
		}
		ticks = append(ticks, MMOrderTick{
			OfferCoinAmount: amm.OfferCoinAmount(ammDir, minPrice, amt),
			Price:           minPrice,
			Amount:          amt,
		})
	}
	return
}

// FormatUint64s returns comma-separated string representation of
// a slice of uint64.
func FormatUint64s(us []uint64) (s string) {
	ss := make([]string, 0, len(us))
	for _, u := range us {
		ss = append(ss, strconv.FormatUint(u, 10))
	}
	return strings.Join(ss, ",")
}

// DeriveAddress derives an address with the given address length type, module name, and
// address derivation name. It is used to derive private plan farming pool address, and staking reserve address.
func DeriveAddress(addressType AddressType, moduleName, name string) sdk.AccAddress {
	switch addressType {
	case AddressType32Bytes:
		return address.Module(moduleName, []byte(name))
	case AddressType20Bytes:
		return sdk.AccAddress(crypto.AddressHash([]byte(moduleName + name)))
	default:
		return sdk.AccAddress{}
	}
}

// ItemExists returns true if item exists in array else false .
func ItemExists(array []string, item string) bool {
	for _, v := range array {
		if v == item {
			return true
		}
	}
	return false
}

// BuildUndirectedGraph builds undirected the graph from the given edges represented as adjacency list .
func BuildUndirectedGraph(edges [][]string) (graph map[string][]string) {
	graph = make(map[string][]string)

	// Loop to iterate over every edge of the graph
	for _, edge := range edges {
		a, b := edge[0], edge[1]

		// Creating the graph as adjacency list
		graph[a] = append(graph[a], b)
		graph[b] = append(graph[b], a)
	}
	return graph
}

// BfsShortestPath returns the shortest path between two nodes in undirected graph.
func BfsShortestPath(undirectedGraph map[string][]string, start string, goal string) ([]string, bool) {
	var explored []string

	// Queue for traversing the graph in the BFS
	queue := [][]string{{start}}

	// If the desired node is reached
	if start == goal {
		return []string{start}, true
	}

	// Loop to traverse the graph with the help of the queue
	for {
		if len(queue) == 0 {
			// empty queue, hence break the loop
			break
		}
		path := queue[0]

		// dequeue opearation
		queue = queue[1:]

		node := path[len(path)-1]

		// Condition to check if the current node is not visited
		if !ItemExists(explored, node) {
			neighbours := undirectedGraph[node]

			// Loop to iterate over the neighbours of the node
			for _, neighbour := range neighbours {
				newPath := path
				newPath = append(newPath, neighbour)

				// enqueue operation
				queue = append(queue, newPath)

				// Condition to check if the neighbour node is the goal
				if neighbour == goal {
					// path found
					return newPath, true
				}
			}
			explored = append(explored, node)
		}
	}
	// return false, if no paths exists between start -> goal
	return []string{}, false
}
