package decnum

/*

#include "mydecquad.h"
*/
import "C"

import (
	"errors"
	"fmt"
	"strings"
)

// Note: in the comment block for cgo above, the path for LDFLAGS must be an "absolute" path.
//       If it is a relative path, it seems that it is relative to the current directory when "go build" is run.
//       As we want to run "go build decnum" from any location, the path must be absolute.
//       See   Issue 5428 in May 2013: cmd/ld: relative #cgo LDFLAGS -L does not work
//             This problem is still not resolved in May 2014.

// DecQuad is just a struct with a C.decQuad value, which is an array of 16 bytes.
type DecQuad struct {
	val C.decQuad
}

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

/************************************************************************/
/*                                                                      */
/*                            init function                             */
/*                                                                      */
/************************************************************************/

var DecQuad_module_MACROS string // macros defined by the C decQuad module

func init() {

	C.mdq_init()

	DecQuad_module_MACROS = fmt.Sprintf("decQuad module: DECDPUN %d, DECSUBSET %d, DECEXTFLAG %d. Constants DECQUAD_Pmax %d, DECQUAD_String %d DECQUAD_Bytes %d.", C.DECDPUN, C.DECSUBSET, C.DECEXTFLAG, C.DECQUAD_Pmax, C.DECQUAD_String, C.DECQUAD_Bytes)

	if DECQUAD_Bytes != 16 { // DECQUAD_Bytes MUST NOT BE > 16, because Append_compressed_bytes() will silently fail if it is not the case
		panic("DECQUAD_Bytes != 16")
	}

	assert(C.DECSUBSET == 0) // because else, we should define Flag_Lost_digits as status flag

	// set global variables g_zero to 0, and g_nan to Nan

	g_zero = zero_for_init()

	g_nan = nan_for_init()
}

/************************************************************************/
/*                                                                      */
/*                          utility functions                           */
/*                                                                      */
/************************************************************************/

func assert(val bool) {
	if val == false {
		panic("assertion failed")
	}
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

type Status_t uint32

const (
	Flag_Conversion_syntax    Status_t = C.DEC_Conversion_syntax
	Flag_Division_by_zero     Status_t = C.DEC_Division_by_zero
	Flag_Division_impossible  Status_t = C.DEC_Division_impossible
	Flag_Division_undefined   Status_t = C.DEC_Division_undefined
	Flag_Insufficient_storage Status_t = C.DEC_Insufficient_storage
	Flag_Inexact              Status_t = C.DEC_Inexact
	Flag_Invalid_context      Status_t = C.DEC_Invalid_context
	Flag_Invalid_operation    Status_t = C.DEC_Invalid_operation
	Flag_Overflow             Status_t = C.DEC_Overflow
	Flag_Clamped              Status_t = C.DEC_Clamped
	Flag_Rounded              Status_t = C.DEC_Rounded
	Flag_Subnormal            Status_t = C.DEC_Subnormal
	Flag_Underflow            Status_t = C.DEC_Underflow // e.g. 1e-6000/1e1000

	//Flag_Lost_digits          Status_t = C.DEC_Lost_digits // exists only if DECSUBSET is set, which is not the case by default
)

const ErrorMask Status_t = C.DEC_Errors // ErrorMask is a bitmask of many of the above flags, ORed together. After a series of operations, if status & decnum.Errors != 0, an error has occured, e.g. division by 0.

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

// Context contains the rounding mode, and a status field that records exceptional conditions, some of which are considered as error, e.g. division by 0, underlow for operations like 1e-6000/1e1000, overflow, etc.
// For decQuad usage, only these two fields are used.
//
// When an error occurs during an operation, the result will probably be NaN or infinite, or a infinitesimal number if underflow.
//
type Context struct {
	set C.decContext
}

type Context_kind_t uint32

const (
	DEFAULT_DECQUAD Context_kind_t = C.DEC_INIT_DECQUAD // default Context settings for decQuad operations
)

// Init is used to initialize a context with default value for rounding mode, and status field is cleared.
//
func (context *Context) Init(kind Context_kind_t) {

	context.set = C.mdq_context_default(context.set, C.uint32_t(kind))
}

// Rounding returns the rounding mode of the context.
//
func (context *Context) Rounding() Round_mode_t {

	return Round_mode_t(C.mdq_context_get_rounding(context.set))
}

// SetRounding sets the rounding mode of the context.
//
func (context *Context) SetRounding(rounding Round_mode_t) {

	context.set = C.mdq_context_set_rounding(context.set, C.int(rounding))
}

// Status returns the status of the context.
//
// After a series of operations, the status contains the accumulated errors or informational flags that occurred during all the operations.
// Beware: the status can contain informational flags, like Flag_Inexact, which is not an error.
// To find the real errors, you must discard the non-error bits of the status as follows:
//      ctx.Status() & decnum.Errors
//
//
// It is easier to use the context.ErrorMask method to check for errors.
//
func (context *Context) Status() Status_t {

	return Status_t(C.mdq_context_get_status(context.set))
}

// SetStatus sets a status bit in the status of the context.
//
// Normally, only library modules use this function. Applications have no reason to set status bits.
//
func (context *Context) SetStatus(flag Status_t) {

	context.set = C.mdq_context_set_status(context.set, C.uint32_t(flag))
}

// ResetStatus clears all bits of the status field of the context.
// You can continue to use this context for a new series of operations.
//
func (context *Context) ResetStatus() {

	context.set = C.mdq_context_zero_status(context.set)
}

// Error checks if status contains a flag that should be considered as an error.
// In this case, the resut of the operations contains Nan or Infinite, or a infinitesimal number if Underflow-
//
func (context *Context) Error() error {
	var status Status_t

	status = context.Status()

	if status&ErrorMask != 0 {
		return fmt.Errorf("decnum error: %s", status.String())
	}

	return nil
}

/************************************************************************/
/*                                                                      */
/*                      arithmetic operations                           */
/*                                                                      */
/************************************************************************/

const (
	CMP_LESS    Cmp_t = C.CMP_LESS    // 1
	CMP_EQUAL   Cmp_t = C.CMP_EQUAL   // 2
	CMP_GREATER Cmp_t = C.CMP_GREATER // 4
	CMP_NAN     Cmp_t = C.CMP_NAN     // 8
)

type Cmp_t uint32 // result of Compare

var (
	g_zero DecQuad // a constant DecQuad with value 0
	g_nan  DecQuad // a constant DecQuad with value Nan
)

// used only by init() to initialize the global variable g_zero.
//
func zero_for_init() (r DecQuad) {
	var val C.decQuad

	val = C.mdq_zero()

	return DecQuad{val: val}
}

// used only by init() to initialize the global variable g_Nan.
//
func nan_for_init() (r DecQuad) {
	var val C.decQuad

	val = C.mdq_nan()

	return DecQuad{val: val}
}

// return a 0 DecQuad value.
//
//     r = Zero()  // assign 0 to the DecQuad r
//
func Zero() (r DecQuad) {

	return g_zero
}

// return a Nan DecQuad value.
//
//     r = Nan()  // assign Nan to the DecQuad r
//
func Nan() (r DecQuad) {

	return g_nan
}

// Minus returns -a.
//
func (context *Context) Minus(a DecQuad) (r DecQuad) {
	var result C.Ret_decQuad_t

	result = C.mdq_minus(a.val, context.set)

	context.set = result.set

	return DecQuad{val: result.val}
}

// Add returns a + b.
//
func (context *Context) Add(a DecQuad, b DecQuad) (r DecQuad) {
	var result C.Ret_decQuad_t

	result = C.mdq_add(a.val, b.val, context.set)

	context.set = result.set

	return DecQuad{val: result.val}
}

// Subtract returns a - b.
//
func (context *Context) Subtract(a DecQuad, b DecQuad) (r DecQuad) {
	var result C.Ret_decQuad_t

	result = C.mdq_subtract(a.val, b.val, context.set)

	context.set = result.set

	return DecQuad{val: result.val}
}

// Multiply returns a * b.
//
func (context *Context) Multiply(a DecQuad, b DecQuad) (r DecQuad) {
	var result C.Ret_decQuad_t

	result = C.mdq_multiply(a.val, b.val, context.set)

	context.set = result.set

	return DecQuad{val: result.val}
}

// Divide returns a/b.
//
func (context *Context) Divide(a DecQuad, b DecQuad) (r DecQuad) {
	var result C.Ret_decQuad_t

	result = C.mdq_divide(a.val, b.val, context.set)

	context.set = result.set

	return DecQuad{val: result.val}
}

// DivideInteger returns the integral part of a/b.
//
func (context *Context) DivideInteger(a DecQuad, b DecQuad) (r DecQuad) {
	var result C.Ret_decQuad_t

	result = C.mdq_divide_integer(a.val, b.val, context.set)

	context.set = result.set

	return DecQuad{val: result.val}
}

// Remainder returns the modulo of a and b.
//
func (context *Context) Remainder(a DecQuad, b DecQuad) (r DecQuad) {
	var result C.Ret_decQuad_t

	result = C.mdq_remainder(a.val, b.val, context.set)

	context.set = result.set

	return DecQuad{val: result.val}
}

// Abs returns the absolute value of a.
//
func (context *Context) Abs(a DecQuad) (r DecQuad) {
	var result C.Ret_decQuad_t

	result = C.mdq_abs(a.val, context.set)

	context.set = result.set

	return DecQuad{val: result.val}
}

// ToIntegral returns the value of a rounded to an integral value.
//
func (context *Context) ToIntegral(a DecQuad, round Round_mode_t) (r DecQuad) {
	var result C.Ret_decQuad_t

	result = C.mdq_to_integral(a.val, context.set, C.int(round))

	context.set = result.set

	return DecQuad{val: result.val}
}

// Quantize rounds a to the same pattern as b.
// b is just a model, its value is not used.
// The result is the value of a, but with the same exponent as the pattern b.
// The rounding of the context is used.
//
// You can use this function with the proper rounding to round (e.g. set context rounding mode to ROUND_HALF_EVEN) or truncate (ROUND_DOWN) 'a'.
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
func (context *Context) Quantize(a DecQuad, b DecQuad) (r DecQuad) {
	var result C.Ret_decQuad_t

	result = C.mdq_quantize(a.val, b.val, context.set)

	context.set = result.set

	return DecQuad{val: result.val}
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
func (context *Context) Compare(a DecQuad, b DecQuad) Cmp_t {
	var result C.Ret_uint32_t

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
func (context *Context) Cmp(a DecQuad, b DecQuad, comp_mask Cmp_t) bool {
	var result C.Ret_uint32_t

	result = C.mdq_compare(a.val, b.val, context.set)

	context.set = result.set

	if Cmp_t(result.val)&comp_mask != 0 {
		return true
	}

	return false
}

// IsFinite returns true if a is not Infinite, nor Nan.
//
func (context *Context) IsFinite(a DecQuad) bool {

	if C.mdq_is_finite(a.val) != 0 {
		return true
	}

	return false
}

// IsInteger returns true if a is finite and has exponent=0.
//
func (context *Context) IsInteger(a DecQuad) bool {

	if C.mdq_is_integer(a.val) != 0 {
		return true
	}

	return false
}

// IsInfinite returns true if a is Infinite.
//
func (context *Context) IsInfinite(a DecQuad) bool {

	if C.mdq_is_infinite(a.val) != 0 {
		return true
	}

	return false
}

// IsNan returns true if a is Nan.
//
func (context *Context) IsNan(a DecQuad) bool {

	if C.mdq_is_nan(a.val) != 0 {
		return true
	}

	return false
}

// IsNegative returns true if a is Nan.
//
func (context *Context) IsNegative(a DecQuad) bool {

	if C.mdq_is_negative(a.val) != 0 {
		return true
	}

	return false
}

// IsZero returns true if a == 0.
//
func (context *Context) IsZero(a DecQuad) bool {

	if C.mdq_is_zero(a.val) != 0 {
		return true
	}

	return false
}

// Max returns the larger of a and b.
// If either a or b is NaN then the other argument is the result.
//
func (context *Context) Max(a DecQuad, b DecQuad) (r DecQuad) {
	var result C.Ret_decQuad_t

	result = C.mdq_max(a.val, b.val, context.set)

	context.set = result.set

	return DecQuad{val: result.val}
}

// Min returns the smaller of a and b.
// If either a or b is NaN then the other argument is the result.
//
func (context *Context) Min(a DecQuad, b DecQuad) (r DecQuad) {
	var result C.Ret_decQuad_t

	result = C.mdq_min(a.val, b.val, context.set)

	context.set = result.set

	return DecQuad{val: result.val}
}

/************************************************************************/
/*                                                                      */
/*                      conversion to string                            */
/*                                                                      */
/************************************************************************/

// AppendQuad appends string representation of decQuad into byte slice.
// This representation is the best to display decQuad, because it shows all numbers having exponent between 0 and -34 (DECQUAD_Pmax), that is, all 34 significant digits, without using exponent notation.
//
// All digits of the coefficient are displayed, e.g. 12344567890.123456789000000000000000
// If the number must have an exponent because it is too large or too small, we can have    1.4E+201     0E-6176    Infinity
// If the number has an internal positive exponent, like 33e4, it will be displayed as 3.3E+5, though. But 33.1234e4 is displayed as 331234.
//
// Method String() calls AppendQuad internally.
//
func AppendQuad(dst []byte, a *DecQuad) []byte {
	var (
		ret_str   C.Ret_str
		str_slice []byte = make([]byte, DECQUAD_String)

		ret               C.Ret_BCD
		d                 byte
		skip_leading_zero bool = true
		inf_nan           C.uint32_t
		exp               int32
		sign              uint32
		BCD_slice         []byte = make([]byte, DECQUAD_Pmax)

		buff [DECQUAD_String]byte // array size is max of DECQUAD_String and DECQUAD_Pmax. DECQUAD_String is larger.
	)

	// fill BCD array

	ret = C.mdq_to_BCD(a.val) // sign will be 1 for negative and non-zero number, else, 0. If Inf or Nan, returns an error.
	for i := 0; i < DECQUAD_Pmax; i++ {
		BCD_slice[i] = byte(ret.BCD[i])
	}
	exp = int32(ret.exp)
	sign = uint32(ret.sign)
	inf_nan = ret.inf_nan

	if exp > 0 || exp < -DECQUAD_Pmax || inf_nan != 0 { // if decQuad value is not in 34 digits range, or Inf or Nan, we want our function to output the number, or Infinity, or NaN.
		ret_str = C.mdq_to_QuadToString(a.val) // may use exponent notation
		for i := 0; i < int(ret_str.length); i++ {
			str_slice[i] = byte(ret_str.s[i])
		}

		dst = append(dst, str_slice[:ret_str.length]...) // write buff into destination and return
		return dst
	}

	// write string. Here, the number is not Inf nor Nan.

	i := 0

	integral_part_length := len(BCD_slice) + int(exp) // here, exp is negative

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

	if i == 0 { // write '0' if integral part is 0
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
func (a DecQuad) String() string {
	var buffer [DECQUAD_String]byte // to avoid reallocation, this capacity is needed to receive result of C.mdq_to_mallocated_QuadToString(), and also big enough to receive [sign] + [DECQUAD_Pmax digits] + [fractional dot]

	ss := AppendQuad(buffer[:0], &a)

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
func (context *Context) ToInt32(a DecQuad, round Round_mode_t) int32 {
	var result C.Ret_int32_t

	result = C.mdq_to_int32(a.val, context.set, C.int(round))

	context.set = result.set

	return int32(result.val)
}

// ToInt64 returns the int64 value from a.
// The rounding passed as argument is used, instead of the rounding mode of context which is ignored.
//
func (context *Context) ToInt64(a DecQuad, round Round_mode_t) int64 {
	var result C.Ret_int64_t

	result = C.mdq_to_int64(a.val, context.set, C.int(round))

	context.set = result.set

	return int64(result.val)
}

/************************************************************************/
/*                                                                      */
/*                   conversion from string and numbers                 */
/*                                                                      */
/************************************************************************/

const MAX_STRING_SIZE = C.MAX_STRING_SIZE

// FromString returns a DecQuad from a string.
//
func (context *Context) FromString(s string) (r DecQuad) {
	var (
		i        int
		strarray C.Strarray_t
		result   C.Ret_decQuad_t
	)

	s = strings.TrimSpace(s)

	if len(s) > MAX_STRING_SIZE {
		context.SetStatus(Flag_Conversion_syntax)
		r = Nan()
		return r
	}

	for i = 0; i < len(s); i++ {
		strarray.arr[i] = C.char(s[i])
	}

	strarray.arr[i] = 0 // terminating 0

	result = C.mdq_from_string(strarray, context.set)

	context.set = result.set

	return DecQuad{val: result.val}
}

// FromInt32 returns a DecQuad from a int32 value.
//
func (context *Context) FromInt32(value int32) (r DecQuad) {

	C.decQuadFromInt32(&r.val, C.int32_t(value))

	return r
}

// FromInt64 returns a DecQuad from a int64 value.
//
func (context *Context) FromInt64(value int64) (r DecQuad) {
	var result C.Ret_decQuad_t

	result = C.mdq_from_int64(C.int64_t(value), context.set)

	context.set = result.set

	return DecQuad{val: result.val}
}
