package ovh

import (
	"fmt"
	"strings"
)

// CloudResourceTaskError is a single error reported on an asynchronous task.
type CloudResourceTaskError struct {
	Message string `json:"message,omitempty"`
}

// CloudResourceTask mirrors apiv2 model.ResourceTask: an asynchronous operation
// reconciling a Public Cloud API v2 resource towards its target spec. Used for
// wait/diagnostics only, never stored in Terraform state.
type CloudResourceTask struct {
	Id     string                   `json:"id,omitempty"`
	Type   string                   `json:"type,omitempty"`
	Status string                   `json:"status,omitempty"`
	Link   string                   `json:"link,omitempty"`
	Errors []CloudResourceTaskError `json:"errors,omitempty"`
}

// cloudResourceErrorFromTasks builds the error surfaced when a resource reaches
// the terminal ERROR status, joining the reasons reported by its failed tasks.
func cloudResourceErrorFromTasks(resourceKind, resourceId string, tasks []CloudResourceTask) error {
	var messages []string
	for _, task := range tasks {
		for _, taskErr := range task.Errors {
			if taskErr.Message != "" {
				messages = append(messages, taskErr.Message)
			}
		}
	}
	if len(messages) == 0 {
		return fmt.Errorf("%s %s entered ERROR state (no error detail reported by the API)", resourceKind, resourceId)
	}
	return fmt.Errorf("%s %s entered ERROR state: %s", resourceKind, resourceId, strings.Join(messages, "; "))
}
