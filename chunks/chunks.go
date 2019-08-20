// The code in this file was largely copy from https://github.com/grafana/metrictank/blob/master/mdata/aggmetric.go
// Copyright 2016-2018 Dieter Plaetinck, Anthony Woods, Jeremy Bingham, Damian Gryski, raintank inc
//
// This software is distributed under the terms of the GNU Affero General Public License.

package chunks

import (
	"fmt"
	"log"
	"sync"

	"github.com/n4mine/cacheserver/config"
)

type CS struct {
	Chunks          []*Chunk
	currentChunkPos int

	sync.RWMutex
}

func NewChunks(c config.CacheConfig) *CS {
	cs := make([]*Chunk, 0, c.NumOfChunks)

	return &CS{Chunks: cs}
}

func (cs *CS) Push(ts int64, value float64) error {
	t0 := uint32(ts - (ts % int64(config.C.Cache.SpanInSeconds)))

	// 尚无chunk
	if len(cs.Chunks) == 0 {
		c := NewChunk(uint32(t0))
		c.FirstTs = uint32(ts)
		cs.Chunks = append(cs.Chunks, c)

		return cs.Chunks[0].Push(uint32(ts), value)
	}

	// push到当前chunk
	currentChunk := cs.getChunk(cs.currentChunkPos)
	if t0 == currentChunk.T0 {
		if currentChunk.Closed {
			return fmt.Errorf("push to closed chunk")
		}

		return currentChunk.Push(uint32(ts), value)
	}

	if t0 < currentChunk.T0 {
		return fmt.Errorf("data @%v, goes back into previous chunk. currentchunk t0: %v\n", t0, currentChunk.T0)
	}

	// 需要新建chunk
	// 先finish掉现有chunk
	if !currentChunk.Closed {
		currentChunk.Finish()
	}

	// 超过chunks限制, pos回绕到0
	cs.currentChunkPos++
	if cs.currentChunkPos >= int(config.C.Cache.NumOfChunks) {
		cs.currentChunkPos = 0
	}

	// chunks未满, 直接append即可
	if len(cs.Chunks) < int(config.C.Cache.NumOfChunks) {
		c := NewChunk(uint32(t0))
		c.FirstTs = uint32(ts)
		cs.Chunks = append(cs.Chunks, c)

		return cs.Chunks[cs.currentChunkPos].Push(uint32(ts), value)
	} else {
		c := NewChunk(uint32(t0))
		c.FirstTs = uint32(ts)
		cs.Chunks[cs.currentChunkPos] = c

		return cs.Chunks[cs.currentChunkPos].Push(uint32(ts), value)
	}

	return nil
}

func (cs *CS) Get(from, to int64) []Iter {
	// 这种case不应该发生
	if from >= to {
		return nil
	}

	cs.RLock()
	defer cs.RUnlock()

	// cache server还没有数据
	if len(cs.Chunks) == 0 {
		return nil
	}

	var iters []Iter

	// from 超出最新chunk可能达到的最新点, 这种case不应该发生
	newestChunk := cs.getChunk(cs.currentChunkPos)
	if from >= int64(newestChunk.T0)+int64(config.C.Cache.SpanInSeconds) {
		return nil
	}

	// 假设共有2个chunk
	// len = 1, currentChunkPos = 0, oldestPos = 0
	// len = 2, currentChunkPos = 0, oldestPos = 1
	// len = 2, currentChunkPos = 1, oldestPos = 0
	oldestPos := cs.currentChunkPos + 1
	if oldestPos >= len(cs.Chunks) {
		oldestPos = 0
	}
	oldestChunk := cs.getChunk(oldestPos)
	if oldestChunk == nil {
		log.Println("unexpected nil chunk")
		return nil
	}

	// to 太老了, 这种case不应发生, 应由query处理
	if to <= int64(oldestChunk.FirstTs) {
		return nil
	}

	// 找from所在的chunk
	for from >= int64(oldestChunk.T0)+int64(config.C.Cache.SpanInSeconds) {
		oldestPos++
		if oldestPos >= len(cs.Chunks) {
			oldestPos = 0
		}
		oldestChunk = cs.getChunk(oldestPos)
		if oldestChunk == nil {
			log.Println("unexpected nil chunk")
			return nil
		}
	}

	// 找to所在的trunk
	newestPos := cs.currentChunkPos
	for to <= int64(newestChunk.T0) {
		newestPos--
		if newestPos < 0 {
			newestPos += len(cs.Chunks)
		}
		newestChunk = cs.getChunk(newestPos)
		if newestChunk == nil {
			log.Println("unexpected nil chunk")
			return nil
		}
	}

	for {
		c := cs.getChunk(oldestPos)
		iters = append(iters, NewIter(c.Iter()))
		if oldestPos == newestPos {
			break
		}
		oldestPos++
		if oldestPos >= len(cs.Chunks) {
			oldestPos = 0
		}
	}

	return iters
}

// GetInfo get oldest ts and newest ts in cache
func (cs *CS) GetInfo() (uint32, uint32) {
	cs.RLock()
	defer cs.RUnlock()

	return cs.GetInfoUnsafe()
}

func (cs *CS) GetInfoUnsafe() (uint32, uint32) {
	var oldestTs, newestTs uint32

	if len(cs.Chunks) == 0 {
		return 0, 0
	}

	newestChunk := cs.getChunk(cs.currentChunkPos)
	if newestChunk == nil {
		newestTs = 0
	} else {
		newestTs = newestChunk.LastTs
	}

	oldestPos := cs.currentChunkPos + 1
	if oldestPos >= len(cs.Chunks) {
		oldestPos = 0
	}

	oldestChunk := cs.getChunk(oldestPos)
	if oldestChunk == nil {
		oldestTs = 0
	} else {
		oldestTs = oldestChunk.FirstTs
	}

	return oldestTs, newestTs
}

func (cs CS) getChunk(pos int) *Chunk {
	if pos < 0 || pos >= len(cs.Chunks) {
		return cs.Chunks[0]
	}

	return cs.Chunks[pos]
}
