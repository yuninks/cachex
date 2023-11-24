package cachex_test

import (
	"fmt"
	"testing"
	"time"

	"code.yun.ink/pkg/cachex"
)

func TestCache(t *testing.T) {
	cachex.NewCache().Set("test", "test", time.Second*5)

	da := cachex.NewCache().Get("test")
	fmt.Println(da)

	time.Sleep(time.Second *5)
	da = cachex.NewCache().Get("test")
	fmt.Println(da)
}
