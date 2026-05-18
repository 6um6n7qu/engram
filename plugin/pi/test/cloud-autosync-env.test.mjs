import test from "node:test";
import assert from "node:assert/strict";
import { mkdirSync, writeFileSync } from "node:fs";
import { join } from "node:path";
import { tmpdir } from "node:os";
import { mkdtempSync } from "node:fs";
import { cloudAutosyncProcessEnv } from "../cloud-autosync-env.js";

function tempDataDir() {
  return mkdtempSync(join(tmpdir(), "engram-cloud-autosync-env-"));
}

test("cloudAutosyncProcessEnv enables autosync with token and env server", () => {
  const env = {
    ENGRAM_CLOUD_TOKEN: "tok",
    ENGRAM_CLOUD_SERVER: "https://cloud.example.test",
  };

  const result = cloudAutosyncProcessEnv(env);

  assert.notEqual(result, env);
  assert.equal(result.ENGRAM_CLOUD_AUTOSYNC, "1");
  assert.equal(result.ENGRAM_CLOUD_TOKEN, "tok");
});

test("cloudAutosyncProcessEnv enables autosync with token and persisted server", () => {
  const dir = tempDataDir();
  writeFileSync(join(dir, "cloud.json"), JSON.stringify({ server_url: "https://cloud.example.test" }), "utf8");

  const result = cloudAutosyncProcessEnv({
    ENGRAM_DATA_DIR: dir,
    ENGRAM_CLOUD_TOKEN: "tok",
  });

  assert.equal(result.ENGRAM_CLOUD_AUTOSYNC, "1");
});

test("cloudAutosyncProcessEnv does not enable without token", () => {
  const env = { ENGRAM_CLOUD_SERVER: "https://cloud.example.test" };

  assert.equal(cloudAutosyncProcessEnv(env), env);
});

test("cloudAutosyncProcessEnv preserves explicit autosync override", () => {
  const env = {
    ENGRAM_CLOUD_AUTOSYNC: "0",
    ENGRAM_CLOUD_TOKEN: "tok",
    ENGRAM_CLOUD_SERVER: "https://cloud.example.test",
  };

  const result = cloudAutosyncProcessEnv(env);

  assert.equal(result, env);
  assert.equal(result.ENGRAM_CLOUD_AUTOSYNC, "0");
});

test("cloudAutosyncProcessEnv ignores malformed persisted config", () => {
  const dir = tempDataDir();
  writeFileSync(join(dir, "cloud.json"), "{not json", "utf8");
  const env = { ENGRAM_DATA_DIR: dir, ENGRAM_CLOUD_TOKEN: "tok" };

  assert.doesNotThrow(() => cloudAutosyncProcessEnv(env));
  assert.equal(cloudAutosyncProcessEnv(env), env);
});

test("cloudAutosyncProcessEnv ignores missing persisted config", () => {
  const dir = tempDataDir();
  mkdirSync(dir, { recursive: true });
  const env = { ENGRAM_DATA_DIR: dir, ENGRAM_CLOUD_TOKEN: "tok" };

  assert.equal(cloudAutosyncProcessEnv(env), env);
});
