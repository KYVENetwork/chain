<!--
order: 1
-->

# Concepts

For policy reasons not everybody can or is allowed to receive, hold and 
control other tokens. Therefore, users need to opt-in for multi-coin rewards.
Because users could miss the opt-in, there is a grace-period in which multi-coin
rewards can still be retrieved after a claim. If a user does not enable 
multi-coin rewards and does not enable it within the grace-period after a claim,
the multi-coin rewards get re-distributed to the existing pools according to
a distribution policy which can be modified by an admin-address which is set
by the governance.

## (Re-)distribution policy

The redistribution policy defines which multi-coin denom gets re-distributed
to which pool depending on a certain weight. To allow quick updates an
admin address (other than the governance) can be specified which can update
the policy. The address can not drain any rewards, but only modify to which
pools the coins get re-distributed to. Therefore, the governance might
set the admin address to a trusted (multi-sig) address.

## Token Flow
Within the withdraw-rewards function inside the CosmosSDK distribution module
the multi-coin-rewards module is called. 
1. User has opted in: All tokens are directly paid out to the user
2. User has not opted in: Only the native token is paid out, the other tokens are
   transferred to the `multi_coin_rewards` module account. A queue entry is
   created and a user has a certain amount of time to enable multi-coin-rewards.
   1. User enables rewards within time: All pending rewards are transferred to the user
   2. User does not enable rewards within time: The rewards are transferred to 
      the `multi_coin_rewards_distribution` module account.

Every 50 blocks all coins in `multi_coin_rewards_distribution` are 
redistributed according to the distribution policy. If tokens are not
covered by the policy they remain inside the module account.
