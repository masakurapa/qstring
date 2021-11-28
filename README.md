# qstringer

qstringer is a Golang module for generating query strings.

## Usage

### Encode the map

```
import "github.com/masakurapa/qstringer"

func main() {
	q, err := qstringer.Encode(qstringer.Q{"key":
		qstringer.Q{"a": "value"},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(q) // returns "?key%5Ba%5D=value"
}
```

### Encode the struct

The fields of the struct need to be public.

```
import "github.com/masakurapa/qstringer"

func main() {
	s := struct {
		Key interface{}
	}{
		Key: struct {
			A string
		}{
			A: "value",
		},
	}

	q, err := qstringer.Encode(s)
	if err != nil {
		panic(err)
	}
	fmt.Println(q) // returns "?key%5Ba%5D=value"
}
```

### Encode the string

```
import "github.com/masakurapa/qstringer"

func main() {
	s := "?key[a]=value"

	q, err := qstringer.Encode(s)
	if err != nil {
		panic(err)
	}
	fmt.Println(q) // returns "?key%5Ba%5D=value"
}
```
