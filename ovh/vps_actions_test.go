package ovh

import (
	"context"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

// vpsTestActions records the action requests of the fake API ("METHOD path",
// then the body when there is one), in the order they arrive.
type vpsTestActions struct {
	mu       sync.Mutex
	requests []string
}

func (a *vpsTestActions) add(request string, body []byte) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if len(body) > 0 {
		request += " " + string(body)
	}
	a.requests = append(a.requests, request)
}

func (a *vpsTestActions) list() []string {
	a.mu.Lock()
	defer a.mu.Unlock()
	return slices.Clone(a.requests)
}

// serveActions answers the action routes of the VPS: PUT /vps/{sn} (void),
// and POST start, stop, reboot, rebuild and automatedBackup/reschedule, each
// with a task done at once. It returns the record of these requests.
func (a *vpsTestAPI) serveActions() *vpsTestActions {
	actions := &vpsTestActions{}
	a.handle("PUT "+vpsTestBase, func(body []byte) any {
		actions.add("PUT "+vpsTestBase, body)
		return nil
	})
	taskId := int64(100)
	for _, action := range []struct{ path, taskType string }{
		{"/start", "startVm"},
		{"/stop", "stopVm"},
		{"/reboot", "rebootVm"},
		{"/rebuild", "reinstallVm"},
		{"/automatedBackup/reschedule", "rescheduleAutoBackup"},
	} {
		taskId++
		id := taskId
		a.handle("POST "+vpsTestBase+action.path, func(body []byte) any {
			actions.add("POST "+vpsTestBase+action.path, body)
			return a.serveTask(id, action.taskType, "done")
		})
	}
	return actions
}

// serveTask serves task id under GET /vps/{sn}/tasks/{id} in the given state,
// and returns it as a POST answers it (in state todo).
func (a *vpsTestAPI) serveTask(id int64, taskType, state string) map[string]any {
	a.set(vpsTestTaskPath(id), map[string]any{"id": id, "type": taskType, "state": state, "progress": 100, "date": "2026-11-27T10:00:00+01:00"})
	return map[string]any{"id": id, "type": taskType, "state": "todo", "progress": 0, "date": "2026-11-27T10:00:00+01:00"}
}

func vpsTestTaskPath(id int64) string {
	return vpsTestBase + "/tasks/" + strconv.FormatInt(id, 10)
}

// vpsTestClient starts the fake API with fast polls and returns its client.
func vpsTestClient(t *testing.T, api *vpsTestAPI) *Config {
	t.Helper()
	vpsTestFastPolls(t)
	return api.config(t)
}

func TestVPSPostActionRetriesConflict(t *testing.T) {
	t.Run("conflict then task", func(t *testing.T) {
		api := newVPSTestAPI()
		calls := 0
		api.handle("POST "+vpsTestBase+"/reboot", func([]byte) any {
			calls++
			if calls <= 2 {
				return vpsTestStatus{code: http.StatusConflict, message: "A task is already running"}
			}
			return map[string]any{"id": 42, "type": "rebootVm", "state": "todo"}
		})
		config := vpsTestClient(t, api)

		task, err := vpsPostAction(context.Background(), config.OVHClient, "/vps/"+vpsTestId+"/reboot", nil, time.Now().Add(time.Minute))
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if got := api.count("POST " + vpsTestBase + "/reboot"); got != 3 {
			t.Errorf("got %d calls, want 3", got)
		}
		if task.Id != 42 || task.Type != "rebootVm" || task.State != "todo" {
			t.Errorf("got task %+v", *task)
		}
	})

	t.Run("conflict until the deadline", func(t *testing.T) {
		api := newVPSTestAPI()
		api.handle("POST "+vpsTestBase+"/reboot", func([]byte) any {
			return vpsTestStatus{code: http.StatusConflict, message: "A task is already running"}
		})
		config := vpsTestClient(t, api)

		_, err := vpsPostAction(context.Background(), config.OVHClient, "/vps/"+vpsTestId+"/reboot", nil, time.Now().Add(50*time.Millisecond))
		if err == nil || !strings.Contains(err.Error(), "409") {
			t.Fatalf("got error %v, want one with 409", err)
		}
	})

	t.Run("other error", func(t *testing.T) {
		api := newVPSTestAPI()
		api.handle("POST "+vpsTestBase+"/reboot", func([]byte) any {
			return vpsTestStatus{code: http.StatusBadRequest, message: "Invalid"}
		})
		config := vpsTestClient(t, api)

		_, err := vpsPostAction(context.Background(), config.OVHClient, "/vps/"+vpsTestId+"/reboot", nil, time.Now().Add(time.Minute))
		if err == nil || !strings.Contains(err.Error(), "400") {
			t.Fatalf("got error %v, want one with 400", err)
		}
		if got := api.count("POST " + vpsTestBase + "/reboot"); got != 1 {
			t.Errorf("got %d calls, want 1", got)
		}
	})
}

func TestVPSWaitTask(t *testing.T) {
	task := &vpsV1Task{Id: 7, Type: "reinstallVm", State: "todo"}
	sequence := func(api *vpsTestAPI, answers ...any) {
		calls := 0
		api.handle("GET "+vpsTestTaskPath(7), func([]byte) any {
			answer := answers[min(calls, len(answers)-1)]
			calls++
			if state, ok := answer.(string); ok {
				return map[string]any{"id": 7, "type": "reinstallVm", "state": state}
			}
			return answer
		})
	}

	for name, tc := range map[string]struct {
		answers []any
		wantErr []string
	}{
		"doing then done": {answers: []any{"doing", "done"}},
		"not found":       {answers: []any{vpsTestStatus{code: http.StatusNotFound, message: "The requested object (id = 7) does not exist"}}},
		"error":           {answers: []any{"doing", "error"}, wantErr: []string{"reinstallVm", "error"}},
		"cancelled":       {answers: []any{"cancelled"}, wantErr: []string{"reinstallVm", "cancelled"}},
		"bad gateway":     {answers: []any{vpsTestStatus{code: http.StatusBadGateway, message: "Bad gateway"}, "done"}},
	} {
		t.Run(name, func(t *testing.T) {
			api := newVPSTestAPI()
			sequence(api, tc.answers...)
			config := vpsTestClient(t, api)

			err := vpsWaitTask(context.Background(), config.OVHClient, vpsTestId, task, time.Now().Add(time.Minute))
			if tc.wantErr == nil {
				if err != nil {
					t.Fatalf("unexpected error: %s", err)
				}
				if api.count("GET "+vpsTestBase) == 0 {
					t.Error("the state of the VPS was not checked after the task")
				}
				return
			}
			if err == nil {
				t.Fatalf("got no error, want one with %q", tc.wantErr)
			}
			for _, want := range tc.wantErr {
				if !strings.Contains(err.Error(), want) {
					t.Errorf("got error %q, want one with %q", err, want)
				}
			}
		})
	}
}

func TestVPSScheduleFromCron(t *testing.T) {
	for cron, want := range map[string]string{
		"0 2 * * *":   "02:00:00",
		"30 23 * * *": "23:30:00",
		"5 0 * * *":   "00:05:00",
	} {
		got, err := vpsScheduleFromCron(cron)
		if err != nil {
			t.Errorf("%q: unexpected error: %s", cron, err)
		}
		if got != want {
			t.Errorf("%q: got %q, want %q", cron, got, want)
		}
		// vpsCronFromV1Schedule is its inverse.
		if back, err := vpsCronFromV1Schedule(got); err != nil || back != cron {
			t.Errorf("%q: vpsCronFromV1Schedule(%q) = %q, %v", cron, got, back, err)
		}
	}
	for _, cron := range []string{"0 2 * * 1", "60 2 * * *", "0 24 * * *", "0 2 * *", "*/5 2 * * *", ""} {
		if got, err := vpsScheduleFromCron(cron); err == nil {
			t.Errorf("%q: got %q, want an error", cron, got)
		}
	}
}

func TestVPSImageIdForName(t *testing.T) {
	api := newVPSTestAPI()
	images := map[string]string{
		"id-debian-12": "Debian 12",
		"id-ubuntu-24": "Ubuntu 24.04",
		"id-alma-9":    "AlmaLinux 9",
	}
	ids := make([]string, 0, len(images))
	for id, name := range images {
		ids = append(ids, id)
		api.set(vpsTestBase+"/images/available/"+id, map[string]any{"id": id, "name": name})
	}
	api.set(vpsTestBase+"/images/available", ids)
	config := vpsTestClient(t, api)

	imageId, names, err := vpsImageIdForName(context.Background(), config.OVHClient, vpsTestId, "Ubuntu 24.04")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if imageId != "id-ubuntu-24" {
		t.Errorf("got image id %q, want id-ubuntu-24", imageId)
	}

	imageId, names, err = vpsImageIdForName(context.Background(), config.OVHClient, vpsTestId, "Windows")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if imageId != "" {
		t.Errorf("got image id %q, want none", imageId)
	}
	if want := []string{"AlmaLinux 9", "Debian 12", "Ubuntu 24.04"}; !slices.Equal(names, want) {
		t.Errorf("got names %q, want %q", names, want)
	}
}

func TestVPSSetPowerState(t *testing.T) {
	for _, tc := range []struct {
		state, netboot, target string
		want                   []string
	}{
		{"running", "local", "SHUTOFF", []string{"POST /stop"}},
		{"stopped", "local", "ACTIVE", []string{"POST /start"}},
		{"running", "local", "RESCUE", []string{`PUT {"netbootMode":"rescue"}`, "POST /reboot"}},
		{"stopped", "local", "RESCUE", []string{`PUT {"netbootMode":"rescue"}`, "POST /start"}},
		{"rescued", "rescue", "ACTIVE", []string{`PUT {"netbootMode":"local"}`, "POST /reboot"}},
		{"rescued", "rescue", "SHUTOFF", []string{`PUT {"netbootMode":"local"}`, "POST /stop"}},
		{"running", "local", "ACTIVE", nil},
		{"rescued", "local", "ACTIVE", []string{"POST /reboot"}},
		{"running", "rescue", "RESCUE", []string{"POST /reboot"}},
	} {
		t.Run(tc.state+" "+tc.netboot+" to "+tc.target, func(t *testing.T) {
			api := newVPSTestAPI()
			api.set(vpsTestBase, map[string]any{"name": vpsTestId, "displayName": "web-01", "state": tc.state, "netbootMode": tc.netboot})
			actions := api.serveActions()
			config := vpsTestClient(t, api)

			if err := vpsSetPowerState(context.Background(), config.OVHClient, vpsTestId, tc.target, time.Now().Add(time.Minute)); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			var got []string
			for _, request := range actions.list() {
				request = strings.Replace(request, vpsTestBase, "", 1)
				request = strings.Replace(request, "PUT  ", "PUT ", 1)
				got = append(got, request)
			}
			if !slices.Equal(got, tc.want) {
				t.Errorf("got requests %q, want %q", got, tc.want)
			}
		})
	}

	t.Run("unknown target", func(t *testing.T) {
		api := newVPSTestAPI()
		actions := api.serveActions()
		config := vpsTestClient(t, api)

		if err := vpsSetPowerState(context.Background(), config.OVHClient, vpsTestId, "UNKNOWN", time.Now().Add(time.Minute)); err == nil {
			t.Fatal("got no error for the target UNKNOWN")
		}
		if got := actions.list(); len(got) != 0 {
			t.Errorf("got requests %q, want none", got)
		}
	})
}
