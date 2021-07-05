package pron

import (
	"bytes"
	"fmt"
	"os/exec"
)

// Interface for external and internal jobs
type job interface {
	Register(schedule string, tab *Tab) error
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

// Table of time maps for a given job
type schedule struct {
	sec   map[int]struct{}
	min   map[int]struct{}
	hour  map[int]struct{}
	day   map[int]struct{}
	month map[int]struct{}
	dow   map[int]struct{}
}

// Native go method
type internalAction struct {
	fn func(...interface{}) ([]byte, error)
}

// External arbitrary program
type externalAction struct {
	cmd *exec.Cmd
}

// Registers an external job to the tab
func (a *externalJob) register(p *Tab) {
	p.j.externalJobs = append(p.j.externalJobs, *a)
}

// Registers an internal job to the tab
func (a *internalJob) register(p *Tab) {
	p.j.internalJobs = append(p.j.internalJobs, *a)
}

// Internal action dispatch
func (j *internalJob) Dispatch() ([]byte, error) {
	return j.action.fn()
}

// TODO: fix + cleanup console routines
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
