#include "mydecquad.h"


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


/* unary minus.
*/
Result_t mdq_unary_minus(decQuad a, decContext set) {

  Result_t     res;

  /* operation */

  decQuadMinus(&res.val, &a, &set);            // raises an error if string is invalid
  res.set = set;

  return res;
}


/* addition.
*/
Result_t mdq_add(decQuad a, decQuad b, decContext set) {

  Result_t     res;

  /* operation */

  decQuadAdd(&res.val, &a, &b, &set);            // raises an error if string is invalid
  res.set = set;

  return res;
}





/* convert a int64 into a decQuad.
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



