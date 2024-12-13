package v2_0

type ValidatorMapping struct {
	Name             string
	ConsensusAddress string
	ProtocolAddress  string
}

var ValidatorMappings = []ValidatorMapping{
	ValidatorMapping{
		Name:             "",
		ConsensusAddress: "",
		ProtocolAddress:  "",
	},
}
