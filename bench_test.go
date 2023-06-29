package pile

import "testing"

func BenchmarkMapFind(b *testing.B) {
	b.Run("empty", func(b *testing.B) {
		benchmarkMapFindSeq(b, 0)
	})
	b.Run("1M-seq", func(b *testing.B) {
		benchmarkMapFindSeq(b, 1e6)
	})
	b.Run("5M-seq", func(b *testing.B) {
		benchmarkMapFindSeq(b, 5e6)
	})
}

func benchmarkMapFindSeq(b *testing.B, fillSize int) {
	var m Map[int, string]
	for i := 0; i < fillSize; i++ {
		m.Put(i, "fill")
	}
	if n := m.Size(); n != fillSize {
		b.Fatalf("fill size is %d, want %d", n, fillSize)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Find(i)
	}
}

func BenchmarkMapAppend(b *testing.B) {
	b.Run("empty", func(b *testing.B) {
		benchmarkMapAppend(b, 0)
	})
	b.Run("1M-seq", func(b *testing.B) {
		benchmarkMapAppend(b, 1e6)
	})
	b.Run("5M-seq", func(b *testing.B) {
		benchmarkMapAppend(b, 5e6)
	})
}

func benchmarkMapAppend(b *testing.B, fillSize int) {
	var m Map[int, string]
	for i := -fillSize; i < 0; i++ {
		m.Put(i, "fill")
	}
	if n := m.Size(); n != fillSize {
		b.Fatalf("fill size is %d, want %d", n, fillSize)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Put(i, "foo")
	}
}
