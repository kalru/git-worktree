package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type model struct {
	list list.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func get_worktrees() []string {
	cmd := exec.Command("sh", "-c", "git worktree list")
	out, err := cmd.Output()
	if err != nil {
		log.Printf("Error running command: %s", err)
	}
	worktrees := strings.Split(string(out[:]), "\n")
	return worktrees[:len(worktrees)-1]
}

func get_files_changed(branch string) string {
	// the last line of the command: git diff --stat
	cmd := exec.Command("sh", "-c", fmt.Sprintf("git diff %s --stat", branch))
	out, err := cmd.Output()
	if err != nil {
		log.Printf("Error running command: %s", err)
	}
	lines := strings.Split(string(out[:]), "\n")
	log.Printf("%d", len(lines))
	return lines[len(lines)-2]
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			fs := m.list.FilterState()
			if fs == list.Unfiltered || fs == list.FilterApplied {
				log.Printf("Index: %v, Selected item: %v", m.list.Index(), m.list.SelectedItem())
			}
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

func main() {
	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	}

	items := []list.Item{}

	for _, tree := range get_worktrees() {
		re_branch := regexp.MustCompile(`\[([^][]*)]`)
		branch := re_branch.FindString(tree)
		branch = branch[1 : len(branch)-1]
		items = append(items, item{title: branch, desc: get_files_changed(branch)})
	}

	m := model{list: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	m.list.Title = "Select Worktree"

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
