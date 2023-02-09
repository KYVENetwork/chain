<!--
order: 4
-->

# EndBlock

The `x/delegation` module end-block hook handles the unbonding queue. After the
`DelegationUnbondingTime` time has passed, delegators will receive the number
of tokens they undelegated. However, if the validator they were delegating to
was slashed during this time, the received amount will be smaller.

Please note that a queue like unbonding doesn't track redelegation. Instead,
the remaining redelegation slots are calculated on demand during transaction
execution.
