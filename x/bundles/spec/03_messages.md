<!--
order: 3
-->

# Messages

## MsgSubmitBundleProposal

With this transaction the uploader of the current proposal round
submits his bundle proposal to the network for others to validate.
The uploader has to be a staker in the storage pool and should be
the designated uploader of this round. The most important property
which gets submitted is the storage id of the proposal. With this
other participants can retrieve the data and validate it for 
themselves. Once the proposal is validated the uploader receives
the bundle reward for his effort.

## MsgVoteBundleProposal

Once other participants see that a new bundle proposal is available
they validate it. Depending on the result they either vote with valid,
invalid or abstain. Abstain is a special vote which implies that
the validator could not make a decision. If the validator votes with
abstain it is impossible to receive a slash for that in the current round,
but the validator won't be chosen as uploader for the next round either.

## MsgClaimUploaderRole

If the storage pool is in genesis state (the pool just got created) or
the last bundle has been dropped for not reaching a required quorum the
fastest participant can claim the current free uploader role. This can 
only be called if the next uploader is not defined and the role for the
current round is free.

## MsgSkipUploaderRole

This transaction gets called when the uploader can't produce a bundle proposal
for whatever reason. Examples could be the storage provider being offline or
the data source not returning any data. With this the uploader skips his role
and lets another participant try to submit a valid bundle proposal.

