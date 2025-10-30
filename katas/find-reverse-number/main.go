package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func main() {
	var start time.Time
	// start = time.Now()
	// fmt.Println(FindReverseNumberV1(100000)) // 2min46s
	// fmt.Println("V1 took ", time.Since(start))

	start = time.Now()
	fmt.Println(FindReverseNumberV2(100000)) // 8,9s
	fmt.Printf("V2 took %s\n", time.Since(start))

	// start = time.Now()
	// fmt.Println(FindReverseNumberV3(100000)) // 10.25s -> slower than V2
	// fmt.Printf("V3 took %s\n", time.Since(start))

	// Final version should do 100000000 in less than 12s
}

func FindReverseNumberV1(n uint64) uint64 {
	var it, m uint64
	var l []string
	var s string
	for {
		s = strconv.FormatUint(it, 10)
		l = strings.Split(s, "")
		for i, j := 0, len(l)-1; i < len(l); i, j = i+1, j-1 {
			if l[i] != l[j] {
				break
			}
			if (len(l)%2 == 1 && i == j && l[i] == l[j]) || (len(l)%2 == 0 && j-i == 1 && l[i] == l[j]) {
				m++
			}
		}
		if m == n {
			return it
		}
		it++
	}
}

func FindReverseNumberV2(n uint64) uint64 {
	var it, m uint64
	var isPalindrome = func(a uint64) bool {
		var z uint64
		tmp := a
		for tmp > 0 {
			z = z*10 + tmp%10
			tmp = tmp / 10
		}
		return z == a
	}

	for {
		if isPalindrome(it) {
			m++
			if m == n {
				return it
			}
		}
		it++
	}
}

func FindReverseNumberV3(n uint64) uint64 {
	var it, m uint64
	var isPalindrome = func(a uint64) bool {
		var z uint64
		tmp := a
		if tmp%10 == 0 && tmp != 0 {
			return false
		}
		for tmp > 0 {
			z = z*10 + tmp%10
			tmp = tmp / 10
		}
		return z == a
	}

	for {
		if isPalindrome(it) {
			m++
			if m == n {
				return it
			}
		}
		it++
	}
}
