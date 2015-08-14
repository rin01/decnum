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


// struct for return values

typedef struct Ret_decQuad {
  uint32_t   mdqerr;
  decQuad    val;
} Ret_decQuad;

typedef struct Ret_decQuad_p_s {
  uint32_t   mdqerr;
  decQuad    val;
  uint16_t   precision;
  uint16_t   scale;
} Ret_decQuad_p_s;

typedef struct Ret_str {
  char      *s;
  size_t     length;
} Ret_str;

typedef struct Ret_BCD {
  uint32_t   mdqerr;
  char      *BCD;
  size_t     capacity;
  int32_t    exp;
  uint32_t   sign;
} Ret_BCD;

typedef struct Ret_int32 {
  uint32_t   mdqerr;
  int32_t    val;
} Ret_int32;

typedef struct Ret_int64 {
  uint32_t   mdqerr;
  int64_t    val;
} Ret_int64;

typedef struct Ret_double {
  uint32_t   mdqerr;
  double     val;
} Ret_double;


//-------


void mdq_init(void);

uint32_t        mdq_hash(decQuad a);

Ret_decQuad     mdq_subtract(uint16_t precision, uint16_t scale, decQuad a, decQuad b);
Ret_decQuad     mdq_multiply(uint16_t precision, uint16_t scale, decQuad a, decQuad b);
Ret_decQuad     mdq_divide(uint16_t precision, uint16_t scale, decQuad a, decQuad b);
int32_t         mdq_compare(decQuad a, decQuad b);
int32_t         mdq_check_equality_FOR_TEST(decQuad a, decQuad b);

decQuad         mdq_zero(uint16_t precision, uint16_t scale);
Ret_decQuad     mdq_copy(uint16_t precision, uint16_t scale, decQuad a);
Ret_decQuad     mdq_abs(uint16_t precision, uint16_t scale, decQuad a);
Ret_decQuad     mdq_ceiling(uint16_t precision, uint16_t scale, decQuad a);
Ret_decQuad     mdq_floor(uint16_t precision, uint16_t scale, decQuad a);
Ret_decQuad     mdq_sign(uint16_t precision, uint16_t scale, decQuad a);
Ret_decQuad     mdq_power(uint16_t precision, uint16_t scale, decQuad a, double b);
Ret_decQuad     mdq_round(uint16_t precision, uint16_t scale, decQuad a, uint16_t a_precision, uint16_t a_scale, int32_t b, uint8_t truncate_flag);
Ret_decQuad     mdq_round_for_formatting(decQuad a, int32_t b);

Ret_decQuad     mdq_from_double_raw(double value);
Ret_decQuad     mdq_from_double(uint16_t precision, uint16_t scale, double value);

Ret_decQuad     mdq_from_bytes_raw_and_free(char *s);
Ret_decQuad     mdq_from_bytes_and_free(uint16_t precision, uint16_t scale, char *s);
Ret_decQuad_p_s mdq_from_bytes_with_implied_p_s_and_free(char *s);

Ret_str         mdq_to_mallocated_QuadToString(decQuad a);
Ret_BCD         mdq_to_mallocated_BCD(decQuad a);
Ret_str         mdq_to_mallocated_string_raw(decQuad a);
void            mdq_print_string_raw(const char *format, decQuad a);

Ret_int32       mdq_to_int32_truncate(decQuad a);
Ret_int32       mdq_to_int32_round(decQuad a);
Ret_int64       mdq_to_int64_truncate(decQuad a);
Ret_int64       mdq_to_int64_round(decQuad a);
Ret_double      mdq_to_double(decQuad a);

decQuad         mdq_decQuadZero(decQuad a);
uint32_t        mdq_decQuadIsZero(decQuad a);
uint32_t        mdq_decQuadIsNegative(decQuad a);






typedef struct Result_t {
  decContext  set;
  decQuad     val;
} Result_t;

Result_t mdq_unary_minus(decQuad a, decContext set);
Result_t mdq_add(decQuad a, decQuad b, decContext set);

Result_t mdq_from_int64(int64_t value, decContext set);

#endif

