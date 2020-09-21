package ovh

import (
	"context"
	"fmt"
	"github.com/ovh/go-ovh/ovh"
	"github.com/terraform-providers/terraform-provider-ovh/ovh/helpers"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const CloudDBEnterpriseSecurityGroupRuleBaseUrl = CloudDBEnterpriseSecurityGroupBaseUrl + "/%s/rule"

func resourceCloudDBEnterpriseSecurityGroupRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudDBEnterpriseSecurityGroupRuleCreate,
		Read:   resourceCloudDBEnterpriseSecurityGroupRuleRead,
		Delete: resourceCloudDBEnterpriseSecurityGroupRuleDelete,
		Importer: &schema.ResourceImporter{
			State: resourceCloudDBEnterpriseSecurityGroupRuleImportState,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeString,
				Computed: false,
				ForceNew: true,
				Required: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Computed: false,
				Required: true,
				ForceNew: true,
			},
			"source": {
				Type:     schema.TypeString,
				Computed: false,
				ForceNew: true,
				Required: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					err := helpers.ValidateIpBlock(v.(string))
					if err != nil {
						errors = append(errors, err)
					}
					return
				},
			},
		},
	}
}

func resourceCloudDBEnterpriseSecurityGroupRuleCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	sg := (&CloudDBEnterpriseSecurityGroupRuleCreateUpdateOpts{}).FromResource(d)
	clusterId := d.Get("cluster_id").(string)
	securityGroupId := d.Get("security_group_id").(string)

	url := fmt.Sprintf(CloudDBEnterpriseSecurityGroupRuleBaseUrl, clusterId, securityGroupId)
	var securityGroupRule CloudDBEnterpriseSecurityGroupRule

	// monitor POST/PUT
	stateCreation := &resource.StateChangeConf{
		Target: []string{"OK"},
		Refresh: func() (interface{}, string, error) {
			if err := config.OVHClient.Post(url, sg, &securityGroupRule); err != nil {
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
		Timeout:    20 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateCreation.WaitForStateContext(context.Background()); err != nil {
		return fmt.Errorf("Error waiting for DB Entreprise Security Group creation:\n\t %q", err)
	}

	getUrl := fmt.Sprintf(CloudDBEnterpriseSecurityGroupRuleBaseUrl+"/%s", clusterId, securityGroupId, securityGroupRule.Id)

	// monitor task execution
	stateConf := &resource.StateChangeConf{
		Target: []string{string(CloudDBEnterpriseSecurityGroupRuleStatusCreated)},
		Refresh: func() (interface{}, string, error) {
			var resp CloudDBEnterpriseSecurityGroupRule
			if err := config.OVHClient.Get(getUrl, &resp); err != nil {
				return nil, "", err
			}
			return d, string(resp.Status), nil
		},
		Timeout:    20 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForStateContext(context.Background()); err != nil {
		return fmt.Errorf("Error waiting for DB Entreprise Security Group creation:\n\t %q", err)
	}

	d.SetId(securityGroupRule.Id)

	return resourceCloudDBEnterpriseSecurityGroupRuleRead(d, meta)
}

func resourceCloudDBEnterpriseSecurityGroupRuleImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	givenId := d.Id()
	splitId := strings.SplitN(givenId, "/", 3)
	if len(splitId) != 3 {
		return nil, fmt.Errorf("import id is not cluster_id/security_group_id/security_group_rule formatted")
	}
	d.Set("cluster_id", splitId[0])
	d.Set("security_group_id", splitId[1])
	d.SetId(splitId[2])

	results := make([]*schema.ResourceData, 1)
	results[0] = d
	return results, nil
}

func resourceCloudDBEnterpriseSecurityGroupRuleRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	clusterId := d.Get("cluster_id").(string)
	securityGroupId := d.Get("security_group_id").(string)
	url := fmt.Sprintf(CloudDBEnterpriseSecurityGroupRuleBaseUrl+"/%s", clusterId, securityGroupId, d.Id())
	resp := &CloudDBEnterpriseSecurityGroupRule{}

	if err := config.OVHClient.Get(url, resp); err != nil {
		id := err.(*ovh.APIError).QueryID
		return fmt.Errorf("Error calling GET (%s) %s:\n\t%q", id, url, err)
	}
	d.Set("source", resp.Source)
	return nil
}

func resourceCloudDBEnterpriseSecurityGroupRuleDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	clusterId := d.Get("cluster_id").(string)
	securityGroupId := d.Get("security_group_id").(string)
	url := fmt.Sprintf(CloudDBEnterpriseSecurityGroupRuleBaseUrl+"/%s", clusterId, securityGroupId, d.Id())

	// monitor POST/PUT
	stateCreation := &resource.StateChangeConf{
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

	if _, err := stateCreation.WaitForStateContext(context.Background()); err != nil {
		return fmt.Errorf("Error execution DB Entreprise Security Group Deletion:\n\t %q", err)
	}

	getUrl := fmt.Sprintf(CloudDBEnterpriseSecurityGroupRuleBaseUrl+"/%s", clusterId, securityGroupId, d.Id())

	// monitor task execution
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
		Timeout:    20 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForStateContext(context.Background()); err != nil {
		return fmt.Errorf("Error waiting for DB Entreprise Security Group deletion:\n\t %q", err)
	}
	d.SetId("")
	return nil
}
