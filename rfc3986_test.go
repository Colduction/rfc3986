package rfc3986_test

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/colduction/rfc3986"
)

func benchPerCoreConfigs(b *testing.B, f func(b *testing.B)) {
	b.Helper()
	coreConfigs := []int{1, 2, 4, 8, 12}
	for _, n := range coreConfigs {
		name := fmt.Sprintf("%d cores", n)
		b.Run(name, func(b *testing.B) {
			runtime.GOMAXPROCS(n)
			f(b)
		})
	}
}

func BenchmarkQueryEscape(b *testing.B) {
	benchPerCoreConfigs(b, func(b *testing.B) {
		b.RunParallel(func(b *testing.PB) {
			for b.Next() {
				rfc3986.QueryEscape(" Hello World! ")
			}
		})
	})
}
func BenchmarkQueryUnescape(b *testing.B) {
	benchPerCoreConfigs(b, func(b *testing.B) {
		b.RunParallel(func(b *testing.PB) {
			for b.Next() {
				rfc3986.QueryUnescape("%20Hello%20World!%20")
			}
		})
	})
}
