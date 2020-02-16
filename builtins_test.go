package tengo

import (
	"errors"
	"reflect"
	"testing"
)

func Test_builtinDelete(t *testing.T) {
	type args struct {
		args []Object
	}
	tests := []struct {
		name      string
		args      args
		want      Object
		wantErr   bool
		wantedErr error
		target    interface{}
	}{
		//Map
		{name: "invalid-arg", args: args{[]Object{&String{}, &String{}}}, wantErr: true,
			wantedErr: ErrInvalidArgumentType{Name: "first", Expected: "map|array", Found: "string"}},
		{name: "no-args", wantErr: true, wantedErr: ErrWrongNumArguments},
		{name: "empty-args", args: args{[]Object{}}, wantErr: true, wantedErr: ErrWrongNumArguments},
		{name: "3-args", args: args{[]Object{(*Map)(nil), (*String)(nil), (*String)(nil)}}, wantErr: true, wantedErr: ErrWrongNumArguments},
		{name: "nil-map-empty-key", args: args{[]Object{&Map{}, &String{}}}, want: UndefinedValue},
		{name: "nil-map-nonstr-key", args: args{[]Object{&Map{}, &Int{}}}, wantErr: true,
			wantedErr: ErrInvalidArgumentType{Name: "second", Expected: "string", Found: "int"}},
		{name: "nil-map-no-key", args: args{[]Object{&Map{}}}, wantErr: true,
			wantedErr: ErrWrongNumArguments},
		{name: "map-missing-key",
			args: args{
				[]Object{
					&Map{Value: map[string]Object{
						"key": &String{Value: "value"},
					}},
					&String{Value: "key1"},
				}},
			want:   UndefinedValue,
			target: &Map{Value: map[string]Object{"key": &String{Value: "value"}}},
		},
		{name: "map-emptied",
			args: args{
				[]Object{
					&Map{Value: map[string]Object{
						"key": &String{Value: "value"},
					}},
					&String{Value: "key"},
				}},
			want:   UndefinedValue,
			target: &Map{Value: map[string]Object{}},
		},
		{name: "map-multi-keys",
			args: args{
				[]Object{
					&Map{Value: map[string]Object{
						"key1": &String{Value: "value1"},
						"key2": &Int{Value: 10},
					}},
					&String{Value: "key1"},
				}},
			want:   UndefinedValue,
			target: &Map{Value: map[string]Object{"key2": &Int{Value: 10}}},
		},
		//Array
		{name: "nil-array-zero-index", args: args{[]Object{&Array{}, &Int{}}}, wantErr: true,
			wantedErr: ErrIndexOutOfBounds},
		{name: "array-str-index", args: args{[]Object{&Array{}, &String{}}}, wantErr: true,
			wantedErr: ErrInvalidArgumentType{Name: "second", Expected: "int", Found: "string"}},
		{name: "array-one", args: args{[]Object{
			&Array{Value: []Object{&Int{Value: 1}}}, &Int{Value: 0}}}, wantErr: false,
			want:   UndefinedValue,
			target: &Array{Value: []Object{}}},
		{name: "array-multi", args: args{[]Object{
			&Array{Value: []Object{&Int{Value: 1}, &String{Value: "xyz"}}}, &Int{Value: 0}}}, wantErr: false,
			want:   UndefinedValue,
			target: &Array{Value: []Object{&String{Value: "xyz"}}}},
		{name: "array-multi2", args: args{[]Object{
			&Array{Value: []Object{&Int{Value: 1}, &String{Value: "xyz"}}}, &Int{Value: 1}}}, wantErr: false,
			want:   UndefinedValue,
			target: &Array{Value: []Object{&Int{Value: 1}}}},
		{name: "array-negative", args: args{[]Object{
			&Array{Value: []Object{&Int{Value: 2}, &String{Value: "xyz"}}}, &Int{Value: -1}}}, wantErr: true,
			wantedErr: ErrIndexOutOfBounds},
		{name: "array-out-of-bounds", args: args{[]Object{
			&Array{Value: []Object{&Int{Value: 3}, &String{Value: "def"}}}, &Int{Value: 2}}}, wantErr: true,
			wantedErr: ErrIndexOutOfBounds},
		{name: "array-out-of-bounds-negative", args: args{[]Object{
			&Array{Value: []Object{&Int{Value: 4}, &String{Value: "ghi"}}}, &Int{Value: -3}}}, wantErr: true,
			wantedErr: ErrIndexOutOfBounds},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := builtinDelete(tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("builtinDelete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && !errors.Is(err, tt.wantedErr) {
				if err.Error() != tt.wantedErr.Error() {
					t.Errorf("builtinDelete() error = %v, wantedErr %v", err, tt.wantedErr)
					return
				}
			}
			if got != tt.want {
				t.Errorf("builtinDelete() = %v, want %v", got, tt.want)
				return
			}
			if !tt.wantErr && tt.target != nil {
				switch v := tt.args.args[0].(type) {
				case *Map, *Array:
					if !reflect.DeepEqual(tt.target, tt.args.args[0]) {
						t.Errorf("builtinDelete() objects are not equal got: %+v, want: %+v", tt.args.args[0], tt.target)
					}
				default:
					t.Errorf("builtinDelete() unsuporrted arg[0] type %s", v.TypeName())
					return
				}
			}
		})
	}
}
