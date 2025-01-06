## Protocol-Consensus Validator Linking

For the Shared Staking update [https://www.mintscan.io/kyve/proposals/43](https://www.mintscan.io/kyve/proposals/43)
it is necessary that every validator links their protocol and consensus validator.

All delegators are then transferred from the protocol validator to the consensus
validator during the upgrade.

If a protocol validator does not link to a chain validator before the upgrade is finalized,
all stake is returned to the original delegators during the migration.

### Steps

1. 	Enter the `mainnet`-directory and copy the `example-validator.json` config file and name it after your validator.

2.  Fill out the `name`, `protocol_address` and `consensus_address`

3.  Send 1 $KYVE from the protocol-address to the consensus-validator-operator address using the memo "Shared-Staking"
    and put the tx-hash in proof_1.

4.  Send 1 $KYVE from the consensus-validator-operator address to the protocol address using the memo "Shared-Staking"
    and put the tx-hash in proof_2.

5.  Submit a Pull-Request to https://github.com/KYVENetwork/chain

6.  (Optional) Perform the same steps for the `kaon` directory with your Kaon validators.

## General Upgrade Procedure

All pending protocol commission will be claimed and returned to the validators
during the upgrade.

During the upgrade, all validators will be kicked out of all pools and
need to manually rejoin after the upgrade. This is due to the fact, that 
every validator now has to specify a commission per pool and a stake-fraction.
The stake-fraction is a percentage of how much of the chain-stake a validator
wants to use for a specific pool.

## Questions

If you have any questions feel free to submit an issue or ask them directly while
creating the pull-request.
