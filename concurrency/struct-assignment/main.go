package main

import (
	"fmt"
	"sync"
)

const InputElements = 10000

type output struct {
	string
}

const MinGoroutines = 5

func main() {
	input := GenerateInput()
	fmt.Printf("Input data: %v\n", input)

	outputsSeq := GenerateSequentially(input)
	fmt.Printf("Output sequentially: %v\n", outputsSeq)

	outputsMap := GenerateConcurrentlyWithMap(input)
	fmt.Printf("Output concurrently	with map: %v\n", outputsMap)

	outputsChannel := GenerateConcurrentlyWithChannel(input)
	fmt.Printf("Output concurrently with channel: %v\n", outputsChannel)
}

func GenerateInput() []string {
	var input = []string{
		"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n",
		"ñ", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
	}
	var inputs []string

	for i := 0; i < InputElements; i++ {
		inputs = append(inputs, input...)
	}

	return inputs
}

func GenerateSequentially(input []string) []output {
	outputs := make([]output, len(input))
	for i, v := range input {
		outputs[i] = output{v}
	}

	return outputs
}

func GenerateConcurrentlyWithMap(input []string) []output {
	var wg sync.WaitGroup
	var mu sync.Mutex // Mutex for synchronizing access to outputs

	size := len(input)
	chunk := size / MinGoroutines
	iterations := MinGoroutines
	if size%chunk != 0 {
		iterations++
	}

	outputs := make([]output, len(input))

	for i := 0; i < iterations; i++ {
		startOffset := i * chunk
		endOffset := chunk * (i + 1)
		if chunk*(i+1) > size {
			endOffset = size
		}

		wg.Add(1)
		go func(part []string, iteration int) {
			defer wg.Done()

			for i, v := range part {
				index := i + (iteration * chunk)
				output := output{v}

				mu.Lock()
				outputs[index] = output
				mu.Unlock()
			}
		}(input[startOffset:endOffset], i)
	}

	wg.Wait()

	return outputs
}

func GenerateConcurrentlyWithChannel(input []string) []output {
	var wg sync.WaitGroup

	size := len(input)
	chunk := size / MinGoroutines
	iterations := MinGoroutines
	if size%chunk != 0 {
		iterations++
	}

	outputsCh := make(chan output, len(input))

	for i := 0; i < iterations; i++ {
		startOffset := i * chunk
		endOffset := chunk * (i + 1)
		if chunk*(i+1) > size {
			endOffset = size
		}

		wg.Add(1)
		go func(part []string, iteration int) {
			defer wg.Done()

			for _, v := range part {
				outputsCh <- output{v}
			}
		}(input[startOffset:endOffset], i)
	}

	go func() {
		wg.Wait()
		close(outputsCh)
	}()

	outputs := make([]output, 0, len(input))
	for o := range outputsCh {
		outputs = append(outputs, o)
	}

	return outputs
}
