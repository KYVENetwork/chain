package v2_0

/*

Protocol-Consensus Validator Linking:

1. 	Fill out the entry below following the example.

2. 	Send 1$KYVE from the protocol-address to the consensus-validator-operator address using the memo "Shared-Staking"
	and put the tx-hash in Proof1.

3.  Send 1$KYVE from the consensus-validator-operator address to the protocol address using the memo "Shared-Staking"
	and put the tx-hash in Proof2.

4.	Submit a Pull-Request to https://github.com/KYVENetwork/chain

*/

type ValidatorMapping struct {
	Name             string
	ConsensusAddress string
	ProtocolAddress  string
	Proof1           string
	Proof2           string
}

var ValidatorMappingsMainnet = []ValidatorMapping{
	{
		// human-readable name, only used for logging
		Name: "",
		// kyvevaloper... address of the chain node
		ConsensusAddress: "",
		// kyve... address of the protocol node
		ProtocolAddress: "",
		// Proof TX-Hash 1, transferring 1 $KYVE from the protocol-address to the operator address
		// using "Shared Staking" as memo.
		Proof1: "",
		// Proof TX-Hash 2, transferring 1 $KYVE from the operator address to the protocol-address
		// using "Shared Staking" as memo.
		Proof2: "",
	},
}

var ValidatorMappingsKaon = []ValidatorMapping{
	{
		// human-readable name, only used for logging
		Name: "",
		// kyvevaloper... address of the chain node
		ConsensusAddress: "",
		// kyve... address of the protocol node
		ProtocolAddress: "",
		// Proof TX-Hash 1, transferring 1 $KYVE from the protocol-address to the operator address
		// using "Shared Staking" as memo.
		Proof1: "",
		// Proof TX-Hash 2, transferring 1 $KYVE from the operator address to the protocol-address
		// using "Shared Staking" as memo.
		Proof2: "",
	},
}
