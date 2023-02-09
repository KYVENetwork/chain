<!--
order: 4
-->

# EndBlock

The `x/stakers` module end-block hook handles the commission-change and
leave-pool queue. After the `CommissionChangeTime` resp `LeavePoolTime` has
passed, the queue entry is executed.
