package gsyncmap

import "sync"

func zero[T any]() T {
	return *new(T)
}

type Map[Key comparable, Value any] sync.Map

func (m *Map[Key, Value]) syncMap() *sync.Map {
	return (*sync.Map)(m)
}

func (m *Map[Key, Value]) Store(key Key, value Value) {
	m.syncMap().Store(key, value)
}

func (m *Map[Key, Value]) Load(key Key) (value Value, ok bool) {
	anyValue, ok := m.syncMap().Load(key)
	if !ok {
		return zero[Value](), false
	}

	return anyValue.(Value), ok
}

func (m *Map[Key, Value]) Delete(key Key) {
	m.syncMap().Delete(key)
}

func (m *Map[Key, Value]) Range(f func(key Key, value Value) bool) {
	m.syncMap().Range(func(key, value any) bool {
		return f(key.(Key), value.(Value))
	})
}

func (m *Map[Key, Value]) LoadOrStore(key Key, value Value) (actual Value, loaded bool) {
	anyActual, loaded := m.syncMap().LoadOrStore(key, value)
	return anyActual.(Value), loaded
}

func (m *Map[Key, Value]) LoadAndDelete(key Key) (value Value, loaded bool) {
	anyValue, loaded := m.syncMap().LoadAndDelete(key)
	if !loaded {
		return zero[Value](), false
	}

	return anyValue.(Value), loaded
}
