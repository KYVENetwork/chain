package v2_0

import (
	"embed"
	"encoding/json"
	"fmt"
)

//go:embed validator-proofs/*
var validatorProofs embed.FS

func init() {
	parseDirectory := func(directory string) []ValidatorMapping {
		dir, err := validatorProofs.ReadDir(directory)
		if err != nil {
			panic(err)
		}

		result := make([]ValidatorMapping, 0)
		for _, file := range dir {
			readFile, err := validatorProofs.ReadFile(fmt.Sprintf("%s/%s", directory, file.Name()))
			if err != nil {
				panic(err)
			}

			var proof ValidatorMapping
			if err = json.Unmarshal(readFile, &proof); err != nil {
				panic(err)
			}
			result = append(result, proof)
		}
		return result
	}

	ValidatorMappingsMainnet = parseDirectory("validator-proofs/mainnet")
	ValidatorMappingsKaon = parseDirectory("validator-proofs/kaon")
	ValidatorMappingsKorellia = parseDirectory("validator-proofs/korellia")
}

type ValidatorMapping struct {
	Name             string `json:"name"`
	ConsensusAddress string `json:"consensus_address"`
	ProtocolAddress  string `json:"protocol_address"`
	Proof1           string `json:"proof_1"`
	Proof2           string `json:"proof_2"`
}

var ValidatorMappingsMainnet []ValidatorMapping

var ValidatorMappingsKaon []ValidatorMapping

var ValidatorMappingsKorellia []ValidatorMapping
