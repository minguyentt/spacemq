package testutil

import (
	"sync"
	"time"
)

type signer struct {
	t time.Time
	mu sync.Mutex
}

func NewSignerTest(t time.Time) *signer {
	return &signer{t:t}
}

func (s *signer) Now() time.Time {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.t
}

func (s *signer) NowInUnix() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.t.Unix()
}
