package main

import (
	"github.com/lburgazzoli/camel-go/core"
	"fmt"
)

func main() {
    context := core.NewCamelContext()
    context.Start()

    comp, err := context.GetComponent("log")
    if err != nil {
		fmt.Printf("Failed to get log component: %v\n", err)
	}
	if comp != nil {
        comp.Process("message")
    }

    context.Stop()
}
