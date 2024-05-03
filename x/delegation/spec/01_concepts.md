<!--
order: 1
-->

# Concepts

The KYVE Protocol layer implements staking and delegation similar to
Tendermint. Validators that want to participate in the protocol need to stake
$KYVE to have voting power and join different KYVE storage pools.

## Code Structure

This module adheres to our global coding structure, defined [here](../../../CodeStructure.md).

## Delegation

Users who want to support a validator can delegate their $KYVE. This validator
is often referred to in our codebase as a `staker`.

If a validator delegates to itself, this is called self-delegation. From a
technical point of view, self and user delegations are treated as the same.

## F1 Distribution

Because there is no limit to the number of validators, a direct payout of each
reward would cost an outrageous amount of gas. We have turned to the
"F1 Fee Distribution" algorithm to solve this issue. It handles the delegation
itself, payouts of rewards, and slashing events.

The main idea is that if there is no change to delegation distribution (in
other words, no new delegations or undelegations), there is no need to pay out
rewards. When users want to withdraw their rewards or update their
delegation amount, their rewards are calculated and correctly distributed.

## Rewards and Slashes

Users who delegate tokens to a validator will receive a portion of its rewards.
These rewards are generated when the validator produces valid data bundles in a
KYVE storage pool. On the other hand, these delegations are also subject to
slashing events if the validator misbehaves.

# References

[1] D. Ohja, C. Goes. F1 Fee Distribution. 
In *International Conference on Blockchain Economics, Security and Protocols*, *pages 10:1-10:6*, 2019. URL: `https://d-nb.info/1208239872/34`
