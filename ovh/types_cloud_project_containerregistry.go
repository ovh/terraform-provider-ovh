package ovh

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

type CloudProjectCapabilitiesContainerRegistry struct {
	RegionName string                                          `json:"regionName"`
	Plans      []CloudProjectCapabilitiesContainerRegistryPlan `json:"plans"`
}

func (v CloudProjectCapabilitiesContainerRegistry) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})

	obj["region_name"] = v.RegionName

	if v.Plans != nil {
		plans := make([]interface{}, len(v.Plans))
		for i, plan := range v.Plans {
			plans[i] = plan.ToMap()
		}
		obj["plans"] = plans
	}

	return obj
}

type CloudProjectCapabilitiesContainerRegistryPlan struct {
	Code           string                                                      `json:"code"`
	CreatedAt      string                                                      `json:"createdAt"`
	Features       CloudProjectCapabilitiesContainerRegistryPlanFeatures       `json:"features"`
	Id             string                                                      `json:"id"`
	Name           string                                                      `json:"name"`
	RegistryLimits CloudProjectCapabilitiesContainerRegistryPlanRegistryLimits `json:"registryLimits"`
	UpdatedAt      string                                                      `json:"updatedAt"`
}

func (v CloudProjectCapabilitiesContainerRegistryPlan) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})

	obj["code"] = v.Code
	obj["created_at"] = v.CreatedAt
	obj["id"] = v.Id
	obj["features"] = []interface{}{v.Features.ToMap()}
	obj["name"] = v.Name
	obj["updated_at"] = v.UpdatedAt
	obj["registry_limits"] = []interface{}{v.RegistryLimits.ToMap()}
	return obj
}

type CloudProjectCapabilitiesContainerRegistryPlanFeatures struct {
	Vulnerability bool `json:"vulnerability"`
}

func (v CloudProjectCapabilitiesContainerRegistryPlanFeatures) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})

	obj["vulnerability"] = v.Vulnerability
	return obj
}

type CloudProjectCapabilitiesContainerRegistryPlanRegistryLimits struct {
	ImageStorage    int64 `json:"imageStorage"`
	ParallelRequest int64 `json:"parallelRequest"`
}

func (v CloudProjectCapabilitiesContainerRegistryPlanRegistryLimits) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})

	obj["image_storage"] = v.ImageStorage
	obj["parallel_request"] = v.ParallelRequest
	return obj
}

type CloudProjectContainerRegistry struct {
	CreatedAt string `json:"createdAt"`
	Id        string `json:"id"`
	Name      string `json:"name"`
	ProjectID string `json:"projectID"`
	Region    string `json:"region"`
	Size      int64  `json:"size"`
	Status    string `json:"status"`
	UpdatedAt string `json:"updatedAt"`
	Url       string `json:"url"`
	Version   string `json:"version"`
}

func (p *CloudProjectContainerRegistry) String() string {
	return fmt.Sprintf(
		"Id: %s, Status: %s, Name: %s, ProjectID: %s",
		p.Id,
		p.Status,
		p.Name,
		p.ProjectID,
	)
}

func (r CloudProjectContainerRegistry) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["created_at"] = r.CreatedAt
	obj["id"] = r.Id
	obj["name"] = r.Name
	obj["project_id"] = r.ProjectID
	obj["region"] = r.Region
	obj["size"] = r.Size
	obj["status"] = r.Status
	obj["updated_at"] = r.UpdatedAt
	obj["url"] = r.Url
	obj["version"] = r.Version
	return obj
}

type CloudProjectContainerRegistryCreateOpts struct {
	Name   string  `json:"name"`
	Region string  `json:"region"`
	PlanId *string `json:"planID,omitempty"`
}

func (opts *CloudProjectContainerRegistryCreateOpts) FromResource(d *schema.ResourceData) *CloudProjectContainerRegistryCreateOpts {
	opts.Name = d.Get("name").(string)
	opts.Region = d.Get("region").(string)
	opts.PlanId = helpers.GetNilStringPointerFromData(d, "plan_id")
	return opts
}

type CloudProjectContainerRegistryUpdateOpts struct {
	Name string `json:"name"`
}

func (opts *CloudProjectContainerRegistryUpdateOpts) FromResource(d *schema.ResourceData) *CloudProjectContainerRegistryUpdateOpts {
	opts.Name = d.Get("name").(string)
	return opts
}

type CloudProjectContainerRegistryPlanUpdateOpts struct {
	PlanId string `json:"planID"`
}

func (opts *CloudProjectContainerRegistryPlanUpdateOpts) FromResource(d *schema.ResourceData) *CloudProjectContainerRegistryPlanUpdateOpts {
	opts.PlanId = d.Get("plan_id").(string)
	return opts
}

type CloudProjectContainerRegistryUser struct {
	Email    string `json:"email"`
	Id       string `json:"id"`
	Password string `json:"password"`
	User     string `json:"user"`
}

func (p *CloudProjectContainerRegistryUser) String() string {
	return fmt.Sprintf(
		"Id: %s, User: %s, Email: %s",
		p.Id,
		p.User,
		p.Email,
	)
}

func (r CloudProjectContainerRegistryUser) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["email"] = r.Email
	obj["id"] = r.Id
	obj["password"] = r.Password
	obj["user"] = r.User
	return obj
}

func (r CloudProjectContainerRegistryUser) ToMapWithKeys(keys []string) map[string]interface{} {
	obj := make(map[string]interface{})
	all := r.ToMap()
	for _, key := range keys {
		obj[key] = all[key]
	}
	return obj
}

type CloudProjectContainerRegistryUserCreateOpts struct {
	Email string `json:"email"`
	Login string `json:"login"`
}

func (opts *CloudProjectContainerRegistryUserCreateOpts) FromResource(d *schema.ResourceData) *CloudProjectContainerRegistryUserCreateOpts {
	opts.Email = d.Get("email").(string)
	opts.Login = d.Get("login").(string)
	return opts
}
