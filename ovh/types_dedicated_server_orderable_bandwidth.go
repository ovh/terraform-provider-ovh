package ovh

type DedicatedServerOrderableBandwidth struct {
	Orderable bool  `json:"orderable"`
	Platinium []int `json:"platinium"`
	Ultimate  []int `json:"ultimate"`
	Premium   []int `json:"premium"`
}

type DedicatedServerOrderableBandwidthVrack struct {
	Orderable bool  `json:"orderable"`
	Vrack     []int `json:"vrack"`
}
