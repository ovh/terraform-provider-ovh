package ovh

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sync"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/ovh/go-ovh/ovh"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/ovhwrap"
)

// vpsBackend reads VPS as the types of the API v2 contract, whatever API
// serves them. The data sources only know this interface.
type vpsBackend interface {
	// Get reads one VPS. A VPS that does not exist is an error.
	Get(ctx context.Context, id string) (*VPSAPIServer, diag.Diagnostics)
	// List reads every VPS of the account.
	List(ctx context.Context) ([]VPSAPIServer, diag.Diagnostics)
}

// newVPSBackend returns the backend the provider uses: API v1 until API v2
// serves VPS.
func newVPSBackend(config *Config) vpsBackend {
	return &vpsV1Backend{client: config.OVHClient}
}

// vpsV1Backend reads VPS through the API v1 routes /vps/{serviceName}/...
// API v1 keeps no target: the target specification is the current state.
type vpsV1Backend struct {
	client *ovhwrap.Client

	// locations caches GET /v2/location once it succeeded.
	locationsMu sync.Mutex
	locations   []vpsLocation
}

// vpsV1 is the part of GET /vps/{serviceName} read by vpsV1Backend.
type vpsV1 struct {
	DisplayName string `json:"displayName"`
	State       string `json:"state"`
	Vcore       int64  `json:"vcore"`
	MemoryLimit int64  `json:"memoryLimit"`
	Model       struct {
		Disk int64 `json:"disk"`
	} `json:"model"`
}

// vpsV1Image is GET /vps/{serviceName}/images/current.
type vpsV1Image struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

// vpsV1Datacenter is the part of GET /vps/{serviceName}/datacenter read by
// vpsV1Backend.
type vpsV1Datacenter struct {
	Name string `json:"name"`
}

func (b *vpsV1Backend) Get(ctx context.Context, id string) (*VPSAPIServer, diag.Diagnostics) {
	var diags diag.Diagnostics
	base := "/vps/" + url.PathEscape(id)

	var vps vpsV1
	if err := b.client.GetWithContext(ctx, base, &vps); err != nil {
		if vpsIsNotFound(err) {
			diags.AddError(
				fmt.Sprintf("VPS %s not found", id),
				fmt.Sprintf("GET %s answered: %s", base, err),
			)
			return nil, diags
		}
		diags.AddError(fmt.Sprintf("Error calling GET %s", base), err.Error())
		return nil, diags
	}

	from, err := b.installationFrom(ctx, base)
	if err != nil {
		diags.AddError(fmt.Sprintf("Error calling GET %s/images/current", base), err.Error())
		return nil, diags
	}

	var ips []string
	endpoint := base + "/ips"
	if err := b.client.GetWithContext(ctx, endpoint, &ips); err != nil {
		diags.AddError(fmt.Sprintf("Error calling GET %s", endpoint), err.Error())
		return nil, diags
	}
	ipAddresses := make([]VPSAPIIPAddress, 0, len(ips))
	for _, ip := range ips {
		version, err := vpsIPVersion(ip)
		if err != nil {
			diags.AddError(fmt.Sprintf("Error reading the IP addresses of VPS %s", id), err.Error())
			return nil, diags
		}
		ipAddresses = append(ipAddresses, VPSAPIIPAddress{Ip: ip, Version: version})
	}

	region, regionDiags := b.region(ctx, id, base)
	diags.Append(regionDiags...)
	if diags.HasError() {
		return nil, diags
	}

	targetPowerState, currentPowerState := vpsPowerStatesFromV1(vps.State)

	return &VPSAPIServer{
		Id: id,
		TargetSpec: VPSAPIReadTargetSpec{
			DisplayName:  vps.DisplayName,
			PowerState:   targetPowerState,
			Installation: VPSAPIReadInstallation{From: from},
		},
		CurrentState: &VPSAPICurrentState{
			DisplayName:  vps.DisplayName,
			PowerState:   currentPowerState,
			Installation: VPSAPICurrentInstallation{From: from},
			Location:     VPSAPILocation{Region: region},
			Model:        VPSAPIHardware{Vcpus: vps.Vcore, Ram: vps.MemoryLimit, Disk: vps.Model.Disk},
			IpAddresses:  ipAddresses,
		},
	}, diags
}

func (b *vpsV1Backend) List(ctx context.Context) ([]VPSAPIServer, diag.Diagnostics) {
	var diags diag.Diagnostics
	diags.AddError("Listing VPS is not implemented", "vpsV1Backend.List is not implemented yet")
	return nil, diags
}

// installationFrom returns the image the VPS is installed from, with null
// fields when it is unknown: GET /images/current answers 404 when API v1
// cannot trace the image back (the first image of the chain, such as a private
// image deleted after a restoration, no longer exists).
func (b *vpsV1Backend) installationFrom(ctx context.Context, base string) (VPSAPIInstallationFrom, error) {
	var image vpsV1Image
	if err := b.client.GetWithContext(ctx, base+"/images/current", &image); err != nil {
		if vpsIsNotFound(err) {
			return VPSAPIInstallationFrom{}, nil
		}
		return VPSAPIInstallationFrom{}, err
	}
	return VPSAPIInstallationFrom{
		ImageId:   vpsNonEmpty(image.Id),
		ImageName: vpsNonEmpty(image.Name),
	}, nil
}

// region returns the region of the VPS, nil with a warning when its datacenter
// code matches no region.
func (b *vpsV1Backend) region(ctx context.Context, id, base string) (*string, diag.Diagnostics) {
	var diags diag.Diagnostics

	var datacenter vpsV1Datacenter
	endpoint := base + "/datacenter"
	if err := b.client.GetWithContext(ctx, endpoint, &datacenter); err != nil {
		diags.AddError(fmt.Sprintf("Error calling GET %s", endpoint), err.Error())
		return nil, diags
	}

	locations, err := b.listLocations(ctx)
	if err != nil {
		diags.AddError("Error calling GET /v2/location", err.Error())
		return nil, diags
	}

	region, ok := vpsRegionFromDatacenterCode(datacenter.Name, locations)
	if !ok {
		diags.AddWarning(
			fmt.Sprintf("Unknown region for VPS %s", id),
			fmt.Sprintf("The datacenter code %q of the VPS matches no single region of GET /v2/location: its region is left empty.", datacenter.Name),
		)
		return nil, diags
	}
	return &region, diags
}

// listLocations returns GET /v2/location, called once per backend. A failure
// is not kept: the next call tries again.
func (b *vpsV1Backend) listLocations(ctx context.Context) ([]vpsLocation, error) {
	b.locationsMu.Lock()
	defer b.locationsMu.Unlock()

	if b.locations != nil {
		return b.locations, nil
	}
	locations, err := helpers.GetAllPagesV2[vpsLocation](ctx, b.client, "/v2/location")
	if err != nil {
		return nil, err
	}
	if locations == nil {
		locations = []vpsLocation{}
	}
	b.locations = locations
	return locations, nil
}

// vpsNonEmpty returns nil for an empty string, a pointer to it otherwise.
func vpsNonEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func vpsIsNotFound(err error) bool {
	var apiErr *ovh.APIError
	return errors.As(err, &apiErr) && apiErr.Code == http.StatusNotFound
}
