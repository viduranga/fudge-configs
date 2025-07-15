package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type state int

const (
	menuState state = iota
	addSyncFileState
	confirmState
	errorState
)

type model struct {
	choices  []string
	cursor   int
	state    state
	input    string
	message  string
	err      error
}

type syncFile struct {
	Paths []string `json:"paths"`
}

func initialModel() model {
	return model{
		choices: []string{"Open config", "Add a Sync file", "Sync config", "Quit"},
		state:   menuState,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.state {
	case menuState:
		return m.updateMenu(msg)
	case addSyncFileState:
		return m.updateAddSyncFile(msg)
	case confirmState:
		return m.updateConfirm(msg)
	case errorState:
		return m.updateError(msg)
	}
	return m, nil
}

func (m model) updateMenu(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter", " ":
			return m.handleMenuSelection()
		}
	}
	return m, nil
}

func (m model) updateAddSyncFile(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			m.state = menuState
			m.input = ""
			m.message = ""
		case "enter":
			return m.addSyncFile()
		case "backspace":
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
			}
		default:
			m.input += msg.String()
		}
	}
	return m, nil
}

func (m model) updateConfirm(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter", "esc":
			m.state = menuState
			m.message = ""
		}
	}
	return m, nil
}

func (m model) updateError(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter", "esc":
			m.state = menuState
			m.err = nil
			m.message = ""
		}
	}
	return m, nil
}

func (m model) handleMenuSelection() (tea.Model, tea.Cmd) {
	switch m.cursor {
	case 0: // Open config
		err := m.openConfig()
		if err != nil {
			m.err = err
			m.state = errorState
		}
		return m, nil
	case 1: // Add a Sync file
		m.state = addSyncFileState
		m.input = ""
		m.message = "Enter path to JSON file:"
		return m, nil
	case 2: // Sync config
		err := m.syncConfig()
		if err != nil {
			m.err = err
			m.state = errorState
		} else {
			m.message = "Config synced successfully!"
			m.state = confirmState
		}
		return m, nil
	case 3: // Quit
		return m, tea.Quit
	}
	return m, nil
}

func (m model) openConfig() error {
	configDir, err := getConfigDir()
	if err != nil {
		return fmt.Errorf("failed to get config directory: %w", err)
	}

	configPath := filepath.Join(configDir, "config.json")
	
	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Create config.json if it doesn't exist
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		defaultConfig := map[string]interface{}{
			"servers": map[string]interface{}{},
		}
		data, err := json.MarshalIndent(defaultConfig, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal default config: %w", err)
		}
		if err := os.WriteFile(configPath, data, 0644); err != nil {
			return fmt.Errorf("failed to create config file: %w", err)
		}
	}

	// Open with $EDITOR
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi" // fallback
	}

	cmd := exec.Command(editor, configPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	return cmd.Run()
}

func (m model) addSyncFile() (tea.Model, tea.Cmd) {
	path := strings.TrimSpace(m.input)
	if path == "" {
		m.message = "Path cannot be empty"
		return m, nil
	}

	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		m.message = "File does not exist: " + path
		return m, nil
	}

	// Check if it's a JSON file
	if !strings.HasSuffix(strings.ToLower(path), ".json") {
		m.message = "File must be a JSON file"
		return m, nil
	}

	// Add to sync.json
	if err := m.addToSyncFile(path); err != nil {
		m.err = err
		m.state = errorState
		return m, nil
	}

	m.message = "File added successfully: " + path
	m.state = confirmState
	m.input = ""
	return m, nil
}

func (m model) addToSyncFile(path string) error {
	configDir, err := getConfigDir()
	if err != nil {
		return fmt.Errorf("failed to get config directory: %w", err)
	}

	syncPath := filepath.Join(configDir, "sync.json")
	
	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	var sync syncFile

	// Read existing sync file if it exists
	if data, err := os.ReadFile(syncPath); err == nil {
		if err := json.Unmarshal(data, &sync); err != nil {
			return fmt.Errorf("failed to parse sync.json: %w", err)
		}
	}

	// Check if path already exists
	for _, existingPath := range sync.Paths {
		if existingPath == path {
			return fmt.Errorf("path already exists in sync file")
		}
	}

	// Add new path
	sync.Paths = append(sync.Paths, path)

	// Write back to file
	data, err := json.MarshalIndent(sync, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal sync data: %w", err)
	}

	return os.WriteFile(syncPath, data, 0644)
}

func (m model) syncConfig() error {
	configDir, err := getConfigDir()
	if err != nil {
		return fmt.Errorf("failed to get config directory: %w", err)
	}

	configPath := filepath.Join(configDir, "config.json")
	syncPath := filepath.Join(configDir, "sync.json")

	// Read config.json
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config.json: %w", err)
	}

	var config map[string]interface{}
	if err := json.Unmarshal(configData, &config); err != nil {
		return fmt.Errorf("failed to parse config.json: %w", err)
	}

	// Read sync.json
	syncData, err := os.ReadFile(syncPath)
	if err != nil {
		return fmt.Errorf("failed to read sync.json: %w", err)
	}

	var sync syncFile
	if err := json.Unmarshal(syncData, &sync); err != nil {
		return fmt.Errorf("failed to parse sync.json: %w", err)
	}

	// Apply config to each sync file
	for _, path := range sync.Paths {
		if err := m.applySyncToFile(path, config); err != nil {
			return fmt.Errorf("failed to sync file %s: %w", path, err)
		}
	}

	return nil
}

func (m model) applySyncToFile(path string, config map[string]interface{}) error {
	// Read target file
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	var target map[string]interface{}
	if err := json.Unmarshal(data, &target); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Merge config into target (overwriting existing keys)
	for key, value := range config {
		target[key] = value
	}

	// Write back to file
	newData, err := json.MarshalIndent(target, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return os.WriteFile(path, newData, 0644)
}

func getConfigDir() (string, error) {
	xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfigHome == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		xdgConfigHome = filepath.Join(homeDir, ".config")
	}
	return filepath.Join(xdgConfigHome, "fudge_configs"), nil
}

func (m model) View() string {
	var s strings.Builder

	// Title
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 1).
		MarginBottom(1)

	s.WriteString(titleStyle.Render("Fudge Config Manager"))
	s.WriteString("\n\n")

	switch m.state {
	case menuState:
		s.WriteString(m.renderMenu())
	case addSyncFileState:
		s.WriteString(m.renderAddSyncFile())
	case confirmState:
		s.WriteString(m.renderConfirm())
	case errorState:
		s.WriteString(m.renderError())
	}

	s.WriteString("\n\nPress 'ctrl+c' to quit")
	return s.String()
}

func (m model) renderMenu() string {
	var s strings.Builder
	s.WriteString("Select an option:\n\n")

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
			choice = lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4")).Render(choice)
		}
		s.WriteString(fmt.Sprintf("%s %s\n", cursor, choice))
	}

	return s.String()
}

func (m model) renderAddSyncFile() string {
	var s strings.Builder
	s.WriteString(m.message + "\n\n")
	s.WriteString("> " + m.input + "█\n\n")
	s.WriteString("Press 'enter' to add, 'esc' to cancel")
	return s.String()
}

func (m model) renderConfirm() string {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00"))
	return style.Render(m.message) + "\n\nPress 'enter' or 'esc' to continue"
}

func (m model) renderError() string {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))
	return style.Render("Error: " + m.err.Error()) + "\n\nPress 'enter' or 'esc' to continue"
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}
} 