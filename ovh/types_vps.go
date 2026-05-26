package ovh

import (
	"fmt"
)

type VPSModel struct {
	Name                 string   `json:"name"`
	Offer                string   `json:"offer"`
	Memory               int      `json:"memory"`
	Vcore                int      `json:"vcore"`
	Version              string   `json:"version"`
	Disk                 int      `json:"disk"`
	Datacenter           []string `json:"datacenter,omitempty"`
	AvailableOptions     []string `json:"availableOptions,omitempty"`
	MaximumAdditionnalIp int      `json:"maximumAdditionnalIp,omitempty"`
}

type VPSOrderRuleDatacenter struct {
	Code               string `json:"code"`
	Datacenter         string `json:"datacenter"`
	DaysBeforeDelivery int    `json:"daysBeforeDelivery"`
	Status             string `json:"status"`
	LinuxStatus        string `json:"linuxStatus"`
	WindowsStatus      string `json:"windowsStatus"`
}

type VPSOrderRuleDatacenters struct {
	Datacenters []VPSOrderRuleDatacenter `json:"datacenters"`
}

type VPSOrderRuleOSChoice struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

type VPSOrderRuleOSChoices struct {
	Choices []VPSOrderRuleOSChoice `json:"choices"`
}

type VPS struct {
	Name               string   `json:"name"`
	Cluster            string   `json:"cluster"`
	Memory             int      `json:"memoryLimit"`
	NetbootMode        string   `json:"netbootMode"`
	Keymap             string   `json:"keymap"`
	Zone               string   `json:"zone"`
	State              string   `json:"state"`
	Vcore              int      `json:"vcore"`
	OfferType          string   `json:"offerType"`
	SlaMonitoring      bool     `json:"slaMonitoring"`
	DisplayName        string   `json:"displayName"`
	Model              VPSModel `json:"model"`
	IamResourceDetails `json:"iam"`
}

type VPSDatacenter struct {
	Name     string `json:"name"`
	Longname string `json:"longname"`
}

type VPSTask struct {
	Id       int64  `json:"id"`
	Date     string `json:"date"`
	Type     string `json:"type"`
	State    string `json:"state"`
	Progress int    `json:"progress"`
}

var vpsTaskStates = []string{
	"blocked",
	"cancelled",
	"doing",
	"done",
	"error",
	"paused",
	"todo",
	"waitingAck",
}

var vpsTaskTypes = []string{
	"addVeeamBackupJob",
	"changeRootPassword",
	"createSnapshot",
	"deleteSnapshot",
	"deliverVm",
	"getConsoleUrl",
	"internalTask",
	"migrate",
	"openConsoleAccess",
	"provisioningAdditionalIp",
	"reOpenVm",
	"rebootVm",
	"reinstallVm",
	"removeVeeamBackup",
	"rescheduleAutoBackup",
	"restoreFullVeeamBackup",
	"restoreVeeamBackup",
	"restoreVm",
	"revertSnapshot",
	"setMonitoring",
	"setNetboot",
	"startVm",
	"stopVm",
	"upgradeVm",
}

// VPSIp represents an IP attached to a VPS (vps.Ip dto).
// Only the `reverse` field is writable through PUT /vps/{sn}/ips/{ipAddress}.
type VPSIp struct {
	IpAddress   string  `json:"ipAddress"`
	Version     string  `json:"version"`
	Type        string  `json:"type"`
	Gateway     *string `json:"gateway,omitempty"`
	MacAddress  *string `json:"macAddress,omitempty"`
	Geolocation string  `json:"geolocation,omitempty"`
	Reverse     *string `json:"reverse,omitempty"`
}

// VPSIpReverseUpdateOpts is the body for PUT /vps/{sn}/ips/{ipAddress}.
// Only `reverse` is writable on the dto, so we only marshal that key.
type VPSIpReverseUpdateOpts struct {
	Reverse string `json:"reverse"`
}

// VPSStatusProbe represents one entry of vps.ip.ServiceStatus.
type VPSStatusProbe struct {
	Service string `json:"service"`
	Port    int    `json:"port,omitempty"`
	State   string `json:"state"`
}

// VPSDistributionTemplate mirrors vps.Template as returned by
// GET /vps/{serviceName}/distribution. Note that the API returns bitFormat
// as a JSON string ("32" or "64") even though it represents an integer.
type VPSDistributionTemplate struct {
	ID                int      `json:"id"`
	Name              string   `json:"name"`
	Distribution      string   `json:"distribution"`
	BitFormat         string   `json:"bitFormat"`
	AvailableLanguage []string `json:"availableLanguage"`
	Locale            string   `json:"locale"`
}

// VPSDistributionSoftware mirrors vps.Software as returned by
// GET /vps/{serviceName}/distribution/software/{softwareId}.
type VPSDistributionSoftware struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Status string `json:"status"`
}

func ovhvps_getType(offertype string, model_name string, model_version string) string {
	var offertypeToOfferPrefix = make(map[string]string)
	offertypeToOfferPrefix["cloud"] = "ceph-nvme"
	offertypeToOfferPrefix["cloudram"] = "ceph-nvme-ram"
	offertypeToOfferPrefix["ssd"] = "ssd"
	offertypeToOfferPrefix["classic"] = "classic"
	return (fmt.Sprintf("vps_%s_%s_%s", offertypeToOfferPrefix[offertype],
		model_name,
		model_version))
}

func strPtr(s string) *string {
	return &s
}
