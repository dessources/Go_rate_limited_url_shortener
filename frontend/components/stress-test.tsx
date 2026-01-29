import { Loader2, Play } from "lucide-react";
import { Button } from "./ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "./ui/card";
import { useState } from "react";

export default function StressTest() {
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
  return (
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
              // className="w-full"
            >
              {stressTestRunning ? (
                <>
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  Running Test...
                </>
              ) : (
                <>
                  <Play className="mr-2 h-4 w-4" />
                  Run Stress Test
                </>
              )}
            </Button>
            <p className="text-sm text-muted-foreground">
              This will simulate high traffic to the server. Metrics will update
              in real-time. The test takes about one minute to complete.
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
  );
}
