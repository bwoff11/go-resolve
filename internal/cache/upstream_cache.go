package cache

type UpstreamCache struct {
}

// NewRemoteCache initializes the remote DNS cache.
func NewUpstreamCache() *UpstreamCache {
	return &UpstreamCache{}
}
