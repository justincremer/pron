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

// Interface for external and internal jobs
type job interface {
	Register(schedule string, tab *Prontab) error
	Dispatch() ([]byte, error)
}

// External job rollup w/ time maps
type externalJob struct {
	s      schedule
	action externalAction
}

// Internal job rollup w/ time maps
type internalJob struct {
	s      schedule
	action internalAction
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

func (p *Prontab) log(t time.Time) {
	currentTime := fmt.Sprintf("%d:%d:%d", t.Hour(), t.Minute(), t.Second())
	// fmt.Printf("%s", currentTime)

	select {
	case s := <-p.outChan:
		fmt.Printf("%s %s\n", currentTime, s)
	case e := <-p.errChan:
		fmt.Printf("%s %v\n", currentTime, e)
	default:
		fmt.Printf("%s\n", currentTime)
	}
}

func (p *Prontab) dispatchJobs(t time.Time, writer chan []byte, err chan error) {
	tick := getTick(t)
	p.log(t)

	for _, j := range p.j.externalJobs {
		if j.Scheduled(tick) {
			go j.Dispatch(writer, err)
		}
	}

	for _, j := range p.j.internalJobs {
		if j.Scheduled(tick) {
			go j.Dispatch(writer, err)
		}
	}
}

// Registers an external job to the tab
func (a *externalJob) Register(p *Prontab) {
	p.j.externalJobs = append(p.j.externalJobs, *a)
}

// Registers an internal job to the tab
func (a *internalJob) Register(p *Prontab) {
	p.j.internalJobs = append(p.j.internalJobs, *a)
}

// Internal action dispatch
func (j *internalJob) Dispatch(writer chan []byte, err chan error) {
	r, e := j.action.fn()
	writer <- r
	err <- e
}

// External action dispatch
func (j *externalJob) Dispatch(writer chan []byte, err chan error) {
	r, e := j.action.cmd.Output()
	writer <- r
	err <- e
}

func (j *externalJob) Scheduled(t tick) bool {
	return scheduled(t, j.s)
}

func (j *internalJob) Scheduled(t tick) bool {
	return scheduled(t, j.s)
}

func scheduled(t tick, s schedule) bool {
	if _, ok := s.sec[t.min]; !ok {
		return false
	}

	if _, ok := s.min[t.min]; !ok {
		return false
	}

	if _, ok := s.hour[t.hour]; !ok {
		return false
	}

	// cummulative day and dayOfWeek, as it should be
	_, day := s.day[t.day]
	_, dow := s.dow[t.dow]
	if !day && !dow {
		return false
	}

	if _, ok := s.month[t.month]; !ok {
		return false
	}

	return true
}
