package main

import (
	"context"
	"fs_monitor/timewheel"
	"log"
	"math/rand"
	"sync"
	"time"
)

func main() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	rcap := r.Intn(80) + 2
	log.Printf("random cap: %d", rcap)
	tw := timewheel.NewTimeWheel(time.Second, rcap)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func(tw *timewheel.TimeWheel, wg *sync.WaitGroup) {
		defer wg.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()
		tw.Run(ctx)
	}(tw, wg)
	n := 300
	for i := 0; i < 300; i++ {
		mt := time.Duration(r.Intn(99)+1) * time.Second
		start := time.Now()
		tw.AddTask(mt, func() {
			if x := time.Since(start).Abs(); x > mt+100*time.Millisecond || x < mt-100*time.Millisecond {
				log.Printf("now-start != mt: want %v got %v", mt, x)
			} else {
				log.Println("callback ok")
			}
			n--
		})
		time.Sleep(time.Duration(r.Intn(10000)+1000) * time.Millisecond)
	}
	wg.Wait()
	if n != 0 {
		log.Printf("n is not 0, got %d !", n)
	}
}
