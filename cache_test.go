package lru

import (
	"fmt"
	"testing"
)

func TestLRU(t *testing.T) {
	func() {
		defer func() {
			if recover() == nil {
				t.Fatal()
			}
		}()
		New(0, nil)
	}()
	var evicted []string
	c := New(5, func(key string, value interface{}) {
		if value != key {
			t.Fatal()
		}
		evicted = append(evicted, key)
	})
	nums := []string{
		"7", "3", "2", "4", "5", "6", "7",
		"1", "2", "3", "4", "6", "1", "4",
	}
	for _, num := range nums {
		c.Set(num, num)
	}
	if c.Len() != 5 {
		t.Fatalf("got %v want %v", c.Len(), 5)
	}
	sres := fmt.Sprintf("%s", evicted)
	if sres != "[7 3 2 4 5 6 7]" {
		t.Fatalf("got %v want %v", sres, "[7 3 2 4 5 6 7]")
	}
	sres = fmt.Sprintf("%s", c.all())
	if sres != "[4 1 6 3 2]" {
		t.Fatalf("got %v want %v", sres, "[4 1 6 3 2]")
	}
	nums = []string{"7", "3", "2", "4", "5", "6", "7"}
	var res []interface{}
	for _, num := range nums {
		res = append(res, c.Get(num))
	}
	sres = fmt.Sprintf("%s", res)
	if sres != "[<nil> 3 2 4 <nil> 6 <nil>]" {
		t.Fatalf("got %v want %v", sres, "[<nil> 3 2 4 <nil> 6 <nil>]")
	}
	sres = fmt.Sprintf("%s", c.all())
	if sres != "[6 4 2 3 1]" {
		t.Fatalf("got %v want %v", sres, "[6 4 2 3 1]")
	}
	nums = []string{"1", "6", "2", "2", "6", "1", "7"}
	for _, num := range nums {
		c.Delete(num)
	}
	sres = fmt.Sprintf("%s", c.all())
	if sres != "[4 3]" {
		t.Fatalf("got %v want %v", sres, "[4 3]")
	}
	if c.Len() != 2 {
		t.Fatalf("got %v want %v", c.Len(), 2)
	}
	c.Delete("4")
	c.Delete("3")
	if c.Len() != 0 {
		t.Fatalf("got %v want %v", c.Len(), 0)
	}
}

// all returns all keys
func (c *Cache) all() []string {
	var a []string
	e := c.head
	for e != nil {
		a = append(a, e.key)
		e = e.next
	}
	return a
}
