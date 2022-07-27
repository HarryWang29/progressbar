package main

import (
	"fmt"
	"github.com/schollz/progressbar/v3"
	"os"
	"sync"
	"time"
)

func main() {
	waitTime := time.Millisecond * 100
	pb := progressbar.NewMultiProgress()
	pb.Start()
	var wg sync.WaitGroup
	bar1 := pb.Add64Bar(20)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 1; i <= 20; i++ {
			bar1.Add64(1)
			time.Sleep(waitTime)
		}
	}()

	bar2 := pb.AddDefaultBar(40)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 1; i <= 40; i++ {
			bar2.Add64(1)
			time.Sleep(waitTime)
		}
	}()

	time.Sleep(time.Second)
	bar3 := progressbar.NewOptions(1000,
		progressbar.OptionSetWriter(os.Stdout),
		progressbar.OptionSetWidth(15),
		progressbar.OptionSetDescription(fmt.Sprintf("downloading %s ...", "test")),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "=",
			SaucerHead:    ">",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))
	pb.AddBar(bar3)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 1; i <= 1000; i++ {
			bar3.Add64(1)
			time.Sleep(waitTime)
		}
	}()
	wg.Wait()
}
