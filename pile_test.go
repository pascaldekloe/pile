package pile

import "testing"

func TestSet_script(t *testing.T) {
	var keys Set[string]
	if n := keys.Size(); n != 0 {
		t.Errorf("got size %d on empty Set, want 0", n)
	}
	if keys.Find("de") {
		t.Fatal("found string in empty Set")
	}
	if !keys.Insert("de") {
		t.Error("new insert got false from empty Set, want true")
	}
	if n := keys.Size(); n != 1 {
		t.Errorf("got size %d after insert, want 1", n)
	}
	if !keys.Find("de") {
		t.Error("single insert not found")
	}
	if keys.Insert("de") {
		t.Error("insert again got true, want false")
	}
	if n := keys.Size(); n != 1 {
		t.Errorf("got size %d after duplicate insert, want 1", n)
	}
	if t.Failed() {
		return
	}

	if keys.Find("de_CH") {
		t.Error("found absent string with matching prefix")
	}
	if !keys.Insert("de_CH") {
		t.Error("second insert got false on matching prefix, want true")
	}
	if !keys.Find("de_CH") {
		t.Error("second insert not found")
	}
	if !keys.Find("de") {
		t.Error("first insert not lost after second insert")
	}
	if n := keys.Size(); n != 2 {
		t.Errorf("got size %d after second insert, want 2", n)
	}
}

func TestMap_script(t *testing.T) {
	var m Map[string, [3]byte]
	if n := m.Size(); n != 0 {
		t.Errorf("got size %d on empty Map, want 0", n)
	}
	if _, ok := m.Find("de"); ok {
		t.Fatal("found string in empty Map")
	}
	if !m.Insert("de", [3]byte{'D', 'E', 'U'}) {
		t.Error("new insert got false from empty Map, want true")
	}
	if n := m.Size(); n != 1 {
		t.Errorf("got size %d after insert, want 1", n)
	}
	if _, ok := m.Find("de"); !ok {
		t.Error("single insert not found")
	}
	if m.Insert("de", [3]byte{'G', 'E', 'R'}) {
		t.Error("insert again got true, want false")
	}
	if n := m.Size(); n != 1 {
		t.Errorf("got size %d after duplicate insert, want 1", n)
	}
	if t.Failed() {
		return
	}

	if _, ok := m.Find("de_CH"); ok {
		t.Error("found absent string with matching prefix")
	}
	if !m.Insert("de_CH", [3]byte{'C', 'H', 'E'}) {
		t.Error("second insert got false on matching prefix, want true")
	}
	if v, ok := m.Find("de_CH"); !ok {
		t.Error("second insert not found")
	} else if string(v[:]) != "CHE" {
		t.Errorf("got second insert value %#x, want CHE", v)
	}
	if v, ok := m.Find("de"); !ok {
		t.Error("first insert not lost after second insert")
	} else if string(v[:]) != "DEU" {
		t.Errorf("got second insert value %#x, want DEU", v)
	}
	if n := m.Size(); n != 2 {
		t.Errorf("got size %d after second insert, want 2", n)
	}
}

func TestAppendKeys(t *testing.T) {
	const n = 99
	var keys Set[int]
	var want []int
	for i := 0; i < n; i++ {
		keys.Insert(i)
		want = append(want, i)
	}

	got := keys.AppendKeys(nil)
	if len(got) != n {
		t.Fatalf("got %d keys, want %d", len(got), n)
	}
	for i := range got {
		if got[i] != want[i] {
			t.Errorf("got keys %d, want %d", got[i], want[i])
		}
	}
}

func TestAppendKeysAndValues(t *testing.T) {
	const n = 99
	var m Map[int, int]
	var wantKeys, wantValues []int
	for i := 0; i < n; i++ {
		m.Put(i, i+1001)
		wantKeys = append(wantKeys, i)
		wantValues = append(wantValues, i+1001)
	}

	gotKeys, gotValues := m.AppendPairs(nil, nil)
	if len(gotKeys) != n || len(gotValues) != n {
		t.Fatalf("got %d keys and %d values, want %d for both", len(gotKeys), len(gotValues), n)
	}
	for i := range gotKeys {
		if gotKeys[i] != wantKeys[i] || gotValues[i] != wantValues[i] {
			t.Errorf("got key—value pair %d–%d, want %d–%d",
				gotKeys[i], gotValues[i], wantValues[i], wantValues[i])
		}
	}

	gotValues = m.AppendValues(gotValues[:0])
	if len(gotValues) != n {
		t.Fatalf("got %d values, want %d", len(gotValues), n)
	}
	for i := range gotValues {
		if gotValues[i] != wantValues[i] {
			t.Errorf("got value %d, want %d", wantValues[i], wantValues[i])
		}
	}
}
