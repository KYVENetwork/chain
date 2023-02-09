<!--
order: 6
-->

# Events

The team module contains the following events:

## EventCreateTeamVestingAccount

EventCreateTeamVestingAccount indicates that a new team vesting account has been
created.

```protobuf
syntax = "proto3";

message EventCreateTeamVestingAccount {
  // id is a unique identify for each vesting account, tied to a single team member.
  uint64 id = 1;
  // total_allocation is the number of tokens reserved for this team member.
  uint64 total_allocation = 2;
  // commencement is the unix timestamp of the member's official start date.
  uint64 commencement = 3;
}
```

It gets thrown from the following actions:

- MsgCreateTeamVestingAccount

## EventClaimedUnlocked

EventClaimedUnlocked indicates that the authority has claimed unlocked $KYVE for a team
member.

```protobuf
syntax = "proto3";

message EventClaimedUnlocked {
  // id is a unique identify for each vesting account, tied to a single team member.
  uint64 id = 1;
  // amount is the number of tokens claimed from the unlocked amount.
  uint64 amount = 2;
  // recipient is the receiver address of the claim.
  string recipient = 3;
}
```

It gets thrown from the following actions:

- MsgClaimUnlocked

## EventClawback

EventClawback indicates that the authority has clawed back the remaining unvested $KYVE of a team
member vesting account.

```protobuf
syntax = "proto3";

message EventClawback {
  // id is a unique identify for each vesting account, tied to a single team member.
  uint64 id = 1;
  // clawback is a unix timestamp of a clawback. If timestamp is zero
  // it means that the account has not received a clawback
  uint64 clawback = 2;
  // amount which got clawed back.
  uint64 amount = 3;
}
```

It gets thrown from the following actions:

- MsgClawback
