package orbitdb

import (
	"fmt"
	"github.com/keks/go-ipfs-colog"
	"io"
	"time"
)

func ExampleFeedStore() {
	s1, err := NewOrbitDB(topic)
	assert(err == nil, err)

	time.Sleep(10 * time.Millisecond)

	s2, err := NewOrbitDB(topic)
	assert(err == nil, err)

	fs1 := NewFeedStore(s1)
	fs2 := NewFeedStore(s2)

	fmt.Println("fs1: add foo")
	fooEntry, err := fs1.Add("foo")
	assert(err == nil, err)

	time.Sleep(10 * time.Millisecond)

	fmt.Println("fs2: add bar")
	_, err = fs2.Add("bar")
	assert(err == nil, err)

	time.Sleep(10 * time.Millisecond)

	fmt.Println("fs2: del foo")
	_, err = fs2.Delete(fooEntry.Hash)
	assert(err == nil, err)

	time.Sleep(10 * time.Millisecond)

	res := fs1.Query(colog.Query{})

	var (
		p EventPayload
	)
	for err == nil {
		p, err = res()
		assert(err == nil || err == io.EOF, err)

		fmt.Println(p.DataString())
	}

	// Output:
	// fs1: add foo
	// fs2: add bar
	// fs2: del foo
	// bar
}
