package defaultimpl_test

import (
	"runtime"
	"strconv"
	"sync"
	"testing"
	"time"

	"gitee.com/ivfzhou/cache/defaultimpl"
)

func TestCache(t *testing.T) {
	cache := defaultimpl.NewCache()
	cache.SetMaxMemory("1KB")

	cache.Set("int", 1, time.Second*3)
	cache.Set("string", "abc", time.Second*3)
	cache.Set("slice", []int{1, 2, 3}, time.Second*3)
	cache.Set("struct", struct {
		Field string
	}{"cba"}, time.Second*3)
	cache.Set("map", map[int]string{1: "hello"}, time.Second*3)

	t.Log(cache.Get("int"))
	t.Log(cache.Get("string"))
	t.Log(cache.Get("slice"))
	t.Log(cache.Get("struct"))
	t.Log(cache.Get("map"))

	time.Sleep(time.Second * 3)
	t.Log(cache.Get("int"))

	for i := 0; i < 162; i++ {
		cache.Set("m"+strconv.FormatInt(int64(i), 10), "a", -1)
	}
	t.Log(cache.Keys())
	cache.Set("m", "a", -1)
	cache.Flush()
	t.Log(cache.Keys())

	wg := sync.WaitGroup{}
	for i := 0; i < 160; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			cache.Set("a"+strconv.FormatInt(int64(i), 10), "b", -1)
		}(i)
	}
	wg.Wait()
	t.Log(cache.Keys())

	cache = nil
	runtime.GC()
}
