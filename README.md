# qstringer

qstringer is a Golang module for generating query strings.

## Quick Start

```
import "github.com/masakurapa/qstringer"

func main() {
	q, err := qstringer.Encode(qstringer.Q{"key": "value"})
	if err != nil {
		panic(err)
	}
	fmt.Println(q) // returns "?key=value
}
```

or any struct can be used.

```
import "github.com/masakurapa/qstringer"

func main() {
	s := struct {
		Key string
	}{
		Key: "string,
	}

	q, err := qstringer.Encode(s)
	if err != nil {
		panic(err)
	}
	fmt.Println(q) // returns "?key=value
}
```
