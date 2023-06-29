package pile

import (
	"fmt"
	"io"
	"strings"
)

// Height returns the number of levels in the B-tree.
// An empty Map has zero height.
func (m *Map[Key, Value]) height() int {
	var h int
	for t := m.top; t != nil; t = t.subs[0] {
		h++
	}
	return h
}

// dumpMap lists nodes per level (max 5) for debugging purposes.
func dumpMap[Key Sortable, Value any](m *Map[Key, Value]) string {
	height := m.height()
	if height == 0 {
		return ""
	}

	// scan line starts at top
	nodeRow := []*node[Key, Value]{m.top}
	var b strings.Builder

	minLevel := 0
	if height > 5 {
		minLevel = height - 5
	}
	for level := height - 1; ; level-- {
		fmt.Fprintf(&b, "Level %d:", level)
		for i := range nodeRow {
			b.WriteByte(' ')
			b.WriteString(nodeRow[i].String())
		}
		b.WriteByte('\n')

		if level == minLevel {
			return b.String()
		}

		var subs []*node[Key, Value]
		for _, t := range nodeRow {
			subs = append(subs, t.subs[:t.pairN+1]...)
		}
		nodeRow = subs
	}
}

// String returns a compact notation for debugging purposes.
func (t node[Key, Value]) String() string {
	var b strings.Builder
	b.WriteByte('[')
	printAsSub(&b, t.subs[0])
	for i := range t.pairs[:t.pairN] {
		fmt.Fprintf(&b, " %#v:%#v ", t.pairs[i].K, t.pairs[i].V)
		printAsSub(&b, t.subs[i+1])
	}
	b.WriteString(" ]")
	return b.String()
}

func printAsSub[Key Sortable, Value any](w io.Writer, t *node[Key, Value]) {
	if t == nil {
		return
	}
	switch t.pairN {
	case 1:
		fmt.Fprintf(w, " #%#v", t.pairs[0].K)
	case 2:
		fmt.Fprintf(w, " #%#v,%v", t.pairs[0].K, t.pairs[1].K)
	case 3:
		fmt.Fprintf(w, " #%#v,%v,%v", t.pairs[0].K, t.pairs[1].K, t.pairs[2].K)
	}
}
