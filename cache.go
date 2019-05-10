/* cache.go -- lru cache type for string keys
 *
 * Copyright 2019, Joshua J Baker
 *
 * Permission to use, copy, modify, and/or distribute this software for any
 * purpose with or without fee is hereby granted, provided that the above
 * copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY
 * SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN ACTION
 * OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF OR IN
 * CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 */

package lru

type entry struct {
	key   string
	value interface{}
	prev  *entry
	next  *entry
}

// Cache is a standard non-thread safe fixed-sized lru cache where the
// keys are strings and the values are interfaces.
type Cache struct {
	entries map[string]*entry
	size    int
	onEvict func(key string, value interface{})
	head    *entry
	tail    *entry
}

// New returns a standard non-thread safe fixed-sized lru cache where the
// keys are strings and the values are interfaces.
// The size must be a positive number. The onEvict param is an optional
// callback function that fires when entries are forced to be evicted.
func New(size int, onEvict func(key string, value interface{})) *Cache {
	if size <= 0 {
		panic("invalid size")
	}
	return &Cache{
		entries: make(map[string]*entry, int(float64(size)*1.5)),
		size:    size,
		onEvict: onEvict,
	}
}

// Set a cache entry.
func (c *Cache) Set(key string, value interface{}) {
	e := c.entries[key]
	if e == nil {
		e = &entry{key: key, value: value}
		c.entries[key] = e
		if c.head == nil {
			c.head = e
			c.tail = e
		} else {
			c.head.prev = e
			e.next = c.head
			c.head = e
		}
	} else {
		e.value = value
		c.promote(e)
	}
	if len(c.entries) > c.size {
		evicted := c.tail
		delete(c.entries, c.tail.key)
		c.tail = c.tail.prev
		c.tail.next = nil
		if c.onEvict != nil {
			c.onEvict(evicted.key, evicted.value)
		}
	}
}

func (c *Cache) promote(e *entry) {
	if c.head != e {
		if c.tail == e {
			c.tail = c.tail.prev
			c.tail.next = nil
		} else {
			e.prev.next = e.next
			e.next.prev = e.prev
		}
		e.prev = nil
		e.next = c.head
		c.head.prev = e
		c.head = e
	}
}

// Len returns the number of entries in cache.
func (c *Cache) Len() int {
	return len(c.entries)
}

// Get an entry from cache.
func (c *Cache) Get(key string) interface{} {
	e := c.entries[key]
	if e == nil {
		return nil
	}
	c.promote(e)
	return e.value
}

// Delete an entry from cache.
func (c *Cache) Delete(key string) {
	e := c.entries[key]
	if e == nil {
		return
	}
	delete(c.entries, key)
	if len(c.entries) == 0 {
		c.head = nil
		c.tail = nil
	} else if e == c.head {
		c.head = c.head.next
		c.head.prev = nil
	} else if e == c.tail {
		c.tail = c.tail.prev
		c.tail.next = nil
	} else {
		e.prev.next = e.next
		e.next.prev = e.prev
	}
}
