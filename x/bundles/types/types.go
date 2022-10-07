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
