package ovh

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"

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
	if i == nil {
		return nil, nil
	}

	// We need to initialize slices to empty slice instead of nil
	template := CloudProjectKubeNodePoolTemplate{
		Metadata: CloudProjectKubeNodePoolTemplateMetadata{
			Annotations: make(map[string]string),
			Finalizers:  make([]string, 0),
			Labels:      make(map[string]string),
		},
		Spec: CloudProjectKubeNodePoolTemplateSpec{
			Taints:        make([]Taint, 0),
			Unschedulable: false,
		},
	}

	templateSet := i.(*schema.Set).List()
	if templateSet == nil || len(templateSet) == 0 {
		return &template, nil
	}

	// Due to this bug https://github.com/hashicorp/terraform-plugin-sdk/pull/1042
	// when updating the 'template' object there is two objects, one is empty, take the not empty one
	templateObject := templateSet[0].(map[string]interface{}) // by default take the first one
	for _, to := range templateSet {
		var (
			annotations map[string]interface{}
			labels      map[string]interface{}
			finalizers  []interface{}
			taints      []interface{}
		)

		metadataSet := to.(map[string]interface{})["metadata"].(*schema.Set).List()
		if len(metadataSet) > 0 {
			metadata := metadataSet[0].(map[string]interface{})
			annotations = metadata["annotations"].(map[string]interface{})
			labels = metadata["labels"].(map[string]interface{})
			finalizers = metadata["finalizers"].([]interface{})
		}

		specSet := to.(map[string]interface{})["spec"].(*schema.Set).List()
		if len(specSet) > 0 {
			spec := specSet[0].(map[string]interface{})
			taints = spec["taints"].([]interface{})
		}

		if len(annotations) == 0 && len(labels) == 0 && len(finalizers) == 0 && len(taints) == 0 {
			continue
		}

		// We found the not empty object
		templateObject = to.(map[string]interface{})
	}

	// metadata
	{
		metadataSet := templateObject["metadata"].(*schema.Set).List()
		if len(metadataSet) > 0 {
			metadata := metadataSet[0].(map[string]interface{})

			// metadata.annotations
			annotations := metadata["annotations"].(map[string]interface{})
			if len(annotations) > 0 {
				for k, v := range annotations {
					template.Metadata.Annotations[k] = v.(string)
				}
			}

			// metadata.finalizers
			finalizers := metadata["finalizers"].([]interface{})
			if len(finalizers) > 0 {
				for _, finalizer := range finalizers {
					template.Metadata.Finalizers = append(template.Metadata.Finalizers, finalizer.(string))
				}
			}

			// metadata.labels
			labels := metadata["labels"].(map[string]interface{})
			if len(labels) > 0 {
				for k, v := range labels {
					template.Metadata.Labels[k] = v.(string)
				}
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

	emptyTemplateResponse := &CloudProjectKubeNodePoolTemplate{
		Metadata: CloudProjectKubeNodePoolTemplateMetadata{
			Annotations: make(map[string]string),
			Finalizers:  make([]string, 0),
			Labels:      make(map[string]string),
		},
		Spec: CloudProjectKubeNodePoolTemplateSpec{
			Taints:        make([]Taint, 0),
			Unschedulable: false,
		},
	}

	// If the template is not nil and not empty, then we need to add it to the map
	if v.Template != nil && !reflect.DeepEqual(v.Template, emptyTemplateResponse) {
		obj["template"] = []map[string]interface{}{{}}

		// template.metadata
		{
			data := map[string]interface{}{
				"finalizers":  v.Template.Metadata.Finalizers,
				"labels":      v.Template.Metadata.Labels,
				"annotations": v.Template.Metadata.Annotations,
			}

			obj["template"].([]map[string]interface{})[0]["metadata"] = []map[string]interface{}{data}
		}

		// template.spec
		{
			data := map[string]interface{}{
				"unschedulable": v.Template.Spec.Unschedulable,
			}

			var taints []map[string]interface{}
			for _, taint := range v.Template.Spec.Taints {
				t := map[string]interface{}{
					"effect": taint.Effect.String(),
					"key":    taint.Key,
					"value":  taint.Value,
				}

				taints = append(taints, t)
			}
			data["taints"] = taints

			obj["template"].([]map[string]interface{})[0]["spec"] = []map[string]interface{}{data}
		}
	}

	return obj
}

func (n *CloudProjectKubeNodePoolResponse) String() string {
	return fmt.Sprintf("%s(%s): %s", n.Name, n.Id, n.Status)
}
