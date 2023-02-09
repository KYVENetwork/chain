<!--
order: 3
-->

# Messages

## `MsgDelegate`

Using this message, a user can delegate a specified amount to a KYVE protocol
validator. The chosen validator must exist in the `x/stakers` module.
Otherwise, the transaction will fail. If the user previously delegated to this
validator, any pending rewards will be withdrawn immediately.

Delegated $KYVE tokens are locked for `DelegationUnbondingTime` seconds. This
is the minimum time users need to wait before they can use their tokens again.

## `MsgWithdrawRewards`

It is impossible to distribute rewards to delegators immediately. This is
because of gas limits. Therefore, all rewards are collected in a pool, and
delegators can use this message to withdraw their pending rewards.

## `MsgUndelegate`

This message starts the undelegation process by creating a new entry in the
unbonding queue. Nothing else happens after that, and users will still receive
rewards and are still subject to slashing. After `DelegationUnbondingTime`
seconds, the actual unbonding is performed via an end-block hook.

After the unbonding time has passed, if the amount requested to undelegate is
higher than the actual amount (because of a slashing event), only the available
amount is returned to the user.

## `MsgRedelegate`

This message allows delegators to switch their delegation between different
KYVE protocol validators. It is only possible to redelegate to active
validators (this means they are participating in at least one storage pool).

Every delegator has `RedelegationMaxAmount` number of spells. Once a spell is
cast, it goes on a cooldown for `RedelegationCooldown` seconds. If all
redelegation slots are used, the user must wait until the first slot is
available again.
