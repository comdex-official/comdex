import subprocess
import requests
import toml
import threading
import datetime
import os
import json
import time

from constants import *
from states import *

def SetupNewChain():
    subprocess.getstatusoutput("rm -rf ~/.comdex")
    print("deleted previos state")
    print("rebuilding binaries.....")
    response = subprocess.getstatusoutput(f"make install --directory={COMDEX_DIR_PATH}")
    print(response[1])
    print("binary re-build done ✔️")
    subprocess.getstatusoutput("sudo mv ~/go/bin/comdex /usr/local/bin")
    subprocess.getstatusoutput(f"comdex init {NODE_MONIKER} --chain-id {CHAIN_ID}")
    subprocess.getstatusoutput(f"comdex keys add {GENESIS_ACCOUNT_NAME} --keyring-backend test")
    subprocess.getstatusoutput(f"comdex add-genesis-account $(comdex keys show cooluser --keyring-backend test -a) {GENESIS_TOKENS}")
    subprocess.getstatusoutput("comdex gentx cooluser 1000000000stake --chain-id test-1 --keyring-backend test")
    subprocess.getstatusoutput("comdex collect-gentxs")
    print("chain initialization done ✔️")

    with open(f"{HOME_DIR}/.comdex/config/genesis.json", "r") as jsonFile:
        data = json.load(jsonFile)

    data["app_state"]["gov"]["deposit_params"]["min_deposit"][0]["denom"] = "ucmdx"
    data["app_state"]["gov"]["deposit_params"]["max_deposit_period"] = str(DEPOSIT_PERIOD_IN_SEC)+"s"
    data["app_state"]["gov"]["voting_params"]["voting_period"] = str(VOTING_PERIOD_IN_SEC)+"s"
    data["app_state"]["gov"]["tally_params"]["quorum"] = "0"
    data["app_state"]["gov"]["tally_params"]["threshold"] = "0"
    data["app_state"]["gov"]["tally_params"]["veto_threshold"] = "0"

    with open(f"{HOME_DIR}/.comdex/config/genesis.json", "w") as jsonFile:
        json.dump(data, jsonFile)
    
    print("genesis update done ✔️")

    rpcConfig = toml.load(f"{HOME_DIR}/.comdex/config/config.toml")
    rpcConfig["rpc"]["laddr"]="tcp://0.0.0.0:26657" 
    rpcConfig["rpc"]["cors_allowed_origins"]= ["*"]

    with open(f"{HOME_DIR}/.comdex/config/config.toml",'w') as f:
        toml.dump(rpcConfig, f)
        f.close()
    
    print("RPC configurations done ✔️")
    
    lcdConfig = toml.load(f"{HOME_DIR}/.comdex/config/app.toml")
    lcdConfig["api"]["enable"]=True 
    lcdConfig["api"]["enabled-unsafe-cors"]= True

    with open(f"{HOME_DIR}/.comdex/config/app.toml",'w') as f:
        toml.dump(lcdConfig, f)
        f.close()
    print("LCD configurations done ✔️")
    print("Test Chain Setup done ✔️")
    print("Waiting for the chain to be started automatically...")
    

def StartChain():
    command = "comdex start"
    subprocess.getstatusoutput(command)

def StartChainIndicator():
    while True:
        output = subprocess.getstatusoutput("lsof -i tcp:26657")[1]
        if output:
            break
    print("chain start  done ✔️")
    time.sleep(6)

def GetLatestPropID():
    proposals = requests.get("http://127.0.0.1:1317/cosmos/gov/v1beta1/proposals").json()
    return int(proposals["proposals"][-1]["proposal_id"])

def GetGenesisAccAddress():
    command = f"comdex keys show {GENESIS_ACCOUNT_NAME} --keyring-backend test --output json"
    output = subprocess.getstatusoutput(command)[1]
    output = json.loads(output)
    return output["address"]

def Vote(option):
    if option not in ["yes", "no"]:
        exit("Invalid voting option")
    latestPropID = GetLatestPropID()
    command = f"comdex tx gov vote {latestPropID} {option} --from {GENESIS_ACCOUNT_NAME} --chain-id {CHAIN_ID} --keyring-backend test -y"
    output = subprocess.getstatusoutput(command)[1]
    output = json.loads(output)
    if int(output["code"]) != 0:
        print(output)
        exit(f"error while voting on prop {latestPropID}")
    print(f"Vote submitted on Prop {latestPropID} ✔️")

def AddApp(name, shortName, minGovDeposit=0, govTimeInSeconds=0):
    command = f"""comdex tx gov submit-proposal add-app {name} {shortName} {minGovDeposit} {govTimeInSeconds} --title "New App" --description "Adding new app on comdex" --deposit 10000000ucmdx --from {GENESIS_ACCOUNT_NAME} --chain-id {CHAIN_ID} --keyring-backend test -y"""
    output = subprocess.getstatusoutput(command)[1]
    output = json.loads(output)
    if int(output["code"]) != 0:
        print(output)
        exit("error in add app prop")
    print(f"New App {name} Proposal  Submitted ✔️")

def AddAsset(name, denom, isOnChain=1, assetOraclePriceRequired=1):
    jsonData = {
        "name" : name,
        "denom" : denom,
        "decimals" :"1000000",
        "is_on_chain" :str(isOnChain),
        "asset_oracle_price" :str(assetOraclePriceRequired),
        "title" :"Add assets for applications to be deployed on comdex chain",
        "description" :f"This proposal it to add asset {name} to be then used on harbor, commodo and cswap apps",
        "deposit" :"1000000000ucmdx"
    }
    fileName = f"newAsset-{name}-{datetime.datetime.now()}.json"
    with open(fileName, "w") as jsonFile:
        json.dump(jsonData, jsonFile)
    
    command = f"""comdex tx gov submit-proposal add-assets --add-assets-file "{fileName}" --from {GENESIS_ACCOUNT_NAME} --chain-id {CHAIN_ID} --keyring-backend test -y"""
    output = subprocess.getstatusoutput(command)[1]
    output = json.loads(output)
    if int(output["code"]) != 0:
        print(output)
        exit("error in add asset prop")
    if os.path.exists(fileName):
        os.remove(fileName)
    print(f"New Asset {name} Proposal  Submitted ✔️")

def AddPair(assetID1, assetID2):
    command = f"""comdex tx gov submit-proposal add-pairs {assetID1} {assetID2}  --title "New Pair" --description "Adding new pair" --deposit 10000000ucmdx --from {GENESIS_ACCOUNT_NAME} --chain-id {CHAIN_ID} --keyring-backend test -y"""
    output = subprocess.getstatusoutput(command)[1]
    output = json.loads(output)
    if int(output["code"]) != 0:
        print(output)
        exit("error in add pairs prop")
    print(f"New Pair ({assetID1}, {assetID2}) Proposal  Submitted ✔️")

def MintToken(appID, assetID):
    print("Minting token for previosly added asset in app..")
    command = f"comdex tx tokenmint tokenmint {appID} {assetID} --from {GENESIS_ACCOUNT_NAME} --chain-id {CHAIN_ID} --keyring-backend test -y"
    output = subprocess.getstatusoutput(command)[1]
    output = json.loads(output)
    if int(output["code"]) != 0:
        print(output)
        exit("error whle minting tokens")
    print(f"Token Minting Done For AssetID {assetID} in App {appID} ✔️")

def AddAssetInAppsAndVote(appID, assetID):
    jsonData = {
        "app_id" : str(appID),
        "asset_id" : str(assetID),
        "genesis_supply" : "1000000000000000",
        "is_gov_token" : "1",
        "recipient" : GetGenesisAccAddress(),
        "title" : "fdv",
        "description" : "dffd",
        "deposit" : "100000000ucmdx"
    }
    fileName = f"assetMap-{datetime.datetime.now()}.json"
    with open(fileName, "w") as jsonFile:
        json.dump(jsonData, jsonFile)

    command = f"""comdex tx gov submit-proposal add-asset-in-app --add-asset-mapping-file "{fileName}" --from {GENESIS_ACCOUNT_NAME} --chain-id {CHAIN_ID} --keyring-backend test -y"""
    output = subprocess.getstatusoutput(command)[1]
    output = json.loads(output)
    if int(output["code"]) != 0:
        print(output)
        exit("error in add asset in app prop")
    print(f"New Add Asset In App (appID - {appID}, assetID - {assetID}) Proposal  Submitted ✔️")
    if os.path.exists(fileName):
        os.remove(fileName)
    Vote("yes")
    print("waiting for prop to pass")
    time.sleep(VOTING_PERIOD_IN_SEC)
    MintToken(appID, assetID)

def CreateLiquidityPair(appID, baseCoinDenom, quoteCoinDenom):
    command = f"""comdex tx gov submit-proposal create-liquidity-pair {appID} {baseCoinDenom} {quoteCoinDenom}  --title "New Liquidity Pair" --description "Adding new liquidity pair" --deposit 10000000ucmdx --from {GENESIS_ACCOUNT_NAME} --chain-id {CHAIN_ID} --keyring-backend test -y"""
    output = subprocess.getstatusoutput(command)[1]
    output = json.loads(output)
    if int(output["code"]) != 0:
        print(output)
        exit("error in add pairs prop")
    print(f"New Liquidity Pair ({baseCoinDenom}, {quoteCoinDenom}) Proposal  Submitted ✔️")

def CreateLiquidityPool(appID, pairID, depositCoins):
    command = f"comdex tx liquidityV1 create-pool {appID} {pairID} {depositCoins} --from {GENESIS_ACCOUNT_NAME} --chain-id {CHAIN_ID} --gas 5000000 --keyring-backend test -y"
    output = subprocess.getstatusoutput(command)[1]
    output = json.loads(output)
    if int(output["code"]) != 0:
        print(output)
        exit("error in create pool")
    print(f"New liquidity pool created for pairID {pairID} in app {appID} with initial deposit of {depositCoins} ✔️")

def StoreAndIntantiateGovernanceWasmContract():
    command = f"comdex tx wasm store {GOVERNANCE_CONTRACT_PATH} --from {GENESIS_ACCOUNT_NAME}  --chain-id {CHAIN_ID} --gas 5000000 --gas-adjustment 1.3 --keyring-backend test  -y  --output json"
    output = subprocess.getstatusoutput(command)[1]
    output = json.loads(output)
    if int(output["code"]) != 0:
        print(output)
        exit("error in adding governance wasm contract")
    print(f"Governance wasm contract added successfully ✔️")

    init = {
        "threshold":{
            "threshold_quorum":{
                "threshold":"0.50",
                "quorum":"0.33"
            }
        },
        "target":"0.0.0.0:9090"
    }
    command = f"""comdex tx wasm instantiate 1 '{json.dumps(init)}' --label "Instantiate Contract Governance" --no-admin --from {GENESIS_ACCOUNT_NAME} --chain-id {CHAIN_ID} --gas 5000000 --gas-adjustment 1.3 --keyring-backend test -y"""
    output = subprocess.getstatusoutput(command)[1]
    output = json.loads(output)
    if int(output["code"]) != 0:
        print(output)
        exit("error in instantiating governance wasm contract")
    print(f"Governance wasm contract instantiated successfully ✔️")

def ExecuteWasmGovernanceProposal(proposalID):
    execute = {
        "execute": {
            "proposal_id":proposalID
        }
    }
    command = f"""comdex tx wasm execute {GOVERNANCE_CONTRACT_ADDRESS} '{json.dumps(execute)}' --from {GENESIS_ACCOUNT_NAME} --chain-id {CHAIN_ID} --gas 5000000 --keyring-backend test -y"""
    output = subprocess.getstatusoutput(command)[1]
    output = json.loads(output)
    if int(output["code"]) != 0:
        print(output)
        exit(f"error while executing wasm prop with id {proposalID} ")
    print(f"Proposal with ID {proposalID} executed successfully ✔️")

def ProposeWasmGovernanceProposal(proposal, proposlID):
    command = f"""comdex tx wasm execute {GOVERNANCE_CONTRACT_ADDRESS}  '{json.dumps(proposal)}' --amount 100000000uharbor --from {GENESIS_ACCOUNT_NAME}  --chain-id {CHAIN_ID} --gas 5000000 --keyring-backend test -y"""
    output = subprocess.getstatusoutput(command)[1]
    output = json.loads(output)
    if int(output["code"]) != 0:
        print(output)
        exit(f"error while proposing prop {proposlID} ")
    print(f"Proposal {proposlID} raised successfully ✔️")


def CreateState():

    for app in APPS:
        if len(app) != 4:
            exit("Invalid app configs")
        AddApp(app[0], app[1], app[2], app[3])
        Vote("yes")
    
    for asset in ASSETS:
        if len(asset) != 4:
            exit("Invalid asset configs")
        AddAsset(asset[0], asset[1], asset[2], asset[3])
        Vote("yes")
    
    for pair in PAIRS:
        if len(pair) != 2:
            exit("Invalid pairs configs")
        AddPair(pair[0], pair[1])
        Vote("yes")
    
    AddAssetInAppsAndVote(1, 9)
    StoreAndIntantiateGovernanceWasmContract()

    for wasmProp in WASM_PROPOSALS:
        propID = wasmProp[0]
        prop = wasmProp[1]
        ProposeWasmGovernanceProposal(prop, propID)
        print(f"waiting for wasm prop {propID}")
        time.sleep(APPS[0][3]) # waiting for proposal duration
        ExecuteWasmGovernanceProposal(propID)

    for liquidityPair in LIQUIDITY_PAIRS:
        if len(liquidityPair) != 3:
            exit("Invalid liquidity pair configs")
        CreateLiquidityPair(liquidityPair[0], liquidityPair[1], liquidityPair[2])
        Vote("yes")
    
    for liquidityPool in LIQUIDITY_POOLS:
        if len(liquidityPool) != 3:
            exit("Invalid liquidity pool configs")
        CreateLiquidityPool(liquidityPool[0], liquidityPool[1], liquidityPool[2])

def main():
    if not os.path.exists(HOME_DIR):
        exit(f"Error - root dir not found {HOME_DIR}")
    if not os.path.exists(COMDEX_DIR_PATH):
        exit(f"Error - invalid comdex repo path {COMDEX_DIR_PATH}")
    SetupNewChain()
    thr = threading.Thread(target=StartChain, args=(), kwargs={})
    thr.start()
    StartChainIndicator()
    CreateState()
    print("Press Ctr+C to stop the chain")

main()