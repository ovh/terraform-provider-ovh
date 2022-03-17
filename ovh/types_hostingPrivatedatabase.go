package ovh

type HostingPrivateDatabase struct {
	ServiceName    string        `json:"serviceName"`
	Cpu            int           `json:"cpu"`
	Datacenter     string        `json:"datacenter"`
	DisplayName    string        `json:"displayName"`
	Hostname       string        `json:"hostname"`
	HostnameFtp    string        `json:"hostnameFtp"`
	Infrastructure string        `json:"infrastructure"`
	Offer          string        `json:"offer"`
	Port           int           `json:"port"`
	PortFtp        int           `json:"portFtp"`
	QuotaSize      *UnitAndValue `json:"quotasize"`
	QuotaUsed      *UnitAndValue `json:"quotaused"`
	Ram            *UnitAndValue `json:"ram"`
	Server         string        `json:"server"`
	State          string        `json:"state"`
	Type           string        `json:"type"`
	Version        string        `json:"version"`
	VersionLabel   string        `json:"versionLabel"`
	VersionNumber  float64       `json:"versionNumber"`
}
