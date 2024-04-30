#!/bin/bash

# Source this file with `source scripts/commands.sh` to use the functions in your shell.
# First argument can be an alternative binary path, e.g. `source scripts/commands.sh /path/to/kyved`.

RED="\e[31m"
GREEN="\e[32m"
ENDCOLOR="\e[0m"

export BINARY=${1:-"kyved"}

# Check if binary is executable
if ! command -v "$BINARY" &> /dev/null; then
  echo -e "$RED✗ $BINARY binary not found. Please install it or use the correct path.$ENDCOLOR"
  return
fi

export CHAIN_ID=$($BINARY status --output json | jq -r '.node_info.network' | tr -d '\"')
if [[ -z "$CHAIN_ID" ]]; then
  export CHAIN_ID=$($BINARY status | jq -r '.NodeInfo.network' | tr -d '\"')
  if [[ -z "$CHAIN_ID" ]]; then
    echo -e "$RED✗ Could not fetch chain ID. Make sure your node is running.$ENDCOLOR"
    return
  fi
fi

export TX="--gas 200000 --fees 5000000tkyve --yes --chain-id $CHAIN_ID"
export ALICE_TX="$TX --from alice"
export TESTBACKEND="--keyring-backend test"
export CHAINHOME="--home ~/.chain"
export JSON="--output json"

check_tx() {
  echo "🔍 Checking transaction..."

  local tx_output
  local tx_hash
  local code
  tx_output=$(cat -)
  tx_hash=$(echo "$tx_output" | jq -r '.txhash' | tr -d '.')
  code=$(echo "$tx_output" | jq -r '.code' | tr -d '.')

  if [[ -z "$tx_hash" || "$code" -ne 0 ]]; then
    local raw_log=$(echo "$tx_output" | jq -r '.raw_log')
    echo -e "${RED}✗ Transaction failed with code $code:\n$raw_log${ENDCOLOR}"
    return "$code"
  fi

  local sleep_time=0.5
  local elapsed_seconds=0
  local progress_bar=""

  # Run a loop to check for tx_output
  while true; do
    local tx_output=$(kyved q tx "$tx_hash" --output json 2>/dev/null)
    local code=$(echo "$tx_output" | jq -r '.code')

    if [[ -z "$code" ]]; then
      printf "\r⏳ Fetching transaction: %02d seconds | Progress: [%-30s]" "$elapsed_seconds" "$progress_bar"
      sleep $sleep_time
      ((elapsed_seconds++))
      progress_bar+="="
    else
      break
    fi
  done

  echo ""

  if [[ "$code" -eq 0 ]]; then
    echo "✏️ logs:"
    echo "$tx_output" | jq '.logs'
    echo ""
    echo -e "${GREEN}✅ Transaction $tx_hash successful!${ENDCOLOR}"
  else
    echo "✏️ raw_log:"
    echo "$tx_output" | jq '.raw_log'
    echo ""
    echo -e "${RED}✗ Transaction $tx_hash failed with code $code.${ENDCOLOR}"
  fi
  return "$code"
}
export -f check_tx

echo "👋 Welcome to the KYVE CLI!"
echo ""
echo "# Environment variables -> override with your own values"
echo "export TX=\"$TX\""
echo "export ALICE_TX=\"$ALICE_TX\""
echo "export TESTBACKEND=\"$TESTBACKEND\""
echo "export CHAINHOME=\"$CHAINHOME\""
echo "export JSON=\"$JSON\""
echo ""
echo "# Add Alice's key"
echo "$BINARY keys add alice --recover \$TESTBACKEND"
echo "# Mnemonic"
echo "expect crisp umbrella hospital firm exhibit future size slot update blood deliver fat happy ghost visa recall usual path purity junior ring ordinary stove"
echo ""
echo "# Send coins"
echo "$BINARY tx bank send alice kyve1jve70gkvgyvdnrxw4q7ry7vq2asu25ac0m48vk 1000tkyve \$CHAINHOME \$TESTBACKEND \$TX \$JSON | check_tx"
echo ""
echo "# Governance"
echo "$BINARY tx gov draft-proposal"
echo "$BINARY tx gov submit-proposal draft_proposal.json \$ALICE_TX \$CHAINHOME \$TESTBACKEND \$TX \$JSON | check_tx"
echo "$BINARY tx gov vote \$($BINARY q gov proposals --page-limit 1 --page-reverse --proposal-status voting-period \$JSON | jq '.proposals[0].id | tonumber') yes \$ALICE_TX \$CHAINHOME \$TESTBACKEND \$JSON | check_tx"
echo ""
echo "# Funders"
echo "$BINARY tx funders create-funder Alice \$CHAINHOME \$TESTBACKEND \$ALICE_TX \$JSON | check_tx"
echo "$BINARY tx funders update-funder Alice \$CHAINHOME \$TESTBACKEND \$ALICE_TX \$JSON | check_tx"
