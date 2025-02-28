package ovh

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

func resourceDomainDsRecords() *schema.Resource {
	return &schema.Resource{
		Description: "Resource to manage a domain name DS records",
		Schema:      resourceDomainDsRecordsSchema(),
		Create:      resourceDomainDsRecordsUpdate,
		Read:        resourceDomainDsRecordsRead,
		Update:      resourceDomainDsRecordsUpdate,
		Delete:      resourceDomainDsRecordsDelete,
		Importer: &schema.ResourceImporter{
			State: func(resourceData *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
				resourceData.Set("domain", resourceData.Id())
				return []*schema.ResourceData{resourceData}, nil
			},
		},

		Timeouts: &schema.ResourceTimeout{
			Default: schema.DefaultTimeout(10 * time.Minute),
			Read:    schema.DefaultTimeout(30 * time.Second),
			Create:  schema.DefaultTimeout(10 * time.Minute),
			Update:  schema.DefaultTimeout(10 * time.Minute),
			Delete:  schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceDomainDsRecordsSchema() map[string]*schema.Schema {
	resourceSchema := map[string]*schema.Schema{
		"domain": {
			Type:        schema.TypeString,
			Description: "Domain name",
			Required:    true,
			ForceNew:    true,
		},
		"ds_records": {
			Type:        schema.TypeList,
			Description: "DS Records for the domain",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"algorithm": {
						Type:         schema.TypeString,
						Description:  "Algorithm name of the DNSSEC key",
						Required:     true,
						ValidateFunc: helpers.ValidateEnum([]string{"RSASHA1", "RSASHA1_NSEC3_SHA1", "RSASHA256", "RSASHA512", "ECDSAP256SHA256", "ECDSAP384SHA384", "ED25519"}),
					},
					"flags": {
						Type:         schema.TypeString,
						Description:  "Flag name of the DNSSEC key",
						Required:     true,
						ValidateFunc: helpers.ValidateEnum([]string{"ZONE_SIGNING_KEY", "KEY_SIGNING_KEY"}),
					},
					"public_key": {
						Type:        schema.TypeString,
						Description: "Public key",
						Required:    true,
					},
					"tag": {
						Type:        schema.TypeInt,
						Description: "Tag of the DNSSEC key",
						Required:    true,
						ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
							if v.(int) <= 0 || v.(int) >= 65536 {
								errors = append(errors, fmt.Errorf(`Field "tag" must be larger than 0 and smaller than 65536.`))
								return
							}

							return
						},
					},
				},
			},
			Required: true,
			MinItems: 1,
			MaxItems: 4,
		},
	}

	return resourceSchema
}

func resourceDomainDsRecordsRead(resourceData *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	domainName := resourceData.Id()

	domainDsRecords := &DomainDsRecords{
		Domain: domainName,
	}

	log.Printf("[DEBUG] Will read domain name DS records: %s\n", domainName)

	var responseData []int
	endpoint := fmt.Sprintf("/domain/%s/dsRecord", url.PathEscape(domainName))

	if err := config.OVHClient.Get(endpoint, &responseData); err != nil {
		return fmt.Errorf("calling GET %s:\n\t %s", endpoint, err.Error())
	}

	for _, dsRecordId := range responseData {
		responseData := DomainDsRecord{}
		endpoint := fmt.Sprintf("/domain/%s/dsRecord/%d", url.PathEscape(domainName), dsRecordId)

		if err := config.OVHClient.Get(endpoint, &responseData); err != nil {
			return helpers.CheckDeleted(resourceData, err, endpoint)
		}

		domainDsRecords.DsRecords = append(domainDsRecords.DsRecords, responseData)
	}

	resourceData.SetId(domainName)
	for k, v := range domainDsRecords.ToMap() {
		resourceData.Set(k, v)
	}

	return nil
}

func resourceDomainDsRecordsUpdate(resourceData *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	domainName := resourceData.Get("domain").(string)
	task := DomainTask{}

	dsRecordsUpdate := &DomainDsRecordsUpdateOpts{
		DsRecords: make([]DomainDsRecord, 0),
	}

	for _, dsRecord := range resourceData.Get("ds_records").([]interface{}) {
		record := dsRecord.(map[string]interface{})

		dsRecordsUpdate.DsRecords = append(dsRecordsUpdate.DsRecords, DomainDsRecord{
			Algorithm: DsRecordAlgorithmValuesMap[record["algorithm"].(string)],
			Flags:     DsRecordFlagValuesMap[record["flags"].(string)],
			PublicKey: record["public_key"].(string),
			Tag:       record["tag"].(int),
		})
	}

	log.Printf("[DEBUG] Will update domain name DS records: %s\n", domainName)

	endpoint := fmt.Sprintf("/domain/%s/dsRecord", url.PathEscape(domainName))

	if err := config.OVHClient.Post(endpoint, dsRecordsUpdate, &task); err != nil {
		return fmt.Errorf("calling POST %s :\n\t %s", endpoint, err.Error())
	}

	if err := waitDomainTask(config.OVHClient, domainName, task.TaskID); err != nil {
		return fmt.Errorf("waiting for %s DS records to be updated: %s", domainName, err.Error())
	}

	resourceData.SetId(domainName)

	return resourceDomainDsRecordsRead(resourceData, meta)
}

func resourceDomainDsRecordsDelete(resourceData *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	domainName := resourceData.Get("domain").(string)
	task := DomainTask{}

	domainDsRecordsUpdateOpts := &DomainDsRecordsUpdateOpts{
		DsRecords: make([]DomainDsRecord, 0),
	}

	log.Printf("[DEBUG] Will remove all domain name DS records: %s\n", domainName)

	endpoint := fmt.Sprintf("/domain/%s/dsRecord", url.PathEscape(domainName))

	if err := config.OVHClient.Post(endpoint, domainDsRecordsUpdateOpts, &task); err != nil {
		return fmt.Errorf("calling POST %s :\n\t %s", endpoint, err.Error())
	}

	if err := waitDomainTask(config.OVHClient, domainName, task.TaskID); err != nil {
		return fmt.Errorf("waiting for %s DS records to be deleted: %s", domainName, err.Error())
	}

	resourceData.SetId("")

	return nil
}
