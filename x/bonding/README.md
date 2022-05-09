# Commands

```sh
# 1 day 100stake lock-tokens command
comdex tx bonding lock-tokens 200stake --duration="86400s" --from=validator --chain-id=testing --keyring-backend=test --yes

# 5s 100stake lock-tokens command
comdex tx bonding lock-tokens 100stake --duration="5s" --from=validator --chain-id=testing --keyring-backend=test --yes

# begin unlock tokens, NOTE: add more gas when unlocking more than two locks in a same command
comdex tx bonding begin-unlock-tokens --from=validator --gas=500000 --chain-id=testing --keyring-backend=test --yes

# unlock tokens, NOTE: add more gas when unlocking more than two locks in a same command
comdex tx bonding unlock-tokens --from=validator --gas=500000 --chain-id=testing --keyring-backend=test --yes

# unlock specific period lock
comdex tx bonding unlock-by-id 1 --from=validator --chain-id=testing --keyring-backend=test --yes

# account balance
comdex query bank balances $(comdex keys show -a validator --keyring-backend=test)

# query module balance
comdex query bonding module-balance

# query locked amount
comdex query bonding module-locked-amount

# query lock by id
comdex query bonding lock-by-id 1

# query account unlockable coins
comdex query bonding account-unlockable-coins $(comdex keys show -a validator --keyring-backend=test)

# query account locks by denom past time
comdex query bonding account-locked-pasttime-denom $(comdex keys show -a validator --keyring-backend=test) 1611879610 stake

# query account locks past time
comdex query bonding account-locked-pasttime $(comdex keys show -a validator --keyring-backend=test) 1611879610

# query account locks by denom with longer duration
comdex query bonding account-locked-longer-duration-denom $(comdex keys show -a validator --keyring-backend=test) 5.1s stake

# query account locks with longer duration
comdex query bonding account-locked-longer-duration $(comdex keys show -a validator --keyring-backend=test) 5.1s

# query account locked coins
comdex query bonding account-locked-coins $(comdex keys show -a validator --keyring-backend=test)

# query account locks before time
comdex query bonding account-locked-beforetime $(comdex keys show -a validator --keyring-backend=test) 1611879610
```