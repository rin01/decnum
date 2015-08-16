package decnum

import (
	"log"
	"testing"
)

var (
	maxquad = "9.999999999999999999999999999999999E+6144"
	minquad = "-9.999999999999999999999999999999999E+6144"

	smallquad  = "9.999999999999999999999999999999999E-6143"
	nsmallquad = "-9.999999999999999999999999999999999E-6143"
)

// converts string to Quad or panics.
//
func must_quad(s string) Quad {
	var (
		ctx Context
		q   Quad
	)

	ctx.InitDefaultQuad()

	q = ctx.FromString(s)

	if err := ctx.Error(); err != nil {
		log.Fatalf("must_quad(%s) error: %s", s, err)
	}

	return q
}

// converts string to Round_mode_t or panics.
//
func must_rounding(s string) Round_mode_t {

	switch s {
	case "ROUND_CEILING":
		return ROUND_CEILING
	case "ROUND_DOWN":
		return ROUND_DOWN
	case "ROUND_FLOOR":
		return ROUND_FLOOR
	case "ROUND_HALF_DOWN":
		return ROUND_HALF_DOWN
	case "ROUND_HALF_EVEN":
		return ROUND_HALF_EVEN
	case "ROUND_HALF_UP":
		return ROUND_HALF_UP
	case "ROUND_UP":
		return ROUND_UP
	case "ROUND_05UP":
		return ROUND_05UP
	default:
		log.Fatalf("Unknown rounding mode %s", s)
	}

	panic("impossible")
}

func must_cmp(s string) Cmp_t {

	switch s {
	case "CMP_LESS":
		return CMP_LESS
	case "CMP_EQUAL":
		return CMP_EQUAL
	case "CMP_GREATER":
		return CMP_GREATER
	case "CMP_NAN":
		return CMP_NAN
	default:
		log.Fatalf("Unknown Cmp_t %s", s)
	}

	panic("impossible")
}

func bool2string(b bool) string {

	if b {
		return "true"
	}

	return "false"
}

func Test_operations(t *testing.T) {

	type Operation_t string

	const (
		T_MINUS          Operation_t = "Minus"
		T_ADD            Operation_t = "Add"
		T_SUBTRACT       Operation_t = "Subtract"
		T_MULTIPLY       Operation_t = "Multiply"
		T_DIVIDE         Operation_t = "Divide"
		T_DIVIDE_INTEGER Operation_t = "DivideInteger"
		T_REMAINDER      Operation_t = "Remainder"
		T_ABS            Operation_t = "Abs"
		T_TOINTEGRAL     Operation_t = "ToIntegral"
		T_QUANTIZE       Operation_t = "Quantize"
		T_COMPARE        Operation_t = "Compare"
		T_CMP_GE         Operation_t = "Cmp_GE" // Cmp(a, b, CMP_GREATER|CMP_EQUAL)
		T_GREATER        Operation_t = "Greater"
		T_GREATEREQUAL        Operation_t = "GreaterEqual"
		T_EQUAL        Operation_t = "Equal"
		T_LESSEQUAL        Operation_t = "LessEqual"
		T_LESS        Operation_t = "Less"
		T_ISFINITE        Operation_t = "IsFinite"
		T_ISINTEGER        Operation_t = "IsInteger"
		T_ISINFINITE        Operation_t = "IsInfinite"
		T_ISNAN        Operation_t = "IsNan"
		T_ISPOSITIVE        Operation_t = "IsPositive"
		T_ISZERO        Operation_t = "IsZero"
		T_ISNEGATIVE        Operation_t = "IsNegative"
	)

	var (
		ctx Context
	)

	var samples = []struct {
		operation       Operation_t // operation to test
		a               string      // first argument of operation to test. Type depends on operation.
		b               string      // 2nd argument of operation to test. Type depends on operation.
		expected_result string      // expected result of operation
		expected_error  bool        // true if Context status contains an error after operation
	}{
		{T_MINUS, "NaN", "", "NaN", false},
		{T_MINUS, "Inf", "", "-Infinity", false},
		{T_MINUS, "-Inf", "", "Infinity", false},
		{T_MINUS, "-13256748.9879878", "", "13256748.9879878", false},
		{T_MINUS, "13256748.9879878", "", "-13256748.9879878", false},
		{T_MINUS, "-13256748.9879878e456", "", "1.32567489879878E+463", false},
		{T_MINUS, maxquad, "", minquad, false},
		{T_MINUS, minquad, "", maxquad, false},
		{T_MINUS, smallquad, "", nsmallquad, false},
		{T_MINUS, nsmallquad, "", smallquad, false},

		{T_ADD, "NaN", "NaN", "NaN", false},
		{T_ADD, "NaN", "123", "NaN", false},
		{T_ADD, "NaN", "Inf", "NaN", false},
		{T_ADD, "123", "NaN", "NaN", false},
		{T_ADD, "Inf", "NaN", "NaN", false},
		{T_ADD, "Inf", "Inf", "Infinity", false},
		{T_ADD, "-Inf", "-Inf", "-Infinity", false},
		{T_ADD, "Inf", "-Inf", "", true}, // Invalid_operation
		{T_ADD, "123", "200", "323", false},
		{T_ADD, "123.1230", "200", "323.1230", false},
		{T_ADD, "123.1230", "Inf", "Infinity", false},
		{T_ADD, "-123456789012345678901234567890.1234", "123456789012345678901234567890.1234", "0.0000", false},
		{T_ADD, "-123456789012345678901234567890.1234e200", "123456789012345678901234567890.1234", "-1.234567890123456789012345678901234E+229", false},
		{T_ADD, "-123456789012345678901234567890.1234e200", "123456789012345678901234567890.1234e206", "1.234566655555566665555556666555555E+235", false},
		{T_ADD, "9999999999999999999999999999999999", "0", "9999999999999999999999999999999999", false},
		{T_ADD, "9999999999999999999999999999999999", "1", "1.000000000000000000000000000000000E+34", false},
		{T_ADD, "9999999999999999999999999999999999", "2", "1.000000000000000000000000000000000E+34", false},
		{T_ADD, maxquad, "0", maxquad, false},
		{T_ADD, maxquad, "1", maxquad, false},
		{T_ADD, maxquad, "1e6111", "", true}, // Overflow
		{T_ADD, maxquad, "Inf", "Infinity", false},
		{T_ADD, maxquad, "-Inf", "-Infinity", false},
		{T_ADD, "142566.645373", "647833330000004.7367", "647833330142571.382073", false},
		{T_ADD, "1425658446.645373", "-647833330000004.7367", "-647831904341558.091327", false},
		{T_ADD, smallquad, "1", "1.000000000000000000000000000000000", false},

		{T_SUBTRACT, "NaN", "NaN", "NaN", false},
		{T_SUBTRACT, "NaN", "123", "NaN", false},
		{T_SUBTRACT, "NaN", "Inf", "NaN", false},
		{T_SUBTRACT, "123", "NaN", "NaN", false},
		{T_SUBTRACT, "Inf", "NaN", "NaN", false},
		{T_SUBTRACT, "Inf", "Inf", "", true},   // Invalid_operation
		{T_SUBTRACT, "-Inf", "-Inf", "", true}, // Invalid_operation
		{T_SUBTRACT, "Inf", "-Inf", "Infinity", false},
		{T_SUBTRACT, "123", "200", "-77", false},
		{T_SUBTRACT, "123.1230", "200", "-76.8770", false},
		{T_SUBTRACT, "123.1230", "Inf", "-Infinity", false},
		{T_SUBTRACT, "-123456789012345678901234567890.1234", "-123456789012345678901234567890.1234", "0.0000", false},
		{T_SUBTRACT, "-123456789012345678901234567890.1234e200", "123456789012345678901234567890.1234", "-1.234567890123456789012345678901234E+229", false},
		{T_SUBTRACT, "-123456789012345678901234567890.1234e200", "123456789012345678901234567890.1234e206", "-1.234569124691346912469134691246913E+235", false},
		{T_SUBTRACT, minquad, "0", minquad, false},
		{T_SUBTRACT, minquad, "1", minquad, false},
		{T_SUBTRACT, minquad, "1e6111", "", true}, // Overflow
		{T_SUBTRACT, minquad, "Inf", "-Infinity", false},
		{T_SUBTRACT, minquad, "-Inf", "Infinity", false},
		{T_SUBTRACT, "142566.645373", "-647833330000004.7367", "647833330142571.382073", false},
		{T_SUBTRACT, "1425658446.645373", "647833330000004.7367", "-647831904341558.091327", false},
		{T_SUBTRACT, smallquad, "1", "-1.000000000000000000000000000000000", false},

		{T_MULTIPLY, "NaN", "NaN", "NaN", false},
		{T_MULTIPLY, "NaN", "123", "NaN", false},
		{T_MULTIPLY, "NaN", "Inf", "NaN", false},
		{T_MULTIPLY, "123", "NaN", "NaN", false},
		{T_MULTIPLY, "Inf", "NaN", "NaN", false},
		{T_MULTIPLY, "Inf", "0", "NaN", true}, // Invalid_operation
		{T_MULTIPLY, "Inf", "Inf", "Infinity", false},
		{T_MULTIPLY, "-Inf", "-Inf", "Infinity", false},
		{T_MULTIPLY, "Inf", "-Inf", "-Infinity", false},
		{T_MULTIPLY, "123.0", "200.0", "24600.00", false},
		{T_MULTIPLY, "123.1230", "200", "24624.6000", false},
		{T_MULTIPLY, "123.1230", "Inf", "Infinity", false},
		{T_MULTIPLY, "-123456789012345678901234567890.1234", "55", "-6790123395679012339567901233956.787", false},
		{T_MULTIPLY, "-123456789012345678901234567890.1234e200", "55", "-6.790123395679012339567901233956787E+230", false},
		{T_MULTIPLY, "-123456789012345678901234567890.1234e200", "55e-205", "-67901233956790123395679012.33956787", false},
		{T_MULTIPLY, "1e6000", "1e6000", "", true}, // Overflow
		{T_MULTIPLY, maxquad, "2", "", true},       // Overflow
		{T_MULTIPLY, maxquad, "1e-6144", "9.999999999999999999999999999999999", false},
		{T_MULTIPLY, smallquad, "0.1", "", true},  // Underflow
		{T_MULTIPLY, smallquad, "1e-1", "", true}, // Underflow
		{T_MULTIPLY, smallquad, "1", smallquad, false},
		{T_MULTIPLY, smallquad, "1.000", smallquad, false},
		{T_MULTIPLY, "435648995.83677856", "15267.748590", "6651379341921.89172958223040", false},

		{T_DIVIDE, "NaN", "NaN", "NaN", false},
		{T_DIVIDE, "NaN", "123", "NaN", false},
		{T_DIVIDE, "NaN", "Inf", "NaN", false},
		{T_DIVIDE, "123", "NaN", "NaN", false},
		{T_DIVIDE, "Inf", "NaN", "NaN", false},
		{T_DIVIDE, "Inf", "0", "Infinity", false},
		{T_DIVIDE, "Inf", "-0", "-Infinity", false},
		{T_DIVIDE, "Inf", "Inf", "", true},   // Invalid_operation
		{T_DIVIDE, "-Inf", "-Inf", "", true}, // Invalid_operation
		{T_DIVIDE, "Inf", "-Inf", "", true},  // Invalid_operation
		{T_DIVIDE, "123", "0", "", true},     // Invalid_operation
		{T_DIVIDE, "123.0", "200.0", "0.615", false},
		{T_DIVIDE, "123.1230", "200", "0.615615", false},
		{T_DIVIDE, "123.1230", "Inf", "0E-6176", false},
		{T_DIVIDE, "-123456789012345678901234567890.1234", "55", "-2244668891133557798204264870.729516", false},
		{T_DIVIDE, "-123456789012345678901234567890.1234e200", "55", "-2.244668891133557798204264870729516E+227", false},
		{T_DIVIDE, "-123456789012345678901234567890.1234e200", "55e-205", "-2.244668891133557798204264870729516E+432", false},
		{T_DIVIDE, "1e6000", "1e6000", "1", false},
		{T_DIVIDE, "1e6000", "1e-6000", "", true}, // Overflow
		{T_DIVIDE, "1e-6000", "1e6000", "", true}, // Underflow
		{T_DIVIDE, maxquad, "0.9999", "", true},   // Overflow
		{T_DIVIDE, maxquad, "1e6144", "9.999999999999999999999999999999999", false},
		{T_DIVIDE, smallquad, "10", "", true},  // Underflow
		{T_DIVIDE, smallquad, "1e1", "", true}, // Underflow
		{T_DIVIDE, smallquad, "1", smallquad, false},
		{T_DIVIDE, smallquad, "1.000", smallquad, false},
		{T_DIVIDE, "435648995.83677856", "15267.748590", "28533.93827313333825337767105582134", false},

		{T_DIVIDE_INTEGER, "NaN", "NaN", "NaN", false},
		{T_DIVIDE_INTEGER, "NaN", "123", "NaN", false},
		{T_DIVIDE_INTEGER, "NaN", "Inf", "NaN", false},
		{T_DIVIDE_INTEGER, "123", "NaN", "NaN", false},
		{T_DIVIDE_INTEGER, "Inf", "NaN", "NaN", false},
		{T_DIVIDE_INTEGER, "Inf", "0", "Infinity", false},
		{T_DIVIDE_INTEGER, "Inf", "-0", "-Infinity", false},
		{T_DIVIDE_INTEGER, "Inf", "Inf", "", true},   // Invalid_operation
		{T_DIVIDE_INTEGER, "-Inf", "-Inf", "", true}, // Invalid_operation
		{T_DIVIDE_INTEGER, "Inf", "-Inf", "", true},  // Invalid_operation
		{T_DIVIDE_INTEGER, "123", "0", "", true},     // Invalid_operation
		{T_DIVIDE_INTEGER, "123.0", "50.0", "2", false},
		{T_DIVIDE_INTEGER, "123.0", "200.0", "0", false},
		{T_DIVIDE_INTEGER, "123.1230", "200", "0", false},
		{T_DIVIDE_INTEGER, "123.1230", "Inf", "0", false},
		{T_DIVIDE_INTEGER, "-123456789012345678901234567890.1234", "55", "-2244668891133557798204264870", false},
		{T_DIVIDE_INTEGER, "-123456789012345678901234567890.1234e200", "55", "", true},                                               // Division_impossible
		{T_DIVIDE_INTEGER, "-123456789012345678901234567890.1234e200", "55e-205", "-2.244668891133557798204264870729516E+432", true}, // Division_impossible
		{T_DIVIDE_INTEGER, "1e6000", "1e6000", "1", false},
		{T_DIVIDE_INTEGER, "1e6000", "1e-6000", "", true}, // Overflow
		{T_DIVIDE_INTEGER, "1e-6000", "1e6000", "0", false},

		{T_REMAINDER, "NaN", "NaN", "NaN", false},
		{T_REMAINDER, "NaN", "123", "NaN", false},
		{T_REMAINDER, "NaN", "Inf", "NaN", false},
		{T_REMAINDER, "123", "NaN", "NaN", false},
		{T_REMAINDER, "Inf", "NaN", "NaN", false},
		{T_REMAINDER, "Inf", "0", "", true},     // Invalid_operation
		{T_REMAINDER, "Inf", "-0", "", true},    // Invalid_operation
		{T_REMAINDER, "Inf", "Inf", "", true},   // Invalid_operation
		{T_REMAINDER, "-Inf", "-Inf", "", true}, // Invalid_operation
		{T_REMAINDER, "Inf", "-Inf", "", true},  // Invalid_operation
		{T_REMAINDER, "123", "0", "", true},     // Invalid_operation
		{T_REMAINDER, "123.0", "200.0", "123.0", false},
		{T_REMAINDER, "123.1230", "200", "123.1230", false},
		{T_REMAINDER, "123.1230", "Inf", "123.1230", false},
		{T_REMAINDER, "1e6000", "1e-6000", "", true}, // Division_impossible
		{T_REMAINDER, "Inf", "2", "", true},          // Invalid_operation

		{T_ABS, "NaN", "", "NaN", false},
		{T_ABS, "Inf", "", "Infinity", false},
		{T_ABS, "-Inf", "", "Infinity", false},
		{T_ABS, "-13256748.9879878", "", "13256748.9879878", false},
		{T_ABS, "13256748.9879878", "", "13256748.9879878", false},
		{T_ABS, "-13256748.9879878e456", "", "1.32567489879878E+463", false},
		{T_ABS, maxquad, "", maxquad, false},
		{T_ABS, minquad, "", maxquad, false},
		{T_ABS, smallquad, "", smallquad, false},
		{T_ABS, nsmallquad, "", smallquad, false},

		{T_TOINTEGRAL, "NaN", "ROUND_HALF_EVEN", "NaN", false},
		{T_TOINTEGRAL, "Inf", "ROUND_HALF_EVEN", "Infinity", false},
		{T_TOINTEGRAL, "-Inf", "ROUND_HALF_EVEN", "-Infinity", false},

		{T_TOINTEGRAL, "13256748.9879878", "ROUND_UP", "13256749", false},
		{T_TOINTEGRAL, "13256748.1879878", "ROUND_UP", "13256749", false},
		{T_TOINTEGRAL, "-13256748.1879878", "ROUND_UP", "-13256749", false},
		{T_TOINTEGRAL, "-13256748.9879878", "ROUND_UP", "-13256749", false},

		{T_TOINTEGRAL, "13256748.9879878", "ROUND_DOWN", "13256748", false},
		{T_TOINTEGRAL, "13256748.1879878", "ROUND_DOWN", "13256748", false},
		{T_TOINTEGRAL, "-13256748.1879878", "ROUND_DOWN", "-13256748", false},
		{T_TOINTEGRAL, "-13256748.9879878", "ROUND_DOWN", "-13256748", false},

		{T_TOINTEGRAL, "13256748.9879878", "ROUND_CEILING", "13256749", false},
		{T_TOINTEGRAL, "13256748.1879878", "ROUND_CEILING", "13256749", false},
		{T_TOINTEGRAL, "-13256748.1879878", "ROUND_CEILING", "-13256748", false},
		{T_TOINTEGRAL, "-13256748.9879878", "ROUND_CEILING", "-13256748", false},

		{T_TOINTEGRAL, "13256748.9879878", "ROUND_FLOOR", "13256748", false},
		{T_TOINTEGRAL, "13256748.1879878", "ROUND_FLOOR", "13256748", false},
		{T_TOINTEGRAL, "-13256748.1879878", "ROUND_FLOOR", "-13256749", false},
		{T_TOINTEGRAL, "-13256748.9879878", "ROUND_FLOOR", "-13256749", false},

		{T_TOINTEGRAL, "13256748.9879878", "ROUND_HALF_EVEN", "13256749", false},
		{T_TOINTEGRAL, "13256748.1879878", "ROUND_HALF_EVEN", "13256748", false},
		{T_TOINTEGRAL, "-13256748.1879878", "ROUND_HALF_EVEN", "-13256748", false},
		{T_TOINTEGRAL, "-13256748.9879878", "ROUND_HALF_EVEN", "-13256749", false},

		{T_TOINTEGRAL, maxquad, "ROUND_HALF_EVEN", maxquad, false},
		{T_TOINTEGRAL, minquad, "ROUND_HALF_EVEN", minquad, false},
		{T_TOINTEGRAL, smallquad, "ROUND_HALF_EVEN", "0", false},
		{T_TOINTEGRAL, nsmallquad, "ROUND_HALF_EVEN", "0", false},



		{T_QUANTIZE, "NaN", "NaN", "NaN", false},
		{T_QUANTIZE, "NaN", "123", "NaN", false},
		{T_QUANTIZE, "NaN", "Inf", "NaN", false},
		{T_QUANTIZE, "123", "NaN", "NaN", false},
		{T_QUANTIZE, "123", "Inf", "", true}, // Invalid_operation
		{T_QUANTIZE, "Inf", "NaN", "NaN", false},
		{T_QUANTIZE, "Inf", "Inf", "Infinity", false},
		{T_QUANTIZE, "-Inf", "-Inf", "-Infinity", false},
		{T_QUANTIZE, "Inf", "-Inf", "Infinity", false},
		{T_QUANTIZE, "Inf", "123", "", true}, // Invalid_operation
		{T_QUANTIZE, "123", "200", "123", false},
		{T_QUANTIZE, "123.132456784", "1", "123", false},
		{T_QUANTIZE, "123.132456784", "10000000000000", "123", false},
		{T_QUANTIZE, "123.132456784", "999999999999999999999", "123", false},
		{T_QUANTIZE, "123.132456784", "1.000000000000", "123.132456784000", false},
		{T_QUANTIZE, "123.132456784", "1.000000", "123.132457", false},
		{T_QUANTIZE, "123.132456784", "1.", "123", false},
		{T_QUANTIZE, "123.1230", "1e2", "1E+2", false},
		{T_QUANTIZE, "12345.1230", "1e2", "1.23E+4", false},


		{T_COMPARE, "NaN", "NaN", "CMP_NAN", false},
		{T_COMPARE, "NaN", "123", "CMP_NAN", false},
		{T_COMPARE, "NaN", "Inf", "CMP_NAN", false},
		{T_COMPARE, "123", "NaN", "CMP_NAN", false},
		{T_COMPARE, "Inf", "NaN", "CMP_NAN", false},
		{T_COMPARE, "Inf", "Inf", "CMP_EQUAL", false},
		{T_COMPARE, "-Inf", "-Inf", "CMP_EQUAL", false},
		{T_COMPARE, "-Inf", "123", "CMP_LESS", false},
		{T_COMPARE, "Inf", "-Inf", "CMP_GREATER", false},
		{T_COMPARE, "Inf", "123", "CMP_GREATER", false},
		{T_COMPARE, "Inf", "Inf", "CMP_EQUAL", false},
		{T_COMPARE, "12345.6700001", "12345.67", "CMP_GREATER", false},
		{T_COMPARE, "12345.67000", "12345.67", "CMP_EQUAL", false},
		{T_COMPARE, "12345.669999", "12345.67", "CMP_LESS", false},
		{T_COMPARE, "-12345.6700001", "-12345.67", "CMP_LESS", false},
		{T_COMPARE, "-12345.67000", "-12345.67", "CMP_EQUAL", false},
		{T_COMPARE, "-12345.669999", "-12345.67", "CMP_GREATER", false},


		{T_CMP_GE, "12345.6700001", "12345.67", "true", false},
		{T_CMP_GE, "12345.67000", "12345.67", "true", false},
		{T_CMP_GE, "12345.669999", "12345.67", "false", false},
		{T_CMP_GE, "-12345.6700001", "-12345.67", "false", false},
		{T_CMP_GE, "-12345.67000", "-12345.67", "true", false},
		{T_CMP_GE, "-12345.669999", "-12345.67", "true", false},

		{T_GREATER, "NaN", "NaN", "false", false},
		{T_GREATER, "NaN", "123", "false", false},
		{T_GREATER, "NaN", "Inf", "false", false},
		{T_GREATER, "123", "NaN", "false", false},
		{T_GREATER, "Inf", "NaN", "false", false},
		{T_GREATER, "Inf", "Inf", "false", false},
		{T_GREATER, "-Inf", "-Inf", "false", false},
		{T_GREATER, "-Inf", "123", "false", false},
		{T_GREATER, "Inf", "-Inf", "true", false},
		{T_GREATER, "Inf", "123", "true", false},
		{T_GREATER, "Inf", "Inf", "false", false},
		{T_GREATER, "12345.6700001", "12345.67", "true", false},
		{T_GREATER, "12345.67000", "12345.67", "false", false},
		{T_GREATER, "12345.669999", "12345.67", "false", false},
		{T_GREATER, "-12345.6700001", "-12345.67", "false", false},
		{T_GREATER, "-12345.67000", "-12345.67", "false", false},
		{T_GREATER, "-12345.669999", "-12345.67", "true", false},

		{T_GREATEREQUAL, "NaN", "NaN", "false", false},
		{T_GREATEREQUAL, "NaN", "123", "false", false},
		{T_GREATEREQUAL, "NaN", "Inf", "false", false},
		{T_GREATEREQUAL, "123", "NaN", "false", false},
		{T_GREATEREQUAL, "Inf", "NaN", "false", false},
		{T_GREATEREQUAL, "Inf", "Inf", "true", false},
		{T_GREATEREQUAL, "-Inf", "-Inf", "true", false},
		{T_GREATEREQUAL, "-Inf", "123", "false", false},
		{T_GREATEREQUAL, "Inf", "-Inf", "true", false},
		{T_GREATEREQUAL, "Inf", "123", "true", false},
		{T_GREATEREQUAL, "Inf", "Inf", "true", false},
		{T_GREATEREQUAL, "12345.6700001", "12345.67", "true", false},
		{T_GREATEREQUAL, "12345.67000", "12345.67", "true", false},
		{T_GREATEREQUAL, "12345.669999", "12345.67", "false", false},
		{T_GREATEREQUAL, "-12345.6700001", "-12345.67", "false", false},
		{T_GREATEREQUAL, "-12345.67000", "-12345.67", "true", false},
		{T_GREATEREQUAL, "-12345.669999", "-12345.67", "true", false},

		{T_EQUAL, "NaN", "NaN", "false", false},
		{T_EQUAL, "NaN", "123", "false", false},
		{T_EQUAL, "NaN", "Inf", "false", false},
		{T_EQUAL, "123", "NaN", "false", false},
		{T_EQUAL, "Inf", "NaN", "false", false},
		{T_EQUAL, "Inf", "Inf", "true", false},
		{T_EQUAL, "-Inf", "-Inf", "true", false},
		{T_EQUAL, "-Inf", "123", "false", false},
		{T_EQUAL, "Inf", "-Inf", "false", false},
		{T_EQUAL, "Inf", "123", "false", false},
		{T_EQUAL, "Inf", "Inf", "true", false},
		{T_EQUAL, "12345.6700001", "12345.67", "false", false},
		{T_EQUAL, "12345.67000", "12345.67", "true", false},
		{T_EQUAL, "12345.669999", "12345.67", "false", false},
		{T_EQUAL, "-12345.6700001", "-12345.67", "false", false},
		{T_EQUAL, "-12345.67000", "-12345.67", "true", false},
		{T_EQUAL, "-12345.669999", "-12345.67", "false", false},

		{T_LESSEQUAL, "NaN", "NaN", "false", false},
		{T_LESSEQUAL, "NaN", "123", "false", false},
		{T_LESSEQUAL, "NaN", "Inf", "false", false},
		{T_LESSEQUAL, "123", "NaN", "false", false},
		{T_LESSEQUAL, "Inf", "NaN", "false", false},
		{T_LESSEQUAL, "Inf", "Inf", "true", false},
		{T_LESSEQUAL, "-Inf", "-Inf", "true", false},
		{T_LESSEQUAL, "-Inf", "123", "true", false},
		{T_LESSEQUAL, "Inf", "-Inf", "false", false},
		{T_LESSEQUAL, "Inf", "123", "false", false},
		{T_LESSEQUAL, "Inf", "Inf", "true", false},
		{T_LESSEQUAL, "12345.6700001", "12345.67", "false", false},
		{T_LESSEQUAL, "12345.67000", "12345.67", "true", false},
		{T_LESSEQUAL, "12345.669999", "12345.67", "true", false},
		{T_LESSEQUAL, "-12345.6700001", "-12345.67", "true", false},
		{T_LESSEQUAL, "-12345.67000", "-12345.67", "true", false},
		{T_LESSEQUAL, "-12345.669999", "-12345.67", "false", false},

		{T_LESS, "NaN", "NaN", "false", false},
		{T_LESS, "NaN", "123", "false", false},
		{T_LESS, "NaN", "Inf", "false", false},
		{T_LESS, "123", "NaN", "false", false},
		{T_LESS, "Inf", "NaN", "false", false},
		{T_LESS, "Inf", "Inf", "false", false},
		{T_LESS, "-Inf", "-Inf", "false", false},
		{T_LESS, "-Inf", "123", "true", false},
		{T_LESS, "Inf", "-Inf", "false", false},
		{T_LESS, "Inf", "123", "false", false},
		{T_LESS, "Inf", "Inf", "false", false},
		{T_LESS, "12345.6700001", "12345.67", "false", false},
		{T_LESS, "12345.67000", "12345.67", "false", false},
		{T_LESS, "12345.669999", "12345.67", "true", false},
		{T_LESS, "-12345.6700001", "-12345.67", "true", false},
		{T_LESS, "-12345.67000", "-12345.67", "false", false},
		{T_LESS, "-12345.669999", "-12345.67", "false", false},

		{T_ISFINITE, "NaN", "", "false", false},
		{T_ISFINITE, "Inf", "", "true", false},
		{T_ISFINITE, "-Inf", "", "true", false},
		{T_ISFINITE, "0.0000", "", "false", false},
		{T_ISFINITE, "1234", "", "false", false},
		{T_ISFINITE, maxquad, "", "false", false},



		T_ISINTEGER        Operation_t = "IsInteger"
		T_ISINFINITE        Operation_t = "IsInfinite"
		T_ISNAN        Operation_t = "IsNan"
		T_ISPOSITIVE        Operation_t = "IsPositive"
		T_ISZERO        Operation_t = "IsZero"
		T_ISNEGATIVE        Operation_t = "IsNegative"

	}

	ctx.InitDefaultQuad()

	for i, sp := range samples {
		var (
			result Quad
			status Status_t
			output string // operation output as string
		)

		ctx.ResetStatus()

		switch sp.operation {
		case T_MINUS:
			result = ctx.Minus(must_quad(sp.a))
			output = result.String()

		case T_ADD:
			result = ctx.Add(must_quad(sp.a), must_quad(sp.b))
			output = result.String()

		case T_SUBTRACT:
			result = ctx.Subtract(must_quad(sp.a), must_quad(sp.b))
			output = result.String()

		case T_MULTIPLY:
			result = ctx.Multiply(must_quad(sp.a), must_quad(sp.b))
			output = result.String()

		case T_DIVIDE:
			result = ctx.Divide(must_quad(sp.a), must_quad(sp.b))
			output = result.String()

		case T_DIVIDE_INTEGER:
			result = ctx.DivideInteger(must_quad(sp.a), must_quad(sp.b))
			output = result.String()

		case T_REMAINDER:
			result = ctx.Remainder(must_quad(sp.a), must_quad(sp.b))
			output = result.String()

		case T_ABS:
			result = ctx.Abs(must_quad(sp.a))
			output = result.String()

		case T_TOINTEGRAL:
			result = ctx.ToIntegral(must_quad(sp.a), must_rounding(sp.b))
			output = result.String()

		case T_QUANTIZE:
			result = ctx.Quantize(must_quad(sp.a), must_quad(sp.b))
			output = result.String()

		case T_COMPARE:
			result_cmp := ctx.Compare(must_quad(sp.a), must_quad(sp.b))
			output = result_cmp.String()

		case T_CMP_GE:
			result_cmp_bool := ctx.Cmp(must_quad(sp.a), must_quad(sp.b), CMP_GREATER|CMP_EQUAL)
			output = bool2string(result_cmp_bool)

		case T_GREATER:
			result_cmp_bool := ctx.Greater(must_quad(sp.a), must_quad(sp.b))
			output = bool2string(result_cmp_bool)

		case T_GREATEREQUAL:
			result_cmp_bool := ctx.GreaterEqual(must_quad(sp.a), must_quad(sp.b))
			output = bool2string(result_cmp_bool)

		case T_EQUAL:
			result_cmp_bool := ctx.Equal(must_quad(sp.a), must_quad(sp.b))
			output = bool2string(result_cmp_bool)

		case T_LESSEQUAL:
			result_cmp_bool := ctx.LessEqual(must_quad(sp.a), must_quad(sp.b))
			output = bool2string(result_cmp_bool)

		case T_LESS:
			result_cmp_bool := ctx.Less(must_quad(sp.a), must_quad(sp.b))
			output = bool2string(result_cmp_bool)

		case T_ISFINITE:
			result_cmp_bool := must_quad(sp.a).IsFinite()
			output = bool2string(result_cmp_bool)

		case T_ISINTEGER:
			result_cmp_bool := must_quad(sp.a).IsInteger()
			output = bool2string(result_cmp_bool)

		case T_ISINFINITE:
			result_cmp_bool := must_quad(sp.a).IsInfinite()
			output = bool2string(result_cmp_bool)

		case T_ISNAN:
			result_cmp_bool := must_quad(sp.a).IsNan()
			output = bool2string(result_cmp_bool)

		case T_ISPOSITIVE:
			result_cmp_bool := must_quad(sp.a).IsPositive()
			output = bool2string(result_cmp_bool)

		case T_ISZERO:
			result_cmp_bool := must_quad(sp.a).IsZero()
			output = bool2string(result_cmp_bool)

		case T_ISNEGATIVE:
			result_cmp_bool := must_quad(sp.a).IsNegative()
			output = bool2string(result_cmp_bool)

		default:
			panic("operation unknown")
		}

		status = ctx.Status() & ErrorMask // discard informational flags, keep only error flags

		switch sp.expected_error {
		case true: // error expected
			if status == 0 { // but got none
				t.Fatalf("sample %d, %s <%s, %s>, error expected. Got %s", i, sp.operation, sp.a, sp.b, output)
			}

		default: // no error expected
			switch {
			case status != 0: // but got error
				t.Fatalf("sample %d, %s <%s, %s>, error not expected. Got error: %s", i, sp.operation, sp.a, sp.b, status)

			case output != sp.expected_result: // no error, but bad result
				t.Fatalf("sample %d, %s <%s, %s>, incorrect result. Got %s, expected %s", i, sp.operation, sp.a, sp.b, output, sp.expected_result)
			}
		}
	}
}
