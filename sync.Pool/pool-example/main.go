package main

import (
	"bytes"
	"fmt"
	"sync"
)

// A sync.Pool to manage reusable byte buffers.
var pool = sync.Pool{
	New: func() interface{} {
		return &bytes.Buffer{}
	},
}

// worker simulates a worker that uses a byte buffer from the pool.
func worker(id int, wg *sync.WaitGroup) {
	defer wg.Done()

	buf := pool.Get().(*bytes.Buffer)
	buf.WriteString(fmt.Sprintf("worker %d\n", id))

	fmt.Print(buf.String())

	buf.Reset()
	pool.Put(buf)
}

func main() {
	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go worker(i, &wg)
	}

	wg.Wait()
}
