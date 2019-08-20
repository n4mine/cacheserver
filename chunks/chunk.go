package chunks

import (
	"fmt"

	"github.com/devtoolkits/go-tsz"
)

type Chunk struct {
	tsz.Series
	FirstTs   uint32
	LastTs    uint32
	NumPoints uint32
	Closed    bool
}

func NewChunk(t0 uint32) *Chunk {
	return &Chunk{
		Series:    *tsz.New(t0),
		FirstTs:   0,
		LastTs:    0,
		NumPoints: 0,
		Closed:    false,
	}
}

func (c *Chunk) Push(t uint32, v float64) error {
	if t <= c.LastTs {
		return fmt.Errorf("Point must be newer than already added points. t:%d lastTs: %d\n", t, c.LastTs)
	}
	c.Series.Push(t, v)
	c.NumPoints += 1
	c.LastTs = t

	return nil
}

func (c *Chunk) Finish() {
	c.Closed = true
	c.Series.Finish()
}
