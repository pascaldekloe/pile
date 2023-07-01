package pile

import (
	"math/rand"
	"testing"
)

func BenchmarkMapFind(b *testing.B) {
	b.Run("Sequential", func(b *testing.B) {
		var seq1Ki Map[int, string]
		for i := 0; i < 1024; i++ {
			seq1Ki.Put(i, "fill")
		}
		var seq1Mi Map[int, string]
		for i := 0; i < 1024*1024; i++ {
			seq1Mi.Put(i, "fill")
		}
		var seq64Mi Map[int, string]
		for i := 0; i < 64*1024*1024; i++ {
			seq64Mi.Put(i, "fill")
		}

		b.Run("1Ki", func(b *testing.B) {
			const mask = 1024 - 1
			for i := 0; i < b.N; i++ {
				if _, ok := seq1Ki.Find(i & mask); !ok {
					b.Fatalf("key %d not found", i&mask)
				}
			}
		})
		b.Run("1Mi", func(b *testing.B) {
			const mask = 1024*1024 - 1
			for i := 0; i < b.N; i++ {
				if _, ok := seq1Mi.Find(i & mask); !ok {
					b.Fatalf("key %d not found", i&mask)
				}
			}
		})
		b.Run("64Mi", func(b *testing.B) {
			const mask = 64*1024*1024 - 1
			for i := 0; i < b.N; i++ {
				if _, ok := seq64Mi.Find(i & mask); !ok {
					b.Fatalf("key %d not found", i&mask)
				}
			}
		})
	})

	b.Run("Random", func(b *testing.B) {
		feed := nRandomInts(64 * 1024 * 1024)
		var rnd1Ki Map[int, string]
		for _, k := range feed[:1024] {
			rnd1Ki.Put(k, "fill")
		}
		var rnd1Mi Map[int, string]
		for _, k := range feed[:1024*1024] {
			rnd1Mi.Put(k, "fill")
		}
		var rnd64Mi Map[int, string]
		for _, k := range feed[:64*1024*1024] {
			rnd64Mi.Put(k, "fill")
		}

		b.Run("1Ki", func(b *testing.B) {
			const mask = 1024 - 1
			for i := 0; i < b.N; i++ {
				if _, ok := rnd1Ki.Find(feed[i&mask]); !ok {
					b.Fatalf("key %d not found", i&mask)
				}
			}
		})
		b.Run("1Mi", func(b *testing.B) {
			const mask = 1024*1024 - 1
			for i := 0; i < b.N; i++ {
				if _, ok := rnd1Mi.Find(feed[i&mask]); !ok {
					b.Fatalf("key %d not found", i&mask)
				}
			}
		})
		b.Run("64Mi", func(b *testing.B) {
			const mask = 64*1024*1024 - 1
			for i := 0; i < b.N; i++ {
				if _, ok := rnd64Mi.Find(feed[i&mask]); !ok {
					b.Fatalf("key %d not found", i&mask)
				}
			}
		})
	})
}

func BenchmarkInsert(b *testing.B) {
	b.Run("Append", func(b *testing.B) {
		b.Run("map", func(b *testing.B) {
			var m Map[int, string]
			for i := 0; i < b.N; i++ {
				if !m.Insert(i, "foo") {
					b.Fatalf("insertion %d denied", i)
				}
			}
		})
		b.Run("set", func(b *testing.B) {
			var keys Set[int]
			for i := 0; i < b.N; i++ {
				if !keys.Insert(i) {
					b.Fatalf("insertion %d denied", i)
				}
			}
		})
	})

	b.Run("Prepend", func(b *testing.B) {
		b.Run("map", func(b *testing.B) {
			var m Map[int, string]
			for i := b.N; i > 0; i-- {
				if !m.Insert(i, "foo") {
					b.Fatalf("insertion %d denied", i)
				}
			}
		})
		b.Run("set", func(b *testing.B) {
			var keys Set[int]
			for i := b.N; i > 0; i-- {
				if !keys.Insert(i) {
					b.Fatalf("insertion %d denied", i)
				}
			}
		})
	})

	b.Run("Random", func(b *testing.B) {
		feed := nRandomInts(10e6)

		b.Run("map", func(b *testing.B) {
			if b.N > len(feed) {
				b.Fatalf("random set of %d keys not enough for bench N %d", len(feed), b.N)
			}
			var m Map[int, string]
			for _, k := range feed[:b.N] {
				if !m.Insert(k, "foo") {
					b.Fatalf("insertion %d denied", k)
				}
			}
		})

		b.Run("set", func(b *testing.B) {
			if b.N > len(feed) {
				b.Fatalf("random set of %d keys not enough for bench N %d", len(feed), b.N)
			}
			var keys Set[int]
			for _, k := range feed[:b.N] {
				if !keys.Insert(k) {
					b.Fatalf("insertion %d denied", k)
				}
			}
		})
	})
}

func nRandomInts(n int) []int {
	ints := make([]int, n)
	have := make(map[int]struct{}, n)

	r := rand.NewSource(42)
	for i := range ints {
		for {
			k := int(r.Int63())
			if _, ok := have[k]; ok {
				continue // dupe
			}
			ints[i] = k
			have[k] = struct{}{}
			break
		}
	}
	return ints
}
