package urlvalues_test

import (
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vcraescu/go-urlvalues"
)

var now = time.Now()

func String(s string) *string {
	return &s
}

type marshalArgs struct {
	in   any
	opts urlvalues.MarshalerOptions
}

func TestMarshal_Struct(t *testing.T) {
	t.Parallel()

	type Nested struct {
		String      string   `url:"string"`
		StringSlice []string `url:"stringSlice"`
		IntSlice    []int    `url:"intSlice"`
		Int         int      `url:"int"`
	}

	type Object struct {
		String              string   `url:"string,omitempty"`
		Object              *Nested  `url:"object,omitempty"`
		StringSlice         []string `url:"stringSlice,omitempty"`
		StringSliceBrackets []string `url:"stringSliceBrackets,omitempty,brackets"`
		IntSlice            []int    `url:"intSlice,omitempty"`
		Int                 int      `url:"int,omitempty"`
	}

	tests := []struct {
		name    string
		args    marshalArgs
		want    url.Values
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "struct",
			args: marshalArgs{
				in: Object{
					String:              "string-value",
					IntSlice:            []int{101, 102, 103},
					Int:                 100,
					StringSlice:         []string{"101", "102", "103"},
					StringSliceBrackets: []string{"101", "102", "103"},
					Object: &Nested{
						String:      "nested-string-value",
						StringSlice: []string{"201", "202", "203"},
						IntSlice:    []int{201, 202, 203},
						Int:         200,
					},
				},
			},
			want: url.Values{
				"string":                []string{"string-value"},
				"intSlice":              []string{"101", "102", "103"},
				"int":                   []string{"100"},
				"stringSlice":           []string{"101", "102", "103"},
				"stringSliceBrackets[]": []string{"101", "102", "103"},
				"object[string]":        []string{"nested-string-value"},
				"object[stringSlice]":   []string{"201", "202", "203"},
				"object[intSlice]":      []string{"201", "202", "203"},
				"object[int]":           []string{"200"},
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			testMarshal(t, tt.args, tt.wantErr, tt.want)
		})
	}
}

func TestMarshal_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		args    marshalArgs
		want    url.Values
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "empty string",
			args: marshalArgs{
				in: "",
			},
		},
		{
			name: "string value",
			args: marshalArgs{
				in: url.Values{
					"seriesName":  []string{"123"},
					"client[cif]": []string{"0123"},
					"precision":   []string{"12"},
				}.Encode(),
			},
			want: url.Values{
				"seriesName":  []string{"123"},
				"client[cif]": []string{"0123"},
				"precision":   []string{"12"},
			},
		},
		{
			name: "string pointer",
			args: marshalArgs{
				in: String(url.Values{
					"seriesName":  []string{"123"},
					"client[cif]": []string{"0123"},
					"precision":   []string{"12"},
				}.Encode()),
			},
			want: url.Values{
				"seriesName":  []string{"123"},
				"client[cif]": []string{"0123"},
				"precision":   []string{"12"},
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			testMarshal(t, tt.args, tt.wantErr, tt.want)
		})
	}
}

func TestMarshal_BasicTypes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		args    marshalArgs
		want    url.Values
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "bool as true/false",
			args: marshalArgs{
				in: map[string]any{
					"trueBool":  true,
					"falseBool": false,
				},
			},
			want: url.Values{
				"trueBool":  []string{"true"},
				"falseBool": []string{"false"},
			},
		},
		{
			name: "int bool",
			args: marshalArgs{
				in: map[string]any{
					"trueBool":  true,
					"falseBool": false,
				},
				opts: urlvalues.MarshalerOptions{
					IntBool: true,
				},
			},
			want: url.Values{
				"trueBool":  []string{"1"},
				"falseBool": []string{"0"},
			},
		},
		{
			name: "unix time",
			args: marshalArgs{
				in: map[string]any{
					"time": now,
				},
				opts: urlvalues.MarshalerOptions{
					TimeUnix: true,
				},
			},
			want: url.Values{
				"time": []string{fmt.Sprint(now.Unix())},
			},
		},
		{
			name: "unix milli time",
			args: marshalArgs{
				in: map[string]any{
					"time": now,
				},
				opts: urlvalues.MarshalerOptions{
					TimeUnixMilli: true,
				},
			},
			want: url.Values{
				"time": []string{fmt.Sprint(now.UnixMilli())},
			},
		},
		{
			name: "unix nano time",
			args: marshalArgs{
				in: map[string]any{
					"time": now,
				},
				opts: urlvalues.MarshalerOptions{
					TimeUnixNano: true,
				},
			},
			want: url.Values{
				"time": []string{fmt.Sprint(now.UnixNano())},
			},
		},
		{
			name: "time with layout",
			args: marshalArgs{
				in: map[string]any{
					"time": now,
				},
				opts: urlvalues.MarshalerOptions{
					TimeLayout: time.RFC822Z,
				},
			},
			want: url.Values{
				"time": []string{now.Format(time.RFC822Z)},
			},
		},
		{
			name: "time with default layout",
			args: marshalArgs{
				in: map[string]any{
					"time": now,
				},
			},
			want: url.Values{
				"time": []string{now.Format(time.RFC3339)},
			},
		},
		{
			name: "time pointer",
			args: marshalArgs{
				in: map[string]any{
					"time": &now,
				},
			},
			want: url.Values{
				"time": []string{now.Format(time.RFC3339)},
			},
		},
		{
			name: "zero time",
			args: marshalArgs{
				in: map[string]any{
					"time": time.Time{},
				},
			},
			want: url.Values{
				"time": []string{time.Time{}.Format(time.RFC3339)},
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			testMarshal(t, tt.args, tt.wantErr, tt.want)
		})
	}
}

func TestMarshal_Slices(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		args    marshalArgs
		want    url.Values
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "empty slice",
			args: marshalArgs{
				in: map[string]any{
					"slice": []string{},
				},
			},
			want: url.Values{},
		},
		{
			name: "nil slice",
			args: marshalArgs{
				in: map[string]any{
					"slice": nil,
				},
			},
			want: url.Values{},
		},
		{
			name: "default slice",
			args: marshalArgs{
				in: map[string]any{
					"slice": []string{"1", "2"},
				},
			},
			want: url.Values{
				"slice": []string{"1", "2"},
			},
		},
		{
			name: "slice with delimiter",
			args: marshalArgs{
				in: map[string]any{
					"slice": []string{"1", "2"},
				},
				opts: urlvalues.MarshalerOptions{
					ArrayDelimiter: "|",
				},
			},
			want: url.Values{
				"slice": []string{"1|2"},
			},
		},
		{
			name: "slice with brackets",
			args: marshalArgs{
				in: map[string]any{
					"slice": []string{"1", "2"},
				},
				opts: urlvalues.MarshalerOptions{
					ArrayBrackets: true,
				},
			},
			want: url.Values{
				"slice[]": []string{"1", "2"},
			},
		},
		{
			name: "slice with empty values",
			args: marshalArgs{
				in: map[string]any{
					"slice": []string{"", "1"},
				},
			},
			want: url.Values{
				"slice": []string{"", "1"},
			},
		},
		{
			name: "empty array",
			args: marshalArgs{
				in: map[string]any{
					"array": [1]string{},
				},
			},
			want: url.Values{
				"array": []string{""},
			},
		},
		{
			name: "default array",
			args: marshalArgs{
				in: map[string]any{
					"array": [2]string{"1", "2"},
				},
			},
			want: url.Values{
				"array": []string{"1", "2"},
			},
		},
		{
			name: "array with delimiter",
			args: marshalArgs{
				in: map[string]any{
					"array": [2]string{"1", "2"},
				},
				opts: urlvalues.MarshalerOptions{
					ArrayDelimiter: "|",
				},
			},
			want: url.Values{
				"array": []string{"1|2"},
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			testMarshal(t, tt.args, tt.wantErr, tt.want)
		})
	}
}

func TestMarshal_NestedTypes(t *testing.T) {
	t.Parallel()

	type Object struct {
		Int   int   `url:"int"`
		Slice []int `url:"slice"`
	}

	tests := []struct {
		name    string
		args    marshalArgs
		want    url.Values
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "empty map",
			args: marshalArgs{
				in: map[string]any{
					"map": map[string]any{},
				},
			},
			want: url.Values{},
		},
		{
			name: "nil value",
			args: marshalArgs{
				in: map[string]any{
					"map": nil,
				},
			},
			want: url.Values{},
		},
		{
			name: "nested map",
			args: marshalArgs{
				in: map[string]any{
					"map1": map[string]any{
						"map2": map[string]any{
							"int":   100,
							"slice": []int{1, 2},
						},
						"object": &Object{
							Int:   100,
							Slice: []int{1, 2},
						},
					},
					"object": Object{
						Int:   100,
						Slice: []int{1, 2},
					},
				},
			},
			want: url.Values{
				"map1[map2][int]":     []string{"100"},
				"map1[map2][slice]":   []string{"1", "2"},
				"map1[object][slice]": []string{"1", "2"},
				"map1[object][int]":   []string{"100"},
				"object[int]":         []string{"100"},
				"object[slice]":       []string{"1", "2"},
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			testMarshal(t, tt.args, tt.wantErr, tt.want)
		})
	}
}

func TestMarshal_InvalidValues(t *testing.T) {
	t.Parallel()

	var p *int

	tests := []struct {
		name    string
		args    marshalArgs
		want    url.Values
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "nil",
			args: marshalArgs{
				in: nil,
			},
			want: url.Values{},
		},
		{
			name: "empty string",
			args: marshalArgs{
				in: "",
			},
		},
		{
			name: "nil pointer",
			args: marshalArgs{
				in: p,
			},
		},
		{
			name: "non-zero int",
			args: marshalArgs{
				in: 10,
			},
			wantErr: func(t assert.TestingT, err error, _ ...any) bool {
				return assert.Error(t, err)
			},
		},
		{
			name: "true",
			args: marshalArgs{
				in: true,
			},
			wantErr: func(t assert.TestingT, err error, _ ...any) bool {
				return assert.Error(t, err)
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			testMarshal(t, tt.args, tt.wantErr, tt.want)
		})
	}
}

var _ urlvalues.Marshaler = (*CustomEncodingInt)(nil)

type CustomEncodingInt int

func (e CustomEncodingInt) EncodeValues(key string, values *url.Values) error {
	if e < 0 {
		return fmt.Errorf("invalid value: %d", e)
	}

	values.Set(key, fmt.Sprintf("$%d", e))

	return nil
}

var _ urlvalues.Marshaler = (*CustomEncodingIntPtr)(nil)

type CustomEncodingIntPtr int

func (e *CustomEncodingIntPtr) EncodeValues(key string, values *url.Values) error {
	values.Set(key, fmt.Sprintf("$%d", *e))

	return nil
}

func TestMarshal_EncodedValues(t *testing.T) {
	t.Parallel()

	i := CustomEncodingIntPtr(10)

	tests := []struct {
		name    string
		args    marshalArgs
		want    url.Values
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "int",
			args: marshalArgs{
				in: map[string]any{
					"int": CustomEncodingInt(10),
				},
			},
			want: url.Values{
				"int": []string{"$10"},
			},
		},
		{
			name: "int pointer",
			args: marshalArgs{
				in: map[string]any{
					"int": &i,
				},
			},
			want: url.Values{
				"int": []string{"$10"},
			},
		},
		{
			name: "encoding returns error",
			args: marshalArgs{
				in: map[string]any{
					"int": CustomEncodingInt(-10),
				},
			},
			wantErr: func(t assert.TestingT, err error, _ ...any) bool {
				return assert.Error(t, err)
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			testMarshal(t, tt.args, tt.wantErr, tt.want)
		})
	}
}

func testMarshal(t *testing.T, args marshalArgs, wantErr assert.ErrorAssertionFunc, want url.Values) {
	t.Parallel()

	got, err := args.opts.Marshal(args.in)

	if wantErr != nil {
		wantErr(t, err)

		return
	}

	require.NoError(t, err)
	require.Equal(t, want, got)
}
