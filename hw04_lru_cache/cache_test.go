package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type entry struct {
	Key   string
	Value int
}

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

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
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

	t.Run("pushing items out of cache due to cache capacity", func(t *testing.T) {
		capacity := 5
		c := NewCache(capacity)

		entries := []entry{
			{"a", 10},
			{"b", 20},
			{"c", 30},
			{"d", 40},
			{"e", 50},
			{"f", 60},
			{"g", 70},
		}

		for _, e := range entries {
			inCache := c.Set(Key(e.Key), e.Value)
			assert.False(t, inCache, "value for key: %s should not be in cache", e.Key)
		}

		pushedOutEntriesCount := len(entries) - capacity
		pushedOutEntries := entries[:pushedOutEntriesCount]
		for _, e := range pushedOutEntries {
			v, inCache := c.Get(Key(e.Key))
			assert.False(t, inCache, "value for key: %s should be pushed out from cache", e.Key)
			assert.Empty(t, v, "value for key: %s should be empty but got: %d", e.Key, v)
		}

		expectedEntries := entries[pushedOutEntriesCount:]
		for _, e := range expectedEntries {
			v, inCache := c.Get(Key(e.Key))
			assert.True(t, inCache, "value for key: %s should be in cache", e.Key)
			assert.Equal(t, e.Value, v, "value for key: %s should be %d but got: %d", e.Key, e.Value, v)
		}
	})

	t.Run("pushing last used items out of cache", func(t *testing.T) {
		capacity := 3
		c := NewCache(capacity)

		entry1 := entry{"a", 10}
		entry2 := entry{"b", 20}
		entry3 := entry{"c", 30}
		entries := []entry{entry1, entry2, entry3}

		entry4 := entry{"d", 40}

		for _, e := range entries {
			inCache := c.Set(Key(e.Key), e.Value)
			assert.False(t, inCache, "value for key: %s should not be in cache", e.Key)
		}

		c.Get(Key(entry1.Key))
		c.Set(Key(entry2.Key), entry2.Value)
		c.Set(Key(entry4.Key), entry4.Value)

		pushedOutEntries := []entry{entry3}
		for _, e := range pushedOutEntries {
			v, inCache := c.Get(Key(e.Key))
			assert.False(t, inCache, "value for key: %s should be pushed out from cache", e.Key)
			assert.Empty(t, v, "value for key: %s should be empty but got: %d", e.Key, v)
		}

		expectedEntries := []entry{entry1, entry2, entry4}
		for _, e := range expectedEntries {
			v, inCache := c.Get(Key(e.Key))
			assert.True(t, inCache, "value for key: %s should be in cache", e.Key)
			assert.Equal(t, e.Value, v, "value for key: %s should be %d but got: %d", e.Key, e.Value, v)
		}
	})
}

func TestCacheMultithreading(_ *testing.T) {
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
}
