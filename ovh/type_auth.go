package ovh

import (
	"time"

	"github.com/ovh/go-ovh/ovh"
)

type OvhAuthCurrentCredential struct {
	OvhSupport    bool             `json:"ovhSupport"`
	Status        string           `json:"status"`
	ApplicationId int64            `json:"applicationId"`
	CredentialId  int64            `json:"credentialId"`
	Rules         []ovh.AccessRule `json:"rules"`
	Expiration    time.Time        `json:"expiration"`
	LastUse       time.Time        `json:"lastUse"`
	Creation      time.Time        `json:"creation"`
}

type OvhAuthDetails struct {
	Account       string           `json:"account"`
	AllowedRoutes []ovh.AccessRule `json:"allowedRoutes"`
	Description   string           `json:"description"`
	Identities    []string         `json:"identities"`
	Method        string           `json:"method"`
	Roles         []string         `json:"roles"`
	User          string           `json:"user"`
}
