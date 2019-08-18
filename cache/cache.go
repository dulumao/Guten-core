package cache

import (
	"github.com/dulumao/Guten-utils/conv"
	"github.com/dulumao/Guten-utils/os/cache"
	"github.com/dulumao/Guten-core/env"
)

var Cache cache.Cache

func New(driver string) {
	var cacheConfig = ``

	// memory
	if driver == "memory" {
		cacheConfig = `{"interval":` + conv.String(env.Value.Cache.Memory.Interval) + `}`
	}

	// file
	if driver == "file" {
		cacheConfig = `{"CachePath":"` + env.Value.Cache.File.Path + `","FileSuffix":"` + env.Value.Cache.File.FileSuffix + `","DirectoryLevel":` + conv.String(env.Value.Cache.File.DirectoryLevel) + `,"EmbedExpiry":` + conv.String(env.Value.Cache.File.EmbedExpiry) + `}`
	}

	// redis
	if driver == "redis" {
		cacheConfig = `{"key":` + env.Value.Cache.Redis.Key + `,"conn":` + env.Value.Cache.Redis.Addr + `,"dbNum":"` + conv.String(env.Value.Cache.Redis.DbNumber) + `","password":` + env.Value.Cache.Redis.Password + `}`
	}

	// memcache
	if driver == "memcache" {
		cacheConfig = `{"conn":` + env.Value.Cache.Memcache.Addr + `}`
	}

	if adapter, err := cache.NewCache(driver, cacheConfig); err == nil {
		Cache = adapter
	}
}
