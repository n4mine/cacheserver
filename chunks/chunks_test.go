package chunks

import (
	"fmt"
	"testing"
	"time"

	"github.com/n4mine/cacheserver/config"
)

func Test(t *testing.T) {
	c := config.LoadConfig("../etc/dev.cfg")

	cs := NewChunks(c.Cache)
	for i := 1; i <= 19; i++ {
		err := cs.Push(int64(i), float64(i))
		if err != nil {
			t.Fatal(err)
		}
		time.Sleep(1e3)
	}
	for _, i := range cs.Get(9, 21) {
		for i.Next() {
			fmt.Println(i.Values())
		}
	}
}

func Benchmark_Push(b *testing.B) {
	c := config.LoadConfig("../etc/dev.cfg")

	cs := NewChunks(c.Cache)
	for i := 0; i < b.N; i++ {
		cs.Push(int64(i+1), float64(i))
	}
}
