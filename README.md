# plugin-dsh

Manage a deployed deepseek-harness (`dsh`) box from the host — the compiled-in
`charly dsh` management CLI plus the declarative `dsh:` check verb.

The plugin serves two capabilities over one provider: a `command:dsh` CLI
(`charly dsh status|profile|plugin …`) and a `verb:dsh` check verb for candy/box
plans. Both reach a **deployed** box over the reverse channel.

## What it provides

| Capability | Surface |
|---|---|
| `command:dsh` | `charly dsh status`, `charly dsh profile list`, `charly dsh plugin list` against a deployed box |
| `verb:dsh` | the `dsh:` check verb — an in-box probe of the built image or running deployment |

The command requires **compiled-in** placement: it drives the `pod-lifecycle`
host-builder's `op="cmd"` over the reverse channel, which is unavailable
out-of-process. The verb is placement-invisible (compiled-in or served over
go-plugin gRPC) and runs `dsh` **inside the venue** via `cc.Exec()`.

## How to use it

Compose the plugin candy in a box or check bed's `candy:` list:

```yaml
- '@github.com/opencharly/plugin-dsh/candy/plugin-dsh:<tag>'
```

The command then targets a deployed box:

```
charly dsh status --box my-dsh-box
charly dsh profile list --box my-dsh-box
charly dsh plugin list --box my-dsh-box --profile web
```

Every leaf takes `--box <name>` (required) and `--instance <name>`; the box's
stdout rides the host-held interactive leg straight to the operator's terminal,
and a non-zero exit in the box propagates as the command's exit code.

The `dsh:` check verb is authored in a candy/box `plan:`:

```yaml
- check: the dsh version is reported
  dsh: version
  context: [runtime]
```

| Method | Meaning |
|---|---|
| `version` | `dsh --version` in the venue |
| `web-running` | probe the token-authenticated web UI on `127.0.0.1:3080` in-box |
| `profile-list` | list profiles under `$DSH_HOME/profiles` |
| `plugin-list` | list a profile's plugins (`profile:`, default `web`) |

`web-running` is a live-service probe and skips under `charly check box` (no
running service on a disposable `podman run --rm`); the other methods probe the
built image and run in-box under both modes.

## Layout

- `candy/plugin-dsh/` — the plugin module: `plugin.go` (provider + `NewProvider()`/`NewMeta()`),
  `control.go` (the `charly dsh` command tree), `verb.go` (the `dsh:` verb),
  `webprobe.go` (the shared authenticated web-UI probe fragment),
  `schema/dsh.cue` (the self-contained `#DshInput`), `params/cue_types_gen.go`,
  `cmd/serve/main.go`.
- `charly.yml` — the root project manifest (`discover: candy`).
- `.github/workflows/tag-on-merge.yml` — CalVer tag + `CHANGELOG/` on merge.

## Related

- Owning skill: `/charly-tools:dsh-cli` — the `charly dsh` CLI reference (this
  candy's `skill:` entity).
- `/charly-tools:dsh` — the dsh candy (npm CLI + socat-exposed web service).
- `/charly-internals:plugin` — the plugin/provider model.
- [`opencharly/charly`](https://github.com/opencharly/charly) — the charly CLI.
