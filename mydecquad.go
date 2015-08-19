package decnum

/*

#include "mydecquad.h"
*/
import "C"

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync"
	"unsafe"
)

// Note: in the comment block for cgo above, if LDFLAGS is used, the path for LDFLAGS must be an "absolute" path.
//       If it is a relative path, it seems that it is relative to the current directory when "go build" is run.
//       As we want to run "go build decnum" from any location, the path must be absolute.
//       See   Issue 5428 in May 2013: cmd/ld: relative #cgo LDFLAGS -L does not work
//             This problem is still not resolved in May 2014.
//
//       *** this note is just a remainder, as LDFLAGS is not used here ***

func assert(val bool) {
	if val == false {
		panic("assertion failed")
	}
}

func assert_sane(context *Context) {
	if context.sane == false {
		panic("context not initialized")
	}
}

// Quad is just a struct with a C.decQuad value, which is an array of 16 bytes.
type Quad struct {
	val C.decQuad // array of 16 bytes
}

/************************************************************************/
/*                                                                      */
/*                 global constants and variables                       */
/*                                                                      */
/************************************************************************/

const (
	DECQUAD_Pmax   = C.DECQUAD_Pmax   // number of digits in coefficient
	DECQUAD_Bytes  = C.DECQUAD_Bytes  // size in bytes of decQuad
	DECQUAD_String = C.DECQUAD_String // buffer capacity for C.decQuadToString()
)

var (
	ERROR_DEC_UNLISTED             = errors.New("decnum: Unlisted")
	ERROR_DEC_INVALID_OPERATION    = errors.New("decnum: Invalid operation")
	ERROR_DEC_DIVISION_BY_ZERO     = errors.New("decnum: Division by zero")
	ERROR_DEC_OVERFLOW             = errors.New("decnum: Overflow")
	ERROR_DEC_UNDERFLOW            = errors.New("decnum: Underflow")
	ERROR_DEC_DIVISION_IMPOSSIBLE  = errors.New("decnum: Division impossible")
	ERROR_DEC_DIVISION_UNDEFINED   = errors.New("decnum: Division undefined")
	ERROR_DEC_CONVERSION_SYNTAX    = errors.New("decnum: Conversion syntax")
	ERROR_DEC_INSUFFICIENT_STORAGE = errors.New("decnum: Insufficient storage")
	ERROR_DEC_INVALID_CONTEXT      = errors.New("decnum: Invalid Context")
)

// g_nan, g_zero and g_one are private variable, because else, a user of the package can change their value by doing decnum.G_ZERO = ...

var (
	g_nan  Quad = nan_for_varinit()     // a constant Quad with value Nan. It runs BEFORE init().
	g_zero Quad = zero_for_varinit()    // a constant Quad with value 0. It runs BEFORE init().
	g_one  Quad = quad_for_varinit("1") // a constant Quad with value 1. It runs BEFORE init().
)

// used only to initialize the global variable g_nan.
//
// So, it runs BEFORE init().
//
func nan_for_varinit() (r Quad) {
	var val C.decQuad

	val = C.mdq_nan()

	return Quad{val: val}
}

// used only to initialize the global variable g_zero.
//
// So, it runs BEFORE init().
//
func zero_for_varinit() (r Quad) {
	var val C.decQuad

	val = C.mdq_zero()

	return Quad{val: val}
}

// used only to initialize some global variables, like g_one.
//
// So, it runs BEFORE init().
//
func quad_for_varinit(s string) (r Quad) {
	var (
		ctx Context
		val Quad
	)

	ctx.InitDefaultQuad()

	val = ctx.FromString(s)

	if err := ctx.Error(); err != nil {
		panic("decnum: initialization error in quad_for_varinit()")
	}

	return val
}

/************************************************************************/
/*                                                                      */
/*                       init and version functions                     */
/*                                                                      */
/************************************************************************/

var (
	decNumber_C_version string = C.GoString(C.decQuadVersion()) // version of the original C decNumber package

	decNumber_C_MACROS string = fmt.Sprintf("decQuad module: DECDPUN %d, DECSUBSET %d, DECEXTFLAG %d. Constants DECQUAD_Pmax %d, DECQUAD_String %d DECQUAD_Bytes %d.",
		C.DECDPUN, C.DECSUBSET, C.DECEXTFLAG, C.DECQUAD_Pmax, C.DECQUAD_String, C.DECQUAD_Bytes) // macros defined by the C decNumber module
)

func init() {
	C.mdq_init()

	if DECQUAD_Bytes != 16 { // 16 bytes == 128 bits
		panic("DECQUAD_Bytes != 16")
	}

	assert(C.DECSUBSET == 0) // because else, we should define Flag_Lost_digits as status flag

	assert(POOL_BUFF_CAPACITY > DECQUAD_Pmax)
	assert(POOL_BUFF_CAPACITY > DECQUAD_String)

}

// DecNumber_C_Version returns the version of the original C decNumber package.
//
func DecNumber_C_Version() string {

	return decNumber_C_version
}

// DecNumber_C_MACROS returns the values of macros defined in the original C decNumber package.
//
func DecNumber_C_MACROS() string {

	return decNumber_C_MACROS
}

/************************************************************************/
/*                                                                      */
/*                              Context                                 */
/*                                                                      */
/************************************************************************/

type Status_t uint32

const (
	Flag_Conversion_syntax    Status_t = C.DEC_Conversion_syntax    // error flag
	Flag_Division_by_zero     Status_t = C.DEC_Division_by_zero     // error flag
	Flag_Division_impossible  Status_t = C.DEC_Division_impossible  // error flag
	Flag_Division_undefined   Status_t = C.DEC_Division_undefined   // error flag
	Flag_Insufficient_storage Status_t = C.DEC_Insufficient_storage // error flag
	Flag_Inexact              Status_t = C.DEC_Inexact              // informational flag
	Flag_Invalid_context      Status_t = C.DEC_Invalid_context      // error flag
	Flag_Invalid_operation    Status_t = C.DEC_Invalid_operation    // error flag
	Flag_Overflow             Status_t = C.DEC_Overflow             // error flag
	Flag_Clamped              Status_t = C.DEC_Clamped              // informational flag
	Flag_Rounded              Status_t = C.DEC_Rounded              // informational flag
	Flag_Subnormal            Status_t = C.DEC_Subnormal            // informational flag
	Flag_Underflow            Status_t = C.DEC_Underflow            // error flag. E.g. 1e-6000/1e1000

	//Flag_Lost_digits          Status_t = C.DEC_Lost_digits        // informational flag. Exists only if DECSUBSET is set, which is not the case by default
)

const ErrorMask Status_t = C.DEC_Errors // ErrorMask is the bitmask of the error flags, ORed together. After a series of operations, if status & decnum.ErrorMask != 0, an error has occured, e.g. division by 0.

// String representation of a single flag (status with one bit set).
//
func (flag Status_t) flag_string() string {

	if flag == 0 {
		return ""
	}

	switch flag {
	case Flag_Conversion_syntax:
		return "Conversion_syntax"
	case Flag_Division_by_zero:
		return "Division_by_zero"
	case Flag_Division_impossible:
		return "Division_impossible"
	case Flag_Division_undefined:
		return "Division_undefined"
	case Flag_Insufficient_storage:
		return "Insufficient_storage"
	case Flag_Inexact:
		return "Inexact"
	case Flag_Invalid_context:
		return "Invalid_context"
	case Flag_Invalid_operation:
		return "Invalid_operation"
	//case Flag_Lost_digits:
	//return "Lost_digits"
	case Flag_Overflow:
		return "Overflow"
	case Flag_Clamped:
		return "Clamped"
	case Flag_Rounded:
		return "Rounded"
	case Flag_Subnormal:
		return "Subnormal"
	case Flag_Underflow:
		return "Underflow"
	default:
		return "Unknown status flag"
	}
}

// String representation of a status.
// status can have many flags set.
//
func (status Status_t) String() string {
	var (
		s    string
		flag Status_t
	)

	for i := Status_t(0); i < 32; i++ {
		flag = Status_t(0x0001 << i)
		if status&flag != 0 {
			if s == "" {
				s = flag.flag_string()
			} else {
				s += ";" + flag.flag_string()
			}
		}
	}

	return s
}

type Round_mode_t int

// Rounding mode is used if rounding is necessary during an operation.
const (
	ROUND_CEILING   Round_mode_t = C.DEC_ROUND_CEILING   // Round towards +Infinity.
	ROUND_DOWN      Round_mode_t = C.DEC_ROUND_DOWN      // Round towards 0 (truncation).
	ROUND_FLOOR     Round_mode_t = C.DEC_ROUND_FLOOR     // Round towards –Infinity.
	ROUND_HALF_DOWN Round_mode_t = C.DEC_ROUND_HALF_DOWN // Round to nearest; if equidistant, round down.
	ROUND_HALF_EVEN Round_mode_t = C.DEC_ROUND_HALF_EVEN // Round to nearest; if equidistant, round so that the final digit is even.
	ROUND_HALF_UP   Round_mode_t = C.DEC_ROUND_HALF_UP   // Round to nearest; if equidistant, round up.
	ROUND_UP        Round_mode_t = C.DEC_ROUND_UP        // Round away from 0.
	ROUND_05UP      Round_mode_t = C.DEC_ROUND_05UP      // The same as DEC_ROUND_UP, except that rounding up only occurs if the digit to be rounded up is 0 or 5 and after Overflow the result is the same as for DEC_ROUND_DOWN.
	ROUND_DEFAULT   Round_mode_t = ROUND_HALF_EVEN       // The same as DEC_ROUND_HALF_EVEN.
)

func (rounding Round_mode_t) String() string {

	switch rounding {
	case ROUND_CEILING:
		return "ROUND_CEILING"
	case ROUND_DOWN:
		return "ROUND_DOWN"
	case ROUND_FLOOR:
		return "ROUND_FLOOR"
	case ROUND_HALF_DOWN:
		return "ROUND_HALF_DOWN"
	case ROUND_HALF_EVEN:
		return "ROUND_HALF_EVEN"
	case ROUND_HALF_UP:
		return "ROUND_HALF_UP"
	case ROUND_UP:
		return "ROUND_UP"
	case ROUND_05UP:
		return "ROUND_05UP"
	default:
		return "Unknown rounding mode"
	}
}

// Context contains the rounding mode, and a status field that records exceptional conditions, some of which are considered as error, e.g. division by 0, underlow for operations like 1e-6000/1e1000, overflow, etc.
// For decQuad usage, only these two fields are used.
//
// When an error occurs during an operation, the result will probably be NaN or infinite, or a infinitesimal number if underflow.
// If conversion error to int32, int64, etc, it will be 0.
//
type Context struct {
	sane bool // if true, it can be used because it has been initialized with ctx.InitDefaultQuad()

	set C.decContext
}

type Context_kind_t uint32

const (
	DEFAULT_DECQUAD Context_kind_t = C.DEC_INIT_DECQUAD // default Context settings for decQuad operations
)

// initialize is used to initialize a context with default value for rounding mode, and clears status field.
//
func (context *Context) initialize(kind Context_kind_t) {

	context.set = C.mdq_context_default(context.set, C.uint32_t(kind))

	context.sane = true
}

// InitDefaultQuad is used to initialize a context with default value for Quad operations. It sets rounding mode, and clears status field.
//
func (context *Context) InitDefaultQuad() {

	context.set = C.mdq_context_default(context.set, C.uint32_t(DEFAULT_DECQUAD))

	context.sane = true
}

// Rounding returns the rounding mode of the context.
//
func (context *Context) Rounding() Round_mode_t {
	assert_sane(context)

	return Round_mode_t(C.mdq_context_get_rounding(context.set))
}

// SetRounding sets the rounding mode of the context.
//
func (context *Context) SetRounding(rounding Round_mode_t) {
	assert_sane(context)

	context.set = C.mdq_context_set_rounding(context.set, C.int(rounding))
}

// Status returns the status of the context.
//
// After a series of operations, the status contains the accumulated errors or informational flags that occurred during all the operations.
//
// Beware: the status can contain informational flags, like Flag_Inexact, which is not an error.
//
// So, to find the real errors, you must discard the non-error bits of the status as follows:
//      status = ctx.Status() & decnum.ErrorMask
//      if status != 0 {
//             ... error occurred
//      }
//
// It is easier to use the context.Error method to check for errors.
//
func (context *Context) Status() Status_t {
	assert_sane(context)

	return Status_t(C.mdq_context_get_status(context.set))
}

// SetStatus sets a status bit in the status of the context.
//
// Normally, only library modules use this function. Applications have no reason to set status bits.
//
func (context *Context) SetStatus(flag Status_t) {
	assert_sane(context)

	context.set = C.mdq_context_set_status(context.set, C.uint32_t(flag))
}

// ResetStatus clears all bits of the status field of the context.
// You can continue to use this context for a new series of operations.
//
func (context *Context) ResetStatus() {
	assert_sane(context)

	context.set = C.mdq_context_zero_status(context.set)
}

// Error checks if status contains a flag that should be considered as an error.
// In this case, the resut of the operations contains Nan or Infinite, or an infinitesimal number if Underflow.
// It contains 0 if conversion to int64, float64, etc failed.
//
// It is not necessary and not usual to check for errors after each operation.
// You can make many arithmetic operations in a row, and check ctx.Error() when you are finished.
//
// If an error occured, the subsequent operations will work on operands that will frequently be Nan, and Nan will propagate.
// But if you convert a Quad to a int32 and overflow occurs, the value returned is 0, making the error not so obvious to detect.
//
// So, don't forget to call ctx.Error at the end of each series of operations.
//
// Errors accumulate in the status field of Context, setting bits but never clearing them. So, an error will never be lost.
//
// Before you begin a new series of operations, you must clear the Context status field with ctx.ResetStatus().
//
func (context *Context) Error() error {
	var status Status_t
	assert_sane(context)

	status = context.Status()

	status = status & ErrorMask // discard informational flags, keep only error flags

	if status != 0 {
		return fmt.Errorf("decnum error: %s", status.String())
	}

	return nil
}

/************************************************************************/
/*                                                                      */
/*                      arithmetic operations                           */
/*                                                                      */
/************************************************************************/

type Cmp_t uint32 // result of Compare

const (
	CMP_LESS    Cmp_t = C.CMP_LESS    // 1
	CMP_EQUAL   Cmp_t = C.CMP_EQUAL   // 2
	CMP_GREATER Cmp_t = C.CMP_GREATER // 4
	CMP_NAN     Cmp_t = C.CMP_NAN     // 8
)

func (cmp Cmp_t) String() string {

	switch cmp {
	case CMP_LESS:
		return "CMP_LESS"
	case CMP_EQUAL:
		return "CMP_EQUAL"
	case CMP_GREATER:
		return "CMP_GREATER"
	case CMP_NAN:
		return "CMP_NAN"
	default:
		return "Unknown Cmp_t"
	}
}

// Zero returns 0 Quad value.
//
//     r = Zero()  // assign 0 to the Quad r
//
func Zero() (r Quad) {

	return g_zero
}

// One returns 1 Quad value.
//
//     r = One()  // assign 1 to the Quad r
//
func One() (r Quad) {

	return g_one
}

// NaN returns NaN Quad value.
//
//     r = NaN()  // assign NaN to the Quad r
//
func NaN() (r Quad) {

	return g_nan
}

// Copy returns a copy of a.
//
// But it is easier to just use '=' :
//
//        a = r
//
func Copy(a Quad) (r Quad) {

	return a
}

// Minus returns -a.
//
func (context *Context) Minus(a Quad) (r Quad) {
	var result C.Ret_decQuad_t
	assert_sane(context)

	result = C.mdq_minus(a.val, context.set)

	context.set = result.set
	return Quad{val: result.val}
}

// Add returns a + b.
//
func (context *Context) Add(a Quad, b Quad) (r Quad) {
	var result C.Ret_decQuad_t
	assert_sane(context)

	result = C.mdq_add(a.val, b.val, context.set)

	context.set = result.set
	return Quad{val: result.val}
}

// Subtract returns a - b.
//
func (context *Context) Subtract(a Quad, b Quad) (r Quad) {
	var result C.Ret_decQuad_t
	assert_sane(context)

	result = C.mdq_subtract(a.val, b.val, context.set)

	context.set = result.set
	return Quad{val: result.val}
}

// Multiply returns a * b.
//
func (context *Context) Multiply(a Quad, b Quad) (r Quad) {
	var result C.Ret_decQuad_t
	assert_sane(context)

	result = C.mdq_multiply(a.val, b.val, context.set)

	context.set = result.set
	return Quad{val: result.val}
}

// Divide returns a/b.
//
func (context *Context) Divide(a Quad, b Quad) (r Quad) {
	var result C.Ret_decQuad_t
	assert_sane(context)

	result = C.mdq_divide(a.val, b.val, context.set)

	context.set = result.set
	return Quad{val: result.val}
}

// DivideInteger returns the integral part of a/b.
//
func (context *Context) DivideInteger(a Quad, b Quad) (r Quad) {
	var result C.Ret_decQuad_t
	assert_sane(context)

	result = C.mdq_divide_integer(a.val, b.val, context.set)

	context.set = result.set
	return Quad{val: result.val}
}

// Remainder returns the modulo of a and b.
//
func (context *Context) Remainder(a Quad, b Quad) (r Quad) {
	var result C.Ret_decQuad_t
	assert_sane(context)

	result = C.mdq_remainder(a.val, b.val, context.set)

	context.set = result.set
	return Quad{val: result.val}
}

// Abs returns the absolute value of a.
//
func (context *Context) Abs(a Quad) (r Quad) {
	var result C.Ret_decQuad_t
	assert_sane(context)

	result = C.mdq_abs(a.val, context.set)

	context.set = result.set
	return Quad{val: result.val}
}

// ToIntegral returns the value of a rounded to an integral value.
//
//      The representation of a number is:
//
//           (-1)^sign  coefficient * 10^exponent
//           where coefficient is an integer storing 34 digits.
//
//       - If exponent < 0, the least significant digits are discarded, so that new exponent becomes 0.
//             Internally, it calls Quantize(a, 1E0) with specified rounding.
//       - If exponent >= 0, the number remains unchanged.
//
//         E.g.     12.345678e2    is     12345678E-4     -->   1235E0
//                  123e5          is     123E5        remains   123E5
//
func (context *Context) ToIntegral(a Quad, round Round_mode_t) (r Quad) {
	var result C.Ret_decQuad_t
	assert_sane(context)

	result = C.mdq_to_integral(a.val, context.set, C.int(round))

	context.set = result.set
	return Quad{val: result.val}
}

// Quantize rounds a to the same pattern as b.
// b is just a model, its sign and coefficient value are ignored. Only its exponent is used.
// The result is the value of a, but with the same exponent as the pattern b.
// The rounding of the context is used.
//
// You can use this function with the proper rounding to round (e.g. set context rounding mode to ROUND_HALF_EVEN) or truncate (ROUND_DOWN) 'a'.
//
//      The representation of a number is:
//
//           (-1)^sign  coefficient * 10^exponent
//           where coefficient is an integer storing 34 digits.
//
// Examples:
//    quantization of 134.6454 with    0.00001    is   134.64540
//                    134.6454 with    0.00000    is   134.64540     the value of b has no importance
//                    134.6454 with 1234.56789    is   134.64540     the value of b has no importance
//                    134.6454 with 0.0001        is   134.6454
//                    134.6454 with 0.01          is   134.65
//                    134.6454 with 1             is   135
//                    134.6454 with 1000000000    is   135           the value of b has no importance
//                    134.6454 with 1E+2          is   1E+2
//
//		        123e32 with 1             sets Invalid_operation error flag in status
//		        123e32 with 1E1           is   1230000000000000000000000000000000E1
//		        123e32 with 10            sets Invalid_operation error flag in status
//
func (context *Context) Quantize(a Quad, b Quad) (r Quad) {
	var result C.Ret_decQuad_t
	assert_sane(context)

	result = C.mdq_quantize(a.val, b.val, context.set)

	context.set = result.set
	return Quad{val: result.val}
}

// Compare compares the value of a and b.
//
//     If a <  b,        returns CMP_LESS
//     If a == b,        returns CMP_GREATER
//     If a >  b,        returns CMP_EQUAL
//     If a or b is Nan, returns CMP_NAN
//
// Compare doesn't set status flag, as no error occurs when just reading numbers.
//
// Example:
//
//     if ctx.Compare(a, b) & (CMP_GREATER|CMP_EQUAL) != 0 { // if a >= b
//         ...
//     }
//
func (context *Context) Compare(a Quad, b Quad) Cmp_t {
	var result C.Ret_uint32_t
	assert_sane(context)

	result = C.mdq_compare(a.val, b.val, context.set)

	context.set = result.set
	return Cmp_t(result.val)
}

// Cmp returns true if comparison of a and b complies with comp_mask.
// It is easier to use than Compare.
//
// Cmp doesn't set status flag, as no error occurs when just reading numbers.
//
// Example:
//
//     if ctx.Cmp(a, b, CMP_GREATER|CMP_EQUAL) { // if a >= b
//         ...
//     }
//
func (context *Context) Cmp(a Quad, b Quad, comp_mask Cmp_t) bool {
	var result C.Ret_uint32_t
	assert_sane(context)

	result = C.mdq_compare(a.val, b.val, context.set)

	context.set = result.set
	if Cmp_t(result.val)&comp_mask != 0 {
		return true
	}

	return false
}

// Greater is same as Cmp(a, b, CMP_GREATER)
//
func (context *Context) Greater(a Quad, b Quad) bool {
	var result C.Ret_uint32_t
	assert_sane(context)

	result = C.mdq_compare(a.val, b.val, context.set)

	context.set = result.set
	if Cmp_t(result.val)&CMP_GREATER != 0 {
		return true
	}

	return false
}

// GreaterEqual is same as Cmp(a, b, CMP_GREATER|CMP_EQUAL)
//
func (context *Context) GreaterEqual(a Quad, b Quad) bool {
	var result C.Ret_uint32_t
	assert_sane(context)

	result = C.mdq_compare(a.val, b.val, context.set)

	context.set = result.set
	if Cmp_t(result.val)&(CMP_GREATER|CMP_EQUAL) != 0 {
		return true
	}

	return false
}

// Equal is same as Cmp(a, b, CMP_EQUAL)
//
func (context *Context) Equal(a Quad, b Quad) bool {
	var result C.Ret_uint32_t
	assert_sane(context)

	result = C.mdq_compare(a.val, b.val, context.set)

	context.set = result.set
	if Cmp_t(result.val)&CMP_EQUAL != 0 {
		return true
	}

	return false
}

// LessEqual is same as Cmp(a, b, CMP_LESS|CMP_EQUAL)
//
func (context *Context) LessEqual(a Quad, b Quad) bool {
	var result C.Ret_uint32_t
	assert_sane(context)

	result = C.mdq_compare(a.val, b.val, context.set)

	context.set = result.set
	if Cmp_t(result.val)&(CMP_LESS|CMP_EQUAL) != 0 {
		return true
	}

	return false
}

// Less is same as Cmp(a, b, CMP_LESS)
//
func (context *Context) Less(a Quad, b Quad) bool {
	var result C.Ret_uint32_t
	assert_sane(context)

	result = C.mdq_compare(a.val, b.val, context.set)

	context.set = result.set
	if Cmp_t(result.val)&CMP_LESS != 0 {
		return true
	}

	return false
}

// IsFinite returns true if a is not Infinite, nor Nan.
//
func (a Quad) IsFinite() bool {

	if C.mdq_is_finite(a.val) != 0 {
		return true
	}

	return false
}

// IsInteger returns true if a is finite and has exponent=0.
//
//      The number representation is:
//
//           (-1)^sign  coefficient * 10^exponent
//           where coefficient is an integer storing 34 digits.
//
//      If the number in the above representation has exponent=0, then IsInteger returns true.
//
//      0              0E+0        returns true
//      1              1E+0        returns true
//      12.34e2     1234E+0        returns true
//
//      0.0000         0E-4        returns false
//      1.0000     10000E-4        returns false
//     -12.34e5    -1234E+3        returns false
//      1e3            1E+3        returns false
//
func (a Quad) IsInteger() bool {

	if C.mdq_is_integer(a.val) != 0 {
		return true
	}

	return false
}

// IsInfinite returns true if a is Infinite.
//
func (a Quad) IsInfinite() bool {

	if C.mdq_is_infinite(a.val) != 0 {
		return true
	}

	return false
}

// IsNaN returns true if a is Nan.
//
func (a Quad) IsNaN() bool {

	if C.mdq_is_nan(a.val) != 0 {
		return true
	}

	return false
}

// IsPositive returns true if a > 0 and not Nan.
//
func (a Quad) IsPositive() bool {

	if C.mdq_is_positive(a.val) != 0 {
		return true
	}

	return false
}

// IsZero returns true if a == 0.
//
func (a Quad) IsZero() bool {

	if C.mdq_is_zero(a.val) != 0 {
		return true
	}

	return false
}

// IsNegative returns true if a < 0 and not Nan.
//
func (a Quad) IsNegative() bool {

	if C.mdq_is_negative(a.val) != 0 {
		return true
	}

	return false
}

// Max returns the larger of a and b.
// If either a or b is NaN then the other argument is the result.
//
func (context *Context) Max(a Quad, b Quad) (r Quad) {
	var result C.Ret_decQuad_t
	assert_sane(context)

	result = C.mdq_max(a.val, b.val, context.set)

	context.set = result.set
	return Quad{val: result.val}
}

// Min returns the smaller of a and b.
// If either a or b is NaN then the other argument is the result.
//
func (context *Context) Min(a Quad, b Quad) (r Quad) {
	var result C.Ret_decQuad_t
	assert_sane(context)

	result = C.mdq_min(a.val, b.val, context.set)

	context.set = result.set
	return Quad{val: result.val}
}

/************************************************************************/
/*                                                                      */
/*                   conversion from string and numbers                 */
/*                                                                      */
/************************************************************************/

// FromString returns a Quad from a string.
//
// Special values "NaN" (also "qNaN"), "sNaN", "NaN123" (NaN with payload), "sNaN123" (sNaN with payload), "Infinity" (or "Inf", "+Inf"), "-Infinity" ( or "-Inf") are accepted.
//
//      Infinity and -Infinity, or Inf and -Inf, represent a value infinitely large.
//
//      NaN or qNaN, which means "Not a Number", represents an undefined result, when an arithmetic operation has failed. E.g. FromString("hello")
//                   NaN propagates to all subsequent operations, because if NaN is passed as argument, the result, will be NaN.
//                   These NaN are called "quiet NaN", because they don't set exceptional condition flag in status when passed as argument to an operation.
//
//      sNaN, or "signaling NaN", are created by FromString("sNaN"). When passed as argument to an operation, the result will be NaN, like with quiet NaN.
//                   But they will set (==signal) an exceptional condition flag in status, "Invalid_operation".
//                   Signaling NaN propagate to subsequent operation as ordinary NaN (quiet NaN), and not as "signaling NaN".
//
// Note that both NaN and sNaN can take an integer payload, e.g. NaN123, created by FromString("NaN123"), and it is up to you to give it a significance.
// sNaN and payload are not used often, and most probably, you won't use them.
//
func (context *Context) FromString(s string) (r Quad) {
	var (
		cs     *C.char
		result C.Ret_decQuad_t
	)
	assert_sane(context)

	s = strings.TrimSpace(s)

	cs = C.CString(s)
	defer C.free(unsafe.Pointer(cs))

	result = C.mdq_from_string(cs, context.set)

	context.set = result.set
	return Quad{val: result.val}
}

// FromInt32 returns a Quad from a int32 value.
//
// No error should occur, and context status will not change.
//
func (context *Context) FromInt32(value int32) (r Quad) {
	var result C.Ret_decQuad_t
	assert_sane(context)

	result = C.mdq_from_int32(C.int32_t(value), context.set)

	context.set = result.set
	return Quad{val: result.val}
}

// FromInt64 returns a Quad from a int64 value.
//
// No error should occur, and context status will not change.
//
func (context *Context) FromInt64(value int64) (r Quad) {
	var result C.Ret_decQuad_t
	assert_sane(context)

	result = C.mdq_from_int64(C.int64_t(value), context.set)

	context.set = result.set
	return Quad{val: result.val}
}

// FromFloat64 returns a Quad from a int64 value.
//
//   DEPRECATED: FromFloat64 function has been removed, because it is impossible to know the desired precision of the result.
//               The user should convert float64 to string, with the desired precision, and pass it to FromString.
//
//func (context *Context) FromFloat64(value float64) (r Quad) {
//	var result C.Ret_decQuad_t
//	assert_sane(context)
//
//	result = C.mdq_from_double(C.double(value), context.set)
//
//	context.set = result.set
//	return Quad{val: result.val}
//}

/************************************************************************/
/*                                                                      */
/*                      conversion to string                            */
/*                                                                      */
/************************************************************************/

const POOL_BUFF_CAPACITY = 50 // capacity of []byte buffer generated by the pool of buffers

// pool is a pool of byte slice, used by AppendQuad and String.
//
// note:
//    DECQUAD_String     = 43         sign, 34 digits, decimal point, E+xxxx, terminal \0   gives 43
//    DECQUAD_Pmax       = 34
//    POOL_BUFF_CAPACITY = 50         just to be sure, it is largely enough
//
// The pool must return []byte with capacity being at least the largest of DECQUAD_String and DECQUAD_Pmax. We Prefer a capacity of POOL_BUFF_CAPACITY to be sure.
//
var pool = sync.Pool{
	New: func() interface{} {
		//fmt.Println("---   POOL")
		return make([]byte, POOL_BUFF_CAPACITY) // POOL_BUFF_CAPACITY is larger than DECQUAD_String and DECQUAD_Pmax. This size is ok for AppendQuad and String methods.
	},
}

// QuadToString returns the string representation of a Quad number.
// It calls the C function QuadToString of the original decNumber package.
//
//       This function uses exponential notation quite often.
//       E.g. 0.0000001 returns "1E-7", which is often not what we want.
//
//       It is better to use the method AppendQuad() or String(), which don't use exponential notation for a wider range.
//       AppendQuad() and String() write a number without exp notation if it can be displayed with at most 34 digits, and an optional fractional point.
//
func (a Quad) QuadToString() string {
	var (
		ret_str   C.Ret_str
		str_slice []byte // capacity must be exactly DECQUAD_String
		s         string
	)

	ret_str = C.mdq_to_QuadToString(a.val) // may use exponent notation

	str_slice = pool.Get().([]byte)[:DECQUAD_String]
	defer pool.Put(str_slice)

	for i := 0; i < int(ret_str.length); i++ {
		str_slice[i] = byte(ret_str.s[i])
	}

	s = string(str_slice[:ret_str.length])

	return s
}

// AppendQuad appends string representation of Quad into byte slice.
// AppendQuad and String are best to display Quad, as exponent notation is used less often than with QuadToString.
//
//       AppendQuad() writes a number without exp notation if it can be displayed with at most 34 digits, and an optional fractional point.
//       Else, falls back on QuadToString(), which will use exponential notation.
//
// See also method String(), which calls AppendQuad internally.
//
func AppendQuad(dst []byte, a Quad) []byte {
	var (
		ret_str   C.Ret_str
		str_slice []byte // length must be exactly DECQUAD_String

		ret               C.Ret_BCD
		d                 byte
		skip_leading_zero bool = true
		inf_nan           uint32
		exp               int32
		sign              uint32
		BCD_slice         []byte // length must be exactly DECQUAD_Pmax

		buff [DECQUAD_String]byte // enough for      sign    optional "0."    34 digits
	)

	// fill BCD array

	ret = C.mdq_to_BCD(a.val) // sign will be 1 for negative and non-zero number, else, 0. If Inf or Nan, returns an error.

	BCD_slice = pool.Get().([]byte)[:DECQUAD_Pmax]
	defer pool.Put(BCD_slice)

	for i := 0; i < DECQUAD_Pmax; i++ {
		BCD_slice[i] = byte(ret.BCD[i])
	}
	inf_nan = uint32(ret.inf_nan)
	exp = int32(ret.exp)
	sign = uint32(ret.sign)

	// if Quad value is not in 34 digits range, or Inf or Nan, we want our function to output the number, or Infinity, or NaN. Falls back on QuadToString.

	if exp > 0 || exp < -DECQUAD_Pmax || inf_nan != 0 {
		ret_str = C.mdq_to_QuadToString(a.val) // may use exponent notation

		str_slice = pool.Get().([]byte)[:DECQUAD_String]
		defer pool.Put(str_slice)

		for i := 0; i < int(ret_str.length); i++ {
			str_slice[i] = byte(ret_str.s[i])
		}

		dst = append(dst, str_slice[:ret_str.length]...) // write buff into destination and return

		return dst
	}

	// write string. Here, the number is not Inf nor Nan.

	i := 0

	integral_part_length := len(BCD_slice) + int(exp) // here, exp is [-DECQUAD_Pmax ... 0]

	BCD_integral_part := BCD_slice[:integral_part_length]
	BCD_fractional_part := BCD_slice[integral_part_length:]

	for _, d = range BCD_integral_part { // ==== write integral part ====
		if skip_leading_zero && d == 0 {
			continue
		} else {
			skip_leading_zero = false
		}
		buff[i] = '0' + d
		i++
	}

	if i == 0 { // write '0' if no digit written for integral part
		buff[i] = '0'
		i++
	}

	if sign != 0 {
		dst = append(dst, '-') // write '-' sign if any into destination
	}

	dst = append(dst, buff[:i]...) // write integral part into destination

	if exp == 0 { // if no fractional part, just return
		return dst
	}

	dst = append(dst, '.') // ==== write fractional part ====

	i = 0
	for _, d = range BCD_fractional_part {
		buff[i] = '0' + d
		i++
	}

	dst = append(dst, buff[:i]...) // write fractional part into destination

	return dst
}

// String is the preferred way to display a decQuad number.
// It calls AppendQuad internally.
//
func (a Quad) String() string {
	var buffer []byte

	buffer = pool.Get().([]byte)[:0] // capacity is enough to receive result of C.mdq_to_QuadToString(), and also big enough to receive [sign] + [DECQUAD_Pmax digits] + [fractional dot]
	defer pool.Put(buffer)

	ss := AppendQuad(buffer[:0], a)

	return string(ss)
}

/************************************************************************/
/*                                                                      */
/*                      conversion to number                            */
/*                                                                      */
/************************************************************************/

// ToInt32 returns the int32 value from a.
// The rounding passed as argument is used, instead of the rounding mode of context which is ignored.
//
func (context *Context) ToInt32(a Quad, round Round_mode_t) int32 {
	var result C.Ret_int32_t
	assert_sane(context)

	result = C.mdq_to_int32(a.val, context.set, C.int(round))

	context.set = result.set
	return int32(result.val)
}

// ToInt64 returns the int64 value from a.
// The rounding passed as argument is used, instead of the rounding mode of context which is ignored.
//
func (context *Context) ToInt64(a Quad, round Round_mode_t) int64 {
	var result C.Ret_int64_t
	assert_sane(context)

	result = C.mdq_to_int64(a.val, context.set, C.int(round))

	context.set = result.set
	return int64(result.val)
}

// ToFloat64 returns the float64 value from a.
//
func (context *Context) ToFloat64(a Quad) float64 {
	var (
		err error
		s   string
		val float64
	)
	assert_sane(context)

	if a.IsNaN() { // because strconv.ParseFloat doesn't parse signaling sNaN
		return math.NaN()
	}

	s = a.String()

	if val, err = strconv.ParseFloat(s, 64); err != nil {
		context.SetStatus(Flag_Conversion_syntax)
		return math.NaN()
	}

	return val
}
