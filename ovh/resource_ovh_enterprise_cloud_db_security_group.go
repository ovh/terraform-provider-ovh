package ovh

import (
	"context"
	"fmt"
	"github.com/ovh/go-ovh/ovh"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const EnterpriseCloudDBSecurityGroupBaseUrl = EnterpriseCloudDBBaseUrl + "/securityGroup"

func resourceEnterpriseCloudDBSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceEnterpriseCloudDBSecurityGroupCreateOrUpdate,
		Read:   resourceEnterpriseCloudDBSecurityGroupRead,
		Delete: resourceEnterpriseCloudDBSecurityGroupDelete,
		Update: resourceEnterpriseCloudDBSecurityGroupCreateOrUpdate,
		Importer: &schema.ResourceImporter{
			State: resourceEnterpriseCloudDBSecurityGroupImportState,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeString,
				Computed: false,
				ForceNew: true,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: false,
				Required: true,
			},
		},
	}
}

func resourceEnterpriseCloudDBSecurityGroupCreateOrUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	sg := (&EnterpriseCloudDBSecurityGroupCreateUpdateOpts{}).FromResource(d)

	clusterId := d.Get("cluster_id").(string)

	getUrl := fmt.Sprintf(EnterpriseCloudDBSecurityGroupBaseUrl, clusterId)

	var securityGroup EnterpriseCloudDBSecurityGroup

	if d.Id() != "" {
		getUrl = fmt.Sprintf("%s/%s", getUrl, d.Id())
		securityGroup.Id = d.Id()
	}

	// monitor POST/PUT
	stateCreation := &resource.StateChangeConf{
		Target: []string{"OK"},
		Refresh: func() (interface{}, string, error) {
			var err error
			if d.Id() != "" {
				err = config.OVHClient.Put(getUrl, sg, nil)
			} else {
				err = config.OVHClient.Post(getUrl, sg, &securityGroup)
			}
			if err != nil {
				apiError := err.(*ovh.APIError)
				// Cluster is pending
				if apiError.Code == http.StatusForbidden {
					return d, "", nil
				}
				return d, "", fmt.Errorf("Error calling (id: %s) %s:\n\t%q", apiError.QueryID, getUrl, err)
			}
			return d, "OK", nil
		},
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateCreation.WaitForStateContext(context.Background()); err != nil {
		return fmt.Errorf("Error execution DB Entreprise Security Group Create/Update:\n\t %q", err)
	}

	getUrl = fmt.Sprintf(EnterpriseCloudDBSecurityGroupBaseUrl+"/%s", clusterId, securityGroup.Id)

	// monitor task execution
	stateConf := &resource.StateChangeConf{
		Target: []string{string(EnterpriseCloudDBStatusCreated)},
		Refresh: func() (interface{}, string, error) {
			var stateResp EnterpriseCloudDBSecurityGroup
			if err := config.OVHClient.Get(getUrl, &stateResp); err != nil {
				return nil, "", err
			}
			return d, string(stateResp.Status), nil
		},
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForStateContext(context.Background()); err != nil {
		return fmt.Errorf("Error waiting for DB Entreprise Security Group creation:\n\t %q", err)
	}

	d.SetId(securityGroup.Id)

	return resourceEnterpriseCloudDBSecurityGroupRead(d, meta)
}

func resourceEnterpriseCloudDBSecurityGroupImportState(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	givenId := d.Id()
	splitId := strings.SplitN(givenId, "/", 2)
	if len(splitId) != 2 {
		return nil, fmt.Errorf("import id is not cluster_id/security_group_id formatted")
	}
	id := splitId[0]
	groupId := splitId[1]
	d.SetId(groupId)
	d.Set("cluster_id", id)

	results := make([]*schema.ResourceData, 1)
	results[0] = d
	return results, nil
}

func resourceEnterpriseCloudDBSecurityGroupRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	clusterId := d.Get("cluster_id").(string)
	url := fmt.Sprintf(EnterpriseCloudDBSecurityGroupBaseUrl+"/%s", clusterId, d.Id())
	resp := &EnterpriseCloudDBSecurityGroup{}

	if err := config.OVHClient.Get(url, resp); err != nil {
		id := err.(*ovh.APIError).QueryID
		return fmt.Errorf("Error calling GET (%s) %s:\n\t%q", id, url, err)
	}
	d.Set("name", resp.Name)
	return nil
}

func resourceEnterpriseCloudDBSecurityGroupDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	clusterId := d.Get("cluster_id").(string)

	url := fmt.Sprintf(EnterpriseCloudDBSecurityGroupBaseUrl+"/%s", clusterId, d.Id())

	// monitor Delete action
	stateDelete := &resource.StateChangeConf{
		Target: []string{"OK"},
		Refresh: func() (interface{}, string, error) {
			if err := config.OVHClient.Delete(url, nil); err != nil {
				apiError := err.(*ovh.APIError)
				// Cluster is pending
				if apiError.Code == http.StatusForbidden {
					return d, "", nil
				}
				id := apiError.QueryID
				return d, "", fmt.Errorf("Error calling DELETE (id: %s) %s:\n\t%q", id, url, err)
			}
			return d, "OK", nil
		},
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateDelete.WaitForStateContext(context.Background()); err != nil {
		return fmt.Errorf("Error execution DB Entreprise Security Group Deletion:\n\t %q", err)
	}

	getUrl := fmt.Sprintf(EnterpriseCloudDBSecurityGroupBaseUrl+"/%s", clusterId, d.Id())
	// monitor task	execution
	stateConf := &resource.StateChangeConf{
		Target: []string{"DONE"},
		Refresh: func() (interface{}, string, error) {
			if err := config.OVHClient.Get(getUrl, nil); err != nil {
				if err.(*ovh.APIError).Code == 404 {
					return d, "DONE", nil
				}
				return nil, "NOK", err
			}
			return d, "NOK", nil
		},
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForStateContext(context.Background()); err != nil {
		return fmt.Errorf("Error waiting for DB Entreprise Security Group deletion:\n\t %q", err)
	}
	d.SetId("")
	return nil
}
