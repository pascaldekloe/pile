package pile_test

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"

	"github.com/pascaldekloe/pile"
)

func ExampleCursor() {
	var m pile.Map[rune, string]
	m.Put('Ⅰ', "一")
	m.Put('Ⅱ', "二")
	m.Put('Ⅲ', "三")

	for c, ok := m.Least(); ok; ok = c.Ascend() {
		fmt.Print(c.Swap(c.Value() + "つ"))
	}
	// Output: 一二三
}

func TestIteration(t *testing.T) {
	var keys pile.Set[int]
	var want []int

	r := rand.NewSource(42)
	for i := 0; i < 99; i++ {
		k := int(r.Int63())
		keys.Insert(k)
		want = append(want, k)
	}
	sort.Ints(want)

	t.Run("Forward", func(t *testing.T) {
		allocN := testing.AllocsPerRun(1, func() {
			verifyForward(t, &keys, want)
		})
		if !t.Failed() && allocN != 0 {
			t.Errorf("cursor allocated %f times, want 0", allocN)
		}
	})

	t.Run("Backward", func(t *testing.T) {
		allocN := testing.AllocsPerRun(1, func() {
			verifyBackward(t, &keys, want)
		})
		if !t.Failed() && allocN != 0 {
			t.Errorf("cursor allocated %f times, want 0", allocN)
		}
	})

	t.Run("JumpIn", func(t *testing.T) {
		allocN := testing.AllocsPerRun(1, func() {
			for i := 1; i < len(want); i++ {
				j := i
				for c, ok := keys.At(want[i]); ok; ok = c.Ascend() {
					if j >= len(want) {
						t.Fatalf("cursor since key № %d exceeds %d entries", i+1, len(want))
					}
					if got := c.Key(); got != want[j] {
						t.Errorf("cursor since key № %d mismatched key № %d; got %d, want %d",
							i+1, j+1, got, want[j])
					}
					j++
				}
				if j != len(want) {
					t.Fatalf("cursor since key № %d missed %d entries", i+1, len(want)-j)
				}
			}
		})
		if !t.Failed() && allocN != 0 {
			t.Errorf("cursor allocated %f times, want 0", allocN)
		}
	})
}

// VerifyForward iterates ascending to validate keys.
func verifyForward(t *testing.T, got *pile.Set[int], want []int) {
	c, ok := got.Least()
	for i := range want {
		if !ok {
			t.Errorf("cursor halted before key № %d, want %d more keys", i+1, len(want)-i)
			return
		}

		if k := c.Key(); k != want[i] {
			t.Errorf("cursor got key № %d value %d, want %d", i+1, k, want[i])
		}

		ok = c.Ascend()
	}
	if ok {
		t.Errorf("cursor got more after all %d keys passed", len(want))
	}
}

// VerifyBackward iterates descending to validate keys.
func verifyBackward(t *testing.T, got *pile.Set[int], want []int) {
	c, ok := got.Most()
	for i := len(want) - 1; i >= 0; i-- {
		if !ok {
			t.Errorf("cursor halted before key № %d", len(want)-i)
			return
		}

		if k := c.Key(); k != want[i] {
			t.Errorf("cursor got key № %d value %d, want %d", i+1, k, want[i])
		}

		ok = c.Descend()
	}
	if ok {
		t.Errorf("cursor got more after all %d keys passed", len(want))
	}
}

func TestRange(t *testing.T) {
	var keys pile.Set[int]
	r := rand.NewSource(42)
	var want []int
	for i := 0; i < 99; i++ {
		k := int(r.Int63())
		keys.Insert(k)

		want = append(want, k)
		sort.Ints(want)

		verifyRange(t, &keys, want[0], want[len(want)-1])
	}
}

func verifyRange[Key pile.Sortable](t *testing.T, got *pile.Set[Key], least, most Key) {
	c, ok := got.Least()
	if !ok {
		t.Errorf("least unvalaible, want %v", least)
	} else if got := c.Key(); got != least {
		t.Errorf("got least %v, want %v", got, least)
	}

	c, ok = got.Most()
	if !ok {
		t.Errorf("most unvalaible, want %v", most)
	} else if got := c.Key(); got != most {
		t.Errorf("got most %v, want %v", got, most)
	}
}

func TestNoCursor(t *testing.T) {
	var keys pile.Map[string, string]
	if c, ok := keys.Least(); ok {
		t.Errorf("got least key %q on empty Set", c.Key())
	} else {
		verifyZeroCursor(t, &c)
	}
	if c, ok := keys.Most(); ok {
		t.Errorf("got most key %q on empty Set", c.Key())
	} else {
		verifyZeroCursor(t, &c)
	}
	if c, ok := keys.At("x"); ok {
		t.Errorf(`got at "x" key %q on empty Set`, c.Key())
	} else {
		verifyZeroCursor(t, &c)
	}
}

func verifyZeroCursor(t *testing.T, c *pile.Cursor[string, string]) {
	if got := c.Key(); got != "" {
		t.Errorf("got key %q from zero iterator, want none", got)
	}
	if got := c.Value(); got != "" {
		t.Errorf("got value %q from zero iterator, want none", got)
	}
	if c.Ascend() {
		t.Error("got ascend from zero iterator")
	}
	if c.Descend() {
		t.Error("got descend from zero iterator")
	}
	if got := c.Swap("foo"); got != "" {
		t.Errorf("got %q from swap on zero iterator, want none", got)
	}
	if got := c.Swap("bar"); got != "" {
		t.Errorf("got %q from second swap on zero iterator, want none", got)
	}
}
