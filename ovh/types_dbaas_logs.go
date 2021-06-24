package ovh

import ()

type DbaasLogsInputEngine struct {
	Id           string `json:"engineId"`
	IsDeprecated bool   `json:"isDeprecated"`
	Name         string `json:"name"`
	Version      string `json:"version"`
}

func (v DbaasLogsInputEngine) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["id"] = v.Id
	obj["is_deprecated"] = v.IsDeprecated
	obj["name"] = v.Name
	obj["version"] = v.Version
	return obj
}

type DbaasLogsOperation struct {
	AliasId     *string `json:"aliasId"`
	CreatedAt   string  `json:"createdAt"`
	DashboardId *string `json:"dashboardId"`
	IndexId     *string `json:"indexId"`
	InputId     *string `json:"inputId"`
	KibanaId    *string `json:"kibanaId"`
	OperationId string  `json:"operationId"`
	RoleId      *string `json:"roleId"`
	State       string  `json:"state"`
	StreamId    *string `json:"streamId"`
	UpdatedAt   string  `json:"updatedAt"`
}
