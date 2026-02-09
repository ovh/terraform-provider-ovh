package ovh

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/go-ovh/ovh"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
	"golang.org/x/exp/slices"
)

// taskExpiresAfter is the duration in seconds after which we'll consider an ongoing task to be expired and we'll allow ourselves to create a new one.
// Usually, time taken for such a task is around 1 minute, here we tolerate 5 minutes
const taskExpiresAfter = 300 * time.Second

// waitingTimeInSecondsBeforeRefreshState number if seconds to wait before making a new API call to refresh ip task state
const waitingTimeInSecondsBeforeRefreshState = 10

var ipTaskUnrecoverableErrors = []int{http.StatusBadRequest, http.StatusUnauthorized, http.StatusForbidden, http.StatusNotFound}

func resourceIpServiceMove() *schema.Resource {
	return &schema.Resource{
		Create: resourceIpMoveCreate,
		Update: resourceIpMoveUpdate,
		Read:   resourceIpRead,
		Delete: resourceIpMoveDelete,

		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				d.Set("ip", d.Id())
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: resourceIpMoveSchema(),
	}
}

func resourceIpMoveSchema() map[string]*schema.Schema {
	schema := map[string]*schema.Schema{
		"description": {
			Type:        schema.TypeString,
			Description: "Custom description on your ip",
			Optional:    true,
			Computed:    true,
		},

		//computed
		"can_be_terminated": {
			Type:     schema.TypeBool,
			Computed: true,
		},

		"country": {
			Type:     schema.TypeString,
			Computed: true,
		},

		"ip": {
			Type:     schema.TypeString,
			Required: true,
		},
		"organisation_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"routed_to": {
			Type:        schema.TypeList,
			MinItems:    1,
			MaxItems:    1,
			Description: "Routage information",
			Required:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"service_name": {
						Type:        schema.TypeString,
						Description: "Service where ip is routed to",
						Required:    true,
					},
				},
			},
		},
		"service_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"type": {
			Type:        schema.TypeString,
			Description: "Possible values for ip type",
			Computed:    true,
		},
		"task_status": {
			Type:        schema.TypeString,
			Description: "Status field of the current IP task that is in charge of changing the service the IP is attached to",
			Computed:    true,
		},
		"task_start_date": {
			Type:        schema.TypeString,
			Description: "Starting date and time field of the current IP task that is in charge of changing the service the IP is attached to",
			Computed:    true,
		},
	}

	return schema
}

func resourceIpMoveCreate(d *schema.ResourceData, meta interface{}) error {
	// later on this ID will be replaced by the task if when we need to create it (see resourceIpMoveUpdate)
	d.SetId(d.Get("ip").(string))

	err := resourceIpMoveUpdate(d, meta)
	if err != nil {
		d.SetId("")
	}

	return err
}

// resourceIpMoveUpdate will move an ip to a provided service name or detach (= park) it otherwise
// if the resource ID is a taskId and the previous task is not done, wait for it to be finished until the previous task is considered expired
// (see recursiveWaitTaskFinished) before trying to do the service move.
// then do the move only if the current service is different from the one given in the inputs
func resourceIpMoveUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	ip := d.Get("ip").(string)

	opts, err := (&IpMoveOpts{}).FromResource(d)
	if err != nil {
		return err
	}

	serviceName, err := helpers.ServiceNameFromIpBlock(ip)
	if err != nil {
		return err
	}
	err = d.Set("service_name", serviceName)
	if err != nil {
		return err
	}

	ipTask := &IpTask{}
	taskId, err := strconv.ParseInt(d.Id(), 10, 64)
	if err == nil {
		// if previous task is not done yet we need to wait for it to be completed
		if d.Get("task_status") != nil {
			ipTask.Status = IpTaskStatusEnum(d.Get("task_status").(string))
			ipTask.TaskId = taskId
			taskStartDate, err := time.Parse(time.RFC3339, d.Get("task_start_date").(string))
			if err != nil {
				return err
			}
			ipTask.StartDate = taskStartDate
			_, err = waitForTaskFinished(d, meta, ipTask, ip, opts)
			if err != nil {
				log.Printf("[WARNING] - waitForTaskFinished on previously registered task return error %s. Will continue nevertheless", err)
			}
		}
	} else {
		log.Printf("[WARNING] - resource ID %s is not an int64/not a task ID. Cannot get last task state", d.Id())
	}
	err = resourceIpReadByIp(d, ip, config)
	if err != nil {
		return err
	}

	currentlyRoutedService := GetRoutedToServiceName(d)
	// no need to update if ip is already routed to the appropriate service
	if stringPtrEqual(currentlyRoutedService, opts.To) {
		log.Printf("[DEBUG] Won't do anything as ip %s (service name = %s) is already routed to service %v", ip, serviceName, currentlyRoutedService)
		return nil
	}

	if opts.To == nil {
		var currentService string
		if currentlyRoutedService != nil {
			currentService = *currentlyRoutedService
		}
		log.Printf("[DEBUG] Will move ip %s (service name = %s) from service %s to IP parking", ip, serviceName, currentService)
		endpoint := fmt.Sprintf("/ip/%s/park",
			url.PathEscape(ip),
		)
		// retrieve the task
		if err := config.OVHClient.Post(endpoint, nil, ipTask); err != nil {
			return fmt.Errorf("calling Post %s: %q", endpoint, err)
		}
	} else {
		log.Printf("[DEBUG] Will move ip %s (service name = %s) from service %v to service %s", ip, serviceName, currentlyRoutedService, *opts.To)
		endpoint := fmt.Sprintf("/ip/%s/move",
			url.PathEscape(ip),
		)
		if err := config.OVHClient.Post(endpoint, opts, ipTask); err != nil {
			return fmt.Errorf("calling Post %s: %q", endpoint, err)
		}
	}
	d.SetId(fmt.Sprint(ipTask.TaskId))
	if err = d.Set("task_start_date", ipTask.StartDate.Format(time.RFC3339)); err != nil {
		return err
	}

	_, err = waitForTaskFinished(d, meta, ipTask, ip, opts)
	if err != nil {
		return err
	}

	return resourceIpRead(d, meta)
}

// waitForTaskFinished queries GET /ip/:ip/task/:taskId route until task state is in a terminal success or error state or until taskExpiresAfter is reached
// and returns :
//   - finishedWithSuccess : true if task ended with success, false if ended with error, nil if not ended at all
//   - err : nil if no error is encountered in the process, any met error otherwise
//
// in any case before returning, "task_status" field of d will be updated with the last known ipTask.Status
func waitForTaskFinished(d *schema.ResourceData, meta interface{}, ipTask *IpTask, ip string, opts *IpMoveOpts) (finishedWithSuccess *bool, err error) {
	if ipTask == nil {
		return helpers.GetNilBoolPointer(false), fmt.Errorf("could not assign IP %s to service %v as Ip task does not exist", ip, opts.To)
	}

	var result *bool
	for {
		switch ipTask.Status {
		case IpTaskStatusDone:
			result = helpers.GetNilBoolPointer(true)
			err = d.Set("task_status", ipTask.Status)
			return result, err
		case IpTaskStatusCancelled, IpTaskStatusOvhError, IpTaskStatusCustomerError:
			result = helpers.GetNilBoolPointer(false)
			errSet := d.Set("task_status", ipTask.Status)
			return result, errors.Join(fmt.Errorf("could not assign IP %s to service %v as Ip task is %s", ip, opts.To, ipTask.Status), errSet)
		}

		timeDiff := time.Since(ipTask.StartDate)
		if timeDiff >= taskExpiresAfter {
			log.Printf("[WARNING] - waitForTaskFinished max time reached without the task having reached a terminal state")
			_ = d.Set("task_status", ipTask.Status)
			return nil, nil
		}

		log.Printf("[DEBUG] ipTask.Status is currently: %s. Waiting %d seconds (we allow %f more seconds for the task to complete)", ipTask.Status, waitingTimeInSecondsBeforeRefreshState, taskExpiresAfter.Seconds()-timeDiff.Seconds())
		time.Sleep(waitingTimeInSecondsBeforeRefreshState * time.Second)

		err = resourceIpTaskRead(d, ipTask, meta)
		if errOvh, ok := err.(*ovh.APIError); ok {
			// bad request, unauthorized, forbidden & not found errors are unrecoverable
			if slices.Contains(ipTaskUnrecoverableErrors, errOvh.Code) {
				_ = d.Set("task_status", ipTask.Status)
				return helpers.GetNilBoolPointer(false), err
			}
		}
	}
}

// stringPtrEqual compares two string pointers for equality
func stringPtrEqual(a, b *string) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

func resourceIpTaskRead(d *schema.ResourceData, ipTask *IpTask, meta interface{}) error {
	config := meta.(*Config)

	endpoint := fmt.Sprintf("/ip/%s/task/%d",
		url.PathEscape(d.Get("ip").(string)),
		ipTask.TaskId,
	)

	return config.OVHClient.Get(endpoint, ipTask)
}

func resourceIpRead(d *schema.ResourceData, meta interface{}) error {
	ip := d.Get("ip").(string)
	serviceName, err := helpers.ServiceNameFromIpBlock(ip)
	if err != nil {
		return err
	}
	err = d.Set("service_name", serviceName)
	if err != nil {
		return err
	}
	config := meta.(*Config)
	return resourceIpReadByIp(d, ip, config)
}

// resourceIpMoveDelete is an empty implementation as move do not actually create API objects but rather updates the underlying ip spec (by modifying its routed_to service)
func resourceIpMoveDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceIpReadByIp(d *schema.ResourceData, ip string, config *Config) error {
	r := &Ip{}
	endpoint := fmt.Sprintf("/ip/%s",
		url.PathEscape(ip),
	)

	if err := config.OVHClient.Get(endpoint, &r); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	// set resource attributes
	for k, v := range r.ToMap() {
		d.Set(k, v)
	}

	return nil
}
