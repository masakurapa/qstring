# qstringer

qstringer is a Golang module for generating query strings.

## Quick Start

```
import "github.com/masakurapa/qstringer"

func main() {
	q := qstringer.Q{"key": "value"}
	qs, err := q.Encode()
	if err != nil {
		panic(err)
	}
	fmt.Println(qs) // returns "?key=value
}
```
