package ovh

import (
	"fmt"
	"net/netip"
	"time"
)

// Types of the API v2 contract of a VPS (public-schemas-v2,
// vps/models/objects/vps.server.*.json), JSON in camelCase.
//
// Only the fields API v1 can serve are present: the v2 envelope
// (resourceStatus, checksum, createdAt, updatedAt), backupId and sshKey are
// added with the API v2 backend.

// VPSAPIServer is vps.server.Server.
type VPSAPIServer struct {
	Id           string               `json:"id"`
	TargetSpec   VPSAPIReadTargetSpec `json:"targetSpec"`
	CurrentState *VPSAPICurrentState  `json:"currentState"`
}

// VPSAPIReadTargetSpec is vps.server.ReadTargetSpec.
type VPSAPIReadTargetSpec struct {
	DisplayName string `json:"displayName"`
	// PowerState is nil when it cannot be told (API v1 transitional state).
	PowerState   *string                `json:"powerState"`
	Installation VPSAPIReadInstallation `json:"installation"`
	Options      VPSAPIReadOptions      `json:"options"`
}

// VPSAPIReadInstallation is vps.server.ReadInstallation.
type VPSAPIReadInstallation struct {
	From VPSAPIInstallationFrom `json:"from"`
}

// VPSAPIInstallationFrom is vps.server.InstallationFrom and
// vps.server.CurrentInstallationFrom.
type VPSAPIInstallationFrom struct {
	ImageId   *string `json:"imageId"`
	ImageName *string `json:"imageName"`
}

// VPSAPIReadOptions is vps.server.ReadOptions.
type VPSAPIReadOptions struct {
	// AutoBackup is nil when you did not request an automated backup.
	AutoBackup *VPSAPIAutoBackup `json:"autoBackup"`
}

// VPSAPIAutoBackup is vps.server.AutoBackup.
type VPSAPIAutoBackup struct {
	Cron string `json:"cron"`
}

// VPSAPICurrentState is vps.server.CurrentState.
type VPSAPICurrentState struct {
	DisplayName string `json:"displayName"`
	// PowerState is a vps.server.CurrentPowerStateEnum: ACTIVE, RESCUE,
	// SHUTOFF or UNKNOWN.
	PowerState   string                    `json:"powerState"`
	Installation VPSAPICurrentInstallation `json:"installation"`
	Location     VPSAPILocation            `json:"location"`
	Model        VPSAPIHardware            `json:"model"`
	IpAddresses  []VPSAPIIPAddress         `json:"ipAddresses"`
	Options      VPSAPIOptions             `json:"options"`
}

// VPSAPICurrentInstallation is vps.server.CurrentInstallation.
type VPSAPICurrentInstallation struct {
	From VPSAPIInstallationFrom `json:"from"`
}

// VPSAPILocation is vps.server.Location.
type VPSAPILocation struct {
	// Region is nil when the datacenter of the VPS matches no known region.
	Region *string `json:"region"`
}

// VPSAPIHardware is vps.server.Model.
type VPSAPIHardware struct {
	Vcpus int64 `json:"vcpus"`
	Ram   int64 `json:"ram"`  // MB
	Disk  int64 `json:"disk"` // GB
}

// VPSAPIIPAddress is vps.server.IpAddress.
type VPSAPIIPAddress struct {
	Ip string `json:"ip"`
	// Version is a vps.server.IpVersionEnum: 4 or 6.
	Version int64 `json:"version"`
}

// VPSAPIOptions is vps.server.Options.
type VPSAPIOptions struct {
	AutoBackup      VPSAPIAutoBackupOption      `json:"autoBackup"`
	Snapshot        VPSAPISnapshotOption        `json:"snapshot"`
	AdditionalDisks VPSAPIAdditionalDisksOption `json:"additionalDisks"`
}

// VPSAPIAutoBackupOption is vps.server.AutoBackupOption.
type VPSAPIAutoBackupOption struct {
	Enabled  bool    `json:"enabled"`
	Cron     *string `json:"cron"`
	Rotation *int64  `json:"rotation"`
}

// VPSAPISnapshotOption is vps.server.SnapshotOption.
type VPSAPISnapshotOption struct {
	Enabled bool `json:"enabled"`
}

// VPSAPIAdditionalDisksOption is vps.server.AdditionalDisksOption.
type VPSAPIAdditionalDisksOption struct {
	Enabled bool   `json:"enabled"`
	Size    *int64 `json:"size"` // GB
}

// vpsPowerStatesFromV1 maps the API v1 state of a VPS to the requested power
// state (vps.server.PowerStateEnum) and the observed one
// (vps.server.CurrentPowerStateEnum).
//
// netbootMode is not read: it only applies at the next reboot, and a VPS booted
// on the rescue system already reads "rescued".
//
// API v1 keeps no target: the requested power state is the observed one, and
// nil when it cannot be told (transitional or maintenance state).
func vpsPowerStatesFromV1(state string) (target *string, current string) {
	switch state {
	case "running", "backuping":
		current = "ACTIVE"
	case "stopped":
		current = "SHUTOFF"
	case "rescued":
		current = "RESCUE"
	default:
		// installing, rebooting, upgrading, maintenance or unknown value
		return nil, "UNKNOWN"
	}

	requested := current
	return &requested, current
}

// vpsCronFromV1Schedule converts the API v1 automated backup schedule
// ("HH:MM:SS", UTC) to the cron expression of the API v2 contract
// ("M H * * *"). Seconds are dropped: the backup runs daily at a minute.
func vpsCronFromV1Schedule(schedule string) (string, error) {
	t, err := time.Parse(time.TimeOnly, schedule)
	if err != nil {
		return "", fmt.Errorf("invalid automated backup schedule %q, expected HH:MM:SS: %w", schedule, err)
	}
	return fmt.Sprintf("%d %d * * *", t.Minute(), t.Hour()), nil
}

// vpsIPVersion returns the vps.server.IpVersionEnum value of an IP address.
func vpsIPVersion(ip string) (int64, error) {
	addr, err := netip.ParseAddr(ip)
	if err != nil {
		return 0, fmt.Errorf("invalid IP address %q: %w", ip, err)
	}
	if addr.Is4() {
		return 4, nil
	}
	return 6, nil
}
