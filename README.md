# decnum

This is a Go binding around C decNumber package, for calculation with base-10 floating point numbers.
Decimal data type is important for financial calculations.

The C decNumber package can be found at:
http://speleotrove.com/decimal/

I downloaded the decNumber package "International Components for Unicode (ICU)".

### Documentation of original C decNumber package

Its documentation is here:
http://speleotrove.com/decimal/decnumber.html

More specifically, you should read these topics:
   - Context: http://speleotrove.com/decimal/dncont.html
   - decQuad: http://speleotrove.com/decimal/dnfloat.html
   - decQuad example: http://speleotrove.com/decimal/dnusers.html#example7


The original C decNumber package contains two kinds of data type:
   - decNumber, which contains arbitrary-precision numbers. Storage will grow as needed.
   - decQuad, decDouble, decSingle, which are fixed-size data types. They are faster than decNumber.


### This Go package
  
__This Go package only uses the decQuad data type__, which is 128 bits long. It can store numbers with 34 significant digits.
It is very much like the float64, except that its precision is better (float64 has a precision of only 15 digits), and it works in base-10 instead of base-2.

I have only written the following files:
   - mydecquad.c
   - mydecquad.h
   - mydecquad.go

They are quite easy to read, and the pattern for calling C function is always the same.
Parameters are always passed from Go to C as value, and in the other direction too.
Strings are passed as array in struct, so they are passed by value too.


### Godoc
https://godoc.org/github.com/rin01/decnum


### Example
You can find an example of use in the directory decnum/example_decnum.


### DISCLAIMER
  I have written this package just to quickly experiment with decimal fixed-size base-10 floating point.
  I think this package works quite well, __but I haven't tested it AT ALL__ !
  So, don't use it in production, or read the code and write tests !
  The code is quite easy to read, and not very long.

  Unfortunately, I really don't have time to work on this package further on.
  So, don't file issues, as I won't read them.


