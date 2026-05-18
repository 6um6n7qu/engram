# Tasks: gentle-engram-cloud-autotick

## Review workload forecast

Estimated changed lines: 250-450 excluding SDD artifacts. Single PR is acceptable under 800-line budget.

## A. Tests first

- [x] A.1 **[RED]** Add `plugin/pi/test/cloud-autosync-env.test.mjs` covering:
  - env token+server enables autosync;
  - token+persisted `cloud.json.server_url` enables autosync;
  - missing token does not enable;
  - explicit `ENGRAM_CLOUD_AUTOSYNC=0` is preserved;
  - malformed/missing `cloud.json` does not throw.
- [x] A.2 **[RED]** Add `plugin/pi/test/index-source.test.mjs` asserting `index.ts` imports/uses `cloudAutosyncProcessEnv` for spawned Engram child processes.
- [x] A.3 **[RED]** Add `plugin/pi/test/cli-source.test.mjs` asserting generated launcher/template contain autosync gate and still launch `engram mcp --tools=agent`.
- [x] A.4 **[REVIEW FIX]** Strengthen MCP launcher tests to execute the generated launcher and assert cli/template parity.

## B. Implementation

- [x] B.1 Add `plugin/pi/cloud-autosync-env.js` helper with built-in `fs/path/os` only.
- [x] B.2 Update `plugin/pi/index.ts` to use helper for `spawnDetached()` child env.
- [x] B.3 Update `plugin/pi/cli.js` launcher string with equivalent env gate.
- [x] B.4 Update `plugin/pi/mcp-template.json` to match launcher behavior.
- [x] B.5 Add helper to `plugin/pi/package.json` `files` list.

## C. Docs

- [x] C.1 Update `plugin/pi/README.md` with Pi Cloud autosync parity section.
- [x] C.2 Update `docs/AGENT-SETUP.md` Cloud Autosync toggle section with Pi-specific behavior.
- [x] C.3 Clarify persisted config path as `${ENGRAM_DATA_DIR:-~/.engram}/cloud.json`.

## D. Verification

- [x] D.1 Run `cd plugin/pi && npm test`.
- [x] D.2 Run focused source/diff review.
- [x] D.3 Run fresh-context reviewer against final diff.
- [x] D.4 Capture verify report in `openspec/changes/gentle-engram-cloud-autotick/verify-report.md`.

## Acceptance

- [x] Pi-launched `engram serve` gets autosync env when token+server configured.
- [x] Pi-configured MCP launcher gets autosync env when token+server configured.
- [x] Explicit env override is respected.
- [x] No plugin-owned 30-second polling loop is introduced.
- [x] Tests and docs are updated.
