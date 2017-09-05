package main

import (
	"github.com/lburgazzoli/camel-go/camel"
	"fmt"
)

type LogComponent struct {
}

func (component *LogComponent) Init(context camel.Context) error {
	fmt.Println("Initialize")
    return nil
}

func (component *LogComponent) Start() {
}

func (component *LogComponent) Stop() {
}

func (component *LogComponent) Process(message string) {
    fmt.Printf("%s\n", message)
}

func Component() (camel.Component, error) {
	return &LogComponent{}, nil
}