import os
from pathlib import Path

HOME_DIR = Path.home() 
COMDEX_DIR_PATH  = os.path.abspath(os.curdir)

NODE_MONIKER = "testdev"
CHAIN_ID = "test-1"
GENESIS_ACCOUNT_NAME = "cooluser"
GENESIS_TOKENS = "1000000000000000000000stake,100000000000000000000000ucmdx"

VOTING_PERIOD_IN_SEC = 10
DEPOSIT_PERIOD_IN_SEC = 10

GOVERNANCE_CONTRACT_PATH = f"{COMDEX_DIR_PATH}/scripts/comdex_local_setup/governance.wasm" 
GOVERNANCE_CONTRACT_ADDRESS = "comdex14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9spunaxy"