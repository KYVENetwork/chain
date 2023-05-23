<!--
order: 3
-->

# Messages

## `CreateStaker`

Using this message, a user can create a staker. This can only be executed once
for each address. The sender can specify an amount which in turn is a direct
self-delegation to the given staker.

## `MsgUpdateMetadata`

This message changes Moniker, Website, Identity, SecurityContact and Details
of the staker. The message fails if the user does not have created a staker yet.

## `MsgUpdateCommission`

This message starts a commission change process by creating a new entry in the
commission change queue. Nothing else happens after that. The upcoming
commission change is shown in the staker query. So that delegators can see that
the given staker is about to change its commission.

After the `CommissionChangeTime` has passed the new commission is applied.

## `MsgClaimCommissionRewards`

This message claims the commission rewards of a protocol node. When a protocol
node receives commission rewards, it is transferred from the pool module to the
stakers module, which can be claimed with this message.

## `MsgJoinPool`

This message allows a staker to join a pool. For joining a pool the staker must
provide the poolId and an address which is operated by the protocol node. This
address is allowed to vote in favor of the staker. If this address misbehaves,
the staker will get slashed. The message also takes an amount as an argument
which is transferred to the valaddress. The valaddress needs a small balance to
pay for fees.

## `MsgLeavePoolResponse`

This message starts a leave pool process by creating a new entry in the leave
pool queue. Nothing else happens after that. The upcoming pool leave is shown in
the staker query. So that delegators can see that the given staker is about to
leave the given pool.

After the `LeavePoolTime` has passed the valaccount is deleted and the staker
can shut down the protocol node.