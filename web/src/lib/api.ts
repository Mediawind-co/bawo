import Client, { Local } from "./client";

const baseUrl = process.env.NEXT_PUBLIC_API_URL || Local;

export function createClient(token?: string) {
  return new Client(baseUrl, {
    auth: token ? { Authorization: `Bearer ${token}` } : undefined,
  });
}

export const publicClient = new Client(baseUrl);
