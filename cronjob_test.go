package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestReadCronJobs tests reading valid cron jobs from a JSON file
func TestReadCronJobs(t *testing.T) {
	// Create a temporary JSON file
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test_cron.json")

	testJobs := []CronJob{
		{
			Name:    "test-job-1",
			Cron:    "*/10 * * * * *",
			Cmd:     "echo",
			Args:    []string{"hello"},
			Enabled: true,
		},
		{
			Name:    "test-job-2",
			Cron:    "0 */5 * * * *",
			Cmd:     "ls",
			Args:    []string{"-la"},
			Enabled: false,
		},
	}

	data, err := json.Marshal(testJobs)
	if err != nil {
		t.Fatalf("Failed to marshal test data: %v", err)
	}

	if err := os.WriteFile(tmpFile, data, 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Read the jobs
	jobs, err := ReadCronJobs(tmpFile)
	if err != nil {
		t.Fatalf("ReadCronJobs failed: %v", err)
	}

	// Verify the results
	if len(jobs) != 2 {
		t.Errorf("Expected 2 jobs, got %d", len(jobs))
	}

	if jobs[0].Name != "test-job-1" {
		t.Errorf("Expected job name 'test-job-1', got '%s'", jobs[0].Name)
	}

	if jobs[0].Cron != "*/10 * * * * *" {
		t.Errorf("Expected cron '*/10 * * * * *', got '%s'", jobs[0].Cron)
	}

	if !jobs[0].Enabled {
		t.Errorf("Expected job to be enabled")
	}

	if jobs[1].Enabled {
		t.Errorf("Expected job to be disabled")
	}
}

// TestReadCronJobsFileNotFound tests error handling for missing file
func TestReadCronJobsFileNotFound(t *testing.T) {
	_, err := ReadCronJobs("/nonexistent/path/to/file.json")
	if err == nil {
		t.Error("Expected error for nonexistent file, got nil")
	}
}

// TestReadCronJobsInvalidJSON tests error handling for invalid JSON
func TestReadCronJobsInvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "invalid.json")

	// Write invalid JSON
	if err := os.WriteFile(tmpFile, []byte("{ invalid json }"), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	_, err := ReadCronJobs(tmpFile)
	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
}

// TestToHumanReadable tests various cron expression conversions
func TestToHumanReadable(t *testing.T) {
	tests := []struct {
		name     string
		cron     string
		expected string
	}{
		{
			name:     "Every second",
			cron:     "* * * * * *",
			expected: "Runs every second",
		},
		{
			name:     "Every 10 seconds",
			cron:     "*/10 * * * * *",
			expected: "Runs every 10 seconds",
		},
		{
			name:     "Every 30 seconds",
			cron:     "*/30 * * * * *",
			expected: "Runs every 30 seconds",
		},
		{
			name:     "Every minute",
			cron:     "0 * * * * *",
			expected: "Runs every minute",
		},
		{
			name:     "Every 5 minutes",
			cron:     "0 */5 * * * *",
			expected: "Runs every 5 minutes",
		},
		{
			name:     "Every 15 minutes",
			cron:     "0 */15 * * * *",
			expected: "Runs every 15 minutes",
		},
		{
			name:     "Every hour",
			cron:     "0 0 * * * *",
			expected: "Runs every hour",
		},
		{
			name:     "Every 2 hours",
			cron:     "0 0 */2 * * *",
			expected: "Runs every 2 hours",
		},
		{
			name:     "Every 6 hours",
			cron:     "0 0 */6 * * *",
			expected: "Runs every 6 hours",
		},
		{
			name:     "Daily at 9:30",
			cron:     "0 30 9 * * *",
			expected: "Runs daily at 9:30",
		},
		{
			name:     "Daily at midnight",
			cron:     "0 0 0 * * *",
			expected: "Runs daily at 0:0",
		},
		{
			name:     "Daily at 14:45",
			cron:     "0 45 14 * * *",
			expected: "Runs daily at 14:45",
		},
		{
			name:     "Invalid cron (too few fields)",
			cron:     "* * *",
			expected: "Invalid cron expression: * * *",
		},
		{
			name:     "Complex expression",
			cron:     "0 15 10 * * 1",
			expected: "Runs daily at 10:15",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			job := CronJob{Cron: tt.cron}
			result := job.ToHumanReadable()
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

// TestToString tests the ToString method
func TestToString(t *testing.T) {
	tests := []struct {
		name     string
		job      CronJob
		contains []string
	}{
		{
			name: "Enabled job with args",
			job: CronJob{
				Name:    "test-task",
				Cron:    "*/10 * * * * *",
				Cmd:     "echo",
				Args:    []string{"hello", "world"},
				Enabled: true,
			},
			contains: []string{"[enabled]", "test-task", "Runs every 10 seconds", "echo hello world"},
		},
		{
			name: "Disabled job without args",
			job: CronJob{
				Name:    "disabled-task",
				Cron:    "0 */5 * * * *",
				Cmd:     "ls",
				Args:    []string{},
				Enabled: false,
			},
			contains: []string{"[disabled]", "disabled-task", "Runs every 5 minutes", "ls"},
		},
		{
			name: "Job with single arg",
			job: CronJob{
				Name:    "single-arg",
				Cron:    "0 0 * * * *",
				Cmd:     "rm",
				Args:    []string{"-rf"},
				Enabled: true,
			},
			contains: []string{"[enabled]", "single-arg", "Runs every hour", "rm -rf"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.job.ToString()
			for _, substr := range tt.contains {
				if !strings.Contains(result, substr) {
					t.Errorf("Expected ToString to contain '%s', got: %s", substr, result)
				}
			}
		})
	}
}

// TestRunSuccessfulCommand tests running a successful command
func TestRunSuccessfulCommand(t *testing.T) {
	var cmd string
	var args []string

	// Use platform-appropriate commands
	if os.PathSeparator == '\\' {
		// Windows
		cmd = "cmd"
		args = []string{"/C", "echo", "test"}
	} else {
		// Unix-like
		cmd = "echo"
		args = []string{"test"}
	}

	job := CronJob{
		Name:    "test-echo",
		Cron:    "* * * * * *",
		Cmd:     cmd,
		Args:    args,
		Enabled: true,
	}

	// Run should not panic
	job.Run()
}

// TestRunInvalidCommand tests running an invalid command
func TestRunInvalidCommand(t *testing.T) {
	job := CronJob{
		Name:    "invalid-cmd",
		Cron:    "* * * * * *",
		Cmd:     "nonexistent_command_12345",
		Args:    []string{},
		Enabled: true,
	}

	// Run should not panic even with invalid command
	job.Run()
}

// BenchmarkToHumanReadable benchmarks the ToHumanReadable function
func BenchmarkToHumanReadable(b *testing.B) {
	job := CronJob{Cron: "*/10 * * * * *"}
	for i := 0; i < b.N; i++ {
		job.ToHumanReadable()
	}
}

// BenchmarkReadCronJobs benchmarks reading cron jobs from file
func BenchmarkReadCronJobs(b *testing.B) {
	// Create a temporary JSON file
	tmpDir := b.TempDir()
	tmpFile := filepath.Join(tmpDir, "bench_cron.json")

	testJobs := []CronJob{
		{
			Name:    "bench-job",
			Cron:    "*/10 * * * * *",
			Cmd:     "echo",
			Args:    []string{"hello"},
			Enabled: true,
		},
	}

	data, _ := json.Marshal(testJobs)
	os.WriteFile(tmpFile, data, 0644)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ReadCronJobs(tmpFile)
	}
}
