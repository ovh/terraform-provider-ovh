package ovh

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/go-ovh/ovh"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

type CloudProjectDatabaseResponse struct {
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
	Status          string                         `json:"status"`
	SubnetId        string                         `json:"subnetId"`
	Version         string                         `json:"version"`
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
	Plan         string                           `json:"plan"`
	SubnetId     string                           `json:"subnetId,omitempty"`
	Version      string                           `json:"version"`
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

	return nil, opts
}

type CloudProjectDatabaseUpdateOpts struct {
	Description string `json:"description,omitempty"`
	Flavor      string `json:"flavor,omitempty"`
	NodeNumber  int    `json:"nodeNumber,omitempty"`
	Plan        string `json:"plan,omitempty"`
	Version     string `json:"version,omitempty"`
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
				}
				return res, "", err
			}

			return res, res.Status, nil
		},
		Timeout:      30 * time.Minute,
		Delay:        30 * time.Second,
		PollInterval: 20 * time.Second,
	}

	_, err := stateConf.WaitForState()
	return err
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

func validateCloudProjectDatabaseUserEngineFunc(v interface{}, k string) (ws []string, errors []error) {
	err := helpers.ValidateStringEnum(v.(string), []string{
		"cassandra",
		"mysql",
		"kafka",
		"kafkaConnect",
	})

	if err != nil {
		errors = append(errors, err)
	}
	return
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

// RedisUser

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

// // Certificates

type CloudProjectDatabaseKafkaCertificatesResponse struct {
	Ca string `json:"ca"`
}

func (p *CloudProjectDatabaseKafkaCertificatesResponse) String() string {
	return fmt.Sprintf(
		"Ca: %s",
		p.Ca,
	)
}

func (v CloudProjectDatabaseKafkaCertificatesResponse) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})

	obj["ca"] = v.Ca

	return obj
}

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
