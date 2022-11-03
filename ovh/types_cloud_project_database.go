package ovh

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/go-ovh/ovh"
	"github.com/ybriffa/rfc3339"
)

type CloudProjectDatabaseResponse struct {
	AclsEnabled     bool                           `json:"aclsEnabled"`
	BackupTime      string                         `json:"backupTime"`
	CreatedAt       string                         `json:"createdAt"`
	Description     string                         `json:"description"`
	Endpoints       []CloudProjectDatabaseEndpoint `json:"endpoints"`
	Flavor          string                         `json:"flavor"`
	Id              string                         `json:"id"`
	MaintenanceTime string                         `json:"maintenanceTime"`
	NetworkId       string                         `json:"networkId"`
	NetworkType     string                         `json:"networkType"`
	Plan            string                         `json:"plan"`
	NodeNumber      int                            `json:"nodeNumber"`
	Region          string                         `json:"region"`
	RestApi         bool                           `json:"restApi"`
	Status          string                         `json:"status"`
	SubnetId        string                         `json:"subnetId"`
	Version         string                         `json:"version"`
	Disk            CloudProjectDatabaseDisk       `json:"disk"`
}

func (s *CloudProjectDatabaseResponse) String() string {
	return fmt.Sprintf("%s(%s): %s", s.Description, s.Id, s.Status)
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
	obj["kafka_rest_api"] = v.RestApi
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

	obj["opensearch_acls_enabled"] = v.AclsEnabled
	obj["plan"] = v.Plan
	obj["status"] = v.Status
	obj["version"] = v.Version
	obj["disk_size"] = v.Disk.Size
	obj["disk_type"] = v.Disk.Type

	return obj
}

type CloudProjectDatabaseEndpoint struct {
	Component string `json:"component"`
	Domain    string `json:"domain"`
	Path      string `json:"path"`
	Port      int    `json:"port"`
	Scheme    string `json:"scheme"`
	Ssl       bool   `json:"ssl"`
	SslMode   string `json:"sslMode"`
	Uri       string `json:"uri"`
}

func (v CloudProjectDatabaseEndpoint) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})

	obj["component"] = v.Component
	obj["domain"] = v.Domain
	obj["path"] = v.Path
	obj["port"] = v.Port
	obj["scheme"] = v.Scheme
	obj["ssl"] = v.Ssl
	obj["ssl_mode"] = v.SslMode
	obj["uri"] = v.Uri

	return obj
}

type CloudProjectDatabaseNodes struct {
	NetworkId string `json:"networkId,omitempty"`
	Region    string `json:"region"`
	SubnetId  string `json:"subnetId,omitempty"`
}

func (opts *CloudProjectDatabaseNodes) FromResourceWithPath(d *schema.ResourceData, path string) *CloudProjectDatabaseNodes {
	opts.Region = d.Get(fmt.Sprintf("%s.region", path)).(string)
	opts.NetworkId = d.Get(fmt.Sprintf("%s.network_id", path)).(string)
	opts.SubnetId = d.Get(fmt.Sprintf("%s.subnet_id", path)).(string)

	return opts
}

func (v CloudProjectDatabaseNodes) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})

	obj["network_id"] = v.NetworkId
	obj["region"] = v.Region
	obj["subnet_id"] = v.SubnetId

	return obj
}

type CloudProjectDatabaseCreateOpts struct {
	Description  string                           `json:"description,omitempty"`
	NetworkId    string                           `json:"networkId,omitempty"`
	NodesPattern CloudProjectDatabaseNodesPattern `json:"nodesPattern,omitempty"`
	Disk         CloudProjectDatabaseDisk         `json:"disk,omitempty"`
	Plan         string                           `json:"plan"`
	SubnetId     string                           `json:"subnetId,omitempty"`
	Version      string                           `json:"version"`
}

type CloudProjectDatabaseDisk struct {
	Type string `json:"type,omitempty"`
	Size int    `json:"size,omitempty"`
}

func validateCloudProjectDatabaseDiskSize(v interface{}, k string) (ws []string, errors []error) {
	errors = validateIsSupEqual(v.(int), 0)
	return
}

type CloudProjectDatabaseNodesPattern struct {
	Flavor string `json:"flavor"`
	Number int    `json:"number"`
	Region string `json:"region"`
}

func (opts *CloudProjectDatabaseCreateOpts) FromResource(d *schema.ResourceData) (error, *CloudProjectDatabaseCreateOpts) {
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

	opts.NodesPattern = CloudProjectDatabaseNodesPattern{
		Flavor: d.Get("flavor").(string),
		Region: nodes[0].Region,
		Number: nbOfNodes,
	}

	opts.NetworkId = nodes[0].NetworkId
	opts.SubnetId = nodes[0].SubnetId
	opts.Version = d.Get("version").(string)
	opts.Disk = CloudProjectDatabaseDisk{Size: d.Get("disk_size").(int)}
	return nil, opts
}

type CloudProjectDatabaseUpdateOpts struct {
	AclsEnabled bool                     `json:"aclsEnabled,omitempty"`
	Description string                   `json:"description,omitempty"`
	Flavor      string                   `json:"flavor,omitempty"`
	Plan        string                   `json:"plan,omitempty"`
	RestApi     bool                     `json:"restApi,omitempty"`
	Version     string                   `json:"version,omitempty"`
	Disk        CloudProjectDatabaseDisk `json:"disk,omitempty"`
}

func (opts *CloudProjectDatabaseUpdateOpts) FromResource(d *schema.ResourceData) (error, *CloudProjectDatabaseUpdateOpts) {
	engine := d.Get("engine").(string)
	if engine == "opensearch" {
		opts.AclsEnabled = d.Get("opensearch_acls_enabled").(bool)
	}
	if engine == "kafka" {
		opts.RestApi = d.Get("kafka_rest_api").(bool)
	}

	opts.Description = d.Get("description").(string)
	opts.Plan = d.Get("plan").(string)
	opts.Flavor = d.Get("flavor").(string)
	opts.Version = d.Get("version").(string)
	opts.Disk = CloudProjectDatabaseDisk{Size: d.Get("disk_size").(int)}
	return nil, opts
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

	networkId := nodes[0].NetworkId

	region := nodes[0].Region

	subnetId := nodes[0].SubnetId

	for _, n := range nodes[1:] {
		if n.NetworkId != networkId {
			return errors.New("network_id is not the same across nodes")
		}
		if region != n.Region {
			return errors.New("region is not the same across nodes")
		}
		if n.SubnetId != subnetId {
			return errors.New("subnet_id is not the same across nodes")
		}
	}
	return nil
}

func waitForCloudProjectDatabaseReady(client *ovh.Client, serviceName, engine string, databaseId string, timeOut time.Duration) error {
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
		Delay:      30 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	_, err := stateConf.WaitForState()
	return err
}

func waitForCloudProjectDatabaseDeleted(client *ovh.Client, serviceName, engine string, databaseId string, timeOut time.Duration) error {
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
				}
				return res, "", err
			}

			return res, res.Status, nil
		},
		Timeout:      timeOut,
		Delay:        30 * time.Second,
		PollInterval: 20 * time.Second,
	}

	_, err := stateConf.WaitForState()
	return err
}

// Capabilities

type CloudProjectDatabaseCapabilitiesEngine struct {
	DefaultVersion string   `json:"defaultVersion"`
	Description    string   `json:"description"`
	Name           string   `json:"name"`
	SslModes       []string `json:"sslModes"`
	Versions       []string `json:"versions"`
}

func (v CloudProjectDatabaseCapabilitiesEngine) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["default_version"] = v.DefaultVersion
	obj["description"] = v.Description
	obj["name"] = v.Name
	obj["ssl_modes"] = v.SslModes
	obj["versions"] = v.Versions
	return obj
}

type CloudProjectDatabaseCapabilitiesFlavor struct {
	Core    int    `json:"core"`
	Memory  int    `json:"memory"`
	Name    string `json:"name"`
	Storage int    `json:"storage"`
}

func (v CloudProjectDatabaseCapabilitiesFlavor) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["core"] = v.Core
	obj["memory"] = v.Memory
	obj["name"] = v.Name
	obj["storage"] = v.Storage
	return obj
}

type CloudProjectDatabaseCapabilitiesOption struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func (v CloudProjectDatabaseCapabilitiesOption) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["name"] = v.Name
	obj["type"] = v.Type
	return obj
}

type CloudProjectDatabaseCapabilitiesPlan struct {
	BackupRetention string `json:"backupRetention"`
	Description     string `json:"description"`
	Name            string `json:"name"`
}

func (v CloudProjectDatabaseCapabilitiesPlan) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["backup_retention"] = v.BackupRetention
	obj["description"] = v.Description
	obj["name"] = v.Name
	return obj
}

type CloudProjectDatabaseCapabilitiesResponse struct {
	Engines []CloudProjectDatabaseCapabilitiesEngine `json:"engines"`
	Flavors []CloudProjectDatabaseCapabilitiesFlavor `json:"flavors"`
	Options []CloudProjectDatabaseCapabilitiesOption `json:"options"`
	Plans   []CloudProjectDatabaseCapabilitiesPlan   `json:"plans"`
}

func (v CloudProjectDatabaseCapabilitiesResponse) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	var engines []map[string]interface{}
	for _, e := range v.Engines {
		engines = append(engines, e.ToMap())
	}
	obj["engines"] = engines

	var flavors []map[string]interface{}
	for _, e := range v.Flavors {
		flavors = append(flavors, e.ToMap())
	}
	obj["flavors"] = flavors

	var options []map[string]interface{}
	for _, e := range v.Options {
		options = append(options, e.ToMap())
	}
	obj["options"] = options

	var plans []map[string]interface{}
	for _, e := range v.Plans {
		plans = append(plans, e.ToMap())
	}
	obj["plans"] = plans

	return obj
}

// IP Restriction

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
	Description string `json:"description,omitempty"`
	Ip          string `json:"ip"`
}

func (opts *CloudProjectDatabaseIpRestrictionCreateOpts) FromResource(d *schema.ResourceData) *CloudProjectDatabaseIpRestrictionCreateOpts {
	opts.Description = d.Get("description").(string)
	opts.Ip = d.Get("ip").(string)
	return opts
}

type CloudProjectDatabaseIpRestrictionUpdateOpts struct {
	Description string `json:"description,omitempty"`
}

func (opts *CloudProjectDatabaseIpRestrictionUpdateOpts) FromResource(d *schema.ResourceData) *CloudProjectDatabaseIpRestrictionUpdateOpts {
	opts.Description = d.Get("description").(string)
	return opts
}

func waitForCloudProjectDatabaseIpRestrictionReady(client *ovh.Client, serviceName, engine string, databaseId string, ip string, timeOut time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "CREATING", "UPDATING"},
		Target:  []string{"READY"},
		Refresh: func() (interface{}, string, error) {
			res := &CloudProjectDatabaseIpRestrictionResponse{}
			endpoint := fmt.Sprintf("/cloud/project/%s/database/%s/%s/ipRestriction/%s",
				url.PathEscape(serviceName),
				url.PathEscape(engine),
				url.PathEscape(databaseId),
				url.PathEscape(ip),
			)
			err := client.Get(endpoint, res)
			if err != nil {
				return res, "", err
			}

			return res, res.Status, nil
		},
		Timeout:    timeOut,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForState()
	return err
}

func waitForCloudProjectDatabaseIpRestrictionDeleted(client *ovh.Client, serviceName, engine string, databaseId string, ip string, timeOut time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"DELETING"},
		Target:  []string{"DELETED"},
		Refresh: func() (interface{}, string, error) {
			res := &CloudProjectDatabaseIpRestrictionResponse{}
			endpoint := fmt.Sprintf("/cloud/project/%s/database/%s/%s/ipRestriction/%s",
				url.PathEscape(serviceName),
				url.PathEscape(engine),
				url.PathEscape(databaseId),
				url.PathEscape(ip),
			)
			err := client.Get(endpoint, res)
			if err != nil {
				if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
					return res, "DELETED", nil
				}
				return res, "", err
			}

			return res, res.Status, nil
		},
		Timeout:    timeOut,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForState()
	return err
}

// User

type CloudProjectDatabaseUserResponse struct {
	CreatedAt string `json:"createdAt"`
	Id        string `json:"id"`
	Password  string `json:"password"`
	Status    string `json:"status"`
	Username  string `json:"username"`
}

func (p *CloudProjectDatabaseUserResponse) String() string {
	return fmt.Sprintf(
		"Id: %s, User: %s, Status: %s",
		p.Id,
		p.Username,
		p.Status,
	)
}

func (v CloudProjectDatabaseUserResponse) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})

	obj["created_at"] = v.CreatedAt
	obj["id"] = v.Id
	obj["name"] = v.Username
	obj["status"] = v.Status

	return obj
}

type CloudProjectDatabaseUserCreateOpts struct {
	Name string `json:"name"`
}

func (opts *CloudProjectDatabaseUserCreateOpts) FromResource(d *schema.ResourceData) *CloudProjectDatabaseUserCreateOpts {
	opts.Name = d.Get("name").(string)
	return opts
}

func postCloudProjectDatabaseUser(d *schema.ResourceData, meta interface{}, engine string, endpoint string, params interface{}, res *CloudProjectDatabaseUserResponse, timeout string) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	clusterId := d.Get("cluster_id").(string)
	err := config.OVHClient.Post(endpoint, params, res)
	if err != nil {
		if errOvh, ok := err.(*ovh.APIError); engine == "mongodb" && ok && (errOvh.Code == 409) {
			return err
		}
		return fmt.Errorf("calling Post %s with params %+v:\n\t %q", endpoint, params, err)
	}

	log.Printf("[DEBUG] Waiting for user %s to be READY", res.Id)
	err = waitForCloudProjectDatabaseUserReady(config.OVHClient, serviceName, engine, clusterId, res.Id, d.Timeout(timeout))
	if err != nil {
		return fmt.Errorf("timeout while waiting user %s to be READY: %w", res.Id, err)
	}
	log.Printf("[DEBUG] user %s is READY", res.Id)

	d.Set("password", res.Password)
	return nil
}

func waitForCloudProjectDatabaseUserReady(client *ovh.Client, serviceName, engine string, databaseId string, userId string, timeOut time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "CREATING", "UPDATING"},
		Target:  []string{"READY"},
		Refresh: func() (interface{}, string, error) {
			res := &CloudProjectDatabaseUserResponse{}
			endpoint := fmt.Sprintf("/cloud/project/%s/database/%s/%s/user/%s",
				url.PathEscape(serviceName),
				url.PathEscape(engine),
				url.PathEscape(databaseId),
				url.PathEscape(userId),
			)
			err := client.Get(endpoint, res)
			if err != nil {
				return res, "", err
			}

			return res, res.Status, nil
		},
		Timeout:    timeOut,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForState()
	return err
}

func waitForCloudProjectDatabaseUserDeleted(client *ovh.Client, serviceName, engine string, databaseId string, userId string, timeOut time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"DELETING"},
		Target:  []string{"DELETED"},
		Refresh: func() (interface{}, string, error) {
			res := &CloudProjectDatabaseUserResponse{}
			endpoint := fmt.Sprintf("/cloud/project/%s/database/%s/%s/user/%s",
				url.PathEscape(serviceName),
				url.PathEscape(engine),
				url.PathEscape(databaseId),
				url.PathEscape(userId),
			)
			err := client.Get(endpoint, res)
			if err != nil {
				if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
					return res, "DELETED", nil
				}
				return res, "", err
			}

			return res, res.Status, nil
		},
		Timeout:    timeOut,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForState()
	return err
}

// Database

type CloudProjectDatabaseDatabaseResponse struct {
	Default bool   `json:"default"`
	Id      string `json:"id"`
	Name    string `json:"name"`
}

func (p *CloudProjectDatabaseDatabaseResponse) String() string {
	return fmt.Sprintf(
		"Id: %s, Database: %s, Default: %t",
		p.Id,
		p.Name,
		p.Default,
	)
}

func (v CloudProjectDatabaseDatabaseResponse) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})

	obj["default"] = v.Default
	obj["id"] = v.Id
	obj["name"] = v.Name

	return obj
}

type CloudProjectDatabaseDatabaseCreateOpts struct {
	Name string `json:"name"`
}

func (opts *CloudProjectDatabaseDatabaseCreateOpts) FromResource(d *schema.ResourceData) *CloudProjectDatabaseDatabaseCreateOpts {
	opts.Name = d.Get("name").(string)
	return opts
}

func waitForCloudProjectDatabaseDatabaseReady(client *ovh.Client, serviceName, engine string, serviceId string, databaseId string, timeOut time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"READY"},
		Refresh: func() (interface{}, string, error) {
			res := &CloudProjectDatabaseDatabaseResponse{}
			endpoint := fmt.Sprintf("/cloud/project/%s/database/%s/%s/database/%s",
				url.PathEscape(serviceName),
				url.PathEscape(engine),
				url.PathEscape(serviceId),
				url.PathEscape(databaseId),
			)
			err := client.Get(endpoint, res)
			if err != nil {
				if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
					return res, "PENDING", nil
				}
				return res, "", err
			}
			return res, "READY", nil
		},
		Timeout:    timeOut,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForState()
	return err
}

func waitForCloudProjectDatabaseDatabaseDeleted(client *ovh.Client, serviceName, engine string, serviceId string, databaseId string, timeOut time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"DELETING"},
		Target:  []string{"DELETED"},
		Refresh: func() (interface{}, string, error) {
			res := &CloudProjectDatabaseDatabaseResponse{}
			endpoint := fmt.Sprintf("/cloud/project/%s/database/%s/%s/database/%s",
				url.PathEscape(serviceName),
				url.PathEscape(engine),
				url.PathEscape(serviceId),
				url.PathEscape(databaseId),
			)
			err := client.Get(endpoint, res)
			if err != nil {
				if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
					return res, "DELETED", nil
				}
				return res, "", err
			}

			return res, "DELETING", nil
		},
		Timeout:    timeOut,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForState()
	return err
}

// Integration

type CloudProjectDatabaseIntegrationResponse struct {
	DestinationServiceId string            `json:"destinationServiceId"`
	Id                   string            `json:"id"`
	Parameters           map[string]string `json:"parameters"`
	SourceServiceId      string            `json:"sourceServiceId"`
	Status               string            `json:"status"`
	Type                 string            `json:"type"`
}

func (p *CloudProjectDatabaseIntegrationResponse) String() string {
	return fmt.Sprintf(
		"Id: %s, Type: %s ,SourceServiceId: %s, DestinationServiceId: %s",
		p.Id,
		p.Type,
		p.SourceServiceId,
		p.DestinationServiceId,
	)
}

func (v CloudProjectDatabaseIntegrationResponse) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})

	obj["destination_service_id"] = v.DestinationServiceId
	obj["id"] = v.Id
	obj["parameters"] = v.Parameters
	obj["source_service_id"] = v.SourceServiceId
	obj["status"] = v.Status
	obj["type"] = v.Type

	return obj
}

type CloudProjectDatabaseIntegrationCreateOpts struct {
	DestinationServiceId string            `json:"destinationServiceId"`
	Parameters           map[string]string `json:"parameters,omitempty"`
	SourceServiceId      string            `json:"sourceServiceId"`
	Type                 string            `json:"type,omitempty"`
}

func (opts *CloudProjectDatabaseIntegrationCreateOpts) FromResource(d *schema.ResourceData) *CloudProjectDatabaseIntegrationCreateOpts {
	opts.DestinationServiceId = d.Get("destination_service_id").(string)
	opts.Parameters = make(map[string]string)
	parameters := d.Get("parameters").(map[string]interface{})
	for k, v := range parameters {
		opts.Parameters[k] = v.(string)
	}
	opts.SourceServiceId = d.Get("source_service_id").(string)
	opts.Type = d.Get("type").(string)
	return opts
}

func waitForCloudProjectDatabaseIntegrationReady(client *ovh.Client, serviceName, engine string, serviceId string, integrationId string, timeOut time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"READY"},
		Refresh: func() (interface{}, string, error) {
			res := &CloudProjectDatabaseDatabaseResponse{}
			endpoint := fmt.Sprintf("/cloud/project/%s/database/%s/%s/integration/%s",
				url.PathEscape(serviceName),
				url.PathEscape(engine),
				url.PathEscape(serviceId),
				url.PathEscape(integrationId),
			)
			err := client.Get(endpoint, res)
			if err != nil {
				if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
					return res, "PENDING", nil
				}
				return res, "", err
			}
			return res, "READY", nil
		},
		Timeout:    timeOut,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForState()
	return err
}

func waitForCloudProjectDatabaseIntegrationDeleted(client *ovh.Client, serviceName, engine string, serviceId string, integrationId string, timeOut time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"DELETING"},
		Target:  []string{"DELETED"},
		Refresh: func() (interface{}, string, error) {
			res := &CloudProjectDatabaseDatabaseResponse{}
			endpoint := fmt.Sprintf("/cloud/project/%s/database/%s/%s/integration/%s",
				url.PathEscape(serviceName),
				url.PathEscape(engine),
				url.PathEscape(serviceId),
				url.PathEscape(integrationId),
			)
			err := client.Get(endpoint, res)
			if err != nil {
				if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
					return res, "DELETED", nil
				}
				return res, "", err
			}

			return res, "DELETING", nil
		},
		Timeout:    timeOut,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForState()
	return err
}

// Certificates

type CloudProjectDatabaseCertificatesResponse struct {
	Ca string `json:"ca"`
}

func (p *CloudProjectDatabaseCertificatesResponse) String() string {
	return fmt.Sprintf(
		"Ca: %s",
		p.Ca,
	)
}

func (v CloudProjectDatabaseCertificatesResponse) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})

	obj["ca"] = v.Ca

	return obj
}

// PostgresqlUser

type CloudProjectDatabasePostgresqlUserResponse struct {
	CreatedAt string   `json:"createdAt"`
	Id        string   `json:"id"`
	Password  string   `json:"password"`
	Roles     []string `json:"roles"`
	Status    string   `json:"status"`
	Username  string   `json:"username"`
}

func (p *CloudProjectDatabasePostgresqlUserResponse) String() string {
	return fmt.Sprintf(
		"Id: %s, User: %s, Status: %s",
		p.Id,
		p.Username,
		p.Status,
	)
}

func (v CloudProjectDatabasePostgresqlUserResponse) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})

	obj["created_at"] = v.CreatedAt
	obj["id"] = v.Id
	obj["name"] = v.Username
	obj["roles"] = v.Roles
	obj["status"] = v.Status

	return obj
}

type CloudProjectDatabasePostgresqlUserCreateOpts struct {
	Name  string   `json:"name"`
	Roles []string `json:"roles,omitempty"`
}

func (opts *CloudProjectDatabasePostgresqlUserCreateOpts) FromResource(d *schema.ResourceData) *CloudProjectDatabasePostgresqlUserCreateOpts {
	opts.Name = d.Get("name").(string)
	roles := d.Get("roles").(*schema.Set).List()
	opts.Roles = make([]string, len(roles))
	for i, e := range roles {
		if e != nil {
			opts.Roles[i] = e.(string)
		}
	}
	return opts
}

type CloudProjectDatabasePostgresqlUserUpdateOpts struct {
	Roles []string `json:"roles,omitempty"`
}

func (opts *CloudProjectDatabasePostgresqlUserUpdateOpts) FromResource(d *schema.ResourceData) *CloudProjectDatabasePostgresqlUserUpdateOpts {
	roles := d.Get("roles").(*schema.Set).List()
	opts.Roles = make([]string, len(roles))
	for i, e := range roles {
		opts.Roles[i] = e.(string)
	}
	return opts
}

// MongoDBUser

type CloudProjectDatabaseMongodbUserResponse struct {
	CreatedAt string   `json:"createdAt"`
	Id        string   `json:"id"`
	Password  string   `json:"password"`
	Roles     []string `json:"roles"`
	Status    string   `json:"status"`
	Username  string   `json:"username"`
}

func (p *CloudProjectDatabaseMongodbUserResponse) String() string {
	return fmt.Sprintf(
		"Id: %s, User: %s, Status: %s",
		p.Id,
		p.Username,
		p.Status,
	)
}

func (v CloudProjectDatabaseMongodbUserResponse) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})

	obj["created_at"] = v.CreatedAt
	obj["id"] = v.Id
	obj["name"] = v.Username
	obj["status"] = v.Status
	for i := range v.Roles {
		v.Roles[i] = strings.TrimSuffix(v.Roles[i], "@admin")
	}
	obj["roles"] = v.Roles

	return obj
}

type CloudProjectDatabaseMongodbUserCreateOpts struct {
	Name  string   `json:"name"`
	Roles []string `json:"roles"`
}

func (p *CloudProjectDatabaseMongodbUserCreateOpts) String() string {
	return fmt.Sprintf(
		"Name: %s, Password: <sensitive>, Roles: %v",
		p.Name,
		p.Roles,
	)
}

func (opts *CloudProjectDatabaseMongodbUserCreateOpts) FromResource(d *schema.ResourceData) *CloudProjectDatabaseMongodbUserCreateOpts {
	opts.Name = d.Get("name").(string)
	roles := d.Get("roles").(*schema.Set).List()
	opts.Roles = make([]string, len(roles))
	for i, e := range roles {
		if e != nil {
			opts.Roles[i] = e.(string)
		}
	}
	return opts
}

type CloudProjectDatabaseMongodbUserUpdateOpts struct {
	Roles []string `json:"roles"`
}

func (p *CloudProjectDatabaseMongodbUserUpdateOpts) String() string {
	return fmt.Sprintf(
		"Password: <sensitive>, Roles: %v",
		p.Roles,
	)
}

func (opts *CloudProjectDatabaseMongodbUserUpdateOpts) FromResource(d *schema.ResourceData) *CloudProjectDatabaseMongodbUserUpdateOpts {
	roles := d.Get("roles").(*schema.Set).List()
	opts.Roles = make([]string, len(roles))
	for i, e := range roles {
		opts.Roles[i] = e.(string)
	}
	return opts
}

// Redis User

type CloudProjectDatabaseRedisUserResponse struct {
	Categories []string `json:"categories"`
	Channels   []string `json:"channels"`
	Commands   []string `json:"commands"`
	CreatedAt  string   `json:"createdAt"`
	Id         string   `json:"id"`
	Keys       []string `json:"keys"`
	Password   string   `json:"password"`
	Status     string   `json:"status"`
	Username   string   `json:"username"`
}

func (p *CloudProjectDatabaseRedisUserResponse) String() string {
	return fmt.Sprintf(
		"Id: %s, User: %s, Status: %s",
		p.Id,
		p.Username,
		p.Status,
	)
}

func (v CloudProjectDatabaseRedisUserResponse) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})

	obj["categories"] = v.Categories
	obj["channels"] = v.Channels
	obj["commands"] = v.Commands
	obj["created_at"] = v.CreatedAt
	obj["id"] = v.Id
	obj["keys"] = v.Keys
	obj["name"] = v.Username
	obj["status"] = v.Status

	return obj
}

type CloudProjectDatabaseRedisUserCreateOpts struct {
	Categories []string `json:"categories,omitempty"`
	Channels   []string `json:"channels,omitempty"`
	Commands   []string `json:"commands,omitempty"`
	Keys       []string `json:"keys,omitempty"`
	Name       string   `json:"name"`
}

func getStringSlice(i interface{}) []string {
	iarr := i.(*schema.Set).List()
	arr := make([]string, len(iarr))
	for i, e := range iarr {
		if e != nil {
			arr[i] = e.(string)
		}
	}
	return arr
}

func (opts *CloudProjectDatabaseRedisUserCreateOpts) FromResource(d *schema.ResourceData) *CloudProjectDatabaseRedisUserCreateOpts {
	opts.Name = d.Get("name").(string)
	opts.Categories = getStringSlice(d.Get("categories"))
	opts.Channels = getStringSlice(d.Get("channels"))
	opts.Commands = getStringSlice(d.Get("commands"))
	opts.Keys = getStringSlice(d.Get("keys"))

	return opts
}

type CloudProjectDatabaseRedisUserUpdateOpts struct {
	Categories []string `json:"categories,omitempty"`
	Channels   []string `json:"channels,omitempty"`
	Commands   []string `json:"commands,omitempty"`
	Keys       []string `json:"keys,omitempty"`
}

func (opts *CloudProjectDatabaseRedisUserUpdateOpts) FromResource(d *schema.ResourceData) *CloudProjectDatabaseRedisUserUpdateOpts {
	opts.Categories = getStringSlice(d.Get("categories"))
	opts.Channels = getStringSlice(d.Get("channels"))
	opts.Commands = getStringSlice(d.Get("commands"))
	opts.Keys = getStringSlice(d.Get("keys"))
	return opts
}

// M3DB

// // User

type CloudProjectDatabaseM3dbUserResponse struct {
	CreatedAt string `json:"createdAt"`
	Group     string `json:"group"`
	Id        string `json:"id"`
	Password  string `json:"password"`
	Status    string `json:"status"`
	Username  string `json:"username"`
}

func (p *CloudProjectDatabaseM3dbUserResponse) String() string {
	return fmt.Sprintf(
		"Id: %s, User: %s, Status: %s",
		p.Id,
		p.Username,
		p.Status,
	)
}

func (v CloudProjectDatabaseM3dbUserResponse) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})

	obj["created_at"] = v.CreatedAt
	obj["group"] = v.Group
	obj["id"] = v.Id
	obj["name"] = v.Username
	obj["status"] = v.Status

	return obj
}

type CloudProjectDatabaseM3dbUserCreateOpts struct {
	Group string `json:"group,omitempty"`
	Name  string `json:"name"`
}

func (opts *CloudProjectDatabaseM3dbUserCreateOpts) FromResource(d *schema.ResourceData) *CloudProjectDatabaseM3dbUserCreateOpts {
	opts.Group = d.Get("group").(string)
	opts.Name = d.Get("name").(string)
	return opts
}

type CloudProjectDatabaseM3dbUserUpdateOpts struct {
	Group string `json:"group,omitempty"`
}

func (opts *CloudProjectDatabaseM3dbUserUpdateOpts) FromResource(d *schema.ResourceData) *CloudProjectDatabaseM3dbUserUpdateOpts {
	opts.Group = d.Get("group").(string)
	return opts
}

// // Namespace

func DiffDurationRfc3339(k, old, new string, d *schema.ResourceData) bool {
	newD, _ := rfc3339.ParseDuration(new)
	oldD, _ := rfc3339.ParseDuration(old)
	return newD == oldD
}

type CloudProjectDatabaseM3dbNamespaceRetention struct {
	BlockDataExpirationDuration string `json:"blockDataExpirationDuration,omitempty"`
	BlockSizeDuration           string `json:"blockSizeDuration,omitempty"`
	BufferFutureDuration        string `json:"bufferFutureDuration,omitempty"`
	BufferPastDuration          string `json:"bufferPastDuration,omitempty"`
	PeriodDuration              string `json:"periodDuration,omitempty"`
}

type CloudProjectDatabaseM3dbNamespaceResponse struct {
	Id                       string                                     `json:"id"`
	Name                     string                                     `json:"name"`
	Resolution               string                                     `json:"resolution"`
	Retention                CloudProjectDatabaseM3dbNamespaceRetention `json:"retention"`
	SnapshotEnabled          bool                                       `json:"snapshotEnabled"`
	Type                     string                                     `json:"type"`
	WritesToCommitLogEnabled bool                                       `json:"writesToCommitLogEnabled"`
}

func (p *CloudProjectDatabaseM3dbNamespaceResponse) String() string {
	return fmt.Sprintf(
		"Id: %s, Namespace: %s, Type: %s",
		p.Id,
		p.Name,
		p.Type,
	)
}

func (v CloudProjectDatabaseM3dbNamespaceResponse) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["id"] = v.Id
	obj["name"] = v.Name
	obj["resolution"] = v.Resolution
	obj["retention_block_data_expiration_duration"] = v.Retention.BlockDataExpirationDuration
	obj["retention_block_size_duration"] = v.Retention.BlockSizeDuration
	obj["retention_buffer_future_duration"] = v.Retention.BufferFutureDuration
	obj["retention_buffer_past_duration"] = v.Retention.BufferPastDuration
	obj["retention_period_duration"] = v.Retention.PeriodDuration
	obj["snapshot_enabled"] = v.SnapshotEnabled
	obj["type"] = v.Type
	obj["writes_to_commit_log_enabled"] = v.WritesToCommitLogEnabled
	return obj
}

type CloudProjectDatabaseM3dbNamespaceCreateOpts struct {
	Name                     string                                     `json:"name"`
	Resolution               string                                     `json:"resolution,omitempty"`
	Retention                CloudProjectDatabaseM3dbNamespaceRetention `json:"retention,omitempty"`
	SnapshotEnabled          bool                                       `json:"snapshotEnabled"`
	Type                     string                                     `json:"type"`
	WritesToCommitLogEnabled bool                                       `json:"writesToCommitLogEnabled"`
}

func (opts *CloudProjectDatabaseM3dbNamespaceCreateOpts) FromResource(d *schema.ResourceData) *CloudProjectDatabaseM3dbNamespaceCreateOpts {
	opts.Name = d.Get("name").(string)
	opts.Resolution = d.Get("resolution").(string)
	opts.Retention = CloudProjectDatabaseM3dbNamespaceRetention{
		BlockDataExpirationDuration: d.Get("retention_block_data_expiration_duration").(string),
		BlockSizeDuration:           d.Get("retention_block_size_duration").(string),
		BufferFutureDuration:        d.Get("retention_buffer_future_duration").(string),
		BufferPastDuration:          d.Get("retention_buffer_past_duration").(string),
		PeriodDuration:              d.Get("retention_period_duration").(string),
	}
	opts.SnapshotEnabled = d.Get("snapshot_enabled").(bool)
	opts.Type = "aggregated"
	opts.WritesToCommitLogEnabled = d.Get("writes_to_commit_log_enabled").(bool)
	return opts
}

type CloudProjectDatabaseM3dbNamespaceUpdateOpts struct {
	Resolution               string                                     `json:"resolution,omitempty"`
	Retention                CloudProjectDatabaseM3dbNamespaceRetention `json:"retention,omitempty"`
	SnapshotEnabled          bool                                       `json:"snapshotEnabled,omitempty"`
	WritesToCommitLogEnabled bool                                       `json:"writesToCommitLogEnabled,omitempty"`
}

func (opts *CloudProjectDatabaseM3dbNamespaceUpdateOpts) FromResource(d *schema.ResourceData) *CloudProjectDatabaseM3dbNamespaceUpdateOpts {
	opts.Resolution = d.Get("resolution").(string)
	opts.Retention = CloudProjectDatabaseM3dbNamespaceRetention{
		BlockDataExpirationDuration: d.Get("retention_block_data_expiration_duration").(string),
		BlockSizeDuration:           d.Get("retention_block_size_duration").(string),
		BufferFutureDuration:        d.Get("retention_buffer_future_duration").(string),
		BufferPastDuration:          d.Get("retention_buffer_past_duration").(string),
		PeriodDuration:              d.Get("retention_period_duration").(string),
	}
	opts.SnapshotEnabled = d.Get("snapshot_enabled").(bool)
	opts.WritesToCommitLogEnabled = d.Get("writes_to_commit_log_enabled").(bool)
	return opts
}

func waitForCloudProjectDatabaseM3dbNamespaceReady(client *ovh.Client, serviceName, databaseId string, namespaceId string, timeOut time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"READY"},
		Refresh: func() (interface{}, string, error) {
			res := &CloudProjectDatabaseM3dbNamespaceResponse{}
			endpoint := fmt.Sprintf("/cloud/project/%s/database/m3db/%s/namespace/%s",
				url.PathEscape(serviceName),
				url.PathEscape(databaseId),
				url.PathEscape(namespaceId),
			)
			err := client.Get(endpoint, res)
			if err != nil {
				if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
					return res, "PENDING", nil
				}
				return res, "", err
			}
			return res, "READY", nil
		},
		Timeout:    timeOut,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForState()
	return err
}

func waitForCloudProjectDatabaseM3dbNamespaceDeleted(client *ovh.Client, serviceName, databaseId string, namespaceId string, timeOut time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"DELETING"},
		Target:  []string{"DELETED"},
		Refresh: func() (interface{}, string, error) {
			res := &CloudProjectDatabaseM3dbNamespaceResponse{}
			endpoint := fmt.Sprintf("/cloud/project/%s/database/m3db/%s/namespace/%s",
				url.PathEscape(serviceName),
				url.PathEscape(databaseId),
				url.PathEscape(namespaceId),
			)
			err := client.Get(endpoint, res)
			if err != nil {
				if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
					return res, "DELETED", nil
				}
				return res, "", err
			}

			return res, "DELETING", nil
		},
		Timeout:    timeOut,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForState()
	return err
}

// Opensearch

// // User

type CloudProjectDatabaseOpensearchUserAcl struct {
	Pattern    string `json:"pattern"`
	Permission string `json:"permission"`
}

func (v CloudProjectDatabaseOpensearchUserAcl) ToMap() map[string]string {
	obj := make(map[string]string)

	obj["pattern"] = v.Pattern
	obj["permission"] = v.Permission

	return obj
}

type CloudProjectDatabaseOpensearchUserResponse struct {
	Acls      []CloudProjectDatabaseOpensearchUserAcl `json:"acls"`
	CreatedAt string                                  `json:"createdAt"`
	Id        string                                  `json:"id"`
	Password  string                                  `json:"password"`
	Status    string                                  `json:"status"`
	Username  string                                  `json:"username"`
}

func (p *CloudProjectDatabaseOpensearchUserResponse) String() string {
	return fmt.Sprintf(
		"Id: %s, User: %s, Status: %s",
		p.Id,
		p.Username,
		p.Status,
	)
}

func (v CloudProjectDatabaseOpensearchUserResponse) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})

	var acls []map[string]string
	for _, e := range v.Acls {
		acls = append(acls, e.ToMap())
	}

	obj["acls"] = acls
	obj["created_at"] = v.CreatedAt
	obj["id"] = v.Id
	obj["name"] = v.Username
	obj["status"] = v.Status

	return obj
}

type CloudProjectDatabaseOpensearchUserCreateOpts struct {
	Acls []CloudProjectDatabaseOpensearchUserAcl `json:"acls,omitempty"`
	Name string                                  `json:"name"`
}

func (opts *CloudProjectDatabaseOpensearchUserCreateOpts) FromResource(d *schema.ResourceData) *CloudProjectDatabaseOpensearchUserCreateOpts {
	opts.Name = d.Get("name").(string)
	acls := d.Get("acls").(*schema.Set).List()
	opts.Acls = make([]CloudProjectDatabaseOpensearchUserAcl, len(acls))
	for i, e := range acls {
		aclMap := e.(map[string]interface{})
		opts.Acls[i] = CloudProjectDatabaseOpensearchUserAcl{
			Pattern:    aclMap["pattern"].(string),
			Permission: aclMap["permission"].(string),
		}
	}
	return opts
}

type CloudProjectDatabaseOpensearchUserUpdateOpts struct {
	Acls []CloudProjectDatabaseOpensearchUserAcl `json:"acls,omitempty"`
}

func (opts *CloudProjectDatabaseOpensearchUserUpdateOpts) FromResource(d *schema.ResourceData) *CloudProjectDatabaseOpensearchUserUpdateOpts {
	acls := d.Get("acls").(*schema.Set).List()
	opts.Acls = make([]CloudProjectDatabaseOpensearchUserAcl, len(acls))
	for i, e := range acls {
		aclMap := e.(map[string]interface{})
		opts.Acls[i] = CloudProjectDatabaseOpensearchUserAcl{
			Pattern:    aclMap["pattern"].(string),
			Permission: aclMap["permission"].(string),
		}
	}
	return opts
}

// // Pattern

type CloudProjectDatabaseOpensearchPatternResponse struct {
	Id            string `json:"id"`
	MaxIndexCount int    `json:"maxIndexCount"`
	Pattern       string `json:"pattern"`
}

func (p *CloudProjectDatabaseOpensearchPatternResponse) String() string {
	return fmt.Sprintf(
		"Id: %s, Pattern: %s, MaxIndexCount: %d",
		p.Id,
		p.Pattern,
		p.MaxIndexCount,
	)
}

func (v CloudProjectDatabaseOpensearchPatternResponse) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})

	obj["id"] = v.Id
	obj["max_index_count"] = v.MaxIndexCount
	obj["pattern"] = v.Pattern

	return obj
}

type CloudProjectDatabaseOpensearchPatternCreateOpts struct {
	MaxIndexCount int    `json:"maxIndexCount,omitempty"`
	Pattern       string `json:"pattern"`
}

func (opts *CloudProjectDatabaseOpensearchPatternCreateOpts) FromResource(d *schema.ResourceData) *CloudProjectDatabaseOpensearchPatternCreateOpts {
	opts.MaxIndexCount = d.Get("max_index_count").(int)
	opts.Pattern = d.Get("pattern").(string)

	return opts
}

func waitForCloudProjectDatabaseOpensearchPatternReady(client *ovh.Client, serviceName, databaseId string, patternId string, timeOut time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"READY"},
		Refresh: func() (interface{}, string, error) {
			res := &CloudProjectDatabaseOpensearchPatternResponse{}
			endpoint := fmt.Sprintf("/cloud/project/%s/database/opensearch/%s/pattern/%s",
				url.PathEscape(serviceName),
				url.PathEscape(databaseId),
				url.PathEscape(patternId),
			)
			err := client.Get(endpoint, res)
			if err != nil {
				if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
					return res, "PENDING", nil
				}
				return res, "", err
			}
			return res, "READY", nil
		},
		Timeout:    timeOut,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForState()
	return err
}

func waitForCloudProjectDatabaseOpensearchPatternDeleted(client *ovh.Client, serviceName, databaseId string, patternId string, timeOut time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"DELETING"},
		Target:  []string{"DELETED"},
		Refresh: func() (interface{}, string, error) {
			res := &CloudProjectDatabaseOpensearchPatternResponse{}
			endpoint := fmt.Sprintf("/cloud/project/%s/database/opensearch/%s/pattern/%s",
				url.PathEscape(serviceName),
				url.PathEscape(databaseId),
				url.PathEscape(patternId),
			)
			err := client.Get(endpoint, res)
			if err != nil {
				if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
					return res, "DELETED", nil
				}
				return res, "", err
			}

			return res, "DELETING", nil
		},
		Timeout:    timeOut,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForState()
	return err
}

// Kafka

// // Topic

type CloudProjectDatabaseKafkaTopicResponse struct {
	Id                string `json:"id"`
	MinInsyncReplicas int    `json:"minInsyncReplicas"`
	Name              string `json:"name"`
	Partitions        int    `json:"partitions"`
	Replication       int    `json:"replication"`
	RetentionBytes    int    `json:"retentionBytes"`
	RetentionHours    int    `json:"retentionHours"`
}

func (p *CloudProjectDatabaseKafkaTopicResponse) String() string {
	return fmt.Sprintf(
		"Id: %s, Name: %s",
		p.Id,
		p.Name,
	)
}

func (v CloudProjectDatabaseKafkaTopicResponse) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})

	obj["id"] = v.Id
	obj["min_insync_replicas"] = v.MinInsyncReplicas
	obj["name"] = v.Name
	obj["partitions"] = v.Partitions
	obj["replication"] = v.Replication
	obj["retention_bytes"] = v.RetentionBytes
	obj["retention_hours"] = v.RetentionHours

	return obj
}

type CloudProjectDatabaseKafkaTopicCreateOpts struct {
	MinInsyncReplicas int    `json:"minInsyncReplicas"`
	Name              string `json:"name,omitempty"`
	Partitions        int    `json:"partitions"`
	Replication       int    `json:"replication"`
	RetentionBytes    int    `json:"retentionBytes"`
	RetentionHours    int    `json:"retentionHours"`
}

func (opts *CloudProjectDatabaseKafkaTopicCreateOpts) FromResource(d *schema.ResourceData) *CloudProjectDatabaseKafkaTopicCreateOpts {
	opts.MinInsyncReplicas = d.Get("min_insync_replicas").(int)
	opts.Name = d.Get("name").(string)
	opts.Partitions = d.Get("partitions").(int)
	opts.Replication = d.Get("replication").(int)
	opts.RetentionBytes = d.Get("retention_bytes").(int)
	opts.RetentionHours = d.Get("retention_hours").(int)

	return opts
}

func validateIsSupEqual(v, min int) (errors []error) {
	if v < min {
		errors = append(errors, fmt.Errorf("Value %d is inferior of min value %d", v, min))
	}
	return
}

func validateCloudProjectDatabaseKafkaTopicMinInsyncReplicasFunc(v interface{}, k string) (ws []string, errors []error) {
	errors = validateIsSupEqual(v.(int), 1)
	return
}

func validateCloudProjectDatabaseKafkaTopicPartitionsFunc(v interface{}, k string) (ws []string, errors []error) {
	errors = validateIsSupEqual(v.(int), 1)
	return
}

func validateCloudProjectDatabaseKafkaTopicReplicationFunc(v interface{}, k string) (ws []string, errors []error) {
	errors = validateIsSupEqual(v.(int), 2)
	return
}

func validateCloudProjectDatabaseKafkaTopicRetentionHoursFunc(v interface{}, k string) (ws []string, errors []error) {
	errors = validateIsSupEqual(v.(int), -1)
	return
}

func waitForCloudProjectDatabaseKafkaTopicReady(client *ovh.Client, serviceName, databaseId string, topicId string, timeOut time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"READY"},
		Refresh: func() (interface{}, string, error) {
			res := &CloudProjectDatabaseKafkaTopicResponse{}
			endpoint := fmt.Sprintf("/cloud/project/%s/database/kafka/%s/topic/%s",
				url.PathEscape(serviceName),
				url.PathEscape(databaseId),
				url.PathEscape(topicId),
			)
			err := client.Get(endpoint, res)
			if err != nil {
				if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
					return res, "PENDING", nil
				}
				return res, "", err
			}
			return res, "READY", nil
		},
		Timeout:    timeOut,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForState()
	return err
}

func waitForCloudProjectDatabaseKafkaTopicDeleted(client *ovh.Client, serviceName, databaseId string, topicId string, timeOut time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"DELETING"},
		Target:  []string{"DELETED"},
		Refresh: func() (interface{}, string, error) {
			res := &CloudProjectDatabaseKafkaTopicResponse{}
			endpoint := fmt.Sprintf("/cloud/project/%s/database/kafka/%s/topic/%s",
				url.PathEscape(serviceName),
				url.PathEscape(databaseId),
				url.PathEscape(topicId),
			)
			err := client.Get(endpoint, res)
			if err != nil {
				if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
					return res, "DELETED", nil
				}
				return res, "", err
			}

			return res, "DELETING", nil
		},
		Timeout:    timeOut,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForState()
	return err
}

// // ACL

type CloudProjectDatabaseKafkaAclResponse struct {
	Id         string `json:"id"`
	Permission string `json:"permission"`
	Topic      string `json:"topic"`
	Username   string `json:"username"`
}

func (p *CloudProjectDatabaseKafkaAclResponse) String() string {
	return fmt.Sprintf(
		"Id: %s, Permission: %s, affected User: %s, affected Topic %s",
		p.Id,
		p.Permission,
		p.Username,
		p.Topic,
	)
}

func (v CloudProjectDatabaseKafkaAclResponse) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})

	obj["id"] = v.Id
	obj["permission"] = v.Permission
	obj["topic"] = v.Topic
	obj["username"] = v.Username

	return obj
}

type CloudProjectDatabaseKafkaAclCreateOpts struct {
	Permission string `json:"permission"`
	Topic      string `json:"topic"`
	Username   string `json:"username"`
}

func (opts *CloudProjectDatabaseKafkaAclCreateOpts) FromResource(d *schema.ResourceData) *CloudProjectDatabaseKafkaAclCreateOpts {
	opts.Permission = d.Get("permission").(string)
	opts.Topic = d.Get("topic").(string)
	opts.Username = d.Get("username").(string)

	return opts
}

func waitForCloudProjectDatabaseKafkaAclReady(client *ovh.Client, serviceName, databaseId string, aclId string, timeOut time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"READY"},
		Refresh: func() (interface{}, string, error) {
			res := &CloudProjectDatabaseKafkaTopicResponse{}
			endpoint := fmt.Sprintf("/cloud/project/%s/database/kafka/%s/acl/%s",
				url.PathEscape(serviceName),
				url.PathEscape(databaseId),
				url.PathEscape(aclId),
			)
			err := client.Get(endpoint, res)
			if err != nil {
				if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
					return res, "PENDING", nil
				}
				return res, "", err
			}
			return res, "READY", nil
		},
		Timeout:    timeOut,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForState()
	return err
}

func waitForCloudProjectDatabaseKafkaAclDeleted(client *ovh.Client, serviceName, databaseId string, aclId string, timeOut time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"DELETING"},
		Target:  []string{"DELETED"},
		Refresh: func() (interface{}, string, error) {
			res := &CloudProjectDatabaseKafkaTopicResponse{}
			endpoint := fmt.Sprintf("/cloud/project/%s/database/kafka/%s/acl/%s",
				url.PathEscape(serviceName),
				url.PathEscape(databaseId),
				url.PathEscape(aclId),
			)
			err := client.Get(endpoint, res)
			if err != nil {
				if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
					return res, "DELETED", nil
				}
				return res, "", err
			}

			return res, "DELETING", nil
		},
		Timeout:    timeOut,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForState()
	return err
}

// // User Access

type CloudProjectDatabaseKafkaUserAccessResponse struct {
	Cert string `json:"cert"`
	Key  string `json:"key"`
}

func (p *CloudProjectDatabaseKafkaUserAccessResponse) String() string {
	return fmt.Sprintf(
		"Cert: %s",
		p.Cert,
	)
}

func (v CloudProjectDatabaseKafkaUserAccessResponse) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})

	obj["cert"] = v.Cert
	obj["key"] = v.Key

	return obj
}
