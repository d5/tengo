package tengo_test

import (
	"math"
	"strings"
	"testing"

	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/require"
)

type formatTest struct {
	format string
	args   ARR
	want   string
}

func expectFormat(t *testing.T, tests []formatTest) {
	t.Helper()
	var format tengo.CallableFunc
	for _, f := range tengo.GetAllBuiltinFunctions() {
		if f.Name == "format" {
			format = f.Value
			break
		}
	}
	if format == nil {
		t.Fatal("builtin format not found")
	}
	for _, tt := range tests {
		t.Run(tt.format, func(t *testing.T) {
			args := []tengo.Object{&tengo.String{Value: tt.format}}
			for _, a := range tt.args {
				args = append(args, toObject(a))
			}
			res, err := format(args...)
			require.NoError(t, err, "format(%q, %v)", tt.format, tt.args)
			require.Equal(t, tt.want, res.(*tengo.String).Value,
				"format(%q, %v)", tt.format, tt.args)
		})
	}
}

func TestFormatBool(t *testing.T) {
	expectFormat(t, []formatTest{
		{"%t", ARR{true}, "true"},
		{"%t", ARR{false}, "false"},
		{"%6t|", ARR{true}, "  true|"},
		{"%-6t|", ARR{false}, "false |"},
		{"%.2t", ARR{true}, "true"},
		{"%6v|", ARR{true}, "  true|"},
		{"%-6v|", ARR{false}, "false |"},
	})
}

func TestFormatInt(t *testing.T) {
	expectFormat(t, []formatTest{
		{"%d", ARR{-42}, "-42"},
		{"%d", ARR{int64(math.MinInt64)}, "-9223372036854775808"},
		{"%d", ARR{int64(math.MaxInt64)}, "9223372036854775807"},
		{"%b", ARR{5}, "101"},
		{"%b", ARR{-5}, "-101"},
		{"%b", ARR{int64(math.MinInt64)}, "-1" + strings.Repeat("0", 63)},
		{"%#b", ARR{5}, "0b101"},
		{"%#b", ARR{0}, "0b0"},
		{"%o", ARR{8}, "10"},
		{"%o", ARR{-8}, "-10"},
		{"%o", ARR{int64(math.MaxInt64)}, "777777777777777777777"},
		{"%#o", ARR{8}, "010"},
		{"%#o", ARR{-8}, "-010"},
		{"%#o", ARR{0}, "0"},
		{"%O", ARR{8}, "0o10"},
		{"%O", ARR{-8}, "-0o10"},
		{"%#O", ARR{8}, "0o010"},
		{"%x", ARR{255}, "ff"},
		{"%x", ARR{-255}, "-ff"},
		{"%x", ARR{int64(3735928559)}, "deadbeef"},
		{"%X", ARR{255}, "FF"},
		{"%#x", ARR{255}, "0xff"},
		{"%#x", ARR{-255}, "-0xff"},
		{"%#x", ARR{0}, "0x0"},
		{"%#X", ARR{255}, "0XFF"},
		{"%#X", ARR{int64(math.MinInt64)}, "-0X8000000000000000"},
		{"%+d", ARR{5}, "+5"},
		{"%+d", ARR{-5}, "-5"},
		{"% d", ARR{5}, " 5"},
		{"% d", ARR{-5}, "-5"},
		{"%+ d", ARR{5}, "+5"},
		{"%5d|", ARR{42}, "   42|"},
		{"%-5d|", ARR{42}, "42   |"},
		{"%-+6d|", ARR{5}, "+5    |"},
		{"% -6d|", ARR{5}, " 5    |"},
		{"%05d", ARR{42}, "00042"},
		{"%05d", ARR{-42}, "-0042"},
		{"%+05d", ARR{42}, "+0042"},
		{"% 05d", ARR{42}, " 0042"},
		{"%-05d|", ARR{42}, "42   |"},
		{"%.3d", ARR{7}, "007"},
		{"%.3d", ARR{-7}, "-007"},
		{"%+.3d", ARR{5}, "+005"},
		{"%8.3d|", ARR{7}, "     007|"},
		{"%08.3d|", ARR{7}, "     007|"},
		{"%.0d", ARR{3}, "3"},
		{"%.0d|", ARR{0}, "|"},
		{"%.d|", ARR{0}, "|"},
		{"%5.0d|", ARR{0}, "     |"},
		{"%#8x|", ARR{255}, "    0xff|"},
		{"%#-8x|", ARR{255}, "0xff    |"},
		{"%#08x", ARR{255}, "0x000000ff"},
		{"%#.4x", ARR{255}, "0x00ff"},
		{"%100d", ARR{1}, strings.Repeat(" ", 99) + "1"},
		{"%.70d", ARR{1}, strings.Repeat("0", 69) + "1"},
	})
}

func TestFormatIntAsChar(t *testing.T) {
	expectFormat(t, []formatTest{
		{"%c", ARR{65}, "A"},
		{"%c", ARR{0x65E5}, "日"},
		{"%c", ARR{0x1F600}, "😀"},
		{"%c", ARR{0x110000}, "\ufffd"},
		{"%c", ARR{-1}, "\ufffd"},
		{"%3c|", ARR{65}, "  A|"},
		{"%-3c|", ARR{65}, "A  |"},
		{"%q", ARR{65}, `'A'`},
		{"%q", ARR{0x263A}, `'☺'`},
		{"%q", ARR{10}, `'\n'`},
		{"%q", ARR{0}, `'\x00'`},
		{"%+q", ARR{0x263A}, `'\u263a'`},
		{"%5q|", ARR{65}, `  'A'|`},
		{"%U", ARR{65}, "U+0041"},
		{"%U", ARR{0x1F600}, "U+1F600"},
		{"%U", ARR{-1}, "U+FFFFFFFFFFFFFFFF"},
		{"%#U", ARR{65}, "U+0041 'A'"},
		{"%#U", ARR{0x65E5}, "U+65E5 '日'"},
		{"%#U", ARR{7}, "U+0007"},
		{"%#U", ARR{0x110000}, "U+110000"},
		{"%.6U", ARR{65}, "U+000041"},
		{"%#.6U", ARR{65}, "U+000041 'A'"},
		{"%.70U", ARR{1}, "U+" + strings.Repeat("0", 69) + "1"},
		{"%8U|", ARR{65}, "  U+0041|"},
		{"%-8U|", ARR{65}, "U+0041  |"},
		{"%08U|", ARR{65}, "  U+0041|"},
	})
}

func TestFormatString(t *testing.T) {
	expectFormat(t, []formatTest{
		{"%s", ARR{"héllo"}, "héllo"},
		{"%.3s", ARR{"abcdef"}, "abc"},
		{"%.0s|", ARR{"abc"}, "|"},
		{"%.s|", ARR{"abc"}, "|"},
		{"%.10s", ARR{"abc"}, "abc"},
		{"%.2s", ARR{"日本語"}, "日本"},
		{"%5s|", ARR{"ab"}, "   ab|"},
		{"%5s|", ARR{"日本"}, "   日本|"},
		{"%-5s|", ARR{"ab"}, "ab   |"},
		{"%5.1s|", ARR{"abc"}, "    a|"},
		{"%05s", ARR{"ab"}, "000ab"},
		{"%s", ARR{'x'}, "x"},
		{"%v", ARR{'x'}, "x"},
		{"%x", ARR{"hello"}, "68656c6c6f"},
		{"%X", ARR{"hello"}, "68656C6C6F"},
		{"% x", ARR{"hello"}, "68 65 6c 6c 6f"},
		{"% x", ARR{"a"}, "61"},
		{"% x", ARR{"é"}, "c3 a9"},
		{"%#x", ARR{"hello"}, "0x68656c6c6f"},
		{"%# x", ARR{"hello"}, "0x68 0x65 0x6c 0x6c 0x6f"},
		{"%# X", ARR{"hello"}, "0X68 0X65 0X6C 0X6C 0X6F"},
		{"%.2x", ARR{"hello"}, "6865"},
		{"%8x|", ARR{"hi"}, "    6869|"},
		{"%-8x|", ARR{"hi"}, "6869    |"},
		{"% 12x|", ARR{"hi"}, "       68 69|"},
		{"%-# 12x|", ARR{"hi"}, "0x68 0x69   |"},
		{"%x|", ARR{""}, "|"},
		{"%#x|", ARR{""}, "|"},
		{"%4x|", ARR{""}, "    |"},
		{"%5.0x|", ARR{"hi"}, "     |"},
		{"%q", ARR{"hello"}, `"hello"`},
		{"%q", ARR{"日本"}, `"日本"`},
		{"%q", ARR{`a"b`}, `"a\"b"`},
		{"%q", ARR{"a\nb"}, `"a\nb"`},
		{"%q", ARR{"\xff"}, `"\xff"`},
		{"%+q", ARR{"abc"}, `"abc"`},
		{"%+q", ARR{"日本"}, `"\u65e5\u672c"`},
		{"%#q", ARR{"hello"}, "`hello`"},
		{"%#q", ARR{"a\tb"}, "`a\tb`"},
		{"%#q", ARR{"a\nb"}, `"a\nb"`},
		{"%#q", ARR{"a`b"}, "\"a`b\""},
		{"%#+q", ARR{"日本"}, "`日本`"},
		{"%.3q", ARR{"abcdef"}, `"abc"`},
		{"%10q|", ARR{"hi"}, `      "hi"|`},
		{"%-10q|", ARR{"hi"}, `"hi"      |`},
	})
}

func TestFormatBytes(t *testing.T) {
	expectFormat(t, []formatTest{
		{"%s", ARR{[]byte("hello")}, "hello"},
		{"%.2s", ARR{[]byte("héllo")}, "hé"},
		{"%.1s", ARR{[]byte("héllo")}, "h"},
		{"%6s|", ARR{[]byte("hé")}, "    hé|"},
		{"%-6s|", ARR{[]byte("hé")}, "hé    |"},
		{"%x", ARR{[]byte("hello")}, "68656c6c6f"},
		{"%X", ARR{[]byte{0xde, 0xad, 0xbe, 0xef}}, "DEADBEEF"},
		{"% x", ARR{[]byte("hi")}, "68 69"},
		{"%#x", ARR{[]byte("hi")}, "0x6869"},
		{"%# x", ARR{[]byte("hi")}, "0x68 0x69"},
		{"%# X", ARR{[]byte{0xab, 0xcd}}, "0XAB 0XCD"},
		{"%.1x", ARR{[]byte("hi")}, "68"},
		{"%6x|", ARR{[]byte("hi")}, "  6869|"},
		{"%-6x|", ARR{[]byte("hi")}, "6869  |"},
		{"%08x", ARR{[]byte("hi")}, "00006869"},
		{"%x|", ARR{[]byte{}}, "|"},
		{"%3x|", ARR{[]byte{}}, "   |"},
		{"%q", ARR{[]byte("hi")}, `"hi"`},
		{"%#q", ARR{[]byte("hi")}, "`hi`"},
	})
}

func TestFormatFloat(t *testing.T) {
	expectFormat(t, []formatTest{
		{"%e", ARR{1234.5678}, "1.234568e+03"},
		{"%e", ARR{0.0}, "0.000000e+00"},
		{"%e", ARR{1e123}, "1.000000e+123"},
		{"%E", ARR{1234.5678}, "1.234568E+03"},
		{"%f", ARR{1234.5678}, "1234.567800"},
		{"%F", ARR{1234.5678}, "1234.567800"},
		{"%g", ARR{1234.5678}, "1234.5678"},
		{"%g", ARR{1e21}, "1e+21"},
		{"%g", ARR{0.0001}, "0.0001"},
		{"%g", ARR{0.00001}, "1e-05"},
		{"%G", ARR{1e21}, "1E+21"},
		{"%G", ARR{1e-10}, "1E-10"},
		{"%b", ARR{1.0}, "4503599627370496p-52"},
		{"%x", ARR{1.0}, "0x1p+00"},
		{"%x", ARR{-1.0}, "-0x1p+00"},
		{"%x", ARR{1234.5678}, "0x1.34a456d5cfaadp+10"},
		{"%X", ARR{-0.5}, "-0X1P-01"},
		{"%.2x", ARR{1.5}, "0x1.80p+00"},
		{"%.2f", ARR{1234.5678}, "1234.57"},
		{"%.0f", ARR{2.5}, "2"},
		{"%.0f", ARR{3.5}, "4"},
		{"%.f", ARR{3.7}, "4"},
		{"%.3e", ARR{1234.5678}, "1.235e+03"},
		{"%.3g", ARR{1234.5678}, "1.23e+03"},
		{"%.10g", ARR{1.0 / 3}, "0.3333333333"},
		{"%10.2f|", ARR{1234.5678}, "   1234.57|"},
		{"%-10.2f|", ARR{1234.5678}, "1234.57   |"},
		{"%5.1f|", ARR{-0.04}, " -0.0|"},
		{"%12.3e|", ARR{1234.5678}, "   1.235e+03|"},
		{"%-12.3e|", ARR{1234.5678}, "1.235e+03   |"},
		{"%8.3g|", ARR{3.14159}, "    3.14|"},
		{"%010.2f", ARR{1234.5678}, "0001234.57"},
		{"%010.2f", ARR{-1234.5678}, "-001234.57"},
		{"%+010.2f", ARR{3.14159}, "+000003.14"},
		{"% 010.2f", ARR{1234.5678}, " 001234.57"},
		{"%+.2f", ARR{1234.5678}, "+1234.57"},
		{"%+.2f", ARR{-1234.5678}, "-1234.57"},
		{"%+e", ARR{1.0}, "+1.000000e+00"},
		{"% f", ARR{1.5}, " 1.500000"},
		{"% f", ARR{-1.5}, "-1.500000"},
		{"% +f", ARR{1.5}, "+1.500000"},
		{"%#.0f", ARR{3.0}, "3."},
		{"%#.0e", ARR{3.0}, "3.e+00"},
		{"%f", ARR{math.NaN()}, "NaN"},
		{"%e", ARR{math.NaN()}, "NaN"},
		{"%g", ARR{math.NaN()}, "NaN"},
		{"%+f", ARR{math.NaN()}, "+NaN"},
		{"% f", ARR{math.NaN()}, " NaN"},
		{"%5f|", ARR{math.NaN()}, "  NaN|"},
		{"%-5f|", ARR{math.NaN()}, "NaN  |"},
		{"%05f|", ARR{math.NaN()}, "  NaN|"},
		{"%f", ARR{math.Inf(1)}, "+Inf"},
		{"%f", ARR{math.Inf(-1)}, "-Inf"},
		{"%e", ARR{math.Inf(1)}, "+Inf"},
		{"%g", ARR{math.Inf(-1)}, "-Inf"},
		{"%+f", ARR{math.Inf(1)}, "+Inf"},
		{"% f", ARR{math.Inf(1)}, " Inf"},
		{"%6f|", ARR{math.Inf(1)}, "  +Inf|"},
		{"%-6f|", ARR{math.Inf(1)}, "+Inf  |"},
		{"%06f|", ARR{math.Inf(-1)}, "  -Inf|"},
	})
}

func TestFormatStarWidthPrecision(t *testing.T) {
	expectFormat(t, []formatTest{
		{"%*d|", ARR{5, 42}, "   42|"},
		{"%-*d|", ARR{5, 42}, "42   |"},
		{"%*d|", ARR{-5, 42}, "42   |"},
		{"%0*d", ARR{5, 42}, "00042"},
		{"%*s|", ARR{4, "ab"}, "  ab|"},
		{"%.*f", ARR{2, 3.14159}, "3.14"},
		{"%.*s", ARR{2, "abc"}, "ab"},
		{"%.*x", ARR{1, "hello"}, "68"},
		{"%*.*f|", ARR{8, 2, 3.14159}, "    3.14|"},
		{"%*.*x|", ARR{6, 2, "hello"}, "  6865|"},
	})
}

func TestFormatArgIndex(t *testing.T) {
	expectFormat(t, []formatTest{
		{"%[2]d %[1]d", ARR{1, 2}, "2 1"},
		{"%[1]d %[1]d", ARR{7, 8}, "7 7"},
		{"%[2]s", ARR{"a", "b"}, "b"},
		{"%[1]d %d", ARR{1, 2}, "1 2"},
		{"%d %d %#[1]x %#x", ARR{16, 17}, "16 17 0x10 0x11"},
		{"%[2]*[1]d|", ARR{12, 5}, "   12|"},
		{"%[3]*.[2]*[1]f|", ARR{12.0, 2, 6}, " 12.00|"},
		{"%[1]*.[1]*[1]d|", ARR{3}, "003|"},
		{"%[1]d", ARR{1, 2}, "1"},
		{"%[2]d %d", ARR{1, 2}, "2 %!d(MISSING)"},
		{"%[3]d", ARR{1, 2}, "%!d(BADINDEX)"},
		{"%[1]d%[3]d", ARR{1, 2}, "1%!d(BADINDEX)"},
		{"%[0]d", ARR{1}, "%!d(BADINDEX)"},
		{"%[-1]d", ARR{1}, "%!d(BADINDEX)"},
		{"%[x]d", ARR{1}, "%!d(BADINDEX)"},
		{"%[]d", ARR{1}, "%!d(BADINDEX)"},
		{"%[d", ARR{1}, "%!d(BADINDEX)"},
		{"%[1d", ARR{1}, "%!d(BADINDEX)"},
		{"%[99999999]d", ARR{1}, "%!d(BADINDEX)"},
		{"%[2]2d", ARR{1, 2}, "%!d(BADINDEX)"},
		{"%[1].2d", ARR{1, 2}, "%!d(BADINDEX)"},
		{"%[3]日", ARR{1}, "%!日(BADINDEX)"},
		{"%[1]", ARR{1}, "%!(NOVERB)"},
	})
}

func TestFormatArgErrors(t *testing.T) {
	expectFormat(t, []formatTest{
		{"%d %d", ARR{1}, "1 %!d(MISSING)"},
		{"%d %s %q", ARR{1}, "1 %!s(MISSING) %!q(MISSING)"},
		{"%d %日", ARR{1}, "1 %!日(MISSING)"},
		{"%d", ARR{1, 2}, "1%!(EXTRA int=2)"},
		{"%d", ARR{1, 2.5, true}, "1%!(EXTRA float=2.5, bool=true)"},
		{"x", ARR{1}, "x%!(EXTRA int=1)"},
		{"%d %", ARR{1}, "1 %!(NOVERB)"},
		{"%", ARR{1}, "%!(NOVERB)%!(EXTRA int=1)"},
		{"%-", ARR{1}, "%!(NOVERB)%!(EXTRA int=1)"},
		{"%5", ARR{1}, "%!(NOVERB)%!(EXTRA int=1)"},
		{"%99999999d", ARR{1}, "%!(NOVERB)%!(EXTRA int=1)"},
		{"%*d", ARR{"x", 42}, "%!(BADWIDTH)42"},
		{"%*d", ARR{10000000, 42}, "%!(BADWIDTH)42"},
		{"%*d", ARR{-10000000, 42}, "%!(BADWIDTH)42"},
		{"%d %*d", ARR{1}, "1 %!(BADWIDTH)%!d(MISSING)"},
		{"%.*d", ARR{"x", 42}, "%!(BADPREC)42"},
		{"%.*d", ARR{-1, 42}, "%!(BADPREC)42"},
		{"%.*d", ARR{10000000, 42}, "%!(BADPREC)42"},
		{"%d %.*d", ARR{1}, "1 %!(BADPREC)%!d(MISSING)"},
	})
}

func TestFormatTypeAndPercent(t *testing.T) {
	expectFormat(t, []formatTest{
		{"%T %T %T %T", ARR{1, 1.5, "s", true}, "int float string bool"},
		{"%T %T %T %T", ARR{
			[]byte("b"),
			&tengo.Array{},
			&tengo.Map{},
			tengo.UndefinedValue,
		}, "bytes array map undefined"},
		{"%5T|", ARR{1}, "  int|"},
		{"%-8T|", ARR{1}, "int     |"},
		{"%d%%%d", ARR{1, 2}, "1%2"},
		{"%d%%", ARR{50}, "50%"},
		{"%%%d", ARR{7}, "%7"},
		{"%d%%%%", ARR{1}, "1%%"},
		{"%5%%d", ARR{7}, "%7"},
		{"%.2%%d", ARR{7}, "%7"},
	})
}
