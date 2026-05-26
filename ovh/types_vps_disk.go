package ovh

// VPSDisk represents a VPS disk as returned by /vps/{serviceName}/disks/{id}.
type VPSDisk struct {
	ID                    int64  `json:"id"`
	ServiceName           string `json:"serviceName"`
	Type                  string `json:"type"`
	State                 string `json:"state"`
	Size                  int64  `json:"size"`
	BandwidthLimit        int64  `json:"bandwidthLimit"`
	Monitoring            bool   `json:"monitoring"`
	LowFreeSpaceThreshold *int64 `json:"lowFreeSpaceThreshold,omitempty"`
}

// VPSDiskUpdateOpts is the PUT body for /vps/{serviceName}/disks/{id}.
// Only `monitoring` and `lowFreeSpaceThreshold` are writable.
type VPSDiskUpdateOpts struct {
	Monitoring            bool   `json:"monitoring"`
	LowFreeSpaceThreshold *int64 `json:"lowFreeSpaceThreshold,omitempty"`
}

// VPSDiskUsage is the response from /vps/{serviceName}/disks/{id}/use.
type VPSDiskUsage struct {
	Unit  string  `json:"unit"`
	Value float64 `json:"value"`
}
