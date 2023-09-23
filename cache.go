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

import "time"

type Cache interface {

	// Set 设置缓存。expire为-1代表永不过期。
	Set(key string, val any, expire time.Duration) bool

	// Get 获取缓存。
	Get(key string) (any, bool)

	// GetDel 获取并删除缓存
	GetDel(key string) (any, bool)

	// Del 删除这个key关联的缓存。
	Del(key string) bool

	// Exists 判断键是否存在与缓存中。
	Exists(key string) bool

	// Size 返回占用内存字节数，
	// 容量仅考虑key和value序列化后的字节大小，不考虑库本身使用的内存占用，也不考虑储存key-value的map的内存对齐占用的内存。
	Size() uint

	// Flush 清除所有键值对。如已有go程在清理则返回false。
	Flush() bool

	// Keys 返回键值对个数。
	Keys() int

	// SetMemoryEvictPolicy 缓存容量超限时的执行策略
	SetMemoryEvictPolicy(MemoryEvictPolicy)

	// SetSerializer 变量序列化策略
	SetSerializer(Serializer)

	// SetLogger 设置log输出对象
	SetLogger(Logger)

	// SetMaxMemory 设置缓存最大大小。默认不限制。
	// 单位：b、kb、mb、gb、tb
	SetMaxMemory(size string) bool
}

// MemoryEvictPolicy 缓存容量超限时的执行策略
type MemoryEvictPolicy interface {
	// Handle 当内存已满时调用，返回true表示清理成功
	Handle(Cache) bool
}

// Logger 日志打印接口
type Logger interface {
	Info(format string, args ...any)
	Warn(format string, args ...any)
	Error(format string, args ...any)
	SetLevel(Level)
}

// Level 日志打印等级
type Level int

const (
	LevelInfo Level = iota
	LevelWarn
	LevelError
	LevelSilence
)

// Serializer 变量序列化策略
type Serializer interface {
	Serialize(any) ([]byte, error)
	Deserialize([]byte, any) error
}

// Get 获取缓存，没有缓存或者类型不匹配返回零值。
func Get[T any](c Cache, key string) (t T) {
	val, ok := c.Get(key)
	if !ok {
		return
	}
	t, ok = val.(T)
	return
}
