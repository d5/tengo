package stdmods_test

import (
	"testing"

	"github.com/d5/tengo/assert"
	"github.com/d5/tengo/compiler/stdmods"
	"github.com/d5/tengo/objects"
)

func TestFuncARF(t *testing.T) {
	uf := stdmods.FuncARF(func() float64 {
		return 10.0
	})
	ret, err := uf.Call()
	assert.NoError(t, err)
	assert.Equal(t, &objects.Float{Value: 10.0}, ret)
	ret, err = uf.Call(objects.TrueValue)
	assert.Error(t, err)
}

func TestFuncAFRF(t *testing.T) {
	uf := stdmods.FuncAFRF(func(a float64) float64 {
		return a
	})
	ret, err := uf.Call(&objects.Float{Value: 10.0})
	assert.NoError(t, err)
	assert.Equal(t, &objects.Float{Value: 10.0}, ret)
	ret, err = uf.Call()
	assert.Error(t, err)
	ret, err = uf.Call(objects.TrueValue, objects.TrueValue)
	assert.Error(t, err)
}

func TestFuncAIRF(t *testing.T) {
	uf := stdmods.FuncAIRF(func(a int) float64 {
		return float64(a)
	})
	ret, err := uf.Call(&objects.Int{Value: 10.0})
	assert.NoError(t, err)
	assert.Equal(t, &objects.Float{Value: 10.0}, ret)
	ret, err = uf.Call()
	assert.Error(t, err)
	ret, err = uf.Call(objects.TrueValue, objects.TrueValue)
	assert.Error(t, err)
}

func TestFuncAFRI(t *testing.T) {
	uf := stdmods.FuncAFRI(func(a float64) int {
		return int(a)
	})
	ret, err := uf.Call(&objects.Float{Value: 10.5})
	assert.NoError(t, err)
	assert.Equal(t, &objects.Int{Value: 10}, ret)
	ret, err = uf.Call()
	assert.Error(t, err)
	ret, err = uf.Call(objects.TrueValue, objects.TrueValue)
	assert.Error(t, err)
}

func TestFuncAFRB(t *testing.T) {
	uf := stdmods.FuncAFRB(func(a float64) bool {
		return a > 0.0
	})
	ret, err := uf.Call(&objects.Float{Value: 0.1})
	assert.NoError(t, err)
	assert.Equal(t, &objects.Bool{Value: true}, ret)
	ret, err = uf.Call()
	assert.Error(t, err)
	ret, err = uf.Call(objects.TrueValue, objects.TrueValue)
	assert.Error(t, err)
}

func TestFuncAFFRF(t *testing.T) {
	uf := stdmods.FuncAFFRF(func(a, b float64) float64 {
		return a + b
	})
	ret, err := uf.Call(&objects.Float{Value: 10.0}, &objects.Float{Value: 20.0})
	assert.NoError(t, err)
	assert.Equal(t, &objects.Float{Value: 30.0}, ret)
	ret, err = uf.Call()
	assert.Error(t, err)
	ret, err = uf.Call(objects.TrueValue)
	assert.Error(t, err)
}

func TestFuncAIFRF(t *testing.T) {
	uf := stdmods.FuncAIFRF(func(a int, b float64) float64 {
		return float64(a) + b
	})
	ret, err := uf.Call(&objects.Int{Value: 10}, &objects.Float{Value: 20.0})
	assert.NoError(t, err)
	assert.Equal(t, &objects.Float{Value: 30.0}, ret)
	ret, err = uf.Call()
	assert.Error(t, err)
	ret, err = uf.Call(objects.TrueValue)
	assert.Error(t, err)
}

func TestFuncAFIRF(t *testing.T) {
	uf := stdmods.FuncAFIRF(func(a float64, b int) float64 {
		return a + float64(b)
	})
	ret, err := uf.Call(&objects.Float{Value: 10.0}, &objects.Int{Value: 20})
	assert.NoError(t, err)
	assert.Equal(t, &objects.Float{Value: 30.0}, ret)
	ret, err = uf.Call()
	assert.Error(t, err)
	ret, err = uf.Call(objects.TrueValue)
	assert.Error(t, err)
}

func TestFuncAFIRB(t *testing.T) {
	uf := stdmods.FuncAFIRB(func(a float64, b int) bool {
		return a < float64(b)
	})
	ret, err := uf.Call(&objects.Float{Value: 10.0}, &objects.Int{Value: 20})
	assert.NoError(t, err)
	assert.Equal(t, &objects.Bool{Value: true}, ret)
	ret, err = uf.Call()
	assert.Error(t, err)
	ret, err = uf.Call(objects.TrueValue)
	assert.Error(t, err)
}
