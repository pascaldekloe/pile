package pile

import "testing"

func TestSet_script(t *testing.T) {
	var c Set[string]
	if n := c.Size(); n != 0 {
		t.Errorf("got size %d on empty Set, want 0", n)
	}
	if c.Find("de") {
		t.Fatal("found string in empty Set")
	}
	if !c.Insert("de") {
		t.Error("new insert got false from empty Set, want true")
	}
	if n := c.Size(); n != 1 {
		t.Errorf("got size %d after insert, want 1", n)
	}
	if !c.Find("de") {
		t.Error("single insert not found")
	}
	if c.Insert("de") {
		t.Error("insert again got true, want false")
	}
	if n := c.Size(); n != 1 {
		t.Errorf("got size %d after duplicate insert, want 1", n)
	}
	if t.Failed() {
		return
	}

	if c.Find("de_CH") {
		t.Error("found absent string with matching prefix")
	}
	if !c.Insert("de_CH") {
		t.Error("second insert got false on matching prefix, want true")
	}
	if !c.Find("de_CH") {
		t.Error("second insert not found")
	}
	if !c.Find("de") {
		t.Error("first insert not lost after second insert")
	}
	if n := c.Size(); n != 2 {
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
