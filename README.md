# go-urlvalues #

![Test Status](https://github.com/vcraescu/go-urlvalues/actions/workflows/go.yml/badge.svg)
[![codecov](https://codecov.io/github/vcraescu/go-urlvalues/branch/master/graph/badge.svg)](https://codecov.io/github/vcraescu/go-urlvalues)

Go library that provides a simple and flexible way to encode Go values into URL query parameters. It
offers a customizable interface to marshal various data types, including maps, slices, structs, and more.

## Usage

```bash
go get -u github.com/vcraescu/go-urlvalues
```

## Marshaler

### Features

- Encode Go values (structs and maps) into URL query parameters.
- Support for encoding maps, slices, structs, and other data types.
- Configurable options for customizing the marshaling process.

### Marshaler Interface

```
type Marshaler interface {
	EncodeValues(key string, values *url.Values) error
}
```

Implement the `EncodeValues` method in your custom types to provide specific encoding logic.

### Custom Options

#### MarshalerOptions

* `ArrayBrackets`: Add brackets to array/slice keys in the URL parameters.
* `ArrayDelimiter`: Specify a custom delimiter for joining slice values in the URL parameters.
* `IntBool`: When set to true, represents boolean values as integers (0 or 1).
* `TimeUnix`: Time encoded in seconds.
* `TimeUnixMilli`: Time encoded in milliseconds.
* `TimeUnixNano`: Time encoded in nanoseconds.
* `TimeLayout`: Allows specifying a custom layout for time values if provided.

Example:

```
marshaler := urlvalues.MarshalerOptions{
        ArrayBrackets: true
        ArrayDelimiter: "|"
}

values, err := marshaler.Marshal(data)
```

### Example

```go
package main

import (
	"fmt"
	"github.com/vcraescu/go-urlvalues"
)

func main() {
	type Object struct {
		Slice  []int  `url:"slice"`
		String string `url:"string"`
	}

	m := map[string]any{
		"slice": []string{"100", "200"},
		"int":   1,
		"map": map[string]any{
			"slice": []string{"100", "200"},
			"int":   1,
		},
		"object": Object{
			Slice:  []int{1, 2},
			String: "example",
		},
	}

	values, err := urlvalues.Marshal(m)
	if err != nil {
		panic(err)
	}

	fmt.Println(values.Encode())
	// Output: int=1&map%5Bint%5D=1&map%5Bslice%5D=100&map%5Bslice%5D=200&object%5Bslice%5D=1&object%5Bslice%5D=2&object%5Bstring%5D=example&slice=100&slice=200
}
```

## Credits

This project utilizes the [google/go-querystring](https://github.com/google/go-querystring) library for encoding structs
to URL.

## License

This library is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

