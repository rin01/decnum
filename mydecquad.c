#include "mydecquad.h"


/************************************************************************/
/*            global constant variables for Round and Truncate          */
/************************************************************************/

/* decNumber constants for decQuad rounding.

   It contains
      - 1e0
      - 1e-1
      - 1e-2
      - ...
      - 1e-DECQUAD_Pmax        (1e-34)
*/
static decQuad G_DECQUAD_QUANTIZER[DECQUAD_Pmax+1];  // 0...34


/* decNumber constants for decQuad rounding.

   It contains
      - 1e0
      - 1e1
      - 1e2
      - ...
      - 1eDECQUAD_Pmax         (1e34)
      - 1e(DECQUAD_Pmax+1)     (1e35)

    NOTE: the max index is 35, because rounding      9234567890123456789012345678901234
                                           with     10000000000000000000000000000000000    (1e34)
                                           gives    10000000000000000000000000000000000

                                   and rounding      9234567890123456789012345678901234
                                           with    100000000000000000000000000000000000    (1e35)
                                           gives                                      0

                                   So, we must allow rounding functions to use an integral part quantizer of 1e35.
*/
static decQuad G_DECQUAD_INTEGRAL_PART_QUANTIZER[DECQUAD_Pmax+2];  // 0...35


/************************************************************************/
/*                          init and context                            */
/************************************************************************/

static decQuad static_one;  // contains 1, only used by mdq_to_int64


/* initialize the global constants used by this library.

   It is called by Go in init() function.

   Exit(1) if an error occurs.
*/
void mdq_init(void) {

  decContext   set;
  const char  *s;
  int          i;


  //----- check DECLITEND -----

  if ( decContextTestEndian(1) ) {  // if argument is 0, a warning message is displayed (using printf) if DECLITEND is set incorrectly. If 1, no message is displayed. Returns 0 if correct.
      fprintf(stderr, "INITIALIZATION mydecquad.c:mdq_init() FAILED: decnum: decContextTestEndian() failed. Change DECLITEND constant (see \"The decNumber Library\")");
      exit(1);
  }

  assert( DECQUAD_Pmax == 34 );             // we have 34 digits max precision.
  assert( DECQUAD_String > DECQUAD_Pmax );  // because Go function quad.AppendQuad()


  //----- put 1 in static_one -----

  decQuadFromInt32(&static_one, 1); // IMPORTANT: this means that mdq_to_int64 can only be called after Go init() has been run, as it uses static_one. ctx.ToInt32() cannot be called to initialize Go global variables.


  //----- fill decContext -----

  decContextDefault(&set, DEC_INIT_DECQUAD);

  if ( decContextGetRounding(&set) != DEC_ROUND_HALF_EVEN ) {
      fprintf(stderr, "INITIALIZATION mydecquad.c:mdq_init() FAILED: decnum: decContextGetRounding(&set) != DEC_ROUND_HALF_EVEN");
      exit(1);
  }


  //----- fill G_DECQUAD_QUANTIZER[] -----

  decQuadFromInt32(&G_DECQUAD_QUANTIZER[0], 1);                       //  store  1e0  in G_DECQUAD_QUANTIZER[0]

  assert( decQuadDigits(     &G_DECQUAD_QUANTIZER[0]) == 1 );
  assert( decQuadGetExponent(&G_DECQUAD_QUANTIZER[0]) == 0 );

  for ( i=1; i<=DECQUAD_Pmax; i++ ) {                                 // in G_DECQUAD_QUANTIZER[1..DECQUAD_Pmax]
      decQuadCopy(&G_DECQUAD_QUANTIZER[i], &G_DECQUAD_QUANTIZER[0]);

      decQuadSetExponent(&G_DECQUAD_QUANTIZER[i], &set, -i);          // store 1e-1 .. 1e-DECQUAD_Pmax
  }

  assert( decQuadDigits(     &G_DECQUAD_QUANTIZER[DECQUAD_Pmax]) == 1 );
  assert( decQuadGetExponent(&G_DECQUAD_QUANTIZER[DECQUAD_Pmax]) == -DECQUAD_Pmax );  // -34


  //----- fill G_DECQUAD_INTEGRAL_PART_QUANTIZER[] -----

  decQuadFromInt32(&G_DECQUAD_INTEGRAL_PART_QUANTIZER[0], 1);              //  store  1e0  in G_DECQUAD_INTEGRAL_PART_QUANTIZER[0]

  assert( decQuadDigits(     &G_DECQUAD_INTEGRAL_PART_QUANTIZER[0]) == 1 );
  assert( decQuadGetExponent(&G_DECQUAD_INTEGRAL_PART_QUANTIZER[0]) == 0 );

  for ( i=1; i<=DECQUAD_Pmax+1; i++ ) {                                    // in G_DECQUAD_INTEGRAL_PART_QUANTIZER[1..DECQUAD_Pmax+1]
      decQuadCopy(&G_DECQUAD_INTEGRAL_PART_QUANTIZER[i], &G_DECQUAD_INTEGRAL_PART_QUANTIZER[0]);

      decQuadSetExponent(&G_DECQUAD_INTEGRAL_PART_QUANTIZER[i], &set, i);  // store 1e1 .. 1e(DECQUAD_Pmax+1)
  }

  assert( decQuadDigits(     &G_DECQUAD_INTEGRAL_PART_QUANTIZER[DECQUAD_Pmax])   == 1 );
  assert( decQuadGetExponent(&G_DECQUAD_INTEGRAL_PART_QUANTIZER[DECQUAD_Pmax])   == DECQUAD_Pmax   );  // 34

  assert( decQuadDigits(     &G_DECQUAD_INTEGRAL_PART_QUANTIZER[DECQUAD_Pmax+1]) == 1 );
  assert( decQuadGetExponent(&G_DECQUAD_INTEGRAL_PART_QUANTIZER[DECQUAD_Pmax+1]) == DECQUAD_Pmax+1 );  // 35


  //----- check for errors or any warning -----

  if ( set.status ) {
      s = decContextStatusToString(&set);
      fprintf(stderr, "INITIALIZATION mydecquad.c:mdq_init() FAILED: decNumber quantizer initialization failed. %s\n", s);
      exit(1);
  }

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


decContext mdq_context_set_status(decContext set, uint32_t flag) {

  decContextSetStatus(&set, flag);

  return set;
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


decQuad mdq_nan() {
  decContext set;
  decQuad    val;

  decContextDefault(&set, DEC_INIT_DECQUAD);

  decQuadFromString(&val, "Nan", &set);

  //assert(set.status & DEC_Errors == 0); // a status bit is set, because the Nan

  return val;
}


/* unary minus.
*/
Ret_decQuad_t mdq_minus(decQuad a, decContext set) {

  Ret_decQuad_t     res;

  /* operation */

  decQuadMinus(&res.val, &a, &set);
  res.set = set;

  return res;
}


/* addition.
*/
Ret_decQuad_t mdq_add(decQuad a, decQuad b, decContext set) {

  Ret_decQuad_t     res;

  /* operation */

  decQuadAdd(&res.val, &a, &b, &set);
  res.set = set;

  return res;
}


/* subtraction.
*/
Ret_decQuad_t mdq_subtract(decQuad a, decQuad b, decContext set) {

  Ret_decQuad_t     res;

  /* operation */

  decQuadSubtract(&res.val, &a, &b, &set);
  res.set = set;

  return res;
}


/* multiplication.
*/
Ret_decQuad_t mdq_multiply(decQuad a, decQuad b, decContext set) {

  Ret_decQuad_t     res;

  /* operation */

  decQuadMultiply(&res.val, &a, &b, &set);
  res.set = set;

  return res;
}


/* division.
*/
Ret_decQuad_t mdq_divide(decQuad a, decQuad b, decContext set) {

  Ret_decQuad_t     res;

  /* operation */

  decQuadDivide(&res.val, &a, &b, &set);
  res.set = set;

  return res;
}


/* integer division.
*/
Ret_decQuad_t mdq_divide_integer(decQuad a, decQuad b, decContext set) {

  Ret_decQuad_t     res;

  /* operation */

  decQuadDivideInteger(&res.val, &a, &b, &set);
  res.set = set;

  return res;
}


/* modulo.
*/
Ret_decQuad_t mdq_remainder(decQuad a, decQuad b, decContext set) {

  Ret_decQuad_t     res;

  /* operation */

  decQuadRemainder(&res.val, &a, &b, &set);
  res.set = set;

  return res;
}


/* absolute value.
*/
Ret_decQuad_t mdq_abs(decQuad a, decContext set) {

  Ret_decQuad_t     res;

  /* operation */

  decQuadAbs(&res.val, &a, &set);
  res.set = set;

  return res;
}


/* to integral.
*/
Ret_decQuad_t mdq_to_integral(decQuad a, decContext set, int round) {

  Ret_decQuad_t     res;

  /* operation */

  decQuadToIntegralValue(&res.val, &a, &set, round);
  res.set = set;

  return res;
}


/* quantize.
*/
Ret_decQuad_t mdq_quantize(decQuad a, decQuad b, decContext set) {

  Ret_decQuad_t     res;

  /* operation */

  decQuadQuantize(&res.val, &a, &b, &set);
  res.set = set;

  return res;
}


/* compare.
*/
Ret_uint32_t mdq_compare(decQuad a, decQuad b, decContext set) {

  decQuad         cmp_val;
  Ret_uint32_t    res;

  /* operation */

  decQuadCompare(&cmp_val, &a, &b, &set);
  res.set = set;

  if ( decQuadIsNaN(&cmp_val) ) {
      res.val = CMP_NAN;
      return res;
  }

  if ( decQuadIsZero(&cmp_val) ) {
      res.val = CMP_EQUAL;
      return res;
  }

  if ( decQuadIsPositive(&cmp_val) ) {
      res.val = CMP_GREATER;
      return res;
  }

  assert( decQuadIsNegative(&cmp_val) );

  res.val = CMP_LESS;
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


/* check if a is < 0 and not Nan.
*/
uint32_t mdq_is_negative(decQuad a) {

  return decQuadIsNegative(&a);
}

/* get exponent.
*/
int32_t mdq_get_exponent(decQuad a) {

  return decQuadGetExponent(&a);
}


/* max.
*/
Ret_decQuad_t mdq_max(decQuad a, decQuad b, decContext set) {

  Ret_decQuad_t     res;

  /* operation */

  decQuadMax(&res.val, &a, &b, &set);
  res.set = set;

  return res;
}


/* min.
*/
Ret_decQuad_t mdq_min(decQuad a, decQuad b, decContext set) {

  Ret_decQuad_t     res;

  /* operation */

  decQuadMin(&res.val, &a, &b, &set);
  res.set = set;

  return res;
}


/************************************************************************/
/*                    conversion from string or numbers                 */
/************************************************************************/


/* conversion from string.
*/
Ret_decQuad_t mdq_from_string(char *s, decContext set) {

  Ret_decQuad_t     res;

  /* operation */

  decQuadFromString(&res.val, s, &set);
  res.set = set;

  return res;
}


/* conversion from int32.
*/
Ret_decQuad_t mdq_from_int32(int32_t value, decContext set) {

  Ret_decQuad_t     res;

  /* operation */

  decQuadFromInt32(&res.val, value); // decQuadFromInt32 doesn't need context, but conversion from string or int64 need it, so I do the same for int32
  res.set = set;

  return res;
}


/* conversion from int64.
*/
Ret_decQuad_t mdq_from_int64(int64_t value, decContext set) {

  char         buff[30]; // more than enough to store a int64     max val: 9,223,372,036,854,775,807
  Ret_decQuad_t     res;

  /* write value into buffer */

  sprintf(buff, "%lld", (long long int)value);

  /* operation */

  decQuadFromString(&res.val, buff, &set);            // raises an error if string is invalid
  res.set = set;

  return res;
}


/* conversion from double.

   DEPRECATED: FromFloat64 function has been removed, because it is impossible to know the desired precision of the result.
               The user should convert float64 to string, with the desired precision, and pass it to FromString.

*/
Ret_decQuad_t mdq_from_double(double value, decContext set) {

  char         buff[40]; // more than enough to store a double in the format specified by sprintf
  Ret_decQuad_t     res;

  /* write value into buffer */

  sprintf(buff, "%.18e", value);
  //printf("mdq_from_double: %s\n", buff);

  /* operation */

  decQuadFromString(&res.val, buff, &set);            // raises an error if string is invalid
  res.set = set;

  return res;
}


/************************************************************************/
/*                        conversion to string                          */
/************************************************************************/


/* write decQuad into byte array.

   A terminating 0 is written in the array.
   Never fails.

   The function decQuadToString() uses exponential notation too often in my opinion. E.g. 0.0000001 returns "1E-7".
*/
Ret_str mdq_to_QuadToString(decQuad a) {

  Ret_str  res = {.length = 0};

  decQuadToString(&a, res.s);

  res.length = strlen(res.s);

  return res;
}


/* write decQuad into BCD_array.

   The returned fields are:
      BCD:       byte array. The coefficient is written one digit per byte.
      exp:       if a is not Inf or Nan, will contain the exponent.
      sign:      if negative and not zero, sign bit is set.
                 THE SIGN IS VALID ALSO IF THE FUNCTION RETURNS MDQ_INFINITE, so that we can know if it is +Inf or -Inf.
*/
Ret_BCD mdq_to_BCD(decQuad a) {

  int32_t     exp;
  uint32_t    sign;
  Ret_BCD     res = {.inf_nan = 0, .exp = 0, .sign = 0};

  // convert to BCD

  decQuadToBCD(&a, &exp, res.BCD);  // this function returns a sign bit, but we don't use it because we don't want -0

  sign = decQuadIsNegative(&a);     // 0 is never negative


  // check that result is not Inf nor Nan

  if ( ! decQuadIsFinite(&a) ) {
      if ( decQuadIsInfinite(&a) ) {
          res.inf_nan = MDQ_INFINITE;
      } else {
          res.inf_nan = MDQ_NAN;
      }
      return res;
  }

  res.exp      = exp;
  res.sign     = sign;

  return res;
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

  decQuad      a_integral;
  decQuad      a_integral_quantized;
  char         a_str[DECQUAD_String];
  char        *tailptr;
  int64_t      r_val;
  Ret_int64_t  res;


  /* operation */

  decQuadToIntegralValue(&a_integral, &a, &set, round); // rounds the number to an integral. Only numbers with exponent<0 are rounded and shifted so that exponent becomes 0.

  decQuadQuantize(&a_integral_quantized, &a_integral, &static_one, &set); // for numbers with exponent>0. E.g. change 1e3 to 1000

  if (set.status & DEC_Errors) {
    res.set = set;
    res.val = 0;
    return res;
  }

  if (! decQuadIsFinite(&a_integral_quantized)) {
    decContextSetStatus(&set, DEC_Invalid_operation);
    res.set = set;
    res.val = 0;
    return res;
  }

  assert(decQuadGetExponent(&a_integral_quantized) == 0); // in the absence of decQuadQuantize error, the exponent of the result is always equal to that of the model 'static_one'

  decQuadToString(&a_integral_quantized, a_str);  // never raises error. Exponential notation never occurs for integral, which allows strtoll() to parse the number.
  //printf("xxxxxxxxxxxxxx  %s\n", a_str);

  errno = 0;
  r_val = strtoll(a_str, &tailptr, 10);  // changes errno if error

  if ( errno ) { // in particular, if a_str is an integer that overflows int64
    decContextSetStatus(&set, DEC_Invalid_operation);
    res.set = set;
    res.val = 0;
    return res;
  }

  if ( *tailptr != 0 ) { // may happen for e.g.  123e10, because it parses up to 'e'
    decContextSetStatus(&set, DEC_Invalid_operation);
    res.set = set;
    res.val = 0;
    return res;
  }

  res.set = set;
  res.val = r_val;
  return res;
}


/************************************************************************/
/*                       rounding and truncating                        */
/************************************************************************/

// rounds a. If rounding < 0, uses context rounding mode.
// Note: other rounding modes like DEC_ROUND_CEILING, etc, are enum, positive values.
//
Ret_decQuad_t mdq_roundM(decQuad a, int32_t n, int rounding, decContext set) {

  Ret_decQuad_t     res;
  decQuad           r;
  decQuad          *operation_quantizer;
  int               rounding_bak;


  /* if n is out-of-range, return Invalid_operation */

  if ( n > 34 || n < -35 ) {
      decContextSetStatus(&set, DEC_Invalid_operation); // add flag to status

      res.val = mdq_nan();
      res.set = set;
      return res;
  }


  /* operation */

  if ( rounding >= 0 ) { // if specific rounding
      rounding_bak = decContextGetRounding(&set);
      decContextSetRounding(&set, rounding);                           // change rounding mode
  }

  if ( n >= 0 ) {   // round or truncate fractional part
      operation_quantizer = &G_DECQUAD_QUANTIZER[n];                   // n is [0..34]

      decQuadQuantize(&res.val, &a, operation_quantizer, &set);        // rounding, e.g. quaantize(1234.5678, 2)  --> 1234.57

  } else {          // n < 0, round or truncate integral part
      operation_quantizer = &G_DECQUAD_INTEGRAL_PART_QUANTIZER[-n];    // -n is [0..35]

      decQuadQuantize(&r, &a, operation_quantizer, &set);              // rounding, e.g. quaantize(1234.5678, -2) --> 12E2
      decQuadQuantize(&res.val, &r, &G_DECQUAD_QUANTIZER[0], &set);    // right-shift the number, adding missing 0s on the left. E.g. 12E2 --> 1200E0
  }

  if ( rounding >= 0 ) { // if specific rounding
      decContextSetRounding(&set, rounding_bak);                       // restore original rounding
  }

  res.set = set;

  return res;
}


