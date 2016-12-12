package orbitdb

import (
	"crypto/sha512"
	"encoding/binary"
	"testing"
	"time"
)

func BenchmarkHash(b *testing.B) {
	for n := 0; n < b.N; n++ {
		t := time.Now().Unix()
		buf := make([]byte, 8)
		binary.PutVarint(buf, t)
		_ = sha512.Sum512(buf)
	}
}

func BenchmarkAdd(b *testing.B) {
	db, err := NewOrbitDB("benchmark")
	if err != nil {
		b.Fatal(err)
	}

	for n := 0; n < b.N; n++ {
		t := time.Now().Unix()
		buf := make([]byte, 8)
		binary.PutVarint(buf, t)
		h := sha512.Sum512(buf)
		db.Add(h)
	}
}
