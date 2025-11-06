package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

const n, max uint64 = 11, 100_000_000_000
const timeout = 12 * time.Second

func main() {
	// FindReverseNumberV1(n=100000) // 2min46s
	// FindReverseNumberV2(n=100000) // 8,9s
	// FindReverseNumberV3(n=100000) // 10.25s -> slower than V2

	// r1, err1 := execWithTimeout(n, FindReverseNumberV1)
	// r2, err2 := execWithTimeout(n, FindReverseNumberV2)
	// r3, err3 := execWithTimeout(n, FindReverseNumberV3)
	r4, err4 := execWithTimeout(n, FindReverseNumberV4)

	// fmt.Printf("V1 finding n=%d: <%d, %v>\n", n, r1, err1)
	// fmt.Printf("V2 finding n=%d: <%d, %v>\n", n, r2, err2)
	// fmt.Printf("V3 finding n=%d: <%d, %v>\n", n, r3, err3)
	fmt.Printf("V4 finding n=%d: <%d, %v>\n", n, r4, err4)

	// Final version should do 100000000 in less than 12s
}

func execWithTimeout(n uint64, fn func(n uint64) uint64) (uint64, error) {
	type result struct {
		v   uint64
		err error
	}
	res := make(chan result, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				res <- result{0, fmt.Errorf("panic: %v", r)}
			}
		}()
		res <- result{v: fn(n)}
	}()

	select {
	case r := <-res:
		return r.v, r.err
	case <-time.After(timeout):
		return 0, fmt.Errorf("timeout after %fs", timeout.Seconds())
	}
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

func FindReverseNumberV4(n uint64) uint64 {
	if n <= 10 {
		return n - 1
	}

	var cumN, length uint64 = 10, 2
	for cumN < n {
		l := int(math.Ceil(float64(length)/2)) - 1
		nL := 9 * uint64(math.Pow10(l))
		if cumN+nL > n {
			break
		}
		cumN += nL
		length++
	}
	offset := n - cumN - 1
	prefixL := int(math.Ceil(float64(length) / 2))
	base := uint64(math.Pow10(prefixL - 1))
	firstDigit := 1 + (offset / base)
	restPadded := fmt.Sprintf("%0*d", prefixL-1, offset%base)
	prefix := fmt.Sprintf("%d%s", firstDigit, restPadded)
	fmt.Println(cumN, prefixL, base, firstDigit, restPadded, prefix)

	reverse := func(s string) string {
		r := []rune(s)
		for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
			r[i], r[j] = r[j], r[i]
		}
		return string(r)
	}

	var palindrome string
	if length%2 == 0 {
		palindrome = prefix + reverse(prefix[:prefixL])
	} else {
		palindrome = prefix + reverse(prefix[:prefixL-1])
	}

	res, err := strconv.ParseUint(palindrome, 10, 64)
	if err != nil {
		panic(err)
	}
	return res
}
