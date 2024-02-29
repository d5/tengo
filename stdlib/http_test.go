package stdlib_test

import (
	"fmt"
	"testing"

	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/stdlib"
)

func TestHTTP(t *testing.T) {
	run := func(required string, opts string, errMsg string) {
		script := tengo.NewScript([]byte(fmt.Sprintf(`
http := import("http")
res := http.do(%s%s)
`, required, opts)))
		script.SetImports(stdlib.GetModuleMap("http"))

		executed, err := script.Run()
		if err != nil {
			t.Error(err)
		}

		res := executed.Get("res").Value()

		err, ok := res.(error)
		if ok {
			if err.Error() != errMsg {
				t.Errorf("unexpected error: %s", err.Error())
			}
			return
		}
		if !ok && errMsg != "" {
			t.Errorf("missing expected error")
		}

		check := func(name string, ok func(i interface{}) bool) {
			if !ok(res.(map[string]interface{})[name]) {
				t.Errorf("unexpected %s value", name)
			}
		}

		check("code", func(i interface{}) bool { v, ok := i.(int64); return ok && v == 200 })
		check("status", func(i interface{}) bool { v, ok := i.(string); return ok && v == "200 OK" })
		check("headers", func(i interface{}) bool { v, ok := i.(map[string]interface{}); return ok && len(v) > 0 })
		check("body", func(i interface{}) bool { v, ok := i.([]byte); return ok && len(v) > 0 })
	}

	required := `"GET", "https://avatars.githubusercontent.com/u/1291934?s=48&v=4"`
	run(required, `, {dnt: 1, "my-header": "yolo"}, bytes("test")`, ``)
	run(required, ``, ``)
	run(`"GET", "tengo"`, ``, `error: "Get \"tengo\": unsupported protocol scheme \"\""`)
}
