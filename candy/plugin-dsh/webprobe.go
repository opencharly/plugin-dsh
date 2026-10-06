package dsh

// webprobe.go holds the shared authenticated web-UI probe fragment used by BOTH
// the `dsh:` web-running verb (verb.go) and the `charly dsh status` command
// (control.go).
//
// `dsh web` authenticates at TWO independent layers, and a probe that checks only
// the first can pass while the web app is not actually serving its API:
//
//  1. the per-process LAUNCH TOKEN: the dsh candy's entrypoint persists it to
//     $DSH_HOME/web-token, and GET /?token=<t> (root only) answers 303 and sets a
//     signed cookie in the jar. A tokenless GET / returns 401.
//  2. the HOST/ORIGIN BROWSER-TRUST FENCE (dsh-client-connection's
//     api-request-trust): every /api/* call passes it BEFORE authentication. The
//     request Host must be loopback or a configured trustedHosts authority. A
//     REVERSE-PROXY deployment (tailscale serve, an SSH -L tunnel) presents a
//     non-loopback Host and therefore needs `dsh web --trusted-host <host:port>` or
//     every /api/* call returns 403 while / still loads — a HOST-SIDE requirement
//     this IN-BOX probe cannot itself observe (it always dials 127.0.0.1).
//
// So the probe does BOTH: it exchanges the token on / (capturing the cookie into a
// temp jar), then makes one /api/* call with that cookie. The exchange failure
// PROPAGATES (a bad/expired token or a down web UI exits non-zero, keeping the
// layer-1 assertion the old fragment had by virtue of being its last command). The
// /api/* call is then rejected on 403 (the fence), 000 (no response), and 5xx —
// everything else is a live API that answered. A missing/empty token file is a HARD
// error (the token is the precondition of the probe) — never a silent tokenless
// probe.
func dshWebTokenProbe(port string) string {
	return `T="$(cat "${DSH_HOME:-$HOME/.dsh}/web-token" 2>/dev/null)"; ` +
		`if [ -z "$T" ]; then echo "no dsh web token at ${DSH_HOME:-$HOME/.dsh}/web-token — the dsh entrypoint captures it on service start" >&2; exit 1; fi; ` +
		`J="$(mktemp)"; ` +
		`curl -fsS -o /dev/null -c "$J" "http://127.0.0.1:` + port + `/?token=$T" || { echo "dsh web token exchange failed on 127.0.0.1:` + port + ` (bad or expired launch token, or the web UI is down)" >&2; rm -f "$J"; exit 1; }; ` +
		`C="$(curl -s -o /dev/null -w '%{http_code}' -b "$J" "http://127.0.0.1:` + port + `/api/models")"; rm -f "$J"; ` +
		`case "$C" in 403) echo "dsh web UI returned 403 on /api/* — the Host/Origin browser-trust fence refused the request" >&2; exit 1;; 000) echo "dsh web UI gave no response on /api/*" >&2; exit 1;; 5??) echo "dsh web UI returned $C on /api/*" >&2; exit 1;; esac`
}
