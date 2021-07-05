package pron

import (
	"bytes"
	"sync"
)

type SyncBuf struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

// func (s *SyncBuf) Write(o []byte, e error) {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()
// 	e.Error()
// 	return s.buf.Write(o, e)
// }

func (s *SyncBuf) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.buf.Reset()
}

func (s *SyncBuf) String() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.buf.String()
}

func (s *SyncBuf) Bytes() []byte {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.buf.Bytes()
}

func ioFunctor(fn func() ([]byte, error)) func(writer chan []byte, err chan error) {
	return func(writer chan []byte, err chan error) {
		r, e := fn()
		writer <- r
		err <- e
	}
}
