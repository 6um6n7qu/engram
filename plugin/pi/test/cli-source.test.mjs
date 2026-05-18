import test from "node:test";
import assert from "node:assert/strict";
import { chmodSync, mkdtempSync, readFileSync, writeFileSync } from "node:fs";
import { tmpdir } from "node:os";
import { dirname, join } from "node:path";
import { spawnSync } from "node:child_process";
import { fileURLToPath } from "node:url";

const __dirname = dirname(fileURLToPath(import.meta.url));
const cliSource = readFileSync(join(__dirname, "..", "cli.js"), "utf8");
const template = JSON.parse(readFileSync(join(__dirname, "..", "mcp-template.json"), "utf8"));

function cliLauncher() {
  const match = cliSource.match(/const MCP_LAUNCHER =\n  ("[\s\S]*?");/);
  assert.ok(match, "MCP_LAUNCHER string should be parseable from cli.js");
  return JSON.parse(match[1]);
}

function templateLauncher() {
  return template.mcpServers.engram.args[1];
}

function runLauncher(launcher, env) {
  const dir = mkdtempSync(join(tmpdir(), "engram-mcp-launcher-"));
  const envFile = join(dir, "env.json");
  const fakeEngram = join(dir, "fake-engram.sh");
  writeFileSync(fakeEngram, `#!/bin/sh\nnode -e 'require("node:fs").writeFileSync(process.argv[1], JSON.stringify({ autosync: process.env.ENGRAM_CLOUD_AUTOSYNC, args: process.argv.slice(2) }))' ${JSON.stringify(envFile)} "$@"\n`, "utf8");
  chmodSync(fakeEngram, 0o755);

  const result = spawnSync(process.execPath, ["-e", launcher], {
    env: { ...process.env, ...env, ENGRAM_BIN: fakeEngram },
    encoding: "utf8",
  });

  assert.equal(result.status, 0, `launcher failed: stdout=${result.stdout} stderr=${result.stderr}`);
  return JSON.parse(readFileSync(envFile, "utf8"));
}

test("pi-engram MCP launcher matches packaged template", () => {
  assert.equal(cliLauncher(), templateLauncher());
});

test("MCP launcher auto-enables Cloud autosync when env token and server are configured", () => {
  const result = runLauncher(cliLauncher(), {
    ENGRAM_CLOUD_AUTOSYNC: "",
    ENGRAM_CLOUD_TOKEN: "tok",
    ENGRAM_CLOUD_SERVER: "https://cloud.example.test",
  });

  assert.equal(result.autosync, "1");
  assert.deepEqual(result.args, ["mcp", "--tools=agent"]);
});

test("MCP launcher preserves explicit autosync override", () => {
  const result = runLauncher(cliLauncher(), {
    ENGRAM_CLOUD_AUTOSYNC: "0",
    ENGRAM_CLOUD_TOKEN: "tok",
    ENGRAM_CLOUD_SERVER: "https://cloud.example.test",
  });

  assert.equal(result.autosync, "0");
});

test("MCP launcher does not auto-enable without token", () => {
  const result = runLauncher(cliLauncher(), {
    ENGRAM_CLOUD_AUTOSYNC: "",
    ENGRAM_CLOUD_TOKEN: "",
    ENGRAM_CLOUD_SERVER: "https://cloud.example.test",
  });

  assert.equal(result.autosync, "");
});

test("MCP launcher auto-enables from persisted cloud.json server", () => {
  const dir = mkdtempSync(join(tmpdir(), "engram-mcp-cloud-config-"));
  writeFileSync(join(dir, "cloud.json"), JSON.stringify({ server_url: "https://cloud.example.test" }), "utf8");

  const result = runLauncher(cliLauncher(), {
    ENGRAM_CLOUD_AUTOSYNC: "",
    ENGRAM_CLOUD_TOKEN: "tok",
    ENGRAM_CLOUD_SERVER: "",
    ENGRAM_DATA_DIR: dir,
  });

  assert.equal(result.autosync, "1");
});
