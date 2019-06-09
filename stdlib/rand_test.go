package stdlib_test

import (
	"math/rand"
	"testing"

	"github.com/d5/tengo"
	"github.com/d5/tengo/assert"
)

func TestRand(t *testing.T) {
	var seed int64 = 1234
	r := rand.New(rand.NewSource(seed))

	module(t, "rand").call("seed", mockInterop{}, seed).expect(tengo.UndefinedValue)
	module(t, "rand").call("int", mockInterop{}).expect(r.Int63())
	module(t, "rand").call("float", mockInterop{}).expect(r.Float64())
	module(t, "rand").call("intn", mockInterop{}, 111).expect(r.Int63n(111))
	module(t, "rand").call("exp_float", mockInterop{}).expect(r.ExpFloat64())
	module(t, "rand").call("norm_float", mockInterop{}).expect(r.NormFloat64())
	module(t, "rand").call("perm", mockInterop{}, 10).expect(r.Perm(10))

	buf1 := make([]byte, 10)
	buf2 := &tengo.Bytes{Value: make([]byte, 10)}
	n, _ := r.Read(buf1)
	module(t, "rand").call("read", mockInterop{}, buf2).expect(n)
	assert.Equal(t, buf1, buf2.Value)

	seed = 9191
	r = rand.New(rand.NewSource(seed))
	randObj := module(t, "rand").call("rand", mockInterop{}, seed)
	randObj.call("seed", mockInterop{}, seed).expect(tengo.UndefinedValue)
	randObj.call("int", mockInterop{}).expect(r.Int63())
	randObj.call("float", mockInterop{}).expect(r.Float64())
	randObj.call("intn", mockInterop{}, 111).expect(r.Int63n(111))
	randObj.call("exp_float", mockInterop{}).expect(r.ExpFloat64())
	randObj.call("norm_float", mockInterop{}).expect(r.NormFloat64())
	randObj.call("perm", mockInterop{}, 10).expect(r.Perm(10))

	buf1 = make([]byte, 12)
	buf2 = &tengo.Bytes{Value: make([]byte, 12)}
	n, _ = r.Read(buf1)
	randObj.call("read", mockInterop{}, buf2).expect(n)
	assert.Equal(t, buf1, buf2.Value)
}
