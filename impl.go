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

package cache

import (
	"reflect"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// New 创建缓存对象实例。
// 内部采用json序列化对象储存，而非直接储存引用值。使用标准json库序列化，那么不能json序列化的字段将忽略。
// 如果设置了最大容量，那么超过容量后将不再进行存储。
func New() Cache {
	c := &defaultImpl{
		data:            make(map[string]*value),
		serializer:      &defaultSerializer{},
		maxMemoryPolicy: &defaultMaxMemoryPolicy{},
		logger:          &defaultLog{},
		stop:            make(chan struct{}),
	}
	go regularClean(c)
	runtime.SetFinalizer(c, stop)
	return c
}

type defaultImpl struct {
	serializer      Serializer
	maxMemoryPolicy MemoryEvictPolicy
	logger          Logger

	data       map[string]*value
	mu         sync.RWMutex
	maxSize    uint
	size       uint
	isFlushing uint32
	stop       chan struct{}
}

type value struct {
	val    []byte
	expire int64
	typ    reflect.Type
}

func (c *defaultImpl) Set(key string, val any, expire time.Duration) bool {
	data, err := c.serializer.Serialize(val)
	if err != nil {
		c.logger.Warn("serialize error key = %s %v", key, err)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.maxSize != 0 && c.size+uint(len(data)) > c.maxSize {
		if !c.maxMemoryPolicy.Handle(c) {
			c.logger.Warn("exceed max memory %d, stop cache", c.maxSize)
			return false
		}
	}

	t := int64(expire)
	if expire != -1 {
		t = time.Now().Add(expire).UnixMilli()
	}

	v, ok := c.data[key]
	if ok {
		oldSize := len(v.val)
		v.val = data
		v.expire = t
		c.size = uint(len(data) - oldSize)
	} else {
		c.size += uint(len(data) + len(key))
		c.data[key] = &value{val: data, expire: t, typ: reflect.TypeOf(val)}
	}
	return true
}

func (c *defaultImpl) Get(key string) (any, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, ok := c.data[key]
	if ok {
		if isTimeValid(value.expire) {
			newVal := reflect.New(value.typ).Interface()
			err := c.serializer.Deserialize(value.val, newVal)
			if err != nil {
				c.logger.Warn("serialize error key: %s, data: %s, type: %v %s", key, string(value.val), value.typ, err)
			}
			return reflect.ValueOf(newVal).Elem().Interface(), true
		}
	}
	return nil, false
}

func (c *defaultImpl) GetDel(key string) (any, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	value, ok := c.data[key]
	if ok {
		if isTimeValid(value.expire) {
			newVal := reflect.New(value.typ).Interface()
			err := c.serializer.Deserialize(value.val, &newVal)
			if err != nil {
				c.logger.Warn("serialize error key: %s, data: %s, type: %v %s", key, string(value.val), value.typ, err)
			}
			return reflect.ValueOf(newVal).Elem().Interface(), true
		}
		delete(c.data, key)
	}
	return nil, false
}

func (c *defaultImpl) Size() uint {
	return c.size
}

func (c *defaultImpl) Del(key string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	val, ok := c.data[key]
	if ok {
		delete(c.data, key)
		c.minusSizeLocked(key, val.val)
		return true
	}
	return false
}

func (c *defaultImpl) Exists(key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, ok := c.data[key]
	return ok
}

func (c *defaultImpl) Flush() bool {
	if !atomic.CompareAndSwapUint32(&c.isFlushing, 0, 1) {
		return false
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.data = make(map[string]*value)
	atomic.StoreUint32(&c.isFlushing, 0)
	c.size = 0

	return true
}

func (c *defaultImpl) Keys() int {
	return len(c.data)
}

func (c *defaultImpl) SetMaxMemory(size string) bool {
	memorySize, err := parseMemorySize(size, c.logger)
	if err != nil {
		return false
	}
	c.maxSize = memorySize
	return true
}

func (c *defaultImpl) SetLogger(logger Logger) {
	c.logger = logger
}

func (c *defaultImpl) SetMemoryEvictPolicy(m MemoryEvictPolicy) {
	c.maxMemoryPolicy = m
}

func (c *defaultImpl) SetSerializer(s Serializer) {
	c.serializer = s
}

func (c *defaultImpl) minusSizeLocked(key string, value []byte) {
	c.size -= uint(len(value) + len(key))
}

func stop(c *defaultImpl) {
	c.stop <- struct{}{}
}

func regularClean(c *defaultImpl) {
	ticker := time.NewTicker(time.Second)
	for range ticker.C {

		select {
		case <-c.stop:
			ticker.Stop()
		default:
		}

		c.mu.Lock()
		for k, v := range c.data {
			if !isTimeValid(v.expire) {
				delete(c.data, k)
				c.minusSizeLocked(k, v.val)
			}
		}
		c.mu.Unlock()
	}
}

func isTimeValid(t int64) bool {
	if t == -1 {
		return true
	}
	if t > time.Now().UnixMilli() {
		return true
	}
	return false
}
