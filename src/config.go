package pron

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Represents a native go method
type internalAction struct {
	fn func(...interface{}) ([]byte, error)
}

// Represents an external shell command
type externalAction struct {
	cmd *exec.Cmd
}

// Parses config file and registers externalJobs
func (p *Prontab) RegisterConfig(location string) []error {
	jobs, errs := parseConfig(location)
	for _, j := range jobs {
		j.Register(p)
	}

	return errs
}

func parseConfig(location string) (jobs []*externalJob, errs []error) {
	file, err := os.Open(location)
	defer file.Close()

	if err != nil {
		errs = append(errs, err)
		return jobs, errs
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		if job, err := parseLine(scanner.Text()); err != nil {
			errs = append(errs, err)
		} else {
			jobs = append(jobs, job)
		}
	}

	return jobs, errs
}

func parseLine(s string) (j *externalJob, err error) {
	line := strings.Fields(s)
	if len(line) < 7 {
		message := fmt.Sprintf("Malformed config line %s", s)
		return j, errors.New(message)
	}

	tLine := strings.Join(line[:6], " ")

	schedule, err := parseSchedule(tLine)

	if err != nil {
		return j, err
	}

	aLine := strings.Join(line[6:], " ")
	action, err := parseExternalAction(aLine)
	if err != nil {
		return j, err
	}

	return &externalJob{s: schedule, action: action}, nil
}

// Parses an external shell command
func parseExternalAction(i string) (a externalAction, err error) {
	switch argv := strings.Fields(i); len(argv) {
	case 0:
		return externalAction{}, errors.New("Command must not be nil")
	case 1:
		return externalAction{cmd: exec.Command(argv[0])}, nil
	default:
		return externalAction{cmd: exec.Command(argv[0], argv[1:]...)}, nil
	}
}
