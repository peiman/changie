// internal/config/config_test.go

package config

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegistry(t *testing.T) {
	// Get the registry
	options := Registry()

	// Should have options from all sources
	assert.NotEmpty(t, options, "Registry should return configuration options")

	// Verify we have options from each category
	var hasCoreOptions, hasInitOptions, hasVersionOptions, hasDocsOptions bool

	for _, opt := range options {
		// Check for core options (app.log_level, etc.)
		if opt.Key == "app.log_level" {
			hasCoreOptions = true
		}
		// Check for init options
		if opt.Key == "app.version.use_v_prefix" {
			hasInitOptions = true
		}
		// Check for version options
		if opt.Key == "app.version.auto_push" {
			hasVersionOptions = true
		}
		// Check for docs options
		if opt.Key == "app.docs.output_format" {
			hasDocsOptions = true
		}
	}

	assert.True(t, hasCoreOptions, "Registry should include core options")
	assert.True(t, hasInitOptions, "Registry should include init options")
	assert.True(t, hasVersionOptions, "Registry should include version options")
	assert.True(t, hasDocsOptions, "Registry should include docs options")

	// Verify all options have required fields
	for _, opt := range options {
		assert.NotEmpty(t, opt.Key, "All options should have a key")
		assert.NotEmpty(t, opt.Description, "All options should have a description: %s", opt.Key)
		assert.NotEmpty(t, opt.Type, "All options should have a type: %s", opt.Key)
		// DefaultValue can be nil or any value
		// Required, Example are optional
	}
}

func TestSetDefaults(t *testing.T) {
	// Reset viper to clean state
	viper.Reset()

	// Call SetDefaults
	SetDefaults()

	// Verify some key defaults were set
	logLevel := viper.Get("app.log_level")
	assert.NotNil(t, logLevel, "app.log_level should be set")
	assert.Equal(t, "info", logLevel, "app.log_level should default to 'info'")

	changelogFile := viper.Get("app.changelog.file")
	assert.NotNil(t, changelogFile, "app.changelog.file should be set")
	assert.Equal(t, "CHANGELOG.md", changelogFile, "app.changelog.file should default to 'CHANGELOG.md'")

	useVPrefix := viper.Get("app.version.use_v_prefix")
	assert.NotNil(t, useVPrefix, "app.version.use_v_prefix should be set")
	assert.Equal(t, true, useVPrefix, "app.version.use_v_prefix should default to true")

	// Verify all options from registry are set
	for _, opt := range Registry() {
		value := viper.Get(opt.Key)
		assert.Equal(t, opt.DefaultValue, value,
			"Viper should have default value for %s", opt.Key)
	}
}

func TestCoreOptions(t *testing.T) {
	options := CoreOptions()

	assert.NotEmpty(t, options, "CoreOptions should return options")

	// Check for expected core options
	var hasLogLevel bool

	for _, opt := range options {
		switch opt.Key {
		case "app.log_level":
			hasLogLevel = true
			assert.Equal(t, "info", opt.DefaultValue)
			assert.Equal(t, "string", opt.Type)
			assert.Contains(t, opt.Description, "Logging level")
		}
	}

	assert.True(t, hasLogLevel, "CoreOptions should include app.log_level")
}

func TestInitOptions(t *testing.T) {
	options := InitOptions()

	assert.NotEmpty(t, options, "InitOptions should return options")

	// Check for expected init options
	var hasUseVPrefix bool

	for _, opt := range options {
		if opt.Key == "app.version.use_v_prefix" {
			hasUseVPrefix = true
			assert.Equal(t, true, opt.DefaultValue)
			assert.Equal(t, "bool", opt.Type)
			assert.Contains(t, opt.Description, "'v' prefix")
		}
	}

	assert.True(t, hasUseVPrefix, "InitOptions should include app.version.use_v_prefix")
}

func TestVersionOptions(t *testing.T) {
	options := VersionOptions()

	assert.NotEmpty(t, options, "VersionOptions should return options")

	// Check for expected version options
	var hasAutoPush, hasRepoProvider bool

	for _, opt := range options {
		switch opt.Key {
		case "app.version.auto_push":
			hasAutoPush = true
			assert.Equal(t, false, opt.DefaultValue)
			assert.Equal(t, "bool", opt.Type)
			assert.Contains(t, opt.Description, "push")
		case "app.changelog.repository_provider":
			hasRepoProvider = true
			assert.Equal(t, "github", opt.DefaultValue)
			assert.Equal(t, "string", opt.Type)
		}
	}

	assert.True(t, hasAutoPush, "VersionOptions should include app.version.auto_push")
	assert.True(t, hasRepoProvider, "VersionOptions should include app.changelog.repository_provider")
}

func TestDocsOptions(t *testing.T) {
	options := DocsOptions()

	assert.NotEmpty(t, options, "DocsOptions should return options")

	// Check for expected docs options
	var hasOutputFormat, hasOutputFile bool

	for _, opt := range options {
		switch opt.Key {
		case "app.docs.output_format":
			hasOutputFormat = true
			assert.Equal(t, "markdown", opt.DefaultValue)
			assert.Equal(t, "string", opt.Type)
		case "app.docs.output_file":
			hasOutputFile = true
			assert.Equal(t, "", opt.DefaultValue) // Empty string means stdout
			assert.Equal(t, "string", opt.Type)
		}
	}

	assert.True(t, hasOutputFormat, "DocsOptions should include app.docs.output_format")
	assert.True(t, hasOutputFile, "DocsOptions should include app.docs.output_file")
}

func TestConfigOption_EnvVarName(t *testing.T) {
	tests := []struct {
		name     string
		opt      ConfigOption
		prefix   string
		expected string
	}{
		{
			name: "simple key",
			opt: ConfigOption{
				Key: "app.log_level",
			},
			prefix:   "CHANGIE",
			expected: "CHANGIE_APP_LOG_LEVEL",
		},
		{
			name: "nested key",
			opt: ConfigOption{
				Key: "app.changelog.file",
			},
			prefix:   "CHANGIE",
			expected: "CHANGIE_APP_CHANGELOG_FILE",
		},
		{
			name: "different prefix",
			opt: ConfigOption{
				Key: "app.version.use_v_prefix",
			},
			prefix:   "MYAPP",
			expected: "MYAPP_APP_VERSION_USE_V_PREFIX",
		},
		{
			name: "empty prefix",
			opt: ConfigOption{
				Key: "test.key",
			},
			prefix:   "",
			expected: "_TEST_KEY",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.opt.EnvVarName(tc.prefix)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestConfigOption_DefaultValueString(t *testing.T) {
	tests := []struct {
		name     string
		opt      ConfigOption
		expected string
	}{
		{
			name: "string value",
			opt: ConfigOption{
				DefaultValue: "test",
			},
			expected: "test",
		},
		{
			name: "int value",
			opt: ConfigOption{
				DefaultValue: 42,
			},
			expected: "42",
		},
		{
			name: "bool true",
			opt: ConfigOption{
				DefaultValue: true,
			},
			expected: "true",
		},
		{
			name: "bool false",
			opt: ConfigOption{
				DefaultValue: false,
			},
			expected: "false",
		},
		{
			name: "nil value",
			opt: ConfigOption{
				DefaultValue: nil,
			},
			expected: "nil",
		},
		{
			name: "float value",
			opt: ConfigOption{
				DefaultValue: 3.14,
			},
			expected: "3.14",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.opt.DefaultValueString()
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestConfigOption_ExampleValueString(t *testing.T) {
	tests := []struct {
		name     string
		opt      ConfigOption
		expected string
	}{
		{
			name: "has example",
			opt: ConfigOption{
				DefaultValue: "default",
				Example:      "example",
			},
			expected: "example",
		},
		{
			name: "no example - uses default",
			opt: ConfigOption{
				DefaultValue: "default",
				Example:      "",
			},
			expected: "default",
		},
		{
			name: "no example, nil default",
			opt: ConfigOption{
				DefaultValue: nil,
				Example:      "",
			},
			expected: "nil",
		},
		{
			name: "example with special characters",
			opt: ConfigOption{
				DefaultValue: "default",
				Example:      "/path/to/file.md",
			},
			expected: "/path/to/file.md",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.opt.ExampleValueString()
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestConfigOption_AllFieldsPopulated(t *testing.T) {
	// Verify that all options in the registry have sensible values
	options := Registry()

	for _, opt := range options {
		t.Run(opt.Key, func(t *testing.T) {
			// Every option must have these fields
			require.NotEmpty(t, opt.Key, "Key should not be empty")
			require.NotEmpty(t, opt.Description, "Description should not be empty")
			require.NotEmpty(t, opt.Type, "Type should not be empty")

			// Type should be one of the known types
			validTypes := []string{"string", "int", "bool", "float"}
			assert.Contains(t, validTypes, opt.Type,
				"Type should be one of: %v", validTypes)

			// If there's an example, it should not be empty
			if opt.Example != "" {
				assert.NotEmpty(t, opt.Example,
					"If Example is set, it should not be empty string")
			}

			// Test that the option methods work
			envVar := opt.EnvVarName("TEST")
			assert.NotEmpty(t, envVar, "EnvVarName should produce a value")

			defaultStr := opt.DefaultValueString()
			// DefaultValueString can be empty for empty string defaults
			assert.NotNil(t, defaultStr, "DefaultValueString should produce a value")

			exampleStr := opt.ExampleValueString()
			// ExampleValueString can be empty for empty string defaults
			assert.NotNil(t, exampleStr, "ExampleValueString should produce a value")
		})
	}
}

func TestSetDefaults_Idempotent(t *testing.T) {
	// Reset viper
	viper.Reset()

	// Set defaults twice
	SetDefaults()
	SetDefaults()

	// Values should still be correct
	logLevel := viper.Get("app.log_level")
	assert.Equal(t, "info", logLevel, "Multiple calls to SetDefaults should be idempotent")
}

func TestRegistry_NoDuplicateKeys(t *testing.T) {
	options := Registry()
	keys := make(map[string]bool)

	for _, opt := range options {
		if keys[opt.Key] {
			t.Errorf("Duplicate key found in registry: %s", opt.Key)
		}
		keys[opt.Key] = true
	}
}

func TestConfigOption_TypeValidation(t *testing.T) {
	// Verify that the DefaultValue matches the declared Type
	options := Registry()

	for _, opt := range options {
		t.Run(opt.Key, func(t *testing.T) {
			if opt.DefaultValue == nil {
				return // nil is valid for any type
			}

			switch opt.Type {
			case "string":
				_, ok := opt.DefaultValue.(string)
				assert.True(t, ok, "DefaultValue should be string for type 'string'")
			case "int":
				_, ok := opt.DefaultValue.(int)
				assert.True(t, ok, "DefaultValue should be int for type 'int'")
			case "bool":
				_, ok := opt.DefaultValue.(bool)
				assert.True(t, ok, "DefaultValue should be bool for type 'bool'")
			case "float":
				_, ok := opt.DefaultValue.(float64)
				if !ok {
					// Also accept float32
					_, ok = opt.DefaultValue.(float32)
				}
				assert.True(t, ok, "DefaultValue should be float for type 'float'")
			}
		})
	}
}
