<!--
order: 1
-->

# Concepts

The bundles module implements the main logic for archiving data on other
storage providers. It handles the submission of bundle proposals, the
finalization of valid bundles and keeps track of validator votes. It
uses the staker and delegation module to determine who can submit/vote
and which participants get slashed for malicious behaviour.

## Code Structure

This module adheres to our global coding structure, defined [here](../../../CodeStructure.md).

## Validating in Rounds

Data which gets validated and archived by KYVE is validated in rounds.
In every round one uploader will get selected in a deterministic way who
is then responsible for submitting a bundle proposal. Every other participant
then has to validate this bundle proposal. If the network agrees on the proposal
the bundle gets finalized and the network moves to the next bundle proposal.

## Bundle Proposals

In order to get data validated and archived by KYVE a participant of a pool
has to package data in a bundle proposal and submit it to the KYVE chain.
The role of the participant is then `Uploader`. Once the bundle proposal with
required metadata like the data range, data size and hash is submitted other
participants can vote on this proposal.

## Voting

All other participants who have not uploaded data to the network because they
were not the designated uploader have to validate the submitted data. The role
of the participant is then `Validator`. They take the storage id the uploader
submitted, retrieve it from the storage provider and then compare it locally
with respect to the used runtime implementation. Furthermore, the metadata is
validated by comparing the data size and hash. The validator then votes on the
bundle proposal accordingly.

## Bundle Evaluation

After a certain timeout (`upload_interval`) the next uploader can submit the next
bundle proposal. While the next bundle proposal is submitted the current one
gets evaluated. If more than 50% voted for valid the bundle gets finalized and gets
saved forever on-chain so that everyone can use that validated data.

## Punishing malicious behaviour

If more than 50% voted invalid the uploader receives a slash and gets removed 
from the storage pool. Furthermore, validators who voted incorrectly also get 
slashed and removed. If an uploader or validator don't upload/vote in a specific
time range they receive points. If they have a certain number of points they 
receive a timeout slash and also get removed.
