package main

import (
	"github.com/num5/env"
)

func main() {
	_, err := env.Load()
	if err != nil {
		panic(err)
	}

	opts := NewOptions()

	Run(opts)
}