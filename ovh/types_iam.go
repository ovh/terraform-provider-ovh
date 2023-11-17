package ovh

type IamReferenceAction struct {
	Action       string   `json:"action"`
	Categories   []string `json:"categories"`
	Description  string   `json:"description"`
	ResourceType string   `json:"resourceType"`
}

func (a *IamReferenceAction) ToMap() map[string]any {
	out := make(map[string]any, 4)

	out["action"] = a.Action
	out["categories"] = a.Categories
	out["description"] = a.Description
	out["resource_type"] = a.ResourceType

	return out
}

type IamPolicy struct {
	Id          string         `json:"id,omitempty"`
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	Identities  []string       `json:"identities"`
	Resources   []IamResource  `json:"resources"`
	Permissions IamPermissions `json:"permissions"`
	CreatedAt   string         `json:"createdAt,omitempty"`
	UpdatedAt   string         `json:"updatedAt,omitempty"`
	ReadOnly    bool           `json:"readOnly,omitempty"`
	Owner       string         `json:"owner,omitempty"`
}

func (p IamPolicy) ToMap() map[string]any {
	out := make(map[string]any, 0)
	out["name"] = p.Name

	out["owner"] = p.Owner
	out["created_at"] = p.CreatedAt
	out["identities"] = p.Identities
	var resources []string
	for _, r := range p.Resources {
		resources = append(resources, r.URN)
	}
	out["resources"] = resources

	// inline allow, except and deny
	allow, except, deny := p.Permissions.ToLists()
	if len(allow) != 0 {
		out["allow"] = allow
	}
	if len(except) != 0 {
		out["except"] = except
	}
	if len(deny) != 0 {
		out["deny"] = deny
	}

	if p.Description != "" {
		out["description"] = p.Description
	}
	if p.ReadOnly {
		out["read_only"] = p.ReadOnly
	}
	if p.UpdatedAt != "" {
		out["updated_at"] = p.UpdatedAt
	}

	return out
}

// IamResource represent a possible information returned when viewing a policy
type IamResource struct {
	// URN is always returned and is the urn of the resource or resource group
	// can also be a globing urn, for example "urn:v1:eu:resource:*"
	URN string `json:"urn,omitempty"`
	// If the urn is a resourceGroup. the details are also present
	Group *IamResourceGroup `json:"group,omitempty"`
	// If the urn is an IAM resource, the details are also present
	Resource *IamResourceDetails `json:"resource,omitempty"`
}

type IamResourceDetails struct {
	Id          string            `json:"id,omitempty"`
	URN         string            `json:"urn,omitempty"`
	Name        string            `json:"name,omitempty"`
	DisplayName string            `json:"displayName,omitempty"`
	Owner       string            `json:"owner,omitempty"`
	Type        string            `json:"type,omitempty"`
	Tags        map[string]string `json:"tags,omitempty"`
}

type IamPermissions struct {
	Allow  []IamAction `json:"allow"`
	Except []IamAction `json:"except"`
	Deny   []IamAction `json:"deny"`
}

func (p IamPermissions) ToLists() ([]string, []string, []string) {
	var allow []string
	var except []string
	var deny []string

	for _, r := range p.Allow {
		allow = append(allow, r.Action)
	}

	for _, r := range p.Except {
		except = append(except, r.Action)
	}

	for _, r := range p.Deny {
		deny = append(deny, r.Action)
	}
	return allow, except, deny
}

type IamAction struct {
	Action string `json:"action"`
}

type IamResourceGroup struct {
	ID        string               `json:"id,omitempty"`
	Name      string               `json:"name"`
	Resources []IamResourceDetails `json:"resources"`
	URN       string               `json:"urn,omitempty"`
	Owner     string               `json:"owner,omitempty"`
	CreatedAt string               `json:"createdAt,omitempty"`
	UpdatedAt string               `json:"updatedAt,omitempty"`
	ReadOnly  bool                 `json:"readOnly,omitempty"`
}

type IamResourceGroupCreate struct {
	Name      string               `json:"name"`
	Resources []IamResourceDetails `json:"resources"`
}
