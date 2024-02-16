package main

func main() {

	for i := 2; i < 1_000_000; i++ {
		var prime bool = true

		for j := 2; j < i; j++ {
			if i%j == 0 {
				prime = false
			}
		}

		if prime {
			println(i)
		}
	}
}
