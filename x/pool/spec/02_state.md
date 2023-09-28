<!--
order: 2
-->

# State

The module is mainly responsible for holding the pools state
and keeping track of pool funders.

## Pools

The pool object is rather large and holds multiple sub-objects grouped
by functionality.

### Pool

Pool is the main type and holds everything a pool needs to know including some
sub-objects which are listed below. Each pool has their own module account, storing
the funds from the inflation split in order to pay those out with the funders to the
pool participants. The pool account is defined by the following: `pool/$ID`

- Pool: `0x01 | PoolId -> ProtocolBuffer(pool)`

```protobuf
syntax = "proto3";

enum PoolStatus {
  option (gogoproto.goproto_enum_prefix) = false;

  // POOL_STATUS_UNSPECIFIED ...
  POOL_STATUS_UNSPECIFIED = 0;
  // POOL_STATUS_ACTIVE ...
  POOL_STATUS_ACTIVE = 1;
  // POOL_STATUS_DISABLED ...
  POOL_STATUS_DISABLED = 2;
  // POOL_STATUS_NO_FUNDS ...
  POOL_STATUS_NO_FUNDS = 3;
  // POOL_STATUS_NOT_ENOUGH_DELEGATION ...
  POOL_STATUS_NOT_ENOUGH_DELEGATION = 4;
  // POOL_STATUS_UPGRADING ...
  POOL_STATUS_UPGRADING = 5;
}

message Protocol {
  // version holds the current software version tag of the pool binaries
  string version = 1;
  // binaries is a stringified json object which holds binaries in the
  // current version for multiple platforms and architectures
  string binaries = 2;
  // last_upgrade is the unix time the pool was upgraded the last time
  uint64 last_upgrade = 3;
}

message UpgradePlan {
  // version is the new software version tag of the upgrade
  string version = 1;
  // binaries is the new stringified json object which holds binaries in the
  // upgrade version for multiple platforms and architectures
  string binaries = 2;
  // scheduled_at is the unix time the upgrade is supposed to be done
  uint64 scheduled_at = 3;
  // duration is the time in seconds how long the pool should halt
  // during the upgrade to give all validators a chance of switching
  // to the new binaries
  uint64 duration = 4;
}

message Funder {
  // address is the address of the funder
  string address = 1;
  // amount is the current amount of funds in ukyve the funder has
  // still funded the pool with
  uint64 amount = 2;
}
```

```protobuf
syntax = "proto3";

message Pool {
  // id ...
  uint64 id = 1;
  // name ...
  string name = 2;
  // runtime ...
  string runtime = 3;
  // logo ...
  string logo = 4;
  // config ...
  string config = 5;

  // start_key ...
  string start_key = 6;
  // current_key ...
  string current_key = 7;
  // current_summary ...
  string current_summary = 8;
  // current_index ...
  uint64 current_index = 9;

  // total_bundles ...
  uint64 total_bundles = 10;

  // upload_interval ...
  uint64 upload_interval = 11;
  // inflation_share_weight ...
  uint64 inflation_share_weight = 12;
  // min_delegation ...
  uint64 min_delegation = 13;
  // max_bundle_size ...
  uint64 max_bundle_size = 14;

  // disabled ...
  bool disabled = 15;

  // old funders and total_funds fields
  reserved 16, 17;

  // protocol ...
  Protocol protocol = 18;
  // upgrade_plan ...
  UpgradePlan upgrade_plan = 19;

  // storage_provider_id ...
  uint32 current_storage_provider_id = 20;
  // compression_id ...
  uint32 current_compression_id = 21;
}
```
