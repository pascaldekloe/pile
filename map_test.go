package pile

import (
	"math/rand"
	"slices"
	"strconv"
	"testing"
)

func TestMap(t *testing.T) {
	// key-only definition of upsert result
	type golden struct {
		feed []int
		want []int
	}

	// single node, insert-only
	var allTripples = []golden{
		{feed: []int{1, 2, 3}, want: []int{1, 2, 3}},
		{feed: []int{1, 3, 2}, want: []int{1, 2, 3}},
		{feed: []int{2, 1, 3}, want: []int{1, 2, 3}},
		{feed: []int{2, 3, 1}, want: []int{1, 2, 3}},
		{feed: []int{3, 2, 1}, want: []int{1, 2, 3}},
		{feed: []int{3, 1, 2}, want: []int{1, 2, 3}},
	}

	var tests []golden
	tests = append(tests, allTripples...)

	for _, tri := range allTripples {
		var l golden
		l.feed = append(l.feed, tri.feed...)
		l.feed = append(l.feed, tri.feed...)
		l.feed = append(l.feed, tri.feed...)
		l.want = append(l.want, tri.want...)
		l.want = append(l.want, tri.want...)
		l.want = append(l.want, tri.want...)
		for i := 3; i < 6; i++ {
			l.feed[i] += 3
			l.want[i] += 3
		}
		for i := 6; i < 9; i++ {
			l.feed[i] += 6
			l.want[i] += 6
		}
		tests = append(tests, l)
	}
	for _, tri := range allTripples {
		var l golden
		l.feed = append(l.feed, tri.feed...)
		l.feed = append(l.feed, tri.feed...)
		l.feed = append(l.feed, tri.feed...)
		l.want = append(l.want, tri.want...)
		l.want = append(l.want, tri.want...)
		l.want = append(l.want, tri.want...)
		for i := 3; i < 6; i++ {
			l.feed[i+3] += 3
			l.want[i] += 3
		}
		for i := 6; i < 9; i++ {
			l.feed[i-6] += 6
			l.want[i] += 6
		}
		tests = append(tests, l)
	}

	var buf []Pair[int, int]
	for _, test := range tests {
		// build expected with some values
		want := make([]Pair[int, int], len(test.want))
		values := make(map[int]int, len(test.want))
		for i, key := range test.want {
			want[i] = Pair[int, int]{K: key, V: key + 42}
			values[key] = want[i].V
		}

		var list []byte
		for i, v := range test.feed {
			if i != 0 {
				list = append(list, ',')
			}
			list = strconv.AppendInt(list, int64(v), 10)
		}

		t.Run("Put"+string(list), func(t *testing.T) {
			var m Map[int, int]
			for _, key := range test.feed {
				value := values[key]
				m.Put(key, value)
			}
			verifyMapEqual(t, "golden", &m, values)
			buf = m.AppendPairs(buf[:0])
			if !slices.Equal(buf, want) {
				t.Errorf("result got pairs %d, want pairs %d", buf, want)
				t.Log(dumpMap(&m))
			}
		})

		t.Run("InsertOrUpdate"+string(list), func(t *testing.T) {
			var m Map[int, int]
			for _, key := range test.feed {
				value := values[key]
				if !m.Insert(key, value) {
					if !m.Update(key, value) {
						t.Errorf("both Insert and Update got false for %d, want either one true", key)
					}
				}
			}
			verifyMapEqual(t, "golden", &m, values)
			buf = m.AppendPairs(buf[:0])
			if !slices.Equal(buf, want) {
				t.Errorf("result got pairs %d, want pairs %d", buf, want)
				t.Log(dumpMap(&m))
			}
		})

		t.Run("UpdateOrInsert"+string(list), func(t *testing.T) {
			var m Map[int, int]
			for _, key := range test.feed {
				value := values[key]
				if !m.Update(key, value) {
					if !m.Insert(key, value) {
						t.Errorf("both Update and Insert got false for %d, want either one true", key)
					}
				}
			}
			verifyMapEqual(t, "golden", &m, values)
			buf = m.AppendPairs(buf[:0])
			if !slices.Equal(buf, want) {
				t.Errorf("result got pairs %d, want pairs %d", buf, want)
				t.Log(dumpMap(&m))
			}
		})

	}
}

func TestMapAppend(t *testing.T) {
	const entryN = 40 // causes 3 levels
	reference := make(map[int]int, entryN)
	var inserts, puts Map[int, int]

	for i := 0; i < entryN && !t.Failed(); i++ {
		key, value := i, i+100
		reference[key] = value

		if !inserts.Insert(key, value) {
			t.Errorf("Insert %d got false", key)
		}
		verifyMapEqual(t, "Insert", &inserts, reference)
		puts.Put(key, value)
		verifyMapEqual(t, "Put", &puts, reference)
	}

	if testing.Verbose() {
		t.Log("Inserts got:\n", dumpMap(&inserts))
		t.Log("Puts got:\n", dumpMap(&puts))
	}
}

func TestMapPrepend(t *testing.T) {
	const entryN = 40 // causes 3 levels
	reference := make(map[int]int, entryN)
	var inserts, puts Map[int, int]

	for i := entryN - 1; i >= 0 && !t.Failed(); i-- {
		key, value := i, i+100
		reference[key] = value

		if !inserts.Insert(key, value) {
			t.Errorf("Insert %d got false", key)
		}
		verifyMapEqual(t, "Insert", &inserts, reference)
		puts.Put(key, value)
		verifyMapEqual(t, "Put", &puts, reference)
	}

	if testing.Verbose() || t.Failed() {
		t.Log("Inserts got:\n", dumpMap(&inserts))
		t.Log("Puts got:\n", dumpMap(&puts))
	}
}

func TestMapRandom(t *testing.T) {
	r, ok := rand.NewSource(1337).(rand.Source64)
	if !ok {
		t.Fatal("non 64-bit random source")
	}

	const entryN = 1000 // causes some dupes with 16-bit space
	reference := make(map[uint16]int, entryN)
	var insertOrUpdates, puts Map[uint16, int]

	for i := 0; i < entryN && !t.Failed(); i++ {
		bits := r.Uint64()
		k1, k2, k3, k4 := uint16(bits), uint16(bits>>16), uint16(bits>>32), uint16(bits>>48)
		v1, v2, v3, v4 := i<<2+0, i<<2+1, i<<2+2, i<<2+3
		reference[k1] = v1
		reference[k2] = v2
		reference[k3] = v3
		reference[k4] = v4

		if !insertOrUpdates.Insert(k1, v1) {
			if !insertOrUpdates.Update(k1, v1) {
				t.Errorf("both Insert and Update got false for %d, want either one true", k1)
			}
		}
		if !insertOrUpdates.Insert(k2, v2) {
			if !insertOrUpdates.Update(k2, v2) {
				t.Errorf("both Insert and Update got false for %d, want either one true", k2)
			}
		}
		if !insertOrUpdates.Insert(k3, v3) {
			if !insertOrUpdates.Update(k3, v3) {
				t.Errorf("both Insert and Update got false for %d, want either one true", k3)
			}
		}
		if !insertOrUpdates.Insert(k4, v4) {
			if !insertOrUpdates.Update(k4, v4) {
				t.Errorf("both Insert and Update got false for %d, want either one true", k4)
			}
		}
		verifyMapEqual(t, "Insert or Update", &insertOrUpdates, reference)
		puts.Put(k1, v1)
		puts.Put(k2, v2)
		puts.Put(k3, v3)
		puts.Put(k4, v4)
		verifyMapEqual(t, "Put", &puts, reference)
	}

	if testing.Verbose() || t.Failed() {
		t.Log("Inserts or Updates got:\n", dumpMap(&insertOrUpdates))
		t.Log("Puts got:\n", dumpMap(&puts))
	}
}

func verifyMapEqual[Key Sortable, Value comparable](t *testing.T, name string, got *Map[Key, Value], want map[Key]Value) {
	for k, v := range want {
		switch actual, found := got.Find(k); {
		case !found:
			t.Errorf("%s result map has key %#v absent, want value %#v", name, k, v)
		case actual != v:
			t.Errorf("%s result map got key %#v value %#v, want value %#v", name, k, actual, v)
		}
	}
	if n := got.Size(); n != len(want) {
		t.Errorf("%s result map got Size %d, want %d", name, n, len(want))
	}
}
