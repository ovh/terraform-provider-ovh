package ovh

type VPSProperties struct {
	NetbootMode    *string `json:"netbootMode"`
	Keymap         *string `json:"keymap"`
	SlaMonitorting bool    `json:"slaMonitoring"`
	DisplayName    *string `json:"displayName"`
}

func strPtr(s string) *string {
	return &s
}
