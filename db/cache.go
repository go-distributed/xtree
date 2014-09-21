package db

type cache struct {
	c map[revpath]Value
}

type revpath struct {
	rev  int
	path string
}

func newCache() *cache {
	return &cache{c: make(map[revpath]Value)}
}

func (c *cache) put(rp revpath, v Value) {
	c.c[rp] = v
}

func (c *cache) get(rp revpath) (Value, bool) {
	v, ok := c.c[rp]
	return v, ok
}
