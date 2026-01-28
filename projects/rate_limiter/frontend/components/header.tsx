import { MAIN_TEXT_LOGO, REPO_LINK } from "@/lib/utils";
import { siGithub } from "simple-icons";
import { ThemeSwitcher } from "./theme-switcher";

export default function Header() {
  return (
    <header className="border-b border-border bg-card">
      <div className="container mx-auto flex items-center justify-between px-4 py-4">
        <h1 className="text-2xl font-bold text-foreground">{MAIN_TEXT_LOGO}</h1>
        <div className="flex items-center gap-2">
          <a
            href={REPO_LINK}
            target="_blank"
            rel="noopener noreferrer"
            className="text-muted-foreground transition-colors hover:text-foreground"
          >
            <svg role="img" viewBox="0 0 24 24" className="h-5 w-5">
              <title>Github</title>
              <path d={siGithub.path}></path>
            </svg>

            <span className="sr-only">GitHub Repository</span>
          </a>
          <ThemeSwitcher />
        </div>
      </div>
    </header>
  );
}
