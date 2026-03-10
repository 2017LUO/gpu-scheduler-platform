package policy

type PlacementPolicy struct {
	SameNodeFirst      bool
	Binpack            bool
	Spread             bool
	TopologyAware      bool
	FragmentationAware bool
	RequireHealthyGPU  bool
}
