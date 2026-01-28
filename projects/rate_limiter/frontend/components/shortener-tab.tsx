import { Loader2, CheckCheck, Copy, AlertTriangle } from "lucide-react";
import { Alert, AlertTitle, AlertDescription } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { TabsContent } from "@/components/ui/tabs";
import React, { Dispatch, JSX, SetStateAction, useState } from "react";

// interface ShortenerTabProps {
//   url: string;
//   setUrl: Dispatch<SetStateAction<string>>;
//   handleShorten: () => void;
//   loading: boolean;
//   shortUrl: string;
//   handleCopy: () => void;
//   copied: boolean;
//   error: string;
// }

export default function ShortenerTab() {
  const [url, setUrl] = useState("");
  const [loading, setLoading] = useState(false);
  const [shortUrl, setShortUrl] = useState("");
  const [error, setError] = useState("");
  const [copied, setCopied] = useState(false);

  const handleShorten = async () => {
    setLoading(true);
    setError("");
    setShortUrl("");

    // Simulate API call
    setTimeout(() => {
      if (!url || !url.startsWith("http")) {
        setError("Invalid URL provided");
        setLoading(false);
        return;
      }

      // Mock shortened URL
      const shortCode = Math.random().toString(36).substring(2, 8);
      setShortUrl(`http://localhost:8080/s/${shortCode}`);
      setLoading(false);
    }, 1000);
  };

  const handleCopy = async () => {
    await navigator.clipboard.writeText(shortUrl);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };
  return (
    <TabsContent value="shortener" className="mt-30">
      <Card className="mx-auto max-w-2xl">
        <CardHeader>
          <CardTitle>Paste a Long Link</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="flex gap-2">
            <Input
              type="url"
              placeholder="https://example.com/very/long/url/to/shorten"
              value={url}
              onChange={(e) => setUrl(e.target.value)}
              onKeyDown={(e) => e.key === "Enter" && handleShorten()}
              className="flex-1"
            />
            <Button onClick={handleShorten} disabled={loading}>
              {loading ? (
                <>
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  Shortening
                </>
              ) : (
                "Shorten"
              )}
            </Button>
          </div>

          {/* Results Section */}
          {shortUrl && (
            <div className="space-y-3">
              <Alert className="border-primary/50 bg-primary/10">
                <CheckCheck className="h-4 w-4 text-primary" />
                <AlertTitle className="text-primary">Success!</AlertTitle>
                <AlertDescription className="text-foreground">
                  Your URL has been shortened successfully.
                </AlertDescription>
              </Alert>
              <div className="flex gap-2">
                <Input value={shortUrl} readOnly className="flex-1" />
                <Button onClick={handleCopy} variant="outline">
                  {copied ? (
                    <>
                      <CheckCheck className="mr-2 h-4 w-4" />
                      Copied
                    </>
                  ) : (
                    <>
                      <Copy className="mr-2 h-4 w-4" />
                      Copy
                    </>
                  )}
                </Button>
              </div>
            </div>
          )}

          {error && (
            <Alert variant="destructive">
              <AlertTriangle className="h-4 w-4" />
              <AlertTitle>Error</AlertTitle>
              <AlertDescription>{error}</AlertDescription>
            </Alert>
          )}
        </CardContent>
      </Card>
    </TabsContent>
  );
}
