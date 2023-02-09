<!--
order: 1
-->

# Concepts

The team module is responsible for distributing the team allocation **16.5%** of the total
genesis supply to all eligible team members. It uses a special mechanism we call
"Two-Layer-Vesting" which is the reason why the team distribution received its own module.

## Code Structure

This module adheres to our global coding structure, defined [here](../../../CodeStructure.md).

## Team Vesting Accounts

Each team member gets a _TeamVestingAccount_ assigned. A TeamVestingAccount is not identified by an address, rather
its identified by an incrementing ID. It is tracked off-chain which Account ID belongs to which team member.

Furthermore, a TeamVestingAccount tracks the _commencement_ (the official start date of working at KYVE) and a
_clawback_ (the official leave date after working at KYVE). The vesting amount depends on those two variables which
are custom for each team member.

Finally, a TeamVestingAccount tracks the amount a team member has already claimed from his vesting account, this is
used later to calculate the current unlocked amount.

## Vesting

The total vesting duration is a constant set to 3 years and won't change anymore. $KYVE will vest 3 years linearly
from the commencement date. During vesting there is a cliff which is a constant set to 1 year and won't change
anymore. So for the first year the vested amount is zero, after the first day the cliff is over the vested amount
is 33.33% of the total allocation since one third of the vesting duration passed.

Vested $KYVE can not be spent by the team member already. Vested $KYVE is just the first of the two layers of vesting.
If $KYVE has vested the team member is just eligible for inflation rewards which will be explained in detail below.

## Unlocking

Once $KYVE has successfully vested for a TeamVestingAccount it is still locked. In order for the team member to claim
his $KYVE they need to unlock. The Unlock starts either exactly 1 year after commencement or exactly 1 year after
TGE, whatever is the latter. The Unlock duration is constant and is set to 2 years. During unlocking there is
no cliff and $KYVE is unlocking at a linear rate based on seconds passed.

## Clawback

If a team members leaves KYVE during his vesting period the authority is allowed to clawback the 
**remaining, unvested** $KYVE from the vesting account. A clawback is a unix timestamp of when the team member
left. The clawback can only be initiated by the team module authority.

## Claim Unlocked $KYVE

Once the $KYVE of a team member have unlocked the team member is allowed to claim them. In order to do that
the team member has to notify the authority that he wants to claim and provide it with a receiver address.
With that the authority can then claim for the team member and claim his $KYVE which will then get transferred
to the receiver address.
