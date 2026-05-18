# Apply Progress: gentle-engram-cloud-autotick

## Implementation summary

Implemented Pi `gentle-engram` Cloud autosync parity by adding a child-process environment helper and wiring it into Pi-managed Engram processes.

## Changes

- Added `plugin/pi/cloud-autosync-env.js`.
  - Enables `ENGRAM_CLOUD_AUTOSYNC=1` only when token+server are configured and no explicit autosync env exists.
  - Supports `ENGRAM_CLOUD_SERVER` and persisted `${ENGRAM_DATA_DIR:-~/.engram}/cloud.json` `server_url`.
- Updated `plugin/pi/index.ts`.
  - `spawnDetached()` now passes helper-generated env to Engram child processes.
- Updated `plugin/pi/cli.js` and `plugin/pi/mcp-template.json`.
  - Generated MCP launcher now applies equivalent auto-enable behavior before running `engram mcp --tools=agent`.
  - `pi-engram init` now migrates only the known old generated Engram MCP launcher without `--force`; custom configs remain preserved.
- Added package file entry for the helper.
- Added Node tests for helper behavior, executable MCP launcher behavior, and safe init migration.
- Updated Pi/setup docs.

## TDD evidence

RED:

```text
cd plugin/pi && npm test
# failed: missing cloud-autosync-env.js, index.ts did not use helper, MCP launcher/template lacked ENGRAM_CLOUD_AUTOSYNC logic
# later failed: known old generated MCP launcher was kept instead of migrated
```

GREEN:

```text
cd plugin/pi && npm test
# pass: 28 tests
```

Additional syntax check:

```text
cd plugin/pi && node --check cli.js && node --check cloud-autosync-env.js
# pass
```

## Decisions

- Kept the 30-second tick in Engram core (`internal/cloud/autosync.Manager`).
- Treated token+server as the Pi-side Cloud opt-in signal.
- Preserved explicit `ENGRAM_CLOUD_AUTOSYNC` values as authoritative.
- Did not alter Go core sync behavior.
- Migrated only the exact old generated launcher shape to avoid overwriting custom user MCP configs.

## Pending

- Round 2 Judgment Day re-review results.
- Final verify report update.
