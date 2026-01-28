import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"


export const MAIN_TEXT_LOGO =  process.env.MAIN_TEXT_LOGO ?? "Pety.io"

export const REPO_LINK  = process.env.REPO_LINK ?? "https://github.com/dessources/go_rate_limiter"




export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}
