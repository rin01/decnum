/*
Package decnum is a Go binding around C decNumber package, for calculation with decimal floating point numbers.
Decimal base-10 data type is important for financial calculations.

Godoc: https://godoc.org/github.com/rin01/decnum


Example of use

	var (
		ctx decnum.Context
		a   decnum.Quad  //    uninitialized value is 0e-6176. It is really zero, but with the highest negative exponent for this type.
		b   decnum.Quad  //    If you prefer a variable to be 0, that is, 0e0, do      x = decnum.Zero()
		r   decnum.Quad
	)

	ctx.InitDefaultQuad()           // initialize context with default settings for Quad operations. Context contains the rounding mode and accumulates errors in its status field.

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


Internal representation of numbers

It is easier to work with this package if you keep in mind the following representation for numbers:

         (-1)^sign  coefficient * 10^exponent
         where coefficient is an integer storing 34 digits.

         12.345678e2    is     12345678E-4
         123e5          is          123E+5
         0              is            0E+0
         1              is            1E+0
         1.00           is          100E+0
         34.560         is        34560E-3

This representation is important to grasp when using functions like ToIntegral, Quantize, IsInteger, etc.


Represention of numbers for display

When numbers are displayed, the functions that convert them to string like ToString use a different format:

         (-1)^sign  c.oefficient * 10^exp
         where c.oefficient is a fractional number with one digit before fractional point

         1234.567e-12       is printed as     1.234567E-9
         650e4              is printed as         6.50E+6

This representation is well suited for displaying numbers, but not to work with other functions in this package.


Test

The functions in this library have been all tested in the file https://github.com/rin01/decnum/blob/master/mydecquad_test.go.
It contains a lot of interesting cases.


Tech note

This package uses cgo to call functions in the C decNumber package.

All parameters are sent and received BY VALUE, because they are small.

Quad is only 128 bits, Context are also sent and received by value because they are small struct (28 bytes). Same for strings, which are small arrays embedded in struct.

This way, there is no need to make complex things with pointers between Go and C world, and it is as fast, or even faster.

*/
package decnum
