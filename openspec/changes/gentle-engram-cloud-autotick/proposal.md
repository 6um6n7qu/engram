# Proposal: gentle-engram-cloud-autotick

## Intent

Give Pi `gentle-engram` parity with Engram Cloud autosync by auto-enabling the existing Engram core 30-second autosync manager for Engram processes launched or configured by the Pi package when Cloud is clearly configured.

The Pi plugin must not implement Cloud sync itself. It should remain a thin adapter that starts/configures Engram processes with the right environment so `internal/cloud/autosync.Manager` owns the actual push/pull tick, lease, backoff, policy handling, and status.

## Scope

### In Scope

- Add a small Pi package helper that decides whether child Engram processes should receive `ENGRAM_CLOUD_AUTOSYNC=1`.
- Use that helper when `plugin/pi/index.ts` auto-starts `engram serve`.
- Update the `pi-engram init` MCP launcher so generated `engram mcp --tools=agent` processes get equivalent auto-enable behavior.
- Update `plugin/pi/mcp-template.json` to match the generated launcher.
- Add deterministic Node tests for the helper and source-level integration guardrails.
- Update Pi and setup docs.

### Out of Scope

- Reimplementing the 30-second autosync loop in TypeScript.
- Changing `internal/cloud/autosync` tick interval or sync algorithm.
- Changing server-side project enrollment, pause, or auth semantics.
- Enabling Cloud sync without token+server configuration.
- Changing raw `engram serve` / `engram mcp` CLI semantics outside Pi package launchers.

## Behavioral Gate

The Pi package should inject `ENGRAM_CLOUD_AUTOSYNC=1` only when all are true:

1. `ENGRAM_CLOUD_AUTOSYNC` is unset/blank.
2. `ENGRAM_CLOUD_TOKEN` is non-empty.
3. A Cloud server is configured via either:
   - `ENGRAM_CLOUD_SERVER`, or
   - `${ENGRAM_DATA_DIR:-~/.engram}/cloud.json` with a non-empty `server_url`.

Explicit `ENGRAM_CLOUD_AUTOSYNC=0`, `false`, or `1` must be preserved.

## Success Criteria

- Pi-launched `engram serve` starts core autosync automatically when Cloud token+server are configured.
- Pi-configured `engram mcp --tools=agent` starts core autosync automatically under the same gate.
- No Cloud token means no auto-enable.
- Explicit autosync env remains authoritative.
- Tests cover env-server, persisted-server, missing-token, and explicit override cases.
- Docs tell users that the 30-second tick remains owned by Engram core.

## Risks and Mitigations

| Risk | Mitigation |
|---|---|
| Surprise Cloud replication | Require token+server and preserve explicit overrides; core enrollment still gates projects. |
| Plugin becomes thick | Only environment derivation lives in plugin; all sync behavior stays in Go core. |
| MCP launcher drift | Keep `cli.js` and `mcp-template.json` aligned and add source tests. |
| Partial cloud config causes noise | Do not inject autosync unless both token and server are present. |

## Review Workload Forecast

Estimated change: ~250-450 lines across helper, tests, launcher string, docs, and SDD artifacts. This is under the session budget of 800 changed lines and should fit a single PR.
