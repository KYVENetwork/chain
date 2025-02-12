<!--
order: 6
-->

# Events

The multi-coin-rewards module contains the following events:

## EventToggleMultiCoinRewards

EventToggleMultiCoinRewards indicates that someone has changed their
multi-coin-settings.

```protobuf
syntax = "proto3";

message EventToggleMultiCoinRewards {
  // address ...
  string address = 1;

  // enabled ...
  bool enabled = 2;

  // pending_rewards_claimed ...
  string pending_rewards_claimed = 3;
}
```

It gets emitted by the following actions:

- SetMultiCoinRewardDistributionPolicy
