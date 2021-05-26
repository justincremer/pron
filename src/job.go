package pron

import (
	"bytes"
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

func (p *Prontab) DispatchJobs(t time.Time, writer chan []byte, err chan error) {
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

// Registers an external job to the tab
func (a *externalJob) Register(p *Prontab) {
	p.j.externalJobs = append(p.j.externalJobs, *a)
}

// Registers an internal job to the tab
func (a *internalJob) Register(p *Prontab) {
	p.j.internalJobs = append(p.j.internalJobs, *a)
}

// Internal action dispatch
func (j *internalJob) Dispatch() ([]byte, error) {
	return j.action.fn()
}

// External action dispatch
func (j *externalJob) Dispatch() ([]byte, error) {
	var buf bytes.Buffer

	cmd := j.action.cmd
	cmd.Stdout = &buf

	err := cmd.Start()
	out := buf.Bytes()

	fmt.Println(buf.String())
	buf.Reset()
	cmd.Process.Kill()

	return out, err

	// return j.action.cmd.Output()
}

func (j *externalJob) scheduled(t tick) bool {
	return scheduled(t, j.s)
}

func (j *internalJob) scheduled(t tick) bool {
	return scheduled(t, j.s)
}

func ioFunctor(fn func() ([]byte, error)) func(writer chan []byte, err chan error) {
	return func(writer chan []byte, err chan error) {
		r, e := fn()
		writer <- r
		err <- e
	}
}
