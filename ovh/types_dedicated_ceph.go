package ovh

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

type DedicatedCephCrushTunable string
type DedicatedCephState string
type DedicatedCephACLType string
type DedicatedCephStatus string

const (
	CrushTunableOptimal  DedicatedCephCrushTunable = "OPTIMAL"
	CrushTunableDefault  DedicatedCephCrushTunable = "DEFAULT"
	CrushTunableLegacy   DedicatedCephCrushTunable = "LEGACY"
	CrushTunableBobtail  DedicatedCephCrushTunable = "BOBTAIL"
	CrushTunableArgonaut DedicatedCephCrushTunable = "ARGONAUT"
	CrushTunableFirefly  DedicatedCephCrushTunable = "FIREFLY"
	CrushTunableHammer   DedicatedCephCrushTunable = "HAMMER"
	CrushTunablEjewel    DedicatedCephCrushTunable = "JEWEL"
	StateActive          DedicatedCephState        = "ACTIVE"
	StateSuspended       DedicatedCephState        = "SUSPENDED"
	ACLTypeIPv4          DedicatedCephACLType      = "IPV4"
	ACLTypeIPv6          DedicatedCephACLType      = "IPV6"
	StatusCreating       DedicatedCephStatus       = "CREATING"
	StatusInstalled      DedicatedCephStatus       = "INSTALLED"
	StatusDeleting       DedicatedCephStatus       = "DELETING"
	StatusDeleted        DedicatedCephStatus       = "DELETED"
	StatusTaskInProgress DedicatedCephStatus       = "TASK_IN_PROGRESS"
)

type DedicatedCeph struct {
	ServiceName   string                    `json:"serviceName"`
	CephMonitors  []string                  `json:"cephMons"`
	CephVersion   string                    `json:"cephVersion"`
	CrushTunables DedicatedCephCrushTunable `json:"crushTunables"`
	Label         string                    `json:"label"`
	Region        string                    `json:"region"`
	Size          float32                   `json:"size"`
	State         DedicatedCephState        `json:"state"`
	Status        DedicatedCephStatus       `json:"status"`
}

type DedicatedCephACL struct {
	Id      int                  `json:"id"`
	Family  DedicatedCephACLType `json:"family"`
	Netmask string               `json:"netmask"`
	Network string               `json:"network"`
}

type DedicatedCephACLCreateOpts struct {
	AclList []string `json:"aclList"`
}

type DedicatedCephTask struct {
	Name       string `json:"name"`
	State      string `json:"state"`
	FinishDate string `json:"finishDate"`
	Type       string `json:"type"`
	CreateDate string `json:"createDate"`
}

func (opts *DedicatedCephACLCreateOpts) FromResource(d *schema.ResourceData) *DedicatedCephACLCreateOpts {
	network := helpers.GetNilStringPointerFromData(d, "network")
	netmask := helpers.GetNilStringPointerFromData(d, "netmask")
	opts.AclList = []string{fmt.Sprintf("%s/%s", *network, *netmask)}
	return opts
}
