package dsh

// webprobe.go holds the shared authenticated web-UI probe fragment used by BOTH
// the `dsh:` web-running verb (verb.go) and the `charly dsh status` command
// (control.go). Since dsh-web-app 0.1.5-rc.x the web UI authenticates every
// request: each process mints a random launch token, the dsh candy's entrypoint
// persists it to $DSH_HOME/web-token, and a tokenless GET / returns 401. This
// fragment reads that token and curls the authenticated URL; a missing/empty
// token file is a HARD error (the token is the precondition of the probe) —
// never a silent tokenless probe. curl -f accepts the token exchange's 303
// (no -L needed) and fails on any non-2xx, so a zero exit means the
// authenticated web UI answered on the given loopback port.
func dshWebTokenProbe(port string) string {
	return `T="$(cat "${DSH_HOME:-$HOME/.dsh}/web-token" 2>/dev/null)"; ` +
		`if [ -z "$T" ]; then echo "no dsh web token at ${DSH_HOME:-$HOME/.dsh}/web-token — the dsh entrypoint captures it on service start" >&2; exit 1; fi; ` +
		`curl -fsS -o /dev/null "http://127.0.0.1:` + port + `/?token=$T"`
}
