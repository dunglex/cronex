package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

// TestMainConfiguration tests the overall configuration flow
func TestMainConfiguration(t *testing.T) {
	// Create a temporary directory and config file
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "test_config.json")

	// Create test configuration
	testJobs := []CronJob{
		{
			Name:    "integration-test-1",
			Cron:    "*/10 * * * * *",
			Cmd:     "echo",
			Args:    []string{"test1"},
			Enabled: true,
		},
		{
			Name:    "integration-test-2",
			Cron:    "0 */5 * * * *",
			Cmd:     "echo",
			Args:    []string{"test2"},
			Enabled: false,
		},
		{
			Name:    "integration-test-3",
			Cron:    "0 0 * * * *",
			Cmd:     "echo",
			Args:    []string{"test3"},
			Enabled: true,
		},
	}

	data, err := json.Marshal(testJobs)
	if err != nil {
		t.Fatalf("Failed to marshal test data: %v", err)
	}

	if err := os.WriteFile(configFile, data, 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Read the configuration
	jobs, err := ReadCronJobs(configFile)
	if err != nil {
		t.Fatalf("Failed to read cron jobs: %v", err)
	}

	// Verify the jobs were loaded correctly
	if len(jobs) != 3 {
		t.Errorf("Expected 3 jobs, got %d", len(jobs))
	}

	// Count enabled jobs (like main.go does)
	enabledCount := 0
	for _, job := range jobs {
		if job.Enabled {
			enabledCount++
		}
	}

	if enabledCount != 2 {
		t.Errorf("Expected 2 enabled jobs, got %d", enabledCount)
	}
}

// TestConfigWithEmptyArgs tests configuration with empty args array
func TestConfigWithEmptyArgs(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "empty_args.json")

	testJobs := []CronJob{
		{
			Name:    "no-args-job",
			Cron:    "* * * * * *",
			Cmd:     "echo",
			Args:    []string{},
			Enabled: true,
		},
	}

	data, _ := json.Marshal(testJobs)
	os.WriteFile(configFile, data, 0644)

	jobs, err := ReadCronJobs(configFile)
	if err != nil {
		t.Fatalf("Failed to read cron jobs: %v", err)
	}

	if len(jobs[0].Args) != 0 {
		t.Errorf("Expected empty args, got %v", jobs[0].Args)
	}
}

// TestConfigWithoutArgs tests configuration with nil args (omitted in JSON)
func TestConfigWithoutArgs(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "no_args.json")

	jsonData := `[
		{
			"name": "no-args-field",
			"cron": "* * * * * *",
			"cmd": "echo",
			"enabled": true
		}
	]`

	os.WriteFile(configFile, []byte(jsonData), 0644)

	jobs, err := ReadCronJobs(configFile)
	if err != nil {
		t.Fatalf("Failed to read cron jobs: %v", err)
	}

	if jobs[0].Args == nil {
		// Nil is acceptable
		return
	}

	if len(jobs[0].Args) != 0 {
		t.Errorf("Expected nil or empty args, got %v", jobs[0].Args)
	}
}

// TestConfigWithMultipleJobs tests handling multiple cron jobs
func TestConfigWithMultipleJobs(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "multiple_jobs.json")

	testJobs := make([]CronJob, 10)
	for i := 0; i < 10; i++ {
		testJobs[i] = CronJob{
			Name:    "job-" + string(rune('0'+i)),
			Cron:    "*/10 * * * * *",
			Cmd:     "echo",
			Args:    []string{"job", string(rune('0' + i))},
			Enabled: i%2 == 0, // Enable every other job
		}
	}

	data, _ := json.Marshal(testJobs)
	os.WriteFile(configFile, data, 0644)

	jobs, err := ReadCronJobs(configFile)
	if err != nil {
		t.Fatalf("Failed to read cron jobs: %v", err)
	}

	if len(jobs) != 10 {
		t.Errorf("Expected 10 jobs, got %d", len(jobs))
	}

	enabledCount := 0
	for _, job := range jobs {
		if job.Enabled {
			enabledCount++
		}
	}

	if enabledCount != 5 {
		t.Errorf("Expected 5 enabled jobs, got %d", enabledCount)
	}
}

// TestEmptyConfigFile tests reading an empty job array
func TestEmptyConfigFile(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "empty.json")

	os.WriteFile(configFile, []byte("[]"), 0644)

	jobs, err := ReadCronJobs(configFile)
	if err != nil {
		t.Fatalf("Failed to read empty config: %v", err)
	}

	if len(jobs) != 0 {
		t.Errorf("Expected 0 jobs, got %d", len(jobs))
	}
}

// TestYAMLConfiguration tests reading YAML configuration file
func TestYAMLConfiguration(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "test_config.yaml")

	testJobs := []CronJob{
		{
			Name:    "yaml-test-1",
			Cron:    "*/10 * * * * *",
			Cmd:     "echo",
			Args:    []string{"yaml test"},
			Enabled: true,
		},
		{
			Name:    "yaml-test-2",
			Cron:    "0 */5 * * * *",
			Cmd:     "echo",
			Args:    []string{"yaml test 2"},
			Enabled: false,
		},
	}

	data, err := yaml.Marshal(testJobs)
	if err != nil {
		t.Fatalf("Failed to marshal YAML data: %v", err)
	}

	if err := os.WriteFile(configFile, data, 0644); err != nil {
		t.Fatalf("Failed to write YAML config file: %v", err)
	}

	jobs, err := ReadCronJobs(configFile)
	if err != nil {
		t.Fatalf("Failed to read YAML cron jobs: %v", err)
	}

	if len(jobs) != 2 {
		t.Errorf("Expected 2 jobs, got %d", len(jobs))
	}

	if jobs[0].Name != "yaml-test-1" {
		t.Errorf("Expected job name 'yaml-test-1', got '%s'", jobs[0].Name)
	}

	if jobs[0].Enabled != true {
		t.Errorf("Expected first job to be enabled")
	}

	if jobs[1].Enabled != false {
		t.Errorf("Expected second job to be disabled")
	}
}

// TestYAMLWithYMLExtension tests reading .yml extension
func TestYAMLWithYMLExtension(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "test_config.yml")

	testJobs := []CronJob{
		{
			Name:    "yml-test",
			Cron:    "* * * * * *",
			Cmd:     "echo",
			Args:    []string{"yml extension"},
			Enabled: true,
		},
	}

	data, _ := yaml.Marshal(testJobs)
	os.WriteFile(configFile, data, 0644)

	jobs, err := ReadCronJobs(configFile)
	if err != nil {
		t.Fatalf("Failed to read .yml file: %v", err)
	}

	if len(jobs) != 1 {
		t.Errorf("Expected 1 job, got %d", len(jobs))
	}

	if jobs[0].Name != "yml-test" {
		t.Errorf("Expected job name 'yml-test', got '%s'", jobs[0].Name)
	}
}

// TestYAMLWithEmptyArgs tests YAML config with empty args
func TestYAMLWithEmptyArgs(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "yaml_empty_args.yaml")

	yamlData := `- name: no-args-yaml
  cron: "* * * * * *"
  cmd: echo
  args: []
  enabled: true
`

	os.WriteFile(configFile, []byte(yamlData), 0644)

	jobs, err := ReadCronJobs(configFile)
	if err != nil {
		t.Fatalf("Failed to read YAML cron jobs: %v", err)
	}

	if len(jobs[0].Args) != 0 {
		t.Errorf("Expected empty args, got %v", jobs[0].Args)
	}
}

// TestYAMLWithoutArgs tests YAML config without args field
func TestYAMLWithoutArgs(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "yaml_no_args.yaml")

	yamlData := `- name: no-args-field-yaml
  cron: "* * * * * *"
  cmd: echo
  enabled: true
`

	os.WriteFile(configFile, []byte(yamlData), 0644)

	jobs, err := ReadCronJobs(configFile)
	if err != nil {
		t.Fatalf("Failed to read YAML cron jobs: %v", err)
	}

	if jobs[0].Args == nil {
		// Nil is acceptable
		return
	}

	if len(jobs[0].Args) != 0 {
		t.Errorf("Expected nil or empty args, got %v", jobs[0].Args)
	}
}

// TestYAMLMultipleJobs tests YAML with multiple jobs
func TestYAMLMultipleJobs(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "yaml_multiple.yaml")

	testJobs := make([]CronJob, 10)
	for i := 0; i < 10; i++ {
		testJobs[i] = CronJob{
			Name:    "yaml-job-" + string(rune('0'+i)),
			Cron:    "*/10 * * * * *",
			Cmd:     "echo",
			Args:    []string{"yaml", "job", string(rune('0' + i))},
			Enabled: i%3 == 0, // Enable every third job
		}
	}

	data, _ := yaml.Marshal(testJobs)
	os.WriteFile(configFile, data, 0644)

	jobs, err := ReadCronJobs(configFile)
	if err != nil {
		t.Fatalf("Failed to read YAML cron jobs: %v", err)
	}

	if len(jobs) != 10 {
		t.Errorf("Expected 10 jobs, got %d", len(jobs))
	}

	enabledCount := 0
	for _, job := range jobs {
		if job.Enabled {
			enabledCount++
		}
	}

	if enabledCount != 4 {
		t.Errorf("Expected 4 enabled jobs, got %d", enabledCount)
	}
}

// TestEmptyYAMLConfigFile tests reading an empty YAML job array
func TestEmptyYAMLConfigFile(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "empty.yaml")

	os.WriteFile(configFile, []byte("[]"), 0644)

	jobs, err := ReadCronJobs(configFile)
	if err != nil {
		t.Fatalf("Failed to read empty YAML config: %v", err)
	}

	if len(jobs) != 0 {
		t.Errorf("Expected 0 jobs, got %d", len(jobs))
	}
}

// TestJSONAndYAMLCompatibility tests that both formats produce same results
func TestJSONAndYAMLCompatibility(t *testing.T) {
	tmpDir := t.TempDir()
	jsonFile := filepath.Join(tmpDir, "test.json")
	yamlFile := filepath.Join(tmpDir, "test.yaml")

	testJobs := []CronJob{
		{
			Name:    "compatibility-test",
			Cron:    "*/10 * * * * *",
			Cmd:     "echo",
			Args:    []string{"test", "arg"},
			Enabled: true,
		},
	}

	// Write JSON
	jsonData, _ := json.Marshal(testJobs)
	os.WriteFile(jsonFile, jsonData, 0644)

	// Write YAML
	yamlData, _ := yaml.Marshal(testJobs)
	os.WriteFile(yamlFile, yamlData, 0644)

	// Read both
	jsonJobs, err := ReadCronJobs(jsonFile)
	if err != nil {
		t.Fatalf("Failed to read JSON: %v", err)
	}

	yamlJobs, err := ReadCronJobs(yamlFile)
	if err != nil {
		t.Fatalf("Failed to read YAML: %v", err)
	}

	// Compare
	if len(jsonJobs) != len(yamlJobs) {
		t.Errorf("Job count mismatch: JSON=%d, YAML=%d", len(jsonJobs), len(yamlJobs))
	}

	if jsonJobs[0].Name != yamlJobs[0].Name {
		t.Errorf("Name mismatch: JSON=%s, YAML=%s", jsonJobs[0].Name, yamlJobs[0].Name)
	}

	if jsonJobs[0].Cron != yamlJobs[0].Cron {
		t.Errorf("Cron mismatch: JSON=%s, YAML=%s", jsonJobs[0].Cron, yamlJobs[0].Cron)
	}

	if jsonJobs[0].Enabled != yamlJobs[0].Enabled {
		t.Errorf("Enabled mismatch: JSON=%v, YAML=%v", jsonJobs[0].Enabled, yamlJobs[0].Enabled)
	}
}

// TestFileExtensionCaseInsensitive tests that file extension detection is case-insensitive
func TestFileExtensionCaseInsensitive(t *testing.T) {
	tmpDir := t.TempDir()

	testCases := []string{
		"test.YAML",
		"test.YML",
		"test.Yaml",
		"test.Yml",
	}

	testJobs := []CronJob{
		{
			Name:    "case-test",
			Cron:    "* * * * * *",
			Cmd:     "echo",
			Args:    []string{"test"},
			Enabled: true,
		},
	}

	data, _ := yaml.Marshal(testJobs)

	for _, fileName := range testCases {
		configFile := filepath.Join(tmpDir, fileName)
		os.WriteFile(configFile, data, 0644)

		jobs, err := ReadCronJobs(configFile)
		if err != nil {
			t.Errorf("Failed to read %s: %v", fileName, err)
			continue
		}

		if len(jobs) != 1 {
			t.Errorf("Expected 1 job from %s, got %d", fileName, len(jobs))
		}
	}
}

// TestUnknownExtensionDefaultsToJSON tests that unknown extensions default to JSON
func TestUnknownExtensionDefaultsToJSON(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "test.conf")

	testJobs := []CronJob{
		{
			Name:    "unknown-ext-test",
			Cron:    "* * * * * *",
			Cmd:     "echo",
			Args:    []string{"test"},
			Enabled: true,
		},
	}

	// Write as JSON even though extension is .conf
	data, _ := json.Marshal(testJobs)
	os.WriteFile(configFile, data, 0644)

	jobs, err := ReadCronJobs(configFile)
	if err != nil {
		t.Fatalf("Failed to read config with unknown extension: %v", err)
	}

	if len(jobs) != 1 {
		t.Errorf("Expected 1 job, got %d", len(jobs))
	}

	if jobs[0].Name != "unknown-ext-test" {
		t.Errorf("Expected job name 'unknown-ext-test', got '%s'", jobs[0].Name)
	}
}
