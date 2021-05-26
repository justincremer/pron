package pron

import (
	"bytes"
	"sync"
)

type IoBuf struct {
	outbuffer SyncBuf
	errbuffer SyncBuf
}

type SyncBuf struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (s *SyncBuf) Write(p []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.buf.Write(p)
}

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
