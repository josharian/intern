// Package intern interns strings.
// Interning is best effort only.
// Interned strings may be removed automatically
// at any time without notification.
// All functions may be called concurrently
// with themselves and each other.
package intern

import (
	"sync"
	"unsafe"
)

var (
	pool sync.Pool = sync.Pool{
		New: func() interface{} {
			return make(map[string]string)
		},
	}
)

// String returns s, interned.
func String(s string) string {
	m := pool.Get().(map[string]string)
	c, ok := m[s]
	if ok {
		pool.Put(m)
		return c
	}
	m[s] = s
	pool.Put(m)
	return s
}

// Bytes returns b converted to a string, interned.
func Bytes(b []byte) string {
	m := pool.Get().(map[string]string)
	c, ok := m[string(b)]
	if ok {
		pool.Put(m)
		return c
	}
	s := string(b)
	m[s] = s
	pool.Put(m)
	return s
}

// StringHeader copy from $(go env GOROOT)/src/reflect/value.go
type StringHeader struct {
	Data uintptr
	Len  int
}

// Equal compare two strings. If s1 and s2 are both interned, it's only 0.9 ns/op
func Equal(s1 string, s2 string) bool {
	x1 := (*StringHeader)(unsafe.Pointer(&s1))
	x2 := (*StringHeader)(unsafe.Pointer(&s2))
	return x1.Len == x2.Len && x1.Data == x2.Data || s1 == s2
}
