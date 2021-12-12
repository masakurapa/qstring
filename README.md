# qstringer

qstringer is a Golang module for generating query strings.

[godoc](https://pkg.go.dev/github.com/masakurapa/qstringer) please check here.

## Usage

### Encode the map

```
import (
	"fmt"

	"github.com/masakurapa/qstringer"
)

func main() {
	q, err := qstringer.Encode(qstringer.Q{"key":
		qstringer.Q{"a": "value"},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(q) // returns "key%5Ba%5D=value"
}
```

### Encode the struct

The fields of the struct need to be public.
And specify `qstringer` for tags.

```
import (
	"fmt"

	"github.com/masakurapa/qstringer"
)

func main() {
	type ss struct {
		A string `qstringer:"a"`
	}
	type s struct {
		Key ss `qstriger:"key"`
	}

	q, err := qstringer.Encode(s{
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

	"github.com/masakurapa/qstringer"
)

func main() {
	s := "?key[a]=value"

	q, err := qstringer.Encode(s)
	if err != nil {
		panic(err)
	}
	fmt.Println(q) // returns "?key%5Ba%5D=value"
}
```
