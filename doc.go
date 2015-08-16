/*
package decnum is a Go binding around C decNumber package, for calculation with base-10 floating point numbers.
Decimal data type is important for financial calculations.

Godoc: https://godoc.org/github.com/rin01/decnum

Example of use:

	var (
		ctx decnum.Context
		a   decnum.Quad  // unlike Go variable, the uninitialized value is not zero, but garbage.
		b   decnum.Quad  //    so, if you need a variable to be 0, do         x = decnum.Zero()
		r   decnum.Quad
	)

	ctx.InitDefaultQuad()           // initialize context with default settings for Quad operations. Essentially, it contains the rounding mode.

	a = ctx.FromString("1234.5678") // convert string to Quad. If error, a status bit in ctx will be set.
	b = ctx.FromString("-45.7")     //   Error bits in status can be tested with ctx.Error() at any time. Errors are cumulative, and only ctx.ResetStatus() will clear them.

	r = ctx.Add(a, b) // r = a + b       If ctx already contains an error in status, the result of any arithmetic operation is undefined, most probably Nan.
	// ...
	// you can put other operations here
	// ...

	fmt.Println("r", r.String())

	if err := ctx.Error(); err != nil { // you can just check for after a series of operations have been done
		log.Fatalf("ERROR OCCURRED !!!!!!!   %v\n", err)
	}


Tech note

This package uses cgo to call functions in the C decNumber package.

    All parameters are sent and received BY VALUE, because they are small.
    Quad is only 128 bits, Context are also sent and received by value because they are small struct (28 bytes). Sane for strings, which are small arrays embedded in struct.
    This way, there is no need to make complex things with pointers between Go and C world, and it is as fast, or even faster.

*/
package decnum
