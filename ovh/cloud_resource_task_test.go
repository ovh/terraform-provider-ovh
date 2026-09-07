package ovh

import "testing"

// TestCloudResourceErrorFromTasks verifies the reasons reported by failed tasks
// are joined into the single error surfaced to the practitioner, and that an
// empty task list still yields an actionable message.
func TestCloudResourceErrorFromTasks(t *testing.T) {
	tasks := []CloudResourceTask{
		{Errors: []CloudResourceTaskError{{Message: "placement resolution failed"}}},
		{Errors: []CloudResourceTaskError{{Message: "quota exceeded"}, {Message: ""}}},
	}

	got := cloudResourceErrorFromTasks("private endpoint", "pe-1", tasks).Error()
	want := "private endpoint pe-1 entered ERROR state: placement resolution failed; quota exceeded"
	if got != want {
		t.Fatalf("unexpected error: got %q, want %q", got, want)
	}

	got = cloudResourceErrorFromTasks("private endpoint", "pe-1", nil).Error()
	want = "private endpoint pe-1 entered ERROR state (no error detail reported by the API)"
	if got != want {
		t.Fatalf("unexpected error: got %q, want %q", got, want)
	}
}
