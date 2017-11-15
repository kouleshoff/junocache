package main

import (
  "time"
)

type Item interface {
  Expiration() int64
  Size() int
  GetFirst() ([]byte,bool)
}

type Cache struct {
  Items map[string]Item
}

// Returns true if the item has expired.
func Expired(item Item) bool {
	if item.Expiration() == 0 {
		return false
	}
	return time.Now().UnixNano() > item.Expiration()
}

func (c *Cache) FindByKey(key string) (Item,bool) {
  item, ok := c.Items[key]
  if ok && Expired(item) {
    ok = false
  }
  return item, ok
}
