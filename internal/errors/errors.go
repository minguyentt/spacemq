package errors

import (
	"fmt"
	"strings"
)

type Error struct {
	Code Code
	Area Area
	Err  error
}

type (
	Area string
	Code int8
)

const (
	Unknown Code = iota
	NotFound
	Exists
	Failed
)

func (c Code) String() string {
	switch c {
	case Unknown:
		return "UNKNOWN"
	case NotFound:
		return "NOT_FOUND"
	case Exists:
		return "ALREADY_EXISTS"
	case Failed:
		return "FAILED_OPERATION"
	}
	panic(fmt.Sprintf("unknown error code %d", c))
}

type RedisCmdError struct {
	Cmd string
	Err Error
}

func (e *RedisCmdError) Error() string {
	return fmt.Sprintf("redis command error: %s msg: %v", strings.ToUpper(e.Cmd), e.Err)
}
