<!--
order: 3
-->

# Messages

All txs of this module can be only called by the authority.

## `MsgCreateTeamVestingAccount`

Using this message, the authority can create a new _TeamVestingAccount_.
For that the authority has to provide the total allocation the team member
receives and the commencement date of the team member. The ID for the new
TeamVestingAccount will be automatically assigned on-chain.

The tx fails  if the team module has not enough funds anymore to create a 
vesting account with the requested allocation, therefore ensuring the 
authority does not overspend $KYVE.

## `MsgClawback`

If a team member leaves during his vesting period and the authority wants 
to clawback the **remaining** unvested $KYVE the authority can call this 
tx. It has to provide the account id and the unix timestamp of when the
clawback should be applied. The authority can update the clawback time of
an account multiple times and even remove it again if the time is `0`.

## `MsgClaimUnlocked`

If a team member wants to claim $KYVE of his unlocked amount he has to notify
the authority to do that for him. The team member has to provide a wallet
address to which the authority should claim the $KYVE to. In order to claim
the authority has to call this tx with the matching account ID and a recipient
address which can be the team members wallet directly or send it to a proxy address
instead to deal with e.g. taxes.
