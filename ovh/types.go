package ovh

import (
	"fmt"
)

type IPPool struct {
	Network string `json:"network"`
	Region  string `json:"region"`
	Dhcp    bool   `json:"dhcp"`
	Start   string `json:"start"`
	End     string `json:"end"`
}

func (p *IPPool) String() string {
	return fmt.Sprintf("IPPool[Network: %s, Region: %s, Dhcp: %v, Start: %s, End: %s]", p.Network, p.Region, p.Dhcp, p.Start, p.End)
}

// Task Opts
type TaskOpts struct {
	ServiceName string `json:"serviceName"`
	TaskId      string `json:"taskId"`
}

type UnitAndValue struct {
	Unit  string `json:"unit"`
	Value int    `json:"value"`
}

func (v UnitAndValue) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["unit"] = v.Unit
	obj["value"] = v.Value

	return obj
}

type PublicCloudKubernetesClusterResponse struct {
	ControlPlaneIsUpToDate bool     `json:"controlPlaneIsUpToDate"`
	Id                     string   `json:"id"`
	IsUpToDate             bool     `json:"isUpToDate"`
	Name                   string   `json:"name"`
	NextUpgradeVersions    []string `json:"nextUpgradeVersions"`
	NodesUrl               string   `json:"nodesUrl"`
	Region                 string   `json:"region"`
	Status                 string   `json:"status"`
	UpdatePolicy           string   `json:"updatePolicy"`
	Url                    string   `json:"url"`
	Version                string   `json:"version"`
}

func (s *PublicCloudKubernetesClusterResponse) String() string {
	return fmt.Sprintf("%s(%s): %s", s.Name, s.Id, s.Status)
}

type PublicCloudKubernetesKubeConfigResponse struct {
	Content string `json:"content"`
}

type PublicCloudKubernetesNodeResponse struct {
	Id         string `json:"id"`
	ProjectId  string `json:"projectId"`
	InstanceId string `json:"instanceId"`
	IsUpToDate bool   `json:"isUpToDate"`
	Name       string `json:"name"`
	Flavor     string `json:"flavor"`
	Status     string `json:"status"`
	Version    string `json:"version"`
}

func (n *PublicCloudKubernetesNodeResponse) String() string {
	return fmt.Sprintf("%s(%s): %s", n.Name, n.Id, n.Status)
}

type PublicCloudKubernetesNodeCreationRequest struct {
	FlavorName string `json:"flavorName"`
	Name       string `json:"name"`
}
