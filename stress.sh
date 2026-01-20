#!/bin/bash

# High-traffic API stress test
# Target: 10,000 req/sec sustained with bursts up to 50,000

go run . &
PID=$!
sleep 3

echo "=== High Traffic API Stress Test ==="
echo "Config: 50,000 bucket capacity, 10,000 tokens/sec refill"
echo ""

# Sustained load: 10,000 req/sec (within refill rate)
echo "--- Test 1: Sustained load (10k req/sec) ---"
hey -n 50000 -c 200 -q 50 http://localhost:8090/ | grep responses

echo ""
echo "--- Test 2: Burst test (exceed capacity) ---"
# Burst: try to send 100k requests as fast as possible
hey -n 100000 -c 500 http://localhost:8090/ | grep responses

echo ""
echo "--- Test 3: Recovery after burst ---"
sleep 3
hey -n 10000 -c 100 http://localhost:8090/ | grep responses

echo ""
# Sustained load After Recovery: 10,000 req/sec
echo "--- Test 4: Sustained load after recovery (10k req/sec) ---"
sleep 10
hey -n 50000 -c 200 -q 50 http://localhost:8090/ | grep responses

curl -s http://localhost:8090/stop
wait $PID 2>/dev/null

