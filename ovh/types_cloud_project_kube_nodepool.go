package ovh

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

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
	if len(templateSet) == 0 {
		return nil, nil
	}
	templateObject := templateSet[0]

	// when updating the nested object template there is two objects, one is empty, take the not empty one
	if len(templateSet) > 1 {
		for _, to := range templateSet {
			metadata := to.(map[string]interface{})["metadata"].(*schema.Set).List()[0]
			annotations := metadata.(map[string]interface{})["annotations"].(map[string]interface{})
			labels := metadata.(map[string]interface{})["labels"].(map[string]interface{})
			finalizers := metadata.(map[string]interface{})["finalizers"].([]interface{})

			spec := templateObject.(map[string]interface{})["spec"].(*schema.Set).List()[0]
			taints := spec.(map[string]interface{})["taints"].([]interface{})
			unschedulable := spec.(map[string]interface{})["unschedulable"].(bool)

			if len(annotations) == 0 && len(labels) == 0 && len(finalizers) == 0 && len(taints) == 0 && unschedulable == false {
				// is empty
			} else {
				templateObject = to
				break
			}
		}
	}
	if len(templateSet) > 2 {
		return nil, errors.New("resource template cannot have more than 2 elements")
	}

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
			effectString := taint.(map[string]interface{})["effect"].(string)
			effect := TaintEffecTypeToID[effectString]
			if effect == NotATaint {
				return nil, errors.New(fmt.Sprintf("Effect: %s is not a allowable taint %#v", effectString, TaintEffecTypeToID))
			}

			template.Spec.Taints = append(template.Spec.Taints, Taint{
				Effect: effect,
				Key:    taint.(map[string]interface{})["key"].(string),
				Value:  taint.(map[string]interface{})["value"].(string),
			})
		}

		template.Spec.Unschedulable = spec.(map[string]interface{})["unschedulable"].(bool)

	}

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

	var taints []map[string]interface{}
	for _, taint := range v.Template.Spec.Taints {
		t := map[string]interface{}{
			"effect": taint.Effect.String(),
			"key":    taint.Key,
			"value":  taint.Value,
		}

		taints = append(taints, t)
	}

	obj["template"] = []map[string]interface{}{
		{
			"metadata": []map[string]interface{}{
				{
					"finalizers":  v.Template.Metadata.Finalizers,
					"labels":      v.Template.Metadata.Labels,
					"annotations": v.Template.Metadata.Annotations,
				},
			},
			"spec": []map[string]interface{}{
				{
					"unschedulable": v.Template.Spec.Unschedulable,
					"taints":        taints,
				},
			},
		},
	}

	if len(obj["template"].([]map[string]interface{})[0]["metadata"].([]map[string]interface{})[0]["finalizers"].([]string)) == 0 {
		obj["template"].([]map[string]interface{})[0]["metadata"].([]map[string]interface{})[0]["finalizers"] = nil
	}
	if len(obj["template"].([]map[string]interface{})[0]["metadata"].([]map[string]interface{})[0]["labels"].(map[string]string)) == 0 {
		obj["template"].([]map[string]interface{})[0]["metadata"].([]map[string]interface{})[0]["labels"] = nil
	}
	if len(obj["template"].([]map[string]interface{})[0]["metadata"].([]map[string]interface{})[0]["annotations"].(map[string]string)) == 0 {
		obj["template"].([]map[string]interface{})[0]["metadata"].([]map[string]interface{})[0]["annotations"] = nil
	}
	if obj["template"].([]map[string]interface{})[0]["metadata"].([]map[string]interface{})[0]["finalizers"] == nil &&
		obj["template"].([]map[string]interface{})[0]["metadata"].([]map[string]interface{})[0]["labels"] == nil &&
		obj["template"].([]map[string]interface{})[0]["metadata"].([]map[string]interface{})[0]["annotations"] == nil {
		obj["template"].([]map[string]interface{})[0]["metadata"] = nil
	}
	if obj["template"].([]map[string]interface{})[0]["spec"].([]map[string]interface{})[0]["unschedulable"].(bool) == false {
		obj["template"].([]map[string]interface{})[0]["spec"].([]map[string]interface{})[0]["unschedulable"] = nil
	}
	if len(obj["template"].([]map[string]interface{})[0]["spec"].([]map[string]interface{})[0]["taints"].([]map[string]interface{})) == 0 {
		obj["template"].([]map[string]interface{})[0]["spec"].([]map[string]interface{})[0]["taints"] = nil
	}

	if obj["template"].([]map[string]interface{})[0]["spec"].([]map[string]interface{})[0]["unschedulable"] == nil &&
		obj["template"].([]map[string]interface{})[0]["spec"].([]map[string]interface{})[0]["taints"] == nil {
		obj["template"].([]map[string]interface{})[0]["spec"] = nil
	}

	if obj["template"].([]map[string]interface{})[0]["metadata"] == nil &&
		obj["template"].([]map[string]interface{})[0]["spec"] == nil {
		obj["template"] = nil
	}

	return obj
}

func (n *CloudProjectKubeNodePoolResponse) String() string {
	return fmt.Sprintf("%s(%s): %s", n.Name, n.Id, n.Status)
}
