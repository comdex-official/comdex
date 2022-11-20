make all

rm -rf ~/.petri

mkdir ~/.petri

petri init --chain-id test test
petri keys add test --recover --keyring-backend test<<<"y
wage thunder live sense resemble foil apple course spin horse glass mansion midnight laundry acoustic rhythm loan scale talent push green direct brick please"
petri add-genesis-account test 100000000000000stake --keyring-backend test
petri gentx test 1000000000stake --chain-id test --keyring-backend test
petri collect-gentxs
petri start
