package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		val, ok := c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("evict one logic", func(t *testing.T) {
		c := NewCache(1)

		c.Set("aaa", 100)
		c.Set("bbb", 200)
		c.Set("ccc", 300)
		c.Set("ddd", 400) //ddd

		val, ok := c.Get("aaa")
		require.False(t, ok)
		require.Nil(t, val)

		val, ok = c.Get("ddd")
		require.True(t, ok)
		require.Equal(t, 400, val)
	})

	t.Run("evict logic", func(t *testing.T) {
		c := NewCache(3)

		c.Set("aaa", 100) // aaa
		c.Set("bbb", 200) // bbb aaa
		c.Set("ccc", 300) // ccc bbb aaa
		c.Set("ddd", 400) //ddd ccc bbb

		val, ok := c.Get("aaa")
		require.False(t, ok)
		require.Nil(t, val)

		c.Get("bbb") // bbb ddd ccc

		c.Get("ddd") // ddd bbb ccc

		val, ok = c.Get("ddd")
		require.True(t, ok)
		require.Equal(t, 400, val) // ddd bbb ccc

		c.Set("eee", 500) // eee ddd bbb

		_, ok = c.Get("bbb")
		require.True(t, ok)

		_, ok = c.Get("ccc")
		require.False(t, ok)

		_, ok = c.Get("eee")
		require.True(t, ok)

		_, ok = c.Get("ddd")
		require.True(t, ok)
	})

	t.Run("purge logic", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		c.Clear()
		val, ok = c.Get("aaa")
		require.False(t, ok)
		require.Nil(t, val)
	})
}

func TestCacheMultithreading(t *testing.T) {
	t.Run("multithreading logic", func(t *testing.T) {
		c := NewCache(10)
		wg := &sync.WaitGroup{}
		wg.Add(2)

		go func() {
			defer wg.Done()
			for i := 0; i < 1_000_000; i++ {
				c.Set(Key(strconv.Itoa(i)), i)
			}
		}()

		go func() {
			defer wg.Done()
			for i := 0; i < 1_000_000; i++ {
				c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
			}
		}()

		wg.Wait()
	})
}
