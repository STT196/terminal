import { ulid } from "ulid";

export const prefixes = {
  user: "usr",
  apiClient: "cli",
  apiSecret: "sec",
  apiPersonal: "pat",
  project: "pjt",
} as const;

export function createID(prefix: keyof typeof prefixes): string {
  return [prefixes[prefix], ulid()].join("_");
}
