package timewheel

import (
	"context"
	"math/rand"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestNewTimeWheel(t *testing.T) {
	type args struct {
		unit time.Duration
		cap  int
	}
	tests := []struct {
		name string
		args args
		want *TimeWheel
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTimeWheel(tt.args.unit, tt.args.cap); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTimeWheel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimeWheel_step(t *testing.T) {
	type fields struct {
		ticker *time.Ticker
		array  []tasks
		unit   time.Duration
		idx    int
	}
	type args struct {
		n int
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantCircle int
		wantRem    int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tw := &TimeWheel{
				array: tt.fields.array,
				unit:  tt.fields.unit,
				idx:   tt.fields.idx,
			}
			gotCircle, gotRem := tw.step(tt.args.n)
			if gotCircle != tt.wantCircle {
				t.Errorf("TimeWheel.step() gotCircle = %v, want %v", gotCircle, tt.wantCircle)
			}
			if gotRem != tt.wantRem {
				t.Errorf("TimeWheel.step() gotRem = %v, want %v", gotRem, tt.wantRem)
			}
		})
	}
}

func TestTimeWheel_pop(t *testing.T) {
	type fields struct {
		ticker *time.Ticker
		array  []tasks
		unit   time.Duration
		idx    int
	}
	tests := []struct {
		name   string
		fields fields
		want   tasks
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tw := &TimeWheel{
				array: tt.fields.array,
				unit:  tt.fields.unit,
				idx:   tt.fields.idx,
			}
			if got := tw.pop(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TimeWheel.pop() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimeWheel_AddTask(t *testing.T) {
	type fields struct {
		ticker *time.Ticker
		array  []tasks
		unit   time.Duration
		idx    int
	}
	type args struct {
		nd time.Duration
		do func()
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tw := &TimeWheel{
				array: tt.fields.array,
				unit:  tt.fields.unit,
				idx:   tt.fields.idx,
			}
			tw.AddTask(tt.args.nd, tt.args.do)
		})
	}
}

func TestTimeWheel_AddTasks(t *testing.T) {
	type fields struct {
		ticker *time.Ticker
		array  []tasks
		unit   time.Duration
		idx    int
	}
	type args struct {
		nd     time.Duration
		doList []func()
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tw := &TimeWheel{
				array: tt.fields.array,
				unit:  tt.fields.unit,
				idx:   tt.fields.idx,
			}
			tw.AddTasks(tt.args.nd, tt.args.doList...)
		})
	}
}

func TestTimeWheel_Run(t *testing.T) {
	type fields struct {
		ticker *time.Ticker
		array  []tasks
		unit   time.Duration
		idx    int
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tw := &TimeWheel{
				array: tt.fields.array,
				unit:  tt.fields.unit,
				idx:   tt.fields.idx,
			}
			tw.Run(tt.args.ctx)
		})
	}
}

const timeFormat = "2006-01-02 15:04:05"

/*
func TestTimeWheel(t *testing.T) {
	t.Run("test1", func(t *testing.T) {
		tw := NewTimeWheel(time.Second, 7)
		r := rand.New(rand.NewSource(int64(time.Now().UnixNano())))
		x := r.Intn(10)
		t.Logf("%v - forloop rand 5 seconds %d times", time.Now().Format(timeFormat), x)
		for i := 0; i < x; i++ {
			n := i
			tw.AddTask(5*time.Second, func() {
				t.Logf("%v - [%d] after 5 second, do print", time.Now().Format(timeFormat), n)
			})
		}
		tw.AddTask(20*time.Second, func() {
			t.Logf("%v - after 20 second, do print", time.Now().Format(timeFormat))
		})
		tw.AddTask(50*time.Second, func() {
			t.Logf("%v - after 50 second, do print", time.Now().Format(timeFormat))
		})
		tw.AddTask(65*time.Second, func() {
			t.Logf("%v - after 65 second, do print", time.Now().Format(timeFormat))
		})
		ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
		defer cancel()
		tw.Run(ctx)
	})
}
*/

func TestRandomTimeWheel(t *testing.T) {
	t.Run("test1", func(t *testing.T) {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		rcap := r.Intn(90) + 10
		tw := NewTimeWheel(time.Millisecond, rcap)
		wg := &sync.WaitGroup{}
		wg.Add(1)
		go func(tw *TimeWheel, wg *sync.WaitGroup) {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
			defer cancel()
			tw.Run(ctx)
		}(tw, wg)
		n := 300
		for i := 0; i < 300; i++ {
			mt := time.Duration(r.Intn(100000)) * time.Millisecond
			start := time.Now()
			tw.AddTask(mt, func() {
				if x := time.Since(start).Abs(); x > mt+1*time.Millisecond || x < mt-1*time.Millisecond {
					t.Errorf("now-start != mt: want %v got %v", mt, x)
				}
				n--
			})
			time.Sleep(time.Duration(r.Intn(1000)+20) * time.Millisecond)
		}
		wg.Wait()
		if n != 0 {
			t.Errorf("n is not 0, got %d !", n)
		}
	})
}
