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

func ExampleOrbitDB_twoParty() {
	s1, err := NewOrbitDB(topic)
	assert(err == nil, err)

	time.Sleep(5 * time.Millisecond)

	s2, err := NewOrbitDB(topic)
	assert(err == nil, err)

	fmt.Println("s1.colog.Watch()")
	eCh := s1.colog.Watch()

	fmt.Println(`s2.Add("foo")`)
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
	// s1.colog.Watch()
	// s2.Add("foo")
	// <-eCh
	// ok.
}

func ExampleOrbitDB_threeParty() {
	s1, err := NewOrbitDB(topic)
	assert(err == nil, err)

	time.Sleep(5 * time.Millisecond)

	s2, err := NewOrbitDB(topic)
	assert(err == nil, err)

	time.Sleep(5 * time.Millisecond)

	s3, err := NewOrbitDB(topic)
	assert(err == nil, err)

	fmt.Println("s1.colog.Watch()")
	eCh1 := s1.colog.Watch()

	fmt.Println("s2.colog.Watch()")
	eCh2 := s2.colog.Watch()

	fmt.Println(`s3.Add("foo")`)
	e1, err := s3.Add("foo")
	if err != nil {
		fmt.Println(err)
	}

	var e2, e3 *colog.Entry

	fmt.Println("<-eCh1")
	e2 = <-eCh1

	fmt.Println("<-eCh2")
	e3 = <-eCh2

	assert(e2.GetString() == e1.GetString(), "e1!=e2", e1, e2)
	assert(e2.GetString() == e3.GetString(), "e2!=e3", e2, e3)
	fmt.Println("ok.")

	// Output:
	// s1.colog.Watch()
	// s2.colog.Watch()
	// s3.Add("foo")
	// <-eCh1
	// <-eCh2
	// ok.
}
