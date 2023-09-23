/*
 * Copyright (c) 2023 ivfzhou
 * cache is licensed under Mulan PSL v2.
 * You can use this software according to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *          http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
 * MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 */

package cache_test

import (
	"reflect"
	"runtime"
	"strconv"
	"sync"
	"testing"
	"time"

	"gitee.com/ivfzhou/cache"
)

func TestCache(t *testing.T) {
	c := cache.New()
	c.SetMaxMemory("1KB")

	c.Set("int", 1, time.Second*3)
	c.Set("string", "abc", time.Second*3)
	c.Set("slice", []int{1, 2, 3}, time.Second*3)
	c.Set("struct", struct {
		Field string
	}{"cba"}, time.Second*3)
	c.Set("map", map[int]string{1: "hello"}, time.Second*3)

	v, ok := c.Get("int")
	if !ok {
		t.Fatal("unexpect get value")
	}
	if i, ok := v.(int); !ok {
		t.Fatal("unexpect get value")
	} else if i != 1 {
		t.Fatal("unexpect get value", i)
	}

	v, ok = c.GetDel("string")
	if !ok {
		t.Fatal("unexpect get value")
	}
	if i, ok := v.(string); !ok {
		t.Fatal("unexpect get value")
	} else if i != "abc" {
		t.Fatal("unexpect get value", i)
	}

	v, ok = c.GetDel("slice")
	if !ok {
		t.Fatal("unexpect get value")
	}
	if i, ok := v.([]int); !ok {
		t.Fatal("unexpect get value")
	} else if !reflect.DeepEqual(i, []int{1, 2, 3}) {
		t.Fatal("unexpect get value", i)
	}

	v, ok = c.GetDel("struct")
	if !ok {
		t.Fatal("unexpect get value")
	}
	if i, ok := v.(struct {
		Field string
	}); !ok {
		t.Fatal("unexpect get value")
	} else if !reflect.DeepEqual(i, struct {
		Field string
	}{"cba"}) {
		t.Fatal("unexpect get value", i)
	}

	v, ok = c.GetDel("map")
	if !ok {
		t.Fatal("unexpect get value")
	}
	if i, ok := v.(map[int]string); !ok {
		t.Fatal("unexpect get value")
	} else if !reflect.DeepEqual(i, map[int]string{1: "hello"}) {
		t.Fatal("unexpect get value", i)
	}

	time.Sleep(time.Second * 3)
	if _, ok = c.GetDel("int"); ok {
		t.Fatal("unexpect get value")
	}

	wg := sync.WaitGroup{}
	for i := 0; i < 162; i++ {
		wg.Add(1)
		go func(i int) {
			c.Set("m"+strconv.FormatInt(int64(i), 10), "a", -1)
			wg.Done()
		}(i)
		c.Get("m" + strconv.FormatInt(int64(i), 10))
	}
	wg.Wait()
	keys := c.Keys()
	if keys != 162 {
		t.Fatal("unexpect keys", keys)
	}
	if c.Set("m", "a", -1) {
		t.Fatal("unexpect set value")
	}
	c.Flush()

	c = nil
	runtime.GC()
}
