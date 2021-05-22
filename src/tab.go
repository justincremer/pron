package pron

import "time"

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

	func() {
		for t := range p.t.C {
			p.dispatchJobs(t, writer, err)
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
		p.dispatchJobs(t, p.outChan, p.errChan)
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
