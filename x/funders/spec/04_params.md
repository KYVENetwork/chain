<!--
order: 4
-->

# Parameters

The pool module contains the following parameters:

| Key                | Type                 | Example |
|--------------------|----------------------|---------|
| CoinWhitelist      | WhitelistCoinEntry[] |         |
| MinFundingMultiple | math.LegacyDec (%)   | 20      |

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
