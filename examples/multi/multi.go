package main

import (
	"fmt"
	"github.com/HarryWang29/progressbar"
	"os"
	"sync"
	"time"
)

func main() {
	waitTime := time.Millisecond * 100
	pb := progressbar.NewMultiProgress()
	pb.Start()
	var wg sync.WaitGroup
	//bar1, _ := pb.Add64Bar("1", 20)
	_, _ = pb.Add64Bar("1", 80)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 1; i <= 80; i++ {
			//bar1.Add64(1)
			pb.BarAdd("1", 1)
			time.Sleep(waitTime)
		}
	}()

	//bar2, _ := pb.AddDefaultBar("2", 40)
	_, _ = pb.AddDefaultBar("2", 40)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 40; i++ {
			//bar2.Add64(1)
			pb.BarAdd("2", 1)
			time.Sleep(waitTime)
		}
	}()

	time.Sleep(time.Second)
	bar3 := progressbar.NewOptions(100,
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
	pb.AddBar("3", bar3)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 1; i <= 100; i++ {
			bar3.Add64(1)
			time.Sleep(waitTime)
		}
	}()
	wg.Wait()
}
