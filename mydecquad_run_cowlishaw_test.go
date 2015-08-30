package decnum

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// This function runs the test files in cowlishaw_test_files/ directory.
// These test files are provided with the original C decNumber package, and have been downloaded from http://speleotrove.com/decimal, topic "Testcases" (http://speleotrove.com/decimal/dectest.zip).
// Only the test files correponding to Quad type, and for operations that have been implemented in our Go decnum package, are run.
//
func Test_cowlishaw(t *testing.T) {
	var (
		current_rounding RoundingMode = RoundHalfEven
	)

	dir := "cowlishaw_test_files"

	filename_list := []string{"dqMinus.decTest", "dqAdd.decTest", "dqSubtract.decTest", "dqMultiply.decTest", "dqDivide.decTest", "dqDivideInt.decTest", "dqRemainder.decTest", "dqAbs.decTest", "dqToIntegral.decTest", "dqQuantize.decTest", "dqCompare.decTest", "dqMax.decTest", "dqMin.decTest"}

	for _, file_path := range filename_list {

		fmt.Printf("--- processing %s ---\n", file_path)

		inputFile, err := os.Open(filepath.Join(dir, file_path))
		if err != nil {
			t.Fatal("Error opening input file:", err)
		}

		defer inputFile.Close()

		scanner := bufio.NewScanner(inputFile)

		for scanner.Scan() {
			process_line(t, &current_rounding, file_path, scanner.Text()) // current_rounding can be modified if "rounding" directive is found
		}

		if err := scanner.Err(); err != nil {
			t.Fatal(scanner.Err())
		}
	}
}

func process_line(t *testing.T, current_rounding *RoundingMode, file_path string, line_original string) {
	var (
		ctx Context
	)

	// initialize context

	ctx.InitDefaultQuad()

	ctx.SetRounding(*current_rounding)

	// analyze line

	line := strings.TrimSpace(line_original)

	if line == "" || // line is empty
		strings.HasPrefix(line, "version") || // line is a directive we don't use
		strings.HasPrefix(line, "extended") ||
		strings.HasPrefix(line, "clamp") ||
		strings.HasPrefix(line, "precision") ||
		strings.HasPrefix(line, "maxExponent") ||
		strings.HasPrefix(line, "minExponent") ||
		strings.HasPrefix(line, "--") || // line is comment
		strings.Contains(line, "#") { // we don't process line with #, tests are too specific
		return
	}

	// if rounding directive, set rounding

	if strings.HasPrefix(line, "rounding") {
		ss := strings.Split(line, ":")
		if len(ss) != 2 {
			t.Fatal("Bad 'rounding' directive in test file %s for line %s", file_path, line_original)
		}

		rounding_mode_string := strings.TrimSpace(ss[1])

		switch rounding_mode_string {
		case "ceiling":
			*current_rounding = RoundCeiling
		case "down":
			*current_rounding = RoundDown
		case "floor":
			*current_rounding = RoundFloor
		case "half_down":
			*current_rounding = RoundHalfDown
		case "half_even":
			*current_rounding = RoundHalfEven
		case "half_up":
			*current_rounding = RoundHalfUp
		case "up":
			*current_rounding = RoundUp
		case "05up":
			*current_rounding = Round05Up
		default:
			t.Fatalf("Unknown rounding mode %s", rounding_mode_string)
		}

		ctx.SetRounding(*current_rounding)

		return
	}

	// get fields from line

	fields := strings.Fields(line)

	test_name := fields[0]
	_ = test_name

	test_operator := fields[1]

	switch test_operator {
	case "apply": // skip, do nothing
		return

	case "minus":
		process_operation_1_operand(t, &ctx, (*Context).Minus, fields, file_path, line_original)

	case "add":
		process_operation_2_operands(t, &ctx, (*Context).Add, fields, file_path, line_original)

	case "subtract":
		process_operation_2_operands(t, &ctx, (*Context).Subtract, fields, file_path, line_original)

	case "multiply":
		process_operation_2_operands(t, &ctx, (*Context).Multiply, fields, file_path, line_original)

	case "divide":
		process_operation_2_operands(t, &ctx, (*Context).Divide, fields, file_path, line_original)

	case "divideint":
		process_operation_2_operands(t, &ctx, (*Context).DivideInteger, fields, file_path, line_original)

	case "remainder":
		process_operation_2_operands(t, &ctx, (*Context).Remainder, fields, file_path, line_original)

	case "abs":
		process_operation_1_operand(t, &ctx, (*Context).Abs, fields, file_path, line_original)

	case "tointegralx":
		a := must_from_string(t, &ctx, fields[2], file_path, line_original)
		if fields[3] != "->" {
			t.Fatalf("Bad -> in test file %s for line %s", file_path, line_original)
		}

		expected_result := must_from_string(t, &ctx, fields[4], file_path, line_original)

		r := ctx.ToIntegral(a, *current_rounding)

		if r.QuadToString() != expected_result.QuadToString() {
			t.Fatalf("Test failed in test file %s for line %s", file_path, line_original)
		}

	case "quantize":
		process_operation_2_operands(t, &ctx, (*Context).Quantize, fields, file_path, line_original)

	case "compare":
		a := must_from_string(t, &ctx, fields[2], file_path, line_original)
		b := must_from_string(t, &ctx, fields[3], file_path, line_original)
		if fields[4] != "->" {
			t.Fatalf("Bad -> in test file %s for line %s", file_path, line_original)
		}

		expected_result := must_from_string(t, &ctx, fields[5], file_path, line_original)
		if expected_result.IsNaN() {
			expected_result = NaN() // because we don't care about sign of NaN, or about its payload
		}

		var r Quad

		r_comp := ctx.Compare(a, b)

		switch r_comp {
		case CmpLess:
			r = ctx.FromInt32(-1)
		case CmpEqual:
			r = Zero()
		case CmpGreater:
			r = ctx.FromInt32(1)
		case CmpNaN:
			r = NaN()
		default:
			t.Fatal("impossible")
		}

		if r.QuadToString() != expected_result.QuadToString() {
			t.Fatalf("Test failed in test file %s for line %s", file_path, line_original)
		}

	case "max":
		process_operation_2_operands(t, &ctx, (*Context).Max, fields, file_path, line_original)

	case "min":
		process_operation_2_operands(t, &ctx, (*Context).Min, fields, file_path, line_original)

	default:
		t.Fatalf("Unknown operator in test file %s for line %s", file_path, line_original)
	}

	//fmt.Println(fields)

}

func process_operation_1_operand(t *testing.T, ctx *Context, f func(*Context, Quad) Quad, fields []string, file_path string, line_original string) {

	a := must_from_string(t, ctx, fields[2], file_path, line_original)
	if fields[3] != "->" {
		t.Fatalf("Bad -> in test file %s for line %s", file_path, line_original)
	}

	expected_result := must_from_string(t, ctx, fields[4], file_path, line_original)

	r := f(ctx, a)

	if r.QuadToString() != expected_result.QuadToString() {
		t.Fatalf("Test failed in test file %s for line %s", file_path, line_original)
	}

	expected_status := get_expected_status(fields[5:])

	if ctx.Status() != expected_status {
		t.Fatalf("Test failed in test file %s for line %s. Status %s != %s.", file_path, line_original, ctx.Status(), expected_status)
	}
}

func process_operation_2_operands(t *testing.T, ctx *Context, f func(*Context, Quad, Quad) Quad, fields []string, file_path string, line_original string) {

	a := must_from_string(t, ctx, fields[2], file_path, line_original)
	b := must_from_string(t, ctx, fields[3], file_path, line_original)
	if fields[4] != "->" {
		t.Fatalf("Bad -> in test file %s for line %s", file_path, line_original)
	}

	expected_result := must_from_string(t, ctx, fields[5], file_path, line_original)

	r := f(ctx, a, b)

	if r.QuadToString() != expected_result.QuadToString() {
		t.Fatalf("Test failed in test file %s for line %s", file_path, line_original)
	}

	expected_status := get_expected_status(fields[6:])

	if ctx.Status() != expected_status {
		t.Fatalf("Test failed in test file %s for line %s. Status %s != %s.", file_path, line_original, ctx.Status(), expected_status)
	}
}

// converts a string into a Quad.
// It is a fatal error if string is invalid, which should never happen with the test files we have.
//
func must_from_string(t *testing.T, ctx *Context, s string, file_path string, line_original string) Quad {

	if ctx.Error() != nil {
		t.Fatalf("Test failed in test file %s for line %s. At entry of must_from_string(%s), there is already an error %s", file_path, line_original, s, ctx.Error())
	}

	if len(s) > 0 && s[0] == '\'' { // delete opening quote if any
		s = s[1:]

		assert(s[len(s)-1] == '\'') // delete closing quote
		s = s[:len(s)-1]
	}

	q := ctx.FromString(s)

	if ctx.Error() != nil {
		t.Fatalf("Test failed in test file %s for line %s. must_from_string(%s) failed. %s", file_path, line_original, s, ctx.Error())
	}

	// we take this occasion to also test the conversion   string --> Quad --> string --> Quad

	q2 := ctx.FromString(q.String())

	if ctx.Error() != nil {
		t.Fatalf("Test failed in test file %s for line %s. must_from_string(%s) failed. %s", file_path, line_original, s, ctx.Error())
	}

	if q2.QuadToString() != q.QuadToString() { // in particular, we can have the case    -0 --> "0" --> 0
		if q2.IsZero() && q.IsZero() { // because FromString() always discard the "-" sign for 0 values. In financial applications, displaying "-0" is strange, and we prefer to avoid it.
			return q
		}

		t.Fatalf("q2 %s != q %s     %v   %v", q2, q, q2.QuadToString(), q.QuadToString())
	}

	return q
}

// return a status value with bits set as described by flags argument.
// If "--" is encountered, it is the start of a comment, and the function stops parsing flags.
//
func get_expected_status(flags []string) Status {
	var status Status

	for _, flag := range flags {

		if strings.HasPrefix(flag, "--") { // comment, no more flags on the line
			return status
		}

		switch flag {
		case "Conversion_syntax":
			status |= FlagConversionSyntax
		case "Division_by_zero":
			status |= FlagDivisionByZero
		case "Division_impossible":
			status |= FlagDivisionImpossible
		case "Division_undefined":
			status |= FlagDivisionUndefined
		case "Insufficient_storage":
			status |= FlagInsufficientStorage
		case "Inexact":
			status |= FlagInexact
		case "Invalid_context":
			status |= FlagInvalidContext
		case "Invalid_operation":
			status |= FlagInvalidOperation
		case "Overflow":
			status |= FlagOverflow
		case "Clamped":
			// status |= FlagClamped
		case "Rounded":
			// status |= FlagRounded
		case "Subnormal":
			// status |= FlagSubnormal
		case "Underflow":
			status |= FlagUnderflow
		default:
			panic("Unknown status flag: " + flag)
		}
	}

	return status
}
