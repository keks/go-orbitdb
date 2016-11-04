package orbitdb

import (
	"fmt"
	"time"
)

func ExampleCounter() {
	s1, err := NewStore(topic)
	assert(err == nil, err)

	time.Sleep(5 * time.Millisecond)

	s2, err := NewStore(topic)
	assert(err == nil, err)

	ctrStore1 := NewCtrStore(s1)
	ctrStore2 := NewCtrStore(s2)

	ctrStore1.Increment(42)
	ctrStore2.Increment(23)

	time.Sleep(5 * time.Millisecond)

	fmt.Println(ctrStore1.Value())
	fmt.Println(ctrStore2.Value())

	// Output:
	// 65
	// 65
}
