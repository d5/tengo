package tengo_test

import (
	"context"
	"testing"

	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/require"
)

func TestEval(t *testing.T) {
	eval := func(
		expr string,
		params map[string]interface{},
		expected interface{},
	) {
		ctx := context.Background()
		actual, err := tengo.Eval(ctx, expr, params)
		require.NoError(t, err)
		require.Equal(t, expected, actual)
	}

	eval(`undefined`, nil, nil)
	eval(`1`, nil, int64(1))
	eval(`19 + 23`, nil, int64(42))
	eval(`"foo bar"`, nil, "foo bar")
	eval(`[1, 2, 3][1]`, nil, int64(2))

	eval(
		`5 + p`,
		map[string]interface{}{
			"p": 7,
		},
		int64(12),
	)
	eval(
		`"seven is " + p`,
		map[string]interface{}{
			"p": 7,
		},
		"seven is 7",
	)
	eval(
		`"" + a + b`,
		map[string]interface{}{
			"a": 7,
			"b": " is seven",
		},
		"7 is seven",
	)

	eval(
		`a ? "success" : "fail"`,
		map[string]interface{}{
			"a": 1,
		},
		"success",
	)
}
