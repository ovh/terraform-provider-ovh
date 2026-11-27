package ovh

import "testing"

func TestVPSPowerStatesFromV1(t *testing.T) {
	cases := []struct {
		state   string
		target  string // "" means nil
		current string
	}{
		{"running", "ACTIVE", "ACTIVE"},
		{"stopped", "SHUTOFF", "SHUTOFF"},
		{"rescued", "RESCUE", "RESCUE"},
		{"backuping", "ACTIVE", "ACTIVE"},
		{"maintenance", "", "UNKNOWN"},
		{"installing", "", "UNKNOWN"},
		{"rebooting", "", "UNKNOWN"},
		{"upgrading", "", "UNKNOWN"},
		{"stopping", "", "UNKNOWN"},
		{"", "", "UNKNOWN"},
	}

	for _, c := range cases {
		target, current := vpsPowerStatesFromV1(c.state)

		if c.target == "" {
			if target != nil {
				t.Errorf("state %q: target = %q, want nil", c.state, *target)
			}
		} else if target == nil || *target != c.target {
			t.Errorf("state %q: target = %v, want %q", c.state, target, c.target)
		}

		if current != c.current {
			t.Errorf("state %q: current = %q, want %q", c.state, current, c.current)
		}
	}
}

func TestVPSCronFromV1Schedule(t *testing.T) {
	cases := map[string]string{
		"03:00:00": "0 3 * * *",
		"23:45:00": "45 23 * * *",
		"00:05:00": "5 0 * * *",
	}
	for schedule, want := range cases {
		got, err := vpsCronFromV1Schedule(schedule)
		if err != nil {
			t.Errorf("schedule %q: unexpected error: %s", schedule, err)
			continue
		}
		if got != want {
			t.Errorf("schedule %q: got %q, want %q", schedule, got, want)
		}
	}

	if got, err := vpsCronFromV1Schedule("3h"); err == nil {
		t.Errorf("schedule %q: expected an error, got %q", "3h", got)
	}
}

func TestVPSIPVersion(t *testing.T) {
	cases := map[string]int64{
		"51.75.1.2":    4,
		"2001:41d0::1": 6,
	}
	for ip, want := range cases {
		got, err := vpsIPVersion(ip)
		if err != nil {
			t.Errorf("ip %q: unexpected error: %s", ip, err)
			continue
		}
		if got != want {
			t.Errorf("ip %q: got %d, want %d", ip, got, want)
		}
	}

	if got, err := vpsIPVersion("nope"); err == nil {
		t.Errorf("ip %q: expected an error, got %d", "nope", got)
	}
}
