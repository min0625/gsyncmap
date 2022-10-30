package gsyncmap_test

import (
	"fmt"

	"github.com/min0625/gsyncmap"
)

func ExampleMap() {
	var m gsyncmap.Map[string, string]

	m.Store("k1", "v1")
	fmt.Println(m.Load("k1"))
	fmt.Println(m.Load("k2"))

	m.Delete("k1")
	fmt.Println(m.Load("k1"))

	// Output:
	// v1 true
	//  false
	//  false
}

func ExampleMap_LoadOrStore() {
	var m gsyncmap.Map[string, string]

	fmt.Println(m.LoadOrStore("k1", "v1"))
	fmt.Println(m.Load("k1"))
	fmt.Println(m.LoadOrStore("k1", "v2"))
	fmt.Println(m.Load("k1"))

	// Output:
	// v1 false
	// v1 true
	// v1 true
	// v1 true
}

func ExampleMap_LoadAndDelete() {
	var m gsyncmap.Map[string, string]

	m.Store("k1", "v1")

	fmt.Println(m.LoadAndDelete("k1"))
	fmt.Println(m.Load("k1"))
	fmt.Println(m.LoadAndDelete("k1"))
	fmt.Println(m.Load("k1"))

	// Output:
	// v1 true
	//  false
	//  false
	//  false
}

func ExampleMap_Range() {
	var m gsyncmap.Map[string, string]

	m.Store("k1", "v1")
	m.Store("k2", "v2")
	m.Store("k3", "v3")

	m.Range(func(key, value string) bool {
		fmt.Println(key, value)
		return true
	})

	// Unordered output:
	// k1 v1
	// k2 v2
	// k3 v3
}

func ExampleMap_Range_break() {
	var m gsyncmap.Map[string, string]

	m.Store("k1", "v1")
	m.Store("k2", "v2")
	m.Store("k3", "v3")

	var cnt int

	m.Range(func(key, value string) bool {
		cnt++
		return false
	})

	fmt.Println(cnt)

	// Output:
	// 1
}
