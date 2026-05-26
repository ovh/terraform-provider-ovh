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
	Datacenter           []string `json:"datacenter"`
	AvailableOptions     []string `json:"availableOptions"`
	MaximumAdditionnalIp int      `json:"maximumAdditionnalIp"`
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
