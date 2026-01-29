package main

import (
	"bytes"
	"fmt"
	"sync"
	"time"
)

const iterations = 500000

// Pool pour réutiliser des buffers
var bufferPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

// ==========================================================
// VERSION SANS POOL
// ==========================================================
func withoutPool() {
	start := time.Now()

	for i := 0; i < iterations; i++ {
		// À chaque itération, on créé un nouveau buffer
		buf := new(bytes.Buffer)
		buf.WriteString("hello world")
		_ = buf.String()
	}

	fmt.Println("Sans pool :", time.Since(start))
}

// ==========================================================
// VERSION AVEC POOL
// ==========================================================
func withPool() {
	start := time.Now()

	for i := 0; i < iterations; i++ {
		// On récupère un buffer dans le pool
		buf := bufferPool.Get().(*bytes.Buffer)
		buf.Reset() // important

		buf.WriteString("hello world")
		_ = buf.String()

		// On remet le buffer dans le pool
		bufferPool.Put(buf)
	}

	fmt.Println("Avec pool :", time.Since(start))
}

func main() {
	fmt.Println("Comparaison avec", iterations, "itérations:")
	withoutPool()
	withPool()
}
