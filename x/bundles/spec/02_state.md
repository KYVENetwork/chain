<!--
order: 2
-->

# State

The module is mainly responsible for handling the current bundle 
proposal state including holding the state for all finalized bundles.

## Bundle Proposals
Bundle proposals have their own prefix in the KV-Store and are defined in
one proto file

### BundleProposal
BundleProposal has all the information of the current bundle proposal and 
also keeps track of votes. One bundle proposal is always linked to one
storage pool with a 1-1 relationship.

- BundleProposal `0x01 | PoolId -> ProtocolBuffer(bundleProposal)`

```go
type BundleProposal struct {
    PoolId uint64
    StorageId string
    Uploader string
    NextUploader string
    DataSize uint64
    BundleSize uint64
    ToKey string
    BundleSummary string
    DataHash string
    UpdatedAt uint64
    VotersValid []string
    VotersInvalid []string
    VotersAbstain []string
    FromKey string
    StorageProviderId uint32
    CompressionId uint32
}
```

## Finalized Bundles
Finalized bundles have their own prefix in the KV-Store.

### FinalizedBundle
FinalizedBundle has all the important information of a bundle which is saved
forever on the KYVE chain.

- FinalizedBundle `0x02 | PoolId | Id -> ProtocolBuffer(finalizedBundle)`

```go
type FinalizedBundle struct {
    PoolId uint64
    Id uint64
    StorageId string
    Uploader string
    FromIndex uint64
    ToIndex uint64
    ToKey string
    BundleSummary string
    DataHash string
    FinalizedAt {
        Height uint64
        Timestamp uint64
    }   
    FromKey string
    StorageProviderId uint32
    CompressionId uint32
    StakeSecurity {
        ValidVotePower uint64
        TotalVotePower uint64
    }
}
```

### BundleVersionMap

The version map keeps track of which protocol version was present at given 
block heights. It is only updated during chain upgrades. It helps the query
handler to probably decode a finalized bundle. Later it might also be important
for on chain computations. 

- BundleVersionMap `0x03 -> ProtocolBuffer(BundleVersionMap)`


## Round-Robin
For correctly determining the next uploader the current round-robin
progress needs to be saved in the KV-Store. Every pool keeps track of its
own round-robin state.

### RoundRobinSingleValidatorProgress
This struct is not stored directly in the KV-Store but used by the
RoundRobinProgress struct.

```go
type RoundRobinSingleValidatorProgress struct {
    // address of the validator
    Address string
    // progress within the current round-robin set
    Progress int64
}
```

### RoundRobinProgress
RoundRobinProgress stores the current state of the round-robin selection for a 
given pool.

- RoundRobinProgress `0x04 | PoolId -> ProtocolBuffer(roundRobinProgress)`

```go
message RoundRobinProgress {
    PoolId uint64
    ProgressList []RoundRobinSingleValidatorProgress
}
```