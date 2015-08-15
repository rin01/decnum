# decnum
This is a Go binding around C decNumber package, for calculation with base-10 floating point numbers.
Decimal data type is important for financial calculations.

The C decNumber package can be found at:
http://speleotrove.com/decimal/

Its documentation is here:
http://speleotrove.com/decimal/decnumber.html

The original C decNumber package contains two kinds of data type:
  - decNumber, which contains arbitrary-precision numbers. Storage will grow as needed.
  - decQuad, decDouble, decSingle, which are fixed-size data types. They are faster than decNumber.
  
__This Go package only uses the decQuad data type__, which is 128 bits long. It can store numbers with 34 significant digits.
It is very much like the float64, except that its precision is better (float64 has a precision of only 15 digits), and it works in base-10 instead of base-2.



