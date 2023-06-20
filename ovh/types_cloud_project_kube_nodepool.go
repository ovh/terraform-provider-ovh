package ovh

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

const (
	NotATaint TaintEffectType = iota
	NoExecute
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
	Effect TaintEffectType `json:"effect"`
	Key    string          `json:"key"`
	Value  string          `json:"value"`
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
	Metadata CloudProjectKubeNodePoolTemplateMetadata `json:"metadata"`
	Spec     CloudProjectKubeNodePoolTemplateSpec     `json:"spec"`
}

type CloudProjectKubeNodePoolUpdateOpts struct {
	Autoscale    *bool                             `json:"autoscale,omitempty"`
	DesiredNodes *int                              `json:"desiredNodes,omitempty"`
	MaxNodes     *int                              `json:"maxNodes,omitempty"`
	MinNodes     *int                              `json:"minNodes,omitempty"`
	Template     *CloudProjectKubeNodePoolTemplate `json:"template,omitempty"`
}

var toString = map[TaintEffectType]string{
	NotATaint:        "",
	NoExecute:        "NoExecute",
	NoSchedule:       "NoSchedule",
	PreferNoSchedule: "PreferNoSchedule",
}

var TaintEffecTypeToID = map[string]TaintEffectType{
	"":                 NotATaint,
	"NoExecute":        NoExecute,
	"NoSchedule":       NoSchedule,
	"PreferNoSchedule": PreferNoSchedule,
}

func (opts *CloudProjectKubeNodePoolCreateOpts) FromResource(d *schema.ResourceData) (*CloudProjectKubeNodePoolCreateOpts, error) {
	opts.Autoscale = helpers.GetNilBoolPointerFromData(d, "autoscale")
	opts.AntiAffinity = helpers.GetNilBoolPointerFromData(d, "anti_affinity")
	opts.DesiredNodes = helpers.GetNilIntPointerFromDataAndNilIfNotPresent(d, "desired_nodes")
	opts.FlavorName = d.Get("flavor_name").(string)
	opts.MaxNodes = helpers.GetNilIntPointerFromDataAndNilIfNotPresent(d, "max_nodes")
	opts.MinNodes = helpers.GetNilIntPointerFromDataAndNilIfNotPresent(d, "min_nodes")
	opts.MonthlyBilled = helpers.GetNilBoolPointerFromData(d, "monthly_billed")
	opts.Name = helpers.GetNilStringPointerFromData(d, "name")
	template, err := loadNodelPoolTemplateFromResource(d.Get("template"))
	if err != nil {
		return nil, err
	}
	opts.Template = template

	return opts, nil
}

func loadNodelPoolTemplateFromResource(i interface{}) (*CloudProjectKubeNodePoolTemplate, error) {
	// initialize map variables to explicit empty map
	template := CloudProjectKubeNodePoolTemplate{
		Metadata: CloudProjectKubeNodePoolTemplateMetadata{},
		Spec:     CloudProjectKubeNodePoolTemplateSpec{},
	}

	templateSet := i.(*schema.Set).List()
	if len(templateSet) > 2 {
		return nil, errors.New("resource template cannot have more than 2 elements")
	}

	if len(templateSet) > 0 {
		templateObject := templateSet[0].(map[string]interface{})

		// metadata
		{
			metadataSet := templateObject["metadata"].(*schema.Set).List()
			if len(metadataSet) > 0 {
				metadata := metadataSet[0].(map[string]interface{})

				// metadata.annotations
				annotations := metadata["annotations"].(map[string]interface{})
				template.Metadata.Annotations = make(map[string]string)
				for k, v := range annotations {
					template.Metadata.Annotations[k] = v.(string)
				}

				// metadata.finalizers
				finalizers := metadata["finalizers"].([]interface{})
				template.Metadata.Finalizers = make([]string, 0)
				for _, finalizer := range finalizers {
					template.Metadata.Finalizers = append(template.Metadata.Finalizers, finalizer.(string))
				}

				// metadata.labels
				labels := metadata["labels"].(map[string]interface{})
				template.Metadata.Labels = make(map[string]string)
				for k, v := range labels {
					template.Metadata.Labels[k] = v.(string)
				}
			}
		}

		// spec
		{
			specSet := templateObject["spec"].(*schema.Set).List()
			if len(specSet) > 0 {
				spec := specSet[0].(map[string]interface{})

				// spec.taints
				taints := spec["taints"].([]interface{})
				template.Spec.Taints = make([]Taint, 0)
				for _, taint := range taints {
					effectString := taint.(map[string]interface{})["effect"].(string)
					effect := TaintEffecTypeToID[effectString]
					if effect == NotATaint {
						return nil, fmt.Errorf("effect: %s is not a allowable taint %#v", effectString, TaintEffecTypeToID)
					}

					template.Spec.Taints = append(template.Spec.Taints, Taint{
						Effect: effect,
						Key:    taint.(map[string]interface{})["key"].(string),
						Value:  taint.(map[string]interface{})["value"].(string),
					})
				}

				// spec.unschedulable
				template.Spec.Unschedulable = spec["unschedulable"].(bool)
			}
		}
	}

	log.Printf("[DEBUG] >>>>>>>>>>>>>>>>>>>>>>>>>>%#+v", templateSet)

	return &template, nil
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

func (opts *CloudProjectKubeNodePoolUpdateOpts) FromResource(d *schema.ResourceData) (*CloudProjectKubeNodePoolUpdateOpts, error) {
	opts.Autoscale = helpers.GetNilBoolPointerFromData(d, "autoscale")
	opts.DesiredNodes = helpers.GetNilIntPointerFromDataAndNilIfNotPresent(d, "desired_nodes")
	opts.MaxNodes = helpers.GetNilIntPointerFromDataAndNilIfNotPresent(d, "max_nodes")
	opts.MinNodes = helpers.GetNilIntPointerFromDataAndNilIfNotPresent(d, "min_nodes")
	template, err := loadNodelPoolTemplateFromResource(d.Get("template"))
	if err != nil {
		return nil, err
	}
	opts.Template = template

	return opts, nil
}

func (s *CloudProjectKubeNodePoolUpdateOpts) String() string {
	return fmt.Sprintf("%d/%d/%d", *s.DesiredNodes, *s.MinNodes, *s.MaxNodes)
}

type CloudProjectKubeNodePoolResponse struct {
	Autoscale      bool                              `json:"autoscale"`
	AntiAffinity   bool                              `json:"antiAffinity"`
	AvailableNodes int                               `json:"availableNodes"`
	CreatedAt      string                            `json:"createdAt"`
	CurrentNodes   int                               `json:"currentNodes"`
	DesiredNodes   int                               `json:"desiredNodes"`
	Flavor         string                            `json:"flavor"`
	Id             string                            `json:"id"`
	MaxNodes       int                               `json:"maxNodes"`
	MinNodes       int                               `json:"minNodes"`
	MonthlyBilled  bool                              `json:"monthlyBilled"`
	Name           string                            `json:"name"`
	ProjectId      string                            `json:"projectId"`
	SizeStatus     string                            `json:"sizeStatus"`
	Status         string                            `json:"status"`
	UpToDateNodes  int                               `json:"upToDateNodes"`
	UpdatedAt      string                            `json:"updatedAt"`
	Template       *CloudProjectKubeNodePoolTemplate `json:"template,omitempty"`
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

	if v.Template != nil {
		obj["template"] = []map[string]interface{}{{}}

		// template.metadata
		data := make(map[string]interface{})
		if vv := v.Template.Metadata.Finalizers; vv != nil && len(vv) > 0 {
			data["finalizers"] = vv
		}

		if vv := v.Template.Metadata.Labels; vv != nil && len(vv) > 0 {
			data["labels"] = vv
		}

		if vv := v.Template.Metadata.Annotations; vv != nil && len(vv) > 0 {
			data["annotations"] = vv
		}

		if len(data) > 0 {
			obj["template"].([]map[string]interface{})[0]["metadata"] = []map[string]interface{}{data}
		}

		// template.spec
		data = make(map[string]interface{})
		if vv := v.Template.Spec.Taints; vv != nil && len(vv) > 0 {
			var taints []map[string]interface{}
			for _, taint := range vv {
				t := map[string]interface{}{
					"effect": taint.Effect.String(),
					"key":    taint.Key,
					"value":  taint.Value,
				}

				taints = append(taints, t)
			}

			data["taints"] = taints
		}

		data["unschedulable"] = v.Template.Spec.Unschedulable

		if len(data) > 0 {
			obj["template"].([]map[string]interface{})[0]["spec"] = []map[string]interface{}{data}
		}

	}

	// Delete the entire template if it's empty
	if len(obj["template"].([]map[string]interface{})[0]) == 0 {
		delete(obj, "template")
	}

	return obj
}

func (n *CloudProjectKubeNodePoolResponse) String() string {
	return fmt.Sprintf("%s(%s): %s", n.Name, n.Id, n.Status)
}
