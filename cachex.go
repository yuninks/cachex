package cachex

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"
)

// 应用内缓存
// 特点
// 1.全局单实例
// 到设定的点就删除

type CacheLocal struct {
	store  sync.Map
	cancel context.CancelFunc
}

type cacheData struct {
	key    string
	data   any
	expire time.Time
}

var ErrorEmpty error = errors.New("empty cache")

func NewCacheLocal() *CacheLocal {
	ctx, cancel := context.WithCancel(context.Background())
	c := &CacheLocal{cancel: cancel}

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				c.store.Range(func(key, value any) bool {
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
func (c *CacheLocal) Set(key string, value any, expire time.Duration) {
	if expire == 0 {
		expire = time.Hour * 24 * 365
	}
	cd := &cacheData{key, value, time.Now().Add(expire)}
	c.store.Store(key, cd)
}

// 读取缓存
func (c *CacheLocal) Get(key string) (any, error) {
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

// 根据结构体体反射读取缓存 dest必须是指针类型
func (c *CacheLocal) GetStruct(key string, dest any) error {

	v, err := c.Get(key)
	if err != nil {
		return err
	}
	jsonBytes, err := json.Marshal(v)
	if err != nil {
		return err
	}

	return json.Unmarshal(jsonBytes, dest)
}

// 删除缓存
func (c *CacheLocal) Delete(key string) {
	c.store.Delete(key)
}

// 清空缓存
func (c *CacheLocal) Clear() {
	c.store.Range(func(key, value any) bool {
		c.store.Delete(key)
		return true
	})
}

// 关闭缓存，停止清理协程
func (c *CacheLocal) Close() {
	if c.cancel != nil {
		c.cancel()
	}
}
