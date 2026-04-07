// internal/ui/ui_test.go

package ui

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// autoQuitModel is a tea.Model that immediately quits on Init.
// Used in tests to exercise the programFactory code path without blocking on TTY.
type autoQuitModel struct{}

func (a autoQuitModel) Init() tea.Cmd                       { return tea.Quit }
func (a autoQuitModel) Update(tea.Msg) (tea.Model, tea.Cmd) { return a, nil }
func (a autoQuitModel) View() string                        { return "" }

func TestGetLipglossColor(t *testing.T) {
	tests := []struct {
		colorName string
		wantErr   bool
	}{
		{"red", false},
		{"green", false},
		{"invalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.colorName, func(t *testing.T) {
			color, err := GetLipglossColor(tt.colorName)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				_, ok := ColorMap[tt.colorName]
				assert.True(t, ok, "Color %s should be valid", tt.colorName)
				expectedColor := ColorMap[tt.colorName]
				assert.Equal(t, expectedColor, color)
			}
		})
	}
}

func TestRunUIWithMock(t *testing.T) {
	tests := []struct {
		name       string
		message    string
		color      string
		mockError  error
		wantErr    bool
		wantCalled bool
	}{
		{
			name:       "Valid message and color",
			message:    "Hello, World!",
			color:      "red",
			mockError:  nil,
			wantErr:    false,
			wantCalled: true,
		},
		{
			name:       "Invalid color",
			message:    "Invalid Color Test",
			color:      "not-a-color",
			mockError:  errors.New("invalid color"),
			wantErr:    true,
			wantCalled: true,
		},
		{
			name:       "Empty message",
			message:    "",
			color:      "blue",
			mockError:  nil,
			wantErr:    false,
			wantCalled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mockRunner := &MockUIRunner{
				ReturnError: tt.mockError,
			}

			err := mockRunner.RunUI(tt.message, tt.color)

			// Check if RunUI was called
			if tt.wantCalled {
				assert.Equal(t, tt.message, mockRunner.CalledWithMessage, "RunUI() message argument mismatch")
				assert.Equal(t, tt.color, mockRunner.CalledWithColor, "RunUI() color argument mismatch")
			}

			// Validate the error returned
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestModelView(t *testing.T) {
	m := model{
		message:    "Test Message",
		colorStyle: lipgloss.NewStyle(),
	}

	expectedOutput := "Test Message\n\nPress 'q' or 'CTRL-C' to exit."

	assert.Equal(t, expectedOutput, m.View())
}

func TestModelUpdate(t *testing.T) {
	m := model{
		message:    "Test Message",
		colorStyle: lipgloss.NewStyle(),
	}

	tests := []struct {
		name    string
		msg     tea.Msg
		wantCmd bool
	}{
		{
			name:    "Key 'q' quits",
			msg:     tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}},
			wantCmd: true,
		},
		{
			name:    "CTRL+C quits",
			msg:     tea.KeyMsg{Type: tea.KeyCtrlC},
			wantCmd: true,
		},
		{
			name:    "Esc quits",
			msg:     tea.KeyMsg{Type: tea.KeyEsc},
			wantCmd: true,
		},
		{
			name:    "Unhandled key",
			msg:     tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}},
			wantCmd: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, cmd := m.Update(tt.msg)

			// Check if a command was returned
			assert.Equal(t, tt.wantCmd, cmd != nil, "Update() cmd presence mismatch")
		})
	}
}

// TestModelInit tests the Init method of the model
func TestModelInit(t *testing.T) {
	m := model{
		message:    "Test Message",
		colorStyle: lipgloss.NewStyle(),
	}

	// Init should return nil as it's a no-op in this implementation
	cmd := m.Init()

	assert.Nil(t, cmd, "Init() should return nil")
}

// TestRunUI tests the basic error path of RunUI
func TestRunUI(t *testing.T) {
	tests := []struct {
		name    string
		message string
		color   string
		wantErr bool
	}{
		{
			name:    "Invalid color error",
			message: "Test Message",
			color:   "invalid-color",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			runner := NewDefaultUIRunner()
			err := runner.RunUI(tt.message, tt.color)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestRunUISuccessPath tests the success path of RunUI using the test runner
func TestRunUISuccessPath(t *testing.T) {
	// Create a test runner with nil program factory
	runner := NewTestUIRunner()

	// Valid color to pass the GetLipglossColor check
	err := runner.RunUI("Test Message", "blue")

	// Should not return an error
	assert.NoError(t, err)
}

// TestDefaultUIRunnerCreation tests creating a DefaultUIRunner
func TestDefaultUIRunnerCreation(t *testing.T) {
	runner := NewDefaultUIRunner()
	require.NotNil(t, runner, "NewDefaultUIRunner() returned nil")
	assert.NotNil(t, runner.newProgram, "NewDefaultUIRunner() returned a runner with nil newProgram")
}

// TestNewTestUIRunner tests creating a test UI runner
func TestNewTestUIRunner(t *testing.T) {
	runner := NewTestUIRunner()
	require.NotNil(t, runner, "NewTestUIRunner() returned nil")
	assert.Nil(t, runner.newProgram, "NewTestUIRunner() returned a runner with non-nil newProgram")
}

// TestColorMap verifies that ColorMap contains the expected color names.
func TestColorMap(t *testing.T) {
	expectedColors := []string{"black", "red", "green", "yellow", "blue", "magenta", "cyan", "white"}
	assert.Len(t, ColorMap, len(expectedColors))
	for _, name := range expectedColors {
		_, ok := ColorMap[name]
		assert.True(t, ok, "ColorMap should contain %s", name)
	}
}

// TestRunUI_SuccessPathWithAutoQuitProgram tests the RunUI success path by
// injecting a programFactory that creates a headless, auto-quitting program.
// This exercises the d.newProgram != nil branch and the p.Run() success path.
func TestRunUI_SuccessPathWithAutoQuitProgram(t *testing.T) {
	var outBuf bytes.Buffer
	runner := &DefaultUIRunner{
		newProgram: func(_ tea.Model) *tea.Program {
			// autoQuitModel.Init() returns tea.Quit so the program exits immediately.
			return tea.NewProgram(
				autoQuitModel{},
				tea.WithOutput(&outBuf),
				tea.WithInput(strings.NewReader("")),
			)
		},
	}
	err := runner.RunUI("Test Message", "blue")
	assert.NoError(t, err)
}

// TestNewDefaultUIRunner_LambdaIsCallable verifies that the newProgram lambda
// in NewDefaultUIRunner is callable and returns a non-nil *tea.Program.
// This exercises the lambda body that is otherwise unreachable in tests using
// NewTestUIRunner or direct DefaultUIRunner struct literals.
func TestNewDefaultUIRunner_LambdaIsCallable(t *testing.T) {
	runner := NewDefaultUIRunner()
	require.NotNil(t, runner)
	require.NotNil(t, runner.newProgram)
	m := model{message: "test", colorStyle: lipgloss.NewStyle()}
	p := runner.newProgram(m)
	assert.NotNil(t, p, "newProgram lambda should return a non-nil *tea.Program")
}
