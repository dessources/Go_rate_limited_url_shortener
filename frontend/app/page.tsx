"use client";

import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs";
import ShortenerTab from "@/components/shortener-tab";
import AdvancedTab from "@/components/advanced-tab";
import Header from "@/components/header";

export default function Page() {
  return (
    <div className="min-h-screen bg-background">
      <Header />

      <main className="container mx-auto px-4 py-8">
        <Tabs defaultValue="shortener" className="mx-auto max-w-4xl">
          <TabsList className="grid w-full grid-cols-2">
            <TabsTrigger value="shortener">Link Shortener</TabsTrigger>
            <TabsTrigger value="advanced">Advanced</TabsTrigger>
          </TabsList>
          <ShortenerTab />
          <AdvancedTab />
        </Tabs>
      </main>
    </div>
  );
}
