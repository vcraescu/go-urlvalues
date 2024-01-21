package urlvalues_test

import (
	"fmt"
	"github.com/vcraescu/go-urlvalues"
)

func ExampleMarshal() {
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
