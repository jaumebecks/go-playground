package main

import (
	"testing"
)

func BenchmarkGenerateFeedSequentially(b *testing.B) {
	b.ResetTimer()
	input := GenerateInput()

	for i := 0; i < b.N; i++ {
		GenerateSequentially(input)
	}
}

func BenchmarkGenerateConcurrentlyWithMap(b *testing.B) {
	b.ResetTimer()
	input := GenerateInput()

	for i := 0; i < b.N; i++ {
		GenerateConcurrentlyWithMap(input)
	}
}

func BenchmarkGenerateConcurrentlyWithChannel(b *testing.B) {
	b.ResetTimer()

	input := GenerateInput()
	for i := 0; i < b.N; i++ {
		GenerateConcurrentlyWithChannel(input)
	}
}
