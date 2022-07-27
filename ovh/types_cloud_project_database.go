package ovh

import (
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/go-ovh/ovh"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

type CloudProjectDatabaseCreateOpts struct {
	Description  *string                          `json:"description,omitempty"`
	NetworkId    *string                          `json:"networkId,omitempty"`
	NodesPattern CloudProjectDatabaseNodesPattern `json:"nodesPattern,omitempty"`
	Plan         string                           `json:"plan"`
	SubnetId     *string                          `json:"subnetId,omitempty"`
	Version      string                           `json:"version"`
}

type CloudProjectDatabaseUpdateOpts struct {
	Description string `json:"description,omitempty"`
	Flavor      string `json:"flavor,omitempty"`
	NodeNumber  int    `json:"nodeNumber,omitempty"`
	Plan        string `json:"plan,omitempty"`
	Version     string `json:"version,omitempty"`
}

type CloudProjectDatabaseNodesPattern struct {
	Flavor string `json:"flavor"`
	Number int    `json:"number"`
	Region string `json:"region"`
}

type CloudProjectDatabaseNodes struct {
	NetworkId *string `json:"networkId,omitempty"`
	Region    string  `json:"region"`
	SubnetId  *string `json:"subnetId,omitempty"`
}

type CloudProjectDatabaseResponse struct {
	BackupTime      string                         `json:"backupTime"`
	CreatedAt       string                         `json:"createdAt"`
	Description     string                         `json:"description"`
	Endpoints       []CloudProjectDatabaseEndpoint `json:"endpoints"`
	Flavor          string                         `json:"flavor"`
	Id              string                         `json:"id"`
	MaintenanceTime string                         `json:"maintenanceTime"`
	NetworkId       *string                        `json:"networkId,omitempty"`
	NetworkType     string                         `json:"networkType"`
	Plan            string                         `json:"plan"`
	NodeNumber      int                            `json:"nodeNumber"`
	Region          string                         `json:"region"`
	Status          string                         `json:"status"`
	SubnetId        *string                        `json:"subnetId,omitempty"`
	Version         string                         `json:"version"`
}

type CloudProjectDatabaseEndpoint struct {
	Component string  `json:"component"`
	Domain    string  `json:"domain"`
	Path      *string `json:"path,omitempty"`
	Port      *int    `json:"port,omitempty"`
	Scheme    *string `json:"scheme,omitempty"`
	Ssl       *bool   `json:"ssl,omitempty"`
	SslMode   *string `json:"sslMode,omitempty"`
	Uri       *string `json:"uri,omitempty"`
}

func (opts *CloudProjectDatabaseCreateOpts) FromResource(d *schema.ResourceData) (error, *CloudProjectDatabaseCreateOpts) {
	opts.Description = helpers.GetNilStringPointerFromData(d, "description")
	opts.Plan = d.Get("plan").(string)

	nodes := []CloudProjectDatabaseNodes{}
	nbOfNodes := d.Get("nodes.#").(int)
	for i := 0; i < nbOfNodes; i++ {
		nodes = append(nodes, *(&CloudProjectDatabaseNodes{}).FromResourceWithPath(d, fmt.Sprintf("nodes.%d", i)))
	}

	if err := checkNodesEquality(nodes); err != nil {
		return err, nil
	}

	opts.NodesPattern = CloudProjectDatabaseNodesPattern{
		Flavor: d.Get("flavor").(string),
		Region: nodes[0].Region,
		Number: nbOfNodes,
	}

	opts.NetworkId = nodes[0].NetworkId
	opts.SubnetId = nodes[0].SubnetId
	opts.Version = d.Get("version").(string)

	return nil, opts
}

func (opts *CloudProjectDatabaseNodes) FromResourceWithPath(d *schema.ResourceData, path string) *CloudProjectDatabaseNodes {
	opts.Region = d.Get(fmt.Sprintf("%s.region", path)).(string)
	opts.NetworkId = helpers.GetNilStringPointerFromData(d, fmt.Sprintf("%s.network_id", path))
	opts.SubnetId = helpers.GetNilStringPointerFromData(d, fmt.Sprintf("%s.subnet_id", path))

	return opts
}

func (opts *CloudProjectDatabaseUpdateOpts) FromResource(d *schema.ResourceData) (error, *CloudProjectDatabaseUpdateOpts) {
	opts.Description = d.Get("description").(string)
	opts.Plan = d.Get("plan").(string)

	nodes := []CloudProjectDatabaseNodes{}
	nbOfNodes := d.Get("nodes.#").(int)
	for i := 0; i < nbOfNodes; i++ {
		nodes = append(nodes, *(&CloudProjectDatabaseNodes{}).FromResourceWithPath(d, fmt.Sprintf("nodes.%d", i)))
	}

	if err := checkNodesEquality(nodes); err != nil {
		return err, nil
	}
	opts.Flavor = d.Get("flavor").(string)
	opts.NodeNumber = nbOfNodes
	opts.Version = d.Get("version").(string)

	return nil, opts
}

func (v CloudProjectDatabaseResponse) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["backup_time"] = v.BackupTime
	obj["created_at"] = v.CreatedAt
	obj["description"] = v.Description
	obj["id"] = v.Id

	var endpoints []map[string]interface{}
	for _, e := range v.Endpoints {
		endpoints = append(endpoints, e.ToMap())
	}
	obj["endpoints"] = endpoints

	obj["flavor"] = v.Flavor
	obj["maintenance_time"] = v.MaintenanceTime
	obj["network_type"] = v.NetworkType

	var nodes []map[string]interface{}
	for i := 0; i < v.NodeNumber; i++ {
		node := CloudProjectDatabaseNodes{
			Region:    v.Region,
			NetworkId: v.NetworkId,
			SubnetId:  v.SubnetId,
		}
		nodes = append(nodes, node.ToMap())
	}
	obj["nodes"] = nodes

	obj["plan"] = v.Plan
	obj["status"] = v.Status
	obj["version"] = v.Version

	return obj
}

func (v CloudProjectDatabaseNodes) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})

	if v.NetworkId != nil {
		obj["network_id"] = v.NetworkId
	}

	obj["region"] = v.Region

	if v.SubnetId != nil {
		obj["subnet_id"] = v.SubnetId
	}

	return obj
}

func (v CloudProjectDatabaseEndpoint) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})

	obj["component"] = v.Component
	obj["domain"] = v.Domain

	if v.Path != nil {
		obj["path"] = v.Path
	}

	if v.Port != nil {
		obj["port"] = v.Port
	}

	if v.Scheme != nil {
		obj["scheme"] = v.Scheme
	}

	if v.Ssl != nil {
		obj["ssl"] = v.Ssl
	}

	if v.SslMode != nil {
		obj["ssl_mode"] = v.SslMode
	}

	if v.Uri != nil {
		obj["uri"] = v.Uri
	}

	return obj
}

func (s *CloudProjectDatabaseResponse) String() string {
	return fmt.Sprintf("%s(%s): %s", s.Description, s.Id, s.Status)
}

func (s *CloudProjectDatabaseCreateOpts) String() string {
	return fmt.Sprintf("%s: %s", *s.Description, s.Version)
}

// This make sure Nodes are homogenous.
// When multi region cluster will be available the check will be done on API side
func checkNodesEquality(nodes []CloudProjectDatabaseNodes) error {
	if len(nodes) == 0 {
		return errors.New("node list empty")
	}
	if len(nodes) == 1 {
		return nil
	}

	networkId := ""
	if nodes[0].NetworkId != nil {
		networkId = *nodes[0].NetworkId
	}

	region := nodes[0].Region

	subnetId := ""
	if nodes[0].SubnetId != nil {
		subnetId = *nodes[0].SubnetId
	}

	for _, n := range nodes[1:] {
		if (n.NetworkId == nil && networkId != "") || (n.NetworkId != nil && networkId != *n.NetworkId) {
			return errors.New("network_id is not the same across nodes")
		}
		if region != n.Region {
			return errors.New("region is not the same across nodes")
		}
		if (n.SubnetId == nil && subnetId != "") || (n.SubnetId != nil && subnetId != *n.SubnetId) {
			return errors.New("subnet_id is not the same across nodes")
		}
	}
	return nil
}

func waitForCloudProjectDatabaseReady(client *ovh.Client, serviceName, engine string, databaseId string, timeOut time.Duration, delay time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "CREATING", "UPDATING"},
		Target:  []string{"READY"},
		Refresh: func() (interface{}, string, error) {
			res := &CloudProjectDatabaseResponse{}
			endpoint := fmt.Sprintf("/cloud/project/%s/database/%s/%s",
				url.PathEscape(serviceName),
				url.PathEscape(engine),
				url.PathEscape(databaseId),
			)
			err := client.Get(endpoint, res)
			if err != nil {
				return res, "", err
			}

			return res, res.Status, nil
		},
		Timeout:    timeOut,
		Delay:      delay,
		MinTimeout: 10 * time.Second,
	}

	_, err := stateConf.WaitForState()
	return err
}

func waitForCloudProjectDatabaseDeleted(client *ovh.Client, serviceName, engine string, databaseId string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"DELETING"},
		Target:  []string{"DELETED"},
		Refresh: func() (interface{}, string, error) {
			res := &CloudProjectDatabaseResponse{}
			endpoint := fmt.Sprintf("/cloud/project/%s/database/%s/%s",
				url.PathEscape(serviceName),
				url.PathEscape(engine),
				url.PathEscape(databaseId),
			)
			err := client.Get(endpoint, res)
			if err != nil {
				if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
					return res, "DELETED", nil
				} else {
					return res, "", err
				}
			}

			return res, res.Status, nil
		},
		Timeout:      30 * time.Minute,
		Delay:        30 * time.Second,
		MinTimeout:   3 * time.Second,
		PollInterval: 20 * time.Second,
	}

	_, err := stateConf.WaitForState()
	return err
}

type CloudProjectDatabaseIpRestrictionResponse struct {
	Description string `json:"description"`
	Ip          string `json:"ip"`
	Status      string `json:"status"`
}

func (p *CloudProjectDatabaseIpRestrictionResponse) String() string {
	return fmt.Sprintf(
		"IP: %s, Status: %s, Description: %s",
		p.Ip,
		p.Status,
		p.Description,
	)
}

func (v CloudProjectDatabaseIpRestrictionResponse) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})

	obj["description"] = v.Description
	obj["ip"] = v.Ip
	obj["status"] = v.Status

	return obj
}

type CloudProjectDatabaseIpRestrictionCreateOpts struct {
	Description string `json:"description"`
	Ip          string `json:"ip"`
}

func (opts *CloudProjectDatabaseIpRestrictionCreateOpts) FromResource(d *schema.ResourceData) *CloudProjectDatabaseIpRestrictionCreateOpts {
	opts.Description = d.Get("description").(string)
	opts.Ip = d.Get("ip").(string)
	return opts
}

type CloudProjectDatabaseIpRestrictionUpdateOpts struct {
	Description string `json:"description"`
}

func (opts *CloudProjectDatabaseIpRestrictionUpdateOpts) FromResource(d *schema.ResourceData) *CloudProjectDatabaseIpRestrictionUpdateOpts {
	opts.Description = d.Get("description").(string)
	return opts
}
