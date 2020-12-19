package timer

import (
	"fmt"
	"github.com/tanyiqin/zack/utils"
	"sync"
	"testing"
)

func TestMain(m *testing.M) {
	m.Run()
}

func checkoutTime(args...interface{}) {
	if !(utils.NowTime() == args[0]) {
		t := args[1].(*testing.T)
		t.Error("not time")
	}
}

func TestTime(t *testing.T) {
	w := NewMWheel()
	w.Run()
	x := utils.NowTime()
	TestCase := []struct{
		u int64
	}{
		{5000000},
	}
	for _, tc := range TestCase {
		t.Run(fmt.Sprintf("begintime%v", tc.u), func(t *testing.T) {
			for i := 0; i < int(tc.u); i++ {
				d := genS(x, i)
				w.AddTimer(d, checkoutTime, d, &t)
			}
		})
	}
	w.Stop()
}

func TestTicker(t *testing.T) {
	w := NewMWheel()
	w.Run()
	mut := sync.WaitGroup{}
	mut.Add(10)
	f := func(...interface{}) {
		mut.Done()
		fmt.Println(utils.NowTime())
	}
	w.AddTicker(1, f)
	mut.Wait()
}

func genS(u int64, i int) int64{
	return u + int64(i%10)
}

func BenchmarkTime(b *testing.B) {
	w := NewMWheel()
	w.Run()
	cases := []struct{
		name string
		u int64
	}{
		{"dd-1", 50000},
	}
	for _, c := range cases {
		b.Run(c.name, func(b *testing.B){
			T := utils.NowTime()
			b.ResetTimer()
			for i:=0; i<int(c.u); i++ {
				w.AddTimer(genS(T, i), func(...interface{}){})
			}
			b.StopTimer()
			w.waitGroup.Wait()
		})
	}
}

