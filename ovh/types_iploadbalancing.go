package ovh

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type IpLoadbalancing struct {
	IPv6             string                          `json:"ipv6,omitempty"`
	IPv4             string                          `json:"ipv4,omitempty"`
	MetricsToken     string                          `json:"metricsToken,omitempty"`
	Zone             []string                        `json:"zone"`
	Offer            string                          `json:"offer"`
	ServiceName      string                          `json:"serviceName"`
	IpLoadbalancing  string                          `json:"ipLoadbalancing"`
	State            string                          `json:"state"`
	OrderableZones   []*IpLoadbalancingOrderableZone `json:"orderableZone"`
	VrackEligibility bool                            `json:"vrackEligibility"`
	VrackName        string                          `json:"vrackName"`
	SslConfiguration string                          `json:"sslConfiguration"`
	DisplayName      string                          `json:"displayName"`
}

type IpLoadbalancingOrderableZone struct {
	Name     string `json:"name"`
	PlanCode string `json:"plan_code"`
}

type IPLoadbalancingRefreshTask struct {
	CreationDate string   `json:"creationDate"`
	Status       string   `json:"status"`
	Progress     int      `json:"progress"`
	Action       string   `json:"action"`
	ID           int      `json:"id"`
	DoneDate     string   `json:"doneDate"`
	Zones        []string `json:"zones"`
}

type IPLoadbalancingRefreshPending struct {
	Number int    `json:"number"`
	Zone   string `json:"zone"`
}

type IPLoadbalancingRefreshPendings []IPLoadbalancingRefreshPending

type IpLoadbalancingTcpFarm struct {
	FarmId         int                              `json:"farmId,omitempty"`
	Zone           string                           `json:"zone,omitempty"`
	VrackNetworkId int                              `json:"vrackNetworkId,omitempty"`
	Port           int                              `json:"port,omitempty"`
	Stickiness     string                           `json:"stickiness,omitempty"`
	Balance        string                           `json:"balance,omitempty"`
	Probe          *IpLoadbalancingFarmBackendProbe `json:"probe,omitempty"`
	DisplayName    string                           `json:"displayName,omitempty"`
}

type IpLoadbalancingFarmBackendProbe struct {
	Match    *string `json:"match,omitempty"`
	Port     *int    `json:"port,omitempty"`
	Interval *int    `json:"interval,omitempty"`
	Negate   *bool   `json:"negate,omitempty"`
	Pattern  *string `json:"pattern,omitempty"`
	ForceSsl *bool   `json:"forceSsl,omitempty"`
	URL      *string `json:"url,omitempty"`
	Method   *string `json:"method,omitempty"`
	Type     *string `json:"type,omitempty"`
}

func (opts *IpLoadbalancingFarmBackendProbe) FromResource(d *schema.ResourceData, parent string) *IpLoadbalancingFarmBackendProbe {
	opts.Match = getNilStringPointerFromData(d, fmt.Sprintf("%s.match", parent))
	opts.Port = getNilIntPointerFromData(d, fmt.Sprintf("%s.port", parent))
	opts.Interval = getNilIntPointerFromData(d, fmt.Sprintf("%s.interval", parent))
	opts.Negate = getNilBoolPointerFromData(d, fmt.Sprintf("%s.negate", parent))
	opts.Pattern = getNilStringPointerFromData(d, fmt.Sprintf("%s.pattern", parent))
	opts.ForceSsl = getNilBoolPointerFromData(d, fmt.Sprintf("%s.force_ssl", parent))
	opts.URL = getNilStringPointerFromData(d, fmt.Sprintf("%s.url", parent))
	opts.Method = getNilStringPointerFromData(d, fmt.Sprintf("%s.method", parent))
	opts.Type = getNilStringPointerFromData(d, fmt.Sprintf("%s.type", parent))
	return opts
}

func (v IpLoadbalancingFarmBackendProbe) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})

	isNil := true
	if v.Match != nil {
		isNil = false
		obj["match"] = *v.Match
	}
	if v.Port != nil {
		isNil = false
		obj["port"] = *v.Port
	}
	if v.Interval != nil {
		isNil = false
		obj["interval"] = *v.Interval
	}
	if v.Negate != nil {
		isNil = false
		obj["negate"] = *v.Negate
	}
	if v.Pattern != nil {
		isNil = false
		obj["pattern"] = *v.Pattern
	}
	if v.ForceSsl != nil {
		isNil = false
		obj["force_ssl"] = *v.ForceSsl
	}
	if v.URL != nil {
		isNil = false
		obj["url"] = *v.URL
	}
	if v.Method != nil {
		isNil = false
		obj["method"] = *v.Method
	}
	if v.Type != nil {
		isNil = false
		obj["type"] = *v.Type
	}

	if isNil {
		return nil
	}

	return obj
}

type IpLoadbalancingHttpFarm struct {
	FarmId         int                              `json:"farmId,omitempty"`
	Zone           string                           `json:"zone,omitempty"`
	VrackNetworkId int                              `json:"vrackNetworkId,omitempty"`
	Port           int                              `json:"port,omitempty"`
	Stickiness     string                           `json:"stickiness,omitempty"`
	Balance        string                           `json:"balance,omitempty"`
	Probe          *IpLoadbalancingFarmBackendProbe `json:"probe,omitempty"`
	DisplayName    string                           `json:"displayName,omitempty"`
}

// IPLoadbalancingRouteHTTPAction Action triggered when all rules match
type IPLoadbalancingRouteHTTPAction struct {
	Target string `json:"target,omitempty"` // Farm ID for "farm" action type or URL template for "redirect" action. You may use ${uri}, ${protocol}, ${host}, ${port} and ${path} variables in redirect target
	Status int    `json:"status,omitempty"` // HTTP status code for "redirect" and "reject" actions
	Type   string `json:"type,omitempty"`   // Action to trigger if all the rules of this route matches
}

//IPLoadbalancingRouteHTTP HTTP Route
type IPLoadbalancingRouteHTTP struct {
	Status      string                          `json:"status,omitempty"`      //Route status. Routes in "ok" state are ready to operate
	Weight      int                             `json:"weight,omitempty"`      //Route priority ([0..255]). 0 if null. Highest priority routes are evaluated first. Only the first matching route will trigger an action
	Action      *IPLoadbalancingRouteHTTPAction `json:"action,omitempty"`      //Action triggered when all rules match
	RouteID     int                             `json:"routeId,omitempty"`     //Id of your route
	DisplayName string                          `json:"displayName,omitempty"` //Human readable name for your route, this field is for you
	FrontendID  int                             `json:"frontendId,omitempty"`  //Route traffic for this frontend
}

type IpLoadbalancingTcpFrontend struct {
	FrontendId    int      `json:"frontendId,omitempty"`
	Port          string   `json:"port"`
	Zone          string   `json:"zone"`
	AllowedSource []string `json:"allowedSource,omitempty"`
	DedicatedIpFo []string `json:"dedicatedIpfo,omitempty"`
	DefaultFarmId *int     `json:"defaultFarmId,omitempty"`
	DefaultSslId  *int     `json:"defaultSslId,omitempty"`
	Disabled      *bool    `json:"disabled"`
	Ssl           *bool    `json:"ssl"`
	DisplayName   string   `json:"displayName,omitempty"`
}

type IpLoadbalancingHttpFrontend struct {
	FrontendId    int      `json:"frontendId,omitempty"`
	Port          string   `json:"port"`
	Zone          string   `json:"zone"`
	AllowedSource []string `json:"allowedSource,omitempty"`
	DedicatedIpFo []string `json:"dedicatedIpfo,omitempty"`
	DefaultFarmId *int     `json:"defaultFarmId,omitempty"`
	DefaultSslId  *int     `json:"defaultSslId,omitempty"`
	Disabled      *bool    `json:"disabled"`
	Ssl           *bool    `json:"ssl"`
	DisplayName   string   `json:"displayName,omitempty"`
}

//IPLoadbalancingRouteHTTPRule HTTP Route Rule
type IPLoadbalancingRouteHTTPRule struct {
	RuleID      int    `json:"ruleId,omitempty"`      //Id of your rule
	RouteID     int    `json:"routeId,omitempty"`     //Id of your route
	DisplayName string `json:"displayName,omitempty"` //Human readable name for your rule
	Field       string `json:"field,omitempty"`       //Name of the field to match like "protocol" or "host". See "/ipLoadbalancing/{serviceName}/availableRouteRules" for a list of available rules
	Match       string `json:"match,omitempty"`       //Matching operator. Not all operators are available for all fields. See "/ipLoadbalancing/{serviceName}/availableRouteRules"
	Negate      bool   `json:"negate,omitempty"`      //Invert the matching operator effect
	Pattern     string `json:"pattern,omitempty"`     //Value to match against this match. Interpretation if this field depends on the match and field
	SubField    string `json:"subField,omitempty"`    //Name of sub-field, if applicable. This may be a Cookie or Header name for instance
}

type IpLoadbalancingTcpFarmServer struct {
	BackendId            int     `json:"backendId,omitempty"`
	ServerId             int     `json:"serverId,omitempty"`
	FarmId               int     `json:"farmId,omitempty"`
	DisplayName          *string `json:"displayName,omitempty"`
	Address              string  `json:"address"`
	Cookie               *string `json:"cookie,omitempty"`
	Port                 *int    `json:"port"`
	ProxyProtocolVersion *string `json:"proxyProtocolVersion"`
	Chain                *string `json:"chain"`
	Weight               *int    `json:"weight"`
	Probe                *bool   `json:"probe"`
	Ssl                  *bool   `json:"ssl"`
	Backup               *bool   `json:"backup"`
	Status               string  `json:"status"`
}

type IpLoadbalancingHttpFarmServer struct {
	BackendId            int     `json:"backendId,omitempty"`
	ServerId             int     `json:"serverId,omitempty"`
	FarmId               int     `json:"farmId,omitempty"`
	DisplayName          *string `json:"displayName,omitempty"`
	Address              string  `json:"address"`
	Cookie               *string `json:"cookie,omitempty"`
	Port                 *int    `json:"port"`
	ProxyProtocolVersion *string `json:"proxyProtocolVersion"`
	Chain                *string `json:"chain"`
	Weight               *int    `json:"weight"`
	Probe                *bool   `json:"probe"`
	Ssl                  *bool   `json:"ssl"`
	Backup               *bool   `json:"backup"`
	Status               string  `json:"status"`
}

type IpLoadbalancingDefinedFarm struct {
	Type string `json:"type"`
	Id   int64  `json:"id"`
}

func (v IpLoadbalancingDefinedFarm) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["type"] = v.Type
	obj["id"] = v.Id
	return obj
}

type IpLoadbalancingVrackNetwork struct {
	Subnet         string                       `json:"subnet"`
	Vlan           int64                        `json:"vlan"`
	VrackNetworkId int64                        `json:"vrackNetworkId"`
	FarmId         []IpLoadbalancingDefinedFarm `json:"farmId`
	DisplayName    *string                      `json:"displayName,omitempty"`
	NatIp          string                       `json:"natIp"`
}

func (v IpLoadbalancingVrackNetwork) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["subnet"] = v.Subnet
	obj["vlan"] = v.Vlan
	obj["nat_ip"] = v.NatIp
	obj["vrack_network_id"] = v.VrackNetworkId

	if v.DisplayName != nil {
		obj["display_name"] = *v.DisplayName
	}

	ids := make([]interface{}, len(v.FarmId))
	for i, farm := range v.FarmId {
		ids[i] = farm.ToMap()
	}
	obj["farm_id"] = ids

	return obj
}

type IpLoadbalancingVrackNetworkCreateOpts struct {
	Subnet      string   `json:"subnet"`
	Vlan        *int64   `json:"vlan,omitempty"`
	FarmId      *[]int64 `json:"farmId,omitempty"`
	DisplayName *string  `json:"displayName,omitempty"`
	NatIp       string   `json:"natIp"`
}

func (opts *IpLoadbalancingVrackNetworkCreateOpts) FromResource(d *schema.ResourceData) *IpLoadbalancingVrackNetworkCreateOpts {
	opts.Subnet = d.Get("subnet").(string)
	opts.NatIp = d.Get("nat_ip").(string)
	opts.DisplayName = getNilStringPointerFromData(d, "display_name")
	opts.Vlan = getNilInt64PointerFromData(d, "vlan")

	if val, ok := d.GetOkExists("farm_id"); ok {
		farmId := val.([]interface{})
		arr := make([]int64, len(farmId))
		for i, id := range farmId {
			arr[i] = int64(id.(int))
		}

		opts.FarmId = &arr
	}

	return opts
}

type IpLoadbalancingVrackNetworkUpdateOpts struct {
	Vlan        *int64  `json:"vlan,omitempty"`
	DisplayName *string `json:"displayName,omitempty"`
	NatIp       *string `json:"natIp,omitempty"`
}

func (opts *IpLoadbalancingVrackNetworkUpdateOpts) FromResource(d *schema.ResourceData) *IpLoadbalancingVrackNetworkUpdateOpts {
	opts.NatIp = getNilStringPointerFromData(d, "nat_ip")
	opts.DisplayName = getNilStringPointerFromData(d, "display_name")
	opts.Vlan = getNilInt64PointerFromData(d, "vlan")

	return opts
}

type IpLoadbalancingVrackNetworkFarmIdUpdateOpts struct {
	FarmId []int64 `json:"farmId"`
}

func (opts *IpLoadbalancingVrackNetworkFarmIdUpdateOpts) FromResource(d *schema.ResourceData) *IpLoadbalancingVrackNetworkFarmIdUpdateOpts {
	opts.FarmId = []int64{}

	if val, ok := d.GetOkExists("farm_id"); ok {
		farmId := val.([]int64)
		opts.FarmId = farmId
	}

	return opts
}
