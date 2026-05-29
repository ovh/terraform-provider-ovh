package ovh

// VPSAutomatedBackup represents the configuration of automated backup on a VPS.
// Schema: vps.AutomatedBackup
type VPSAutomatedBackup struct {
	State               string `json:"state"`
	Schedule            string `json:"schedule"`
	Rotation            int64  `json:"rotation"`
	ServiceResourceName string `json:"serviceResourceName"`
}

// VPSAutomatedBackupAttachedAccess represents access information for an attached
// automated backup restore point. Schema: vps.automatedBackup.AttachedAccess
type VPSAutomatedBackupAttachedAccess struct {
	NFS            string `json:"nfs,omitempty"`
	SMB            string `json:"smb,omitempty"`
	AdditionalDisk string `json:"additionalDisk,omitempty"`
}

// VPSAutomatedBackupAttached represents a single attached restore point on a VPS.
// Schema: vps.automatedBackup.Attached
type VPSAutomatedBackupAttached struct {
	RestorePoint string                           `json:"restorePoint"`
	Access       VPSAutomatedBackupAttachedAccess `json:"access"`
}

// VPSAutomatedBackupRescheduleOpts is the payload used to reschedule the
// automated backup window on a VPS.
type VPSAutomatedBackupRescheduleOpts struct {
	Schedule string `json:"schedule"`
}

// VPSAutomatedBackupRestoreOpts is the payload used to trigger a restore from
// an automated backup restore point.
type VPSAutomatedBackupRestoreOpts struct {
	RestorePoint   string `json:"restorePoint"`
	Type           string `json:"type"`
	ChangePassword *bool  `json:"changePassword,omitempty"`
}

// VPSAutomatedBackupDetachOpts is the payload used to detach a previously
// attached automated backup restore point.
type VPSAutomatedBackupDetachOpts struct {
	RestorePoint string `json:"restorePoint"`
}
