package ovh

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
)

const (
	NotATaint TaintEffectType = iota
	NoExecute
	NoSchedule
	PreferNoSchedule
)

type CloudProjectKubeNodePoolCreateOpts struct {
	AntiAffinity      *bool                                      `json:"antiAffinity,omitempty"`
	AttachFloatingIps *CloudProjectKubeNodepoolAttachFloatingIps `json:"attachFloatingIps,omitempty"`
	Autoscale         *bool                                      `json:"autoscale,omitempty"`
	AvailabilityZones *[]string                                  `json:"availabilityZones,omitempty"`
	DesiredNodes      *int                                       `json:"desiredNodes,omitempty"`
	FlavorName        string                                     `json:"flavorName"`
	MaxNodes          *int                                       `json:"maxNodes,omitempty"`
	MinNodes          *int                                       `json:"minNodes,omitempty"`
	MonthlyBilled     *bool                                      `json:"monthlyBilled,omitempty"`
	Name              *string                                    `json:"name,omitempty"`
	Autoscaling       *CloudProjectKubeNodePoolAutoscaling       `json:"autoscaling,omitempty"`
	Template          *CloudProjectKubeNodePoolTemplate          `json:"template,omitempty"`
}

type TaintEffectType int

type Taint struct {
	Effect TaintEffectType `json:"effect"`
	Key    string          `json:"key"`
	Value  string          `json:"value"`
}

type CloudProjectKubeNodepoolAttachFloatingIps struct {
	Enabled bool `json:"enabled"`
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

type CloudProjectKubeNodePoolAutoscaling struct {
	ScaleDownUtilizationThreshold *float64 `json:"scaleDownUtilizationThreshold,omitempty"`
	ScaleDownUnneededTimeSeconds  *int     `json:"scaleDownUnneededTimeSeconds,omitempty"`
	ScaleDownUnreadyTimeSeconds   *int     `json:"scaleDownUnreadyTimeSeconds,omitempty"`
}

type CloudProjectKubeNodePoolUpdateOpts struct {
	AttachFloatingIps *CloudProjectKubeNodepoolAttachFloatingIps `json:"attachFloatingIps,omitempty"`
	Autoscale         *bool                                      `json:"autoscale,omitempty"`
	DesiredNodes      *int                                       `json:"desiredNodes,omitempty"`
	MaxNodes          *int                                       `json:"maxNodes,omitempty"`
	MinNodes          *int                                       `json:"minNodes,omitempty"`
	Autoscaling       *CloudProjectKubeNodePoolAutoscaling       `json:"autoscaling,omitempty"`
	Template          *CloudProjectKubeNodePoolTemplate          `json:"template,omitempty"`
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

func GetAutoscalingOpts(d *schema.ResourceData) (*CloudProjectKubeNodePoolAutoscaling, error) {
	var autoscaling CloudProjectKubeNodePoolAutoscaling
	var e error
	autoscaling.ScaleDownUtilizationThreshold, e = helpers.GetNilFloat64PointerFromData(d, kubeNodePoolAutoscalingScaleDownUtilizationThresholdKey)
	autoscaling.ScaleDownUnneededTimeSeconds = helpers.GetNilIntPointerFromData(d, kubeNodePoolAutoscalingScaleDownUnneededTimeSecondsKey)
	autoscaling.ScaleDownUnreadyTimeSeconds = helpers.GetNilIntPointerFromData(d, kubeNodePoolAutoscalingScaleDownUnreadyTimeSecondsKey)
	return &autoscaling, e
}

func (opts *CloudProjectKubeNodePoolCreateOpts) FromResource(d *schema.ResourceData) (*CloudProjectKubeNodePoolCreateOpts, error) {
	opts.Autoscale = helpers.GetNilBoolPointerFromData(d, kubeNodePoolAutoscaleKey)
	opts.AntiAffinity = helpers.GetNilBoolPointerFromData(d, kubeNodePoolAntiAffinityKey)
	opts.DesiredNodes = helpers.GetNilIntPointerFromDataAndNilIfNotPresent(d, kubeNodePoolDesiredNodesKey)
	opts.FlavorName = d.Get(kubeNodePoolFlavorNameKey).(string)
	opts.MaxNodes = helpers.GetNilIntPointerFromDataAndNilIfNotPresent(d, kubeNodePoolMaxNodesKey)
	opts.MinNodes = helpers.GetNilIntPointerFromDataAndNilIfNotPresent(d, kubeNodePoolMinNodesKey)
	opts.MonthlyBilled = helpers.GetNilBoolPointerFromData(d, kubeNodePoolMonthlyBilledKey)
	opts.Name = helpers.GetNilStringPointerFromData(d, kubeNameKey)

	template, err := loadNodelPoolTemplateFromResource(d.Get(kubeNodePoolTemplateKey))

	if err != nil {
		return nil, err
	}

	autoscaling, err := GetAutoscalingOpts(d)

	if err != nil {
		return nil, err
	}

	availabilityZones, err := helpers.StringsFromSchema(d, kubeNodePoolAvailabilityZonesKey)
	if err != nil {
		return nil, err
	}
	opts.AvailabilityZones = &availabilityZones

	attachFloatingIps, err := loadNodePoolAttachFloatingIpsFromResource(d.Get(kubeNodePoolAttachFloatingIpsKey))
	if err != nil {
		return nil, err
	}
	opts.AttachFloatingIps = attachFloatingIps

	opts.Autoscaling = autoscaling
	opts.Template = template

	return opts, nil
}

func loadNodePoolAttachFloatingIpsFromResource(i interface{}) (*CloudProjectKubeNodepoolAttachFloatingIps, error) {
	if i == nil {
		return nil, nil
	}

	attachFloatingIpsSet := i.(*schema.Set).List()
	if len(attachFloatingIpsSet) == 0 {
		return nil, nil
	}

	// Due to this bug https://github.com/hashicorp/terraform-plugin-sdk/pull/1042
	// when updating the 'template' object there is two objects, one is empty, take the not empty one
	attachFIPsObject := attachFloatingIpsSet[0].(map[string]interface{}) // by default take the first one
	for _, to := range attachFloatingIpsSet {
		empty := true

		object := to.(map[string]interface{})
		if object[kubeClusterCustomizationEnabledKey].(bool) == true {
			empty = false
		}

		if empty {
			continue
		}

		// We found the not empty object
		attachFIPsObject = object
	}

	return &CloudProjectKubeNodepoolAttachFloatingIps{
		Enabled: attachFIPsObject[kubeClusterCustomizationEnabledKey].(bool),
	}, nil
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
	if len(templateSet) == 0 {
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

		metadataSet := to.(map[string]interface{})[kubeNodePoolTemplateMetadataKey].(*schema.Set).List()
		if len(metadataSet) > 0 {
			metadata := metadataSet[0].(map[string]interface{})
			annotations = metadata[kubeNodePoolTemplateAnnotationsKey].(map[string]interface{})
			labels = metadata[kubeNodePoolTemplateLabelsKey].(map[string]interface{})
			finalizers = metadata[kubeNodePoolTemplateFinalizersKey].([]interface{})
		}

		specSet := to.(map[string]interface{})[kubeNodePoolTemplateSpecKey].(*schema.Set).List()
		if len(specSet) > 0 {
			spec := specSet[0].(map[string]interface{})
			taints = spec[kubeNodePoolTemplateTaintsKey].([]interface{})
		}

		if len(annotations) == 0 && len(labels) == 0 && len(finalizers) == 0 && len(taints) == 0 {
			continue
		}

		// We found the not empty object
		templateObject = to.(map[string]interface{})
	}

	// metadata
	{
		metadataSet := templateObject[kubeNodePoolTemplateMetadataKey].(*schema.Set).List()
		if len(metadataSet) > 0 {
			metadata := metadataSet[0].(map[string]interface{})

			// metadata.annotations
			annotations := metadata[kubeNodePoolTemplateAnnotationsKey].(map[string]interface{})
			if len(annotations) > 0 {
				for k, v := range annotations {
					template.Metadata.Annotations[k] = v.(string)
				}
			}

			// metadata.finalizers
			finalizers := metadata[kubeNodePoolTemplateFinalizersKey].([]interface{})
			if len(finalizers) > 0 {
				for _, finalizer := range finalizers {
					template.Metadata.Finalizers = append(template.Metadata.Finalizers, finalizer.(string))
				}
			}

			// metadata.labels
			labels := metadata[kubeNodePoolTemplateLabelsKey].(map[string]interface{})
			if len(labels) > 0 {
				for k, v := range labels {
					template.Metadata.Labels[k] = v.(string)
				}
			}
		}
	}

	// spec
	{
		specSet := templateObject[kubeNodePoolTemplateSpecKey].(*schema.Set).List()
		if len(specSet) > 0 {
			spec := specSet[0].(map[string]interface{})

			// spec.taints
			taints := spec[kubeNodePoolTemplateTaintsKey].([]interface{})
			for _, taint := range taints {
				effectString := taint.(map[string]interface{})[kubeNodePoolTemplateTaintEffectKey].(string)
				effect := TaintEffecTypeToID[effectString]
				if effect == NotATaint {
					return nil, fmt.Errorf("effect: %s is not a allowable taint %#v", effectString, TaintEffecTypeToID)
				}

				taintObject := Taint{
					Effect: effect,
					Key:    taint.(map[string]interface{})[kubeNodePoolTemplateTaintKeyKey].(string),
				}
				if taint.(map[string]interface{})[kubeNodePoolTemplateTaintValueKey] != nil {
					taintObject.Value = taint.(map[string]interface{})[kubeNodePoolTemplateTaintValueKey].(string)
				}

				template.Spec.Taints = append(template.Spec.Taints, taintObject)
			}

			// spec.unschedulable
			template.Spec.Unschedulable = spec[kubeNodePoolTemplateUnschedulableKey].(bool)
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
	opts.Autoscale = helpers.GetNilBoolPointerFromData(d, kubeNodePoolAutoscaleKey)
	opts.DesiredNodes = helpers.GetNilIntPointerFromDataAndNilIfNotPresent(d, kubeNodePoolDesiredNodesKey)
	opts.MaxNodes = helpers.GetNilIntPointerFromDataAndNilIfNotPresent(d, kubeNodePoolMaxNodesKey)
	opts.MinNodes = helpers.GetNilIntPointerFromDataAndNilIfNotPresent(d, kubeNodePoolMinNodesKey)
	var autoscaling CloudProjectKubeNodePoolAutoscaling
	var e error
	autoscaling.ScaleDownUtilizationThreshold, e = helpers.GetNilFloat64PointerFromData(d, kubeNodePoolAutoscalingScaleDownUtilizationThresholdKey)
	autoscaling.ScaleDownUnneededTimeSeconds = helpers.GetNilIntPointerFromData(d, kubeNodePoolAutoscalingScaleDownUnneededTimeSecondsKey)
	autoscaling.ScaleDownUnreadyTimeSeconds = helpers.GetNilIntPointerFromData(d, kubeNodePoolAutoscalingScaleDownUnreadyTimeSecondsKey)
	opts.Autoscaling = &autoscaling

	template, err := loadNodelPoolTemplateFromResource(d.Get(kubeNodePoolTemplateKey))
	if err != nil {
		return nil, err
	} else if e != nil {
		return nil, e
	}
	opts.Template = template

	attachFloatingIPs, err := loadNodePoolAttachFloatingIpsFromResource(d.Get(kubeNodePoolAttachFloatingIpsKey))
	if err != nil {
		return nil, err
	}
	opts.AttachFloatingIps = attachFloatingIPs

	return opts, nil
}

func (s *CloudProjectKubeNodePoolUpdateOpts) String() string {
	return fmt.Sprintf("%d/%d/%d", *s.DesiredNodes, *s.MinNodes, *s.MaxNodes)
}

type CloudProjectKubeNodePoolResponse struct {
	AttachFloatingIPs *CloudProjectKubeNodepoolAttachFloatingIps `json:"attachFloatingIps"`
	Autoscale         bool                                       `json:"autoscale"`
	AntiAffinity      bool                                       `json:"antiAffinity"`
	AvailabilityZones []string                                   `json:"availabilityZones"`
	AvailableNodes    int                                        `json:"availableNodes"`
	CreatedAt         string                                     `json:"createdAt"`
	CurrentNodes      int                                        `json:"currentNodes"`
	DesiredNodes      int                                        `json:"desiredNodes"`
	Flavor            string                                     `json:"flavor"`
	Id                string                                     `json:"id"`
	MaxNodes          int                                        `json:"maxNodes"`
	MinNodes          int                                        `json:"minNodes"`
	MonthlyBilled     bool                                       `json:"monthlyBilled"`
	Name              string                                     `json:"name"`
	ProjectId         string                                     `json:"projectId"`
	SizeStatus        string                                     `json:"sizeStatus"`
	Status            string                                     `json:"status"`
	UpToDateNodes     int                                        `json:"upToDateNodes"`
	UpdatedAt         string                                     `json:"updatedAt"`
	Autoscaling       CloudProjectKubeNodePoolAutoscaling        `json:"autoscaling"`
	Template          *CloudProjectKubeNodePoolTemplate          `json:"template,omitempty"`
}

func (v CloudProjectKubeNodePoolResponse) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj[kubeNodePoolAntiAffinityKey] = v.AntiAffinity
	obj[kubeNodePoolAutoscaleKey] = v.Autoscale
	obj[kubeNodePoolAvailabilityZonesKey] = v.AvailabilityZones
	obj[kubeNodePoolAvailableNodesKey] = v.AvailableNodes
	obj[kubeCreatedAtKey] = v.CreatedAt
	obj[kubeNodePoolCurrentNodesKey] = v.CurrentNodes
	obj[kubeNodePoolDesiredNodesKey] = v.DesiredNodes
	obj[kubeFlavorKey] = v.Flavor
	obj[kubeNodePoolFlavorNameKey] = v.Flavor
	obj[kubeNodeIdKey] = v.Id
	obj[kubeNodePoolMaxNodesKey] = v.MaxNodes
	obj[kubeNodePoolMinNodesKey] = v.MinNodes
	obj[kubeNodePoolMonthlyBilledKey] = v.MonthlyBilled
	obj[kubeNameKey] = v.Name
	obj[kubeProjectIdKey] = v.ProjectId
	obj[kubeNodePoolSizeStatusKey] = v.SizeStatus
	obj[kubeStatusKey] = v.Status
	obj[kubeNodePoolUpToDateNodesKey] = v.UpToDateNodes
	obj[kubeUpdatedAtKey] = v.UpdatedAt
	obj[kubeNodePoolAutoscalingScaleDownUtilizationThresholdKey] = v.Autoscaling.ScaleDownUtilizationThreshold
	obj[kubeNodePoolAutoscalingScaleDownUnneededTimeSecondsKey] = v.Autoscaling.ScaleDownUnneededTimeSeconds
	obj[kubeNodePoolAutoscalingScaleDownUnreadyTimeSecondsKey] = v.Autoscaling.ScaleDownUnreadyTimeSeconds

	if v.AttachFloatingIPs != nil {
		obj[kubeNodePoolAttachFloatingIpsKey] = []map[string]interface{}{
			{
				kubeClusterCustomizationEnabledKey: v.AttachFloatingIPs.Enabled,
			},
		}
	}

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
		// template.spec
		specData := map[string]interface{}{
			kubeNodePoolTemplateUnschedulableKey: v.Template.Spec.Unschedulable,
		}

		var taints []map[string]interface{}
		for _, taint := range v.Template.Spec.Taints {
			t := map[string]interface{}{
				kubeNodePoolTemplateTaintEffectKey: taint.Effect.String(),
				kubeNodePoolTemplateTaintKeyKey:    taint.Key,
				kubeNodePoolTemplateTaintValueKey:  taint.Value,
			}

			taints = append(taints, t)
		}
		specData[kubeNodePoolTemplateTaintsKey] = taints

		obj[kubeNodePoolTemplateKey] = []map[string]interface{}{
			{
				// template.metadata
				kubeNodePoolTemplateMetadataKey: []map[string]interface{}{
					{
						kubeNodePoolTemplateFinalizersKey:  v.Template.Metadata.Finalizers,
						kubeNodePoolTemplateLabelsKey:      v.Template.Metadata.Labels,
						kubeNodePoolTemplateAnnotationsKey: v.Template.Metadata.Annotations,
					},
				},
				kubeNodePoolTemplateSpecKey: []map[string]interface{}{specData},
			},
		}
	}

	return obj
}

func (n *CloudProjectKubeNodePoolResponse) String() string {
	return fmt.Sprintf("%s(%s): %s", n.Name, n.Id, n.Status)
}
