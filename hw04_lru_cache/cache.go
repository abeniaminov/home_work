package hw04lrucache

import (
	"sync"
)

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
	sync.Mutex
}

func (lru *lruCache) Set(key Key, value interface{}) bool {
	lru.Lock()
	defer lru.Unlock()
	if len(lru.items) == lru.capacity && lru.capacity > 0 {
		lru.queue.Remove(lru.queue.Back())
	}
	if listItem, ok := lru.items[key]; ok {
		listItem.Value = value
		lru.queue.MoveToFront(listItem)
		return true
	} else {
		listItem = &ListItem{value, key, nil, nil}
		lru.queue.PushFront(listItem, key)
		lru.items[key] = listItem
	}
    return false
}

func (lru *lruCache) Get(key Key) (interface{},  bool) {
	lru.Lock()
	defer lru.Unlock()
	if listItem, ok := lru.items[key]; ok {
		lru.queue.MoveToFront(listItem)
		return listItem.Value, true
	}
	return nil, false
}

func (lru *lruCache) Clear()  {
	lru.Lock()
	defer lru.Unlock()
	lru.queue = NewList()
	lru.items = make(map[Key]*ListItem, lru.capacity)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
