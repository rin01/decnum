package decnum

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const DEBUG_PRINT_PROCESSED_LINES bool = false // #####    set to true if you want to list all the lines in test files that have been processed    #####

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
		process_operation_1_operand(t, Quad.Neg, fields, file_path, line_original, *current_rounding)

	case "add":
		if *current_rounding != RoundHalfEven {
			return
		}
		process_operation_2_operands(t, Quad.Add, fields, file_path, line_original, *current_rounding)

	case "subtract":
		if *current_rounding != RoundHalfEven {
			return
		}
		process_operation_2_operands(t, Quad.Sub, fields, file_path, line_original, *current_rounding)

	case "multiply":
		process_operation_2_operands(t, Quad.Mul, fields, file_path, line_original, *current_rounding)

	case "divide":
		process_operation_2_operands(t, Quad.Div, fields, file_path, line_original, *current_rounding)

	case "divideint":
		process_operation_2_operands(t, Quad.DivInt, fields, file_path, line_original, *current_rounding)

	case "remainder":
		process_operation_2_operands(t, Quad.Mod, fields, file_path, line_original, *current_rounding)

	case "abs":
		process_operation_1_operand(t, Quad.Abs, fields, file_path, line_original, *current_rounding)

	case "tointegralx": // status is not checked because "The DEC_Inexact flag is not set by decQuadToIntegralValue, even if rounding ocurred".
		a := must_from_string(t, fields[2], file_path, line_original)
		if fields[3] != "->" {
			t.Fatalf("Bad -> in test file %s for line %s", file_path, line_original)
		}

		expected_result := must_from_string(t, fields[4], file_path, line_original)

		r := a.ToIntegral(*current_rounding)

		if r.QuadToString() != expected_result.QuadToString() {
			t.Fatalf("Test failed in test file %s for line %s. Result %s != %s. Rounding mode is %s.", file_path, line_original, r.QuadToString(), expected_result, *current_rounding)
		}

	case "quantize":
		process_operation_2_operands_and_rounding(t, Quad.Quantize, fields, file_path, line_original, *current_rounding)

	case "compare":
		var err error
		var expected_result_int32 int32 = 123456 // initialize with invalid value

		a := must_from_string(t, fields[2], file_path, line_original)
		b := must_from_string(t, fields[3], file_path, line_original)
		if fields[4] != "->" {
			t.Fatalf("Bad -> in test file %s for line %s", file_path, line_original)
		}

		expected_result := must_from_string(t, fields[5], file_path, line_original)

		r_greater := a.Greater(b)
		r_greater_equal := a.GreaterEqual(b)
		r_equal := a.Equal(b)
		r_less_equal := a.LessEqual(b)
		r_less := a.Less(b)

		if expected_result.IsNaN() == false {
			if expected_result_int32, err = expected_result.ToInt32(RoundHalfEven); err != nil {
				panic("impossible")
			}
		}

		failed_flag := false
		switch {
		case expected_result.IsNaN():
			if !(r_greater == false && r_greater_equal == false && r_equal == false && r_less_equal == false && r_less == false) {
				failed_flag = true
			}

		case expected_result_int32 == -1:
			if !(r_greater == false && r_greater_equal == false && r_equal == false && r_less_equal == true && r_less == true) {
				failed_flag = true
			}

		case expected_result_int32 == 0:
			if !(r_greater == false && r_greater_equal == true && r_equal == true && r_less_equal == true && r_less == false) {
				failed_flag = true
			}

		case expected_result_int32 == 1:
			if !(r_greater == true && r_greater_equal == true && r_equal == false && r_less_equal == false && r_less == false) {
				failed_flag = true
			}

		default:
			t.Fatal("impossible")
		}

		if failed_flag {
			t.Fatalf("Test failed in test file %s for line %s", file_path, line_original)
		}

	case "max":
		process_operation_2_operands(t, Max, fields, file_path, line_original, *current_rounding)

	case "min":
		process_operation_2_operands(t, Min, fields, file_path, line_original, *current_rounding)

	default:
		t.Fatalf("Unknown operator in test file %s for line %s", file_path, line_original)
	}

	if DEBUG_PRINT_PROCESSED_LINES {
		fmt.Printf("%-20s  %s\n", *current_rounding, line_original)
	}
}

func process_operation_1_operand(t *testing.T, f func(Quad) Quad, fields []string, file_path string, line_original string, rounding_mode RoundingMode) {

	a := must_from_string(t, fields[2], file_path, line_original)
	if fields[3] != "->" {
		t.Fatalf("Bad -> in test file %s for line %s", file_path, line_original)
	}

	expected_result := must_from_string(t, fields[4], file_path, line_original)

	r := f(a)

	if r.QuadToString() != expected_result.QuadToString() {
		t.Fatalf("Test failed in test file %s for line %s. Result %s != %s. Rounding mode is %s.", file_path, line_original, r.QuadToString(), expected_result, rounding_mode)
	}

	expected_status := get_expected_status(fields[5:])

	if r.Status() != expected_status {
		t.Fatalf("Test failed in test file %s for line %s. Status %s != %s. Rounding mode is %s.", file_path, line_original, r.Status(), expected_status, rounding_mode)
	}
}

func process_operation_1_operand_and_rounding(t *testing.T, f func(Quad, RoundingMode) Quad, fields []string, file_path string, line_original string, rounding_mode RoundingMode) {

	a := must_from_string(t, fields[2], file_path, line_original)
	if fields[3] != "->" {
		t.Fatalf("Bad -> in test file %s for line %s", file_path, line_original)
	}

	expected_result := must_from_string(t, fields[4], file_path, line_original)

	r := f(a, rounding_mode)

	if r.QuadToString() != expected_result.QuadToString() {
		t.Fatalf("Test failed in test file %s for line %s. Result %s != %s. Rounding mode is %s.", file_path, line_original, r.QuadToString(), expected_result, rounding_mode)
	}

	expected_status := get_expected_status(fields[5:])

	if r.Status() != expected_status {
		t.Fatalf("Test failed in test file %s for line %s. Status %s != %s. Rounding mode is %s.", file_path, line_original, r.Status(), expected_status, rounding_mode)
	}
}

func process_operation_2_operands(t *testing.T, f func(Quad, Quad) Quad, fields []string, file_path string, line_original string, rounding_mode RoundingMode) {

	a := must_from_string(t, fields[2], file_path, line_original)
	b := must_from_string(t, fields[3], file_path, line_original)
	if fields[4] != "->" {
		t.Fatalf("Bad -> in test file %s for line %s", file_path, line_original)
	}

	expected_result := must_from_string(t, fields[5], file_path, line_original)

	r := f(a, b)

	if r.QuadToString() != expected_result.QuadToString() {
		t.Fatalf("Test failed in test file %s for line %s. Result %s != %s. Rounding mode is %s.", file_path, line_original, r.QuadToString(), expected_result, rounding_mode)
	}

	expected_status := get_expected_status(fields[6:])

	if r.Status() != expected_status {
		t.Fatalf("Test failed in test file %s for line %s. Status %s != %s. Rounding mode is %s.", file_path, line_original, r.Status(), expected_status, rounding_mode)
	}
}

func process_operation_2_operands_and_rounding(t *testing.T, f func(Quad, Quad, RoundingMode) Quad, fields []string, file_path string, line_original string, rounding_mode RoundingMode) {

	a := must_from_string(t, fields[2], file_path, line_original)
	b := must_from_string(t, fields[3], file_path, line_original)
	if fields[4] != "->" {
		t.Fatalf("Bad -> in test file %s for line %s", file_path, line_original)
	}

	expected_result := must_from_string(t, fields[5], file_path, line_original)

	r := f(a, b, rounding_mode)

	if r.QuadToString() != expected_result.QuadToString() {
		t.Fatalf("Test failed in test file %s for line %s. Result %s != %s. Rounding mode is %s.", file_path, line_original, r.QuadToString(), expected_result, rounding_mode)
	}

	expected_status := get_expected_status(fields[6:])

	if r.Status() != expected_status {
		t.Fatalf("Test failed in test file %s for line %s. Status %s != %s. Rounding mode is %s.", file_path, line_original, r.Status(), expected_status, rounding_mode)
	}
}

// converts a string into a Quad.
// It is a fatal error if string is invalid, which should never happen with the test files we have.
//
func must_from_string(t *testing.T, s string, file_path string, line_original string) Quad {

	if len(s) > 0 && s[0] == '\'' { // delete opening quote if any
		s = s[1:]

		assert(s[len(s)-1] == '\'') // delete closing quote
		s = s[:len(s)-1]
	}

	q, _ := FromString(s)
	if q.Error() != nil {
		t.Fatalf("Test failed in test file %s for line %s. must_from_string(%s) failed. %s", file_path, line_original, s, q.Error())
	}

	// we take this occasion to also test the conversion   string --> Quad --> string --> Quad

	q2, _ := FromString(q.String())
	if q2.Error() != nil {
		t.Fatalf("Test failed in test file %s for line %s. must_from_string(%s) failed. %s", file_path, line_original, s, q2.Error())
	}

	if q2.QuadToString() != q.QuadToString() {
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
			status |= ConversionSyntax
		case "Division_by_zero":
			status |= DivisionByZero
		case "Division_impossible":
			status |= DivisionImpossible
		case "Division_undefined":
			status |= DivisionUndefined
		case "Insufficient_storage":
			status |= InsufficientStorage
		case "Inexact":
			status |= Inexact
		case "Invalid_context":
			status |= InvalidContext
		case "Invalid_operation":
			status |= InvalidOperation
		case "Overflow":
			status |= Overflow
		case "Clamped":
			// status |= Clamped
		case "Rounded":
			// status |= Rounded
		case "Subnormal":
			// status |= Subnormal
		case "Underflow":
			status |= Underflow
		default:
			panic("Unknown status flag: " + flag)
		}
	}

	return status
}
