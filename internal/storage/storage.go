package storage

import "sync"

type Cache struct {
	passwords []string
	sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{
		passwords: []string{},
	}
}

func (c *Cache) Get() []string {
	c.RLock()
	defer c.RUnlock()

	return c.passwords
}

func (c *Cache) Add(pwd string) {
	c.Lock()
	defer c.Unlock()

	c.passwords = append(c.passwords, pwd)
}

func (c *Cache) Exists(p string) bool {
	c.RLock()
	defer c.RUnlock()

	for _, pwd := range c.passwords {
		if p == pwd {
			return true
		}
	}

	return false
}
