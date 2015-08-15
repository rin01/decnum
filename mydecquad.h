#ifndef MYDECQUAD_H
#define MYDECQUAD_H

#include <errno.h>
#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <assert.h>
#include "decQuad.h"      // this header includes "decContext.h"
#include "decimal128.h"   // interface to decNumber, used for decNumberPower()


#define MDQ_INFINITE                     1     // result is Inf or -Inf
#define MDQ_NAN                          2     // result is Nan


void mdq_init(void);

#define MAX_STRING_SIZE  50

// struct used to retrieve both a decQuad result and a decContext, from an operation.
// This way, the result of an operation (value and context) are returned to the caller as value.
// No need to fuss with pointers.
//
typedef struct Ret_decQuad_t {
  decContext  set;
  decQuad     val;
} Ret_decQuad_t;

// struct used to pass strings from Go to C and vice-versa.
// This way, strings are just passed as value, no need to fuss with pointers.
//
typedef struct Strarray_t {
  char arr[MAX_STRING_SIZE];
} Strarray_t;

typedef struct Ret_BCD {
  uint32_t   inf_nan;
  char      *BCD;
  size_t     capacity;
  int32_t    exp;
  uint32_t   sign;
} Ret_BCD;

typedef struct Ret_str {
  char      *s;
  size_t     length;
} Ret_str;

typedef struct Ret_int32_t {
  decContext  set;
  int32_t     val;
} Ret_int32_t;

typedef struct Ret_int64_t {
  decContext  set;
  int64_t     val;
} Ret_int64_t;


decContext       mdq_context_default(decContext set, uint32_t kind);
int              mdq_context_get_rounding(decContext set);
decContext       mdq_context_set_rounding(decContext set, int rounding);
uint32_t         mdq_context_get_status(decContext set);
decContext       mdq_context_zero_status(decContext set);

decQuad          mdq_zero();
Ret_decQuad_t    mdq_minus(decQuad a, decContext set);
Ret_decQuad_t    mdq_add(decQuad a, decQuad b, decContext set);
Ret_decQuad_t    mdq_subtract(decQuad a, decQuad b, decContext set);
Ret_decQuad_t    mdq_multiply(decQuad a, decQuad b, decContext set);
Ret_decQuad_t    mdq_divide(decQuad a, decQuad b, decContext set);
Ret_decQuad_t    mdq_divide_integer(decQuad a, decQuad b, decContext set);
Ret_decQuad_t    mdq_remainder(decQuad a, decQuad b, decContext set);
Ret_decQuad_t    mdq_abs(decQuad a, decContext set);
Ret_decQuad_t    mdq_to_integral(decQuad a, decContext set, int round);
Ret_decQuad_t    mdq_quantize(decQuad a, decQuad b, decContext set);
Ret_decQuad_t    mdq_compare(decQuad a, decQuad b, decContext set);
uint32_t         mdq_is_finite(decQuad a);
uint32_t         mdq_is_integer(decQuad a);
uint32_t         mdq_is_infinite(decQuad a);
uint32_t         mdq_is_nan(decQuad a);
uint32_t         mdq_is_negative(decQuad a);
uint32_t         mdq_is_positive(decQuad a);
uint32_t         mdq_is_zero(decQuad a);
Ret_decQuad_t    mdq_max(decQuad a, decQuad b, decContext set);
Ret_decQuad_t    mdq_min(decQuad a, decQuad b, decContext set);

Ret_int32_t      mdq_to_int32(decQuad a, decContext set, int round);
Ret_int64_t      mdq_to_int64(decQuad a, decContext set, int round);
Ret_str          mdq_to_mallocated_QuadToString(decQuad a);
Ret_BCD          mdq_to_mallocated_BCD(decQuad a);

Ret_decQuad_t    mdq_from_string(Strarray_t strarray, decContext set);
Ret_decQuad_t    mdq_from_int64(int64_t value, decContext set);

#endif
