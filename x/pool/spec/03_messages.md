<!--
order: 4
-->

# Messages

## MsgFundPool

With this transaction stakeholders of the pool can provide funds to a storage pool, so it can continue
validating and archiving data. Funding a pool does not earn any rewards, actually the opposite is the case.
By funding a pool the funders pay for the rewards the validators receive. If all funder 
slots are occupied a user needs to fund more than the current lowest funder in order for
the transaction to succeed.

## MsgDefundPool

When a funder has funded a pool he can of course withdraw his funds again. If the full amount is defunded the funder
gets completely removed from the pool. Also funds can be partially defunded.

## MsgCreatePool

MsgCreatePool is a gov transaction and can be only called by the governance authority. To submit this transaction
someone has to create a MsgCreatePool governance proposal.

This will create a new storage pool in the KYVE network where other participants can join.

## MsgUpdatePool

MsgUpdatePool is a gov transaction and can be only called by the governance authority. To submit this transaction
someone has to create a MsgUpdatePool governance proposal.

This will update an existing storage pool based on the given parameters.

## MsgDisablePool

MsgDisablePool is a gov transaction and can be only called by the governance authority. To submit this transaction
someone has to create a MsgDisablePool governance proposal.

This will disable a currently active pool. Once a pool is disabled it will not use any funds and therefore will not
validate or archive any data.

## MsgEnablePool

MsgEnablePool is a gov transaction and can be only called by the governance authority. To submit this transaction
someone has to create a MsgEnablePool governance proposal.

This will enable a currently disabled pool. Once a pool is enabled it can continue to validate and archive data again.

## MsgScheduleRuntimeUpgrade

MsgScheduleRuntimeUpgrade is a gov transaction and can be only called by the governance authority. To submit 
this transaction someone has to create a MsgScheduleRuntimeUpgrade governance proposal.

This will schedule an upgrade for the specified runtime. A runtime upgrade contains an upgrade version and the 
associated upgrade binaries. If the scheduled upgrade time is reached the upgrade will be performed with the
specified upgrade duration.

## MsgCancelRuntimeUpgrade

MsgCancelRuntimeUpgrade is a gov transaction and can be only called by the governance authority. To submit 
this transaction someone has to create a MsgCancelRuntimeUpgrade governance proposal.

This will cancel a scheduled runtime upgrade if it has not been reached yet. If the upgrade was already performed
it is not possible to cancel anymore. But it is still possible to downgrade a runtime by simply "upgrading" to the
prior version.
