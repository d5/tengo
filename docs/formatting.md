# Formatting

The format 'verbs' are derived from Go's but are simpler.

## The verbs

## General

```
%v  the value in a default format
%T  a Go-syntax representation of the type of the value
%%  a literal percent sign; consumes no value
```

## Boolean

```
%t  the word true or false
```

## Integer

```
%b  base 2
%c  the character represented by the corresponding Unicode code point
%d  base 10
%o  base 8
%O  base 8 with 0o prefix
%q  a single-quoted character literal safely escaped with Go syntax.
%x  base 16, with lower-case letters for a-f
%X  base 16, with upper-case letters for A-F
%U  Unicode format: U+1234; same as "U+%04X"
```

## Float

```
%b  decimalless scientific notation with exponent a power of two,
in the manner of Go's strconv.FormatFloat with the 'b' format,
e.g. -123456p-78
%e  scientific notation, e.g. -1.234456e+78
%E  scientific notation, e.g. -1.234456E+78
%f  decimal point but no exponent, e.g. 123.456
%F  synonym for %f
%g  %e for large exponents, %f otherwise. Precision is discussed below.
%G  %E for large exponents, %F otherwise
%x  hexadecimal notation (with decimal power of two exponent), e.g. -0x1.23abcp+20
%X  upper-case hexadecimal notation, e.g. -0X1.23ABCP+20
```

## String and Bytes

```
%s  the uninterpreted bytes of the string or slice
%q  a double-quoted string safely escaped with Go syntax
%x  base 16, lower-case, two characters per byte
%X  base 16, upper-case, two characters per byte
```

## Default format for %v

```
Bool:                    %t
Int:                     %d
Float:                   %g
String:                  %s
```

## Compound Objects

```
Array:              [elem0 elem1 ...]
Maps:               {key1:value1 key2:value2 ...}
```

## Width and Precision

Width is specified by an optional decimal number immediately preceding the verb.
If absent, the width is whatever is necessary to represent the value.

Precision is specified after the (optional) width by a period followed by a
decimal number. If no period is present, a default precision is used. A period
with no following number specifies a precision of zero.
Examples:
```
%f     default width, default precision
%9f    width 9, default precision
%.2f   default width, precision 2
%9.2f  width 9, precision 2
%9.f   width 9, precision 0
```

Width and precision are measured in units of Unicode code points.  Either or
both of the flags may be replaced with the character '*', causing their values
to be obtained from the next operand (preceding the one to format), which must
be of type Int.

For most values, width is the minimum number of runes to output, padding the
formatted form with spaces if necessary.

For Strings and Bytes, however, precision limits the length of the input to be
formatted (not the size of the output), truncating if necessary. Normally it is
measured in units of Unicode code points, but for these types when formatted
with the %x or %X format it is measured in bytes.

For floating-point values, width sets the minimum width of the field and
precision sets the number of places after the decimal, if appropriate, except
that for %g/%G precision sets the maximum number of significant digits
(trailing zeros are removed).

For example, given 12.345 the format %6.3f prints 12.345 while %.3g prints 12.3.

The default precision for %e, %f and %#g is 6; for %g it is the smallest number
of digits necessary to identify the value uniquely.

For complex numbers, the width and precision apply to the two components
independently and the result is parenthesized, so %f applied to 1.2+3.4i
produces (1.200000+3.400000i).

## Other flags

```
+   always print a sign for numeric values;
guarantee ASCII-only output for %q (%+q)
-   pad with spaces on the right rather than the left (left-justify the field)
#   alternate format: add leading 0b for binary (%#b), 0 for octal (%#o),
0x or 0X for hex (%#x or %#X);
for %q, print a raw (backquoted) string if strconv.CanBackquote returns true;
always print a decimal point for %e, %E, %f, %F, %g and %G;
do not remove trailing zeros for %g and %G;
write e.g. U+0078 'x' if the character is printable for %U (%#U).
' ' (space) leave a space for elided sign in numbers (% d);
put spaces between bytes printing strings or slices in hex (% x, % X)
0   pad with leading zeros rather than spaces;
for numbers, this moves the padding after the sign
```

Flags are ignored by verbs that do not expect them.
For example there is no alternate decimal format, so %#d and %d behave
identically.

