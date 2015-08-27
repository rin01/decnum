/*
  This is just a little example to show how to use the decnum package.

  Pass two numbers, or Inf, or Nan, as arguments, and the program will make some computation with them.
*/
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/rin01/decnum"
)

func assert(val bool) {
	if val == false {
		panic("assertion failed")
	}
}

func main() {
	var (
		ctx decnum.Context
		a   decnum.Quad
		b   decnum.Quad
		r   decnum.Quad
	)

	if len(os.Args) != 3 {
		log.Fatal("2 numbers are required as argument")
	}

	fmt.Println(decnum.DecNumber_C_Version())
	fmt.Println(decnum.DecNumber_C_MACROS()) // just for info about the macros defined in C decQuad module

	//========================= display a, b and r ==================================

	fmt.Println("")
	fmt.Println("========= display a, b and r ==========")

	ctx.InitDefaultQuad() // initialize context with default settings for Quad operations. Essentially, it contains the rounding mode.

	fmt.Printf("rounding: %s\n", ctx.Rounding()) // display default rounding mode

	ctx.SetRounding(decnum.ROUND_UP) // we can change it
	fmt.Printf("rounding: %s\n", ctx.Rounding())

	ctx.SetRounding(decnum.ROUND_HALF_EVEN) // we can change it again
	fmt.Printf("rounding: %s\n", ctx.Rounding())

	a = ctx.FromString(os.Args[1]) // convert first argument to Quad
	b = ctx.FromString(os.Args[2]) // convert 2nd argument to Quad

	if err := ctx.Error(); err != nil { // check if string conversion succeeded
		fmt.Println("ERROR: incorrect string input...")
	}

	fmt.Println("")
	fmt.Printf("a is:  %s\n", a)
	fmt.Printf("b is:  %s\n", b)

	fmt.Println("")
	fmt.Printf("r is:   %s    %s\n", r, "// we see that an uninitialized Quad is 0e-6176, that is, 0.00000000.........000000000.")
	fmt.Printf("5 + r = %s \n", ctx.Add(ctx.FromString("5"), r)) // 5.00000000..........00000000
	fmt.Println("")

	r = decnum.Zero()
	fmt.Println("r = decnum.Zero()") // if you want to have 5 without all fractional 0s, you should initialize r as     r = decnum.Zero()
	fmt.Printf("r is:   %s\n", r)
	fmt.Printf("5  + r = %s \n", ctx.Add(ctx.FromString("5"), r)) // result is 5
	fmt.Println("")


	//========================= division a/b ==================================

	fmt.Println("")
	fmt.Println("========= division a/b ==========")

	ctx.ResetStatus() // clear the status

	r = ctx.Divide(a, b)
	// ...
	// you can put other operations here, you will check for error after the series of operations
	// ...

	fmt.Println("r = a/b; r = ", r)

	status := ctx.Status()
	fmt.Printf("status: %s\n", status)

	if err := ctx.Error(); err != nil { // check if an error flag has been set. No need to check for error after each operation, wee can just check it after a series of operations have been done.
		log.Printf("ERROR OCCURED !!!!!!!   %v\n", err)
	}

	//=========================== convert 'a' to int64 =================================

	fmt.Println("")
	fmt.Println("========= convert 'a' to int64 ==========")

	ctx.ResetStatus() // clear the status

	// you can put another series of operations here

	// ...

	var x int64

	x = ctx.ToInt64(a, decnum.ROUND_HALF_EVEN)

	if err := ctx.Error(); err != nil { // check for errors
		log.Printf("ERROR OCCURED !!!!!!!   %v\n", err)
	}

	fmt.Printf("%s converted to int64 is %d\n", a, x) // you can always print a Quad, it always contains a valid value, even after errors

	//============================ compare 'a' and 'b' ================================

	fmt.Println("")
	fmt.Println("========= compare 'a' and 'b' ==========")

	ctx.ResetStatus() // clear the status

	// you can put another series of operations here

	// ...

	var comp decnum.Cmp_t

	comp = ctx.Compare(a, b) // note: Compare doesn't set status flag

	if err := ctx.Error(); err != nil {
		log.Fatalf("ERROR OCCURED !!!!!!!   %v\n", err)
	}

	fmt.Printf("comparison of %s and %s is %s\n", a, b, comp)

	//============================ quantize 'a' with pattern 'b' ================================

	fmt.Println("")
	fmt.Println("========= quantize 'a' with pattern 'b' ==========")

	ctx.ResetStatus() // clear the status

	// you can put another series of operations here

	// ...

	var q decnum.Quad

	q = ctx.Quantize(a, b)

	if err := ctx.Error(); err != nil {
		log.Printf("ERROR OCCURED !!!!!!!   %v\n", err)
	}

	fmt.Printf("quantization of %s with %s is %s\n", a, b, q)

	//============================ loop ================================

	fmt.Println("")
	fmt.Println("========= loop ==========")

	ctx.ResetStatus() // clear the status

	var h decnum.Quad = decnum.Zero()
	var hh decnum.Quad = ctx.FromInt32(1000000000)

	for i := 0; i < 7; i++ {
		h = ctx.Add(h, hh)
		fmt.Printf("%d   %s\n", i, h)
	}

	if err := ctx.Error(); err != nil {
		log.Printf("ERROR OCCURED !!!!!!!   %v\n", err)
	}

	//============================ rounding ================================

	fmt.Println("")
	fmt.Println("========= rounding (quantizing) ==========")

	ctx.ResetStatus() // clear the status

	var g_up decnum.Quad
	var g_down decnum.Quad

	ctx.SetRounding(decnum.ROUND_UP)
	assert(ctx.Rounding() == decnum.ROUND_UP)
	g_up = ctx.Quantize(a, decnum.One())

	ctx.SetRounding(decnum.ROUND_DOWN)
	assert(ctx.Rounding() == decnum.ROUND_DOWN)
	g_down = ctx.Quantize(a, decnum.One())

	ctx.SetRounding(decnum.ROUND_DEFAULT)
	assert(ctx.Rounding() == decnum.ROUND_HALF_EVEN)

	if err := ctx.Error(); err != nil {
		log.Printf("ERROR OCCURED !!!!!!!   %v\n", err)
	}

	fmt.Printf("%s rounded up   is %s\n", a, g_up)
	fmt.Printf("%s rounded down is %s\n", a, g_down)

}
