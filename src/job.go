package pron

import (
	"fmt"
	"time"
)

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

// Registers an external job to the tab
func (a *externalJob) Register(p *Prontab) {
	p.j.externalJobs = append(p.j.externalJobs, *a)
}

// Registers an internal job to the tab
func (a *internalJob) Register(p *Prontab) {
	p.j.internalJobs = append(p.j.internalJobs, *a)
}

func (p *Prontab) dispatchJobs(t time.Time, writer chan []byte, err chan error) {
	tick := getTick(t)
	p.log(t)

	for _, j := range p.j.externalJobs {
		if j.Scheduled(tick) {
			go ioFunctor(j.Dispatch)(writer, err)
		}
	}

	for _, j := range p.j.internalJobs {
		if j.Scheduled(tick) {
			go ioFunctor(j.Dispatch)(writer, err)
		}
	}
}

func ioFunctor(fn func() ([]byte, error)) func(writer chan []byte, err chan error) {
	return func(writer chan []byte, err chan error) {
		r, e := fn()
		writer <- r
		err <- e
	}
}

// Internal action dispatch
func (j *internalJob) Dispatch() ([]byte, error) {
	return j.action.fn()
}

// External action dispatch
func (j *externalJob) Dispatch() ([]byte, error) {
	return j.action.cmd.Output()
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
