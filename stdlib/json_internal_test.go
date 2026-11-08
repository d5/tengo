package stdlib

import (
	"bytes"
	gojson "encoding/json"
	"testing"

	"github.com/d5/tengo/v2/require"
)

func TestIndentedLen(t *testing.T) {
	for _, src := range []string{
		`1`, `"a"`, `[]`, `{}`, `[1]`, `[[]]`, `[{}]`, `{"a":[]}`,
		`{"a":1,"b":[1,2,{"c":"d,e:[]{}\"\\"}],"f":{}}`,
		`[[[1,2],[3]],[],{"a":{"b":[true,false,null]}}]`,
	} {
		for _, pi := range [][2]string{{"", ""}, {"", "  "}, {">", "\t"}} {
			var dst bytes.Buffer
			require.NoError(t, gojson.Indent(&dst, []byte(src), pi[0], pi[1]))
			require.Equal(t, int64(dst.Len()),
				indentedLen([]byte(src), pi[0], pi[1]), src)
		}
	}

	// spaces outside strings count as kept, so the length is a bound
	var dst bytes.Buffer
	src := " [ 1 , 2 ] "
	require.NoError(t, gojson.Indent(&dst, []byte(src), "", " "))
	require.True(t, indentedLen([]byte(src), "", " ") >= int64(dst.Len()))
}
