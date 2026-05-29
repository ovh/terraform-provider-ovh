package ovh

// VPSBackupFtpUnitAndValue represents a unit/value pair returned by the
// /vps/{serviceName}/backupftp endpoint for quota and usage.
type VPSBackupFtpUnitAndValue struct {
	Unit  string  `json:"unit"`
	Value float64 `json:"value"`
}

// VPSBackupFtp represents the response payload from
// GET /vps/{serviceName}/backupftp.
type VPSBackupFtp struct {
	FtpBackupName string                   `json:"ftpBackupName"`
	Quota         VPSBackupFtpUnitAndValue `json:"quota"`
	Usage         VPSBackupFtpUnitAndValue `json:"usage"`
	ReadOnlyDate  *string                  `json:"readOnlyDate,omitempty"`
	Type          string                   `json:"type"`
}

// VPSBackupFtpAcl represents the response payload from
// GET /vps/{serviceName}/backupftp/access/{ipBlock}.
type VPSBackupFtpAcl struct {
	IpBlock    string `json:"ipBlock"`
	Cifs       bool   `json:"cifs"`
	Ftp        bool   `json:"ftp"`
	Nfs        bool   `json:"nfs"`
	IsApplied  bool   `json:"isApplied"`
	LastUpdate string `json:"lastUpdate"`
}

// VPSBackupFtpAclCreateOpts is the body sent on
// POST /vps/{serviceName}/backupftp/access.
type VPSBackupFtpAclCreateOpts struct {
	IpBlock string `json:"ipBlock"`
	Cifs    bool   `json:"cifs"`
	Nfs     bool   `json:"nfs"`
	Ftp     *bool  `json:"ftp,omitempty"`
}

// VPSBackupFtpAclUpdateOpts is the body sent on
// PUT /vps/{serviceName}/backupftp/access/{ipBlock}.
type VPSBackupFtpAclUpdateOpts struct {
	IpBlock string `json:"ipBlock"`
	Cifs    bool   `json:"cifs"`
	Ftp     bool   `json:"ftp"`
	Nfs     bool   `json:"nfs"`
}
