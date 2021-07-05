package pron

import (
	"fmt"
	"time"
)

// Top level pron struct
type Tab struct {
	ticker  *time.Ticker
	jobs    jobs
	outChan chan []byte
	errChan chan error
}

type jobs struct {
	externalJobs []externalJob
	internalJobs []internalJob
}

// Initializes the tab and registers jobs
func Create(t time.Duration, file string) *Tab {
	writer := make(chan []byte)
	err := make(chan error)

	p := &Tab{t: time.NewTicker(t), outChan: writer, errChan: err}
	p.initialize(file)

	go func() {
		for t := range p.t.C {
			p.DispatchJobs(t, writer, err)
		}
	}()
	return p
}

// Reads config file and propogates the Tab jobs slice
func (p *Tab) initialize(file string) {
	if errs := p.registerConfig(file); len(errs) != 0 {
		for _, e := range errs {
			panic(e)
		}
	}
}

// Starts auto dispatching commands
func (p *Tab) Startup() {
	for t := range p.t.C {
		p.DispatchJobs(t, p.outChan, p.errChan)
	}
}

// Emptys the job buffer, stops the clock, and closes channels
func (p *Tab) Shutdown() {
	defer close(p.outChan)
	defer close(p.errChan)

	p.j.externalJobs = []externalJob{}
	p.j.internalJobs = []internalJob{}
	p.t.Stop()
}

func (p *Tab) DispatchJobs(t time.Time, writer chan []byte, err chan error) {
	tick := getTick(t)
	for _, j := range p.j.externalJobs {
		if j.scheduled(tick) {
			go ioFunctor(j.Dispatch)(writer, err)
		}
	}

	for _, j := range p.j.internalJobs {
		if j.scheduled(tick) {
			go ioFunctor(j.Dispatch)(writer, err)
		}
	}
	p.log(t)
}

func (p *Tab) RegisterJobs(path string, internalJobs *[]internalJob) *[]error {
	errs := p.registerConfig(path)
	// register internal jobs
	return &errs
}

func (p *Tab) log(t time.Time) {
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
