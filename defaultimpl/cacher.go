package defaultimpl

import (
	"log"
	"reflect"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"gitee.com/ivfzhou/cache"
	"gitee.com/ivfzhou/cache/defaultimpl/internal"
)

type defaultCacher struct {
	data            map[string]*defaultValue
	serializer      *defaultSerializer
	mu              sync.RWMutex
	maxmemoryPolicy *defaultMaxMemoryPolicy
	maxSize         uint64
	size            uint64
	isFlushing      uint32
	stop            chan struct{}
}

// NewCache 创建缓存对象实例。
// 内部采用json序列化对象储存，而非直接储存引用值。使用标准json库序列化，那么不能json序列化的字段将忽略。
// 如果设置了最大容量，那么超过容量后将不再进行存储。
// 容量仅考虑key和value序列化后的字节大小，不考虑库本身使用的内存占用，也不考虑储存key-value的map的内存对齐占用的内存。
// 使用原生map存储，而非sync.Map。
func NewCache() cache.Cache {
	c := &defaultCacher{
		data:            make(map[string]*defaultValue),
		serializer:      &defaultSerializer{},
		maxmemoryPolicy: &defaultMaxMemoryPolicy{},
		stop:            make(chan struct{}),
	}
	go regularClean(c)
	runtime.SetFinalizer(c, stopCacher)
	return c
}

// SetMaxMemory 设置缓存最大大小。
func (c *defaultCacher) SetMaxMemory(size string) bool {
	memorySize, err := internal.ParseMemorySize(size)
	if err != nil {
		log.Printf("paramter size parsing error, %s\n", err)
		return false
	}
	c.maxSize = memorySize
	return true
}

// Set 设置缓存。expire为-1代表永不过期。
func (c *defaultCacher) Set(key string, val interface{}, expire time.Duration) {
	data, err := c.serializer.serialize(val)
	if err != nil {
		log.Printf("val doesn't supported, [%s]\n", err)
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.maxSize != 0 && c.size+uint64(len(data)) > c.maxSize {
		log.Printf("cacher capacity is full, [%d]\n", c.maxSize)
		if !c.maxmemoryPolicy.handle(c) {
			return
		}
	}

	t := int64(expire)
	if expire != -1 {
		t = time.Now().Add(expire).UnixMilli()
	}

	value, ok := c.data[key]
	if ok {
		oldSize := len(value.val)
		value.val = data
		value.expire = t
		c.size = uint64(len(data) - oldSize)
	} else {
		c.size += uint64(len(data) + len(key))
		c.data[key] = &defaultValue{val: data, expire: t, typ: reflect.TypeOf(val)}
	}
}

// Get 获取缓存。
func (c *defaultCacher) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	value, ok := c.data[key]
	if ok {
		if isTimeValid(value.expire) {
			newVal := value.newVal()
			err := c.serializer.deserialize(value.val, newVal)
			if err != nil {
				log.Printf("unmarshal falied, key: [%s], data: [%s], type: [%v] %s\n", key, value.val, value.typ, err)
			}
			return newVal, true
		}
		delete(c.data, key)
	}
	return nil, false
}

// Del 删除这个key关联的缓存。
func (c *defaultCacher) Del(key string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	val, ok := c.data[key]
	if ok {
		log.Printf("deleting cache key [%s]\n", key)
		delete(c.data, key)
		c.minusSizeLocked(key, val.val)
		return true
	}
	return false
}

// Exists 判断键是否存在与缓存中。
func (c *defaultCacher) Exists(key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, ok := c.data[key]
	return ok
}

// Flush 清除所有键值对。如已有go程在清理则返回false。
func (c *defaultCacher) Flush() bool {
	if !atomic.CompareAndSwapUint32(&c.isFlushing, 0, 1) {
		log.Println("already in flushing")
		return false
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.data = make(map[string]*defaultValue)
	log.Println("cacher is flushed")
	atomic.StoreUint32(&c.isFlushing, 0)
	c.size = 0

	return true
}

// Keys 返回键值对个数。
func (c *defaultCacher) Keys() int64 {
	return int64(len(c.data))
}

func (c *defaultCacher) minusSizeLocked(key, value string) {
	c.size -= uint64(len(value) + len(key))
}

func stopCacher(c *defaultCacher) {
	c.stop <- struct{}{}
}

func regularClean(c *defaultCacher) {
	ticker := time.NewTicker(time.Second)
	for range ticker.C {

		select {
		case <-c.stop:
			log.Println("cacher stopping")
			ticker.Stop()
		default:
		}

		c.mu.Lock()
		log.Printf("regular clean cache begin\n")
		for k, v := range c.data {
			if !isTimeValid(v.expire) {
				delete(c.data, k)
				c.minusSizeLocked(k, v.val)
				log.Printf("regular clean cache, [%s]\n", k)
			}
		}
		c.mu.Unlock()
		log.Printf("regular clean cache end\n")
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
