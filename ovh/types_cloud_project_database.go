package ovh

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/go-ovh/ovh"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/ovhwrap"
	"github.com/ybriffa/rfc3339"
	"golang.org/x/exp/slices"
)

// Helper
func diagnosticsToError(diags diag.Diagnostics) error {
	if diags.HasError() {
		return fmt.Errorf(diags[slices.IndexFunc(diags, func(d diag.Diagnostic) bool { return d.Severity == diag.Error })].Summary)
	}
	return nil
}

// enginesWithoutBackupTime is the list of engines
// for which "backup_time" is not customizable
var enginesWithoutBackupTime = []string{"m3db", "grafana", "kafka", "kafkaconnect", "kafkamirrormaker", "opensearch", "m3aggregator"}

type CloudProjectDatabaseBackups struct {
	Regions []string `json:"regions,omitempty"`
	Time    string   `json:"time,omitempty"`
}

type CloudProjectDatabaseIPRestriction struct {
	Description string `json:"description"`
	IP          string `json:"ip"`
}

func (opts *CloudProjectDatabaseIPRestriction) FromResourceWithPath(d *schema.ResourceData, path string) *CloudProjectDatabaseIPRestriction {
	opts.Description = d.Get(fmt.Sprintf("%s.description", path)).(string)
	opts.IP = d.Get(fmt.Sprintf("%s.ip", path)).(string)

	return opts
}

type CloudProjectDatabaseIPRestrictionResponse struct {
	CloudProjectDatabaseIPRestriction
	Status string `json:"status"`
}

func (ir CloudProjectDatabaseIPRestrictionResponse) toMap() map[string]interface{} {
	obj := make(map[string]interface{})

	obj["description"] = ir.Description
	obj["ip"] = ir.IP
	obj["status"] = ir.Status

	return obj
}

type CloudProjectDatabaseResponse struct {
	AclsEnabled           bool                                        `json:"aclsEnabled"`
	Backups               CloudProjectDatabaseBackups                 `json:"backups"`
	CreatedAt             string                                      `json:"createdAt"`
	Description           string                                      `json:"description"`
	Endpoints             []CloudProjectDatabaseEndpoint              `json:"endpoints"`
	Flavor                string                                      `json:"flavor"`
	ID                    string                                      `json:"id"`
	IPRestrictions        []CloudProjectDatabaseIPRestrictionResponse `json:"ipRestrictions"`
	MaintenanceTime       string                                      `json:"maintenanceTime"`
	NetworkID             string                                      `json:"networkId"`
	NetworkType           string                                      `json:"networkType"`
	Plan                  string                                      `json:"plan"`
	NodeNumber            int                                         `json:"nodeNumber"`
	Region                string                                      `json:"region"`
	RestAPI               bool                                        `json:"restApi"`
	SchemaRegistry        bool                                        `json:"schemaRegistry"`
	Status                string                                      `json:"status"`
	SubnetID              string                                      `json:"subnetId"`
	Version               string                                      `json:"version"`
	Disk                  CloudProjectDatabaseDisk                    `json:"disk"`
	AdvancedConfiguration map[string]string                           `json:"advancedConfiguration"`
	EnablePrometheus      bool                                        `json:"enablePrometheus"`
}

func (r CloudProjectDatabaseResponse) toMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["backup_regions"] = r.Backups.Regions
	obj["backup_time"] = r.Backups.Time
	obj["created_at"] = r.CreatedAt
	obj["description"] = r.Description
	obj["id"] = r.ID

	var ipRests []map[string]interface{}
	for _, ir := range r.IPRestrictions {
		ipRests = append(ipRests, ir.toMap())
	}
	obj["ip_restrictions"] = ipRests

	var endpoints []map[string]interface{}
	for _, e := range r.Endpoints {
		endpoints = append(endpoints, e.ToMap())
	}
	obj["endpoints"] = endpoints

	obj["flavor"] = r.Flavor
	obj["kafka_rest_api"] = r.RestAPI
	obj["kafka_schema_registry"] = r.SchemaRegistry
	obj["maintenance_time"] = r.MaintenanceTime
	obj["network_type"] = r.NetworkType

	var nodes []map[string]interface{}
	for i := 0; i < r.NodeNumber; i++ {
		node := CloudProjectDatabaseNodes{
			Region:    r.Region,
			NetworkId: r.NetworkID,
			SubnetId:  r.SubnetID,
		}
		nodes = append(nodes, node.ToMap())
	}
	obj["nodes"] = nodes

	obj["opensearch_acls_enabled"] = r.AclsEnabled
	obj["plan"] = r.Plan
	obj["status"] = r.Status
	obj["version"] = r.Version
	obj["disk_size"] = r.Disk.Size
	obj["disk_type"] = r.Disk.Type
	obj["advanced_configuration"] = r.AdvancedConfiguration

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
	Backups         *CloudProjectDatabaseBackups        `json:"backups,omitempty"`
	Description     string                              `json:"description,omitempty"`
	Disk            CloudProjectDatabaseDisk            `json:"disk,omitempty"`
	IPRestrictions  []CloudProjectDatabaseIPRestriction `json:"ipRestrictions,omitempty"`
	MaintenanceTime string                              `json:"maintenanceTime,omitempty"`
	NetworkId       string                              `json:"networkId,omitempty"`
	NodesPattern    CloudProjectDatabaseNodesPattern    `json:"nodesPattern,omitempty"`
	Plan            string                              `json:"plan"`
	SubnetId        string                              `json:"subnetId,omitempty"`
	Version         string                              `json:"version"`
}

type CloudProjectDatabaseDisk struct {
	Type string `json:"type,omitempty"`
	Size int    `json:"size,omitempty"`
}

func validateCloudProjectDatabaseDiskSize(v any, p cty.Path) (diags diag.Diagnostics) {
	diags = validateIsSupEqual(v.(int), 0)
	return
}

type CloudProjectDatabaseNodesPattern struct {
	Flavor string `json:"flavor"`
	Number int    `json:"number"`
	Region string `json:"region"`
}

func (opts *CloudProjectDatabaseCreateOpts) fromResource(d *schema.ResourceData) (*CloudProjectDatabaseCreateOpts, error) {
	opts.Description = d.Get("description").(string)
	opts.Plan = d.Get("plan").(string)

	nodes := []CloudProjectDatabaseNodes{}
	nbOfNodes := d.Get("nodes.#").(int)
	for i := 0; i < nbOfNodes; i++ {
		nodes = append(nodes, *(&CloudProjectDatabaseNodes{}).FromResourceWithPath(d, fmt.Sprintf("nodes.%d", i)))
	}

	ipRests := d.Get("ip_restrictions").(*schema.Set).List()
	opts.IPRestrictions = make([]CloudProjectDatabaseIPRestriction, len(ipRests))
	for i, ir := range ipRests {
		ipRestMap := ir.(map[string]interface{})
		opts.IPRestrictions[i] = CloudProjectDatabaseIPRestriction{
			Description: ipRestMap["description"].(string),
			IP:          ipRestMap["ip"].(string),
		}
	}

	if err := checkNodesEquality(nodes); err != nil {
		return nil, err
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

	regions, err := helpers.StringsFromSchema(d, "backup_regions")
	if err != nil {
		return nil, err
	}
	time := d.Get("backup_time").(string)
	engine := d.Get("engine").(string)
	if time != "" && slices.Contains(enginesWithoutBackupTime, engine) {
		return nil, fmt.Errorf("backup_time is not customizable for engine %q", engine)
	}

	if len(regions) != 0 || time != "" {
		opts.Backups = &CloudProjectDatabaseBackups{
			Regions: regions,
			Time:    time,
		}
	}

	opts.MaintenanceTime = d.Get("maintenance_time").(string)

	return opts, nil
}

type CloudProjectDatabaseUpdateOpts struct {
	AclsEnabled     bool                                `json:"aclsEnabled,omitempty"`
	Backups         *CloudProjectDatabaseBackups        `json:"backups,omitempty"`
	Description     string                              `json:"description,omitempty"`
	Disk            CloudProjectDatabaseDisk            `json:"disk,omitempty"`
	Flavor          string                              `json:"flavor,omitempty"`
	IPRestrictions  []CloudProjectDatabaseIPRestriction `json:"ipRestrictions,omitempty"`
	MaintenanceTime string                              `json:"maintenanceTime,omitempty"`
	Plan            string                              `json:"plan,omitempty"`
	RestAPI         bool                                `json:"restApi,omitempty"`
	SchemaRegistry  bool                                `json:"schemaRegistry,omitempty"`
	Version         string                              `json:"version,omitempty"`
}

func (opts *CloudProjectDatabaseUpdateOpts) fromResource(d *schema.ResourceData) (*CloudProjectDatabaseUpdateOpts, error) {
	engine := d.Get("engine").(string)
	if engine == "opensearch" {
		opts.AclsEnabled = d.Get("opensearch_acls_enabled").(bool)
	}
	if engine == "kafka" {
		opts.RestAPI = d.Get("kafka_rest_api").(bool)
		opts.SchemaRegistry = d.Get("kafka_schema_registry").(bool)
	}

	opts.Description = d.Get("description").(string)
	opts.Plan = d.Get("plan").(string)
	opts.Flavor = d.Get("flavor").(string)
	opts.Version = d.Get("version").(string)
	opts.Disk = CloudProjectDatabaseDisk{Size: d.Get("disk_size").(int)}

	ipRests := d.Get("ip_restrictions").(*schema.Set).List()
	opts.IPRestrictions = make([]CloudProjectDatabaseIPRestriction, len(ipRests))
	for i, ir := range ipRests {
		ipRestMap := ir.(map[string]interface{})
		opts.IPRestrictions[i] = CloudProjectDatabaseIPRestriction{
			Description: ipRestMap["description"].(string),
			IP:          ipRestMap["ip"].(string),
		}

	}

	regions, err := helpers.StringsFromSchema(d, "backup_regions")
	if err != nil {
		return nil, err
	}

	if d.HasChange("backup_time") && slices.Contains(enginesWithoutBackupTime, engine) {
		return nil, fmt.Errorf("backup_time is not customizable for engine %q", engine)
	}

	time := d.Get("backup_time").(string)
	if engine != "kafka" && (len(regions) != 0 || time != "") {
		opts.Backups = &CloudProjectDatabaseBackups{
			Regions: regions,
			Time:    time,
		}
	}

	opts.MaintenanceTime = d.Get("maintenance_time").(string)

	return opts, nil
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

func waitForCloudProjectDatabaseReady(ctx context.Context, client *ovhwrap.Client, serviceName, engine string, databaseId string, timeOut time.Duration) error {
	stateConf := &retry.StateChangeConf{
		Pending: []string{"PENDING", "CREATING", "UPDATING"},
		Target:  []string{"READY"},
		Refresh: func() (interface{}, string, error) {
			res := &CloudProjectDatabaseResponse{}
			endpoint := fmt.Sprintf("/cloud/project/%s/database/%s/%s",
				url.PathEscape(serviceName),
				url.PathEscape(engine),
				url.PathEscape(databaseId),
			)
			err := client.GetWithContext(ctx, endpoint, res)
			if err != nil {
				return res, "", err
			}

			return res, res.Status, nil
		},
		Timeout:    timeOut,
		Delay:      30 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func waitForCloudProjectDatabaseDeleted(ctx context.Context, client *ovhwrap.Client, serviceName, engine string, databaseId string, timeOut time.Duration) error {
	stateConf := &retry.StateChangeConf{
		Pending: []string{"DELETING"},
		Target:  []string{"DELETED"},
		Refresh: func() (interface{}, string, error) {
			res := &CloudProjectDatabaseResponse{}
			endpoint := fmt.Sprintf("/cloud/project/%s/database/%s/%s",
				url.PathEscape(serviceName),
				url.PathEscape(engine),
				url.PathEscape(databaseId),
			)
			err := client.GetWithContext(ctx, endpoint, res)
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

	_, err := stateConf.WaitForStateContext(ctx)
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

// Log Subscrition

type CloudProjectDatabaseLogSubscriptionResource struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type CloudProjectDatabaseLogSubscriptionResponse struct {
	CreatedAt      string                                      `json:"createdAt"`
	Kind           string                                      `json:"kind"`
	OperationID    string                                      `json:"operationId"`
	Resource       CloudProjectDatabaseLogSubscriptionResource `json:"resource"`
	LDPServiceName string                                      `json:"serviceName"`
	StreamID       string                                      `json:"streamId"`
	SubscriptionID string                                      `json:"subscriptionId"`
	UpdatedAt      string                                      `json:"updatedAt"`
}

func (r *CloudProjectDatabaseLogSubscriptionResponse) string() string {
	return fmt.Sprintf(
		"Operation ID: %s, Subscription ID: %s, Stream ID: %s",
		r.OperationID,
		r.SubscriptionID,
		r.StreamID,
	)
}

func (r CloudProjectDatabaseLogSubscriptionResponse) toMap() map[string]interface{} {
	obj := make(map[string]interface{})

	obj["created_at"] = r.CreatedAt
	obj["id"] = r.SubscriptionID
	obj["ldp_service_name"] = r.LDPServiceName
	obj["kind"] = r.Kind
	obj["operation_id"] = r.OperationID
	obj["resource_name"] = r.Resource.Name
	obj["resource_type"] = r.Resource.Type
	obj["stream_id"] = r.StreamID
	obj["updated_at"] = r.UpdatedAt

	return obj
}

type CloudProjectDatabaseLogSubscriptionCreateOpts struct {
	StreamID string `json:"streamId"`
	Kind     string `json:"kind"`
}

func (opts *CloudProjectDatabaseLogSubscriptionCreateOpts) fromResource(d *schema.ResourceData) *CloudProjectDatabaseLogSubscriptionCreateOpts {
	opts.StreamID = d.Get("stream_id").(string)
	opts.Kind = d.Get("kind").(string)
	return opts
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

func waitForCloudProjectDatabaseIpRestrictionReady(ctx context.Context, client *ovhwrap.Client, serviceName, engine string, databaseId string, ip string, timeOut time.Duration) error {
	stateConf := &retry.StateChangeConf{
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
			err := client.GetWithContext(ctx, endpoint, res)
			if err != nil {
				return res, "", err
			}

			return res, res.Status, nil
		},
		Timeout:    timeOut,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func waitForCloudProjectDatabaseIpRestrictionDeleted(ctx context.Context, client *ovhwrap.Client, serviceName, engine string, databaseId string, ip string, timeOut time.Duration) error {
	stateConf := &retry.StateChangeConf{
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
			err := client.GetWithContext(ctx, endpoint, res)
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

	_, err := stateConf.WaitForStateContext(ctx)
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

func importCloudProjectDatabaseUser(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	givenId := d.Id()
	n := 3
	splitId := strings.SplitN(givenId, "/", n)
	if len(splitId) != n {
		return nil, fmt.Errorf("import Id is not service_name/cluster_id/id formatted")
	}
	serviceName := splitId[0]
	clusterId := splitId[1]
	id := splitId[2]
	d.SetId(id)
	d.Set("cluster_id", clusterId)
	d.Set("service_name", serviceName)

	results := make([]*schema.ResourceData, 1)
	results[0] = d
	return results, nil
}

func postCloudProjectDatabaseUser(ctx context.Context, d *schema.ResourceData, meta interface{}, engine string, dsReadFunc, readFunc schema.ReadContextFunc, updateFunc schema.UpdateContextFunc, f func() interface{}) diag.Diagnostics {
	name := d.Get("name").(string)
	if name == "avnadmin" && engine != "redis" && engine != "valkey" {
		diags := dsReadFunc(ctx, d, meta)
		if diags.HasError() {
			return diags
		}
		return updateFunc(ctx, d, meta)
	}
	if engine == "grafana" && name != "avnadmin" {
		return diag.FromErr(fmt.Errorf("the Grafana engine does not allow to create a user resource other than avnadmin"))
	}

	serviceName := d.Get("service_name").(string)
	clusterId := d.Get("cluster_id").(string)

	endpoint := fmt.Sprintf("/cloud/project/%s/database/%s/%s/user",
		url.PathEscape(serviceName),
		url.PathEscape(engine),
		url.PathEscape(clusterId),
	)

	params := f()
	res := &CloudProjectDatabaseUserResponse{}

	log.Printf("[DEBUG] Will create user: %+v for cluster %s from project %s", params, clusterId, serviceName)
	err := postFuncCloudProjectDatabaseUser(ctx, d, meta, engine, endpoint, params, res, schema.TimeoutCreate)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(res.Id)
	return readFunc(ctx, d, meta)
}

func postFuncCloudProjectDatabaseUser(ctx context.Context, d *schema.ResourceData, meta interface{}, engine string, endpoint string, params interface{}, res *CloudProjectDatabaseUserResponse, timeout string) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	clusterId := d.Get("cluster_id").(string)
	err := config.OVHClient.PostWithContext(ctx, endpoint, params, res)
	if err != nil {
		if errOvh, ok := err.(*ovh.APIError); engine == "mongodb" && ok && (errOvh.Code == 409) {
			return err
		}
		return fmt.Errorf("calling Post %s with params %+v:\n\t %q", endpoint, params, err)
	}

	log.Printf("[DEBUG] Waiting for user %s to be READY", res.Id)
	err = waitForCloudProjectDatabaseUserReady(ctx, config.OVHClient, serviceName, engine, clusterId, res.Id, d.Timeout(timeout))
	if err != nil {
		return fmt.Errorf("timeout while waiting user %s to be READY: %w", res.Id, err)
	}
	log.Printf("[DEBUG] user %s is READY", res.Id)

	d.Set("password", res.Password)
	return nil
}

func updateCloudProjectDatabaseUser(ctx context.Context, d *schema.ResourceData, meta interface{}, engine string, readFunc schema.ReadContextFunc, f func() interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	clusterId := d.Get("cluster_id").(string)
	isAvnAdmin := d.Get("name").(string) == "avnadmin"
	// The M3DB condition must be remove when avnadmin password reset will be possible on this engine
	passwordReset := d.HasChange("password_reset") || (d.IsNewResource() && isAvnAdmin && engine != "m3db")
	id := d.Id()

	endpoint := fmt.Sprintf("/cloud/project/%s/database/%s/%s/user/%s",
		url.PathEscape(serviceName),
		url.PathEscape(engine),
		url.PathEscape(clusterId),
		url.PathEscape(id),
	)

	if !(isAvnAdmin && engine == "postgresql") {
		params := f()

		log.Printf("[DEBUG] Will update user: %+v from cluster %s from project %s", params, clusterId, serviceName)
		err := config.OVHClient.Put(endpoint, params, nil)
		if err != nil {
			return diag.Errorf("calling Put %s with params %+v:\n\t %q", endpoint, params, err)
		}

		log.Printf("[DEBUG] Waiting for user %s to be READY", id)
		err = waitForCloudProjectDatabaseUserReady(ctx, config.OVHClient, serviceName, engine, clusterId, id, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.Errorf("timeout while waiting user %s to be READY: %s", id, err.Error())
		}
		log.Printf("[DEBUG] user %s is READY", id)
	}

	if passwordReset {
		pwdResetEndpoint := endpoint + "/credentials/reset"
		res := &CloudProjectDatabaseUserResponse{}
		log.Printf("[DEBUG] Will update user password for cluster %s from project %s", clusterId, serviceName)
		err := postFuncCloudProjectDatabaseUser(ctx, d, meta, engine, pwdResetEndpoint, nil, res, schema.TimeoutUpdate)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return readFunc(ctx, d, meta)
}

func deleteCloudProjectDatabaseUser(ctx context.Context, d *schema.ResourceData, meta interface{}, engine string) diag.Diagnostics {
	name := d.Get("name").(string)
	if name == "avnadmin" && engine != "redis" && engine != "valkey" {
		d.SetId("")
		return nil
	}

	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	clusterId := d.Get("cluster_id").(string)
	id := d.Id()

	endpoint := fmt.Sprintf("/cloud/project/%s/database/%s/%s/user/%s",
		url.PathEscape(serviceName),
		url.PathEscape(engine),
		url.PathEscape(clusterId),
		url.PathEscape(id),
	)

	log.Printf("[DEBUG] Will delete user %s from cluster %s from project %s", id, clusterId, serviceName)
	err := config.OVHClient.DeleteWithContext(ctx, endpoint, nil)
	if err != nil {
		return diag.FromErr(helpers.CheckDeleted(d, err, endpoint))
	}

	log.Printf("[DEBUG] Waiting for user %s to be DELETED", id)
	err = waitForCloudProjectDatabaseUserDeleted(ctx, config.OVHClient, serviceName, engine, clusterId, id, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.Errorf("timeout while waiting user %s to be DELETED: %s", id, err.Error())
	}
	log.Printf("[DEBUG] user %s is DELETED", id)

	d.SetId("")

	return nil
}

func waitForCloudProjectDatabaseUserReady(ctx context.Context, client *ovhwrap.Client, serviceName, engine string, databaseId string, userId string, timeOut time.Duration) error {
	stateConf := &retry.StateChangeConf{
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
			err := client.GetWithContext(ctx, endpoint, res)
			if err != nil {
				return res, "", err
			}

			return res, res.Status, nil
		},
		Timeout:    timeOut,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func waitForCloudProjectDatabaseUserDeleted(ctx context.Context, client *ovhwrap.Client, serviceName, engine string, databaseId string, userId string, timeOut time.Duration) error {
	stateConf := &retry.StateChangeConf{
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
			err := client.GetWithContext(ctx, endpoint, res)
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

	_, err := stateConf.WaitForStateContext(ctx)
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

func waitForCloudProjectDatabaseDatabaseReady(ctx context.Context, client *ovhwrap.Client, serviceName, engine string, serviceId string, databaseId string, timeOut time.Duration) error {
	stateConf := &retry.StateChangeConf{
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
			err := client.GetWithContext(ctx, endpoint, res)
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

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func waitForCloudProjectDatabaseDatabaseDeleted(ctx context.Context, client *ovhwrap.Client, serviceName, engine string, serviceId string, databaseId string, timeOut time.Duration) error {
	stateConf := &retry.StateChangeConf{
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
			err := client.GetWithContext(ctx, endpoint, res)
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

	_, err := stateConf.WaitForStateContext(ctx)
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

func validateCloudProjectDatabaseIntegrationEngine(v any, p cty.Path) (diags diag.Diagnostics) {
	value := v.(string)
	if value == "mongodb" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Invalid engine",
			Detail:   fmt.Sprintf("value %s is not a valid engine for integration", value),
		})
	}
	return
}

func waitForCloudProjectDatabaseIntegrationReady(ctx context.Context, client *ovhwrap.Client, serviceName, engine string, serviceId string, integrationId string, timeOut time.Duration) error {
	stateConf := &retry.StateChangeConf{
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
			err := client.GetWithContext(ctx, endpoint, res)
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

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func waitForCloudProjectDatabaseIntegrationDeleted(ctx context.Context, client *ovhwrap.Client, serviceName, engine string, serviceId string, integrationId string, timeOut time.Duration) error {
	stateConf := &retry.StateChangeConf{
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
			err := client.GetWithContext(ctx, endpoint, res)
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

	_, err := stateConf.WaitForStateContext(ctx)
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
	ID        string   `json:"id"`
	Password  string   `json:"password"`
	Roles     []string `json:"roles"`
	Status    string   `json:"status"`
	Username  string   `json:"username"`
}

func (r *CloudProjectDatabaseMongodbUserResponse) String() string {
	return fmt.Sprintf(
		"Id: %s, User: %s, Status: %s",
		r.ID,
		r.Username,
		r.Status,
	)
}

func (r CloudProjectDatabaseMongodbUserResponse) toMap() map[string]interface{} {
	obj := make(map[string]interface{})

	obj["created_at"] = r.CreatedAt
	obj["id"] = r.ID
	obj["name"] = r.Username
	obj["status"] = r.Status
	obj["roles"] = r.Roles

	return obj
}

type CloudProjectDatabaseMongodbUserCreateOpts struct {
	Name  string   `json:"name"`
	Roles []string `json:"roles"`
}

func (opts *CloudProjectDatabaseMongodbUserCreateOpts) fromResource(d *schema.ResourceData) *CloudProjectDatabaseMongodbUserCreateOpts {
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

func (opts *CloudProjectDatabaseMongodbUserUpdateOpts) fromResource(d *schema.ResourceData) *CloudProjectDatabaseMongodbUserUpdateOpts {
	roles := d.Get("roles").(*schema.Set).List()
	opts.Roles = make([]string, len(roles))
	for i, e := range roles {
		opts.Roles[i] = e.(string)
	}
	return opts
}

func validateCloudProjectDatabaseMongodbUserAuthenticationDatabase(v any, p cty.Path) (diags diag.Diagnostics) {
	value := v.(string)
	if !strings.Contains(value, "@") {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Missing authentication database",
			Detail:   fmt.Sprintf("value %s does not have authentication database", value),
		})
	}
	return
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
	SnapshotEnabled          *bool                                      `json:"snapshotEnabled"`
	Type                     string                                     `json:"type"`
	WritesToCommitLogEnabled *bool                                      `json:"writesToCommitLogEnabled"`
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
	snapshotEnabled, ok := d.GetOkExists("snapshot_enabled")
	if ok {
		snapshotEnabledBool := snapshotEnabled.(bool)
		opts.SnapshotEnabled = &snapshotEnabledBool
	}
	opts.Type = "aggregated"
	writesToCommitLogEnabled, ok := d.GetOkExists("writes_to_commit_log_enabled")
	if ok {
		writesToCommitLogEnabledBool := writesToCommitLogEnabled.(bool)
		opts.WritesToCommitLogEnabled = &writesToCommitLogEnabledBool
	}

	return opts
}

type CloudProjectDatabaseM3dbNamespaceUpdateOpts struct {
	Resolution               string                                     `json:"resolution,omitempty"`
	Retention                CloudProjectDatabaseM3dbNamespaceRetention `json:"retention,omitempty"`
	SnapshotEnabled          *bool                                      `json:"snapshotEnabled,omitempty"`
	WritesToCommitLogEnabled *bool                                      `json:"writesToCommitLogEnabled,omitempty"`
}

func (opts *CloudProjectDatabaseM3dbNamespaceUpdateOpts) FromResource(d *schema.ResourceData) *CloudProjectDatabaseM3dbNamespaceUpdateOpts {
	opts.Resolution = d.Get("resolution").(string)
	opts.Retention = CloudProjectDatabaseM3dbNamespaceRetention{
		BlockDataExpirationDuration: d.Get("retention_block_data_expiration_duration").(string),
		BufferFutureDuration:        d.Get("retention_buffer_future_duration").(string),
		BufferPastDuration:          d.Get("retention_buffer_past_duration").(string),
		PeriodDuration:              d.Get("retention_period_duration").(string),
	}
	if d.HasChange("snapshot_enabled") {
		snapshotEnabledBool := d.Get("snapshot_enabled").(bool)
		opts.SnapshotEnabled = &snapshotEnabledBool
	}
	if d.HasChange("writes_to_commit_log_enabled") {
		writesToCommitLogEnabledBool := d.Get("writes_to_commit_log_enabled").(bool)
		opts.WritesToCommitLogEnabled = &writesToCommitLogEnabledBool
	}

	return opts
}

func waitForCloudProjectDatabaseM3dbNamespaceReady(ctx context.Context, client *ovhwrap.Client, serviceName, databaseId string, namespaceId string, timeOut time.Duration) error {
	stateConf := &retry.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"READY"},
		Refresh: func() (interface{}, string, error) {
			res := &CloudProjectDatabaseM3dbNamespaceResponse{}
			endpoint := fmt.Sprintf("/cloud/project/%s/database/m3db/%s/namespace/%s",
				url.PathEscape(serviceName),
				url.PathEscape(databaseId),
				url.PathEscape(namespaceId),
			)
			err := client.GetWithContext(ctx, endpoint, res)
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

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func waitForCloudProjectDatabaseM3dbNamespaceDeleted(ctx context.Context, client *ovhwrap.Client, serviceName, databaseId string, namespaceId string, timeOut time.Duration) error {
	stateConf := &retry.StateChangeConf{
		Pending: []string{"DELETING"},
		Target:  []string{"DELETED"},
		Refresh: func() (interface{}, string, error) {
			res := &CloudProjectDatabaseM3dbNamespaceResponse{}
			endpoint := fmt.Sprintf("/cloud/project/%s/database/m3db/%s/namespace/%s",
				url.PathEscape(serviceName),
				url.PathEscape(databaseId),
				url.PathEscape(namespaceId),
			)
			err := client.GetWithContext(ctx, endpoint, res)
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

	_, err := stateConf.WaitForStateContext(ctx)
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

func waitForCloudProjectDatabaseOpensearchPatternReady(ctx context.Context, client *ovhwrap.Client, serviceName, databaseId string, patternId string, timeOut time.Duration) error {
	stateConf := &retry.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"READY"},
		Refresh: func() (interface{}, string, error) {
			res := &CloudProjectDatabaseOpensearchPatternResponse{}
			endpoint := fmt.Sprintf("/cloud/project/%s/database/opensearch/%s/pattern/%s",
				url.PathEscape(serviceName),
				url.PathEscape(databaseId),
				url.PathEscape(patternId),
			)
			err := client.GetWithContext(ctx, endpoint, res)
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

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func waitForCloudProjectDatabaseOpensearchPatternDeleted(ctx context.Context, client *ovhwrap.Client, serviceName, databaseId string, patternId string, timeOut time.Duration) error {
	stateConf := &retry.StateChangeConf{
		Pending: []string{"DELETING"},
		Target:  []string{"DELETED"},
		Refresh: func() (interface{}, string, error) {
			res := &CloudProjectDatabaseOpensearchPatternResponse{}
			endpoint := fmt.Sprintf("/cloud/project/%s/database/opensearch/%s/pattern/%s",
				url.PathEscape(serviceName),
				url.PathEscape(databaseId),
				url.PathEscape(patternId),
			)
			err := client.GetWithContext(ctx, endpoint, res)
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

	_, err := stateConf.WaitForStateContext(ctx)
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

func validateIsSupEqual(v, min int) (diags diag.Diagnostics) {
	if v < min {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Inferior of Min value",
			Detail:   fmt.Sprintf("value %d is inferior of min value %d", v, min),
		})
	}
	return
}

func validateCloudProjectDatabaseKafkaTopicMinInsyncReplicasFunc(v any, p cty.Path) (diags diag.Diagnostics) {
	diags = validateIsSupEqual(v.(int), 1)
	return
}

func validateCloudProjectDatabaseKafkaTopicPartitionsFunc(v any, p cty.Path) (diags diag.Diagnostics) {
	diags = validateIsSupEqual(v.(int), 1)
	return
}

func validateCloudProjectDatabaseKafkaTopicReplicationFunc(v any, p cty.Path) (diags diag.Diagnostics) {
	diags = validateIsSupEqual(v.(int), 2)
	return
}

func validateCloudProjectDatabaseKafkaTopicRetentionHoursFunc(v any, p cty.Path) (diags diag.Diagnostics) {
	diags = validateIsSupEqual(v.(int), -1)
	return
}

func waitForCloudProjectDatabaseKafkaTopicReady(ctx context.Context, client *ovhwrap.Client, serviceName, databaseId string, topicId string, timeOut time.Duration) error {
	stateConf := &retry.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"READY"},
		Refresh: func() (interface{}, string, error) {
			res := &CloudProjectDatabaseKafkaTopicResponse{}
			endpoint := fmt.Sprintf("/cloud/project/%s/database/kafka/%s/topic/%s",
				url.PathEscape(serviceName),
				url.PathEscape(databaseId),
				url.PathEscape(topicId),
			)
			err := client.GetWithContext(ctx, endpoint, res)
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

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func waitForCloudProjectDatabaseKafkaTopicDeleted(ctx context.Context, client *ovhwrap.Client, serviceName, databaseId string, topicId string, timeOut time.Duration) error {
	stateConf := &retry.StateChangeConf{
		Pending: []string{"DELETING"},
		Target:  []string{"DELETED"},
		Refresh: func() (interface{}, string, error) {
			res := &CloudProjectDatabaseKafkaTopicResponse{}
			endpoint := fmt.Sprintf("/cloud/project/%s/database/kafka/%s/topic/%s",
				url.PathEscape(serviceName),
				url.PathEscape(databaseId),
				url.PathEscape(topicId),
			)
			err := client.GetWithContext(ctx, endpoint, res)
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

	_, err := stateConf.WaitForStateContext(ctx)
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

func waitForCloudProjectDatabaseKafkaAclReady(ctx context.Context, client *ovhwrap.Client, serviceName, databaseId string, aclId string, timeOut time.Duration) error {
	stateConf := &retry.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"READY"},
		Refresh: func() (interface{}, string, error) {
			res := &CloudProjectDatabaseKafkaTopicResponse{}
			endpoint := fmt.Sprintf("/cloud/project/%s/database/kafka/%s/acl/%s",
				url.PathEscape(serviceName),
				url.PathEscape(databaseId),
				url.PathEscape(aclId),
			)
			err := client.GetWithContext(ctx, endpoint, res)
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

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func waitForCloudProjectDatabaseKafkaAclDeleted(ctx context.Context, client *ovhwrap.Client, serviceName, databaseId string, aclId string, timeOut time.Duration) error {
	stateConf := &retry.StateChangeConf{
		Pending: []string{"DELETING"},
		Target:  []string{"DELETED"},
		Refresh: func() (interface{}, string, error) {
			res := &CloudProjectDatabaseKafkaTopicResponse{}
			endpoint := fmt.Sprintf("/cloud/project/%s/database/kafka/%s/acl/%s",
				url.PathEscape(serviceName),
				url.PathEscape(databaseId),
				url.PathEscape(aclId),
			)
			err := client.GetWithContext(ctx, endpoint, res)
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

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

// // Schema Registry ACL

type CloudProjectDatabaseKafkaSchemaRegistryAclResponse struct {
	Id         string `json:"id"`
	Permission string `json:"permission"`
	Resource   string `json:"resource"`
	Username   string `json:"username"`
}

func (p *CloudProjectDatabaseKafkaSchemaRegistryAclResponse) String() string {
	return fmt.Sprintf(
		"Id: %s, Permission: %s, affected User: %s, affected Topic %s",
		p.Id,
		p.Permission,
		p.Username,
		p.Resource,
	)
}

func (v CloudProjectDatabaseKafkaSchemaRegistryAclResponse) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})

	obj["id"] = v.Id
	obj["permission"] = v.Permission
	obj["resource"] = v.Resource
	obj["username"] = v.Username

	return obj
}

type CloudProjectDatabaseKafkaSchemaRegistryAclCreateOpts struct {
	Permission string `json:"permission"`
	Resource   string `json:"resource"`
	Username   string `json:"username"`
}

func (opts *CloudProjectDatabaseKafkaSchemaRegistryAclCreateOpts) FromResource(d *schema.ResourceData) *CloudProjectDatabaseKafkaSchemaRegistryAclCreateOpts {
	opts.Permission = d.Get("permission").(string)
	opts.Resource = d.Get("resource").(string)
	opts.Username = d.Get("username").(string)

	return opts
}

func waitForCloudProjectDatabaseKafkaSchemaRegistryAclReady(ctx context.Context, client *ovhwrap.Client, serviceName, databaseId string, schemaRegistryAclId string, timeOut time.Duration) error {
	stateConf := &retry.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"READY"},
		Refresh: func() (interface{}, string, error) {
			res := &CloudProjectDatabaseKafkaTopicResponse{}
			endpoint := fmt.Sprintf("/cloud/project/%s/database/kafka/%s/schemaRegistryAcl/%s",
				url.PathEscape(serviceName),
				url.PathEscape(databaseId),
				url.PathEscape(schemaRegistryAclId),
			)
			err := client.GetWithContext(ctx, endpoint, res)
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

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func waitForCloudProjectDatabaseKafkaSchemaRegistryAclDeleted(ctx context.Context, client *ovhwrap.Client, serviceName, databaseId string, schemaRegistryAclId string, timeOut time.Duration) error {
	stateConf := &retry.StateChangeConf{
		Pending: []string{"DELETING"},
		Target:  []string{"DELETED"},
		Refresh: func() (interface{}, string, error) {
			res := &CloudProjectDatabaseKafkaTopicResponse{}
			endpoint := fmt.Sprintf("/cloud/project/%s/database/kafka/%s/schemaRegistryAcl/%s",
				url.PathEscape(serviceName),
				url.PathEscape(databaseId),
				url.PathEscape(schemaRegistryAclId),
			)
			err := client.GetWithContext(ctx, endpoint, res)
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

	_, err := stateConf.WaitForStateContext(ctx)
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

// Connection Pool

type CloudProjectDatabasePostgresqlConnectionPoolCreateOpts struct {
	DatabaseID string `json:"databaseId"`
	Mode       string `json:"mode"`
	Name       string `json:"name"`
	Size       int    `json:"size"`
	UserID     string `json:"userId,omitempty"`
}

func (opts *CloudProjectDatabasePostgresqlConnectionPoolCreateOpts) fromResource(d *schema.ResourceData) *CloudProjectDatabasePostgresqlConnectionPoolCreateOpts {
	opts.DatabaseID = d.Get("database_id").(string)
	opts.Mode = d.Get("mode").(string)
	opts.Name = d.Get("name").(string)
	opts.Size = d.Get("size").(int)
	opts.UserID = d.Get("user_id").(string)
	return opts
}

type CloudProjectDatabasePostgresqlConnectionPoolUpdateOpts struct {
	DatabaseID string `json:"databaseId,omitempty"`
	Mode       string `json:"mode,omitempty"`
	Size       int    `json:"size,omitempty"`
	UserID     string `json:"userId,omitempty"`
}

func (opts *CloudProjectDatabasePostgresqlConnectionPoolUpdateOpts) fromResource(d *schema.ResourceData) *CloudProjectDatabasePostgresqlConnectionPoolUpdateOpts {
	opts.DatabaseID = d.Get("database_id").(string)
	opts.Mode = d.Get("mode").(string)
	opts.Size = d.Get("size").(int)
	opts.UserID = d.Get("user_id").(string)
	return opts
}

type CloudProjectDatabasePostgresqlConnectionPoolResponse struct {
	DatabaseID string `json:"databaseId"`
	ID         string `json:"id"`
	Mode       string `json:"mode"`
	Name       string `json:"name"`
	Port       int64  `json:"port"`
	Size       int    `json:"size"`
	SslMode    string `json:"sslMode"`
	URI        string `json:"uri"`
	UserID     string `json:"userId"`
}

func (r CloudProjectDatabasePostgresqlConnectionPoolResponse) toMap() map[string]interface{} {
	obj := make(map[string]interface{})

	obj["database_id"] = r.DatabaseID
	obj["id"] = r.ID
	obj["mode"] = r.Mode
	obj["name"] = r.Name
	obj["port"] = r.Port
	obj["size"] = r.Size
	obj["ssl_mode"] = r.SslMode
	obj["uri"] = r.URI
	obj["user_id"] = r.UserID

	return obj
}

func waitForCloudProjectDatabasePostgresqlConnectionPoolReady(ctx context.Context, client *ovhwrap.Client, serviceName, databaseID, connectionPoolID string, timeOut time.Duration) error {
	stateConf := &retry.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"READY"},
		Refresh: func() (interface{}, string, error) {
			res := &CloudProjectDatabasePostgresqlConnectionPoolResponse{}
			endpoint := fmt.Sprintf("/cloud/project/%s/database/postgresql/%s/connectionPool/%s",
				url.PathEscape(serviceName),
				url.PathEscape(databaseID),
				url.PathEscape(connectionPoolID),
			)
			err := client.GetWithContext(ctx, endpoint, res)
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

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func waitForCloudProjectDatabasePostgresqlConnectionPoolDeleted(ctx context.Context, client *ovhwrap.Client, serviceName, databaseID, connectionPoolID string, timeOut time.Duration) error {
	stateConf := &retry.StateChangeConf{
		Pending: []string{"DELETING"},
		Target:  []string{"DELETED"},
		Refresh: func() (interface{}, string, error) {
			res := &CloudProjectDatabasePostgresqlConnectionPoolResponse{}
			endpoint := fmt.Sprintf("/cloud/project/%s/database/postgresql/%s/connectionPool/%s",
				url.PathEscape(serviceName),
				url.PathEscape(databaseID),
				url.PathEscape(connectionPoolID),
			)
			err := client.GetWithContext(ctx, endpoint, res)
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

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

// Prometheus

type CloudProjectDatabasePrometheusCreateOpts struct {
	EnablePrometheus bool `json:"enablePrometheus"`
}

type CloudProjectDatabasePrometheusAccessResponse struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

func (r CloudProjectDatabasePrometheusAccessResponse) toMap() map[string]interface{} {
	obj := make(map[string]interface{})

	obj["password"] = r.Password
	obj["username"] = r.Username

	return obj
}

type CloudProjectDatabasePrometheusEndpointResponse struct {
	Username string                                         `json:"username"`
	Targets  []CloudProjectDatabasePrometheusEndpointTarget `json:"targets"`
}

type CloudProjectDatabasePrometheusEndpointTarget struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

func (r CloudProjectDatabasePrometheusEndpointResponse) toMap() map[string]interface{} {
	obj := make(map[string]interface{})

	obj["username"] = r.Username

	var targets []map[string]interface{}
	for _, t := range r.Targets {
		targets = append(targets, t.toMap())
	}
	obj["targets"] = targets

	return obj
}

func (r CloudProjectDatabasePrometheusEndpointTarget) toMap() map[string]interface{} {
	obj := make(map[string]interface{})

	obj["host"] = r.Host
	obj["port"] = r.Port

	return obj
}

type CloudProjectDatabaseMongodbPrometheusEndpointResponse struct {
	Username  string `json:"username"`
	SrvDomain string `json:"srvDomain"`
}

func (r CloudProjectDatabaseMongodbPrometheusEndpointResponse) toMap() map[string]interface{} {
	obj := make(map[string]interface{})

	obj["username"] = r.Username
	obj["srv_domain"] = r.SrvDomain

	return obj
}

func enableCloudProjectDatabasePrometheus(ctx context.Context, d *schema.ResourceData, meta interface{}, engine string, enablePrometheus bool, updateFunc schema.UpdateContextFunc) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	clusterID := d.Get("cluster_id").(string)

	serviceEndpoint := fmt.Sprintf("/cloud/project/%s/database/%s/%s",
		url.PathEscape(serviceName),
		url.PathEscape(engine),
		url.PathEscape(clusterID),
	)

	params := &CloudProjectDatabasePrometheusCreateOpts{
		EnablePrometheus: enablePrometheus,
	}
	res := &CloudProjectDatabaseResponse{}
	log.Printf("[DEBUG] Will update database: %+v", params)
	err := config.OVHClient.PutWithContext(ctx, serviceEndpoint, params, res)
	if err != nil {
		return diag.Errorf("calling Put %s with params %v:\n\t %q", serviceEndpoint, params, err)
	}

	log.Printf("[DEBUG] Waiting for database %s to be READY", res.ID)
	err = waitForCloudProjectDatabaseReady(ctx, config.OVHClient, serviceName, engine, res.ID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("timeout while waiting database %s to be READY: %s", res.ID, err.Error())
	}
	log.Printf("[DEBUG] database %s is READY", res.ID)

	if enablePrometheus {
		d.SetId(res.ID)
		return updateFunc(ctx, d, meta)
	}
	d.SetId("")
	return nil
}

func updateCloudProjectDatabasePrometheus(ctx context.Context, d *schema.ResourceData, meta interface{}, engine string, readFunc schema.ReadContextFunc) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	id := d.Id()

	endpoint := fmt.Sprintf("/cloud/project/%s/database/%s/%s/prometheus/credentials/reset",
		url.PathEscape(serviceName),
		url.PathEscape(engine),
		url.PathEscape(id),
	)

	res := &CloudProjectDatabasePrometheusAccessResponse{}
	err := config.OVHClient.PostWithContext(ctx, endpoint, nil, res)
	if err != nil {
		return diag.Errorf("calling Post %s:\n\t %q", endpoint, err)
	}

	d.Set("password", res.Password)

	return readFunc(ctx, d, meta)
}
