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
		err error
		ctx decnum.Context
		a   decnum.DecQuad
		b   decnum.DecQuad
		r   decnum.DecQuad
	)

	if len(os.Args) != 3 {
		log.Fatal("2 numbers are required as argument")
	}

	fmt.Println(decnum.DecQuad_module_MACROS) // just for info about the macros defined in C decQuad module

	//========================= division a/b ==================================

	fmt.Println("")
	fmt.Println("========= division a/b ==========")

	ctx.Init(decnum.DEFAULT_DECQUAD) // initialize context with default settings for DecQuad operations. Essentially, it contains the rounding mode.

	fmt.Printf("rounding: %s\n", ctx.Rounding()) // display default rounding mode

	ctx.SetRounding(decnum.ROUND_UP) // we can change it
	fmt.Printf("rounding: %s\n", ctx.Rounding())

	ctx.SetRounding(decnum.ROUND_HALF_EVEN) // we can change it again
	fmt.Printf("rounding: %s\n", ctx.Rounding())

	if a, err = ctx.FromString(os.Args[1]); err != nil { // convert first argument to DecQuad
		log.Fatal(err)
	}
	if b, err = ctx.FromString(os.Args[2]); err != nil { // convert 2nd argument to DecQuad
		log.Fatal(err)
	}

	fmt.Println("")
	fmt.Println("a is:  ", a.String())
	fmt.Println("b is:  ", b.String())

	fmt.Println("")
	fmt.Println("r is:  ", r.String(), "we see that an uninitialized DecQuad contains garbage.")

	r = ctx.Divide(a, b) // but no need to initialize r with decnum.Zero(), because its value is overwritten by the operation
	// ...
	// you can put other operations here, you will check for error after the series of operations
	// ...

	fmt.Println("r", r.String())

	status := ctx.Status()
	fmt.Printf("status: %d\n", status)

	if err := ctx.Error(); err != nil { // check if an error flag has been set. No need to check for error after each operation, wee can just check it after a series of operations have been done.
		log.Printf("ERROR OCCURS !!!!!!!   %v\n", err)
	}

	//=========================== convert 'a' to int64 =================================

	fmt.Println("")
	fmt.Println("========= convert 'a' to int64 ==========")

	ctx.ResetStatus() // clear the status

	// you can put another series of operations here

	// ...

	var x int64

	x = ctx.ToInt64(a, decnum.ROUND_DOWN)

	if err := ctx.Error(); err != nil { // check for errors
		log.Printf("ERROR OCCURS !!!!!!!   %v\n", err)
	}

	fmt.Printf("%s converted to int64 is %d\n", a.String(), x) // you can always print a DecQuad, it always contains a valid value, even after errors

	//============================ compare 'a' and 'b' ================================

	fmt.Println("")
	fmt.Println("========= compare 'a' and 'b' ==========")

	ctx.ResetStatus() // clear the status

	// you can put another series of operations here

	// ...

	var comp decnum.DecQuad

	comp = ctx.Compare(a, b)

	if err := ctx.Error(); err != nil {
		log.Fatalf("ERROR OCCURS !!!!!!!   %v\n", err)
	}

	fmt.Printf("comparison of %s and %s is %s\n", a.String(), b.String(), comp.String())

	//============================ quantize 'a' with pattern 'b' ================================

	fmt.Println("")
	fmt.Println("========= quantize 'a' with pattern 'b' ==========")

	ctx.ResetStatus() // clear the status

	// you can put another series of operations here

	// ...

	var q decnum.DecQuad

	q = ctx.Quantize(a, b)

	if err := ctx.Error(); err != nil {
		log.Printf("ERROR OCCURS !!!!!!!   %v\n", err)
	}

	fmt.Printf("quantization of %s with %s is %s\n", a.String(), b.String(), q.String())
}
