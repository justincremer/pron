package pron

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

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

	schedule, err := parseSchedule(strings.Join(line[:5], " "))
	if err != nil {
		return j, err
	}

	action, err := parseExternalAction(strings.Join(line[6:], " "))
	if err != nil {
		return j, err
	}

	return &externalJob{s: schedule, action: action}, nil
}
