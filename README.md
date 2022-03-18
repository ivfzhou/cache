## 说明
本地缓存工具

## 使用
```golang
go get gitee.com/ivfzhou/cache@latest

cache := defaultimpl.NewCache()
cache.Set("key", value, time.Second*10)
val, ok := cache.Get("key")
fmt.Println(val)
```
