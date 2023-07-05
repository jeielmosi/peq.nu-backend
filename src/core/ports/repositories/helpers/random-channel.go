package repository_helpers

import (
	"math/rand"
	"time"
)

func NewRandChanIdxs(size uint) <-chan int {
	rnd := rand.New(
		rand.NewSource(
			time.Now().UTC().UnixNano(),
		),
	)
	perm := rnd.Perm(int(size))

	ch := make(chan int, size)
	for _, val := range perm {
		ch <- val
	}
	close(ch)

	return ch
}
