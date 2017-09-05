package main

import (
	"github.com/lburgazzoli/camel-go/camel"
	"fmt"
)

type LogComponent struct {
    context camel.Context
}

func (component *LogComponent) Start() {
}

func (component *LogComponent) Stop() {
}

func (component *LogComponent) Process(message string) {
    fmt.Printf("%s\n", message)
}

// ========================================
// plugin entry-pooint
// ========================================

func CreateComponent(context camel.Context) (camel.Component) {
	return &LogComponent{
        context: context,
    }
}
