package pron

import (
	"time"
)

// Top level pron struct
type Prontab struct {
	t *time.Ticker
	j jobs
	e []error
}

type jobs struct {
	externalJobs []*externalJob
	internalJobs []*internalJob
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
	p := &Prontab{t: time.NewTicker(t)}
	if errs := p.RegisterConfig(file); len(errs) != 0 {
		for _, e := range errs {
			panic(e)
		}
	}

	return p
}

func (p *Prontab) Test() (results [][]byte, errs []error) {
	for _, j := range p.j.externalJobs {
		res, err := j.Dispatch()
		results = append(results, res)
		errs = append(errs, err)
	}

	return results, errs
}

// Emptys the job buffer and stops the clock
func (p *Prontab) Shutdown() {
	p.j.externalJobs = nil
	p.j.internalJobs = nil
	p.t.Stop()
}

// Registers an external job to the tab
func (a *externalJob) Register(tab *Prontab) {
	tab.j.externalJobs = append(tab.j.externalJobs, a)
}

// Registers an internal job to the tab
func (a *internalJob) Register(tab *Prontab) {
	tab.j.internalJobs = append(tab.j.internalJobs, a)
}

// Internal action dispatch
func (j *internalJob) Dispatch() ([]byte, error) {
	return j.action.fn()
}

// External action dispatch
func (j *externalJob) Dispatch() ([]byte, error) {
	return j.action.cmd.Output()
}
