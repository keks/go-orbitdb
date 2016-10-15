package orbitdb

import (
	"fmt"
	"time"

	"github.com/keks/go-ipfs-colog"
)

const topic = "testTopic"

func assert(check bool, args ...interface{}) {
	if !check {
		fmt.Print("assert: ")
		fmt.Println(args...)
	}
}

func ExampleStore() {
	s1, err := NewStore(topic)
	assert(err == nil, err)

	time.Sleep(5 * time.Millisecond)

	s2, err := NewStore(topic)
	assert(err == nil, err)

	fmt.Println("s1.Watch()")
	eCh := s1.Watch()

	fmt.Println("s2.Add(\"foo\")")
	e1, err := s2.Add("foo")
	if err != nil {
		fmt.Println(err)
	}

	var e2 *colog.Entry

	fmt.Println("<-eCh")
	e2 = <-eCh

	assert(e2.GetString() == e1.GetString(), "e1!=e2", e1, e2)
	fmt.Println("ok.")

	// Output:
	// s1.Watch()
	// s2.Add("foo")
	// <-eCh
	// ok.
}
