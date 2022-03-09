make all

rm -rf ~/.comdex

mkdir ~/.comdex

comdex init --chain-id test test
comdex keys add test --recover<<<"y
wage thunder live sense resemble foil apple course spin horse glass mansion midnight laundry acoustic rhythm loan scale talent push green direct brick please"
comdex add-genesis-account test 100000000000000stake
comdex gentx test 1000000000stake --chain-id test
comdex collect-gentxs
# Make sure to add the admin account address in params after running this script and before starting the chain.
# use this account to generate assets and asset pairs comdex14acvd72sw34x5kpdns658s2wwqwq99janhrd36