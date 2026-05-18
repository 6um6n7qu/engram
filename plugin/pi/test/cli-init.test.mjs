import test from "node:test";
import assert from "node:assert/strict";
import { mkdtempSync, readFileSync, writeFileSync } from "node:fs";
import { tmpdir } from "node:os";
import { join } from "node:path";
import { spawnSync } from "node:child_process";
import { fileURLToPath } from "node:url";
import { dirname } from "node:path";

const __dirname = dirname(fileURLToPath(import.meta.url));
const cliPath = join(__dirname, "..", "cli.js");
const oldLauncher = "const { spawn } = require('node:child_process'); const bin = process.env.ENGRAM_BIN || 'engram'; const child = spawn(bin, ['mcp', '--tools=agent'], { stdio: 'inherit' }); child.on('error', () => process.exit(127)); child.on('exit', (code, signal) => { if (typeof code === 'number') process.exit(code); process.kill(process.pid, signal || 'SIGTERM'); });";

function runInit(agentDir) {
  const result = spawnSync(process.execPath, [cliPath, "init"], {
    env: { ...process.env, PI_CODING_AGENT_DIR: agentDir },
    encoding: "utf8",
  });
  assert.equal(result.status, 0, `init failed: stdout=${result.stdout} stderr=${result.stderr}`);
  return result;
}

function readMcp(agentDir) {
  return JSON.parse(readFileSync(join(agentDir, "mcp.json"), "utf8"));
}

test("pi-engram init migrates the known old generated Engram MCP launcher", () => {
  const agentDir = mkdtempSync(join(tmpdir(), "pi-engram-init-"));
  writeFileSync(join(agentDir, "mcp.json"), JSON.stringify({
    mcpServers: {
      engram: {
        command: "node",
        args: ["-e", oldLauncher],
        lifecycle: "lazy",
        directTools: false,
      },
    },
  }, null, 2));

  const result = runInit(agentDir);
  const config = readMcp(agentDir);

  assert.match(result.stdout, /Migrated Engram MCP server/);
  assert.match(config.mcpServers.engram.args[1], /ENGRAM_CLOUD_AUTOSYNC/);
});

test("pi-engram init preserves custom Engram MCP server without force", () => {
  const agentDir = mkdtempSync(join(tmpdir(), "pi-engram-init-"));
  const custom = {
    command: "custom-engram-wrapper",
    args: ["mcp"],
    lifecycle: "lazy",
    directTools: false,
  };
  writeFileSync(join(agentDir, "mcp.json"), JSON.stringify({ mcpServers: { engram: custom } }, null, 2));

  const result = runInit(agentDir);
  const config = readMcp(agentDir);

  assert.match(result.stdout, /Kept existing Engram MCP server/);
  assert.deepEqual(config.mcpServers.engram, custom);
});
