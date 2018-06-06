package ovh

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
)

type OvhDomainZoneRedirection struct {
	Id          int    `json:"id,omitempty"`
	Zone        string `json:"zone,omitempty"`
	Target      string `json:"target,omitempty"`
	SubDomain   string `json:"subDomain"`
	Type        string `json:"type,omitempty"`
	Description string `json:"description"`
	Keywords    string `json:"keywords"`
	Title       string `json:"title"`
}

func resourceOvhDomainZoneRedirection() *schema.Resource {
	return &schema.Resource{
		Create: resourceOvhDomainZoneRedirectionCreate,
		Read:   resourceOvhDomainZoneRedirectionRead,
		Update: resourceOvhDomainZoneRedirectionUpdate,
		Delete: resourceOvhDomainZoneRedirectionDelete,

		Schema: map[string]*schema.Schema{
			"zone": {
				Type:     schema.TypeString,
				Required: true,
			},
			"target": {
				Type:     schema.TypeString,
				Required: true,
			},
			"subdomain": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"keywords": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"title": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceOvhDomainZoneRedirectionCreate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*Config)

	// Create the new redirection
	newRedirection := &OvhDomainZoneRedirection{
		Type:        d.Get("type").(string),
		SubDomain:   d.Get("subdomain").(string),
		Target:      d.Get("target").(string),
		Description: d.Get("description").(string),
		Keywords:    d.Get("keywords").(string),
		Title:       d.Get("title").(string),
	}

	log.Printf("[DEBUG] OVH Redirection create configuration: %#v", newRedirection)

	resultRedirection := OvhDomainZoneRedirection{}

	err := provider.OVHClient.Post(
		fmt.Sprintf("/domain/zone/%s/redirection", d.Get("zone").(string)),
		newRedirection,
		&resultRedirection,
	)

	if err != nil {
		return fmt.Errorf("Failed to create OVH Redirection: %s", err)
	}

	d.SetId(strconv.Itoa(resultRedirection.Id))

	log.Printf("[INFO] OVH Redirection ID: %s", d.Id())

	OvhDomainZoneRefresh(d, meta)

	return resourceOvhDomainZoneRedirectionRead(d, meta)
}

func resourceOvhDomainZoneRedirectionRead(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*Config)

	redirection := OvhDomainZoneRedirection{}
	err := provider.OVHClient.Get(
		fmt.Sprintf("/domain/zone/%s/redirection/%s", d.Get("zone").(string), d.Id()),
		&redirection,
	)

	if err != nil {
		d.SetId("")
		return nil
	}

	d.Set("zone", redirection.Zone)
	d.Set("type", redirection.Type)
	d.Set("subdomain", redirection.SubDomain)
	d.Set("description", redirection.Description)
	d.Set("target", redirection.Target)
	d.Set("keywords", redirection.Keywords)
	d.Set("title", redirection.Title)

	return nil
}

func resourceOvhDomainZoneRedirectionUpdate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*Config)

	redirection := OvhDomainZoneRedirection{}

	if attr, ok := d.GetOk("subdomain"); ok {
		redirection.SubDomain = attr.(string)
	}
	if attr, ok := d.GetOk("type"); ok {
		redirection.Type = attr.(string)
	}
	if attr, ok := d.GetOk("target"); ok {
		redirection.Target = attr.(string)
	}
	if attr, ok := d.GetOk("description"); ok {
		redirection.Description, _ = attr.(string)
	}
	if attr, ok := d.GetOk("keywords"); ok {
		redirection.Keywords, _ = attr.(string)
	}
	if attr, ok := d.GetOk("title"); ok {
		redirection.Title, _ = attr.(string)
	}

	log.Printf("[DEBUG] OVH Redirection update configuration: %#v", redirection)

	err := provider.OVHClient.Put(
		fmt.Sprintf("/domain/zone/%s/redirection/%s", d.Get("zone").(string), d.Id()),
		redirection,
		nil,
	)
	if err != nil {
		return fmt.Errorf("Failed to update OVH Redirection: %s", err)
	}

	OvhDomainZoneRefresh(d, meta)

	return resourceOvhDomainZoneRedirectionRead(d, meta)
}

func resourceOvhDomainZoneRedirectionDelete(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*Config)

	log.Printf("[INFO] Deleting OVH Redirection: %s.%s, %s", d.Get("zone").(string), d.Get("subdomain").(string), d.Id())

	err := provider.OVHClient.Delete(
		fmt.Sprintf("/domain/zone/%s/redirection/%s", d.Get("zone").(string), d.Id()),
		nil,
	)

	if err != nil {
		return fmt.Errorf("Error deleting OVH Redirection: %s", err)
	}

	OvhDomainZoneRefresh(d, meta)

	return nil
}
