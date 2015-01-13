package shared

import "testing"

func TestNewMapCache(t *testing.T) {
	var c Cacher = NewMapCache(10)
	testCacher(c, t)
}

func TestNewArrayCache(t *testing.T) {
	var c Cacher = NewArrayCache(10)
	testCacher(c, t)
}

func testCacher(c Cacher, t testing.TB) {
	if v := c.Get(1, 2); v != -1 {
		t.Errorf("%T: wrong cacher value, got %v", c, v)
	}

	c.Set(1, 2, 1)
	if v := c.Get(1, 2); v != 1 {
		t.Errorf("%T: wrong cacher value, got %v", c, v)
	}

	// symmetric
	if v := c.Get(2, 1); v != 1 {
		t.Errorf("%T: wrong cacher value, got %v", c, v)
	}

	c.Set(9, 9, 5)
	if v := c.Get(9, 9); v != 5 {
		t.Errorf("%T: wrong cacher value, got %v", c, v)
	}
}
