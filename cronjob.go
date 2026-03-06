package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type CronJob struct {
	Name    string   `yaml:"name" json:"name"`
	Cron    string   `yaml:"cron" json:"cron"`
	Cmd     string   `yaml:"cmd" json:"cmd"`
	Args    []string `yaml:"args" json:"args"`
	Enabled bool     `yaml:"enabled" json:"enabled"`
}

func ReadCronJobs(path string) ([]CronJob, error) {

	var tasks []CronJob
	file, err := os.Open(path)
	if err != nil {
		return tasks, err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&tasks)
	return tasks, err
}

func (job CronJob) Run() {
	cmd := exec.Command(job.Cmd, job.Args...)
	out := bytes.NewBuffer(nil)
	cmd.Stdout = out
	cmd.Stderr = out

	if err := cmd.Run(); err != nil {
		fmt.Printf("[%s] Error: %s\n", job.Name, err.Error())
	}

	output := out.String()
	if output != "" {
		fmt.Printf("[%s] %s\n", job.Name, strings.TrimSpace(output))
	}
}

func (job CronJob) ToHumanReadable() string {
	parts := strings.Fields(job.Cron)
	if len(parts) != 6 {
		return fmt.Sprintf("Invalid cron expression: %s", job.Cron)
	}

	second, minute, hour := parts[0], parts[1], parts[2]

	var description string

	// Check for every second
	if second == "*" && minute == "*" && hour == "*" {
		return "Runs every second"
	}

	// Check for intervals in seconds
	if strings.HasPrefix(second, "*/") && minute == "*" && hour == "*" {
		interval := strings.TrimPrefix(second, "*/")
		return fmt.Sprintf("Runs every %s seconds", interval)
	}

	// Check for every minute
	if second == "0" && minute == "*" && hour == "*" {
		return "Runs every minute"
	}

	// Check for intervals in minutes
	if second == "0" && strings.HasPrefix(minute, "*/") && hour == "*" {
		interval := strings.TrimPrefix(minute, "*/")
		return fmt.Sprintf("Runs every %s minutes", interval)
	}

	// Check for every hour
	if second == "0" && minute == "0" && hour == "*" {
		return "Runs every hour"
	}

	// Check for intervals in hours
	if second == "0" && minute == "0" && strings.HasPrefix(hour, "*/") {
		interval := strings.TrimPrefix(hour, "*/")
		return fmt.Sprintf("Runs every %s hours", interval)
	}

	// Check for specific time daily
	if second == "0" && minute != "*" && hour != "*" && !strings.Contains(minute, "/") && !strings.Contains(hour, "/") {
		return fmt.Sprintf("Runs daily at %s:%s", hour, minute)
	}

	// Default to showing the raw cron expression
	description = fmt.Sprintf("Cron: %s", job.Cron)
	return description
}

func (job CronJob) ToString() string {
	status := "enabled"
	if !job.Enabled {
		status = "disabled"
	}
	fullCmd := job.Cmd
	if len(job.Args) > 0 {
		fullCmd = fmt.Sprintf("%s %s", job.Cmd, strings.Join(job.Args, " "))
	}
	return fmt.Sprintf("[%s] %s - %s - Command: %s", status, job.Name, job.ToHumanReadable(), fullCmd)
}
