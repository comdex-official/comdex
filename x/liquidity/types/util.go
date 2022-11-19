package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/tendermint/tendermint/crypto"

	"github.com/petrichormoney/petri/x/liquidity/amm"
	"github.com/petrichormoney/petri/x/liquidity/expected"
)

// BulkSendCoinsOperation holds a list of SendCoins operations for bulk execution.
type BulkSendCoinsOperation struct {
	Inputs  []banktypes.Input
	Outputs []banktypes.Output
}

// NewBulkSendCoinsOperation returns an empty BulkSendCoinsOperation.
func NewBulkSendCoinsOperation() *BulkSendCoinsOperation {
	return &BulkSendCoinsOperation{
		Inputs:  []banktypes.Input{},
		Outputs: []banktypes.Output{},
	}
}

// QueueSendCoins queues a BankKeeper.SendCoins operation for later execution.
func (op *BulkSendCoinsOperation) QueueSendCoins(
	fromAddr, toAddr sdk.AccAddress,
	amt sdk.Coins,
) {
	if amt.IsValid() && !amt.IsZero() {
		op.Inputs = append(op.Inputs, banktypes.NewInput(fromAddr, amt))
		op.Outputs = append(op.Outputs, banktypes.NewOutput(toAddr, amt))
	}
}

// Run runs BankKeeper.InputOutputCoins once for queued operations.
func (op *BulkSendCoinsOperation) Run(ctx sdk.Context, bankKeeper expected.BankKeeper) error {
	if len(op.Inputs) > 0 && len(op.Outputs) > 0 {
		return bankKeeper.InputOutputCoins(ctx, op.Inputs, op.Outputs)
	}
	return nil
}

// IsTooSmallOrderAmount returns whether the order amount is too small for
// matching, based on the order price.
func IsTooSmallOrderAmount(amt sdk.Int, price sdk.Dec) bool {
	return amt.LT(amm.MinCoinAmount) || price.MulInt(amt).LT(amm.MinCoinAmount.ToDec())
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
