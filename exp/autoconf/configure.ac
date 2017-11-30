#
# This is the configure.ac, which is used by the M4 preprocessor when running autoconf
# to generate configure script.
#
# Note that the 'dnl' is the M4 way to start a comment..
#
echo FIXMEH: wut wut
AC_INIT(program, 1.0)

#
# The following shows up when doing ./configure, right at the start.
#
AC_MSG_NOTICE([Hello1, world.])

dnl
dnl Check for zlib.h
dnl
AC_CHECK_HEADERS(zlib.h)
AC_SEARCH_LIBS(deflate, z)

AC_CHECK_HEADERS(zlib.h, [], [AC_ERROR([A working zlib is required])])
AC_SEARCH_LIBS(deflate, z, [], [AC_ERROR([A working zlib is required])])

dnl Switch to a C++ compiler, and check if it works.
AC_LANG(C++)
AC_PROG_CXX

dnl
dnl Process Makefile.in to create Makefile
dnl
AC_OUTPUT(Makefile)

