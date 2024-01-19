package main

import (
	"fmt"
	"math"
)

func main() {
	// Code
	for i := 2; i < 100; i++ {
		var prime bool = true

		for j := 2; j < i; j++ {
			if i%j == 0 {
				prime = false
			}
		}

		if prime {
			fmt.Println(i)
		}
	}

	fmt.Println(2)
	for i := 3; i < 100; i += 2 {
		var prime bool = true

		for j := 3; j < i; j += 2 {
			if i%j == 0 {
				prime = false
			}
		}

		if prime {
			fmt.Println(i)
		}
	}

	fmt.Println(2)
	for i := 3; i < 100; i += 2 {
		var prime bool = true

		for j := 3; j < int(math.Sqrt(float64(i)))+1; j += 2 {
			if i%j == 0 {
				prime = false
			}
		}

		if prime {
			fmt.Println(i)
		}
	}
}
