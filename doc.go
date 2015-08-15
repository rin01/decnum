/*
package decnum is a Go binding around C decNumber package, for calculation with base-10 floating point numbers.
Decimal data type is important for financial calculations.

Godoc: https://godoc.org/github.com/rin01/decnum

Example of use:

	var (
		ctx decnum.Context
		a   decnum.DecQuad
		b   decnum.DecQuad
		r   decnum.DecQuad
	)

	ctx.Init(decnum.DEFAULT_DECQUAD) // initialize context with default settings for DecQuad operations. Essentially, it contains the rounding mode.

	if a, err = ctx.FromString("1234.5678"); err != nil { // convert string to DecQuad
		log.Fatal(err)
	}
	if b, err = ctx.FromString("-45.7"); err != nil { // convert string to DecQuad
		log.Fatal(err)
	}

	r = ctx.Add(a, b) // r = a + b
	// ...
	// you can put other operations here
	// ...

	fmt.Println("r", r.String())

	if err := ctx.Error(); err != nil { // you can just check for after a series of operations have been done
		log.Fatalf("ERROR OCCURRED !!!!!!!   %v\n", err)
	}

*/
package decnum



