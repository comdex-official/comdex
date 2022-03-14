make all

rm -rf ~/.comdex

mkdir ~/.comdex

comdex init --chain-id test test
comdex keys add test --keyring-backend test
comdex add-genesis-account test 100000000000000stake,100000000000000ucmdx,100000000000000ucgold,100000000000000ucsilver,100000000000000ucoil,100000000000000ibc/4294C3DB67564CF4A0B2BFACC8415A59B38243F6FF9E288FBA34F9B4823BA16E,10000000000ibc/C4CFF46FD6DE35CA4CF4CE031E643C8FDC9BA4B99AE598E9B0ED98FE3A2319F9 --keyring-backend test
comdex gentx test 1000000000stake --chain-id test --keyring-backend test
comdex collect-gentxs
# Make sure to add the admin account address in params after running this script and before starting the chain.
# use this account to generate assets and asset pairs comdex14acvd72sw34x5kpdns658s2wwqwq99janhrd36
