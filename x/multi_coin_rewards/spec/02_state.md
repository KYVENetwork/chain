<!--
order: 2
-->

# State

The module is mainly responsible for holding the policy itself and keeping
track of who has opted in for multi-coin rewards.

## MultiCoinRewardsEnabled

The users who have opted in for multi-coin rewards are stored as a key in
the IAVL tree. There is no value associated with it. If the key exists,
the users has opted in.

- MultiCoinRewardsEnabled: `0x01 | AccAddress -> {}`

## MultiCoinDistributionPolicy

The MultiCoinDistributionPolicy stores for every denom a list of pools and
weights. The weights determine on how the rewards of a given denom are 
re-distributed under the pools.

```protobuf
syntax = "proto3";

// MultiCoinDistributionPolicy ...
message MultiCoinDistributionPolicy {
  repeated MultiCoinDistributionDenomEntry entries = 1;
}

// MultiCoinDistributionDenomEntry ...
message MultiCoinDistributionDenomEntry {
  string denom = 1;
  repeated MultiCoinDistributionPoolWeightEntry pool_weights = 2;
}

// MultiCoinDistributionPoolWeightEntry ...
message MultiCoinDistributionPoolWeightEntry {
  uint64 pool_id = 1;
  string weight = 2 [
    (gogoproto.customtype) = "cosmossdk.io/math.LegacyDec",
    (gogoproto.nullable) = false
  ];
}

```
