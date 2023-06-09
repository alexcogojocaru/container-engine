package main

import (
	"math/rand"
	"time"
)

func main() {
	for {
		vector := make([]int, 100)
		for i := 0; i < 100; i++ {
			vector[i] = rand.Int()
		}

		mean := 0
		for i := 0; i < 100; i++ {
			mean += vector[i]
		}

		time.Sleep(2 * time.Second)
	}
}
