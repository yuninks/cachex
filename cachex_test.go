package cachex_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/yuninks/cachex"
)

func TestCache(t *testing.T) {
	c := cachex.NewCache()

	c.Set("test", "test", time.Second*5)

	da, err := c.Get("test")
	fmt.Println(da, err)

	time.Sleep(time.Second * 5)
	da, err = c.Get("test")
	fmt.Println(da, err)
}
