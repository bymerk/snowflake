package showflake

import (
	"testing"
)

func BenchmarkSnowflake_Generate(b *testing.B) {
	s, _ := NewSnowflake(10)
	for b.Loop() {
		_ = s.Generate()
	}
}
