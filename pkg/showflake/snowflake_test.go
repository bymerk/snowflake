package showflake

import (
	"testing"
)

func BenchmarkSnowflake_Generate(b *testing.B) {
	sf, err := NewSnowflake(10, 10)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = sf.Generate()
		}
	})
}
