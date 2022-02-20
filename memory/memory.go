package memory

import (
	"context"
	"strings"
	"sync"
	"time"
)

type Memory struct {
	db     map[string]time.Time
	ttl    time.Duration
	maxKey int
	lock   *sync.RWMutex
}

func New(ctx context.Context, ttl time.Duration) *Memory {
	m := &Memory{
		db:     make(map[string]time.Time),
		ttl:    ttl,
		maxKey: 100,
		lock:   &sync.RWMutex{},
	}
	go m.garbage(ctx)
	return m
}

func (m *Memory) garbage(ctx context.Context) {
	t := time.NewTicker(5 * time.Minute)
	c := ctx.Done()
	for {
		select {
		case <-c:
			m.db = nil
			return
		case <-t.C:
			m.sync()
		}
	}
}

func (m *Memory) sync() {
	tombstones := make([]string, 0)
	death := time.Now().Add(-m.ttl)
	m.lock.Lock()
	for k, v := range m.db {
		if v.Before(death) {
			tombstones = append(tombstones, k)
		}
	}
	for _, tomb := range tombstones {
		delete(m.db, tomb)
	}
	m.lock.Unlock()
}

func (m *Memory) Set(keys []string, ts time.Time) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.db[strings.Join(keys, " ")] = ts
}

func (m *Memory) HasKey(keys []string) bool {
	m.lock.RLock()
	defer m.lock.RUnlock()
	ts, ok := m.db[strings.Join(keys, " ")]
	if !ok {
		return false
	}
	return ts.Add(-m.ttl).Before(time.Now())
}
