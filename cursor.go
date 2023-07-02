package pile

// Least returns a new Cursor located at the Key which is less than all others
// in the Set. The return is false when Map is empty. A Delete or Insert renders
// the Cursor invalid.
func (keys *Set[Key]) Least() (Cursor[Key, struct{}], bool) { return keys.m.Least() }

// Most returns a new Cursor located at the Key which is more than all others in
// the Set. The return is false when Map is empty. A Delete or Insert renders
// the Cursor invalid.
func (keys *Set[Key]) Most() (Cursor[Key, struct{}], bool) { return keys.m.Most() }

// Least returns a new Cursor located at the Key which is less than all others
// in the Set. The return is false when Map is empty. A Delete or Insert renders
// the Cursor invalid.
func (m *Map[Key, Value]) Least() (Cursor[Key, Value], bool) {
	t := m.top
	if t == nil {
		return Cursor[Key, Value]{}, false
	}
	for t.subs[0] != nil {
		t = t.subs[0]
	}
	return Cursor[Key, Value]{t: t, pairI: 0}, true
}

// Most returns a new Cursor located at the Key which is more than all others in
// the Map. The return is false when Map is empty. A Delete or Insert renders
// the Cursor invalid.
func (m *Map[Key, Value]) Most() (Cursor[Key, Value], bool) {
	t := m.top
	if t == nil {
		return Cursor[Key, Value]{}, false
	}
	for t.subs[t.pairN&3] != nil {
		t = t.subs[t.pairN&3]
	}
	return Cursor[Key, Value]{t: t, pairI: t.pairN - 1}, true
}

// Cursor navigates over Sortable content.
type Cursor[Key Sortable, Value any] struct {
	t     *node[Key, Value]
	pairI int
}

// Key returns the Key at the current position.
func (c *Cursor[Key, Value]) Key() Key {
	if c.t == nil {
		var zero Key
		return zero
	}
	return c.t.pairs[c.pairI%3].K
}

// Value returns the Value at the current position.
func (c *Cursor[Key, Value]) Value() Value {
	if c.t == nil {
		var zero Value
		return zero
	}
	return c.t.pairs[c.pairI%3].V
}

// Swap sets the Value and it returns the previous one.
func (c *Cursor[Key, Value]) Swap(v Value) (previous Value) {
	if c.t == nil {
		var zero Value
		return zero
	}
	p := &c.t.pairs[c.pairI%3].V
	previous = *p
	*p = v
	return
}

// Ascend moves the Cursor one key closer to Most, up to Most itself.
func (c *Cursor[Key, Value]) Ascend() bool {
	if c.t == nil {
		return false
	}
	sub := c.t.subs[(c.pairI+1)&3]
	if sub != nil {
		// down to bottom level, left side
		for sub.subs[0] != nil {
			sub = sub.subs[0]
		}
		c.t = sub
		c.pairI = 0
		return true
	}

	if c.pairI+1 < c.t.pairN {
		c.pairI++
		return true
	}

	// move up
	t := c.t
	for t.above != nil {
		switch t {
		case t.above.subs[0]:
			c.t = t.above
			c.pairI = 0
			return true
		case t.above.subs[1]:
			if t.above.pairN < 2 {
				break
			}
			c.t = t.above
			c.pairI = 1
			return true
		case t.above.subs[2]:
			if t.above.pairN < 3 {
				break
			}
			c.t = t.above
			c.pairI = 2
			return true
		}
		t = t.above
	}
	return false
}

// Descend moves the Cursor one Key closer to Least, up to Least itself.
func (c *Cursor[Key, Value]) Descend() bool {
	if c.t == nil {
		return false
	}
	sub := c.t.subs[c.pairI&3]
	if sub != nil {
		// down to bottom level, right side
		for sub.subs[sub.pairN&3] != nil {
			sub = sub.subs[sub.pairN&3]
		}
		c.t = sub
		c.pairI = sub.pairN - 1
		return true
	}

	if c.pairI > 0 {
		c.pairI--
		return true
	}

	// move up
	t := c.t
	for t.above != nil {
		switch t {
		case t.above.subs[1]:
			c.t = t.above
			c.pairI = 0
			return true
		case t.above.subs[2]:
			c.t = t.above
			c.pairI = 1
			return true
		case t.above.subs[3]:
			c.t = t.above
			c.pairI = 2
			return true
		}
		t = t.above
	}
	return false
}
