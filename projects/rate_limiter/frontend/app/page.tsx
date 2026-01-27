'use client'

import { useState } from 'react'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Progress } from '@/components/ui/progress'
import { Github, Loader2, Copy, CheckCheck, AlertTriangle } from 'lucide-react'
import { ThemeSwitcher } from '@/components/theme-switcher'

export default function Page() {
  const [url, setUrl] = useState('')
  const [loading, setLoading] = useState(false)
  const [shortUrl, setShortUrl] = useState('')
  const [error, setError] = useState('')
  const [copied, setCopied] = useState(false)
  const [stressTestRunning, setStressTestRunning] = useState(false)
  const [stressTestOutput, setStressTestOutput] = useState('Waiting for test to start...')

  // Mock metrics data (would come from API in real implementation)
  const [metrics] = useState({
    globalLimit: { current: 82, max: 100 },
    activeClients: 15,
    totalUrls: 1245,
  })

  const handleShorten = async () => {
    setLoading(true)
    setError('')
    setShortUrl('')

    // Simulate API call
    setTimeout(() => {
      if (!url || !url.startsWith('http')) {
        setError('Invalid URL provided')
        setLoading(false)
        return
      }

      // Mock shortened URL
      const shortCode = Math.random().toString(36).substring(2, 8)
      setShortUrl(`http://localhost:8080/s/${shortCode}`)
      setLoading(false)
    }, 1000)
  }

  const handleCopy = async () => {
    await navigator.clipboard.writeText(shortUrl)
    setCopied(true)
    setTimeout(() => setCopied(false), 2000)
  }

  const handleStressTest = () => {
    setStressTestRunning(true)
    setStressTestOutput('Initializing stress test...\n')

    // Simulate stress test output
    const outputs = [
      'Starting stress test with 1000 concurrent requests...',
      'Spawning worker threads...',
      'Workers: [=========>] 100%',
      'Sending requests to /api/shorten...',
      'Progress: [==>       ] 25% (250/1000)',
      'Progress: [=====>    ] 50% (500/1000)',
      'Progress: [========> ] 75% (750/1000)',
      'Progress: [==========] 100% (1000/1000)',
      '',
      'Results:',
      '  Total Requests: 1000',
      '  Successful: 982',
      '  Failed: 18',
      '  Average Response Time: 45ms',
      '  Max Response Time: 312ms',
      '  Min Response Time: 12ms',
      '',
      'Rate Limiter Performance:',
      '  Global hits: 1000/1000',
      '  Per-client rejections: 18',
      '',
      'Test completed successfully!',
    ]

    let index = 0
    const interval = setInterval(() => {
      if (index < outputs.length) {
        setStressTestOutput((prev) => prev + '\n' + outputs[index])
        index++
      } else {
        clearInterval(interval)
        setStressTestRunning(false)
      }
    }, 300)
  }

  return (
    <div className="min-h-screen bg-background">
      {/* Header */}
      <header className="border-b border-border bg-card">
        <div className="container mx-auto flex items-center justify-between px-4 py-4">
          <h1 className="text-2xl font-bold text-foreground">Go-Shortener</h1>
          <div className="flex items-center gap-2">
            <a
              href="https://github.com"
              target="_blank"
              rel="noopener noreferrer"
              className="text-muted-foreground transition-colors hover:text-foreground"
            >
              <Github className="h-5 w-5" />
              <span className="sr-only">GitHub Repository</span>
            </a>
            <ThemeSwitcher />
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="container mx-auto px-4 py-8">
        <Tabs defaultValue="shortener" className="mx-auto max-w-4xl">
          <TabsList className="grid w-full grid-cols-2">
            <TabsTrigger value="shortener">URL Shortener</TabsTrigger>
            <TabsTrigger value="advanced">Advanced</TabsTrigger>
          </TabsList>

          {/* URL Shortener Tab */}
          <TabsContent value="shortener" className="mt-6">
            <Card className="mx-auto max-w-2xl">
              <CardHeader>
                <CardTitle>Shorten a Long URL</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="flex gap-2">
                  <Input
                    type="url"
                    placeholder="https://example.com/very/long/url/to/shorten"
                    value={url}
                    onChange={(e) => setUrl(e.target.value)}
                    onKeyDown={(e) => e.key === 'Enter' && handleShorten()}
                    className="flex-1"
                  />
                  <Button onClick={handleShorten} disabled={loading}>
                    {loading ? (
                      <>
                        <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                        Shortening
                      </>
                    ) : (
                      'Shorten'
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

          {/* Advanced Tab */}
          <TabsContent value="advanced" className="mt-6 space-y-8">
            {/* Live Metrics Section */}
            <div>
              <h2 className="mb-4 text-2xl font-semibold text-foreground">Live Metrics</h2>
              <div className="grid gap-4 md:grid-cols-3">
                {/* Global Limiter Card */}
                <Card>
                  <CardHeader>
                    <CardTitle className="text-base">Global Rate Limit</CardTitle>
                  </CardHeader>
                  <CardContent className="space-y-3">
                    <Progress value={(metrics.globalLimit.current / metrics.globalLimit.max) * 100} />
                    <p className="text-center text-sm text-muted-foreground">
                      {metrics.globalLimit.current} / {metrics.globalLimit.max} requests
                    </p>
                  </CardContent>
                </Card>

                {/* Active Clients Card */}
                <Card>
                  <CardHeader>
                    <CardTitle className="text-base">Active Clients</CardTitle>
                  </CardHeader>
                  <CardContent className="flex flex-col items-center justify-center">
                    <p className="text-4xl font-bold text-primary">{metrics.activeClients}</p>
                    <p className="mt-2 text-sm text-muted-foreground">Clients in the last 30 minutes</p>
                  </CardContent>
                </Card>

                {/* Total URLs Card */}
                <Card>
                  <CardHeader>
                    <CardTitle className="text-base">Total URLs Stored</CardTitle>
                  </CardHeader>
                  <CardContent className="flex flex-col items-center justify-center">
                    <p className="text-4xl font-bold text-primary">{metrics.totalUrls.toLocaleString()}</p>
                    <p className="mt-2 text-sm text-muted-foreground">URLs created</p>
                  </CardContent>
                </Card>
              </div>
            </div>

            {/* System Stress Test Section */}
            <div>
              <h2 className="mb-4 text-2xl font-semibold text-foreground">System Stress Test</h2>
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
                      This will simulate high traffic to the server. Metrics will update in real-time. The test takes
                      about one minute to complete.
                    </p>
                    <Card className="bg-black/50">
                      <CardHeader>
                        <CardTitle className="text-sm font-mono text-green-400">Test Output</CardTitle>
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
        </Tabs>
      </main>
    </div>
  )
}
