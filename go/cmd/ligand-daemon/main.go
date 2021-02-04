package main

import (
	"fmt"
	"os"

	"github.com/latchai/latch/pkg/infra/grpc"
)

func main() {

	socket := os.Args[1]

	err := grpc.ListenAndServe(socket)
	if err != nil {
		fmt.Errorf("\n%s", err)
	}
}
