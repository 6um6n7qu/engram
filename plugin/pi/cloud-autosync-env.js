import { existsSync, readFileSync } from "node:fs";
import { homedir } from "node:os";
import { join } from "node:path";

function trimmed(value) {
  return typeof value === "string" ? value.trim() : "";
}

function cloudConfigPath(env) {
  return join(trimmed(env.ENGRAM_DATA_DIR) || join(homedir(), ".engram"), "cloud.json");
}

function hasPersistedCloudServer(env) {
  const path = cloudConfigPath(env);
  if (!existsSync(path)) return false;
  try {
    const parsed = JSON.parse(readFileSync(path, "utf8"));
    return trimmed(parsed?.server_url).length > 0;
  } catch {
    return false;
  }
}

function hasCloudServer(env) {
  return trimmed(env.ENGRAM_CLOUD_SERVER).length > 0 || hasPersistedCloudServer(env);
}

export function shouldAutoEnableCloudAutosync(env = process.env) {
  if (trimmed(env.ENGRAM_CLOUD_AUTOSYNC).length > 0) return false;
  if (trimmed(env.ENGRAM_CLOUD_TOKEN).length === 0) return false;
  return hasCloudServer(env);
}

export function cloudAutosyncProcessEnv(env = process.env) {
  if (!shouldAutoEnableCloudAutosync(env)) return env;
  return { ...env, ENGRAM_CLOUD_AUTOSYNC: "1" };
}
