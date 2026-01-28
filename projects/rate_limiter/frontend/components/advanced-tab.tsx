import { Loader2, AlertTriangle } from "lucide-react";
import { Progress } from "@/components/ui/progress";
import { Dispatch, SetStateAction, useState } from "react";
import { Button } from "@/components/ui/button";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { TabsContent } from "@/components/ui/tabs";

interface ShortenerTabProps {
  url: string;
  setUrl: Dispatch<SetStateAction<string>>;
  handleShorten: () => void;
  loading: boolean;
  shortUrl: string;
  handleCopy: () => void;
  copied: boolean;
  error: string;
}

export default function AdvancedTab() {
  const [stressTestRunning, setStressTestRunning] = useState(false);
  const [stressTestOutput, setStressTestOutput] = useState(
    "Waiting for test to start...",
  );

  const handleStressTest = () => {
    setStressTestRunning(true);
    setStressTestOutput("Initializing stress test...\n");

    // Simulate stress test output
    const outputs = [
      "Starting stress test with 1000 concurrent requests...",
      "Spawning worker threads...",
      "Workers: [=========>] 100%",
      "Sending requests to /api/shorten...",
      "Progress: [==>       ] 25% (250/1000)",
      "Progress: [=====>    ] 50% (500/1000)",
      "Progress: [========> ] 75% (750/1000)",
      "Progress: [==========] 100% (1000/1000)",
      "",
      "Results:",
      "  Total Requests: 1000",
      "  Successful: 982",
      "  Failed: 18",
      "  Average Response Time: 45ms",
      "  Max Response Time: 312ms",
      "  Min Response Time: 12ms",
      "",
      "Rate Limiter Performance:",
      "  Global hits: 1000/1000",
      "  Per-client rejections: 18",
      "",
      "Test completed successfully!",
    ];

    let index = 0;
    const interval = setInterval(() => {
      if (index < outputs.length) {
        setStressTestOutput((prev) => prev + "\n" + outputs[index]);
        index++;
      } else {
        clearInterval(interval);
        setStressTestRunning(false);
      }
    }, 300);
  };

  // Mock metrics data (would come from API in real implementation)
  const [metrics] = useState({
    globalLimit: { current: 82, max: 100 },
    activeClients: 15,
    totalUrls: 1245,
  });

  return (
    <TabsContent value="advanced" className="mt-6 space-y-8">
      {/* Live Metrics Section */}
      <div>
        <h2 className="mb-4 text-2xl font-semibold text-foreground">
          Live Metrics
        </h2>
        <div className="grid gap-4 md:grid-cols-3">
          {/* Global Limiter Card */}
          <Card>
            <CardHeader>
              <CardTitle className="text-base text-center">
                Global Rate Limit
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-3">
              <Progress
                value={
                  (metrics.globalLimit.current / metrics.globalLimit.max) * 100
                }
              />
              <p className="mt-10 text-center text-sm text-muted-foreground">
                {metrics.globalLimit.current} / {metrics.globalLimit.max}{" "}
                requests
              </p>
            </CardContent>
          </Card>

          {/* Active Users Card */}
          <Card>
            <CardHeader>
              <CardTitle className="text-base text-center">
                Active Users
              </CardTitle>
            </CardHeader>
            <CardContent className="flex flex-col items-center justify-center">
              <p className="text-4xl font-bold text-primary">
                {metrics.activeClients}
              </p>
              <p className="mt-2 text-sm text-muted-foreground">
                Users in the last 30 minutes
              </p>
            </CardContent>
          </Card>

          {/* Total URLs Card */}
          <Card>
            <CardHeader>
              <CardTitle className="text-base text-center">
                Total URLs Stored
              </CardTitle>
            </CardHeader>
            <CardContent className="flex flex-col items-center justify-center">
              <p className="text-4xl font-bold text-primary">
                {metrics.totalUrls.toLocaleString()}
              </p>
              <p className="mt-2 text-sm text-muted-foreground">URLs created</p>
            </CardContent>
          </Card>
        </div>
      </div>

      {/* System Stress Test Section */}
      <div>
        <h2 className="mb-4 text-2xl font-semibold text-foreground">
          System Stress Test
        </h2>
        <Card>
          <CardContent className="pt-6">
            <div className="space-y-4">
              <Button
                variant="destructive"
                onClick={handleStressTest}
                disabled={stressTestRunning}
                className="w-full"
              >
                {stressTestRunning ? (
                  <>
                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                    Running Test...
                  </>
                ) : (
                  <>
                    <AlertTriangle className="mr-2 h-4 w-4" />
                    Run Stress Test
                  </>
                )}
              </Button>
              <p className="text-sm text-muted-foreground">
                This will simulate high traffic to the server. Metrics will
                update in real-time. The test takes about one minute to
                complete.
              </p>
              <Card className="bg-black">
                <CardHeader>
                  <CardTitle className="text-sm font-mono text-green-400">
                    Test Output
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <pre className="max-h-80 overflow-auto font-mono text-xs text-green-400/90">
                    {stressTestOutput}
                  </pre>
                </CardContent>
              </Card>
            </div>
          </CardContent>
        </Card>
      </div>
    </TabsContent>
  );
}
