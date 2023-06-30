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
	m.Put('一', "いち")
	m.Put('二', "に")
	m.Put('三', "さん")

	c, more := m.Last()
	for ; more; more = c.Backward() {
		c.Update(string([]rune{c.Fetch().K, 'つ'}))
		fmt.Println(c.Fetch().V)
	}
	// Output:
	// 二つ
	// 三つ
	// 一つ
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
			t.Errorf("cursor iteration allocated %f times, want 0", allocN)
		}
	})

	t.Run("Backward", func(t *testing.T) {
		allocN := testing.AllocsPerRun(1, func() {
			verifyBackward(t, &keys, want)
		})
		if !t.Failed() && allocN != 0 {
			t.Errorf("cursor iteration allocated %f times, want 0", allocN)
		}
	})
}

// VerifyForward uses an ascending Cursor to validate keys.
func verifyForward(t *testing.T, got *pile.Set[int], want []int) {
	c, more := got.First()
	for i := range want {
		if !more {
			t.Errorf("cursor halted before key № %d, want %d more keys", i+1, len(want)-i)
			return
		}

		if k := c.Fetch().K; k != want[i] {
			t.Errorf("cursor got key № %d value %d, want %d", i+1, k, want[i])
		}

		more = c.Forward()
	}
	if more {
		t.Errorf("cursor got more after all %d keys passed", len(want))
	}
}

// VerifyBackward uses an ascending Cursor to validate keys.
func verifyBackward(t *testing.T, got *pile.Set[int], want []int) {
	c, more := got.Last()
	for i := len(want) - 1; i >= 0; i-- {
		if !more {
			t.Errorf("cursor halted before key № %d", len(want)-i)
			return
		}

		if k := c.Fetch().K; k != want[i] {
			t.Errorf("cursor got key № %d value %d, want %d", i+1, k, want[i])
		}

		more = c.Backward()
	}
	if more {
		t.Errorf("cursor got more after all %d keys passed", len(want))
	}
}

func TestRange(t *testing.T) {
	var keys pile.Set[int]
	if c, ok := keys.First(); ok {
		t.Errorf("got first key %d on empty Set", c.Fetch().K)
	}
	if c, ok := keys.Last(); ok {
		t.Errorf("got last key %d on empty Set", c.Fetch().K)
	}

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

func verifyRange[Key pile.Sortable](t *testing.T, got *pile.Set[Key], first, last Key) {
	c, ok := got.First()
	if !ok {
		t.Errorf("first unvalaible, want key %v", first)
	} else if got := c.Fetch().K; got != first {
		t.Errorf("got first key %v, want %v", got, first)
	}

	c, ok = got.Last()
	if !ok {
		t.Errorf("last unvalaible, want key %v", last)
	} else if got := c.Fetch().K; got != last {
		t.Errorf("got last key %v, want %v", got, last)
	}
}
