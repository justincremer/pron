package main

import (
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
	p.Startup()
}
