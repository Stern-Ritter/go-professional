package hw04lrucache

import "sync"

type Key string

type Entry struct {
	key   Key
	value interface{}
}

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	muQueue  sync.Mutex
	muItems  sync.RWMutex
	queue    List
	items    map[Key]*ListItem
	capacity int
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	if c.capacity <= 0 {
		return false
	}

	c.muItems.Lock()
	defer c.muItems.Unlock()
	c.muQueue.Lock()
	defer c.muQueue.Unlock()

	if item, ok := c.items[key]; ok {
		c.queue.MoveToFront(item)
		entry := item.Value.(*Entry)
		entry.value = value
		return true
	}

	if len(c.items) >= c.capacity {
		last := c.queue.Back()
		entry := last.Value.(*Entry)
		delete(c.items, entry.key)
		c.queue.Remove(last)
	}

	item := c.queue.PushFront(&Entry{key: key, value: value})
	c.items[key] = item
	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.muItems.RLock()
	defer c.muItems.RUnlock()

	if item, ok := c.items[key]; ok {
		c.muQueue.Lock()
		defer c.muQueue.Unlock()

		c.queue.MoveToFront(item)
		entry := item.Value.(*Entry)
		return entry.value, true
	}

	return nil, false
}

func (c *lruCache) Clear() {
	c.muItems.Lock()
	defer c.muItems.Unlock()
	c.muQueue.Lock()
	defer c.muQueue.Unlock()

	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
