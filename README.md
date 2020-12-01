# go-flake

This is a simple go implementation of the twitter [snowflake](https://blog.twitter.com/engineering/en_us/a/2010/announcing-snowflake.html) id generator.

The id breakdown is:
* 48 bits of a epoch time in milliseconds 
* 5 bits of a worker id 
* 12 bits of sequence

There are two ways to use the generator
* straight call to generate the id
    * this method will allow the user to pass the worker and seqeunce
* use the generator struct 
    * this is for one worker id and will handle auto update of the sequence and check for rollover
    * Errors can help with identifying possible id collision

## Examples
The following are the two implementation examples.

### Function Generation
```go

import (
	"fmt"

	"github.com/g8rswimmer/go-flake"
)

const (
	worker = uint64(10)
	seq    = uint64(1000)
)

func main() {
	id, err := flake.Generate(worker, seq)
	if err != nil {
		panic(err)
	}
	fmt.Printf("The generated id %v\n", id)
}
```

### Generator 
```go
package main

import (
	"fmt"

	"github.com/g8rswimmer/go-flake"
)

const (
	worker = uint64(10)
)

func main() {
	gen, err := flake.New(worker)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 10; i++ {
		id, err := gen.Generate()
		if err != nil {
			panic(err)
		}
		fmt.Printf("The generated id %v\n", id)
	}
}
```