package ovh

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
)

type NewRecord struct {
	Target    string `json:"target"`
	Ttl       int    `json:"ttl"`
	FieldType string `json:"fieldType"`
	SubDomain string `json:"subDomain"`
}

type Record struct {
	Id        int    `json:"id"`
	Zone      string `json:"zone"`
	Target    string `json:"target"`
	Ttl       int    `json:"ttl"`
	FieldType string `json:"fieldType"`
	SubDomain string `json:"subDomain"`
}

func resourceOVHDomainZoneRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceOVHRecordCreate,
		Read:   resourceOVHRecordRead,
		Update: resourceOVHRecordUpdate,
		Delete: resourceOVHRecordDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"zone": {
				Type:     schema.TypeString,
				Required: true,
			},
			"target": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ttl": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  3600,
			},
			"fieldtype": {
				Type:     schema.TypeString,
				Required: true,
			},
			"subdomain": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceOVHRecordCreate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*Config)

	// Create the new record
	newRecord := &NewRecord{
		FieldType: d.Get("fieldtype").(string),
		SubDomain: d.Get("subdomain").(string),
		Target:    d.Get("target").(string),
		Ttl:       d.Get("ttl").(int),
	}

	log.Printf("[DEBUG] OVH Record create configuration: %#v", newRecord)

	resultRecord := Record{}

	err := provider.OVHClient.Post(
		fmt.Sprintf("/domain/zone/%s/record", d.Get("zone").(string)),
		newRecord,
		&resultRecord,
	)

	if err != nil {
		return fmt.Errorf("Failed to create OVH Record: %s", err)
	}

	d.Set("id", resultRecord.Id)
	d.SetId(strconv.Itoa(resultRecord.Id))

	log.Printf("[INFO] OVH Record ID: %s", d.Id())

    OVHZoneRefresh(d, meta)

	return resourceOVHRecordRead(d, meta)
}

func resourceOVHRecordRead(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*Config)

	record := Record{}
	err := provider.OVHClient.Get(
		fmt.Sprintf("/domain/zone/%s/record/%s", d.Get("zone").(string), d.Id()),
		&record,
	)

	if err != nil {
		d.SetId("")
		return nil
	}

	d.Set("id", record.Id)
	d.Set("zone", record.Zone)
	d.Set("fieldtype", record.FieldType)
	d.Set("subdomain", record.SubDomain)
	d.Set("ttl", record.Ttl)
	d.Set("target", record.Target)

	return nil
}

func resourceOVHRecordUpdate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*Config)

	record := NewRecord{}

	if attr, ok := d.GetOk("subdomain"); ok {
		record.SubDomain = attr.(string)
	}
	if attr, ok := d.GetOk("fieldtype"); ok {
		record.FieldType = attr.(string)
	}
	if attr, ok := d.GetOk("target"); ok {
		record.Target = attr.(string)
	}
	if attr, ok := d.GetOk("ttl"); ok {
		record.Ttl, _ = attr.(int)
	}

	log.Printf("[DEBUG] OVH Record update configuration: %#v", record)

	err := provider.OVHClient.Put(
		fmt.Sprintf("/domain/zone/%s/record/%s", d.Get("zone").(string), d.Id()),
		record,
		nil,
	)
	if err != nil {
		return fmt.Errorf("Failed to update OVH Record: %s", err)
	}

    OVHZoneRefresh(d, meta)

	return resourceOVHRecordRead(d, meta)
}

func resourceOVHRecordDelete(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*Config)

	log.Printf("[INFO] Deleting OVH Record: %s.%s, %s", d.Get("zone").(string), d.Get("subdomain").(string), d.Id())

	err := provider.OVHClient.Delete(
		fmt.Sprintf("/domain/zone/%s/record/%s", d.Get("zone").(string), d.Id()),
		nil,
	)

	if err != nil {
		return fmt.Errorf("Error deleting OVH Record: %s", err)
	}

    OVHZoneRefresh(d, meta)

	return nil
}

func OVHZoneRefresh(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*Config)

	log.Printf("[INFO] Refresh OVH Zone: %s", d.Get("zone").(string))

	err := provider.OVHClient.Post(
		fmt.Sprintf("/domain/zone/%s/refresh", d.Get("zone").(string)),
		nil,
        nil,
	)

	if err != nil {
		return fmt.Errorf("Error refresh OVH Zone: %s", err)
	}

	return nil
}
