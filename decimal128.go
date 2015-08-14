package decnum




/*

#include "mydecquad.h"
*/
import "C"

import (
	"log"
	"errors"
)

// Note: in the comment block for cgo above, the path for LDFLAGS must be an "absolute" path.
//       If it is a relative path, it seems that it is relative to the current directory when "go build" is run.
//       As we want to run "go build decnum" from any location, the path must be absolute.
//       See   Issue 5428 in May 2013: cmd/ld: relative #cgo LDFLAGS -L does not work
//             This problem is still not resolved in May 2014.

// MyQuad is a struct that just contains a C.decQuad value.
type MyQuad struct {
	val C.decQuad
}

const (
	S_DECQUAD_Pmax        = C.DECQUAD_Pmax          // number of digits in coefficient
	S_DECQUAD_Bytes       = C.DECQUAD_Bytes         // size in bytes of decQuad
	S_DECQUAD_String      = C.DECQUAD_String        // buffer capacity for C.decQuadToString()
	S_STRING_RAW_CAPACITY = C.S_STRING_RAW_CAPACITY // buffer capacity for C.mdq_to_mallocated_string_raw()
)

var (
		ERROR_DEC_UNLISTED = errors.New("decnum: Unlisted")
		ERROR_DEC_INVALID_OPERATION = errors.New("decnum: Invalid operation")
		ERROR_DEC_DIVISION_BY_ZERO = errors.New("decnum: Division by zero")
		ERROR_DEC_OVERFLOW = errors.New("decnum: Overflow")
		ERROR_DEC_UNDERFLOW = errors.New("decnum: Underflow")
		ERROR_DEC_DIVISION_IMPOSSIBLE = errors.New("decnum: Division impossible")
		ERROR_DEC_DIVISION_UNDEFINED = errors.New("decnum: Division undefined")
		ERROR_DEC_CONVERSION_SYNTAX = errors.New("decnum: Conversion syntax")
		ERROR_DEC_INSUFFICIENT_STORAGE = errors.New("decnum: Insufficient storage")
		ERROR_DEC_INVALID_CONTEXT = errors.New("decnum: Invalid Context")
)

/************************************************************************/
/*                                                                      */
/*                            init function                             */
/*                                                                      */
/************************************************************************/

func init() {

	C.mdq_init()

	log.Printf("decQuad module: DECDPUN %d, DECSUBSET %d, DECEXTFLAG %d. Constants DECQUAD_Pmax %d, DECQUAD_String %d DECQUAD_Bytes %d.", C.DECDPUN, C.DECSUBSET, C.DECEXTFLAG, C.DECQUAD_Pmax, C.DECQUAD_String, C.DECQUAD_Bytes)

	if S_DECQUAD_Bytes != 16 { // S_DECQUAD_Bytes MUST NOT BE > 16, because Append_compressed_bytes() will silently fail if it is not the case
		panic("S_DECQUAD_Bytes != 16")
	}

}

/************************************************************************/
/*                                                                      */
/*                          utility functions                           */
/*                                                                      */
/************************************************************************/

// get_rsql_error_message_id converts C.xxx error code into rsql.ERROR_xxx error code.
// All C.MDQ_ERROR_DEC_XXX errors come from the C decNumber library.
//
func get_error(mdqerr C.uint32_t) error {

	switch mdqerr {
	case C.MDQ_ERROR_DEC_UNLISTED:
		return ERROR_DEC_UNLISTED
	case C.MDQ_ERROR_DEC_INVALID_OPERATION:
		return ERROR_DEC_INVALID_OPERATION
	case C.MDQ_ERROR_DEC_DIVISION_BY_ZERO:
		return ERROR_DEC_DIVISION_BY_ZERO
	case C.MDQ_ERROR_DEC_OVERFLOW:
		return ERROR_DEC_OVERFLOW
	case C.MDQ_ERROR_DEC_UNDERFLOW:
		return ERROR_DEC_UNDERFLOW
	case C.MDQ_ERROR_DEC_DIVISION_IMPOSSIBLE:
		return ERROR_DEC_DIVISION_IMPOSSIBLE
	case C.MDQ_ERROR_DEC_DIVISION_UNDEFINED:
		return ERROR_DEC_DIVISION_UNDEFINED
	case C.MDQ_ERROR_DEC_CONVERSION_SYNTAX:
		return ERROR_DEC_CONVERSION_SYNTAX
	case C.MDQ_ERROR_DEC_INSUFFICIENT_STORAGE:
		return ERROR_DEC_INSUFFICIENT_STORAGE
	case C.MDQ_ERROR_DEC_INVALID_CONTEXT:
		return ERROR_DEC_INVALID_CONTEXT
	}

	panic("never get here")
}

type Context struct {
	set C.decContext
}



/************************************************************************/
/*                                                                      */
/*                      arithmetic operations                           */
/*                                                                      */
/************************************************************************/

func (context *Context) Unary_minus(a MyQuad) (r MyQuad) {
        var result C.Result_t

        result = C.mdq_unary_minus(a.val, context.set)

	context.set = result.set
	
	return MyQuad{val: result.val}
}

func (context *Context) Add(a MyQuad, b MyQuad) (r MyQuad) {
        var result C.Result_t

        result = C.mdq_add(a.val, b.val, context.set)

	context.set = result.set
	
	return MyQuad{val: result.val}
}

/************************************************************************/
/*                                                                      */
/*                      conversion operations                           */
/*                                                                      */
/************************************************************************/

func (context *Context) From_int32(value int32) (r MyQuad) {

        C.decQuadFromInt32(&r.val, C.int32_t(value))

	return r
}

func (context *Context) From_int64(value int64) (r MyQuad) {
        var result C.Result_t

        result = C.mdq_from_int64(C.int64_t(value), context.set)

	context.set = result.set
	
	return MyQuad{val: result.val}
}



