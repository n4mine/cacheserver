package models

import (
	"errors"
)

var (
	GitVer    = "unknown"
	BuildTime = "unknown"
)

var (
	ErrNonEnoughData  = errors.New("not enough data")
	ErrNonExistSeries = errors.New("not exist series")
	ErrInternalError  = errors.New("internal error")
)

const (
	CodeSucc = iota
	CodeUserErr
	CodeNonEnoughErr
	CodeNonExistSeries
	CodeInternalErr
)

type Point struct {
	Key       string  `msg:"key"`
	Timestamp int64   `msg:"timestamp"`
	Value     float64 `msg:"value"`
}
