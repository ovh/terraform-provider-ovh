package ovh

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

type IpLoadbalancing struct {
	DisplayName      *string                         `json:"displayName"`
	IPv4             *string                         `json:"ipv4,omitempty"`
	IPv6             *string                         `json:"ipv6,omitempty"`
	IpLoadbalancing  string                          `json:"ipLoadbalancing"`
	MetricsToken     *string                         `json:"metricsToken,omitempty"`
	Offer            string                          `json:"offer"`
	OrderableZones   []*IpLoadbalancingOrderableZone `json:"orderableZone"`
	ServiceName      string                          `json:"serviceName"`
	SslConfiguration *string                         `json:"sslConfiguration"`
	State            string                          `json:"state"`
	VrackEligibility bool                            `json:"vrackEligibility"`
	VrackName        *string                         `json:"vrackName"`
	Zone             []string                        `json:"zone"`
}

func (v IpLoadbalancing) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})

	obj["zone"] = v.Zone
	obj["offer"] = v.Offer
	obj["service_name"] = v.ServiceName
	obj["ip_loadbalancing"] = v.IpLoadbalancing
	obj["state"] = v.State
	obj["vrack_eligibility"] = v.VrackEligibility
	obj["display_name"] = v.DisplayName

	if v.IPv4 != nil {
		obj["ipv4"] = *v.IPv4
	}

	if v.IPv6 != nil {
		obj["ipv6"] = *v.IPv6
	}

	if v.DisplayName != nil {
		obj["display_name"] = *v.DisplayName
	}

	if v.MetricsToken != nil {
		obj["metrics_token"] = *v.MetricsToken
	}

	if v.VrackName != nil {
		obj["vrack_name"] = *v.VrackName
	}

	if v.SslConfiguration != nil {
		obj["ssl_configuration"] = *v.SslConfiguration
	}

	if v.OrderableZones != nil {
		var orderableZone []map[string]interface{}
		for _, z := range v.OrderableZones {
			orderableZone = append(orderableZone, z.ToMap())
		}

		obj["orderable_zone"] = orderableZone
	}

	return obj
}

type IpLoadbalancingOrderableZone struct {
	Name     string `json:"name"`
	PlanCode string `json:"plan_code"`
}

func (v IpLoadbalancingOrderableZone) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["name"] = v.Name
	obj["plan_code"] = v.PlanCode
	return obj
}

type IpLoadbalancingUpdateOpts struct {
	DisplayName      *string `json:"displayName,omitempty"`
	SslConfiguration *string `json:"sslConfiguration,omitempty"`
}

func (opts *IpLoadbalancingUpdateOpts) FromResource(d *schema.ResourceData) *IpLoadbalancingUpdateOpts {
	opts.DisplayName = helpers.GetNilStringPointerFromData(d, "display_name")
	opts.SslConfiguration = helpers.GetNilStringPointerFromData(d, "ssl_configuration")

	return opts
}

type IpLoadbalancingConfirmTerminationOpts struct {
	Token string `json:"token"`
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

type IpLoadbalancingFarmCreateOrUpdateOpts struct {
	Balance        *string                          `json:"balance,omitempty"`
	DisplayName    *string                          `json:"displayName,omitempty"`
	Port           *int                             `json:"port,omitempty"`
	Probe          *IpLoadbalancingFarmBackendProbe `json:"probe,omitempty"`
	Stickiness     *string                          `json:"stickiness,omitempty"`
	VrackNetworkId *int64                           `json:"vrackNetworkId,omitempty"`
	Zone           string                           `json:"zone"`
}

func (opts *IpLoadbalancingFarmCreateOrUpdateOpts) FromResource(d *schema.ResourceData) *IpLoadbalancingFarmCreateOrUpdateOpts {
	opts.Balance = helpers.GetNilStringPointerFromData(d, "balance")
	opts.DisplayName = helpers.GetNilStringPointerFromData(d, "display_name")
	opts.Port = helpers.GetNilIntPointerFromData(d, "port")
	opts.Stickiness = helpers.GetNilStringPointerFromData(d, "stickiness")
	opts.VrackNetworkId = helpers.GetNilInt64PointerFromData(d, "vrack_network_id")
	opts.Zone = d.Get("zone").(string)

	probe := d.Get("probe").([]interface{})
	if probe != nil && len(probe) == 1 {
		opts.Probe = (&IpLoadbalancingFarmBackendProbe{}).FromResource(d, "probe.0")
	}

	return opts
}

type IpLoadbalancingFarm struct {
	Balance        *string                          `json:"balance,omitempty"`
	DisplayName    *string                          `json:"displayName,omitempty"`
	FarmId         int                              `json:"farmId"`
	Port           *int                             `json:"port,omitempty"`
	Probe          *IpLoadbalancingFarmBackendProbe `json:"probe,omitempty"`
	Stickiness     *string                          `json:"stickiness,omitempty"`
	VrackNetworkId *int64                           `json:"vrackNetworkId,omitempty"`
	Zone           string                           `json:"zone"`
}

type IpLoadbalancingFarmBackendProbe struct {
	ForceSsl *bool   `json:"forceSsl"`
	Interval *int    `json:"interval,omitempty"`
	Match    *string `json:"match,omitempty"`
	Method   *string `json:"method,omitempty"`
	Negate   *bool   `json:"negate"`
	Pattern  *string `json:"pattern,omitempty"`
	Port     *int    `json:"port,omitempty"`
	Type     *string `json:"type,omitempty"`
	URL      *string `json:"url,omitempty"`
}

func (opts *IpLoadbalancingFarmBackendProbe) FromResource(d *schema.ResourceData, parent string) *IpLoadbalancingFarmBackendProbe {
	opts.Match = helpers.GetNilStringPointerFromData(d, fmt.Sprintf("%s.match", parent))
	opts.Port = helpers.GetNilIntPointerFromData(d, fmt.Sprintf("%s.port", parent))
	opts.Interval = helpers.GetNilIntPointerFromData(d, fmt.Sprintf("%s.interval", parent))
	opts.Negate = helpers.GetNilBoolPointerFromData(d, fmt.Sprintf("%s.negate", parent))
	opts.Pattern = helpers.GetNilStringPointerFromData(d, fmt.Sprintf("%s.pattern", parent))
	opts.ForceSsl = helpers.GetNilBoolPointerFromData(d, fmt.Sprintf("%s.force_ssl", parent))
	opts.URL = helpers.GetNilStringPointerFromData(d, fmt.Sprintf("%s.url", parent))
	opts.Method = helpers.GetNilStringPointerFromData(d, fmt.Sprintf("%s.method", parent))
	opts.Type = helpers.GetNilStringPointerFromData(d, fmt.Sprintf("%s.type", parent))

	// Error 400: "A probe can not negate without a match"
	if opts.Match == nil {
		opts.Negate = nil
	}

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

// IPLoadbalancingHttpRouteAction Action triggered when all rules match
type IPLoadbalancingHttpRouteActionOpts struct {
	Target *string `json:"target,omitempty"` // Farm ID for "farm" action type or URL template for "redirect" action. You may use ${uri}, ${protocol}, ${host}, ${port} and ${path} variables in redirect target
	Status *int64  `json:"status,omitempty"` // HTTP status code for "redirect" and "reject" actions
	Type   string  `json:"type"`             // Action to trigger if all the rules of this route matches
}

func (opts *IPLoadbalancingHttpRouteActionOpts) FromResource(d *schema.ResourceData, parent string) *IPLoadbalancingHttpRouteActionOpts {
	opts.Target = helpers.GetNilStringPointerFromData(d, fmt.Sprintf("%s.target", parent))
	opts.Status = helpers.GetNilInt64PointerFromData(d, fmt.Sprintf("%s.status", parent))
	opts.Type = d.Get(fmt.Sprintf("%s.type", parent)).(string)
	return opts
}

//IPLoadbalancingHttpRoute HTTP Route
type IPLoadbalancingHttpRouteOpts struct {
	Action      IPLoadbalancingHttpRouteActionOpts `json:"action"`                //Action triggered when all rules match
	DisplayName *string                            `json:"displayName,omitempty"` //Human readable name for your route, this field is for you
	FrontendId  *int64                             `json:"frontendId,omitempty"`  //Route traffic for this frontend
	Weight      *int64                             `json:"weight,omitempty"`      //Route priority ([0..255]). 0 if null. Highest priority routes are evaluated first. Only the first matching route will trigger an action
}

func (opts *IPLoadbalancingHttpRouteOpts) FromResource(d *schema.ResourceData) *IPLoadbalancingHttpRouteOpts {
	opts.DisplayName = helpers.GetNilStringPointerFromData(d, "display_name")

	actions := d.Get("action").([]interface{})
	if actions != nil && len(actions) == 1 {
		opts.Action = *((&IPLoadbalancingHttpRouteActionOpts{}).FromResource(d, "action.0"))
	}

	opts.FrontendId = helpers.GetNilInt64PointerFromData(d, "frontend_id")
	opts.Weight = helpers.GetNilInt64PointerFromData(d, "weight")

	return opts
}

// IPLoadbalancingHttpRouteAction Action triggered when all rules match
type IPLoadbalancingHttpRouteAction struct {
	Target *string `json:"target"` // Farm ID for "farm" action type or URL template for "redirect" action. You may use ${uri}, ${protocol}, ${host}, ${port} and ${path} variables in redirect target
	Status *int64  `json:"status"` // HTTP status code for "redirect" and "reject" actions
	Type   string  `json:"type"`   // Action to trigger if all the rules of this route matches
}

func (v IPLoadbalancingHttpRouteAction) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})

	if v.Target != nil {
		obj["target"] = *v.Target
	}

	if v.Status != nil {
		obj["status"] = *v.Status
	}

	obj["type"] = v.Type

	return obj
}

//IPLoadbalancingHttpRoute HTTP Route
type IPLoadbalancingHttpRoute struct {
	Action      IPLoadbalancingHttpRouteAction `json:"action"`      //Action triggered when all rules match
	DisplayName *string                        `json:"displayName"` //Human readable name for your route, this field is for you
	FrontendId  *int64                         `json:"frontendId"`  //Route traffic for this frontend
	RouteId     int64                          `json:"routeId"`     //Id of your route
	Rules       []*IPLoadbalancingRouteRule    `json:"rules"`       //List of rules to match to trigger action
	Status      string                         `json:"status"`      //Route status. Routes in "ok" state are ready to operate
	Weight      int64                          `json:"weight"`      //Route priority ([0..255]). 0 if null. Highest priority routes are evaluated first. Only the first matching route will trigger an action
}

func (v IPLoadbalancingHttpRoute) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})

	obj["action"] = []interface{}{v.Action.ToMap()}

	if v.DisplayName != nil {
		obj["display_name"] = *v.DisplayName
	}

	if v.FrontendId != nil {
		obj["frontend_id"] = *v.FrontendId
	}

	obj["route_id"] = v.RouteId

	var rules []map[string]interface{}
	for _, r := range v.Rules {
		rules = append(rules, r.ToMapForRoutes())
	}
	obj["rules"] = rules

	obj["status"] = v.Status
	obj["weight"] = v.Weight

	return obj
}

// IPLoadbalancingTcpRouteAction Action triggered when all rules match
type IPLoadbalancingTcpRouteActionOpts struct {
	Target *string `json:"target,omitempty"` // Farm ID for "farm" action type, empty for others
	Type   string  `json:"type"`             // Action to trigger if all the rules of this route matches
}

func (opts *IPLoadbalancingTcpRouteActionOpts) FromResource(d *schema.ResourceData, parent string) *IPLoadbalancingTcpRouteActionOpts {
	opts.Target = helpers.GetNilStringPointerFromData(d, fmt.Sprintf("%s.target", parent))
	opts.Type = d.Get(fmt.Sprintf("%s.type", parent)).(string)
	return opts
}

//IPLoadbalancingTcpRoute HTTP Route
type IPLoadbalancingTcpRouteOpts struct {
	Action      IPLoadbalancingTcpRouteActionOpts `json:"action"`                //Action triggered when all rules match
	DisplayName *string                           `json:"displayName,omitempty"` //Human readable name for your route, this field is for you
	FrontendId  *int64                            `json:"frontendId,omitempty"`  //Route traffic for this frontend
	Weight      *int64                            `json:"weight,omitempty"`      //Route priority ([0..255]). 0 if null. Highest priority routes are evaluated first. Only the first matching route will trigger an action
}

func (opts *IPLoadbalancingTcpRouteOpts) FromResource(d *schema.ResourceData) *IPLoadbalancingTcpRouteOpts {
	opts.DisplayName = helpers.GetNilStringPointerFromData(d, "display_name")

	actions := d.Get("action").([]interface{})
	if actions != nil && len(actions) == 1 {
		opts.Action = *((&IPLoadbalancingTcpRouteActionOpts{}).FromResource(d, "action.0"))
	}

	opts.FrontendId = helpers.GetNilInt64PointerFromData(d, "frontend_id")
	opts.Weight = helpers.GetNilInt64PointerFromData(d, "weight")

	return opts
}

// IPLoadbalancingTcpRouteAction Action triggered when all rules match
type IPLoadbalancingTcpRouteAction struct {
	Target *string `json:"target"` // Farm ID for "farm" action type or URL template for "redirect" action. You may use ${uri}, ${protocol}, ${host}, ${port} and ${path} variables in redirect target
	Type   string  `json:"type"`   // Action to trigger if all the rules of this route matches
}

func (v IPLoadbalancingTcpRouteAction) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})

	if v.Target != nil {
		obj["target"] = *v.Target
	}

	obj["type"] = v.Type

	return obj
}

//IPLoadbalancingTcpRoute HTTP Route
type IPLoadbalancingTcpRoute struct {
	Action      IPLoadbalancingTcpRouteAction `json:"action"`      //Action triggered when all rules match
	DisplayName *string                       `json:"displayName"` //Human readable name for your route, this field is for you
	FrontendId  *int64                        `json:"frontendId"`  //Route traffic for this frontend
	RouteId     int64                         `json:"routeId"`     //Id of your route
	Rules       []*IPLoadbalancingRouteRule   `json:"rules"`       //List of rules to match to trigger action
	Status      string                        `json:"status"`      //Route status. Routes in "ok" state are ready to operate
	Weight      int64                         `json:"weight"`      //Route priority ([0..255]). 0 if null. Highest priority routes are evaluated first. Only the first matching route will trigger an action
}

func (v IPLoadbalancingTcpRoute) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})

	obj["action"] = []interface{}{v.Action.ToMap()}

	if v.DisplayName != nil {
		obj["display_name"] = *v.DisplayName
	}

	if v.FrontendId != nil {
		obj["frontend_id"] = *v.FrontendId
	}

	obj["route_id"] = v.RouteId

	var rules []map[string]interface{}
	for _, r := range v.Rules {
		rules = append(rules, r.ToMapForRoutes())
	}
	obj["rules"] = rules

	obj["status"] = v.Status
	obj["weight"] = v.Weight

	return obj
}

//IPLoadbalancingRouteRule Route Rule
type IPLoadbalancingRouteRule struct {
	DisplayName *string `json:"displayName"` //Human readable name for your rule
	Field       string  `json:"field"`       //Name of the field to match like "protocol" or "host". See "/ipLoadbalancing/{serviceName}/availableRouteRules" for a list of available rules
	Match       string  `json:"match"`       //Matching operator. Not all operators are available for all fields. See "/ipLoadbalancing/{serviceName}/availableRouteRules"
	Negate      bool    `json:"negate"`      //Invert the matching operator effect
	Pattern     *string `json:"pattern"`     //Value to match against this match. Interpretation if this field depends on the match and field
	RuleId      int64   `json:"ruleId"`      //Id of your rule
	SubField    *string `json:"subField"`    //Name of sub-field, if applicable. This may be a Cookie or Header name for instance
}

func (v IPLoadbalancingRouteRule) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})

	if v.DisplayName != nil {
		obj["display_name"] = *v.DisplayName
	}

	obj["field"] = v.Field
	obj["match"] = v.Match
	obj["negate"] = v.Negate

	if v.Pattern != nil {
		obj["pattern"] = *v.Pattern
	}

	obj["rule_id"] = v.RuleId
	obj["sub_field"] = v.SubField

	return obj
}

func (v IPLoadbalancingRouteRule) ToMapForRoutes() map[string]interface{} {
	obj := make(map[string]interface{})

	obj["field"] = v.Field
	obj["match"] = v.Match
	obj["negate"] = v.Negate

	if v.Pattern != nil {
		obj["pattern"] = *v.Pattern
	}

	obj["rule_id"] = v.RuleId
	obj["sub_field"] = v.SubField

	return obj
}

//IPLoadbalancingRouteRule Route Rule
type IPLoadbalancingRouteRuleOpts struct {
	DisplayName *string `json:"displayName,omitempty"` //Human readable name for your rule
	Field       string  `json:"field"`                 //Name of the field to match like "protocol" or "host". See "/ipLoadbalancing/{serviceName}/availableRouteRules" for a list of available rules
	Match       string  `json:"match"`                 //Matching operator. Not all operators are available for all fields. See "/ipLoadbalancing/{serviceName}/availableRouteRules"
	Negate      *bool   `json:"negate,omitempty"`      //Invert the matching operator effect
	Pattern     *string `json:"pattern,omitempty"`     //Value to match against this match. Interpretation if this field depends on the match and field
	SubField    *string `json:"subField,omitempty"`    //Name of sub-field, if applicable. This may be a Cookie or Header name for instance
}

func (opts *IPLoadbalancingRouteRuleOpts) FromResource(d *schema.ResourceData) *IPLoadbalancingRouteRuleOpts {
	opts.DisplayName = helpers.GetNilStringPointerFromData(d, "display_name")
	opts.Field = d.Get("field").(string)
	opts.Match = d.Get("match").(string)
	opts.Negate = helpers.GetNilBoolPointerFromData(d, "negate")
	opts.Pattern = helpers.GetNilStringPointerFromData(d, "pattern")
	opts.SubField = helpers.GetNilStringPointerFromData(d, "sub_field")

	return opts
}

type IpLoadbalancingTcpFrontend struct {
	FrontendId    int      `json:"frontendId,omitempty"`
	Port          string   `json:"port"`
	Zone          string   `json:"zone"`
	AllowedSource []string `json:"allowedSource"`
	DedicatedIpFo []string `json:"dedicatedIpfo"`
	DefaultFarmId *int     `json:"defaultFarmId,omitempty"`
	DefaultSslId  *int     `json:"defaultSslId,omitempty"`
	Disabled      bool     `json:"disabled"`
	Ssl           bool     `json:"ssl"`
	DisplayName   string   `json:"displayName"`
}

type IpLoadbalancingHttpFrontend struct {
	FrontendId       int      `json:"frontendId,omitempty"`
	Port             string   `json:"port"`
	Zone             string   `json:"zone"`
	AllowedSource    []string `json:"allowedSource"`
	DedicatedIpFo    []string `json:"dedicatedIpfo"`
	DefaultFarmId    *int     `json:"defaultFarmId,omitempty"`
	DefaultSslId     *int     `json:"defaultSslId,omitempty"`
	Disabled         bool     `json:"disabled"`
	Ssl              bool     `json:"ssl"`
	RedirectLocation string   `json:"redirectLocation,omitempty"`
	DisplayName      string   `json:"displayName,omitempty"`
}

type IpLoadbalancingFarmServerCreateOpts struct {
	Address              string  `json:"address"`
	Backup               *bool   `json:"backup"`
	Chain                *string `json:"chain,omitempty"`
	Cookie               *string `json:"cookie,omitempty"`
	DisplayName          *string `json:"displayName,omitempty"`
	Port                 *int    `json:"port,omitempty"`
	Probe                *bool   `json:"probe"`
	ProxyProtocolVersion *string `json:"proxyProtocolVersion,omitempty"`
	Ssl                  *bool   `json:"ssl"`
	Status               string  `json:"status"`
	Weight               *int    `json:"weight,omitempty"`
}

type IpLoadbalancingFarmServerUpdateOpts struct {
	Address              *string `json:"address"`
	Backup               *bool   `json:"backup"`
	Chain                *string `json:"chain"`
	Cookie               *string `json:"cookie,omitempty"`
	DisplayName          *string `json:"displayName"`
	Port                 *int    `json:"port,omitempty"`
	Probe                *bool   `json:"probe"`
	ProxyProtocolVersion *string `json:"proxyProtocolVersion"`
	Ssl                  *bool   `json:"ssl"`
	Status               *string `json:"status"`
	Weight               *int    `json:"weight,omitempty"`
}

type IpLoadbalancingFarmServer struct {
	Address              string  `json:"address"`
	Backup               *bool   `json:"backup"`
	Chain                *string `json:"chain"`
	Cookie               *string `json:"cookie"`
	DisplayName          *string `json:"displayName"`
	FarmId               int     `json:"backendId"`
	Port                 *int    `json:"port"`
	Probe                *bool   `json:"probe"`
	ProxyProtocolVersion *string `json:"proxyProtocolVersion"`
	ServerId             int     `json:"serverId"`
	Ssl                  *bool   `json:"ssl"`
	Status               string  `json:"status"`
	Weight               *int    `json:"weight"`
}

func (v IpLoadbalancingFarmServer) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})

	obj["address"] = v.Address
	obj["backup"] = v.Backup
	obj["farm_id"] = v.FarmId
	obj["probe"] = v.Probe
	obj["ssl"] = v.Ssl
	obj["status"] = v.Status

	if v.Chain != nil {
		obj["chain"] = *v.Chain
	}

	if v.Cookie != nil {
		obj["cookie"] = *v.Cookie
	}

	if v.DisplayName != nil {
		obj["display_name"] = *v.DisplayName
	}

	if v.Port != nil {
		obj["port"] = *v.Port
	}

	if v.ProxyProtocolVersion != nil {
		obj["proxy_protocol_version"] = *v.ProxyProtocolVersion
	}

	if v.Weight != nil {
		obj["weight"] = *v.Weight
	}

	return obj
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
	Subnet         string  `json:"subnet"`
	Vlan           int64   `json:"vlan"`
	VrackNetworkId int64   `json:"vrackNetworkId"`
	DisplayName    *string `json:"displayName,omitempty"`
	NatIp          string  `json:"natIp"`
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

	return obj
}

type IpLoadbalancingVrackNetworkCreateOpts struct {
	Subnet      string  `json:"subnet"`
	Vlan        *int64  `json:"vlan,omitempty"`
	DisplayName *string `json:"displayName,omitempty"`
	NatIp       string  `json:"natIp"`
}

func (opts *IpLoadbalancingVrackNetworkCreateOpts) FromResource(d *schema.ResourceData) *IpLoadbalancingVrackNetworkCreateOpts {
	opts.Subnet = d.Get("subnet").(string)
	opts.NatIp = d.Get("nat_ip").(string)
	opts.DisplayName = helpers.GetNilStringPointerFromData(d, "display_name")
	opts.Vlan = helpers.GetNilInt64PointerFromData(d, "vlan")

	return opts
}

type IpLoadbalancingVrackNetworkUpdateOpts struct {
	Vlan        *int64  `json:"vlan,omitempty"`
	DisplayName *string `json:"displayName,omitempty"`
	NatIp       *string `json:"natIp,omitempty"`
}

func (opts *IpLoadbalancingVrackNetworkUpdateOpts) FromResource(d *schema.ResourceData) *IpLoadbalancingVrackNetworkUpdateOpts {
	opts.NatIp = helpers.GetNilStringPointerFromData(d, "nat_ip")
	opts.DisplayName = helpers.GetNilStringPointerFromData(d, "display_name")
	opts.Vlan = helpers.GetNilInt64PointerFromData(d, "vlan")

	return opts
}
