module github.com/comdex-official/comdex

go 1.16

require (
	github.com/CosmWasm/wasmd v0.22.0
	github.com/bandprotocol/bandchain-packet v0.0.2
	github.com/cosmos/cosmos-sdk v0.45.1
	github.com/cosmos/ibc-go v1.2.2
	github.com/cosmos/ibc-go/v2 v2.0.2
	github.com/fsnotify/fsnotify v1.5.1 // indirect
	github.com/gogo/protobuf v1.3.3
	github.com/golang/protobuf v1.5.2
	github.com/golang/snappy v0.0.3 // indirect
	github.com/gorilla/mux v1.8.0
	github.com/gravity-devs/liquidity v1.4.0
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/pelletier/go-toml v1.9.4 // indirect
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.12.0
	github.com/spf13/cast v1.4.1
	github.com/spf13/cobra v1.2.1
	github.com/stretchr/testify v1.7.0
	github.com/swaggo/swag v1.7.4
	github.com/tendermint/spm v0.1.9
	github.com/tendermint/tendermint v0.34.15
	github.com/tendermint/tm-db v0.6.6
	google.golang.org/genproto v0.0.0-20211208223120-3a66f561d7aa
	google.golang.org/grpc v1.43.0
	google.golang.org/protobuf v1.27.1
	gopkg.in/yaml.v2 v2.4.0
)

replace (
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
	google.golang.org/grpc => google.golang.org/grpc v1.33.2
)
