package types

import (
	"fmt"

	"cosmossdk.io/math"
)

type (
	ComplianceMultiCoinMap                 map[string][]ComplainceMultiCoinPoolNormalizedEntry
	ComplainceMultiCoinPoolNormalizedEntry struct {
		PoolId           uint64
		NormalizedWeight math.LegacyDec
	}
)

func ParseMultiCoinComplianceMap(policy MultiCoinRefundPolicy) (ComplianceMultiCoinMap, error) {
	var compliance ComplianceMultiCoinMap = make(map[string][]ComplainceMultiCoinPoolNormalizedEntry)

	for _, denomEntry := range policy.Entries {
		complianceWeightsDuplicateCheck := make(map[uint64]struct{})

		totalWeight := math.LegacyNewDec(0)
		for _, weights := range denomEntry.PoolWeights {
			totalWeight = totalWeight.Add(weights.Weight)
			if _, ok := complianceWeightsDuplicateCheck[weights.PoolId]; ok {
				return nil, fmt.Errorf("duplicate compliance weight for pool id %d", weights.PoolId)
			}
		}

		normalizedWeights := make([]ComplainceMultiCoinPoolNormalizedEntry, 0)
		for _, weight := range denomEntry.PoolWeights {
			normalizedWeights = append(normalizedWeights, ComplainceMultiCoinPoolNormalizedEntry{
				PoolId:           weight.PoolId,
				NormalizedWeight: weight.Weight.Quo(totalWeight),
			})
		}

		if _, ok := compliance[denomEntry.Denom]; !ok {
			compliance[denomEntry.Denom] = normalizedWeights
		} else {
			return nil, fmt.Errorf("duplicate entry for denom %s", denomEntry.Denom)
		}
	}

	return compliance, nil
}
