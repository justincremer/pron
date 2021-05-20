package main

import (
	"fmt"
	"time"

	pron "github.com/justincremer/pron/src"
)

const (
	t          time.Duration = time.Second
	configFile string        = "/home/xiuxiu/.config/pron/prontab"
)

func main() {
	p := pron.Create(t, configFile)
	defer p.Shutdown()

	results, errors := p.Test()
	for i := range results {
		fmt.Printf("Result: %v\n", string(results[i]))
	}
	for i := range errors {
		fmt.Printf("Error: %v\n", errors[i])
	}
}
