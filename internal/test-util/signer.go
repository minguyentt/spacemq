package testutil

import (
	"sync"
	"time"
)

type signerTest struct {
	t time.Time
	mu sync.Mutex
}

func NewSignerTest(t time.Time) *signerTest {
	return &signerTest{t:t}
}

func (s *signerTest) Now() time.Time {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.t
}

func (s *signerTest) NowInUnix() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.t.Unix()
}
