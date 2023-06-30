package pile

// First either returns a new Cursor at the Key that is more than all other keys
// from Map, or it returns false when the Map is empty. The Cursor becomes
// invalid after any insertion or deletion.
func (keys *Set[Key]) First() (Cursor[Key, struct{}], bool) { return keys.m.First() }

// Last either returns a new Cursor at the Key that is more than all other keys
// from Map, or it returns false when the Map is empty. The Cursor becomes
// invalid after any insertion or deletion.
func (keys *Set[Key]) Last() (Cursor[Key, struct{}], bool) { return keys.m.Last() }

// First either returns a new Cursor at the Key that is more than all other keys
// from Map, or it returns false when the Map is empty. The Cursor becomes
// invalid after any insertion or deletion.
func (m *Map[Key, Value]) First() (Cursor[Key, Value], bool) {
	t := m.top
	if t == nil {
		return Cursor[Key, Value]{}, false
	}
	for t.subs[0] != nil {
		t = t.subs[0]
	}
	return Cursor[Key, Value]{t: t, pairI: 0}, true
}

// Last either returns a new Cursor at the Key that is more than all other keys
// from Map, or it returns false when the Map is empty. The Cursor becomes
// invalid after any insertion or deletion.
func (m *Map[Key, Value]) Last() (Cursor[Key, Value], bool) {
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

// Fetch returns the current position.
func (c *Cursor[Key, Value]) Fetch() Pair[Key, Value] {
	return c.t.pairs[c.pairI%3]
}

// Update assigns the Value to the Key.
func (c *Cursor[Key, Value]) Update(v Value) {
	c.t.pairs[c.pairI%3].V = v
}

// Forward moves the Cursor one step closer towards the Last if possible.
func (c *Cursor[Key, Value]) Forward() bool {
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

// Backward moves the Cursor one step closer towards the First if possible.
func (c *Cursor[Key, Value]) Backward() bool {
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
