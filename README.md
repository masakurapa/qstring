# qstring

`qstring` is a Golang module for generating query strings.

[godoc](https://pkg.go.dev/github.com/masakurapa/qstring) please check here.

## Usage

### Encode the map

```
import (
	"fmt"

	"github.com/masakurapa/qstring"
)

func main() {
	q, err := qstring.Encode(qstring.Q{"key":
		qstring.Q{"a": "value"},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(q) // returns "key%5Ba%5D=value"
}
```

### Encode the struct

The fields of the struct need to be public.
And specify `qstring` for tags.

```
import (
	"fmt"

	"github.com/masakurapa/qstring"
)

func main() {
	type ss struct {
		A string `qstring:"a"`
	}
	type s struct {
		Key ss `qstriger:"key"`
	}

	q, err := qstring.Encode(s{
		Key: ss{
			A: "value",
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(q) // returns "key%5Ba%5D=value"
}
```

### Encode the string

```
import (
	"fmt"

	"github.com/masakurapa/qstring"
)

func main() {
	s := "?key[a]=value"

	q, err := qstring.Encode(s)
	if err != nil {
		panic(err)
	}
	fmt.Println(q) // returns "?key%5Ba%5D=value"
}
```
