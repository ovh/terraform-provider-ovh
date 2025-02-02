package ovh

type DomainTask struct {
	CanAccelerate bool   `json:"canAccelerate"`
	CanCancel     bool   `json:"canCancel"`
	CanRelaunch   bool   `json:"canRelaunch"`
	Comment       string `json:"comment"`
	CreationDate  string `json:"creationDate"`
	Domain        string `json:"domain"`
	DoneDate      string `json:"doneDate"`
	Function      string `json:"function"`
	TaskId        int    `json:"id"`
	LastUpdate    string `json:"lastUpdate"`
	Status        string `json:"status"`
	TodoDate      string `json:"todoDate"`
}
