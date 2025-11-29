package cachex

import (
	"context"
	"errors"
	"sync"
	"time"
)

// 应用内缓存
// 特点
// 1.全局单实例
// 到设定的点就删除

type Cache struct {
	store  sync.Map
	cancel context.CancelFunc
}

type cacheData struct {
	key    string
	data   interface{}
	expire time.Time
}

var ErrorEmpty error = errors.New("empty cache")

func NewCache() *Cache {
	ctx, cancel := context.WithCancel(context.Background())
	c := &Cache{cancel: cancel}

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				c.store.Range(func(key, value interface{}) bool {
					cd, ok := value.(*cacheData)
					if !ok {
						return true
					}
					if cd.expire.Before(time.Now()) {
						c.store.Delete(key)
					}
					return true
				})
			}
		}
	}()

	return c
}

// 设置缓存
func (c *Cache) Set(key string, value interface{}, expire time.Duration) {
	if expire == 0 {
		expire = time.Hour * 24 * 365
	}
	cd := &cacheData{key, value, time.Now().Add(expire)}
	c.store.Store(key, cd)
}

// 读取缓存
func (c *Cache) Get(key string) (interface{}, error) {
	if v, ok := c.store.Load(key); ok {
		cc, ok := v.(*cacheData)
		if !ok {
			return nil, ErrorEmpty
		}
		if cc.expire.Before(time.Now()) {
			c.store.Delete(key)
			return nil, ErrorEmpty
		}
		return cc.data, nil
	}
	return nil, ErrorEmpty
}

// 删除缓存
func (c *Cache) Delete(key string) {
	c.store.Delete(key)
}

// 清空缓存
func (c *Cache) Clear() {
	c.store.Range(func(key, value interface{}) bool {
		c.store.Delete(key)
		return true
	})
}

// 关闭缓存，停止清理协程
func (c *Cache) Close() {
	if c.cancel != nil {
		c.cancel()
	}
}
