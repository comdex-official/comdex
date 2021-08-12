#pre-requisites
1. # In order to build and run the relayer you need to install and configure Rust on your machine.
https://hermes.informal.systems/pre_requisites.html


2. #install hermis
cargo install ibc-relayer-cli --bin hermes --locked

3. #configure ~/.hermes/config.toml

add below configs to above file
[global]
strategy = 'packets'
log_level = 'info'

[[chains]]
id = 'test'
rpc_addr = 'http://127.0.0.1:26657'
grpc_addr = 'http://127.0.0.1:9090'
websocket_addr = 'ws://127.0.0.1:26657/websocket'
rpc_timeout = '10s'
account_prefix = 'cosmos'
key_name = 'comdexkey'
store_prefix = 'ibc'
max_gas = 100000000
gas_price = { price = 0.001, denom = 'satom' }
gas_adjustment = 0.1
clock_drift = '5s'
trusting_period = '14days'
trust_threshold = { numerator = '1', denominator = '3' }

[[chains]]
id = 'testhub'
rpc_addr = 'http://127.0.0.1:36657'
grpc_addr = 'http://127.0.0.1:9091'
websocket_addr = 'ws://127.0.0.1:36657/websocket'
rpc_timeout = '10s'
account_prefix = 'cosmos'
key_name = 'gaiakey'
store_prefix = 'ibc'
max_gas = 100000000
gas_price = { price = 0.001, denom = 'uatom' }
gas_adjustment = 0.1
clock_drift = '5s'
trusting_period = '14days'
trust_threshold = { numerator = '1', denominator = '3' }


4. #Start the the two chains. chain names in below case is test testhub.

5. configure keyring for the accounts on chains
hermes -c ~/.hermes/config.toml keys add test -f $jsonname.json
hermes -c ~/.hermes/config.toml keys add testhub -f $jsonname1.json
hermes -c ~/.hermes/config.toml keys add test -f $jsonname2.json
hermes -c ~/.hermes/config.toml keys add testhub -f $jsonname3.json

6. #create client creation, connection and channel handshake b/w chains

hermes create channel test testhub --port-a transfer --port-b transfer -o unordered

7. #start hermes - command ::
hermes start

8. #test transaction

hermes tx raw ft-transfer $DESTINATION_CHAIN $SOURCE_CHAIN transfer $CHANNEL 9999 -o 1000 -n 2