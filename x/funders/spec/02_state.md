<!--
order: 2
-->

# State

The module is mainly responsible for holding the funders state
and keeping track of each funded balance of a funder.

## Funder

A funder can exist independently of a pool and stores account information like address, moniker, description and so on.

- Funder: `0x01 | FunderAddr -> ProtocolBuffer(funder)`

```protobuf
syntax = "proto3";

message Funder {
  // address ...
  string address = 1;
  // moniker ...
  string moniker = 2;
  // identity is the 64 bit keybase.io identity string
  string identity = 3;
  // website ...
  string website = 4;
  // contact ...
  string contact = 5;
  // description are some additional notes the funder finds important
  string description = 6;
}
```

## Funding

Since funders and pools have a many-to-many relation we track the funding status in the `Funding` object containing
the information about what the funder funded in the specific pool like the amount of coins and the corresponding
amount per bundle. We also track how much the funder spent in total in the pool.

- Funding: `0x02 | 0x00 | PoolId | FunderAddr -> ProtocolBuffer(funding)`
- Funding: `0x02 | 0x01 | FunderAddr | PoolId -> ProtocolBuffer(funding)`

```protobuf
syntax = "proto3";

message Funding {
  // funder_id is the id of the funder
  string funder_address = 1;
  // pool_id is the id of the pool this funding is for
  uint64 pool_id = 2;
  // amounts is a list of coins the funder wants to fund the pool with
  repeated cosmos.base.v1beta1.Coin amounts = 3 [
    (gogoproto.nullable)     = false,
    (amino.dont_omitempty)   = true,
    (gogoproto.castrepeated) = "github.com/cosmos/cosmos-sdk/types.Coins"
  ];
  // amounts_per_bundle defines the amount of each coin that are distributed
  // per finalized bundle
  repeated cosmos.base.v1beta1.Coin amounts_per_bundle = 4 [
    (gogoproto.nullable)     = false,
    (amino.dont_omitempty)   = true,
    (gogoproto.castrepeated) = "github.com/cosmos/cosmos-sdk/types.Coins"
  ];
  // total_funded is the total amount of coins that the funder has funded
  repeated cosmos.base.v1beta1.Coin total_funded = 5 [
    (gogoproto.nullable)     = false,
    (amino.dont_omitempty)   = true,
    (gogoproto.castrepeated) = "github.com/cosmos/cosmos-sdk/types.Coins"
  ];
}
```

## FundingState

`FundingState` is an object that keeps track of the funding state of the pool. It contains a list of all active 
funders for each pool.

- FundingState: `0x03 | PoolId -> ProtocolBuffer(fundingState)`

```protobuf
syntax = "proto3";

message FundingState {
  // pool_id is the id of the pool this funding is for
  uint64 pool_id = 1;
  // active_funder_addresses is the list of all active fundings
  repeated string active_funder_addresses = 2;
}
```

## WhitelistCoinEntry

With multiple coin funding being possible we also have to limit the amount of coin types funders can fund or 
else a user could spam coins and dramatically increase the gas costs for protocol node operators. Therefore,
we have a coin whitelist so a funder can only fund a coin if it is included in the whitelist. For each coin there are
additional requirements like the minimum funding amount to also prevent spam. Note that the native $KYVE coin
is always included in the whitelist and can't be removed.

```protobuf
syntax = "proto3";

message WhitelistCoinEntry {
  // coin_denom is the denom of a coin which is allowed to be funded, this value
  // needs to be unique
  string coin_denom = 1;
  // min_funding_amount is the minimum required amount of this denom that needs
  // to be funded
  uint64 min_funding_amount = 2;
  // min_funding_amount_per_bundle is the minimum required amount of this denom
  // that needs to be funded per bundle
  uint64 min_funding_amount_per_bundle = 3;
  // coin_weight is a factor used to sort funders after their funding amounts
  string coin_weight = 4 [
    (gogoproto.customtype) = "cosmossdk.io/math.LegacyDec",
    (gogoproto.nullable) = false
  ];
}
```
