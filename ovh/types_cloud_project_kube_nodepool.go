package ovh

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

const (
	NoExecute TaintEffectType = iota
	NoSchedule
	PreferNoSchedule
)

type CloudProjectKubeNodePoolCreateOpts struct {
	AntiAffinity  *bool                             `json:"antiAffinity,omitempty"`
	Autoscale     *bool                             `json:"autoscale,omitempty"`
	DesiredNodes  *int                              `json:"desiredNodes,omitempty"`
	FlavorName    string                            `json:"flavorName"`
	MaxNodes      *int                              `json:"maxNodes,omitempty"`
	MinNodes      *int                              `json:"minNodes,omitempty"`
	MonthlyBilled *bool                             `json:"monthlyBilled,omitempty"`
	Name          *string                           `json:"name,omitempty"`
	Template      *CloudProjectKubeNodePoolTemplate `json:"template,omitempty"`
}

type TaintEffectType int

type Taint struct {
	Effect TaintEffectType `json:"effect,omitempty"`
	Key    string          `json:"key,omitempty"`
	Value  string          `json:"value,omitempty"`
}

type CloudProjectKubeNodePoolTemplateMetadata struct {
	Annotations map[string]string `json:"annotations"`
	Finalizers  []string          `json:"finalizers"`
	Labels      map[string]string `json:"labels"`
}

type CloudProjectKubeNodePoolTemplateSpec struct {
	Taints        []Taint `json:"taints"`
	Unschedulable bool    `json:"unschedulable"`
}

type CloudProjectKubeNodePoolTemplate struct {
	Metadata *CloudProjectKubeNodePoolTemplateMetadata `json:"metadata,omitempty"`
	Spec     *CloudProjectKubeNodePoolTemplateSpec     `json:"spec,omitempty"`
}

type CloudProjectKubeNodePoolUpdateOpts struct {
	Autoscale    *bool                             `json:"autoscale,omitempty"`
	DesiredNodes *int                              `json:"desiredNodes,omitempty"`
	MaxNodes     *int                              `json:"maxNodes,omitempty"`
	MinNodes     *int                              `json:"minNodes,omitempty"`
	Template     *CloudProjectKubeNodePoolTemplate `json:"template,omitempty"`
}

var toString = map[TaintEffectType]string{
	NoExecute:        "NoExecute",
	NoSchedule:       "NoSchedule",
	PreferNoSchedule: "PreferNoSchedule",
}

var TaintEffecTypeToID = map[string]TaintEffectType{
	"NoExecute":        NoExecute,
	"NoSchedule":       NoSchedule,
	"PreferNoSchedule": PreferNoSchedule,
}

func (opts *CloudProjectKubeNodePoolCreateOpts) FromResource(d *schema.ResourceData) *CloudProjectKubeNodePoolCreateOpts {
	opts.Autoscale = helpers.GetNilBoolPointerFromData(d, "autoscale")
	opts.AntiAffinity = helpers.GetNilBoolPointerFromData(d, "anti_affinity")
	opts.DesiredNodes = helpers.GetNilIntPointerFromData(d, "desired_nodes")
	opts.FlavorName = d.Get("flavor_name").(string)
	opts.MaxNodes = helpers.GetNilIntPointerFromData(d, "max_nodes")
	opts.MinNodes = helpers.GetNilIntPointerFromData(d, "min_nodes")
	opts.MonthlyBilled = helpers.GetNilBoolPointerFromData(d, "monthly_billed")
	opts.Name = helpers.GetNilStringPointerFromData(d, "name")
	opts.Template = loadNodelPoolTemplateFromResource(d.Get("template"))

	return opts
}

func loadNodelPoolTemplateFromResource(i interface{}) *CloudProjectKubeNodePoolTemplate {
	template := CloudProjectKubeNodePoolTemplate{
		Metadata: &CloudProjectKubeNodePoolTemplateMetadata{
			Annotations: map[string]string{},
			Finalizers:  []string{},
			Labels:      map[string]string{},
		},
		Spec: &CloudProjectKubeNodePoolTemplateSpec{
			Taints:        []Taint{},
			Unschedulable: false,
		},
	}

	templateSet := i.(*schema.Set).List()
	for _, templateObject := range templateSet {

		metadataSet := templateObject.(map[string]interface{})["metadata"].(*schema.Set).List()
		for _, meta := range metadataSet {

			annotations := meta.(map[string]interface{})["annotations"].(map[string]interface{})
			template.Metadata.Annotations = make(map[string]string)
			for k, v := range annotations {
				template.Metadata.Annotations[k] = v.(string)
			}

			labels := meta.(map[string]interface{})["labels"].(map[string]interface{})
			template.Metadata.Labels = make(map[string]string)
			for k, v := range labels {
				template.Metadata.Labels[k] = v.(string)
			}

			finalizers := meta.(map[string]interface{})["finalizers"].([]interface{})
			for _, finalizer := range finalizers {
				template.Metadata.Finalizers = append(template.Metadata.Finalizers, finalizer.(string))
			}

		}

		specSet := templateObject.(map[string]interface{})["spec"].(*schema.Set).List()
		for _, spec := range specSet {

			taints := spec.(map[string]interface{})["taints"].([]interface{})
			for _, taint := range taints {
				template.Spec.Taints = append(template.Spec.Taints, Taint{
					Effect: TaintEffecTypeToID[taint.(map[string]interface{})["effect"].(string)],
					Key:    taint.(map[string]interface{})["key"].(string),
					Value:  taint.(map[string]interface{})["value"].(string),
				})
			}

			unschedulable := spec.(map[string]interface{})["unschedulable"].(bool)
			template.Spec.Unschedulable = unschedulable

		}
	}

	return &template
}

func (s *CloudProjectKubeNodePoolCreateOpts) String() string {
	return fmt.Sprintf("%s(%s): %d/%d/%d", *s.Name, s.FlavorName, *s.DesiredNodes, *s.MinNodes, *s.MaxNodes)
}

func (e TaintEffectType) String() string {
	return toString[e]
}

// MarshalJSON marshals the enum as a quoted json string
func (e TaintEffectType) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(toString[e])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (e *TaintEffectType) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	*e = TaintEffecTypeToID[j]
	return nil
}

func (opts *CloudProjectKubeNodePoolUpdateOpts) FromResource(d *schema.ResourceData) *CloudProjectKubeNodePoolUpdateOpts {
	opts.Autoscale = helpers.GetNilBoolPointerFromData(d, "autoscale")
	opts.DesiredNodes = helpers.GetNilIntPointerFromData(d, "desired_nodes")
	opts.MaxNodes = helpers.GetNilIntPointerFromData(d, "max_nodes")
	opts.MinNodes = helpers.GetNilIntPointerFromData(d, "min_nodes")
	opts.Template = loadNodelPoolTemplateFromResource(d.Get("template"))

	return opts
}

func (s *CloudProjectKubeNodePoolUpdateOpts) String() string {
	return fmt.Sprintf("%d/%d/%d", *s.DesiredNodes, *s.MinNodes, *s.MaxNodes)
}

type CloudProjectKubeNodePoolResponse struct {
	Autoscale      bool   `json:"autoscale"`
	AntiAffinity   bool   `json:"antiAffinity"`
	AvailableNodes int    `json:"availableNodes"`
	CreatedAt      string `json:"createdAt"`
	CurrentNodes   int    `json:"currentNodes"`
	DesiredNodes   int    `json:"desiredNodes"`
	Flavor         string `json:"flavor"`
	Id             string `json:"id"`
	MaxNodes       int    `json:"maxNodes"`
	MinNodes       int    `json:"minNodes"`
	MonthlyBilled  bool   `json:"monthlyBilled"`
	Name           string `json:"name"`
	ProjectId      string `json:"projectId"`
	SizeStatus     string `json:"sizeStatus"`
	Status         string `json:"status"`
	UpToDateNodes  int    `json:"upToDateNodes"`
	UpdatedAt      string `json:"updatedAt"`
}

func (v CloudProjectKubeNodePoolResponse) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["anti_affinity"] = v.AntiAffinity
	obj["autoscale"] = v.Autoscale
	obj["available_nodes"] = v.AvailableNodes
	obj["created_at"] = v.CreatedAt
	obj["current_nodes"] = v.CurrentNodes
	obj["desired_nodes"] = v.DesiredNodes
	obj["flavor"] = v.Flavor
	obj["flavor_name"] = v.Flavor
	obj["id"] = v.Id
	obj["max_nodes"] = v.MaxNodes
	obj["min_nodes"] = v.MinNodes
	obj["monthly_billed"] = v.MonthlyBilled
	obj["name"] = v.Name
	obj["project_id"] = v.ProjectId
	obj["size_status"] = v.SizeStatus
	obj["status"] = v.Status
	obj["up_to_date_nodes"] = v.UpToDateNodes
	obj["updated_at"] = v.UpdatedAt

	return obj
}

func (n *CloudProjectKubeNodePoolResponse) String() string {
	return fmt.Sprintf("%s(%s): %s", n.Name, n.Id, n.Status)
}
