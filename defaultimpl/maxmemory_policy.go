package defaultimpl

import (
	"log"

	"gitee.com/ivfzhou/cache"
)

type defaultMaxMemoryPolicy struct{}

func (m *defaultMaxMemoryPolicy) handle(_ cache.Cache) (can bool) {
	log.Println("cache capacity is full, just discard")
	return false
}
