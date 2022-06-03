make install

rm -rf ~/.comdex

mkdir ~/.comdex

comdex init --chain-id test test
comdex keys add test --recover --keyring-backend test<<<"y
wage thunder live sense resemble foil apple course spin horse glass mansion midnight laundry acoustic rhythm loan scale talent push green direct brick please"
comdex add-genesis-account test 100000000000000stake --keyring-backend test
comdex gentx test 1000000000stake --chain-id test --keyring-backend test
comdex collect-gentxs
comdex start

comdex tx gov submit-proposal add-app-mapping composite com 10000 1000000 --title "Adding New Assets" --description "adding cmdx and atom" --deposit 1000000000stake --from test --chain-id test --keyring-backend test --gas auto -y

comdex tx gov vote 1  --from test --chain-id test --keyring-backend test --gas auto -y

# wait for proposal to pass
comdex tx gov submit-proposal add-assets CMDX,OSMO,ATOM,CMST,HARBOR ucmdx,uosmo,uatom,ucmst,uharbor 1000000,1000000,1000000,1000000,1000000 1,1,1,1,0 --title "Adding New Assets" --description "adding cmdx and atom" --deposit 1000000000stake --from test --chain-id test --keyring-backend test --gas auto -y

comdex tx gov vote 2  --from test --chain-id test --keyring-backend test --gas auto -y

# wait for proposal to pass
comdex tx gov submit-proposal add-pairs 1,1,1,1 2,3,4,5 --title "Adding New Assets" --description "adding cmdx and atom" --deposit 1000000000stake --from test --chain-id test --keyring-backend test --gas auto -y

comdex tx gov vote 3  --from test --chain-id test --keyring-backend test --gas auto -y

# wait for proposal to pass
comdex tx gov submit-proposal add-pairs-vault 1 1 1.4 0.1 0.01 0.0013 0.01 1 1000000 10000 0 1.5 CMDX-A 0 1 --title addAsset --description addingAsset --from test --deposit 10000000stake --chain-id test --keyring-backend test -y

comdex tx gov vote 4  --from test --chain-id test --keyring-backend test --gas auto -y

# wait for proposal to pass
comdex tx gov submit-proposal add-asset-mapping 1 1 10000 1 comdex1pkkayn066msg6kn33wnl5srhdt3tnu2v9jjqu0 --title "Adding New Assets" --description "adding cmdx and atom" --deposit 1000000000stake --from test --chain-id test --keyring-backend test --gas auto

comdex tx gov vote 5  --from test --chain-id test --keyring-backend test --gas auto -y

# wait for proposal to pass
comdex tx liquidation  whitelist-app-id 1 --from test --chain-id test --keyring-backend test -y

# create vault
comdex tx locker whitelist-asset-locker 0 1 --from test --chain-id test --keyring-backend test -y

# whitelist asset locker
comdex tx vault create 1 1 5000 10000 --from test --chain-id test --keyring-backend test -y

#mint cmst
comdex tx tokenmint tokenmint 1 1 --from test --chain-id test --keyring-backend test -y