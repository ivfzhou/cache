package defaultimpl

import "reflect"

// defaultValue 默认实现，缓存value对象。
type defaultValue struct {
	// val JSON字符串。
	val string
	// expire unix毫秒时间戳。
	expire int64
	typ    reflect.Type
}

func (v *defaultValue) newVal() interface{} {
	return reflect.New(v.typ).Interface()
}
