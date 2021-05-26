package pron

import (
	"fmt"
	"time"
)

// Top level pron struct
type Prontab struct {
	t       *time.Ticker
	j       jobs
	outChan chan []byte
	errChan chan error
}

type jobs struct {
	externalJobs []externalJob
	internalJobs []internalJob
}

// Initializes the tab and registers jobs
func Create(t time.Duration, file string) *Prontab {
	writer := make(chan []byte)
	err := make(chan error)

	p := &Prontab{t: time.NewTicker(t), outChan: writer, errChan: err}
	p.initialize(file)

	go func() {
		for t := range p.t.C {
			p.DispatchJobs(t, writer, err)
		}
	}()
	return p
}

// Reads config file and propogates the Prontab jobs slice
func (p *Prontab) initialize(file string) {
	if errs := p.RegisterConfig(file); len(errs) != 0 {
		for _, e := range errs {
			panic(e)
		}
	}
}

// Starts auto dispatching commands
func (p *Prontab) Startup() {
	for t := range p.t.C {
		p.DispatchJobs(t, p.outChan, p.errChan)
	}
}

// Emptys the job buffer, stops the clock, and closes channels
func (p *Prontab) Shutdown() {
	defer close(p.outChan)
	defer close(p.errChan)

	p.j.externalJobs = []externalJob{}
	p.j.internalJobs = []internalJob{}
	p.t.Stop()
}

func (p *Prontab) log(t time.Time) {
	currentTime := fmt.Sprintf("%d:%d:%d", t.Hour(), t.Minute(), t.Second())

	select {
	case s := <-p.outChan:
		fmt.Printf("%s %s\n", currentTime, s)
	case e := <-p.errChan:
		fmt.Printf("%s %v\n", currentTime, e)
	default:
		fmt.Printf("%s\n", currentTime)
	}
}
