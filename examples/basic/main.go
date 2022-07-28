package main

import (
	"time"

	"github.com/HarryWang29/progressbar/v4"
)

func main() {
	bar := progressbar.Default(100)
	for i := 0; i < 100; i++ {
		bar.Add(1)
		time.Sleep(40 * time.Millisecond)
	}
}
