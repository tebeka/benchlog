package main

import (
	"testing"
)

func memsetLoop(a []byte, v byte) {
	for i := range a {
		a[i] = v
	}
}

func memsetRepeat(a []byte, v byte) {
	if len(a) == 0 {
		return
	}
	a[0] = v
	for bp := 1; bp < len(a); bp *= 2 {
		copy(a[bp:], a[:bp])
	}
}

func BenchmarkLoop(b *testing.B) {
	buf := make([]byte, 1027)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		memsetLoop(buf, 0)
	}
}

func BenchmarkRepeat(b *testing.B) {
	buf := make([]byte, 1027)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		memsetRepeat(buf, 0)
	}
}
