package keeper

type msgServer struct {
	Keeper
}

func NewMsgServerImpl(keeper Keeper) types.MsgServiceServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServiceServer = msgServer{}
