// internal/config/core_options.go
//
// Core application configuration options
//
// This file contains application-wide configuration options that apply across
// all commands and are not specific to any particular command.
// These are fundamental settings like logging level that affect the entire application.

package config

// CoreOptions returns core application configuration options
// These settings affect the overall behavior of the application
func CoreOptions() []ConfigOption {
	return []ConfigOption{
		{
			Key:          "app.log_level",
			DefaultValue: "info",
			Description:  "Logging level for the application (trace, debug, info, warn, error, fatal, panic)",
			Type:         "string",
			Required:     false,
			Example:      "debug",
		},
		{
			Key:          "app.log_format",
			DefaultValue: "auto",
			Description:  "Log output format (json, console, auto). Auto uses console when TTY is detected, JSON otherwise",
			Type:         "string",
			Required:     false,
			Example:      "json",
		},
		{
			Key:          "app.log_caller",
			DefaultValue: false,
			Description:  "Include caller information (file:line) in log output",
			Type:         "bool",
			Required:     false,
			Example:      "true",
		},
		{
			Key:          "app.json_output",
			DefaultValue: false,
			Description:  "Output command results in JSON format for machine consumption",
			Type:         "bool",
			Required:     false,
			Example:      "true",
		},
	}
}
