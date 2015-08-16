/*
  This is just a little example to show how to use the decnum package.

  Pass two numbers, or Inf, or Nan, as arguments, and the program will make some computation with it.
*/
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/rin01/decnum"
)

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

	fmt.Println(decnum.Version())
	fmt.Println(decnum.DecQuad_module_MACROS) // just for info about the macros defined in C decQuad module

	//========================= division a/b ==================================

	fmt.Println("")
	fmt.Println("========= division a/b ==========")

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
	fmt.Println("a is:  ", a)
	fmt.Println("b is:  ", b)

	fmt.Println("")
	fmt.Println("r is:  ", r, "we see that an uninitialized Quad contains garbage.")

	r = ctx.Divide(a, b) // but no need to initialize r with decnum.Zero(), because its value is overwritten by the operation
	// ...
	// you can put other operations here, you will check for error after the series of operations
	// ...

	fmt.Println("r", r)

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

	fmt.Printf("comparison of %s and %s is %d\n", a, b, comp)

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

	for i:=0; i<50; i++ {
		h = ctx.Add(h, hh)
		fmt.Printf("%d   %s\n", i, h)
	}

	if err := ctx.Error(); err != nil {
		log.Printf("ERROR OCCURED !!!!!!!   %v\n", err)
	}




	fmt.Println("beuu", a)
	fmt.Printf("beuu %s\n", a)
}
