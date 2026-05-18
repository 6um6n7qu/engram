# Design: gentle-engram-cloud-autotick

## Summary

Add a tiny Pi-package environment helper and use it anywhere the Pi package starts/configures an Engram child process. The helper turns Cloud token+server configuration into `ENGRAM_CLOUD_AUTOSYNC=1` for the child only when the user has not set an explicit autosync value.

The 30-second tick stays in `internal/cloud/autosync.Manager`.

## Architecture

```text
Pi session
  └─ gentle-engram extension
      ├─ auto-starts `engram serve`
      │    └─ child env includes ENGRAM_CLOUD_AUTOSYNC=1 when token+server configured
      └─ pi-engram init writes MCP launcher
           └─ launcher starts `engram mcp --tools=agent` with same env gate

Engram core process
  └─ cmd/engram.tryStartAutosync
       └─ internal/cloud/autosync.Manager
            └─ 30s poll ticker + dirty debounce + push/pull + lease/backoff
```

## Files

| File | Action | Purpose |
|---|---|---|
| `plugin/pi/cloud-autosync-env.js` | NEW | Pure helper for child-process env derivation. |
| `plugin/pi/index.ts` | MODIFY | Use helper in `spawnDetached()` for `engram serve` / sync import child processes. |
| `plugin/pi/cli.js` | MODIFY | Generate MCP launcher with equivalent env auto-enable logic. |
| `plugin/pi/mcp-template.json` | MODIFY | Keep packaged template aligned with generated launcher. |
| `plugin/pi/package.json` | MODIFY | Include helper in package files. |
| `plugin/pi/test/cloud-autosync-env.test.mjs` | NEW | Helper unit tests. |
| `plugin/pi/test/index-source.test.mjs` | NEW | Source guard for helper usage in extension. |
| `plugin/pi/test/cli-source.test.mjs` | NEW | Source guard for launcher behavior/template alignment. |
| `plugin/pi/README.md` | MODIFY | Document Pi Cloud autosync parity. |
| `docs/AGENT-SETUP.md` | MODIFY | Clarify Pi package behavior vs raw CLI behavior. |

## Helper contract

```js
cloudAutosyncProcessEnv(baseEnv = process.env) -> env
```

Rules:

1. If `baseEnv.ENGRAM_CLOUD_AUTOSYNC` is non-blank, return `baseEnv` unchanged.
2. If `baseEnv.ENGRAM_CLOUD_TOKEN` is blank, return unchanged.
3. If `baseEnv.ENGRAM_CLOUD_SERVER` is non-blank, return a shallow copy with `ENGRAM_CLOUD_AUTOSYNC: "1"`.
4. Else inspect `${ENGRAM_DATA_DIR || ~/.engram}/cloud.json` for non-empty `server_url`; if present, return shallow copy with autosync `1`.
5. On absent/malformed config, return unchanged.

Returning unchanged by reference for no-op paths makes tests easy and avoids unnecessary env object copies. Returning a copy for enabled paths avoids mutating `process.env`.

## MCP launcher

The generated `node -e` launcher cannot import package files reliably, so it should embed equivalent minimal logic. This is acceptable because the launcher is already a generated adapter boundary. Keep it dependency-free and Node built-in only.

## Boundary rationale

- Plugin owns process setup, not Cloud sync behavior.
- Core `tryStartAutosync()` remains the final authority for startup validation.
- Core manager remains the only owner of tick interval, network push/pull, policy failure handling, and status.
- Server-side enrollment/pause remains authoritative.

## Validation

Targeted:

```sh
cd plugin/pi && npm test
```

Recommended broader check if time allows:

```sh
go test ./cmd/engram ./internal/cloud/autosync/...
```

No Go behavior should be changed by this SDD.
