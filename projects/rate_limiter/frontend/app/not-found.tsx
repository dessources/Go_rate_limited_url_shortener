import Link from "next/link";
import { Button } from "@/components/ui/button";
import { Home, ArrowLeft, Link2Off } from "lucide-react";
import Header from "@/components/header";

export default function NotFound() {
  return (
    <div className="max-h-screen">
      <Header></Header>
      <div className="flex mt-8 flex-col items-center justify-center bg-background px-4">
        <div className="mx-auto flex max-w-2xl flex-col items-center text-center">
          {/* Broken Link Illustration */}
          <div className="mb-8">
            <Link2Off width={200} height={200} className="rounded-2xl" />
          </div>

          {/* 404 Header */}
          <div className="mb-4">
            <h1 className="mb-2 text-7xl font-bold text-foreground">404</h1>
            <h2 className="text-3xl font-semibold text-foreground">
              Link Not Found
            </h2>
          </div>

          {/* Message */}
          <p className="mb-8 max-w-md text-pretty text-lg leading-relaxed text-muted-foreground">
            {"Looks like this short link has gone missing or never existed."}
          </p>

          {/* Action Buttons */}
          <div className="flex flex-col gap-3 sm:flex-row">
            <Button asChild size="lg" className="gap-2">
              <Link href="/">
                <Home className="h-5 w-5" />
                Go to Homepage
              </Link>
            </Button>
            <Button
              asChild
              variant="outline"
              size="lg"
              className="gap-2 bg-transparent"
            >
              <Link href="/">
                <ArrowLeft className="h-5 w-5" />
                Create a New Link
              </Link>
            </Button>
          </div>

          {/* Additional Help Text */}
          <p className="mt-8 text-sm text-muted-foreground">
            {"Need help? Make sure you copied the full short URL correctly."}
          </p>
        </div>
      </div>
    </div>
  );
}
