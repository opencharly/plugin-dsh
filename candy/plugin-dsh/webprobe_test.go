package dsh

import (
	"strings"
	"testing"
)

// TestDshWebTokenProbeAuthenticates locks in the two-layer probe shape both call
// sites (verbWebRunning, DshStatusCmd.Run) share. Layer 1: it MUST read the launch
// token file, exchange it on / with ?token=$T, and PROPAGATE that exchange's failure
// (a bad/expired token or a down UI must exit non-zero — the old fragment got this
// only because the exchange was its last command). Layer 2: it MUST then exercise an
// /api/* call and reject 403 (the Host/Origin fence), 000 (no response), and 5xx.
func TestDshWebTokenProbeAuthenticates(t *testing.T) {
	got := dshWebTokenProbe("3080")
	for _, want := range []string{
		`"${DSH_HOME:-$HOME/.dsh}/web-token"`,
		`http://127.0.0.1:3080/?token=$T`,
		`|| { echo "dsh web token exchange failed`, // layer-1 failure propagates
		`http://127.0.0.1:3080/api/models`,
		`403)`,
		`000)`,
		`5??)`,
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
