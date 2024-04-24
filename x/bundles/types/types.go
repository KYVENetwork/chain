package types

type VoteDistribution struct {
	// valid ...
	Valid uint64
	// invalid ...
	Invalid uint64
	// abstain ...
	Abstain uint64
	// total ...
	Total uint64
	// status ...
	Status BundleStatus
}

type BundleReward struct {
	// treasury ...
	Treasury uint64
	// uploader ...
	Uploader uint64
	// delegation ...
	Delegation uint64
	// total ...
	Total uint64
}

// GetMap converts to array to a go map which return the upgrade-height for each version.
// e.g. the schema changed from v1 to v2 at block 1,000.
// then: GetMap()[2] = 1000
// Version 1 start at 0 and is not encoded in the map
func (bundleVersionMap BundleVersionMap) GetMap() (versionMap map[int32]uint64) {
	versionMap = make(map[int32]uint64, 0)
	for _, entry := range bundleVersionMap.Versions {
		versionMap[entry.Version] = entry.Height
	}
	return
}

type TallyResultStatus uint32

const (
	TallyResultValid TallyResultStatus = iota
	TallyResultInvalid
	TallyResultNoQuorum
)

type TallyResult struct {
	Status           TallyResultStatus
	VoteDistribution VoteDistribution
	FundersPayout    uint64
	InflationPayout  uint64
	BundleReward     BundleReward
}
