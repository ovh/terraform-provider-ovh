package ovh

// VpsVeeam represents the veeam configuration of a VPS.
type VpsVeeam struct {
	Backup bool `json:"backup"`
}

// VpsVeeamRestorePoint represents a single Veeam restore point.
type VpsVeeamRestorePoint struct {
	Id           int64  `json:"id"`
	CreationTime string `json:"creationTime"`
}

// VpsVeeamRestoredBackupAccessInfos represents the access information
// of a currently mounted Veeam restore.
type VpsVeeamRestoredBackupAccessInfos struct {
	Nfs string `json:"nfs"`
	Smb string `json:"smb"`
}

// VpsVeeamRestoredBackup represents the currently mounted Veeam
// restore on a VPS (if any).
type VpsVeeamRestoredBackup struct {
	RestorePointId int64                             `json:"restorePointId"`
	State          string                            `json:"state"`
	AccessInfos    VpsVeeamRestoredBackupAccessInfos `json:"accessInfos"`
}

// VpsVeeamRestoreOpts is the body of POST /vps/{sn}/veeam/restorePoints/{id}/restore.
type VpsVeeamRestoreOpts struct {
	Full           bool   `json:"full"`
	Export         string `json:"export"`
	ChangePassword *bool  `json:"changePassword,omitempty"`
}
