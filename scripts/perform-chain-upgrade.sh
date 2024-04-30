#!/bin/bash
# This script is only used for local testing. Do not use it in production.

RED="\e[31m"
ENDCOLOR="\e[0m"

# Check if at least upgrade name is provided
if [ $# -lt 1 ]; then
    echo -e "$RED✗ Usage: perform-chain-upgrade.sh [upgrade_name] [wait_blocks (default=30)]$ENDCOLOR"
    exit 1
fi

# Check if wait_blocks is provided
if [ -z "$2" ]; then
    WAIT_BLOCKS=30
else
    WAIT_BLOCKS=$2
fi

echo "🚀 Upgrading chain to $1 in $WAIT_BLOCKS blocks"

# Check if binary env is set
if [ -z "$BINARY" ]; then
    echo -e "$RED✗ BINARY env not set. Please source commands.sh before running this script$ENDCOLOR"
    exit 1
fi

get_height() {
  local height=
  height=$($BINARY status | jq '.sync_info.latest_block_height' | tr -d '\"')
  if [ "$height" == "null" ]; then
      height=$($BINARY status | jq '.SyncInfo.latest_block_height' | tr -d '\"')
      if [[ -z "$height" ]]; then
          exit 1
      fi
  fi
  echo "$height"
}

PROPOSAL='{
  "messages": [
    {
      "@type": "/cosmos.upgrade.v1beta1.MsgSoftwareUpgrade",
      "authority": "kyve10d07y265gmmuvt4z0w9aw880jnsr700jdv7nah",
      "plan": {
        "name": "MY_UPGRADE_NAME",
        "height": "MY_UPGRADE_height",
        "info": ""
      }
    }
  ],
  "metadata": "ipfs://QmSbXwSgAnhhkDk5LtvgEwKUf6iUjww5FbMngYs86uGxdG",
  "deposit": "50000000000MY_DENOM",
  "title": "Upgrade to MY_UPGRADE_NAME",
  "summary": "This is an upgrade to MY_UPGRADE_NAME",
  "expedited": false
}'

NAME=$1
UPGRADE_HEIGHT=`expr "$(get_height)" + $WAIT_BLOCKS`
DENOM="tkyve"

PROPOSAL="${PROPOSAL//MY_UPGRADE_height/$UPGRADE_HEIGHT}"
PROPOSAL="${PROPOSAL//MY_UPGRADE_NAME/$NAME}"
PROPOSAL="${PROPOSAL//MY_DENOM/$DENOM}"
printf "$PROPOSAL" > /tmp/proposal.json

# Submit proposal
$BINARY tx gov submit-proposal /tmp/proposal.json $ALICE_TX $CHAINHOME $TESTBACKEND $TX $JSON | check_tx

# Exit if proposal submission failed
if [ $? -ne 0 ]; then
    exit 1
fi
echo ""
echo "Submitted proposal"
echo ""

# Get proposal ID
ID=$($BINARY q gov proposals --page-limit 1 --page-reverse --proposal-status voting-period $JSON | jq '.proposals[0].id | tonumber' | tee /dev/null)
if [[ -z "$ID" ]]; then
    ID=$($BINARY q gov proposals --limit 1 --reverse --status voting_period $JSON | jq '.proposals[0].id | tonumber')
    if [[ -z "$ID" ]]; then
        echo -e "$RED✗ Could not fetch proposal ID. Make sure your node is running.$ENDCOLOR"
        exit 1
    fi
fi

# Vote on proposal
$BINARY tx gov vote $ID yes $ALICE_TX $CHAINHOME $TESTBACKEND $JSON | check_tx
echo ""
echo "Voted yes on proposal $ID"
echo "Scheduled upgrade for height $UPGRADE_HEIGHT"

did_proposal_pass() {
  local result=$1
  local status=$(echo "$result" | jq -r '.proposal.status')

  if [ "$status" == "null" ]; then
    status=$(echo "$result" | jq -r '.status')
  fi

  # status 3 or "PROPOSAL_STATUS_PASSED" is passed
  if [[ "$status" -eq 3 || "$status" == "PROPOSAL_STATUS_PASSED" ]]; then
    return 0
  else
    return 1
  fi
}

did_proposal_fail() {
  local result=$1
  local status=$(echo "$result" | jq -r '.proposal.status')

  if [ "$status" == "null" ]; then
    status=$(echo "$result" | jq -r '.status')
  fi

  # everything except status 2 (PROPOSAL_STATUS_VOTING_PERIOD) or status 3 (PROPOSAL_STATUS_PASSED) is failed
  if [[ "$status" -ne 2 && "$status" -ne 3 && "$status" != "PROPOSAL_STATUS_VOTING_PERIOD" && "$status" != "PROPOSAL_STATUS_PASSED" ]]; then
    return 0
  else
    return 1
  fi
}

# Wait for upgrade and poll status
elapsed_seconds=0
progress_bar_length=$(expr $WAIT_BLOCKS \* 1)
while true; do
  result=$($BINARY q gov proposal $ID $JSON)

  if did_proposal_pass "$result"; then
    printf "\n\r✅ Upgrade successful! Took %02d seconds\n" "$elapsed_seconds"
    break
  elif did_proposal_fail "$result"; then
    printf "\n\r${RED}✗ Upgrade failed!\nCheck if your voting parameters are correct. They need to have a short voting time and a low quorum threshold.${ENDCOLOR}\n"
    break
  fi

  height=$(get_height)
  remaining_blocks=$(expr $UPGRADE_HEIGHT - $height)

  passed_blocks=$(expr $WAIT_BLOCKS - $remaining_blocks)
  progress_bar=$(printf "%0.s=" $(seq 1 $passed_blocks))

  printf "\r⏳ Waiting for upgrade: %02d seconds | Progress: [%-*s] %02d blocks left" "$elapsed_seconds" "$progress_bar_length" "$progress_bar" "$remaining_blocks"

  if [[ "$remaining_blocks" -le 0 ]]; then
    printf "\n\r${RED}✗ Upgrade failed!\nCheck if your voting parameters are correct. They need to have a short voting time and a low quorum threshold.${ENDCOLOR}\n"
    exit 1
  fi

  sleep 1
  ((elapsed_seconds++))
done