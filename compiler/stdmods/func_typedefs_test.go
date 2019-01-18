package stdmods_test

import (
	"errors"
	"testing"

	"github.com/d5/tengo/assert"
	"github.com/d5/tengo/compiler/stdmods"
	"github.com/d5/tengo/objects"
)

func TestFuncAIR(t *testing.T) {
	uf := stdmods.FuncAIR(func(int) {})
	ret, err := uf.Call(&objects.Int{Value: 10})
	assert.NoError(t, err)
	assert.Equal(t, &objects.Undefined{}, ret)
	ret, err = uf.Call()
	assert.Equal(t, objects.ErrWrongNumArguments, err)
}

func TestFuncAR(t *testing.T) {
	uf := stdmods.FuncAR(func() {})
	ret, err := uf.Call()
	assert.NoError(t, err)
	assert.Equal(t, &objects.Undefined{}, ret)
	ret, err = uf.Call(objects.TrueValue)
	assert.Equal(t, objects.ErrWrongNumArguments, err)
}

func TestFuncARI(t *testing.T) {
	uf := stdmods.FuncARI(func() int { return 10 })
	ret, err := uf.Call()
	assert.NoError(t, err)
	assert.Equal(t, &objects.Int{Value: 10}, ret)
	ret, err = uf.Call(objects.TrueValue)
	assert.Equal(t, objects.ErrWrongNumArguments, err)
}

func TestFuncARIsE(t *testing.T) {
	uf := stdmods.FuncARIsE(func() ([]int, error) { return []int{1, 2, 3}, nil })
	ret, err := uf.Call()
	assert.NoError(t, err)
	assert.Equal(t, array(&objects.Int{Value: 1}, &objects.Int{Value: 2}, &objects.Int{Value: 3}), ret)
	uf = stdmods.FuncARIsE(func() ([]int, error) { return nil, errors.New("some error") })
	ret, err = uf.Call()
	assert.NoError(t, err)
	assert.Equal(t, &objects.Error{Value: &objects.String{Value: "some error"}}, ret)
	ret, err = uf.Call(objects.TrueValue)
	assert.Equal(t, objects.ErrWrongNumArguments, err)
}

func TestFuncARS(t *testing.T) {
	uf := stdmods.FuncARS(func() string { return "foo" })
	ret, err := uf.Call()
	assert.NoError(t, err)
	assert.Equal(t, &objects.String{Value: "foo"}, ret)
	ret, err = uf.Call(objects.TrueValue)
	assert.Equal(t, objects.ErrWrongNumArguments, err)
}

func TestFuncARSE(t *testing.T) {
	uf := stdmods.FuncARSE(func() (string, error) { return "foo", nil })
	ret, err := uf.Call()
	assert.NoError(t, err)
	assert.Equal(t, &objects.String{Value: "foo"}, ret)
	uf = stdmods.FuncARSE(func() (string, error) { return "", errors.New("some error") })
	ret, err = uf.Call()
	assert.NoError(t, err)
	assert.Equal(t, &objects.Error{Value: &objects.String{Value: "some error"}}, ret)
	ret, err = uf.Call(objects.TrueValue)
	assert.Equal(t, objects.ErrWrongNumArguments, err)
}

func TestFuncARSs(t *testing.T) {
	uf := stdmods.FuncARSs(func() []string { return []string{"foo", "bar"} })
	ret, err := uf.Call()
	assert.NoError(t, err)
	assert.Equal(t, array(&objects.String{Value: "foo"}, &objects.String{Value: "bar"}), ret)
	ret, err = uf.Call(objects.TrueValue)
	assert.Equal(t, objects.ErrWrongNumArguments, err)
}

func TestFuncASRE(t *testing.T) {
	uf := stdmods.FuncASRE(func(a string) error { return nil })
	ret, err := uf.Call(&objects.String{Value: "foo"})
	assert.NoError(t, err)
	assert.Equal(t, objects.TrueValue, ret)
	uf = stdmods.FuncASRE(func(a string) error { return errors.New("some error") })
	ret, err = uf.Call(&objects.String{Value: "foo"})
	assert.NoError(t, err)
	assert.Equal(t, &objects.Error{Value: &objects.String{Value: "some error"}}, ret)
	ret, err = uf.Call()
	assert.Equal(t, objects.ErrWrongNumArguments, err)
}

func TestFuncASRS(t *testing.T) {
	uf := stdmods.FuncASRS(func(a string) string { return a })
	ret, err := uf.Call(&objects.String{Value: "foo"})
	assert.NoError(t, err)
	assert.Equal(t, &objects.String{Value: "foo"}, ret)
	ret, err = uf.Call()
	assert.Equal(t, objects.ErrWrongNumArguments, err)
}

func TestFuncASI64RE(t *testing.T) {
	uf := stdmods.FuncASI64RE(func(a string, b int64) error { return nil })
	ret, err := uf.Call(&objects.String{Value: "foo"}, &objects.Int{Value: 5})
	assert.NoError(t, err)
	assert.Equal(t, objects.TrueValue, ret)
	uf = stdmods.FuncASI64RE(func(a string, b int64) error { return errors.New("some error") })
	ret, err = uf.Call(&objects.String{Value: "foo"}, &objects.Int{Value: 5})
	assert.NoError(t, err)
	assert.Equal(t, &objects.Error{Value: &objects.String{Value: "some error"}}, ret)
	ret, err = uf.Call()
	assert.Equal(t, objects.ErrWrongNumArguments, err)
}

func TestFuncASIIRE(t *testing.T) {
	uf := stdmods.FuncASIIRE(func(a string, b, c int) error { return nil })
	ret, err := uf.Call(&objects.String{Value: "foo"}, &objects.Int{Value: 5}, &objects.Int{Value: 7})
	assert.NoError(t, err)
	assert.Equal(t, objects.TrueValue, ret)
	uf = stdmods.FuncASIIRE(func(a string, b, c int) error { return errors.New("some error") })
	ret, err = uf.Call(&objects.String{Value: "foo"}, &objects.Int{Value: 5}, &objects.Int{Value: 7})
	assert.NoError(t, err)
	assert.Equal(t, &objects.Error{Value: &objects.String{Value: "some error"}}, ret)
	ret, err = uf.Call()
	assert.Equal(t, objects.ErrWrongNumArguments, err)
}

func TestFuncASRSE(t *testing.T) {
	uf := stdmods.FuncASRSE(func(a string) (string, error) { return a, nil })
	ret, err := uf.Call(&objects.String{Value: "foo"})
	assert.NoError(t, err)
	assert.Equal(t, &objects.String{Value: "foo"}, ret)
	uf = stdmods.FuncASRSE(func(a string) (string, error) { return a, errors.New("some error") })
	ret, err = uf.Call(&objects.String{Value: "foo"})
	assert.NoError(t, err)
	assert.Equal(t, &objects.Error{Value: &objects.String{Value: "some error"}}, ret)
	ret, err = uf.Call()
	assert.Equal(t, objects.ErrWrongNumArguments, err)
}

func TestFuncASSRE(t *testing.T) {

}

func TestFuncARF(t *testing.T) {
	uf := stdmods.FuncARF(func() float64 {
		return 10.0
	})
	ret, err := uf.Call()
	assert.NoError(t, err)
	assert.Equal(t, &objects.Float{Value: 10.0}, ret)
	ret, err = uf.Call(objects.TrueValue)
	assert.Equal(t, objects.ErrWrongNumArguments, err)
}

func TestFuncAFRF(t *testing.T) {
	uf := stdmods.FuncAFRF(func(a float64) float64 {
		return a
	})
	ret, err := uf.Call(&objects.Float{Value: 10.0})
	assert.NoError(t, err)
	assert.Equal(t, &objects.Float{Value: 10.0}, ret)
	ret, err = uf.Call()
	assert.Equal(t, objects.ErrWrongNumArguments, err)
	ret, err = uf.Call(objects.TrueValue, objects.TrueValue)
	assert.Equal(t, objects.ErrWrongNumArguments, err)
}

func TestFuncAIRF(t *testing.T) {
	uf := stdmods.FuncAIRF(func(a int) float64 {
		return float64(a)
	})
	ret, err := uf.Call(&objects.Int{Value: 10.0})
	assert.NoError(t, err)
	assert.Equal(t, &objects.Float{Value: 10.0}, ret)
	ret, err = uf.Call()
	assert.Equal(t, objects.ErrWrongNumArguments, err)
	ret, err = uf.Call(objects.TrueValue, objects.TrueValue)
	assert.Equal(t, objects.ErrWrongNumArguments, err)
}

func TestFuncAFRI(t *testing.T) {
	uf := stdmods.FuncAFRI(func(a float64) int {
		return int(a)
	})
	ret, err := uf.Call(&objects.Float{Value: 10.5})
	assert.NoError(t, err)
	assert.Equal(t, &objects.Int{Value: 10}, ret)
	ret, err = uf.Call()
	assert.Equal(t, objects.ErrWrongNumArguments, err)
	ret, err = uf.Call(objects.TrueValue, objects.TrueValue)
	assert.Equal(t, objects.ErrWrongNumArguments, err)
}

func TestFuncAFRB(t *testing.T) {
	uf := stdmods.FuncAFRB(func(a float64) bool {
		return a > 0.0
	})
	ret, err := uf.Call(&objects.Float{Value: 0.1})
	assert.NoError(t, err)
	assert.Equal(t, &objects.Bool{Value: true}, ret)
	ret, err = uf.Call()
	assert.Equal(t, objects.ErrWrongNumArguments, err)
	ret, err = uf.Call(objects.TrueValue, objects.TrueValue)
	assert.Equal(t, objects.ErrWrongNumArguments, err)
}

func TestFuncAFFRF(t *testing.T) {
	uf := stdmods.FuncAFFRF(func(a, b float64) float64 {
		return a + b
	})
	ret, err := uf.Call(&objects.Float{Value: 10.0}, &objects.Float{Value: 20.0})
	assert.NoError(t, err)
	assert.Equal(t, &objects.Float{Value: 30.0}, ret)
	ret, err = uf.Call()
	assert.Equal(t, objects.ErrWrongNumArguments, err)
	ret, err = uf.Call(objects.TrueValue)
	assert.Equal(t, objects.ErrWrongNumArguments, err)
}

func TestFuncAIFRF(t *testing.T) {
	uf := stdmods.FuncAIFRF(func(a int, b float64) float64 {
		return float64(a) + b
	})
	ret, err := uf.Call(&objects.Int{Value: 10}, &objects.Float{Value: 20.0})
	assert.NoError(t, err)
	assert.Equal(t, &objects.Float{Value: 30.0}, ret)
	ret, err = uf.Call()
	assert.Equal(t, objects.ErrWrongNumArguments, err)
	ret, err = uf.Call(objects.TrueValue)
	assert.Equal(t, objects.ErrWrongNumArguments, err)
}

func TestFuncAFIRF(t *testing.T) {
	uf := stdmods.FuncAFIRF(func(a float64, b int) float64 {
		return a + float64(b)
	})
	ret, err := uf.Call(&objects.Float{Value: 10.0}, &objects.Int{Value: 20})
	assert.NoError(t, err)
	assert.Equal(t, &objects.Float{Value: 30.0}, ret)
	ret, err = uf.Call()
	assert.Equal(t, objects.ErrWrongNumArguments, err)
	ret, err = uf.Call(objects.TrueValue)
	assert.Equal(t, objects.ErrWrongNumArguments, err)
}

func TestFuncAFIRB(t *testing.T) {
	uf := stdmods.FuncAFIRB(func(a float64, b int) bool {
		return a < float64(b)
	})
	ret, err := uf.Call(&objects.Float{Value: 10.0}, &objects.Int{Value: 20})
	assert.NoError(t, err)
	assert.Equal(t, &objects.Bool{Value: true}, ret)
	ret, err = uf.Call()
	assert.Equal(t, objects.ErrWrongNumArguments, err)
	ret, err = uf.Call(objects.TrueValue)
	assert.Equal(t, objects.ErrWrongNumArguments, err)
}

func array(elements ...objects.Object) *objects.Array {
	return &objects.Array{Value: elements}
}
