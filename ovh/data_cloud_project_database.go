package ovh

import (
	"fmt"
	"log"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCloudProjectDatabase() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudProjectDatabaseRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_CLOUD_PROJECT_SERVICE", nil),
			},
			"engine": {
				Type:        schema.TypeString,
				Description: "Name of the engine of the service",
				Required:    true,
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Description: "Cluster ID",
				Required:    true,
			},

			//Computed
			"backup_time": {
				Type:        schema.TypeString,
				Description: "Time on which backups start every day",
				Computed:    true,
			},
			"created_at": {
				Type:        schema.TypeString,
				Description: "Date of the creation of the cluster",
				Computed:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Description of the cluster",
				Computed:    true,
			},
			"endpoints": {
				Type:        schema.TypeList,
				Description: "List of all endpoints of the service",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"component": {
							Type:        schema.TypeString,
							Description: "Type of component the URI relates to",
							Computed:    true,
						},
						"domain": {
							Type:        schema.TypeString,
							Description: "Domain of the cluster",
							Computed:    true,
						},
						"path": {
							Type:        schema.TypeString,
							Description: "Path of the endpoint",
							Computed:    true,
						},
						"port": {
							Type:        schema.TypeInt,
							Description: "Connection port for the endpoint",
							Computed:    true,
						},
						"scheme": {
							Type:        schema.TypeString,
							Description: "Scheme used to generate the URI",
							Computed:    true,
						},
						"ssl": {
							Type:        schema.TypeBool,
							Description: "Defines whether the endpoint uses SSL",
							Computed:    true,
						},
						"ssl_mode": {
							Type:        schema.TypeString,
							Description: "SSL mode used to connect to the service if the SSL is enabled",
							Computed:    true,
						},
						"uri": {
							Type:        schema.TypeString,
							Description: "URI of the endpoint",
							Computed:    true,
						},
					},
				},
			},
			"flavor": {
				Type:        schema.TypeString,
				Description: "The node flavor used for this cluster",
				Computed:    true,
			},
			"maintenance_time": {
				Type:        schema.TypeString,
				Description: "Time on which maintenances can start every day",
				Computed:    true,
			},
			"network_type": {
				Type:        schema.TypeString,
				Description: "Type of network of the cluster",
				Computed:    true,
			},
			"nodes": {
				Type:        schema.TypeList,
				Description: "List of nodes composing the service",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"network_id": {
							Type:        schema.TypeString,
							Description: "Private network ID in which the node is",
							Computed:    true,
						},
						"region": {
							Type:        schema.TypeString,
							Description: "Region of the node",
							Computed:    true,
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Description: "Private subnet ID in which the node is",
							Computed:    true,
						},
					},
				},
			},
			"plan": {
				Type:        schema.TypeString,
				Description: "Plan of the cluster",
				Computed:    true,
			},
			"status": {
				Type:        schema.TypeString,
				Description: "Current status of the cluster",
				Computed:    true,
			},
			"version": {
				Type:        schema.TypeString,
				Description: "Version of the engine deployed on the cluster",
				Computed:    true,
			},
		},
	}
}

func dataSourceCloudProjectDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	engine := d.Get("engine").(string)
	id := d.Get("cluster_id").(string)

	serviceEndpoint := fmt.Sprintf("/cloud/project/%s/database/%s/%s",
		url.PathEscape(serviceName),
		url.PathEscape(engine),
		url.PathEscape(id),
	)
	res := &CloudProjectDatabaseResponse{}

	log.Printf("[DEBUG] Will read database %s from project: %s", id, serviceName)
	if err := config.OVHClient.Get(serviceEndpoint, res); err != nil {
		return fmt.Errorf("Error calling GET %s:\n\t %q", serviceEndpoint, err)
	}

	nodesEndpoint := fmt.Sprintf("%s/node", serviceEndpoint)
	nodeList := &[]string{}
	if err := config.OVHClient.Get(nodesEndpoint, nodeList); err != nil {
		return fmt.Errorf("unable to get database %s nodes: %v", res.Id, err)
	}

	if len(*nodeList) == 0 {
		return fmt.Errorf("no node found for database %s", res.Id)
	}
	nodeEndpoint := fmt.Sprintf("%s/%s", nodesEndpoint, url.PathEscape((*nodeList)[0]))
	node := &CloudProjectDatabaseNodes{}
	if err := config.OVHClient.Get(nodeEndpoint, node); err != nil {
		return fmt.Errorf("unable to get database %s node %s: %v", res.Id, (*nodeList)[0], err)
	}

	res.Region = node.Region

	for k, v := range res.ToMap() {
		if k != "id" {
			d.Set(k, v)
		} else {
			d.SetId(fmt.Sprint(v))
		}
	}

	log.Printf("[DEBUG] Read Database %+v", res)
	return nil
}
