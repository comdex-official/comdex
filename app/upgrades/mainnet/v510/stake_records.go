package v510

// jq '.delegation_responses | map({address:.delegation.delegator_address,amount:((.balance.amount | tonumber)*0.05*((0.42/365)*13+1) | floor) | tostring})' DAN.JSON > to_mint.json

// Slash was 5%
// Lost APR is 42% for 13 days

var recordsJSONString = `[
 {
    "address": "comdex1g9wqptyaxlkzaryt8dezq4eed566kkfpswxjm9",
    "amount": "15000000"
  },
]`
