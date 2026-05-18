# Exploration: gentle-engram-cloud-autotick

## Status

- SDD mode: auto
- Artifact store: OpenSpec
- Worktree: `../engram-gentle-engram-cloud-autotick`

## Question

Pi `gentle-engram` currently connects Pi sessions to Engram memory, but Cloud autosync only starts when the underlying Engram process sees `ENGRAM_CLOUD_AUTOSYNC=1`. The user wants Pi package parity with Engram Cloud's 30-second autosync loop.

## Evidence

### Pi plugin behavior

- `plugin/pi/index.ts` registers Pi-native `mem_*` tools and lifecycle handlers.
- `initOnce()` starts `engram serve` when `/health` is not reachable and `ENGRAM_URL` is not set.
- The spawned `engram serve` process inherits the parent environment, but the plugin does not derive or inject `ENGRAM_CLOUD_AUTOSYNC=1` from Cloud configuration.
- `initOnce()` runs one `engram sync --import` for git-synced chunks when `.engram/manifest.json` exists. This is not Cloud autosync.

### Engram core autosync behavior

- `internal/cloud/autosync/manager.go` owns durable background sync behavior.
- `DefaultConfig()` sets `PollInterval: 30 * time.Second`.
- `Run()` uses `time.NewTicker(m.cfg.PollInterval)` and runs `safeRun()` on each poll tick.
- This manager also supports dirty notification debounce, push/pull, lease, backoff, and status.

### Core process startup gates

- `cmd/engram/main.go::tryStartAutosync()` starts autosync only when `ENGRAM_CLOUD_AUTOSYNC` is exactly `"1"`.
- It also requires a Cloud server URL and `ENGRAM_CLOUD_TOKEN`.
- The function is intentionally non-fatal: missing config logs and leaves Engram running without autosync.
- `engram serve` and `engram mcp` call `tryStartAutosync()`.

### Pi MCP setup behavior

- `plugin/pi/cli.js` writes an MCP server launcher that runs `engram mcp --tools=agent`.
- The launcher currently inherits the Pi process environment but does not infer `ENGRAM_CLOUD_AUTOSYNC=1` from Cloud token/server configuration.

### Docs behavior

- `docs/AGENT-SETUP.md` documents the current core toggle: users must export `ENGRAM_CLOUD_AUTOSYNC=1`, `ENGRAM_CLOUD_TOKEN`, and `ENGRAM_CLOUD_SERVER`.
- `plugin/pi/README.md` says Cloud is opt-in and project-scoped, but does not document Pi package autosync parity.

## Findings

### What is missing in Pi parity today?

Pi does not automatically enable the existing Engram core autosync loop for Engram processes it launches/configures. Users who run Pi with Cloud token/server configured still need to know and export `ENGRAM_CLOUD_AUTOSYNC=1`; otherwise the 30-second core ticker never starts.

### Where should the 30-second tick live?

The tick must stay in `internal/cloud/autosync`. The Pi plugin should not implement its own polling or Cloud sync loop. `gentle-engram` should only provide thin process environment wiring so the existing Engram core process starts its own autosync manager.

### What gates preserve opt-in cloud sync?

Recommended plugin-side auto-enable gate:

1. Never override an explicit `ENGRAM_CLOUD_AUTOSYNC` value.
2. Require `ENGRAM_CLOUD_TOKEN` to be non-empty.
3. Require either `ENGRAM_CLOUD_SERVER` or persisted local `cloud.json.server_url`.
4. Let Engram core continue enforcing project enrollment, server policy, pause controls, lease, retry, and status.

This keeps local-first behavior because no token means no Cloud replication; configured token+server is a strong Cloud opt-in signal.

### Tests that should fail first

1. `plugin/pi/test/cloud-autosync-env.test.mjs`
   - auto-enables when token + `ENGRAM_CLOUD_SERVER` are present;
   - auto-enables when token + persisted `cloud.json.server_url` are present;
   - does not enable without token;
   - does not override explicit `ENGRAM_CLOUD_AUTOSYNC=0` or `1`.
2. `plugin/pi/test/index-source.test.mjs`
   - verifies `spawnDetached` uses the cloud-autosync env helper when starting Engram processes.
3. `plugin/pi/test/cli-source.test.mjs` or CLI helper tests
   - verifies generated MCP launcher includes equivalent auto-enable env logic for `engram mcp`.

### Docs to update

- `plugin/pi/README.md`: add Pi Cloud autosync parity section.
- `docs/AGENT-SETUP.md`: clarify that Pi `gentle-engram` can auto-enable autosync for its launched/configured Engram processes when token+server are present, while raw `engram serve`/`engram mcp` still support the explicit env toggle.

## Risks

| Risk | Mitigation |
|---|---|
| Plugin violates thin-adapter boundary | Only set child process environment; do not implement sync/polling in plugin. |
| Surprise Cloud sync | Require token + server and preserve explicit `ENGRAM_CLOUD_AUTOSYNC` overrides. Project enrollment remains authoritative. |
| Duplicate autosync workers | Existing SQLite lease in core manager prevents duplicate workers from syncing concurrently. |
| Noisy logs when Cloud partially configured | Gate on token+server before injecting autosync. |

## Recommendation

Implement a small Pi helper that builds child-process env for Engram processes. Use it for `engram serve`, generated `engram mcp` launcher, and any plugin-spawned one-shot sync command where harmless. The helper should only set `ENGRAM_CLOUD_AUTOSYNC=1` when the user has Cloud token + server configured and did not set an explicit autosync value.
