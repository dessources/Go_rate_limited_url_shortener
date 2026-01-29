import { Loader2, AlertTriangle, TriangleRight, Play } from "lucide-react";
import { Progress } from "@/components/ui/progress";
import { Dispatch, SetStateAction, useState } from "react";
import { Button } from "@/components/ui/button";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { TabsContent } from "@/components/ui/tabs";
import StressTest from "./stress-test";
import LiveMetrics from "./live-metrics";

import type { Metrics } from "@/app/page";
import type { LiveMetricsProps } from "@/components/live-metrics";

export default function AdvancedTab(props: LiveMetricsProps & any) {
  return (
    <TabsContent value="advanced" className="mt-6 space-y-8">
      <LiveMetrics metrics={props.metrics} setMetrics={props.setMetrics} />
      <StressTest />
    </TabsContent>
  );
}
