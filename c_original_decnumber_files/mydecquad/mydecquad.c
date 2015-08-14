#include "mydecquad.h"


/****===========================================================================****
 ****                                                                           ****
 **** IMOPRTANT: if a function in this module returns an error code,            ****
 ****            the result is undefined and should not be used.                ****
 ****                                                                           ****
 ****===========================================================================****/


/* decNumber constants for decQuad quantization.

   It contains
      - 1e0
      - 1e-1
      - 1e-2
      - ...
      - 1e-DECQUAD_Pmax        (1e-34)
*/
static decQuad G_DECQUAD_QUANTIZER[DECQUAD_Pmax+1];  // 0...34


/* decNumber constants for decQuad rounding, used by data.Sysfunc_round_NUMERIC().

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

                                   So, we must allow data.Sysfunc_round_NUMERIC to use an integral part quantizer of 1e35.
*/
static decQuad G_DECQUAD_INTEGRAL_PART_QUANTIZER[DECQUAD_Pmax+2];  // 0...35


/************************************************************************/
/*                                                                      */
/*                            init function                             */
/*                                                                      */
/************************************************************************/


/* initialize the global constants used by this library.

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

  assert( DECQUAD_Pmax == 34 );             // for NUMERIC, we have 34 digits max precision.
  assert( DECQUAD_String > DECQUAD_Pmax );  // because Go function quad.AppendQuad()


  //----- print decNumber settings -----

  fprintf(stderr, "info: decNumber settings are DECDPUN %d, DECSUBSET %d, DECEXTFLAG %d. Constants DECQUAD_Pmax %d, DECQUAD_String %d.\n", DECDPUN, DECSUBSET, DECEXTFLAG, DECQUAD_Pmax, DECQUAD_String);


  //----- fill decContext -----

  decContextDefault(&set, DEC_INIT_DECQUAD);
  decContextSetRounding(&set, DEC_ROUND_HALF_UP);

  if ( decContextGetRounding(&set) != DEC_ROUND_HALF_UP ) {
      fprintf(stderr, "INITIALIZATION mydecquad.c:mdq_init() FAILED: decnum: decContextGetRounding(&set) != DEC_ROUND_HALF_UP");
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


/************************************************************************/
/*                                                                      */
/*                        error check functions                         */
/*                                                                      */
/************************************************************************/


/* ****===========================================================================================================****
   **** IMPORTANT: In our code, Overflow and Underflow are not considered as errors.                              ****
   ****            Operations that raise Overflow set result to +-Inf,                                            ****
   ****            and operations that raise Underflow set result to 0 or a subnormal number close to 0.          ****
   ****            For Underflow, we accept 0 or number close to 0 as valid result.                               ****
   ****            For Overflow, +-Inf is catched by the code that checks for Nan or Inf.                         ****
   ****                                                                                                           ****
   ****            We consider Overflow and Underflow like information flags.                                     ****
   ****===========================================================================================================****
*/


/* translate decNumber error code to mydecquad error code.

   This code is a modified copy of decContextStatusToString() in decContext.c.

   Only error flags in status are translated. Information flags as DEC_Inexact are skipped.

   Example: DEC_Invalid_operation error is translated to MDQ_ERROR_DEC_INVALID_OPERATION.

   Even if 'status' may contain multiple error flags which are ORed, this function only returns the most explicit error.

   Status MUST CONTAIN AN ERROR CODE (status & DEC_Errors != 0)
*/
static uint32_t mdq_get_status_error(uint32_t status) {

  status = (status & DEC_Errors);  // keep only real errors, not information flags

  assert( status != 0 );


  // the most explicit error are put first

  if (status & DEC_Division_by_zero     ) return MDQ_ERROR_DEC_DIVISION_BY_ZERO;     // error: division by zero
  if (status & DEC_Overflow             ) return MDQ_ERROR_DEC_OVERFLOW;             // error: result exponent is too large. E.g.  1e6000 * 1e6000 = Inf
  if (status & DEC_Underflow            ) return MDQ_ERROR_DEC_UNDERFLOW;            // error: result is subnormal and digits have been lost. E.g.  189e-6170 * 1e-7 = 19e-6176
  if (status & DEC_Conversion_syntax    ) return MDQ_ERROR_DEC_CONVERSION_SYNTAX;    // error: conversion string to decNumber
//  if (status & DEC_Inexact              ) return MDQ_DEC_INEXACT;                  // info flag

  if (status & DEC_Division_impossible  ) return MDQ_ERROR_DEC_DIVISION_IMPOSSIBLE;  // error: result of decQuadDivideInteger() or decQuadRemainder() is larger than an integral value with exponent 0.
  if (status & DEC_Division_undefined   ) return MDQ_ERROR_DEC_DIVISION_UNDEFINED;   // error: 0/0
  if (status & DEC_Invalid_operation    ) return MDQ_ERROR_DEC_INVALID_OPERATION;    // error: e.g. Inf*0 or Inf/Inf
//  if (status & DEC_Rounded              ) return MDQ_DEC_ROUNDED;                  // info flag
//  if (status & DEC_Clamped              ) return MDQ_DEC_CLAMPED;                  // info flag
//  if (status & DEC_Subnormal            ) return MDQ_DEC_SUBNORMAL;                // info flag
  if (status & DEC_Insufficient_storage ) return MDQ_ERROR_DEC_INSUFFICIENT_STORAGE;
  if (status & DEC_Invalid_context      ) return MDQ_ERROR_DEC_INVALID_CONTEXT;
  #if DECSUBSET
//  if (status & DEC_Lost_digits          ) return MDQ_DEC_LOST_DIGITS;              // info flag
  #endif

  return MDQ_ERROR_DEC_UNLISTED;
}


/* This function checks that no error occured in decContext, and that result is finite (that is, not Inf nor Nan).

   It checks that : 
      - status has no decNumber error bit
      - result is finite number (not Inf nor Nan)

   Returns 0 if success.
   Else, the returned value indicates a real error :
      - MDQ_ERROR_INFINITE
      - MDQ_ERROR_NAN
      - MDQ_ERROR_DEC_DIVISION_BY_ZERO, MDQ_ERROR_DEC_INVALID_OPERATION, etc
      - MDQ_ERROR_DEC_UNLISTED (should not happen)
*/
uint32_t mdq_check_error(decQuad *r, decContext *set) {

  /* check for decNumber errors so far */

  if ( set->status & MYDECQUAD_Errors ) {          // check only real errors. DEC_Inexact and DEC_Rounded are filtered out. Overflow and Underflow are not considered as errors.
      return mdq_get_status_error(set->status);        // translate real decNumber error code into MDQ_ERROR_XXX error code.
  }


  /* result must be a finite number ( not Inf nor Nan ) */

  if ( ! decQuadIsFinite(r) ) {                    // +Inf -Inf Nan are forbidden in the result.
      if ( decQuadIsInfinite(r) ) {
          return MDQ_ERROR_INFINITE;
      } else {
          return MDQ_ERROR_NAN;
      }
  }

  return 0;  // no error
}


/* This function adjusts the scale of the result, ensures that no error occured, that result is finite, and does not overflow the target precision.

   It checks that : 
      - status has no decNumber error bit
      - result is finite number (not Inf nor Nan)
      - result number of digits fits in the precision of the NUMERIC target

   Returns 0 if success.
   Else, the returned value indicates a real error :
      - MDQ_ERROR_INFINITE
      - MDQ_ERROR_NAN
      - MDQ_ERROR_OVERFLOW
      - MDQ_ERROR_DEC_DIVISION_BY_ZERO, MDQ_ERROR_DEC_INVALID_OPERATION, etc
      - MDQ_ERROR_DEC_UNLISTED (should not happen)
*/
uint32_t mdq_adjust_p_s_and_check_error(decQuad *r, uint16_t precision, uint16_t scale, decContext *set) {

  uint32_t    r_nb_of_digits;

  assert(scale <= DECQUAD_Pmax);


  /* check for decNumber errors so far */

  if ( set->status & MYDECQUAD_Errors) {           // check only real errors. DEC_Inexact and DEC_Rounded are filtered out.
      return mdq_get_status_error(set->status);        // translate real decNumber error code into MDQ_ERROR_XXX error code.
  }


  /* r must be a finite number ( not Inf nor Nan ) */

  if ( ! decQuadIsFinite(r) ) {                    // +Inf -Inf Nan are forbidden in the result.
      if ( decQuadIsInfinite(r) ) {
          return MDQ_ERROR_INFINITE;
      } else {
          return MDQ_ERROR_NAN;
      }
  }


  /* adjust the scale of the result */

  decQuadQuantize(r, r, &G_DECQUAD_QUANTIZER[scale], set);  // number r is rounded if necessary. If no error, the exponent of the result is always equal to that of the rhs (right-hand-side operand).
                                                                // if number is too large to be "flattened" to the given scale, or is +-Inf, DEC_Invalid_operation error occurs.
  if ( set->status & MYDECQUAD_Errors) {                    // check only real errors. DEC_Inexact and DEC_Rounded are filtered out.
      return MDQ_ERROR_OVERFLOW;                                // if finite number cannot be quantized, it is an overflow
  }


  /* result must be a finite number ( not Inf nor Nan ) */

  // documentation of decNumberQuantize() says: If adjusting the exponent would mean that more than context.digits would be needed in the coefficient, then the DEC_Invalid_operation condition is raised.
  //                                            This guarantees that in the absence of error the exponent of number is always equal to that of the rhs.
  //                                            If either operand is a special value (that is, Nan or Inf) then the usual rules apply [...]
  // As we catch Inf and Nan before quantizing, we can be sure that if decQuadQuantize() doesn't raise DEC_Invalid_operation, it has succeeded.

  assert(decQuadIsFinite(r));


  /* precision of the result should not exceed precision of the target.
       ( if result is 0, decQuadDigits() returns 1. For 0.0001, it returns 1. )

     pppppp pp
          0.00         0e-2       decQuadDigits() = 1
          0.01         1e-2       decQuadDigits() = 1
          0.08         8e-2       decQuadDigits() = 1
          0.12        12e-2       decQuadDigits() = 2
          1.12       112e-2       decQuadDigits() = 3
       1000.25    100025e-2       decQuadDigits() = 6
       1000.00    100000e-2       decQuadDigits() = 6

     We see that decQuadDigits() is the number of significant digits of a number.

     It must fit in the 'precision' number of digits available in the numeric('precision', 'scale') datatype of the target.
  */

  r_nb_of_digits = decQuadDigits(r);

  if ( r_nb_of_digits > precision ) {
      return MDQ_ERROR_OVERFLOW;
  }


  return 0;  // no error
}


/************************************************************************/
/*                                                                      */
/*                      arithmetic operations                           */
/*                                                                      */
/************************************************************************/


/* Unary minus.

   Returns 0 if success, or MDQ_ERROR_xxx if error.
*/
uint32_t mdq_unary_minus(decQuad *r, uint16_t precision, uint16_t scale, decQuad *a) {

  decContext   set;
  uint32_t     mdqerr;


  /* operation */

  decContextDefault(&set, DEC_INIT_DECQUAD);
  decContextSetRounding(&set, DEC_ROUND_HALF_UP);

  decQuadMinus(r, a, &set);

  mdqerr = mdq_adjust_p_s_and_check_error(r, precision, scale, &set);

  return mdqerr;
}


/* Addition.

   Returns 0 if success, or MDQ_ERROR_xxx if error.
*/
uint32_t mdq_add(decQuad *r, uint16_t precision, uint16_t scale, decQuad *a, decQuad *b) {

  decContext   set;
  uint32_t     mdqerr;


  /* operation */

  decContextDefault(&set, DEC_INIT_DECQUAD);
  decContextSetRounding(&set, DEC_ROUND_HALF_UP);

  decQuadAdd(r, a, b, &set);

      // if you want to print the arguments and result, uncomment the line below:
      // mdq_print_string_raw("a = %s\n", a); mdq_print_string_raw("b = %s\n", b); mdq_print_string_raw("r = %s\n\n", r);

  mdqerr = mdq_adjust_p_s_and_check_error(r, precision, scale, &set);

  return mdqerr;
}


/* Subtraction.

   Returns 0 if success, or MDQ_ERROR_xxx if error.
*/
uint32_t mdq_subtract(decQuad *r, uint16_t precision, uint16_t scale, decQuad *a, decQuad *b) {

  decContext   set;
  uint32_t     mdqerr;


  /* operation */

  decContextDefault(&set, DEC_INIT_DECQUAD);
  decContextSetRounding(&set, DEC_ROUND_HALF_UP);

  decQuadSubtract(r, a, b, &set);

  mdqerr = mdq_adjust_p_s_and_check_error(r, precision, scale, &set);

  return mdqerr;
}


/* Multiplication.

   Returns 0 if success, or MDQ_ERROR_xxx if error.
*/
uint32_t mdq_multiply(decQuad *r, uint16_t precision, uint16_t scale, decQuad *a, decQuad *b) {

  decContext   set;
  uint32_t     mdqerr;


  /* operation */

  decContextDefault(&set, DEC_INIT_DECQUAD);
  decContextSetRounding(&set, DEC_ROUND_HALF_UP);

  decQuadMultiply(r, a, b, &set);

  mdqerr = mdq_adjust_p_s_and_check_error(r, precision, scale, &set);

  return mdqerr;
}


/* Division.

   Returns 0 if success, or MDQ_ERROR_xxx if error.
*/
uint32_t mdq_divide(decQuad *r, uint16_t precision, uint16_t scale, decQuad *a, decQuad *b) {

  decContext   set;
  uint32_t     mdqerr;


  /* operation */

  decContextDefault(&set, DEC_INIT_DECQUAD);
  decContextSetRounding(&set, DEC_ROUND_HALF_UP);

  decQuadDivide(r, a, b, &set);

  mdqerr = mdq_adjust_p_s_and_check_error(r, precision, scale, &set); // may return DEC_Division_by_zero or DEC_Division_undefined (for 0/0)

  return mdqerr;
}


/* Compares a with b.

   Returns 1 (greater), 0 (equal), or -1 (less)

   Never fails.
*/
int32_t mdq_compare(decQuad *a, decQuad *b) {

  decContext   set;
  decQuad      r_compare_decnum;


  /* operation */

  decContextDefault(&set, DEC_INIT_DECQUAD);
  decContextSetRounding(&set, DEC_ROUND_HALF_UP);


  // if a or b is Nan

  if ( decQuadIsNaN(b) ) { // if b is Nan
      if ( decQuadIsNaN(a) ) { // if a is also Nan, a (Nan) == b (Nan)
              return 0;
      }
      return 1; // else, a (not-Nan) > b (Nan)
  }

  if ( decQuadIsNaN(a) ) {  // if a is Nan, ( and here, b is not-Nan ), a < b
      return -1; // a (Nan) < b (not-Nan)
  }


  // normal comparison (here, a or b are not Nan, but can be Inf)

  decQuadCompare(&r_compare_decnum, a, b, &set);   // -1 0 1   or Nan if a or b is Nan
  assert( decQuadIsNaN(&r_compare_decnum) == 0 );  // but here, Nan cannot appear
  assert( (set.status & DEC_Errors) == 0 );        // there should be no error

  if ( decQuadIsZero(&r_compare_decnum) ) {
      return 0;
  }

  if ( decQuadIsNegative(&r_compare_decnum) ) {
      return -1;
  }

  return 1;
}


/* Check if values of a and b are equal, and also that their exponents are equal.

   *** THIS FUNCTION MUST BE USED ONLY FOR TESTS ***

   In tests, we want to compare not only the value of a and b, but also their exponents.
   E.g. "12.5" != "12.50"

   Returns 1 or 0.

   Never fails.
*/
int32_t mdq_check_equality_FOR_TEST(decQuad *a, decQuad *b) {

  int32_t res;

  res = mdq_compare(a, b); // compare values

  if ( res != 0 ) {        // if values are not equal, return 0
      return 0;
  }

  // here, values are equal. For finite numbers, check exponents.

  if ( decQuadIsFinite(a) && decQuadIsFinite(b) ) {
      if ( decQuadGetExponent(a) != decQuadGetExponent(b) ) {
          return 0;
      }
      return 1; // two finite values are equal, and their exponent too.
  }

  return 1; // two infinite values or two Nans are equal
}


/* Set to zero.
*/
void mdq_zero(decQuad *r, uint16_t precision, uint16_t scale) {

  decContext   set;
  uint32_t     mdqerr;


  /* operation */

  decContextDefault(&set, DEC_INIT_DECQUAD);
  decContextSetRounding(&set, DEC_ROUND_HALF_UP);

  decQuadZero(r); // never raises an error

  mdqerr = mdq_adjust_p_s_and_check_error(r, precision, scale, &set); // should never returns an error

  assert(mdqerr == 0);
}


/* Copy.

   Returns 0 if success, or MDQ_ERROR_xxx if error.
*/
uint32_t mdq_copy(decQuad *r, uint16_t precision, uint16_t scale, decQuad *a) {

  decContext   set;
  uint32_t     mdqerr;


  /* operation */

  decContextDefault(&set, DEC_INIT_DECQUAD);
  decContextSetRounding(&set, DEC_ROUND_HALF_UP);

  decQuadCopy(r, a); // never raises an error

  mdqerr = mdq_adjust_p_s_and_check_error(r, precision, scale, &set);

  return mdqerr;
}


/* Absolute value.

   Returns 0 if success, or MDQ_ERROR_xxx if error.
*/
uint32_t mdq_abs(decQuad *r, uint16_t precision, uint16_t scale, decQuad *a) {

  decContext   set;
  uint32_t     mdqerr;


  /* operation */

  decContextDefault(&set, DEC_INIT_DECQUAD);
  decContextSetRounding(&set, DEC_ROUND_HALF_UP);

  decQuadAbs(r, a, &set);

  mdqerr = mdq_adjust_p_s_and_check_error(r, precision, scale, &set);

  return mdqerr;
}


/* Ceiling.

   Returns 0 if success, or MDQ_ERROR_xxx if error.
*/
uint32_t mdq_ceiling(decQuad *r, uint16_t precision, uint16_t scale, decQuad *a) {

  decContext   set;
  uint32_t     mdqerr;


  /* operation */

  decContextDefault(&set, DEC_INIT_DECQUAD);
  decContextSetRounding(&set, DEC_ROUND_HALF_UP);

  decQuadToIntegralValue(r, a, &set, DEC_ROUND_CEILING);  // negative exponent becomes 0. (positive exponent are unchanged, but such numbers have been quantized and don't exist in our case)

  mdqerr = mdq_adjust_p_s_and_check_error(r, precision, scale, &set);

  return mdqerr;
}


/* Floor.

   Returns 0 if success, or MDQ_ERROR_xxx if error.
*/
uint32_t mdq_floor(decQuad *r, uint16_t precision, uint16_t scale, decQuad *a) {

  decContext   set;
  uint32_t     mdqerr;


  /* operation */

  decContextDefault(&set, DEC_INIT_DECQUAD);
  decContextSetRounding(&set, DEC_ROUND_HALF_UP);

  decQuadToIntegralValue(r, a, &set, DEC_ROUND_FLOOR);  // negative exponent becomes 0. (positive exponent are unchanged, but such numbers have been quantized and don't exist in our case)

  mdqerr = mdq_adjust_p_s_and_check_error(r, precision, scale, &set);

  return mdqerr;
}


/* Sign.

   Returns 0 if success, or MDQ_ERROR_xxx if error.
*/
uint32_t mdq_sign(decQuad *r, uint16_t precision, uint16_t scale, decQuad *a) {

  decContext   set;
  uint32_t     mdqerr;


  /* check if operand is Nan or Inf */

  if ( ! decQuadIsFinite(a) ) {
      if ( decQuadIsInfinite(a) ) {
          return MDQ_ERROR_INFINITE;
      } else {
          return MDQ_ERROR_NAN;
      }
  }


  /* operation */

  decContextDefault(&set, DEC_INIT_DECQUAD);
  decContextSetRounding(&set, DEC_ROUND_HALF_UP);

  if ( decQuadIsZero(a) ) {
      decQuadZero(r);
  } else if ( decQuadIsPositive(a) ) {  // 1 if a is greater than zero and not a NaN
      decQuadFromInt32(r, 1);
  } else {
      decQuadFromInt32(r, -1);
  }

  mdqerr = mdq_adjust_p_s_and_check_error(r, precision, scale, &set);

  return mdqerr;
}


/* Power.

   Returns 0 if success, or MDQ_ERROR_xxx if error.
*/
uint32_t mdq_power(decQuad *r, uint16_t precision, uint16_t scale, decQuad *a, decQuad *b) {

  decContext   set;
  uint32_t     mdqerr;
  decNumber    num_a;   // working decNumber
  decNumber    num_b;   // working decNumber


  /* operation */

  decContextDefault(&set, DEC_INIT_DECQUAD);
  decContextSetRounding(&set, DEC_ROUND_HALF_UP);


  if ( decQuadIsZero(a) && decQuadIsZero(b) ) { // 0 power 0 is special case

      decQuadFromInt32(r, 1);                               // we can't use decNumberPower(0,0), because it gives an error: invalid operation. MS SQL Server gives 1.

  } else {                                      // if normal operands

      decQuadToNumber(a, &num_a);                           // convert decQuad to decNumber
      decQuadToNumber(b, &num_b);

      decNumberPower(&num_a, &num_a, &num_b, &set);         // we use decNumberPower() because there is no decQuadPower(). Underflow, Overflow or Invalid_operation (e.g. -1**Inf) may occurs.

      decQuadFromNumber(r, &num_a, &set);
  }


  mdqerr = mdq_adjust_p_s_and_check_error(r, precision, scale, &set);

  return mdqerr;
}


uint32_t mdq_round(decQuad *r, uint16_t precision, uint16_t scale, decQuad *a, uint16_t a_precision, uint16_t a_scale, int32_t b, uint8_t truncate_flag) {

  decContext   set;
  uint32_t     mdqerr;
  int64_t      b_val;
  decQuad     *operation_quantizer;


  /* process b */

  b_val = b;

  if ( b_val >= 0 ) {   // round or truncate fractional part

      if ( b_val > a_scale ) {
          b_val = a_scale;      //   a_scale is <= DECQUAD_Pmax          [0..34]
      }

      operation_quantizer = &G_DECQUAD_QUANTIZER[b_val];

  } else {              // b_val < 0, round or truncate integral part

      b_val = -b_val;   // no overflow because b_val is int64_t and b is int32_t

      if ( b_val > (a_precision - a_scale + 1) ) {
          b_val = (a_precision - a_scale + 1);      // [0..35]  see comment for G_DECQUAD_INTEGRAL_PART_QUANTIZER
      }

      operation_quantizer = &G_DECQUAD_INTEGRAL_PART_QUANTIZER[b_val];
  }


  /* operation */

  decContextDefault(&set, DEC_INIT_DECQUAD);
  decContextSetRounding(&set, DEC_ROUND_HALF_UP);

  if ( truncate_flag == 0 ) {

      decQuadQuantize(r, a, operation_quantizer, &set);  // round

  } else {

      decContextSetRounding(&set, DEC_ROUND_DOWN);       // change rounding to truncation mode
      decQuadQuantize(r, a, operation_quantizer, &set);  // truncate
      decContextSetRounding(&set, DEC_ROUND_HALF_UP);    // restore original rounding, because mdq_adjust_p_s_and_check_error() uses default rounding
  }

  mdqerr = mdq_adjust_p_s_and_check_error(r, precision, scale, &set);

  return mdqerr;
}


/* round the fractional part of a number.

   If fractional part length <= b, there is nothing to round, and a copy of the number is returned.

   THIS FUNCTION IS USED ONLY BY SQL 'FORMAT' FUNCTION. See Go rsql/format module.
       The result r is not normalized to any precision or scale, as it is usually the case for the other functions.
*/
uint32_t mdq_round_for_formatting(decQuad *r, decQuad *a, int32_t b) {

  decContext   set;
  int64_t      b_val;
  decQuad     *operation_quantizer;


  /* process b */

  b_val = b;

  if ( b_val < 0 ) {
      b_val = 0;
  }

  if ( b_val > DECQUAD_Pmax ) {
      b_val = DECQUAD_Pmax;      //   b_val is <= DECQUAD_Pmax          [0..34]
  }

  if ( b_val >= -decQuadGetExponent(a) ) { // if there is no need to round fractional part, return an unchanged copy of argument
      decQuadCopy(r, a);
      return 0;
  }

  operation_quantizer = &G_DECQUAD_QUANTIZER[b_val];


  /* operation */

  decContextDefault(&set, DEC_INIT_DECQUAD);
  decContextSetRounding(&set, DEC_ROUND_HALF_UP);

  decQuadQuantize(r, a, operation_quantizer, &set);  // round

  if ( set.status & MYDECQUAD_Errors ) {             // check only real errors. DEC_Inexact and DEC_Rounded are filtered out.
      return mdq_get_status_error(set.status);           // translate real decNumber error code into MDQ_ERROR_XXX error code.
  }


  /* r must be a finite number ( not Inf nor Nan ) */

  if ( ! decQuadIsFinite(r) ) {                      // +Inf -Inf Nan are forbidden in the result.
      if ( decQuadIsInfinite(r) ) {
          return MDQ_ERROR_INFINITE;
      } else {
          return MDQ_ERROR_NAN;
      }
  }

  return 0;
}


/************************************************************************/
/*                                                                      */
/*                      conversion operations                           */
/*                                                                      */
/************************************************************************/


/* fill in a decQuad with given precision and scale, from a int32_t value.

   Returns 0 if success, or MDQ_ERROR_xxx if error.
*/
uint32_t mdq_from_int32(decQuad *r, uint16_t precision, uint16_t scale, int32_t value) {

  decContext   set;
  uint32_t     mdqerr;


  /* operation */

  decContextDefault(&set, DEC_INIT_DECQUAD);
  decContextSetRounding(&set, DEC_ROUND_HALF_UP);

  decQuadFromInt32(r, value);

  mdqerr = mdq_adjust_p_s_and_check_error(r, precision, scale, &set);

  return mdqerr;
}


/* Fill in a decQuad, from a byte array containing a number in ascii.

   val byte array doesn't need to be terminated by 0.

   len is the length of string in val array.

   Argument can be e.g. "-123", "123.45", "12.345e300", "Inf", "-Inf", "Nan".

   The result number is not constrained by any precision or scale.

   Returns 0 if success, or MDQ_ERROR_xxx if error.
*/
uint32_t mdq_from_bytes_raw(decQuad *r, uint8_t *val, int32_t len) {

  decContext  set;
  int32_t     i;
  char        buff[len+1];  // + 1 for terminating 0


  /* operation */

  decContextDefault(&set, DEC_INIT_DECQUAD);
  decContextSetRounding(&set, DEC_ROUND_HALF_UP);

  for ( i=0; i<len; i++ ) {
      buff[i] = (char)val[i];  // copy argument bytes array into local bytes array, because we append a terminating 0
  }
  buff[i] = 0;  // write the terminating 0

  decQuadFromString(r, buff, &set);              // raises an error if string is invalid
  if ( set.status & MYDECQUAD_Errors) {          // check only real errors. DEC_Inexact and DEC_Rounded are filtered out. We don't catch Overflow or Underflow, as result is valid.
      return mdq_get_status_error(set.status);        // translate real decNumber error code into MDQ_ERROR_XXX error code.
  }

  return 0;
}


/* Fill in a decQuad with given precision and scale, from a byte array containing a number in ascii.

   val byte array doesn't need to be terminated by 0

   len is the length of string in val array.

   Returns 0 if success, or MDQ_ERROR_xxx if error.
*/
uint32_t mdq_from_bytes(decQuad *r, uint16_t precision, uint16_t scale, uint8_t *val, int32_t len) {

  decContext  set;
  uint32_t    mdqerr;
  int32_t     i;
  char        buff[len+1];  // + 1 for terminating 0


  /* operation */

  decContextDefault(&set, DEC_INIT_DECQUAD);
  decContextSetRounding(&set, DEC_ROUND_HALF_UP);

  for ( i=0; i<len; i++ ) {
      buff[i] = (char)val[i];  // copy argument bytes array into local bytes array, because we append a terminating 0
  }
  buff[i] = 0;  // write the terminating 0


  decQuadFromString(r, buff, &set); // raises an error if string is invalid

  mdqerr = mdq_adjust_p_s_and_check_error(r, precision, scale, &set);

  return mdqerr;
}


/* Fill in a decQuad from a byte array containing a number in ascii.

   val byte array doesn't need to be terminated by 0

   len is the length of string in val array.

   precision and scale of the resulting decQuad are written and returned in corresponding output parameters 'out_precision' and 'out_scale'.

   If precision exceeds max precision of a NUMERIC (that is, rsql.DATATYPE_NUMERIC_PRECISION_MAX), an error is returned.

   Returns 0 if success, or MDQ_ERROR_xxx if error.
*/
uint32_t mdq_from_bytes_with_implied_p_s(decQuad *r, uint16_t *out_precision, uint16_t *out_scale, uint8_t *val, int32_t len) {

  decContext  set;
  int32_t     exponent;
  int32_t     precision;    // decQuadDigits() returns uint32_t, but the number cannot be negative.
  int32_t     scale;
  uint32_t    mdqerr;
  int32_t     i;
  char        buff[len+1];  // + 1 for terminating 0


  /* operation */

  decContextDefault(&set, DEC_INIT_DECQUAD);
  decContextSetRounding(&set, DEC_ROUND_HALF_UP);

  for ( i=0; i<len; i++ ) {
      buff[i] = (char)val[i];  // copy argument bytes array into local bytes array, because we append a terminating 0
  }
  buff[i] = 0;  // write the terminating 0


  decQuadFromString(r, buff, &set);              // raises an error if string is invalid

  mdqerr = mdq_check_error(r, &set);             // check for decNumber errors, and that result is a finite number ( not Inf nor Nan )
  if ( mdqerr != 0 ) {
    return mdqerr;
  }


  /* adjust the scale of the result. Exponent of 22e3 is 3, of 22.456e3 is 0, of 22.45e3 is 1. Number with exponent > 0 must be quantized to exponent = 0 */

  // 1e35            1e35  ->  10..0e0    must be quantized but will fail because more than 34 digits on left of dot
  // 1e34            1e34  ->  10..0e0    must be quantized
  // 22e3           22e3   ->  22000e0    must be quantized
  // 22.45e3      2245e1   ->  22450e0    must be quantized

  // 22.456e3    22456e0   ->  22456e0
  // 22.4567e3  224567e-1  -> 224567e-1
  // 1e-3            1e-3  ->      1e-3
  // 123e-35       123e-35 ->     12e-34  must be quantized

  // The quantization keeps the numeric value of the number, just allowing a rounding difference. 

  // ****** Only number with exponent out of NUMERIC range [-DECQUAD_Pmax .. 0] must be quantized, to adjust the exponent into this range. ******

  // DECQUAD_Pmax is defined in decQuad.h:         #define DECQUAD_Pmax     34      /* maximum precision (digits)*/
  // The max precision of NUMERIC type, rsql.DATATYPE_NUMERIC_PRECISION_MAX, is also DECQUAD_Pmax.
  // By quantizing, we "flatten" the number, so that it fits in 34 digits precision of a NUMERIC type.

  exponent = decQuadGetExponent(r);  // very large number for Inf, -Inf and Nan. But this is impossible here, because the number is finite.

  if ( exponent > 0 ) {
      decQuadQuantize(r, r, &G_DECQUAD_QUANTIZER[0], &set);  // an error is set if Inf or -Inf (which is impossible here),
                                                             //    or if number cannot be quantized to the desired exponent because it is too large (more than 34 digits before dot)
  } else if ( exponent < -DECQUAD_Pmax ) {
      decQuadQuantize(r, r, &G_DECQUAD_QUANTIZER[DECQUAD_Pmax], &set);
  }

  if ( set.status & MYDECQUAD_Errors) {                     // catch error when quantizing
      return MDQ_ERROR_OVERFLOW;                                // if finite number cannot be quantized, it is an overflow
  }


  /* compute precision and scale of the literal number */

  scale     = -decQuadGetExponent(r);  // always 0 or positive ( decQuadGetExponent(r) always 0 or negative )

  precision = decQuadDigits(r);

  if ( precision < scale ) {
    precision = scale;      // precision of 0.00123 is 5
  }

  assert(precision >= 1 && precision <= DECQUAD_Pmax);
  assert(scale     >= 0 && scale     <= DECQUAD_Pmax);

  *out_precision = precision;  // precision and scale are uint32_t but always fits in a uint16_t
  *out_scale     = scale;

  return 0;
}


/* write decQuad into byte array.

   NOTE: we need this function only because we want a "Go" array of bytes (uint8) as argument, and decQuadToString() requires an array of chars.

   Capacity of byte array must be >= DECQUAD_String.

   A terminating 0 is written in the array.

   Never fails.

   The function decQuadToString() uses exponential notation if number < 0 and too many 0 after decimal point.
*/
void mdq_QuadToString(uint8_t *byte_array, int32_t capacity, decQuad *a) {

  assert(capacity >= DECQUAD_String);

  decQuadToString(a, (char *)byte_array); // cast uint8_t* to char*
}


/* write decQuad into BCD_array.

   IMPORTANT: BCD_array capacity must be DECQUAD_Pmax.

   The output arguments are:
      BCD_array: the coefficient is written one digit per byte.
      exp:       if a is not Inf or Nan, will contain the exponent.
      sign:      if negative and not zero, sign bit is set.
                 THE SIGN IS VALID ALSO IF THE FUNCTION RETURNS MDQ_ERROR_INFINITE, so that we can know if it is +Inf or -Inf.

   Returns 0 if success, or MDQ_ERROR_INFINITE or MDQ_ERROR_NAN.
*/
uint32_t mdq_to_BCD(uint8_t *BCD_array, int32_t *exp, uint32_t *sign, decQuad *a) {

  // convert to BCD

  decQuadToBCD(a, exp, BCD_array);  // this function returns a sign bit, but we don't use it because we don't want -0

  *sign = decQuadIsNegative(a);     // 0 is never negative


  // check that result is not Inf nor Nan

  if ( ! decQuadIsFinite(a) ) {
      if ( decQuadIsInfinite(a) ) {
          return MDQ_ERROR_INFINITE;
      } else {
          return MDQ_ERROR_NAN;
      }
  }

  return 0;
}


/* write decQuad to byte_array.

   This function is like decQuadToString(), except that it outputs the coefficient as integer (no decimal dot), followed by exponent if not 0. E.g. -12345e-4     12345     12345e2      +Inf      -Inf      Nan

   byte_array capacity must be S_STRING_RAW_CAPACITY.

   The function returns the length of the result string (excluding the terminating 0).

   This function is used for testing, to see exactly the coefficient and exponent stored in the decQuad.
   IT IS NOT USEFUL FOR NORMAL USE AND SHOULD BE AVOIDED.
*/
int mdq_to_string_raw(uint8_t *byte_array, int32_t capacity, decQuad *a) {

  uint8_t    *p;
  uint8_t    *p_BCD;
  uint8_t    *BCD_array_sentinel;
  int         n = 0;

  int32_t     exp;
  uint32_t    sign;
  uint8_t     BCD_array[DECQUAD_Pmax];

  assert(capacity >= S_STRING_RAW_CAPACITY);


  // check that result is not Inf nor Nan

  if ( ! decQuadIsFinite(a) ) {
      if ( decQuadIsInfinite(a) ) {
          if ( decQuadIsPositive(a) ) {
              n = sprintf((char*)byte_array, "%s", "+Inf");
          } else {
              n = sprintf((char*)byte_array, "%s", "-Inf");
          }
          return n;

      } else {
          n = sprintf((char*)byte_array, "%s", "Nan");
          return n;
      }
  }


  // convert to BCD

  decQuadToBCD(a, &exp, BCD_array);  // this function returns a sign bit, but we don't use it because we don't want -0

  sign = decQuadIsNegative(a);       // 0 is never negative


  // copy number to byte_array, as integer coefficient and exponent, e.g. 123456e-3     123456e100

  BCD_array_sentinel = BCD_array + sizeof(BCD_array);

  p_BCD = BCD_array;
  p     = byte_array;

  if ( sign ) {
      *p = '-';
      p++;
  }

  while ( p_BCD < BCD_array_sentinel ) { // skip all leading '0'
      if ( *p_BCD != 0 ) {
          break;
      }
      p_BCD++;
  }

  if ( p_BCD == BCD_array_sentinel ) {   // only '0' in BCD_array
      p_BCD--;                           // keep one 0 to print
  }

  while ( p_BCD < BCD_array_sentinel ) { // print out coefficient
      *p = '0' + *p_BCD;
      p_BCD++;
      p++;
  }

  *p = 0;                                // print out trailing 0

  if ( exp != 0 ) {                      // print out exponent if any
      *p = 'e';
      p++;

      n = sprintf((char*)p, "%d", (int)exp);
      p += n;
  }

  assert(*p == 0);

  return p - byte_array;
}


/* print decQuad to stderr, useful for debugging.

   format argument is e.g. "The decQuad is <%s>.\n" with %s as placeholder for the string representation of the number. 
*/
void mdq_print_string_raw(const char *format, decQuad *a) {

  char          buff[S_STRING_RAW_CAPACITY];


  mdq_to_string_raw((uint8_t*)buff, sizeof(buff), a);

  fprintf(stderr, format, buff);
}


/* convert decQuad to int32_t

   Returns 0 if success, or MDQ_ERROR_XXX if error.
*/
uint32_t mdq_to_int32_truncate(int32_t *dest, decQuad *a) {

  decContext   set;
  int32_t      r_val;


  /* operation */

  decContextDefault(&set, DEC_INIT_DECQUAD);
  decContextSetRounding(&set, DEC_ROUND_HALF_UP);


  if ( ! decQuadIsFinite(a) ) {                    // check that a is not Inf nor Nan
      if ( decQuadIsInfinite(a) ) {
          return MDQ_ERROR_INFINITE;
      } else {
          return MDQ_ERROR_NAN;
      }
  }


  r_val = decQuadToInt32(a, &set, DEC_ROUND_DOWN); // truncate. raise error if overflow

  if ( set.status & DEC_Errors ) {                 // check only real errors. DEC_Inexact and DEC_Rounded are filtered out.
      if ( set.status & DEC_Invalid_operation )        // out of int32 range
          return MDQ_ERROR_OUT_OF_RANGE;

      return mdq_get_status_error(set.status);         // translate real decNumber error code into MDQ_ERROR_XXX error code.
  }

  *dest = r_val;

  return 0;
}


/* convert decQuad to int32_t

   Returns 0 if success, or MDQ_ERROR_XXX if error.
*/
uint32_t mdq_to_int32_round(int32_t *dest, decQuad *a) {

  decContext   set;
  int32_t      r_val;


  /* operation */

  decContextDefault(&set, DEC_INIT_DECQUAD);
  decContextSetRounding(&set, DEC_ROUND_HALF_UP);


  if ( ! decQuadIsFinite(a) ) {                    // check that a is not Inf nor Nan
      if ( decQuadIsInfinite(a) ) {
          return MDQ_ERROR_INFINITE;
      } else {
          return MDQ_ERROR_NAN;
      }
  }


  r_val = decQuadToInt32(a, &set, DEC_ROUND_HALF_UP); // <--- ROUNDING MODE. raise error if overflow

  if ( set.status & DEC_Errors ) {                 // check only real errors. DEC_Inexact and DEC_Rounded are filtered out.
      if ( set.status & DEC_Invalid_operation )        // out of int32 range
          return MDQ_ERROR_OUT_OF_RANGE;

      return mdq_get_status_error(set.status);         // translate real decNumber error code into MDQ_ERROR_XXX error code.
  }

  *dest = r_val;

  return 0;
}


/* convert decQuad to int64_t

   Returns 0 if success, or MDQ_ERROR_XXX if error.
*/
uint32_t mdq_to_int64_truncate(int64_t *dest, decQuad *a) {

  decContext   set;
  uint32_t     mdqerr;
  decQuad      a_rounded;
  char         a_str[DECQUAD_String];
  char        *tailptr;
  int64_t      r_val;


  /* operation */

  decContextDefault(&set, DEC_INIT_DECQUAD);
  decContextSetRounding(&set, DEC_ROUND_HALF_UP);

  decQuadToIntegralValue(&a_rounded, a, &set, DEC_ROUND_DOWN);  // truncate. negative exponent becomes 0. (positive exponent are unchanged, but such numbers have been quantized and don't exist in our case)

  mdqerr = mdq_check_error(&a_rounded, &set);   // check for decNumber errors, and that result is a finite number ( not Inf nor Nan )
  if ( mdqerr != 0 ) {
    return mdqerr;
  }


  decQuadToString(&a_rounded, a_str);  // never raises error. Exponential notation never occurs with our NUMERIC numbers, which allows strtoll() to parse the number.


  errno = 0;
  r_val = strtoll(a_str, &tailptr, 10);  // changes errno if error

  if ( errno )
    return MDQ_ERROR_OUT_OF_RANGE;

  if ( *tailptr != 0 )  // may happen for e.g.  123e10, because it parses up to 'e'
    return MDQ_ERROR_OUT_OF_RANGE;

  *dest = r_val;

  return 0;
}


/* convert decQuad to int64_t

   Returns 0 if success, or MDQ_ERROR_XXX if error.
*/
uint32_t mdq_to_int64_round(int64_t *dest, decQuad *a) {

  decContext   set;
  uint32_t     mdqerr;
  decQuad      a_rounded;
  char         a_str[DECQUAD_String];
  char        *tailptr;
  int64_t      r_val;


  /* operation */

  decContextDefault(&set, DEC_INIT_DECQUAD);
  decContextSetRounding(&set, DEC_ROUND_HALF_UP);

  decQuadToIntegralValue(&a_rounded, a, &set, DEC_ROUND_HALF_UP); // <--- ROUNDING MODE. negative exponent becomes 0. (positive exponent are unchanged, but such numbers have been quantized and don't exist in our case)

  mdqerr = mdq_check_error(&a_rounded, &set);   // check for decNumber errors, and that result is a finite number ( not Inf nor Nan )
  if ( mdqerr != 0 ) {
    return mdqerr;
  }


  decQuadToString(&a_rounded, a_str);  // never raises error. Exponential notation never occurs with our NUMERIC numbers, which allows strtoll() to parse the number.


  errno = 0;
  r_val = strtoll(a_str, &tailptr, 10);  // changes errno if error
  if ( errno )
    return MDQ_ERROR_OUT_OF_RANGE;

  if ( *tailptr != 0 )  // may happen for e.g.  123e10, because it parses up to 'e'
    return MDQ_ERROR_OUT_OF_RANGE;

  *dest = r_val;

  return 0;
}


/* convert decQuad to double

   Returns 0 if success, or MDQ_ERROR_XXX if error.
*/
uint32_t mdq_to_double(double *dest, decQuad *a) {

  char         a_str[DECQUAD_String];
  char        *tailptr;
  double       r_val;


  /* operation */

  if ( ! decQuadIsFinite(a) ) {          // check that number is not Inf nor Nan
      if ( decQuadIsInfinite(a) ) {
          return MDQ_ERROR_INFINITE;
      } else {
          return MDQ_ERROR_NAN;
      }
  }


  decQuadToString(a, a_str);  // never raises error


  errno = 0;
  r_val = strtod(a_str, &tailptr);  // changes errno if error (ERANGE if overflow)
  if ( errno )
    return MDQ_ERROR_OUT_OF_RANGE;

  if ( *tailptr != 0 )  // should never happen
    return MDQ_ERROR_DEC_UNLISTED;

  *dest = r_val;

  return 0;
}




