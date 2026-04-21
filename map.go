// Package gsyncmap provides a generic, type-safe wrapper around [sync.Map].
//
// There are two map types available:
//
//   - [Map]: accepts any value type. Suitable for most use cases.
//     CompareAndDelete and CompareAndSwap are available but deprecated due to
//     potential runtime panic when the value type is not comparable.
//
//   - [ComparableMap]: requires the value type to implement comparable.
//     Provides safe CompareAndDelete and CompareAndSwap without risk of panic.
package gsyncmap

import "sync"

func zero[T any]() T {
	return *new(T)
}

// Map is a generic, type-safe concurrent map wrapping [sync.Map].
//
// The zero value of Map is valid and ready to use. Map must not be copied
// after first use (same restriction as [sync.Map]).
//
// If the value type is known to be comparable and CompareAndDelete or
// CompareAndSwap is required, prefer using [ComparableMap] instead.
type Map[Key comparable, Value any] sync.Map

func (m *Map[Key, Value]) syncMap() *sync.Map {
	return (*sync.Map)(m)
}

// Store sets the value for a key.
func (m *Map[Key, Value]) Store(key Key, value Value) {
	m.syncMap().Store(key, value)
}

// Clear deletes all the entries, resulting in an empty Map.
func (m *Map[Key, Value]) Clear() {
	m.syncMap().Clear()
}

// Load returns the value stored in the map for a key, or the zero value if no
// value is present. The ok result indicates whether value was found in the map.
func (m *Map[Key, Value]) Load(key Key) (value Value, ok bool) {
	anyValue, ok := m.syncMap().Load(key)
	if !ok {
		return zero[Value](), false
	}

	return anyValue.(Value), ok
}

// Delete deletes the value for a key.
func (m *Map[Key, Value]) Delete(key Key) {
	m.syncMap().Delete(key)
}

// Range calls f sequentially for each key and value present in the map.
// If f returns false, range stops the iteration.
//
// Range does not necessarily correspond to any consistent snapshot of the Map's
// contents: no key will be visited more than once, but if the value for any key
// is stored or deleted concurrently, Range may reflect any mapping for that key
// from any point during the Range call.
func (m *Map[Key, Value]) Range(f func(key Key, value Value) bool) {
	m.syncMap().Range(func(key, value any) bool {
		return f(key.(Key), value.(Value))
	})
}

// LoadOrStore returns the existing value for the key if present.
// Otherwise, it stores and returns the given value.
// The loaded result is true if the value was loaded, false if stored.
func (m *Map[Key, Value]) LoadOrStore(key Key, value Value) (actual Value, loaded bool) {
	anyActual, loaded := m.syncMap().LoadOrStore(key, value)
	return anyActual.(Value), loaded
}

// LoadAndDelete deletes the value for a key, returning the previous value if any.
// The loaded result reports whether the key was present.
func (m *Map[Key, Value]) LoadAndDelete(key Key) (value Value, loaded bool) {
	anyValue, loaded := m.syncMap().LoadAndDelete(key)
	if !loaded {
		return zero[Value](), false
	}

	return anyValue.(Value), loaded
}

// CompareAndDelete deletes the entry for key if its value is equal to oldVal.
//
// Deprecated: CompareAndDelete panics at runtime if Value is not a comparable
// type (e.g. slice, map, or func). Use [ComparableMap.CompareAndDelete] instead
// to enforce comparability at compile time.
func (m *Map[Key, Value]) CompareAndDelete(key Key, oldVal Value) (deleted bool) {
	return m.syncMap().CompareAndDelete(key, oldVal)
}

// CompareAndSwap swaps the oldVal and newVal values for key if the value stored in
// the map is equal to oldVal.
//
// Deprecated: CompareAndSwap panics at runtime if Value is not a comparable
// type (e.g. slice, map, or func). Use [ComparableMap.CompareAndSwap] instead
// to enforce comparability at compile time.
func (m *Map[Key, Value]) CompareAndSwap(key Key, oldVal, newVal Value) bool {
	return m.syncMap().CompareAndSwap(key, oldVal, newVal)
}

// Swap stores the value for a key and returns the previous value if any.
// The loaded result reports whether the key was present.
func (m *Map[Key, Value]) Swap(key Key, value Value) (previous Value, loaded bool) {
	anyPrevious, loaded := m.syncMap().Swap(key, value)
	if !loaded {
		return zero[Value](), false
	}

	return anyPrevious.(Value), loaded
}
