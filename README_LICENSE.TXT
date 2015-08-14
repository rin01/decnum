
This directory c_files/decNumber contains files from:

        decNumber package, International Components for Unicode (ICU)
        downloaded from "General Decimal Arithmetic" page at http://speleotrove.com/decimal/
        The license for deNumber package is ICU-license.html.

	This directory is just to have the complete decNumber package files in pristine state, for information.
	This directory is untouched by compilation process.


The directory "quad" contains:
        - a subset of the files in "decNumber" package, needed for our program, under the ICU license
        - our own code in files: "mydecquad.c", "mydecquad.h", which are released under the same license as the main program ("rsql" database server software)
	. all these files will be compiled by using cgo.

	- Just run "go build" in "quad" directory, or "go build rsql/quad", to compile the "quad" package. (to install, you must use of course "go install")


Note:

By typing "go build -x", we have the confirmation that the proper flags are used ( -fPIC and -pthread ), as shown in the following line from the output:

	...
	gcc -I . -g -O2 -fPIC -m64 -pthread -I $WORK/rsql/quad/_obj/ -o $WORK/rsql/quad/_obj/decNumber.o -c ./decNumber.c
	...
	...

