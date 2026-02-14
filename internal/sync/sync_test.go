// Copyright (C) 2024 BuildBureau team
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sync

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

func TestErrGroup(t *testing.T) {
	t.Run("all succeed", func(t *testing.T) {
		ctx := context.Background()
		ctx, g := NewErrGroup(ctx)

		var count int32
		for i := 0; i < 10; i++ {
			g.Go(func() error {
				atomic.AddInt32(&count, 1)
				return nil
			})
		}

		if err := g.Wait(); err != nil {
			t.Errorf("expected no error, got: %v", err)
		}

		if count != 10 {
			t.Errorf("expected count=10, got: %d", count)
		}
	})

	t.Run("one error", func(t *testing.T) {
		ctx := context.Background()
		ctx, g := NewErrGroup(ctx)

		expectedErr := errors.New("test error")
		g.Go(func() error {
			return expectedErr
		})
		g.Go(func() error {
			time.Sleep(10 * time.Millisecond)
			return nil
		})

		err := g.Wait()
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !errors.Is(err, expectedErr) {
			t.Errorf("expected error %v, got: %v", expectedErr, err)
		}
	})

	t.Run("context cancellation", func(t *testing.T) {
		ctx := context.Background()
		ctx, g := NewErrGroup(ctx)

		g.Go(func() error {
			return errors.New("test error")
		})

		g.Go(func() error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(100 * time.Millisecond):
				return nil
			}
		})

		err := g.Wait()
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestOnce(t *testing.T) {
	var once Once
	var count int32

	for i := 0; i < 10; i++ {
		go once.Do(func() {
			atomic.AddInt32(&count, 1)
		})
	}

	time.Sleep(10 * time.Millisecond)

	if count != 1 {
		t.Errorf("expected count=1, got: %d", count)
	}
}

func TestPool(t *testing.T) {
	pool := NewPool(func() interface{} {
		return &struct{ value int }{value: 42}
	})

	item := pool.Get().(*struct{ value int })
	if item.value != 42 {
		t.Errorf("expected value=42, got: %d", item.value)
	}

	item.value = 100
	pool.Put(item)

	item2 := pool.Get().(*struct{ value int })
	if item2.value != 100 {
		t.Errorf("expected reused item with value=100, got: %d", item2.value)
	}
}

func TestMap(t *testing.T) {
	m := NewMap()

	t.Run("store and load", func(t *testing.T) {
		m.Store("key", "value")
		val, ok := m.Load("key")
		if !ok {
			t.Error("expected key to exist")
		}
		if val != "value" {
			t.Errorf("expected value='value', got: %v", val)
		}
	})

	t.Run("load or store", func(t *testing.T) {
		actual, loaded := m.LoadOrStore("key2", "value2")
		if loaded {
			t.Error("expected key to not exist")
		}
		if actual != "value2" {
			t.Errorf("expected value='value2', got: %v", actual)
		}

		actual, loaded = m.LoadOrStore("key2", "value3")
		if !loaded {
			t.Error("expected key to exist")
		}
		if actual != "value2" {
			t.Errorf("expected value='value2', got: %v", actual)
		}
	})

	t.Run("delete", func(t *testing.T) {
		m.Store("key3", "value3")
		m.Delete("key3")
		_, ok := m.Load("key3")
		if ok {
			t.Error("expected key to not exist after deletion")
		}
	})

	t.Run("range", func(t *testing.T) {
		testMap := NewMap()
		testMap.Store("a", 1)
		testMap.Store("b", 2)
		testMap.Store("c", 3)

		count := 0
		testMap.Range(func(key, value interface{}) bool {
			count++
			return true
		})

		if count != 3 {
			t.Errorf("expected count=3, got: %d", count)
		}
	})
}

func TestMutex(t *testing.T) {
	var mu Mutex
	var count int

	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			mu.Lock()
			count++
			mu.Unlock()
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}

	if count != 10 {
		t.Errorf("expected count=10, got: %d", count)
	}
}

func TestRWMutex(t *testing.T) {
	var mu RWMutex
	var count int

	// Test write lock
	mu.Lock()
	count++
	mu.Unlock()

	// Test read lock
	mu.RLock()
	_ = count
	mu.RUnlock()

	if count != 1 {
		t.Errorf("expected count=1, got: %d", count)
	}
}
