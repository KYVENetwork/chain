package types

import (
	"fmt"

	"cosmossdk.io/math"
)

type (
	MultiCoinDistributionMap                 map[string][]MultiCoinDistributionPoolNormalizedEntry
	MultiCoinDistributionPoolNormalizedEntry struct {
		PoolId           uint64
		NormalizedWeight math.LegacyDec
	}
)

// ParseAndNormalizeMultiCoinDistributionMap turns the list structured policy map into a go map and
// normalizes the pool weights. Due to go maps indeterminism one can not store go-maps directly in state.
func ParseAndNormalizeMultiCoinDistributionMap(policy MultiCoinDistributionPolicy) (MultiCoinDistributionMap, error) {
	var distributionMap MultiCoinDistributionMap = make(map[string][]MultiCoinDistributionPoolNormalizedEntry)

	for _, denomEntry := range policy.Entries {
		distributionWeightsDuplicateCheck := make(map[uint64]struct{})

		totalWeight := math.LegacyNewDec(0)
		for _, weights := range denomEntry.PoolWeights {
			if !weights.Weight.IsPositive() {
				return nil, fmt.Errorf("invalid weight for pool id %d", weights.PoolId)
			}
			totalWeight = totalWeight.Add(weights.Weight)
			if _, ok := distributionWeightsDuplicateCheck[weights.PoolId]; ok {
				return nil, fmt.Errorf("duplicate distribution weight for pool id %d", weights.PoolId)
			}
			distributionWeightsDuplicateCheck[weights.PoolId] = struct{}{}
		}

		normalizedWeights := make([]MultiCoinDistributionPoolNormalizedEntry, 0)
		for _, weight := range denomEntry.PoolWeights {
			normalizedWeights = append(normalizedWeights, MultiCoinDistributionPoolNormalizedEntry{
				PoolId:           weight.PoolId,
				NormalizedWeight: weight.Weight.Quo(totalWeight),
			})
		}

		if _, ok := distributionMap[denomEntry.Denom]; !ok {
			distributionMap[denomEntry.Denom] = normalizedWeights
		} else {
			return nil, fmt.Errorf("duplicate entry for denom %s", denomEntry.Denom)
		}
	}

	return distributionMap, nil
}
