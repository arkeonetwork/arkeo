<!--
order: 2
-->

# State

### Claim Records

```protobuf
// A Claim Records is the metadata of claim data per address
message ClaimRecord {

  Chain chain = 1;

  // arkeo address of claim user
  string address = 2 [ (gogoproto.moretags) = "yaml:\"address\"" ];

  // total initial claimable amount for the user
  cosmos.base.v1beta1.Coin initial_claimable_amount = 3 [
    (gogoproto.nullable) = false,
    (gogoproto.moretags) = "yaml:\"initial_claimable_amount\""
  ];

  // true if action is completed
  // index of bool in array refers to action enum #
  repeated bool action_completed = 4 [ (gogoproto.moretags) = "yaml:\"action_completed\"" ];
}
```
ClaimRecords will be populated on genesis for all users and updated as a users takes actions to recieve additional airdrop tokens.

### State

```protobuf
// GenesisState defines the claim module's genesis state.
message GenesisState {
  
  // balance of the claim module's account
  cosmos.base.v1beta1.Coin module_account_balance = 1 [
    (gogoproto.moretags) = "yaml:\"module_account_balance\"",
    (gogoproto.nullable) = false
  ];
  
  Params params = 2 [(gogoproto.nullable) = false];

  // list of claim records, one for every airdrop recipient
  repeated ClaimRecord claim_records = 3 [
    (gogoproto.moretags) = "yaml:\"claim_records\"",
    (gogoproto.nullable) = false
  ];
}
```

Claim module's state consists of `params`, `claim_records`, and `module_account_balance`.
