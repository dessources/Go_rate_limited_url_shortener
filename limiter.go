package main

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"time"
)

type LimiterBucketType int

const (
	InMemory LimiterBucketType = iota
	Redis
)

// --------------Token Store Definitions --------------
type TokenStore interface {
	Debit(count int) bool
	Count() int
	Capacity() int
	AddTokens(count int)
}

type MemoryBucket struct {
	count int
	cap   int
	mu    sync.Mutex
}

func NewMemoryBucket(count int, cap int) (*MemoryBucket, error) {

	if cap <= 0 {
		return nil, errors.New("Capacity must be a non-zero positive integer.")
	}

	if count < 0 {
		return nil, errors.New("count must be a non-negative integer if provided.")
	}

	if count > cap {
		return nil, errors.New("count must be less than or equal to capacity if provided.")
	}

	bucket := MemoryBucket{count: count, cap: cap}
	return &bucket, nil
}

func (b *MemoryBucket) AddTokens(count int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.count+count >= b.cap {
		b.count = b.cap
	} else {
		b.count += count
	}
}

func (b *MemoryBucket) Debit(count int) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.count >= count {
		b.count -= count
		return true
	}
	return false
}

func (b *MemoryBucket) Count() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.count
}

func (b *MemoryBucket) Capacity() int {
	return b.cap
}

// ----------------Limiter definition-----------
type Limiter struct {
	bucket TokenStore
	rate   int //token refill rate per second
	done   chan struct{}
}

func newLimiter(rate int, store TokenStore) *Limiter {

	done := make(chan struct{})
	limiter := Limiter{store, rate, done}
	go limiter.AddTokens()
	return &limiter
}

// func (l *Limiter) FillBucket() {
// 	elapsed := time.Since(l.lastRequestTime).Milliseconds()
// 	tokenCount := (elapsed * int64(l.rate)) / 1000
// 	l.bucket.AddTokens(int(tokenCount))
// }

func (l *Limiter) Allow(size int) bool {
	// l.FillBucket()
	if l.bucket.Debit(size) {
		return true
	}
	return false
}

func (l *Limiter) Stop() {
	close(l.done)
}

func (l *Limiter) AddTokens() {
	ticker := time.NewTicker(time.Second / time.Duration(l.rate))
	for {
		select {
		case <-ticker.C:
			l.bucket.AddTokens(1)
		case <-l.done:
			ticker.Stop()
			return
		}
	}
}

func NewLimiter(t LimiterBucketType, count, cap, rate int) *Limiter {
	var limiter *Limiter

	switch t {
	case InMemory:
		// initialize bucket
		bucket, err := NewMemoryBucket(count, cap)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		limiter = newLimiter(rate, bucket)
	}

	return limiter
}
