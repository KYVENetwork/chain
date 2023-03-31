<!--
order: 6
-->

# Parameters

The `x/delegation` module relies on the following parameters:

| Key                       | Type            | Default Value |
|---------------------------|-----------------|---------------|
| `UnbondingDelegationTime` | uint64 (time s) | 432000        |
| `RedelegationCooldown`    | uint64 (time s) | 432000        |
| `RedelegationMaxAmount`   | uint64 (time s) | 5             |
| `VoteSlash`               | sdk.Dec (%)     | 0.1           |
| `UploadSlash`             | sdk.Dec (%)     | 0.2           |
| `TimeoutSlash`            | sdk.Dec (%)     | 0.02          |
