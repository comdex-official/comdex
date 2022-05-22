package keeper

import (
	"time"

	assettypes "github.com/comdex-official/comdex/x/asset/types"
	auctiontypes "github.com/comdex-official/comdex/x/auction/types"
	liquidationtypes "github.com/comdex-official/comdex/x/liquidation/types"
	vaulttypes "github.com/comdex-official/comdex/x/vault/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	protobuftypes "github.com/gogo/protobuf/types"
)

func (k *Keeper) GetCollateralAuctionID(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.CollateralAuctionIdKey
		value = store.Get(key)
	)
	if value == nil {
		return 0
	}
	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue()
}

func (k *Keeper) SetCollateralAuctionID(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.CollateralAuctionIdKey
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)

	store.Set(key, value)
}

func (k *Keeper) GetDebtAuctionID(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.DebtAuctionIdKey
		value = store.Get(key)
	)
	if value == nil {
		return 0
	}
	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue()
}

func (k *Keeper) SetDebtAuctionID(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.DebtAuctionIdKey
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)

	store.Set(key, value)
}

func (k *Keeper) GetDutchAuctionID(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.DutchAuctionIdKey
		value = store.Get(key)
	)
	if value == nil {
		return 0
	}
	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue()
}

func (k *Keeper) SetDutchAuctionID(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.DutchAuctionIdKey
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)

	store.Set(key, value)
}

func (k *Keeper) SetCollateralAuction(ctx sdk.Context, auction auctiontypes.CollateralAuction) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.CollateralAuctionKey(auction.Id)
		value = k.cdc.MustMarshal(&auction)
	)
	store.Set(key, value)
}

func (k *Keeper) DeleteCollateralAuction(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.CollateralAuctionKey(id)
	)
	store.Delete(key)
}

func (k *Keeper) SetDebtAuction(ctx sdk.Context, auction auctiontypes.DebtAuction) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.DebtAuctionKey(auction.AuctionId)
		value = k.cdc.MustMarshal(&auction)
	)
	store.Set(key, value)
}

func (k *Keeper) DeleteDebtAuction(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.DebtAuctionKey(id)
	)
	store.Delete(key)
}

func (k *Keeper) SetDutchAuction(ctx sdk.Context, auction auctiontypes.DutchAuction) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.DutchAuctionKey(auction.AuctionId)
		value = k.cdc.MustMarshal(&auction)
	)
	store.Set(key, value)
}

func (k *Keeper) DeleteDutchAuction(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.DutchAuctionKey(id)
	)
	store.Delete(key)
}

func (k *Keeper) GetCollateralAuction(ctx sdk.Context, id uint64) (auction auctiontypes.CollateralAuction, found bool) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.CollateralAuctionKey(id)
		value = store.Get(key)
	)

	if value == nil {
		return auction, false
	}

	k.cdc.MustUnmarshal(value, &auction)
	return auction, true
}

func (k *Keeper) GetCollateralAuctions(ctx sdk.Context) (auctions []auctiontypes.CollateralAuction) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, auctiontypes.CollateralAuctionKeyPrefix)
	)

	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var auction auctiontypes.CollateralAuction
		k.cdc.MustUnmarshal(iter.Value(), &auction)
		auctions = append(auctions, auction)
	}

	return auctions
}

func (k *Keeper) GetDebtAuction(ctx sdk.Context, id uint64) (auction auctiontypes.DebtAuction, found bool) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.DebtAuctionKey(id)
		value = store.Get(key)
	)

	if value == nil {
		return auction, false
	}

	k.cdc.MustUnmarshal(value, &auction)
	return auction, true
}

func (k *Keeper) GetDebtAuctions(ctx sdk.Context) (auctions []auctiontypes.DebtAuction) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, auctiontypes.DebtAuctionKeyPrefix)
	)

	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var auction auctiontypes.DebtAuction
		k.cdc.MustUnmarshal(iter.Value(), &auction)
		auctions = append(auctions, auction)
	}

	return auctions
}

func (k *Keeper) GetDutchAuction(ctx sdk.Context, id uint64) (auction auctiontypes.DutchAuction, found bool) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.DutchAuctionKey(id)
		value = store.Get(key)
	)

	if value == nil {
		return auction, false
	}

	k.cdc.MustUnmarshal(value, &auction)
	return auction, true
}

func (k *Keeper) GetDutchAuctions(ctx sdk.Context) (auctions []auctiontypes.DutchAuction) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, auctiontypes.DutchAuctionKeyPrefix)
	)

	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var auction auctiontypes.DutchAuction
		k.cdc.MustUnmarshal(iter.Value(), &auction)
		auctions = append(auctions, auction)
	}

	return auctions
}

func (k *Keeper) GetBiddingID(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.BiddingsIdKey
		value = store.Get(key)
	)
	if value == nil {
		return 0
	}
	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue()
}

func (k *Keeper) SetBiddingID(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.BiddingsIdKey
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)

	store.Set(key, value)
}

func (k *Keeper) GetDebtBiddingID(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.DebtBiddingsIdKey
		value = store.Get(key)
	)
	if value == nil {
		return 0
	}
	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue()
}

func (k *Keeper) SetDebtBiddingID(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.DebtBiddingsIdKey
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)

	store.Set(key, value)
}

func (k *Keeper) GetDutchBiddingID(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.DutchBiddingsIdKey
		value = store.Get(key)
	)
	if value == nil {
		return 0
	}
	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue()
}

func (k *Keeper) SetDutchBiddingID(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.DutchBiddingsIdKey
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)

	store.Set(key, value)
}
func (k *Keeper) SetBidding(ctx sdk.Context, bidding auctiontypes.Biddings) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.BiddingsKey(bidding.Id)
		value = k.cdc.MustMarshal(&bidding)
	)
	store.Set(key, value)
}

func (k *Keeper) GetBidding(ctx sdk.Context, id uint64) (bidding auctiontypes.Biddings, found bool) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.BiddingsKey(id)
		value = store.Get(key)
	)

	if value == nil {
		return bidding, false
	}

	k.cdc.MustUnmarshal(value, &bidding)
	return bidding, true
}

func (k *Keeper) SetDebtBidding(ctx sdk.Context, bidding auctiontypes.Biddings) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.DebtBiddingsKey(bidding.Id)
		value = k.cdc.MustMarshal(&bidding)
	)
	store.Set(key, value)
}

func (k *Keeper) GetDebtBidding(ctx sdk.Context, id uint64) (bidding auctiontypes.Biddings, found bool) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.DebtBiddingsKey(id)
		value = store.Get(key)
	)

	if value == nil {
		return bidding, false
	}

	k.cdc.MustUnmarshal(value, &bidding)
	return bidding, true
}

func (k *Keeper) SetDutchBidding(ctx sdk.Context, bidding auctiontypes.DutchBiddings) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.DutchBiddingsKey(bidding.BiddingId)
		value = k.cdc.MustMarshal(&bidding)
	)
	store.Set(key, value)
}

func (k *Keeper) GetDutchBidding(ctx sdk.Context, id uint64) (bidding auctiontypes.DutchBiddings, found bool) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.DutchBiddingsKey(id)
		value = store.Get(key)
	)

	if value == nil {
		return bidding, false
	}

	k.cdc.MustUnmarshal(value, &bidding)
	return bidding, true
}

func (k *Keeper) GetBiddings(ctx sdk.Context) (biddings []auctiontypes.Biddings) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, auctiontypes.BiddingsKeyPrefix)
	)

	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var bidding auctiontypes.Biddings
		k.cdc.MustUnmarshal(iter.Value(), &bidding)
		biddings = append(biddings, bidding)
	}

	return biddings
}

func (k *Keeper) GetDebtBiddings(ctx sdk.Context) (biddings []auctiontypes.Biddings) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, auctiontypes.DebtAuctionKeyPrefix)
	)

	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var bidding auctiontypes.Biddings
		k.cdc.MustUnmarshal(iter.Value(), &bidding)
		biddings = append(biddings, bidding)
	}

	return biddings
}

func (k *Keeper) GetDutchBiddings(ctx sdk.Context) (biddings []auctiontypes.DutchBiddings) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, auctiontypes.DutchAuctionKeyPrefix)
	)

	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var bidding auctiontypes.DutchBiddings
		k.cdc.MustUnmarshal(iter.Value(), &bidding)
		biddings = append(biddings, bidding)
	}

	return biddings
}

func (k *Keeper) GetUserBiddingID(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.UserBiddingsIdKey
		value = store.Get(key)
	)
	if value == nil {
		return 0
	}
	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue()
}

func (k *Keeper) SetUserBiddingID(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.UserBiddingsIdKey
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)

	store.Set(key, value)
}

func (k *Keeper) GetDebtUserBiddingID(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.DebtUserBiddingsIdKey
		value = store.Get(key)
	)
	if value == nil {
		return 0
	}
	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue()
}

func (k *Keeper) SetDebtUserBiddingID(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.DebtUserBiddingsIdKey
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)

	store.Set(key, value)
}

func (k *Keeper) GetDutchUserBiddingID(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.DutchUserBiddingsIdKey
		value = store.Get(key)
	)
	if value == nil {
		return 0
	}
	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue()
}

func (k *Keeper) SetDutchUserBiddingID(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.DutchUserBiddingsIdKey
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)

	store.Set(key, value)
}

func (k *Keeper) SetUserBidding(ctx sdk.Context, userBiddings auctiontypes.UserBiddings) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.UserBiddingsKey(userBiddings.Bidder)
		value = k.cdc.MustMarshal(&userBiddings)
	)
	store.Set(key, value)
}

func (k *Keeper) GetUserBiddings(ctx sdk.Context, bidder string) (userBiddings auctiontypes.UserBiddings, found bool) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.UserBiddingsKey(bidder)
		value = store.Get(key)
	)

	if value == nil {
		return userBiddings, false
	}

	k.cdc.MustUnmarshal(value, &userBiddings)
	return userBiddings, true
}

func (k *Keeper) SetDebtUserBidding(ctx sdk.Context, userBiddings auctiontypes.UserBiddings) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.DebtUserBiddingsKey(userBiddings.Bidder)
		value = k.cdc.MustMarshal(&userBiddings)
	)
	store.Set(key, value)
}

func (k *Keeper) GetDebtUserBiddings(ctx sdk.Context, bidder string) (userBiddings auctiontypes.UserBiddings, found bool) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.DebtUserBiddingsKey(bidder)
		value = store.Get(key)
	)

	if value == nil {
		return userBiddings, false
	}

	k.cdc.MustUnmarshal(value, &userBiddings)
	return userBiddings, true
}

func (k *Keeper) SetDutchUserBidding(ctx sdk.Context, userBiddings auctiontypes.UserBiddings) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.DutchUserBiddingsKey(userBiddings.Bidder)
		value = k.cdc.MustMarshal(&userBiddings)
	)
	store.Set(key, value)
}

func (k *Keeper) GetDutchUserBiddings(ctx sdk.Context, bidder string) (userBiddings auctiontypes.UserBiddings, found bool) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.DutchUserBiddingsKey(bidder)
		value = store.Get(key)
	)

	if value == nil {
		return userBiddings, false
	}

	k.cdc.MustUnmarshal(value, &userBiddings)
	return userBiddings, true
}

func (k Keeper) CreateNewAuctions(ctx sdk.Context) {
	locked_vaults := k.GetLockedVaults(ctx)
	for _, locked_vault := range locked_vaults {
		pair, found := k.GetPair(ctx, locked_vault.PairId)
		if !found {
			continue
		}
		assetIn, found := k.GetAsset(ctx, pair.AssetIn)
		if !found {
			continue
		}

		assetOut, found := k.GetAsset(ctx, pair.AssetOut)
		if !found {
			continue
		}
		collateralizationRatio, err := k.CalculateCollaterlizationRatio(ctx, locked_vault.AmountIn, assetIn, locked_vault.AmountOut, assetOut)
		if err != nil {
			continue
		}
		if sdk.Dec.LT(collateralizationRatio, pair.LiquidationRatio) && !locked_vault.IsAuctionInProgress {
			k.StartCollateralAuction(ctx, locked_vault, pair, assetIn, assetOut)
		}
	}
}

func (k Keeper) CloseAuctions(ctx sdk.Context) {
	collateral_auctions := k.GetCollateralAuctions(ctx)
	for _, collateral_auction := range collateral_auctions {
		if ctx.BlockTime().After(collateral_auction.EndTime) {
			k.CloseCollateralAuction(ctx, collateral_auction)
		}
	}
}

func (k Keeper) CloseDebtAuctions(ctx sdk.Context) {
	debtAuctions := k.GetDebtAuctions(ctx)
	for _, debtAuction := range debtAuctions {
		if ctx.BlockTime().After(debtAuction.EndTime) {
			k.CloseDebtAuction(ctx, debtAuction)
		}
	}
}

func (k Keeper) RestartDutchAuctions(ctx sdk.Context) {
	dutchAuctions := k.GetDutchAuctions(ctx)
	auctionParams := k.GetParams(ctx)
	// SET current price of inflow token and outflow token
	for _, dutchAuction := range dutchAuctions {
		//TODO get price of dai from market module and set it in InflowTokenCurrentPrice of auction

		inFlowTokenCurrentPrice := sdk.MustNewDecFromStr("1")
		dutchAuction.InflowTokenCurrentPrice = inFlowTokenCurrentPrice
		tau := sdk.NewInt(int64(auctionParams.AuctionDurationSeconds))
		dur := ctx.BlockTime().Sub(dutchAuction.StartTime)
		seconds := sdk.NewInt(int64(dur.Seconds()))
		outFlowTokenCurrentPrice := k.getPriceFromLinearDecreaseFunction(dutchAuction.OutflowTokenInitialPrice, tau, seconds)
		dutchAuction.OutflowTokenCurrentPrice = outFlowTokenCurrentPrice

		//check if auction need to be restarted
		if ctx.BlockTime().After(dutchAuction.EndTime) || outFlowTokenCurrentPrice.LT(dutchAuction.OutflowTokenEndPrice) {
			//SET initial price fetched from market module and also end price , start time , end time
			//TODO get price of outflowtoken from market module
			outFlowTokenCurrentPrice := sdk.MustNewDecFromStr("10")
			timeNow := ctx.BlockTime()
			dutchAuction.StartTime = timeNow
			dutchAuction.EndTime = timeNow.Add(time.Second * time.Duration(auctionParams.AuctionDurationSeconds))
			outFlowTokenInitialPrice := k.getOutflowTokenInitialPrice(outFlowTokenCurrentPrice, auctionParams.Buffer)
			outFlowTokenEndPrice := k.getOutflowTokenEndPrice(outFlowTokenInitialPrice, auctionParams.Cusp)
			dutchAuction.OutflowTokenInitialPrice = outFlowTokenInitialPrice
			dutchAuction.OutflowTokenEndPrice = outFlowTokenEndPrice
			dutchAuction.OutflowTokenCurrentPrice = outFlowTokenInitialPrice
		}
		k.SetDutchAuction(ctx, dutchAuction)
	}
}

func (k Keeper) StartCollateralAuction(
	ctx sdk.Context,
	locked_vault liquidationtypes.LockedVault,
	pair assettypes.Pair,
	assetIn assettypes.Asset,
	assetOut assettypes.Asset,
) error {

	assetInPrice, found := k.GetPriceForAsset(ctx, pair.AssetIn)
	if !found {
		return assettypes.ErrorAssetDoesNotExist
	}

	assetOutPrice, found := k.GetPriceForAsset(ctx, pair.AssetOut)
	if !found {
		return assettypes.ErrorAssetDoesNotExist
	}

	if locked_vault.CollateralToBeAuctioned.LTE(sdk.NewDec(0)) {
		return auctiontypes.ErrorInvalidAuctioningCollateral
	}

	auctionParams := k.GetParams(ctx)

	liquidatedQuantity := sdk.NewDec(locked_vault.CollateralToBeAuctioned.Quo(sdk.NewDec(int64(assetInPrice))).RoundInt64())

	penaltyQuantity := liquidatedQuantity.Mul(sdk.MustNewDecFromStr(auctionParams.LiquidationPenaltyPercent).Mul(sdk.NewDec(100))).Quo(sdk.NewDec(100))
	discountedQuantity := liquidatedQuantity.Mul(sdk.MustNewDecFromStr(auctionParams.AuctionDiscountPercent).Mul(sdk.NewDec(100))).Quo(sdk.NewDec(100))
	auctioningQuantity := liquidatedQuantity.Sub(penaltyQuantity.Add(discountedQuantity))
	minBid := auctioningQuantity.Mul(sdk.NewDec(int64(assetInPrice))).Quo(sdk.NewDec(int64(assetOutPrice))).Ceil().RoundInt()
	maxBid := auctioningQuantity.Add(discountedQuantity).Mul(sdk.NewDec(int64(assetInPrice))).Quo(sdk.NewDec(int64(assetOutPrice))).Ceil().RoundInt()

	auction := auctiontypes.CollateralAuction{
		LockedVaultId:       locked_vault.LockedVaultId,
		AuctionedCollateral: sdk.NewCoin(assetIn.Denom, sdk.NewInt(auctioningQuantity.RoundInt64())),
		DiscountQuantity:    sdk.NewCoin(assetIn.Denom, sdk.NewInt(discountedQuantity.RoundInt64())),
		ActiveBiddingId:     0,
		Bidder:              nil,
		Bid:                 sdk.NewCoin(assetOut.Denom, sdk.NewInt(0)),
		MinBid:              sdk.NewCoin(assetOut.Denom, minBid),
		MaxBid:              sdk.NewCoin(assetOut.Denom, maxBid),
		EndTime:             ctx.BlockTime().Add(time.Second * time.Duration(auctionParams.AuctionDurationSeconds)),
		Pair:                pair,
		BiddingIds:          []uint64{},
	}
	auction.Id = k.GetCollateralAuctionID(ctx) + 1
	k.SetCollateralAuctionID(ctx, auction.Id)
	k.SetCollateralAuction(ctx, auction)
	k.SetFlagIsAuctionInProgress(ctx, locked_vault.LockedVaultId, true)
	return nil
}

func (k Keeper) StartDebtAuction(
	ctx sdk.Context,
	auctionToken sdk.Coin,
	expectedUserToken sdk.Coin,
	sourceAddress string,
) error {

	sourceAddress1, err := sdk.AccAddressFromBech32(sourceAddress)
	if err != nil {
		return auctiontypes.ErrorInvalidAddress
	}
	if err := k.bank.SendCoinsFromAccountToModule(ctx, sourceAddress1, liquidationtypes.ModuleName, sdk.NewCoins(auctionToken)); err != nil {
		return err
	}
	auctionParams := k.GetParams(ctx)
	auction := auctiontypes.DebtAuction{
		AuctionedToken:      auctionToken,
		ExpectedMintedToken: auctionToken,
		ExpectedUserToken:   expectedUserToken,
		ActiveBiddingId:     0,
		Bidder:              nil,
		EndTime:             ctx.BlockTime().Add(time.Second * time.Duration(auctionParams.AuctionDurationSeconds)),
		CurrentBidAmount:    sdk.NewCoin(auctionToken.Denom, sdk.NewInt(0)),
		AuctionStatus:       auctiontypes.AuctionStartNoBids,
	}
	auction.AuctionId = k.GetDebtAuctionID(ctx) + 1
	k.SetDebtAuctionID(ctx, auction.AuctionId)
	k.SetDebtAuction(ctx, auction)
	return nil
}

func (k Keeper) StartDutchAuction(
	ctx sdk.Context,
	outFlowToken sdk.Coin,
	inFlowToken sdk.Coin,
	outFlowTokenAddress string,
	inFlowTokenAddress string,
) error {

	outFlowTokenAddress1, err := sdk.AccAddressFromBech32(outFlowTokenAddress)
	inFlowTokenAddress1, err := sdk.AccAddressFromBech32(inFlowTokenAddress)
	if err != nil {
		return auctiontypes.ErrorInvalidAddress
	}
	outFlowTokenBalance := k.bank.GetBalance(ctx, outFlowTokenAddress1, outFlowToken.Denom)
	if outFlowTokenBalance.Amount.LT(outFlowToken.Amount) {
		return auctiontypes.ErrorInvalidOutFlowTokenBalance
	}
	auctionParams := k.GetParams(ctx)
	//TODO need to get real price instead of hard coding
	outFlowTokenPrice := sdk.MustNewDecFromStr("100")
	inFlowTokenPrice := sdk.MustNewDecFromStr("1")
	outFlowTokenInitialPrice := k.getOutflowTokenInitialPrice(outFlowTokenPrice, auctionParams.Buffer)
	outFlowTokenEndPrice := k.getOutflowTokenEndPrice(outFlowTokenInitialPrice, auctionParams.Cusp)
	timeNow := ctx.BlockTime()
	auction := auctiontypes.DutchAuction{
		OutflowTokenInitAmount:    outFlowToken,
		OutflowTokenCurrentAmount: outFlowToken,
		InflowTokenTargetAmount:   inFlowToken,
		InflowTokenCurrentAmount:  inFlowToken,
		OutflowTokenInitialPrice:  outFlowTokenInitialPrice,
		OutflowTokenCurrentPrice:  outFlowTokenInitialPrice,
		OutflowTokenEndPrice:      outFlowTokenEndPrice,
		InflowTokenCurrentPrice:   inFlowTokenPrice,
		StartTime:                 timeNow,
		EndTime:                   timeNow.Add(time.Second * time.Duration(auctionParams.AuctionDurationSeconds)),
		AuctionStatus:             auctiontypes.AuctionStartNoBids,
		InflowTokenAddress:        outFlowTokenAddress1,
		OutflowTokenAddress:       inFlowTokenAddress1,
	}
	auction.AuctionId = k.GetDutchAuctionID(ctx) + 1
	k.SetDebtAuctionID(ctx, auction.AuctionId)
	k.SetDutchAuction(ctx, auction)
	return nil
}

func (k Keeper) CloseCollateralAuction(
	ctx sdk.Context,
	collateral_auction auctiontypes.CollateralAuction,
) error {

	if collateral_auction.Bidder != nil && collateral_auction.Bid.Amount.GTE(collateral_auction.MinBid.Amount) {

		assetIn, found := k.GetAsset(ctx, collateral_auction.Pair.AssetIn)
		if !found {
			return assettypes.ErrorAssetDoesNotExist
		}
		assetOut, found := k.GetAsset(ctx, collateral_auction.Pair.AssetOut)
		if !found {
			return assettypes.ErrorAssetDoesNotExist
		}

		highestBidReceived := collateral_auction.Bid

		collateralQuantity := collateral_auction.AuctionedCollateral.Amount.Add(collateral_auction.DiscountQuantity.Amount)

		err := k.bank.SendCoinsFromModuleToAccount(ctx, vaulttypes.ModuleName, collateral_auction.Bidder, sdk.NewCoins(sdk.NewCoin(assetIn.Denom, collateralQuantity)))
		if err != nil {
			return err
		}
		bidding, _ := k.GetBidding(ctx, collateral_auction.ActiveBiddingId)
		bidding.BiddingStatus = auctiontypes.SuccessBiddingStatus
		k.SetBidding(ctx, bidding)
		k.BurnCAssets(ctx, liquidationtypes.ModuleName, assetIn.Denom, assetOut.Denom, highestBidReceived.Amount)
		k.UpdateAssetQuantitiesInLockedVault(ctx, collateral_auction, collateralQuantity, assetIn, highestBidReceived.Amount, assetOut)

		for _, biddingId := range collateral_auction.BiddingIds {
			bidding, found := k.GetBidding(ctx, biddingId)
			if !found {
				continue
			}
			bidding.AuctionStatus = auctiontypes.ClosedAuctionStatus
			k.SetBidding(ctx, bidding)
		}
	}
	k.SetFlagIsAuctionComplete(ctx, collateral_auction.LockedVaultId, true)
	k.SetFlagIsAuctionInProgress(ctx, collateral_auction.LockedVaultId, false)
	k.DeleteCollateralAuction(ctx, collateral_auction.Id)
	return nil
}

func (k Keeper) CloseDebtAuction(
	ctx sdk.Context,
	debtAuction auctiontypes.DebtAuction,
) error {

	if debtAuction.AuctionStatus != auctiontypes.AuctionStartNoBids {

		err := k.bank.SendCoinsFromModuleToAccount(ctx, liquidationtypes.ModuleName, debtAuction.Bidder, sdk.NewCoins(debtAuction.AuctionedToken))
		if err != nil {
			return err
		}
		bidding, _ := k.GetDebtBidding(ctx, debtAuction.ActiveBiddingId)
		bidding.BiddingStatus = auctiontypes.SuccessBiddingStatus
		k.SetBidding(ctx, bidding)

		//for _, biddingId := range debtAuction.BiddingIds {
		//	bidding, found := k.GetBidding(ctx, biddingId)
		//	if !found {
		//		continue
		//	}
		//	bidding.AuctionStatus = auctiontypes.ClosedAuctionStatus
		//	k.SetBidding(ctx, bidding)
		//}
	}
	k.bank.BurnCoins(ctx, liquidationtypes.ModuleName, sdk.NewCoins(debtAuction.ExpectedUserToken))
	k.DeleteCollateralAuction(ctx, debtAuction.AuctionId)
	return nil
}

func (k Keeper) CloseDutchAuction(
	ctx sdk.Context,
	dutchtAuction auctiontypes.DutchAuction,
) error {

	if dutchtAuction.AuctionStatus != auctiontypes.AuctionStartNoBids {
		//for _, biddingId := range debtAuction.BiddingIds {
		//	bidding, found := k.GetBidding(ctx, biddingId)
		//	if !found {
		//		continue
		//	}
		//	bidding.AuctionStatus = auctiontypes.ClosedAuctionStatus
		//	k.SetBidding(ctx, bidding)
		//}
	}
	k.DeleteDutchAuction(ctx, dutchtAuction.AuctionId)
	return nil
}

func (k Keeper) CreateNewBid(ctx sdk.Context, auctionId uint64, bidder sdk.AccAddress, bid sdk.Coin) (biddingId uint64, err error) {
	auction, found := k.GetCollateralAuction(ctx, auctionId)
	if !found {
		return 0, auctiontypes.ErrorInvalidAuctionId
	}
	bidding := auctiontypes.Biddings{
		Id:                  k.GetBiddingID(ctx) + 1,
		AuctionId:           auctionId,
		AuctionStatus:       auctiontypes.ActiveAuctionStatus,
		AuctionedCollateral: auction.AuctionedCollateral,
		Bidder:              bidder.String(),
		Bid:                 bid,
		BiddingTimestamp:    ctx.BlockTime(),
		BiddingStatus:       auctiontypes.PlacedBiddingStatus,
	}
	k.SetBiddingID(ctx, bidding.Id)
	k.SetBidding(ctx, bidding)

	userBiddings, found := k.GetUserBiddings(ctx, bidder.String())
	if !found {
		userBiddings = auctiontypes.UserBiddings{
			Id:         k.GetUserBiddingID(ctx) + 1,
			Bidder:     bidder.String(),
			BiddingIds: []uint64{},
		}
		k.SetUserBiddingID(ctx, userBiddings.Id)
	}
	userBiddings.BiddingIds = append(userBiddings.BiddingIds, bidding.Id)
	k.SetUserBidding(ctx, userBiddings)
	return bidding.Id, nil
}

func (k Keeper) CreateNewDebtBid(ctx sdk.Context, auctionId uint64, bidder sdk.AccAddress, bid sdk.Coin) (biddingId uint64, err error) {
	bidding := auctiontypes.Biddings{
		Id:               k.GetDebtBiddingID(ctx) + 1,
		AuctionId:        auctionId,
		AuctionStatus:    auctiontypes.ActiveAuctionStatus,
		Bidder:           bidder.String(),
		Bid:              bid,
		BiddingTimestamp: ctx.BlockTime(),
		BiddingStatus:    auctiontypes.PlacedBiddingStatus,
	}
	k.SetDebtBiddingID(ctx, bidding.Id)
	k.SetDebtBidding(ctx, bidding)

	userBiddings, found := k.GetDebtUserBiddings(ctx, bidder.String())
	if !found {
		userBiddings = auctiontypes.UserBiddings{
			Id:         k.GetDebtUserBiddingID(ctx) + 1,
			Bidder:     bidder.String(),
			BiddingIds: []uint64{},
		}
		k.SetDebtUserBiddingID(ctx, userBiddings.Id)
	}
	userBiddings.BiddingIds = append(userBiddings.BiddingIds, bidding.Id)
	k.SetDebtUserBidding(ctx, userBiddings)
	return bidding.Id, nil
}

func (k Keeper) CreateNewDutchBid(ctx sdk.Context, auctionId uint64, bidder sdk.AccAddress, outFlowTokenCoin sdk.Coin, inFlowTokenCoin sdk.Coin) (biddingId uint64) {
	bidding := auctiontypes.DutchBiddings{
		BiddingId:          k.GetDutchBiddingID(ctx) + 1,
		AuctionId:          auctionId,
		AuctionStatus:      auctiontypes.ActiveAuctionStatus,
		Bidder:             bidder.String(),
		OutflowTokenAmount: outFlowTokenCoin,
		InflowTokenAmount:  inFlowTokenCoin,
		BiddingTimestamp:   ctx.BlockTime(),
		BiddingStatus:      auctiontypes.SuccessBiddingStatus,
	}
	k.SetDutchBiddingID(ctx, bidding.BiddingId)
	k.SetDutchBidding(ctx, bidding)

	userDutchBiddings, found := k.GetDutchUserBiddings(ctx, bidder.String())
	if !found {
		userDutchBiddings = auctiontypes.UserBiddings{
			Id:         k.GetDutchUserBiddingID(ctx) + 1,
			Bidder:     bidder.String(),
			BiddingIds: []uint64{},
		}
		k.SetDutchUserBiddingID(ctx, userDutchBiddings.Id)
	}
	userDutchBiddings.BiddingIds = append(userDutchBiddings.BiddingIds, bidding.BiddingId)
	k.SetDutchUserBidding(ctx, userDutchBiddings)
	return bidding.BiddingId
}

func (k Keeper) PlaceBid(ctx sdk.Context, auctionId uint64, bidder sdk.AccAddress, bid sdk.Coin) error {
	auction, found := k.GetCollateralAuction(ctx, auctionId)
	if !found {
		return auctiontypes.ErrorInvalidAuctionId
	}
	if bid.Denom != auction.MinBid.Denom {
		return auctiontypes.ErrorInvalidBiddingDenom
	}
	if bid.Amount.LT(auction.MinBid.Amount) {
		return auctiontypes.ErrorLowBidAmount
	}
	if bid.Amount.GT(auction.MaxBid.Amount) {
		return auctiontypes.ErrorMaxBidAmount
	}
	if bid.Amount.LT(auction.Bid.Amount.Add(sdk.NewInt(1))) {
		return auctiontypes.ErrorBidAlreadyExists
	}
	err := k.SendCoinsFromAccountToModule(ctx, bidder, liquidationtypes.ModuleName, sdk.NewCoins(bid))
	if err != nil {
		return err
	}
	biddingId, err := k.CreateNewBid(ctx, auctionId, bidder, bid)
	if err != nil {
		return err
	}
	// auction.Bidder as previous bidder
	//TODO Move this in the next if condition
	err = k.bank.SendCoinsFromModuleToAccount(ctx, liquidationtypes.ModuleName, auction.Bidder, sdk.NewCoins(auction.Bid))
	if err != nil {
		return err
	}
	if auction.ActiveBiddingId != 0 {
		bidding, _ := k.GetBidding(ctx, auction.ActiveBiddingId)
		bidding.BiddingStatus = auctiontypes.RejectedBiddingStatus
		k.SetBidding(ctx, bidding)
	}
	auction.ActiveBiddingId = biddingId
	auction.BiddingIds = append(auction.BiddingIds, biddingId)
	auction.Bidder = bidder
	auction.Bid = bid
	k.SetCollateralAuction(ctx, auction)
	return nil
}

func (k Keeper) PlaceDebtBid(ctx sdk.Context, auctionId uint64, bidder sdk.AccAddress, bid sdk.Coin, expectedUserToken sdk.Coin) error {
	auction, found := k.GetDebtAuction(ctx, auctionId)
	if !found {
		return auctiontypes.ErrorInvalidDebtAuctionId
	}
	if expectedUserToken.Denom != auction.ExpectedUserToken.Denom {
		return auctiontypes.ErrorInvalidDebtUserExpectedDenom
	}
	if expectedUserToken.Amount.Equal(auction.ExpectedUserToken.Amount) == false {
		return auctiontypes.ErrorDebtExpectedUserAmount
	}
	if bid.Denom != auction.ExpectedMintedToken.Denom {
		return auctiontypes.ErrorInvalidDebtMintedDenom
	}
	if bid.Amount.GT(auction.ExpectedMintedToken.Amount) {
		return auctiontypes.ErrorDebtMoreBidAmount
	}

	// auction.Bidder as previous bidder
	err := k.SendCoinsFromAccountToModule(ctx, bidder, liquidationtypes.ModuleName, sdk.NewCoins(expectedUserToken))
	if err != nil {
		return err
	}
	biddingId, err := k.CreateNewDebtBid(ctx, auctionId, bidder, bid)
	if err != nil {
		return err
	}
	//If auction gets bid from second time onwards
	if auction.AuctionStatus != auctiontypes.AuctionStartNoBids {
		err = k.bank.SendCoinsFromModuleToAccount(ctx, liquidationtypes.ModuleName, auction.Bidder, sdk.NewCoins(auction.ExpectedUserToken))
		if err != nil {
			return err
		}
		bidding, _ := k.GetDebtBidding(ctx, auction.ActiveBiddingId)
		bidding.BiddingStatus = auctiontypes.RejectedBiddingStatus
		k.SetDebtBidding(ctx, bidding)
	} else {
		auction.AuctionStatus = auctiontypes.AuctionGoingOn
	}
	auction.ActiveBiddingId = biddingId
	auction.Bidder = bidder
	auction.CurrentBidAmount = bid
	//decreasing expected minted token for next bid
	expectedMintedToken := sdk.NewDecFromInt(auction.ExpectedMintedToken.Amount)
	decreaseAmount := expectedMintedToken.Mul(auctiontypes.DefaultDebtMintTokenDecreasePercentage)
	expectedMintedToken = expectedMintedToken.Sub(decreaseAmount).Ceil() // As of now ceiling is done
	auction.ExpectedMintedToken = sdk.NewCoin(auction.ExpectedMintedToken.Denom, expectedMintedToken.TruncateInt())
	k.SetDebtAuction(ctx, auction)
	return nil
}

func (k Keeper) PlaceDutchBid(ctx sdk.Context, auctionId uint64, bidder sdk.AccAddress, bid sdk.Coin, max sdk.Dec) error {
	auction, found := k.GetDutchAuction(ctx, auctionId)
	if !found {
		return auctiontypes.ErrorInvalidDutchAuctionId
	}

	if bid.Denom != auction.OutflowTokenCurrentAmount.Denom {
		return auctiontypes.ErrorInvalidDutchUserbidDenom
	}

	if max.LT(auction.OutflowTokenCurrentPrice) {
		return auctiontypes.ErrorInvalidDutchPrice
	}
	// slice tells amount of collateral user should be given
	auctionParams := k.GetParams(ctx)
	outFlowTokenCurrentPrice := sdk.NewIntFromBigInt(auction.OutflowTokenCurrentPrice.BigInt())
	inFlowTokenCurrentPrice := sdk.NewIntFromBigInt(auction.InflowTokenCurrentPrice.BigInt())
	slice := sdk.MinInt(bid.Amount, auction.OutflowTokenCurrentAmount.Amount)
	owe := slice.Mul(outFlowTokenCurrentPrice)
	tab := auction.InflowTokenCurrentAmount.Amount.Mul(inFlowTokenCurrentPrice)
	//check if bid is greater than target dai
	if owe.GT(tab) {
		slice = tab.Quo(sdk.NewIntFromBigInt(auction.OutflowTokenCurrentPrice.BigInt()))
	} else if auction.OutflowTokenCurrentAmount.Amount.Sub(slice).Mul(outFlowTokenCurrentPrice).LT(sdk.NewIntFromBigInt(auctionParams.Chost.BigInt())) {
		//(outflowtokenavailableamount-slice) in usd < chost in usd
		//see if user has balance to buy whole collateral
		user_balance_usd := k.bank.GetBalance(ctx, bidder, bid.Denom).Amount.Mul(outFlowTokenCurrentPrice)
		collateral_available_usd := auction.OutflowTokenCurrentAmount.Amount.Mul(outFlowTokenCurrentPrice)
		if user_balance_usd.LT(collateral_available_usd) {
			return auctiontypes.ErrorDutchinsufficientUserBalance
		}
		slice = auction.OutflowTokenCurrentAmount.Amount
	}

	inFlowTokenToCharge := slice.Mul(outFlowTokenCurrentPrice).Quo(inFlowTokenCurrentPrice)
	inFlowTokenCoin := sdk.NewCoin(auction.InflowTokenTargetAmount.Denom, inFlowTokenToCharge)
	outFlowTokenCoin := sdk.NewCoin(auction.OutflowTokenInitAmount.Denom, slice)
	err := k.SendCoinsFromAccountToModule(ctx, bidder, liquidationtypes.ModuleName, sdk.NewCoins(inFlowTokenCoin))
	if err != nil {
		return err
	}
	err = k.SendCoinsFromModuleToAccount(ctx, liquidationtypes.ModuleName, bidder, sdk.NewCoins(outFlowTokenCoin))
	if err != nil {
		//refund tokens to user as he is not getting collateral
		k.SendCoinsFromModuleToAccount(ctx, liquidationtypes.ModuleName, bidder, sdk.NewCoins(inFlowTokenCoin))
		return err
	}

	k.CreateNewDutchBid(ctx, auctionId, bidder, inFlowTokenCoin, outFlowTokenCoin)
	auction.OutflowTokenCurrentAmount = auction.OutflowTokenCurrentAmount.Sub(outFlowTokenCoin)
	auction.InflowTokenCurrentAmount = auction.InflowTokenCurrentAmount.Add(inFlowTokenCoin)
	if auction.InflowTokenTargetAmount.IsEqual(inFlowTokenCoin) {
		//TODO return collateral to vault as target dai reached
		//remove dutch auction
		k.CloseDutchAuction(ctx, auction)
		return nil
	}
	if auction.OutflowTokenCurrentAmount.Amount.IsZero() {
		k.CloseDutchAuction(ctx, auction)
		//TODO send target dai to the module specified
		return nil
	}
	k.SetDutchAuction(ctx, auction)
	return nil
}
