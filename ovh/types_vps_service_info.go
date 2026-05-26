package ovh

// VPSServiceInfoRenew represents the writable nested `renew`
// object returned/accepted by /vps/{serviceName}/serviceInfos.
type VPSServiceInfoRenew struct {
	Automatic          bool  `json:"automatic"`
	DeleteAtExpiration bool  `json:"deleteAtExpiration"`
	Forced             bool  `json:"forced"`
	ManualPayment      *bool `json:"manualPayment,omitempty"`
	Period             *int  `json:"period,omitempty"`
}

// VPSServiceInfo mirrors services.Service returned by
// /vps/{serviceName}/serviceInfos.
type VPSServiceInfo struct {
	ServiceID             int64               `json:"serviceId"`
	Status                string              `json:"status"`
	Creation              string              `json:"creation"`
	Expiration            string              `json:"expiration"`
	EngagedUpTo           *string             `json:"engagedUpTo,omitempty"`
	RenewalType           string              `json:"renewalType"`
	ContactAdmin          string              `json:"contactAdmin"`
	ContactBilling        string              `json:"contactBilling"`
	ContactTech           string              `json:"contactTech"`
	Domain                string              `json:"domain"`
	CanDeleteAtExpiration bool                `json:"canDeleteAtExpiration"`
	PossibleRenewPeriod   []int               `json:"possibleRenewPeriod"`
	Renew                 VPSServiceInfoRenew `json:"renew"`
}

// VPSChangeContactOpts is the body of POST /vps/{serviceName}/changeContact.
type VPSChangeContactOpts struct {
	ContactAdmin   *string `json:"contactAdmin,omitempty"`
	ContactBilling *string `json:"contactBilling,omitempty"`
	ContactTech    *string `json:"contactTech,omitempty"`
}
