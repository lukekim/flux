package compiler_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/influxdata/flux"
	"github.com/influxdata/flux/compiler"
	"github.com/influxdata/flux/runtime"
	"github.com/influxdata/flux/semantic"
	"github.com/influxdata/flux/semantic/semantictest"
	"github.com/influxdata/flux/values"
)

var CmpOptions = semantictest.CmpOptions

func TestCompileAndEval(t *testing.T) {
	testCases := []struct {
		name           string
		fn             string
		inType         semantic.MonoType
		input          values.Object
		want           values.Value
		wantCompileErr bool
		wantEvalErr    bool
	}{
		{
			name: "interpolated string expression",
			fn:   `(r) => "n = ${r.n}"`,
			inType: semantic.NewObjectType([]semantic.PropertyType{
				{Key: []byte("r"), Value: semantic.NewObjectType([]semantic.PropertyType{
					{Key: []byte("n"), Value: semantic.BasicString},
				})},
			}),
			input: values.NewObjectWithValues(map[string]values.Value{
				"r": values.NewObjectWithValues(map[string]values.Value{
					"n": values.NewString("2"),
				}),
			}),
			want: values.NewString("n = 2"),
		},
		{
			name: "interpolated string expression with int",
			fn:   `(r) => "n = ${r.n}"`,
			inType: semantic.NewObjectType([]semantic.PropertyType{
				{Key: []byte("r"), Value: semantic.NewObjectType([]semantic.PropertyType{
					{Key: []byte("n"), Value: semantic.BasicInt},
				})},
			}),
			input: values.NewObjectWithValues(map[string]values.Value{
				"r": values.NewObjectWithValues(map[string]values.Value{
					"n": values.NewInt(2),
				}),
			}),
			want: values.NewString("n = 2"),
		},
		{
			name: "interpolated string expression with duration type",
			fn:   `(r) => "n = ${r.n}"`,
			inType: semantic.NewObjectType([]semantic.PropertyType{
				{Key: []byte("r"), Value: semantic.NewObjectType([]semantic.PropertyType{
					{Key: []byte("n"), Value: semantic.BasicDuration},
				})},
			}),
			input: values.NewObjectWithValues(map[string]values.Value{
				"r": values.NewObjectWithValues(map[string]values.Value{
					"n": values.NewDuration(flux.ConvertDuration(time.Minute)),
				}),
			}),
			want: values.NewString("n = 1m"),
		},
		{
			name: "interpolated string expression error",
			fn:   `(r) => "n = ${r.n}"`,
			inType: semantic.NewObjectType([]semantic.PropertyType{
				{Key: []byte("r"), Value: semantic.NewObjectType([]semantic.PropertyType{
					{Key: []byte("n"), Value: semantic.BasicBytes},
				})},
			}),
			input: values.NewObjectWithValues(map[string]values.Value{
				"r": values.NewObjectWithValues(map[string]values.Value{
					"n": values.NewBytes([]byte("abc")),
				}),
			}),
			wantEvalErr: true,
		},
		{
			name: "interpolated string expression null",
			fn:   `(r) => "n = ${r.n}"`,
			inType: semantic.NewObjectType([]semantic.PropertyType{
				{Key: []byte("r"), Value: semantic.NewObjectType([]semantic.PropertyType{
					{Key: []byte("n"), Value: semantic.BasicString},
				})},
			}),
			input: values.NewObjectWithValues(map[string]values.Value{
				"r": values.NewObjectWithValues(map[string]values.Value{}),
			}),
			wantEvalErr: true,
		},
		{
			name: "simple ident return",
			fn:   `(r) => r`,
			inType: semantic.NewObjectType([]semantic.PropertyType{
				{Key: []byte("r"), Value: semantic.BasicInt},
			}),
			input: values.NewObjectWithValues(map[string]values.Value{
				"r": values.NewInt(4),
			}),
			want: values.NewInt(4),
		},
		{
			name: "call function",
			fn:   `(r) => ((a,b) => a + b)(a:1, b:r)`,
			inType: semantic.NewObjectType([]semantic.PropertyType{
				{Key: []byte("r"), Value: semantic.BasicInt},
			}),
			input: values.NewObjectWithValues(map[string]values.Value{
				"r": values.NewInt(4),
			}),
			want: values.NewInt(5),
		},
		{
			name: "call function with defaults",
			fn:   `(r) => ((a=0,b) => a + b)(b:r)`,
			inType: semantic.NewObjectType([]semantic.PropertyType{
				{Key: []byte("r"), Value: semantic.BasicInt},
			}),
			input: values.NewObjectWithValues(map[string]values.Value{
				"r": values.NewInt(4),
			}),
			want: values.NewInt(4),
		},
		{
			name: "call function via identifier",
			fn: `(r) => {
				f = (a,b) => a + b
				return f(a:1, b:r)
			}`,
			inType: semantic.NewObjectType([]semantic.PropertyType{
				{Key: []byte("r"), Value: semantic.BasicInt},
			}),
			input: values.NewObjectWithValues(map[string]values.Value{
				"r": values.NewInt(4),
			}),
			want: values.NewInt(5),
		},
		{
			name: "call function via identifier with different types",
			fn: `(r) => {
				i = (x) => x
				return i(x:i)(x:r+1)
			}`,
			inType: semantic.NewObjectType([]semantic.PropertyType{
				{Key: []byte("r"), Value: semantic.BasicInt},
			}),
			input: values.NewObjectWithValues(map[string]values.Value{
				"r": values.NewInt(4),
			}),
			want: values.NewInt(5),
		},
		{
			name: "call filter function with index expression",
			fn:   `(r) => r[2] == 3`,
			inType: semantic.NewObjectType([]semantic.PropertyType{
				{Key: []byte("r"), Value: semantic.NewArrayType(semantic.BasicInt)},
			}),
			input: values.NewObjectWithValues(map[string]values.Value{
				"r": values.NewArrayWithBacking(semantic.NewArrayType(semantic.BasicInt), []values.Value{
					values.NewInt(5),
					values.NewInt(6),
					values.NewInt(3),
				}),
			}),
			want: values.NewBool(true),
		},
		{
			name: "call filter function with complex index expression",
			fn:   `(r) => r[((x) => 2)(x: "anything")] == 3`,
			inType: semantic.NewObjectType([]semantic.PropertyType{
				{Key: []byte("r"), Value: semantic.NewArrayType(semantic.BasicInt)},
			}),
			input: values.NewObjectWithValues(map[string]values.Value{
				"r": values.NewArrayWithBacking(semantic.NewArrayType(semantic.BasicInt), []values.Value{
					values.NewInt(5),
					values.NewInt(6),
					values.NewInt(3),
				}),
			}),
			want: values.NewBool(true),
		},
		{
			name: "call with pipe argument",
			fn: `(n) => {
				f = (v=<-) => v + n
				return 5 |> f()
			}`,
			inType: semantic.NewObjectType([]semantic.PropertyType{
				{Key: []byte("n"), Value: semantic.BasicInt},
			}),
			input: values.NewObjectWithValues(map[string]values.Value{
				"n": values.NewInt(6),
			}),
			want: values.NewInt(11),
		},
		{
			name: "conditional",
			fn:   `(t, c, a) => if t then c else a`,
			inType: semantic.NewObjectType([]semantic.PropertyType{
				{Key: []byte("t"), Value: semantic.BasicBool},
				{Key: []byte("c"), Value: semantic.BasicString},
				{Key: []byte("a"), Value: semantic.BasicString},
			}),
			input: values.NewObjectWithValues(map[string]values.Value{
				"t": values.NewBool(true),
				"c": values.NewString("cats"),
				"a": values.NewString("dogs"),
			}),
			want: values.NewString("cats"),
		},
		{
			name: "unary logical operator - not",
			fn:   `(a, b) => not a or b`,
			inType: semantic.NewObjectType([]semantic.PropertyType{
				{Key: []byte("a"), Value: semantic.BasicBool},
				{Key: []byte("b"), Value: semantic.BasicBool},
			}),
			input: values.NewObjectWithValues(map[string]values.Value{
				"a": values.NewBool(true),
				"b": values.NewBool(true),
			}),
			want: values.NewBool(true),
		},
		{
			name: "unary logical operator - exists with null",
			fn:   `(a) => exists a`,
			inType: semantic.NewObjectType([]semantic.PropertyType{
				{Key: []byte("a"), Value: semantic.BasicString},
			}),
			input: values.NewObjectWithValues(map[string]values.Value{
				"a": values.NewNull(semantic.BasicString),
			}),
			want: values.NewBool(false),
		},
		{
			name: "unary logical operator - exists without null",
			fn:   `(a, b) => not a and exists b`,
			inType: semantic.NewObjectType([]semantic.PropertyType{
				{Key: []byte("a"), Value: semantic.BasicBool},
				{Key: []byte("b"), Value: semantic.BasicString},
			}),
			input: values.NewObjectWithValues(map[string]values.Value{
				"a": values.NewBool(true),
				"b": values.NewString("I exist"),
			}),
			want: values.NewBool(false),
		},
		{
			name: "unary operator",
			fn:   `(a) => if a < 0 then -a else +a`,
			inType: semantic.NewObjectType([]semantic.PropertyType{
				{Key: []byte("a"), Value: semantic.BasicInt},
			}),
			input: values.NewObjectWithValues(map[string]values.Value{
				"a": values.NewInt(5),
			}),
			want: values.NewInt(5),
		},
		{
			name: "filter with member expression",
			fn:   `(r) => r.m == "cpu"`,
			inType: semantic.NewObjectType([]semantic.PropertyType{
				{Key: []byte("r"), Value: semantic.NewObjectType([]semantic.PropertyType{
					{Key: []byte("m"), Value: semantic.BasicString},
				})},
			}),
			input: values.NewObjectWithValues(map[string]values.Value{
				"r": values.NewObjectWithValues(map[string]values.Value{
					"m": values.NewString("cpu"),
				}),
			}),
			want: values.NewBool(true),
		},
		{
			name: "regex literal filter",
			fn:   `(r) => r =~ /^(c|g)pu$/`,
			inType: semantic.NewObjectType([]semantic.PropertyType{
				{Key: []byte("r"), Value: semantic.BasicString},
			}),
			input: values.NewObjectWithValues(map[string]values.Value{
				"r": values.NewString("cpu"),
			}),
			want: values.NewBool(true),
		},
		{
			name: "block statement with conditional",
			fn: `(r) => {
				v = if r < 0 then -r else r
				return v * v
			}`,
			inType: semantic.NewObjectType([]semantic.PropertyType{
				{Key: []byte("r"), Value: semantic.BasicInt},
			}),
			input: values.NewObjectWithValues(map[string]values.Value{
				"r": values.NewInt(-3),
			}),
			want: values.NewInt(9),
		},
		{
			name:   "array literal",
			fn:     `() => [1.0, 2.0, 3.0]`,
			inType: semantic.NewObjectType(nil),
			input:  values.NewObjectWithValues(nil),
			want: values.NewArrayWithBacking(
				semantic.NewArrayType(semantic.BasicFloat),
				[]values.Value{
					values.NewFloat(1),
					values.NewFloat(2),
					values.NewFloat(3),
				},
			),
		},
		{
			name: "dict literal",
			fn: `() => {
				a = "a"
				b = "b"
				return [a: "a", b: "b"]
			}`,
			inType: semantic.NewObjectType(nil),
			input:  values.NewObjectWithValues(nil),
			want: func() values.Value {
				builder := values.NewDictBuilder(semantic.NewDictType(semantic.BasicString, semantic.BasicString))
				builder.Insert(values.NewString("a"), values.NewString("a"))
				builder.Insert(values.NewString("b"), values.NewString("b"))
				return builder.Dict()
			}(),
		},
		{
			name: "array access",
			fn:   `(values) => values[0]`,
			inType: semantic.NewObjectType([]semantic.PropertyType{
				{Key: []byte("values"), Value: semantic.NewArrayType(semantic.BasicFloat)},
			}),
			input: values.NewObjectWithValues(map[string]values.Value{
				"values": values.NewArrayWithBacking(
					semantic.NewArrayType(semantic.BasicFloat),
					[]values.Value{
						values.NewFloat(1),
						values.NewFloat(2),
						values.NewFloat(3),
					},
				),
			}),
			want: values.NewFloat(1),
		},
		{
			name: "array access out of bounds low",
			fn:   `(values) => values[-1]`,
			inType: semantic.NewObjectType([]semantic.PropertyType{
				{Key: []byte("values"), Value: semantic.NewArrayType(semantic.BasicFloat)},
			}),
			input: values.NewObjectWithValues(map[string]values.Value{
				"values": values.NewArrayWithBacking(
					semantic.NewArrayType(semantic.BasicFloat),
					[]values.Value{
						values.NewFloat(1),
						values.NewFloat(2),
						values.NewFloat(3),
					},
				),
			}),
			wantEvalErr: true,
		},
		{
			name: "array access out of bounds high",
			fn:   `(values) => values[3]`,
			inType: semantic.NewObjectType([]semantic.PropertyType{
				{Key: []byte("values"), Value: semantic.NewArrayType(semantic.BasicFloat)},
			}),
			input: values.NewObjectWithValues(map[string]values.Value{
				"values": values.NewArrayWithBacking(
					semantic.NewArrayType(semantic.BasicFloat),
					[]values.Value{
						values.NewFloat(1),
						values.NewFloat(2),
						values.NewFloat(3),
					},
				),
			}),
			wantEvalErr: true,
		},
		{
			name: "logical expression",
			fn:   `(a, b) => a or b`,
			inType: semantic.NewObjectType([]semantic.PropertyType{
				{Key: []byte("a"), Value: semantic.BasicBool},
				{Key: []byte("b"), Value: semantic.BasicBool},
			}),
			input: values.NewObjectWithValues(map[string]values.Value{
				"a": values.NewBool(true),
				"b": values.NewBool(false),
			}),
			want: values.NewBool(true),
		},
		{
			name: "call with nonexistant value",
			fn:   `(r) => r.a + r.b`,
			inType: semantic.NewObjectType([]semantic.PropertyType{
				{Key: []byte("r"), Value: semantic.NewObjectType([]semantic.PropertyType{
					{Key: []byte("a"), Value: semantic.BasicInt},
				})},
			}),
			input: values.NewObjectWithValues(map[string]values.Value{
				"r": values.NewObjectWithValues(map[string]values.Value{
					"a": values.NewInt(4),
				}),
			}),
			want: values.Null,
		},
		{
			name: "call with null value",
			fn:   `(r) => r.a + r.b`,
			inType: semantic.NewObjectType([]semantic.PropertyType{
				{Key: []byte("r"), Value: semantic.NewObjectType([]semantic.PropertyType{
					{Key: []byte("a"), Value: semantic.BasicInt},
					{Key: []byte("b"), Value: semantic.BasicInt},
				})},
			}),
			input: values.NewObjectWithValues(map[string]values.Value{
				"r": values.NewObjectWithValues(map[string]values.Value{
					"a": values.NewInt(4),
					// The object is typed as an integer,
					// but it doesn't have an actual value
					// because it is null.
					"b": values.Null,
				}),
			}),
			want: values.Null,
		},
		{
			name: "call with null parameter",
			fn: `(r) => {
				eval = (a, b) => a + b
				return eval(a: r.a, b: r.b)
			}`,
			inType: semantic.NewObjectType([]semantic.PropertyType{
				{Key: []byte("r"), Value: semantic.NewObjectType([]semantic.PropertyType{
					{Key: []byte("a"), Value: semantic.BasicInt},
					// "b": semantic.Nil,
				})},
			}),
			input: values.NewObjectWithValues(map[string]values.Value{
				"r": values.NewObjectWithValues(map[string]values.Value{
					"a": values.NewInt(4),
					// "b": values.Null,
				}),
			}),
			want: values.Null,
		},
		{
			name: "return nonexistant value",
			fn:   `(r) => r.b`,
			inType: semantic.NewObjectType([]semantic.PropertyType{
				{Key: []byte("r"), Value: semantic.NewObjectType([]semantic.PropertyType{
					{Key: []byte("a"), Value: semantic.BasicInt},
				})},
			}),
			input: values.NewObjectWithValues(map[string]values.Value{
				"r": values.NewObjectWithValues(map[string]values.Value{
					"a": values.NewInt(4),
				}),
			}),
			want: values.Null,
		},
		{
			name: "return nonexistant and used parameter",
			fn: `(r) => {
				b = (r) => r.b
				return r.a + b(r: r)
			}`,
			inType: semantic.NewObjectType([]semantic.PropertyType{
				{Key: []byte("r"), Value: semantic.NewObjectType([]semantic.PropertyType{
					{Key: []byte("a"), Value: semantic.BasicInt},
					// "b": semantic.Nil,
				})},
			}),
			input: values.NewObjectWithValues(map[string]values.Value{
				"r": values.NewObjectWithValues(map[string]values.Value{
					"a": values.NewInt(4),
					// "b": values.Null,
				}),
			}),
			want: values.Null,
		},
		{
			name: "two null values are not equal",
			fn:   `(a, b) => a == b`,
			inType: semantic.NewObjectType([]semantic.PropertyType{
				{Key: []byte("a"), Value: semantic.BasicInt},
				{Key: []byte("b"), Value: semantic.BasicInt},
			}),
			input: values.NewObjectWithValues(map[string]values.Value{
				"a": values.Null,
				"b": values.Null,
			}),
			want: values.Null,
		},
		{
			name: "superseding record field type",
			fn: `
				(str) => {
					m = (s) => ({s with v: 10.0})
					f = (t=<-) => t.v == 10.0
					return m(s: {v: str}) |> f()
				}`,
			inType: semantic.NewObjectType([]semantic.PropertyType{
				{Key: []byte("str"), Value: semantic.BasicString},
			}),
			input: values.NewObjectWithValues(map[string]values.Value{
				"str": values.NewString("foo"),
			}),
			want: values.NewBool(true),
		},
		{
			name: "null array",
			fn:   `(r) => r.a[0]`,
			inType: semantic.NewObjectType([]semantic.PropertyType{
				{Key: []byte("r"), Value: semantic.NewObjectType([]semantic.PropertyType{
					{Key: []byte("a"), Value: semantic.NewArrayType(semantic.BasicString)},
				})},
			}),
			input: values.NewObjectWithValues(map[string]values.Value{
				"r": values.NewObjectWithValues(nil),
			}),
			wantEvalErr: true,
		},
		{
			name: "null record",
			fn:   `(r) => r.a["b"]`,
			inType: semantic.NewObjectType([]semantic.PropertyType{
				{Key: []byte("r"), Value: semantic.NewObjectType([]semantic.PropertyType{
					{Key: []byte("a"), Value: semantic.NewObjectType([]semantic.PropertyType{
						{Key: []byte("b"), Value: semantic.BasicString},
					})},
				})},
			}),
			input: values.NewObjectWithValues(map[string]values.Value{
				"r": values.NewObjectWithValues(nil),
			}),
			wantEvalErr: true,
		},
		// TODO(jsternberg): We presently have not implemented dictionary support for
		// runtime functions. There aren't any builtins that use this functionality,
		// but when we do, this test will need to be uncommented to ensure that
		// a null dictionary does not sneak in.
		// {
		// 	name: "null dict",
		//	fn: `import "dict"
		// (r) => dict.get(dict: r.a, key: "b", default: "")`,
		//	inType: semantic.NewObjectType([]semantic.PropertyType{
		//		{Key: []byte("r"), Value: semantic.NewObjectType([]semantic.PropertyType{
		//			{Key: []byte("a"), Value: semantic.NewDictType(semantic.BasicString, semantic.BasicString)},
		//		})},
		//	}),
		//	input: values.NewObjectWithValues(map[string]values.Value{
		//		"r": values.NewObjectWithValues(nil),
		//	}),
		//	wantEvalErr: true,
		// },
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			pkg, err := runtime.AnalyzeSource(tc.fn)
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			stmt := pkg.Files[0].Body[0].(*semantic.ExpressionStatement)
			fn := stmt.Expression.(*semantic.FunctionExpression)
			f, err := compiler.Compile(nil, fn, tc.inType)
			if err != nil {
				if !tc.wantCompileErr {
					t.Fatalf("unexpected error: %s", err)
				}
				return
			} else if tc.wantCompileErr {
				t.Fatal("wanted error but got nothing")
			}

			// ctx := dependenciestest.Default().Inject(context.Background())
			got, err := f.Eval(context.TODO(), tc.input)
			if tc.wantEvalErr != (err != nil) {
				t.Errorf("unexpected error: %s", err)
			}

			if !cmp.Equal(tc.want, got, CmpOptions...) {
				t.Errorf("unexpected value -want/+got\n%s", cmp.Diff(tc.want, got, CmpOptions...))
			}
		})
	}
}

func TestCompiler_ReturnType(t *testing.T) {
	testCases := []struct {
		name   string
		fn     string
		inType semantic.MonoType
		want   string
	}{
		{
			name: "with",
			fn:   `(r) => ({r with _value: r._value * 2.0})`,
			inType: semantic.NewObjectType([]semantic.PropertyType{
				{Key: []byte("r"), Value: semantic.NewObjectType([]semantic.PropertyType{
					{Key: []byte("_value"), Value: semantic.BasicFloat},
					{Key: []byte("_time"), Value: semantic.BasicTime},
				})},
			}),
			want: `{_time: time, _value: float, _value: float}`,
		},
		{
			name: "array access",
			fn:   `(values) => values[0]`,
			inType: semantic.NewObjectType([]semantic.PropertyType{
				{Key: []byte("values"), Value: semantic.NewArrayType(semantic.BasicFloat)},
			}),
			want: `float`,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			pkg, err := runtime.AnalyzeSource(tc.fn)
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			stmt := pkg.Files[0].Body[0].(*semantic.ExpressionStatement)
			fn := stmt.Expression.(*semantic.FunctionExpression)
			f, err := compiler.Compile(nil, fn, tc.inType)
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if got, want := f.Type().String(), tc.want; got != want {
				t.Fatalf("unexpected return type -want/+got:\n\t- %s\n\t+ %s", want, got)
			}
		})
	}
}

func TestToScopeNil(t *testing.T) {
	if compiler.ToScope(nil) != nil {
		t.Fatal("ToScope made non-nil scope from a nil base")
	}
}
