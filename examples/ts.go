package main

import "fmt"

func main() {
	ss := []string{"site", "user", "info", "0"}

	fmt.Printf("%v -> ", ss)

	sliceReverse(ss)

	fmt.Printf("%v\n", ss)
}

func sliceReverse(ss []string) {
	ln := len(ss)

	for i := 0; i < int(ln/2); i++ {
		li := ln - i - 1
		// fmt.Println(i, "<=>", li)
		ss[i], ss[li] = ss[li], ss[i]
	}
}
