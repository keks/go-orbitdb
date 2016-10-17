package orbitdb

import (
	"fmt"
	"time"
)

func ExampleKVStore() {
	s1, err := NewStore(topic)
	assert(err == nil, err)

	time.Sleep(5 * time.Millisecond)

	s2, err := NewStore(topic)
	assert(err == nil, err)

	kv1 := NewKVStore(s1)
	kv2 := NewKVStore(s2)

	fmt.Println("kv1: foo=bar")
	err = kv1.Put("foo", "bar")
	assert(err == nil, err)

	time.Sleep(5 * time.Millisecond)

	val2, err := kv2.Get("foo")
	assert(err == nil, err)
	assert(val2 == "bar", "val2 != \"bar\"", val2)
	fmt.Println("kv2: get foo -> bar")

	val1, err := kv2.Get("foo")
	assert(err == nil, err)
	assert(val1 == "bar", "val1 != \"bar\"", val1)
	fmt.Println("kv1: get foo -> bar")

	//Output:
	// kv1: foo=bar
	// kv2: get foo -> bar
	// kv1: get foo -> bar
}
