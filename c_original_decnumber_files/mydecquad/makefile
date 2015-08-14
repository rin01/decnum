
# To see what -pthread does, type this:  gcc -dumpspecs | grep pthread
#     It defines -D_REENTRANT flag, and library -lpthread

# info gcc
# ========
#
# -pthread
#     Adds support for multithreading with the "pthreads" library.  This
#     option sets flags for both the preprocessor and linker.
#
# -fPIC
#     If supported for the target machine, emit position-independent
#     code, suitable for dynamic linking and avoiding any limit on the
#
#     There are some code generation flags preserved by GCC when
#     generating bytecodes, as they need to be used during the final link
#     stage.  Currently, the following options are saved into the GIMPLE
#     bytecode files: `-fPIC', `-fcommon' and all the `-m' target flags.
#     At link time, these options are read in and reapplied.

# info make
# =========
#
# Recipe Execution:
#     When it is time to execute recipes to update a target, they are executed by invoking a new sub-shell for each line of the recipe.
#
# Catalogue of Implicit Rules:
#
#     Compiling C programs
#         n.o is made automatically from n.c with a recipe of the form ‘$(CC) $(CPPFLAGS) $(CFLAGS) -c’.
#
#     Linking a single object file
#         n is made automatically from n.o by running the linker (usually called ld) via the C compiler. The precise recipe used is ‘$(CC) $(LDFLAGS) n.o $(LOADLIBES) $(LDLIBS)’.


CFLAGS = -Wall -Wextra -fPIC -O2 -pthread


objects = mydecquad.o decContext.o decQuad.o decNumber.o decimal128.o decimal64.o

myheaders = mydecquad.h

all: /tmp/libmydecquad.a

/tmp/libmydecquad.a: libmydecquad.a
	cp -t /tmp libmydecquad.a

# we copy the library into /tmp, because cgo needs an absolute path to find it (see golang Issue 5428)
#
libmydecquad.a: $(objects)
	ar rcs $@ $(objects)

mydecquad.o: $(myheaders)


# mytest is a little program that calls some function in libmydecquad.a.
# It is used to debug the library, because it is not possible to use valgrind with Go.
# We can just run "valgrind ./mytest" to see if there is a memory error, or use the cgdb debugger.
# You can modifiy it like you want.
#
mytest: mytest.o libmydecquad.a
	$(CC) $(CPPFLAGS) $(CFLAGS) -o mytest mytest.o libmydecquad.a


.PHONY: clean
clean:
	-rm -f $(objects) libmydecquad.a;
	@echo 'voila, it is cleaned.'



