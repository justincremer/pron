package main

import (
	"time"

	pron "github.com/justincremer/pron/src"
)

const (
	configFile string        = "/home/xiuxiu/.config/pron/prontab"
	t          time.Duration = time.Second
)

func main() {
	p := pron.Create(t, configFile)
	defer p.Shutdown()
	p.Startup()
}
