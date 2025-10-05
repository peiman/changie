// cmd/root_test.go

package cmd

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func TestInitConfig_InvalidConfigFile(t *testing.T) {
	cfgFile = "/invalid/path/to/config.yaml"
	defer func() { cfgFile = "" }()

	buf := new(bytes.Buffer)
	log.Logger = zerolog.New(buf)

	err := initConfig()

	if err == nil {
		t.Errorf("Expected initConfig() to return an error for invalid config file")
	}

	// Actual error message includes "failed to read config file"
	if !strings.Contains(err.Error(), "failed to read config file") {
		t.Errorf("Expected error message to contain 'failed to read config file', got '%v'", err)
	}
}

func TestInitConfig_NoConfigFile(t *testing.T) {
	viper.Reset()
	cfgFile = ""
	err := initConfig()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestExecute_ErrorPropagation(t *testing.T) {
	// Create a temporary root command for testing
	origRoot := RootCmd
	defer func() { RootCmd = origRoot }()

	testRoot := &cobra.Command{Use: "test-root"}
	testRoot.RunE = func(cmd *cobra.Command, args []string) error {
		return errors.New("some error")
	}

	// Replace the global rootCmd with testRoot
	RootCmd = testRoot

	// Execute should now produce the error "some error"
	err := Execute()
	if err == nil || !strings.Contains(err.Error(), "some error") {
		t.Errorf("Expected 'some error', got %v", err)
	}
}

func TestGetConfigValue_String(t *testing.T) {
	viper.Reset()
	defer viper.Reset()

	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().String("test-flag", "default", "test flag")

	// Test 1: Get value from viper (not flag)
	viper.Set("test.key", "viper-value")
	value := getConfigValue[string](cmd, "test-flag", "test.key")
	if value != "viper-value" {
		t.Errorf("Expected 'viper-value', got '%s'", value)
	}

	// Test 2: Flag overrides viper
	cmd.Flags().Set("test-flag", "flag-value")
	value = getConfigValue[string](cmd, "test-flag", "test.key")
	if value != "flag-value" {
		t.Errorf("Expected 'flag-value', got '%s'", value)
	}

	// Test 3: No viper value, no flag set
	viper.Reset()
	cmd2 := &cobra.Command{Use: "test2"}
	cmd2.Flags().String("test-flag", "default", "test flag")
	value = getConfigValue[string](cmd2, "test-flag", "nonexistent.key")
	if value != "" {
		t.Errorf("Expected empty string, got '%s'", value)
	}
}

func TestGetConfigValue_Bool(t *testing.T) {
	viper.Reset()
	defer viper.Reset()

	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().Bool("test-flag", false, "test flag")

	// Test 1: Get value from viper
	viper.Set("test.bool", true)
	value := getConfigValue[bool](cmd, "test-flag", "test.bool")
	if value != true {
		t.Errorf("Expected true, got %v", value)
	}

	// Test 2: Flag overrides viper
	cmd.Flags().Set("test-flag", "false")
	value = getConfigValue[bool](cmd, "test-flag", "test.bool")
	if value != false {
		t.Errorf("Expected false, got %v", value)
	}
}

func TestGetConfigValue_Int(t *testing.T) {
	viper.Reset()
	defer viper.Reset()

	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().Int("test-flag", 0, "test flag")

	// Test 1: Get value from viper
	viper.Set("test.int", 42)
	value := getConfigValue[int](cmd, "test-flag", "test.int")
	if value != 42 {
		t.Errorf("Expected 42, got %d", value)
	}

	// Test 2: Flag overrides viper
	cmd.Flags().Set("test-flag", "100")
	value = getConfigValue[int](cmd, "test-flag", "test.int")
	if value != 100 {
		t.Errorf("Expected 100, got %d", value)
	}
}

func TestGetConfigValue_Float64(t *testing.T) {
	viper.Reset()
	defer viper.Reset()

	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().Float64("test-flag", 0.0, "test flag")

	// Test 1: Get value from viper
	viper.Set("test.float", 3.14)
	value := getConfigValue[float64](cmd, "test-flag", "test.float")
	if value != 3.14 {
		t.Errorf("Expected 3.14, got %f", value)
	}

	// Test 2: Flag overrides viper
	cmd.Flags().Set("test-flag", "2.71")
	value = getConfigValue[float64](cmd, "test-flag", "test.float")
	if value != 2.71 {
		t.Errorf("Expected 2.71, got %f", value)
	}
}

func TestGetConfigValue_StringSlice(t *testing.T) {
	viper.Reset()
	defer viper.Reset()

	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().StringSlice("test-flag", []string{}, "test flag")

	// Test 1: Get value from viper
	viper.Set("test.slice", []string{"one", "two", "three"})
	value := getConfigValue[[]string](cmd, "test-flag", "test.slice")
	if len(value) != 3 || value[0] != "one" || value[1] != "two" || value[2] != "three" {
		t.Errorf("Expected [one two three], got %v", value)
	}

	// Test 2: Flag overrides viper
	cmd.Flags().Set("test-flag", "a,b,c")
	value = getConfigValue[[]string](cmd, "test-flag", "test.slice")
	if len(value) != 3 || value[0] != "a" || value[1] != "b" || value[2] != "c" {
		t.Errorf("Expected [a b c], got %v", value)
	}
}

func TestSetupCommandConfig_WithError(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}

	// Set an existing PreRunE that returns an error
	originalCalled := false
	cmd.PreRunE = func(c *cobra.Command, args []string) error {
		originalCalled = true
		return errors.New("original error")
	}

	// Setup command config
	setupCommandConfig(cmd)

	// Execute PreRunE
	err := cmd.PreRunE(cmd, []string{})

	// Should propagate the original error
	if err == nil || !strings.Contains(err.Error(), "original error") {
		t.Errorf("Expected 'original error', got %v", err)
	}

	// Original PreRunE should have been called
	if !originalCalled {
		t.Error("Expected original PreRunE to be called")
	}
}

func TestEnvPrefix(t *testing.T) {
	prefix := EnvPrefix()
	if prefix != "CHANGIE" {
		t.Errorf("Expected 'CHANGIE', got '%s'", prefix)
	}
}

func TestConfigPaths(t *testing.T) {
	paths := ConfigPaths()
	if paths.DefaultName != ".changie" {
		t.Errorf("Expected '.changie', got '%s'", paths.DefaultName)
	}
	if paths.Extension != "yaml" {
		t.Errorf("Expected 'yaml', got '%s'", paths.Extension)
	}
}
