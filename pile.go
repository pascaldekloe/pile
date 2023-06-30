// Package pile provides sorted memory structures.
package pile

// Sortable is a key constraint.
type Sortable interface {
	~string |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~int | ~int8 | ~int16 | ~int32 | ~int64
}

// Pair is a Map entry.
type Pair[Key Sortable, Value any] struct {
	K Key
	V Value
}

// A node holds up to tree Pairs in ascending Key order.
// Nodes on ground level do not have any subnodes.
// Higher nodes stack on pairN plus one subnodes.
type node[Key Sortable, Value any] struct {
	above *node[Key, Value]
	pairN int                  // actual pairs count
	subs  [4]*node[Key, Value] // directly under
	pairs [3]Pair[Key, Value]  // own entries
}

func newNode[Key Sortable, Value any](above *node[Key, Value], p Pair[Key, Value]) *node[Key, Value] {
	t := node[Key, Value]{above: above, pairN: 1}
	t.pairs[0] = p
	return &t
}

func newNodeOfTwo[Key Sortable, Value any](above *node[Key, Value], p1, p2 Pair[Key, Value]) *node[Key, Value] {
	t := node[Key, Value]{above: above, pairN: 2}
	t.pairs[0] = p1
	t.pairs[1] = p2
	return &t
}

// Map provides sorted Key–Value registration. The zero Map is empty and ready
// for use. Do not copy the Map struct.
//
// The best-case for storage-overhead per Pair is 16 bytes on 64-bit platforms.
// The worst-case per Pair is 48 bytes plus the size of another 2 Pairs. E.g.,
// the Pair[int, string] costs 24 bytes, which makes an overhead of between ⅔
// and 4 times the Pair size.
type Map[Key Sortable, Value any] struct {
	check noCopy

	top *node[Key, Value]

	// reusable buffer for level push
	split Pair[Key, Value]
}

// Size returns the number of Keys in the Map.
func (m *Map[Key, Value]) Size() int {
	return m.top.size() // nil safe
}

func (t *node[Key, Value]) size() int {
	if t == nil {
		return 0
	}
	size := t.pairN
	switch t.pairN {
	case 3:
		size += t.subs[3].size()
		fallthrough
	case 2:
		size += t.subs[2].size()
		fallthrough
	default:
		size += t.subs[1].size()
		size += t.subs[0].size()
	}
	return size
}

func (m *Map[Key, Value]) AppendPairs(dst []Pair[Key, Value]) []Pair[Key, Value] {
	return m.top.appendPairs(dst) // nil safe
}

func (t *node[Key, Value]) appendPairs(dst []Pair[Key, Value]) []Pair[Key, Value] {
	if t != nil {
		dst = t.subs[0].appendPairs(dst)
		dst = append(dst, t.pairs[0])
		dst = t.subs[1].appendPairs(dst)
		if t.pairN > 1 {
			dst = append(dst, t.pairs[1])
			dst = t.subs[2].appendPairs(dst)
			if t.pairN > 2 {
				dst = append(dst, t.pairs[2])
				dst = t.subs[3].appendPairs(dst)
			}
		}
	}
	return dst
}

// Set provides sorted Key registration. The zero Set is empty and ready for
// use. Do not copy the Set struct.
//
// The best-case for storage-overhead per Key is 16 bytes on 64-bit platforms.
// The worst-case per Key is 48 bytes plus the size of another 2 Keys. E.g., the
// uint costs 8 bytes, which makes an overhead of between 2 and 8 times the Key
// size.
type Set[Key Sortable] struct {
	m Map[Key, struct{}]
}

// Size returns the number of Keys in the Set.
func (keys *Set[Key]) Size() int {
	return keys.m.Size()
}

// Find returns the Key's presence in the Set.
func (keys *Set[Key]) Find(k Key) bool {
	return keys.m.FindPointer(k) != nil
}

// Insert adds the Key to the Set if and only if the Key is absent.
func (keys *Set[Key]) Insert(entry Key) bool {
	return keys.m.Insert(entry, struct{}{})
}

// NoCopy triggers go(1) vet when copied after the first use.
// See https://golang.org/issues/8005#issuecomment-190753527 for details.
type noCopy struct{}

// Lock triggers the -copylocks checker from go(1) vet.
func (*noCopy) Lock() {}

// Unlock triggers the -copylocks checker from go(1) vet.
func (*noCopy) Unlock() {}
