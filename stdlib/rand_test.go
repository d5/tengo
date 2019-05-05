package stdlib_test

import (
	"math/rand"
	"testing"

	"github.com/d5/tengo/assert"
	"github.com/d5/tengo/objects"
)

func TestRand(t *testing.T) {
	var seed int64 = 1234
	r := rand.New(rand.NewSource(seed))

	module(t, "rand").call("seed", mockRuntime{}, seed).expect(objects.UndefinedValue)
	module(t, "rand").call("int", mockRuntime{}).expect(r.Int63())
	module(t, "rand").call("float", mockRuntime{}).expect(r.Float64())
	module(t, "rand").call("intn", mockRuntime{}, 111).expect(r.Int63n(111))
	module(t, "rand").call("exp_float", mockRuntime{}).expect(r.ExpFloat64())
	module(t, "rand").call("norm_float", mockRuntime{}).expect(r.NormFloat64())
	module(t, "rand").call("perm", mockRuntime{}, 10).expect(r.Perm(10))

	buf1 := make([]byte, 10)
	buf2 := &objects.Bytes{Value: make([]byte, 10)}
	n, _ := r.Read(buf1)
	module(t, "rand").call("read", mockRuntime{}, buf2).expect(n)
	assert.Equal(t, buf1, buf2.Value)

	seed = 9191
	r = rand.New(rand.NewSource(seed))
	randObj := module(t, "rand").call("rand", mockRuntime{}, seed)
	randObj.call("seed", mockRuntime{}, seed).expect(objects.UndefinedValue)
	randObj.call("int", mockRuntime{}).expect(r.Int63())
	randObj.call("float", mockRuntime{}).expect(r.Float64())
	randObj.call("intn", mockRuntime{}, 111).expect(r.Int63n(111))
	randObj.call("exp_float", mockRuntime{}).expect(r.ExpFloat64())
	randObj.call("norm_float", mockRuntime{}).expect(r.NormFloat64())
	randObj.call("perm", mockRuntime{}, 10).expect(r.Perm(10))

	buf1 = make([]byte, 12)
	buf2 = &objects.Bytes{Value: make([]byte, 12)}
	n, _ = r.Read(buf1)
	randObj.call("read", mockRuntime{}, buf2).expect(n)
	assert.Equal(t, buf1, buf2.Value)
}
