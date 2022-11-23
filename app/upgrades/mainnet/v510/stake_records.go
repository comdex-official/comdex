package v510

// jq '.delegation_responses | map({address:.delegation.delegator_address,amount:((.balance.amount | tonumber)*0.05*((0.42/365)*13+1) | floor) | tostring})' DAN.JSON > to_mint.json

// Slash was 5%
// Lost APR is 42% for 13 days

var recordsJSONString = `[
 {
    "address": "comdex1g9wqptyaxlkzaryt8dezq4eed566kkfpswxjm9",
    "amount": "50000"
  },
  {
    "address": "comdex1s6pkt43rjzjqpcwaz6se9ajzr2wq7kfwyhy3nl",
    "amount": "5000000"
  },
  {
    "address": "comdex1phq8yxpagpcmtdv7jymvvyxeekk5y2nx5arlfu",
    "amount": "5000000"
  },
  {
    "address": "comdex1phq8yxpagpcmtdv7jymvvyxeekk5y2nx5arlfu",
    "amount": "5000000"
  },
  {
    "address": "comdex10gh5u4qey4xxewyjwcq6jcp0h8vfqusw3y7nku",
    "amount": "5000000"
  },
  {
    "address": "comdex1t4f654txl5hdzghz932cesjz47j7088kka96m9",
    "amount": "5000000"
  }
]`
