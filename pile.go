// Package pile provides sorted memory structures.
package pile

// Sortable is a key constraint.
type Sortable interface {
	~string |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~int | ~int8 | ~int16 | ~int32 | ~int64
}

// Pair is a Map entry.
type pair[Key Sortable, Value any] struct {
	K Key
	V Value
}

// A node holds up to tree pairs in ascending Key order.
// Nodes on ground level do not have any subnodes.
// Higher nodes stack on pairN plus one subnodes.
type node[Key Sortable, Value any] struct {
	above *node[Key, Value]
	pairN int                  // actual pairs count
	subs  [4]*node[Key, Value] // directly under
	pairs [3]pair[Key, Value]  // own entries
}

func (m *Map[Key, Value]) newNodeWith1(above *node[Key, Value], p pair[Key, Value]) *node[Key, Value] {
	t := m.newNode()
	t.above = above
	t.pairs[0] = p
	t.pairN = 1
	return t
}

func (m *Map[Key, Value]) newNodeWith2(above *node[Key, Value], p1, p2 pair[Key, Value]) *node[Key, Value] {
	t := m.newNode()
	t.above = above
	t.pairs[0] = p1
	t.pairs[1] = p2
	t.pairN = 2
	return t
}

// NodeBatchN sets the number of nodes allocated together.
const nodeBatchN = 512 // must be a power of two

func (m *Map[Key, Value]) newNode() *node[Key, Value] {
	if m.nodeN == 0 || m.nodeQ == nil {
		m.nodeQ = new([nodeBatchN]node[Key, Value])
		m.nodeN = nodeBatchN
	}
	m.nodeN--
	return &m.nodeQ[m.nodeN&(nodeBatchN-1)]
}

// Map provides sorted Key–Value registration. The zero Map is empty and ready
// for use. Do not copy the Map struct.
//
// The best-case for storage-overhead per pair is 16 bytes on 64-bit platforms.
// The worst-case per pair is 48 bytes plus the size of another 2 pairs. E.g.,
// the pair[int, string] costs 24 bytes, which makes an overhead of between ⅔
// and 4 times the pair size.
type Map[Key Sortable, Value any] struct {
	check noCopy

	top *node[Key, Value]

	// reusable buffer for level push
	split pair[Key, Value]

	// allocation pool
	nodeN int
	nodeQ *[nodeBatchN]node[Key, Value]
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

// AppendKeys appends each Key in the Map to dst, ascending in Key order, and it
// returns the extended buffer.
func (m *Map[Key, Value]) AppendKeys(dst []Key) []Key {
	return m.top.appendKeys(dst)
}

func (t *node[Key, Value]) appendKeys(dst []Key) []Key {
	if t != nil {
		dst = t.subs[0].appendKeys(dst)
		dst = append(dst, t.pairs[0].K)
		dst = t.subs[1].appendKeys(dst)
		if t.pairN > 1 {
			dst = append(dst, t.pairs[1].K)
			dst = t.subs[2].appendKeys(dst)
			if t.pairN > 2 {
				dst = append(dst, t.pairs[2].K)
				dst = t.subs[3].appendKeys(dst)
			}
		}
	}
	return dst
}

// Append appends each Value in the Map to dst, ascending in Key order, and it
// returns the extended buffer.
func (m *Map[Key, Value]) AppendValues(dst []Value) []Value {
	return m.top.appendValues(dst)
}

func (t *node[Key, Value]) appendValues(dst []Value) []Value {
	if t != nil {
		dst = t.subs[0].appendValues(dst)
		dst = append(dst, t.pairs[0].V)
		dst = t.subs[1].appendValues(dst)
		if t.pairN > 1 {
			dst = append(dst, t.pairs[1].V)
			dst = t.subs[2].appendValues(dst)
			if t.pairN > 2 {
				dst = append(dst, t.pairs[2].V)
				dst = t.subs[3].appendValues(dst)
			}
		}
	}
	return dst
}

// Append appends each Key–Value pair in the Map to keys and values, ascending
// in Key order, and it returns the extended buffers.
func (m *Map[Key, Value]) Append(keys []Key, values []Value) ([]Key, []Value) {
	return m.top.appendKeysAndValues(keys, values)
}

func (t *node[Key, Value]) appendKeysAndValues(keys []Key, values []Value) ([]Key, []Value) {
	if t != nil {
		keys, values = t.subs[0].appendKeysAndValues(keys, values)
		keys = append(keys, t.pairs[0].K)
		values = append(values, t.pairs[0].V)
		keys, values = t.subs[1].appendKeysAndValues(keys, values)
		if t.pairN > 1 {
			keys = append(keys, t.pairs[1].K)
			values = append(values, t.pairs[1].V)
			keys, values = t.subs[2].appendKeysAndValues(keys, values)
			if t.pairN > 2 {
				keys = append(keys, t.pairs[2].K)
				values = append(values, t.pairs[2].V)
				keys, values = t.subs[3].appendKeysAndValues(keys, values)
			}
		}
	}
	return keys, values
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

// Append appends each Key in the Sap to dst, ascending in Key order, and it
// returns the extended buffer.
func (keys *Set[Key]) Append(dst []Key) []Key {
	return keys.m.AppendKeys(dst)
}

// Find returns the Key's presence in the Set.
func (keys *Set[Key]) Find(k Key) bool {
	return keys.m.FindPointer(k) != nil
}

// Insert adds the Key to the Set if and only if the Key is absent.
func (keys *Set[Key]) Insert(entry Key) bool {
	return keys.m.Insert(entry, struct{}{})
}

// At returns a new Cursor at located the Key, with false for none. A Delete or
// Insert renders the Cursor invalid.
func (keys *Set[Key]) At(k Key) (Cursor[Key, struct{}], bool) {
	return keys.m.At(k)
}

// NoCopy triggers go(1) vet when copied after the first use.
// See https://golang.org/issues/8005#issuecomment-190753527 for details.
type noCopy struct{}

// Lock triggers the -copylocks checker from go(1) vet.
func (*noCopy) Lock() {}

// Unlock triggers the -copylocks checker from go(1) vet.
func (*noCopy) Unlock() {}
