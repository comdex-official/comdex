package locker

type LockerQuery struct {
	State *State `json:"balance_at_height,omitempty"`
}

type State struct {
	Address string `json:"address"`
	Denom   string `json:"denom"`
	Height  string `json:"height"`
	Target  string `json:"target"`
}
