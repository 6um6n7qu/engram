import test from "node:test";
import assert from "node:assert/strict";
import { readFileSync } from "node:fs";
import { fileURLToPath } from "node:url";
import { dirname, join } from "node:path";

const __dirname = dirname(fileURLToPath(import.meta.url));
const source = readFileSync(join(__dirname, "..", "index.ts"), "utf8");

test("Pi extension uses cloud autosync env helper for Engram child processes", () => {
  assert.match(source, /cloudAutosyncProcessEnv/);
  assert.match(source, /env:\s*cloudAutosyncProcessEnv\(\)/);
});
