<!--
order: 2
-->

# State

## Staker
Every address can create one single staker. Once the staker is created
people can delegate to it and the staker can start joining pools
(if the stake is high enough).

- Staker: `0x01 | StakerAddr -> ProtocolBuffer(staker)`

```go
type Staker struct {
    Address string
    // Needs to be a valid decimal representation
    Commission sdk.Dec 
    Moniker string 
    Website string
    Identity string 
    SecurityContact string 
    Details string 
}
```

## PoolAccount
The PoolAccount represents the membership of the staker in a given pool.
It contains the address of the protocol node which is allowed to vote
in favor of the staker and stores the poolId as well as a counter for 
penalty-points.

- PoolAccount: `0x02 | 0x00 | PoolId | StakerAddr -> ProtocolBuffer(poolAccount)`

One additional index is maintained to query for all valaccounts of a staker. 
For this index only the key is used as StakerAddr and PoolId contain all 
information to fetch the pool account using the main key.

- ValaccountIndex2: `0x02 | 0x01 | StakerAddr | PoolId -> (empty)`

```go
type PoolAccount struct {
    // PoolId defines the pool in which the address
    // is allowed to vote in.
    PoolId uint64
    // Staker is the address the pool account is voting for.
    Staker string
    // pool address is the account stored on the protocol
    // node which votes for the staker in the given pool
    PoolAccount string
    // When a node is inactive (does not vote at all)
    // a point is added. After a certain amount of points
    // is reached, the node gets kicked out.
    Points uint64
    // isLeaving indicates if a staker is leaving the given pool.
    IsLeaving bool
}
```

## Queue

The staker module contains two queues managing commission changes and
the leaving of pools.

### QueueState
For the queue the module needs to keep track of the head (HighIndex) and
tail (LowIndex) of the queue. New entries are appended to the
head. The EndBlocker checks the tail if entries are due and processes them.
There are two queues distinguished by the queue identifier.

- QueueState: `0x1E | 0x02 -> ProtocolBuffer(commissionQueueState)`
- QueueState: `0x1E | 0x03 -> ProtocolBuffer(leaveQueueState)`

```go
type QueueState struct {
    LowIndex uint64
    HighIndex uint64
}
```

### CommissionChangeQueueEntry
Every time a user starts a commission change, an entry is created
and appended to the head of the queue. I.e. the current HighIndex is
incremented and assigned to the entry.
The order of the queue is automatically provided by the KV-Store.

- CommissionChangeQueueEntry: `0x04 | 0x00 | Index  -> ProtocolBuffer(commissionChangeQueueEntry)`

A second index is provided so that users can query their own pending entries
without iterating the entire queue. The key is unique as there can only be
one commission change entry per staker. If a staker performs another
commission change the current pending entry is overwritten.

- UndelegationQueueEntryIndex2: `0x04 | 0x01 | StakerAddr  -> ProtocolBuffer(commissionChangeQueueEntry)`


```go
type CommissionChangeEntry struct {
    // Index is needed for the queue-algorithm which
    // processes the commission changes
    Index uint64
    // Staker is the address of the affected staker
    Staker string
    // Commission is the new commission which will
    // be applied after the waiting time is over.
    Commission sdk.Dec
    // CreationDate is the UNIX-timestamp in seconds
    // when the entry was created.
    CreationDate uint64
}
```


### LeavePoolQueueEntry
Every time a user initiates a pool leave, an entry is created
and appended to the head of the queue, i.e. the current HighIndex is
incremented and assigned to the entry.
The order of the queue is automatically provided by the KV-Store.

- LeavePoolEntry: `0x05 | 0x00 | Index  -> ProtocolBuffer(leavePoolEntry)`

A second index is provided so that users can query their own pending entries
without iterating the entire queue. 

- LeavePoolEntryIndex2: `0x05 | 0x01 | StakerAddr | PoolId  -> ProtocolBuffer(leavePoolEntry)`


```go
type CommissionChangeEntry struct {
    // Index is needed for the queue-algorithm which
    // processes the commission changes
    Index uint64
    // Staker is the address of the affected staker
    Staker string
    // Commission is the new commission which will
    // be applied after the waiting time is over.
    Commission String
    // CreationDate is the UNIX-timestamp in seconds
    // when the entry was created.
    CreationDate uint64
}
```