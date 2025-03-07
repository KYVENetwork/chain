import os
import json
import sys

import requests
import bech32


def verify_tx(api_endpoint, tx_hash, expect_from_address, expect_to_address):
    x = requests.get(api_endpoint + "/cosmos/tx/v1beta1/txs/" + tx_hash)
    if x.status_code != 200:
        raise Exception("transaction does not exist: ", tx_hash)

    tx = x.json()
    if tx["tx"]["body"]["memo"] != "Shared-Staking":
        raise Exception("incorrect memo for transaction: ", tx_hash)

    if tx["tx"]["body"]["messages"][0]["from_address"] != expect_from_address:
        raise Exception("Incorrect from_address. Expected: {}, got: {}"
                        .format(expect_from_address, tx["tx"]["body"]["messages"][0]["from_address"]))

    if tx["tx"]["body"]["messages"][0]["to_address"] != expect_to_address:
        raise Exception("Incorrect to_address. Expected: {}, got: {}"
                        .format(expect_to_address, tx["tx"]["body"]["messages"][0]["to_address"]))


def verify_proof(api_endpoint, entry):
    x = requests.get(api_endpoint + "/cosmos/staking/v1beta1/validators/" + entry["consensus_address"])
    if x.status_code != 200:
        raise Exception("Consensus validator does not exist: ", entry["consensus_address"])

    x = requests.get(api_endpoint + "/kyve/query/v1beta1/staker/" + entry["protocol_address"])
    if x.status_code != 200:
        raise Exception("Protocol validator does not exist: ", entry["protocol_address"])

    prefix, address_bytes = bech32.bech32_decode(entry["consensus_address"])
    validator_acc_address = bech32.bech32_encode("kyve", address_bytes)

    verify_tx(api_endpoint, entry["proof_1"], entry["protocol_address"], validator_acc_address)
    verify_tx(api_endpoint, entry["proof_2"], validator_acc_address, entry["protocol_address"])


def verify_network(name, api_endpoint):
    status = {"correct": 0, "error": 0}
    for file in os.listdir("./" + name):
        try:
            proof = json.load(open("./{}/{}".format(name, file)))
            verify_proof(api_endpoint, proof)
            print("[{}]".format(name.title()), file, "✅")
            status["correct"] += 1

        except Exception as e:
            print("[{}]".format(name.title()), file, "❌")
            print(e)
            status["error"] += 1

    return status


status_mainnet = verify_network("mainnet", "https://api.kyve.network")
print("\n[Mainnet] Correct: {}, Error: {}".format(status_mainnet["correct"], status_mainnet["error"]))

if status_kaon["error"] != 0 or status_mainnet["error"] != 0:
    sys.exit(1)
