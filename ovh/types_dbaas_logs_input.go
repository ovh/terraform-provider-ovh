package ovh

import (
	"fmt"
	"html"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

type DbaasLogsInput struct {
	AllowedNetworks   []string `json:"allowedNetworks"`
	CreatedAt         string   `json:createdAt`
	Description       string   `json:"description"`
	EngineId          string   `json:"engineId"`
	ExposedPort       *string  `json:"exposedPort"`
	Hostname          string   `json:hostname`
	InputId           string   `json:inputId`
	IsRestartRequired bool     `json:isRestartRequired`
	NbInstance        *int64   `json:"nbInstance"`
	PublicAddress     string   `json:publicAddress`
	SslCertificate    string   `json:sslCertificate`
	Status            string   `json:status`
	StreamId          string   `json:"streamId"`
	Title             string   `json:"title"`
	UpdatedAt         string   `json:updatedAt`
}

func (v DbaasLogsInput) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["created_at"] = v.CreatedAt
	obj["description"] = v.Description
	obj["engine_id"] = v.EngineId
	obj["hostname"] = v.Hostname
	obj["input_id"] = v.InputId
	obj["is_restart_required"] = v.IsRestartRequired
	obj["public_address"] = v.PublicAddress
	obj["ssl_certificate"] = v.SslCertificate
	obj["status"] = v.Status
	obj["stream_id"] = v.StreamId
	obj["title"] = v.Title
	obj["updated_at"] = v.UpdatedAt

	if v.AllowedNetworks != nil {
		obj["allowed_networks"] = v.AllowedNetworks
	}
	if v.ExposedPort != nil {
		obj["exposed_port"] = *v.ExposedPort
	}
	if v.NbInstance != nil {
		obj["nb_instance"] = *v.NbInstance
	}

	return obj
}

type DbaasLogsInputOpts struct {
	Description     string   `json:"description"`
	EngineId        string   `json:"engineId"`
	StreamId        string   `json:"streamId"`
	Title           string   `json:"title"`
	AllowedNetworks []string `json:"allowedNetworks,omitempty"`
	ExposedPort     *string  `json:"exposedPort,omitempty"`
	NbInstance      *int64   `json:"nbInstance,omitempty"`
}

func (opts *DbaasLogsInputOpts) FromResource(d *schema.ResourceData) *DbaasLogsInputOpts {
	opts.Description = d.Get("description").(string)
	opts.EngineId = d.Get("engine_id").(string)
	opts.StreamId = d.Get("stream_id").(string)
	opts.Title = d.Get("title").(string)

	networks := d.Get("allowed_networks").([]interface{})
	if networks != nil && len(networks) > 0 {
		networksString := make([]string, len(networks))
		for i, net := range networks {
			networksString[i] = net.(string)
		}
		opts.AllowedNetworks = networksString
	}

	opts.ExposedPort = helpers.GetNilStringPointerFromData(d, "exposed_port")
	opts.NbInstance = helpers.GetNilInt64PointerFromData(d, "nb_instance")
	return opts
}

type DbaasLogsInputConfigurationLogstash struct {
	FilterSection  *string `json:"filterSection"`
	InputSection   string  `json:"inputSection"`
	PatternSection *string `json:"patternSection"`
}

func (v DbaasLogsInputConfigurationLogstash) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["input_section"] = html.UnescapeString(v.InputSection)

	if v.FilterSection != nil {
		obj["filter_section"] = html.UnescapeString(*v.FilterSection)
	}
	if v.PatternSection != nil {
		obj["pattern_section"] = html.UnescapeString(*v.PatternSection)
	}
	return obj
}

func (opts *DbaasLogsInputConfigurationLogstash) FromResourceWithPath(d *schema.ResourceData, path string) *DbaasLogsInputConfigurationLogstash {
	opts.InputSection = d.Get(fmt.Sprintf("%s.input_section", path)).(string)

	filterSection := d.Get(fmt.Sprintf("%s.filter_section", path)).(string)
	opts.FilterSection = &filterSection

	patternSection := d.Get(fmt.Sprintf("%s.pattern_section", path)).(string)
	opts.PatternSection = &patternSection

	return opts
}

type DbaasLogsInputConfigurationFlowgger struct {
	LogFormat  string `json:"logFormat"`
	LogFraming string `json:"logFraming"`
}

func (v DbaasLogsInputConfigurationFlowgger) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["log_format"] = html.UnescapeString(v.LogFormat)
	obj["log_framing"] = html.UnescapeString(v.LogFraming)
	return obj
}

func (opts *DbaasLogsInputConfigurationFlowgger) FromResourceWithPath(d *schema.ResourceData, path string) *DbaasLogsInputConfigurationFlowgger {
	opts.LogFormat = d.Get(fmt.Sprintf("%s.log_framing", path)).(string)
	opts.LogFraming = d.Get(fmt.Sprintf("%s.log_format", path)).(string)
	return opts
}
