/*
  Example:

	package main

	import (
		"fmt"
		"log"
		"os"

		"github.com/rin01/decnum"
	)

	func main() {
		var (
			a decnum.Quad
			b decnum.Quad
			r decnum.Quad
		)

		a, _ = decnum.FromString(os.Args[1])
		b, _ = decnum.FromString(os.Args[2])

		r = a.Add(b)

		if err := r.Error(); err != nil {
			log.Fatalf("Error: %s", err)
		}

		fmt.Printf("result:    %s\n", r)

	}


  ===== WARNING: THIS PACKAGE WORKS, BUT API MAY CHANGE A LITTLE. Also, I must add some more tests. =====

*/
package decnum
