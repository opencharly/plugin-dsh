# AGENTS.md — plugin-dsh

Standalone plugin repo for the `dsh` capability (`command:dsh` + `verb:dsh`).
The plugin is a Go module at `candy/plugin-dsh/` (module path
`github.com/opencharly/plugin-dsh/candy/plugin-dsh`); the root `charly.yml` only
declares `discover: candy` so the repo is a project and its candy is scanned.

Canonical files:

- `candy/plugin-dsh/charly.yml` — the `plugin-dsh:` candy entity (`plugin:`
  block, `plan:` check) and the `dsh-cli-skill:` entity (the owning `skill:`).
- `candy/plugin-dsh/plugin.go` — the provider: `NewProvider()` + `NewMeta()`
  (a lazy `Describe` reflecting the kong CLIModel), command + verb dispatch.
- `candy/plugin-dsh/control.go` — the `charly dsh` command tree (the
  `pod-lifecycle` `op="cmd"` reverse leg).
- `candy/plugin-dsh/verb.go` — the in-box `dsh:` check verb.
- `candy/plugin-dsh/schema/dsh.cue` — the self-contained `#DshInput`.
- `.github/workflows/tag-on-merge.yml` — CalVer tag + `CHANGELOG/` on merge.
- `README.md` — user overview only; never agent guidance.

## Load these skills first (R0)

- `/charly-internals:plugin` — the plugin authoring reference: the `plugin:`
  block, the unified Provider model, the per-plugin CUE-schema contract,
  placement. Load before touching the provider or schema.
- `/charly-tools:dsh-cli` — the owning skill for the `charly dsh` CLI this
  candy projects.
- `/charly-tools:dsh` — the dsh candy the CLI manages.
- `/charly-internals:git-workflow` — before any git/PR action.

## Build / validate / test

- `go build ./...` in `candy/plugin-dsh/` — compile the plugin module.
- `go test ./...` in `candy/plugin-dsh/` — the plugin's Go tests
  (`webprobe_test.go`).
- `charly box validate` at the repo root — the structural check (the candy +
  `plugin:` block, CUE schema).
- The merge gate is the **org-wide** `charly/pr-validator` (required check
  `validate / validate`, defined in `opencharly/.github`); this repo has **no**
  per-repo candy gate.
- The full `charly dsh` CLI end-to-end is exercised by the Go e2e + the live
  R10 bed roster, not by the candy's build-context `plan:` check.

## Modify this repo

- Edit the `plugin-dsh:` candy entity, the Go source, and `schema/dsh.cue`
  **together** — the schema is the single source for the `params/` struct, so a
  field change not mirrored in the schema desyncs the generated types.
- The command is **compiled-in only** (the `pod-lifecycle` reverse channel is
  unavailable out-of-process); do not describe it as a general out-of-process
  command.

## Landing

- PR-only. Every change lands through a pull request; the org-required
  `charly/pr-validator` validates the diff and body and arms native auto-merge on
  PASS. Direct pushes to `main` are blocked.
- History lives in `CHANGELOG/` (written by `tag-on-merge` at merge time); the PR
  body IS the changelog.
- The authoritative rulebook is the umbrella `AGENTS.md` in
  `opencharly/opencharly` and `charly/AGENTS.md` in the charly repo. Do not
  restate its rules here.
