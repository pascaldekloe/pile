package pile

// FindPointer returns the Value assigned to the Key, with nil for none. The
// return becomes undefined after any mutation to the Map. Use with caution.
func (m *Map[Key, Value]) FindPointer(k Key) *Value {
	t := m.top
	for t != nil {
		switch t.pairN {
		case 3:
			switch {
			case k > t.pairs[1].K:
				switch {
				case k > t.pairs[2].K:
					t = t.subs[3]
				case k < t.pairs[2].K:
					t = t.subs[2]
				default:
					return &t.pairs[2].V
				}
			case k <= t.pairs[0].K:
				if k < t.pairs[0].K {
					t = t.subs[0]
				} else {
					return &t.pairs[0].V
				}
			case k < t.pairs[1].K:
				t = t.subs[1]
			default:
				return &t.pairs[1].V
			}

		case 2:
			switch {
			case k >= t.pairs[1].K:
				if k > t.pairs[1].K {
					t = t.subs[2]
				} else {
					return &t.pairs[1].V
				}
			case k <= t.pairs[0].K:
				if k < t.pairs[0].K {
					t = t.subs[0]
				} else {
					return &t.pairs[0].V
				}
			default:
				t = t.subs[1]
			}

		default:
			switch {
			case k > t.pairs[0].K:
				t = t.subs[1]
			case k < t.pairs[0].K:
				t = t.subs[0]
			default:
				return &t.pairs[0].V
			}
		}
	}

	return nil // not found
}

// Find returns the Value assigned to the Key.
func (m *Map[Key, Value]) Find(k Key) (Value, bool) {
	vp := m.FindPointer(k)
	if vp == nil {
		var zero Value
		return zero, false
	}
	return *vp, true
}

// Update assigns the Value to the Key if and only if the Key is present.
func (m *Map[Key, Value]) Update(k Key, v Value) bool {
	vp := m.FindPointer(k)
	if vp == nil {
		return false
	}
	*vp = v
	return true
}
