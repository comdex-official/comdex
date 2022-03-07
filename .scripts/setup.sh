make all

rm -rf ~/.comdex

mkdir ~/.comdex

comdex init --chain-id test test
comdex keys add test --keyring-backend test
comdex add-genesis-account test 100000000000000stake,100000000000000ucmdx,100000000000000ucgold,100000000000000ucsilver,100000000000000ucoil,100000000000000ibc/4294C3DB67564CF4A0B2BFACC8415A59B38243F6FF9E288FBA34F9B4823BA16E,10000000000ibc/C4CFF46FD6DE35CA4CF4CE031E643C8FDC9BA4B99AE598E9B0ED98FE3A2319F9 --keyring-backend test
comdex gentx test 1000000000stake --chain-id test --keyring-backend test
comdex collect-gentxs
# Make sure to add the admin account address in params after running this script and before starting the chain.
# use this account to generate assets and asset pairs

comdex tx asset add-asset CMDX ucmdx 1000000 --from test --chain-id test --keyring-backend test -y
comdex tx asset add-asset STAKE stake 1000000 --from test --chain-id test --keyring-backend test -y
comdex tx asset add-asset XAU ucgold 1000000 --from test --chain-id test --keyring-backend test -y
comdex tx asset add-asset XAG ucsilver 1000000 --from test --chain-id test --keyring-backend test -y


comdex tx asset add-asset ATOM ibc/C4CFF46FD6DE35CA4CF4CE031E643C8FDC9BA4B99AE598E9B0ED98FE3A2319F9 1000000 --from test --chain-id test --keyring-backend test -y
comdex tx asset add-asset XAU ucgold 1000000 --from test --chain-id test --keyring-backend test -y
comdex tx asset add-asset XAG ucsilver 1000000 --from test --chain-id test --keyring-backend test -y
comdex tx asset add-asset OIL ucoil 1000000 --from test --chain-id test --keyring-backend test -y
comdex tx asset add-asset UST ibc/4294C3DB67564CF4A0B2BFACC8415A59B38243F6FF9E288FBA34F9B4823BA16E 1000000 --from test --chain-id test --keyring-backend test -y
comdex tx asset add-asset CMDX ucmdx 1000000 --from test --chain-id test --keyring-backend test -y


comdex tx asset add-pair 6 2 1.5 --from cooluser --chain-id test-1 --keyring-backend test -y 
comdex tx asset add-pair 6 3 1.5 --from cooluser --chain-id test-1 --keyring-backend test -y 
comdex tx asset add-pair 6 4 1.5 --from cooluser --chain-id test-1 --keyring-backend test -y 

comdex tx asset add-pair 1 3 1.5 --from test --chain-id test --keyring-backend test -y 
comdex tx asset add-pair 1 4 1.5 --from test --chain-id test --keyring-backend test -y 
comdex tx asset add-pair 2 3 1.5 --from test --chain-id test --keyring-backend test -y 
comdex tx asset add-pair 2 4 1.5 --from test --chain-id test --keyring-backend test -y 
comdex tx asset add-pair 2 1 1.5 --from test --chain-id test --keyring-backend test -y 
comdex tx asset add-pair 3 1 1.5 --from test --chain-id test --keyring-backend test -y 

comdex tx vault create 1 200 100 --from test2 --chain-id test --keyring-backend test -y
comdex tx vault create 2 100 50 --from test2 --chain-id test --keyring-backend test -y
comdex tx vault create 3 100 50 --from test2 --chain-id test --keyring-backend test -y
comdex tx vault create 4 100 50 --from test2 --chain-id test --keyring-backend test -y
comdex tx vault create 5 100 50 --from test2 --chain-id test --keyring-backend test -y
comdex tx vault create 6 100 50 --from test2 --chain-id test --keyring-backend test -y

comdex tx asset add-asset ATOM uatom 1000000 --from test --chain-id test --keyring-backend test
comdex tx asset add-asset XAU ucgold 1000000 --from test --chain-id test --keyring-backend test

comdex tx asset add-asset ATOM uatom 1000000 --from test --chain-id test --keyring-backend test -y
comdex tx asset add-asset XAU ucgold 1000000 --from test --chain-id test --keyring-backend test -y
comdex tx asset add-pair 1 2 1.5 --from test --chain-id test --keyring-backend test -y 

comdex tx vault create 1 30 10 --from test --chain-id test --keyring-backend test -y



markets:
- rates: "30880000"
  script_id: "112"
  symbol: ATOM
- rates: "2660000"
  script_id: "112"
  symbol: CMDX
- rates: "92810000"
  script_id: "112"
  symbol: OIL
- rates: "1000000"
  script_id: "112"
  symbol: UST
- rates: "22500000"
  script_id: "112"
  symbol: XAG
- rates: "1807950000"
  script_id: "112"
  symbol: XAU


comdex tx asset add-pair 1 2 1.5 --from cooluser --chain-id test-1 --keyring-backend test -y 
comdex tx asset add-pair 1 3 1.5 --from cooluser --chain-id test-1 --keyring-backend test -y 
comdex tx asset add-pair 1 4 1.5 --from cooluser --chain-id test-1 --keyring-backend test -y 
comdex tx asset add-pair 5 2 1.5 --from cooluser --chain-id test-1 --keyring-backend test -y 
comdex tx asset add-pair 5 3 1.5 --from cooluser --chain-id test-1 --keyring-backend test -y 
comdex tx asset add-pair 5 4 1.5 --from cooluser --chain-id test-1 --keyring-backend test -y 

comdex tx asset add-pair 6 3 1.5 --from cooluser --chain-id test-1 --keyring-backend test -y 



script hero outside lecture region leopard fortune much live galaxy chunk payment awkward enhance figure ill execute trash current wage roof where wine busy

1) remove existing data 
rm -rf ~/.comdex
rm -rf ~/.starport

2) build binary and move into path
make install
sudo mv ~/go/bin/comdex /usr/local/bin

3) add congiguration
comdex init --chain-id test test
comdex keys add test --keyring-backend test
comdex add-genesis-account test 100000000000000stake,100000000000000ucmdx,100000000000000ucgold,100000000000000ucsilver,100000000000000ucoil,100000000000000ibc/4294C3DB67564CF4A0B2BFACC8415A59B38243F6FF9E288FBA34F9B4823BA16E,10000000000ibc/C4CFF46FD6DE35CA4CF4CE031E643C8FDC9BA4B99AE598E9B0ED98FE3A2319F9 --keyring-backend test
comdex gentx test 1000000000stake --chain-id test --keyring-backend test
comdex collect-gentxs

4) Add admin address in genesis
sudo vim ~/.comdex/config/genesis.json

5) start chain
comdex start

6) configure relayer
starport relayer configure -a \
--source-rpc "http://rpc.laozi-testnet4.bandchain.org:80" \
--source-faucet "https://laozi-testnet4.bandchain.org/faucet" \
--source-port "oracle" \
--source-gasprice "0uband" \
--source-gaslimit 5000000 \
--source-prefix "band" \
--source-version "bandchain-1" \
--target-rpc "http://localhost:26657" \
--target-faucet "http://localhost:4500" \
--target-port "bandoracle" \
--target-gasprice "0.0ucmdx" \
--target-gaslimit 300000 \
--target-prefix "comdex"  \
--target-version "bandchain-1"

7) transfer cmdx tokens in relayer generated account
comdex tx bank send test comdex.... 1000000ucmdx --from test --chain-id test --keyring-backend test -y

8) connect the relayer
starport relayer connect

9) Add assets in chain
comdex tx asset add-asset ATOM ibc/C4CFF46FD6DE35CA4CF4CE031E643C8FDC9BA4B99AE598E9B0ED98FE3A2319F9 1000000 --from test --chain-id test --keyring-backend test -y
comdex tx asset add-asset XAU ucgold 1000000 --from test --chain-id test --keyring-backend test -y
comdex tx asset add-asset XAG ucsilver 1000000 --from test --chain-id test --keyring-backend test -y
comdex tx asset add-asset OIL ucoil 1000000 --from test --chain-id test --keyring-backend test -y
comdex tx asset add-asset UST ibc/4294C3DB67564CF4A0B2BFACC8415A59B38243F6FF9E288FBA34F9B4823BA16E 1000000 --from test --chain-id test --keyring-backend test -y
comdex tx asset add-asset CMDX ucmdx 1000000 --from test --chain-id test --keyring-backend test -y

10) Trigger price fetch from oracle
comdex tx bandoracle fetch-price-data 112 4 3 --channel channel-0 --symbols "ATOM" --multiplier 1000000 --fee-limit 250000uband --prepare-gas 600000 --execute-gas 600000 --from test --chain-id test --keyring-backend test -y

11) Add asset pairs
# CMDX as collateral 
comdex tx asset add-pair 6 2 1.5 --from test --chain-id test --keyring-backend test -y 
comdex tx asset add-pair 6 3 1.5 --from test --chain-id test --keyring-backend test -y 
comdex tx asset add-pair 6 4 1.5 --from test --chain-id test --keyring-backend test -y 

# ATOM as collateral 
comdex tx asset add-pair 1 2 1.5 --from test --chain-id test --keyring-backend test -y 
comdex tx asset add-pair 1 3 1.5 --from test --chain-id test --keyring-backend test -y 
comdex tx asset add-pair 1 4 1.5 --from test --chain-id test --keyring-backend test -y 

# UST as collateral 
comdex tx asset add-pair 5 2 1.5 --from test --chain-id test --keyring-backend test -y 
comdex tx asset add-pair 5 3 1.5 --from test --chain-id test --keyring-backend test -y 
comdex tx asset add-pair 5 4 1.5 --from test --chain-id test --keyring-backend test -y 

12) create pools # calculate quantity based on 50%-50% ratio
comdex tx liquidity create-pool 1 100000000ucmdx,100000000ucgold --from test --chain-id test --keyring-backend test -y
comdex tx liquidity create-pool 1 100000000ucmdx,100000000ucgold --from test --chain-id test --keyring-backend test -y
comdex tx liquidity create-pool 1 100000000ucmdx,100000000ucgold --from test --chain-id test --keyring-backend test -y