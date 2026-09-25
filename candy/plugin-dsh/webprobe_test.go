package dsh

import (
	"strings"
	"testing"
)

// TestDshWebTokenProbeAuthenticates locks in the token-authenticated probe shape
// both call sites (verbWebRunning, DshStatusCmd.Run) share: it MUST read the
// launch token file, authenticate the URL with ?token=$T, and fail loudly on a
// missing token. A regression to the tokenless probe (which returns 401 on
// dsh-web-app >= 0.1.5-rc.x) fails this test.
func TestDshWebTokenProbeAuthenticates(t *testing.T) {
	got := dshWebTokenProbe("3080")
	for _, want := range []string{
		`"${DSH_HOME:-$HOME/.dsh}/web-token"`,
		`http://127.0.0.1:3080/?token=$T`,
		"no dsh web token",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("probe missing %q; got: %s", want, got)
		}
	}
	if strings.Contains(got, `127.0.0.1:3080/"`) {
		t.Errorf("probe must not issue a tokenless request; got: %s", got)
	}
}

// TestDshWebTokenProbePort proves the fragment targets the requested loopback
// port (3080 in-box, 3081 through the socat forwarder).
func TestDshWebTokenProbePort(t *testing.T) {
	if got := dshWebTokenProbe("3081"); !strings.Contains(got, "http://127.0.0.1:3081/?token=$T") {
		t.Errorf("probe must target the requested port; got: %s", got)
	}
}
