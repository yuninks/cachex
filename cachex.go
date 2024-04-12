package cachex

import (
	"errors"
	"sync"
	"time"
)

// 应用内缓存
// 特点
// 1.全局单实例
// 到设定的点就删除

type cache struct {
	prefix string
	store  sync.Map
}

type cacheData struct {
	key    string
	data   interface{}
	expire time.Time
}

var ErrorEmpty error = errors.New("empty cache")

func NewCache(prefix string) *cache {

	c := &cache{
		prefix: prefix,
	}

	go func() {
		for {
			c.store.Range(func(key, value interface{}) bool {
				if value.(*cacheData).expire.Before(time.Now()) {
					c.store.Delete(key)
				}
				return true
			})
			time.Sleep(time.Second * 5)
		}
	}()

	return c
}

// 设置缓存
func (c *cache) Set(key string, value interface{}, expire time.Duration) {
	if expire == 0 {
		expire = time.Hour * 24 * 365
	}
	cd := &cacheData{key, value, time.Now().Add(expire)}
	key = c.prefix + key
	c.store.Store(key, cd)
}

// 读取缓存
func (c *cache) Get(key string) (interface{}, error) {
	key = c.prefix + key
	if v, ok := c.store.Load(key); ok {

		cc := v.(*cacheData)
		if cc.expire.Before(time.Now()) {
			c.store.Delete(key)
			return nil, ErrorEmpty
		}
		return cc.data, nil
	}
	return nil, ErrorEmpty
}

// 删除缓存
func (c *cache) Delete(key string) {
	key = c.prefix + key
	c.store.Delete(key)
}

// 清空缓存
func (c *cache) Clear() {
	c.store.Range(func(key, value interface{}) bool {
		// 根据前缀删除
		if key.(string)[:len(c.prefix)] == c.prefix {
			c.store.Delete(key)
		}
		return true
	})
}
