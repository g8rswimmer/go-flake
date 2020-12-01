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
