package tengo_test

import (
	"errors"
	"math"
	"reflect"
	"testing"
	"time"

	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/require"
)

func Test_builtinDelete(t *testing.T) {
	var builtinDelete func(args ...tengo.Object) (tengo.Object, error)
	for _, f := range tengo.GetAllBuiltinFunctions() {
		if f.Name == "delete" {
			builtinDelete = f.Value
			break
		}
	}
	if builtinDelete == nil {
		t.Fatal("builtin delete not found")
	}
	type args struct {
		args []tengo.Object
	}
	tests := []struct {
		name      string
		args      args
		want      tengo.Object
		wantErr   bool
		wantedErr error
		target    interface{}
	}{
		{name: "invalid-arg", args: args{[]tengo.Object{&tengo.String{},
			&tengo.String{}}}, wantErr: true,
			wantedErr: tengo.ErrInvalidArgumentType{
				Name:     "first",
				Expected: "map",
				Found:    "string"},
		},
		{name: "no-args",
			wantErr: true, wantedErr: tengo.ErrWrongNumArguments},
		{name: "empty-args", args: args{[]tengo.Object{}}, wantErr: true,
			wantedErr: tengo.ErrWrongNumArguments,
		},
		{name: "3-args", args: args{[]tengo.Object{
			(*tengo.Map)(nil), (*tengo.String)(nil), (*tengo.String)(nil)}},
			wantErr: true, wantedErr: tengo.ErrWrongNumArguments,
		},
		{name: "nil-map-empty-key",
			args: args{[]tengo.Object{&tengo.Map{}, &tengo.String{}}},
			want: tengo.UndefinedValue,
		},
		{name: "nil-map-nonstr-key",
			args: args{[]tengo.Object{
				&tengo.Map{}, &tengo.Int{}}}, wantErr: true,
			wantedErr: tengo.ErrInvalidArgumentType{
				Name: "second", Expected: "string", Found: "int"},
		},
		{name: "nil-map-no-key",
			args: args{[]tengo.Object{&tengo.Map{}}}, wantErr: true,
			wantedErr: tengo.ErrWrongNumArguments,
		},
		{name: "map-missing-key",
			args: args{
				[]tengo.Object{
					&tengo.Map{Value: map[string]tengo.Object{
						"key": &tengo.String{Value: "value"},
					}},
					&tengo.String{Value: "key1"}}},
			want: tengo.UndefinedValue,
			target: &tengo.Map{
				Value: map[string]tengo.Object{
					"key": &tengo.String{
						Value: "value"}}},
		},
		{name: "map-emptied",
			args: args{
				[]tengo.Object{
					&tengo.Map{Value: map[string]tengo.Object{
						"key": &tengo.String{Value: "value"},
					}},
					&tengo.String{Value: "key"}}},
			want:   tengo.UndefinedValue,
			target: &tengo.Map{Value: map[string]tengo.Object{}},
		},
		{name: "map-multi-keys",
			args: args{
				[]tengo.Object{
					&tengo.Map{Value: map[string]tengo.Object{
						"key1": &tengo.String{Value: "value1"},
						"key2": &tengo.Int{Value: 10},
					}},
					&tengo.String{Value: "key1"}}},
			want: tengo.UndefinedValue,
			target: &tengo.Map{Value: map[string]tengo.Object{
				"key2": &tengo.Int{Value: 10}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := builtinDelete(tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("builtinDelete() error = %v, wantErr %v",
					err, tt.wantErr)
				return
			}
			if tt.wantErr && !errors.Is(err, tt.wantedErr) {
				if err.Error() != tt.wantedErr.Error() {
					t.Errorf("builtinDelete() error = %v, wantedErr %v",
						err, tt.wantedErr)
					return
				}
			}
			if got != tt.want {
				t.Errorf("builtinDelete() = %v, want %v", got, tt.want)
				return
			}
			if !tt.wantErr && tt.target != nil {
				switch v := tt.args.args[0].(type) {
				case *tengo.Map, *tengo.Array:
					if !reflect.DeepEqual(tt.target, tt.args.args[0]) {
						t.Errorf("builtinDelete() objects are not equal "+
							"got: %+v, want: %+v", tt.args.args[0], tt.target)
					}
				default:
					t.Errorf("builtinDelete() unsuporrted arg[0] type %s",
						v.TypeName())
					return
				}
			}
		})
	}
}

func Test_builtinSplice(t *testing.T) {
	var builtinSplice func(args ...tengo.Object) (tengo.Object, error)
	for _, f := range tengo.GetAllBuiltinFunctions() {
		if f.Name == "splice" {
			builtinSplice = f.Value
			break
		}
	}
	if builtinSplice == nil {
		t.Fatal("builtin splice not found")
	}
	tests := []struct {
		name      string
		args      []tengo.Object
		deleted   tengo.Object
		Array     *tengo.Array
		wantErr   bool
		wantedErr error
	}{
		{name: "no args", args: []tengo.Object{}, wantErr: true,
			wantedErr: tengo.ErrWrongNumArguments,
		},
		{name: "invalid args", args: []tengo.Object{&tengo.Map{}},
			wantErr: true,
			wantedErr: tengo.ErrInvalidArgumentType{
				Name: "first", Expected: "array", Found: "map"},
		},
		{name: "invalid args",
			args:    []tengo.Object{&tengo.Array{}, &tengo.String{}},
			wantErr: true,
			wantedErr: tengo.ErrInvalidArgumentType{
				Name: "second", Expected: "int", Found: "string"},
		},
		{name: "negative index",
			args:      []tengo.Object{&tengo.Array{}, &tengo.Int{Value: -1}},
			wantErr:   true,
			wantedErr: tengo.ErrIndexOutOfBounds},
		{name: "non int count",
			args: []tengo.Object{
				&tengo.Array{}, &tengo.Int{Value: 0},
				&tengo.String{Value: ""}},
			wantErr: true,
			wantedErr: tengo.ErrInvalidArgumentType{
				Name: "third", Expected: "int", Found: "string"},
		},
		{name: "negative count",
			args: []tengo.Object{
				&tengo.Array{Value: []tengo.Object{
					&tengo.Int{Value: 0},
					&tengo.Int{Value: 1},
					&tengo.Int{Value: 2}}},
				&tengo.Int{Value: 0},
				&tengo.Int{Value: -1}},
			wantErr:   true,
			wantedErr: tengo.ErrIndexOutOfBounds,
		},
		{name: "insert with zero count",
			args: []tengo.Object{
				&tengo.Array{Value: []tengo.Object{
					&tengo.Int{Value: 0},
					&tengo.Int{Value: 1},
					&tengo.Int{Value: 2}}},
				&tengo.Int{Value: 0},
				&tengo.Int{Value: 0},
				&tengo.String{Value: "b"}},
			deleted: &tengo.Array{Value: []tengo.Object{}},
			Array: &tengo.Array{Value: []tengo.Object{
				&tengo.String{Value: "b"},
				&tengo.Int{Value: 0},
				&tengo.Int{Value: 1},
				&tengo.Int{Value: 2}}},
		},
		{name: "insert",
			args: []tengo.Object{
				&tengo.Array{Value: []tengo.Object{
					&tengo.Int{Value: 0},
					&tengo.Int{Value: 1},
					&tengo.Int{Value: 2}}},
				&tengo.Int{Value: 1},
				&tengo.Int{Value: 0},
				&tengo.String{Value: "c"},
				&tengo.String{Value: "d"}},
			deleted: &tengo.Array{Value: []tengo.Object{}},
			Array: &tengo.Array{Value: []tengo.Object{
				&tengo.Int{Value: 0},
				&tengo.String{Value: "c"},
				&tengo.String{Value: "d"},
				&tengo.Int{Value: 1},
				&tengo.Int{Value: 2}}},
		},
		{name: "insert with zero count",
			args: []tengo.Object{
				&tengo.Array{Value: []tengo.Object{
					&tengo.Int{Value: 0},
					&tengo.Int{Value: 1},
					&tengo.Int{Value: 2}}},
				&tengo.Int{Value: 1},
				&tengo.Int{Value: 0},
				&tengo.String{Value: "c"},
				&tengo.String{Value: "d"}},
			deleted: &tengo.Array{Value: []tengo.Object{}},
			Array: &tengo.Array{Value: []tengo.Object{
				&tengo.Int{Value: 0},
				&tengo.String{Value: "c"},
				&tengo.String{Value: "d"},
				&tengo.Int{Value: 1},
				&tengo.Int{Value: 2}}},
		},
		{name: "insert with delete",
			args: []tengo.Object{
				&tengo.Array{Value: []tengo.Object{
					&tengo.Int{Value: 0},
					&tengo.Int{Value: 1},
					&tengo.Int{Value: 2}}},
				&tengo.Int{Value: 1},
				&tengo.Int{Value: 1},
				&tengo.String{Value: "c"},
				&tengo.String{Value: "d"}},
			deleted: &tengo.Array{
				Value: []tengo.Object{&tengo.Int{Value: 1}}},
			Array: &tengo.Array{Value: []tengo.Object{
				&tengo.Int{Value: 0},
				&tengo.String{Value: "c"},
				&tengo.String{Value: "d"},
				&tengo.Int{Value: 2}}},
		},
		{name: "insert with delete multi",
			args: []tengo.Object{
				&tengo.Array{Value: []tengo.Object{
					&tengo.Int{Value: 0},
					&tengo.Int{Value: 1},
					&tengo.Int{Value: 2}}},
				&tengo.Int{Value: 1},
				&tengo.Int{Value: 2},
				&tengo.String{Value: "c"},
				&tengo.String{Value: "d"}},
			deleted: &tengo.Array{Value: []tengo.Object{
				&tengo.Int{Value: 1},
				&tengo.Int{Value: 2}}},
			Array: &tengo.Array{
				Value: []tengo.Object{
					&tengo.Int{Value: 0},
					&tengo.String{Value: "c"},
					&tengo.String{Value: "d"}}},
		},
		{name: "delete all with positive count",
			args: []tengo.Object{
				&tengo.Array{Value: []tengo.Object{
					&tengo.Int{Value: 0},
					&tengo.Int{Value: 1},
					&tengo.Int{Value: 2}}},
				&tengo.Int{Value: 0},
				&tengo.Int{Value: 3}},
			deleted: &tengo.Array{Value: []tengo.Object{
				&tengo.Int{Value: 0},
				&tengo.Int{Value: 1},
				&tengo.Int{Value: 2}}},
			Array: &tengo.Array{Value: []tengo.Object{}},
		},
		{name: "delete all with big count",
			args: []tengo.Object{
				&tengo.Array{Value: []tengo.Object{
					&tengo.Int{Value: 0},
					&tengo.Int{Value: 1},
					&tengo.Int{Value: 2}}},
				&tengo.Int{Value: 0},
				&tengo.Int{Value: 5}},
			deleted: &tengo.Array{Value: []tengo.Object{
				&tengo.Int{Value: 0},
				&tengo.Int{Value: 1},
				&tengo.Int{Value: 2}}},
			Array: &tengo.Array{Value: []tengo.Object{}},
		},
		{name: "nothing2",
			args: []tengo.Object{
				&tengo.Array{Value: []tengo.Object{
					&tengo.Int{Value: 0},
					&tengo.Int{Value: 1},
					&tengo.Int{Value: 2}}}},
			Array: &tengo.Array{Value: []tengo.Object{}},
			deleted: &tengo.Array{Value: []tengo.Object{
				&tengo.Int{Value: 0},
				&tengo.Int{Value: 1},
				&tengo.Int{Value: 2}}},
		},
		{name: "pop without count",
			args: []tengo.Object{
				&tengo.Array{Value: []tengo.Object{
					&tengo.Int{Value: 0},
					&tengo.Int{Value: 1},
					&tengo.Int{Value: 2}}},
				&tengo.Int{Value: 2}},
			deleted: &tengo.Array{Value: []tengo.Object{&tengo.Int{Value: 2}}},
			Array: &tengo.Array{Value: []tengo.Object{
				&tengo.Int{Value: 0}, &tengo.Int{Value: 1}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := builtinSplice(tt.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("builtinSplice() error = %v, wantErr %v",
					err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.deleted) {
				t.Errorf("builtinSplice() = %v, want %v", got, tt.deleted)
			}
			if tt.wantErr && tt.wantedErr.Error() != err.Error() {
				t.Errorf("builtinSplice() error = %v, wantedErr %v",
					err, tt.wantedErr)
			}
			if tt.Array != nil && !reflect.DeepEqual(tt.Array, tt.args[0]) {
				t.Errorf("builtinSplice() arrays are not equal expected"+
					" %s, got %s", tt.Array, tt.args[0].(*tengo.Array))
			}
		})
	}
}

func Test_builtinRange(t *testing.T) {
	var builtinRange func(args ...tengo.Object) (tengo.Object, error)
	for _, f := range tengo.GetAllBuiltinFunctions() {
		if f.Name == "range" {
			builtinRange = f.Value
			break
		}
	}
	if builtinRange == nil {
		t.Fatal("builtin range not found")
	}
	tests := []struct {
		name      string
		args      []tengo.Object
		result    *tengo.Array
		wantErr   bool
		wantedErr error
	}{
		{name: "no args", args: []tengo.Object{}, wantErr: true,
			wantedErr: tengo.ErrWrongNumArguments,
		},
		{name: "single args", args: []tengo.Object{&tengo.Map{}},
			wantErr:   true,
			wantedErr: tengo.ErrWrongNumArguments,
		},
		{name: "4 args", args: []tengo.Object{&tengo.Map{}, &tengo.String{}, &tengo.String{}, &tengo.String{}},
			wantErr:   true,
			wantedErr: tengo.ErrWrongNumArguments,
		},
		{name: "invalid start",
			args:    []tengo.Object{&tengo.String{}, &tengo.String{}},
			wantErr: true,
			wantedErr: tengo.ErrInvalidArgumentType{
				Name: "start", Expected: "int", Found: "string"},
		},
		{name: "invalid stop",
			args:    []tengo.Object{&tengo.Int{}, &tengo.String{}},
			wantErr: true,
			wantedErr: tengo.ErrInvalidArgumentType{
				Name: "stop", Expected: "int", Found: "string"},
		},
		{name: "invalid step",
			args:    []tengo.Object{&tengo.Int{}, &tengo.Int{}, &tengo.String{}},
			wantErr: true,
			wantedErr: tengo.ErrInvalidArgumentType{
				Name: "step", Expected: "int", Found: "string"},
		},
		{name: "zero step",
			args:      []tengo.Object{&tengo.Int{}, &tengo.Int{}, &tengo.Int{}}, //must greate than 0
			wantErr:   true,
			wantedErr: tengo.ErrInvalidRangeStep,
		},
		{name: "negative step",
			args:      []tengo.Object{&tengo.Int{}, &tengo.Int{}, intObject(-2)}, //must greate than 0
			wantErr:   true,
			wantedErr: tengo.ErrInvalidRangeStep,
		},
		{name: "same bound",
			args:    []tengo.Object{&tengo.Int{}, &tengo.Int{}},
			wantErr: false,
			result: &tengo.Array{
				Value: nil,
			},
		},
		{name: "positive range",
			args:    []tengo.Object{&tengo.Int{}, &tengo.Int{Value: 5}},
			wantErr: false,
			result: &tengo.Array{
				Value: []tengo.Object{
					intObject(0),
					intObject(1),
					intObject(2),
					intObject(3),
					intObject(4),
				},
			},
		},
		{name: "negative range",
			args:    []tengo.Object{&tengo.Int{}, &tengo.Int{Value: -5}},
			wantErr: false,
			result: &tengo.Array{
				Value: []tengo.Object{
					intObject(0),
					intObject(-1),
					intObject(-2),
					intObject(-3),
					intObject(-4),
				},
			},
		},

		{name: "positive with step",
			args:    []tengo.Object{&tengo.Int{}, &tengo.Int{Value: 5}, &tengo.Int{Value: 2}},
			wantErr: false,
			result: &tengo.Array{
				Value: []tengo.Object{
					intObject(0),
					intObject(2),
					intObject(4),
				},
			},
		},

		{name: "negative with step",
			args:    []tengo.Object{&tengo.Int{}, &tengo.Int{Value: -10}, &tengo.Int{Value: 2}},
			wantErr: false,
			result: &tengo.Array{
				Value: []tengo.Object{
					intObject(0),
					intObject(-2),
					intObject(-4),
					intObject(-6),
					intObject(-8),
				},
			},
		},

		{name: "large range",
			args:    []tengo.Object{intObject(-10), intObject(10), &tengo.Int{Value: 3}},
			wantErr: false,
			result: &tengo.Array{
				Value: []tengo.Object{
					intObject(-10),
					intObject(-7),
					intObject(-4),
					intObject(-1),
					intObject(2),
					intObject(5),
					intObject(8),
				},
			},
		},
		{name: "step larger than span does not overflow",
			args:    []tengo.Object{intObject(0), intObject(10), intObject(math.MaxInt64 - 1)},
			wantErr: false,
			result: &tengo.Array{
				Value: []tengo.Object{intObject(0)},
			},
		},
		{name: "span overflows int64",
			args:      []tengo.Object{intObject(math.MinInt64), intObject(math.MaxInt64)},
			wantErr:   true,
			wantedErr: tengo.ErrRangeLimit,
		},
		{name: "exceeds size limit",
			args:      []tengo.Object{intObject(0), intObject(int64(tengo.MaxRangeLen) + 1)},
			wantErr:   true,
			wantedErr: tengo.ErrRangeLimit,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := builtinRange(tt.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("builtinRange() error = %v, wantErr %v",
					err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.wantedErr.Error() != err.Error() {
				t.Errorf("builtinRange() error = %v, wantedErr %v",
					err, tt.wantedErr)
			}
			if tt.result != nil && !reflect.DeepEqual(tt.result, got) {
				t.Errorf("builtinRange() arrays are not equal expected"+
					" %s, got %s", tt.result, got.(*tengo.Array))
			}
		})
	}
}

func Benchmark_builtinRange(b *testing.B) {
	var builtinRange func(args ...tengo.Object) (tengo.Object, error)
	for _, f := range tengo.GetAllBuiltinFunctions() {
		if f.Name == "range" {
			builtinRange = f.Value
			break
		}
	}
	if builtinRange == nil {
		b.Fatal("builtin range not found")
	}
	benchmarks := []struct {
		name string
		args []tengo.Object
	}{
		{"asc-10", []tengo.Object{
			&tengo.Int{Value: 0}, &tengo.Int{Value: 10}}},
		{"asc-1000", []tengo.Object{
			&tengo.Int{Value: 0}, &tengo.Int{Value: 1000}}},
		{"asc-100000", []tengo.Object{
			&tengo.Int{Value: 0}, &tengo.Int{Value: 100000}}},
		{"asc-100000-step-10", []tengo.Object{
			&tengo.Int{Value: 0}, &tengo.Int{Value: 100000},
			&tengo.Int{Value: 10}}},
		{"desc-1000", []tengo.Object{
			&tengo.Int{Value: 1000}, &tengo.Int{Value: 0}}},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				if _, err := builtinRange(bm.args...); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func TestBuiltinFreeze(t *testing.T) {
	var freeze func(args ...tengo.Object) (tengo.Object, error)
	for _, f := range tengo.GetAllBuiltinFunctions() {
		if f.Name == "freeze" {
			freeze = f.Value
			break
		}
	}
	if freeze == nil {
		t.Fatal("builtin freeze not found")
	}

	t.Run("array becomes immutable array", func(t *testing.T) {
		arr := &tengo.Array{Value: []tengo.Object{
			&tengo.Int{Value: 1}, &tengo.Int{Value: 2},
		}}
		got, err := freeze(arr)
		if err != nil {
			t.Fatal(err)
		}
		ia, ok := got.(*tengo.ImmutableArray)
		if !ok {
			t.Fatalf("expected *ImmutableArray, got %T", got)
		}
		if len(ia.Value) != 2 || !reflect.DeepEqual(ia.Value[0], &tengo.Int{Value: 1}) {
			t.Fatalf("unexpected value: %v", ia)
		}
	})

	t.Run("map becomes immutable map", func(t *testing.T) {
		m := &tengo.Map{Value: map[string]tengo.Object{
			"x": &tengo.Int{Value: 42},
		}}
		got, err := freeze(m)
		if err != nil {
			t.Fatal(err)
		}
		im, ok := got.(*tengo.ImmutableMap)
		if !ok {
			t.Fatalf("expected *ImmutableMap, got %T", got)
		}
		if !reflect.DeepEqual(im.Value["x"], &tengo.Int{Value: 42}) {
			t.Fatalf("unexpected value: %v", im)
		}
	})

	t.Run("nested structures frozen recursively", func(t *testing.T) {
		inner := &tengo.Array{Value: []tengo.Object{&tengo.Int{Value: 7}}}
		outer := &tengo.Map{Value: map[string]tengo.Object{"inner": inner}}

		got, err := freeze(outer)
		if err != nil {
			t.Fatal(err)
		}
		im := got.(*tengo.ImmutableMap)
		if _, ok := im.Value["inner"].(*tengo.ImmutableArray); !ok {
			t.Fatalf("inner array not frozen, got %T", im.Value["inner"])
		}
	})

	t.Run("already-immutable array with no mutable children returns same pointer", func(t *testing.T) {
		ia := &tengo.ImmutableArray{Value: []tengo.Object{&tengo.Int{Value: 1}}}
		got, err := freeze(ia)
		if err != nil {
			t.Fatal(err)
		}
		if got != ia {
			t.Fatal("expected same pointer for already-frozen array")
		}
	})

	t.Run("already-immutable map with no mutable children returns same pointer", func(t *testing.T) {
		im := &tengo.ImmutableMap{Value: map[string]tengo.Object{"a": &tengo.Int{Value: 1}}}
		got, err := freeze(im)
		if err != nil {
			t.Fatal(err)
		}
		if got != im {
			t.Fatal("expected same pointer for already-frozen map")
		}
	})

	t.Run("immutable array with mutable element is re-frozen", func(t *testing.T) {
		inner := &tengo.Array{Value: []tengo.Object{&tengo.Int{Value: 5}}}
		ia := &tengo.ImmutableArray{Value: []tengo.Object{inner}}
		got, err := freeze(ia)
		if err != nil {
			t.Fatal(err)
		}
		result := got.(*tengo.ImmutableArray)
		if result == ia {
			t.Fatal("expected a new ImmutableArray since element was mutable")
		}
		if _, ok := result.Value[0].(*tengo.ImmutableArray); !ok {
			t.Fatalf("inner element not frozen, got %T", result.Value[0])
		}
	})

	t.Run("primitives pass through unchanged", func(t *testing.T) {
		for _, obj := range []tengo.Object{
			&tengo.Int{Value: 1},
			&tengo.Float{Value: 3.14},
			tengo.TrueValue,
			&tengo.String{Value: "hello"},
			tengo.UndefinedValue,
		} {
			got, err := freeze(obj)
			if err != nil {
				t.Fatal(err)
			}
			if got != obj {
				t.Fatalf("primitive %T was not returned as-is", obj)
			}
		}
	})

	t.Run("self-referential array does not infinite-loop", func(t *testing.T) {
		arr := &tengo.Array{Value: []tengo.Object{&tengo.Int{Value: 1}}}
		arr.Value = append(arr.Value, arr) // arr contains itself
		got, err := freeze(arr)
		if err != nil {
			t.Fatal(err)
		}
		if _, ok := got.(*tengo.ImmutableArray); !ok {
			t.Fatalf("expected *ImmutableArray, got %T", got)
		}
	})

	t.Run("wrong number of arguments", func(t *testing.T) {
		_, err := freeze()
		if err == nil {
			t.Fatal("expected error for no arguments")
		}
		_, err = freeze(&tengo.Int{Value: 1}, &tengo.Int{Value: 2})
		if err == nil {
			t.Fatal("expected error for two arguments")
		}
	})
}

func TestBuiltinIsType(t *testing.T) {
	expectRun(t, `out = is_float(1.5)`, nil, true)
	expectRun(t, `out = is_float(1)`, nil, false)
	expectRun(t, `out = is_float("1.5")`, nil, false)

	expectRun(t, `out = is_bool(true)`, nil, true)
	expectRun(t, `out = is_bool(false)`, nil, true)
	expectRun(t, `out = is_bool(0)`, nil, false)
	expectRun(t, `out = is_bool(undefined)`, nil, false)

	expectRun(t, `out = is_char('a')`, nil, true)
	expectRun(t, `out = is_char(char(97))`, nil, true)
	expectRun(t, `out = is_char("a")`, nil, false)
	expectRun(t, `out = is_char(97)`, nil, false)

	expectRun(t, `out = is_bytes(bytes("a"))`, nil, true)
	expectRun(t, `out = is_bytes(bytes(0))`, nil, true)
	expectRun(t, `out = is_bytes("a")`, nil, false)
	expectRun(t, `out = is_bytes([97])`, nil, false)

	expectRun(t, `out = is_time(time(0))`, nil, true)
	expectRun(t, `out = is_time(0)`, nil, false)
	expectRun(t, `out = is_time("2009-11-10")`, nil, false)

	for _, iterable := range []string{
		`[1]`, `[]`, `immutable([1])`, `{a: 1}`, `{}`, `immutable({a: 1})`,
		`"abc"`, `""`, `bytes("a")`,
	} {
		expectRun(t, `out = is_iterable(`+iterable+`)`, nil, true)
	}
	for _, notIterable := range []string{
		`1`, `1.5`, `true`, `'a'`, `error("e")`, `time(0)`, `func() {}`,
		`len`,
	} {
		expectRun(t, `out = is_iterable(`+notIterable+`)`, nil, false)
	}
}

func TestBuiltinWrongNumArguments(t *testing.T) {
	one := &tengo.Int{Value: 1}
	for _, name := range []string{
		"type_name", "len", "copy", "bool",
		"is_int", "is_float", "is_string", "is_bool", "is_char", "is_bytes",
		"is_array", "is_immutable_array", "is_map", "is_immutable_map",
		"is_iterable", "is_time", "is_error", "is_undefined", "is_function",
		"is_callable",
	} {
		fn := builtinFunc(t, name)
		for _, args := range [][]tengo.Object{nil, {one, one}} {
			_, err := fn(args...)
			require.Equal(t, tengo.ErrWrongNumArguments, err,
				"%s with %d args", name, len(args))
		}
	}

	for _, name := range []string{
		"string", "int", "float", "char", "bytes", "time",
	} {
		fn := builtinFunc(t, name)
		for _, args := range [][]tengo.Object{nil, {one, one, one}} {
			_, err := fn(args...)
			require.Equal(t, tengo.ErrWrongNumArguments, err,
				"%s with %d args", name, len(args))
		}
	}

	_, err := builtinFunc(t, "append")()
	require.Equal(t, tengo.ErrWrongNumArguments, err)

	expectError(t, `type_name()`, nil,
		"wrong number of arguments in call to 'builtin-function:type_name'")
	expectError(t, `is_iterable(1, 2)`, nil,
		"wrong number of arguments in call to 'builtin-function:is_iterable'")
}

func TestBuiltinTime(t *testing.T) {
	expectRun(t, `out = time(1257894000) - time(1257893999)`, nil,
		int64(time.Second))
	expectRun(t, `out = time(1257894000) == time(1257894000)`, nil, true)
	expectRun(t, `t := time(1257894000); out = time(t) == t`, nil, true)
	expectRun(t, `out = time(1257894000, "x") == time(1257894000)`, nil, true)
	expectRun(t, `out = time("1257894000")`, nil, tengo.UndefinedValue)
	expectRun(t, `out = time(1.5)`, nil, tengo.UndefinedValue)
	expectRun(t, `out = time(undefined)`, nil, tengo.UndefinedValue)
	expectRun(t, `out = time(undefined, "fallback")`, nil, "fallback")
	expectRun(t, `out = time([], 5)`, nil, 5)

	res, err := builtinFunc(t, "time")(&tengo.Int{Value: 1257894000})
	require.NoError(t, err)
	require.True(t, res.(*tengo.Time).Value.Equal(time.Unix(1257894000, 0)))
}

func TestBuiltinAppendInvalidArgument(t *testing.T) {
	expectError(t, `append()`, nil,
		"wrong number of arguments in call to 'builtin-function:append'")
	expectError(t, `append(1, 2)`, nil,
		"invalid type for argument 'first' in call to "+
			"'builtin-function:append': expected array, found int")
	expectError(t, `append({}, 2)`, nil,
		"invalid type for argument 'first' in call to "+
			"'builtin-function:append': expected array, found map")
	expectError(t, `append("abc", "d")`, nil,
		"invalid type for argument 'first' in call to "+
			"'builtin-function:append': expected array, found string")
	expectError(t, `append(undefined, 1)`, nil,
		"invalid type for argument 'first' in call to "+
			"'builtin-function:append': expected array, found undefined")
}

func TestBuiltinLenBytes(t *testing.T) {
	expectRun(t, `out = len(bytes(""))`, nil, 0)
	expectRun(t, `out = len(bytes("abc"))`, nil, 3)
	expectRun(t, `out = len(bytes(5))`, nil, 5)
	expectRun(t, `out = len(bytes("가"))`, nil, 3)
}

func TestBuiltinTypeName(t *testing.T) {
	expectRun(t, `out = type_name(immutable([1]))`, nil, "immutable-array")
	expectRun(t, `out = type_name(immutable({a: 1}))`, nil, "immutable-map")
	expectRun(t, `out = type_name(time(0))`, nil, "time")
	expectRun(t, `out = type_name(len)`, nil, "builtin-function:len")
	expectRun(t, `out = type_name(x)`,
		Opts().Symbol("x", &tengo.UserFunction{Name: "fn"}).Skip2ndPass(),
		"user-function:fn")
}

func builtinFunc(t *testing.T, name string) tengo.CallableFunc {
	for _, f := range tengo.GetAllBuiltinFunctions() {
		if f.Name == name {
			return f.Value
		}
	}
	t.Fatalf("builtin %s not found", name)
	return nil
}
