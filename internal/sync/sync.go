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

// Package sync provides concurrency utilities
// inspired by vdaas/vald sync patterns
package sync

import (
	"context"
	"sync"
)

// ErrGroup is a collection of goroutines working on subtasks that are part of
// the same overall task. It's similar to sync.WaitGroup but with error handling
type ErrGroup interface {
	Go(f func() error)
	Wait() error
}

// errGroup implements ErrGroup
type errGroup struct {
	cancel  context.CancelFunc
	wg      sync.WaitGroup
	errOnce sync.Once
	err     error
}

// NewErrGroup creates a new error group
func NewErrGroup(ctx context.Context) (context.Context, ErrGroup) {
	ctx, cancel := context.WithCancel(ctx)
	return ctx, &errGroup{cancel: cancel}
}

// Go runs the given function in a new goroutine
func (g *errGroup) Go(f func() error) {
	g.wg.Add(1)

	go func() {
		defer g.wg.Done()

		if err := f(); err != nil {
			g.errOnce.Do(func() {
				g.err = err
				if g.cancel != nil {
					g.cancel()
				}
			})
		}
	}()
}

// Wait blocks until all goroutines have completed and returns the first error
func (g *errGroup) Wait() error {
	g.wg.Wait()
	if g.cancel != nil {
		g.cancel()
	}
	return g.err
}

// Once is a wrapper around sync.Once with additional functionality
type Once struct {
	once sync.Once
}

// Do calls the function f if and only if Do is being called for the
// first time for this instance of Once
func (o *Once) Do(f func()) {
	o.once.Do(f)
}

// Pool is a wrapper around sync.Pool
type Pool struct {
	pool sync.Pool
}

// NewPool creates a new pool with the given constructor function
func NewPool(new func() interface{}) *Pool {
	return &Pool{
		pool: sync.Pool{
			New: new,
		},
	}
}

// Get returns an item from the pool
func (p *Pool) Get() interface{} {
	return p.pool.Get()
}

// Put adds an item to the pool
func (p *Pool) Put(x interface{}) {
	p.pool.Put(x)
}

// Map is a concurrent safe map
type Map struct {
	m sync.Map
}

// NewMap creates a new concurrent safe map
func NewMap() *Map {
	return &Map{}
}

// Load returns the value stored in the map for a key
func (m *Map) Load(key interface{}) (value interface{}, ok bool) {
	return m.m.Load(key)
}

// Store sets the value for a key
func (m *Map) Store(key, value interface{}) {
	m.m.Store(key, value)
}

// LoadOrStore returns the existing value for the key if present
func (m *Map) LoadOrStore(key, value interface{}) (actual interface{}, loaded bool) {
	return m.m.LoadOrStore(key, value)
}

// Delete deletes the value for a key
func (m *Map) Delete(key interface{}) {
	m.m.Delete(key)
}

// Range calls f sequentially for each key and value present in the map
func (m *Map) Range(f func(key, value interface{}) bool) {
	m.m.Range(f)
}

// Mutex is a wrapper around sync.Mutex
type Mutex struct {
	mu sync.Mutex
}

// Lock locks the mutex
func (m *Mutex) Lock() {
	m.mu.Lock()
}

// Unlock unlocks the mutex
func (m *Mutex) Unlock() {
	m.mu.Unlock()
}

// RWMutex is a wrapper around sync.RWMutex
type RWMutex struct {
	mu sync.RWMutex
}

// Lock locks for writing
func (m *RWMutex) Lock() {
	m.mu.Lock()
}

// Unlock unlocks for writing
func (m *RWMutex) Unlock() {
	m.mu.Unlock()
}

// RLock locks for reading
func (m *RWMutex) RLock() {
	m.mu.RLock()
}

// RUnlock unlocks for reading
func (m *RWMutex) RUnlock() {
	m.mu.RUnlock()
}
