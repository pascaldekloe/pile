package pile

// Insert assigns the Value to the Key if and only if the Key is absent.
func (m *Map[Key, Value]) Insert(k Key, v Value) bool {
	if m.top == nil {
		m.top = m.newNodeWith1(nil, pair[Key, Value]{K: k, V: v})
		return true
	}
	t := m.top

	var splitRight *node[Key, Value]

	// search t and below
	for {
		if t.pairN == 3 {
			// insert overflows
			switch {
			case k > t.pairs[1].K:
				switch {
				case k > t.pairs[2].K:
					if t.subs[3] != nil {
						t = t.subs[3]
						continue
					}
					goto InsertFourthOverflow
				case k < t.pairs[2].K:
					if t.subs[2] != nil {
						t = t.subs[2]
						continue
					}
					goto InsertThirdOverflow
				}
				return false
			case k <= t.pairs[0].K:
				if k < t.pairs[0].K {
					if t.subs[0] != nil {
						t = t.subs[0]
						continue
					}
					goto InsertFirstOverflow
				}
				return false
			case k < t.pairs[1].K:
				if t.subs[1] != nil {
					t = t.subs[1]
					continue
				}
				goto InsertSecondOverflow
			}
			return false
		}

		switch {
		case t.pairN == 2 && k >= t.pairs[1].K:
			if k > t.pairs[1].K {
				if t.subs[2] != nil {
					t = t.subs[2]
					continue
				}
				t.pairs[2].K = k
				t.pairs[2].V = v
				t.pairN++
				return true
			}
			return false

		case k > t.pairs[0].K:
			if t.subs[1] != nil {
				t = t.subs[1]
				continue
			}
			t.pairs[2] = t.pairs[1] // redundant if pairN is 1
			t.pairs[1].K = k
			t.pairs[1].V = v
			t.pairN++
			return true

		case k < t.pairs[0].K:
			if t.subs[0] != nil {
				t = t.subs[0]
				continue
			}
			t.pairs[2] = t.pairs[1] // redundant if pairN is 1
			t.pairs[1] = t.pairs[0]
			t.pairs[0].K = k
			t.pairs[0].V = v
			t.pairN++
			return true
		}
		return false
	}

InsertFirstOverflow:
	m.split = t.pairs[0]
	t.pairs[0].K = k
	t.pairs[0].V = v
	goto Branche2Right
InsertSecondOverflow:
	m.split.K = k
	m.split.V = v
Branche2Right:
	t.pairN = 1
	splitRight = m.newNodeWith2(t.above, t.pairs[1], t.pairs[2])
	goto Overflow

InsertThirdOverflow:
	t.pairN = 2
	m.split.K = k
	m.split.V = v
	splitRight = m.newNodeWith1(t.above, t.pairs[2])
	goto Overflow

InsertFourthOverflow:
	t.pairN = 2
	m.split = t.pairs[2]
	splitRight = m.newNodeWith1(t.above, pair[Key, Value]{K: k, V: v})

Overflow:
	for t.above != nil {
		above := t.above
		splitRight = m.takeSplit(above, t, splitRight, &m.split)
		if splitRight == nil {
			return true
		}
		t = above
	}

	grow := m.newNodeWith1(nil, m.split)
	m.top.above = grow
	splitRight.above = grow
	grow.subs[0] = m.top
	grow.subs[1] = splitRight
	m.top = grow
	return true
}

// Put assigns the Value to the Key regardles whether the Key is present or not.
// The result is equivalent to both m.Insert(k, v) || m.Update(k, v), and to
// m.Update(k, v) || m.Insert(k, v).
func (m *Map[Key, Value]) Put(k Key, v Value) {
	if m.top == nil {
		m.top = m.newNodeWith1(nil, pair[Key, Value]{K: k, V: v})
		return
	}
	t := m.top

	var splitRight *node[Key, Value]

	// search t and below
	for {
		if t.pairN == 3 {
			// insert overflows
			switch {
			case k > t.pairs[1].K:
				switch {
				case k > t.pairs[2].K:
					if t.subs[3] != nil {
						t = t.subs[3]
						continue
					}
					goto InsertFourthOverflow
				case k < t.pairs[2].K:
					if t.subs[2] != nil {
						t = t.subs[2]
						continue
					}
					goto InsertThirdOverflow
				}
				t.pairs[2].V = v // update
				return
			case k <= t.pairs[0].K:
				if k < t.pairs[0].K {
					if t.subs[0] != nil {
						t = t.subs[0]
						continue
					}
					goto InsertFirstOverflow
				}
				t.pairs[0].V = v // update
				return
			case k < t.pairs[1].K:
				if t.subs[1] != nil {
					t = t.subs[1]
					continue
				}
				goto InsertSecondOverflow
			}
			t.pairs[1].V = v // update
			return
		}

		switch {
		case t.pairN == 2 && k >= t.pairs[1].K:
			if k > t.pairs[1].K {
				if t.subs[2] != nil {
					t = t.subs[2]
					continue
				}
				t.pairs[2].K = k
				t.pairs[2].V = v
				t.pairN++
				return
			}
			t.pairs[1].V = v // update
			return

		case k > t.pairs[0].K:
			if t.subs[1] != nil {
				t = t.subs[1]
				continue
			}
			t.pairs[2] = t.pairs[1] // redundant if pairN is 1
			t.pairs[1].K = k
			t.pairs[1].V = v
			t.pairN++
			return

		case k < t.pairs[0].K:
			if t.subs[0] != nil {
				t = t.subs[0]
				continue
			}
			t.pairs[2] = t.pairs[1] // redundant if pairN is 1
			t.pairs[1] = t.pairs[0]
			t.pairs[0].K = k
			t.pairs[0].V = v
			t.pairN++
			return
		}
		t.pairs[0].V = v // update
		return
	}

InsertFirstOverflow:
	m.split = t.pairs[0]
	t.pairs[0].K = k
	t.pairs[0].V = v
	goto Branche2Right
InsertSecondOverflow:
	m.split.K = k
	m.split.V = v
Branche2Right:
	t.pairN = 1
	splitRight = m.newNodeWith2(t.above, t.pairs[1], t.pairs[2])
	goto Overflow

InsertThirdOverflow:
	t.pairN = 2
	m.split.K = k
	m.split.V = v
	splitRight = m.newNodeWith1(t.above, t.pairs[2])
	goto Overflow

InsertFourthOverflow:
	t.pairN = 2
	m.split = t.pairs[2]
	splitRight = m.newNodeWith1(t.above, pair[Key, Value]{K: k, V: v})

Overflow:
	for t.above != nil {
		above := t.above
		splitRight = m.takeSplit(above, t, splitRight, &m.split)
		if splitRight == nil {
			return
		}
		t = above
	}

	grow := m.newNodeWith1(nil, m.split)
	m.top.above = grow
	splitRight.above = grow
	grow.subs[0] = m.top
	grow.subs[1] = splitRight
	m.top = grow
}

// TakeSplit adds node rightInsert next to fromSub in t, separated by the split.
// The operation may cause another split (pointer update) with a new splitRight
// (relative to t).
func (m *Map[Key, Value]) takeSplit(t, fromSub, rightInsert *node[Key, Value], split *pair[Key, Value]) (splitRight *node[Key, Value]) {
	if t.pairN < 3 { // fits in node
		t.pairN++
		switch fromSub {
		case t.subs[0]:
			t.subs[3] = t.subs[2]
			t.subs[2] = t.subs[1]
			t.subs[1] = rightInsert
			t.pairs[2] = t.pairs[1]
			t.pairs[1] = t.pairs[0]
			t.pairs[0] = *split
		case t.subs[1]:
			t.subs[3] = t.subs[2]
			t.subs[2] = rightInsert
			t.pairs[2] = t.pairs[1]
			t.pairs[1] = *split
		case t.subs[2]:
			t.subs[3] = rightInsert
			t.pairs[2] = *split
		}

		return nil
	}
	// node has no place for insert

	switch fromSub {
	case t.subs[0]: // rightInsert goes into second spot
		splitRight = m.newNodeWith2(t.above, t.pairs[1], t.pairs[2])
		if t.subs[1] != nil {
			splitRight.subs[0] = t.subs[1]
			splitRight.subs[0].above = splitRight
		}
		if t.subs[2] != nil {
			splitRight.subs[1] = t.subs[2]
			splitRight.subs[1].above = splitRight
		}
		if t.subs[3] != nil {
			splitRight.subs[2] = t.subs[3]
			splitRight.subs[2].above = splitRight
		}
		t.subs[1] = rightInsert
		t.pairN = 1
		*split, t.pairs[0] = t.pairs[0], *split
	case t.subs[1]: // rightInsert goes into third sport
		t.pairN = 1
		splitRight = m.newNodeWith2(t.above, t.pairs[1], t.pairs[2])
		splitRight.subs[0] = rightInsert
		splitRight.subs[0].above = splitRight
		if t.subs[2] != nil {
			splitRight.subs[1] = t.subs[2]
			splitRight.subs[1].above = splitRight
		}
		if t.subs[3] != nil {
			splitRight.subs[2] = t.subs[3]
			splitRight.subs[2].above = splitRight
		}
		// pass split to upper level
	case t.subs[2]: // rightInsert goes into fourth spot
		t.pairN = 2
		splitRight = m.newNodeWith1(t.above, t.pairs[2])
		splitRight.subs[0] = rightInsert
		splitRight.subs[0].above = splitRight
		if t.subs[3] != nil {
			splitRight.subs[1] = t.subs[3]
			splitRight.subs[1].above = splitRight
		}
		// pass split to upper level
	case t.subs[3]: // rightInsert goes into fifth spot
		t.pairN = 2
		splitRight = m.newNodeWith1(t.above, *split)
		if t.subs[3] != nil {
			splitRight.subs[0] = t.subs[3]
			splitRight.subs[0].above = splitRight
		}
		splitRight.subs[1] = rightInsert
		splitRight.subs[1].above = splitRight
		*split = t.pairs[2]
	}
	return splitRight
}
