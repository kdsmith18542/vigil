// Copyright (c) 2024 The Vigil developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package lru

import (
	"container/list"
	"sync"
)

// Cache implements a least-recently-updated (LRU) cache with nearly O(1)
// lookups and inserts.  Items are added to the cache up to a limit, at which
// point further additions will evict the least recently added item.  The zero
// value is not valid and Caches must be created with NewCache.  All Cache
// methods are concurrent safe.
type Cache[T comparable] struct {
	mu    sync.Mutex
	m     map[T]*list.Element
	list  *list.List
	limit int
}

// NewCache creates an initialized and empty LRU cache.
func NewCache[T comparable](limit int) Cache[T] {
	return Cache[T]{
		m:     make(map[T]*list.Element, limit),
		list:  list.New(),
		limit: limit,
	}
}

// Add adds an item to the LRU cache, removing the oldest item if the new item
// is not already a member, or marking item as the most recently added item if
// it is already present.
func (c *Cache[T]) Add(item T) {
	defer c.mu.Unlock()
	c.mu.Lock()

	// Move this item to front of list if already present
	elem, ok := c.m[item]
	if ok {
		c.list.MoveToFront(elem)
		return
	}

	// If necessary, make room by popping an item off from the back
	if len(c.m) > c.limit {
		elem := c.list.Back()
		if elem != nil {
			v := c.list.Remove(elem)
			delete(c.m, v.(T))
		}
	}

	// Add new item to the LRU
	elem = c.list.PushFront(item)
	c.m[item] = elem
}

// Contains checks whether v is a member of the LRU cache.
func (c *Cache[T]) Contains(v T) bool {
	c.mu.Lock()
	_, ok := c.m[v]
	c.mu.Unlock()
	return ok
}
