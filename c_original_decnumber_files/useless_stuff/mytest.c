#include <assert.h>
#include "mydecquad.h"

/* This little program just calls some functions of libmydecquad.a

   It is used just for debugging.
   You can modify and play with it.

   To compile, just type "make mytest", and run it with "./mytest"
*/

int main(int argc, char *argv[]) {

    uint32_t     mdqerr;
    decQuad      a;
    decQuad      b;
    decQuad      c;
    decQuad      r;
    int32_t      compres;
    double       d;


    assert(argc); // to silence compiler warning "unused parameter"
    assert(argv); // same

    mdq_init();

    mdqerr = mdq_from_bytes_raw(&a, (uint8_t*)"123.45", sizeof("123.45"));
    assert(mdqerr == 0);

    mdqerr = mdq_from_bytes_raw(&b, (uint8_t*)"12345678901234567890123456789.78", sizeof("12345678901234567890123456789.78"));
    assert(mdqerr == 0);

    mdqerr = mdq_from_bytes_raw(&c, (uint8_t*)"12345678901234567890123456789012.78", sizeof("12345678901234567890123456789012.78"));
    assert(mdqerr == 0);

    mdq_print_string_raw("a is: %s\n", &a);
    mdq_print_string_raw("b is: %s\n", &b);

    // unary minus
    mdqerr = mdq_unary_minus(&r, 34, 2, &a);
    assert(mdqerr == 0);
    mdq_print_string_raw("r is: %s\n", &r);

    // add
    mdqerr = mdq_add(&r, 34, 2, &a, &b);
    assert(mdqerr == 0);
    mdq_print_string_raw("r is: %s\n", &r);

    // subtract
    mdqerr = mdq_subtract(&r, 34, 2, &a, &b);
    assert(mdqerr == 0);
    mdq_print_string_raw("r is: %s\n", &r);

    // multiply
    mdqerr = mdq_multiply(&r, 34, 2, &a, &b);
    assert(mdqerr == 0);
    mdq_print_string_raw("r is: %s\n", &r);

    // divide
    mdqerr = mdq_divide(&r, 34, 2, &a, &b);
    assert(mdqerr == 0);
    mdq_print_string_raw("r is: %s\n", &r);

    // compare
    compres = mdq_compare(&a, &b);
    fprintf(stderr, "compres is: %d\n", compres);

    // to double
    mdqerr = mdq_to_double(&d, &c);
    assert(mdqerr == 0);
    fprintf(stderr, "d is: %g\n", d);


    return 0;
}


