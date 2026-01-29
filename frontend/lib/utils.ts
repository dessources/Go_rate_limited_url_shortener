import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";

export const MAIN_TEXT_LOGO = process.env.MAIN_TEXT_LOGO ?? "Pety.io";

export const REPO_LINK =
  process.env.REPO_LINK ?? "https://github.com/dessources/go_rate_limiter";

export const MAX_URL_LENGTH = process.env.MAX_URL_LENGTH
  ? parseInt(process.env.MAX_URL_LENGTH)
  : 4096;

export const ALLOWED_PROTOCOLS =
  process.env.ALLOWED_PROTOCOLS ?? "http:,https:";

export const BASE_URL =
  // process.env.NODE_ENV == "production" ? "" : "http://localhost:8090";
  process.env.NODE_ENV == "production" ? "" : "http://localhost:8090";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}
