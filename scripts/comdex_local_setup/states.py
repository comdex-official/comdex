APPS = [
    # [name, shortName, minGovDeposit, govTimeInSeconds]
    ["harbor", "hbr", 1000000, 5],              # ID - 1
    ["cswap", "cswap", 0, 0],                   # ID - 2
    ["commodo", "comdo", 0, 0],                 # ID - 3
]

ASSETS = [
    # [name, denom, isOnChain, assetOraclePriceRequired]
    ["ATOM", "uatom", 0, 1],                    # ID - 1
    ["CMDX", "ucmdx", 0, 1],                    # ID - 2
    ["CMST", "ucmst", 0, 0],                    # ID - 3
    ["OSMO", "uosmo", 0, 1],                    # ID - 4
    ["cATOM", "ucatom", 0, 0],                  # ID - 5
    ["cCMDX", "uccmdx", 0, 0],                  # ID - 6
    ["cCMST", "uccmst", 0, 0],                  # ID - 7
    ["cOSMO", "ucosmo", 0, 0],                  # ID - 8
    ["HARBOR", "uharbor", 1, 0],                # ID - 9
]

PAIRS = [
    # [assetID1, assetID2]
    [1, 3],                                     # ID - 1
    [2, 3],                                     # ID - 2
    [4, 3]                                      # ID - 3
]

LIQUIDITY_PAIRS = [
    # [appID, baseCoinDenom, quoteCoinDenom]
    [2, ASSETS[1][1], ASSETS[8][1]],            # ID - 1
    [2, ASSETS[1][1], ASSETS[2][1]]             # ID - 2
]

LIQUIDITY_POOLS = [
    # [appID, pairID, depositCoins]
    [2, 1, f"1000000000000{ASSETS[1][1]},2000000000000{ASSETS[8][1]}"],     # ID - 1
    [2, 2, f"1000000000000{ASSETS[1][1]},250000000000{ASSETS[2][1]}"]       # ID - 2
]

WASM_PROPOSALS = [
    # [proposalID, proposal]
    [
        1, 
        {
            "propose":{
                "propose": {
                    "title":"New proposal for add vault pair for CMDX C - CMST",
                    "description" :"This is a base execution proposal to add CMDX C - CMST vault pair with given Vault properties a. Liquidation ratio : 140 % b. Stability Fee : 1%  c. Liquidation Penalty : 12% d. DrawDown Fee : 1% e. Debt Cieling : 100000000 CMST f. Debt Floor : 100 CMST ",
                    "msgs" : [{
                        "msg_add_extended_pairs_vault":{
                            "app_id":1,
                            "pair_id":1,
                            "stability_fee":"0.025",
                            "closing_fee":"0.00",
                            "liquidation_penalty":"0.12",
                            "draw_down_fee":"0.001",
                            "is_vault_active":True,
                            "debt_ceiling":1000000000000,
                            "debt_floor":100000000,
                            "is_stable_mint_vault":False,
                            "min_cr":"1.7",
                            "pair_name":"ATOM-A",
                            "asset_out_oracle_price":False,
                            "asset_out_price":1000000,
                            "min_usd_value_left":100000
                        }
                    }],
                    "app_id_param" :1
                }
            }
        }
    ],
    [
        2,
        {
            "propose":{
                "propose":{
                    "title":"New proposal to initialise collector param for stability fee and auction (surplus and debt)threshold data",
                    "description" :"This is an base  execution proposal to initialise CMST and HARBOR pair for Surplus and dutch auction with Debt Threshold being  1000 CMST and Surplus Threshold as 100000000 CMST  ",
                    "msgs" : [{
                        "msg_set_collector_lookup_table":{
                            "app_id" :1,
                            "collector_asset_id" :3,
                            "secondary_asset_id" :9,
                            "surplus_threshold" : 10000000000,
                            "debt_threshold":10000000,
                            "locker_saving_rate":"0.06",
                            "lot_size":20000,
                            "bid_factor":"0.01",
                            "debt_lot_size":200
                        }
                    }],
                    "app_id_param" :1
                }
            }
        }
    ],
    [
        3,
        {
            "propose":{
                "propose": {
                    "title":"New proposal for setting auction params for collateral auctions.",
                    "description" :"This is an base proposal to set collateral auction params with auction duration being 600 seconds.",
                    "msgs" : [{
                        "msg_add_auction_params":{
                            "app_id": 1,
                            "auction_duration_seconds":20,
                            "buffer":"1.2",
                            "cusp":"0.4",
                            "step":360,
                            "price_function_type":1,
                            "surplus_id":1,
                            "debt_id":2,
                            "dutch_id":3,
                            "bid_duration_seconds":10
                        }
                    }],
                    "app_id_param" :1
                }
            }
        }
    ],
    [
        4,
        {
            "propose":{
                "propose": {
                    "title":"New proposal to whitelist CMST for locker",
                    "description" :"This is an base  execution proposal to add use CMST as locker deposit asset.",
                    "msgs" : [{
                        "msg_set_auction_mapping_for_app":{
                            "app_id" :1,
                            "asset_id" :[3],
                            "is_surplus_auction" :[False],
                            "is_debt_auction":[True],
                            "is_distributor":[True],
                            "asset_out_oracle_price":[False],
                            "asset_out_price":[1000000]
                        }
                    }],
                    "app_id_param" :1
                }
            }
        }
    ],
    [
        5,
        {
            "propose":{
                "propose":{
                    "title":"New proposal to whitelist CMST for locker",
                    "description" :"This is an base  execution proposal to add use CMST as locker deposit asset.",
                    "msgs" : [{
                        "msg_white_list_asset_locker":{
                            "app_id":1,
                            "asset_id":3
                        }
                    }],
                    "app_id_param" :1
                }
            }
        }
    ],
    [
        6,
        {
            "propose":{
                "propose": {
                    "title":"New proposal to whitelist CMST for locker",
                    "description" :"This is an base  execution proposal to add use CMST as locker deposit asset.",
                    "msgs" : [{
                        "msg_whitelist_app_id_locker_rewards":{
                            "app_id":1,
                            "asset_ids":[3]
                        }
                    }],
                    "app_id_param" :1
                }
            }
        }
    ],
    [
        7,
        {
            "propose":{
                "propose":{
                    "title":"New proposal to whitelist CMST for locker",
                    "description" :"This is an base  execution proposal to add use CMST as locker deposit asset.",
                    "msgs" : [{
                        "msg_whitelist_app_id_vault_interest":{
                            "app_id":1
                        }
                    }],
                    "app_id_param" :1
                }
            }
        }
    ],
    [
        8,
        {
            "propose":{
                "propose":{
                    "title":"New proposal for add pair for CMDX",
                    "description" :"This is an base proposal execution proposal to add CMDX-CMST n.",
                    "msgs" : [{
                        "msg_whitelist_app_id_liquidation":{
                            "app_id":1
                        }
                    }],
                    "app_id_param" :1
                }
            }
        }
    ],
    [
        9,
        {
            "propose":{
                "propose":{
                    "title":"New proposal for setting auction params for collateral auctions.",
                    "description" :"This is an base proposal to set collateral auction params with auction duration being 600 seconds.",
                    "msgs" : [{
                        "msg_add_e_s_m_trigger_params":{
                            "app_id": 1,
                            "target_value":{"amount":"200", "denom":"uharbor"},
                            "cool_off_period":60,
                            "asset_id": [3],
                            "rates": [1000000]
                        }
                    }],
                    "app_id_param" :1
                }
            }
        }
    ],
    [
        10,
        {
            "propose":{
                "propose":{
                    "title":"New proposal for add vault pair for CMDX C - CMST",
                    "description" :"This is a base execution proposal to add CMDX C - CMST vault pair with given Vault properties a. Liquidation ratio : 140 % b. Stability Fee : 1%  c. Liquidation Penalty : 12% d. DrawDown Fee : 1% e. Debt Cieling : 100000000 CMST f. Debt Floor : 100 CMST ",
                    "msgs" : [{
                        "msg_add_extended_pairs_vault":{
                            "app_id":1,
                            "pair_id":2,
                            "stability_fee":"0.025",
                            "closing_fee":"0.00",
                            "liquidation_penalty":"0.12",
                            "draw_down_fee":"0.001",
                            "is_vault_active":True,
                            "debt_ceiling":1000000000000,
                            "debt_floor":100000000,
                            "is_stable_mint_vault":False,
                            "min_cr":"1.7",
                            "pair_name":"CMDX-A",
                            "asset_out_oracle_price":False,
                            "asset_out_price":1000000,
                            "min_usd_value_left":100000
                        }
                    }],
                    "app_id_param" :1
                }
            }
        }
    ]
]