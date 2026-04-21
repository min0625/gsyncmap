package gsyncmap

// ComparableMap is a generic, type-safe concurrent map wrapping [sync.Map].
//
// ComparableMap is a drop-in replacement for [Map] when the value type is
// comparable. It overrides CompareAndDelete and CompareAndSwap with safe
// implementations that cannot panic, because the Value type constraint
// guarantees comparability at compile time.
//
// The zero value of ComparableMap is valid and ready to use. ComparableMap
// must not be copied after first use (same restriction as [sync.Map]).
type ComparableMap[Key, Value comparable] struct {
	Map[Key, Value]
}

// CompareAndDelete deletes the entry for key if its value is equal to oldVal.
// Returns true if the entry was deleted.
func (m *ComparableMap[Key, Value]) CompareAndDelete(key Key, oldVal Value) (deleted bool) {
	return m.Map.syncMap().CompareAndDelete(key, oldVal)
}

// CompareAndSwap swaps the oldVal and newVal values for key if the value stored in
// the map is equal to oldVal. Returns true if the swap was performed.
func (m *ComparableMap[Key, Value]) CompareAndSwap(key Key, oldVal, newVal Value) (swapped bool) {
	return m.Map.syncMap().CompareAndSwap(key, oldVal, newVal)
}
