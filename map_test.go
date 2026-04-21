package gsyncmap_test

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/min0625/gsyncmap"
)

// --- Map tests ---

func TestMap_StoreAndLoad(t *testing.T) {
	t.Parallel()

	var m gsyncmap.Map[string, int]

	// Load non-existing key returns zero value and false.
	v, ok := m.Load("missing")
	if ok || v != 0 {
		t.Fatalf("Load missing key: got (%v, %v), want (0, false)", v, ok)
	}

	m.Store("k", 42)

	v, ok = m.Load("k")
	if !ok || v != 42 {
		t.Fatalf("Load after Store: got (%v, %v), want (42, true)", v, ok)
	}
}

func TestMap_StoreOverwrite(t *testing.T) {
	t.Parallel()

	var m gsyncmap.Map[string, int]

	m.Store("k", 1)
	m.Store("k", 2)

	v, ok := m.Load("k")
	if !ok || v != 2 {
		t.Fatalf("Load after overwrite: got (%v, %v), want (2, true)", v, ok)
	}
}

func TestMap_Delete(t *testing.T) {
	t.Parallel()

	var m gsyncmap.Map[string, int]

	// Delete non-existing key should not panic.
	m.Delete("missing")

	m.Store("k", 1)
	m.Delete("k")

	_, ok := m.Load("k")
	if ok {
		t.Fatal("expected key to be deleted")
	}
}

func TestMap_LoadOrStore(t *testing.T) {
	t.Parallel()

	var m gsyncmap.Map[string, int]

	// Key absent: stores and returns the given value, loaded=false.
	actual, loaded := m.LoadOrStore("k", 10)
	if loaded || actual != 10 {
		t.Fatalf("LoadOrStore new key: got (%v, %v), want (10, false)", actual, loaded)
	}

	// Key present: returns the existing value, loaded=true.
	actual, loaded = m.LoadOrStore("k", 99)
	if !loaded || actual != 10 {
		t.Fatalf("LoadOrStore existing key: got (%v, %v), want (10, true)", actual, loaded)
	}
}

func TestMap_LoadAndDelete(t *testing.T) {
	t.Parallel()

	var m gsyncmap.Map[string, int]

	// Key absent: returns zero value and false.
	v, loaded := m.LoadAndDelete("missing")
	if loaded || v != 0 {
		t.Fatalf("LoadAndDelete missing: got (%v, %v), want (0, false)", v, loaded)
	}

	m.Store("k", 7)

	v, loaded = m.LoadAndDelete("k")
	if !loaded || v != 7 {
		t.Fatalf("LoadAndDelete existing: got (%v, %v), want (7, true)", v, loaded)
	}

	_, ok := m.Load("k")
	if ok {
		t.Fatal("key should be absent after LoadAndDelete")
	}
}

func TestMap_Swap(t *testing.T) {
	t.Parallel()

	var m gsyncmap.Map[string, int]

	// Key absent: previous is zero, loaded=false.
	prev, loaded := m.Swap("k", 1)
	if loaded || prev != 0 {
		t.Fatalf("Swap missing key: got (%v, %v), want (0, false)", prev, loaded)
	}

	// Key present: previous is old value, loaded=true.
	prev, loaded = m.Swap("k", 2)
	if !loaded || prev != 1 {
		t.Fatalf("Swap existing key: got (%v, %v), want (1, true)", prev, loaded)
	}

	v, _ := m.Load("k")
	if v != 2 {
		t.Fatalf("Load after Swap: got %v, want 2", v)
	}
}

func TestMap_Clear(t *testing.T) {
	t.Parallel()

	var m gsyncmap.Map[string, int]

	m.Store("a", 1)
	m.Store("b", 2)
	m.Clear()

	var count int

	m.Range(func(_ string, _ int) bool {
		count++
		return true
	})

	if count != 0 {
		t.Fatalf("expected empty map after Clear, got %d entries", count)
	}
}

func TestMap_Range(t *testing.T) {
	t.Parallel()

	var m gsyncmap.Map[string, int]

	want := map[string]int{"a": 1, "b": 2, "c": 3}
	for k, v := range want {
		m.Store(k, v)
	}

	got := make(map[string]int)

	m.Range(func(k string, v int) bool {
		got[k] = v
		return true
	})

	if len(got) != len(want) {
		t.Fatalf("Range visited %d entries, want %d", len(got), len(want))
	}

	for k, v := range want {
		if got[k] != v {
			t.Errorf("Range: key %q got %v, want %v", k, got[k], v)
		}
	}
}

func TestMap_Range_EmptyMap(t *testing.T) {
	t.Parallel()

	var m gsyncmap.Map[string, int]

	var visited int

	m.Range(func(_ string, _ int) bool {
		visited++
		return true
	})

	if visited != 0 {
		t.Fatalf("Range on empty map should not call f, called %d times", visited)
	}
}

func TestMap_Range_EarlyStop(t *testing.T) {
	t.Parallel()

	var m gsyncmap.Map[string, int]
	m.Store("a", 1)
	m.Store("b", 2)
	m.Store("c", 3)

	var visited int

	m.Range(func(_ string, _ int) bool {
		visited++
		return false
	})

	if visited != 1 {
		t.Fatalf("Range with early stop should visit exactly 1 entry, got %d", visited)
	}
}

func TestMap_CompareAndDelete_Deprecated(t *testing.T) {
	t.Parallel()

	var m gsyncmap.Map[string, string]

	m.Store("k", "v1")

	// Wrong value: should not delete.
	if m.CompareAndDelete("k", "wrong") {
		t.Fatal("CompareAndDelete with wrong value should return false")
	}

	v, ok := m.Load("k")
	if !ok || v != "v1" {
		t.Fatalf("key should still exist: got (%v, %v)", v, ok)
	}

	// Correct value: should delete.
	if !m.CompareAndDelete("k", "v1") {
		t.Fatal("CompareAndDelete with correct value should return true")
	}

	_, ok = m.Load("k")
	if ok {
		t.Fatal("key should be deleted")
	}

	// Key absent: should return false.
	if m.CompareAndDelete("k", "v1") {
		t.Fatal("CompareAndDelete on missing key should return false")
	}
}

func TestMap_CompareAndSwap_Deprecated(t *testing.T) {
	t.Parallel()

	var m gsyncmap.Map[string, string]

	m.Store("k", "old")

	// Wrong old value: should not swap.
	if m.CompareAndSwap("k", "wrong", "new") {
		t.Fatal("CompareAndSwap with wrong old value should return false")
	}

	// Correct old value: should swap.
	if !m.CompareAndSwap("k", "old", "new") {
		t.Fatal("CompareAndSwap with correct old value should return true")
	}

	v, _ := m.Load("k")
	if v != "new" {
		t.Fatalf("value after swap: got %v, want new", v)
	}
}

// --- ComparableMap tests ---

func TestComparableMap_BasicOperations(t *testing.T) {
	t.Parallel()

	var m gsyncmap.ComparableMap[int, int]

	m.Store(1, 100)

	v, ok := m.Load(1)
	if !ok || v != 100 {
		t.Fatalf("Load: got (%v, %v), want (100, true)", v, ok)
	}

	m.Delete(1)

	_, ok = m.Load(1)
	if ok {
		t.Fatal("key should be deleted")
	}
}

func TestComparableMap_CompareAndDelete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		stored      string
		compareWith string
		wantDeleted bool
	}{
		{"match", "v1", "v1", true},
		{"mismatch", "v1", "v2", false},
		{"key absent", "", "", false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var m gsyncmap.ComparableMap[string, string]

			if tc.name != "key absent" {
				m.Store("k", tc.stored)
			}

			deleted := m.CompareAndDelete("k", tc.compareWith)
			if deleted != tc.wantDeleted {
				t.Fatalf("CompareAndDelete: got %v, want %v", deleted, tc.wantDeleted)
			}
		})
	}
}

func TestComparableMap_CompareAndSwap(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		stored    string
		oldVal    string
		newVal    string
		wantSwap  bool
		wantValue string
	}{
		{"match", "v1", "v1", "v2", true, "v2"},
		{"mismatch", "v1", "wrong", "v2", false, "v1"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var m gsyncmap.ComparableMap[string, string]
			m.Store("k", tc.stored)

			swapped := m.CompareAndSwap("k", tc.oldVal, tc.newVal)
			if swapped != tc.wantSwap {
				t.Fatalf("CompareAndSwap: got %v, want %v", swapped, tc.wantSwap)
			}

			v, _ := m.Load("k")
			if v != tc.wantValue {
				t.Fatalf("value after swap: got %v, want %v", v, tc.wantValue)
			}
		})
	}
}

func TestComparableMap_LoadOrStore(t *testing.T) {
	t.Parallel()

	var m gsyncmap.ComparableMap[int, int]

	actual, loaded := m.LoadOrStore(1, 10)
	if loaded || actual != 10 {
		t.Fatalf("LoadOrStore new: got (%v, %v), want (10, false)", actual, loaded)
	}

	actual, loaded = m.LoadOrStore(1, 99)
	if !loaded || actual != 10 {
		t.Fatalf("LoadOrStore existing: got (%v, %v), want (10, true)", actual, loaded)
	}
}

func TestComparableMap_Clear(t *testing.T) {
	t.Parallel()

	var m gsyncmap.ComparableMap[int, int]

	for i := range 5 {
		m.Store(i, i)
	}

	m.Clear()

	var count int

	m.Range(func(_, _ int) bool {
		count++
		return true
	})

	if count != 0 {
		t.Fatalf("expected empty map after Clear, got %d entries", count)
	}
}

// --- Concurrent tests ---

func TestMap_Concurrent_StoreLoad(t *testing.T) {
	t.Parallel()

	const (
		goroutines = 50
		iterations = 100
	)

	var (
		m  gsyncmap.Map[int, int]
		wg sync.WaitGroup
	)

	wg.Add(goroutines)

	for g := range goroutines {
		go func(id int) {
			defer wg.Done()

			for i := range iterations {
				m.Store(id*iterations+i, i)
				m.Load(id*iterations + i)
			}
		}(g)
	}

	wg.Wait()
}

func TestMap_Concurrent_LoadOrStore(t *testing.T) {
	t.Parallel()

	const goroutines = 50

	var (
		m      gsyncmap.Map[string, int]
		stored atomic.Int64
	)

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for range goroutines {
		go func() {
			defer wg.Done()

			_, loaded := m.LoadOrStore("shared", 1)
			if !loaded {
				stored.Add(1)
			}
		}()
	}

	wg.Wait()

	// Only one goroutine should have stored the value.
	if stored.Load() != 1 {
		t.Fatalf("LoadOrStore: expected exactly 1 store, got %d", stored.Load())
	}
}

func TestMap_Concurrent_Range(t *testing.T) {
	t.Parallel()

	const (
		entries    = 100
		goroutines = 10
	)

	var m gsyncmap.Map[int, int]
	for i := range entries {
		m.Store(i, i)
	}

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for range goroutines {
		go func() {
			defer wg.Done()

			m.Range(func(_, _ int) bool {
				return true
			})
		}()
	}

	wg.Wait()
}

func TestComparableMap_Concurrent_CompareAndSwap(t *testing.T) {
	t.Parallel()

	const goroutines = 50

	var m gsyncmap.ComparableMap[string, int]
	m.Store("counter", 0)

	// Each goroutine races to be the first to swap 0 -> 1.
	var (
		swapped atomic.Int64
		wg      sync.WaitGroup
	)

	wg.Add(goroutines)

	for range goroutines {
		go func() {
			defer wg.Done()

			if m.CompareAndSwap("counter", 0, 1) {
				swapped.Add(1)
			}
		}()
	}

	wg.Wait()

	if swapped.Load() != 1 {
		t.Fatalf("expected exactly one successful CompareAndSwap, got %d", swapped.Load())
	}
}
