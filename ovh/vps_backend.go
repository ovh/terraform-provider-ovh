package ovh

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/netip"
	"net/url"
	"slices"
	"strconv"
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
	cache, _ := vpsLocationCaches.LoadOrStore(config.OVHClient, &vpsLocationCache{})
	return &vpsV1Backend{client: config.OVHClient, locations: cache.(*vpsLocationCache)}
}

// vpsLocationCaches holds one vpsLocationCache per API client, so per provider
// configuration: every data source read of a provider shares it.
var vpsLocationCaches sync.Map // *ovhwrap.Client -> *vpsLocationCache

// vpsLocationCache caches the answer of GET /v2/location: the locations, or a
// refusal (vpsLocationsDenied), which a key whose access rules do not allow the
// call gets for the whole run, so that it is called once, not once per VPS.
// Any other failure is not kept: the next call tries again.
type vpsLocationCache struct {
	mu        sync.Mutex
	fetched   bool
	locations []vpsLocation
	err       error
}

func (c *vpsLocationCache) get(ctx context.Context, client *ovhwrap.Client) ([]vpsLocation, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.fetched {
		return c.locations, c.err
	}
	locations, err := helpers.GetAllPagesV2[vpsLocation](ctx, client, "/v2/location")
	if err != nil && !vpsLocationsDenied(err) {
		return nil, err
	}
	c.locations, c.err, c.fetched = locations, err, true
	return locations, err
}

// vpsLocationsDenied tells whether GET /v2/location failed because the
// credentials are not allowed to call it.
func vpsLocationsDenied(err error) bool {
	status := vpsAPIStatus(err)
	return status == http.StatusUnauthorized || status == http.StatusForbidden
}

// vpsV1StatusExpired is the HTTP status API v1 answers on the routes of an
// expired VPS ("This service is expired").
const vpsV1StatusExpired = 460

// vpsV1ListParallelism is the number of VPS List reads at once.
const vpsV1ListParallelism = 8

// errVPSNotFound wraps the error of a VPS whose GET /vps/{serviceName}
// answers 404.
var errVPSNotFound = errors.New("VPS not found")

// errVPSExpired wraps the error of a VPS whose GET /vps/{serviceName}
// answers vpsV1StatusExpired.
var errVPSExpired = errors.New("VPS expired")

// vpsV1Backend reads VPS through the API v1 routes /vps/{serviceName}/...
// API v1 keeps no target: the target specification is the current state.
type vpsV1Backend struct {
	client    *ovhwrap.Client
	locations *vpsLocationCache
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

// vpsV1AutomatedBackup is the part of GET /vps/{serviceName}/automatedBackup
// read by vpsV1Backend.
type vpsV1AutomatedBackup struct {
	State    string  `json:"state"`    // enabled or disabled
	Schedule *string `json:"schedule"` // HH:MM:SS, UTC
	Rotation int64   `json:"rotation"`
}

// vpsV1Disk is the part of GET /vps/{serviceName}/disks/{id} read by
// vpsV1Backend.
type vpsV1Disk struct {
	Type string `json:"type"` // primary or additional
	Size int64  `json:"size"` // GB
}

func (b *vpsV1Backend) Get(ctx context.Context, id string) (*VPSAPIServer, diag.Diagnostics) {
	srv, diags, err := b.read(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, errVPSNotFound):
			diags.AddError(fmt.Sprintf("VPS %s not found", id), err.Error())
		case vpsAPIStatus(err) == vpsV1StatusExpired:
			diags.AddError(fmt.Sprintf("VPS %s is expired", id), err.Error())
		default:
			diags.AddError(fmt.Sprintf("Error reading VPS %s", id), err.Error())
		}
		return nil, diags
	}
	return srv, diags
}

// List reads every VPS GET /vps returns, vpsV1ListParallelism at once, and
// returns them sorted by id, whatever the order of the API. A VPS deleted
// meanwhile is left out; an expired one too, with a warning.
func (b *vpsV1Backend) List(ctx context.Context) ([]VPSAPIServer, diag.Diagnostics) {
	var diags diag.Diagnostics

	var ids []string
	if err := b.client.GetWithContext(ctx, "/vps", &ids); err != nil {
		diags.AddError("Error calling GET /vps", err.Error())
		return nil, diags
	}

	type result struct {
		srv   *VPSAPIServer
		diags diag.Diagnostics
		err   error
	}
	results := make([]result, len(ids))
	slots := make(chan struct{}, vpsV1ListParallelism)
	var wg sync.WaitGroup
	for i, id := range ids {
		slots <- struct{}{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() { <-slots }()
			srv, readDiags, err := b.read(ctx, id)
			results[i] = result{srv: srv, diags: readDiags, err: err}
		}()
	}
	wg.Wait()

	servers := make([]VPSAPIServer, 0, len(ids))
	for i, id := range ids {
		srv, err := results[i].srv, results[i].err
		diags.Append(results[i].diags...)
		switch {
		case err == nil:
			servers = append(servers, *srv)
		case errors.Is(err, errVPSNotFound):
			// Deleted between GET /vps and its own GET.
		case vpsAPIStatus(err) == vpsV1StatusExpired:
			diags.AddWarning(
				fmt.Sprintf("VPS %s is expired", id),
				fmt.Sprintf("It is left out of the list: %s", err),
			)
		default:
			diags.AddError(fmt.Sprintf("Error reading VPS %s", id), err.Error())
			return nil, diags
		}
	}
	slices.SortFunc(servers, func(a, b VPSAPIServer) int { return cmp.Compare(a.Id, b.Id) })
	return servers, diags
}

// read reads one VPS. The diagnostics only hold warnings; err is the first
// failure, wrapping errVPSNotFound when GET /vps/{serviceName} answers 404,
// and errVPSExpired when it answers vpsV1StatusExpired.
func (b *vpsV1Backend) read(ctx context.Context, id string) (*VPSAPIServer, diag.Diagnostics, error) {
	var diags diag.Diagnostics
	base := "/vps/" + url.PathEscape(id)

	var vps vpsV1
	if err := b.client.GetWithContext(ctx, base, &vps); err != nil {
		if vpsAPIStatus(err) == http.StatusNotFound {
			return nil, diags, fmt.Errorf("%w: GET %s answered: %w", errVPSNotFound, base, err)
		}
		if vpsAPIStatus(err) == vpsV1StatusExpired {
			return nil, diags, fmt.Errorf("%w: GET %s answered: %w", errVPSExpired, base, err)
		}
		return nil, diags, fmt.Errorf("calling GET %s: %w", base, err)
	}

	from, err := b.installationFrom(ctx, base)
	if err != nil {
		return nil, diags, err
	}

	ipAddresses, err := b.ipAddresses(ctx, base)
	if err != nil {
		return nil, diags, err
	}

	region, regionDiags, err := b.region(ctx, id, base)
	diags.Append(regionDiags...)
	if err != nil {
		return nil, diags, err
	}

	options, targetAutoBackup, err := b.options(ctx, base)
	if err != nil {
		return nil, diags, err
	}

	targetPowerState, currentPowerState := vpsPowerStatesFromV1(vps.State)

	return &VPSAPIServer{
		Id: id,
		TargetSpec: VPSAPIReadTargetSpec{
			DisplayName:  vps.DisplayName,
			PowerState:   targetPowerState,
			Installation: VPSAPIReadInstallation{From: from},
			Options:      VPSAPIReadOptions{AutoBackup: targetAutoBackup},
		},
		CurrentState: &VPSAPICurrentState{
			DisplayName:  vps.DisplayName,
			PowerState:   currentPowerState,
			Installation: VPSAPICurrentInstallation{From: from},
			Location:     VPSAPILocation{Region: region},
			Model:        VPSAPIHardware{Vcpus: vps.Vcore, Ram: vps.MemoryLimit, Disk: vps.Model.Disk},
			IpAddresses:  ipAddresses,
			Options:      options,
		},
	}, diags, nil
}

// installationFrom returns the image the VPS is installed from, with null
// fields when it is unknown: GET /images/current answers 404 when API v1
// cannot trace the image back (the first image of the chain, such as a private
// image deleted after a restoration, no longer exists).
func (b *vpsV1Backend) installationFrom(ctx context.Context, base string) (VPSAPIInstallationFrom, error) {
	var image vpsV1Image
	endpoint := base + "/images/current"
	if err := b.client.GetWithContext(ctx, endpoint, &image); err != nil {
		if vpsAPIStatus(err) == http.StatusNotFound {
			return VPSAPIInstallationFrom{}, nil
		}
		return VPSAPIInstallationFrom{}, fmt.Errorf("calling GET %s: %w", endpoint, err)
	}
	return VPSAPIInstallationFrom{
		ImageId:   vpsNonEmpty(image.Id),
		ImageName: vpsNonEmpty(image.Name),
	}, nil
}

// ipAddresses returns the IP addresses of the VPS, IPv4 first, then sorted by
// address, whatever the order of the API.
func (b *vpsV1Backend) ipAddresses(ctx context.Context, base string) ([]VPSAPIIPAddress, error) {
	var ips []string
	endpoint := base + "/ips"
	if err := b.client.GetWithContext(ctx, endpoint, &ips); err != nil {
		return nil, fmt.Errorf("calling GET %s: %w", endpoint, err)
	}

	ipAddresses := make([]VPSAPIIPAddress, 0, len(ips))
	for _, ip := range ips {
		version, err := vpsIPVersion(ip)
		if err != nil {
			return nil, fmt.Errorf("reading GET %s: %w", endpoint, err)
		}
		ipAddresses = append(ipAddresses, VPSAPIIPAddress{Ip: ip, Version: version})
	}
	// The addresses were parsed by vpsIPVersion.
	slices.SortFunc(ipAddresses, func(a, b VPSAPIIPAddress) int {
		return cmp.Or(
			cmp.Compare(a.Version, b.Version),
			netip.MustParseAddr(a.Ip).Compare(netip.MustParseAddr(b.Ip)),
		)
	})
	return ipAddresses, nil
}

// region returns the region of the VPS, nil with a warning when the
// credentials are not allowed to call GET /v2/location (such as a key whose
// access rules only allow GET /vps/*) or when the datacenter code matches no
// region. Any other failure of GET /v2/location is an error.
func (b *vpsV1Backend) region(ctx context.Context, id, base string) (*string, diag.Diagnostics, error) {
	var diags diag.Diagnostics

	var datacenter vpsV1Datacenter
	endpoint := base + "/datacenter"
	if err := b.client.GetWithContext(ctx, endpoint, &datacenter); err != nil {
		return nil, diags, fmt.Errorf("calling GET %s: %w", endpoint, err)
	}

	locations, err := b.locations.get(ctx, b.client)
	if err != nil {
		if !vpsLocationsDenied(err) {
			return nil, diags, fmt.Errorf("calling GET /v2/location: %w", err)
		}
		// The same for every VPS, so that the diagnostics of a list hold it once.
		diags.AddWarning(
			"Cannot read the region of the VPS",
			fmt.Sprintf("Your credentials are not allowed to call GET /v2/location (%s): the region of the VPS is left empty. Allow GET /v2/location to read it.", err),
		)
		return nil, diags, nil
	}

	region, ok := vpsRegionFromDatacenterCode(datacenter.Name, locations)
	if !ok {
		diags.AddWarning(
			fmt.Sprintf("Unknown region for VPS %s", id),
			fmt.Sprintf("The datacenter code %q of the VPS matches no single region of GET /v2/location: its region is left empty.", datacenter.Name),
		)
		return nil, diags, nil
	}
	return &region, diags, nil
}

// options returns the options of the VPS, and the automated backup schedule
// of its target specification (nil when there is none).
//
// GET /option only lists the active options (manager/l1/cloud/VPSOption.pm:66
// keeps the ones whose ginapi state is "active"): a terminated option is never
// listed, so being listed means being subscribed.
func (b *vpsV1Backend) options(ctx context.Context, base string) (VPSAPIOptions, *VPSAPIAutoBackup, error) {
	var names []string
	endpoint := base + "/option"
	if err := b.client.GetWithContext(ctx, endpoint, &names); err != nil {
		return VPSAPIOptions{}, nil, fmt.Errorf("calling GET %s: %w", endpoint, err)
	}

	var options VPSAPIOptions
	for _, name := range names {
		switch name {
		case "automatedBackup":
			options.AutoBackup.Enabled = true
		case "snapshot":
			options.Snapshot.Enabled = true
		case "additionalDisk":
			options.AdditionalDisks.Enabled = true
		}
	}

	var target *VPSAPIAutoBackup
	if options.AutoBackup.Enabled {
		var autoBackup vpsV1AutomatedBackup
		endpoint := base + "/automatedBackup"
		err := b.client.GetWithContext(ctx, endpoint, &autoBackup)
		switch {
		case vpsAPIStatus(err) == http.StatusNotFound:
			// Subscribed but not set up yet: no schedule nor rotation.
		case err != nil:
			return VPSAPIOptions{}, nil, fmt.Errorf("calling GET %s: %w", endpoint, err)
		default:
			rotation := autoBackup.Rotation
			options.AutoBackup.Rotation = &rotation
			// A disabled automated backup has no schedule applied.
			if autoBackup.State == "enabled" && autoBackup.Schedule != nil {
				cron, err := vpsCronFromV1Schedule(*autoBackup.Schedule)
				if err != nil {
					return VPSAPIOptions{}, nil, fmt.Errorf("reading GET %s: %w", endpoint, err)
				}
				options.AutoBackup.Cron = &cron
				target = &VPSAPIAutoBackup{Cron: cron}
			}
		}
	}

	if options.AdditionalDisks.Enabled {
		size, err := b.additionalDiskSize(ctx, base)
		if err != nil {
			return VPSAPIOptions{}, nil, err
		}
		options.AdditionalDisks.Size = size
	}

	return options, target, nil
}

// additionalDiskSize returns the size of the additional disk of the VPS, nil
// when there is none. A VPS has at most one additional disk today: the first
// one found is the one.
func (b *vpsV1Backend) additionalDiskSize(ctx context.Context, base string) (*int64, error) {
	var ids []int64
	endpoint := base + "/disks"
	err := b.client.GetWithContext(ctx, endpoint, &ids)
	if vpsAPIStatus(err) == http.StatusNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("calling GET %s: %w", endpoint, err)
	}

	for _, id := range ids {
		var disk vpsV1Disk
		endpoint := base + "/disks/" + url.PathEscape(strconv.FormatInt(id, 10))
		err := b.client.GetWithContext(ctx, endpoint, &disk)
		if vpsAPIStatus(err) == http.StatusNotFound {
			// Removed since GET /disks.
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("calling GET %s: %w", endpoint, err)
		}
		if disk.Type == "additional" {
			size := disk.Size
			return &size, nil
		}
	}
	return nil, nil
}

// vpsNonEmpty returns nil for an empty string, a pointer to it otherwise.
func vpsNonEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// vpsAPIStatus returns the HTTP status of an API error, 0 for any other error
// or no error.
func vpsAPIStatus(err error) int {
	var apiErr *ovh.APIError
	if errors.As(err, &apiErr) {
		return apiErr.Code
	}
	return 0
}
