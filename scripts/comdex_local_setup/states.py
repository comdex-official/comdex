from constants import *

APPS = [
    # [name, shortName, minGovDeposit, govTimeInSeconds]
    ["harbor", "hbr", 1000000, 5],  # ID - 1
    ["cswap", "cswap", 0, 0],  # ID - 2
    ["commodo", "comdo", 0, 0],  # ID - 3
]

ASSETS = [
    # [name, denom, decimals, isOnChain, assetOraclePriceRequired, isCdpMintable]
    ["ATOM", "uatom", 1000000, 0, 1, 0],  # ID - 1
    ["PETRI", "upetri",1000000, 0, 1, 0],  # ID - 2
    ["FUST", "ufust",1000000, 1, 1, 1],  # ID - 3
    ["OSMO", "uosmo",1000000, 0, 1, 0],  # ID - 4
    ["sATOM", "usatom",1000000, 1, 0, 1],  # ID - 5
    ["sPETRI", "uspetri",1000000, 1, 0, 1],  # ID - 6
    ["sFUST", "usfust",1000000, 1, 0, 1],  # ID - 7
    ["sOSMO", "ucosmo",1000000, 1, 0, 1],  # ID - 8
    ["HARBOR", "uharbor",1000000, 1, 1, 0],  # ID - 9
    ["WETH", "weth-wei",1000000000000000000, 0, 1, 0],  # ID - 10
    ["CANTO", "ucant",10000000000000000000000000, 0, 1, 0],  # ID - 11
    ["sGOLD", "usgold",1000000, 1, 1, 1],  # ID - 12
    ["USDC", "usdc",1000000, 0, 1, 0],  # ID - 13

]

PAIRS = [
    # [assetID1, assetID2]
    [1, 3],  # ID - 1
    [2, 3],  # ID - 2
    [4, 3],  # ID - 3
    [10, 3],  # ID - 4
    [2, 12],  # ID - 5
    [13,3], # ID - 6

]

LIQUIDITY_PAIRS = [
    # [appID, baseCoinDenom, quoteCoinDenom]
    [2, ASSETS[1][1], ASSETS[8][1]],  # ID - 1
    [2, ASSETS[1][1], ASSETS[2][1]],  # ID - 2
]

LIQUIDITY_POOLS = [
    # [appID, pairID, depositCoins]
    [2, 1, f"1000000000000{ASSETS[1][1]},2000000000000{ASSETS[8][1]}"],  # ID - 1
]

ADD_ASSET_RATES = [
    # [ assetName, jsonObject]
    [
        "FUST",
        {
            "asset_id": "3",
            "u_optimal": "0.8",
            "base": "0.002",
            "slope_1": "0.06",
            "slope_2": "0.6",
            "enable_stable_borrow": "1",
            "stable_base": "0.04",
            "stable_slope_1": "0.04",
            "stable_slope_2": "0.6",
            "ltv": "0.8",
            "liquidation_threshold": "0.85",
            "liquidation_penalty": "0.025",
            "liquidation_bonus": "0.025",
            "reserve_factor": "0.1",
            "c_asset_id": "7",
            "title": "Add Asset Rates Stats FUST",
            "description": "adding asset rates stats for FUST",
            "deposit": "10000000upetri",
        },
    ],
    [
        "ATOM",
        {
            "asset_id": "1",
            "u_optimal": "0.75",
            "base": "0.002",
            "slope_1": "0.07",
            "slope_2": "1.25",
            "enable_stable_borrow": "0",
            "stable_base": "0.0",
            "stable_slope_1": "0.0",
            "stable_slope_2": "0.0",
            "ltv": "0.7",
            "liquidation_threshold": "0.75",
            "liquidation_penalty": "0.05",
            "liquidation_bonus": "0.05",
            "reserve_factor": "0.2",
            "c_asset_id": "5",
            "title": "Add Asset Rates Stats ATOM",
            "description": "adding asset rates stats ATOM",
            "deposit": "10000000upetri",
        },
    ],
    [
        "OSMO",
        {
            "asset_id": "4",
            "u_optimal": "0.65",
            "base": "0.002",
            "slope_1": "0.08",
            "slope_2": "1.5",
            "enable_stable_borrow": "0",
            "stable_base": "0.0",
            "stable_slope_1": "0.0",
            "stable_slope_2": "0.0",
            "ltv": "0.6",
            "liquidation_threshold": "0.65",
            "liquidation_penalty": "0.05",
            "liquidation_bonus": "0.05",
            "reserve_factor": "0.2",
            "c_asset_id": "8",
            "title": "Add Asset Rates Stats OSMO",
            "description": "adding asset rates stats OSMO",
            "deposit": "10000000upetri",
        },
    ],
    [
        "PETRI",
        {
            "asset_id": "2",
            "u_optimal": "0.5",
            "base": "0.002",
            "slope_1": "0.08",
            "slope_2": "2.0",
            "enable_stable_borrow": "0",
            "stable_base": "0.0",
            "stable_slope_1": "0.0",
            "stable_slope_2": "0.0",
            "ltv": "0.5",
            "liquidation_threshold": "0.55",
            "liquidation_penalty": "0.05",
            "liquidation_bonus": "0.05",
            "reserve_factor": "0.2",
            "c_asset_id": "6",
            "title": "Add Asset Rates Stats PETRI",
            "description": "adding asset rates stats PETRI",
            "deposit": "10000000upetri",
        },
    ],
]

ADD_LEND_POOL = [
    {
        "module_name": "cmdx",
        "asset_id": "1,2,3",
        "asset_transit_type": "3,1,2",
        "supply_cap": "5000000000000000000,1000000000000000000,5000000000000000000",
        "c_pool_name": "PETRI-ATOM-FUST",
        "reserve_funds": "100000000",
        "title": "Add pool",
        "description": "adding pool",
        "deposit": "10000000upetri",
    },
    {
        "module_name": "osmo",
        "asset_id": "1,4,3",
        "asset_transit_type": "3,1,2",
        "supply_cap": "5000000000000000000,3000000000000000000,5000000000000000000",
        "c_pool_name": "OSMO-ATOM-FUST",
        "reserve_funds": "100000000",
        "title": "Add pool",
        "description": "adding pool",
        "deposit": "10000000upetri",
    },
]

ADD_LEND_PAIR = [
    [
        "PETRI-FUST",
        {
            "asset_in": "2",
            "asset_out": "3",
            "is_inter_pool": "0",
            "asset_out_pool_id": "1",
            "min_usd_value_left": "100000000000",
            "title": "Add Extended pair PETRI-FUST",
            "description": "adding extended pairs for PETRI-FUST same pool",
            "deposit": "10000000upetri",
        },
    ],
    [
        "PETRI-ATOM",
        {
            "asset_in": "2",
            "asset_out": "1",
            "is_inter_pool": "0",
            "asset_out_pool_id": "1",
            "min_usd_value_left": "100000000000",
            "title": "Add Extended pair PETRI-ATOM",
            "description": "adding extended pairs PETRI-ATOM same pool",
            "deposit": "10000000upetri",
        },
    ],
    [
        "ATOM-PETRI",
        {
            "asset_in": "1",
            "asset_out": "2",
            "is_inter_pool": "0",
            "asset_out_pool_id": "1",
            "min_usd_value_left": "100000000000",
            "title": "Add Extended pair ATOM-PETRI",
            "description": "adding extended pairs ATOM-PETRI same pool",
            "deposit": "10000000upetri",
        },
    ],
    [
        "ATOM-FUST",
        {
            "asset_in": "1",
            "asset_out": "3",
            "is_inter_pool": "0",
            "asset_out_pool_id": "1",
            "min_usd_value_left": "100000000000",
            "title": "Add Extended pair ATOM-FUST",
            "description": "adding extended pairs ATOM-FUST same pool",
            "deposit": "10000000upetri",
        },
    ],
    [
        "FUST-PETRI",
        {
            "asset_in": "3",
            "asset_out": "2",
            "is_inter_pool": "0",
            "asset_out_pool_id": "1",
            "min_usd_value_left": "100000000000",
            "title": "Add Extended pair FUST-PETRI",
            "description": "adding extended pairs FUST-PETRI same pool",
            "deposit": "10000000upetri",
        },
    ],
    [
        "FUST-ATOM",
        {
            "asset_in": "3",
            "asset_out": "1",
            "is_inter_pool": "0",
            "asset_out_pool_id": "1",
            "min_usd_value_left": "100000000000",
            "title": "Add Extended pair FUST-ATOM",
            "description": "adding extended pairs FUST-ATOM same pool",
            "deposit": "10000000upetri",
        },
    ],
    [
        "OSMO-FUST",
        {
            "asset_in": "4",
            "asset_out": "3",
            "is_inter_pool": "0",
            "asset_out_pool_id": "2",
            "min_usd_value_left": "100000000000",
            "title": "Add Extended pair OSMO-FUST",
            "description": "adding extended pairs OSMO-FUST same pool",
            "deposit": "10000000upetri",
        },
    ],
    [
        "OSMO-ATOM",
        {
            "asset_in": "4",
            "asset_out": "1",
            "is_inter_pool": "0",
            "asset_out_pool_id": "2",
            "min_usd_value_left": "100000000000",
            "title": "Add Extended pair OSMO-ATOM",
            "description": "adding extended pairs OSMO-ATOM same pool",
            "deposit": "10000000upetri",
        },
    ],
    [
        "ATOM-OSMO",
        {
            "asset_in": "1",
            "asset_out": "4",
            "is_inter_pool": "0",
            "asset_out_pool_id": "2",
            "min_usd_value_left": "100000000000",
            "title": "Add Extended pair ATOM-OSMO",
            "description": "adding extended pairs ATOM-OSMO same pool",
            "deposit": "10000000upetri",
        },
    ],
    [
        "ATOM-FUST",
        {
            "asset_in": "1",
            "asset_out": "3",
            "is_inter_pool": "0",
            "asset_out_pool_id": "2",
            "min_usd_value_left": "100000000000",
            "title": "Add Extended pair ATOM-FUST",
            "description": "adding extended pairs ATOM-FUST same pool",
            "deposit": "10000000upetri",
        },
    ],
    [
        "FUST-OSMO",
        {
            "asset_in": "3",
            "asset_out": "4",
            "is_inter_pool": "0",
            "asset_out_pool_id": "2",
            "min_usd_value_left": "100000000000",
            "title": "Add Extended pair FUST-OSMO",
            "description": "adding extended pairs FUST-OSMO same pool",
            "deposit": "10000000upetri",
        },
    ],
    [
        "FUST-ATOM",
        {
            "asset_in": "3",
            "asset_out": "1",
            "is_inter_pool": "0",
            "asset_out_pool_id": "2",
            "min_usd_value_left": "100000000000",
            "title": "Add Extended pair FUST-ATOM",
            "description": "adding extended pairs FUST-ATOM same pool",
            "deposit": "10000000upetri",
        },
    ],
    [
        "PETRI-OSMO",
        {
            "asset_in": "2",
            "asset_out": "4",
            "is_inter_pool": "1",
            "asset_out_pool_id": "2",
            "min_usd_value_left": "100000000000",
            "title": "Add Extended pair PETRI-OSMO",
            "description": "adding extended pairs PETRI-OSMO cross pool",
            "deposit": "10000000upetri",
        },
    ],
    [
        "FUST-OSMO",
        {
            "asset_in": "3",
            "asset_out": "4",
            "is_inter_pool": "1",
            "asset_out_pool_id": "2",
            "min_usd_value_left": "100000000000",
            "title": "Add Extended pair FUST-OSMO",
            "description": "adding extended pairs FUST-OSMO cross pool",
            "deposit": "10000000upetri",
        },
    ],
    [
        "ATOM-OSMO",
        {
            "asset_in": "1",
            "asset_out": "4",
            "is_inter_pool": "1",
            "asset_out_pool_id": "2",
            "min_usd_value_left": "100000000000",
            "title": "Add Extended pair ATOM-OSMO",
            "description": "adding extended pairs ATOM-OSMO cross pool",
            "deposit": "10000000upetri",
        },
    ],
    [
        "OSMO-PETRI",
        {
            "asset_in": "4",
            "asset_out": "2",
            "is_inter_pool": "1",
            "asset_out_pool_id": "1",
            "min_usd_value_left": "100000000000",
            "title": "Add Extended pair OSMO-PETRI",
            "description": "adding extended pairs OSMO-PETRI cross pool",
            "deposit": "10000000upetri",
        },
    ],
    [
        "FUST-PETRI",
        {
            "asset_in": "3",
            "asset_out": "2",
            "is_inter_pool": "1",
            "asset_out_pool_id": "1",
            "min_usd_value_left": "100000000000",
            "title": "Add Extended pair FUST-PETRI",
            "description": "adding extended pairs FUST-PETRI cross pool",
            "deposit": "10000000upetri",
        },
    ],
    [
        "ATOM-PETRI",
        {
            "asset_in": "1",
            "asset_out": "2",
            "is_inter_pool": "1",
            "asset_out_pool_id": "1",
            "min_usd_value_left": "100000000000",
            "title": "Add Extended pair ATOM-PETRI",
            "description": "adding extended pairs ATOM-PETRI cross pool",
            "deposit": "10000000upetri",
        },
    ],
]

LEND_ASSET_PAIR_MAPPING = [
    # [assetID, poolID, pairIDs]
    [1, 1, [3, 4, 15]],
    [2, 1, [1, 2, 13]],
    [3, 1, [5, 6, 14]],
    [4, 2, [7, 8, 16]],
    [1, 2, [9, 10, 18]],
    [3, 2, [11, 12, 17]],
]

WASM_CONTRACTS = [
    {
        "name": "Vesting Contract",
        "contractAddressKey": "vesting_contract",
        "contractLink": "https://github.com/petri-official/test-wasm-artifacts/raw/main/token_vesting.wasm",
        "contractPath": f"./token_vesting.wasm",
        "initator": {},
        "formatKeys": []
    },
    {
        "name": "Locking Contract",
        "contractAddressKey": "locking_contract",
        "contractLink": "https://github.com/petri-official/test-wasm-artifacts/raw/main/locking_contract.wasm",
        "contractPath": f"./locking_contract.wasm",
        "initator": {
            "t1": {"period": 500, "weight": "0.25"},
            "t2": {"period": 1000, "weight": "0.50"},
            "voting_period": 122500,
            "vesting_contract": "",
            "foundation_addr": ["petri1rljg3wwgv6qezu3p05vxny9pwk3mdwl0ja407z"],
            "foundation_percentage": "0.2",
            "surplus_asset_id": 3,
            "emission": {
                "app_id": 1,
                "total_rewards": "10000000000000",
                "rewards_pending": "10000000000000",
                "emission_rate": "0.01",
                "distributed_rewards": "0",
            },
            "min_lock_amount" : "4",
            "admin":"petri1663kc7kwlqxg5s35wuq4nleuqvy5j2tstlkeg2"
        },
        "formatKeys": ['vesting_contract']
    },
    {
        "name": "Governance Contract",
        "contractAddressKey": "governance_contract",
        "contractLink": "https://github.com/petri-official/test-wasm-artifacts/raw/main/governance.wasm",
        "contractPath": f"./governance.wasm",
        "initator": {
            "threshold": {"threshold_quorum": {"threshold": "0.50", "quorum": "0.33"}},
            "target": "0.0.0.0:9090",
            "locking_contract": "",
        },
        "formatKeys": ['locking_contract']
    },
]

WASM_PROPOSALS = [
    {
        "proposalID": 0,
        "isProposal": False,
        "contractAddressKey": "locking_contract",
        "content": {"lock": {"app_id": 1, "locking_period": "t2"}},
    },
    {
        "proposalID": 1,
        "isProposal": True,
        "contractAddressKey": "governance_contract",
        "content": {
            "propose": {
                "propose": {
                    "title": "New proposal for add vault pair for PETRI C - FUST",
                    "description": "This is a base execution proposal to add PETRI C - FUST vault pair with given Vault properties a. Liquidation ratio : 140 % b. Stability Fee : 1%  c. Liquidation Penalty : 12% d. DrawDown Fee : 1% e. Debt Cieling : 100000000 FUST f. Debt Floor : 100 FUST ",
                    "msgs": [
                        {
                            "msg_add_extended_pairs_vault": {
                                "app_id": 1,
                                "pair_id": 1,
                                "stability_fee": "0.025",
                                "closing_fee": "0.00",
                                "liquidation_penalty": "0.12",
                                "draw_down_fee": "0.001",
                                "is_vault_active": True,
                                "debt_ceiling": "1000000000000000000",
                                "debt_floor": "100000000",
                                "is_stable_mint_vault": False,
                                "min_cr": "1.7",
                                "pair_name": "ATOM-A",
                                "asset_out_oracle_price": False,
                                "asset_out_price": 1000000,
                                "min_usd_value_left": 100000,
                            }
                        }
                    ],
                    "app_id_param": 1,
                }
            }
        },
    },
    {
        "proposalID": 2,
        "isProposal": True,
        "contractAddressKey": "governance_contract",
        "content": {
            "propose": {
                "propose": {
                    "title": "New proposal to initialise collector param for stability fee and auction (surplus and debt)threshold data",
                    "description": "This is an base  execution proposal to initialise FUST and HARBOR pair for Surplus and dutch auction with Debt Threshold being  1000 FUST and Surplus Threshold as 100000000 FUST  ",
                    "msgs": [
                        {
                            "msg_set_collector_lookup_table": {
                                "app_id": 1,
                                "collector_asset_id": 3,
                                "secondary_asset_id": 9,
                                "surplus_threshold": "1000000000000",
                                "debt_threshold": "10000000",
                                "locker_saving_rate": "0.06",
                                "lot_size": "20000",
                                "bid_factor": "0.01",
                                "debt_lot_size": "200",
                            }
                        }
                    ],
                    "app_id_param": 1,
                }
            }
        },
    },
    {
        "proposalID": 3,
        "isProposal": True,
        "contractAddressKey": "governance_contract",
        "content": {
            "propose": {
                "propose": {
                    "title": "New proposal for setting auction params for collateral auctions.",
                    "description": "This is an base proposal to set collateral auction params with auction duration being 600 seconds.",
                    "msgs": [
                        {
                            "msg_add_auction_params": {
                                "app_id": 1,
                                "auction_duration_seconds": 3600,
                                "buffer": "1.2",
                                "cusp": "0.85",
                                "step": 360,
                                "price_function_type": 1,
                                "surplus_id": 1,
                                "debt_id": 2,
                                "dutch_id": 3,
                                "bid_duration_seconds": 1800,
                            }
                        }
                    ],
                    "app_id_param": 1,
                }
            }
        },
    },
    {
        "proposalID": 4,
        "isProposal": True,
        "contractAddressKey": "governance_contract",
        "content": {
            "propose": {
                "propose": {
                    "title": "New proposal to whitelist FUST for locker",
                    "description": "This is an base  execution proposal to add use FUST as locker deposit asset.",
                    "msgs": [
                        {
                            "msg_set_auction_mapping_for_app": {
                                "app_id": 1,
                                "asset_id": 3,
                                "is_surplus_auction": False,
                                "is_debt_auction": False,
                                "is_distributor": False,
                                "asset_out_oracle_price": False,
                                "asset_out_price": 1000000,
                            }
                        }
                    ],
                    "app_id_param": 1,
                }
            }
        },
    },
    {
        "proposalID": 5,
        "isProposal": True,
        "contractAddressKey": "governance_contract",
        "content": {
            "propose": {
                "propose": {
                    "title": "New proposal to whitelist FUST for locker",
                    "description": "This is an base  execution proposal to add use FUST as locker deposit asset.",
                    "msgs": [
                        {"msg_white_list_asset_locker": {"app_id": 1, "asset_id": 3}}
                    ],
                    "app_id_param": 1,
                }
            }
        },
    },
    {
        "proposalID": 6,
        "isProposal": True,
        "contractAddressKey": "governance_contract",
        "content": {
            "propose": {
                "propose": {
                    "title": "New proposal to whitelist FUST for locker",
                    "description": "This is an base  execution proposal to add use FUST as locker deposit asset.",
                    "msgs": [
                        {
                            "msg_whitelist_app_id_locker_rewards": {
                                "app_id": 1,
                                "asset_id": 3,
                            }
                        }
                    ],
                    "app_id_param": 1,
                }
            }
        },
    },
    {
        "proposalID": 7,
        "isProposal": True,
        "contractAddressKey": "governance_contract",
        "content": {
            "propose": {
                "propose": {
                    "title": "New proposal to whitelist FUST for locker",
                    "description": "This is an base  execution proposal to add use FUST as locker deposit asset.",
                    "msgs": [{"msg_whitelist_app_id_vault_interest": {"app_id": 1}}],
                    "app_id_param": 1,
                }
            }
        },
    },
    {
        "proposalID": 8,
        "isProposal": True,
        "contractAddressKey": "governance_contract",
        "content": {
            "propose": {
                "propose": {
                    "title": "New proposal for add pair for PETRI",
                    "description": "This is an base proposal execution proposal to add PETRI-FUST n.",
                    "msgs": [{"msg_whitelist_app_id_liquidation": {"app_id": 1}}],
                    "app_id_param": 1,
                }
            }
        },
    },
    {
        "proposalID": 9,
        "isProposal": True,
        "contractAddressKey": "governance_contract",
        "content": {
            "propose": {
                "propose": {
                    "title": "New proposal for setting auction params for collateral auctions.",
                    "description": "This is an base proposal to set collateral auction params with auction duration being 600 seconds.",
                    "msgs": [
                        {
                            "msg_add_e_s_m_trigger_params": {
                                "app_id": 1,
                                "target_value": {"amount": "200", "denom": "uharbor"},
                                "cool_off_period": 7200,
                                "asset_id": [4,3,13],
                                "rates": [1000000,1000000,1000000],
                            }
                        }
                    ],
                    "app_id_param": 1,
                }
            }
        },
    },
    {
        "proposalID": 10,
        "isProposal": True,
        "contractAddressKey": "governance_contract",
        "content": {
            "propose": {
                "propose": {
                    "title": "New proposal for add vault pair for PETRI C - FUST",
                    "description": "This is a base execution proposal to add PETRI C - FUST vault pair with given Vault properties a. Liquidation ratio : 140 % b. Stability Fee : 1%  c. Liquidation Penalty : 12% d. DrawDown Fee : 1% e. Debt Cieling : 100000000 FUST f. Debt Floor : 100 FUST ",
                    "msgs": [
                        {
                            "msg_add_extended_pairs_vault": {
                                "app_id": 1,
                                "pair_id": 2,
                                "stability_fee": "0.025",
                                "closing_fee": "0.00",
                                "liquidation_penalty": "0.12",
                                "draw_down_fee": "0.001",
                                "is_vault_active": True,
                                "debt_ceiling": "1000000000000",
                                "debt_floor": "100000000",
                                "is_stable_mint_vault": False,
                                "min_cr": "1.7",
                                "pair_name": "PETRI-A",
                                "asset_out_oracle_price": False,
                                "asset_out_price": 1000000,
                                "min_usd_value_left": 100000,
                            }
                        }
                    ],
                    "app_id_param": 1,
                }
            }
        },
    },
    {
        "proposalID": 11,
        "isProposal": True,
        "contractAddressKey": "governance_contract",
        "content": {
            "propose": {
                "propose": {
                    "title": "New proposal for add vault pair for PETRI C - FUST",
                    "description": "This is a base execution proposal to add PETRI C - FUST vault pair with given Vault properties a. Liquidation ratio : 140 % b. Stability Fee : 1%  c. Liquidation Penalty : 12% d. DrawDown Fee : 1% e. Debt Cieling : 100000000 FUST f. Debt Floor : 100 FUST ",
                    "msgs": [
                        {
                            "msg_add_extended_pairs_vault": {
                                "app_id": 1,
                                "pair_id": 4,
                                "stability_fee": "0.025",
                                "closing_fee": "0.00",
                                "liquidation_penalty": "0.12",
                                "draw_down_fee": "0.001",
                                "is_vault_active": True,
                                "debt_ceiling": "100000000000000",
                                "debt_floor": "100000000",
                                "is_stable_mint_vault": False,
                                "min_cr": "1.7",
                                "pair_name": "WETH-A",
                                "asset_out_oracle_price": False,
                                "asset_out_price": 1000000,
                                "min_usd_value_left": 100000,
                            }
                        }
                    ],
                    "app_id_param": 1,
                }
            }
        },
    },
    {
        "proposalID": 12,
        "isProposal": True,
        "contractAddressKey": "governance_contract",
        "content": {
            "propose": {
                "propose": {
                    "title": "New proposal for add vault pair for PETRI C - FUST",
                    "description": "This is a base execution proposal to add PETRI C - FUST vault pair with given Vault properties a. Liquidation ratio : 140 % b. Stability Fee : 1%  c. Liquidation Penalty : 12% d. DrawDown Fee : 1% e. Debt Cieling : 100000000 FUST f. Debt Floor : 100 FUST ",
                    "msgs": [
                        {
                            "msg_add_extended_pairs_vault": {
                                "app_id": 1,
                                "pair_id": 6,
                                "stability_fee": "0.0025",
                                "closing_fee": "0.00",
                                "liquidation_penalty": "0.001",
                                "draw_down_fee": "0.001",
                                "is_vault_active": True,
                                "debt_ceiling": "100000000000000",
                                "debt_floor": "100000000",
                                "is_stable_mint_vault": True,
                                "min_cr": "1.01",
                                "pair_name": "USDC-A",
                                "asset_out_oracle_price": False,
                                "asset_out_price": 1000000,
                                "min_usd_value_left": 100000,
                            }
                        }
                    ],
                    "app_id_param": 1,
                }
            }
        },
    },
    
]
