package chunks

import (
	"github.com/devtoolkits/go-tsz"
)

type Iter struct {
	*tsz.Iter
}

func NewIter(i *tsz.Iter) Iter {
	return Iter{
		i,
	}
}
