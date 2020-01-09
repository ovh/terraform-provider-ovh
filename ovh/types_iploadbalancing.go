package ovh

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

type IpLoadbalancingTcpFarmBackendProbe struct {
	Match    string `json:"match,omitempty"`
	Port     int    `json:"port,omitempty"`
	Interval int    `json:"interval,omitempty"`
	Negate   bool   `json:"negate,omitempty"`
	Pattern  string `json:"pattern,omitempty"`
	ForceSsl bool   `json:"forceSsl,omitempty"`
	URL      string `json:"url,omitempty"`
	Method   string `json:"method,omitempty"`
	Type     string `json:"type,omitempty"`
}

type IpLoadbalancingTcpFarm struct {
	FarmId         int                                 `json:"farmId,omitempty"`
	Zone           string                              `json:"zone,omitempty"`
	VrackNetworkId int                                 `json:"vrackNetworkId,omitempty"`
	Port           int                                 `json:"port,omitempty"`
	Stickiness     string                              `json:"stickiness,omitempty"`
	Balance        string                              `json:"balance,omitempty"`
	Probe          *IpLoadbalancingTcpFarmBackendProbe `json:"probe,omitempty"`
	DisplayName    string                              `json:"displayName,omitempty"`
}

type IpLoadbalancingHttpFarmBackendProbe struct {
	Match    string `json:"match,omitempty"`
	Port     int    `json:"port,omitempty"`
	Interval int    `json:"interval,omitempty"`
	Negate   bool   `json:"negate,omitempty"`
	Pattern  string `json:"pattern,omitempty"`
	ForceSsl bool   `json:"forceSsl,omitempty"`
	URL      string `json:"url,omitempty"`
	Method   string `json:"method,omitempty"`
	Type     string `json:"type,omitempty"`
}

type IpLoadbalancingHttpFarm struct {
	FarmId         int                                  `json:"farmId,omitempty"`
	Zone           string                               `json:"zone,omitempty"`
	VrackNetworkId int                                  `json:"vrackNetworkId,omitempty"`
	Port           int                                  `json:"port,omitempty"`
	Stickiness     string                               `json:"stickiness,omitempty"`
	Balance        string                               `json:"balance,omitempty"`
	Probe          *IpLoadbalancingHttpFarmBackendProbe `json:"probe,omitempty"`
	DisplayName    string                               `json:"displayName,omitempty"`
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
