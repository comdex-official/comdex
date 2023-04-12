import os
from pathlib import Path

HOME_DIR = Path.home() 
COMDEX_DIR_PATH  = os.path.abspath(os.curdir)

NODE_MONIKER = "testdev"
CHAIN_ID = "test-1"
GENESIS_ACCOUNT_NAME = "cooluser"

GENESIS_MINTING_TEST_TOKENS = [
    "100000000000000000000000000stake",
    "100000000000000000000000000ucmdx",
    "100000000000000000000000000ucmst",
    "100000000000000000000000000ibc/C4CFF46FD6DE35CA4CF4CE031E643C8FDC9BA4B99AE598E9B0ED98FE3A2319F9", # ATOM
    "100000000000000000000000000ibc/05AC4BBA78C5951339A47DD1BC1E7FC922A9311DF81C85745B1C162F516FF2F1", # OSMO
    "100000000000000000000000000ibc/96CF88731A06654F8510473865C75200005B81A2A9FECB90389034F787015CA3", # AXLUSDC 
    "100000000000000000000000000ibc/065AD21492A7749E283A37AC469426562F988B7BA59712C0C545F381865A66BA", # AXLWETH
    "100000000000000000000000000ibc/50EF138042B553362774A8A9DE967F610E52CAEB3BA864881C9A1436DED98075", # AXLDAI
]

GENESIS_TOKENS = ",".join(GENESIS_MINTING_TEST_TOKENS)

VOTING_PERIOD_IN_SEC = 10
DEPOSIT_PERIOD_IN_SEC = 10