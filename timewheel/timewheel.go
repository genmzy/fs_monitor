package timewheel

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

type task struct {
	do     func()
	circle int
	wait   time.Duration
}

type tasks map[*task]bool

type TimeWheel struct {
	lock  *sync.Mutex
	array []tasks
	unit  time.Duration
	cap   int
	idx   int

	DoDebug bool
}

// every unit time check, cap specify the inner slice length
func NewTimeWheel(unit time.Duration, cap int) *TimeWheel {
	if cap < 2 {
		panic(fmt.Sprintf("invalid timewheel cap: %d", cap))
	}
	tw := &TimeWheel{
		idx:   0,
		cap:   cap,
		array: make([]tasks, cap),
		unit:  unit,
		lock:  &sync.Mutex{},

		DoDebug: true,
	}
	for i := range tw.array {
		tw.array[i] = make(map[*task]bool)
	}
	return tw
}

// FIXME: error circle count, not found reason yet
func (tw *TimeWheel) step(n int) (circle int, rem int) {
	rem = tw.idx + n
	circle = rem / tw.cap
	rem = rem % tw.cap
	return
}

func (tw *TimeWheel) pop() tasks {
	tw.lock.Lock()
	_, rem := tw.step(1)
	res := tw.array[tw.idx]
	tw.idx = rem
	tw.lock.Unlock()
	return res
}

func (tw *TimeWheel) AddTask(nd time.Duration, do func()) {
	if tw.DoDebug {
		log.Printf("add %v when idx: %d", nd, tw.idx)
	}
	if nd == 0 {
		do()
		return
	}
	n := int(nd / tw.unit)
	tw.lock.Lock()
	circle, rem := tw.step(n)
	t := &task{
		do:     do,
		circle: circle,
		wait:   nd,
	}
	tw.array[rem][t] = true
	tw.lock.Unlock()
}

func (tw *TimeWheel) AddTasks(nd time.Duration, doList ...func()) {
	for _, do := range doList {
		tw.AddTask(nd, do)
	}
}

func (tw *TimeWheel) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		tasks := tw.pop()
		for task := range tasks {
			if task.circle > 1 {
				if tw.DoDebug {
					log.Printf("[%v] task got circle when idx: %d : %d, waiting ...\n", task.wait, tw.idx, task.circle)
				}
				task.circle--
				continue
			}
			task.do()
			delete(tasks, task)
		}
		time.Sleep(tw.unit)
	}
}
