#ifndef MYDECQUAD_H
#define MYDECQUAD_H

#include <errno.h>
#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <assert.h>
#include "decQuad.h"      // this header includes "decContext.h"
#include "decimal128.h"   // interface to decNumber, used for decNumberPower()


// error description is in "Decimal Arithmetic Specification, Exceptional conditions" at http://speleotrove.com/decimal/daexcep.html

#define MDQ_ERROR_INFINITE                     1     // result is Inf or -Inf
#define MDQ_ERROR_NAN                          2     // result is Nan
#define MDQ_ERROR_OVERFLOW                     3     // result doesn't fit in target decQuad of precision p
#define MDQ_ERROR_OUT_OF_RANGE                 4     // conversion between decQuad and other type failed because out of range
#define MDQ_ERROR_DEC_UNLISTED                 5     // decQuad error: a decNumer error occurred, but we haven't listed it in mdq_get_error()
#define MDQ_ERROR_DEC_INVALID_OPERATION        6     // decQuad error: result is Nan, for many invalid operations. E.g.  Inf*0 or Inf/Inf, etc.
#define MDQ_ERROR_DEC_DIVISION_BY_ZERO         7     // decQuad error: result is +-Inf, for division by 0
#define MDQ_ERROR_DEC_OVERFLOW                 8     // decQuad error: result is +-Inf, when exponent is too large. E.g.  1e6000 * 1e6000 = Inf
#define MDQ_ERROR_DEC_UNDERFLOW                9     // decQuad error: result is 0 or subnormal number close to 0. It occurs when result is subnormal and digits have been lost. E.g.  189e-6170 * 1e-7 = 19e-6176
#define MDQ_ERROR_DEC_DIVISION_IMPOSSIBLE     10     // decQuad error: result is Nan, for decQuadDivideInteger() or decQuadRemainder() is larger than an integral value with exponent 0.
#define MDQ_ERROR_DEC_DIVISION_UNDEFINED      11     // decQuad error: result is Nan for 0/0
#define MDQ_ERROR_DEC_CONVERSION_SYNTAX       12     // decQuad error: result is Nan, when conversion from string to number failed.
#define MDQ_ERROR_DEC_INSUFFICIENT_STORAGE    13     // decQuad error: insufficient storage
#define MDQ_ERROR_DEC_INVALID_CONTEXT         14     // decQuad error: invalid context


#define MYDECQUAD_Errors  (DEC_Errors & (~(DEC_Overflow | DEC_Underflow)))  // replace DEC_Errors for error checking, because we don't want to catch Overflow and Underflow. The operation can continue with result set to 0 or +-Inf.

#if MYDECQUAD_Errors != (DEC_Division_by_zero | DEC_Conversion_syntax | DEC_Division_impossible | DEC_Division_undefined | DEC_Insufficient_storage | DEC_Invalid_context | DEC_Invalid_operation)
  #error "MYDECQUAD_Errors is unexpected."
#endif


#define S_STRING_RAW_CAPACITY  (DECQUAD_Pmax + 20)    // more than enough to receive    sign + 34 digits + 'e' + exponent (int32_t) + '\0'


void mdq_init(void);

uint32_t mdq_adjust_p_s_and_check_error(decQuad *r, uint16_t precision, uint16_t scale, decContext *set);

uint32_t mdq_unary_minus(decQuad *r, uint16_t precision, uint16_t scale, decQuad *a);
uint32_t mdq_add(decQuad *r, uint16_t precision, uint16_t scale, decQuad *a, decQuad *b);
uint32_t mdq_subtract(decQuad *r, uint16_t precision, uint16_t scale, decQuad *a, decQuad *b);
uint32_t mdq_multiply(decQuad *r, uint16_t precision, uint16_t scale, decQuad *a, decQuad *b);
uint32_t mdq_divide(decQuad *r, uint16_t precision, uint16_t scale, decQuad *a, decQuad *b);
int32_t  mdq_compare(decQuad *a, decQuad *b);
int32_t  mdq_check_equality_FOR_TEST(decQuad *a, decQuad *b);

void     mdq_zero(decQuad *r, uint16_t precision, uint16_t scale);
uint32_t mdq_copy(decQuad *r, uint16_t precision, uint16_t scale, decQuad *a);
uint32_t mdq_abs(decQuad *r, uint16_t precision, uint16_t scale, decQuad *a);
uint32_t mdq_ceiling(decQuad *r, uint16_t precision, uint16_t scale, decQuad *a);
uint32_t mdq_floor(decQuad *r, uint16_t precision, uint16_t scale, decQuad *a);
uint32_t mdq_sign(decQuad *r, uint16_t precision, uint16_t scale, decQuad *a);
uint32_t mdq_power(decQuad *r, uint16_t precision, uint16_t scale, decQuad *a, decQuad *b);
uint32_t mdq_round(decQuad *r, uint16_t precision, uint16_t scale, decQuad *a, uint16_t b_precision, uint16_t b_scale, int32_t b, uint8_t truncate_flag);
uint32_t mdq_round_for_formatting(decQuad *r, decQuad *a, int32_t b);

uint32_t mdq_from_int32(decQuad *r, uint16_t precision, uint16_t scale, int32_t value);
uint32_t mdq_from_bytes_raw(decQuad *r, uint8_t *val, int32_t len);
uint32_t mdq_from_bytes(decQuad *r, uint16_t precision, uint16_t scale, uint8_t *val, int32_t len);
uint32_t mdq_from_bytes_with_implied_p_s(decQuad *r, uint16_t *out_precision, uint16_t *out_scale, uint8_t *val, int32_t len);

void     mdq_QuadToString(uint8_t *byte_array, int32_t capacity, decQuad *a);
uint32_t mdq_to_BCD(uint8_t *BCD_array, int32_t *exp, uint32_t *sign, decQuad *a);
int      mdq_to_string_raw(uint8_t *byte_array, int32_t capacity, decQuad *a);;
void     mdq_print_string_raw(const char *format, decQuad *a);
uint32_t mdq_to_int32_truncate(int32_t *dest, decQuad *a);
uint32_t mdq_to_int64_truncate(int64_t *dest, decQuad *a);
uint32_t mdq_to_int32_round(int32_t *dest, decQuad *a);
uint32_t mdq_to_int64_round(int64_t *dest, decQuad *a);
uint32_t mdq_to_double(double *dest, decQuad *a);


#endif

