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
    FinalizedAt uint64
    FromKey string
    StorageProviderId uint32
    CompressionId uint32
}
```
