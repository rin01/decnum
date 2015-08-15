
ALL THE SOURCE CODE IS UNDER ICU LICENSE (see ICU-license.html)
----------------------------------------


This directory c_original_decnumber_files/decNumber contains files from:

        decNumber package, International Components for Unicode (ICU)
        downloaded from "General Decimal Arithmetic" page at http://speleotrove.com/decimal/
        The license for deNumber package is ICU-license.html.

	This directory is just to have the complete decNumber package files in pristine state, for information.
	This directory is untouched by compilation process and we don't use it.


The directory "decnum" contains:
        - a subset of the files of "decNumber" package, needed for our program, under the ICU license
        - our own code in files: "mydecquad.c", "mydecquad.h", and "mydecquad.go", also under the ICU license
	. all these files will be compiled by using cgo.

	- Just run "go build" in "decnum" directory to compile the "decnum" package. (to install, you must use of course "go install")





