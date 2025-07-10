package signer

import "time"

type Signer interface {
	Now() time.Time
	NowInUnix() int64
}

type signer struct{}
func NewSigner() Signer {
	return &signer{}
}

func (s *signer) Now() time.Time {
	return time.Now()
}

func (s *signer) NowInUnix() int64 {
	return time.Now().Unix()
}
