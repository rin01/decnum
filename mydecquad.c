#include "mydecquad.h"


/************************************************************************/
/*                                xmalloc                               */
/************************************************************************/

inline static void *xmalloc (size_t size) {
  void *p;

  p = malloc(size);
  if ( p == NULL ) {
    fprintf(stderr, "malloc(%d) failed\n", (int)size);
    abort();
  }

  return p;
}


/************************************************************************/
/*                          init and context                            */
/************************************************************************/

/* initialize the global constants used by this library.

   Exit(1) if an error occurs.
*/
void mdq_init(void) {

  //----- check DECLITEND -----

  if ( decContextTestEndian(1) ) {  // if argument is 0, a warning message is displayed (using printf) if DECLITEND is set incorrectly. If 1, no message is displayed. Returns 0 if correct.
      fprintf(stderr, "INITIALIZATION mydecquad.c:mdq_init() FAILED: decnum: decContextTestEndian() failed. Change DECLITEND constant (see \"The decNumber Library\")");
      exit(1);
  }

  assert( DECQUAD_Pmax == 34 );             // we have 34 digits max precision.
  assert( DECQUAD_String > DECQUAD_Pmax );  // because Go function quad.AppendQuad()

}


decContext mdq_context_default(decContext set, uint32_t kind) {

  decContextDefault(&set, kind);

  return set;
}


int mdq_context_get_rounding(decContext set) {

  return decContextGetRounding(&set);
}


decContext mdq_context_set_rounding(decContext set, int rounding) {

  decContextSetRounding(&set, rounding);

  return set;
}


uint32_t mdq_context_get_status(decContext set) {

  return decContextGetStatus(&set);
}


decContext mdq_context_zero_status(decContext set) {

  decContextZeroStatus(&set);

  return set;
}


/************************************************************************/
/*                        arithmetic operations                         */
/************************************************************************/


decQuad mdq_zero() {
  decQuad  val;

  decQuadZero(&val);

  return val;
}


/* unary minus.
*/
Result_t mdq_minus(decQuad a, decContext set) {

  Result_t     res;

  /* operation */

  decQuadMinus(&res.val, &a, &set);
  res.set = set;

  return res;
}


/* addition.
*/
Result_t mdq_add(decQuad a, decQuad b, decContext set) {

  Result_t     res;

  /* operation */

  decQuadAdd(&res.val, &a, &b, &set);
  res.set = set;

  return res;
}


/* subtraction.
*/
Result_t mdq_subtract(decQuad a, decQuad b, decContext set) {

  Result_t     res;

  /* operation */

  decQuadSubtract(&res.val, &a, &b, &set);
  res.set = set;

  return res;
}


/* multiplication.
*/
Result_t mdq_multiply(decQuad a, decQuad b, decContext set) {

  Result_t     res;

  /* operation */

  decQuadMultiply(&res.val, &a, &b, &set);
  res.set = set;

  return res;
}


/* division.
*/
Result_t mdq_divide(decQuad a, decQuad b, decContext set) {

  Result_t     res;

  /* operation */

  decQuadDivide(&res.val, &a, &b, &set);
  res.set = set;

  return res;
}


/* integer division.
*/
Result_t mdq_divide_integer(decQuad a, decQuad b, decContext set) {

  Result_t     res;

  /* operation */

  decQuadDivideInteger(&res.val, &a, &b, &set);
  res.set = set;

  return res;
}


/* modulo.
*/
Result_t mdq_remainder(decQuad a, decQuad b, decContext set) {

  Result_t     res;

  /* operation */

  decQuadRemainder(&res.val, &a, &b, &set);
  res.set = set;

  return res;
}


/* absolute value.
*/
Result_t mdq_abs(decQuad a, decContext set) {

  Result_t     res;

  /* operation */

  decQuadAbs(&res.val, &a, &set);
  res.set = set;

  return res;
}


/* to integral.
*/
Result_t mdq_to_integral(decQuad a, decContext set, int round) {

  Result_t     res;

  /* operation */

  decQuadToIntegralValue(&res.val, &a, &set, round);
  res.set = set;

  return res;
}


/* quantize.
*/
Result_t mdq_quantize(decQuad a, decQuad b, decContext set) {

  Result_t     res;

  /* operation */

  decQuadQuantize(&res.val, &a, &b, &set);
  res.set = set;

  return res;
}


/* compare.
*/
Result_t mdq_compare(decQuad a, decQuad b, decContext set) {

  Result_t     res;

  /* operation */

  decQuadCompare(&res.val, &a, &b, &set);
  res.set = set;

  return res;
}


/* check if a is Finite number.
*/
uint32_t mdq_is_finite(decQuad a) {

  return decQuadIsFinite(&a);
}


/* check if a is integer number.
*/
uint32_t mdq_is_integer(decQuad a) {

  return decQuadIsInteger(&a);
}


/* check if a is Infinite.
*/
uint32_t mdq_is_infinite(decQuad a) {

  return decQuadIsInfinite(&a);
}


/* check if a is Nan.
*/
uint32_t mdq_is_nan(decQuad a) {

  return decQuadIsNaN(&a);
}


/* check if a is < 0 and not Nan.
*/
uint32_t mdq_is_negative(decQuad a) {

  return decQuadIsNegative(&a);
}


/* check if a is > 0 and not Nan.
*/
uint32_t mdq_is_positive(decQuad a) {

  return decQuadIsPositive(&a);
}


/* check if a is == 0.
*/
uint32_t mdq_is_zero(decQuad a) {

  return decQuadIsZero(&a);
}


/* max.
*/
Result_t mdq_max(decQuad a, decQuad b, decContext set) {

  Result_t     res;

  /* operation */

  decQuadMax(&res.val, &a, &b, &set);
  res.set = set;

  return res;
}


/* min.
*/
Result_t mdq_min(decQuad a, decQuad b, decContext set) {

  Result_t     res;

  /* operation */

  decQuadMin(&res.val, &a, &b, &set);
  res.set = set;

  return res;
}


/************************************************************************/
/*                        conversion to string                          */
/************************************************************************/


/* write decQuad into byte array.

   A terminating 0 is written in the array.
   Never fails.

   The function decQuadToString() uses exponential notation if number < 0 and too many 0 after decimal point.

   IMPORTANT: the caller must free the returned buffer when he is finished with it. Else, memory leaks occur.
*/
Ret_str mdq_to_mallocated_QuadToString(decQuad a) {

  Ret_str  ret = {.s = NULL, .length = 0};

  ret.s = (char *)xmalloc(DECQUAD_String);

  decQuadToString(&a, ret.s);

  ret.length = strlen(ret.s);

  return ret;
}


/* write decQuad into BCD_array.

   The returned fields are:
      BCD:       byte array. The coefficient is written one digit per byte.
      capacity:  size of BCD byte array (always DECQUAD_Pmax)
      exp:       if a is not Inf or Nan, will contain the exponent.
      sign:      if negative and not zero, sign bit is set.
                 THE SIGN IS VALID ALSO IF THE FUNCTION RETURNS MDQ_ERROR_INFINITE, so that we can know if it is +Inf or -Inf.

   Returns ret.mdqerr == 0 if success, or MDQ_ERROR_INFINITE or MDQ_ERROR_NAN.

   IMPORTANT: the caller must free the returned buffer when he is finished with it. Else, memory leaks occur.
*/
Ret_BCD mdq_to_mallocated_BCD(decQuad a) {

  int32_t     exp;
  uint32_t    sign;
  Ret_BCD     ret = {.mdqerr = 0, .BCD = NULL, .capacity = 0, .exp = 0, .sign = 0};

  ret.BCD = (char *)xmalloc(DECQUAD_Pmax);

  // convert to BCD

  decQuadToBCD(&a, &exp, ret.BCD);  // this function returns a sign bit, but we don't use it because we don't want -0

  sign = decQuadIsNegative(&a);     // 0 is never negative


  // check that result is not Inf nor Nan

  if ( ! decQuadIsFinite(&a) ) {
      if ( decQuadIsInfinite(&a) ) {
          ret.mdqerr = MDQ_ERROR_INFINITE;
      } else {
          ret.mdqerr = MDQ_ERROR_NAN;
      }
      return ret;
  }

  ret.capacity = DECQUAD_Pmax;
  ret.exp      = exp;
  ret.sign     = sign;

  return ret;
}


/************************************************************************/
/*                         conversion to numbers                        */
/************************************************************************/


/* convert decQuad to int32_t
*/
Ret_int32_t mdq_to_int32(decQuad a, decContext set, int round) {

  Ret_int32_t     res;

  /* operation */

  res.val = decQuadToInt32(&a, &set, round);
  res.set = set;

  return res;
}


/* convert decQuad to int64_t
*/
Ret_int64_t mdq_to_int64(decQuad a, decContext set, int round) {

  decQuad      zero;
  decQuad      a_integral;
  decQuad      a_integral_quantized;
  char         a_str[DECQUAD_String];
  char        *tailptr;
  int64_t      r_val;
  Ret_int64_t  ret;


  /* operation */

  decQuadZero(&zero);

  decQuadToIntegralValue(&a_integral, &a, &set, round);

  decQuadQuantize(&a_integral_quantized, &a_integral, &zero, &set); // because 1e3 remains 1e3

  if (set.status & DEC_Errors) {
    ret.set = set;
    ret.val = 0;
    return ret;
  }

  if (! decQuadIsFinite(&a_integral_quantized)) {
    decContextSetStatus(&set, DEC_Invalid_operation);
    ret.set = set;
    ret.val = 0;
    return ret;
  }

//  assert(decQuadGetExponent(&a_integral) == 0);

  decQuadToString(&a_integral_quantized, a_str);  // never raises error. Exponential notation never occurs for integral, which allows strtoll() to parse the number.
printf("xxxxxxxxxxxxxx  %s\n", a_str);

  errno = 0;
  r_val = strtoll(a_str, &tailptr, 10);  // changes errno if error

  if ( errno ) {
    decContextSetStatus(&set, DEC_Invalid_operation);
    ret.set = set;
    ret.val = 0;
    return ret;
  }

  if ( *tailptr != 0 ) { // may happen for e.g.  123e10, because it parses up to 'e'
    decContextSetStatus(&set, DEC_Invalid_operation);
    ret.set = set;
    ret.val = 0;
    return ret;
  }

  ret.set = set;
  ret.val = r_val;
  return ret;
}


/************************************************************************/
/*                    conversion from string or numbers                 */
/************************************************************************/


/* conversion from string.
*/
Result_t mdq_from_string(Strarray_t strarray, decContext set) {

  Result_t     res;

  /* operation */

  decQuadFromString(&res.val, strarray.arr, &set);
  res.set = set;

  return res;
}


/* conversion from int64.
*/
Result_t mdq_from_int64(int64_t value, decContext set) {

  char         buff[30]; // more than enough to store a int64
  Result_t     res;

  /* write value into buffer */

  sprintf(buff, "%lld", (long long int)value);


  /* operation */

  decQuadFromString(&res.val, buff, &set);            // raises an error if string is invalid
  res.set = set;

  return res;
}



