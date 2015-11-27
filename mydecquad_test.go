package decnum

import (
	"log"
	"strconv"
	"testing"
)

var (
	maxquad = "9.999999999999999999999999999999999E+6144"
	minquad = "-9.999999999999999999999999999999999E+6144"

	smallquad  = "9.999999999999999999999999999999999E-6143"
	nsmallquad = "-9.999999999999999999999999999999999E-6143"
)

// converts string to Quad or aborts.
//
func must_quad(s string) Quad {
	var q Quad

	q, _ = FromString(s)

	if err := q.Error(); err != nil {
		log.Fatalf("must_quad(\"%s\") failed: %s", s, err)
	}

	return q
}

// converts string to RoundingMode or aborts.
//
func must_rounding(s string) RoundingMode {

	switch s {
	case "RoundCeiling":
		return RoundCeiling
	case "RoundDown":
		return RoundDown
	case "RoundFloor":
		return RoundFloor
	case "RoundHalfDown":
		return RoundHalfDown
	case "RoundHalfEven":
		return RoundHalfEven
	case "RoundHalfUp":
		return RoundHalfUp
	case "RoundUp":
		return RoundUp
	case "Round05Up":
		return Round05Up
	default:
		log.Fatalf("must_rounding(\"%s\") failed: unknown rounding mode", s)
	}

	panic("impossible")
}

// converts string to CmpFlag or aborts.
//
func must_cmp(s string) CmpFlag {

	switch s {
	case "CmpLess":
		return CmpLess
	case "CmpEqual":
		return CmpEqual
	case "CmpGreater":
		return CmpGreater
	case "CmpNaN":
		return CmpNaN
	default:
		log.Fatalf("must_cmp(\"%s\") failed: unknown CmpFlag", s)
	}

	panic("impossible")
}

// converts string to int32 or aborts.
//
func must_int32(s string) int32 {
	var (
		err error
		i   int64
	)

	if i, err = strconv.ParseInt(s, 10, 32); err != nil {
		log.Fatalf("must_int32(\"%s\") failed: %s", s, err)
	}

	return int32(i)
}

// converts string to int64 or aborts.
//
func must_int64(s string) int64 {
	var (
		err error
		i   int64
	)

	if i, err = strconv.ParseInt(s, 10, 64); err != nil {
		log.Fatalf("must_int64(\"%s\") failed: %s", s, err)
	}

	return i
}

// converts string to int or aborts.
//
func must_int(s string) int {
	var (
		err error
		i   int64
	)

	if i, err = strconv.ParseInt(s, 10, 64); err != nil {
		log.Fatalf("must_int64(\"%s\") failed: %s", s, err)
	}

	return int(i)
}

// converts string to float64 or aborts.
//
func must_float64(s string) float64 {
	var (
		err error
		f   float64
	)

	if f, err = strconv.ParseFloat(s, 64); err != nil {
		log.Fatalf("must_float64(\"%s\") failed: %s", s, err)
	}

	return f
}

func bool2string(b bool) string {

	if b {
		return "true"
	}

	return "false"
}

func Test_simple_functions(t *testing.T) {
	var (
		a Quad
		b Quad
	)

	// Zero

	a = Zero()

	if a.String() != "0" || a.Error() != nil {
		t.Fatal("a = Zero() failed")
	}

	// One

	a = One()

	if a.String() != "1" || a.Error() != nil {
		t.Fatal("a = One() failed")
	}

	// NaN

	a = NaN()

	if a.String() != "NaN" || a.Error() != nil {
		t.Fatal("a = Nan() failed")
	}

	// assignment

	a, _ = FromString("123.45")

	b = a

	if b.String() != "123.45" || b.Error() != nil {
		t.Fatal("b = a failed")
	}

	// Copy

	a, _ = FromString("567890.245")

	b = Copy(a)

	if b.String() != "567890.245" || b.Error() != nil {
		t.Fatal("b = Copy(a) failed")
	}
}

func Test_status_set_clear(t *testing.T) {

	a := One()

	r := a.SetStatusFlags(DivisionByZero | Underflow | InvalidOperation)
	if r.Status() != DivisionByZero|Underflow|InvalidOperation {
		t.Fatalf("incorrect status: %s", r.Status())
	}

	r = r.ClearStatusFlags(DivisionByZero | InvalidOperation)
	if r.Status() != Underflow {
		t.Fatalf("incorrect status: %s", r.Status())
	}

	r = r.ClearStatus()
	if r.Status() != 0 {
		t.Fatalf("incorrect status: %s", r.Status())
	}
}

func Test_status_propagation(t *testing.T) {

	a := One()
	a = a.SetStatusFlags(DivisionByZero)

	b := One()
	b = b.SetStatusFlags(Underflow)

	r := a.Neg()
	if r.Status() != DivisionByZero {
		t.Fatalf("incorrect status: %s", r.Status())
	}

	r = a.Add(b)
	if r.Status() != DivisionByZero|Underflow {
		t.Fatalf("incorrect status: %s", r.Status())
	}

	r = a.Sub(b)
	if r.Status() != DivisionByZero|Underflow {
		t.Fatalf("incorrect status: %s", r.Status())
	}

	r = a.Mul(b)
	if r.Status() != DivisionByZero|Underflow {
		t.Fatalf("incorrect status: %s", r.Status())
	}

	r = a.Div(b)
	if r.Status() != DivisionByZero|Underflow {
		t.Fatalf("incorrect status: %s", r.Status())
	}

	r = a.DivInt(b)
	if r.Status() != DivisionByZero|Underflow {
		t.Fatalf("incorrect status: %s", r.Status())
	}

	r = a.Mod(b)
	if r.Status() != DivisionByZero|Underflow {
		t.Fatalf("incorrect status: %s", r.Status())
	}

	r = Max(a, b)
	if r.Status() != DivisionByZero|Underflow {
		t.Fatalf("incorrect status: %s", r.Status())
	}

	r = Min(a, b)
	if r.Status() != DivisionByZero|Underflow {
		t.Fatalf("incorrect status: %s", r.Status())
	}

	r = a.ToIntegral(RoundHalfEven)
	if r.Status() != DivisionByZero {
		t.Fatalf("incorrect status: %s", r.Status())
	}

	r = a.Quantize(b, RoundHalfEven)
	if r.Status() != DivisionByZero|Underflow {
		t.Fatalf("incorrect status: %s", r.Status())
	}

	r = a.Abs()
	if r.Status() != DivisionByZero {
		t.Fatalf("incorrect status: %s", r.Status())
	}

	r = a.Add(b).Sub(b).Mul(b).Div(b).DivInt(b).Mod(b).ToIntegral(RoundHalfEven).Quantize(b, RoundHalfEven).Abs()
	if r.Status() != DivisionByZero|Underflow {
		t.Fatalf("incorrect status: %s", r.Status())
	}

	r = b.Add(a).Sub(b).Mul(b).Div(b).DivInt(b).Mod(b).ToIntegral(RoundHalfEven).Quantize(b, RoundHalfEven).Abs()
	if r.Status() != DivisionByZero|Underflow {
		t.Fatalf("incorrect status: %s", r.Status())
	}

	r = a.Round(2)
	if r.Status() != DivisionByZero {
		t.Fatalf("incorrect status: %s", r.Status())
	}

	r = a.Truncate(2)
	if r.Status() != DivisionByZero {
		t.Fatalf("incorrect status: %s", r.Status())
	}

}

func Test_operations(t *testing.T) {

	type Operation_t string

	const (
		T_NEG          Operation_t = "Neg"
		T_ADD          Operation_t = "Add"
		T_SUB          Operation_t = "Sub"
		T_MUL          Operation_t = "Mul"
		T_DIV          Operation_t = "Div"
		T_DIVINT       Operation_t = "DivInt"
		T_MOD          Operation_t = "Mod"
		T_MAX          Operation_t = "Max"
		T_MIN          Operation_t = "Min"
		T_TOINTEGRAL   Operation_t = "ToIntegral"
		T_QUANTIZE     Operation_t = "Quantize"
		T_ABS          Operation_t = "Abs"
		T_ISFINITE     Operation_t = "IsFinite"
		T_ISINTEGER    Operation_t = "IsInteger"
		T_ISINFINITE   Operation_t = "IsInfinite"
		T_ISNAN        Operation_t = "IsNan"
		T_ISPOSITIVE   Operation_t = "IsPositive"
		T_ISZERO       Operation_t = "IsZero"
		T_ISNEGATIVE   Operation_t = "IsNegative"
		T_GREATER      Operation_t = "Greater"
		T_GREATEREQUAL Operation_t = "GreaterEqual"
		T_EQUAL        Operation_t = "Equal"
		T_LESSEQUAL    Operation_t = "LessEqual"
		T_LESS         Operation_t = "Less"
		T_FROMSTRING   Operation_t = "FromString"
		T_FROMINT32    Operation_t = "FromInt32"
		T_FROMINT64    Operation_t = "FromInt64"
		T_QUADTOSTRING Operation_t = "QuadToString"
		T_STRING       Operation_t = "String"
		T_TOINT32      Operation_t = "ToInt32"
		T_TOINT64      Operation_t = "ToInt64"
		T_TOFLOAT64    Operation_t = "ToFloat64"
		T_ROUND        Operation_t = "Round"
		T_TRUNCATE     Operation_t = "Truncate"
	)

	// A decimal number can also represents three special values: Infinity, NaN, and signaling NaN.
	//
	//    Infinity and -Infinity, or Inf and -Inf, represent a value infinitely large.
	//
	//    NaN or qNaN, which means "Not a Number", represents an undefined result, when an arithmetic operation has failed. E.g. FromString("hello")
	//                 NaN propagates to all subsequent operations, because if NaN is passed as argument, the result, will be NaN.
	//                 These NaN are called "quiet NaN", because they don't set exceptional condition flag in status when passed as argument to an operation.
	//
	//    sNaN, or "signaling NaN", are created by FromString("sNaN"). When passed as argument to an operation, the result will be NaN, like with quiet NaN.
	//                 But they will set (==signal) an exceptional condition flag in status, "Invalid_operation".
	//                 Signaling NaN propagate to subsequent operation as ordinary NaN (quiet NaN), and not as "signaling NaN".
	//
	// Note that both NaN and sNaN can take an integer payload, e.g. NaN123, created by FromString("NaN123"), and it is up to you to give it a significance.
	// sNaN and payload are not used often, and most probably, you won't use them. Nan's payload propagates to subsequent operations.
	//

	var samples = []struct {
		operation             Operation_t // operation to test
		a                     string      // first argument of operation to test. Type depends on operation.
		b                     string      // 2nd argument of operation to test. Type depends on operation.
		expected_result       string      // expected result of operation
		expected_error_status Status      // error status flags expected after an operation
	}{
		{T_NEG, "sNaN", "", "NaN", InvalidOperation},       // Invalid_operation      because of sNan (signaling NaN)
		{T_NEG, "sNaN123", "", "NaN123", InvalidOperation}, // Invalid_operation      because of sNan (signaling NaN)
		{T_NEG, "NaN", "", "NaN", 0},
		{T_NEG, "-NaN", "", "-NaN", 0},
		{T_NEG, "NaN123", "", "NaN123", 0},
		{T_NEG, "-NaN123", "", "-NaN123", 0},
		{T_NEG, "Inf", "", "-Infinity", 0},
		{T_NEG, "-Inf", "", "Infinity", 0},
		{T_NEG, "0.00", "", "0.00", 0},
		{T_NEG, "0.00E50", "", "0E+48", 0},
		{T_NEG, "-13256748.9879878", "", "13256748.9879878", 0},
		{T_NEG, "13256748.9879878", "", "-13256748.9879878", 0},
		{T_NEG, "-13256748.9879878e456", "", "1.32567489879878E+463", 0},
		{T_NEG, maxquad, "", minquad, 0},
		{T_NEG, minquad, "", maxquad, 0},
		{T_NEG, smallquad, "", nsmallquad, 0},
		{T_NEG, nsmallquad, "", smallquad, 0},

		{T_ADD, "1", "sNaN", "NaN", InvalidOperation},       // Invalid_operation      because of sNan (signaling NaN)
		{T_ADD, "1", "sNaN123", "NaN123", InvalidOperation}, // Invalid_operation      because of sNan (signaling NaN)
		{T_ADD, "1", "NaN", "NaN", 0},
		{T_ADD, "1", "NaN123", "NaN123", 0},
		{T_ADD, "NaN", "NaN", "NaN", 0},
		{T_ADD, "NaN", "123", "NaN", 0},
		{T_ADD, "NaN", "Inf", "NaN", 0},
		{T_ADD, "123", "NaN", "NaN", 0},
		{T_ADD, "Inf", "NaN", "NaN", 0},
		{T_ADD, "Inf", "Inf", "Infinity", 0},
		{T_ADD, "-Inf", "-Inf", "-Infinity", 0},
		{T_ADD, "Inf", "-Inf", "NaN", InvalidOperation}, // Invalid_operation
		{T_ADD, "123", "200", "323", 0},
		{T_ADD, "123.1230", "200", "323.1230", 0},
		{T_ADD, "123.1230", "Inf", "Infinity", 0},
		{T_ADD, "-123456789012345678901234567890.1234", "123456789012345678901234567890.1234", "0.0000", 0},
		{T_ADD, "-123456789012345678901234567890.1234e200", "123456789012345678901234567890.1234", "-1.234567890123456789012345678901234E+229", 0},
		{T_ADD, "-123456789012345678901234567890.1234e200", "123456789012345678901234567890.1234e206", "1.234566655555566665555556666555555E+235", 0},
		{T_ADD, "9999999999999999999999999999999999", "0", "9999999999999999999999999999999999", 0},
		{T_ADD, "9999999999999999999999999999999999", "1", "1.000000000000000000000000000000000E+34", 0},
		{T_ADD, "9999999999999999999999999999999999", "2", "1.000000000000000000000000000000000E+34", 0},
		{T_ADD, maxquad, "0", maxquad, 0},
		{T_ADD, maxquad, "1", maxquad, 0},
		{T_ADD, maxquad, "1e6111", "Infinity", Overflow}, // Overflow
		{T_ADD, maxquad, "Inf", "Infinity", 0},
		{T_ADD, maxquad, "-Inf", "-Infinity", 0},
		{T_ADD, "142566.645373", "647833330000004.7367", "647833330142571.382073", 0},
		{T_ADD, "1425658446.645373", "-647833330000004.7367", "-647831904341558.091327", 0},
		{T_ADD, smallquad, "1", "1.000000000000000000000000000000000", 0},

		{T_SUB, "1", "sNaN", "NaN", InvalidOperation},       // Invalid_operation      because of sNan (signaling NaN)
		{T_SUB, "1", "sNaN456", "NaN456", InvalidOperation}, // Invalid_operation      because of sNan (signaling NaN)
		{T_SUB, "NaN", "NaN", "NaN", 0},
		{T_SUB, "NaN", "123", "NaN", 0},
		{T_SUB, "NaN456", "123", "NaN456", 0},
		{T_SUB, "NaN", "Inf", "NaN", 0},
		{T_SUB, "123", "NaN", "NaN", 0},
		{T_SUB, "Inf", "NaN", "NaN", 0},
		{T_SUB, "Inf", "Inf", "NaN", InvalidOperation},   // Invalid_operation
		{T_SUB, "-Inf", "-Inf", "NaN", InvalidOperation}, // Invalid_operation
		{T_SUB, "Inf", "-Inf", "Infinity", 0},
		{T_SUB, "123", "200", "-77", 0},
		{T_SUB, "123.1230", "200", "-76.8770", 0},
		{T_SUB, "123.1230", "Inf", "-Infinity", 0},
		{T_SUB, "-123456789012345678901234567890.1234", "-123456789012345678901234567890.1234", "0.0000", 0},
		{T_SUB, "-123456789012345678901234567890.1234e200", "123456789012345678901234567890.1234", "-1.234567890123456789012345678901234E+229", 0},
		{T_SUB, "-123456789012345678901234567890.1234e200", "123456789012345678901234567890.1234e206", "-1.234569124691346912469134691246913E+235", 0},
		{T_SUB, minquad, "0", minquad, 0},
		{T_SUB, minquad, "1", minquad, 0},
		{T_SUB, minquad, "1e6111", "-Infinity", Overflow}, // Overflow
		{T_SUB, minquad, "Inf", "-Infinity", 0},
		{T_SUB, minquad, "-Inf", "Infinity", 0},
		{T_SUB, "142566.645373", "-647833330000004.7367", "647833330142571.382073", 0},
		{T_SUB, "1425658446.645373", "647833330000004.7367", "-647831904341558.091327", 0},
		{T_SUB, smallquad, "1", "-1.000000000000000000000000000000000", 0},

		{T_MUL, "1", "sNaN", "NaN", InvalidOperation},       // Invalid_operation      because of sNan (signaling NaN)
		{T_MUL, "1", "sNaN456", "NaN456", InvalidOperation}, // Invalid_operation      because of sNan (signaling NaN)
		{T_MUL, "NaN", "NaN", "NaN", 0},
		{T_MUL, "NaN", "123", "NaN", 0},
		{T_MUL, "NaN456", "123", "NaN456", 0},
		{T_MUL, "NaN", "Inf", "NaN", 0},
		{T_MUL, "123", "NaN", "NaN", 0},
		{T_MUL, "Inf", "NaN", "NaN", 0},
		{T_MUL, "Inf", "0", "NaN", InvalidOperation}, // Invalid_operation
		{T_MUL, "Inf", "Inf", "Infinity", 0},
		{T_MUL, "-Inf", "-Inf", "Infinity", 0},
		{T_MUL, "Inf", "-Inf", "-Infinity", 0},
		{T_MUL, "123.0", "200.0", "24600.00", 0},
		{T_MUL, "123.1230", "200", "24624.6000", 0},
		{T_MUL, "123.1230", "Inf", "Infinity", 0},
		{T_MUL, "-123456789012345678901234567890.1234", "55", "-6790123395679012339567901233956.787", 0},
		{T_MUL, "-123456789012345678901234567890.1234e200", "55", "-6.790123395679012339567901233956787E+230", 0},
		{T_MUL, "-123456789012345678901234567890.1234e200", "55e-205", "-67901233956790123395679012.33956787", 0},
		{T_MUL, "1e6000", "1e6000", "Infinity", Overflow}, // Overflow
		{T_MUL, maxquad, "2", "Infinity", Overflow},       // Overflow
		{T_MUL, maxquad, "1e-6144", "9.999999999999999999999999999999999", 0},
		{T_MUL, smallquad, "0.1", "1.000000000000000000000000000000000E-6143", Underflow},  // Underflow
		{T_MUL, smallquad, "1e-1", "1.000000000000000000000000000000000E-6143", Underflow}, // Underflow
		{T_MUL, smallquad, "1", smallquad, 0},
		{T_MUL, smallquad, "1.000", smallquad, 0},
		{T_MUL, "435648995.83677856", "15267.748590", "6651379341921.89172958223040", 0},

		{T_DIV, "1", "sNaN", "NaN", InvalidOperation},       // Invalid_operation      because of sNan (signaling NaN)
		{T_DIV, "1", "sNaN456", "NaN456", InvalidOperation}, // Invalid_operation      because of sNan (signaling NaN)
		{T_DIV, "NaN", "NaN", "NaN", 0},
		{T_DIV, "NaN", "123", "NaN", 0},
		{T_DIV, "NaN456", "123", "NaN456", 0},
		{T_DIV, "NaN", "Inf", "NaN", 0},
		{T_DIV, "123", "NaN", "NaN", 0},
		{T_DIV, "Inf", "NaN", "NaN", 0},
		{T_DIV, "Inf", "0", "Infinity", 0},
		{T_DIV, "Inf", "-0", "-Infinity", 0},
		{T_DIV, "Inf", "Inf", "NaN", InvalidOperation},   // Invalid_operation
		{T_DIV, "-Inf", "-Inf", "NaN", InvalidOperation}, // Invalid_operation
		{T_DIV, "Inf", "-Inf", "NaN", InvalidOperation},  // Invalid_operation
		{T_DIV, "123", "0", "Infinity", DivisionByZero},  // Division_by_zero
		{T_DIV, "123.0", "200.0", "0.615", 0},
		{T_DIV, "123.1230", "200", "0.615615", 0},
		{T_DIV, "123.1230", "Inf", "0E-6176", 0},
		{T_DIV, "-123456789012345678901234567890.1234", "55", "-2244668891133557798204264870.729516", 0},
		{T_DIV, "-123456789012345678901234567890.1234e200", "55", "-2.244668891133557798204264870729516E+227", 0},
		{T_DIV, "-123456789012345678901234567890.1234e200", "55e-205", "-2.244668891133557798204264870729516E+432", 0},
		{T_DIV, "1e6000", "1e6000", "1", 0},
		{T_DIV, "1e6000", "1e-6000", "Infinity", Overflow}, // Overflow
		{T_DIV, "1e-6000", "1e6000", "0E-6176", Underflow}, // Underflow
		{T_DIV, maxquad, "0.9999", "Infinity", Overflow},   // Overflow
		{T_DIV, maxquad, "1e6144", "9.999999999999999999999999999999999", 0},
		{T_DIV, smallquad, "10", "1.000000000000000000000000000000000E-6143", Underflow},  // Underflow
		{T_DIV, smallquad, "1e1", "1.000000000000000000000000000000000E-6143", Underflow}, // Underflow
		{T_DIV, smallquad, "1", smallquad, 0},
		{T_DIV, smallquad, "1.000", smallquad, 0},
		{T_DIV, "435648995.83677856", "15267.748590", "28533.93827313333825337767105582134", 0},
		{T_DIV, "1", "Inf", "0E-6176", 0},

		{T_DIVINT, "1", "sNaN", "NaN", InvalidOperation},       // Invalid_operation      because of sNan (signaling NaN)
		{T_DIVINT, "1", "sNaN456", "NaN456", InvalidOperation}, // Invalid_operation      because of sNan (signaling NaN)
		{T_DIVINT, "NaN", "NaN", "NaN", 0},
		{T_DIVINT, "NaN", "123", "NaN", 0},
		{T_DIVINT, "NaN456", "123", "NaN456", 0},
		{T_DIVINT, "NaN", "Inf", "NaN", 0},
		{T_DIVINT, "123", "NaN", "NaN", 0},
		{T_DIVINT, "Inf", "NaN", "NaN", 0},
		{T_DIVINT, "Inf", "0", "Infinity", 0},
		{T_DIVINT, "Inf", "-0", "-Infinity", 0},
		{T_DIVINT, "Inf", "Inf", "NaN", InvalidOperation},   // Invalid_operation
		{T_DIVINT, "-Inf", "-Inf", "NaN", InvalidOperation}, // Invalid_operation
		{T_DIVINT, "Inf", "-Inf", "NaN", InvalidOperation},  // Invalid_operation
		{T_DIVINT, "123", "0", "Infinity", DivisionByZero},  // Division_by_zero
		{T_DIVINT, "123.0", "50.0", "2", 0},
		{T_DIVINT, "123.0", "200.0", "0", 0},
		{T_DIVINT, "123.1230", "200", "0", 0},
		{T_DIVINT, "123.1230", "Inf", "0", 0},
		{T_DIVINT, "-123456789012345678901234567890.1234", "55", "-2244668891133557798204264870", 0},
		{T_DIVINT, "-123456789012345678901234567890.1234e200", "55", "NaN", DivisionImpossible},      // Division_impossible
		{T_DIVINT, "-123456789012345678901234567890.1234e200", "55e-205", "NaN", DivisionImpossible}, // Division_impossible
		{T_DIVINT, "1e6000", "1e6000", "1", 0},
		{T_DIVINT, "1e6000", "1e-6000", "NaN", DivisionImpossible}, // Division_impossible
		{T_DIVINT, "1e-6000", "1e6000", "0", 0},

		{T_MOD, "1", "sNaN", "NaN", InvalidOperation},       // Invalid_operation      because of sNan (signaling NaN)
		{T_MOD, "1", "sNaN456", "NaN456", InvalidOperation}, // Invalid_operation      because of sNan (signaling NaN)
		{T_MOD, "NaN", "NaN", "NaN", 0},
		{T_MOD, "NaN", "123", "NaN", 0},
		{T_MOD, "NaN", "1", "NaN", 0},
		{T_MOD, "NaN456", "123", "NaN456", 0},
		{T_MOD, "NaN", "Inf", "NaN", 0},
		{T_MOD, "123", "NaN", "NaN", 0},
		{T_MOD, "Inf", "NaN", "NaN", 0},
		{T_MOD, "Inf", "0", "NaN", InvalidOperation},     // Invalid_operation
		{T_MOD, "Inf", "-0", "NaN", InvalidOperation},    // Invalid_operation
		{T_MOD, "Inf", "Inf", "NaN", InvalidOperation},   // Invalid_operation
		{T_MOD, "-Inf", "-Inf", "NaN", InvalidOperation}, // Invalid_operation
		{T_MOD, "Inf", "-Inf", "NaN", InvalidOperation},  // Invalid_operation
		{T_MOD, "123", "0", "NaN", InvalidOperation},     // Invalid_operation
		{T_MOD, "123.0", "200.0", "123.0", 0},
		{T_MOD, "123.1230", "200", "123.1230", 0},
		{T_MOD, "123.1230", "Inf", "123.1230", 0},
		{T_MOD, "1e6000", "1e-6000", "NaN", DivisionImpossible}, // Division_impossible
		{T_MOD, "Inf", "2", "NaN", InvalidOperation},            // Invalid_operation

		{T_MAX, "sNaN", "1", "NaN", InvalidOperation},       // Invalid_operation      because of sNan (signaling NaN)
		{T_MAX, "sNaN456", "1", "NaN456", InvalidOperation}, // Invalid_operation      because of sNan (signaling NaN)
		{T_MAX, "NaN", "NaN", "NaN", 0},
		{T_MAX, "NaN456", "NaN", "NaN456", 0},
		{T_MAX, "NaN", "NaN456", "NaN", 0},
		{T_MAX, "NaN789", "NaN456", "NaN789", 0},
		{T_MAX, "NaN456", "NaN789", "NaN456", 0},
		{T_MAX, "NaN", "123", "123", 0},
		{T_MAX, "NaN456", "123", "123", 0},
		{T_MAX, "NaN", "Inf", "Infinity", 0},
		{T_MAX, "123", "NaN", "123", 0},
		{T_MAX, "Inf", "NaN", "Infinity", 0},
		{T_MAX, "Inf", "Inf", "Infinity", 0},
		{T_MAX, "-Inf", "-Inf", "-Infinity", 0},
		{T_MAX, "-Inf", "123", "123", 0},
		{T_MAX, "Inf", "-Inf", "Infinity", 0},
		{T_MAX, "Inf", "123", "Infinity", 0},
		{T_MAX, "123", "Inf", "Infinity", 0},
		{T_MAX, "12345.6700001", "12345.67", "12345.6700001", 0},
		{T_MAX, "12345.67000", "12345.67", "12345.67", 0},
		{T_MAX, "12345.669999", "12345.67", "12345.67", 0},
		{T_MAX, "-12345.6700001", "-12345.67", "-12345.67", 0},
		{T_MAX, "-12345.67000", "-12345.67", "-12345.67000", 0},
		{T_MAX, "-12345.669999", "-12345.67", "-12345.669999", 0},

		{T_MIN, "sNaN", "1", "NaN", InvalidOperation},       // Invalid_operation      because of sNan (signaling NaN)
		{T_MIN, "sNaN456", "1", "NaN456", InvalidOperation}, // Invalid_operation      because of sNan (signaling NaN)
		{T_MIN, "NaN", "NaN", "NaN", 0},
		{T_MIN, "NaN", "123", "123", 0},
		{T_MIN, "NaN456", "123", "123", 0},
		{T_MIN, "NaN", "Inf", "Infinity", 0},
		{T_MIN, "123", "NaN", "123", 0},
		{T_MIN, "Inf", "NaN", "Infinity", 0},
		{T_MIN, "Inf", "Inf", "Infinity", 0},
		{T_MIN, "-Inf", "-Inf", "-Infinity", 0},
		{T_MIN, "-Inf", "123", "-Infinity", 0},
		{T_MIN, "Inf", "-Inf", "-Infinity", 0},
		{T_MIN, "Inf", "123", "123", 0},
		{T_MIN, "123", "Inf", "123", 0},
		{T_MIN, "12345.6700001", "12345.67", "12345.67", 0},
		{T_MIN, "12345.67000", "12345.67", "12345.67000", 0},
		{T_MIN, "12345.669999", "12345.67", "12345.669999", 0},
		{T_MIN, "-12345.6700001", "-12345.67", "-12345.6700001", 0},
		{T_MIN, "-12345.67000", "-12345.67", "-12345.67", 0},
		{T_MIN, "-12345.669999", "-12345.67", "-12345.67", 0},

		{T_TOINTEGRAL, "sNaN", "RoundHalfEven", "NaN", InvalidOperation},       // Invalid_operation      because of sNan (signaling NaN)
		{T_TOINTEGRAL, "sNaN456", "RoundHalfEven", "NaN456", InvalidOperation}, // Invalid_operation      because of sNan (signaling NaN)
		{T_TOINTEGRAL, "NaN", "RoundHalfEven", "NaN", 0},
		{T_TOINTEGRAL, "NaN456", "RoundHalfEven", "NaN456", 0},
		{T_TOINTEGRAL, "Inf", "RoundHalfEven", "Infinity", 0},
		{T_TOINTEGRAL, "-Inf", "RoundHalfEven", "-Infinity", 0},
		{T_TOINTEGRAL, "12e3", "RoundHalfEven", "1.2E+4", 0},
		{T_TOINTEGRAL, "12e-3", "RoundHalfEven", "0", 0},
		{T_TOINTEGRAL, "12e-103", "RoundHalfEven", "0", 0},

		{T_TOINTEGRAL, "13256748.9879878", "RoundUp", "13256749", 0},
		{T_TOINTEGRAL, "13256748.1879878", "RoundUp", "13256749", 0},
		{T_TOINTEGRAL, "-13256748.1879878", "RoundUp", "-13256749", 0},
		{T_TOINTEGRAL, "-13256748.9879878", "RoundUp", "-13256749", 0},

		{T_TOINTEGRAL, "13256748.9879878", "RoundDown", "13256748", 0},
		{T_TOINTEGRAL, "13256748.1879878", "RoundDown", "13256748", 0},
		{T_TOINTEGRAL, "-13256748.1879878", "RoundDown", "-13256748", 0},
		{T_TOINTEGRAL, "-13256748.9879878", "RoundDown", "-13256748", 0},

		{T_TOINTEGRAL, "13256748.9879878", "RoundCeiling", "13256749", 0},
		{T_TOINTEGRAL, "13256748.1879878", "RoundCeiling", "13256749", 0},
		{T_TOINTEGRAL, "-13256748.1879878", "RoundCeiling", "-13256748", 0},
		{T_TOINTEGRAL, "-13256748.9879878", "RoundCeiling", "-13256748", 0},

		{T_TOINTEGRAL, "13256748.9879878", "RoundFloor", "13256748", 0},
		{T_TOINTEGRAL, "13256748.1879878", "RoundFloor", "13256748", 0},
		{T_TOINTEGRAL, "-13256748.1879878", "RoundFloor", "-13256749", 0},
		{T_TOINTEGRAL, "-13256748.9879878", "RoundFloor", "-13256749", 0},

		{T_TOINTEGRAL, "13256748.9879878", "RoundHalfEven", "13256749", 0},
		{T_TOINTEGRAL, "13256748.1879878", "RoundHalfEven", "13256748", 0},
		{T_TOINTEGRAL, "-13256748.1879878", "RoundHalfEven", "-13256748", 0},
		{T_TOINTEGRAL, "-13256748.9879878", "RoundHalfEven", "-13256749", 0},

		{T_TOINTEGRAL, maxquad, "RoundHalfEven", maxquad, 0},
		{T_TOINTEGRAL, minquad, "RoundHalfEven", minquad, 0},
		{T_TOINTEGRAL, smallquad, "RoundHalfEven", "0", 0},
		{T_TOINTEGRAL, nsmallquad, "RoundHalfEven", "0", 0},

		{T_TOINTEGRAL, "1234567890123456789012345678901234", "RoundHalfEven", "1234567890123456789012345678901234", 0},
		{T_TOINTEGRAL, "12345678901234567890123456789012341", "RoundHalfEven", "1.234567890123456789012345678901234E+34", 0},

		{T_QUANTIZE, "sNaN", "1", "NaN", InvalidOperation},       // Invalid_operation      because of sNan (signaling NaN)
		{T_QUANTIZE, "sNaN456", "1", "NaN456", InvalidOperation}, // Invalid_operation      because of sNan (signaling NaN)
		{T_QUANTIZE, "NaN", "NaN", "NaN", 0},
		{T_QUANTIZE, "NaN", "123", "NaN", 0},
		{T_QUANTIZE, "NaN456", "123", "NaN456", 0},
		{T_QUANTIZE, "NaN", "Inf", "NaN", 0},
		{T_QUANTIZE, "123", "NaN", "NaN", 0},
		{T_QUANTIZE, "123", "Inf", "NaN", InvalidOperation}, // Invalid_operation
		{T_QUANTIZE, "Inf", "NaN", "NaN", 0},
		{T_QUANTIZE, "Inf", "Inf", "Infinity", 0},
		{T_QUANTIZE, "-Inf", "-Inf", "-Infinity", 0},
		{T_QUANTIZE, "Inf", "-Inf", "Infinity", 0},
		{T_QUANTIZE, "Inf", "123", "NaN", InvalidOperation}, // Invalid_operation
		{T_QUANTIZE, "123", "200", "123", 0},
		{T_QUANTIZE, "123.132456784", "1", "123", 0},
		{T_QUANTIZE, "123.132456784", "10000000000000", "123", 0},
		{T_QUANTIZE, "123.132456784", "999999999999999999999", "123", 0},
		{T_QUANTIZE, "123.132456784", "1.000000000000", "123.132456784000", 0},
		{T_QUANTIZE, "123.132456784", "1.000000", "123.132457", 0},
		{T_QUANTIZE, "123.132456784", "1.", "123", 0},
		{T_QUANTIZE, "123.1230", "1e2", "1E+2", 0},
		{T_QUANTIZE, "12345.1230", "1e2", "1.23E+4", 0},
		{T_QUANTIZE, "123e31", "1", "1230000000000000000000000000000000", 0},
		{T_QUANTIZE, "123e32", "1", "NaN", InvalidOperation}, // Invalid_operation
		{T_QUANTIZE, "123e32", "1E1", "1.230000000000000000000000000000000E+34", 0},
		{T_QUANTIZE, "123e32", "10", "NaN", InvalidOperation}, // Invalid_operation

		{T_ABS, "sNaN", "", "NaN", InvalidOperation},       // Invalid_operation      because of sNan (signaling NaN)
		{T_ABS, "sNaN456", "", "NaN456", InvalidOperation}, // Invalid_operation      because of sNan (signaling NaN)
		{T_ABS, "NaN", "", "NaN", 0},
		{T_ABS, "NaN456", "", "NaN456", 0},
		{T_ABS, "Inf", "", "Infinity", 0},
		{T_ABS, "-Inf", "", "Infinity", 0},
		{T_ABS, "-13256748.9879878", "", "13256748.9879878", 0},
		{T_ABS, "13256748.9879878", "", "13256748.9879878", 0},
		{T_ABS, "-13256748.9879878e456", "", "1.32567489879878E+463", 0},
		{T_ABS, maxquad, "", maxquad, 0},
		{T_ABS, minquad, "", maxquad, 0},
		{T_ABS, smallquad, "", smallquad, 0},
		{T_ABS, nsmallquad, "", smallquad, 0},

		{T_ISFINITE, "sNaN", "", "false", 0},
		{T_ISFINITE, "sNaN456", "", "false", 0},
		{T_ISFINITE, "NaN", "", "false", 0},
		{T_ISFINITE, "NaN456", "", "false", 0},
		{T_ISFINITE, "Inf", "", "false", 0},
		{T_ISFINITE, "-Inf", "", "false", 0},
		{T_ISFINITE, "0.0000", "", "true", 0},
		{T_ISFINITE, "-0.0000", "", "true", 0},
		{T_ISFINITE, "1234", "", "true", 0},
		{T_ISFINITE, "1234.5", "", "true", 0},
		{T_ISFINITE, "-12.34e5", "", "true", 0},
		{T_ISFINITE, "12.34e5", "", "true", 0},
		{T_ISFINITE, maxquad, "", "true", 0},

		/* not used for the moment
		{T_ISINTEGER, "sNaN", "", "false", 0},
		{T_ISINTEGER, "sNaN456", "", "false", 0},
		{T_ISINTEGER, "NaN", "", "false", 0},
		{T_ISINTEGER, "NaN456", "", "false", 0},
		{T_ISINTEGER, "Inf", "", "false", 0},
		{T_ISINTEGER, "-Inf", "", "false", 0},
		{T_ISINTEGER, "0", "", "true", 0},
		{T_ISINTEGER, "0.0000", "", "false", 0},
		{T_ISINTEGER, "12.34e2", "", "true", 0},
		{T_ISINTEGER, "12.34e3", "", "false", 0},
		{T_ISINTEGER, "1", "", "true", 0},
		{T_ISINTEGER, "1.0000", "", "false", 0},
		{T_ISINTEGER, "-0.0000", "", "false", 0},
		{T_ISINTEGER, "1234", "", "true", 0},
		{T_ISINTEGER, "1234.5", "", "false", 0},
		{T_ISINTEGER, "-12.34e5", "", "false", 0},
		{T_ISINTEGER, "12.34e5", "", "false", 0},
		{T_ISINTEGER, maxquad, "", "false", 0},
		{T_ISINTEGER, "1e3", "", "false", 0},
		*/

		{T_ISINFINITE, "sNaN", "", "false", 0},
		{T_ISINFINITE, "sNaN456", "", "false", 0},
		{T_ISINFINITE, "NaN", "", "false", 0},
		{T_ISINFINITE, "NaN456", "", "false", 0},
		{T_ISINFINITE, "Inf", "", "true", 0},
		{T_ISINFINITE, "-Inf", "", "true", 0},
		{T_ISINFINITE, "0.0000", "", "false", 0},
		{T_ISINFINITE, "-0.0000", "", "false", 0},
		{T_ISINFINITE, "1234", "", "false", 0},
		{T_ISINFINITE, "1234.5", "", "false", 0},
		{T_ISINFINITE, "-12.34e5", "", "false", 0},
		{T_ISINFINITE, "12.34e5", "", "false", 0},
		{T_ISINFINITE, maxquad, "", "false", 0},

		{T_ISNAN, "sNaN", "", "true", 0},
		{T_ISNAN, "sNaN456", "", "true", 0},
		{T_ISNAN, "NaN", "", "true", 0},
		{T_ISNAN, "NaN456", "", "true", 0},
		{T_ISNAN, "Inf", "", "false", 0},
		{T_ISNAN, "-Inf", "", "false", 0},
		{T_ISNAN, "0.0000", "", "false", 0},
		{T_ISNAN, "-0.0000", "", "false", 0},
		{T_ISNAN, "1234", "", "false", 0},
		{T_ISNAN, "1234.5", "", "false", 0},
		{T_ISNAN, "-12.34e5", "", "false", 0},
		{T_ISNAN, "12.34e5", "", "false", 0},
		{T_ISNAN, maxquad, "", "false", 0},

		{T_ISPOSITIVE, "sNaN", "", "false", 0},
		{T_ISPOSITIVE, "sNaN456", "", "false", 0},
		{T_ISPOSITIVE, "NaN", "", "false", 0},
		{T_ISPOSITIVE, "NaN456", "", "false", 0},
		{T_ISPOSITIVE, "Inf", "", "true", 0},
		{T_ISPOSITIVE, "-Inf", "", "false", 0},
		{T_ISPOSITIVE, "0.0000", "", "false", 0},
		{T_ISPOSITIVE, "-0.0000", "", "false", 0},
		{T_ISPOSITIVE, "1234", "", "true", 0},
		{T_ISPOSITIVE, "1234.5", "", "true", 0},
		{T_ISPOSITIVE, "-12.34e5", "", "false", 0},
		{T_ISPOSITIVE, "12.34e5", "", "true", 0},
		{T_ISPOSITIVE, maxquad, "", "true", 0},

		{T_ISZERO, "sNaN", "", "false", 0},
		{T_ISZERO, "sNaN456", "", "false", 0},
		{T_ISZERO, "NaN", "", "false", 0},
		{T_ISZERO, "NaN456", "", "false", 0},
		{T_ISZERO, "Inf", "", "false", 0},
		{T_ISZERO, "-Inf", "", "false", 0},
		{T_ISZERO, "0.0000", "", "true", 0},
		{T_ISZERO, "-0.0000", "", "true", 0},
		{T_ISZERO, "1234", "", "false", 0},
		{T_ISZERO, "1234.5", "", "false", 0},
		{T_ISZERO, "-12.34e5", "", "false", 0},
		{T_ISZERO, "12.34e5", "", "false", 0},
		{T_ISZERO, maxquad, "", "false", 0},

		{T_ISNEGATIVE, "sNaN", "", "false", 0},
		{T_ISNEGATIVE, "sNaN456", "", "false", 0},
		{T_ISNEGATIVE, "NaN", "", "false", 0},
		{T_ISNEGATIVE, "NaN456", "", "false", 0},
		{T_ISNEGATIVE, "-NaN", "", "false", 0},
		{T_ISNEGATIVE, "Inf", "", "false", 0},
		{T_ISNEGATIVE, "-Inf", "", "true", 0},
		{T_ISNEGATIVE, "0.0000", "", "false", 0},
		{T_ISNEGATIVE, "-0.0000", "", "false", 0},
		{T_ISNEGATIVE, "1234", "", "false", 0},
		{T_ISNEGATIVE, "1234.5", "", "false", 0},
		{T_ISNEGATIVE, "-12.34e5", "", "true", 0},
		{T_ISNEGATIVE, "12.34e5", "", "false", 0},
		{T_ISNEGATIVE, maxquad, "", "false", 0},

		{T_GREATER, "sNaN", "1", "false", 0},    // Invalid_operation      because of sNan (signaling NaN)
		{T_GREATER, "sNaN456", "1", "false", 0}, // Invalid_operation      because of sNan (signaling NaN)
		{T_GREATER, "NaN", "NaN", "false", 0},
		{T_GREATER, "NaN", "123", "false", 0},
		{T_GREATER, "NaN456", "123", "false", 0},
		{T_GREATER, "NaN", "Inf", "false", 0},
		{T_GREATER, "123", "NaN", "false", 0},
		{T_GREATER, "Inf", "NaN", "false", 0},
		{T_GREATER, "Inf", "Inf", "false", 0},
		{T_GREATER, "-Inf", "-Inf", "false", 0},
		{T_GREATER, "-Inf", "123", "false", 0},
		{T_GREATER, "Inf", "-Inf", "true", 0},
		{T_GREATER, "Inf", "123", "true", 0},
		{T_GREATER, "Inf", "Inf", "false", 0},
		{T_GREATER, "12345.6700001", "12345.67", "true", 0},
		{T_GREATER, "12345.67000", "12345.67", "false", 0},
		{T_GREATER, "12345.669999", "12345.67", "false", 0},
		{T_GREATER, "-12345.6700001", "-12345.67", "false", 0},
		{T_GREATER, "-12345.67000", "-12345.67", "false", 0},
		{T_GREATER, "-12345.669999", "-12345.67", "true", 0},

		{T_GREATEREQUAL, "sNaN", "1", "false", 0},    // Invalid_operation      because of sNan (signaling NaN)
		{T_GREATEREQUAL, "sNaN456", "1", "false", 0}, // Invalid_operation      because of sNan (signaling NaN)
		{T_GREATEREQUAL, "NaN", "NaN", "false", 0},
		{T_GREATEREQUAL, "NaN", "123", "false", 0},
		{T_GREATEREQUAL, "NaN456", "123", "false", 0},
		{T_GREATEREQUAL, "NaN", "Inf", "false", 0},
		{T_GREATEREQUAL, "123", "NaN", "false", 0},
		{T_GREATEREQUAL, "Inf", "NaN", "false", 0},
		{T_GREATEREQUAL, "Inf", "Inf", "true", 0},
		{T_GREATEREQUAL, "-Inf", "-Inf", "true", 0},
		{T_GREATEREQUAL, "-Inf", "123", "false", 0},
		{T_GREATEREQUAL, "Inf", "-Inf", "true", 0},
		{T_GREATEREQUAL, "Inf", "123", "true", 0},
		{T_GREATEREQUAL, "Inf", "Inf", "true", 0},
		{T_GREATEREQUAL, "12345.6700001", "12345.67", "true", 0},
		{T_GREATEREQUAL, "12345.67000", "12345.67", "true", 0},
		{T_GREATEREQUAL, "12345.669999", "12345.67", "false", 0},
		{T_GREATEREQUAL, "-12345.6700001", "-12345.67", "false", 0},
		{T_GREATEREQUAL, "-12345.67000", "-12345.67", "true", 0},
		{T_GREATEREQUAL, "-12345.669999", "-12345.67", "true", 0},

		{T_EQUAL, "sNaN", "1", "false", 0},    // Invalid_operation      because of sNan (signaling NaN)
		{T_EQUAL, "sNaN456", "1", "false", 0}, // Invalid_operation      because of sNan (signaling NaN)
		{T_EQUAL, "NaN", "NaN", "false", 0},
		{T_EQUAL, "NaN", "123", "false", 0},
		{T_EQUAL, "NaN456", "123", "false", 0},
		{T_EQUAL, "NaN", "Inf", "false", 0},
		{T_EQUAL, "123", "NaN", "false", 0},
		{T_EQUAL, "Inf", "NaN", "false", 0},
		{T_EQUAL, "Inf", "Inf", "true", 0},
		{T_EQUAL, "-Inf", "-Inf", "true", 0},
		{T_EQUAL, "-Inf", "123", "false", 0},
		{T_EQUAL, "Inf", "-Inf", "false", 0},
		{T_EQUAL, "Inf", "123", "false", 0},
		{T_EQUAL, "Inf", "Inf", "true", 0},
		{T_EQUAL, "12345.6700001", "12345.67", "false", 0},
		{T_EQUAL, "12345.67000", "12345.67", "true", 0},
		{T_EQUAL, "12345.669999", "12345.67", "false", 0},
		{T_EQUAL, "-12345.6700001", "-12345.67", "false", 0},
		{T_EQUAL, "-12345.67000", "-12345.67", "true", 0},
		{T_EQUAL, "-12345.669999", "-12345.67", "false", 0},

		{T_LESSEQUAL, "sNaN", "1", "false", 0},    // Invalid_operation      because of sNan (signaling NaN)
		{T_LESSEQUAL, "sNaN456", "1", "false", 0}, // Invalid_operation      because of sNan (signaling NaN)
		{T_LESSEQUAL, "NaN", "NaN", "false", 0},
		{T_LESSEQUAL, "NaN", "123", "false", 0},
		{T_LESSEQUAL, "NaN456", "123", "false", 0},
		{T_LESSEQUAL, "NaN", "Inf", "false", 0},
		{T_LESSEQUAL, "123", "NaN", "false", 0},
		{T_LESSEQUAL, "Inf", "NaN", "false", 0},
		{T_LESSEQUAL, "Inf", "Inf", "true", 0},
		{T_LESSEQUAL, "-Inf", "-Inf", "true", 0},
		{T_LESSEQUAL, "-Inf", "123", "true", 0},
		{T_LESSEQUAL, "Inf", "-Inf", "false", 0},
		{T_LESSEQUAL, "Inf", "123", "false", 0},
		{T_LESSEQUAL, "Inf", "Inf", "true", 0},
		{T_LESSEQUAL, "12345.6700001", "12345.67", "false", 0},
		{T_LESSEQUAL, "12345.67000", "12345.67", "true", 0},
		{T_LESSEQUAL, "12345.669999", "12345.67", "true", 0},
		{T_LESSEQUAL, "-12345.6700001", "-12345.67", "true", 0},
		{T_LESSEQUAL, "-12345.67000", "-12345.67", "true", 0},
		{T_LESSEQUAL, "-12345.669999", "-12345.67", "false", 0},

		{T_LESS, "sNaN", "1", "false", 0},    // Invalid_operation      because of sNan (signaling NaN)
		{T_LESS, "sNaN456", "1", "false", 0}, // Invalid_operation      because of sNan (signaling NaN)
		{T_LESS, "NaN", "NaN", "false", 0},
		{T_LESS, "NaN", "123", "false", 0},
		{T_LESS, "NaN456", "123", "false", 0},
		{T_LESS, "NaN", "Inf", "false", 0},
		{T_LESS, "123", "NaN", "false", 0},
		{T_LESS, "Inf", "NaN", "false", 0},
		{T_LESS, "Inf", "Inf", "false", 0},
		{T_LESS, "-Inf", "-Inf", "false", 0},
		{T_LESS, "-Inf", "123", "true", 0},
		{T_LESS, "Inf", "-Inf", "false", 0},
		{T_LESS, "Inf", "123", "false", 0},
		{T_LESS, "Inf", "Inf", "false", 0},
		{T_LESS, "12345.6700001", "12345.67", "false", 0},
		{T_LESS, "12345.67000", "12345.67", "false", 0},
		{T_LESS, "12345.669999", "12345.67", "true", 0},
		{T_LESS, "-12345.6700001", "-12345.67", "true", 0},
		{T_LESS, "-12345.67000", "-12345.67", "false", 0},
		{T_LESS, "-12345.669999", "-12345.67", "false", 0},

		{T_FROMSTRING, "sNaN", "", "sNaN", 0},         // sNaN is returned without setting status error flag
		{T_FROMSTRING, "sNaN456", "", "sNaN456", 0},   // sNaN is returned without setting status error flag
		{T_FROMSTRING, "-sNaN456", "", "-sNaN456", 0}, // sNaN is returned without setting status error flag
		{T_FROMSTRING, "NaN", "", "NaN", 0},
		{T_FROMSTRING, "-NaN", "", "-NaN", 0},
		{T_FROMSTRING, "-NaN777", "", "-NaN777", 0},
		{T_FROMSTRING, "NaN123", "", "NaN123", 0},
		{T_FROMSTRING, "Inf", "", "Infinity", 0},
		{T_FROMSTRING, "Infinity", "", "Infinity", 0},
		{T_FROMSTRING, "-Inf", "", "-Infinity", 0},
		{T_FROMSTRING, "", "", "NaN", ConversionSyntax},     // Conversion_syntax
		{T_FROMSTRING, "aaa", "", "NaN", ConversionSyntax},  // Conversion_syntax
		{T_FROMSTRING, "qNaN", "", "NaN", ConversionSyntax}, // Conversion_syntax
		{T_FROMSTRING, "0", "", "0", 0},
		{T_FROMSTRING, "1.0", "", "1.0", 0},
		{T_FROMSTRING, "123.45e-45", "", "1.2345E-43", 0},
		{T_FROMSTRING, "   123.45e-45      ", "", "1.2345E-43", 0},
		{T_FROMSTRING, "                                                                    123.45e-45                                 ", "", "1.2345E-43", 0},
		{T_FROMSTRING, "12.3450000000000000000000000000000000000000000000000000000000000000000000000000000000001", "", "12.34500000000000000000000000000000", 0},
		{T_FROMSTRING, "5368487.87676533e-3546", "", "5.36848787676533E-3540", 0},
		{T_FROMSTRING, maxquad, "", maxquad, 0},
		{T_FROMSTRING, minquad, "", minquad, 0},
		{T_FROMSTRING, smallquad, "", smallquad, 0},
		{T_FROMSTRING, nsmallquad, "", nsmallquad, 0},
		{T_FROMSTRING, "9223372036854775807", "", "9223372036854775807", 0},
		{T_FROMSTRING, "-9223372036854775808", "", "-9223372036854775808", 0},
		{T_FROMSTRING, "2147483647", "", "2147483647", 0},
		{T_FROMSTRING, "9223372.036854775807", "", "9223372.036854775807", 0},
		{T_FROMSTRING, "-9223372.036854775808", "", "-9223372.036854775808", 0},
		{T_FROMSTRING, "214748.3647", "", "214748.3647", 0},

		{T_FROMINT32, "2147483647", "", "2147483647", 0},
		{T_FROMINT32, "-2147483648", "", "-2147483648", 0},
		{T_FROMINT32, "-12345", "", "-12345", 0},
		{T_FROMINT32, "53674956", "", "53674956", 0},
		{T_FROMINT32, "0", "", "0", 0},
		{T_FROMINT32, "-230", "", "-230", 0},
		{T_FROMINT32, "127", "", "127", 0},
		{T_FROMINT32, "128", "", "128", 0},
		{T_FROMINT32, "32767", "", "32767", 0},
		{T_FROMINT32, "32768", "", "32768", 0},

		{T_FROMINT64, "9223372036854775807", "", "9223372036854775807", 0},
		{T_FROMINT64, "-9223372036854775808", "", "-9223372036854775808", 0},
		{T_FROMINT64, "2147483647", "", "2147483647", 0},
		{T_FROMINT64, "-2147483648", "", "-2147483648", 0},
		{T_FROMINT64, "-12345", "", "-12345", 0},
		{T_FROMINT64, "53674956", "", "53674956", 0},
		{T_FROMINT64, "0", "", "0", 0},
		{T_FROMINT64, "-230", "", "-230", 0},
		{T_FROMINT64, "127", "", "127", 0},
		{T_FROMINT64, "128", "", "128", 0},
		{T_FROMINT64, "32767", "", "32767", 0},
		{T_FROMINT64, "32768", "", "32768", 0},

		{T_TOINT32, "sNan", "RoundHalfUp", "0", InvalidOperation}, // Invalid_operation
		{T_TOINT32, "Nan", "RoundHalfUp", "0", InvalidOperation},  // Invalid_operation
		{T_TOINT32, "Inf", "RoundHalfUp", "0", InvalidOperation},  // Invalid_operation
		{T_TOINT32, "2147483647", "RoundHalfUp", "2147483647", 0},
		{T_TOINT32, "2147483647.1", "RoundHalfUp", "2147483647", 0},
		{T_TOINT32, "2147483647.49999999999", "RoundHalfUp", "2147483647", 0},
		{T_TOINT32, "2147483647.5", "RoundHalfUp", "0", InvalidOperation}, // Invalid_operation
		{T_TOINT32, "2147483647.5", "RoundHalfDown", "2147483647", 0},
		{T_TOINT32, "12345.6452", "RoundHalfUp", "12346", 0},
		{T_TOINT32, "12345.6452", "RoundFloor", "12345", 0},
		{T_TOINT32, "0.000000", "RoundHalfUp", "0", 0},
		{T_TOINT32, "-2147483648", "RoundHalfDown", "-2147483648", 0},
		{T_TOINT32, "-2147483648.1", "RoundHalfDown", "-2147483648", 0},
		{T_TOINT32, "-2147483648.49999999999", "RoundHalfDown", "-2147483648", 0},
		{T_TOINT32, "-2147483648.5", "RoundHalfUp", "0", InvalidOperation}, // Invalid_operation
		{T_TOINT32, "-2147483648.5", "RoundHalfDown", "-2147483648", 0},
		{T_TOINT32, "-2147483648.5e-4", "RoundHalfDown", "-214748", 0},
		{T_TOINT32, "-2147483648.5e-400", "RoundHalfDown", "0", 0},

		{T_TOINT64, "sNan", "RoundHalfUp", "0", InvalidOperation}, // Invalid_operation
		{T_TOINT64, "Nan", "RoundHalfUp", "0", InvalidOperation},  // Invalid_operation
		{T_TOINT64, "Inf", "RoundHalfUp", "0", InvalidOperation},  // Invalid_operation
		{T_TOINT64, "9223372036854775807", "RoundHalfUp", "9223372036854775807", 0},
		{T_TOINT64, "9223372036854775807.1", "RoundHalfUp", "9223372036854775807", 0},
		{T_TOINT64, "9223372036854775807.49999999999", "RoundHalfUp", "9223372036854775807", 0},
		{T_TOINT64, "9223372036854775807.5", "RoundHalfUp", "0", InvalidOperation}, // Invalid_operation
		{T_TOINT64, "9223372036854775807.5", "RoundHalfDown", "9223372036854775807", 0},
		{T_TOINT64, "12345.6452", "RoundHalfUp", "12346", 0},
		{T_TOINT64, "12345.6452", "RoundFloor", "12345", 0},
		{T_TOINT64, "0.000000", "RoundHalfUp", "0", 0},
		{T_TOINT64, "-9223372036854775808", "RoundHalfDown", "-9223372036854775808", 0},
		{T_TOINT64, "-9223372036854775808.1", "RoundHalfDown", "-9223372036854775808", 0},
		{T_TOINT64, "-9223372036854775808.49999999999", "RoundHalfDown", "-9223372036854775808", 0},
		{T_TOINT64, "-9223372036854775808.5", "RoundHalfUp", "0", InvalidOperation}, // Invalid_operation
		{T_TOINT64, "-9223372036854775808.5", "RoundHalfDown", "-9223372036854775808", 0},
		{T_TOINT64, "-2147483648.5e-4", "RoundHalfDown", "-214748", 0},
		{T_TOINT64, "-2147483648.5e-400", "RoundHalfDown", "0", 0},

		{T_TOFLOAT64, "sNan", "", "NaN", 0},
		{T_TOFLOAT64, "Nan", "", "NaN", 0},
		{T_TOFLOAT64, "Inf", "", "+Inf", 0},
		{T_TOFLOAT64, "12345.250", "", "12345.250000", 0},
		{T_TOFLOAT64, "-12345.250", "", "-12345.250000", 0},
		{T_TOFLOAT64, "12345678901234", "", "12345678901234.000000", 0},
		{T_TOFLOAT64, "12.345e8", "", "1234500000.000000", 0},
		{T_TOFLOAT64, "1.23e2000", "", "NaN", InvalidOperation}, // InvalidOperation, because float64 doen's support exponent this large

		{T_QUADTOSTRING, "sNan456", "", "sNaN456", 0},
		{T_QUADTOSTRING, "-sNan456", "", "-sNaN456", 0},
		{T_QUADTOSTRING, "sNan", "", "sNaN", 0},
		{T_QUADTOSTRING, "sNan123", "", "sNaN123", 0},
		{T_QUADTOSTRING, "Nan", "", "NaN", 0},
		{T_QUADTOSTRING, "Nan123", "", "NaN123", 0},
		{T_QUADTOSTRING, "-Nan123", "", "-NaN123", 0},
		{T_QUADTOSTRING, "Inf", "", "Infinity", 0},
		{T_QUADTOSTRING, "-Inf", "", "-Infinity", 0},
		{T_QUADTOSTRING, "-0", "", "0", 0},
		{T_QUADTOSTRING, "28799.234235", "", "28799.234235", 0},
		{T_QUADTOSTRING, "28799.234235e1000", "", "2.8799234235E+1004", 0},
		{T_QUADTOSTRING, "0", "", "0", 0},
		{T_QUADTOSTRING, "0.0000001", "", "1E-7", 0},
		{T_QUADTOSTRING, "-123786954.4695460934e-5", "", "-1237.869544695460934", 0},
		{T_QUADTOSTRING, maxquad, "", maxquad, 0},
		{T_QUADTOSTRING, minquad, "", minquad, 0},
		{T_QUADTOSTRING, smallquad, "", smallquad, 0},
		{T_QUADTOSTRING, nsmallquad, "", nsmallquad, 0},
		{T_QUADTOSTRING, "1234567890123456789012345678901234", "", "1234567890123456789012345678901234", 0},
		{T_QUADTOSTRING, "12345678901234567890123.45678901234", "", "12345678901234567890123.45678901234", 0},
		{T_QUADTOSTRING, "1.234567890123456789012345678901234", "", "1.234567890123456789012345678901234", 0},
		{T_QUADTOSTRING, ".1234567890123456789012345678901234", "", "0.1234567890123456789012345678901234", 0},
		{T_QUADTOSTRING, ".01234567890123456789012345678901234", "", "0.01234567890123456789012345678901234", 0},
		{T_QUADTOSTRING, "1234567890123456", "", "1234567890123456", 0},
		{T_QUADTOSTRING, "1234567890123456e6", "", "1.234567890123456E+21", 0},
		{T_QUADTOSTRING, "1234567890.0000000", "", "1234567890.0000000", 0},
		{T_QUADTOSTRING, "0.0000000123456789000000000000000000", "", "1.23456789000000000000000000E-8", 0},
		{T_QUADTOSTRING, "1e-34", "", "1E-34", 0},
		{T_QUADTOSTRING, "1.7465e-34", "", "1.7465E-34", 0},
		{T_QUADTOSTRING, "12.3e5", "", "1.23E+6", 0},
		{T_QUADTOSTRING, "123836700e-5", "", "1238.36700", 0},
		{T_QUADTOSTRING, "0.0000000000", "", "0E-10", 0},

		{T_STRING, "sNan", "", "sNaN", 0},
		{T_STRING, "-sNan", "", "-sNaN", 0},
		{T_STRING, "sNan123", "", "sNaN123", 0},
		{T_STRING, "-sNan123", "", "-sNaN123", 0},
		{T_STRING, "Nan", "", "NaN", 0},
		{T_STRING, "-Nan", "", "-NaN", 0},
		{T_STRING, "Nan123", "", "NaN123", 0},
		{T_STRING, "-Nan123", "", "-NaN123", 0},
		{T_STRING, "Inf", "", "Infinity", 0},
		{T_STRING, "-Inf", "", "-Infinity", 0},
		{T_STRING, "-0", "", "0", 0},
		{T_STRING, "28799.234235", "", "28799.234235", 0},
		{T_STRING, "28799.234235e1000", "", "2.8799234235E+1004", 0},
		{T_STRING, "0", "", "0", 0},
		{T_STRING, "0.0000001", "", "0.0000001", 0}, // different than QuadToString
		{T_STRING, "-123786954.4695460934e-5", "", "-1237.869544695460934", 0},
		{T_STRING, maxquad, "", maxquad, 0},
		{T_STRING, minquad, "", minquad, 0},
		{T_STRING, smallquad, "", smallquad, 0},
		{T_STRING, nsmallquad, "", nsmallquad, 0},
		{T_STRING, "1234567890123456789012345678901234", "", "1234567890123456789012345678901234", 0},
		{T_STRING, "12345678901234567890123.45678901234", "", "12345678901234567890123.45678901234", 0},
		{T_STRING, "1.234567890123456789012345678901234", "", "1.234567890123456789012345678901234", 0},
		{T_STRING, ".1234567890123456789012345678901234", "", "0.1234567890123456789012345678901234", 0},
		{T_STRING, ".01234567890123456789012345678901234", "", "0.01234567890123456789012345678901234", 0},
		{T_STRING, "1234567890123456", "", "1234567890123456", 0},
		{T_STRING, "1234567890123456e6", "", "1.234567890123456E+21", 0},
		{T_STRING, "1234567890.0000000", "", "1234567890.0000000", 0},
		{T_STRING, "0.0000000123456789000000000000000000", "", "0.0000000123456789000000000000000000", 0}, // different than QuadToString
		{T_STRING, "1e-34", "", "0.0000000000000000000000000000000001", 0},                                // different than QuadToString
		{T_STRING, "1.7465e-34", "", "1.7465E-34", 0},
		{T_STRING, "12.3e5", "", "1.23E+6", 0},
		{T_STRING, "123836700e-5", "", "1238.36700", 0},
		{T_STRING, "0.0000000000", "", "0.0000000000", 0},

		{T_ROUND, "sNaN", "0", "NaN", InvalidOperation},       // Invalid_operation      because of sNan (signaling NaN)
		{T_ROUND, "sNaN456", "0", "NaN456", InvalidOperation}, // Invalid_operation      because of sNan (signaling NaN)
		{T_ROUND, "NaN", "0", "NaN", 0},
		{T_ROUND, "NaN456", "0", "NaN456", 0},
		{T_ROUND, "Inf", "0", "NaN", InvalidOperation},  // Invalid_operation
		{T_ROUND, "-Inf", "0", "NaN", InvalidOperation}, // Invalid_operation
		{T_ROUND, "-13256748.9879878", "0", "-13256749", 0},
		{T_ROUND, "13256748.9879878", "0", "13256749", 0},
		{T_ROUND, "9999999999999999999999999999999999", "0", "9999999999999999999999999999999999", 0},
		{T_ROUND, "8999999999999999999999999999999999", "-1", "9000000000000000000000000000000000", 0},
		{T_ROUND, "8999999999999999999999999999999999", "-33", "9000000000000000000000000000000000", 0},
		{T_ROUND, "8999999999999999999999999999999999", "-34", "NaN", InvalidOperation}, // Invalid_operation
		{T_ROUND, "9999999999999999999999999999999999", "1", "NaN", InvalidOperation},   // Invalid_operation
		{T_ROUND, "999999999999999999999999999999999", "1", "999999999999999999999999999999999.0", 0},
		{T_ROUND, "0.9999999999999999999999999999999999", "34", "0.9999999999999999999999999999999999", 0},
		{T_ROUND, "-13256748.9879878", "34", "NaN", InvalidOperation}, // Invalid_operation
		{T_ROUND, "-13256748.9879878", "26", "-13256748.98798780000000000000000000", 0},
		{T_ROUND, "-13256748.9879878", "27", "NaN", InvalidOperation}, // Invalid_operation
		{T_ROUND, "-13256748.9879878", "-8", "0", 0},
		{T_ROUND, "-13256748.9879878", "-7", "-10000000", 0},
		{T_ROUND, "-13256748.9879878", "-2", "-13256700", 0},
		{T_ROUND, "-1325674898.79878e-2", "-2", "-13256700", 0},
		{T_ROUND, "123e5", "2", "12300000.00", 0},
		{T_ROUND, "123e5", "-2", "12300000", 0},
		{T_ROUND, "-13256748.9879878", "-1", "-13256750", 0},
		{T_ROUND, "-13256748.9879878", "0", "-13256749", 0},
		{T_ROUND, "-13256748.9879878", "1", "-13256749.0", 0},
		{T_ROUND, "-13256748.9879878", "2", "-13256748.99", 0},
		{T_ROUND, "-13256748.9879878", "3", "-13256748.988", 0},
		{T_ROUND, "-13256748.495", "-8", "0", 0},
		{T_ROUND, "-13256748.495", "-7", "-10000000", 0},
		{T_ROUND, "-13256748.495", "-2", "-13256700", 0},
		{T_ROUND, "-13256748.495", "-1", "-13256750", 0},
		{T_ROUND, "-13256748.495", "0", "-13256748", 0},
		{T_ROUND, "-13256748.495", "1", "-13256748.5", 0},
		{T_ROUND, "-13256748.495", "2", "-13256748.50", 0},
		{T_ROUND, "-13256748.495", "3", "-13256748.495", 0},
		{T_ROUND, "-13256748.495", "4", "-13256748.4950", 0},
		{T_ROUND, "-13256748.9879878e456", "0", "NaN", InvalidOperation}, // Invalid_operation
		{T_ROUND, maxquad, "0", "NaN", InvalidOperation},                 // Invalid_operation
		{T_ROUND, minquad, "0", "NaN", InvalidOperation},                 // Invalid_operation
		{T_ROUND, smallquad, "0", "0", 0},
		{T_ROUND, nsmallquad, "0", "0", 0},
		{T_ROUND, "0.2", "34", "0.2000000000000000000000000000000000", 0},
		{T_ROUND, "0.2", "35", "NaN", InvalidOperation}, // Invalid_operation
		{T_ROUND, "1111111111111111111111111111111111", "-34", "0", 0},
		{T_ROUND, "1111111111111111111111111111111111", "-35", "0", 0},
		{T_ROUND, "1111111111111111111111111111111111", "-36", "NaN", InvalidOperation}, // Invalid_operation
		{T_ROUND, "9999999999999999999999999999999999", "-34", "NaN", InvalidOperation}, // Invalid_operation
		{T_ROUND, "9999999999999999999999999999999999", "-35", "0", 0},
		{T_ROUND, "9999999999999999999999999999999999", "-36", "NaN", InvalidOperation}, // Invalid_operation

		{T_TRUNCATE, "sNaN", "0", "NaN", InvalidOperation},       // Invalid_operation      because of sNan (signaling NaN)
		{T_TRUNCATE, "sNaN456", "0", "NaN456", InvalidOperation}, // Invalid_operation      because of sNan (signaling NaN)
		{T_TRUNCATE, "NaN", "0", "NaN", 0},
		{T_TRUNCATE, "NaN456", "0", "NaN456", 0},
		{T_TRUNCATE, "Inf", "0", "NaN", InvalidOperation},  // Invalid_operation
		{T_TRUNCATE, "-Inf", "0", "NaN", InvalidOperation}, // Invalid_operation
		{T_TRUNCATE, "-13256748.9879878", "0", "-13256748", 0},
		{T_TRUNCATE, "13256748.9879878", "0", "13256748", 0},
		{T_TRUNCATE, "9999999999999999999999999999999999", "0", "9999999999999999999999999999999999", 0},
		{T_TRUNCATE, "8999999999999999999999999999999999", "-1", "8999999999999999999999999999999990", 0},
		{T_TRUNCATE, "8999999999999999999999999999999999", "-33", "8000000000000000000000000000000000", 0},
		{T_TRUNCATE, "8999999999999999999999999999999999", "-34", "0", 0},
		{T_TRUNCATE, "9999999999999999999999999999999999", "1", "NaN", InvalidOperation}, // Invalid_operation
		{T_TRUNCATE, "999999999999999999999999999999999", "1", "999999999999999999999999999999999.0", 0},
		{T_TRUNCATE, "0.9999999999999999999999999999999999", "34", "0.9999999999999999999999999999999999", 0},
		{T_TRUNCATE, "0.9999999999999999999999999999999999", "33", "0.999999999999999999999999999999999", 0},
		{T_TRUNCATE, "-13256748.9879878", "34", "NaN", InvalidOperation}, // Invalid_operation
		{T_TRUNCATE, "-13256748.9879878", "26", "-13256748.98798780000000000000000000", 0},
		{T_TRUNCATE, "-13256748.9879878", "27", "NaN", InvalidOperation}, // Invalid_operation
		{T_TRUNCATE, "-13256748.9879878", "-8", "0", 0},
		{T_TRUNCATE, "-13256748.9879878", "-7", "-10000000", 0},
		{T_TRUNCATE, "-13256748.9879878", "-2", "-13256700", 0},
		{T_TRUNCATE, "-1325674898.79878e-2", "-2", "-13256700", 0},
		{T_TRUNCATE, "123e5", "2", "12300000.00", 0},
		{T_TRUNCATE, "123e5", "-2", "12300000", 0},
		{T_TRUNCATE, "-13256748.9879878", "-1", "-13256740", 0},
		{T_TRUNCATE, "-13256748.9879878", "0", "-13256748", 0},
		{T_TRUNCATE, "-13256748.9879878", "1", "-13256748.9", 0},
		{T_TRUNCATE, "-13256748.9879878", "2", "-13256748.98", 0},
		{T_TRUNCATE, "-13256748.9879878", "3", "-13256748.987", 0},
		{T_TRUNCATE, "-13256748.495", "-8", "0", 0},
		{T_TRUNCATE, "-13256748.495", "-7", "-10000000", 0},
		{T_TRUNCATE, "-13256748.495", "-2", "-13256700", 0},
		{T_TRUNCATE, "-13256748.495", "-1", "-13256740", 0},
		{T_TRUNCATE, "-13256748.495", "0", "-13256748", 0},
		{T_TRUNCATE, "-13256748.495", "1", "-13256748.4", 0},
		{T_TRUNCATE, "-13256748.495", "2", "-13256748.49", 0},
		{T_TRUNCATE, "-13256748.495", "3", "-13256748.495", 0},
		{T_TRUNCATE, "-13256748.495", "4", "-13256748.4950", 0},
		{T_TRUNCATE, "-13256748.9879878e456", "0", "NaN", InvalidOperation}, // Invalid_operation
		{T_TRUNCATE, maxquad, "0", "NaN", InvalidOperation},                 // Invalid_operation
		{T_TRUNCATE, minquad, "0", "NaN", InvalidOperation},                 // Invalid_operation
		{T_TRUNCATE, smallquad, "0", "0", 0},
		{T_TRUNCATE, nsmallquad, "0", "0", 0},
		{T_TRUNCATE, "0.2", "34", "0.2000000000000000000000000000000000", 0},
		{T_TRUNCATE, "0.2", "35", "NaN", InvalidOperation}, // Invalid_operation
		{T_TRUNCATE, "1111111111111111111111111111111111", "-34", "0", 0},
		{T_TRUNCATE, "1111111111111111111111111111111111", "-35", "0", 0},
		{T_TRUNCATE, "9999999999999999999999999999999999", "-34", "0", 0},
		{T_TRUNCATE, "9999999999999999999999999999999999", "-35", "0", 0},
		{T_TRUNCATE, "9999999999999999999999999999999999", "-36", "NaN", InvalidOperation}, // Invalid_operation

	}

	for i, sp := range samples {
		var (
			a      Quad
			b      Quad
			result Quad
			err    error
			status Status
			output string // operation output as string
		)

		switch sp.operation {
		case T_NEG:
			a = must_quad(sp.a)
			result = a.Neg()
			status = result.ErrorStatus()
			output = result.String()

		case T_ADD:
			a = must_quad(sp.a)
			b = must_quad(sp.b)
			result = a.Add(b)
			status = result.ErrorStatus()
			output = result.String()

		case T_SUB:
			a = must_quad(sp.a)
			b = must_quad(sp.b)
			result = a.Sub(b)
			status = result.ErrorStatus()
			output = result.String()

		case T_MUL:
			a = must_quad(sp.a)
			b = must_quad(sp.b)
			result = a.Mul(b)
			status = result.ErrorStatus()
			output = result.String()

		case T_DIV:
			a = must_quad(sp.a)
			b = must_quad(sp.b)
			result = a.Div(b)
			status = result.ErrorStatus()
			output = result.String()

		case T_DIVINT:
			a = must_quad(sp.a)
			b = must_quad(sp.b)
			result = a.DivInt(b)
			status = result.ErrorStatus()
			output = result.String()

		case T_MOD:
			a = must_quad(sp.a)
			b = must_quad(sp.b)
			result = a.Mod(b)
			status = result.ErrorStatus()
			output = result.String()

		case T_MAX:
			a = must_quad(sp.a)
			b = must_quad(sp.b)
			result = Max(a, b)
			status = result.ErrorStatus()
			output = result.String()

		case T_MIN:
			a = must_quad(sp.a)
			b = must_quad(sp.b)
			result = Min(a, b)
			status = result.ErrorStatus()
			output = result.String()

		case T_TOINTEGRAL:
			a = must_quad(sp.a)
			result = a.ToIntegral(must_rounding(sp.b))
			status = result.ErrorStatus()
			output = result.String()

		case T_QUANTIZE:
			a = must_quad(sp.a)
			b = must_quad(sp.b)
			result = a.Quantize(b, RoundHalfEven)
			status = result.ErrorStatus()
			output = result.String()

		case T_ABS:
			a = must_quad(sp.a)
			result = a.Abs()
			status = result.ErrorStatus()
			output = result.String()

		case T_ISFINITE:
			a = must_quad(sp.a)
			result_cmp_bool := a.IsFinite()
			output = bool2string(result_cmp_bool)

		/*
			case T_ISINTEGER:
				a = must_quad(sp.a)
				result_cmp_bool := a.IsInteger()
				output = bool2string(result_cmp_bool)
		*/

		case T_ISINFINITE:
			a = must_quad(sp.a)
			result_cmp_bool := a.IsInfinite()
			output = bool2string(result_cmp_bool)

		case T_ISNAN:
			a = must_quad(sp.a)
			result_cmp_bool := a.IsNaN()
			output = bool2string(result_cmp_bool)

		case T_ISPOSITIVE:
			a = must_quad(sp.a)
			result_cmp_bool := a.IsPositive()
			output = bool2string(result_cmp_bool)

		case T_ISZERO:
			a = must_quad(sp.a)
			result_cmp_bool := a.IsZero()
			output = bool2string(result_cmp_bool)

		case T_ISNEGATIVE:
			a = must_quad(sp.a)
			result_cmp_bool := a.IsNegative()
			output = bool2string(result_cmp_bool)

		case T_GREATER:
			a = must_quad(sp.a)
			b = must_quad(sp.b)
			result_cmp_bool := a.Greater(b)
			output = bool2string(result_cmp_bool)

		case T_GREATEREQUAL:
			a = must_quad(sp.a)
			b = must_quad(sp.b)
			result_cmp_bool := a.GreaterEqual(b)
			output = bool2string(result_cmp_bool)

		case T_EQUAL:
			a = must_quad(sp.a)
			b = must_quad(sp.b)
			result_cmp_bool := a.Equal(b)
			output = bool2string(result_cmp_bool)

		case T_LESSEQUAL:
			a = must_quad(sp.a)
			b = must_quad(sp.b)
			result_cmp_bool := a.LessEqual(b)
			output = bool2string(result_cmp_bool)

		case T_LESS:
			a = must_quad(sp.a)
			b = must_quad(sp.b)
			result_cmp_bool := a.Less(b)
			output = bool2string(result_cmp_bool)

		case T_FROMSTRING:
			result, err = FromString(sp.a)
			output = result.String()
			if err != nil {
				status = Status(err.(QuadError))
			}

		case T_FROMINT32:
			result = FromInt32(must_int32(sp.a))
			output = result.String()

		case T_FROMINT64:
			result = FromInt64(must_int64(sp.a))
			output = result.String()

		case T_QUADTOSTRING:
			a = must_quad(sp.a)
			output = a.QuadToString()

		case T_STRING:
			a = must_quad(sp.a)
			output = a.String()

		case T_TOINT32:
			a = must_quad(sp.a)
			result_int32, err := a.ToInt32(must_rounding(sp.b))
			if err != nil {
				status = Status(err.(QuadError))
			}
			output = strconv.Itoa(int(result_int32))

		case T_TOINT64:
			a = must_quad(sp.a)
			result_int64, err := a.ToInt64(must_rounding(sp.b))
			if err != nil {
				status = Status(err.(QuadError))
			}
			output = strconv.Itoa(int(result_int64))

		case T_TOFLOAT64:
			a = must_quad(sp.a)
			result_float64, err := a.ToFloat64()
			if err != nil {
				status = Status(err.(QuadError))
			}
			output = strconv.FormatFloat(result_float64, 'f', 6, 64)

		case T_ROUND:
			a = must_quad(sp.a)
			n := must_int32(sp.b)
			result = a.Round(n)
			status = result.ErrorStatus()
			output = result.String()

		case T_TRUNCATE:
			a = must_quad(sp.a)
			n := must_int32(sp.b)
			result = a.Truncate(n)
			status = result.ErrorStatus()
			output = result.String()

		default:
			panic("operation unknown")
		}

		// check status and output

		if status != sp.expected_error_status {
			t.Fatalf("sample %d, %s <%s, %s>:  \"%s\" (status) != \"%s\" (expected status)", i, sp.operation, sp.a, sp.b, status, sp.expected_error_status)
		}

		if output != sp.expected_result {
			t.Fatalf("sample %d, %s <%s, %s>:  \"%s\" (output) != \"%s\" (expected result)", i, sp.operation, sp.a, sp.b, output, sp.expected_result)
		}
	}
}
