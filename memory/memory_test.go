package memory

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestMemoryCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	m := New(ctx, 3*time.Hour)
	assert.NotNil(t, m.db)
	cancel()
	time.Sleep(10 * time.Millisecond)
	assert.Nil(t, m.db)
}

func TestMemoryTimeout(t *testing.T) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()
	m := New(ctx, 200*time.Millisecond)
	m.Set([]string{"a", "b"}, time.Now())
	_, ok := m.db["a b"]
	assert.True(t, ok)
	time.Sleep(300 * time.Millisecond)
	m.sync()
	_, ok = m.db["a b"]
	fmt.Println(m.db)
	assert.False(t, ok)
}

func TestMemoryKey(t *testing.T) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()
	m := New(ctx, 200*time.Millisecond)
	m.Set([]string{"a", "b"}, time.Now())
	assert.True(t, m.HasKey([]string{"a", "b"}))
	assert.False(t, m.HasKey([]string{"a", "a"}))
}
