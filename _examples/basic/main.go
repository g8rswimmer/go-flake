package main

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
