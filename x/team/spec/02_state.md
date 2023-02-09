<!--
order: 2
-->

# State

The module is mainly responsible for handling the team vesting account states.

## TeamVestingAccounts

The state is defined in one main proto file.

### TeamVestingAccount

Each team member gets a _TeamVestingAccount_ assigned. A TeamVestingAccount is not identified by an address, rather
its identified by an incrementing ID. It is tracked off-chain which Account ID belongs to which team member.

The TeamVestingAccount stores the total amount of $KYVE the team member has and the commencement of the team member
(a unix timestamp of when the team member official joined KYVE). Furthermore, the clawback time (if the 
team member leaves KYVE) and the already claimed $KYVE is stored. If clawback is zero the member did not receive
a clawback

- TeamVestingAccountKey: `0x02 | Id -> ProtocolBuffer(teamVestingAccount)`
- TeamVestingAccountCountKey: `0x03 | Count -> ProtocolBuffer(teamVestingAccountCount)`

```protobuf
syntax = "proto3";

message TeamVestingAccount {
    // id is a unique identify for each vesting account, tied to a single team member.
    uint64 id = 1;
    // total_allocation is the number of tokens reserved for this team member.
    uint64 total_allocation = 2;
    // claimed is the amount of tokens already claimed by the account holder
    uint64 claimed = 3;
    // clawback is a unix timestamp of a clawback. If timestamp is zero
    // it means that the account has not received a clawback
    uint64 clawback = 4;
    // commencement is the unix timestamp of the member's official start date.
    uint64 commencement = 5;
}
```