package cachex

import (
	"sync"
	"time"
)

// 应用内缓存
// 特点
// 1.全局单实例
// 到设定的点就删除

type cache struct{
	store sync.Map
}
var one sync.Once
var c cache

type cacheData struct {
	key    string
	data   interface{}
	expire time.Time
}

func NewCache() *cache {
	one.Do(func() {
		c = cache{}

		go func(){
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
	})
	return &c
}

// 设置缓存
func (c *cache) Set(key string, value interface{}, expire time.Duration) {
	cd := &cacheData{key, value, time.Now().Add(expire)}
	c.store.Store(key, cd)
}

// 读取缓存
func (c *cache) Get(key string) interface{} {
	if v, ok := c.store.Load(key); ok {

		cc := v.(*cacheData)
		if cc.expire.Before(time.Now()) {
			c.store.Delete(key)
			return nil
		}
		return cc.data
	}
	return nil
}

// 删除缓存
func (c *cache) Delete(key string) {
	c.store.Delete(key)
}
